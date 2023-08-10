# HUB resolvers spec
Pre condition:
  * Validate Operator should be installed
## Test the functionality of hub resolvers
Tags: e2e, sanity
Component: Resolvers
Level: Integration
Type: Functional
Importance: High

Steps:
    * Verify ServiceAccount "pipeline" exist
    * Apply  
      |S.NO|resource_dir                                     |
      |----|-------------------------------------------------|
      |1   |testdata/resolvers/pipelines/git-cli-hub.yaml    |
      |2   |testdata/pvc/pvc.yaml                            |
      |3   |testdata/resolvers/pipelineruns/git-cli-hub.yaml |
    * Verify pipelinerun
      |S.NO|pipeline_run_name|status    |check_label_propagation|
      |----|-----------------|----------|-----------------------|
      |1   |hub-git-cli-run  |successful|no                     |