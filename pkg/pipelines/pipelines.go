package pipelines

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/assert"
	"github.com/openshift-pipelines/release-tests/pkg/clients"
	"github.com/openshift-pipelines/release-tests/pkg/k8s"
	"github.com/openshift-pipelines/release-tests/pkg/wait"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	"github.com/tektoncd/pipeline/pkg/reconciler/pipelinerun/resources"
	"gomodules.xyz/jsonpatch/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
)

var prGroupResource = schema.GroupVersionResource{Group: "tekton.dev", Resource: "pipelineruns"}

func validatePipelineRunForSuccessStatus(c *clients.Clients, prname, labelCheck, namespace string) {
	var err error
	// Verify status of PipelineRun (wait for it)
	err = wait.WaitForPipelineRunState(c, prname, wait.PipelineRunSucceed(prname), "PipelineRunCompleted")
	assert.NoError(err, fmt.Sprintf("Error waiting for PipelineRun %s to finish", prname))
	log.Printf("pipelineRun: %s is successfull under namespace : %s", prname, namespace)

	if strings.ToLower(labelCheck) == "yes" || strings.ToLower(labelCheck) == "y" {
		log.Println("Check for events, labels & annotations")
		actualTaskrunList, err := c.TaskRunClient.List(metav1.ListOptions{LabelSelector: fmt.Sprintf("tekton.dev/pipelineRun=%s", prname)})
		assert.NoError(err, fmt.Sprintf("Error listing TaskRuns for PipelineRun %s: %s", prname, err))
		actualTaskRunNames := []string{}
		for _, tr := range actualTaskrunList.Items {
			actualTaskRunNames = append(actualTaskRunNames, tr.GetName())
			log.Printf("Checking that labels were propagated correctly for TaskRun %s", tr.Name)
			checkLabelPropagation(c, namespace, prname, &tr)
			log.Printf("Checking that annotations were propagated correctly for TaskRun %s", tr.Name)
			checkAnnotationPropagation(c, namespace, prname, &tr)
		}

		matchKinds := map[string][]string{"PipelineRun": {prname}, "TaskRun": actualTaskRunNames}

		log.Printf("Making sure %d events were created from taskrun and pipelinerun with kinds %v", len(actualTaskRunNames)+1, matchKinds)

		events, err := collectMatchingEvents(c, namespace, matchKinds, "Succeeded")

		assert.NoError(err, fmt.Sprintf("Failed to collect matching events: %q", err))
		if len(events) != len(actualTaskRunNames)+1 {
			testsuit.T.Errorf(fmt.Sprintf("Expected %d number of successful events from pipelinerun and taskrun but got %d; list of receieved events : %#v", len(actualTaskRunNames)+1, len(events), events))
		}
	}
}

func validatePipelineRunForFailedStatus(c *clients.Clients, prname, namespace string) {
	var err error
	log.Printf("Waiting for PipelineRun in namespace %s to fail", namespace)
	err = wait.WaitForPipelineRunState(c, prname, wait.PipelineRunFailed(prname), "BuildValidationFailed")
	assert.NoError(err, fmt.Sprintf("Failed to finish PipelineRun: %s", prname))
}

