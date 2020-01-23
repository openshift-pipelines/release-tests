package pipeline

import (
	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/openshift-pipelines/release-tests/pkg/pipelines"
	"github.com/openshift-pipelines/release-tests/steps"
)

var _ = gauge.Step("Create Task", func() {
	pipelines.CreateTask(steps.GetClient(), steps.GetNameSpace())
})

var _ = gauge.Step("Run Task with <serviceAccount> SA", func(serviceAccount string) {
	pipelines.CreateTaskRunWithSA(steps.GetClient(), steps.GetNameSpace(), serviceAccount)
})

var _ = gauge.Step("Validate TaskRun for failed status", func() {
	pipelines.ValidateTaskRunForFailedStatus(steps.GetClient(), steps.GetNameSpace())
})

//=======================================================//

var _ = gauge.Step("Create pipeline", func() {
	pipelines.CreatePipeline(steps.GetClient(), steps.GetNameSpace())
})

var _ = gauge.Step("Run pipeline with <serviceAccount> SA", func(serviceAccount string) {
	pipelines.CreatePipelineRunWithSA(steps.GetClient(), steps.GetNameSpace(), serviceAccount)
})

var _ = gauge.Step("Validate pipelineRun for failed status", func() {
	pipelines.ValidatePipelineRunForFailedStatus(steps.GetClient(), steps.GetNameSpace())
})
