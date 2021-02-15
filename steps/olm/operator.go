package olm

import (
	"fmt"
	"sync"

	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/openshift-pipelines/release-tests/pkg/assert"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	"github.com/openshift-pipelines/release-tests/pkg/olm"
	"github.com/openshift-pipelines/release-tests/pkg/operator"
	"github.com/openshift-pipelines/release-tests/pkg/store"
)

var once sync.Once
var _ = gauge.Step("Validate Operator should be installed", func() {
	// TODO (praveen): why only once ?
	once.Do(func() {
		operator.ValidateInstall(store.Clients())
	})
})

var _ = gauge.Step("Subscribe to operator", func() {
	// Creates subscription yaml with configured details from env/test/test.properties
	_, err := olm.SubscribeAndWaitForOperatorToBeReady(store.Clients(), "openshift-pipelines-operator-rh", config.Flags.Channel, config.Flags.CatalogSource)
	assert.NoError(err, fmt.Sprintf("failed to Subscribe :%s", err))
})

var _ = gauge.Step("Wait for Cluster CR availability", func() {
	operator.WaitForClusterCR(store.Clients(), config.ClusterCRName)
})

var _ = gauge.Step("Upgrade operator subscription", func() {
	// Creates subscription yaml with configured details from env/test/test.properties
	_, err := olm.UptadeSubscriptionAndWaitForOperatorToBeReady(store.Clients(), "openshift-pipelines-operator-rh", config.Flags.Channel)
	assert.NoError(err, fmt.Sprintf("failed to Subscribe :%s", err))
})

var _ = gauge.Step("Validate SCC", func() {
	operator.ValidateSCC(store.Clients())
})

var _ = gauge.Step("Validate pipelines deployment", func() {
	operator.ValidatePipelineDeployments(store.Clients())
})

var _ = gauge.Step("Validate pipeline version <version>", func(version string) {
	operator.VerifyPipelineVersion(store.Clients(), version)
})

var _ = gauge.Step("Validate triggers deployment", func() {
	operator.ValidateTriggerDeployments(store.Clients())
})

var _ = gauge.Step("Uninstall Operator", func() {
	//cleanup operator Traces
	operator.Uninstall(store.Clients())
})
