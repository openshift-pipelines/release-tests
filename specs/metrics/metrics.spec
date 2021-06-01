# Verify openshift pipelines monitoring

Contains scenarios that exercises running openshift pipelines monitoring acceptance tests

Precondition:
  * Validate Operator should be installed

## Openshift pipelines metrics acceptance tests
Tags: e2e, metrics, admin

Steps:
  * Verify job health status metrics
    |S.NO|Job_name                   |Expected_value|
    |----|---------------------------|--------------|
    |1   |node-exporter              |1             |
    |2   |kube-state-metrics         |1             |
    |3   |prometheus-k8s             |1             |
    |4   |prometheus-operator        |1             |
    |5   |alertmanager-main          |1             |
    |6   |tekton-pipelines-controller|1             |
  * Verify pipelines controlPlane metrics