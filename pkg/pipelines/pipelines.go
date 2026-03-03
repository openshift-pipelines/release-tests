package pipelines

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/clients"
	"github.com/openshift-pipelines/release-tests/pkg/cmd"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	"github.com/openshift-pipelines/release-tests/pkg/k8s"
	"github.com/openshift-pipelines/release-tests/pkg/store"
	"github.com/openshift-pipelines/release-tests/pkg/wait"
	"github.com/tektoncd/cli/pkg/cli"
	clipr "github.com/tektoncd/cli/pkg/cmd/pipelinerun"
	"github.com/tektoncd/cli/pkg/options"
	prsort "github.com/tektoncd/cli/pkg/pipelinerun/sort"
	v1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	"knative.dev/pkg/apis"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	w "k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
)

var prGroupResource = schema.GroupVersionResource{Group: "tekton.dev", Resource: "pipelineruns"}

func validatePipelineRunForSuccessStatus(c *clients.Clients, prname, namespace string) {
	// Verify status of PipelineRun (wait for it)
	err := wait.WaitForPipelineRunState(c, prname, wait.PipelineRunSucceed(prname), "PipelineRunCompleted")
	if err != nil {
		buf, logsErr := getPipelinerunLogs(c, prname, namespace)
		events, eventError := k8s.GetWarningEvents(c, namespace)
		if logsErr != nil {
			if eventError != nil {
				testsuit.T.Errorf("error waiting for pipeline run %s to finish \n%v \npipelinerun logs error: \n%v \npipelinerun events error: \n%v", prname, err, logsErr, eventError)
			} else {
				testsuit.T.Errorf("error waiting for pipeline run %s to finish \n%v \npipelinerun logs error: \n%v \npipelinerun events: \n%v", prname, err, logsErr, events)
			}
		} else {
			if eventError != nil {
				testsuit.T.Errorf("error waiting for pipeline run %s to finish \n%v \npipelinerun logs: \n%v \npipelinerun events error: \n%v", prname, err, buf.String(), eventError)
			} else {
				testsuit.T.Errorf("error waiting for pipeline run %s to finish \n%v \npipelinerun logs: \n%v \npipelinerun events: \n%v", prname, err, buf.String(), events)
			}
		}
	}

	log.Printf("pipelineRun: %s is successful under namespace : %s", prname, namespace)
}

func validatePipelineRunForFailedStatus(c *clients.Clients, prname, namespace string) {
	var err error
	log.Printf("Waiting for PipelineRun in namespace %s to fail", namespace)
	err = wait.WaitForPipelineRunState(c, prname, wait.PipelineRunFailed(prname), "BuildValidationFailed")
	if err != nil {
		buf, logsErr := getPipelinerunLogs(c, prname, namespace)
		events, eventError := k8s.GetWarningEvents(c, namespace)
		if logsErr != nil {
			if eventError != nil {
				testsuit.T.Errorf("error waiting for pipeline run %s to finish \n%v \npipelinerun logs error: \n%v \npipelinerun events error: \n%v", prname, err, logsErr, eventError)
			} else {
				testsuit.T.Errorf("error waiting for pipeline run %s to finish \n%v \npipelinerun logs error: \n%v \npipelinerun events: \n%v", prname, err, logsErr, events)
			}
		} else {
			if eventError != nil {
				testsuit.T.Errorf("error waiting for pipeline run %s to finish \n%v \npipelinerun logs: \n%v \npipelinerun events error: \n%v", prname, err, buf.String(), eventError)
			} else {
				testsuit.T.Errorf("error waiting for pipeline run %s to finish \n%v \npipelinerun logs: \n%v \npipelinerun events: \n%v", prname, err, buf.String(), events)
			}
		}
	}
}

