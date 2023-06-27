PIPELINES-25
# Bundles resolver spec

Precondition:
  * Validate Operator should be installed

## Test the functionality of bundles resolver#1: PIPELINES-25-TC01
Tags: e2e, sanity
Component: Resolvers
Level: Integration
Type: Functional
Importance: High

Steps:
    * Create
      |S.NO|resource_dir                                                            |
      |----|------------------------------------------------------------------------|
      |1   |testdata/resolvers/pipelineruns/bundles-resolver-pipelinerun.yaml       |
    * Verify ServiceAccount "pipeline" exist
    * Verify pipelinerun
      |S.NO|pipeline_run_name                  |status      |check_label_propagation  |
      |----|-----------------------------------|--------------------------------------|
      |1   |bundles-resolver-pipelinerun       |successful  |no                       |

## Test the functionality of bundles resolver with parametr#2: PIPELINES-25-TC02
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
    * Verify ServiceAccount "pipeline" exist
    * Verify pipelinerun
      |S.NO|pipeline_run_name                  |status      |check_label_propagation  |
      |----|-----------------------------------|--------------------------------------|
      |1   |bundles-resolver-pipelinerun-param |successful  |no                       |