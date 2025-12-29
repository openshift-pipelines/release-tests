package pac

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/google/go-github/v74/github"
	pacv1alpha1 "github.com/openshift-pipelines/pipelines-as-code/pkg/apis/pipelinesascode/v1alpha1"
	"github.com/openshift-pipelines/release-tests/pkg/clients"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	"github.com/openshift-pipelines/release-tests/pkg/k8s"
	"github.com/openshift-pipelines/release-tests/pkg/pipelines"
	"github.com/openshift-pipelines/release-tests/pkg/store"
	"golang.org/x/oauth2"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const githubWebhookConfigName = "github-webhook-config"

const pacEventTypeAnnotationKey = "pipelinesascode.tekton.dev/event-type"

var ghClient *github.Client

func SetGitHubClient(c *github.Client) {
	ghClient = c
}

// InitGitHubClient initializes a GitHub client for GitHub
func InitGitHubClient() *github.Client {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		testsuit.T.Fail(fmt.Errorf("GITHUB_TOKEN was not exported as a system variable"))
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(context.Background(), ts)
	return github.NewClient(tc)
}

func randWebhookSecret() (string, error) {
	b := make([]byte, 30)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func ensureWebhookSecret(c *clients.Clients, namespace, token, webhookSecret string) error {
	secretsClient := c.KubeClient.Kube.CoreV1().Secrets(namespace)

	want := map[string]string{
		"provider.token": token,
		"webhook.secret": webhookSecret,
	}

	existing, err := secretsClient.Get(context.Background(), githubWebhookConfigName, metav1.GetOptions{})
	if err == nil {
		if existing.StringData == nil {
			existing.StringData = map[string]string{}
		}
		for k, v := range want {
			existing.StringData[k] = v
		}
		_, err = secretsClient.Update(context.Background(), existing, metav1.UpdateOptions{})
		return err
	}
	if !apierrors.IsNotFound(err) {
		return err
	}

	sec := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      githubWebhookConfigName,
			Namespace: namespace,
		},
		Type:       corev1.SecretTypeOpaque,
		StringData: want,
	}
	_, err = secretsClient.Create(context.Background(), sec, metav1.CreateOptions{})
	return err
}

func createGitHubRepositoryCR(c *clients.Clients, repoName, repoURL, namespace string) error {
	repo := &pacv1alpha1.Repository{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "pipelinesascode.tekton.dev/v1alpha1",
			Kind:       "Repository",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      repoName,
			Namespace: namespace,
		},
		Spec: pacv1alpha1.RepositorySpec{
			URL: repoURL,
			Settings: &pacv1alpha1.Settings{
				PipelineRunProvenance: "source",
			},
			GitProvider: &pacv1alpha1.GitProvider{
				Secret: &pacv1alpha1.Secret{
					Name: githubWebhookConfigName,
					Key:  "provider.token",
				},
				WebhookSecret: &pacv1alpha1.Secret{
					Name: githubWebhookConfigName,
					Key:  "webhook.secret",
				},
			},
		},
	}

	if _, err := c.PacClientset.Repositories(namespace).Create(context.Background(), repo, metav1.CreateOptions{}); err != nil {
		return err
	}
	store.PutScenarioData("PAC_REPOSITORY_CR_NAME", repoName)
	return nil
}

func waitForRepoReady(ctx context.Context, owner, repo string) error {
	deadline := time.Now().Add(60 * time.Second)
	for time.Now().Before(deadline) {
		_, resp, err := ghClient.Repositories.Get(ctx, owner, repo)
		if err == nil {
			return nil
		}
		if resp != nil && resp.StatusCode == 404 {
			time.Sleep(2 * time.Second)
			continue
		}
		return err
	}
	return fmt.Errorf("timed out waiting for github repo %s/%s to be ready", owner, repo)
}

func ensureDefaultBranchMain(ctx context.Context, owner, repo, defaultBranch string) error {
	if defaultBranch == "" || defaultBranch == "main" {
		return nil
	}
	_, _, err := ghClient.Repositories.RenameBranch(ctx, owner, repo, defaultBranch, "main")
	return err
}

