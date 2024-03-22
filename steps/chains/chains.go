package chains

import (
	"os"

	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/oc"
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

var _ = gauge.Step("Create secret with image registry credentials for SA", func() {
	if os.Getenv("CHAINS_DOCKER_CONFIG_JSON") == "" {
		testsuit.T.Errorf("'CHAINS_DOCKER_CONFIG_JSON' robot credentials environment variable is not exported")
	} else {
		dockerConfig := os.Getenv("CHAINS_DOCKER_CONFIG_JSON")
		oc.CreateChainsImageRegistrySecret(dockerConfig)
	}
})

var _ = gauge.Step("Update the TektonConfig with taskrun format as <format> taskrun storage as <r_storage> oci storage as <oci_storage> transparency mode as <mode>", func(format, runStorage, ociStorage, mode string) {
	patch_data := "{\"spec\":{\"chain\":{\"artifacts.taskrun.format\":\"" + format + "\",\"artifacts.taskrun.storage\":\"" + runStorage + "\",\"artifacts.oci.storage\":\"" + ociStorage + "\",\"transparency.enabled\":\"" + mode + "\"}}}"
	oc.UpdateTektonConfig(patch_data)
})
