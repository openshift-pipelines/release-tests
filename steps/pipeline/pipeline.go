package pipeline

import (
	"fmt"
	"log"
	"strings"

	"github.com/getgauge-contrib/gauge-go/gauge"
	m "github.com/getgauge-contrib/gauge-go/models"
	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/cmd"
	"github.com/openshift-pipelines/release-tests/pkg/pipelines"
	"github.com/openshift-pipelines/release-tests/pkg/store"
)

var _ = gauge.Step("Verify taskrun <table>", func(table *m.Table) {
	for _, row := range table.Rows {
		trname := row.Cells[1]
		status := row.Cells[2]
		pipelines.ValidateTaskRun(store.Clients(), trname, status, store.Namespace())
	}
})

var _ = gauge.Step("Verify pipelinerun <table>", func(table *m.Table) {
	for _, row := range table.Rows {
		prname := row.Cells[1]
		status := row.Cells[2]
		pipelines.ValidatePipelineRun(store.Clients(), prname, status, store.Namespace())
	}
})

var _ = gauge.Step("Watch for pipelinerun resources", func() {
	pipelines.WatchForPipelineRun(store.Clients(), store.Namespace())
})

var _ = gauge.Step("Verify taskrun <trname> label propagation", func(trname string) {
	pipelines.ValidateTaskRunLabelPropogation(store.Clients(), trname, store.Namespace())
})

var _ = gauge.Step("Assert no new pipelineruns created", func() {
	pipelines.AssertForNoNewPipelineRunCreation(store.Clients(), store.Namespace())
})

var _ = gauge.Step("<numberOfPr> pipelinerun(s) should be present within <timeoutSeconds> seconds", func(numberOfPr, timeoutSeconds string) {
	pipelines.AssertNumberOfPipelineruns(store.Clients(), store.Namespace(), numberOfPr, timeoutSeconds)
})

var _ = gauge.Step("<numberOfPr> pipelinerun(s) with status <status> should be present within <timeoutSeconds> seconds", func(numberOfPr, status, timeoutSeconds string) {
	pipelines.AssertNumberOfPipelinerunsWithStatus(store.Clients(), store.Namespace(), numberOfPr, status, timeoutSeconds)
})

var _ = gauge.Step("<numberOfTr> taskrun(s) should be present within <timeoutSeconds> seconds", func(numberOfTr, timeoutSeconds string) {
	pipelines.AssertNumberOfTaskruns(store.Clients(), store.Namespace(), numberOfTr, timeoutSeconds)
})

var _ = gauge.Step("Tasks <ts> are <status> in namespace <namespace>", func(ts, status string, namespace string) {
	log.Printf("Checking if tasks %v is/are %v in namespace %v", ts, status, namespace)
	tsList := strings.Split(ts, ",")
	if status == "present" {
		for _, c := range tsList {
			pipelines.AssertTaskPresent(store.Clients(), namespace, c)
		}
	} else {
		for _, c := range tsList {
			pipelines.AssertTaskNotPresent(store.Clients(), namespace, c)
		}
	}
})

var _ = gauge.Step("StepActions <stepActions> are <status> in namespace <namespace>", func(stepActions, status string, namespace string) {
	log.Printf("Checking if stepactions %v is/are %v in namespace %v", stepActions, status, namespace)
	saList := strings.Split(stepActions, ",")
	if status == "present" {
		for _, c := range saList {
			pipelines.AssertStepActionPresent(store.Clients(), namespace, c)
		}
	} else {
		for _, c := range saList {
			pipelines.AssertStepActionNotPresent(store.Clients(), namespace, c)
		}
	}
})

var _ = gauge.Step("Assert pipelines are <status> in <namespace> namespace", func(status, namespace string) {
	if status == "present" {
		pipelines.AssertPipelinesPresent(store.Clients(), namespace)
	} else {
		pipelines.AssertPipelinesNotPresent(store.Clients(), namespace)
	}
})

var _ = gauge.Step("Verify the latest pipelinerun for <state> state", func(state string) {
	namespace := store.Namespace()
	prname, err := pipelines.GetLatestPipelinerun(store.Clients(), namespace)
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("failed to get pipelinerun from %s: %v", namespace, err))
	}
	pipelines.ValidatePipelineRun(store.Clients(), prname, state, namespace)
})

var _ = gauge.Step("Validate pipelinerun stored in variable <prname> with task <taskname> logs contains <expectedLogs>", func(prname, taskname, expectedLogs string) {
	logs := cmd.MustSucceed("oc", "logs", "-l", "tekton.dev/pipelineRun="+store.GetScenarioData(prname)+",tekton.dev/pipelineTask="+taskname, "-n", store.Namespace()).Stdout()
	logsLower := strings.ToLower(logs)
	log.Printf("Logs output: %s\n", logsLower)
	expectedLogsLower := strings.ToLower(expectedLogs)
	if !strings.Contains(logsLower, expectedLogsLower) {
		testsuit.T.Errorf("Logs validation failed: Logs did not contain expected content")
	} else {
		log.Print("Logs validated successfully")
	}
})