func validatePipelineRunTimeoutFailure(c *clients.Clients, prname, namespace string) {
	var err error
	pipelineRun, err := c.PipelineRunClient.Get(prname, metav1.GetOptions{})
	assert.NoError(err, fmt.Sprintf("Error Getting PipelineRun %s under namespace %s ", prname, namespace))

	log.Printf("Waiting for Pipelinerun %s in namespace %s to be started", pipelineRun.Name, namespace)
	if err := wait.WaitForPipelineRunState(c, pipelineRun.Name, wait.Running(pipelineRun.Name), "PipelineRunRunning"); err != nil {
		testsuit.T.Errorf(fmt.Sprintf("Error waiting for PipelineRun %s to be running: %s", pipelineRun.Name, err))
	}

	taskrunList, err := c.TaskRunClient.List(metav1.ListOptions{LabelSelector: fmt.Sprintf("tekton.dev/pipelineRun=%s", pipelineRun.Name)})
	if err != nil {
		testsuit.T.Errorf(fmt.Sprintf("Error listing TaskRuns for PipelineRun %s: %v", pipelineRun.Name, err))
	}

	log.Printf("Waiting for TaskRuns from PipelineRun %s in namespace %s to be running", pipelineRun.Name, namespace)
	errChan := make(chan error, len(taskrunList.Items))
	defer close(errChan)

	for _, taskrunItem := range taskrunList.Items {
		go func(name string) {
			err := wait.WaitForTaskRunState(c, name, wait.Running(name), "TaskRunRunning")
			errChan <- err
		}(taskrunItem.Name)
	}

	for i := 1; i <= len(taskrunList.Items); i++ {
		if <-errChan != nil {
			testsuit.T.Errorf(fmt.Sprintf("Error waiting for TaskRun %s to be running: %v", taskrunList.Items[i-1].Name, err))
		}
	}

	if _, err := c.PipelineRunClient.Get(pipelineRun.Name, metav1.GetOptions{}); err != nil {
		testsuit.T.Errorf(fmt.Sprintf("Failed to get PipelineRun `%s`: %s", pipelineRun.Name, err))
	}

	log.Printf("Waiting for PipelineRun %s in namespace %s to be timed out", pipelineRun.Name, namespace)
	if err := wait.WaitForPipelineRunState(c, pipelineRun.Name, wait.FailedWithReason(resources.ReasonTimedOut, pipelineRun.Name), "PipelineRunTimedOut"); err != nil {
		testsuit.T.Errorf(fmt.Sprintf("Error waiting for PipelineRun %s to finish: %s", pipelineRun.Name, err))
	}

	log.Printf("Waiting for TaskRuns from PipelineRun %s in namespace %s to be cancelled", pipelineRun.Name, namespace)
	var wg sync.WaitGroup
	for _, taskrunItem := range taskrunList.Items {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			err := wait.WaitForTaskRunState(c, name, wait.FailedWithReason("TaskRunTimeout", name), "TaskRunTimeout")
			assert.NoError(err, fmt.Sprintf("Error waiting for TaskRun %s to timeout: %s", name, err))
		}(taskrunItem.Name)
	}
	wg.Wait()

	if _, err := c.PipelineRunClient.Get(pipelineRun.Name, metav1.GetOptions{}); err != nil {
		testsuit.T.Errorf(fmt.Sprintf("Failed to get PipelineRun `%s`: %s", pipelineRun.Name, err))
	}
}

func validatePipelineRunCancel(c *clients.Clients, prname, namespace string) {
	var err error

	log.Printf("Waiting for Pipelinerun %s in namespace %s to be started", prname, namespace)
	if err := wait.WaitForPipelineRunState(c, prname, wait.Running(prname), "PipelineRunRunning"); err != nil {
		testsuit.T.Errorf(fmt.Sprintf("Error waiting for PipelineRun %s to be running: %s", prname, err))
	}

	taskrunList, err := c.TaskRunClient.List(metav1.ListOptions{LabelSelector: fmt.Sprintf("tekton.dev/pipelineRun=%s", prname)})
	if err != nil {
		testsuit.T.Errorf(fmt.Sprintf("Error listing TaskRuns for PipelineRun %s: %s", prname, err))
	}

	var wg sync.WaitGroup
	var trName []string
	log.Printf("Waiting for TaskRuns from PipelineRun %s in namespace %s to be running", prname, namespace)
	for _, taskrunItem := range taskrunList.Items {
		trName = append(trName, taskrunItem.Name)
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			err := wait.WaitForTaskRunState(c, name, wait.Running(name), "TaskRunRunning")
			assert.NoError(err, fmt.Sprintf("Error waiting for TaskRun %s to be running", name))
		}(taskrunItem.Name)
	}
	wg.Wait()

	pr, err := c.PipelineRunClient.Get(prname, metav1.GetOptions{})
	assert.NoError(err, fmt.Sprintf("Error Getting PipelineRun %s under namespace %s ", prname, namespace))

	patches := []jsonpatch.JsonPatchOperation{{
		Operation: "add",
		Path:      "/spec/status",
		Value:     v1beta1.PipelineRunSpecStatusCancelled,
	}}
	patchBytes, err := json.Marshal(patches)
	assert.NoError(err, fmt.Sprintf("failed to marshal patch bytes in order to cancel"))

	if _, err := c.PipelineRunClient.Patch(pr.Name, types.JSONPatchType, patchBytes, ""); err != nil {
		testsuit.T.Errorf(fmt.Sprintf("Failed to patch PipelineRun `%s` with cancellation", prname))
	}

	log.Printf("Waiting for PipelineRun %s in namespace %s to be cancelled", prname, namespace)
	if err := wait.WaitForPipelineRunState(c, prname, wait.FailedWithReason("PipelineRunCancelled", prname), "PipelineRunCancelled"); err != nil {
		testsuit.T.Errorf(fmt.Sprintf("Error waiting for PipelineRun `pear` to finished: %s", err))
	}

	log.Printf("Waiting for TaskRuns in PipelineRun %s in namespace %s to be cancelled", prname, namespace)
	for _, taskrunItem := range taskrunList.Items {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			err := wait.WaitForTaskRunState(c, name, wait.FailedWithReason("TaskRunCancelled", name), "TaskRunCancelled")
			assert.NoError(err, fmt.Sprintf("Error waiting for TaskRun %s to be finished: %v", name, err))
		}(taskrunItem.Name)
	}
	wg.Wait()
}

