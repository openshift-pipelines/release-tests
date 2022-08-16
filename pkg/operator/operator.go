package operator

import (
	"log"

	"github.com/openshift-pipelines/release-tests/pkg/clients"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	"github.com/openshift-pipelines/release-tests/pkg/k8s"
	"github.com/openshift-pipelines/release-tests/pkg/olm"
	"github.com/openshift-pipelines/release-tests/pkg/store"
)

func WaitForTektonConfigCR(cs *clients.Clients, rnames config.ResourceNames) {
	EnsureTektonConfigExists(cs.TektonConfig(), rnames)
}

func ValidateRBAC(cs *clients.Clients, rnames config.ResourceNames) {
	log.Printf("Verifying that TektonConfig status is \"installed\"\n")
	EnsureTektonConfigStatusInstalled(cs.TektonConfig(), rnames)

	AssertServiceAccountPesent(cs, store.Namespace(), "pipeline")
	AssertClusterRolePresent(cs, "pipelines-scc-clusterrole")
	AssertConfigMapPresent(cs, store.Namespace(), "config-service-cabundle")
	AssertConfigMapPresent(cs, store.Namespace(), "config-trusted-cabundle")
	AssertRoleBindingPresent(cs, store.Namespace(), "openshift-pipelines-edit")
	AssertRoleBindingPresent(cs, store.Namespace(), "pipelines-scc-rolebinding")
}

func ValidateRBACAfterDisable(cs *clients.Clients, rnames config.ResourceNames) {
	EnsureTektonConfigStatusInstalled(cs.TektonConfig(), rnames)
	//Verify `pipelineSa` exists in the existing namespace
	AssertServiceAccountPesent(cs, store.Namespace(), "pipeline")
	//Verify clusterrole does not create
	AssertClusterRoleNotPresent(cs, "pipelines-scc-clusterrole")
	//Verify configmaps is not created in any namespace
	AssertConfigMapNotPresent(cs, store.Namespace(), "config-service-cabundle")
	AssertConfigMapNotPresent(cs, store.Namespace(), "config-trusted-cabundle")
	//Verify roleBindings is not created in any namespace
	AssertRoleBindingNotPresent(cs, store.Namespace(), "edit")
	AssertRoleBindingNotPresent(cs, store.Namespace(), "pipelines-scc-rolebinding")
}

func ValidatePipelineDeployments(cs *clients.Clients, rnames config.ResourceNames) {
	EnsureTektonPipelineExists(cs.TektonPipeline(), rnames)
	k8s.ValidateDeployments(cs, rnames.TargetNamespace,
		config.PipelineControllerName, config.PipelineWebhookName)
}

func ValidateTriggerDeployments(cs *clients.Clients, rnames config.ResourceNames) {
	EnsureTektonTriggerExists(cs.TektonTrigger(), rnames)
	k8s.ValidateDeployments(cs, rnames.TargetNamespace,
		config.TriggerControllerName, config.TriggerWebhookName)
}

func ValidateOperatorInstallStatus(cs *clients.Clients, rnames config.ResourceNames) {
	log.Printf("Waiting for operator to be up and running....\n")
	EnsureTektonConfigStatusInstalled(cs.TektonConfig(), rnames)
	log.Printf("Operator is up\n")
}

func DeleteTektonConfigCR(cs *clients.Clients, rnames config.ResourceNames) {
	EnsureTektonConfigStatusInstalled(cs.TektonConfig(), rnames)
	TektonConfigCRDelete(cs, rnames)
}

// Unistall helps you to delete operator and it's traces if any from cluster
func Uninstall(cs *clients.Clients, rnames config.ResourceNames) {
	DeleteTektonConfigCR(cs, rnames)
	k8s.ValidateDeploymentDeletion(cs,
		rnames.TargetNamespace,
		config.PipelineControllerName,
		config.PipelineWebhookName,
		config.TriggerControllerName,
		config.TriggerWebhookName,
	)
	k8s.ValidateSCCRemoved(cs, rnames.TargetNamespace, config.PipelineControllerName)
	olm.OperatorCleanup(cs, config.Flags.SubscriptionName)
}
