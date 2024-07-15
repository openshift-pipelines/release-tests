package chains

import (
	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/openshift-pipelines/release-tests/pkg/operator"
)

var _ = gauge.Step("Start the manual-approval-pipeline pipeline", func() {
	operator.StartApprovalGatePipeline()
})

var _ = gauge.Step("Approve the manual-approval-pipeline", func() {
	taskname := operator.GetApprovaltasklist()
	operator.ApproveApprovalGatePipeline(taskname)
})

var _ = gauge.Step("Reject the manual-approval-pipeline", func() {
	taskname := operator.GetApprovaltasklist()
	operator.RejectApprovalGatePipeline(taskname)
})
