PIPELINES-12

# Verify auto-prune E2E
Tags: auto-prune

Pre condition:
  * Validate Operator should be installed

## Verify auto prune for taskrun: PIPELINES-12-TC01
Tags: e2e, integration, auto-prune, admin, sanity
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

This scenario tests auto prune functionality for taskrun resource

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Remove auto pruner configuration from config CR
* Create
      |S.NO|resource_dir                                        |
      |----|----------------------------------------------------|
      |1   |testdata/pruner/pipeline/pipeline-for-pruner.yaml   |
      |2   |testdata/pruner/pipeline/pipelinerun-for-pruner.yaml|
      |3   |testdata/pruner/task/task-for-pruner.yaml           |
      |4   |testdata/pruner/task/taskrun-for-pruner.yaml        |
  * Update pruner config "with" keep "2" schedule "*/1 * * * *" resources "taskrun" and "without" keep-since ""
  * Assert if cronjob with prefix "tekton-resource-pruner" is "present" in target namespace
  * "2" taskrun(s) should be present within "120" seconds
  * "5" pipelinerun(s) should be present within "120" seconds
  * Remove auto pruner configuration from config CR
  * Assert if cronjob with prefix "tekton-resource-pruner" is "not present" in target namespace

## Verify auto prune for pipelinerun: PIPELINES-12-TC02
Tags: e2e, integration, auto-prune, admin
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

This scenario tests auto prune functionality for pipelinerun resource

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Remove auto pruner configuration from config CR
  * Create
      |S.NO|resource_dir                                        |
      |----|----------------------------------------------------|
      |1   |testdata/pruner/pipeline/pipeline-for-pruner.yaml   |
      |2   |testdata/pruner/pipeline/pipelinerun-for-pruner.yaml|
      |3   |testdata/pruner/task/task-for-pruner.yaml           |
      |4   |testdata/pruner/task/taskrun-for-pruner.yaml        |
  * Update pruner config "with" keep "2" schedule "*/1 * * * *" resources "pipelinerun" and "without" keep-since ""
  * Assert if cronjob with prefix "tekton-resource-pruner" is "present" in target namespace
  * "2" pipelinerun(s) should be present within "120" seconds
  * "7" taskrun(s) should be present within "120" seconds
  * Remove auto pruner configuration from config CR
  * Assert if cronjob with prefix "tekton-resource-pruner" is "not present" in target namespace


## Verify auto prune for pipelinerun and taskrun: PIPELINES-12-TC03
Tags: e2e, integration, auto-prune, admin
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

This scenario tests auto prune functionality for pipelinerun and taskrun resources

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Remove auto pruner configuration from config CR
  * Create
      |S.NO|resource_dir                                        |
      |----|----------------------------------------------------|
      |1   |testdata/pruner/pipeline/pipeline-for-pruner.yaml   |
      |2   |testdata/pruner/pipeline/pipelinerun-for-pruner.yaml|
      |3   |testdata/pruner/task/task-for-pruner.yaml           |
      |4   |testdata/pruner/task/taskrun-for-pruner.yaml        |
  * Update pruner config "with" keep "2" schedule "*/1 * * * *" resources "pipelinerun,taskrun" and "without" keep-since ""
  * Assert if cronjob with prefix "tekton-resource-pruner" is "present" in target namespace
  * "2" pipelinerun(s) should be present within "120" seconds
  * "2" taskrun(s) should be present within "120" seconds
  * Remove auto pruner configuration from config CR
  * Assert if cronjob with prefix "tekton-resource-pruner" is "not present" in target namespace

## Verify auto prune with keep-since: PIPELINES-12-TC04
Tags: e2e, integration, auto-prune, admin, sanity
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

This scenario tests auto prune functionality with global strategy keep-since

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Remove auto pruner configuration from config CR
  * Create
      |S.NO|resource_dir                                        |
      |----|----------------------------------------------------|
      |1   |testdata/pruner/pipeline/pipeline-for-pruner.yaml   |
      |2   |testdata/pruner/pipeline/pipelinerun-for-pruner.yaml|
      |3   |testdata/pruner/task/task-for-pruner.yaml           |
      |4   |testdata/pruner/task/taskrun-for-pruner.yaml        |
  * Sleep for "120" seconds
  * Create
      |S.NO|resource_dir                                        |
      |----|----------------------------------------------------|
      |1   |testdata/pruner/pipeline/pipelinerun-for-pruner.yaml|
      |2   |testdata/pruner/task/taskrun-for-pruner.yaml        |
  * Update pruner config "without" keep "" schedule "*/1 * * * *" resources "pipelinerun,taskrun" and "with" keep-since "2"
  * Assert if cronjob with prefix "tekton-resource-pruner" is "present" in target namespace
  * "5" pipelinerun(s) should be present within "120" seconds
  * "10" taskrun(s) should be present within "120" seconds
  * Remove auto pruner configuration from config CR
  * Assert if cronjob with prefix "tekton-resource-pruner" is "not present" in target namespace

