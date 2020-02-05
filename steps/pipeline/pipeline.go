package pipeline

import (
	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/openshift-pipelines/release-tests/pkg/pipelines"
	"github.com/openshift-pipelines/release-tests/steps"
)

var _ = gauge.Step("Create sample pipeline", func() {
	pipelines.CreateSamplePipeline(steps.Clients(), steps.Namespace())
})

var _ = gauge.Step("Run sample pipeline", func() {
	pipelines.RunSamplePipeline(steps.Clients(), steps.Namespace())
})

var _ = gauge.Step("Verify sample pipelinerun is successfull", func() {
	pipelines.ValidatePipelineRunStatus(steps.Clients(), steps.Namespace())
})

var _ = gauge.Step("Create task", func() {
	pipelines.CreateTask(steps.Clients(), steps.Namespace())
})

var _ = gauge.Step("Run task using <sa> ServiceAccount", func(serviceAccount string) {
	pipelines.CreateTaskRunWithSA(steps.Clients(), steps.Namespace(), serviceAccount)
})

var _ = gauge.Step("Verify taskrun has failed", func() {
	pipelines.ValidateTaskRunForFailedStatus(steps.Clients(), steps.Namespace())
})

//=======================================================//

var _ = gauge.Step("Create pipeline", func() {
	pipelines.CreatePipeline(steps.Clients(), steps.Namespace())
})

var _ = gauge.Step("Run pipeline using <sa> SA", func(serviceAccount string) {
	pipelines.CreatePipelineRunWithSA(steps.Clients(), steps.Namespace(), serviceAccount)
})

var _ = gauge.Step("Verify pipelinerun has failed", func() {
	pipelines.ValidatePipelineRunForFailedStatus(steps.Clients(), steps.Namespace())
})
