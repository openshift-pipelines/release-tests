package triggers

import (
	"strconv"

	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/openshift-pipelines/release-tests/pkg/store"
	"github.com/openshift-pipelines/release-tests/pkg/triggers"
)

var _ = gauge.Step("Expose Event listener <elname>", func(elname string) {
	routeurl := triggers.ExposeEventListner(store.Clients(), elname, store.Namespace())
	store.PutScenarioData("route", routeurl)
	store.PutScenarioData("elname", elname)
})

var _ = gauge.Step("Expose Event listener for TLS <elname>", func(elname string) {
	routeurl := triggers.ExposeEventListnerForTLS(store.Clients(), elname, store.Namespace())
	store.PutScenarioData("route", routeurl)
	store.PutScenarioData("elname", elname)
})

var _ = gauge.Step("Mock post event with empty payload", func() {
	gauge.GetScenarioStore()["response"] = triggers.MockPostEventWithEmptyPayload(store.GetScenarioData("route"))
})

var _ = gauge.Step("Assert eventlistener response", func() {
	triggers.AssertElResponse(store.Clients(), store.HttpResponse(), store.GetScenarioData("elname"), store.Namespace())
})

var _ = gauge.Step("Cleanup Triggers", func() {
	triggers.CleanupTriggers(store.Clients(), store.GetScenarioData("elname"), store.Namespace())
})

var _ = gauge.Step("Mock post event to <interceptor> interceptor with event-type <eventType>, payload <payload>, with TLS <tls>", func(interceptor, eventType, payload, tls string) {
	isTLS, _ := strconv.ParseBool(tls)
	gauge.GetScenarioStore()["response"] = triggers.MockPostEvent(store.GetScenarioData("route"), interceptor, eventType, payload, isTLS)
})
