package openshift

import (
	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/openshift-pipelines/release-tests/pkg/openshift"
	"github.com/openshift-pipelines/release-tests/pkg/store"
)

var _ = gauge.Step("Get tags of the imagestream <imageStream> from namespace <namespace> and store to variable <variableName>", func(imageStream, namespace, variableName string) {
	tagNames := openshift.GetImageStreamTags(store.Clients(), namespace, imageStream)
	store.PutScenarioDataSlice(variableName, tagNames)
})