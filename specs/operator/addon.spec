PIPELINES-15
# Verify Addon E2E spec

Pre condition:
  * Validate Operator should be installed

## Disable/Enable resolverTasks: PIPELINES-15-TC06
Tags: e2e, integration, resolverTasks, admin, addon, sanity
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Update addon config with resolverTasks as "true" and expect message ""
  * "s2i-java" tasks are "present" in namespace "openshift-pipelines"
  * Update addon config with resolverTasks as "false" and expect message ""
  * "s2i-java" tasks are "not present" in namespace "openshift-pipelines"

## Disable/Enable resolverTasks with additional Tasks: PIPELINES-15-TC07
Tags: e2e, integration, resolverTasks, admin, addon, sanity
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Update addon config with resolverTasks as "true" and expect message ""
  * "s2i-java" tasks are "present" in namespace "openshift-pipelines"
  * Create task hello in namespace "openshift-pipelines"
  * "hello" tasks are "present" in namespace "openshift-pipelines"
  * Update addon config with resolverTasks as "false" and expect message ""
  * "s2i-java" tasks are "not present" in namespace "openshift-pipelines"
  * "hello" tasks are "present" in namespace "openshift-pipelines"

