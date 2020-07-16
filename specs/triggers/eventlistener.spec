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

## Create Eventlistener with gitlab interceptor
Tags: triggers

This scenario helps you to create eventlistner with gitlab interceptor, listens to gitlab events, on each event it creates/triggers
openshift-pipeline Resources defined under triggers-template, to deploy example app

Steps:
  * Create "./testdata/triggers/eventlisteners/triggertemplate/template.yaml"
  * Create "./testdata/triggers/eventlisteners/triggerbinding/binding.yaml"
  * Create "./testdata/triggers/eventlisteners/triggerbinding/binding-message.yaml"
  * Create "./testdata/triggers/eventlisteners/role-resources/rbac.yaml"
  * Add Event listener with gitlab interceptor
  * Expose event listener service
  * Mock gitlab push event "./testdata/triggers/gitlab/gitlab-push-event.json"
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
