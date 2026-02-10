PIPELINES-35

# Verify Tekton Pruner Functionality
Tags: pruner, tekton-pruner

Pre condition:
  * Validate Operator should be installed

## Enable Tekton Pruner & Validate Deployment Status: TC-01
Tags: e2e, integration, pruner, admin, deployment
Component: Operator
Level: Integration
Type: Deployment
Importance: Critical

This scenario tests pruner migration from legacy to new tekton-pruner and verifies deployment status

Steps:
  * "Disable" legacy pruner
  * "Enable" tekton-pruner
  * Validate tekton-pruner deployment

## Webhook: Negative Values & Invalid Type: TC-02
Tags: e2e, integration, pruner, admin, validation
Component: Operator
Level: Integration
Type: Validation
Importance: Critical

This scenario tests webhook validation for pruner global-config: negative values (e.g. ttlSecondsAfterFinished) and invalid types (string instead of int).

Steps:
  * Update tekton-pruner config with "ttlSecondsAfterFinished" as "-1" and expect message "ttlSecondsAfterFinished cannot be negative"
  * Update tekton-pruner config with "ttlSecondsAfterFinished" as "60s" and expect message "cannot unmarshal string"
  * Update tekton-pruner config with "successfulHistoryLimit" as "not-a-number" and expect message "cannot unmarshal string"

## Global TTL Expiry for PipelineRuns: TC-03
Tags: e2e, integration, pruner, admin, functional
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

This scenario tests global TTL expiry functionality for PipelineRuns

Steps:
  * Update tekton-pruner config with "ttlSecondsAfterFinished" as "60" and expect message ""
  * Create
      |S.NO|resource_dir                                        |
      |----|----------------------------------------------------|
      |1   |testdata/pruner/pipeline/pipeline-for-pruner.yaml   |
      |2   |testdata/pruner/pipeline/pipelinerun-for-pruner.yaml|
  * "5" pipelinerun(s) with status "Succeeded" should be present within "60" seconds
  * Sleep for "60" seconds
  * "0" pipelinerun(s) should be present within "30" seconds


## Global TTL Expiry for TaskRuns: TC-04
Tags: integration, pruner, admin, functional
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

This scenario tests global TTL expiry functionality for TaskRuns

Steps:
  * Update tekton-pruner config with "ttlSecondsAfterFinished" as "60" and expect message ""
  * Create
      |S.NO|resource_dir                                |
      |----|--------------------------------------------|
      |1   |testdata/pruner/task/task-for-pruner.yaml   |
      |2   |testdata/pruner/task/taskrun-for-pruner.yaml|
  * "5" taskrun(s) should be present within "60" seconds
  * Sleep for "60" seconds
  * "0" taskrun(s) should be present within "30" seconds

## Successful History Limit: TC-05
Tags: e2e, integration, pruner, admin, functional
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

This scenario tests successfulHistoryLimit: only the N most recent successful PipelineRuns are kept.

Steps:
  * Create
      |S.NO|resource_dir                                        |
      |----|----------------------------------------------------|
      |1   |testdata/pruner/pipeline/pipeline-for-pruner.yaml   |
      |2   |testdata/pruner/pipeline/pipelinerun-for-pruner.yaml|
  * "5" pipelinerun(s) with status "Succeeded" should be present within "60" seconds
  * Update tekton-pruner config with "successfulHistoryLimit" as "2" and expect message ""
  * "2" pipelinerun(s) with status "Succeeded" should be present within "60" seconds

## Failed History Limit: TC-06
Tags: e2e, integration, pruner, admin, functional
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

This scenario tests failedHistoryLimit: only the N most recent failed PipelineRuns are kept.

Steps:
  * Create
      |S.NO|resource_dir                                                |
      |----|------------------------------------------------------------|
      |1   |testdata/pruner/pipeline/pipeline-fail-for-pruner.yaml      |
      |2   |testdata/pruner/pipeline/pipelinerun-fail-for-pruner.yaml   |
  * "5" pipelinerun(s) with status "Failed" should be present within "60" seconds
  * Update tekton-pruner config with "failedHistoryLimit" as "3" and expect message ""
  * "3" pipelinerun(s) with status "Failed" should be present within "60" seconds

## Mixed History Limits: TC-07
Tags: e2e, integration, pruner, admin, functional
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

This scenario tests successfulHistoryLimit and failedHistoryLimit together: 
set both limits to 5, create 5 successful and 5 failed PipelineRuns (10 total). 
Then set successfulHistoryLimit to 2 and failedHistoryLimit to 3. After the pruner runs, verify exactly 5 PipelineRuns remain: 2 Succeeded and 3 Failed.

