PIPELINES-28
# ManualApprovalGate Pipelines operator specs

Pre condition:
  * Validate manual approval gate deployment

## Approve Manual Approval gate pipeline: PIPELINES-28-TC01
Tags: approvalgate, e2e, taskrun, sanity
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

Steps:
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
  * Create
     | S.NO | resource_dir                                             |
     |------|----------------------------------------------------------|
     | 1    | testdata/manualapprovalgate/manual-approval-pipeline.yaml|
  * Start the "manual-approval-pipeline" pipeline with workspace "name=source,claimName=shared-pvc"
  * Reject the manual-approval-pipeline
  * Validate the manual-approval-pipeline for "Rejected" state
  * Verify the latest pipelinerun for "failed" state

## Single Group / Single Approval: PIPELINES-28-TC03
Tags: approvalgate, approvalgate-users, e2e, taskrun, sanity
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Ensure approval group "group1" has members "user1"
  * Create manual approval gate pipelinerun with approvers "group:group1" required "1" Should "success"
  * User "user1" performs "approve" on the manual approval gate task
  * Validate manual approval gate task for "Approved" state
  * Verify the latest pipelinerun for "successful" state

## Quorum: Partial to Complete: PIPELINES-28-TC04
Tags: approvalgate, approvalgate-users, e2e, taskrun, sanity
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Ensure approval group "group1" has members "user1,user2"
  * Create manual approval gate pipelinerun with approvers "group:group1" required "2" Should "success"
  * User "user1" performs "approve" on the manual approval gate task
  * Validate manual approval gate task list state numberOfApprovalsRequired "2" pending "1" rejected "0" status "Pending"
  * User "user2" performs "approve" on the manual approval gate task
  * Validate manual approval gate task list state numberOfApprovalsRequired "2" pending "0" rejected "0" status "Approved"
  * Verify the latest pipelinerun for "successful" state

## Mixed Entities (User + Group): PIPELINES-28-TC05
Tags: approvalgate, approvalgate-users, e2e, taskrun, sanity
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Ensure approval group "group1" has members "user1,user2"
  * Create manual approval gate pipelinerun with approvers "user5,group:group1" required "2" Should "success"
  * User "user1" performs "approve" on the manual approval gate task
  * User "user5" performs "approve" on the manual approval gate task
  * Verify the latest pipelinerun for "successful" state

## Rejection Authority: PIPELINES-28-TC06
Tags: approvalgate, approvalgate-users, e2e, taskrun, sanity
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Ensure approval group "group1" has members "user1,user2"
  * Create manual approval gate pipelinerun with approvers "group:group1" required "2" Should "fail-fast"
  * User "user1" performs "approve" on the manual approval gate task
  * User "user2" performs "reject" on the manual approval gate task
  * Validate manual approval gate task for "Rejected" state
  * Verify the latest pipelinerun for "failed" state

## Change of Mind (Approve to Reject): PIPELINES-28-TC07
Tags: approvalgate, approvalgate-users, e2e, taskrun, sanity
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Ensure approval group "group1" has members "user1,user2"
  * Create manual approval gate pipelinerun with approvers "group:group1" required "2" Should "fail-fast"
  * User "user1" performs "approve" on the manual approval gate task
  * User "user1" performs "reject" on the manual approval gate task
  * Validate manual approval gate task for "Rejected" state
  * Verify the latest pipelinerun for "failed" state

## Non-Member Block: PIPELINES-28-TC08
Tags: approvalgate, approvalgate-users, e2e, taskrun, sanity
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Ensure approval group "group1" has members "user1"
  * Create manual approval gate pipelinerun with approvers "group:group1" required "1" Should "success"
  * User "user4" performs "approve-expect-fail" on the manual approval gate task
  * User "user1" performs "approve" on the manual approval gate task
  * Verify the latest pipelinerun for "successful" state

## Multi-Group Consensus: PIPELINES-28-TC09
Tags: approvalgate, approvalgate-users, e2e, taskrun, sanity
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Ensure approval group "group1" has members "user1"
  * Ensure approval group "group2" has members "user4"
  * Create manual approval gate pipelinerun with approvers "group:group1,group:group2" required "2" Should "success"
  * User "user1" performs "approve" on the manual approval gate task
  * User "user4" performs "approve" on the manual approval gate task
  * Verify the latest pipelinerun for "successful" state

## Overlapping Membership: PIPELINES-28-TC10
Tags: approvalgate, approvalgate-users, e2e, taskrun, sanity
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Ensure approval group "group1" has members "user1,user2"
  * Ensure approval group "group2" has members "user2"
  * Create manual approval gate pipelinerun with approvers "group:group1,group:group2" required "2" Should "success"
  * User "user2" performs "approve" on the manual approval gate task
  * User "user1" performs "approve" on the manual approval gate task
  * Verify the latest pipelinerun for "successful" state

## Timeout Expiry (Short Timeout): PIPELINES-28-TC11
Tags: approvalgate, approvalgate-users, e2e, taskrun, sanity
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Ensure approval group "group1" has members "user1,user2"
  * Create manual approval gate pipelinerun with approvers "group:group1" required "2" Should "timeout"
  * Verify the latest pipelinerun for "failed" state

