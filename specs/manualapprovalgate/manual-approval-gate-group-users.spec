PIPELINES-35
# ManualApprovalGate Group Users specs

Pre condition:
  * Validate manual approval gate deployment

## Single Group / Single Approval: PIPELINES-35-TC01
Tags: approvalgate, approvalgate-users, mag-group-user, e2e, taskrun, sanity
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

## Quorum: Partial to Complete: PIPELINES-35-TC02
Tags: approvalgate, approvalgate-users, mag-group-user, e2e, taskrun, sanity
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

## Mixed Entities (User + Group): PIPELINES-35-TC03
Tags: approvalgate, approvalgate-users, mag-group-user, e2e, taskrun, sanity
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

## Rejection Authority: PIPELINES-35-TC04
Tags: approvalgate, approvalgate-users, mag-group-user, e2e, taskrun, sanity
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

## Change of Mind (Approve to Reject): PIPELINES-35-TC05
Tags: approvalgate, approvalgate-users, mag-group-user, e2e, taskrun, sanity
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

## Non-Member Block: PIPELINES-35-TC06
Tags: approvalgate, approvalgate-users, mag-group-user, e2e, taskrun, sanity
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

## Multi-Group Consensus: PIPELINES-35-TC07
Tags: approvalgate, approvalgate-users, mag-group-user, e2e, taskrun, sanity
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

## Overlapping Membership: PIPELINES-35-TC08
Tags: approvalgate, approvalgate-users, mag-group-user, e2e, taskrun, sanity
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

## Timeout Expiry (Short Timeout): PIPELINES-35-TC09
Tags: approvalgate, approvalgate-users, mag-group-user, e2e, taskrun, sanity
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Ensure approval group "group1" has members "user1,user2"
  * Create manual approval gate pipelinerun with approvers "group:group1" required "2" Should "timeout"
  * Verify the latest pipelinerun for "failed" state

## Multi-Group Race:Any One can approve: PIPELINES-35-TC10
Tags: approvalgate, approvalgate-users, mag-group-user, e2e, taskrun, sanity
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

## Mixed Entity Change-of-Mind: PIPELINES-35-TC11
Tags: approvalgate, approvalgate-users, mag-group-user, e2e, taskrun, sanity
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

## Re-approve Completed Task: PIPELINES-35-TC12
Tags: approvalgate, approvalgate-users, mag-group-user, e2e, taskrun, sanity
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

## Impossible Quorum (Short Timeout): PIPELINES-35-TC13
Tags: approvalgate, approvalgate-users, mag-group-user, e2e, taskrun, sanity
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

## Invalid Group Name: PIPELINES-35-TC14
Tags: approvalgate, approvalgate-users, mag-group-user, e2e, taskrun, sanity
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Create manual approval gate pipelinerun with approvers "group:invalid-group" required "1" Should "timeout"
  * User "user1" performs "approve-expect-fail" on the manual approval gate task
  * Verify the latest pipelinerun for "failed" state

## Re-approve Rejected Task: PIPELINES-35-TC15
Tags: approvalgate, approvalgate-users, mag-group-user, e2e, taskrun, sanity
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

## The Late Joiner: PIPELINES-35-TC16
Tags: approvalgate, approvalgate-users, mag-group-user, e2e, taskrun, sanity
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

## The Evicted User: PIPELINES-35-TC17
Tags: approvalgate, approvalgate-users, mag-group-user, e2e, taskrun, sanity
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

## The Switcheroo: PIPELINES-35-TC18
Tags: approvalgate, approvalgate-users, mag-group-user, e2e, taskrun, sanity
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

## Dynamic Quorum: PIPELINES-35-TC19
Tags: approvalgate, approvalgate-users, mag-group-user, e2e, taskrun, sanity
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

## Approval Message Audit: PIPELINES-35-TC20
Tags: approvalgate, approvalgate-users, mag-group-user, e2e, taskrun, sanity
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Ensure approval group "group1" has members "user1"
  * Create manual approval gate pipelinerun with approvers "group:group1" required "1" Should "success"
  * User "user1" performs "approve" on the manual approval gate task with message "PIPELINES-35-TC20 custom approve message from user1"
  * Verify manual approval gate task message contains "PIPELINES-35-TC20 custom approve message from user1"
  * Verify the latest pipelinerun for "successful" state