## Verify auto prune skip namespace with annotation: PIPELINES-12-TC05
Tags: e2e, integration, auto-prune, admin
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

This scenario tests auto prune functionality for a namespace with annotation "operator.tekton.dev/prune.skip=true".
Pruning should not happen for the resources of a namespace with annotation operator.tekton.dev/prune.skip=true


Steps:
  * Verify ServiceAccount "pipeline" exist
  * Remove auto pruner configuration from config CR
  * Create
      |S.NO|resource_dir                                        |
      |----|----------------------------------------------------|
      |1   |testdata/pruner/pipeline/pipeline-for-pruner.yaml   |
      |2   |testdata/pruner/pipeline/pipelinerun-for-pruner.yaml|
      |3   |testdata/pruner/task/task-for-pruner.yaml           |
      |4   |testdata/pruner/task/taskrun-for-pruner.yaml        |
  * Annotate namespace with "operator.tekton.dev/prune.skip=true"
  * Update pruner config "with" keep "2" schedule "*/1 * * * *" resources "pipelinerun,taskrun" and "without" keep-since ""
  * Assert if cronjob with prefix "tekton-resource-pruner" is "present" in target namespace
  * "5" pipelinerun(s) should be present within "120" seconds
  * "10" taskrun(s) should be present within "120" seconds
  * Remove annotation "operator.tekton.dev/prune.skip" from namespace
  * "2" pipelinerun(s) should be present within "120" seconds
  * "2" taskrun(s) should be present within "120" seconds
  * Remove auto pruner configuration from config CR
  * Assert if cronjob with prefix "tekton-resource-pruner" is "not present" in target namespace


## Verify auto prune add resources taskrun per namespace with annotation: PIPELINES-12-TC06
Tags: e2e, integration, auto-prune, admin
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

This scenario tests auto prune functionality for a namespace with annotation operator.tekton.dev/prune.resources=taskrun.
Only taskruns should get pruned for a namespace with annotation operator.tekton.dev/prune.resources=taskrun

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Remove auto pruner configuration from config CR
  * Create
      |S.NO|resource_dir                                        |
      |----|----------------------------------------------------|
      |1   |testdata/pruner/pipeline/pipeline-for-pruner.yaml   |
      |2   |testdata/pruner/pipeline/pipelinerun-for-pruner.yaml|
      |3   |testdata/pruner/task/task-for-pruner.yaml           |
      |4   |testdata/pruner/task/taskrun-for-pruner.yaml        |
  * Annotate namespace with "operator.tekton.dev/prune.resources=taskrun"
  * Update pruner config "with" keep "2" schedule "*/1 * * * *" resources "pipelinerun,taskrun" and "without" keep-since ""
  * Assert if cronjob with prefix "tekton-resource-pruner" is "present" in target namespace
  * "5" pipelinerun(s) should be present within "120" seconds
  * "2" taskrun(s) should be present within "120" seconds
  * Remove annotation "operator.tekton.dev/prune.resources" from namespace
  * "2" pipelinerun(s) should be present within "120" seconds
  * Remove auto pruner configuration from config CR
  * Assert if cronjob with prefix "tekton-resource-pruner" is "not present" in target namespace


## Verify auto prune add resources taskrun and pipelinerun per namespace with annotation: PIPELINES-12-TC07
Tags: e2e, integration, auto-prune, admin, sanity
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

