package operator

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"

	"github.com/openshift-pipelines/release-tests/pkg/assert"
	"github.com/openshift-pipelines/release-tests/pkg/clients"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	"github.com/tektoncd/operator/test/utils"
	"knative.dev/pkg/test/logging"

	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/tektoncd/operator/pkg/apis/operator/v1alpha1"
	operatorv1alpha1 "github.com/tektoncd/operator/pkg/client/clientset/versioned/typed/operator/v1alpha1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EnsureTektonHubExists creates a TektonHub with the name names.TektonHub, if it does not exist.
func EnsureTektonHubCRDExists(clients operatorv1alpha1.TektonHubInterface, names utils.ResourceNames) (*v1alpha1.TektonHub, error) {
	// If this function is called by the upgrade tests, we only create the custom resource, if it does not exist.
	ks, err := clients.Get(context.TODO(), names.TektonHub, metav1.GetOptions{})
	err = wait.Poll(config.APIRetry, config.APITimeout, func() (bool, error) {
		ks, err = clients.Get(context.TODO(), names.TektonHub, metav1.GetOptions{})
		if err != nil {
			if apierrs.IsNotFound(err) {
				log.Printf("Waiting for availability of %s CR\n", names.TektonHub)
				return false, nil
			}
			return false, err
		}
		return true, nil
	})
	return ks, err
}

// WaitForTektonHubState polls the status of the TektonHub called name
// from client every `interval` until `inState` returns `true` indicating it
// is done, returns an error or timeout.
func WaitForTektonHubState(clients operatorv1alpha1.TektonHubInterface, name string,
	inState func(s *v1alpha1.TektonHub, err error) (bool, error)) (*v1alpha1.TektonHub, error) {
	span := logging.GetEmitableSpan(context.Background(), fmt.Sprintf("WaitForTektonChainState/%s/%s", name, "TektonHubIsReady"))
	defer span.End()

	var lastState *v1alpha1.TektonHub
	waitErr := wait.PollImmediate(config.APIRetry, config.APITimeout, func() (bool, error) {
		lastState, err := clients.Get(context.TODO(), name, metav1.GetOptions{})
		return inState(lastState, err)
	})

	if waitErr != nil {
		return lastState, fmt.Errorf("tektonhub %s is not in desired state, got: %+v: %w", name, lastState, waitErr)
	}
	return lastState, nil
}

// IsTektonHubReady will check the status conditions of the TektonHub and return true if the TektonHub is ready.
func IsTektonHubReady(s *v1alpha1.TektonHub, err error) (bool, error) {
	return s.Status.IsReady(), err
}

// AssertTektonHubCRReadyStatus verifies if the TektonHub reaches the READY status.
func AssertTektonHubCRReadyStatus(clients *clients.Clients, names utils.ResourceNames) {
	if _, err := WaitForTektonHubState(clients.TektonHub(), names.TektonHub,
		IsTektonHubReady); err != nil {
		assert.FailOnError(fmt.Errorf("TektonHubCR %q failed to get to the READY status: %v", names.TektonHub, err))
	}
}

func VerifyNoTektonHubCR(clients *clients.Clients) error {
	log.Print("Verifying that TektonHub CR is not available")
	hubs, err := clients.TektonHub().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	if len(hubs.Items) > 0 {
		return errors.New("unable to verify cluster-scoped resources are deleted if any TektonHub exists")
	}
	return nil
}

func VerifyTektonHubURLs(clients *clients.Clients) (string, string, error) {

	crdList, err := clients.TektonHub().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return "", "", err
	}

	tektonHubCRDExists := false
	var apiURL, uiURL string
	for _, crd := range crdList.Items {
		if crd.Name == "hub" {
			if err := VerifyURL(crd.Status.ApiRouteUrl); err != nil {
				return "", "", err
			}
			if err := VerifyURL(crd.Status.UiRouteUrl); err != nil {
				return "", "", err
			}
			tektonHubCRDExists = true
			apiURL = crd.Status.ApiRouteUrl
			uiURL = crd.Status.UiRouteUrl
			break
		}
	}

	if tektonHubCRDExists {
		return apiURL, uiURL, nil
	} else {
		return "", "", fmt.Errorf("Tekton Hub CRD does not exist")
	}
}

func VerifyURL(rawurl string) error {
	u, err := url.Parse(rawurl)
	if err != nil {
		return fmt.Errorf("invalid URL: %s", err.Error())
	}
	if u.Scheme == "" || u.Host == "" {
		return fmt.Errorf("missing scheme or host")
	}
	return nil
}
