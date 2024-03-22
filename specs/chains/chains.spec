PIPELINES-27
# Tekton Chains tests

Precondition:
* Validate Operator should be installed

## Using Tekton Chains to create and verify task run signatures: PIPELINES-27-TC01
Tags: chains, e2e, taskrun, sanity
Component: Chains
Level: Integration
Type: Functional
Importance: Critical

Steps: 
    * Update the TektonConfig with taskrun format as "in-toto" taskrun storage as "tekton" oci storage as "" transparency mode as "false"
    * Apply
     | S.NO | resource_dir                          |
     |------|---------------------------------------|
     | 1    | testdata/chains/task-output-image.yaml|
    * Verify "taskrun" signature

## Using Tekton Chains to sign and verify image and provenance : PIPELINES-27-TC02
Tags: chains, e2e, image
Component: Chains
Level: Integration
Type: Functional
Importance: Critical
Steps:
    * Update the TektonConfig with taskrun format as "in-toto" taskrun storage as "oci" oci storage as "oci" transparency mode as "true"
    * Verify that image registry variable is exported
    * Create secret with image registry credentials for SA
    * Apply
     | S.NO | resource_dir                          |
     |------|---------------------------------------|
     | 1    | testdata/pvc/chains-pvc.yaml          |
     | 2    | testdata/chains/kaniko.yaml           |
    * Start the kaniko-chains task
    * Verify image signature
    * Check attestation exists
    * Verify attestation