# Verify Triggers with cronJob 

Pre condition:
  * Operator should be installed

## Create Triggers using k8s cronJob
Tags: triggers, focus, ignore

This scenario helps you to Trigger pipelineRun, using a k8s CronJob, to implement a basic cron trigger that runs every minute

Steps:
  * Create "./testdata/triggers/eventlisteners/cron/template.yaml"
  * Create "./testdata/triggers/eventlisteners/cron/trigger-binding.yaml"
  * Create "./testdata/triggers/eventlisteners/role-resources/rbac.yaml"
  * Add Event listener
  * Expose event listener service
  * Create Cron Job to trigger eventlistener, every 1 minute
  * Verify creation of openshift-pipeline-resources
  * Verify resources are created with labels & event-id
  * Verify pipelinerun is "successfull"