func ValidatePipelineRun(c *clients.Clients, prname, status, labelCheck, namespace string) {
	var err error
	pr, err := c.PipelineRunClient.Get(prname, metav1.GetOptions{})
	assert.NoError(err, fmt.Sprintf("Error Getting PipelineRun %s under namespace %s ", prname, namespace))

	// Verify status of PipelineRun (wait for it)
	switch {
	case strings.Contains(strings.ToLower(status), "success"):
		log.Printf("validating pipeline run for success state...")
		validatePipelineRunForSuccessStatus(c, pr.GetName(), labelCheck, namespace)
	case strings.Contains(strings.ToLower(status), "fail"):
		log.Printf("validating pipeline run for failure state...")
		validatePipelineRunForFailedStatus(c, pr.GetName(), namespace)
	case strings.Contains(strings.ToLower(status), "timeout"):
		log.Printf("validating pipeline run timeout...")
		validatePipelineRunTimeoutFailure(c, pr.GetName(), namespace)
	case strings.Contains(strings.ToLower(status), "cancel"):
		log.Printf("validating pipeline run timeout...")
		validatePipelineRunCancel(c, pr.GetName(), namespace)
	default:
		testsuit.T.Errorf("Error: %s ", "Not valid input")
	}
}

func WatchForPipelineRun(c *clients.Clients, namespace string) {
	var prnames = []string{}
	watchRun, err := k8s.Watch(prGroupResource, c, namespace, metav1.ListOptions{})
	assert.NoError(err, fmt.Sprintf("failed to pipelineruns on a namespace %s", namespace))
	ch := watchRun.ResultChan()
	go func() {
		for event := range ch {
			run, err := cast2pipelinerun(event.Object)
			assert.NoError(err, fmt.Sprintf("failed to convert to v1beta1 pipelinerun on a namespace %s", namespace))
			switch event.Type {
			case watch.Added:
				log.Printf("pipeline run : %s", run.Name)
				prnames = append(prnames, run.Name)
			}

		}
	}()
	time.Sleep(5 * time.Minute)
	gauge.GetScenarioStore()["prcount"] = len(prnames)
	gauge.WriteMessage("%+v", prnames)
}

func cast2pipelinerun(obj runtime.Object) (*v1beta1.PipelineRun, error) {
	var run *v1beta1.PipelineRun
	unstruct, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return nil, err
	}
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(unstruct, &run); err != nil {
		return nil, err
	}
	return run, nil
}

func AssertForNoNewPipelineRunCreation(c *clients.Clients, namespace string) {
	count := 0
	expectedCount := gauge.GetScenarioStore()["prcount"].(int)
	watchRun, err := k8s.Watch(prGroupResource, c, namespace, metav1.ListOptions{})
	assert.NoError(err, fmt.Sprintf("failed to get tekton resources on a namespace %s", namespace))
	ch := watchRun.ResultChan()
	go func() {
		for event := range ch {
			switch event.Type {
			case watch.Added:
				count++
			}
		}
	}()
	time.Sleep(1 * time.Minute)
	if count != expectedCount {
		testsuit.T.Errorf("Error:  Expected: %+v (tekton resources add newly in namespace %s), \n Actual: %+v ", expectedCount, namespace, count)
	}
}
