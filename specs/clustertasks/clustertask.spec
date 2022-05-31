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
      |S.NO|pipeline_run_name|status     |check_lable_propagation|
      |----|-----------------|-----------|-----------------------|
      |1   |nodejs-ex-git-pr |successfull|no                     |

## Disable/Enable community cluster tasks: PIPELINES-14-TC02
Tags: e2e, integration, clustertasks, admin, addon
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Update addon config param "clusterTasks" with value "true" and expect message ""
  * Update addon config param "communityClusterTasks" with value "true" and expect message ""
  * Update addon config param "pipelineTemplates" with value "true" and expect message ""
  * Sleep for "10" seconds
  * "community" clustertasks are "present"
  * "tkn,openshift-client" clustertasks are "present"
  * Update addon config param "communityClusterTasks" with value "false" and expect message ""
  * Sleep for "10" seconds
  * "community" clustertasks are "not present"
  * "tkn,openshift-client" clustertasks are "present"
  * Update addon config param "communityClusterTasks" with value "true" and expect message ""
  * Sleep for "10" seconds
  * "community" clustertasks are "present"
  * "tkn,openshift-client" clustertasks are "present"