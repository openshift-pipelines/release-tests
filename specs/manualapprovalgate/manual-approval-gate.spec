PIPELINES-28
# ManualApprovalGate Pipelines operator specs

## Approve Manual Approval gate pipeline: PIPELINES-28-TC01
Tags: approvalgate, e2e, taskrun, sanity
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Validate manual approval gate deployment
  * Create
     | S.NO | resource_dir                                             |
     |------|----------------------------------------------------------|
     | 1    | testdata/manualapprovalgate/manual-approval-pipeline.yaml|
  * Start the "manual-approval-pipeline" pipeline with workspace "name=source,claimName=shared-pvc"
  * Approve the manual-approval-pipeline
  * Validate the manual-approval-pipeline for "Approved" state
  * Verify the latest pipelinerun for "successful" state

## Reject Manual Approval gate pipeline: PIPELINES-28-TC02
Tags: approvalgate, e2e, taskrun, sanity
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Validate manual approval gate deployment
  * Create
     | S.NO | resource_dir                                             |
     |------|----------------------------------------------------------|
     | 1    | testdata/manualapprovalgate/manual-approval-pipeline.yaml|
  * Start the "manual-approval-pipeline" pipeline with workspace "name=source,claimName=shared-pvc"
  * Reject the manual-approval-pipeline
  * Validate the manual-approval-pipeline for "Rejected" state
  * Verify the latest pipelinerun for "failed" state
