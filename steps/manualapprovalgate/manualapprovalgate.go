package approvalgate

import (
	"errors"
	"strings"

	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/getgauge-contrib/gauge-go/testsuit"
	approvalgate "github.com/openshift-pipelines/release-tests/pkg/manualapprovalgate"
	"github.com/openshift-pipelines/release-tests/pkg/store"
	"github.com/openshift-pipelines/release-tests/pkg/tkn"
)

var _ = gauge.Step("Start the <pipelineName> pipeline with workspace <workspaceValue>", func(pipelineName, workspaceValue string) {
	params := make(map[string]string)
	workspaces := make(map[string]string)
	workspaces[strings.Split(workspaceValue, ",")[0]] = strings.Split(workspaceValue, ",")[1]
	tkn.StartPipeline(pipelineName, params, workspaces, store.Namespace(), "--use-param-defaults")
})

var _ = gauge.Step("Approve the manual-approval-pipeline", func() {
	tasks, err := approvalgate.ListApprovalTask(store.Clients())
	if err != nil {
		testsuit.T.Errorf("Error while listing approval gate tasks: %v", err)
		return
	}

	for _, task := range tasks {
		approvalgate.ApproveApprovalGatePipeline(task.Name)
	}
})

var _ = gauge.Step("Reject the manual-approval-pipeline", func() {
	tasks, err := approvalgate.ListApprovalTask(store.Clients())
	if err != nil {
		testsuit.T.Errorf("Error while listing approval gate tasks: %v", err)
		return
	}

	for _, task := range tasks {
		approvalgate.RejectApprovalGatePipeline(task.Name)
	}
})

var _ = gauge.Step("Validate the manual-approval-pipeline for <status> state", func(status string) {
	success, err := approvalgate.ValidateApprovalGatePipeline(status)
	if err != nil {
		testsuit.T.Fail(err)
		return
	}

	if !success {
		testsuit.T.Fail(errors.New("validation failed: no approvaltasks matched the expected status"))
	}
})
