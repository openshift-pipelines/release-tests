PIPELINES-04
# Verify Triggers with cronjob

Pre condition:
  * Validate Operator should be installed

## Create Triggers using k8s cronJob: PIPELINES-04-TC01
Tags: e2e,triggers,non-admin, sanity
Component: Triggers
Level: Integration
Type: Functional
Importance: Critical

This scenario helps you to Trigger pipelineRun, using a k8s CronJob, to implement a basic cron trigger that runs every minute

Steps:
  * Create
    |S.NO|resource_dir                                 |
    |----|---------------------------------------------|
    |1   |testdata/triggers/cron/example-pipeline.yaml |
    |2   |testdata/triggers/cron/triggerbinding.yaml   |
    |3   |testdata/triggers/cron/triggertemplate.yaml  |
    |4   |testdata/triggers/cron/eventlistener.yaml    |
  * Expose Event listener "cron-listener"
  * Verify that image stream "golang" exists
  * Create cron job with schedule "*/1 * * * *"
  * Watch for pipelinerun resources
  * Delete cron job
  * Assert no new pipelineruns created
  * Cleanup Triggers