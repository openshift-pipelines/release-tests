package triggers

import (
	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/openshift-pipelines/release-tests/pkg/store"
	"github.com/openshift-pipelines/release-tests/pkg/triggers"
)

var _ = gauge.Step("Expose Event listener <elname>", func(elname string) {
	routeurl := triggers.ExposeEventListner(store.Clients(), elname, store.Namespace())
	store.PutScenarioData("route", routeurl)
})

var _ = gauge.Step("Mock get event", func() {
	gauge.GetScenarioStore()["response"] = triggers.MockGetEvent(store.GetScenarioData("route"))
})

var _ = gauge.Step("Mock push event <payload>", func(payload string) {
	gauge.GetScenarioStore()["response"] = triggers.MockPushEvent(store.GetScenarioData("route"), payload)
})

var _ = gauge.Step("Mock push event <payload> to gitlab interceptor", func(payload string) {
	gauge.GetScenarioStore()["response"] = triggers.MockPushEventToGitlabInterceptor(store.GetScenarioData("route"), payload)
})

var _ = gauge.Step("Assert eventlistener <elname> response", func(elname string) {
	triggers.AssertElResponse(store.HttpResponse(), elname, store.Namespace())
})