func validatePipelineRunTimeoutFailure(c *clients.Clients, prname, namespace string) {
	var err error
	pipelineRun, err := c.PipelineRunClient.Get(c.Ctx, prname, metav1.GetOptions{})
	if err != nil {
		testsuit.T.Errorf("failed to get pipeline run %s in namespaces %s \n %v", prname, namespace, err)
	}

	log.Printf("Waiting for Pipelinerun %s in namespace %s to be started", pipelineRun.Name, namespace)
	if err := wait.WaitForPipelineRunState(c, pipelineRun.Name, wait.Running(pipelineRun.Name), "PipelineRunRunning"); err != nil {
		testsuit.T.Errorf("Error waiting for PipelineRun %s to be running: %s", pipelineRun.Name, err)
	}

	taskrunList, err := c.TaskRunClient.List(c.Ctx, metav1.ListOptions{LabelSelector: fmt.Sprintf("tekton.dev/pipelineRun=%s", pipelineRun.Name)})
	if err != nil {
		testsuit.T.Errorf("Error listing TaskRuns for PipelineRun %s: %v", pipelineRun.Name, err)
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
			testsuit.T.Errorf("Error waiting for TaskRun %s to be running: %v", taskrunList.Items[i-1].Name, err)
		}
	}

	if _, err := c.PipelineRunClient.Get(c.Ctx, pipelineRun.Name, metav1.GetOptions{}); err != nil {
		testsuit.T.Errorf("Failed to get PipelineRun `%s`: %s", pipelineRun.Name, err)
	}

	log.Printf("Waiting for PipelineRun %s in namespace %s to be timed out", pipelineRun.Name, namespace)
	if err := wait.WaitForPipelineRunState(c, pipelineRun.Name, wait.FailedWithReason(v1.PipelineRunReasonTimedOut.String(), pipelineRun.Name), "PipelineRunTimedOut"); err != nil {
		testsuit.T.Errorf("Error waiting for PipelineRun %s to finish: %s", pipelineRun.Name, err)
	}

	log.Printf("Waiting for TaskRuns from PipelineRun %s in namespace %s to be cancelled", pipelineRun.Name, namespace)
	var wg sync.WaitGroup
	for _, taskrunItem := range taskrunList.Items {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			err := wait.WaitForTaskRunState(c, name, wait.FailedWithReason(v1.TaskRunReasonCancelled.String(), name), v1.TaskRunReasonCancelled.String())
			if err != nil {
				testsuit.T.Errorf("error waiting for task run %s to be cancelled on pipeline timeout \n %v", name, err)
			}
		}(taskrunItem.Name)
	}
	wg.Wait()

	if _, err := c.PipelineRunClient.Get(c.Ctx, pipelineRun.Name, metav1.GetOptions{}); err != nil {
		testsuit.T.Errorf("Failed to get PipelineRun `%s`: %s", pipelineRun.Name, err)
	}
}

func validatePipelineRunCancel(c *clients.Clients, prname, namespace string) {
	var err error

	log.Printf("Waiting for Pipelinerun %s in namespace %s to be started", prname, namespace)
	if err := wait.WaitForPipelineRunState(c, prname, wait.Running(prname), "PipelineRunRunning"); err != nil {
		testsuit.T.Errorf("Error waiting for PipelineRun %s to be running: %s", prname, err)
	}

	taskrunList, err := c.TaskRunClient.List(c.Ctx, metav1.ListOptions{LabelSelector: fmt.Sprintf("tekton.dev/pipelineRun=%s", prname)})
	if err != nil {
		testsuit.T.Errorf("Error listing TaskRuns for PipelineRun %s: %s", prname, err)
	}

	var wg sync.WaitGroup
	log.Printf("Canceling pipeline run: %s\n", cmd.MustSucceed("opc", "pipelinerun", "cancel", prname, "-n", namespace).Stdout())

	if err := wait.WaitForPipelineRunState(c, prname, wait.FailedWithReason("Cancelled", prname), "Cancelled"); err != nil {
		testsuit.T.Errorf("Error waiting for PipelineRun `%s` to finished: %s", prname, err)
	}

	log.Printf("Waiting for TaskRuns in PipelineRun %s in namespace %s to be cancelled", prname, namespace)
	for _, taskrunItem := range taskrunList.Items {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			err := wait.WaitForTaskRunState(c, name, wait.FailedWithReason(v1.TaskRunReasonCancelled.String(), name), "TaskRunCancelled")
			if err != nil {
				testsuit.T.Errorf("task run %s failed to finish \n %v", name, err)
			}
		}(taskrunItem.Name)
	}
	wg.Wait()
}

