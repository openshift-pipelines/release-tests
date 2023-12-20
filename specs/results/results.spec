PIPELINES-26
# Results pvc tests

## Test Tekton Results: PIPELINES-26-TC01
Tags: results, sanity
Component: Results
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Ensure that tekton results is ready
  * Create project "results-testing" 
  * Apply in namespace "results-testing"
    | S.NO | resource_dir                             |
    |------|------------------------------------------|
    | 1    | testdata/results/task-output-image.yaml  |
  * Get "tr" logs and annotations
  * Apply in namespace "results-testing"
    | S.NO | resource_dir                   |
    |------|--------------------------------|
    | 1    | testdata/results/pipeline.yaml |
  * Get "pr" logs and annotations  