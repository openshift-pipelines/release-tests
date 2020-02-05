package olm

import (
	"sync"

	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	"github.com/openshift-pipelines/release-tests/pkg/operator"
	"github.com/openshift-pipelines/release-tests/steps"
)

var once sync.Once
var _ = gauge.Step("Operator should be installed", func() {
	// TODO: why only once ?
	once.Do(func() {
		operator.ValidateInstall(steps.OperatorClient())
	})
})

var _ = gauge.Step("Wait for Cluster CR availability", func() {
	operator.WaitForClusterCR(steps.OperatorClient(), config.ClusterCRName)
})

var _ = gauge.Step("Validate SCC", func() {
	operator.ValidateSCC(steps.OperatorClient())
})

var _ = gauge.Step("Validate pipelines deployment into target namespace (openshift-pipelines)", func() {
	operator.ValidatePipelineDeployments(steps.OperatorClient())
})

var _ = gauge.Step("Validate pipeline version <version>", func(version string) {
	operator.VerifyPipelineVersion(steps.OperatorClient(), version)
})

var _ = gauge.Step("Validate Triggers deployment into target namespace (openshift-pipelines)", func() {
	operator.ValidateTriggerDeployments(steps.OperatorClient())
})

var _ = gauge.Step("Validate opeartor setup status", func() {
	operator.ValidateInstalledStatus(steps.OperatorClient())
})
