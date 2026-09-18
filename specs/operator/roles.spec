PIPELINES-34

# Verify Roles for OSP

Pre condition:
  * Validate Operator should be installed

## Verify Roles in openshift-pipelines ns: PIPELINES-11-TC02
Tags: e2e, admin, sanity
Component: Operator
Level: Integration
Type: Functional
Importance: High

This scenario helps you to verify all required roles by OSP are present.

Steps:
  * Verify the roles are present in "openshift-pipelines" namespace:
    | Role Name                                 |
    | manual-approval-gate-controller           |
    | manual-approval-gate-info                 |
    | manual-approval-gate-webhook              |
    | openshift-pipelines-read                  |
    | pipelines-as-code-controller-role         | 
    | pipelines-as-code-info                    |
    | pipelines-as-code-monitoring              |
    | pipelines-as-code-watcher-role            |
    | pipelines-as-code-webhook-role            |
    | tekton-chains-info                        |
    | tekton-chains-leader-election             |
    | tekton-default-openshift-pipelines-view   |
    | tekton-ecosystem-stepaction-list-role     |
    | tekton-ecosystem-task-list-role           |
    | tekton-hub-info                           |
    | tekton-operators-proxy-admin              |
    | tekton-pipelines-controller               |
    | tekton-pipelines-events-controller        |
    | tekton-pipelines-info                     |
    | tekton-pipelines-leader-election          |
    | tekton-pipelines-resolvers-namespace-rbac |
    | tekton-pipelines-webhook                  |
    | tekton-results-info                       |
    | tekton-triggers-admin-webhook             |
    | tekton-triggers-core-interceptors         |
    | tekton-triggers-info                      |
  * Verify the total number of roles in "openshift-pipelines" namespace matches the table
