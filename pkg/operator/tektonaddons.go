/*
Copyright 2020 The Tekton Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package operator

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/clients"
	"github.com/openshift-pipelines/release-tests/pkg/config"

	"knative.dev/pkg/test/logging"

	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/tektoncd/operator/pkg/apis/operator/v1alpha1"
	operatorv1alpha1 "github.com/tektoncd/operator/pkg/client/clientset/versioned/typed/operator/v1alpha1"
	"github.com/tektoncd/operator/test/utils"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EnsureTektonAddonExists creates a TektonAddon with the name names.TektonAddon, if it does not exist.
func EnsureTektonAddonExists(clients operatorv1alpha1.TektonAddonInterface, names utils.ResourceNames) (*v1alpha1.TektonAddon, error) {
	// If this function is called by the upgrade tests, we only create the custom resource, if it does not exist.
	ks, err := clients.Get(context.TODO(), names.TektonAddon, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	err = wait.PollUntilContextTimeout(context.TODO(), config.APIRetry, config.APITimeout, false, func(context.Context) (bool, error) {
		ks, err = clients.Get(context.TODO(), names.TektonAddon, metav1.GetOptions{})
		if err != nil {
			if apierrs.IsNotFound(err) {
				log.Printf("Waiting for availability of %s cr\n", names.TektonAddon)
				return false, nil
			}
			return false, err
		}
		return true, nil
	})
	return ks, err
}

// WaitForTektonAddonState polls the status of the TektonAddon called name
// from client every `interval` until `inState` returns `true` indicating it
// is done, returns an error or timeout.
func WaitForTektonAddonState(clients operatorv1alpha1.TektonAddonInterface, name string,
	inState func(s *v1alpha1.TektonAddon, err error) (bool, error)) (*v1alpha1.TektonAddon, error) {
	span := logging.GetEmitableSpan(context.Background(), fmt.Sprintf("WaitForTektonAddonState/%s/%s", name, "TektonAddonIsReady"))
	defer span.End()

	var lastState *v1alpha1.TektonAddon
	waitErr := wait.PollUntilContextTimeout(context.TODO(), config.APIRetry, config.APITimeout, true, func(context.Context) (bool, error) {
		lastState, err := clients.Get(context.TODO(), name, metav1.GetOptions{})
		return inState(lastState, err)
	})

	if waitErr != nil {
		return lastState, fmt.Errorf("tektonaddon %s is not in desired state, got: %+v: %w", name, lastState, waitErr)
	}
	return lastState, nil
}

// IsTektonAddonReady will check the status conditions of the TektonAddon and return true if the TektonAddon is ready.
func IsTektonAddonReady(s *v1alpha1.TektonAddon, err error) (bool, error) {
	return s.Status.IsReady(), err
}

func EnsureTektonAddonsStatusInstalled(clients operatorv1alpha1.TektonAddonInterface, names utils.ResourceNames) {
	err := wait.PollUntilContextTimeout(context.TODO(), config.APIRetry, config.APITimeout, true, func(context.Context) (bool, error) {
		// Refresh Cluster CR
		cr, err := EnsureTektonAddonExists(clients, names)
		if err != nil {
			testsuit.T.Fail(err)
		}
		for _, ac := range cr.Status.Conditions {
			if ac.Type != "InstallSucceeded" && ac.Status != "True" {
				log.Printf("Waiting for %s cr InstalledStatus Actual: [%s] Expected: [True]\n", names.TektonAddon, ac.Status)
				return false, nil
			}
		}
		return true, nil
	})
	if err != nil {
		testsuit.T.Fail(err)
	}
}

// AssertTektonAddonCRReadyStatus verifies if the TektonAddon reaches the READY status.
func AssertTektonAddonCRReadyStatus(clients *clients.Clients, names utils.ResourceNames) {
	if _, err := WaitForTektonAddonState(clients.TektonAddon(), names.TektonAddon,
		IsTektonAddonReady); err != nil {
		testsuit.T.Fail(fmt.Errorf("TektonAddonCR %q failed to get to the READY status: %v", names.TektonAddon, err))
	}
}

// TektonAddonCRDelete deletes tha TektonAddon to see if all resources will be deleted
func TektonAddonCRDelete(clients *clients.Clients, crNames utils.ResourceNames) {
	if err := clients.TektonAddon().Delete(context.TODO(), crNames.TektonAddon, metav1.DeleteOptions{}); err != nil {
		testsuit.T.Fail(fmt.Errorf("TektonAddon %q failed to delete: %v", crNames.TektonAddon, err))
	}
	err := wait.PollUntilContextTimeout(clients.Ctx, config.APIRetry, config.APITimeout, true, func(context.Context) (bool, error) {
		_, err := clients.TektonAddon().Get(context.TODO(), crNames.TektonAddon, metav1.GetOptions{})
		if apierrs.IsNotFound(err) {
			return true, nil
		}
		return false, err
	})
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("Timed out waiting on TektonAddon to delete, Error: %v", err))
	}

	err = verifyNoTektonAddonCR(clients)
	if err != nil {
		testsuit.T.Fail(err)
	}
}

func verifyNoTektonAddonCR(clients *clients.Clients) error {
	addons, err := clients.TektonAddon().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	if len(addons.Items) > 0 {
		return errors.New("Unable to verify cluster-scoped resources are deleted if any TektonAddon exists")
	}
	return nil
}
