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

	"knative.dev/pkg/test/logging"

	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/clients"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/tektoncd/operator/pkg/apis/operator/v1alpha1"
	configv1alpha1 "github.com/tektoncd/operator/pkg/client/clientset/versioned/typed/operator/v1alpha1"
	"github.com/tektoncd/operator/test/utils"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EnsureTektonConfigExists creates a TektonConfig with the name names.TektonConfig, if it does not exist.
func EnsureTektonConfigExists(clients configv1alpha1.TektonConfigInterface, names utils.ResourceNames) (*v1alpha1.TektonConfig, error) {
	// If this function is called by the upgrade tests, we only create the custom resource, if it does not exist.
	tcCR, err := clients.Get(context.TODO(), names.TektonConfig, metav1.GetOptions{})
	if err == nil {
		return tcCR, err
	}

	err = wait.PollUntilContextTimeout(context.TODO(), config.APIRetry, config.APITimeout, false, func(context.Context) (bool, error) {
		tcCR, err = clients.Get(context.TODO(), names.TektonConfig, metav1.GetOptions{})
		if err != nil {
			if apierrs.IsNotFound(err) {
				log.Printf("Waiting for availability of %s cr\n", names.TektonConfig)
				return false, nil
			}
			return false, err
		}
		return true, nil
	})
	return tcCR, err
}

// WaitForTektonConfigState polls the status of the TektonConfig called name
// from client every `interval` until `inState` returns `true` indicating it
// is done, returns an error or timeout.
func WaitForTektonConfigState(clients configv1alpha1.TektonConfigInterface, name string,
	inState func(s *v1alpha1.TektonConfig, err error) (bool, error)) (*v1alpha1.TektonConfig, error) {
	span := logging.GetEmitableSpan(context.Background(), fmt.Sprintf("WaitForTektonConfigState/%s/%s", name, "TektonConfigIsReady"))
	defer span.End()

	var lastState *v1alpha1.TektonConfig
	waitErr := wait.PollUntilContextTimeout(context.TODO(), config.APIRetry, config.APITimeout, true, func(context.Context) (bool, error) {
		lastState, err := clients.Get(context.TODO(), name, metav1.GetOptions{})
		return inState(lastState, err)
	})

	if waitErr != nil {
		return lastState, fmt.Errorf("tektonconfig %s is not in desired state, got: %+v: %w", name, lastState, waitErr)
	}
	return lastState, nil
}

// IsTektonConfigReady will check the status conditions of the TektonConfig and return true if the TektonConfig is ready.
func IsTektonConfigReady(s *v1alpha1.TektonConfig, err error) (bool, error) {
	return s.Status.IsReady(), err
}

func EnsureTektonConfigStatusInstalled(clients configv1alpha1.TektonConfigInterface, names utils.ResourceNames) {
	err := wait.PollUntilContextTimeout(context.TODO(), config.APIRetry, config.APITimeout, true, func(context.Context) (bool, error) {
		// Refresh Cluster CR
		cr, err := EnsureTektonConfigExists(clients, names)
		if err != nil {
			testsuit.T.Fail(err)
		}
		for _, cc := range cr.Status.Conditions {
			if cc.Type != "InstallSucceeded" && cc.Status != "True" {
				log.Printf("Waiting for %s cr InstalledStatus Actual: [%s] Expected: [True]\n", names.TektonConfig, cc.Status)
				return false, nil
			}
		}
		return true, nil
	})
	if err != nil {
		testsuit.T.Fail(err)
	}
}

// AssertTektonConfigCRReadyStatus verifies if the TektonConfig reaches the READY status.
func AssertTektonConfigCRReadyStatus(clients *clients.Clients, names utils.ResourceNames) {
	if _, err := WaitForTektonConfigState(clients.TektonConfig(), names.TektonConfig, IsTektonConfigReady); err != nil {
		testsuit.T.Fail(fmt.Errorf("TektonConfigCR %q failed to get to the READY status: %v", names.TektonConfig, err))
	}
}

// TektonConfigCRDelete deletes tha TektonConfig to see if all resources will be deleted
func TektonConfigCRDelete(clients *clients.Clients, crNames utils.ResourceNames) {
	if err := clients.TektonConfig().Delete(context.TODO(), crNames.TektonConfig, metav1.DeleteOptions{}); err != nil {
		testsuit.T.Fail(fmt.Errorf("TektonConfigCR %q failed to delete: %v", crNames.TektonConfig, err))
	}
	err := wait.PollUntilContextTimeout(clients.Ctx, config.APIRetry, config.APITimeout, true, func(context.Context) (bool, error) {
		_, err := clients.TektonConfig().Get(context.TODO(), crNames.TektonConfig, metav1.GetOptions{})
		if apierrs.IsNotFound(err) {
			return true, nil
		}
		return false, err
	})
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("Timed out waiting on TektonConfigCR to delete, Error: %v", err))
	}
	err = verifyNoTektonConfigCR(clients)
	if err != nil {
		testsuit.T.Fail(err)
	}
}

func verifyNoTektonConfigCR(clients *clients.Clients) error {
	configs, err := clients.TektonConfig().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	if len(configs.Items) > 0 {
		return errors.New("Unable to verify cluster-scoped resources are deleted if any TektonConfig exists")
	}
	return nil
}
