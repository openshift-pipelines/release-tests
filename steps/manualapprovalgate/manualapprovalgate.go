package approvalgate

import (
	"errors"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/getgauge-contrib/gauge-go/gauge_messages"
	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/cmd"
	approvalgate "github.com/openshift-pipelines/release-tests/pkg/manualapprovalgate"
	"github.com/openshift-pipelines/release-tests/pkg/opc"
	"github.com/openshift-pipelines/release-tests/pkg/store"
)

func splitList(s string) []string {
	raw := strings.TrimSpace(s)
	if raw == "" || raw == "-" {
		return []string{}
	}
	parts := strings.FieldsFunc(raw, func(r rune) bool {
		return r == ',' || r == ' ' || r == '\t' || r == '\n' || r == '\r'
	})
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		out = append(out, p)
	}
	return out
}

func resolveMAGGroupName(alias string) string {
	ns := store.Namespace()
	// If alias is a literal group name (e.g. "system:authenticated"), keep it.
	if strings.Contains(strings.TrimSpace(alias), ":") || strings.HasPrefix(strings.TrimSpace(alias), "mag-") {
		return strings.TrimSpace(alias)
	}
	return approvalgate.MAGGroupName(ns, alias)
}

func resolveApprovalTimeout(timeout string) string {
	t := strings.TrimSpace(strings.ToLower(timeout))
	switch t {
	case "success":
		if v := strings.TrimSpace(os.Getenv("APPROVAL_TIMEOUT_SUCCESS")); v != "" {
			return v
		}
		// Safe default for success cases (allows enough time for user actions on slower clusters).
		return "5m"
	case "fail-fast":
		if v := strings.TrimSpace(os.Getenv("APPROVAL_TIMEOUT_FAIL_FAST")); v != "" {
			return v
		}
		// Default for negative cases where we expect quick completion (reject / logic).
		return "2m"
	case "timeout":
		if v := strings.TrimSpace(os.Getenv("APPROVAL_TIMEOUT_TIMEOUT")); v != "" {
			return v
		}
		// Default for timeout-based negative cases (keeps the suite fast).
		return "30s"
	default:
		return strings.TrimSpace(timeout)
	}
}

func getCurrentApprovalTask() string {
	if raw, ok := gauge.GetScenarioStore()["mag.approvaltask"]; ok {
		if s, ok := raw.(string); ok && strings.TrimSpace(s) != "" {
			return strings.TrimSpace(s)
		}
	}
	testsuit.T.Fail(errors.New("no current approval task found in scenario store (expected key: mag.approvaltask)"))
	return ""
}

func currentMAGCaseID() string {
	raw, ok := gauge.GetScenarioStore()["scenario.name"]
	if !ok {
		testsuit.T.Fail(errors.New("scenario.name not found in scenario store"))
		return ""
	}

	name, ok := raw.(string)
	if !ok || strings.TrimSpace(name) == "" {
		testsuit.T.Fail(errors.New("scenario.name is empty or not a string"))
		return ""
	}

	re := regexp.MustCompile(`(?i)\bTC-?(\d{1,3})\b`)
	m := re.FindStringSubmatch(name)
	if len(m) < 2 {
		testsuit.T.Fail(errors.New("could not derive testcase id from scenario name: " + name))
		return ""
	}

	id := m[1]
	if len(id) == 1 {
		id = "0" + id
	}
	return id
}

func addMAGGroupForCleanup(groupName string) {
	groupName = strings.TrimSpace(groupName)
	if groupName == "" {
		return
	}
	storeMap := gauge.GetScenarioStore()
	existing, ok := storeMap["mag.groups"].([]string)
	if !ok {
		storeMap["mag.groups"] = []string{groupName}
		return
	}
	for _, g := range existing {
		if g == groupName {
			return
		}
	}
	storeMap["mag.groups"] = append(existing, groupName)
}

var _ = gauge.Step("Start the <pipelineName> pipeline with workspace <workspaceValue>", func(pipelineName, workspaceValue string) {
	params := make(map[string]string)
	workspaces := make(map[string]string)
	workspaces[strings.Split(workspaceValue, ",")[0]] = strings.Split(workspaceValue, ",")[1]
	opc.StartPipeline(pipelineName, params, workspaces, store.Namespace(), "--use-param-defaults")
})

var _ = gauge.Step("Approve the manual-approval-pipeline", func() {
	tasks, err := approvalgate.ListApprovalTask(store.Clients())
	if err != nil {
		testsuit.T.Errorf("Error while listing approval gate tasks: %v", err)
		return
	}

	for _, task := range tasks {
		approvalgate.ApproveApprovalGatePipeline(task.Name)
	}
})

var _ = gauge.Step("Reject the manual-approval-pipeline", func() {
	tasks, err := approvalgate.ListApprovalTask(store.Clients())
	if err != nil {
		testsuit.T.Errorf("Error while listing approval gate tasks: %v", err)
		return
	}

	for _, task := range tasks {
		approvalgate.RejectApprovalGatePipeline(task.Name)
	}
})

