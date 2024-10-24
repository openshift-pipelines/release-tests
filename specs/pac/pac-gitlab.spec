PIPELINES-30
# Pipelines As Code tests

## Configure PAC in GitLab Project: PIPELINES-30-TC01
Tags: pac, sanity, e2e
Component: PAC
Level: Integration
Type: Functional
Importance: Critical

This scenario tests configuring PAC in Public GitLab project

Steps:
  * Configure GitLab token for PAC tests
  * Verify ServiceAccount "pipeline" exist
  * Create Smee deployment
  * Configure GitLab repo and validate pipelinerun
