package olm

import (
	"log"

	"github.com/openshift-pipelines/release-tests/pkg/client"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	"github.com/openshift-pipelines/release-tests/pkg/helper"
)

// DeleteClusterCR deletes Cluster config from the cluster
func DeleteClusterCR() {

	log.Printf("output: %s\n",
		helper.CmdShouldPass("oc",
			"delete",
			"config.operator.tekton.dev",
			"cluster"))
}

// DeleteCSV deletes cluster Service Version(v0.9.*) Resource from TargetNamespace
func DeleteCSV(version string) {
	log.Printf("output: %s\n",
		helper.CmdShouldPass("oc",
			"delete",
			"csv",
			"openshift-pipelines-operator."+version,
			"-n",
			"openshift-operators"))
}

// DeleteInstallPlan deletes installation plan
func DeleteInstallPlan() {

	installPlan := helper.CmdShouldPass(
		"oc", "get", "-n", "openshift-operators",
		"subscription", "openshift-pipelines-operator",
		`-o=jsonpath={.status.installplan.name}`,
	)
	log.Printf("install paln %s\n", installPlan)
	res := helper.CmdShouldPass(
		"oc", "delete",
		"-n", "openshift-operators",
		"installplan",
		installPlan)
	log.Printf("Deleted install plan : %s\n", res)
}

// DeleteSubscription deletes operator subscription from cluster
func DeleteSubscription() {
	log.Printf("Output %s \n", helper.CmdShouldPass(
		"oc", "delete",
		"-n", "openshift-operators",
		"subscription",
		"openshift-pipelines-operator",
	))
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
