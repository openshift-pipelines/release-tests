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
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/getgauge-contrib/gauge-go/testsuit"
	atv1alpha1 "github.com/openshift-pipelines/manual-approval-gate/pkg/apis/approvaltask/v1alpha1"
	"github.com/openshift-pipelines/release-tests/pkg/clients"
	"github.com/openshift-pipelines/release-tests/pkg/cmd"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	"github.com/openshift-pipelines/release-tests/pkg/store"
	operatorv1alpha1 "github.com/tektoncd/operator/pkg/apis/operator/v1alpha1"
	mag "github.com/tektoncd/operator/pkg/client/clientset/versioned/typed/operator/v1alpha1"
	"github.com/tektoncd/operator/test/utils"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

type ApprovalTaskInfo struct {
	Name   string
	Status string
}

func EnsureManualApprovalGateExists(clients mag.ManualApprovalGateInterface, names utils.ResourceNames) (*operatorv1alpha1.ManualApprovalGate, error) {
	var magCR *operatorv1alpha1.ManualApprovalGate

	err := wait.PollUntilContextTimeout(context.TODO(), config.APIRetry, config.APITimeout, false, func(ctx context.Context) (bool, error) {
		cr, err := clients.Get(ctx, names.ManualApprovalGate, metav1.GetOptions{})
		if err != nil {
			if apierrs.IsNotFound(err) {
				log.Printf("Waiting for availability of manual approval gate cr [%s]\n", names.ManualApprovalGate)
				return false, nil
			}
			return false, err
		}
		magCR = cr
		return true, nil
	})

	return magCR, err
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
			tasks = append(tasks, ApprovalTaskInfo{
				Name:   item.Name,
				Status: item.Status.State,
			})
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

var (
	magAPIServerOnce sync.Once
	magAPIServer     string
	magAPIServerErr  error

	magUserKubeconfigsMu sync.Mutex
	magUserKubeconfigs   = map[string]string{}

	magUserAuthDirtyMu sync.Mutex
	magUserAuthDirty   = map[string]bool{}
)

func markUsersAuthDirty(users []string) {
	magUserAuthDirtyMu.Lock()
	defer magUserAuthDirtyMu.Unlock()
	for _, u := range users {
		u = strings.TrimSpace(u)
		if u == "" {
			continue
		}
		magUserAuthDirty[u] = true
	}
}

func popUserAuthDirty(user string) bool {
	magUserAuthDirtyMu.Lock()
	defer magUserAuthDirtyMu.Unlock()
	dirty := magUserAuthDirty[user]
	if dirty {
		magUserAuthDirty[user] = false
	}
	return dirty
}

func userPassword(user string) string {
	envVar := strings.ToUpper(user) + "_PASS"
	if v := strings.TrimSpace(os.Getenv(envVar)); v != "" {
		return v
	}
	// default: password == username
	return user
}

// EnsureGroupMembers ensures the OpenShift Group exists with exactly the provided users.
// This uses the current kubeconfig (admin context) and is intentionally cluster-scoped.
func EnsureGroupMembers(group string, users []string) {
	if group == "" {
		testsuit.T.Fail(fmt.Errorf("group name is empty"))
	}

	// Capture old membership so we can mark removed users dirty too.
	oldUsers := []string{}
	getUsersCmd := cmd.Run("oc", "get", "group", group, "-o", "jsonpath={.users[*]}")
	groupExists := getUsersCmd.ExitCode == 0
	if groupExists {
		out := strings.TrimSpace(getUsersCmd.Stdout())
		if out != "" {
			oldUsers = strings.Fields(out)
		}
	}

	usersJSON, err := json.Marshal(users)
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("failed to marshal group users: %v", err))
	}
	patch := fmt.Sprintf("{\"users\":%s}", string(usersJSON))

	// Ensure the group exists in an idempotent way (avoid races with other processes).
	// Try creating the group; ignore if it already exists.
	createRes := cmd.Run("oc", "adm", "groups", "new", group)
	if createRes.ExitCode != 0 {
		stderr := strings.ToLower(createRes.Stderr())
		if !strings.Contains(stderr, "already exists") && !strings.Contains(stderr, "alreadyexists") {
			testsuit.T.Fail(fmt.Errorf("failed to create group %s: %s", group, createRes.Stderr()))
		}
	}

	// Patch the group and capture the resulting user list in a single call.
	out := strings.TrimSpace(cmd.MustSucceed("oc", "patch", "group", group, "--type=merge", "-p", patch, "-o", "jsonpath={.users[*]}").Stdout())
	actual := []string{}
	if out != "" {
		actual = strings.Fields(out)
	}

	expected := append([]string{}, users...)
	sort.Strings(expected)
	sort.Strings(actual)
	if strings.Join(expected, ",") != strings.Join(actual, ",") {
		testsuit.T.Fail(fmt.Errorf("group %s membership mismatch: expected [%s], got [%s]", group, strings.Join(expected, " "), strings.Join(actual, " ")))
	}

	// OpenShift group membership is reflected in user tokens; refresh user auth after membership changes.
	// Mark both old and current members dirty so removed users also get a fresh token on next action.
	oldSorted := append([]string{}, oldUsers...)
	sort.Strings(oldSorted)
	if !groupExists || strings.Join(oldSorted, ",") != strings.Join(actual, ",") {
		union := append([]string{}, oldUsers...)
		union = append(union, actual...)
		markUsersAuthDirty(union)
	}
}

