PIPELINES-29
# Pipelines As Code tests

## Enable/Disable PAC: PIPELINES-29-TC01
Tags: pac, sanity, e2e
Component: PAC
Level: Integration
Type: Functional
Importance: Critical

This scenario tests enable/disable of pipelines as code from tektonconfig custom resource

Steps:
  * Configure Gitlab token for PAC tests
  * Verify ServiceAccount "pipeline" exist
  * Create Smee Deployment
  * Configure & Validate Gitlab repo for pipelinerun
