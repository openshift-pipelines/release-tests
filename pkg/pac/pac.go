package pac

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/clients"
	"github.com/openshift-pipelines/release-tests/pkg/cmd"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	"github.com/openshift-pipelines/release-tests/pkg/k8s"
	"github.com/openshift-pipelines/release-tests/pkg/oc"
	"github.com/openshift-pipelines/release-tests/pkg/store"
	"github.com/openshift-pipelines/release-tests/pkg/triggers"
	"github.com/openshift-pipelines/release-tests/pkg/wait"
	eventReconciler "github.com/tektoncd/triggers/pkg/reconciler/eventlistener"
	"github.com/xanzy/go-gitlab"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ConfigureGitlabToken() {
	secretData := os.Getenv("GITLAB_TOKEN")
	if secretData == "" {
		log.Printf("Token for authorization to the Gitlab repository was not exported as a system variable")
	} else {
		if !oc.SecretExists("gitlab-auth-secret", "openshift-pipelines") {
			oc.CreateSecretForGitLab(secretData)
		} else {
			log.Printf("Secret \"gitlab-auth-secret\" already exists")
		}
		store.PutScenarioData("gitlabToken", secretData)
	}
}

func getNewSmeeURL() (string, error) {
	curlCommand := `curl -Ls -o /dev/null -w %{url_effective} https://smee.io/new`

	cmd := exec.Command("sh", "-c", curlCommand)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to create SmeeURL: %v", err)
	}

	smeeURL := strings.TrimSpace(string(output))

	if smeeURL == "" {
		return "", fmt.Errorf("failed to retrieve Smee URL: no URL found")
	}

	return smeeURL, nil
}

func createSmeeDeployment(c *clients.Clients, namespace, smeeURL, targetURL string) error {
	deployment := &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "gosmee-client",
		},
		Spec: v1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "gosmee-client",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "gosmee-client",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "gosmee-client",
							Image: "ghcr.io/chmouel/gosmee:latest",
							Command: []string{
								"gosmee",
								"client",
								smeeURL,
								targetURL,
							},
							Env: []corev1.EnvVar{
								{
									Name:  "SMEE_URL",
									Value: smeeURL,
								},
								{
									Name:  "TARGET_URL",
									Value: targetURL,
								},
							},
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 8080,
								},
							},
						},
					},
				},
			},
		},
	}

	kc := c.KubeClient.Kube
	deploymentsClient := kc.AppsV1().Deployments(namespace)
	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create deployment: %v", err)
	}

	log.Printf("Created deployment %q in namespace %q.\n", result.GetObjectMeta().GetName(), namespace)
	return nil
}

func int32Ptr(i int32) *int32 { return &i }

func forkProject(client *gitlab.Client, projectID, targetGroupNamespace string) (*gitlab.Project, error) {
	var project *gitlab.Project
	var response *gitlab.Response
	var err error
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		randomSuffix := strconv.FormatInt(time.Now().UnixNano(), 10)[:8]
		projectName := "openshift-pipelines-test-fork-" + randomSuffix

		forkOptions := &gitlab.ForkProjectOptions{
			Namespace: &targetGroupNamespace,
			Name:      &projectName,
			Path:      &projectName,
		}

		project, response, err = client.Projects.ForkProject(projectID, forkOptions)
		if err != nil {
			log.Printf("Attempt %d: Failed to fork project: %v\n", i+1, err)
			if response != nil {
				if body, readErr := io.ReadAll(response.Body); readErr == nil {
					log.Printf("Response body: %s\n", string(body))
				}
				if response.StatusCode == http.StatusConflict {
					log.Printf("Project name or path already exists, trying again with a new name.")
					continue
				}
			}
			return nil, err
		}

		return project, nil
	}

	return nil, fmt.Errorf("failed to fork project after %d attempts", maxRetries)
}

func addWebhook(client *gitlab.Client, projectID int, webhookURL, token string) error {
	pushEvents := true
	mergeRequestsEvents := true
	hookOptions := &gitlab.AddProjectHookOptions{
		URL:                 &webhookURL,
		PushEvents:          &pushEvents,
		MergeRequestsEvents: &mergeRequestsEvents,
		Token:               &token,
	}
	_, _, err := client.Projects.AddProjectHook(projectID, hookOptions)
	if err != nil {
		return fmt.Errorf("failed to add webhook: %w", err)
	}
	return nil
}

func createBranch(client *gitlab.Client, projectID int, branchName string) error {
	_, _, err := client.Branches.CreateBranch(projectID, &gitlab.CreateBranchOptions{
		Branch: gitlab.Ptr(branchName),
		Ref:    gitlab.Ptr("main"),
	})
	return err
}

func createCommit(client *gitlab.Client, projectID int, branchName string, commitMessage string) error {
	actionValue := gitlab.FileActionValue("create")
	filePath := "preview.md"
	content := "Preview changes for the new feature"

	actions := []*gitlab.CommitActionOptions{{
		Action:   &actionValue,
		FilePath: &filePath,
		Content:  &content,
	}}
	commitOptions := gitlab.CreateCommitOptions{
		Branch:        &branchName,
		CommitMessage: &commitMessage,
		Actions:       actions,
	}
	_, _, err := client.Commits.CreateCommit(projectID, &commitOptions)
	return err
}

