PIPELINES-31
# http resolvers spec

Pre condition:
  * Validate Operator should be installed
  
## Test the functionality of http resolvers: PIPELINES-31-TC01
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
      |1   |testdata/resolvers/pipelineruns/http-resolver-pipelinerun.yaml   |
    * Verify pipelinerun
      |S.NO|pipeline_run_name                  |status      |check_label_propagation  |
      |----|-----------------------------------|--------------------------------------|
      |1   |http-resolver-pipelinerun           |successful  |no                      |