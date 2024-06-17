package operator

import (
	"fmt"
	"log"
	"strings"

	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/clients"
	"github.com/openshift-pipelines/release-tests/pkg/cmd"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	"github.com/openshift-pipelines/release-tests/pkg/k8s"
	"github.com/openshift-pipelines/release-tests/pkg/olm"
	"github.com/openshift-pipelines/release-tests/pkg/store"
	"github.com/tektoncd/operator/test/utils"
)

func WaitForTektonConfigCR(cs *clients.Clients, rnames utils.ResourceNames) {
	if _, err := EnsureTektonConfigExists(cs.TektonConfig(), rnames); err != nil {
		testsuit.T.Fail(fmt.Errorf("TektonConfig doesn't exists\n %v", err))
	}
}

func ValidateRBAC(cs *clients.Clients, rnames utils.ResourceNames) {
	log.Printf("Verifying that TektonConfig status is \"installed\"\n")
	EnsureTektonConfigStatusInstalled(cs.TektonConfig(), rnames)

	AssertServiceAccountPresent(cs, store.Namespace(), "pipeline")
	AssertClusterRolePresent(cs, "pipelines-scc-clusterrole")
	AssertConfigMapPresent(cs, store.Namespace(), "config-service-cabundle")
	AssertConfigMapPresent(cs, store.Namespace(), "config-trusted-cabundle")
	AssertRoleBindingPresent(cs, store.Namespace(), "openshift-pipelines-edit")
	AssertRoleBindingPresent(cs, store.Namespace(), "pipelines-scc-rolebinding")
	AssertSCCPresent(cs, "pipelines-scc")
}

func ValidateRBACAfterDisable(cs *clients.Clients, rnames utils.ResourceNames) {
	EnsureTektonConfigStatusInstalled(cs.TektonConfig(), rnames)
	//Verify `pipelineSa` exists in the existing namespace
	AssertServiceAccountPresent(cs, store.Namespace(), "pipeline")
	//Verify clusterrole does not create
	AssertClusterRoleNotPresent(cs, "pipelines-scc-clusterrole")
	//Verify configmaps is not created in any namespace
	AssertConfigMapNotPresent(cs, store.Namespace(), "config-service-cabundle")
	AssertConfigMapNotPresent(cs, store.Namespace(), "config-trusted-cabundle")
	//Verify roleBindings is not created in any namespace
	AssertRoleBindingNotPresent(cs, store.Namespace(), "edit")
	AssertRoleBindingNotPresent(cs, store.Namespace(), "pipelines-scc-rolebinding")
	AssertSCCNotPresent(cs, "pipelines-scc")
}

func ValidatePipelineDeployments(cs *clients.Clients, rnames utils.ResourceNames) {
	if _, err := EnsureTektonPipelineExists(cs.TektonPipeline(), rnames); err != nil {
		testsuit.T.Fail(fmt.Errorf("TektonPipelines doesn't exists\n %v", err))
	}
	k8s.ValidateDeployments(cs, rnames.TargetNamespace,
		config.PipelineControllerName, config.PipelineWebhookName)
}

func ValidateTriggerDeployments(cs *clients.Clients, rnames utils.ResourceNames) {
	if _, err := EnsureTektonTriggerExists(cs.TektonTrigger(), rnames); err != nil {
		testsuit.T.Fail(fmt.Errorf("TektonTriggers doesn't exists\n %v", err))
	}
	k8s.ValidateDeployments(cs, rnames.TargetNamespace,
		config.TriggerControllerName, config.TriggerWebhookName)
}

func ValidateChainsDeployments(cs *clients.Clients, rnames utils.ResourceNames) {
	if _, err := EnsureTektonChainsExists(cs.TektonChains(), rnames); err != nil {
		testsuit.T.Fail(fmt.Errorf("TektonChains doesn't exists\n %v", err))
	}
	k8s.ValidateDeployments(cs, rnames.TargetNamespace,
		config.ChainsControllerName)
}

func ValidateHubDeployments(cs *clients.Clients, rnames utils.ResourceNames) {
	if _, err := EnsureTektonHubsExists(cs.TektonHub(), rnames); err != nil {
		testsuit.T.Fail(fmt.Errorf("TektonHub doesn't exists\n %v", err))
	}
	k8s.ValidateDeployments(cs, rnames.TargetNamespace,
		config.HubApiName, config.HubDbName, config.HubUiName)
}

func ValidateManualApprovalGateDeployments(cs *clients.Clients, rnames utils.ResourceNames) {
	if _, err := EnsureManualApprovalGateExists(cs.ManualApprovalGate(), rnames); err != nil {
		testsuit.T.Fail(fmt.Errorf("Manual approval gate doesn't exists\n %v", err))
	}
	k8s.ValidateDeployments(cs, rnames.TargetNamespace,
		config.MAGController, config.MAGWebHook)
}

func ValidateOperatorInstallStatus(cs *clients.Clients, rnames utils.ResourceNames) {
	operatorVersion := cmd.MustSucceed("tkn", "version", "--component", "operator").Stdout()
	if strings.Contains(operatorVersion, "unknown") {
		testsuit.T.Errorf("Operator is not installed")
		return
	}
	log.Printf("Waiting for operator to be up and running....\n")
	EnsureTektonConfigStatusInstalled(cs.TektonConfig(), rnames)
	log.Printf("Operator is up\n")
}

func DeleteTektonConfigCR(cs *clients.Clients, rnames utils.ResourceNames) {
	TektonConfigCRDelete(cs, rnames)
}

// Uninstall helps you to delete operator and it's traces if any from cluster
func Uninstall(cs *clients.Clients, rnames utils.ResourceNames) {
	log.Printf("output: %s\n", cmd.MustSucceed("oc", "delete", "--ignore-not-found", "TektonHub", "hub").Stdout())
	log.Printf("output: %s\n", cmd.MustSucceed("oc", "delete", "--ignore-not-found", "tektonresults", "result").Stdout())
	DeleteTektonConfigCR(cs, rnames)
	k8s.ValidateDeploymentDeletion(cs,
		rnames.TargetNamespace,
		config.PipelineControllerName,
		config.PipelineWebhookName,
		config.TriggerControllerName,
		config.TriggerWebhookName,
		config.ChainsControllerName,
	)
	k8s.ValidateSCCRemoved(cs, rnames.TargetNamespace, config.PipelineControllerName)
	olm.OperatorCleanup(cs, config.Flags.SubscriptionName)
}