Steps:
  * Update tekton-pruner config with "successfulHistoryLimit" as "5" and expect message ""
  * Update tekton-pruner config with "failedHistoryLimit" as "5" and expect message ""
  * Create
      |S.NO|resource_dir                                                |
      |----|------------------------------------------------------------|
      |1   |testdata/pruner/pipeline/pipeline-for-pruner.yaml           |
      |2   |testdata/pruner/pipeline/pipelinerun-for-pruner.yaml        |
      |3   |testdata/pruner/pipeline/pipeline-fail-for-pruner.yaml      |
      |4   |testdata/pruner/pipeline/pipelinerun-fail-for-pruner.yaml   |
  * "10" pipelinerun(s) should be present within "60" seconds
  * Update tekton-pruner config with "successfulHistoryLimit" as "2" and expect message ""
  * Update tekton-pruner config with "failedHistoryLimit" as "3" and expect message ""
  * "5" pipelinerun(s) should be present within "60" seconds
  * "2" pipelinerun(s) with status "Succeeded" should be present within "60" seconds
  * "3" pipelinerun(s) with status "Failed" should be present within "60" seconds

## Namespace Config Override Error: TC-08
Tags: e2e, integration, pruner, admin, hierarchy, validation
Component: Operator
Level: Integration
Type: Hierarchy
Importance: Critical

This scenario tests that the webhook rejects a namespace pruner config when it exceeds the global config (e.g. namespace TTL greater than global TTL).

Steps:
  * Update tekton-pruner config with "ttlSecondsAfterFinished" as "60" and expect message ""
  * Update tekton-pruner config with "namespaces.dev.ttlSecondsAfterFinished" as "300" and expect message "ttlSecondsAfterFinished (300) cannot exceed global limit"


## Label Selector Match & Mismatch: TC-09
Tags: e2e, integration, pruner, admin, selectors, functional
Component: Operator
Level: Integration
Type: Selectors
Importance: Critical

This scenario tests label-selector TTL: 
(1) Match: default TTL=60s, label type=ci has short TTL — PipelineRun with type: ci is deleted after 30s. 
(2) Mismatch: PipelineRun with type: nightly does not match selector — retained for global TTL.

Steps:
  * Update tekton-pruner config with "ttlSecondsAfterFinished" as "60" and expect message ""
  * Update tekton-pruner config with "enforcedConfigLevel" as "namespace" and expect message ""
  * Create
      |S.NO|resource_dir                                                 |
      |----|-------------------------------------------------------------|
      |1   |testdata/pruner/pipeline/pipeline-for-pruner.yaml            |
      |2   |testdata/pruner/configmap/label-prune-ns-cm.yaml             |
      |3   |testdata/pruner/pipeline/pipelinerun-label-for-pruner.yaml   |
  * "1" pipelinerun(s) with status "Succeeded" should be present within "30" seconds
  * Sleep for "30" seconds
  * "0" pipelinerun(s) should be present within "15" seconds
  * Create
      |S.NO|resource_dir                                                      |
      |----|----------------------------------------------------------------  |
      |1   |testdata/pruner/pipeline/pipelinerun-nightly-label-for-pruner.yaml|
  * "1" pipelinerun(s) with status "Succeeded" should be present within "30" seconds
  * Sleep for "30" seconds
  * "1" pipelinerun(s) should be present within "15" seconds

## Annotation Selector: TC-10
Tags: e2e, integration, pruner, admin, selectors, functional
Component: Operator
Level: Integration
Type: Selectors
Importance: Critical

This scenario tests annotation selector: config matches annotation prune=true. 
A PipelineRun with annotation prune: true is pruned according to the specific rule.

Steps:
  * Update tekton-pruner config with "ttlSecondsAfterFinished" as "60" and expect message ""
  * Update tekton-pruner config with "enforcedConfigLevel" as "namespace" and expect message ""
  * Create
      |S.NO|resource_dir                                                      |
      |----|------------------------------------------------------------------|
      |1   |testdata/pruner/pipeline/pipeline-for-pruner.yaml                 |
      |2   |testdata/pruner/configmap/annotation-prune-ns-cm.yaml             |
      |3   |testdata/pruner/pipeline/pipelinerun-annotation-for-pruner.yaml   |
  * "1" pipelinerun(s) with status "Succeeded" should be present within "30" seconds
  * Sleep for "30" seconds
  * "0" pipelinerun(s) should be present within "15" seconds

## AND Logic (Label + Annotation): TC-11
Tags: e2e, integration, pruner, admin, selectors, functional
Component: Operator
Level: Integration
Type: Selectors
Importance: Critical

This scenario tests AND logic for selectors: the rule matches only when both label type=ci and annotation prune=true are present. 
A PipelineRun with both that label and annotation is pruned according to the rule; verify it is deleted after the configured TTL.

Steps:
  * Update tekton-pruner config with "ttlSecondsAfterFinished" as "60" and expect message ""
  * Update tekton-pruner config with "enforcedConfigLevel" as "namespace" and expect message ""
  * Create
      |S.NO|resource_dir                                                     |
      |----|-----------------------------------------------------------------|
      |1   |testdata/pruner/pipeline/pipeline-for-pruner.yaml                |
      |2   |testdata/pruner/configmap/label-and-annotation-ns-cm.yaml        |
      |3   |testdata/pruner/pipeline/pipelinerun-label-annotation.yaml       |
  * "1" pipelinerun(s) with status "Succeeded" should be present within "30" seconds
  * Sleep for "30" seconds
  * "0" pipelinerun(s) should be present within "15" seconds

Teardown:
  * "Enable" legacy pruner
  * "Disable" tekton-pruner