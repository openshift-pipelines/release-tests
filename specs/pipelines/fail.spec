PIPELINES-02
# Verify Pipeline Failures

Contains negative scenarios that exercises running pipeline

Precondition:
  * Validate Operator should be installed

## Run Pipeline with a non-existent ServiceAccount: PIPELINES-02-TC01
Tags: e2e, pipeline, negative, non-admin, sanity
Component: Pipelines
Pos/Neg: Negative
Level: Integration
Type: Functional
Importance: Critical

Running a pipeline using a ServiceAccount that does not exist must fail

Steps:
  * Verify ServiceAccount "foobar" does not exist
  * Create
      |S.NO|resource_dir                               |
      |----|-------------------------------------------|
      |1   |testdata/negative/v1beta1/pipelinerun.yaml |
  * Verify pipelinerun
       |S.NO|pipeline_run_name     |status |check_label_propagation|
       |----|----------------------|-------|-----------------------|
       |1   |output-pipeline-run-vb|Failure|no                     |

## Run Task with a non-existent ServiceAccount: PIPELINES-02-TC02
Tags: e2e, tasks, negative, non-admin
Component: Pipelines
Pos/Neg: Negative
Level: Integration
Type: Functional
Importance: Critical

Running a task using a ServiceAccount that does not exist must fail

Steps:
  * Verify ServiceAccount "foobar" does not exist
  * Create
      |S.NO|resource_dir                                |
      |----|--------------------------------------------|
      |1   |testdata/negative/v1beta1/pull-request.yaml |
  * Verify taskrun
       |S.NO|task_run_name |status |
       |----|--------------|-------|
       |1   |pullrequest-vb|Failure|
