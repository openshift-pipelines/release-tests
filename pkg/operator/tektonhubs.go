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
	"log"

	"github.com/openshift-pipelines/release-tests/pkg/config"

	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/tektoncd/operator/pkg/apis/operator/v1alpha1"
	hubv1alpha1 "github.com/tektoncd/operator/pkg/client/clientset/versioned/typed/operator/v1alpha1"
	"github.com/tektoncd/operator/test/utils"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EnsureTektonPipelineExists creates a TektonPipeline with the name names.TektonPipeline, if it does not exist.
func EnsureTektonHubExists(clients hubv1alpha1.TektonHubInterface, names utils.ResourceNames) (*v1alpha1.TektonHub, error) {
	// If this function is called by the upgrade tests, we only create the custom resource, if it does not exist.
	tpCR, err := clients.Get(context.TODO(), names.TektonHub, metav1.GetOptions{})
	if err == nil {
		return tpCR, err
	}
	err = wait.Poll(config.APIRetry, config.APITimeout, func() (bool, error) {
		tpCR, err = clients.Get(context.TODO(), names.TektonHub, metav1.GetOptions{})
		if err != nil {
			if apierrs.IsNotFound(err) {
				log.Printf("Waiting for availability of hub cr [%s]\n", names.TektonHub)
				return false, nil
			}
			return false, err
		}
		return true, nil
	})
	return tpCR, err
}