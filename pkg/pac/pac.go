package pac

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/getgauge-contrib/gauge-go/testsuit"
	pacv1alpha1 "github.com/openshift-pipelines/pipelines-as-code/pkg/apis/pipelinesascode/v1alpha1"
	"github.com/openshift-pipelines/release-tests/pkg/clients"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	"github.com/openshift-pipelines/release-tests/pkg/k8s"
	"github.com/openshift-pipelines/release-tests/pkg/oc"
	"github.com/openshift-pipelines/release-tests/pkg/pipelines"
	"github.com/openshift-pipelines/release-tests/pkg/store"
	"github.com/xanzy/go-gitlab"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func createNewRepository(c *clients.Clients, projectName, targetGroupNamespace, namespace string) error {

	repo := &pacv1alpha1.Repository{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "pipelinesascode.tekton.dev/v1alpha1",
			Kind:       "Repository",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      projectName,
			Namespace: namespace,
		},
		Spec: pacv1alpha1.RepositorySpec{
			URL: fmt.Sprintf("https://gitlab.com/%s/%s", targetGroupNamespace, projectName),
			GitProvider: &pacv1alpha1.GitProvider{
				URL: "https://gitlab.com",
				Secret: &pacv1alpha1.Secret{
					Name: "gitlab-webhook-config",
					Key:  "provider.token",
				},
				WebhookSecret: &pacv1alpha1.Secret{
					Name: "gitlab-webhook-config",
					Key:  "webhook.secret",
				},
			},
		},
	}

	repo, err := c.PacClientset.Repositories(namespace).Create(context.Background(), repo, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create repository: %v", err)
	}

	log.Printf("Repository '%s' created successfully in namespace '%s'", repo.GetName(), repo.GetNamespace())
	return nil
}

func ConfigureGitlabToken() {
	tokenSecretData := os.Getenv("GITLAB_TOKEN")
	webhookSecretData := os.Getenv("WEBHOOK_TOKEN")
	if tokenSecretData == "" && webhookSecretData == "" {
		testsuit.T.Fail(fmt.Errorf("Token for authorization to the Gitlab repository was not exported as a system variable"))
	} else {
		if !oc.SecretExists("gitlab-webhook-config", store.Namespace()) {
			oc.CreateSecretForWebhook(tokenSecretData, webhookSecretData, store.Namespace())
		} else {
			log.Printf("Secret \"gitlab-webhook-config\" already exists")
		}
		store.PutScenarioData("gitlabToken", tokenSecretData)
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

func createSmeeDeployment(c *clients.Clients, namespace, smeeURL string) error {
	replicas := int32(1)
	targetURL := "http://pipelines-as-code-controller.openshift-pipelines:8080"
	deployment := &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "gosmee-client",
		},
		Spec: v1.DeploymentSpec{
			Replicas: &replicas,
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
	noteEvents := true
	tagPushEvents := true

	hookOptions := &gitlab.AddProjectHookOptions{
		URL:                 &webhookURL,
		PushEvents:          &pushEvents,
		MergeRequestsEvents: &mergeRequestsEvents,
		NoteEvents:          &noteEvents,
		TagPushEvents:       &tagPushEvents,
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

func extractMergeRequestID(mrURL string) (int, error) {
	parsedURL, err := url.Parse(mrURL)
	if err != nil {
		return 0, fmt.Errorf("failed to parse merge request URL: %w", err)
	}
	segments := strings.Split(parsedURL.Path, "/")
	mrIDStr := segments[len(segments)-1]
	mrID, err := strconv.Atoi(mrIDStr)
	if err != nil {
		return 0, fmt.Errorf("failed to convert MR ID to integer: %w", err)
	}
	return mrID, nil
}

func isTerminalStatus(status string) bool {
	switch status {
	case "success", "failed", "canceled":
		return true
	default:
		return false
	}
}

func checkPipelineStatus(client *gitlab.Client, projectID, mergeRequestID int) error {

	const maxRetries = 10
	var retryCount int

	for {
		pipelines, _, err := client.MergeRequests.ListMergeRequestPipelines(projectID, mergeRequestID)
		if err != nil {
			return fmt.Errorf("failed to list merge request pipelines: %w", err)
		}

		if len(pipelines) == 0 {
			if retryCount >= maxRetries {
				log.Printf("No pipelines found for the MR id %d after %d retries\n", mergeRequestID, maxRetries)
				return nil
			}
			log.Println("No pipelines found, retrying...")
			retryCount++
			time.Sleep(time.Duration(2^retryCount) * time.Second)
			continue
		}

		latestPipeline := pipelines[0]

		if isTerminalStatus(latestPipeline.Status) {
			log.Printf("Latest pipeline status for MR #%d: %s\n", mergeRequestID, latestPipeline.Status)
			return nil
		} else {
			log.Println("waiting for Pipeline status to be updated...")
			time.Sleep(30 * time.Second)
		}
	}
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

func SetupSmeeDeployment() {
	var err error
	smeeDeploymentName := "gosmee-client"
	store.PutScenarioData("smee_deployment_name", smeeDeploymentName)

	smeeURL, err := getNewSmeeURL()
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("Failed to get a new Smee URL: %v", err))
	}
	store.PutScenarioData("SMEE_URL", smeeURL)

	if err = createSmeeDeployment(store.Clients(), store.Namespace(), smeeURL); err != nil {
		testsuit.T.Fail(fmt.Errorf("Failed to create deployment: %v", err))
	}
}

func SetupGitLabProject(client *gitlab.Client) *gitlab.Project {
	gitlabGroupNamespace := os.Getenv("GITLAB_GROUP_NAMESPACE")
	projectIDOrPath := os.Getenv("GITLAB_PROJECT_ID")

	if gitlabGroupNamespace == "" || projectIDOrPath == "" {
		testsuit.T.Fail(fmt.Errorf("Failed to get system variables"))
	}

	smeeURL := store.GetScenarioData("SMEE_URL")
	token := config.TriggersSecretToken

	project, err := forkProject(client, projectIDOrPath, gitlabGroupNamespace)
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("error during project forking: %w", err))
	}

	err = addWebhook(client, project.ID, smeeURL, token)
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("failed to add webhook: %w", err))
	}

	err = createNewRepository(store.Clients(), project.Name, gitlabGroupNamespace, store.Namespace())
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("failed to create repository"))
	}

	return project
}

func ConfigurePreviewChanges(client *gitlab.Client, projectID int) string {

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

	mrID, err := extractMergeRequestID(mrURL)
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("Failed to extract merge request ID: %v", err))
	}

	err = checkPipelineStatus(client, projectID, mrID)
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("Failed to check pipeline status: %v", err))
	}

	pipelineName, err := pipelines.GetLatestPipelinerun(store.Clients(), store.Namespace())
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("Failed to get the latest Pipelinerun: %v", err))
	}
	return pipelineName
}

func CleanupPAC(client *gitlab.Client, c *clients.Clients, projectID int, smeeDeploymentName, namespace string) {

	// Remove Created Project
	if cleanupErr := deleteGitlabProject(client, projectID); cleanupErr != nil {
		testsuit.T.Fail(fmt.Errorf("cleanup failed: %v", cleanupErr))
	}

	// Delete Smee Deployment
	err := k8s.DeleteDeployment(c, namespace, smeeDeploymentName)
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("Failed to Delete Smee Deployment: %v", err))
	}
}