func addGitHubWebhook(ctx context.Context, owner, repo, smeeURL, webhookSecret string) error {
	hook := &github.Hook{
		Active: github.Ptr(true),
		Config: &github.HookConfig{
			URL:         github.Ptr(smeeURL),
			ContentType: github.Ptr("json"),
			Secret:      github.Ptr(webhookSecret),
			InsecureSSL: github.Ptr("0"),
		},
		Events: []string{"commit_comment", "issue_comment", "pull_request", "push"},
	}
	_, _, err := ghClient.Repositories.CreateHook(ctx, owner, repo, hook)
	return err
}

func sanitizeK8sName(in string) string {
	s := strings.ToLower(in)
	out := make([]byte, 0, len(s))
	for i := range s {
		ch := s[i]
		if (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') || ch == '-' {
			out = append(out, ch)
		} else {
			out = append(out, '-')
		}
	}
	res := strings.Trim(string(out), "-")
	if res == "" {
		return "pac-repo"
	}
	if len(res) > 63 {
		return res[:63]
	}
	return res
}

func createBranchWithGitHub(ctx context.Context, owner, repo, baseBranch, newBranch, message string, files map[string]string) error {
	baseRef, _, err := ghClient.Git.GetRef(ctx, owner, repo, "refs/heads/"+baseBranch)
	if err != nil {
		return err
	}
	baseCommitSHA := baseRef.GetObject().GetSHA()
	if baseCommitSHA == "" {
		return fmt.Errorf("base branch %q has empty SHA", baseBranch)
	}

	baseCommit, _, err := ghClient.Git.GetCommit(ctx, owner, repo, baseCommitSHA)
	if err != nil {
		return err
	}
	baseTreeSHA := ""
	if baseCommit.Tree != nil {
		baseTreeSHA = baseCommit.Tree.GetSHA()
	}
	if baseTreeSHA == "" {
		return fmt.Errorf("base commit %s has empty tree SHA", baseCommitSHA)
	}

	entries := make([]*github.TreeEntry, 0, len(files))
	for p, c := range files {
		path := p
		content := c
		entries = append(entries, &github.TreeEntry{
			Path:    &path,
			Mode:    github.Ptr("100644"),
			Type:    github.Ptr("blob"),
			Content: &content,
		})
	}

	newTree, _, err := ghClient.Git.CreateTree(ctx, owner, repo, baseTreeSHA, entries)
	if err != nil {
		return err
	}

	commit := &github.Commit{
		Message: github.Ptr(message),
		Tree:    newTree,
		Parents: []*github.Commit{{SHA: github.Ptr(baseCommitSHA)}},
	}
	newCommit, _, err := ghClient.Git.CreateCommit(ctx, owner, repo, commit, nil)
	if err != nil {
		return err
	}
	newCommitSHA := newCommit.GetSHA()
	if newCommitSHA == "" {
		return fmt.Errorf("created commit has empty SHA")
	}

	ref := &github.Reference{
		Ref: github.Ptr("refs/heads/" + newBranch),
		Object: &github.GitObject{
			SHA: github.Ptr(newCommitSHA),
		},
	}
	_, _, err = ghClient.Git.CreateRef(ctx, owner, repo, ref)
	return err
}

