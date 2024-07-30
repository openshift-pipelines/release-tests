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

func ListApprovalTask(cs *clients.Clients) ([]ApprovalTaskInfo, error) {
	var tasks []ApprovalTaskInfo

	err := wait.PollUntilContextTimeout(cs.Ctx, config.APIRetry, config.APITimeout, false, func(ctx context.Context) (bool, error) {
		at, err := cs.ApprovalTask.List(ctx, metav1.ListOptions{})
		if err != nil {
			log.Printf("Failed to list approval tasks, retrying...: %v", err)
			return false, err
		}

		if len(at.Items) == 0 {
			log.Printf("No approval tasks found, retrying...")
			return false, nil
		}

		tasks = make([]ApprovalTaskInfo, 0, len(at.Items))
		for _, item := range at.Items {
			info := ApprovalTaskInfo{
				Name:                      item.Name,
				NumberOfApprovalsRequired: item.Spec.NumberOfApprovalsRequired,
				PendingApprovals:          item.Spec.NumberOfApprovalsRequired - len(item.Status.ApproversResponse),
				Status:                    item.Status.State,
			}
			tasks = append(tasks, info)
		}

		return true, nil
	})

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func ValidateApprovalGatePipeline(expectedStatus string) (bool, error) {
	tasks, err := ListApprovalTask(store.Clients())
	if err != nil {
		return false, fmt.Errorf("error fetching approval tasks: %v", err)
	}

	for _, task := range tasks {
		actualStatus := checkApprovalTaskStatus(task)
		if actualStatus == expectedStatus {
			return true, nil
		}
	}

	return false, errors.New("no approval tasks were found in the specified state")
}

func checkApprovalTaskStatus(task ApprovalTaskInfo) string {
	switch task.Status {
	case "pending":
		return "Pending"
	case "rejected":
		return "Rejected"
	case "approved":
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
