package pac

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/getgauge-contrib/gauge-go/testsuit"
	pacv1alpha1 "github.com/openshift-pipelines/pipelines-as-code/pkg/apis/pipelinesascode/v1alpha1"
	"github.com/openshift-pipelines/pipelines-as-code/pkg/cli"
	pacgenerate "github.com/openshift-pipelines/pipelines-as-code/pkg/cmd/tknpac/generate"
	"github.com/openshift-pipelines/pipelines-as-code/pkg/git"
	"github.com/openshift-pipelines/pipelines-as-code/pkg/params/info"
	"github.com/openshift-pipelines/release-tests/pkg/clients"
	"github.com/openshift-pipelines/release-tests/pkg/k8s"
	"github.com/openshift-pipelines/release-tests/pkg/oc"
	"github.com/openshift-pipelines/release-tests/pkg/pipelines"
	"github.com/openshift-pipelines/release-tests/pkg/store"
	gitlab "github.com/xanzy/go-gitlab"
	yaml "gopkg.in/yaml.v2"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	initialBackoffDuration   = 5 * time.Second
	maxRetriesForkProject    = 5
	maxRetriesPipelineStatus = 10
	targetURL                = "http://pipelines-as-code-controller.openshift-pipelines:8080"
	webhookConfigName        = "gitlab-webhook-config"
)

var client *gitlab.Client

func SetGitLabClient(c *gitlab.Client) {
	client = c
}

// Initialize Gitlab Client
func InitGitLabClient() *gitlab.Client {
	tokenSecretData := os.Getenv("GITLAB_TOKEN")
	webhookSecretData := os.Getenv("GITLAB_WEBHOOK_TOKEN")
	if tokenSecretData == "" && webhookSecretData == "" {
		testsuit.T.Fail(fmt.Errorf("token for authorization to the GitLab repository was not exported as a system variable"))
	} else {
		if !oc.SecretExists(webhookConfigName, store.Namespace()) {
			oc.CreateSecretForWebhook(tokenSecretData, webhookSecretData, store.Namespace())
		} else {
			log.Printf("Secret \"%s\" already exists", webhookConfigName)
		}
	}
	store.PutScenarioData("GITLAB_WEBHOOK_TOKEN", webhookSecretData)
	client, err := gitlab.NewClient(tokenSecretData)
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("failed to initialize GitLab client: %v", err))
	}

	return client
}

func getNewSmeeURL() (string, error) {
	// CURL cmd to retrieve a new smeeURL
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
	/*
		Reference for gosmee.yaml
			https://github.com/openshift-pipelines/pipelines-as-code
			/blob/main/pkg/cmd/tknpac/bootstrap/templates/gosmee.yaml
	*/
	replicas := int32(1)
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

func SetupSmeeDeployment() {
	var err error
	smeeDeploymentName := "gosmee-client"
	store.PutScenarioData("smeeDeploymentName", smeeDeploymentName)

	smeeURL, err := getNewSmeeURL()
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("failed to get a new Smee URL: %v", err))
	}
	store.PutScenarioData("SMEE_URL", smeeURL)

	if err = createSmeeDeployment(store.Clients(), store.Namespace(), smeeURL); err != nil {
		testsuit.T.Fail(fmt.Errorf("failed to create deployment: %v", err))
	}
}

// Specified Gitlab Project ID is forked into Group Namespace
func forkProject(projectID, targetNamespace string) (*gitlab.Project, error) {
	for i := 0; i < maxRetriesForkProject; i++ {
		projectName := fmt.Sprintf("release-tests-fork-%08d", time.Now().UnixNano()%1e8)
		project, _, err := client.Projects.ForkProject(projectID, &gitlab.ForkProjectOptions{
			Namespace: &targetNamespace,
			Name:      &projectName,
			Path:      &projectName,
		})
		if err == nil {
			store.PutScenarioData("PROJECT_URL", project.WebURL)
			return project, nil
		}
		log.Printf("Retry %d: failed to fork project: %v", i+1, err)
		time.Sleep(time.Duration(i+1) * time.Second)
	}
	return nil, fmt.Errorf("failed to fork project after %d attempts", maxRetriesForkProject)
}

