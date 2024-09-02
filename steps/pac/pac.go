package pac

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/openshift-pipelines/release-tests/pkg/k8s"
	"github.com/openshift-pipelines/release-tests/pkg/pac"
	"github.com/openshift-pipelines/release-tests/pkg/store"
	"github.com/openshift-pipelines/release-tests/pkg/triggers"
)

var (
	smeeURL string
)

var _ = gauge.Step("Create Smee Deployment", func() {

	smeeURL = os.Getenv("SMEE_URL")
	targetURL := os.Getenv("TARGET_URL")

	if smeeURL == "" {
		var err error
		smeeURL, err = pac.GetNewSmeeURL()
		if err != nil {
			log.Fatalf("Failed to get a new Smee URL: %v", err)
		}
		os.Setenv("SMEE_URL", smeeURL)
	}

	if targetURL == "" {
		namespace := "openshift-pipelines"
		routeName := "pipelines-as-code-controller"

		targetURL = triggers.GetRouteURL(routeName, namespace)
		if targetURL == "" {
			log.Fatalf("Received an empty Route URL for route %s in namespace %s", routeName, namespace)
		}
	}

	err := pac.CreateSmeeDeployment(store.Clients(), store.Namespace(), smeeURL, targetURL)
	if err != nil {
		log.Fatalf("Failed to create deployment: %v", err)
	}
	k8s.ValidateDeployments(store.Clients(), store.Namespace(), "smee-client")

})

var _ = gauge.Step("Configure Gitlab repo", func() {

	smeeURL := os.Getenv("SMEE_URL")
	projectIDOrPath := "<ProjectID>"
	targetGroupNamespace := "<GroupNamespace>"
	privateToken := "<Private Token>"

	client, err := pac.InitGitLabClient(privateToken)
	if err != nil {
		log.Fatal(err)
	}

	project, err := pac.ForkProject(client, projectIDOrPath, targetGroupNamespace)
	if err != nil {
		log.Println("Error during project forking:", err)
		log.Fatal(err)
	}
	log.Printf("Project successfully forked: %s (Project ID: %d)\n", project.Name, project.ID)
	defer func() {
		if err := pac.DeleteGitlabProject(client, project.ID); err != nil {
			log.Printf("Cleanup failed: %v", err)
		}
	}()

	err = pac.AddWebhook(client, project.ID, smeeURL)
	if err != nil {
		log.Printf("Failed to add webhook: %v\n", err)
		log.Fatal(err)
	}
	log.Printf("Webhook added to %s\n", project.Name)

	randomSuffix := strconv.FormatInt(time.Now().UnixNano(), 10)[:8]
	branchName := "preview-branch-" + randomSuffix
	commitMessage := "Add preview changes for feature"

	if err := pac.CreateBranch(client, project.ID, branchName); err != nil {
		fmt.Println("Failed to create branch:", err)
		log.Fatal(err)
	}

	if err := pac.CreateCommit(client, project.ID, branchName, commitMessage); err != nil {
		fmt.Println("Failed to create commit:", err)
		log.Fatal(err)
	}

	mrURL, err := pac.CreateMergeRequest(client, project.ID, branchName, "main", "Add preview changes for feature")
	if err != nil {
		fmt.Println("Failed to create merge request:", err)
		log.Fatal(err)
	}

	fmt.Printf("Merge Request Created: %s\n", mrURL)

	mrID, err := pac.ExtractMergeRequestID(mrURL)
	if err != nil {
		log.Fatal(err)
	}
	// WIP: check the Pipelinerun status
	// if err := pac.CheckPipelineStatus(client, project.ID, mrID); err != nil {
	// 	fmt.Println("Error checking pipeline status:", err)
	// }
	// log.Fatal("stop")

})
