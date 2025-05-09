PIPELINES-03
# Verify Pipeline E2E spec

Pre condition:
  * Validate Operator should be installed

## Run sample pipeline: PIPELINES-03-TC01
Tags: e2e, pipelines, non-admin
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Run a sample pipeline that has 2 tasks:
  1. create a file
  2. read file content created by above task
and verify that it runs succesfully

Steps:
  * Verify that image stream "golang" exists
  * Create
      |S.NO|resource_dir                                  |
      |----|----------------------------------------------|
      |1   |testdata/pvc/pvc.yaml                         |
      |2   |testdata/v1beta1/pipelinerun/pipelinerun.yaml |
  * Verify pipelinerun
      |S.NO|pipeline_run_name       |status    |
      |----|------------------------|----------|
      |1   |output-pipeline-run-v1b1|successful|

## Pipelinerun Timeout failure Test: PIPELINES-03-TC04
Tags: e2e, pipelines, non-admin, sanity
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Create
      |S.NO|resource_dir                                         |
      |----|-----------------------------------------------------|
      |1   |testdata/v1beta1/pipelinerun/pipelineruntimeout.yaml |
  * Verify pipelinerun
      |S.NO|pipeline_run_name|status             |
      |----|-----------------|-------------------|
      |1   |pear             |timeout            |

## Configure execution results at the Task level Test: PIPELINES-03-TC05
Tags: e2e, integration, pipelines, non-admin, sanity
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Create
      |S.NO|resource_dir                                          |
      |----|------------------------------------------------------|
      |1   |testdata/v1beta1/pipelinerun/task_results_example.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name |status    |
      |----|------------------|----------|
      |1   |task-level-results|successful|

## Cancel pipelinerun Test: PIPELINES-03-TC06
Tags: e2e, integration, pipelines, non-admin, sanity
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Create
      |S.NO|resource_dir                                 |
      |----|---------------------------------------------|
      |1   |testdata/pvc/pvc.yaml                        |
      |2   |testdata/v1beta1/pipelinerun/pipelinerun.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name       |status   |
      |----|------------------------|---------|
      |1   |output-pipeline-run-v1b1|cancelled|

## Pipelinerun with pipelinespec and taskspec (embedded pipelinerun tests): PIPELINES-03-TC07
Tags: e2e, integration, pipelines, non-admin
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Create
      |S.NO|resource_dir                                                                |
      |----|----------------------------------------------------------------------------|
      |1   |testdata/v1beta1/pipelinerun/pipelinerun-with-pipelinespec-and-taskspec.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name                        |status    |
      |----|-----------------------------------------|----------|
      |1   |pipelinerun-with-pipelinespec-taskspec-vb|successful|

## Pipelinerun with large result: PIPELINES-03-TC08
Tags: e2e, integration, pipelines, non-admin, results, sanity
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical
CustomerScenario: yes

Steps:
  * Create
      |S.NO|resource_dir                                                   |
      |----|---------------------------------------------------------------|
      |1   |testdata/v1beta1/pipelinerun/pipelinerun-with-large-result.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name|status    |
      |----|-----------------|----------|
      |1   |result-test-run  |successful|