func CreateApprovalPipelineRun(id, description string, approvers []string, required int, timeout, namespace string) (string, string) {
	if len(approvers) == 0 {
		testsuit.T.Fail(fmt.Errorf("approvers list is empty"))
	}
	if required <= 0 {
		testsuit.T.Fail(fmt.Errorf("numberOfApprovalsRequired must be > 0; got %d", required))
	}

	approverLines := strings.Builder{}
	for _, a := range approvers {
		a = strings.TrimSpace(a)
		if a == "" {
			continue
		}
		approverLines.WriteString("- ")
		approverLines.WriteString(a)
		approverLines.WriteString("\n")
	}
	approversRaw := strings.TrimSuffix(approverLines.String(), "\n")
	if strings.TrimSpace(approversRaw) == "" {
		testsuit.T.Fail(fmt.Errorf("approvers list is empty after trimming"))
	}
	// Indent subsequent list items (the first line is indented by the YAML template itself).
	approversYAML := strings.ReplaceAll(approversRaw, "\n", "\n              ")

	idLower := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(id), " ", "-"))
	idLower = strings.ReplaceAll(idLower, "_", "-")
	genName := fmt.Sprintf("approva-grp-plr-%s-", idLower)

	prYAML := fmt.Sprintf(`apiVersion: tekton.dev/v1
kind: PipelineRun
metadata:
  generateName: %s
  namespace: %s
spec:
  pipelineSpec:
    tasks:
      - name: wait
        timeout: %s
        taskRef:
          apiVersion: openshift-pipelines.org/v1alpha1
          kind: ApprovalTask
        params:
          - name: approvers
            value:
              %s
          - name: numberOfApprovalsRequired
            value: "%d"
          - name: description
            value: "%s"
`, genName, namespace, timeout, approversYAML, required, description)

	prName := strings.TrimSpace(cmd.MustSucceedWithStdin(strings.NewReader(prYAML), "oc", "create", "-n", namespace, "-f", "-", "-o", "jsonpath={.metadata.name}").Stdout())
	if prName == "" {
		testsuit.T.Fail(fmt.Errorf("failed to create PipelineRun: got empty name"))
	}

	taskName := WaitForSingleApprovalTaskName(prName, namespace, 2*time.Minute)
	log.Printf("[MAG users] %s: created PipelineRun=%s ApprovalTask=%s", id, prName, taskName)
	return prName, taskName
}

func WaitForSingleApprovalTaskName(prName, namespace string, timeout time.Duration) string {
	cs := store.Clients()
	if cs == nil {
		testsuit.T.Fail(fmt.Errorf("clients not initialized"))
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var picked *atv1alpha1.ApprovalTask
	err := wait.PollUntilContextTimeout(ctx, time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		list, err := cs.ApprovalTask.List(ctx, metav1.ListOptions{})
		if err != nil {
			return false, err
		}
		if len(list.Items) == 0 {
			return false, nil
		}

		// In per-scenario namespaces we expect a single task, but if multiple exist pick the newest.
		sort.SliceStable(list.Items, func(i, j int) bool {
			return list.Items[i].CreationTimestamp.After(list.Items[j].CreationTimestamp.Time)
		})
		picked = &list.Items[0]
		return true, nil
	})
	if err != nil || picked == nil {
		testsuit.T.Fail(fmt.Errorf("timed out waiting for ApprovalTask for pipelinerun %s in namespace %s: %v", prName, namespace, err))
	}
	return picked.Name
}

