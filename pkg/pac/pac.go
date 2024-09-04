package pac

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/openshift-pipelines/release-tests/pkg/clients"
	"github.com/xanzy/go-gitlab"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetNewSmeeURL() (string, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var smeeURL string
	err := chromedp.Run(ctx,
		chromedp.Navigate(`https://smee.io/new`),
		chromedp.WaitVisible(`body`),
		chromedp.Location(&smeeURL),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create SmeeURL: %v", err)
	}
	return smeeURL, nil
}

func CreateSmeeDeployment(c *clients.Clients, namespace, smeeURL, targetURL string) error {
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

func InitGitLabClient(privateToken string) (*gitlab.Client, error) {
	client, err := gitlab.NewClient(privateToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create GitLab client: %w", err)
	}
	return client, nil
}

func ForkProject(client *gitlab.Client, projectID, targetGroupNamespace string) (*gitlab.Project, error) {
	randomSuffix := strconv.FormatInt(time.Now().UnixNano(), 10)[:8]
	projectName := "openshift-pipelines-test-fork-" + randomSuffix
	projectPath := "openshift-pipelines-test-fork-" + randomSuffix

	forkOptions := &gitlab.ForkProjectOptions{
		Namespace: &targetGroupNamespace,
		Name:      &projectName,
		Path:      &projectPath,
	}

	project, response, err := client.Projects.ForkProject(projectID, forkOptions)
	if err != nil {
		log.Printf("Failed to fork project: %v\n", err)
		if response != nil {
			if body, readErr := io.ReadAll(response.Body); readErr == nil {
				log.Printf("Response body: %s\n", string(body))
			}
		}
		return nil, err
	}
	return project, nil
}

func AddWebhook(client *gitlab.Client, projectID int, webhookURL string) error {
	pushEvents := true
	mergeRequestsEvents := true
	hookOptions := &gitlab.AddProjectHookOptions{
		URL:                 &webhookURL,
		PushEvents:          &pushEvents,
		MergeRequestsEvents: &mergeRequestsEvents,
	}
	_, _, err := client.Projects.AddProjectHook(projectID, hookOptions)
	if err != nil {
		return fmt.Errorf("failed to add webhook: %w", err)
	}
	return nil
}

func CreateBranch(client *gitlab.Client, projectID int, branchName string) error {
	_, _, err := client.Branches.CreateBranch(projectID, &gitlab.CreateBranchOptions{
		Branch: gitlab.Ptr(branchName),
		Ref:    gitlab.Ptr("main"),
	})
	return err
}

func CreateCommit(client *gitlab.Client, projectID int, branchName string, commitMessage string) error {
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

func CreateMergeRequest(client *gitlab.Client, projectID int, sourceBranch, targetBranch, title string) (string, error) {
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

func ExtractMergeRequestID(mrURL string) (int, error) {
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

func CheckPipelineStatus(client *gitlab.Client, projectID int, mergeRequestID int) error {

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

func DeleteGitlabProject(client *gitlab.Client, projectID int) error {
	_, err := client.Projects.DeleteProject(projectID)
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}
	log.Println("Project successfully deleted.")
	return nil
}