func createMergeRequest(client *gitlab.Client, projectID int, sourceBranch, targetBranch, title string) (string, error) {
	mrOptions := &gitlab.CreateMergeRequestOptions{
		SourceBranch: &sourceBranch,
		TargetBranch: &targetBranch,
		Title:        &title,
	}
	mr, _, err := client.MergeRequests.CreateMergeRequest(projectID, mrOptions)
	if err != nil {
		return "", err
	}
	return mr.WebURL, nil
}

func deleteGitlabProject(client *gitlab.Client, projectID int) error {
	_, err := client.Projects.DeleteProject(projectID)
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}
	log.Println("Project successfully deleted.")
	return nil
}

func InitGitLabClient() *gitlab.Client {
	privateToken := store.GetScenarioData("gitlabToken")
	client, err := gitlab.NewClient(privateToken)
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("failed to create GitLab client: %w", err))
	}
	return client
}

func SmeeDeployment(elname string) {
	var err error
	smeeDeploymentName := "gosmee-client"
	store.PutScenarioData("smee_deployment_name", smeeDeploymentName)

	smeeURL, err := getNewSmeeURL()
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("Failed to get a new Smee URL: %v", err))
	}
	store.PutScenarioData("SMEE_URL", smeeURL)

	routeurl := triggers.GetRoute(elname, store.Namespace())
	store.PutScenarioData("route", routeurl)
	store.PutScenarioData("elname", elname)

	if err = createSmeeDeployment(store.Clients(), store.Namespace(), smeeURL, routeurl); err != nil {
		testsuit.T.Fail(fmt.Errorf("Failed to create deployment: %v", err))
	}
}

func SetupGitLabProject(client *gitlab.Client) *gitlab.Project {

	gitlabGroupNamespace := config.GitlabGroupNamespace
	projectIDOrPath := config.ProjectID
	smeeURL := store.GetScenarioData("SMEE_URL")
	token := config.TriggersSecretToken

	project, err := forkProject(client, projectIDOrPath, gitlabGroupNamespace)
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("error during project forking: %w", err))
	}

	defer func() {
		if err != nil {
			if cleanupErr := deleteGitlabProject(client, project.ID); cleanupErr != nil {
				testsuit.T.Fail(fmt.Errorf("cleanup failed: %v", cleanupErr))
			}
		}
	}()

	err = addWebhook(client, project.ID, smeeURL, token)
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("failed to add webhook: %w", err))
	}

	return project
}

func ConfigurePreviewChanges(client *gitlab.Client, projectID int) {

	randomSuffix := strconv.FormatInt(time.Now().UnixNano(), 10)[:8]
	branchName := "preview-branch-" + randomSuffix
	commitMessage := "Add preview changes for feature"

	if err := createBranch(client, projectID, branchName); err != nil {
		testsuit.T.Fail(fmt.Errorf("Failed to create branch: %v", err))
	}

	if err := createCommit(client, projectID, branchName, commitMessage); err != nil {
		testsuit.T.Fail(fmt.Errorf("Failed to create commit: %v", err))
	}

	mrURL, err := createMergeRequest(client, projectID, branchName, "main", "Add preview changes for feature")
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("Failed to create merge request: %v", err))
	}

	log.Printf("Merge Request Created: %s\n", mrURL)
}

func CleanupPAC(c *clients.Clients, elName, smeeDeploymentName, namespace string) {
	// Delete EventListener
	err := c.TriggersClient.TriggersV1alpha1().EventListeners(namespace).Delete(c.Ctx, elName, metav1.DeleteOptions{})
	if err != nil {
		testsuit.T.Fail(err)
	}

	log.Println("Deleted EventListener")

	// Verify the EventListener's Deployment is deleted
	err = wait.WaitFor(c.Ctx, wait.DeploymentNotExist(c, namespace, fmt.Sprintf("%s-%s", eventReconciler.GeneratedResourcePrefix, elName)))
	if err != nil {
		testsuit.T.Fail(err)
	}

	log.Println("EventListener's Deployment was deleted")

	// Verify the EventListener's Service is deleted
	err = wait.WaitFor(c.Ctx, wait.ServiceNotExist(c, namespace, fmt.Sprintf("%s-%s", eventReconciler.GeneratedResourcePrefix, elName)))
	if err != nil {
		testsuit.T.Fail(err)
	}

	log.Println("EventListener's Service was deleted")

	// Delete Route exposed earlier
	err = c.Route.Routes(namespace).Delete(c.Ctx, fmt.Sprintf("%s-%s", eventReconciler.GeneratedResourcePrefix, elName), metav1.DeleteOptions{})
	if err != nil {
		testsuit.T.Fail(err)
	}

	// Verify the EventListener's Route is deleted
	err = wait.WaitFor(c.Ctx, wait.RouteNotExist(c, namespace, fmt.Sprintf("%s-%s", eventReconciler.GeneratedResourcePrefix, elName)))
	if err != nil {
		testsuit.T.Fail(err)
	}
	log.Println("EventListener's Route got deleted successfully...")

	// Delete Smee Deployment
	err = k8s.DeleteDeployment(c, namespace, smeeDeploymentName)
	if err != nil {
		testsuit.T.Fail(err)
	}

	log.Println("Deleted Smee Deployment")

	// This is required when EL runs as TLS
	cmd.MustSucceed("rm", "-rf", os.Getenv("GOPATH")+"/src/github.com/openshift-pipelines/release-tests/testdata/triggers/certs")
}
