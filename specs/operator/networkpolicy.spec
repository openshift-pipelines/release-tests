PIPELINES-38
# Verify NetworkPolicy for all components

Pre condition:
  * Validate Operator should be installed

## Verify NetworkPolicies are created by default: PIPELINES-38-TC01
Tags: e2e, networkpolicy, admin, sanity
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify NetworkPolicies are "present" in namespace "openshift-pipelines" for component "pipelines"
  * Verify NetworkPolicies are "present" in namespace "openshift-pipelines" for component "triggers"
  * Verify NetworkPolicies are "present" in namespace "openshift-pipelines" for component "chains"
  * Verify NetworkPolicies are "present" in namespace "openshift-pipelines" for component "results"
  * Verify NetworkPolicies are "present" in namespace "openshift-pipelines" for component "pruner"
  * Verify NetworkPolicies are "present" in namespace "openshift-pipelines" for component "manual-approval-gate"
  * Verify NetworkPolicies are "present" in namespace "openshift-pipelines" for component "pac"

## Disable NetworkPolicy via TektonConfig: PIPELINES-38-TC02
Tags: e2e, networkpolicy, admin
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Disable NetworkPolicy on TektonConfig
  * Wait for TektonConfig to be ready
  * Verify NetworkPolicies are "not present" in namespace "openshift-pipelines" for component "pipelines"
  * Verify NetworkPolicies are "not present" in namespace "openshift-pipelines" for component "triggers"
  * Verify NetworkPolicies are "not present" in namespace "openshift-pipelines" for component "chains"
  * Verify NetworkPolicies are "not present" in namespace "openshift-pipelines" for component "results"
  * Verify NetworkPolicies are "not present" in namespace "openshift-pipelines" for component "pruner"
  * Verify NetworkPolicies are "not present" in namespace "openshift-pipelines" for component "manual-approval-gate"
  * Verify NetworkPolicies are "not present" in namespace "openshift-pipelines" for component "pac"

## Re-enable NetworkPolicy via TektonConfig: PIPELINES-38-TC03
Tags: e2e, networkpolicy, admin
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Enable NetworkPolicy on TektonConfig
  * Wait for TektonConfig to be ready
  * Verify NetworkPolicies are "present" in namespace "openshift-pipelines" for component "pipelines"
  * Verify NetworkPolicies are "present" in namespace "openshift-pipelines" for component "triggers"
  * Verify NetworkPolicies are "present" in namespace "openshift-pipelines" for component "chains"
  * Verify NetworkPolicies are "present" in namespace "openshift-pipelines" for component "results"
  * Verify NetworkPolicies are "present" in namespace "openshift-pipelines" for component "pruner"
  * Verify NetworkPolicies are "present" in namespace "openshift-pipelines" for component "manual-approval-gate"
  * Verify NetworkPolicies are "present" in namespace "openshift-pipelines" for component "pac"

## Add custom NetworkPolicy via TektonConfig: PIPELINES-38-TC04
Tags: e2e, networkpolicy, admin
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Add custom NetworkPolicy "custom-test-deny-all" on TektonConfig
  * Wait for TektonConfig to be ready
  * Verify NetworkPolicy "custom-test-deny-all" is "present" in namespace "openshift-pipelines"
  * Verify NetworkPolicies are "present" in namespace "openshift-pipelines" for component "pipelines"
  * Remove custom NetworkPolicy from TektonConfig
  * Disable NetworkPolicy on TektonConfig
  * Wait for TektonConfig to be ready
  * Verify NetworkPolicy "custom-test-deny-all" is "not present" in namespace "openshift-pipelines"
  * Enable NetworkPolicy on TektonConfig
  * Wait for TektonConfig to be ready
  * Verify NetworkPolicy "custom-test-deny-all" is "not present" in namespace "openshift-pipelines"
  * Verify NetworkPolicies are "present" in namespace "openshift-pipelines" for component "pipelines"

## Add custom NetworkPolicy on standalone component: PIPELINES-38-TC05
Tags: e2e, networkpolicy, admin
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Add custom NetworkPolicy "mag-custom-test" on ManualApprovalGate
  * Wait for TektonConfig to be ready
  * Verify NetworkPolicy "mag-custom-test" is "present" in namespace "openshift-pipelines"
  * Remove custom NetworkPolicy from ManualApprovalGate
  * Disable NetworkPolicy on TektonConfig
  * Wait for TektonConfig to be ready
  * Verify NetworkPolicy "mag-custom-test" is "not present" in namespace "openshift-pipelines"
  * Enable NetworkPolicy on TektonConfig
  * Wait for TektonConfig to be ready
  * Verify NetworkPolicy "mag-custom-test" is "not present" in namespace "openshift-pipelines"
  * Verify NetworkPolicies are "present" in namespace "openshift-pipelines" for component "manual-approval-gate"

Teardown:
  * Remove custom NetworkPolicy from TektonConfig
  * Remove custom NetworkPolicy from ManualApprovalGate
  * Enable NetworkPolicy on TektonConfig
  * Wait for TektonConfig to be ready
