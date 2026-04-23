// Package ecosystem provides Gauge test steps for ecosystem tasks like jib-maven
package ecosystem

import (
	"os"

	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/cmd"
)

var repo = os.Getenv("JIB_MAVEN_REPOSITORY")

var _ = gauge.Step("Verify that jib-maven image registry variable is exported", func() {
	if repo == "" {
		testsuit.T.Errorf("'JIB_MAVEN_REPOSITORY' environment variable is not exported")
	}
})

var _ = gauge.Step("Create secret with jib-maven image registry credentials", func() {
	if os.Getenv("JIB_MAVEN_DOCKER_CONFIG_JSON") == "" {
		testsuit.T.Errorf("'JIB_MAVEN_DOCKER_CONFIG_JSON' robot credentials environment variable is not exported")
	} else {
		dockerConfig := os.Getenv("JIB_MAVEN_DOCKER_CONFIG_JSON")
		// Create secret for jib-maven image registry
		cmd.MustSucceed(
			"oc", "create", "secret", "generic", "jib-maven-image-registry-credentials",
			"--from-literal=.dockerconfigjson="+dockerConfig,
			"--from-literal=config.json="+dockerConfig,
			"--type=kubernetes.io/dockerconfigjson",
		)
		// Link secret to pipeline service account
		cmd.MustSucceed("oc", "secrets", "link", "pipeline", "jib-maven-image-registry-credentials", "--for=pull,mount")
	}
})
