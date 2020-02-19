package pipelines

import (
	"fmt"
	"log"

	"github.com/openshift-pipelines/release-tests/pkg/assert"
	"github.com/openshift-pipelines/release-tests/pkg/clients"
	"github.com/openshift-pipelines/release-tests/pkg/wait"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	tb "github.com/tektoncd/pipeline/test/builder"
)

var (
	createFileTaskName        = "create-file"
	readFileTaskName          = "check-stuff-file-exists"
	tePipelineName            = "output-pipeline"
	tePipelineGitResourceName = "skaffold-git"
	teTaskName                = "output-task"
	teTaskRunName             = "output-task-run"
	tePipelineRunName         = "output-pipeline-run"
)

func newGitResource(rname string, namespace string) *v1alpha1.PipelineResource {
	return tb.PipelineResource(rname, namespace, tb.PipelineResourceSpec(
		v1alpha1.PipelineResourceTypeGit,
		tb.PipelineResourceSpecParam("url", "https://github.com/GoogleContainerTools/skaffold"),
		tb.PipelineResourceSpecParam("revision", "master"),
	))
}

func newCreateFileTask(taskname string, namespace string) *v1alpha1.Task {

	taskSpecOps := []tb.TaskSpecOp{
		tb.TaskInputs(tb.InputsResource("workspace", v1alpha1.PipelineResourceTypeGit, tb.ResourceTargetPath("damnworkspace"))),
		tb.TaskOutputs(tb.OutputsResource("workspace", v1alpha1.PipelineResourceTypeGit)),
		tb.Step("ubuntu", tb.StepName("read-docs-old"), tb.StepCommand("/bin/bash"), tb.StepArgs("-c", "ls -la /workspace/damnworkspace/docs/README.md")),
		tb.Step("ubuntu", tb.StepName("write-new-stuff"), tb.StepCommand("bash"), tb.StepArgs("-c", "ln -s /workspace/damnworkspace /workspace/output/workspace && echo some stuff > /workspace/output/workspace/stuff")),
	}

	return tb.Task(taskname, namespace, tb.TaskSpec(taskSpecOps...))
}

func newReadFileTask(taskname string, namespace string) *v1alpha1.Task {

	taskSpecOps := []tb.TaskSpecOp{
		tb.TaskInputs(tb.InputsResource("workspace", v1alpha1.PipelineResourceTypeGit, tb.ResourceTargetPath("newworkspace"))),
		tb.Step("ubuntu", tb.StepName("read"), tb.StepCommand("/bin/bash"), tb.StepArgs("-c", "cat /workspace/newworkspace/stuff")),
	}

	return tb.Task(taskname, namespace, tb.TaskSpec(taskSpecOps...))
}

func newOutputPipeline(pipelineName string, namespace string, createFiletaskName string, readFileTaskName string) *v1alpha1.Pipeline {

	pipelineSpec := []tb.PipelineSpecOp{
		tb.PipelineDeclaredResource("source-repo", "git"),
		tb.PipelineTask("first-create-file", createFiletaskName,
			tb.PipelineTaskInputResource("workspace", "source-repo"),
			tb.PipelineTaskOutputResource("workspace", "source-repo"),
		),
		tb.PipelineTask("then-check", readFileTaskName,
			tb.PipelineTaskInputResource("workspace", "source-repo", tb.From("first-create-file")),
		),
	}

	return tb.Pipeline(pipelineName, namespace, tb.PipelineSpec(pipelineSpec...))
}

func CreateSamplePipeline(c *clients.Clients, namespace string) {
	var err error

	log.Printf("Creating Git PipelineResource %s", tePipelineGitResourceName)

	_, err = c.PipelineResourceClient.Create(newGitResource(tePipelineGitResourceName, namespace))
	assert.NoError(err, fmt.Sprintf("Failed to create Pipeline Resource `%s`", tePipelineGitResourceName))

	log.Printf("Creating Task  %s", createFileTaskName)

	_, err = c.TaskClient.Create(newCreateFileTask(createFileTaskName, namespace))
	assert.NoError(err, fmt.Sprintf("Failed to create Task `%s`", createFileTaskName))

	log.Printf("Creating Task  %s", readFileTaskName)

	_, err = c.TaskClient.Create(newReadFileTask(readFileTaskName, namespace))
	assert.NoError(err, fmt.Sprintf("Failed to create Task `%s`", readFileTaskName))

	log.Printf("Create Pipeline %s", tePipelineName)

	_, err = c.PipelineClient.Create(newOutputPipeline(tePipelineName, namespace, createFileTaskName, readFileTaskName))
	assert.NoError(err, fmt.Sprintf("Failed to create pipeline `%s`", tePipelineName))

	log.Println("Created sample pipeline successfully....")
}

