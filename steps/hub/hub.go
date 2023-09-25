package TektonHub

import (
	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/openshift-pipelines/release-tests/pkg/store"
	"github.com/openshift-pipelines/release-tests/pkg/tektonhub"
)

var _ = gauge.Step("Create TektinHub CR", func() {
	tektonhub.CreateHubCR(store.Clients())
})

var _ = gauge.Step("Verify that the TektonHub elements like Kind, Platform, Catalog, Category are available", func() {
	tektonhub.GetTektonHubElements()
})
