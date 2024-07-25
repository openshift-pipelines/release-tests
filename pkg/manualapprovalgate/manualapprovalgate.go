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

package approvalgate

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/openshift-pipelines/release-tests/pkg/clients"
	"github.com/openshift-pipelines/release-tests/pkg/cmd"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	"github.com/openshift-pipelines/release-tests/pkg/store"
	"github.com/tektoncd/operator/pkg/apis/operator/v1alpha1"
	mag "github.com/tektoncd/operator/pkg/client/clientset/versioned/typed/operator/v1alpha1"
	"github.com/tektoncd/operator/test/utils"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

type ApprovalTaskInfo struct {
	Name                      string
	NumberOfApprovalsRequired int
	PendingApprovals          int
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

func ListApprovalTask(cs *clients.Clients) []ApprovalTaskInfo {
	var tasks []ApprovalTaskInfo
	at, err := cs.ApprovalTask.List(cs.Ctx, v1.ListOptions{})
	if err != nil {
		fmt.Errorf("Failed to List approval task")
		return tasks
	}

	for _, item := range at.Items {
		info := ApprovalTaskInfo{
			Name:                      item.Name,
			NumberOfApprovalsRequired: item.Spec.NumberOfApprovalsRequired,
			PendingApprovals:          len(item.Spec.Approvers),
			Status:                    item.Status.State,
		}
		tasks = append(tasks, info)
	}

	return tasks
}

func ValidateApprovalGatePipeline(expectedStatus string) (bool, error) {
	tasks := ListApprovalTask(store.Clients())
	if len(tasks) == 0 {
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

func checkApprovalTaskStatus(task ApprovalTaskInfo) string {
	switch {
	case task.Status == "pending":
		return "Pending"
	case task.Status == "rejected":
		return "Rejected"
	case task.Status == "approved":
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
