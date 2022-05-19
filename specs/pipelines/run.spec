PIPELINES-03
# Verify Pipeline E2E spec

Pre condition:
  * Validate Operator should be installed

## Run sample pipeline: PIPELINES-03-TC01
Tags: e2e, integration, pipelines, non-admin
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Run a sample pipeline that has 2 tasks:
  1. create a file
  2. read file content created by above task
and verify that it runs succesfully

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                  |
      |----|----------------------------------------------|
      |1   |testdata/v1alpha1/pipelinerun/pipelinerun.yaml|
      |2   |testdata/v1beta1/pipelinerun/pipelinerun.yaml |
  * Verify pipelinerun
      |S.NO|pipeline_run_name     |status     |check_label_propagation|
      |----|----------------------|-----------|-----------------------|
      |1   |output-pipeline-run-va|successfull|yes                    |
      |2   |output-pipeline-run-vb|successfull|yes                    |

## Conditional pipeline run: PIPELINES-03-TC02
Tags: e2e, integration, pipelines, non-admin
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                    |
      |----|------------------------------------------------|
      |1   |testdata/v1beta1/pipelinerun/conditional-pr.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name|status     |check_lable_propagation|
      |----|-----------------|-----------|-----------------------|
      |1   |condtional-pr-vb |successfull|no                     |


## Conditional pipeline runs without optional resources: PIPELINES-03-TC03
Tags: e2e, integration, pipelines, non-admin
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                                                     |
      |----|---------------------------------------------------------------------------------|
      |1   |testdata/v1beta1/pipelinerun/conditional-pipelinerun-with-optional-resources.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name                       |status     |check_lable_propagation|
      |----|----------------------------------------|-----------|-----------------------|
      |1   |condtional-pr-without-condition-resource|successfull|no                     |


## Pipelinerun Timeout failure Test: PIPELINES-03-TC04
Tags: e2e, integration, pipelines, non-admin
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                         |
      |----|-----------------------------------------------------|
      |1   |testdata/v1beta1/pipelinerun/pipelineruntimeout.yaml |
  * Verify pipelinerun
      |S.NO|pipeline_run_name|status             |check_lable_propagation|
      |----|-----------------|-------------------|-----------------------|
      |1   |pear             |timeout            |no                     |

## Configure execution results at the Task level Test: PIPELINES-03-TC05
Tags: e2e, integration, pipelines, non-admin
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                          |
      |----|------------------------------------------------------|
      |1   |testdata/v1beta1/pipelinerun/task_results_example.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name |status     |check_lable_propagation|
      |----|------------------|-----------|-----------------------|
      |1   |task-level-results|successfull|no                     |

## Cancel pipelinerun Test: PIPELINES-03-TC06
Tags: e2e, integration, pipelines, non-admin
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                 |
      |----|---------------------------------------------|
      |1   |testdata/v1beta1/pipelinerun/pipelinerun.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name     |status   |check_lable_propagation|
      |----|----------------------|---------|-----------------------|
      |1   |output-pipeline-run-vb|cancelled|no                     |

## Pipelinerun with pipelinespec and taskspec(embedded pipelinerun tests): PIPELINES-03-TC07
Tags: e2e, integration, pipelines, non-admin
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                                                |
      |----|----------------------------------------------------------------------------|
      |1   |testdata/v1beta1/pipelinerun/pipelinerun-with-pipelinespec-and-taskspec.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name                        |status     |check_lable_propagation|
      |----|-----------------------------------------|-----------|-----------------------|
      |1   |pipelinerun-with-pipelinespec-taskspec-vb|successfull|no                     |

## S2I nodejs pipelinerun: PIPELINES-03-TC08
Tags: e2e, integration, pipelines, non-admin, s2i
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                            |
      |----|--------------------------------------------------------|
      |1   |testdata/v1beta1/pipelinerun/s2i-nodejs-pipelinerun.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name|status     |check_lable_propagation|
      |----|-----------------|-----------|-----------------------|
      |1   |nodejs-ex-git-pr |successfull|no                     |