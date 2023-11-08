PIPELINES-05
# Verify eventlisteners spec

Pre condition:
  * Validate Operator should be installed

## Create Eventlistener: PIPELINES-05-TC01
Tags: triggers, to-do
Component: Triggers
Level: Integration
Type: Functional
Importance: High

This scenario helps you to create eventlistner, listens to github events by default, on each github event it creates/triggers
openshift-pipeline Resources defined under triggers-template, to deploy example app

Steps:
  * Create "./testdata/triggers/eventlisteners/triggertemplate/template.yaml"
  * Create "./testdata/triggers/eventlisteners/triggerbinding/binding.yaml"
  * Create "./testdata/triggers/eventlisteners/triggerbinding/binding-message.yaml"
  * Create "./testdata/triggers/eventlisteners/role-resources/rbac.yaml"
  * Add Event listener
  * Expose event listener service
  * Mock push event
  * Verify creation of openshift-pipeline-resources
  * Verify resources are created with labels & event-id
  * Verify pipelinerun is "successful"

## Create Eventlistener with github interceptor: PIPELINES-05-TC02
Tags: triggers, to-do
Component: Triggers
Level: Integration
Type: Functional
Importance: High

This scenario helps you to create eventlistner with github interceptor, listens to github events, on each github event it creates/triggers
openshift-pipeline Resources defined under triggers-template, to deploy example app

Steps:
  * Create "./testdata/triggers/eventlisteners/triggertemplate/template.yaml"
  * Create "./testdata/triggers/eventlisteners/triggerbinding/binding.yaml"
  * Create "./testdata/triggers/eventlisteners/triggerbinding/binding-message.yaml"
  * Create "./testdata/triggers/eventlisteners/role-resources/rbac.yaml"
  * Add Event listener with github interceptor
  * Expose event listener service
  * Mock push event
  * Verify creation of openshift-pipeline-resources
  * Verify resources are created with labels & event-id
  * Verify pipelinerun is "successful"

## Create EventListener with custom interceptor: PIPELINES-05-TC03
Tags: triggers, to-do
Component: Triggers
Level: Integration
Type: Functional
Importance: High

This scenario helps you to create eventlistner with custom interceptor, listens to custom events, on each event to custom service it creates/triggers
openshift-pipeline Resources defined under triggers-template, to deploy example app

Steps:
  * Create "./testdata/triggers/eventlisteners/triggertemplate/template.yaml"
  * Create "./testdata/triggers/eventlisteners/triggerbinding/binding.yaml"
  * Create "./testdata/triggers/eventlisteners/triggerbinding/binding-message.yaml"
  * Create "./testdata/triggers/eventlisteners/role-resources/rbac.yaml"
  * Create gh-validator service "./testdata/triggers/eventlisteners/custom-interceptor/gh-validate-service.yaml"
  * Add Event listener with custom interceptor
  * Expose event listener service
  * Mock push event
  * Verify creation of openshift-pipeline-resources
  * Verify resources are created with labels & event-id
  * Verify pipelinerun is "successful"

## Create EventListener with CEL interceptor with filter: PIPELINES-05-TC04
Tags: triggers, to-do
Component: Triggers
Level: Integration
Type: Functional
Importance: High

This scenario helps you to create eventlistner with CEL interceptor with filter, listens to filtered CEL events, on each event it creates/triggers
openshift-pipeline Resources defined under triggers-template, to deploy example app

Steps:
  * Create "./testdata/triggers/eventlisteners/triggertemplate/template.yaml"
  * Create "./testdata/triggers/eventlisteners/triggerbinding/binding.yaml"
  * Create "./testdata/triggers/eventlisteners/role-resources/rbac.yaml"
  * Add Event listener with CEL interceptor with filter
  * Expose event listener service
  * Mock CEL push/pr event
  * Verify creation of openshift-pipeline-resources
  * Verify resources are created with labels & event-id
  * Verify pipelinerun is "successful"

