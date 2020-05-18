package operator

import (
	"log"

	"github.com/openshift-pipelines/release-tests/pkg/cmd"
)

// DeleteSubscription deletes operator subscription from cluster
func DeleteSubscription() {
	log.Printf("Output %s \n", cmd.MustSucceed(
		"oc", "delete", "-n", "openshift-operators",
		"subscription", "openshift-pipelines-operator",
	).Stdout())
}
