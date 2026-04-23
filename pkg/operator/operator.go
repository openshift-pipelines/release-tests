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
	approvalgate "github.com/openshift-pipelines/release-tests/pkg/manualapprovalgate"
	"github.com/openshift-pipelines/release-tests/pkg/olm"
	"github.com/openshift-pipelines/release-tests/pkg/opc"
	"github.com/openshift-pipelines/release-tests/pkg/store"
	"github.com/tektoncd/operator/test/utils"
)

func WaitForTektonConfigCR(c *clients.Clients, rnames utils.ResourceNames) {
	if _, err := EnsureTektonConfigExists(c.TektonConfig(), rnames); err != nil {
		testsuit.T.Fail(fmt.Errorf("TektonConfig doesn't exists\n %v", err))
	}
}

func ValidateRBAC(c *clients.Clients, rnames utils.ResourceNames) {
	log.Printf("Verifying that TektonConfig status is \"installed\"\n")
	EnsureTektonConfigStatusInstalled(c.TektonConfig(), rnames)

	AssertServiceAccountPresent(c, store.Namespace(), "pipeline")
	AssertClusterRolePresent(c, "pipelines-scc-clusterrole")
	AssertConfigMapPresent(c, store.Namespace(), "config-service-cabundle")
	AssertConfigMapPresent(c, store.Namespace(), "config-trusted-cabundle")
	AssertRoleBindingPresent(c, store.Namespace(), "openshift-pipelines-edit")
	AssertRoleBindingPresent(c, store.Namespace(), "pipelines-scc-rolebinding")
	AssertSCCPresent(c, "pipelines-scc")
}

func ValidateRBACAfterDisable(c *clients.Clients, rnames utils.ResourceNames) {
	EnsureTektonConfigStatusInstalled(c.TektonConfig(), rnames)
	// Verify `pipelineSa` exists in the existing namespace
	AssertServiceAccountPresent(c, store.Namespace(), "pipeline")
	// Verify clusterrole does not create
	AssertClusterRoleNotPresent(c, "pipelines-scc-clusterrole")
	// Verify roleBindings is not created in any namespace
	AssertRoleBindingNotPresent(c, store.Namespace(), "edit")
	AssertRoleBindingNotPresent(c, store.Namespace(), "pipelines-scc-rolebinding")
	AssertSCCNotPresent(c, "pipelines-scc")
}

func ValidateCABundleConfigMaps(c *clients.Clients, rnames utils.ResourceNames) {
	log.Printf("Verifying that TektonConfig status is \"installed\"\n")
	EnsureTektonConfigStatusInstalled(c.TektonConfig(), rnames)
	// Verify CA Bundle ConfigMaps are created
	AssertConfigMapPresent(c, store.Namespace(), "config-service-cabundle")
	AssertConfigMapPresent(c, store.Namespace(), "config-trusted-cabundle")
}

func ValidatePipelineDeployments(c *clients.Clients, rnames utils.ResourceNames) {
	if _, err := EnsureTektonPipelineExists(c.TektonPipeline(), rnames); err != nil {
		testsuit.T.Fail(fmt.Errorf("TektonPipelines doesn't exists\n %v", err))
	}
	k8s.ValidateDeployments(c, rnames.TargetNamespace,
		config.PipelineControllerName, config.PipelineWebhookName)
}

func ValidateTriggerDeployments(c *clients.Clients, rnames utils.ResourceNames) {
	if _, err := EnsureTektonTriggerExists(c.TektonTrigger(), rnames); err != nil {
		testsuit.T.Fail(fmt.Errorf("TektonTriggers doesn't exists\n %v", err))
	}
	k8s.ValidateDeployments(c, rnames.TargetNamespace,
		config.TriggerControllerName, config.TriggerWebhookName)
}

func ValidateChainsDeployments(c *clients.Clients, rnames utils.ResourceNames) {
	if _, err := EnsureTektonChainsExists(c.TektonChains(), rnames); err != nil {
		testsuit.T.Fail(fmt.Errorf("TektonChains doesn't exists\n %v", err))
	}
	k8s.ValidateDeployments(c, rnames.TargetNamespace,
		config.ChainsControllerName)
}

func ValidateHubDeployments(c *clients.Clients, rnames utils.ResourceNames) {
	if _, err := EnsureTektonHubsExists(c.TektonHub(), rnames); err != nil {
		testsuit.T.Fail(fmt.Errorf("TektonHub doesn't exists\n %v", err))
	}
	k8s.ValidateDeployments(c, rnames.TargetNamespace,
		config.HubApiName, config.HubDbName, config.HubUiName)
}

func ValidateManualApprovalGateDeployments(c *clients.Clients, rnames utils.ResourceNames) {
	if _, err := approvalgate.EnsureManualApprovalGateExists(c.ManualApprovalGate(), rnames); err != nil {
		testsuit.T.Fail(fmt.Errorf("manual approval gate doesn't exists\n %v", err))
	}
	k8s.ValidateDeployments(c, rnames.TargetNamespace,
		config.MAGController, config.MAGWebHook)
}

func ValidateOperatorInstallStatus(c *clients.Clients, rnames utils.ResourceNames) {
	operatorVersion := opc.GetOPCServerVersion("operator")
	if strings.Contains(operatorVersion, "unknown") {
		testsuit.T.Errorf("Operator is not installed")
	}
	log.Printf("Waiting for operator to be up and running....\n")
	EnsureTektonConfigStatusInstalled(c.TektonConfig(), rnames)
	log.Printf("Operator is up\n")
}

func DeleteTektonConfigCR(c *clients.Clients, rnames utils.ResourceNames) {
	TektonConfigCRDelete(c, rnames)
}

// Uninstall helps you to delete operator and it's traces if any from cluster
func Uninstall(c *clients.Clients, rnames utils.ResourceNames) {
	log.Printf("output: %s\n", cmd.MustSucceed("oc", "delete", "--ignore-not-found", "TektonHub", "hub").Stdout())
	log.Printf("output: %s\n", cmd.MustSucceed("oc", "delete", "--ignore-not-found", "tektonresults", "result").Stdout())
	log.Printf("output: %s\n", cmd.MustSucceed("oc", "delete", "--ignore-not-found", "manualapprovalgate", "manual-approval-gate").Stdout())
	DeleteTektonConfigCR(c, rnames)
	k8s.ValidateDeploymentDeletion(c,
		rnames.TargetNamespace,
		config.PipelineControllerName,
		config.PipelineWebhookName,
		config.TriggerControllerName,
		config.TriggerWebhookName,
		config.ChainsControllerName,
	)
	k8s.ValidateSCCRemoved(c, rnames.TargetNamespace, config.PipelineControllerName)
	olm.OperatorCleanup(c, config.Flags.SubscriptionName)
}
