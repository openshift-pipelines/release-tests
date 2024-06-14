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
	triggerv1alpha1 "github.com/tektoncd/operator/pkg/client/clientset/versioned/typed/operator/v1alpha1"
	"github.com/tektoncd/operator/test/utils"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EnsureTektonTriggerExists creates a TektonTrigger with the name names.TektonTrigger, if it does not exist.
func EnsureTektonTriggerExists(clients triggerv1alpha1.TektonTriggerInterface, names utils.ResourceNames) (*v1alpha1.TektonTrigger, error) {
	// If this function is called by the upgrade tests, we only create the custom resource, if it does not exist.
	ks, err := clients.Get(context.TODO(), names.TektonTrigger, metav1.GetOptions{})
	err = wait.PollUntilContextTimeout(context.TODO(), config.APIRetry, config.APITimeout, false, func(context.Context) (bool, error) {
		ks, err = clients.Get(context.TODO(), names.TektonTrigger, metav1.GetOptions{})
		if err != nil {
			if apierrs.IsNotFound(err) {
				log.Printf("Waiting for availability of triggers cr [%s]\n", names.TektonTrigger)
				return false, nil
			}
			return false, err
		}
		return true, nil
	})
	return ks, err
}

// WaitForTektonTriggerState polls the status of the TektonTrigger called name
// from client every `interval` until `inState` returns `true` indicating it
// is done, returns an error or timeout.
func WaitForTektonTriggerState(clients triggerv1alpha1.TektonTriggerInterface, name string,
	inState func(s *v1alpha1.TektonTrigger, err error) (bool, error)) (*v1alpha1.TektonTrigger, error) {
	span := logging.GetEmitableSpan(context.Background(), fmt.Sprintf("WaitForTektonTriggerState/%s/%s", name, "TektonTriggerIsReady"))
	defer span.End()

	var lastState *v1alpha1.TektonTrigger
	waitErr := wait.PollUntilContextTimeout(context.TODO(), config.APIRetry, config.APITimeout, true, func(context.Context) (bool, error) {
		lastState, err := clients.Get(context.TODO(), name, metav1.GetOptions{})
		return inState(lastState, err)
	})

	if waitErr != nil {
		return lastState, fmt.Errorf("tektonpipeline %s is not in desired state, got: %+v: %w", name, lastState, waitErr)
	}
	return lastState, nil
}

// IsTektonTriggerReady will check the status conditions of the TektonTrigger and return true if the TektonTrigger is ready.
func IsTektonTriggerReady(s *v1alpha1.TektonTrigger, err error) (bool, error) {
	return s.Status.IsReady(), err
}

// AssertTektonTriggerCRReadyStatus verifies if the TektonTrigger reaches the READY status.
func AssertTektonTriggerCRReadyStatus(clients *clients.Clients, names utils.ResourceNames) {
	if _, err := WaitForTektonTriggerState(clients.TektonTrigger(), names.TektonTrigger,
		IsTektonTriggerReady); err != nil {
		testsuit.T.Fail(fmt.Errorf("TektonTriggerCR %q failed to get to the READY status: %v", names.TektonTrigger, err))
	}
}

// TektonTriggerCRDelete deletes tha TektonTrigger to see if all resources will be deleted
func TektonTriggerCRDelete(clients *clients.Clients, crNames utils.ResourceNames) {
	if err := clients.TektonTrigger().Delete(context.TODO(), crNames.TektonTrigger, metav1.DeleteOptions{}); err != nil {
		testsuit.T.Fail(fmt.Errorf("TektonTrigger %q failed to delete: %v", crNames.TektonTrigger, err))
	}
	err := wait.PollUntilContextTimeout(clients.Ctx, config.APIRetry, config.APITimeout, true, func(context.Context) (bool, error) {
		_, err := clients.TektonTrigger().Get(context.TODO(), crNames.TektonTrigger, metav1.GetOptions{})
		if apierrs.IsNotFound(err) {
			return true, nil
		}
		return false, err
	})
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("Timed out waiting on TektonTrigger to delete, Error: %v", err))
	}

	if err := verifyNoTektonTriggerCR(clients); err != nil {
		testsuit.T.Fail(err)
	}
}

func verifyNoTektonTriggerCR(clients *clients.Clients) error {
	triggers, err := clients.TektonTrigger().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	if len(triggers.Items) > 0 {
		return errors.New("Unable to verify cluster-scoped resources are deleted if any TektonTrigger exists")
	}
	return nil
}
