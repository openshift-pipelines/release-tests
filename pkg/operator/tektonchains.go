package operator

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/openshift-pipelines/release-tests/pkg/assert"
	"github.com/openshift-pipelines/release-tests/pkg/clients"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	"knative.dev/pkg/test/logging"

	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/tektoncd/operator/pkg/apis/operator/v1alpha1"
	operatorv1alpha1 "github.com/tektoncd/operator/pkg/client/clientset/versioned/typed/operator/v1alpha1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EnsureTektonChainExists creates a TektonChain with the name names.TektonChain, if it does not exist.
func EnsureTektonChainExists(clients operatorv1alpha1.TektonChainInterface, names config.ResourceNames) (*v1alpha1.TektonChain, error) {
	// If this function is called by the upgrade tests, we only create the custom resource, if it does not exist.
	ks, err := clients.Get(context.TODO(), names.TektonChain, metav1.GetOptions{})
	err = wait.Poll(config.APIRetry, config.APITimeout, func() (bool, error) {
		ks, err = clients.Get(context.TODO(), names.TektonChain, metav1.GetOptions{})
		if err != nil {
			if apierrs.IsNotFound(err) {
				log.Printf("Waiting for availability of %s CR\n", names.TektonChain)
				return false, nil
			}
			return false, err
		}
		return true, nil
	})
	return ks, err
}

// WaitForTektonChainState polls the status of the TektonChain called name
// from client every `interval` until `inState` returns `true` indicating it
// is done, returns an error or timeout.
func WaitForTektonChainState(clients operatorv1alpha1.TektonChainInterface, name string,
	inState func(s *v1alpha1.TektonChain, err error) (bool, error)) (*v1alpha1.TektonChain, error) {
	span := logging.GetEmitableSpan(context.Background(), fmt.Sprintf("WaitForTektonChainState/%s/%s", name, "TektonChainIsReady"))
	defer span.End()

	var lastState *v1alpha1.TektonChain
	waitErr := wait.PollImmediate(config.APIRetry, config.APITimeout, func() (bool, error) {
		lastState, err := clients.Get(context.TODO(), name, metav1.GetOptions{})
		return inState(lastState, err)
	})

	if waitErr != nil {
		return lastState, fmt.Errorf("TektonChain %s is not in desired state, got: %+v: %w", name, lastState, waitErr)
	}
	return lastState, nil
}

// IsTektonChainReady will check the status conditions of the TektonChain and return true if the TektonChain is ready.
func IsTektonChainReady(s *v1alpha1.TektonChain, err error) (bool, error) {
	return s.Status.IsReady(), err
}

// AssertTektonChainCRReadyStatus verifies if the TektonChain reaches the READY status.
func AssertTektonChainCRReadyStatus(clients *clients.Clients, names config.ResourceNames) {
	if _, err := WaitForTektonChainState(clients.TektonChain(), names.TektonChain,
		IsTektonChainReady); err != nil {
		assert.FailOnError(fmt.Errorf("TektonChainCR %q failed to get to the READY status: %v", names.TektonChain, err))
	}
}

func VerifyNoTektonChainCR(clients *clients.Clients) error {
	log.Print("Verifying that TektonChain CR is not available")
	addons, err := clients.TektonChain().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	if len(addons.Items) > 0 {
		return errors.New("unable to verify cluster-scoped resources are deleted if any TektonChain exists")
	}
	return nil
}