## Multi-Group Race:Any One can approve: PIPELINES-28-TC12
Tags: approvalgate, approvalgate-users, e2e, taskrun, sanity
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Ensure approval group "group1" has members "user1"
  * Ensure approval group "group2" has members "user5"
  * Create manual approval gate pipelinerun with approvers "group:group1,group:group2" required "1" Should "success"
  * User "user5" performs "approve" on the manual approval gate task
  * Verify the latest pipelinerun for "successful" state

## Mixed Entity Change-of-Mind: PIPELINES-28-TC13
Tags: approvalgate, approvalgate-users, e2e, taskrun, sanity
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Ensure approval group "group1" has members "user1"
  * Create manual approval gate pipelinerun with approvers "user1,group:group1" required "2" Should "fail-fast"
  * User "user1" performs "approve" on the manual approval gate task
  * User "user1" performs "reject" on the manual approval gate task
  * Validate manual approval gate task for "Rejected" state
  * Verify the latest pipelinerun for "failed" state

## Re-approve Completed Task: PIPELINES-28-TC14
Tags: approvalgate, approvalgate-users, e2e, taskrun, sanity
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Ensure approval group "group1" has members "user1"
  * Create manual approval gate pipelinerun with approvers "group:group1" required "1" Should "success"
  * User "user1" performs "approve" on the manual approval gate task
  * Validate manual approval gate task for "Approved" state
  * User "user1" performs "approve-expect-fail" on the manual approval gate task
  * Verify the latest pipelinerun for "successful" state

## Impossible Quorum (Short Timeout): PIPELINES-28-TC15
Tags: approvalgate, approvalgate-users, e2e, taskrun, sanity
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Ensure approval group "group1" has members "user1,user2,user3"
  * Create manual approval gate pipelinerun with approvers "group:group1" required "4" Should "timeout"
  * User "user1" performs "approve" on the manual approval gate task
  * User "user2" performs "approve-allow-final-state" on the manual approval gate task
  * User "user3" performs "approve-allow-final-state" on the manual approval gate task
  * Verify the latest pipelinerun for "failed" state

## Invalid Group Name: PIPELINES-28-TC16
Tags: approvalgate, approvalgate-users, e2e, taskrun, sanity
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Create manual approval gate pipelinerun with approvers "group:invalid-group" required "1" Should "timeout"
  * User "user1" performs "approve-expect-fail" on the manual approval gate task
  * Verify the latest pipelinerun for "failed" state

## Re-approve Rejected Task: PIPELINES-28-TC17
Tags: approvalgate, approvalgate-users, e2e, taskrun, sanity
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Ensure approval group "group1" has members "user1"
  * Create manual approval gate pipelinerun with approvers "group:group1" required "1" Should "fail-fast"
  * User "user1" performs "reject" on the manual approval gate task
  * Validate manual approval gate task for "Rejected" state
  * User "user1" performs "approve-expect-fail" on the manual approval gate task
  * Verify the latest pipelinerun for "failed" state

## The Late Joiner: PIPELINES-28-TC18
Tags: approvalgate, approvalgate-users, e2e, taskrun, sanity
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Ensure approval group "group1" has members "-"
  * Create manual approval gate pipelinerun with approvers "group:group1" required "1" Should "success"
  * Ensure approval group "group1" has members "user1"
  * User "user1" performs "approve" on the manual approval gate task
  * Verify the latest pipelinerun for "successful" state

## The Evicted User: PIPELINES-28-TC19
Tags: approvalgate, approvalgate-users, e2e, taskrun, sanity
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Ensure approval group "group1" has members "user1"
  * Create manual approval gate pipelinerun with approvers "group:group1" required "1" Should "timeout"
  * Ensure approval group "group1" has members "-"
  * User "user1" performs "approve-expect-fail" on the manual approval gate task
  * Verify the latest pipelinerun for "failed" state

## The Switcheroo: PIPELINES-28-TC20
Tags: approvalgate, approvalgate-users, e2e, taskrun, sanity
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Ensure approval group "group1" has members "user1"
  * Create manual approval gate pipelinerun with approvers "group:group1" required "1" Should "success"
  * Ensure approval group "group1" has members "user2"
  * User "user1" performs "approve-expect-fail" on the manual approval gate task
  * User "user2" performs "approve" on the manual approval gate task
  * Verify the latest pipelinerun for "successful" state

## Dynamic Quorum: PIPELINES-28-TC21
Tags: approvalgate, approvalgate-users, e2e, taskrun, sanity
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Ensure approval group "group1" has members "user1"
  * Create manual approval gate pipelinerun with approvers "group:group1" required "2" Should "success"
  * User "user1" performs "approve" on the manual approval gate task
  * Ensure approval group "group1" has members "user1,user2"
  * User "user2" performs "approve" on the manual approval gate task
  * Verify the latest pipelinerun for "successful" state

## Approval Message Audit: PIPELINES-28-TC22
Tags: approvalgate, approvalgate-users, e2e, taskrun, sanity
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Ensure approval group "group1" has members "user1"
  * Create manual approval gate pipelinerun with approvers "group:group1" required "1" Should "success"
  * User "user1" performs "approve" on the manual approval gate task with message "PIPELINES-28-TC22 custom approve message from user1"
  * Verify manual approval gate task message contains "PIPELINES-28-TC22 custom approve message from user1"
  * Verify the latest pipelinerun for "successful" state
