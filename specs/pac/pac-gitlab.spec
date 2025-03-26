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
  * Setup Gitlab Client
  * Validate PAC Info Install
  * Verify ServiceAccount "pipeline" exist
  * Create Smee deployment
  * Configure GitLab repo for "pull_request" in "main"
  * Configure PipelineRun
  * Validate PipelineRun for "success"
  * Cleanup PAC

## Configure PAC in GitLab Project: PIPELINES-30-TC02
Tags: pac, e2e
Component: PAC
Level: Integration
Type: Functional
Importance: Critical

This scenario tests on-label annotation in PAC

Steps:
  * Setup Gitlab Client
  * Verify ServiceAccount "pipeline" exist
  * Create Smee deployment
  * Configure GitLab repo for "pull_request" in "main"
  * Update Annotation "pipelinesascode.tekton.dev/on-label" with "[bug]"
  * Configure PipelineRun
  * "0" pipelinerun(s) should be present within "10" seconds
  * Add Label Name "bug" with "red" color with description "Identify a Issue"
  * Validate PipelineRun for "success"
  * Cleanup PAC

## Configure PAC in GitLab Project: PIPELINES-30-TC03
Tags: pac, e2e
Component: PAC
Level: Integration
Type: Functional
Importance: Critical

This scenario tests on-comment annotation in PAC

Steps:
  * Setup Gitlab Client
  * Verify ServiceAccount "pipeline" exist
  * Create Smee deployment
  * Configure GitLab repo for "pull_request" in "main"
  * Update Annotation "pipelinesascode.tekton.dev/on-comment" with "^/hello-world"
  * Configure PipelineRun
  * Validate PipelineRun for "success"
  * Add Comment "/hello-world" in MR 
  * "2" pipelinerun(s) should be present within "10" seconds
  * Validate PipelineRun for "success"
  * Cleanup PAC
