PIPELINES-11

# Verify auto-prune E2E

Pre condition:
  * Validate Operator should be installed

## Verify auto prune for taskrun: PIPELINES-11-TC01
Tags: e2e, integration, pipelines, non-admin
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                             |
      |----|-----------------------------------------|
      |1   |testdata/pruner/task/task-for-pruner.yaml|
  * Update pruner config with keep "2" schedule "*/1 * * * *" resouces "taskrun"
  * Sleep for "120" seconds
  * Assert if cronjob "resource-pruner" is "present" in namespace "openshift-pipelines"
  * "2" of taskruns should be present
  * Update pruner config to default
  * Assert if cronjob "resource-pruner" is "not present" in namespace "openshift-pipelines"

## Verify auto prune for pipelinerun: PIPELINES-11-TC02
Tags: e2e, integration, pipelines, non-admin
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                     |
      |----|-------------------------------------------------|
      |1   |testdata/pruner/pipeline/pipeline-for-pruner.yaml|
  * Update pruner config with keep "2" schedule "*/1 * * * *" resouces "pipelinerun"
  * Sleep for "120" seconds
  * Assert if cronjob "resource-pruner" is "present" in namespace "openshift-pipelines"
  * "2" of pipelineruns should be present
  * Update pruner config to default
  * Assert if cronjob "resource-pruner" is "not present" in namespace "openshift-pipelines"


## Verify auto prune for pipelinerun and taskrun: PIPELINES-11-TC03
Tags: e2e, integration, pipelines, non-admin
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                     |
      |----|-------------------------------------------------|
      |1   |testdata/pruner/pipeline/pipeline-for-pruner.yaml|
      |2   |testdata/pruner/task/task-for-pruner.yaml        |
  * Update pruner config with keep "2" schedule "*/1 * * * *" resouces "pipelinerun,taskrun"
  * Sleep for "120" seconds
  * Assert if cronjob "resource-pruner" is "present" in namespace "openshift-pipelines"
  * "2" of pipelineruns should be present
  * "2" of taskruns should be present
  * Update pruner config to default
  * Assert if cronjob "resource-pruner" is "not present" in namespace "openshift-pipelines"