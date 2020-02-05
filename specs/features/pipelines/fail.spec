# Verify pipeline failure status

Contains negative scenarios that exercises running pipeline

Precondition:
  * Operator should be installed

## Run Task with a non-existent ServiceAccount
Tags: e2e, integration, pipelines, negative, focus

Running a task using a ServiceAccount that does not exist must fail

Steps:
  * Create task
  * Verify ServiceAccount "foobar" does not exist
  * Run task using "foobar" ServiceAccount
  * Verify taskrun has failed


## Run Pipeline with a non-existent ServiceAccount

Tags: e2e, pipeline

Running a Pipeline using a ServiceAccount that does not exist must fail

Steps:
  * Create pipeline
  * Run pipeline with "non-existance" SA
  * Verify pipelineRun has failed
