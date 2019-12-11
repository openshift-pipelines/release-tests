package olm

import (
	"log"
	"path/filepath"

	"github.com/openshift-pipelines/release-tests/pkg/client"
	"github.com/openshift-pipelines/release-tests/pkg/helper"
)

// Subscribe helps you to subscribe specific version pipelines operator from canary channel to OCP cluster
func Subscribe(version string) (*client.Clients, string, func()) {
	path := filepath.Join(helper.RootDir(), "../config/subscription.yaml")
	log.Printf("output: %s\n",
		helper.CmdShouldPass("oc", "apply", "-f", path))

	cs, ns, cleanupNs := helper.NewClientSet()

	cleanup := func() {
		DeleteOperator(cs, version)
		cleanupNs()
	}

	return cs, ns, cleanup
}