func WaitForApprovalTaskState(task, expectedState string, timeout time.Duration) {
	cs := store.Clients()
	exp := strings.ToLower(expectedState)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := wait.PollUntilContextTimeout(ctx, time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		at, err := cs.ApprovalTask.Get(ctx, task, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		return strings.ToLower(at.Status.State) == exp, nil
	})
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("timed out waiting for ApprovalTask %s state=%s: %v", task, expectedState, err))
	}
}

func approvalTaskMessageContains(at *atv1alpha1.ApprovalTask, text string) bool {
	if at == nil || text == "" {
		return false
	}

	// Spec-level message (updated immediately by the approve/reject CLI).
	for _, a := range at.Spec.Approvers {
		if strings.Contains(a.Message, text) {
			return true
		}
	}

	// Status-level message (may be populated asynchronously by the controller).
	for _, r := range at.Status.ApproversResponse {
		if strings.Contains(r.Message, text) {
			return true
		}
		for _, m := range r.GroupMembers {
			if strings.Contains(m.Message, text) {
				return true
			}
		}
	}

	return false
}

func WaitForApprovalTaskMessageContains(task, text string, timeout time.Duration) {
	cs := store.Clients()
	if cs == nil {
		testsuit.T.Fail(fmt.Errorf("clients not initialized"))
		return
	}

	text = strings.TrimSpace(text)
	if text == "" {
		testsuit.T.Fail(fmt.Errorf("message text is empty"))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var last *atv1alpha1.ApprovalTask
	err := wait.PollUntilContextTimeout(ctx, time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		at, err := cs.ApprovalTask.Get(ctx, task, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		last = at
		return approvalTaskMessageContains(at, text), nil
	})
	if err != nil {
		if last == nil {
			testsuit.T.Fail(fmt.Errorf("timed out waiting for ApprovalTask %s to contain message %q: %v", task, text, err))
			return
		}
		testsuit.T.Fail(fmt.Errorf("timed out waiting for ApprovalTask %s to contain message %q; last status=%q", task, text, last.Status.State))
	}
}

func WaitForAndAssertApprovalTaskListState(task string, expectedNum, expectedPending, expectedRejected int, expectedStatus string, timeout time.Duration) {
	cs := store.Clients()
	expStatus := expectedStatus

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var last *atv1alpha1.ApprovalTask
	err := wait.PollUntilContextTimeout(ctx, time.Second, timeout, true, func(ctx context.Context) (bool, error) {
		at, err := cs.ApprovalTask.Get(ctx, task, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		last = at

		num := at.Spec.NumberOfApprovalsRequired
		pending := pendingApprovals(at)
		rejected := rejectedCount(at)
		status := stateHuman(at)

		return num == expectedNum && pending == expectedPending && rejected == expectedRejected && status == expStatus, nil
	})
	if err != nil {
		if last == nil {
			testsuit.T.Fail(fmt.Errorf("failed to read ApprovalTask %s for state assertion: %v", task, err))
		}
		testsuit.T.Fail(fmt.Errorf("approvaltask %s list-state mismatch: expected num=%d pending=%d rejected=%d status=%s; got num=%d pending=%d rejected=%d status=%s",
			task,
			expectedNum, expectedPending, expectedRejected, expStatus,
			last.Spec.NumberOfApprovalsRequired, pendingApprovals(last), rejectedCount(last), stateHuman(last),
		))
	}
}

// pendingApprovals matches the CLI calculation (see manual-approval-gate/pkg/cli/cmd/list).
func pendingApprovals(at *atv1alpha1.ApprovalTask) int {
	respondedUsers := make(map[string]bool)

	for _, approver := range at.Status.ApproversResponse {
		switch atv1alpha1.DefaultedApproverType(approver.Type) {
		case "User":
			respondedUsers[approver.Name] = true
		case "Group":
			for _, member := range approver.GroupMembers {
				if member.Response == "approved" || member.Response == "rejected" {
					respondedUsers[member.Name] = true
				}
			}
		}
	}

	return at.Spec.NumberOfApprovalsRequired - len(respondedUsers)
}

// rejectedCount matches the CLI calculation (see manual-approval-gate/pkg/cli/cmd/list).
func rejectedCount(at *atv1alpha1.ApprovalTask) int {
	count := 0
	rejectedUsers := make(map[string]bool)

	for _, approver := range at.Status.ApproversResponse {
		if atv1alpha1.DefaultedApproverType(approver.Type) == "User" && approver.Response == "rejected" {
			if !rejectedUsers[approver.Name] {
				rejectedUsers[approver.Name] = true
				count++
			}
		} else if atv1alpha1.DefaultedApproverType(approver.Type) == "Group" {
			for _, member := range approver.GroupMembers {
				if member.Response == "rejected" {
					if !rejectedUsers[member.Name] {
						rejectedUsers[member.Name] = true
						count++
					}
				}
			}
		}
	}

	return count
}

func stateHuman(at *atv1alpha1.ApprovalTask) string {
	switch at.Status.State {
	case "approved":
		return "Approved"
	case "rejected":
		return "Rejected"
	case "pending":
		return "Pending"
	default:
		return at.Status.State
	}
}

// MAGGroupName returns a unique group name for a given alias in a scenario namespace.
// This avoids collisions since OpenShift Groups are cluster-scoped.
func MAGGroupName(namespace, alias string) string {
	a := strings.TrimSpace(alias)
	if a == "" {
		return ""
	}

	// If the caller passes a literal group name (e.g. "system:authenticated"), keep it as-is.
	if strings.Contains(a, ":") {
		return a
	}

	a = strings.ToLower(a)
	a = strings.ReplaceAll(a, "_", "-")
	a = strings.ReplaceAll(a, " ", "-")

	// Keep only [a-z0-9-] to avoid invalid group names; collapse other chars into '-'.
	b := strings.Builder{}
	b.Grow(len(a))
	lastDash := false
	for _, r := range a {
		isAZ := r >= 'a' && r <= 'z'
		is09 := r >= '0' && r <= '9'
		if isAZ || is09 {
			b.WriteRune(r)
			lastDash = false
			continue
		}
		if r == '-' {
			if !lastDash {
				b.WriteRune('-')
				lastDash = true
			}
			continue
		}
		if !lastDash {
			b.WriteRune('-')
			lastDash = true
		}
	}
	safeAlias := strings.Trim(b.String(), "-")
	if safeAlias == "" {
		safeAlias = "group"
	}
	return fmt.Sprintf("mag-%s-%s", namespace, safeAlias)
}

func ensureMAGAPIServer() string {
	magAPIServerOnce.Do(func() {
		api := strings.TrimSpace(cmd.MustSucceed("oc", "whoami", "--show-server").Stdout())
		if api == "" {
			magAPIServerErr = fmt.Errorf("failed to detect cluster API server via `oc whoami --show-server`")
			return
		}
		magAPIServer = api
	})
	if magAPIServerErr != nil {
		testsuit.T.Fail(magAPIServerErr)
		return ""
	}
	return magAPIServer
}

func ensureUserKubeconfig(user string) string {
	magUserKubeconfigsMu.Lock()
	if v, ok := magUserKubeconfigs[user]; ok && strings.TrimSpace(v) != "" {
		magUserKubeconfigsMu.Unlock()
		// If group membership changed since the last login, refresh the token so group-based approvers work.
		if popUserAuthDirty(user) {
			apiServer := ensureMAGAPIServer()
			pass := userPassword(user)
			cmd.MustSucceed("oc", "login", apiServer, "-u", user, "-p", pass, "--kubeconfig", v, "--insecure-skip-tls-verify=true")
		}
		return v
	}
	magUserKubeconfigsMu.Unlock()

	apiServer := ensureMAGAPIServer()
	pass := userPassword(user)

	tmp, err := os.CreateTemp("", fmt.Sprintf("mag-kubeconfig-%s-", user))
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("failed to create temp kubeconfig for %s: %v", user, err))
	}
	_ = tmp.Close()

	kcPath := tmp.Name()
	cmd.MustSucceed("oc", "login", apiServer, "-u", user, "-p", pass, "--kubeconfig", kcPath, "--insecure-skip-tls-verify=true")

	magUserKubeconfigsMu.Lock()
	magUserKubeconfigs[user] = kcPath
	magUserKubeconfigsMu.Unlock()
	// Fresh login done; clear dirty flag if it was set.
	_ = popUserAuthDirty(user)
	return kcPath
}