func RunSamplePipeline(c *clients.Clients, namespace string) {
	var err error
	pr := tb.PipelineRun(tePipelineRunName, namespace, tb.PipelineRunSpec(
		tePipelineName,
		tb.PipelineRunResourceBinding("source-repo", tb.PipelineResourceBindingRef(tePipelineGitResourceName)),
	))
	_, err = c.PipelineRunClient.Create(pr)
	assert.NoError(err, fmt.Sprintf("Failed to create PipelineRun `%s`", tePipelineRunName))

}

func ValidatePipelineRunStatus(c *clients.Clients, namespace string) {
	var err error
	// Verify status of PipelineRun (wait for it)
	err = wait.WaitForPipelineRunState(c, tePipelineRunName, wait.PipelineRunSucceed(tePipelineRunName), "PipelineRunCompleted")
	assert.NoError(err, fmt.Sprintf("Error waiting for PipelineRun %s to finish", tePipelineRunName))
	log.Printf("pipelineRun: %s is successfull under namespace : %s", tePipelineRunName, namespace)
}

func newTask(taskname string, namespace string) *v1alpha1.Task {

	taskSpecOps := []tb.TaskSpecOp{
		tb.Step("busybox", tb.StepName("foo"), tb.StepCommand("ls", "-la")),
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

func CreatePipeline(c *clients.Clients, namespace string) {
	var err error
	log.Printf("Creating Task  %s", teTaskName)
	_, err = c.TaskClient.Create(newTask(teTaskName, namespace))
	assert.NoError(err, fmt.Sprintf("Failed to create Task Resource `%s`", teTaskName))

	log.Printf("Create Pipeline %s", tePipelineName)
	_, err = c.PipelineClient.Create(newPipeline(tePipelineName, namespace, teTaskName))
	assert.NoError(err, fmt.Sprintf("Failed to create pipeline `%s`", tePipelineName))

	log.Println("Created pipeline successfully....")
}

func CreateTask(c *clients.Clients, namespace string) {
	var err error
	log.Printf("Creating Task  %s", teTaskName)
	_, err = c.TaskClient.Create(newTask(teTaskName, namespace))
	assert.NoError(err, fmt.Sprintf("Failed to create Task Resource `%s`", teTaskName))

	log.Println("Created Task successfully....")
}

func CreateTaskRunWithSA(c *clients.Clients, namespace string, sa string) {
	var err error
	log.Printf("Starting TaskRun with Service Account %s", sa)
	_, err = c.TaskRunClient.Create(newTaskRunWithSA(teTaskRunName, namespace, teTaskName, sa))
	assert.NoError(err, fmt.Sprintf("Failed to create TaskRun: %s", teTaskRunName))
}

func ValidateTaskRunForFailedStatus(c *clients.Clients, namespace string) {
	var err error
	log.Printf("Waiting for TaskRun in namespace %s to fail", namespace)
	err = wait.WaitForTaskRunState(c, teTaskRunName, wait.TaskRunFailed(teTaskRunName), "BuildValidationFailed")
	assert.NoError(err, fmt.Sprintf("Failed to TaskRun: %s", teTaskRunName))

}

func CreatePipelineRunWithSA(c *clients.Clients, namespace string, sa string) {
	var err error
	log.Printf("Starting PipelineRun : %s, with Service Account %s", tePipelineRunName, sa)
	_, err = c.PipelineRunClient.Create(newPipelineRunWithSA(tePipelineRunName, namespace, tePipelineName, sa))
	assert.NoError(err, fmt.Sprintf("Failed to create PipelineRun `%s`", tePipelineRunName))
}

func ValidatePipelineRunForFailedStatus(c *clients.Clients, namespace string) {
	var err error
	log.Printf("Waiting for PipelineRun in namespace %s to fail", namespace)
	err = wait.WaitForPipelineRunState(c, tePipelineRunName, wait.PipelineRunFailed(tePipelineRunName), "BuildValidationFailed")
	assert.NoError(err, fmt.Sprintf("Failed to finish PipelineRun: %s", tePipelineRunName))

}
