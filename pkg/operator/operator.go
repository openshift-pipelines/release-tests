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
	EnsureTektonConfigStatusInstalled(cs.TektonConfig(), rnames)
	//Verify `pipelineSa` is created in any namespace
	AssertServiceAccount(cs, store.Namespace(), "pipeline")
	//Verify roleBindings are created in any namespace
	AssertRoleBinding(cs, store.Namespace(), "edit")
	AssertRoleBinding(cs, store.Namespace(), "pipelines-scc-rolebinding")
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
	ValidatePipelineDeployments(cs, rnames)
	ValidateTriggerDeployments(cs, rnames)
	EnsureTektonAddonsStatusInstalled(cs.TektonAddon(), rnames)
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