This scenario tests auto prune functionality for a namespace with annotation operator.tekton.dev/prune.resources=taskrun,pipelinerun.
Both taskruns and pipelineruns should get pruned for a namespace with annotation operator.tekton.dev/prune.resources=taskrun.

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Remove auto pruner configuration from config CR
  * Create
      |S.NO|resource_dir                                        |
      |----|----------------------------------------------------|
      |1   |testdata/pruner/pipeline/pipeline-for-pruner.yaml   |
      |2   |testdata/pruner/pipeline/pipelinerun-for-pruner.yaml|
      |3   |testdata/pruner/task/task-for-pruner.yaml           |
      |4   |testdata/pruner/task/taskrun-for-pruner.yaml        |
  * Annotate namespace with "operator.tekton.dev/prune.resources=pipelinerun,taskrun"
  * Update pruner config "with" keep "2" schedule "*/1 * * * *" resources "taskrun" and "without" keep-since ""
  * Assert if cronjob with prefix "tekton-resource-pruner" is "present" in target namespace
  * "2" pipelinerun(s) should be present within "120" seconds
  * "2" taskrun(s) should be present within "120" seconds
  * Remove annotation "operator.tekton.dev/prune.resources" from namespace
  * Create
      |S.NO|resource_dir                                        |
      |----|----------------------------------------------------|
      |1   |testdata/pruner/pipeline/pipelinerun-for-pruner.yaml|
      |2   |testdata/pruner/task/taskrun-for-pruner.yaml        |
  * "7" pipelinerun(s) should be present within "120" seconds
  * "2" taskrun(s) should be present within "120" seconds
  * Remove auto pruner configuration from config CR
  * Assert if cronjob with prefix "tekton-resource-pruner" is "not present" in target namespace



## Verify auto prune add keep per namespace with annotation with global strategy keep: PIPELINES-12-TC08
Tags: e2e, integration, auto-prune, admin
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

This scenario tests auto prune functionality for a namespace with annotation operator.tekton.dev/prune.keep and the global stratergy keep.
If the globaly strategy and the strategy of namespace is same, no need to define operator.tekton.dev/prune.strategy

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Remove auto pruner configuration from config CR
  * Create
      |S.NO|resource_dir                                        |
      |----|----------------------------------------------------|
      |1   |testdata/pruner/pipeline/pipeline-for-pruner.yaml   |
      |2   |testdata/pruner/pipeline/pipelinerun-for-pruner.yaml|
      |3   |testdata/pruner/task/task-for-pruner.yaml           |
      |4   |testdata/pruner/task/taskrun-for-pruner.yaml        |
  * Annotate namespace with "operator.tekton.dev/prune.keep=3"
  * Update pruner config "with" keep "2" schedule "*/1 * * * *" resources "pipelinerun,taskrun" and "without" keep-since ""
  * Assert if cronjob with prefix "tekton-resource-pruner" is "present" in target namespace
  * "3" pipelinerun(s) should be present within "120" seconds
  * "3" taskrun(s) should be present within "120" seconds
  * Remove annotation "operator.tekton.dev/prune.keep" from namespace
  * "2" pipelinerun(s) should be present within "120" seconds
  * Remove auto pruner configuration from config CR
  * Assert if cronjob with prefix "tekton-resource-pruner" is "not present" in target namespace


## Verify auto prune with keep-since per namespace with global stratergy keep-since: PIPELINES-12-TC09
Tags: e2e, integration, auto-prune, admin
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

This scenario tests auto prune functionality for a namespace with annotation operator.tekton.dev/prune.keep-since and the global stratergy keep.
If the globaly strategy and the strategy of namespace is same, no need to define operator.tekton.dev/prune.strategy

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Remove auto pruner configuration from config CR
  * Create
      |S.NO|resource_dir                                        |
      |----|----------------------------------------------------|
      |1   |testdata/pruner/pipeline/pipeline-for-pruner.yaml   |
      |2   |testdata/pruner/pipeline/pipelinerun-for-pruner.yaml|
      |3   |testdata/pruner/task/task-for-pruner.yaml           |
      |4   |testdata/pruner/task/taskrun-for-pruner.yaml        |
  * Sleep for "120" seconds
  * Create
      |S.NO|resource_dir                                        |
      |----|----------------------------------------------------|
      |1   |testdata/pruner/pipeline/pipelinerun-for-pruner.yaml|
      |2   |testdata/pruner/task/taskrun-for-pruner.yaml        |
  * Annotate namespace with "operator.tekton.dev/prune.keep-since=2"
  * Update pruner config "without" keep "" schedule "*/1 * * * *" resources "pipelinerun,taskrun" and "with" keep-since "10"
  * Assert if cronjob with prefix "tekton-resource-pruner" is "present" in target namespace
  * "5" pipelinerun(s) should be present within "120" seconds
  * "10" taskrun(s) should be present within "120" seconds
  * Remove annotation "operator.tekton.dev/prune.keep-since" from namespace
  * Remove auto pruner configuration from config CR
  * Assert if cronjob with prefix "tekton-resource-pruner" is "not present" in target namespace


