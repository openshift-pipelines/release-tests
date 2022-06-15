package metrics

import (
	"github.com/getgauge-contrib/gauge-go/gauge"
	m "github.com/getgauge-contrib/gauge-go/models"
	"github.com/openshift-pipelines/release-tests/pkg/assert"
	"github.com/openshift-pipelines/release-tests/pkg/monitoring"
	"github.com/openshift-pipelines/release-tests/pkg/store"
)

var _ = gauge.Step("Verify job health status metrics <table>", func(table *m.Table) {
	for _, row := range table.Rows {
		ts := monitoring.TargetService{Job: row.Cells[1], ExpectedValue: row.Cells[2]}
		err := monitoring.VerifyHealthStatusMetric(store.Clients(), ts)
		assert.FailOnError(err)
	}
})

var _ = gauge.Step("Verify pipelines controlPlane metrics", func() {
	err := monitoring.VerifyPipelinesControlPlaneMetrics(store.Clients())
	assert.FailOnError(err)
})