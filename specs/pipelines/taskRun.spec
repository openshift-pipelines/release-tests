# Verify Task Run Spec
Pre condition:
  * Operator should be installed

## Create Cluster Task
Tags: e2e, taskrun

Steps:
  * Create "./testdata/taskrun/clusterTask.yaml"

## Run Cluster Task
Tags: e2e, taskrun

Steps:
 * Run "./testdata/taskrun/clusterTaskRun.yaml"
 * Verify taskrun is "successfull"

## Run Task to build, push using kaniko
Tags:  e2e, taskrun

Steps:
 * Run "./testdata/taskrun/build-push-kaniko.yaml" 
 * Verify taskrun is "successfull"
 * Verify image stream

## Run custom-volume task
Tags:  e2e, taskrun

Steps:
 * Run "./testdata/taskrun/custom-volume.yaml"
 * Verify taskrun is "successfull"

## Run pulling private image using docker-credentials
Tags:  e2e, taskrun

Steps:
 * Run "./testdata/taskrun/docker-cred.yaml"
 * Verify taskrun is "successfull"

## Run side car task
Tags:  e2e, taskrun

Steps:
 * Run "./testdata/taskrun/sidecar.yaml"
 * Verify taskrun is "successfull"

## Run task-step as script
Tags:  e2e, taskrun

Steps:
 * Run "./testdata/taskrun/step-script.yaml"
 * Verify taskrun is "successfull"

## Run workspace task
Tags:  e2e, taskrun

Steps:
 * Run "./testdata/taskrun/workspace.yaml"
 * Verify taskrun is "successfull"

## Run podTemplate task (Eg: Includes SCC/privilages to running pod)
Tags:  e2e, taskrun

Steps:
 * Run "./testdata/taskrun/podtemplate.yaml"
 * Verify taskrun is "successfull"

## Cancel Task Run
Tags:  e2e, taskrun

Steps:
 * Create "./testdata/taskrun/cancel.yaml"
 * Cancel taskrun
 * Verify taskrun is "cancelled"