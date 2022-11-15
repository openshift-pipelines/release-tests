PIPELINES-18
# Verify Chains

Contains scenarios that verify image signing using Chains

Pre condition:
  * Validate Operator should be installed

## Enable Chains: PIPELINES-18-TC01
Tags: e2e, chains, admin
Component: Chains
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create Chains CR with format "in-toto", storage "oci" and transparency enabled "true"
