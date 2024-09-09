package pac

import (
	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/openshift-pipelines/release-tests/pkg/k8s"
	"github.com/openshift-pipelines/release-tests/pkg/pac"
	"github.com/openshift-pipelines/release-tests/pkg/pipelines"
	"github.com/openshift-pipelines/release-tests/pkg/store"
)

var _ = gauge.Step("Configure Gitlab token for PAC tests", func() {

	pac.ConfigureGitlabToken()
})

var _ = gauge.Step("Create Smee Deployment with <elname>", func(elname string) {

	pac.SmeeDeployment(elname)
	k8s.ValidateDeployments(store.Clients(), store.Namespace(), store.GetScenarioData("smee_deployment_name"))

})

var _ = gauge.Step("Configure & Validate Gitlab repo for pipelinerun", func() {

	client := pac.InitGitLabClient()
	project := pac.SetupGitLabProject(client)
	pac.ConfigurePreviewChanges(client, project.ID)
	pipelines.ValidatePipelineRun(store.Clients(), "gitlab-run", "successful", "no", store.Namespace())
})

var _ = gauge.Step("Cleanup PAC", func() {
	pac.CleanupPAC(store.Clients(), store.GetScenarioData("elname"), store.GetScenarioData("smee_deployment_name"), store.Namespace())
})
