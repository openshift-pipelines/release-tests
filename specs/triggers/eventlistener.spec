# Verify eventlisteners spec

Pre condition:
  * Operator should be installed

## Create Eventlistener
Tags: triggers

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
  * Verify pipelinerun is "successfull"

## Create Eventlistener with github interceptor
Tags: triggers

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
  * Verify pipelinerun is "successfull"

## Create EventListener with custom interceptor
Tags: triggers

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
  * Verify pipelinerun is "successfull"

## Create EventListener with CEL interceptor with filter
Tags: triggers

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
  * Verify pipelinerun is "successfull"

## Create EventListener with CEL interceptor without filter
Tags: triggers

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
  * Verify pipelinerun is "successfull"


## Create EventListener with multiple interceptors
Tags: triggers

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
  * Verify pipelinerun is "successfull"

## Create Eventlistener embedded TriggersBindings specs
Tags: e2e, triggers

This scenario helps you to create eventlistner with embedded triggerbinding spec, listens to events forwards request to validator service -> parsed response to other validators and so on, on each event it creates/triggers
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
  * Mock push event "testdata/push.json"
  * Assert eventlistener "listener-embed-binding" response
  * Verify pipelinerun
    |S.NO|pipeline_run_name  |status     |check_lable_propagation|
    |----|-------------------|-----------|-----------------------|
    |1   |simple-pipeline-run|successfull|no                     |

## Create embedded TriggersTemplate
Tags: e2e, triggers

This scenario helps you to create embedded triggertemplate spec, listens to events forwards request to validator service -> parsed response to other validators and so on, on each event it creates/triggers
openshift-pipeline Resources defined under triggers-template, which helps you to deploy example app

Steps:
  * Create
    |S.NO|resource_dir                                                       |
    |----|-------------------------------------------------------------------|
    |1   |testdata/triggers/triggerbindings/triggerbinding.yaml              |
    |2   |testdata/triggers/triggertemplate/embed-triggertemplate.yaml       |
    |3   |testdata/triggers/eventlisteners/eventlistener-embeded-binding.yaml|
  * Expose Event listener "listener-embed-binding"
  * Mock push event "testdata/push.json"
  * Assert eventlistener "listener-embed-binding" response
  * Verify pipelinerun
    |S.NO|pipeline_run_name                        |status     |check_lable_propagation|
    |----|-----------------------------------------|-----------|-----------------------|
    |1   |pipelinerun-with-taskspec-to-echo-message|successfull|no                     |

## Create Eventlistener with gitlab interceptor
Tags: e2e, triggers

This scenario helps you to create eventlistner with gitlab interceptor, listens to gitlab events, on each event it creates/triggers
openshift-pipeline Resources defined under triggers-template, to deploy example app

Steps:
  * Create
    |S.NO|resource_dir                                      |
    |----|--------------------------------------------------|
    |1   |testdata/triggers/gitlab/gitlab-push-listener.yaml|
    |2   |testdata/triggers/gitlab/secret.yaml              |
  * Link secret "gitlab-secret" to service account pipeline  
  * Expose Event listener "gitlab-listener"
  * Mock push event "testdata/triggers/gitlab/gitlab-push-event.json" to gitlab interceptor
  * Assert eventlistener "gitlab-listener" response
  * Verify taskrun
    |S.NO|task_run_name |status     |
    |----|--------------|-----------|
    |1   |gitlab-run    |successfull|