// Add WebhookURL to forked Project
func addWebhook(projectID int, webhookURL, token string) error {
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

// Create a new Repository under current namespace
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
					Name: webhookConfigName,
					Key:  "provider.token",
				},
				WebhookSecret: &pacv1alpha1.Secret{
					Name: webhookConfigName,
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

// addLabelToProject adds a label to a GitLab project
func addLabelToProject(projectID int, labelName, color, description string) error {
	// Check if the label already exists
	labels, _, err := client.Labels.ListLabels(projectID, &gitlab.ListLabelsOptions{})
	if err != nil {
		return fmt.Errorf("failed to fetch project labels: %w", err)
	}

	for _, label := range labels {
		if label.Name == labelName {
			fmt.Printf("Label '%s' already exists in project ID %d\n", labelName, projectID)
			return nil
		}
	}

	// Create label if it doesn't exist
	_, _, err = client.Labels.CreateLabel(projectID, &gitlab.CreateLabelOptions{
		Name:        gitlab.Ptr(labelName),
		Color:       gitlab.Ptr(color),
		Description: gitlab.Ptr(description),
	})

	if err != nil {
		return fmt.Errorf("failed to create label '%s': %w", labelName, err)
	}

	fmt.Printf("Successfully added label '%s' to project ID %d\n", labelName, projectID)
	return nil
}

func SetupGitLabProject() *gitlab.Project {
	gitlabGroupNamespace := os.Getenv("GITLAB_GROUP_NAMESPACE")
	projectIDOrPath := os.Getenv("GITLAB_PROJECT_ID")

	if gitlabGroupNamespace == "" || projectIDOrPath == "" {
		testsuit.T.Fail(fmt.Errorf("failed to get system variables"))
	}

	smeeURL := store.GetScenarioData("SMEE_URL")
	webhookToken := store.GetScenarioData("GITLAB_WEBHOOK_TOKEN")

	project, err := forkProject(projectIDOrPath, gitlabGroupNamespace)
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("error during project forking: %w", err))
	}

	err = addWebhook(project.ID, smeeURL, webhookToken)
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("failed to add webhook: %w", err))
	}

	err = createNewRepository(store.Clients(), project.Name, gitlabGroupNamespace, store.Namespace())
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("failed to create repository"))
	}
	store.PutScenarioData("projectID", strconv.Itoa(project.ID))

	return project
}

// adds a comment to the specified merge request.
func AddComment(comment string) error {

	projectID, _ := strconv.Atoi(store.GetScenarioData("projectID"))
	mrID, _ := strconv.Atoi(store.GetScenarioData("mrID"))
	opts := &gitlab.CreateMergeRequestNoteOptions{
		Body: gitlab.Ptr(comment),
	}

	_, _, err := client.Notes.CreateMergeRequestNote(projectID, mrID, opts)
	if err != nil {
		return fmt.Errorf("failed to add comment to MR %d in project %d: %v", mrID, projectID, err)
	}

	return nil
}

func AddLabel(label, color, description string) error {

	projectID, _ := strconv.Atoi(store.GetScenarioData("projectID"))
	mrID, _ := strconv.Atoi(store.GetScenarioData("mrID"))

	// Add a label to the project
	err := addLabelToProject(projectID, label, color, description)
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("failed to add label to project: %w", err))
	}
	// Create a LabelOptions instance
	addLabels := gitlab.LabelOptions{label}

	// Update the merge request to add the label
	_, _, err = client.MergeRequests.UpdateMergeRequest(projectID, mrID, &gitlab.UpdateMergeRequestOptions{
		AddLabels: &addLabels,
	})
	if err != nil {
		return fmt.Errorf("failed to update merge request with label 'bug': %w", err)
	}

	fmt.Printf("Successfully added label %s to merge request %d\n", label, mrID)
	return nil
}

// Create new branch to push the commit
func createBranch(projectID int, branchName string) error {
	_, _, err := client.Branches.CreateBranch(projectID, &gitlab.CreateBranchOptions{
		Branch: gitlab.Ptr(branchName),
		Ref:    gitlab.Ptr("main"),
	})
	return err
}

