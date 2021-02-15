package pipelines

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/assert"
	"github.com/openshift-pipelines/release-tests/pkg/clients"
	"github.com/openshift-pipelines/release-tests/pkg/wait"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ValidateTaskRun(c *clients.Clients, trname, status, namespace string) {
	matched_trname := getTaskRunNameMatches(c, trname, namespace)
	if matched_trname == "" {
		testsuit.T.Errorf("Error: Nothing matched with Taskrun name: %s in namespace %s", trname, namespace)
	}
	// Verify status of TaskRun (wait for it)
	switch {
	case strings.Contains(strings.ToLower(status), "success"):
		validateTaskRunForSuccessStatus(c, matched_trname, namespace)
	case strings.Contains(strings.ToLower(status), "fail"):
		validateTaskRunForFailedStatus(c, matched_trname, namespace)
	case strings.Contains(strings.ToLower(status), "timeout"):
		validateTaskRunTimeOutFailure(c, matched_trname, namespace)
	default:
		testsuit.T.Errorf("Error: %s ", "Not valid input")
	}
}

func validateTaskRunForFailedStatus(c *clients.Clients, trname, namespace string) {
	var err error
	log.Printf("Waiting for TaskRun %s in namespace %s to fail", trname, namespace)
	err = wait.WaitForTaskRunState(c, trname, wait.TaskRunFailed(trname), "BuildValidationFailed")
	assert.NoError(err, fmt.Sprintf("Failed to TaskRun: %s", trname))
}

func validateTaskRunForSuccessStatus(c *clients.Clients, trname, namespace string) {
	var err error
	log.Printf("Waiting for TaskRun %s in namespace %s to succeed", trname, namespace)
	err = wait.WaitForTaskRunState(c, trname, wait.TaskRunSucceed(trname), "TaskRunSucceed")
	assert.NoError(err, fmt.Sprintf("Failed to TaskRun: %s", trname))
}

func validateTaskRunTimeOutFailure(c *clients.Clients, trname, namespace string) {
	var err error
	log.Printf("Waiting for TaskRun %s in namespace %s to complete", trname, namespace)
	err = wait.WaitForTaskRunState(c, "run-giraffe", wait.FailedWithReason("TaskRunTimeout", trname), "TaskRunTimeout")
	assert.NoError(err, fmt.Sprintf("Error waiting for TaskRun %s to finish", trname))
}

func getTaskRunNameMatches(c *clients.Clients, trname, namespace string) string {
	trlist, err := c.TaskRunClient.List(metav1.ListOptions{})
	assert.NoError(err, fmt.Sprintf("Error Getting TaskRun list under namespace %s ", namespace))
	var matched_tr string
	match, _ := regexp.Compile(trname + ".*?")
	for _, tr := range trlist.Items {
		if match.MatchString(tr.Name) {
			matched_tr = tr.Name
			break
		}
	}
	return matched_tr
}

func ValidateTaskRunLabelPropogation(c *clients.Clients, trname, namespace string) {
	trlist, err := c.TaskRunClient.List(metav1.ListOptions{})
	assert.NoError(err, fmt.Sprintf("Error Getting TaskRun list under namespace %s ", namespace))
	var matched_tr string
	match, _ := regexp.Compile(trname + ".*?")
	for _, tr := range trlist.Items {
		if match.MatchString(tr.Name) {
			matched_tr = tr.Name
			break
		}
	}
	labels := make(map[string]string)

	tr, err := c.TaskRunClient.Get(matched_tr, metav1.GetOptions{})
	assert.NoError(err, fmt.Sprintf("Error Getting TaskRun %s under namespace %s ", matched_tr, namespace))

	for key, val := range tr.ObjectMeta.Labels {
		labels[key] = val
	}

	AssertLabelsMatch(labels, tr.ObjectMeta.Labels)
	if tr.Status.PodName != "" {
		pod := GetPodForTaskRun(c, namespace, tr)
		// This label is added to every Pod by the TaskRun controller
		labels[pipeline.GroupName+pipeline.TaskRunLabelKey] = tr.Name
		AssertLabelsMatch(labels, pod.ObjectMeta.Labels)
		gauge.WriteMessage("Labels: \n\n %+v", createKeyValuePairs(labels))
	}
}
