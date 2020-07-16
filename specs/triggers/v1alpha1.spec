# Verify Triggers

Pre condition:
  * Operator should be installed

## v1alpha1 resources creation Test
Tags: e2e, triggers

Steps:
  * Create
    |S.NO|resource_dir                                 |
    |----|---------------------------------------------|
    |1   |testdata/triggers/v1alpha1-task-listener.yaml|
  * Expose Event listener "v1alpha1-task-listener"
  * Mock get event
  * Verify taskrun
       |S.NO|task_run_name    |status |
       |----|-----------------|-------|
       |1   |v1alpha1-task-run|Success|
  * Verify taskrun "v1alpha1-task-run" label propagation
