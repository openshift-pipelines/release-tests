PIPELINES-14
# Verify Clustertasks E2E spec

Pre condition:
  * Validate Operator should be installed


## S2I nodejs pipelinerun: PIPELINES-14-TC01
Tags: e2e, integration, clustertasks, non-admin, s2i
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                            |
      |----|--------------------------------------------------------|
      |1   |testdata/v1beta1/pipelinerun/s2i-nodejs-pipelinerun.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name|status    |check_label_propagation|
      |----|-----------------|----------|-----------------------|
      |1   |nodejs-ex-git-pr |successful|no                     |
