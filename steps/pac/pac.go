package pac

import (
	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/openshift-pipelines/release-tests/pkg/k8s"
	"github.com/openshift-pipelines/release-tests/pkg/pac"
	"github.com/openshift-pipelines/release-tests/pkg/pipelines"
	"github.com/openshift-pipelines/release-tests/pkg/store"
)

var _ = gauge.Step("Setup Gitlab Client", func() {
	c := pac.InitGitLabClient()
	pac.SetGitLabClient(c)
})

var _ = gauge.Step("Create Smee deployment", func() {
	pac.SetupSmeeDeployment()
	k8s.ValidateDeployments(store.Clients(), store.Namespace(), store.GetScenarioData("smeeDeploymentName"))
})

var _ = gauge.Step("Configure GitLab repo for <eventType> in <branch>", func(eventType, branch string) {
	pac.SetupGitLabProject()
	pac.GeneratePipelineRunYaml(eventType, branch)
})

var _ = gauge.Step("Configure PipelineRun", func() {
	pac.ConfigurePreviewChanges()
})

var _ = gauge.Step("Validate PipelineRun for <state>", func(state string) {
	pipelineName := pac.GetPipelineNameFromMR()
	pipelines.ValidatePipelineRun(store.Clients(), pipelineName, state, "no", store.Namespace())
})

var _ = gauge.Step("Update Annotation <annotationKey> with <annotationValue>", func(annotationKey, annotationValue string) {
	pac.UpdateAnnotation(annotationKey, annotationValue)
})

var _ = gauge.Step("Add Comment <comment> in MR", func(comment string) {
	pac.AddComment(comment)
})

var _ = gauge.Step("Add Label Name <labelName> with <color> color with description <description>", func(labelName, color, description string) {
	pac.AddLabel(labelName, color, description)
})

var _ = gauge.Step("Cleanup PAC", func() {
	pac.CleanupPAC(store.Clients(), store.GetScenarioData("smeeDeploymentName"), store.Namespace())
})