// SetupGitHubProject creates a new GitHub repository
func SetupGitHubProject() *github.Repository {
	if ghClient == nil {
		testsuit.T.Fail(fmt.Errorf("github client not initialized; call InitGitHubClient/SetGitHubClient first"))
	}

	ctx := context.Background()
	org := os.Getenv("GITHUB_ORG")
	smeeURL := store.GetScenarioData("SMEE_URL")
	token := os.Getenv("GITHUB_TOKEN")

	webhookSecret := os.Getenv("GITHUB_WEBHOOK_TOKEN")
	if webhookSecret == "" {
		sec, err := randWebhookSecret()
		if err != nil {
			testsuit.T.Fail(fmt.Errorf("failed generating github webhook secret: %v", err))
		}
		webhookSecret = sec
	}

	repoName := fmt.Sprintf("release-tests-pac-%08d", time.Now().UnixNano()%1e8)
	createReq := &github.Repository{
		Name:                github.Ptr(repoName),
		Visibility:          github.Ptr("public"),
		AutoInit:            github.Ptr(true),
		AllowSquashMerge:    github.Ptr(true),
		DeleteBranchOnMerge: github.Ptr(true),
	}
	created, _, err := ghClient.Repositories.Create(ctx, org, createReq)
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("failed to create github repository: %v", err))
	}

	owner := ""
	if org != "" {
		owner = org
	} else if created.Owner != nil && created.Owner.Login != nil {
		owner = created.GetOwner().GetLogin()
	} else {
		u, _, uerr := ghClient.Users.Get(ctx, "")
		if uerr != nil {
			testsuit.T.Fail(fmt.Errorf("failed to determine github username: %v", uerr))
		}
		owner = u.GetLogin()
	}

	if err := waitForRepoReady(ctx, owner, repoName); err != nil {
		testsuit.T.Fail(err)
	}
	if err := ensureDefaultBranchMain(ctx, owner, repoName, created.GetDefaultBranch()); err != nil {
		testsuit.T.Fail(fmt.Errorf("failed to rename default branch to main: %v", err))
	}

	repoURL := created.GetHTMLURL()
	if repoURL == "" {
		repoURL = fmt.Sprintf("https://github.com/%s/%s", owner, repoName)
	}

	store.PutScenarioData("PROJECT_URL", repoURL)
	store.PutScenarioData("GITHUB_REPO_OWNER", owner)
	store.PutScenarioData("GITHUB_REPO_NAME", repoName)

	// Create the local webhook+token secret and Repository CR in the scenario namespace.
	namespace := store.Namespace()
	if err := ensureWebhookSecret(store.Clients(), namespace, token, webhookSecret); err != nil {
		testsuit.T.Fail(fmt.Errorf("failed to ensure github webhook secret: %v", err))
	}
	if err := createGitHubRepositoryCR(store.Clients(), sanitizeK8sName(repoName), repoURL, namespace); err != nil {
		testsuit.T.Fail(fmt.Errorf("failed to create PAC Repository CR: %v", err))
	}

	// Configure GitHub webhook to smee.io (gosmee forwards to controller service in-cluster).
	if err := addGitHubWebhook(ctx, owner, repoName, smeeURL, webhookSecret); err != nil {
		testsuit.T.Fail(fmt.Errorf("failed to add github webhook: %v", err))
	}

	log.Printf("GitHub repo created: %s", repoURL)
	return created
}

func waitForPRMergeable(ctx context.Context, owner, repo string, number int) error {
	deadline := time.Now().Add(60 * time.Second)
	for time.Now().Before(deadline) {
		pr, _, err := ghClient.PullRequests.Get(ctx, owner, repo, number)
		if err != nil {
			return err
		}
		if pr.Mergeable != nil && *pr.Mergeable {
			return nil
		}
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("timed out waiting for PR #%d to become mergeable", number)
}

func ConfigurePreviewChangesGitHub() {
	owner := store.GetScenarioData("GITHUB_REPO_OWNER")
	repo := store.GetScenarioData("GITHUB_REPO_NAME")
	ctx := context.Background()

	prData, err := os.ReadFile(filepath.Clean(pullRequestFileName))
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("read %s: %v", pullRequestFileName, err))
	}
	pushData, pushErr := os.ReadFile(filepath.Clean(pushFileName))
	hasPush := pushErr == nil

	branchName := fmt.Sprintf("preview-%08d", time.Now().UnixNano()%1e8)
	files := map[string]string{
		".tekton/pull-request.yaml": string(prData),
	}
	if hasPush {
		files[".tekton/push.yaml"] = string(pushData)
	}
	if err := createBranchWithGitHub(ctx, owner, repo, "main", branchName, "ci(pac): add pipelines-as-code definitions", files); err != nil {
		testsuit.T.Fail(fmt.Errorf("failed to create single-commit branch %q: %v", branchName, err))
	}

	newPR := &github.NewPullRequest{
		Title: github.Ptr("Add preview changes for feature"),
		Head:  github.Ptr(owner + ":" + branchName),
		Base:  github.Ptr("main"),
	}
	pr, _, err := ghClient.PullRequests.Create(ctx, owner, repo, newPR)
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("failed to create PR: %v", err))
	}
	store.PutScenarioData("prURL", pr.GetHTMLURL())
	store.PutScenarioData("prNumber", strconv.Itoa(pr.GetNumber()))
	log.Printf("Pull Request Created: %s", pr.GetHTMLURL())
}