func createPacGenerateOpts(eventType, branch, fileName string) *pacgenerate.Opts {
	// Initialize PAC generate options
	opts := pacgenerate.MakeOpts()

	// Set Event information
	opts.Event = &info.Event{
		EventType:  eventType,
		BaseBranch: branch,
	}
	// Set Project URL and Branch name to GitInfo
	// ProjectURL is used as PipelineRun name with suffix
	opts.GitInfo = &git.Info{
		URL:    store.GetScenarioData("PROJECT_URL"),
		Branch: branch,
	}
	// Initialize I/O streams
	var outputBuffer bytes.Buffer
	opts.IOStreams = &cli.IOStreams{
		Out:    &outputBuffer,
		ErrOut: os.Stderr,
		In:     os.Stdin,
	}
	// Specify the FileName of the pipelinerun yaml
	opts.FileName = fileName

	return opts
}

// Generate sample PipelineRun, pull-request.yaml
func generatePipelineRun(eventType, branch, fileName string) error {
	opts := createPacGenerateOpts(eventType, branch, fileName)

	err := pacgenerate.Generate(opts, true)
	if err != nil {
		return fmt.Errorf("failed to generate PipelineRun: %v", err)
	}

	return nil
}

// Validate generated yaml file from pac generate cmd
func validateYAML(yamlContent []byte) error {
	var content map[string]interface{}

	err := yaml.Unmarshal(yamlContent, &content)
	if err != nil {
		return fmt.Errorf("invalid YAML format: %v", err)
	}

	return nil
}

func GeneratePipelineRunYaml(eventType, branch string) error {
	fileName := eventType + ".yaml"

	// Generate the PipelineRun YAML.
	if err := generatePipelineRun(eventType, branch, fileName); err != nil {
		return fmt.Errorf("could not generate PipelineRun in %s: %v", fileName, err)
	}

	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("could not read file %s: %v", fileName, err)
	}

	if err := validateYAML(fileContent); err != nil {
		return fmt.Errorf("invalid YAML content: %v", err)
	}
	store.PutScenarioData("fileContent", string(fileContent))
	store.PutScenarioData("branch", string(branch))
	store.PutScenarioData("fileName", string(fileName))

	return nil
}

// updateAnnotation updates the specified annotation in the pull-request.yaml file
func UpdateAnnotation(annotationKey, annotationValue string) error {
	fileName := store.GetScenarioData("fileName")
	data, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("failed to read YAML file: %v", err)
	}

	var content map[string]interface{}
	if err := yaml.Unmarshal(data, &content); err != nil {
		return fmt.Errorf("failed to unmarshal YAML: %v", err)
	}

	meta := content["metadata"].(map[interface{}]interface{})
	anns := meta["annotations"].(map[interface{}]interface{})

	// If the annotation exists, append the new value; otherwise, set it.
	if currValue, exists := anns[annotationKey].(string); exists {
		anns[annotationKey] = currValue + " " + annotationValue
	} else {
		anns[annotationKey] = annotationValue
	}

	out, err := yaml.Marshal(content)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %v", err)
	}

	if err := os.WriteFile(fileName, out, 0644); err != nil {
		return fmt.Errorf("failed to write YAML file: %v", err)
	}

	if err := validateYAML(out); err != nil {
		return fmt.Errorf("invalid YAML content: %v", err)
	}

	store.PutScenarioData("fileContent", string(out))

	fmt.Println("Annotation updated successfully")
	return nil
}

func createCommit(projectID int, branch, commitMessage, fileDesPath string) error {
	// Commit the PAC generated PLR
	fileContent := store.GetScenarioData("fileContent")
	action := gitlab.FileCreate
	commitOpts := &gitlab.CreateCommitOptions{
		Branch:        &branch,
		CommitMessage: &commitMessage,
		Actions: []*gitlab.CommitActionOptions{
			{Action: &action, FilePath: gitlab.Ptr(fileDesPath), Content: gitlab.Ptr(fileContent)},
		},
	}

	_, _, err := client.Commits.CreateCommit(projectID, commitOpts)
	if err != nil {
		return fmt.Errorf("failed to create commit: %v", err)
	}
	return nil
}

// Creates MR to a forked project with a PipelineRun YAML under .tekton directory
func createMergeRequest(projectID int, sourceBranch, targetBranch, title string) (string, error) {
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

// Extract the MR ID to check the Pipeline status
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
	return status == "success" || status == "failed" || status == "canceled"
}