var _ = gauge.Step("Validate the manual-approval-pipeline for <status> state", func(status string) {
	success, err := approvalgate.ValidateApprovalGatePipeline(status)
	if err != nil {
		testsuit.T.Fail(err)
		return
	}

	if !success {
		testsuit.T.Fail(errors.New("validation failed: no approvaltasks matched the expected status"))
	}
})

var _ = gauge.Step("Ensure approval group <groupAlias> has members <members>", func(groupAlias, members string) {
	groupName := resolveMAGGroupName(groupAlias)
	addMAGGroupForCleanup(groupName)
	approvalgate.EnsureGroupMembers(groupName, splitList(members))
})

// Preferred step: testcase id is derived from the scenario name (e.g. PIPELINES-28-TC01).
var _ = gauge.Step("Create manual approval gate pipelinerun with approvers <approvers> required <required> Should <timeout>", func(approvers, required, timeout string) {
	tcID := currentMAGCaseID()
	if tcID == "" {
		return
	}

	desc, _ := gauge.GetScenarioStore()["scenario.name"].(string)
	desc = strings.TrimSpace(desc)
	if desc == "" {
		desc = "manual approval gate users"
	}

	requiredInt, err := strconv.Atoi(strings.TrimSpace(required))
	if err != nil {
		testsuit.T.Fail(err)
		return
	}

	timeoutVal := resolveApprovalTimeout(timeout)
	if timeoutVal == "" {
		// Fall back to a safe default to avoid empty timeout in YAML.
		timeoutVal = "5m"
	}

	rawApprovers := splitList(approvers)
	finalApprovers := make([]string, 0, len(rawApprovers))
	for _, a := range rawApprovers {
		if strings.HasPrefix(a, "group:") {
			alias := strings.TrimPrefix(a, "group:")
			groupName := resolveMAGGroupName(alias)
			finalApprovers = append(finalApprovers, "group:"+groupName)
			continue
		}
		finalApprovers = append(finalApprovers, a)
	}

	pr, task := approvalgate.CreateApprovalPipelineRun(tcID, desc, finalApprovers, requiredInt, timeoutVal, store.Namespace())
	store.PutScenarioData("mag.pipelinerun", pr)
	store.PutScenarioData("mag.approvaltask", task)
})

var _ = gauge.Step("User <user> performs <action> on the manual approval gate task", func(user, action string) {
	approvalgate.PerformApprovalTaskActionAsUser(user, action, getCurrentApprovalTask(), store.Namespace(), "")
})

var _ = gauge.Step("User <user> performs <action> on the manual approval gate task with message <message>", func(user, action, message string) {
	approvalgate.PerformApprovalTaskActionAsUser(user, action, getCurrentApprovalTask(), store.Namespace(), message)
})

var _ = gauge.Step("Validate manual approval gate task for <status> state", func(status string) {
	approvalgate.WaitForApprovalTaskState(getCurrentApprovalTask(), status, 2*time.Minute)
})

var _ = gauge.Step("Validate manual approval gate task list state numberOfApprovalsRequired <num> pending <pending> rejected <rejected> status <status>", func(num, pending, rejected, status string) {
	numInt, err := strconv.Atoi(strings.TrimSpace(num))
	if err != nil {
		testsuit.T.Fail(err)
		return
	}
	pendingInt, err := strconv.Atoi(strings.TrimSpace(pending))
	if err != nil {
		testsuit.T.Fail(err)
		return
	}
	rejectedInt, err := strconv.Atoi(strings.TrimSpace(rejected))
	if err != nil {
		testsuit.T.Fail(err)
		return
	}
	approvalgate.WaitForAndAssertApprovalTaskListState(getCurrentApprovalTask(), numInt, pendingInt, rejectedInt, status, 2*time.Minute)
})

var _ = gauge.Step("Verify manual approval gate task message contains <text>", func(text string) {
	approvalgate.WaitForApprovalTaskMessageContains(getCurrentApprovalTask(), text, 60*time.Second)
})

// Cleanup for approval-gate user/group test scenarios.
// - Deletes any groups created for the scenario.
var _ = gauge.AfterScenario(func(exInfo *gauge_messages.ExecutionInfo) {
	scenarioStore := gauge.GetScenarioStore()

	if raw, ok := scenarioStore["mag.groups"]; ok {
		if groups, ok := raw.([]string); ok {
			for _, g := range groups {
				g = strings.TrimSpace(g)
				if g == "" {
					continue
				}
				cmd.Run("oc", "delete", "group", g, "--ignore-not-found")
			}
		}
	}
}, []string{"approvalgate-users"}, testsuit.AND)

// Cleanup any temp kubeconfigs created for per-user logins after the full Gauge suite execution.
var _ = gauge.AfterSuite(func(exInfo *gauge_messages.ExecutionInfo) {
	approvalgate.CleanupUserKubeconfigs()
}, []string{}, testsuit.AND)
