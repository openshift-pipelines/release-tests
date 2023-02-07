PIPELINES-19
# Olm Openshift Pipelines operator post upgrade tests

## Verify environment after upgrade: PIPELINES-19-TC01
Tags: post-upgrade, admin
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

Creates Openshift Pipelines resources before upgrade

Steps:
  * Switch to project "releasetest-upgrade"
  * Get route for eventlistener "listener-ctb-github-push"
  * Mock post event to "github" interceptor with event-type "push", payload "testdata/triggers/github-ctb/push.json", with TLS "false"
  * Assert eventlistener response
  * Verify pipelinerun
    |S.NO|pipeline_run_name       |status    |check_label_propagation|
    |----|------------------------|----------|-----------------------|
    |1   |pipelinerun-git-push-ctb|successful|no                     |
  * Get route for eventlistener "listener-triggerref"
  * Mock post event to "github" interceptor with event-type "pull_request", payload "testdata/triggers/triggersCRD/pull-request.json", with TLS "false"
  * Assert eventlistener response
  * Verify pipelinerun
    |S.NO|pipeline_run_name       |status    |check_label_propagation|
    |----|------------------------|----------|-----------------------|
    |1   |parallel-pipelinerun    |successful|no                     |
  * Get route for eventlistener "bitbucket-listener"
  * Mock post event to "bitbucket" interceptor with event-type "refs_changed", payload "testdata/triggers/bitbucket/refs-change-event.json", with TLS "false"
  * Assert eventlistener response
  * Verify taskrun
    |S.NO|task_run_name    |status |
    |----|-----------------|-------|
    |1   |bitbucket-run    |Failure|

## Verify S2I nodejs pipeline after upgrade: PIPELINES-19-TC02
Tags: post-upgrade, e2e, clustertasks, non-admin, s2i
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Switch to project "releasetest-upgrade-s2i"
  * Get tags of the imagestream "nodejs" from namespace "openshift" and store to variable "nodejs-tags"
  * Start and verify pipeline "s2i-nodejs-pipeline" with param "VERSION" with values stored in variable "nodejs-tags" with workspace "name=source,claimName=shared-pvc"

## Verify Event listener with TLS after upgrade: PIPELINES-19-TC03
Tags: pre-upgrade, tls, triggers, admin, e2e, sanity
Component: Triggers
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Switch to project "releasetest-upgrade-tls"
  * Get route for eventlistener "listener-embed-binding"
  * Mock post event to "github" interceptor with event-type "push", payload "testdata/push.json", with TLS "true"
  * Assert eventlistener response
  * Verify pipelinerun
    |S.NO|pipeline_run_name  |status    |check_label_propagation|
    |----|-------------------|----------|-----------------------|
    |1   |simple-pipeline-run|successful|no                     |