# Verify pipeline failure status

Contains negative scenarios that exercises running pipeline

Precondition:
  * Operator should be installed

## Run Task with non-existance SA
Tags: e2e, integration, pipelines, focus

  Creates a simple Task
  Validate for failure status
  when we try to run Task with `non-existent` SA

Steps:
  * Create task
  * Ensure "foobar" SA is absent
  * Run task with "foobar" SA
  * Verify taskrun has failed


## Run Pipeline with non-existance SA

Tags: e2e, pipeline

  Creates a simple pipeline
  Validate for failure status
  when we try to run pipeline with `non-existance` SA

Steps:
  * Create pipeline
  * Run pipeline with "non-existance" SA
  * Verify pipelineRun has failed