func checkPipelineStatus(projectID, mergeRequestID int) error {
	var retryCount int

	// Fetch pipelines for the specified MergeRequest ID
	for {
		pipelines, _, err := client.MergeRequests.ListMergeRequestPipelines(projectID, mergeRequestID)
		if err != nil {
			return fmt.Errorf("failed to list merge request pipelines: %w", err)
		}

		// Retry, If no pipelines found
		if len(pipelines) == 0 {
			if retryCount >= maxRetriesPipelineStatus {
				log.Printf("No pipelines found for the MR id %d after %d retries\n", mergeRequestID, maxRetriesPipelineStatus)
				return nil
			}
			log.Println("No pipelines found, retrying...")
			retryCount++
			time.Sleep(time.Duration(2^retryCount) * time.Second)
			continue
		}

		// Check the status of the latest pipeline
		latestPipeline := pipelines[0]
		if isTerminalStatus(latestPipeline.Status) {
			log.Printf("Latest pipeline status for MR #%d: %s\n", mergeRequestID, latestPipeline.Status)
			return nil
		} else {
			log.Println("waiting for Pipeline status to be updated...")
			time.Sleep(maxRetriesPipelineStatus * time.Second)
		}
	}
}

func ConfigurePreviewChanges() {
	randomSuffix := strconv.FormatInt(time.Now().UnixNano(), 10)[:8]
	branchName := "preview-branch-" + randomSuffix
	commitMessage := "Add preview changes for feature"
	fileDesPath := ".tekton/pull-request.yaml"

	projectID, err := strconv.Atoi(store.GetScenarioData("projectID"))
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("failed to convert project ID to integer: %v", err))
	}

	if err := createBranch(projectID, branchName); err != nil {
		testsuit.T.Fail(fmt.Errorf("failed to create branch: %v", err))
	}

	if err := createCommit(projectID, branchName, commitMessage, fileDesPath); err != nil {
		testsuit.T.Fail(fmt.Errorf("failed to create commit: %v", err))
	}

	mrURL, err := createMergeRequest(projectID, branchName, "main", "Add preview changes for feature")
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("failed to create merge request: %v", err))
	}

	log.Printf("Merge Request Created: %s\n", mrURL)

	mrID, err := extractMergeRequestID(mrURL)
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("failed to extract merge request ID: %v", err))
	}

	store.PutScenarioData("mrID", strconv.Itoa(mrID))
}

func GetPipelineNameFromMR() (pipelineName string) {
	projectID, err := strconv.Atoi(store.GetScenarioData("projectID"))
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("failed to convert project ID to integer: %v", err))
	}
	mrID, err := strconv.Atoi(store.GetScenarioData("mrID"))
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("failed to convert project ID to integer: %v", err))
	}

	err = checkPipelineStatus(projectID, mrID)
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("failed to check pipeline status: %v", err))
	}

	pipelineName, err = pipelines.GetLatestPipelinerun(store.Clients(), store.Namespace())
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("failed to get the latest Pipelinerun: %v", err))
	}
	return pipelineName
}

func deleteGitlabProject(projectID int) error {
	_, err := client.Projects.DeleteProject(projectID)
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}
	log.Println("Project successfully deleted.")
	return nil
}

func CleanupPAC(c *clients.Clients, smeeDeploymentName, namespace string) {
	// Remove the generated PipelineRun YAML file
	fileName := store.GetScenarioData("fileName")
	if fileName != "" {
		if err := os.Remove(fileName); err != nil {
			testsuit.T.Fail(fmt.Errorf("failed to remove file %s: %v", fileName, err))
		}
	}

	projectID, err := strconv.Atoi(store.GetScenarioData("projectID"))
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("failed to convert project ID to integer: %v", err))
	}
	// Remove Forked Project
	if cleanupErr := deleteGitlabProject(projectID); cleanupErr != nil {
		testsuit.T.Fail(fmt.Errorf("cleanup failed: %v", cleanupErr))
	}

	// Delete Smee Deployment
	err = k8s.DeleteDeployment(c, namespace, smeeDeploymentName)
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("failed to Delete Smee Deployment: %v", err))
	}
}
