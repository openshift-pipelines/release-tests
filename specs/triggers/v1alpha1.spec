PIPELINES-07
# Verify Triggers

Pre condition:
  * Validate Operator should be installed

## v1alpha1 resources creation Test: PIPELINES-07-TC01
Tags: e2e, triggers, non-admin
Component: Triggers
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Create
    |S.NO|resource_dir                                 |
    |----|---------------------------------------------|
    |1   |testdata/triggers/v1alpha1-task-listener.yaml|
  * Expose Event listener "v1alpha1-task-listener"
  * Mock post event with empty payload
  * Assert eventlistener response
  * Verify taskrun
       |S.NO|task_run_name    |status |
       |----|-----------------|-------|
       |1   |v1alpha1-task-run|Success|
  * Verify taskrun "v1alpha1-task-run" label propagation
  * Cleanup Triggers
