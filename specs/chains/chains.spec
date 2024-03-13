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
    * Patch tekton config to sign and verify "taskrun" with Tekton Chains
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
    * Patch tekton config to sign and verify "image" with Tekton Chains
    * Import image registry variable 
    * Create quay secret for Tekton Chains
    * Apply
     | S.NO | resource_dir                          |
     |------|---------------------------------------|
     | 1    | testdata/pvc/chains-pvc.yaml          |
     | 2    | testdata/chains/kaniko.yaml           |
    * Start task
    * Verify image signature
    * Check Attestation
    * Verify Attestation