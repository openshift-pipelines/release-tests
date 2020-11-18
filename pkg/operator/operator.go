package operator

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/assert"
	"github.com/openshift-pipelines/release-tests/pkg/clients"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	"github.com/openshift-pipelines/release-tests/pkg/k8s"
	"github.com/openshift-pipelines/release-tests/pkg/olm"
	"github.com/openshift-pipelines/release-tests/pkg/store"
	op "github.com/tektoncd/operator/pkg/apis/operator/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
)

// WaitForClusterCR waits for cluster CR to be created
// the function returns an error if Cluster CR is not created within timeout
func WaitForClusterCR(cs *clients.Clients, name string) *op.Config {
	cs, err := clients.InitTestingFramework(store.Clients())
	assert.NoError(err, fmt.Sprintf("Error :%s", err))
	objKey := types.NamespacedName{Name: name}
	cr := &op.Config{}

	err = wait.Poll(config.APIRetry, config.APITimeout, func() (bool, error) {
		err := cs.Client.Get(context.TODO(), objKey, cr)
		if err != nil {
			if errors.IsNotFound(err) {
				log.Printf("Waiting for availability of %s cr\n", name)
				return false, nil
			}
			return false, err
		}
		return true, nil
	})
	assert.NoError(err, fmt.Sprintf("CR: %s is not avaialble\n", name))
	return cr
}

func VerifyPipelineVersion(cs *clients.Clients, version string) {
	cr := WaitForClusterCR(cs, config.ClusterCRName)
	if !strings.HasPrefix(cr.Status.Conditions[0].Version, version) {
		testsuit.T.Errorf("Error: Invalid pipeline version \n Expected: %s , Got: %s", version, cr.Status.Conditions[0].Version)
	}
	log.Printf("Pipeline versions from CR %s", cr.Status.Conditions[0].Version)
}

func ValidateSCC(cs *clients.Clients) {
	cr := WaitForClusterCR(cs, config.ClusterCRName)
	k8s.ValidateSCCAdded(cs, cr.Spec.TargetNamespace, config.PipelineControllerName)
}

func ValidatePipelineDeployments(cs *clients.Clients) {
	cr := WaitForClusterCR(cs, config.ClusterCRName)
	k8s.ValidateDeployments(cs, cr.Spec.TargetNamespace,
		config.PipelineControllerName, config.PipelineWebhookName)
}
func ValidateTriggerDeployments(cs *clients.Clients) {
	cr := WaitForClusterCR(cs, config.ClusterCRName)
	k8s.ValidateDeployments(cs, cr.Spec.TargetNamespace,
		config.TriggerControllerName, config.TriggerWebhookName)
}

func ValidateInstalledStatus(cs *clients.Clients) {
	err := wait.PollImmediate(config.Interval, config.Timeout, func() (bool, error) {
		// Refresh Cluster CR
		cr := WaitForClusterCR(cs, config.ClusterCRName)
		if cr.Status.Conditions[0].Code != op.InstalledStatus {
			log.Printf("config.operator.tekton.dev status [%s] \n", cr.Status.Conditions[0].Code)
			return false, nil
		}
		return true, nil
	})
	assert.FailOnError(err)
}

func ValidateInstall(cs *clients.Clients) {
	log.Printf("Waiting for operator to be up and running....\n")
	ValidateInstalledStatus(cs)
	log.Printf("Operator is up\n")
}

func DeleteClusterCR(cs *clients.Clients, name string) {
	var err error
	// ensure object exists before deletion
	objKey := types.NamespacedName{Name: name}
	cr := &op.Config{}
	err = cs.Client.Get(context.TODO(), objKey, cr)

	assert.NoError(err, fmt.Sprintf("Failed to find cluster CR: %s : %s\n", name, err))

	err = wait.Poll(config.APIRetry, config.APITimeout, func() (bool, error) {
		err := cs.Client.Delete(context.TODO(), cr)
		if err != nil {
			log.Printf("Deletion of CR %s failed %s \n", name, err)
			return false, err
		}

		return true, nil
	})

	assert.NoError(err, fmt.Sprintf("%s cluster CR deletion failed\n", name))
}

// Unistall  helps you to delete operator and it's traces if any from cluster
func Uninstall(cs *clients.Clients) {
	cr := WaitForClusterCR(cs, config.ClusterCRName)

	DeleteClusterCR(cs, config.ClusterCRName)

	ns := cr.Spec.TargetNamespace
	k8s.ValidateDeploymentDeletion(cs,
		ns,
		config.PipelineControllerName,
		config.PipelineWebhookName,
		config.TriggerControllerName,
		config.TriggerWebhookName,
	)
	k8s.ValidateSCCRemoved(cs, ns, config.PipelineControllerName)

	olm.DeleteCSV()
	olm.DeleteInstallPlan()
	olm.Unsubscribe()
}
