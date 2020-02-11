# Verify Pipeline Failures

Contains negative scenarios that exercises running pipeline

Precondition:
  * Operator should be installed

## Run Task with a non-existent ServiceAccount
Tags: e2e, tasks, negative, focus

Running a task using a ServiceAccount that does not exist must fail

Steps:
  * Create task
  * Verify ServiceAccount "foobar" does not exist
  * Run task using "foobar" ServiceAccount
  * Verify taskrun has failed


## Run Pipeline with a non-existent ServiceAccount
Tags: e2e, pipelines, negative, focus

Running a Pipeline using a ServiceAccount that does not exist must fail

Steps:
  * Create pipeline
  * Verify ServiceAccount "foobar" does not exist
  * Run pipeline using "foobar" ServiceAccount
  * Verify pipelinerun has failed
