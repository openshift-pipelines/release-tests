PIPELINES-29
# Pipelines As Code tests

## Enable/Disable PAC: PIPELINES-29-TC01
Tags: pac, sanity
Component: PAC
Level: Integration
Type: Functional
Importance: Critical

This scenario tests enable/disable of pipelines as code from tektonconfig custom resource

Steps:
  * Create
    |S.NO|resource_dir                                      |
    |----|--------------------------------------------------|
    |1   |testdata/triggers/gitlab/gitlab-push-listener.yaml|
  * Verify ServiceAccount "pipeline" exist
  * Create & Link secret "gitlab-secret" to service account "pipeline"
  * Expose Event listener "gitlab-listener"
  * Create Smee Deployment with "gitlab-listener"
  * Configure Gitlab repo
