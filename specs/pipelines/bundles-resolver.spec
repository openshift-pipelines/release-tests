PIPELINES-25
# Bundles resolver spec

Precondition:
  * Validate Operator should be installed

## Test the functionality of bundles resolver: PIPELINES-25-TC01
Tags: e2e
Component: Resolvers
Level: Integration
Type: Functional
Importance: High

Steps:
    * Create
      |S.NO|resource_dir                                                            |
      |----|------------------------------------------------------------------------|
      |1   |testdata/resolvers/pipelineruns/bundles-resolver-pipelinerun.yaml       |
    * Verify pipelinerun
      |S.NO|pipeline_run_name                  |status      |
      |----|-----------------------------------|------------|
      |1   |bundles-resolver-pipelinerun       |successful  |

## Test the functionality of bundles resolver with parameter: PIPELINES-25-TC02
Tags: e2e, sanity
Component: Resolvers
Level: Integration
Type: Functional
Importance: High

Steps:
    * Create
      |S.NO|resource_dir                                                               |
      |----|---------------------------------------------------------------------------|
      |1   |testdata/resolvers/pipelineruns/bundles-resolver-pipelinerun-param.yaml    |
    * Verify pipelinerun
      |S.NO|pipeline_run_name                  |status      |
      |----|-----------------------------------|------------|
      |1   |bundles-resolver-pipelinerun-param |successful  |