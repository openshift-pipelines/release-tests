package pipeline

import (
	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/openshift-pipelines/release-tests/pkg/pipelines"
	"github.com/openshift-pipelines/release-tests/steps"
)

var _ = gauge.Step("Create sample pipeline", func() {
	pipelines.CreateSamplePipeline(steps.GetClient(), steps.GetNameSpace())
})

var _ = gauge.Step("Run pipeline", func() {
	pipelines.RunSamplePipeline(steps.GetClient(), steps.GetNameSpace())
})

var _ = gauge.Step("Validate pipelinerun for success status", func() {
	pipelines.ValidatePipelineRunStatus(steps.GetClient(), steps.GetNameSpace())
})
