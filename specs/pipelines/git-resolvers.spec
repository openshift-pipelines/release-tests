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
    * Create 
      |S.NO|resource_dir                                                     |
      |----|-----------------------------------------------------------------|
      |1   |testdata/resolvers/pipelineruns/git-resolver-pipelinerun.yaml    |
    * Verify pipelinerun
      |S.NO|pipeline_run_name                  |status      |
      |----|-----------------------------------|------------|
      |1   |git-resolver-pipelinerun           |successful  |

## Test the functionality of git resolvers with authentication: PIPELINES-24-TC01
Tags: e2e
Component: Resolvers
Level: Integration
Type: Functional
Importance: High

Steps:
    * Create 
      |S.NO|resource_dir                                                     |
      |----|-----------------------------------------------------------------|
      |1   |testdata/resolvers/pipelineruns/git-resolver-pipelinerun-private.yaml        | 
    * Verify pipelinerun
      |S.NO|pipeline_run_name                  |status      |
      |----|-----------------------------------|------------|
      |1   |git-resolver-pipelinerun-private   |successful  |