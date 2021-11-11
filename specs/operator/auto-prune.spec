PIPELINES-11

# Verify auto-prune E2E

Pre condition:
  * Validate Operator should be installed

## Verify auto prune of schedule per namespace: PIPELINES-11-TC06
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
  * Update pruner config with keep "2" schedule "*/2 * * * *" resouces "pipelinerun"
  * Annotate namespace with "operator.tekton.dev/prune.schedule=*/1 * * * *"
  * Sleep for "100" seconds
  * Assert if cronjob "resource-pruner" is "present" in namespace "targetNamespace"
  * Assert if cronjob "tekton-resource-pruner" is "present" in namespace "currentNamespace"
  * "2" of pipelineruns should be present
  * "7" of taskruns should be present
  * Update pruner config to default
  * Assert if cronjob "resource-pruner" is "not present" in namespace "targetNamespace"
  * Assert if cronjob "tekton-resource-pruner" is "not present" in namespace "currentNamespace"