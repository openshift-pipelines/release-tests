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
	createFileTaskName        = "create-file"
	readFileTaskName          = "check-stuff-file-exists"
	tePipelineName            = "output-pipeline"
	tePipelineGitResourceName = "skaffold-git"
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
		tb.Step("read-docs-old", "ubuntu", tb.StepCommand("/bin/bash"), tb.StepArgs("-c", "ls -la /workspace/damnworkspace/docs/README.md")),
		tb.Step("write-new-stuff", "ubuntu", tb.StepCommand("bash"), tb.StepArgs("-c", "ln -s /workspace/damnworkspace /workspace/output/workspace && echo some stuff > /workspace/output/workspace/stuff")),
	}

	return tb.Task(taskname, namespace, tb.TaskSpec(taskSpecOps...))
}

func newReadFileTask(taskname string, namespace string) *v1alpha1.Task {

	taskSpecOps := []tb.TaskSpecOp{
		tb.TaskInputs(tb.InputsResource("workspace", v1alpha1.PipelineResourceTypeGit, tb.ResourceTargetPath("newworkspace"))),
		tb.Step("read", "ubuntu", tb.StepCommand("/bin/bash"), tb.StepArgs("-c", "cat /workspace/newworkspace/stuff")),
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

func CreateSamplePipeline(c *client.Clients, namespace string) {
	var err error

	log.Printf("Creating Git PipelineResource %s", tePipelineGitResourceName)

	_, err = c.PipelineResourceClient.Create(newGitResource(tePipelineGitResourceName, namespace))
	helper.AssertNoError(err, fmt.Sprintf("Failed to create Pipeline Resource `%s`", tePipelineGitResourceName))

	log.Printf("Creating Task  %s", createFileTaskName)

	_, err = c.TaskClient.Create(newCreateFileTask(createFileTaskName, namespace))
	helper.AssertNoError(err, fmt.Sprintf("Failed to create Task `%s`", createFileTaskName))

	log.Printf("Creating Task  %s", readFileTaskName)

	_, err = c.TaskClient.Create(newReadFileTask(readFileTaskName, namespace))
	helper.AssertNoError(err, fmt.Sprintf("Failed to create Task `%s`", readFileTaskName))

	log.Printf("Create Pipeline %s", tePipelineName)

	_, err = c.PipelineClient.Create(newOutputPipeline(tePipelineName, namespace, createFileTaskName, readFileTaskName))
	helper.AssertNoError(err, fmt.Sprintf("Failed to create pipeline `%s`", tePipelineName))

	log.Println("Created sample pipeline successfully....")
}

func RunSamplePipeline(c *client.Clients, namespace string) {
	var err error
	pr := tb.PipelineRun(tePipelineRunName, namespace, tb.PipelineRunSpec(
		tePipelineName,
		tb.PipelineRunResourceBinding("source-repo", tb.PipelineResourceBindingRef(tePipelineGitResourceName)),
	))
	_, err = c.PipelineRunClient.Create(pr)
	helper.AssertNoError(err, fmt.Sprintf("Failed to create PipelineRun `%s`", tePipelineRunName))

}

func ValidatePipelineRunStatus(c *client.Clients, namespace string){
	var err error
	// Verify status of PipelineRun (wait for it)
	err = wait.WaitForPipelineRunState(c, tePipelineRunName, wait.PipelineRunSucceed(tePipelineRunName), "PipelineRunCompleted")
	helper.AssertNoError(err, fmt.Sprintf("Error waiting for PipelineRun %s to finish", tePipelineRunName))
	log.Printf("pipelineRun: %s is successfull under namespace : %s", tePipelineRunName, namespace)
}