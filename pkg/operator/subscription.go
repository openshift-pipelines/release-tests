package operator

import (
	"log"

	"github.com/openshift-pipelines/release-tests/pkg/cmd"
	"gotest.tools/v3/icmd"
)

// // DeleteClusterCR deletes Cluster config from the cluster
// func DeleteClusterCR() {
//
// 	log.Printf("output: %s\n",
// 		cmd.AssertOutput(
// 			&cmd.Cmd{
// 				Args: []string{"oc", "delete", "config.operator.tekton.dev", "cluster"},
// 				Expected: icmd.Expected{
// 					ExitCode: 0,
// 					Err:      icmd.None,
// 				},
// 			}).Stdout())
//
// }

// DeleteInstallPlan deletes installation plan
func DeleteInstallPlan() {

	installPlan := cmd.AssertOutput(
		&cmd.Cmd{
			Args: []string{"oc", "get", "-n", "openshift-operators",
				"subscription", "openshift-pipelines-operator",
				`-o=jsonpath={.status.installplan.name}`},
			Expected: icmd.Success,
		}).Stdout()

	log.Printf("install paln %s\n", installPlan)
	res := cmd.AssertOutput(
		&cmd.Cmd{
			Args: []string{"oc", "delete",
				"-n", "openshift-operators",
				"installplan",
				installPlan},
			Expected: icmd.Success,
		})
	log.Printf("Deleted install plan : %s\n", res.Stdout())
}

// DeleteSubscription deletes operator subscription from cluster
func DeleteSubscription() {
	log.Printf("Output %s \n", cmd.AssertOutput(
		&cmd.Cmd{
			Args: []string{
				"oc", "delete", "-n", "openshift-operators",
				"subscription", "openshift-pipelines-operator"},
			Expected: icmd.Success,
		}).Stdout())

}