## Create EventListener with CEL interceptor without filter: PIPELINES-05-TC05
Tags: triggers, to-do
Component: Triggers
Level: Integration
Type: Functional
Importance: High

This scenario helps you to create eventlistner with CEL interceptor, listens to all CEL events, on each event it creates/triggers
openshift-pipeline Resources defined under triggers-template, to deploy example app

Steps:
  * Create "./testdata/triggers/eventlisteners/triggertemplate/template.yaml"
  * Create "./testdata/triggers/eventlisteners/triggerbinding/binding.yaml"
  * Create "./testdata/triggers/eventlisteners/role-resources/rbac.yaml"
  * Add Event listener with CEL interceptor without filter
  * Expose event listener service
  * Mock CEL push/pr event
  * Verify creation of openshift-pipeline-resources
  * Verify resources are created with labels & event-id
  * Verify pipelinerun is "successful"


## Create EventListener with multiple interceptors: PIPELINES-05-TC06
Tags: triggers, to-do
Component: Triggers
Level: Integration
Type: Functional
Importance: High

This scenario helps you to create eventlistner with multiple interceptors, listens to events forwards request to validator service -> parsed response to other validators and so on, on each event it creates/triggers
openshift-pipeline Resources defined under triggers-template, which helps you to deploy example app

Steps:
  * Create "./testdata/triggers/eventlisteners/triggertemplate/template.yaml"
  * Create "./testdata/triggers/eventlisteners/triggerbinding/binding.yaml"
  * Create "./testdata/triggers/eventlisteners/role-resources/rbac.yaml"
  * Add Event listener with multiple interceptors
  * Expose event listener service
  * Mock push event
  * Verify creation of openshift-pipeline-resources
  * Verify resources are created with labels & event-id
  * Verify pipelinerun is "successful"

## Create Eventlistener with TLS enabled: PIPELINES-05-TC07
Tags: tls, triggers, admin, e2e
Component: Triggers
Level: Integration
Type: Functional
Importance: Critical

This scenario tests the creation of eventLister with TLS enabled, listens to events forwards request to validator service -> parsed response to other validators and so on, on each event it creates/triggers
openshift-pipeline Resources defined under triggers-template, which helps you to deploy example app

Steps:
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
    |S.NO|pipeline_run_name  |status    |check_label_propagation|
    |----|-------------------|----------|-----------------------|
    |1   |simple-pipeline-run|successful|no                     |
  * Cleanup Triggers

## Create Eventlistener embedded TriggersBindings specs: PIPELINES-05-TC08
Tags: e2e, triggers, non-admin, sanity
Component: Triggers
Level: Integration
Type: Functional
Importance: Critical

This scenario tests the creation of eventLister with embedded triggerbinding spec, listens to events forwards request to validator service -> parsed response to other validators and so on, on each event it creates/triggers
openshift-pipeline Resources defined under triggers-template, which helps you to deploy example app

Steps:
  * Create
    |S.NO|resource_dir                                                        |
    |----|--------------------------------------------------------------------|
    |1   |testdata/triggers/sample-pipeline.yaml                              |
    |2   |testdata/triggers/triggerbindings/triggerbinding.yaml               |
    |3   |testdata/triggers/triggertemplate/triggertemplate.yaml              |
    |4   |testdata/triggers/eventlisteners/eventlistener-embeded-binding.yaml |
  * Expose Event listener "listener-embed-binding"
  * Mock post event to "github" interceptor with event-type "push", payload "testdata/push.json", with TLS "false"
  * Assert eventlistener response
  * Verify pipelinerun
    |S.NO|pipeline_run_name  |status    |check_label_propagation|
    |----|-------------------|----------|-----------------------|
    |1   |simple-pipeline-run|successful|no                     |
  * Cleanup Triggers

