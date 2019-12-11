package pipelines

import (
	"fmt"
	"log"

	"github.com/openshift-pipelines/release-tests/pkg/client"
	"github.com/openshift-pipelines/release-tests/pkg/helper"
	"github.com/openshift-pipelines/release-tests/pkg/wait"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	tb "github.com/tektoncd/pipeline/test/builder"
)

var (
	teTaskName        = "output-task"
	teTaskRunName     = "output-task-run"
	tePipelineRunName = "output-pipeline-run"
)

func newTask(taskname string, namespace string) *v1alpha1.Task {

	taskSpecOps := []tb.TaskSpecOp{
		tb.Step("foo", "busybox", tb.StepCommand("ls", "-la")),
	}

	return tb.Task(taskname, namespace, tb.TaskSpec(taskSpecOps...))
}

func newTaskRunWithSA(taskrunname string, namespace string, taskname string, sa string) *v1alpha1.TaskRun {

	return tb.TaskRun(taskrunname, namespace, tb.TaskRunSpec(
		tb.TaskRunTaskRef(taskname), tb.TaskRunServiceAccountName(sa),
	))
}

func newPipeline(pipelineName string, namespace string, taskname string) *v1alpha1.Pipeline {

	pipelineSpec := []tb.PipelineSpecOp{
		tb.PipelineTask("foo", taskname),
	}
	return tb.Pipeline(pipelineName, namespace, tb.PipelineSpec(pipelineSpec...))
}

func newPipelineRunWithSA(pipelineRunName string, namespace string, pipelineName string, sa string) *v1alpha1.PipelineRun {
	return tb.PipelineRun(pipelineRunName, namespace, tb.PipelineRunSpec(
		pipelineName, tb.PipelineRunServiceAccountName(sa),
	))
}

func CreatePipeline(c *client.Clients, namespace string) {
	var err error
	log.Printf("Creating Task  %s", teTaskName)
	_, err = c.TaskClient.Create(newTask(teTaskName, namespace))
	helper.AssertNoError(err, fmt.Sprintf("Failed to create Task Resource `%s`", teTaskName))

	log.Printf("Create Pipeline %s", tePipelineName)
	_, err = c.PipelineClient.Create(newPipeline(tePipelineName, namespace, teTaskName))
	helper.AssertNoError(err, fmt.Sprintf("Failed to create pipeline `%s`", tePipelineName))

	log.Println("Created pipeline successfully....")
}

func CreateTask(c *client.Clients, namespace string) {
	var err error
	log.Printf("Creating Task  %s", teTaskName)
	_, err = c.TaskClient.Create(newTask(teTaskName, namespace))
	helper.AssertNoError(err, fmt.Sprintf("Failed to create Task Resource `%s`", teTaskName))

	log.Println("Created Task successfully....")
}

func CreateTaskRunWithSA(c *client.Clients, namespace string, sa string) {
	var err error
	log.Printf("Starting TaskRun with Service Account %s", sa)
	_, err = c.TaskRunClient.Create(newTaskRunWithSA(teTaskRunName, namespace, teTaskName, sa))
	helper.AssertNoError(err, fmt.Sprintf("Failed to create TaskRun: %s", teTaskRunName))
}

func ValidateTaskRunForFailedStatus(c *client.Clients, namespace string) {
	var err error
	log.Printf("Waiting for TaskRun in namespace %s to fail", namespace)
	err = wait.WaitForTaskRunState(c, teTaskRunName, wait.TaskRunFailed(teTaskRunName), "BuildValidationFailed")
	helper.AssertNoError(err, fmt.Sprintf("Failed to TaskRun: %s", teTaskRunName))

}

func CreatePipelineRunWithSA(c *client.Clients, namespace string, sa string) {
	var err error
	log.Printf("Starting PipelineRun : %s, with Service Account %s", tePipelineRunName, sa)
	_, err = c.PipelineRunClient.Create(newPipelineRunWithSA(tePipelineRunName, namespace, tePipelineName, sa))
	helper.AssertNoError(err, fmt.Sprintf("Failed to create PipelineRun `%s`", tePipelineRunName))
}

func ValidatePipelineRunForFailedStatus(c *client.Clients, namespace string) {
	var err error
	log.Printf("Waiting for PipelineRun in namespace %s to fail", namespace)
	err = wait.WaitForPipelineRunState(c, tePipelineRunName, wait.PipelineRunFailed(tePipelineRunName), "BuildValidationFailed")
	helper.AssertNoError(err, fmt.Sprintf("Failed to finish PipelineRun: %s", tePipelineRunName))

}
