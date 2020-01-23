package olm

import (
	"log"

	"github.com/openshift-pipelines/release-tests/pkg/client"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	"github.com/openshift-pipelines/release-tests/pkg/helper"
	"gotest.tools/v3/icmd"
)

// DeleteClusterCR deletes Cluster config from the cluster
func DeleteClusterCR() {

	log.Printf("output: %s\n",
		helper.RunCmd(
			&helper.TknCmd{
				Args: []string{"oc", "delete", "config.operator.tekton.dev", "cluster"},
				Expected: icmd.Expected{
					ExitCode: 0,
					Err:      icmd.None,
				},
			}).Stdout())

}

// DeleteCSV deletes cluster Service Version(v0.9.*) Resource from TargetNamespace
func DeleteCSV(version string) {
	log.Printf("output: %s\n",
		helper.RunCmd(
			&helper.TknCmd{
				Args: []string{"oc", "delete",
					"csv",
					"openshift-pipelines-operator." + version,
					"-n",
					"openshift-operators"},
				Expected: icmd.Expected{
					ExitCode: 0,
					Err:      icmd.None,
				},
			}).Stdout())

}

// DeleteInstallPlan deletes installation plan
func DeleteInstallPlan() {

	installPlan := helper.RunCmd(
		&helper.TknCmd{
			Args: []string{"oc", "get", "-n", "openshift-operators",
				"subscription", "openshift-pipelines-operator",
				`-o=jsonpath={.status.installplan.name}`},
			Expected: icmd.Expected{
				ExitCode: 0,
				Err:      icmd.None,
			},
		}).Stdout()

	log.Printf("install paln %s\n", installPlan)
	res := helper.RunCmd(
		&helper.TknCmd{
			Args: []string{"oc", "delete",
				"-n", "openshift-operators",
				"installplan",
				installPlan},
			Expected: icmd.Expected{
				ExitCode: 0,
				Err:      icmd.None,
			},
		})
	log.Printf("Deleted install plan : %s\n", res.Stdout())
}

// DeleteSubscription deletes operator subscription from cluster
func DeleteSubscription() {
	log.Printf("Output %s \n", helper.RunCmd(
		&helper.TknCmd{
			Args: []string{"oc", "delete",
				"-n", "openshift-operators",
				"subscription",
				"openshift-pipelines-operator"},
			Expected: icmd.Expected{
				ExitCode: 0,
				Err:      icmd.None,
			},
		}).Stdout())

}

// DeleteOperator helps you to delete operator and it's traces if any from cluster
func DeleteOperator(cs *client.Clients, version string) {
	cr := helper.WaitForClusterCR(cs, config.ClusterCRName)

	helper.DeleteClusterCR(cs, config.ClusterCRName)

	ns := cr.Spec.TargetNamespace
	helper.ValidateDeploymentDeletion(cs,
		ns,
		config.PipelineControllerName,
		config.PipelineWebhookName,
		config.TriggerControllerName,
		config.TriggerWebhookName,
	)

	helper.ValidateSCCRemoved(cs, ns, config.PipelineControllerName)
	DeleteCSV(version)
	DeleteInstallPlan()
	DeleteSubscription()
}
