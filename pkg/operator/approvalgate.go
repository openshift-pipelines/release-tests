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
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/cmd"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	"github.com/tektoncd/operator/pkg/apis/operator/v1alpha1"
	mag "github.com/tektoncd/operator/pkg/client/clientset/versioned/typed/operator/v1alpha1"
	"github.com/tektoncd/operator/test/utils"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

func EnsureManualApprovalGateExists(clients mag.ManualApprovalGateInterface, names utils.ResourceNames) (*v1alpha1.ManualApprovalGate, error) {
	ks, err := clients.Get(context.TODO(), names.ManualApprovalGate, metav1.GetOptions{})
	err = wait.PollUntilContextTimeout(context.TODO(), config.APIRetry, config.APITimeout, false, func(context.Context) (bool, error) {
		ks, err = clients.Get(context.TODO(), names.ManualApprovalGate, metav1.GetOptions{})
		if err != nil {
			if apierrs.IsNotFound(err) {
				log.Printf("Waiting for availability of manual approval gate cr [%s]\n", names.ManualApprovalGate)
				return false, nil
			}
			return false, err
		}
		return true, nil
	})
	return ks, err
}

func StartApprovalGatePipeline() {
	cmd.MustSucceed("tkn", "pipeline", "start", "manual-approval-pipeline")
	log.Println("Waiting 10sec to start the pipeline")
	cmd.MustSuccedIncreasedTimeout(time.Second*130, "sleep", "10")
}

func GetApprovaltasklist() string {
	output := cmd.MustSucceed("opc", "approvaltask", "list").Stdout()
	tasklist := strings.Split(output, "\n")
	var taskname string
	for _, line := range tasklist {
		if strings.Contains(line, "manual-approval-pipeline") {
			fields := strings.Fields(line)
			if len(fields) >= 1 {
				taskname = fields[0]
				break
			}
		}
	}
	if len(taskname) == 0 {
		testsuit.T.Fail(fmt.Errorf("Manual approval gate Pipeline doesn't exists\n"))
	}
	return taskname
}
func ApproveApprovalGatePipeline(taskname string) {
	cmd.MustSucceed("opc", "approvaltask", "approve", taskname).Stdout()
}

func RejectApprovalGatePipeline(taskname string) {
	cmd.MustSucceed("opc", "approvaltask", "reject", taskname).Stdout()
}