## Create embedded TriggersTemplate: PIPELINES-05-TC09
Tags: e2e, triggers, non-admin, sanity
Component: Triggers
Level: Integration
Type: Functional
Importance: Critical

This scenario tests the creation of embedded triggertemplate spec, listens to events forwards request to validator service -> parsed response to other validators and so on, on each event it creates/triggers
openshift-pipeline Resources defined under triggers-template, which helps you to deploy example app

Steps:
  * Create
    |S.NO|resource_dir                                                       |
    |----|-------------------------------------------------------------------|
    |1   |testdata/triggers/triggerbindings/triggerbinding.yaml              |
    |2   |testdata/triggers/triggertemplate/embed-triggertemplate.yaml       |
    |3   |testdata/triggers/eventlisteners/eventlistener-embeded-binding.yaml|
  * Expose Event listener "listener-embed-binding"
  * Mock post event to "github" interceptor with event-type "push", payload "testdata/push.json", with TLS "false"
  * Assert eventlistener response
  * Verify pipelinerun
    |S.NO|pipeline_run_name                        |status    |check_label_propagation|
    |----|-----------------------------------------|----------|-----------------------|
    |1   |pipelinerun-with-taskspec-to-echo-message|successful|no                     |
  * Cleanup Triggers

## Create Eventlistener with gitlab interceptor: PIPELINES-05-TC10
Tags: e2e, triggers, non-admin
Component: Triggers
Level: Integration
Type: Functional
Importance: Critical

This scenario tests the creation of eventLister with gitlab interceptor, listens to gitlab events, on each event it creates/triggers
openshift-pipeline Resources defined under triggers-template, to deploy example app

Steps:
  * Create
    |S.NO|resource_dir                                      |
    |----|--------------------------------------------------|
    |1   |testdata/triggers/gitlab/gitlab-push-listener.yaml|
  * Verify ServiceAccount "pipeline" exist
  * Create & Link secret "gitlab-secret" to service account "pipeline"
  * Expose Event listener "gitlab-listener"
  * Mock post event to "gitlab" interceptor with event-type "Push Hook", payload "testdata/triggers/gitlab/gitlab-push-event.json", with TLS "false"
  * Assert eventlistener response
  * Verify pipelinerun
    |S.NO|pipeline_run_name |status    |check_label_propagation|
    |----|------------------|----------|-----------------------|
    |1   |gitlab-run        |successful|no                     |
  * Cleanup Triggers

## Create Eventlistener with bitbucket interceptor: PIPELINES-05-TC11
Tags: e2e, triggers, non-admin
Component: Triggers
Level: Integration
Type: Functional
Importance: Critical

This scenario tests the creation of eventLister with bitbucket interceptor, listens to bitbucket events, on each event it creates/triggers
openshift-pipeline Resources defined under triggers-template

Steps:
  * Create
    |S.NO|resource_dir                                                        |
    |----|--------------------------------------------------------------------|
    |1   |testdata/triggers/bitbucket/bitbucket-eventlistener-interceptor.yaml|
  * Verify ServiceAccount "pipeline" exist
  * Create & Link secret "bitbucket-secret" to service account "pipeline"
  * Expose Event listener "bitbucket-listener"
  * Mock post event to "bitbucket" interceptor with event-type "refs_changed", payload "testdata/triggers/bitbucket/refs-change-event.json", with TLS "false"
  * Assert eventlistener response
  * Verify taskrun
    |S.NO|task_run_name    |status |
    |----|-----------------|-------|
    |1   |bitbucket-run    |Failure|
  * Cleanup Triggers

## Verify Github push event with Embbeded TriggerTemplate using Github-CTB: PIPELINES-05-TC12
Tags: e2e, triggers, non-admin, sanity
Component: Triggers
Level: Integration
Type: Functional
Importance: Critical

This scenario tests Github `push` event via CTB, on each event it creates/triggers
openshift-pipeline Resources defined under triggers-template

