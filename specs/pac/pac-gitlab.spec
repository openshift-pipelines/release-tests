PIPELINES-30
# Pipelines As Code tests

## Configure PAC in GitLab Project: PIPELINES-30-TC01
Tags: pac, sanity, e2e
Component: PAC
Level: Integration
Type: Functional
Importance: Critical

This scenario tests configuring PAC with push and pull_request events

Steps:
  * Validate PAC Info Install
  * Setup Gitlab Client
  * Create Smee deployment
  * Configure GitLab repo for "pull_request" in "main"
  * Configure GitLab repo for "push" in "main"
  * Configure PipelineRun
  * Validate "pull_request" PipelineRun for "success"
  * Trigger push event on main branch
  * Validate "push" PipelineRun for "success"
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
  * Create Smee deployment
  * Configure GitLab repo for "pull_request" in "main"
  * Update Annotation "pipelinesascode.tekton.dev/on-label" with "[bug]"
  * Configure PipelineRun
  * "0" pipelinerun(s) should be present within "10" seconds
  * Add Label Name "bug" with "red" color with description "Identify a Issue"
  * Validate "pull_request" PipelineRun for "success"
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
  * Create Smee deployment
  * Configure GitLab repo for "pull_request" in "main"
  * Update Annotation "pipelinesascode.tekton.dev/on-comment" with "^/hello-world"
  * Configure PipelineRun
  * Validate Pull Request PipelineRun for "success"
  * Add Comment "/hello-world" in MR 
  * "2" pipelinerun(s) should be present within "10" seconds
  * Validate "pull_request" PipelineRun for "success"
  * Cleanup PAC
