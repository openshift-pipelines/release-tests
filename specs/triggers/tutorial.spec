PIPELINES-06
# Verify triggers tutorial

Pre condition:
  * Validate Operator should be installed

## Run Triggers tutorial (by Automatically configuring users webhook to git repo): PIPELINES-06-TC01
Tags: e2e, integration, triggers, non-admin
Component: Triggers
Level: Integration
Type: Functional
Importance: Critical

This scenario helps you to configure webhook & listens to github events, on each github event it creates/triggers
openshift-pipeline Resources which helps you to deploy application (vote-app)

Steps:
  * Setup openshift-pipeline resources to create vote-app
  * Setup pipeline triggers
  * Add Event listener with github interceptor
  * Expose event listener service
  * Setup custom github webhook tasks
  * Configure webhooks
     |GitHubOrg |GitUser   |GitRepo |
     |----------|----------|--------|
     |praveen4g0|praveen4g0|vote-api|
     |praveen4g0|praveen4g0|vote-ui |
  * Mock Github push event
     |sha1                                         |head_commit   |repository                               |repositroy   |
     |---------------------------------------------|--------------|-----------------------------------------|-------------|
     |sha1=32b07065424610cff8025eb0deb12ca50088a44d|id=master|url=https://github.com/openshift-pipelines/vote-api.git|name=vote-api|
     |sha1=229cdf873cf63caf73f04ce12e7c5841462de38e|id=master|url=https://github.com/openshift-pipelines/vote-ui.git |name=vote-ui |
  * Verify creation of openshift-pipeline-resources
  * Verify resources are created with labels & event-id
  * Verify pipelinerun is "successfull"