Steps:
  * Create
    |S.NO|resource_dir                                                      |
    |----|------------------------------------------------------------------|
    |1   |testdata/triggers/github-ctb/Embeddedtriggertemplate-git-push.yaml|
    |2   |testdata/triggers/github-ctb/eventlistener-ctb-git-push.yaml      |
  * Verify ServiceAccount "pipeline" exist
  * Create & Link secret "github-secret" to service account "pipeline"
  * Expose Event listener "listener-ctb-github-push"
  * Mock post event to "github" interceptor with event-type "push", payload "testdata/triggers/github-ctb/push.json", with TLS "false"
  * Assert eventlistener response
  * Verify pipelinerun
    |S.NO|pipeline_run_name       |status    |check_label_propagation|
    |----|------------------------|----------|-----------------------|
    |1   |pipelinerun-git-push-ctb|successful|no                     |
  * Cleanup Triggers

## Verify Github pull_request event with Embbeded TriggerTemplate using Github-CTB: PIPELINES-05-TC13
Tags: e2e, triggers, non-admin, sanity
Component: Triggers
Level: Integration
Type: Functional
Importance: Critical

This scenario tests Github `pull_request` event via CTB, on each event it creates/triggers
openshift-pipeline Resources defined under triggers-template

Steps:
  * Create
    |S.NO|resource_dir                                                    |
    |----|----------------------------------------------------------------|
    |1   |testdata/triggers/github-ctb/Embeddedtriggertemplate-git-pr.yaml|
    |2   |testdata/triggers/github-ctb/eventlistener-ctb-git-pr.yaml      |
  * Verify ServiceAccount "pipeline" exist
  * Create & Link secret "github-secret" to service account "pipeline"
  * Expose Event listener "listener-clustertriggerbinding-github-pr"
  * Mock post event to "github" interceptor with event-type "pull_request", payload "testdata/triggers/github-ctb/pr.json", with TLS "false"
  * Assert eventlistener response
  * Verify pipelinerun
    |S.NO|pipeline_run_name       |status    |check_label_propagation|
    |----|------------------------|----------|-----------------------|
    |1   |pipelinerun-git-pr-ctb  |successful|no                     |
  * Cleanup Triggers

## Verify Github pr_review event with Embbeded TriggerTemplate using Github-CTB: PIPELINES-05-TC14
Tags: e2e, triggers, non-admin
Component: Triggers
Level: Integration
Type: Functional
Importance: Critical

This scenario tests Github `issue_comment` event via CTB, on each event it creates/triggers
openshift-pipeline Resources defined under triggers-template

Steps:
  * Create
    |S.NO|resource_dir                                                           |
    |----|-----------------------------------------------------------------------|
    |1   |testdata/triggers/github-ctb/Embeddedtriggertemplate-git-pr-review.yaml|
    |2   |testdata/triggers/github-ctb/eventlistener-ctb-git-pr-review.yaml      |
  * Verify ServiceAccount "pipeline" exist
  * Create & Link secret "github-secret" to service account "pipeline"
  * Expose Event listener "listener-ctb-github-pr-review"
  * Mock post event to "github" interceptor with event-type "issue_comment", payload "testdata/triggers/github-ctb/issue-comment.json", with TLS "false"
  * Assert eventlistener response
  * Verify pipelinerun
    |S.NO|pipeline_run_name                |status    |check_label_propagation|
    |----|---------------------------------|----------|-----------------------|
    |1   |pipelinerun-git-pr-review-ctb    |successful|no                     |
  * Cleanup Triggers

## Create TriggersCRD resource with CEL interceptors (overlays): PIPELINES-05-TC15
Tags: e2e, triggers, non-admin, sanity
Component: Triggers
Level: Integration
Type: Functional
Importance: Critical

This scenario tests the creation of Trigger resource which is combination of TriggerTemplate, TriggerBindings and interceptors. The Trigger is processed by EventListener, and listens to events, on each event it creates/triggers
openshift-pipeline Resources defined under triggers-template

