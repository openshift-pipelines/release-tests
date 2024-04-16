package openshift

import (
	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/openshift-pipelines/release-tests/pkg/openshift"
	"github.com/openshift-pipelines/release-tests/pkg/store"
	"github.com/openshift-pipelines/release-tests/pkg/triggers"
)

var _ = gauge.Step("Get tags of the imagestream <imageStream> from namespace <namespace> and store to variable <variableName>", func(imageStream, namespace, variableName string) {
	tagNames := openshift.GetImageStreamTags(store.Clients(), namespace, imageStream)
	store.PutScenarioDataSlice(variableName, tagNames)
})

var _ = gauge.Step("Verify that image stream <is> exists", func(is string) {
	openshift.VerifyImageStreamExists(store.Clients(), is, "openshift")
})

var _ = gauge.Step("Get route url of the route <routeName>", func(routeName string) {
	routeurl := triggers.GetRouteURL(routeName, store.Namespace())
	store.PutScenarioData("routeurl", routeurl)
})
