package olm

import (
	"sync"

	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	"github.com/openshift-pipelines/release-tests/pkg/operator"
	"github.com/openshift-pipelines/release-tests/pkg/store"
)

var once sync.Once
var _ = gauge.Step("Operator should be installed", func() {
	// TODO (praveen): why only once ?
	once.Do(func() {
		operator.ValidateInstall(store.Clients())
	})
})

var _ = gauge.Step("Wait for Cluster CR availability", func() {
	operator.WaitForClusterCR(store.Clients(), config.ClusterCRName)
})

var _ = gauge.Step("Validate SCC", func() {
	operator.ValidateSCC(store.Clients())
})

var _ = gauge.Step("Validate pipelines deployment into target namespace (openshift-pipelines)", func() {
	operator.ValidatePipelineDeployments(store.Clients())
})

var _ = gauge.Step("Validate pipeline version <version>", func(version string) {
	operator.VerifyPipelineVersion(store.Clients(), version)
})

var _ = gauge.Step("Validate Triggers deployment into target namespace (openshift-pipelines)", func() {
	operator.ValidateTriggerDeployments(store.Clients())
})

var _ = gauge.Step("Validate opeartor setup status", func() {
	operator.ValidateInstalledStatus(store.Clients())
})
