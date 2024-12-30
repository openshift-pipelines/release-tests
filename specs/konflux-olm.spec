PIPELINES-10
# Olm Openshift Pipelines operator specs for Konflux

## Configure GitHub token for git resolver in TektonConfig: PIPELINES-10-TC01
Tags: konflux
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

Configure `openshift-pipelines` operator
Steps:
  * Wait for TektonConfig CR availability
  * Validate Operator should be installed
  * Verify namespace "openshift-pipelines" exists
  * Configure GitHub token for git resolver in TektonConfig

