package operator

import (
	"log"
	"strings"

	. "github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/clients"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	"github.com/openshift-pipelines/release-tests/pkg/helper"
	"github.com/openshift-pipelines/release-tests/pkg/olm"
	op "github.com/tektoncd/operator/pkg/apis/operator/v1alpha1"
)

func VerifyPipelineVersion(cs *clients.Clients, version string) {
	cr := helper.WaitForClusterCR(cs, config.ClusterCRName)
	if strings.HasPrefix(cr.Status.Conditions[0].Version, version) {
		log.Printf("Pipeline versions from CR %s", cr.Status.Conditions[0].Version)
	} else {
		T.Errorf("Error: Invalid pipeline version %s", cr.Status.Conditions[0].Version)
	}
}

func ValidateSCC(cs *clients.Clients) {
	cr := helper.WaitForClusterCR(cs, config.ClusterCRName)
	helper.ValidateSCCAdded(cs, cr.Spec.TargetNamespace, config.PipelineControllerName)
}

func ValidatePipelineDeployments(cs *clients.Clients) {
	cr := helper.WaitForClusterCR(cs, config.ClusterCRName)
	helper.ValidateDeployments(cs, cr.Spec.TargetNamespace,
		config.PipelineControllerName, config.PipelineWebhookName)
}
func ValidateTriggerDeployments(cs *clients.Clients) {
	cr := helper.WaitForClusterCR(cs, config.ClusterCRName)
	helper.ValidateDeployments(cs, cr.Spec.TargetNamespace,
		config.TriggerControllerName, config.TriggerWebhookName)
}

func ValidateInstalledStatus(cs *clients.Clients) {
	// Refresh Cluster CR
	cr := helper.WaitForClusterCR(cs, config.ClusterCRName)

	if code := cr.Status.Conditions[0].Code; code != op.InstalledStatus {
		T.Errorf("Expected code to be %s but got %s", op.InstalledStatus, code)
	}
}

func ValidateInstall(cs *clients.Clients) {
	log.Printf("Waiting for operator to be up and running....\n")

	ValidatePipelineDeployments(cs)
	ValidateTriggerDeployments(cs)

	// Refresh Cluster CR
	cr := helper.WaitForClusterCR(cs, config.ClusterCRName)

	if code := cr.Status.Conditions[0].Code; code != op.InstalledStatus {
		T.Errorf("Expected code to be %s but got %s", op.InstalledStatus, code)
	}
	log.Printf("Operator is up\n")

}

func InstallOperator(version string) (*clients.Clients, string, func()) {
	cs, ns, cleanup := olm.Subscribe(version)
	ValidateInstall(cs)
	helper.VerifyServiceAccountExists(cs.KubeClient, ns)
	return cs, ns, cleanup
}
