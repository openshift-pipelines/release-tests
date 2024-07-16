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
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/openshift-pipelines/release-tests/pkg/cmd"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	"github.com/tektoncd/operator/pkg/apis/operator/v1alpha1"
	mag "github.com/tektoncd/operator/pkg/client/clientset/versioned/typed/operator/v1alpha1"
	"github.com/tektoncd/operator/test/utils"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

type TaskInfo struct {
	Name                      string
	NumberOfApprovalsRequired int
	PendingApprovals          int
	Rejected                  int
	Status                    string
}

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

func GetApprovalTaskList() []TaskInfo {
	output := cmd.MustSucceed("opc", "approvaltask", "list").Stdout()
	tasklist := strings.Split(output, "\n")
	var tasks []TaskInfo

	headers := strings.Fields(tasklist[0])
	if len(headers) < 5 {
		return nil
	}

	for _, line := range tasklist[1:] {
		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}

		requiredApprovals, err := strconv.Atoi(fields[1])
		if err != nil {
			continue
		}

		pendingApprovals, err := strconv.Atoi(fields[2])
		if err != nil {
			continue
		}

		rejected, err := strconv.Atoi(fields[3])
		if err != nil {
			continue
		}

		task := TaskInfo{
			Name:                      fields[0],
			NumberOfApprovalsRequired: requiredApprovals,
			PendingApprovals:          pendingApprovals,
			Rejected:                  rejected,
			Status:                    fields[4],
		}
		tasks = append(tasks, task)
	}

	if len(tasks) == 0 {
		return nil
	}

	return tasks
}

func ValidateApprovalGatePipeline(expectedStatus string) (bool, error) {
	tasks := GetApprovalTaskList()
	if tasks == nil {
		return false, errors.New("no approval gate tasks found")
	}

	found := false
	for _, task := range tasks {
		actualStatus := checkApprovalTaskStatus(task)
		if actualStatus == expectedStatus {
			found = true
			break
		}
	}

	if !found {
		return false, errors.New("no approval tasks were found in the specified state")
	}
	return true, nil
}

func checkApprovalTaskStatus(task TaskInfo) string {
	switch {
	case task.PendingApprovals > 0:
		return "Pending"
	case task.Rejected > 0:
		return "Rejected"
	case task.Status == "Approved" && task.PendingApprovals == 0 && task.Rejected == 0:
		return "Approved"
	default:
		return "Unknown Error: Check Details"
	}
}

func ApproveApprovalGatePipeline(taskname string) {
	cmd.MustSucceed("opc", "approvaltask", "approve", taskname)
}

func RejectApprovalGatePipeline(taskname string) {
	cmd.MustSucceed("opc", "approvaltask", "reject", taskname)
}