func ValidatePipelineRun(c *clients.Clients, prname, status, namespace string) {
	var err error
	pr, err := c.PipelineRunClient.Get(c.Ctx, prname, metav1.GetOptions{})
	if err != nil {
		testsuit.T.Errorf("failed to get pipeline run %s in namespace %s \n %v", prname, namespace, err)
	}

	// Verify status of PipelineRun (wait for it)
	switch {
	case strings.Contains(strings.ToLower(status), "success"):
		log.Printf("validating pipeline run %s for success state...", prname)
		validatePipelineRunForSuccessStatus(c, pr.GetName(), namespace)
	case strings.Contains(strings.ToLower(status), "fail"):
		log.Printf("validating pipeline run %s for failure state...", prname)
		validatePipelineRunForFailedStatus(c, pr.GetName(), namespace)
	case strings.Contains(strings.ToLower(status), "timeout"):
		log.Printf("validating pipeline run %s to time out...", prname)
		validatePipelineRunTimeoutFailure(c, pr.GetName(), namespace)
	case strings.Contains(strings.ToLower(status), "cancel"):
		log.Printf("validating pipeline run %s to be cancelled...", prname)
		validatePipelineRunCancel(c, pr.GetName(), namespace)
	default:
		testsuit.T.Errorf("Error: %s ", "Not valid input")
	}
}

func WatchForPipelineRun(c *clients.Clients, namespace string) {
	var prnames = []string{}
	watchRun, err := k8s.Watch(c.Ctx, prGroupResource, c, namespace, metav1.ListOptions{})
	if err != nil {
		testsuit.T.Errorf("failed to watch pipeline runs in namespace %s \n %v", namespace, err)
	}

	ch := watchRun.ResultChan()
	go func() {
		for event := range ch {
			run, err := cast2pipelinerun(event.Object)
			if err != nil {
				testsuit.T.Errorf("failed to convert pipeline run to v1beta1 in namespace %s \n %v", namespace, err)
			}
			if event.Type == watch.Added {
				log.Printf("pipeline run : %s", run.Name)
				prnames = append(prnames, run.Name)
			}

		}
	}()
	time.Sleep(5 * time.Minute)
	gauge.GetScenarioStore()["prcount"] = len(prnames)
	gauge.WriteMessage("%+v", prnames)
}

func cast2pipelinerun(obj runtime.Object) (*v1.PipelineRun, error) {
	var run *v1.PipelineRun
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
	watchRun, err := k8s.Watch(c.Ctx, prGroupResource, c, namespace, metav1.ListOptions{})
	if err != nil {
		testsuit.T.Errorf("failed to watch pipeline runs in namespace %s \n %v", namespace, err)
	}
	ch := watchRun.ResultChan()
	go func() {
		for event := range ch {
			if event.Type == watch.Added {
				count++
			}
		}
	}()
	time.Sleep(1 * time.Minute)
	if count < expectedCount {
		testsuit.T.Errorf("Error:  Expected: %+v (tekton resources add newly in namespace %s), \n Actual: %+v ", expectedCount, namespace, count)
	}
}

