PIPELINES-18
# Openshift Pipelines pre upgrade specs

## Setup environment for upgrade test: PIPELINES-18-TC01
Tags: pre-upgrade, admin
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

Creates Openshift Pipelines resources before upgrade

Steps:
  * Create project "releasetest-upgrade-triggers"
  * Verify ServiceAccount "pipeline" exist
  * Create
    |S.NO|resource_dir                                                      |
    |----|------------------------------------------------------------------|
    |1   |testdata/triggers/github-ctb/Embeddedtriggertemplate-git-push.yaml|
    |2   |testdata/triggers/github-ctb/eventlistener-ctb-git-push.yaml      |
  * Verify that image stream "golang" exists
  * Create & Link secret "github-secret" to service account "pipeline"
  * Expose Event listener "listener-ctb-github-push"
  * Mock post event to "github" interceptor with event-type "push", payload "testdata/triggers/github-ctb/push.json", with TLS "false"
  * Assert eventlistener response
  * Verify pipelinerun
    |S.NO|pipeline_run_name       |status    |
    |----|------------------------|----------|
    |1   |pipelinerun-git-push-ctb|successful|
  * Delete "pipelinerun" named "pipelinerun-git-push-ctb"
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
    |S.NO|pipeline_run_name       |status    |
    |----|------------------------|----------|
    |1   |parallel-pipelinerun    |successful|
  * Delete "pipelinerun" named "parallel-pipelinerun"
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
  * Delete "taskrun" named "bitbucket-run"

## Setup Eventlistener with TLS enabled pre upgrade: PIPELINES-18-TC03
Tags: pre-upgrade, tls, triggers, admin, e2e, sanity
Component: Triggers
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Create project "releasetest-upgrade-tls"
  * Enable TLS config for eventlisteners
  * Create
    |S.NO|resource_dir                                                        |
    |----|--------------------------------------------------------------------|
    |1   |testdata/triggers/sample-pipeline.yaml                              |
    |2   |testdata/triggers/triggerbindings/triggerbinding.yaml               |
    |3   |testdata/triggers/triggertemplate/triggertemplate.yaml              |
    |4   |testdata/triggers/eventlisteners/eventlistener-embeded-binding.yaml |
  * Expose Event listener for TLS "listener-embed-binding"
  * Mock post event to "github" interceptor with event-type "push", payload "testdata/push.json", with TLS "true"
  * Assert eventlistener response
  * Verify pipelinerun
    |S.NO|pipeline_run_name  |status    |
    |----|-------------------|----------|
    |1   |simple-pipeline-run|successful|
  * Delete "pipelinerun" named "simple-pipeline-run"

## Setup link secret to pipeline SA PIPELINES-18-TC04
Tags: pre-upgrade, e2e, clustertasks, non-admin, git-clone, sanity
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical
CustomerScenario: yes

Steps:
  * Create project "releasetest-upgrade-pipelines"
  * Verify ServiceAccount "pipeline" exist
  * Create
      | S.NO | resource_dir                                                       |
      |------|--------------------------------------------------------------------|
      | 1    | testdata/ecosystem/pipelines/git-clone-read-private.yaml |
      | 2    | testdata/pvc/pvc.yaml                                              |
      | 3    | testdata/ecosystem/secrets/ssh-key.yaml                  |
  * Link secret "ssh-key" to service account "pipeline"
  * Create
      | S.NO | resource_dir                                                          |
      |------|-----------------------------------------------------------------------|
      | 1    | testdata/ecosystem/pipelineruns/git-clone-read-private.yaml |
  * Verify pipelinerun
      | S.NO | pipeline_run_name                   | status     |
      |------|-------------------------------------|------------|
      | 1    | git-clone-read-private-pipeline-run | successful |
  * Delete "pipelinerun" named "git-clone-read-private-pipeline-run"

## Setup S2I golang pipeline pre upgrade: PIPELINES-18-TC05
Tags: pre-upgrade, e2e, clustertasks, non-admin, s2i
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Create project "releasetest-upgrade-s2i"
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                          |
      |----|------------------------------------------------------|
      |1   |testdata/ecosystem/pipelines/s2i-go.yaml|
      |2   |testdata/pvc/pvc.yaml                                 |

## Validate olm skiprange pre upgrade: PIPELINES-18-TC06
Tags: pre-upgrade, olm
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
 * Get olm-skip-range "pre-upgrade" and save to field "pre-upgrade-olm-skip-range" in file "testdata/olm/skiprange.json"