func TriggerPushOnGitHubMain() {
	owner := store.GetScenarioData("GITHUB_REPO_OWNER")
	repo := store.GetScenarioData("GITHUB_REPO_NAME")
	ctx := context.Background()

	prNum, err := strconv.Atoi(store.GetScenarioData("prNumber"))
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("bad prNumber: %v", err))
	}

	if err := waitForPRMergeable(ctx, owner, repo, prNum); err != nil {
		testsuit.T.Fail(err)
	}

	_, _, err = ghClient.PullRequests.Merge(ctx, owner, repo, prNum, "", &github.PullRequestOptions{
		MergeMethod: "squash",
	})
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("failed to merge PR #%d: %v", prNum, err))
	}
}

func WaitForNewPipelineRunName(previousName string) string {
	deadline := time.Now().Add(config.APITimeout)
	for time.Now().Before(deadline) {
		name, err := pipelines.GetLatestPipelinerun(store.Clients(), store.Namespace())
		if err == nil {
			if previousName == "" || name != previousName {
				return name
			}
		}
		time.Sleep(config.APIRetry)
	}

	if previousName == "" {
		testsuit.T.Fail(fmt.Errorf("timed out waiting for a PipelineRun to be created in namespace %q", store.Namespace()))
	} else {
		testsuit.T.Fail(fmt.Errorf("timed out waiting for a new PipelineRun in namespace %q (previous=%q)", store.Namespace(), previousName))
	}
	return ""
}

// WaitForNewPipelineRunNameByEventType waits for a new PipelineRun with the given PaC event type
func WaitForNewPipelineRunNameByEventType(previousName, eventType string) string {
	deadline := time.Now().Add(config.APITimeout)
	for time.Now().Before(deadline) {
		prs, err := store.Clients().PipelineRunClient.List(store.Clients().Ctx, metav1.ListOptions{})
		if err == nil {
			var bestName string
			var bestStart time.Time
			var bestFound bool
			for _, pr := range prs.Items {
				if pr.Annotations == nil {
					continue
				}
				if pr.Annotations[pacEventTypeAnnotationKey] != eventType {
					continue
				}
				if previousName != "" && pr.Name == previousName {
					continue
				}
				start := pr.CreationTimestamp.Time
				if pr.Status.StartTime != nil {
					start = pr.Status.StartTime.Time
				}
				if !bestFound || start.After(bestStart) {
					bestFound = true
					bestStart = start
					bestName = pr.Name
				}
			}
			if bestFound {
				return bestName
			}
		}
		time.Sleep(config.APIRetry)
	}

	testsuit.T.Fail(fmt.Errorf("timed out waiting for a new PipelineRun with event-type=%q in namespace %q (previous=%q)", eventType, store.Namespace(), previousName))
	return ""
}

func CleanupPACGitHub(c *clients.Clients, smeeDeploymentName, namespace string) {
	_ = os.Remove(pullRequestFileName)
	_ = os.Remove(pushFileName)

	owner := store.GetScenarioData("GITHUB_REPO_OWNER")
	repo := store.GetScenarioData("GITHUB_REPO_NAME")
	if owner != "" && repo != "" && ghClient != nil {
		if _, err := ghClient.Repositories.Delete(context.Background(), owner, repo); err != nil {
			testsuit.T.Fail(fmt.Errorf("failed to delete github repository %s/%s: %v", owner, repo, err))
		}
	}

	if err := k8s.DeleteDeployment(c, namespace, smeeDeploymentName); err != nil {
		testsuit.T.Fail(fmt.Errorf("failed to delete smee deployment: %v", err))
	}
}
