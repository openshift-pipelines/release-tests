package triggers

import (
	"github.com/getgauge-contrib/gauge-go/gauge"
	m "github.com/getgauge-contrib/gauge-go/models"
	"github.com/openshift-pipelines/release-tests/pkg/store"
	"github.com/openshift-pipelines/release-tests/pkg/triggers"
)

var _ = gauge.Step("Expose Event listener <elname>", func(elname string) {
	routeurl := triggers.ExposeEventListner(store.Clients(), elname, store.Namespace())
	store.PutScenarioData("route", routeurl)
	store.PutScenarioData("elname", elname)
})

var _ = gauge.Step("Mock get event", func() {
	gauge.GetScenarioStore()["response"] = triggers.MockGetEvent(store.GetScenarioData("route"))
})

var _ = gauge.Step("Assert eventlistener response", func() {
	triggers.AssertElResponse(store.Clients(), store.HttpResponse(), store.GetScenarioData("elname"), store.Namespace())
})

var _ = gauge.Step("Cleanup Triggers", func() {
	triggers.CleanupTriggers(store.Clients(), store.GetScenarioData("elname"), store.Namespace())
})

var _ = gauge.Step("Mock post event <eventpayloads>", func(eventpayloads *m.Table) {
	for _, row := range eventpayloads.Rows {
		payload := row.Cells[1]
		headers := row.Cells[2]
		gauge.GetScenarioStore()["response"] = triggers.MockPostEvent(store.GetScenarioData("route"), headers, payload)
	}
})
