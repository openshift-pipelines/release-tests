package pipelines

import (
	"log"
	"testing"
	"time"

	"github.com/openshift-pipelines/release-tests/pkg/client"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	tb "github.com/tektoncd/pipeline/test/builder"
	"gotest.tools/v3/icmd"
)

const (
	TaskName1                 = "create-file"
	TaskName2                 = "check-stuff-file-exists"
	tePipelineName            = "output-pipeline"
	tePipelineGitResourceName = "skaffold-git"
)

func getGitResourceForOutPutPipeline(rname string, namespace string) *v1alpha1.PipelineResource {
	return tb.PipelineResource(rname, namespace, tb.PipelineResourceSpec(
		v1alpha1.PipelineResourceTypeGit,
		tb.PipelineResourceSpecParam("url", "https://github.com/GoogleContainerTools/skaffold"),
		tb.PipelineResourceSpecParam("revision", "master"),
	))
}

func getCreateFileTask(taskname string, namespace string) *v1alpha1.Task {

	taskSpecOps := []tb.TaskSpecOp{
		tb.TaskInputs(tb.InputsResource("workspace", v1alpha1.PipelineResourceTypeGit, tb.ResourceTargetPath("damnworkspace"))),
		tb.TaskOutputs(tb.OutputsResource("workspace", v1alpha1.PipelineResourceTypeGit)),
		tb.Step("read-docs-old", "ubuntu", tb.StepCommand("/bin/bash"), tb.StepArgs("-c", "ls -la /workspace/damnworkspace/docs/README.md")),
		tb.Step("write-new-stuff", "ubuntu", tb.StepCommand("bash"), tb.StepArgs("-c", "ln -s /workspace/damnworkspace /workspace/output/workspace && echo some stuff > /workspace/output/workspace/stuff")),
	}

	return tb.Task(taskname, namespace, tb.TaskSpec(taskSpecOps...))
}

func getReadFileTask(taskname string, namespace string) *v1alpha1.Task {

	taskSpecOps := []tb.TaskSpecOp{
		tb.TaskInputs(tb.InputsResource("workspace", v1alpha1.PipelineResourceTypeGit, tb.ResourceTargetPath("newworkspace"))),
		tb.Step("read", "ubuntu", tb.StepCommand("/bin/bash"), tb.StepArgs("-c", "cat /workspace/newworkspace/stuff")),
	}

	return tb.Task(taskname, namespace, tb.TaskSpec(taskSpecOps...))
}

func getPipeline(pipelineName string, namespace string, createFiletaskName string, readFileTaskName string) *v1alpha1.Pipeline {

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

func CreateSamplePiplines(c *client.Clients, namespace string) {

	log.Printf("Creating Git PipelineResource %s", tePipelineGitResourceName)
	if _, err := c.PipelineResourceClient.Create(getGitResourceForOutPutPipeline(tePipelineGitResourceName, namespace)); err != nil {
		log.Fatalf("Failed to create Pipeline Resource `%s`: %s", tePipelineGitResourceName, err)
	}

	log.Printf("Creating Task  %s", TaskName1)
	if _, err := c.TaskClient.Create(getCreateFileTask(TaskName1, namespace)); err != nil {
		log.Fatalf("Failed to create Task Resource `%s`: %s", TaskName1, err)
	}

	log.Printf("Creating Task  %s", TaskName2)
	if _, err := c.TaskClient.Create(getReadFileTask(TaskName2, namespace)); err != nil {
		log.Fatalf("Failed to create Task Resource `%s`: %s", TaskName2, err)
	}

	log.Printf("Create Pipeline %s", tePipelineName)
	if _, err := c.PipelineClient.Create(getPipeline(tePipelineName, namespace, TaskName1, TaskName2)); err != nil {
		log.Fatalf("Failed to create pipeline `%s`: %s", tePipelineName, err)
	}

	time.Sleep(1 * time.Second)
	log.Println("Created sample pipleine successfully....")

}

func StartSamplePipelineUsingTkn(t *testing.T, namespace string) {
	log.Printf("starting pipeline %s using tkn \n", tePipelineName)
	res := icmd.RunCmd(icmd.Cmd{Command: append([]string{"tkn"},
		"pipeline",
		"start",
		tePipelineName,
		"-r=source-repo="+tePipelineGitResourceName,
		"--showlog",
		"true",
		"-n",
		namespace),
		Timeout: 10 * time.Minute})
	res.Assert(t, icmd.Expected{
		ExitCode: 0,
		Err:      icmd.None,
	})
	log.Printf("Output: %+v", res.Stdout())
}
