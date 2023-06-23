PIPELINES-24
# Git resolvers spec

Pre condition:
  * Validate Operator should be installed
## Test the functionality of git resolvers: PIPELINES-24-TC01
Tags: e2e, sanity
Component: Resolvers
Level: Integration
Type: Functional
Importance: High

Steps:
    * Verify ServiceAccount "pipeline" exist
    * Create 
      |S.NO|resource_dir                                                     |
      |----|-----------------------------------------------------------------|
      |1   |testdata/resolvers/pipelineruns/git-resolver-pipelinerun.yaml    |
    * Verify pipelinerun
      |S.NO|pipeline_run_name                  |status      |check_label_propagation  |
      |----|-----------------------------------|--------------------------------------|
      |1   |git-resolver-pipelinerun           |successful  |no                       |

## Test the functionality of git resolvers with authentication: PIPELINES-24-TC01
Tags: e2e
Component: Resolvers
Level: Integration
Type: Functional
Importance: High

Steps:
    * Verify ServiceAccount "pipeline" exist
    * Create 
      |S.NO|resource_dir                                                     |
      |----|-----------------------------------------------------------------|
      |1   |testdata/resolvers/pipelineruns/git-resolver-pipelinerun-private.yaml        | 
    * Verify pipelinerun
      |S.NO|pipeline_run_name                  |status      |check_label_propagation  |
      |----|-----------------------------------|--------------------------------------|
      |1   |git-resolver-pipelinerun-private   |successful  |no                       |