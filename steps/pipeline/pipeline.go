package pipeline

import (
	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/openshift-pipelines/release-tests/pkg/pipelines"
	"github.com/openshift-pipelines/release-tests/pkg/store"
)

var _ = gauge.Step("Create sample pipeline", func() {
	pipelines.CreateSamplePipeline(store.Clients(), store.Namespace())
})

var _ = gauge.Step("Run sample pipeline", func() {
	pipelines.RunSamplePipeline(store.Clients(), store.Namespace())
})

var _ = gauge.Step("Verify sample pipelinerun is successfull", func() {
	pipelines.ValidatePipelineRunStatus(store.Clients(), store.Namespace())
})

var _ = gauge.Step("Create task", func() {
	pipelines.CreateTask(store.Clients(), store.Namespace())
})

var _ = gauge.Step("Run task using <sa> ServiceAccount", func(serviceAccount string) {
	pipelines.CreateTaskRunWithSA(store.Clients(), store.Namespace(), serviceAccount)
})

var _ = gauge.Step("Verify taskrun has failed", func() {
	pipelines.ValidateTaskRunForFailedStatus(store.Clients(), store.Namespace())
})

//=======================================================//

var _ = gauge.Step("Create pipeline", func() {
	pipelines.CreatePipeline(store.Clients(), store.Namespace())
})

var _ = gauge.Step("Run pipeline using <sa> ServiceAccount", func(sa string) {
	pipelines.CreatePipelineRunWithSA(store.Clients(), store.Namespace(), sa)
})

var _ = gauge.Step("Verify pipelinerun has failed", func() {
	pipelines.ValidatePipelineRunForFailedStatus(store.Clients(), store.Namespace())
})
