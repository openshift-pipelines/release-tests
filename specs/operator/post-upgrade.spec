PIPELINES-19
# Olm Openshift Pipelines operator post upgrade tests
Pre condition:
  * Validate Operator should be installed

## Verify environment after upgrade: PIPELINES-19-TC01
Tags: post-upgrade, admin
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

Creates Openshift Pipelines resources before upgrade

Steps:
  * Switch to project "releasetest-upgrade-triggers"
  * Get route for eventlistener "listener-ctb-github-push"
  * Mock post event to "github" interceptor with event-type "push", payload "testdata/triggers/github-ctb/push.json", with TLS "false"
  * Assert eventlistener response
  * Verify pipelinerun
    |S.NO|pipeline_run_name       |status    |
    |----|------------------------|----------|
    |1   |pipelinerun-git-push-ctb|successful|
  * Get route for eventlistener "listener-triggerref"
  * Mock post event to "github" interceptor with event-type "pull_request", payload "testdata/triggers/triggersCRD/pull-request.json", with TLS "false"
  * Assert eventlistener response
  * Verify pipelinerun
    |S.NO|pipeline_run_name       |status    |
    |----|------------------------|----------|
    |1   |parallel-pipelinerun    |successful|
  * Get route for eventlistener "bitbucket-listener"
  * Mock post event to "bitbucket" interceptor with event-type "refs_changed", payload "testdata/triggers/bitbucket/refs-change-event.json", with TLS "false"
  * Assert eventlistener response
  * Verify taskrun
    |S.NO|task_run_name    |status |
    |----|-----------------|-------|
    |1   |bitbucket-run    |Failure|

## Verify Event listener with TLS after upgrade: PIPELINES-19-TC03
Tags: post-upgrade, tls, triggers, admin, e2e, sanity
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
    |S.NO|pipeline_run_name  |status    |
    |----|-------------------|----------|
    |1   |simple-pipeline-run|successful|

## Verify secret is linked to SA even after upgrade: PIPELINES-19-TC04
Tags: post-upgrade, e2e, clustertasks, non-admin, git-clone, sanity
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical
CustomerScenario: yes

Steps:
  * Switch to project "releasetest-upgrade-pipelines"
  * Verify ServiceAccount "pipeline" exist
  * Create
      | S.NO | resource_dir                                                          |
      |------|-----------------------------------------------------------------------|
      | 1    | testdata/ecosystem/pipelineruns/git-clone-read-private.yaml |
  * Verify pipelinerun
      | S.NO | pipeline_run_name                   | status     |
      |------|-------------------------------------|------------|
      | 1    | git-clone-read-private-pipeline-run | successful |
    
## Verify S2I golang pipeline after upgrade: PIPELINES-19-TC05
Tags: post-upgrade, e2e, clustertasks, non-admin, s2i
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Switch to project "releasetest-upgrade-s2i"
  * Get tags of the imagestream "golang" from namespace "openshift" and store to variable "golang-tags"
  * Start and verify pipeline "s2i-go-pipeline" with param "VERSION" with values stored in variable "golang-tags" with workspace "name=source,claimName=shared-pvc"

## Validate olm skiprange post upgrade: PIPELINES-19-TC06
Tags: post-upgrade, olm
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
 * Get olm-skip-range "post-upgrade" and save to field "post-upgrade-olm-skip-range" in file "testdata/olm/skiprange.json"
 * Validate skipRange diff between fields "pre-upgrade-olm-skip-range" and "post-upgrade-olm-skip-range" in file "testdata/olm/skiprange.json"