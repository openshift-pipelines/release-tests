# Verify triggerbindings spec

Pre condition:
  * Validate Operator should be installed

## Verify CEL marshaljson function Test: PIPELINES-10-TC01
Tags: e2e,triggers,non-admin
Component: Triggers
Level: Integration
Type: Functional
Importance: Critical

This scenario helps you to verify whethere event message body parsed correctly, using CEL marshalJson function
Steps:
  * Create
    |S.NO|resource_dir                                          |
    |----|------------------------------------------------------|
    |1   |testdata/triggers/triggerbindings/cel-marshalJson.yaml|
  * Expose Event listener "cel-marshaljson"
  * Mock post event to "github" interceptor with event-type "push", payload "testdata/push.json", with TLS "false"
  * Assert eventlistener response
  * Verify taskrun
       |S.NO|task_run_name       |status |
       |----|--------------------|-------|
       |1   |cel-trig-marshaljson|Success|
  * Cleanup Triggers

## Verify event message body parsing with old annotation Test: PIPELINES-10-TC02
Tags: e2e,triggers,non-admin
Component: Triggers
Level: Integration
Type: Functional
Importance: Critical

This scenario helps you to verify whethere event message body parsed correctly, using old annoations in triggertemplate
Steps:
  * Create
    |S.NO|resource_dir                                                          |
    |----|----------------------------------------------------------------------|
    |1   |testdata/triggers/triggerbindings/parse-json-body-with-annotation.yaml|
  * Expose Event listener "parse-json-body-with-annotation"
  * Mock post event to "github" interceptor with event-type "push", payload "testdata/push.json", with TLS "false"
  * Assert eventlistener response
  * Verify taskrun
       |S.NO|task_run_name                       |status |
       |----|------------------------------------|-------|
       |1   |trig-parse-json-body-with-annotation|Success|
  * Cleanup Triggers

## Verify event message body marshalling error Test: PIPELINES-10-TC03
Tags: bug-to-fix,non-admin
Component: Triggers
Level: Integration
Type: Functional
Importance: Critical

This scenario helps you to verify whethere event message body parsed correctly.
Steps:
  * Create
    |S.NO|resource_dir                                          |
    |----|------------------------------------------------------|
    |1   |testdata/triggers/triggerbindings/parse-json-body.yaml|
  * Expose Event listener "parse-json-body"
  * Mock post event to "github" interceptor with event-type "push", payload "testdata/push.json", with TLS "false"
  * Assert eventlistener response
  * Verify taskrun
       |S.NO|task_run_name       |status |
       |----|--------------------|-------|
       |1   |trig-parse-json-body|Success|
  * Cleanup Triggers