// CleanupUserKubeconfigs removes any temp kubeconfig files created for per-user logins.
// It is safe to ignore errors during cleanup.
func CleanupUserKubeconfigs() {
	magUserKubeconfigsMu.Lock()
	defer magUserKubeconfigsMu.Unlock()

	for user, path := range magUserKubeconfigs {
		if strings.TrimSpace(path) != "" {
			_ = os.Remove(path)
		}
		delete(magUserKubeconfigs, user)
	}

	magUserAuthDirtyMu.Lock()
	magUserAuthDirty = map[string]bool{}
	magUserAuthDirtyMu.Unlock()
}

func ApproveApprovalTaskAsUser(user, task, namespace, message string) {
	kc := ensureUserKubeconfig(user)
	args := []string{"opc", "approvaltask", "approve", task, "-n", namespace}
	if strings.TrimSpace(message) != "" {
		args = append(args, "-m", message)
	}
	cmd.MustSucceedWithEnv([]string{"KUBECONFIG=" + kc}, args...)
}

func RejectApprovalTaskAsUser(user, task, namespace, message string) {
	kc := ensureUserKubeconfig(user)
	args := []string{"opc", "approvaltask", "reject", task, "-n", namespace}
	if strings.TrimSpace(message) != "" {
		args = append(args, "-m", message)
	}
	cmd.MustSucceedWithEnv([]string{"KUBECONFIG=" + kc}, args...)
}

