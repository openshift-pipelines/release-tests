package chains

import (
	"os"

	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/operator"
)

var _ = gauge.Step("Verify <resourceType> signature", func(resourceType string) {
	operator.VerifySignature(resourceType)
})

var _ = gauge.Step("Start the kaniko-chains task", func() {
	operator.StartKanikoTask()
})

var _ = gauge.Step("Verify image signature", func() {
	operator.VerifyImageSignature()
})

var _ = gauge.Step("Check attestation exists", func() {
	operator.CheckAttestationExists()
})

var _ = gauge.Step("Verify attestation", func() {
	operator.VerifyAttestation()
})

var _ = gauge.Step("Verify that image registry variable is exported", func() {
	if os.Getenv("CHAINS_REPOSITORY") == "" {
		testsuit.T.Errorf("'CHAINS_REPOSITORY' environment variable is not exported")
	}
})
