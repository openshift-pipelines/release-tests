package ecosystem

import (
	"os"

	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/oc"
)

var _ = gauge.Step("Verify that jib-maven image registry variable is exported", func() {
	if os.Getenv("JIB_MAVEN_REPOSITORY") == "" {
		testsuit.T.Errorf("'JIB_MAVEN_REPOSITORY' environment variable is not exported")
	}
})

var _ = gauge.Step("Create secret with image registry credentials for jib-maven", func() {
	if os.Getenv("JIB_MAVEN_DOCKER_CONFIG_JSON") == "" {
		testsuit.T.Errorf("'JIB_MAVEN_DOCKER_CONFIG_JSON' credentials environment variable is not exported")
	} else {
		dockerConfig := os.Getenv("JIB_MAVEN_DOCKER_CONFIG_JSON")
		oc.CreateJibMavenImageRegistrySecret(dockerConfig)
	}
})
