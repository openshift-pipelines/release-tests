PIPELINES-09
# Olm Openshift Pipelines operator specs

## Install openshift-pipelines operator: PIPELINES-09-TC01
Tags: install, admin, sanity
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

Installs `openshift-pipelines` operator using olm

Steps:
  * Subscribe to operator
  * Wait for TektonConfig CR availability
  * Validate pipelines deployment
  * Validate triggers deployment
  * Validate PAC deployment
  * Validate tkn server cli deployment
  * Verify TektonAddons Install status
  * Validate RBAC
  * Validate default auto prune cronjob in target namespace
  * Validate tektoninstallersets status
  * Validate tektoninstallersets names

## Upgrade openshift-pipelines operator: PIPELINES-09-TC02
Tags: upgrade, admin
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

Installs `openshift-pipelines` operator using olm

Steps:
  * Upgrade operator subscription
  * Wait for TektonConfig CR availability
  * Validate Operator should be installed
  * Validate RBAC

## Uninstall openshift-pipelines operator: PIPELINES-09-TC03
Tags: uninstall, admin
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

Uninstalls `openshift-pipelines` operator using olm
Steps:
  * Validate Operator should be installed
  * Uninstall Operator

## Setup environment for upgrade test: PIPELINES-09-TC04
Tags: pre-upgrade, admin
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

Creates Openshift Pipelines resources before upgrade

Steps:
  * Create project "test-upgrade" and update clients
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                          |
      |----|------------------------------------------------------|
      |1   |testdata/v1beta1/pipelinerun/task_results_example.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name |status    |check_label_propagation|
      |----|------------------|----------|-----------------------|
      |1   |task-level-results|successful|no                     |
  * Create
    |S.NO|resource_dir                                                      |
    |----|------------------------------------------------------------------|
    |1   |testdata/triggers/github-ctb/Embeddedtriggertemplate-git-push.yaml|
    |2   |testdata/triggers/github-ctb/eventlistener-ctb-git-push.yaml      |
  * Verify ServiceAccount "pipeline" exist
  * Create & Link secret "github-secret" to service account "pipeline"
  * Expose Event listener "listener-clustertriggerbinding-github-push"
  * Mock post event to "github" interceptor with event-type "push", payload "testdata/triggers/github-ctb/push.json", with TLS "false"
  * Assert eventlistener response
  * Verify pipelinerun
    |S.NO|pipeline_run_name       |status    |check_label_propagation|
    |----|------------------------|----------|-----------------------|
    |1   |pipelinerun-git-push-ctb|successful|no                     |
  * Create
    |S.NO|resource_dir                                               |
    |----|-----------------------------------------------------------|
    |1   |testdata/triggers/triggersCRD/eventlistener-triggerref.yaml|
    |2   |testdata/triggers/triggersCRD/trigger.yaml                 |
    |3   |testdata/triggers/triggersCRD/triggerbindings.yaml         |
    |4   |testdata/triggers/triggersCRD/triggertemplate.yaml         |
    |5   |testdata/triggers/triggersCRD/pipeline.yaml                |
  * Expose Event listener "listener-triggerref"
  * Mock post event to "github" interceptor with event-type "pull_request", payload "testdata/triggers/triggersCRD/pull-request.json", with TLS "false"
  * Assert eventlistener response
  * Verify pipelinerun
    |S.NO|pipeline_run_name       |status    |check_label_propagation|
    |----|------------------------|----------|-----------------------|
    |1   |parallel-pipelinerun    |successful|no                     |
  * Create
    |S.NO|resource_dir                                      |
    |----|--------------------------------------------------|
    |1   |testdata/triggers/gitlab/gitlab-push-listener.yaml|
  * Create & Link secret "gitlab-secret" to service account "pipeline"
  * Expose Event listener "gitlab-listener"
  * Mock post event to "gitlab" interceptor with event-type "Push Hook", payload "testdata/triggers/gitlab/gitlab-push-event.json", with TLS "false"
  * Assert eventlistener response
  * Verify taskrun
    |S.NO|task_run_name |status    |
    |----|--------------|----------|
    |1   |gitlab-run    |successful|
  * Create
    |S.NO|resource_dir                                                        |
    |----|--------------------------------------------------------------------|
    |1   |testdata/triggers/bitbucket/bitbucket-eventlistener-interceptor.yaml|
  * Create & Link secret "bitbucket-secret" to service account "pipeline"
  * Expose Event listener "bitbucket-listener"
  * Mock post event to "bitbucket" interceptor with event-type "refs_changed", payload "testdata/triggers/bitbucket/refs-change-event.json", with TLS "false"
  * Assert eventlistener response
  * Verify taskrun
    |S.NO|task_run_name    |status |
    |----|-----------------|-------|
    |1   |bitbucket-run    |Failure|
  
  