## Verify auto prune with keep per namespace with global stratergy keep-since: PIPELINES-12-TC10
Tags: e2e, integration, auto-prune, admin
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

This scenario tests auto prune functionality for a namespace with annotation operator.tekton.dev/prune.keep-since and the global stratergy keep-since.
If the globaly strategy and the strategy of namespace is different, the operator.tekton.dev/prune.strategy=strategy is must

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Remove auto pruner configuration from config CR
  * Create
      |S.NO|resource_dir                                        |
      |----|----------------------------------------------------|
      |1   |testdata/pruner/pipeline/pipeline-for-pruner.yaml   |
      |2   |testdata/pruner/pipeline/pipelinerun-for-pruner.yaml|
      |3   |testdata/pruner/task/task-for-pruner.yaml           |
      |4   |testdata/pruner/task/taskrun-for-pruner.yaml        |
  * Annotate namespace with "operator.tekton.dev/prune.keep=2"
  * Update pruner config "without" keep "" schedule "*/1 * * * *" resources "pipelinerun,taskrun" and "with" keep-since "10"
  * Assert if cronjob with prefix "tekton-resource-pruner" is "present" in target namespace
  * "5" pipelinerun(s) should be present within "120" seconds
  * "10" taskrun(s) should be present within "120" seconds
  * Annotate namespace with "operator.tekton.dev/prune.strategy=keep"
  * "2" pipelinerun(s) should be present within "120" seconds
  * "2" taskrun(s) should be present within "120" seconds
  * Remove auto pruner configuration from config CR
  * Assert if cronjob with prefix "tekton-resource-pruner" is "not present" in target namespace

## Verify auto prune shcedule per namespace: PIPELINES-12-TC11
Tags: e2e, integration, auto-prune, admin
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

This scenrio tests auto prune functionality for a namespace with different schedule by annotating namespace with operator.tekton.dev/prune.schedule
Steps:
  * Verify ServiceAccount "pipeline" exist
  * Remove auto pruner configuration from config CR
  * Create
      |S.NO|resource_dir                                        |
      |----|----------------------------------------------------|
      |1   |testdata/pruner/pipeline/pipeline-for-pruner.yaml   |
      |2   |testdata/pruner/pipeline/pipelinerun-for-pruner.yaml|
      |3   |testdata/pruner/task/task-for-pruner.yaml           |
      |4   |testdata/pruner/task/taskrun-for-pruner.yaml        |
  * Annotate namespace with "operator.tekton.dev/prune.schedule=*/1 * * * *"
  * Update pruner config "with" keep "2" schedule "*/8 * * * *" resources "pipelinerun,taskrun" and "without" keep-since ""
  * Assert if cronjob with prefix "tekton-resource-pruner" is "present" in target namespace
  * "2" pipelinerun(s) should be present within "120" seconds
  * "2" taskrun(s) should be present within "120" seconds
  * Remove annotation "operator.tekton.dev/prune.schedule" from namespace
  * Sleep for "60" seconds
  * Create
      |S.NO|resource_dir                                        |
      |----|----------------------------------------------------|
      |1   |testdata/pruner/pipeline/pipelinerun-for-pruner.yaml|
      |2   |testdata/pruner/task/taskrun-for-pruner.yaml        |
  * "7" pipelinerun(s) should be present within "120" seconds
  * "12" taskrun(s) should be present within "120" seconds
  * Remove auto pruner configuration from config CR
  * Assert if cronjob with prefix "tekton-resource-pruner" is "not present" in target namespace

## Verify auto prune validation: PIPELINES-12-TC12
Tags: e2e, integration, auto-prune, admin
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

