# Verify Pipelines

Pre condition:
  * Operator should be installed

## Create and run Pipeline using cli (without resources)
Tags: pipelines, e2e

Steps:
  * Create task from tkn "../../testdata/cli/shell-script-task.yaml"
  * Verify taks creation status "successfull"
  * Create task from tkn "../../testdata/cli/python-script-task.yaml"
  * Verify task creation status "successfull"
  * Create pipeline from tkn "../../testdata/cli/pipeline-script.yaml"
  * Start the pipeline from tkn "pipeline-test"
  * List the pipeline
  * Verify the status of pipeline "running"
  * List the taskruns from tkn
  * Verify the status of taskrun "running"
  * Wait for the pipelinerun to complete
  * Verify the pipelinerun status "successfull"


## Create and run Pipeline using cli (with resources) (manual)
Tags: pipelines, e2e

Steps:
  * Create resource from tkn ""
  * Create resource from tkn ""
  * Create task from tkn ""
  * Create task from tkn ""
  * Create pipeline from tkn ""
  * Start the pipeline manually by providing the inputs manually
  * List the pipeline from tkn
  * Verify the status of pipeline "running"
  * Wait for pipelinerun to complete 
  * Verify the pipelinerun status "successfull"


## Create and run Pipeline using cli (with resources) (Automated)
Tags: pipelines, e2e

Steps:
  * Create resource from tkn ""
  * Create resource from tkn ""
  * Create task from tkn ""
  * Create task from tkn ""
  * Create pipeline from tkn ""
  * Start the pipeline from tkn ""
  * List the pipeline from tkn
  * Verify the status of pipeline "running"
  * Wait for pipelinerun to complete 
  * Verify the pipelinerun status "successfull"


## Check pipelinerun logs
Tags: Pipelines, e2e

Steps:
  * Create task from tkn "../../testdata/tasks/shell-script-task.yaml"
  * Verify taks creation status "successfull"
  * Create task from tkn "../../testdata/tasks/python-script-task.yaml"
  * Verify task creation status "successfull"
  * Create pipeline from tkn "../../testdata/cli/pipeline-script.yaml"
  * Start the pipeline from tkn
  * List the pipeline
  * Verify the status of pipeline "running"
  * List the taskruns from tkn
  * Verify the status of taskrun "running"
  * Wait for the pipelinerun to complete
  * Verify the pipelinerun status "successfull"
  * Verif the pipeline logs


## Chek pipelinerun running logs
Tags: pipelines, e2e

Steps:
  * Create task from tkn "../../testdata/tasks/shell-script-task.yaml"
  * Verify taks creation status "successfull"
  * Create task from tkn "../../testdata/tasks/python-script-task.yaml"
  * Verify task creation status "successfull"
  * Create pipeline from tkn "../../testdata/cli/pipeline-script.yaml"
  * Start the pipeline from tkn 
  * List the pipeline
  * Verify the status of pipeline "running"
  * List the taskruns from tkn
  * Verify the status of taskrun "running"
  * Verify the pipeline running logs
  * Wait for the pipelinerun to complete
  * Verify the pipelinerun status "successfull"


## Delete pipelinerun
Tags: pipelines, e2e

Steps:
  * Create task from tkn "../../testdata/tasks/shell-script-task.yaml"
  * Verify taks creation status "successfull"
  * Create task from tkn "../../testdata/tasks/python-script-task.yaml"
  * Verify task creation status "successfull"
  * Create pipeline from tkn "../../testdata/cli/pipeline-script.yaml"
  * Start the pipeline from tkn
  * List the pipeline
  * Verify the status of pipeline "running"
  * List the taskruns from tkn
  * Verify the status of taskrun "running"
  * Verify the pipeline running logs
  * Wait for the pipelinerun to complete
  * Delete pipelinerun
  * Verify pipelinerun is deleted
  * Verfiy all taskruns of pipelinerun is deleted


## Cancel pipelinerun
Tags: pipelines, e2e

Steps:
  * Create task from tkn "../../testdata/tasks/shell-script-task.yaml"
  * Verify taks creation status "successfull"
  * Create task from tkn "../../testdata/tasks/python-script-task.yaml"
  * Verify task creation status "successfull"
  * Create pipeline from tkn "../../testdata/cli/pipeline-script.yaml"
  * Start the pipeline from tkn
  * List the pipeline
  * Verify the status of pipeline "running"
  * Cancel the pipelinerun from tkn""
  * Verify pipeline run status "Cancelled"
  * Verify the taskrun status "Cancelled"


## Describe pipeline
Tags: Pipelines, e2e

Steps:
  * Create task from tkn "../../testdata/tasks/shell-script-task.yaml"
  * Verify taks creation status "successfull"
  * Create task from tkn "../../testdata/tasks/python-script-task.yaml"
  * Verify task creation status "successfull"
  * Create pipeline from tkn "../../testdata/cli/pipeline-script.yaml"
  * Start the pipeline from tkn
  * Describe the pipeline from tkn  
  * Verify the taskrun description


