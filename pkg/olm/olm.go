package olm

import (
	"log"
	"path/filepath"

	"github.com/openshift-pipelines/release-tests/pkg/clients"
	"github.com/openshift-pipelines/release-tests/pkg/helper"
	"gotest.tools/v3/icmd"
)

// Subscribe helps you to subscribe specific version pipelines operator from canary channel to OCP cluster
func Subscribe(version string) (*clients.Clients, string, func()) {
	path := filepath.Join(helper.RootDir(), "../config/subscription.yaml")
	log.Printf("output: %s\n",
		helper.RunCmd(
			&helper.TknCmd{
				Args: []string{"oc", "apply", "-f", path},
				Expected: icmd.Expected{
					ExitCode: 0,
					Err:      icmd.None,
				},
			}))

	cs, ns, cleanupNs := helper.NewClientSet()

	cleanup := func() {
		DeleteOperator(cs, version)
		cleanupNs()
	}

	return cs, ns, cleanup
}
