package pipeline

import (
	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/openshift-pipelines/release-tests/pkg/pipelines"
	"github.com/openshift-pipelines/release-tests/steps"
)

var _ = gauge.Step("Create sample pipeline", func() {
	pipelines.CreateSamplePipeline(steps.GetClient(), steps.GetNameSpace())
})

var _ = gauge.Step("Run sample pipeline", func() {
	pipelines.RunSamplePipeline(steps.GetClient(), steps.GetNameSpace())
})

var _ = gauge.Step("Verify sample pipelinerun is successfull", func() {
	pipelines.ValidatePipelineRunStatus(steps.GetClient(), steps.GetNameSpace())
})

var _ = gauge.Step("Create task", func() {
	pipelines.CreateTask(steps.GetClient(), steps.GetNameSpace())
})

var _ = gauge.Step("Run task using <serviceAccount> SA", func(serviceAccount string) {
	pipelines.CreateTaskRunWithSA(steps.GetClient(), steps.GetNameSpace(), serviceAccount)
})

var _ = gauge.Step("Verify taskrun has failed", func() {
	pipelines.ValidateTaskRunForFailedStatus(steps.GetClient(), steps.GetNameSpace())
})

//=======================================================//

var _ = gauge.Step("Create pipeline", func() {
	pipelines.CreatePipeline(steps.GetClient(), steps.GetNameSpace())
})

var _ = gauge.Step("Run pipeline using <sa> SA", func(serviceAccount string) {
	pipelines.CreatePipelineRunWithSA(steps.GetClient(), steps.GetNameSpace(), serviceAccount)
})

var _ = gauge.Step("Verify pipelinerun has failed", func() {
	pipelines.ValidatePipelineRunForFailedStatus(steps.GetClient(), steps.GetNameSpace())
})
