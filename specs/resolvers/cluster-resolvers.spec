PIPELINES-23
# Cluster resolvers spec

## Checking the functionality of cluster resolvers: PIPELINES-23-TC01
Steps:
    * Validate Operator should be installed
    * Create project "resolver-test-tasks"
    * Create
      |S.NO|resource_dir                                                |
      |----|------------------------------------------------------------|
      |1   |testdata/resolvers/tasks/resolver-test-task.yaml            |
      |2   |testdata/resolvers/tasks/resolver-test-task2.yaml           |
    * Create project "resolver-test-pipelines"
    * Create
      |S.NO|resource_dir                                                |
      |----|------------------------------------------------------------|
      |1   |testdata/resolvers/pipelines/resolver-test-pipeline.yaml    |
    * Create project "resolver-test-pipelineruns"
    * Create
      |S.NO|resource_dir                                                              |
      |----|--------------------------------------------------------------------------|
      |1   |testdata/resolvers/pipelineruns/resolver-test-pipelinerun.yaml            |
      |2   |testdata/resolvers/pipelines/resolver-test-pipeline-same-ns.yaml          |
      |3   |testdata/resolvers/pipelineruns/resolver-test-pipelinerun-same-ns.yaml    |
    * Verify pipelinerun
      |S.NO|pipeline_run_name                  |status      |check_label_propagation  |
      |----|-----------------------------------|--------------------------------------|
      |1   |resolver-test-pipelinerun          |successful  |no                       |
      |2   |resolver-test-pipelinerun-same-ns  |successful  |no                       |
    * Delete projects
      


