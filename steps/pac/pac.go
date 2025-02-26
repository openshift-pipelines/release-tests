package pac

import (
	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/openshift-pipelines/release-tests/pkg/k8s"
	"github.com/openshift-pipelines/release-tests/pkg/pac"
	"github.com/openshift-pipelines/release-tests/pkg/pipelines"
	"github.com/openshift-pipelines/release-tests/pkg/store"
)

var _ = gauge.Step("Configure GitLab token for PAC tests", func() {
	pac.CreateGitLabSecret()
})

var _ = gauge.Step("Create Smee deployment", func() {
	pac.SetupSmeeDeployment()
	k8s.ValidateDeployments(store.Clients(), store.Namespace(), store.GetScenarioData("smee_deployment_name"))
})

var _ = gauge.Step("Configure GitLab repo and validate pipelinerun", func() {
	client := pac.InitGitLabClient()
	project := pac.SetupGitLabProject(client)
	pipelineName := pac.ConfigurePreviewChanges(client, project.ID)
	pipelines.ValidatePipelineRun(store.Clients(), pipelineName, "successful", "no", store.Namespace())
	pac.AssertPACInfoInstall()
	pac.CleanupPAC(client, store.Clients(), project.ID, store.GetScenarioData("smee_deployment_name"), store.Namespace())
})
