PIPELINES-35
# Pipelines As Code tests

## Configure PAC in GitHub Project: PIPELINES-35-TC01
Tags: pac, sanity, e2e
Component: PAC
Level: Integration
Type: Functional
Importance: Critical

This scenario tests configuring PAC with push and pull_request events in GitHub.

Steps:
//   * Validate PAC Info Install
  * Setup Github Client
  * Create Smee deployment
  * Configure GitHub repo for "pull_request" in "main"
  * Configure GitHub repo for "push" in "main"
  * Configure PipelineRun
  * Validate "pull_request" PipelineRun for "success"
  * Trigger push event on main branch
  * Validate "push" PipelineRun for "success"
  * Cleanup PAC


