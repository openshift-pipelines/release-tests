package pipeline

import (
	"fmt"
	"log"
	"strings"

	"github.com/getgauge-contrib/gauge-go/gauge"
	m "github.com/getgauge-contrib/gauge-go/models"
	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/config"
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
		labelCheck := row.Cells[3]
		pipelines.ValidatePipelineRun(store.Clients(), prname, status, labelCheck, store.Namespace())
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

var _ = gauge.Step("<numberOfTr> taskrun(s) should be present within <timeoutSeconds> seconds", func(numberOfTr, timeoutSeconds string) {
	pipelines.AssertNumberOfTaskruns(store.Clients(), store.Namespace(), numberOfTr, timeoutSeconds)
})

var _ = gauge.Step("<cts> clustertasks are <status>", func(cts, status string) {
	if cts == "community" {
		cts = config.CommunityClustertasks
	}
	log.Printf("Checking if clustertasks %v is/are %v", cts, status)
	ctsList := strings.Split(cts, ",")
	if status == "present" {
		for _, c := range ctsList {
			pipelines.AssertClustertaskPresent(store.Clients(), c)
		}
	} else {
		for _, c := range ctsList {
			pipelines.AssertClustertaskNotPresent(store.Clients(), c)
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
		testsuit.T.Fail(fmt.Errorf("failed to get pipelinerun from %s", namespace))
	}
	pipelines.ValidatePipelineRun(store.Clients(), prname, state, "no", namespace)
})
