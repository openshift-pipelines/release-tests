# Verify Pipeline Failures

Contains negative scenarios that exercises running pipeline

Precondition:
  * Operator should be installed

## Run Pipeline with a non-existent ServiceAccount
Tags: e2e, pipeline, negative

Running a pipeline using a ServiceAccount that does not exist must fail

Steps:
  * Verify ServiceAccount "foobar" does not exist
  * Create
      |S.NO|resource_dir                      |
      |----|----------------------------------|
      |1   |testdata/negative/v1alpha1/pipelinerun.yaml|
      |2   |testdata/negative/v1beta1/pipelinerun.yaml |
  * Verify pipelinerun
       |S.NO|pipeline_run_name     |status |check_lable_propagation|
       |----|----------------------|-------|-----------------------|
       |1   |output-pipeline-run-va|Failure|no                     |
       |2   |output-pipeline-run-vb|Failure|no                     |

## Run Task with a non-existent ServiceAccount
Tags: e2e, tasks, negative

Running a task using a ServiceAccount that does not exist must fail

Steps:
  * Verify ServiceAccount "foobar" does not exist
  * Create
      |S.NO|resource_dir                                |
      |----|--------------------------------------------|
      |1   |testdata/negative/v1alpha1/pull-request.yaml|
      |2   |testdata/negative/v1beta1/pull-request.yaml |
  * Verify taskrun
       |S.NO|task_run_name |status |
       |----|--------------|-------|
       |1   |pullrequest-va|Failure|
       |2   |pullrequest-vb|Failure|
