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

const (
	pacProviderKey        = "pac.provider"
	pacLastPipelineRunKey = "pac.lastPipelinerun"
)

var _ = gauge.Step("Setup Gitlab Client", func() {
	c := pac.InitGitLabClient()
	pac.SetGitLabClient(c)
	store.PutScenarioData(pacProviderKey, "gitlab")
})

var _ = gauge.Step("Setup Github Client", func() {
	c := pac.InitGitHubClient()
	pac.SetGitHubClient(c)
	store.PutScenarioData(pacProviderKey, "github")
})

var _ = gauge.Step("Create Smee deployment", func() {
	pac.SetupSmeeDeployment()
	k8s.ValidateDeployments(store.Clients(), store.Namespace(), store.GetScenarioData("smeeDeploymentName"))
	switch store.GetScenarioData(pacProviderKey) {
	case "gitlab":
		pac.SetupGitLabProject()
	case "github":
		pac.SetupGitHubProject()
	default:
		testsuit.T.Fail(fmt.Errorf("unknown pac provider %q", store.GetScenarioData(pacProviderKey)))
	}
})

var _ = gauge.Step("Configure GitLab repo for <eventType> in <branch>", func(eventType, branch string) {
	pac.GeneratePipelineRunYaml(eventType, branch)
})

var _ = gauge.Step("Configure GitHub repo for <eventType> in <branch>", func(eventType, branch string) {
	pac.GeneratePipelineRunYaml(eventType, branch)
})

var _ = gauge.Step("Configure PipelineRun", func() {
	switch store.GetScenarioData(pacProviderKey) {
	case "gitlab":
		pac.ConfigurePreviewChanges()
	case "github":
		pac.ConfigurePreviewChangesGitHub()
	default:
		testsuit.T.Fail(fmt.Errorf("unknown pac provider %q", store.GetScenarioData(pacProviderKey)))
	}
})

var _ = gauge.Step("Trigger push event on main branch", func() {
	switch store.GetScenarioData(pacProviderKey) {
	case "gitlab":
		pac.TriggerPushOnForkMain()
	case "github":
		pac.TriggerPushOnGitHubMain()
	default:
		testsuit.T.Fail(fmt.Errorf("unknown pac provider %q", store.GetScenarioData(pacProviderKey)))
	}
})

var _ = gauge.Step("Validate PipelineRun for <state>", func(state string) {
	switch store.GetScenarioData(pacProviderKey) {
	case "gitlab":
		pipelineName := pac.GetPipelineNameFromMR()
		pipelines.ValidatePipelineRun(store.Clients(), pipelineName, state, store.Namespace())
	case "github":
		pipelineName := pac.WaitForNewPipelineRunName("")
		pipelines.ValidatePipelineRun(store.Clients(), pipelineName, state, store.Namespace())
		gauge.GetScenarioStore()[pacLastPipelineRunKey] = pipelineName
	default:
		testsuit.T.Fail(fmt.Errorf("unknown pac provider %q", store.GetScenarioData(pacProviderKey)))
	}
})

var _ = gauge.Step("Validate <event_type> PipelineRun for <state>", func(event_type, state string) {
	last := ""
	if v, ok := gauge.GetScenarioStore()[pacLastPipelineRunKey].(string); ok {
		last = v
	}

	switch event_type {
	case "pull_request":
		switch store.GetScenarioData(pacProviderKey) {
		case "gitlab":
			pipelineName := pac.GetPipelineNameFromMR()
			pipelines.ValidatePipelineRun(store.Clients(), pipelineName, state, store.Namespace())
			gauge.GetScenarioStore()[pacLastPipelineRunKey] = pipelineName
		case "github":
			pipelineName := pac.WaitForNewPipelineRunName(last)
			pipelines.ValidatePipelineRun(store.Clients(), pipelineName, state, store.Namespace())
			gauge.GetScenarioStore()[pacLastPipelineRunKey] = pipelineName
		default:
			testsuit.T.Fail(fmt.Errorf("unknown pac provider %q", store.GetScenarioData(pacProviderKey)))
		}
	case "push":
		var pipelineName string
		switch store.GetScenarioData(pacProviderKey) {
		case "gitlab":
			pipelineName = pac.GetPushPipelineNameFromMain()
		case "github":
			pipelineName = pac.WaitForNewPipelineRunName(last)
		default:
			testsuit.T.Fail(fmt.Errorf("unknown pac provider %q", store.GetScenarioData(pacProviderKey)))
		}
		pipelines.ValidatePipelineRun(store.Clients(), pipelineName, state, store.Namespace())
		gauge.GetScenarioStore()[pacLastPipelineRunKey] = pipelineName
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
	switch store.GetScenarioData(pacProviderKey) {
	case "gitlab":
		pac.CleanupPAC(store.Clients(), store.GetScenarioData("smeeDeploymentName"), store.Namespace())
	case "github":
		pac.CleanupPACGitHub(store.Clients(), store.GetScenarioData("smeeDeploymentName"), store.Namespace())
	default:
		testsuit.T.Fail(fmt.Errorf("unknown pac provider %q", store.GetScenarioData(pacProviderKey)))
	}
})
