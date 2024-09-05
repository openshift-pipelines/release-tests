package pac

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	"github.com/openshift-pipelines/release-tests/pkg/k8s"
	"github.com/openshift-pipelines/release-tests/pkg/pac"
	"github.com/openshift-pipelines/release-tests/pkg/pipelines"
	"github.com/openshift-pipelines/release-tests/pkg/store"
	"github.com/openshift-pipelines/release-tests/pkg/triggers"
)

var _ = gauge.Step("Create Smee Deployment with <elname>", func(elname string) {

	var err error
	smeeURL, err := pac.GetNewSmeeURL()
	if err != nil {
		log.Fatalf("Failed to get a new Smee URL: %v", err)
	}
	store.PutScenarioData("SMEE_URL", smeeURL)

	routeurl := triggers.GetRoute(elname, store.Namespace())
	store.PutScenarioData("route", routeurl)

	err = pac.CreateSmeeDeployment(store.Clients(), store.Namespace(), smeeURL, routeurl)
	if err != nil {
		log.Fatalf("Failed to create deployment: %v", err)
	}
	k8s.ValidateDeployments(store.Clients(), store.Namespace(), "gosmee-client")

})

var _ = gauge.Step("Configure Gitlab repo", func() {

	smeeURL := store.GetScenarioData("SMEE_URL")
	projectIDOrPath := ">Project ID<"
	targetGroupNamespace := "<Target Namespace>"
	privateToken := "<PAC Token>"

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
	// defer func() {
	// 	if err := pac.DeleteGitlabProject(client, project.ID); err != nil {
	// 		log.Printf("Cleanup failed: %v", err)
	// 	}
	// }()

	token := config.TriggersSecretToken
	err = pac.AddWebhook(client, project.ID, smeeURL, token)
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
		log.Fatal(mrID, err)
	}

	pipelines.ValidatePipelineRun(store.Clients(), "gitlab-run", "successful", "no", store.Namespace())
})