func AssertNumberOfPipelineruns(c *clients.Clients, namespace, numberOfPr, timeoutSeconds string) {
	log.Printf("Verifying if %s pipelineruns are present", numberOfPr)
	timeoutSecondsInt, _ := strconv.Atoi(timeoutSeconds)
	err := w.PollUntilContextTimeout(c.Ctx, config.APIRetry, time.Second*time.Duration(timeoutSecondsInt), false, func(context.Context) (bool, error) {
		prlist, err := c.PipelineRunClient.List(c.Ctx, metav1.ListOptions{})
		numberOfPrInt, _ := strconv.Atoi(numberOfPr)
		if len(prlist.Items) == numberOfPrInt {
			return true, nil
		}
		return false, err
	})
	if err != nil {
		prlist, _ := c.PipelineRunClient.List(c.Ctx, metav1.ListOptions{})
		testsuit.T.Fail(fmt.Errorf("error: Expected %v pipelineruns but found %v pipelineruns: %s", numberOfPr, len(prlist.Items), err))
	}
}

// countPipelinerunsByStatus returns the number of PipelineRuns in the list that match the given status ("Succeeded" or "Failed").
func countPipelinerunsByStatus(prlist *v1.PipelineRunList, status string) int {
	var match corev1.ConditionStatus
	switch strings.ToLower(status) {
	case "succeeded":
		match = corev1.ConditionTrue
	case "failed":
		match = corev1.ConditionFalse
	default:
		return 0
	}
	count := 0
	for i := range prlist.Items {
		pr := &prlist.Items[i]
		c := pr.Status.GetCondition(apis.ConditionSucceeded)
		if c != nil && c.Status == match {
			count++
		}
	}
	return count
}

