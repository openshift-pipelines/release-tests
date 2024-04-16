PIPELINES-06
# Verify triggers tutorial

Pre condition:
  * Validate Operator should be installed

## Run Triggers tutorial (by Automatically configuring users webhook to git repo): PIPELINES-06-TC01
Tags: e2e, integration, triggers, non-admin, to-do
Component: Triggers
Level: Integration
Type: Functional
Importance: Critical

This scenario tests the pipeline tutorials (https://github.com/openshift/pipelines-tutorial) pipelines related resources 

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create remote
      |S.NO|resource_dir                                                                                                                             |
      |----|-----------------------------------------------------------------------------------------------------------------------------------------|
      |1   |https://raw.githubusercontent.com/openshift/pipelines-tutorial/{OSP_TUTORIAL_BRANCH}/01_pipeline/01_apply_manifest_task.yaml             |
      |2   |https://raw.githubusercontent.com/openshift/pipelines-tutorial/{OSP_TUTORIAL_BRANCH}/01_pipeline/02_update_deployment_task.yaml          |
      |3   |https://raw.githubusercontent.com/openshift/pipelines-tutorial/{OSP_TUTORIAL_BRANCH}/01_pipeline/03_persistent_volume_claim.yaml         |
      |4   |https://raw.githubusercontent.com/openshift/pipelines-tutorial/{OSP_TUTORIAL_BRANCH}/01_pipeline/04_pipeline.yaml                        |
      |5   |https://raw.githubusercontent.com/openshift/pipelines-tutorial/{OSP_TUTORIAL_BRANCH}/02_pipelinerun/01_build_deploy_api_pipelinerun.yaml |
   * Verify pipelinerun
      |S.NO|pipeline_run_name           |status    |check_label_propagation|
      |----|----------------------------|----------|-----------------------|
      |1   |build-deploy-api-pipelinerun|successful|no                     |
   * Create remote
      |S.NO|resource_dir                                                                                                                           |
      |----|---------------------------------------------------------------------------------------------------------------------------------------|
      |1   |https://raw.githubusercontent.com/openshift/pipelines-tutorial/{OSP_TUTORIAL_BRANCH}/02_pipelinerun/02_build_deploy_ui_pipelinerun.yaml|
   * Verify pipelinerun
      |S.NO|pipeline_run_name          |status    |check_label_propagation|
      |----|---------------------------|----------|-----------------------|
      |1   |build-deploy-ui-pipelinerun|successful|no                     |
   * Get route url of the route "pipelines-vote-ui"
   * Validate route url
