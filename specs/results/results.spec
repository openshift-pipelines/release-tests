PIPELINES-26
# Results pvc tests

Precondition:
* Validate Operator should be installed

## Test Tekton results with TaskRun: PIPELINES-26-TC01
Tags: results, e2e, taskrun
Component: Results
Level: Integration
Type: Functional
Importance: Critical

Steps:
* Apply
   |S.NO|resource_dir                 |
   |----|-----------------------------|
   |1   |testdata/results/taskrun.yaml|
* Verify taskrun
   |S.NO|pipeline_run_name|status    |
   |----|-----------------|----------|
   |1   |results-task     |successful|
* Verify "taskrun" Results records
* Verify "taskrun" Results logs

## Test Tekton results with PipelineRun: PIPELINES-26-TC02
Tags: results, e2e, pipelinerun
Component: Results
Level: Integration
Type: Functional
Importance: Critical

Steps:
* Apply
   |S.NO|resource_dir                     |
   |----|---------------------------------|
   |1   |testdata/results/pipeline.yaml   |
   |2   |testdata/results/pipelinerun.yaml|
* Verify pipelinerun
   |S.NO|pipeline_run_name|status    |check_label_propagation|
   |----|-----------------|----------|-----------------------|
   |1   |pipeline-results |successful|no                     |
* Verify "pipelinerun" Results records
* Verify "pipelinerun" Results logs