package pipelines

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/clients"
	"github.com/openshift-pipelines/release-tests/pkg/k8s"
	"github.com/openshift-pipelines/release-tests/pkg/wait"
	"github.com/tektoncd/cli/pkg/cli"
	clitr "github.com/tektoncd/cli/pkg/cmd/taskrun"
	"github.com/tektoncd/cli/pkg/options"
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
	if err != nil {
		var printMsg string
		buf, logsErr := getTaskrunLogs(c, trname, namespace)
		events, eventError := k8s.GetWarningEvents(c, namespace)
		if logsErr != nil {
			if eventError != nil {
				printMsg = fmt.Sprintf("task run %s was expected to be in TaskRunFailed state \n%v \ntaskrun logs error: \n%v \ntaskrun events error: \n%v", trname, err, logsErr, eventError)
			} else {
				printMsg = fmt.Sprintf("task run %s was expected to be in TaskRunFailed state \n%v \ntaskrun logs error: \n%v \ntaskrun events: \n%v", trname, err, logsErr, events)
			}
		} else {
			if eventError != nil {
				printMsg = fmt.Sprintf("task run %s was expected to be in TaskRunFailed state \n%v \ntaskrun logs: \n%v \ntaskrun events error: \n%v", trname, err, buf.String(), eventError)
			} else {
				printMsg = fmt.Sprintf("task run %s was expected to be in TaskRunFailed state \n %v \ntaskrun logs: \n%v \ntaskrun events: \n%v", trname, err, buf.String(), events)
			}
		}
		testsuit.T.Errorf(printMsg)
	}
}

func validateTaskRunForSuccessStatus(c *clients.Clients, trname, namespace string) {
	var err error
	log.Printf("Waiting for TaskRun %s in namespace %s to succeed", trname, namespace)
	err = wait.WaitForTaskRunState(c, trname, wait.TaskRunSucceed(trname), "TaskRunSucceed")
	if err != nil {
		var printMsg string
		buf, logsErr := getTaskrunLogs(c, trname, namespace)
		events, eventError := k8s.GetWarningEvents(c, namespace)
		if logsErr != nil {
			if eventError != nil {
				printMsg = fmt.Sprintf("task run %s was expected to be in TaskRunSucceed state \n%v \ntaskrun logs error: \n%v \ntaskrun events error: \n%v", trname, err, logsErr, eventError)
			} else {
				printMsg = fmt.Sprintf("task run %s was expected to be in TaskRunSucceed state \n%v \ntaskrun logs error: \n%v \ntaskrun events: \n%v", trname, err, logsErr, events)
			}
		} else {
			if eventError != nil {
				printMsg = fmt.Sprintf("task run %s was expected to be in TaskRunSucceed state \n%v \ntaskrun logs: \n%v \ntaskrun events error: \n%v", trname, err, buf.String(), eventError)
			} else {
				printMsg = fmt.Sprintf("task run %s was expected to be in TaskRunSucceed state \n%v \ntaskrun logs: \n%v \ntaskrun events: \n%v", trname, err, buf.String(), events)
			}
		}
		testsuit.T.Errorf(printMsg)
	}
}

func validateTaskRunTimeOutFailure(c *clients.Clients, trname, namespace string) {
	var err error
	log.Printf("Waiting for TaskRun %s in namespace %s to complete", trname, namespace)
	err = wait.WaitForTaskRunState(c, "run-giraffe", wait.FailedWithReason("TaskRunTimeout", trname), "TaskRunTimeout")
	if err != nil {
		testsuit.T.Errorf("task run %s was expected to be in TaskRunTimeout state \n %v", trname, err)
	}
}

func getTaskRunNameMatches(c *clients.Clients, trname, namespace string) string {
	trlist, err := c.TaskRunClient.List(c.Ctx, metav1.ListOptions{})
	if err != nil {
		testsuit.T.Errorf("failed to list task runs in namespace %s \n %v", namespace, err)
	}

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
	trlist, err := c.TaskRunClient.List(c.Ctx, metav1.ListOptions{})
	if err != nil {
		testsuit.T.Errorf("failed to list task runs in namespace %s \n %v", namespace, err)
	}

	var matched_tr string
	match, _ := regexp.Compile(trname + ".*?")
	for _, tr := range trlist.Items {
		if match.MatchString(tr.Name) {
			matched_tr = tr.Name
			break
		}
	}
	labels := make(map[string]string)

	tr, err := c.TaskRunClient.Get(c.Ctx, matched_tr, metav1.GetOptions{})
	if err != nil {
		testsuit.T.Errorf("failed to get task run %s in namespace %s \n %v", matched_tr, namespace, err)
	}

	for key, val := range tr.ObjectMeta.Labels {
		labels[key] = val
	}

	AssertLabelsMatch(labels, tr.ObjectMeta.Labels)
	if tr.Status.PodName != "" {
		pod := GetPodForTaskRun(c, namespace, tr)
		// This label is added to every Pod by the TaskRun controller
		labels[pipeline.TaskRunLabelKey] = tr.Name
		AssertLabelsMatch(labels, pod.ObjectMeta.Labels)
		gauge.WriteMessage("Labels: \n\n %+v", createKeyValuePairs(labels))
	}
}

func getTaskrunLogs(c *clients.Clients, trname, namespace string) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)

	// Set params
	params := cli.TektonParams{}
	params.Clients(c.KubeConfig)
	params.SetNamespace(namespace)

	// Set options for the CLI
	lopts := options.LogOptions{
		TaskrunName: trname,
		Stream: &cli.Stream{
			In:  os.Stdin,
			Out: buf,
			Err: buf,
		},
		Params:    &params,
		Prefixing: true,
	}

	// Get the logs
	err := clitr.Run(&lopts)
	return buf, err
}
