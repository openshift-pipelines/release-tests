PIPELINES-15
# Verify Addon E2E spec

Pre condition:
  * Validate Operator should be installed

## Disable/Enable resolverTasks: PIPELINES-15-TC06
Tags: e2e, integration, resolvertasks, admin, addon, sanity
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Update addon config with resolverTasks as "false" and expect message ""
  * Tasks "s2i-java" are "not present" in namespace "openshift-pipelines"
  * Update addon config with resolverTasks as "true" and expect message ""
  * Tasks "s2i-java" are "present" in namespace "openshift-pipelines"

## Disable/Enable resolverTasks with additional Tasks: PIPELINES-15-TC07
Tags: e2e, integration, resolvertasks, admin, addon
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Update addon config with resolverTasks as "true" and expect message ""
  * Tasks "s2i-java" are "present" in namespace "openshift-pipelines"
  * Apply in namespace "openshift-pipelines"
      |S.NO|resource_dir                                         |
      |----|-----------------------------------------------------|
      |1   |testdata/ecosystem/tasks/hello.yaml                  |
  * Tasks "hello" are "present" in namespace "openshift-pipelines"
  * Update addon config with resolverTasks as "false" and expect message ""
  * Tasks "s2i-java" are "not present" in namespace "openshift-pipelines"
  * Tasks "hello" are "present" in namespace "openshift-pipelines"
  * Update addon config with resolverTasks as "true" and expect message ""
  * Tasks "s2i-java" are "present" in namespace "openshift-pipelines"
  * Tasks "hello" are "present" in namespace "openshift-pipelines"

## Disable/Enable pipeline templates: PIPELINES-15-TC08
Tags: e2e, integration, resolvertasks, admin, addon, sanity
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Update addon config with resolverTasks as "true" and pipelineTemplates as "true" and expect message ""
  * Assert pipelines are "present" in "openshift" namespace
  * Update addon config with resolverTasks as "true" and pipelineTemplates as "false" and expect message ""
  * Assert pipelines are "not present" in "openshift" namespace
  * Update addon config with resolverTasks as "true" and pipelineTemplates as "true" and expect message ""
  * Assert pipelines are "present" in "openshift" namespace

## Enable pipeline templates when clustertask is disabled: PIPELINES-15-TC05
Tags: e2e, integration, negative, admin, addon
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical
Pos/Neg: Negative

Steps:
  * Update addon config with resolverTasks as "false" and pipelineTemplates as "true" and expect message "validation failed: pipelineTemplates cannot be true if resolverTask is false"

## Verify versioned ecosystem tasks: PIPELINES-15-TC09
Tags: e2e, integration, addon
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify versioned ecosystem tasks

## Verify versioned stepaction tasks: PIPELINES-15-TC010
Tags: e2e, integration, addon
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify versioned ecosystem step actions
