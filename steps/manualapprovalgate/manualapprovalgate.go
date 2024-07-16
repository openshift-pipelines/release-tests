package operator

import (
	"errors"

	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/operator"
)

var _ = gauge.Step("Start the manual-approval-pipeline pipeline", func() {
	operator.StartApprovalGatePipeline()
})

var _ = gauge.Step("Approve the manual-approval-pipeline", func() {
	tasks := operator.GetApprovalTaskList()
	if tasks == nil {
		testsuit.T.Errorf("No Approval Gate Tasks Found")
	}

	for _, task := range tasks {
		operator.ApproveApprovalGatePipeline(task.Name)
	}
})

var _ = gauge.Step("Reject the manual-approval-pipeline", func() {
	tasks := operator.GetApprovalTaskList()
	if tasks == nil {
		testsuit.T.Errorf("No Approval Gate Tasks Found")
	}
	for _, task := range tasks {
		operator.RejectApprovalGatePipeline(task.Name)
	}
})

var _ = gauge.Step("Validate the manual-approval-pipeline for <status> state", func(status string) {
	success, err := operator.ValidateApprovalGatePipeline(status)
	if err != nil {
		testsuit.T.Fail(err)
		return
	}

	if !success {
		testsuit.T.Fail(errors.New("validation failed: no approvaltasks match the expected status"))
	}
})