func AssertNumberOfPipelinerunsWithStatus(c *clients.Clients, namespace, numberOfPr, status, timeoutSeconds string) {
	log.Printf("Verifying if %s pipelineruns with status %s are present", numberOfPr, status)
	timeoutSecondsInt, _ := strconv.Atoi(timeoutSeconds)
	numberOfPrInt, _ := strconv.Atoi(numberOfPr)
	err := w.PollUntilContextTimeout(c.Ctx, config.APIRetry, time.Second*time.Duration(timeoutSecondsInt), false, func(context.Context) (bool, error) {
		prlist, err := c.PipelineRunClient.List(c.Ctx, metav1.ListOptions{})
		if err != nil {
			return false, err
		}
		count := countPipelinerunsByStatus(prlist, status)
		if count == numberOfPrInt {
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		prlist, _ := c.PipelineRunClient.List(c.Ctx, metav1.ListOptions{})
		count := countPipelinerunsByStatus(prlist, status)
		testsuit.T.Fail(fmt.Errorf("error: Expected %v pipelineruns with status %s but found %v (total pipelineruns: %v): %s", numberOfPr, status, count, len(prlist.Items), err))
	}
}

func AssertNumberOfTaskruns(c *clients.Clients, namespace, numberOfTr, timeoutSeconds string) {
	log.Printf("Verifying if %s taskruns are present", numberOfTr)
	timeoutSecondsInt, _ := strconv.Atoi(timeoutSeconds)
	err := w.PollUntilContextTimeout(c.Ctx, config.APIRetry, time.Second*time.Duration(timeoutSecondsInt), false, func(context.Context) (bool, error) {
		trlist, err := c.TaskRunClient.List(c.Ctx, metav1.ListOptions{})
		numberOfPrInt, _ := strconv.Atoi(numberOfTr)
		if len(trlist.Items) == numberOfPrInt {
			return true, nil
		}
		return false, err
	})
	if err != nil {
		trlist, _ := c.TaskRunClient.List(c.Ctx, metav1.ListOptions{})
		testsuit.T.Fail(fmt.Errorf("error: Expected %v taskruns but found %v taskruns: %s", numberOfTr, len(trlist.Items), err))
	}
}
func AssertPipelinesPresent(c *clients.Clients, namespace string) {
	pclient := c.Tekton.TektonV1beta1().Pipelines(namespace)
	expectedNumberOfPipelines := len(config.PrefixesOfDefaultPipelines)
	if config.Flags.ClusterArch == "arm64" {
		expectedNumberOfPipelines *= 2
	} else {
		expectedNumberOfPipelines *= 3
	}

	err := w.PollUntilContextTimeout(c.Ctx, config.APIRetry, config.ResourceTimeout, false, func(context.Context) (bool, error) {
		log.Printf("Verifying that %v pipelines are present in namespace %v", expectedNumberOfPipelines, namespace)
		p, _ := pclient.List(c.Ctx, metav1.ListOptions{})
		if len(p.Items) == expectedNumberOfPipelines {
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		p, _ := pclient.List(c.Ctx, metav1.ListOptions{})
		testsuit.T.Fail(fmt.Errorf("expected: %v pipelines present in namespace %v, Actual: %v pipelines present in namespace %v , Error: %v", expectedNumberOfPipelines, namespace, len(p.Items), namespace, err))
	}
	log.Printf("Pipelines are present in namespace %v", namespace)
}

func AssertPipelinesNotPresent(c *clients.Clients, namespace string) {
	pclient := c.Tekton.TektonV1beta1().Pipelines(namespace)
	err := w.PollUntilContextTimeout(c.Ctx, config.APIRetry, config.ResourceTimeout, false, func(context.Context) (bool, error) {
		log.Printf("Verifying if 0 pipelines are not present in namespace %v", namespace)
		p, _ := pclient.List(c.Ctx, metav1.ListOptions{})
		if len(p.Items) == 0 {
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		p, _ := pclient.List(c.Ctx, metav1.ListOptions{})
		testsuit.T.Fail(fmt.Errorf("expected: %v number of pipelines present in namespace %v, Actual: %v number of pipelines present in namespace %v , Error: %v", 0, namespace, len(p.Items), namespace, err))
	}
	log.Printf("Pipelines are present in namespace %v", namespace)
}

func getPipelinerunLogs(c *clients.Clients, prname, namespace string) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)

	// Set params
	params := cli.TektonParams{}
	_, err := params.Clients(c.KubeConfig)
	if err != nil {
		log.Printf("Client Initialization Failed\n %v", err)
		return nil, err
	}
	params.SetNamespace(namespace)

	// Set options for the CLI
	lopts := options.LogOptions{
		PipelineRunName: prname,
		Stream: &cli.Stream{
			In:  os.Stdin,
			Out: buf,
			Err: buf,
		},
		Params:    &params,
		Prefixing: true,
	}

	// Get the logs
	err = clipr.Run(&lopts)
	return buf, err
}

func GetLatestPipelinerun(c *clients.Clients, namespace string) (string, error) {
	prs, err := c.PipelineRunClient.List(c.Ctx, metav1.ListOptions{})
	if err != nil {
		return "", err
	}
	if len(prs.Items) == 0 {
		return "", fmt.Errorf("no pipelineruns found in the namespace %s", namespace)
	}
	prsort.SortByStartTime(prs.Items)
	return prs.Items[0].Name, nil

}

func CheckLogVersion(c *clients.Clients, binary, namespace string) {
	prname, err := GetLatestPipelinerun(store.Clients(), store.Namespace())
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("failed to get PipelineRun: %v", err))
		return
	}
	// Get PipelineRun logs
	logsBuffer, err := getPipelinerunLogs(c, prname, namespace)
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("failed to get PipelineRun logs: %v", err))
		return
	}

	switch binary {
	case "tkn-pac":
		expectedVersion := os.Getenv("PAC_VERSION")
		if !strings.Contains(logsBuffer.String(), expectedVersion) {
			testsuit.T.Fail(fmt.Errorf("tkn-pac Version %s not found in logs:\n%s ", expectedVersion, logsBuffer))
		}
	case "tkn":
		expectedVersion := os.Getenv("TKN_CLIENT_VERSION")
		if !strings.Contains(logsBuffer.String(), "Client version:") {
			testsuit.T.Fail(fmt.Errorf("tkn client version not found! \nlogs:%s", logsBuffer))
			return
		}
		if !strings.Contains(logsBuffer.String(), expectedVersion) {
			testsuit.T.Fail(fmt.Errorf("tkn Version %s not found in logs:\n%s ", expectedVersion, logsBuffer))
		}
	default:
		testsuit.T.Fail(fmt.Errorf("unknown binary or client"))
	}
}