This scenrio tests validation of auto pruner config
Steps:
  * Verify ServiceAccount "pipeline" exist
  * Remove auto pruner configuration from config CR
  * Update pruner config with invalid data "with" keep "2" schedule "*/8 * * * *" resources "pipelinerun,taskrun" and "with" keep-since "2" and expect error message "validation failed: expected exactly one, got both: spec.pruner.keep, spec.pruner.keep-since"
  * Update pruner config with invalid data "with" keep "2" schedule "*/8 * * * *" resources "pipelinerun,taskrunas" and "without" keep-since "" and expect error message "validation failed: invalid value: taskrunas: spec.pruner.resources[1]"
  * Update pruner config with invalid data "with" keep "2" schedule "*/8 * * * *" resources "pipelinerunas,taskrun" and "without" keep-since "" and expect error message "validation failed: invalid value: pipelinerunas: spec.pruner.resources[0]"

## Verify auto prune cronjob re-creation for addition of random annotation/label to namespace: PIPELINES-12-TC13
Tags: e2e, integration, auto-prune, admin
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

This scenrio tests if auto prune job is not getting re-created for addition of random annotation to namespace.
Test case fails if the cronjob gets re-created for addition of random annotation to namepsace.
Steps:
  * Verify ServiceAccount "pipeline" exist
  * Remove auto pruner configuration from config CR
  * Update pruner config "with" keep "2" schedule "10 * * * *" resources "pipelinerun,taskrun" and "without" keep-since ""
  * Assert if cronjob with prefix "tekton-resource-pruner" is "present" in target namespace
  * Store name of the cronjob in target namespace with schedule "10 * * * *" to variable "pre-annotation-name"
  * Annotate namespace with "random-annotation=true"
  * Sleep for "5" seconds
  * Store name of the cronjob in target namespace with schedule "10 * * * *" to variable "post-annotation-name"
  * Assert if values stored in variable "pre-annotation-name" and variable "post-annotation-name" are "equal"
  * Remove annotation "random-annotation" from namespace
  * Store name of the cronjob in target namespace with schedule "10 * * * *" to variable "post-annotation-removal-name"
  * Assert if values stored in variable "pre-annotation-name" and variable "post-annotation-removal-name" are "equal"
  * Add label "random=true" to namespace
  * Store name of the cronjob in target namespace with schedule "10 * * * *" to variable "post-label-name"
  * Assert if values stored in variable "pre-annotation-name" and variable "post-label-name" are "equal"
  * Remove label "random" from the namespace
  * Store name of the cronjob in target namespace with schedule "10 * * * *" to variable "post-label-removal-name"
  * Assert if values stored in variable "pre-annotation-name" and variable "post-label-removal-name" are "equal"

## Verify auto prune cronjob contains single container: PIPELINES-12-TC14
Tags: e2e, integration, auto-prune, admin, cronjob, sanity
Component: Operator
Level: Integration
Type: Functional
Importance: Critical
CustomerScenario: yes

Steps:
  * Update pruner config "with" keep "2" schedule "20 * * * *" resources "taskrun" and "without" keep-since ""
  * Create project "test-project-1"
  * Create project "test-project-2"
  * Sleep for "10" seconds
  * Assert pruner cronjob(s) in namespace "target namespace" contains "1" number of container(s)
  * Delete project "test-project-1"
  * Delete project "test-project-2"
  * Remove auto pruner configuration from config CR

## Verify that the operator is up and running after deleting namespace with pruner annotation: PIPELINES-12-TC15
Tags: e2e, integration, auto-prune, admin, cronjob, sanity
Component: Operator
Level: Integration
Type: Functional
Importance: Critical
CustomerScenario: yes

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Remove auto pruner configuration from config CR
  * Assert if cronjob with prefix "tekton-resource-pruner" is "not present" in target namespace
  * Create
      |S.NO|resource_dir                                 |
      |----|---------------------------------------------|
      |1   |testdata/pruner/namespaces/namespace-one.yaml|
  * Assert if cronjob with prefix "tekton-resource-pruner" is "present" in target namespace
  * Delete project "namespace-one"
  * Sleep for "5" seconds
  * Assert if cronjob with prefix "tekton-resource-pruner" is "not present" in target namespace
  * Validate Operator should be installed
  * Create
      |S.NO|resource_dir                                 |
      |----|---------------------------------------------|
      |1   |testdata/pruner/namespaces/namespace-two.yaml|
  * Assert if cronjob with prefix "tekton-resource-pruner" is "present" in target namespace
  * Delete project "namespace-two"
  * Sleep for "5" seconds
  * Assert if cronjob with prefix "tekton-resource-pruner" is "not present" in target namespace