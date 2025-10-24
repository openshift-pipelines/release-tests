package pac

import (
	"fmt"

	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/getgauge-contrib/gauge-go/testsuit"
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
	pac.SetupGitLabProject()
})

var _ = gauge.Step("Configure GitLab repo for <eventType> in <branch>", func(eventType, branch string) {
	pac.GeneratePipelineRunYaml(eventType, branch)
})

var _ = gauge.Step("Configure PipelineRun", func() {
	pac.ConfigurePreviewChanges()
})

var _ = gauge.Step("Trigger push event on main branch", func() {
	pac.TriggerPushOnForkMain()
})

var _ = gauge.Step("Validate PipelineRun for <state>", func(state string) {
	pipelineName := pac.GetPipelineNameFromMR()
	pipelines.ValidatePipelineRun(store.Clients(), pipelineName, state, store.Namespace())
})

var _ = gauge.Step("Validate <event_type> PipelineRun for <state>", func(event_type, state string) {
	switch event_type {
	case "pull_request":
		pipelineName := pac.GetPipelineNameFromMR()
		pipelines.ValidatePipelineRun(store.Clients(), pipelineName, state, store.Namespace())
	case "push":
		pipelineName := pac.GetPushPipelineNameFromMain()
		pipelines.ValidatePipelineRun(store.Clients(), pipelineName, state, store.Namespace())
	default:
		testsuit.T.Fail(fmt.Errorf("invalid event type: %s", event_type))
	}
})

var _ = gauge.Step("Validate PAC Info Install", func() {
	pac.AssertPACInfoInstall()
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
