PIPELINES-12

# Verify auto-prune E2E

Pre condition:
  * Validate Operator should be installed

## Verify auto prune for taskrun: PIPELINES-12-TC01
Tags: e2e, integration, auto-prune, admin
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

This scenario tests auto prune functionality for taskrun resouce

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
  * Update pruner config "with" keep "2" schedule "*/1 * * * *" resouces "taskrun" and "without" keep-since ""
  * Sleep for "5" seconds
  * Assert if cronjob with prefix "tekton-resource-pruner" is "present" in target namespace
  * "2" taskrun(s) should be present within "120" seconds
  * "5" pipelinerun(s) should be present within "120" seconds
  * Remove auto pruner configuration from config CR
  * Sleep for "5" seconds
  * Assert if cronjob with prefix "tekton-resource-pruner" is "not present" in target namespace

## Verify auto prune for pipelinerun: PIPELINES-12-TC02
Tags: e2e, integration, auto-prune, admin
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

This scenario tests auto prune functionality for pipelinerun resouce

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
  * Update pruner config "with" keep "2" schedule "*/1 * * * *" resouces "pipelinerun" and "without" keep-since ""
  * Sleep for "5" seconds
  * Assert if cronjob with prefix "tekton-resource-pruner" is "present" in target namespace
  * "2" pipelinerun(s) should be present within "120" seconds
  * "7" taskrun(s) should be present within "120" seconds
  * Remove auto pruner configuration from config CR
  * Sleep for "5" seconds
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
  * Update pruner config "with" keep "2" schedule "*/1 * * * *" resouces "pipelinerun,taskrun" and "without" keep-since ""
  * Sleep for "5" seconds
  * Assert if cronjob with prefix "tekton-resource-pruner" is "present" in target namespace
  * "2" pipelinerun(s) should be present within "120" seconds
  * "2" taskrun(s) should be present within "120" seconds
  * Remove auto pruner configuration from config CR
  * Sleep for "5" seconds
  * Assert if cronjob with prefix "tekton-resource-pruner" is "not present" in target namespace

## Verify auto prune with keep-since: PIPELINES-12-TC04
Tags: e2e, integration, auto-prune, admin
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
  * Update pruner config "without" keep "" schedule "*/1 * * * *" resouces "pipelinerun,taskrun" and "with" keep-since "2"
  * Sleep for "5" seconds
  * Assert if cronjob with prefix "tekton-resource-pruner" is "present" in target namespace
  * "5" pipelinerun(s) should be present within "120" seconds
  * "10" taskrun(s) should be present within "120" seconds
  * Remove auto pruner configuration from config CR
  * Sleep for "5" seconds
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
  * Update pruner config "with" keep "2" schedule "*/1 * * * *" resouces "pipelinerun,taskrun" and "without" keep-since ""
  * Sleep for "5" seconds
  * Assert if cronjob with prefix "tekton-resource-pruner" is "present" in target namespace
  * "5" pipelinerun(s) should be present within "120" seconds
  * "10" taskrun(s) should be present within "120" seconds
  * Remove annotation "operator.tekton.dev/prune.skip" from namespace
  * "2" pipelinerun(s) should be present within "120" seconds
  * "2" taskrun(s) should be present within "120" seconds
  * Remove auto pruner configuration from config CR
  * Sleep for "5" seconds
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
  * Update pruner config "with" keep "2" schedule "*/1 * * * *" resouces "pipelinerun,taskrun" and "without" keep-since ""
  * Sleep for "5" seconds
  * Assert if cronjob with prefix "tekton-resource-pruner" is "present" in target namespace
  * "5" pipelinerun(s) should be present within "120" seconds
  * "2" taskrun(s) should be present within "120" seconds
  * Remove annotation "operator.tekton.dev/prune.resources" from namespace
  * "2" pipelinerun(s) should be present within "120" seconds
  * Remove auto pruner configuration from config CR
  * Sleep for "5" seconds
  * Assert if cronjob with prefix "tekton-resource-pruner" is "not present" in target namespace


## Verify auto prune add resources taskrun and pipelinerun per namespace with annotation: PIPELINES-12-TC07
Tags: e2e, integration, auto-prune, admin
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
  * Update pruner config "with" keep "2" schedule "*/1 * * * *" resouces "taskrun" and "without" keep-since ""
  * Sleep for "5" seconds
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
  * Sleep for "5" seconds
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
  * Update pruner config "with" keep "2" schedule "*/1 * * * *" resouces "pipelinerun,taskrun" and "without" keep-since ""
  * Sleep for "5" seconds
  * Assert if cronjob with prefix "tekton-resource-pruner" is "present" in target namespace
  * "3" pipelinerun(s) should be present within "120" seconds
  * "3" taskrun(s) should be present within "120" seconds
  * Remove annotation "operator.tekton.dev/prune.keep" from namespace
  * "2" pipelinerun(s) should be present within "120" seconds
  * Remove auto pruner configuration from config CR
  * Sleep for "5" seconds
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
  * Update pruner config "without" keep "" schedule "*/1 * * * *" resouces "pipelinerun,taskrun" and "with" keep-since "10"
  * Sleep for "5" seconds
  * Assert if cronjob with prefix "tekton-resource-pruner" is "present" in target namespace
  * "5" pipelinerun(s) should be present within "120" seconds
  * "10" taskrun(s) should be present within "120" seconds
  * Remove annotation "operator.tekton.dev/prune.keep-since" from namespace
  * Remove auto pruner configuration from config CR
  * Sleep for "5" seconds
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
  * Update pruner config "without" keep "" schedule "*/1 * * * *" resouces "pipelinerun,taskrun" and "with" keep-since "10"
  * Sleep for "5" seconds
  * Assert if cronjob with prefix "tekton-resource-pruner" is "present" in target namespace
  * "5" pipelinerun(s) should be present within "120" seconds
  * "10" taskrun(s) should be present within "120" seconds
  * Annotate namespace with "operator.tekton.dev/prune.strategy=keep"
  * "2" pipelinerun(s) should be present within "120" seconds
  * "2" taskrun(s) should be present within "120" seconds
  * Remove auto pruner configuration from config CR
  * Sleep for "5" seconds
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
  * Update pruner config "with" keep "2" schedule "*/8 * * * *" resouces "pipelinerun,taskrun" and "without" keep-since ""
  * Sleep for "5" seconds
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
  * Sleep for "5" seconds
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
  * Update pruner config with invalid data "with" keep "2" schedule "*/8 * * * *" resouces "pipelinerun,taskrun" and "with" keep-since "2" and expect error message "validation failed: expected exactly one, got both: spec.pruner.keep, spec.pruner.keep-since"
  * Update pruner config with invalid data "with" keep "2" schedule "*/8 * * * *" resouces "pipelinerun,taskrunas" and "without" keep-since "" and expect error message "validation failed: invalid value: taskrunas: spec.pruner.resources[1]"
  * Update pruner config with invalid data "with" keep "2" schedule "*/8 * * * *" resouces "pipelinerunas,taskrun" and "without" keep-since "" and expect error message "validation failed: invalid value: pipelinerunas: spec.pruner.resources[0]"