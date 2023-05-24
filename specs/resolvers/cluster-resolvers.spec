PIPELINES-23
# Cluster resolvers spec

## Checking the functionality of cluster resolvers#1: PIPELINES-23-TC01
Tags: e2e, sanity
Component: Resolvers
Level: Integration
Type: Functional
Importance: High

Steps:
    * Check if project "resolver-tasks" exists
    * Switch to project "resolver-tasks"
    * Create
      |S.NO|resource_dir                                                |
      |----|------------------------------------------------------------|
      |1   |testdata/resolvers/tasks/resolver-task2.yaml                |
    * Check if project "resolver-pipelines" exists
    * Switch to project "resolver-pipelines"
    * Create
      |S.NO|resource_dir                                                |
      |----|------------------------------------------------------------|
      |1   |testdata/resolvers/pipelines/resolver-pipeline.yaml         |
    * Check if project "resolver-pipelineruns" exists
    * Switch to project "resolver-pipelineruns"
    * Create
      |S.NO|resource_dir                                                              |
      |----|--------------------------------------------------------------------------|
      |1   |testdata/resolvers/pipelineruns/resolver-pipelinerun.yaml                 |
    * Verify pipelinerun
      |S.NO|pipeline_run_name                  |status      |check_label_propagation  |
      |----|-----------------------------------|--------------------------------------|
      |1   |resolver-pipelinerun               |successful  |no                       |

## Checking the functionality of cluster resolvers#2: PIPELINES-23-TC02
Tags: e2e, sanity
Component: Resolvers
Level: Integration
Type: Functional
Importance: High

Steps: 
    * Check if project "resolver-tasks" exists
    * Switch to project "resolver-tasks"
    * Create 
      |S.NO|resource_dir                                                |
      |----|------------------------------------------------------------|            
      |1   |testdata/resolvers/tasks/resolver-task.yaml                 |
    * Check if project "resolver-pipelineruns" exists
    * Switch to project "resolver-pipelineruns"
    * Create
      |S.NO|resource_dir                                                              |
      |----|--------------------------------------------------------------------------|
      |1   |testdata/resolvers/pipelines/resolver-pipeline-same-ns.yaml               |
      |2   |testdata/resolvers/pipelineruns/resolver-pipelinerun-same-ns.yaml         |
    * Verify pipelinerun
      |S.NO|pipeline_run_name                  |status      |check_label_propagation  |
      |----|-----------------------------------|--------------------------------------|
      |1   |resolver-pipelinerun-same-ns       |successful  |no                       |

      