Steps:
  * Create
    |S.NO|resource_dir                                               |
    |----|-----------------------------------------------------------|
    |1   |testdata/triggers/triggersCRD/eventlistener-triggerref.yaml|
    |2   |testdata/triggers/triggersCRD/trigger.yaml                 |
    |3   |testdata/triggers/triggersCRD/triggerbindings.yaml         |
    |4   |testdata/triggers/triggersCRD/triggertemplate.yaml         |
    |5   |testdata/triggers/triggersCRD/pipeline.yaml                |
  * Verify ServiceAccount "pipeline" exist
  * Create & Link secret "github-secret" to service account "pipeline"
  * Expose Event listener "listener-triggerref"
  * Mock post event to "github" interceptor with event-type "pull_request", payload "testdata/triggers/triggersCRD/pull-request.json", with TLS "false"
  * Assert eventlistener response
  * Verify pipelinerun
    |S.NO|pipeline_run_name       |status    |check_label_propagation|
    |----|------------------------|----------|-----------------------|
    |1   |parallel-pipelinerun    |successful|no                     |
  * Cleanup Triggers

## Create multiple Eventlistener with TLS enabled: PIPELINES-05-TC16
Tags: e2e, tls, triggers, admin, sanity
Component: Triggers
Level: Integration
Type: Functional
Importance: Critical
CustomerScenario: yes

This scenario tests the creation of multiple eventLister with TLS enabled, listens to events forwards request to validator service -> parsed response to other validators and so on, on each event it creates/triggers
openshift-pipeline Resources defined under triggers-template, which helps you to deploy example app

Steps:
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
    |S.NO|pipeline_run_name  |status    |check_label_propagation|
    |----|-------------------|----------|-----------------------|
    |1   |simple-pipeline-run|successful|no                     |
  * Create
    |S.NO|resource_dir                                                          |
    |----|----------------------------------------------------------------------|
    |1   |testdata/triggers/triggertemplate/triggertemplate-2.yaml              |
    |2   |testdata/triggers/eventlisteners/eventlistener-embeded-binding-2.yaml |
  * Expose Event listener for TLS "listener-embed-binding-2"
  * Mock post event to "github" interceptor with event-type "push", payload "testdata/push.json", with TLS "true"
  * Assert eventlistener response
  * Verify pipelinerun
    |S.NO|pipeline_run_name    |status    |check_label_propagation|
    |----|---------------------|----------|-----------------------|
    |1   |simple-pipeline-run-2|successful|no                     |
  * Cleanup Triggers

## Create Eventlistener with github interceptor And verify Kuberenetes Events: PIPELINES-05-TC17
Tags: e2e, events, triggers, admin, sanity
Component: Triggers
Level: Integration
Type: Functional
Importance: Medium

This scenario helps you to create eventlistner with github interceptor, listens to github events, on each github event it creates/triggers
openshift-pipeline Resources defined under triggers-template, to deploy example app and emit events for done/started/success/failure.

Steps:
  * Create
    |S.NO|resource_dir                                                        |
    |----|--------------------------------------------------------------------|
    |1   |testdata/triggers/sample-pipeline.yaml                              |
    |2   |testdata/triggers/triggerbindings/triggerbinding.yaml               |
    |3   |testdata/triggers/triggertemplate/triggertemplate.yaml              |
    |4   |testdata/triggers/eventlisteners/eventlistener-embeded-binding.yaml |
  * Expose Event listener "listener-embed-binding"
  * Mock post event to "github" interceptor with event-type "push", payload "testdata/push.json", with TLS "false"
  * Assert eventlistener response
  * Verify kubernetes events for eventlistener
  * Verify pipelinerun
    |S.NO|pipeline_run_name  |status    |check_label_propagation|
    |----|-------------------|----------|-----------------------|
    |1   |simple-pipeline-run|successful|no                     |
  * Cleanup Triggers
