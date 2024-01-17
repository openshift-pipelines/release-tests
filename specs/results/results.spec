PIPELINES-26
# Results pvc tests

## Test Tekton Results: PIPELINES-26-TC01
Tags: results, e2e
Component: Results
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Create Results route
  * Ensure that Tekton Results is ready
  * Create project "results-testing" 
  * Apply in namespace "results-testing"
    | S.NO | resource_dir                             |
    |------|------------------------------------------|
    | 1    | testdata/results/taskrun.yaml  |
  * Verify taskrun
    |S.NO|pipeline_run_name           |status    |
    |----|----------------------------|----------|
    |1   |results-task                |successful|
  * Verify "taskrun" Results records
  * Verify "taskrun" Results logs 
  * Apply in namespace "results-testing"
    | S.NO | resource_dir                      |
    |------|-----------------------------------|
    | 1    | testdata/results/pipeline.yaml    |
    | 2    | testdata/results/pipelinerun.yaml |
  * Verify pipelinerun
    |S.NO|pipeline_run_name     |status    |check_label_propagation|
    |----|----------------------|----------|-----------------------|
    |1   |pipeline-results      |successful|no                     |
  * Verify "pipelinerun" Results records
  * Verify "pipelinerun" Results logs 
  * Delete "taskrun" named "results-task"
  * Delete "pipelinerun" named "pipeline-results"