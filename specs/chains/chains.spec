# Results pvc tests

Precondition:
* Validate Operator should be installed
* Create signing-secrets for tekton chains

## Test Tekton chains verify taskrun signature: PIPELINES-26-TC01
Tags: results, e2e, taskrun
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

## Test Tekton chains verify image signature: PIPELINES-26-TC01
Tags: results, e2e, image
Component: Chains
Level: Integration
Type: Functional
Importance: Critical
Steps:
    * Patch tekton config to sign and verify "image" with Tekton Chains
    * Create quay secret for tekton chains
    * Apply
     | S.NO | resource_dir                          |
     |------|---------------------------------------|
     | 1    | testdata/pvc/chains-pvc.yaml          |
     | 2    | testdata/chains/kaniko.yaml           |
    * Start task
    * Verify image signature
    * Check Attestation
    * Verify Attestation