func ApproveApprovalTaskExpectFailAsUser(user, task, namespace, message string) {
	kc := ensureUserKubeconfig(user)
	args := []string{"opc", "approvaltask", "approve", task, "-n", namespace}
	if strings.TrimSpace(message) != "" {
		args = append(args, "-m", message)
	}
	res := cmd.RunWithEnv([]string{"KUBECONFIG=" + kc}, args...)
	if res.ExitCode == 0 {
		testsuit.T.Fail(fmt.Errorf("expected approval by %s on %s to fail, but it succeeded", user, task))
	}
}

func ApproveApprovalTaskAllowFinalStateAsUser(user, task, namespace, message string) {
	kc := ensureUserKubeconfig(user)
	args := []string{"opc", "approvaltask", "approve", task, "-n", namespace}
	if strings.TrimSpace(message) != "" {
		args = append(args, "-m", message)
	}
	res := cmd.RunWithEnv([]string{"KUBECONFIG=" + kc}, args...)
	if res.ExitCode == 0 {
		return
	}
	out := strings.ToLower(res.Stdout() + "\n" + res.Stderr())
	if strings.Contains(out, "already reached") && strings.Contains(out, "final state") {
		return
	}
	testsuit.T.Fail(fmt.Errorf("unexpected approval failure for %s on %s: %s", user, task, res.Stderr()))
}

type approvalTaskActionFn func(user, task, namespace, message string)

var approvalTaskActionDispatch = map[string]approvalTaskActionFn{
	// Keep actions strict and spec-driven (one canonical key per action).
	"approve":                   ApproveApprovalTaskAsUser,
	"reject":                    RejectApprovalTaskAsUser,
	"approve-expect-fail":       ApproveApprovalTaskExpectFailAsUser,
	"approve-allow-final-state": ApproveApprovalTaskAllowFinalStateAsUser,
}

func PerformApprovalTaskActionAsUser(user, action, task, namespace, message string) {
	a := strings.ToLower(strings.TrimSpace(action))
	if a == "" {
		testsuit.T.Fail(fmt.Errorf("approval task action is empty"))
		return
	}

	fn, ok := approvalTaskActionDispatch[a]
	if !ok {
		testsuit.T.Fail(fmt.Errorf("unsupported approval gate action: %s", action))
		return
	}

	fn(user, task, namespace, message)
}
