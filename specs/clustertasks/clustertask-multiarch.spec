PIPELINES-14
# Verify Clustertasks E2E spec

Pre condition:
  * Validate Operator should be installed

## jib-maven pipelinerun: PIPELINES-17-TC01
Tags: e2e, linux/amd64, clustertasks, non-admin, jib-maven, sanity
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                          |
      |----|------------------------------------------------------|
      |1   |testdata/v1/clustertask/pipelines/jib-maven.yaml      |
      |2   |testdata/pvc/pvc.yaml                                 |
      |3   |testdata/v1/clustertask/pipelineruns/jib-maven.yaml   |
  * Verify pipelinerun
      |S.NO|pipeline_run_name|status    |check_label_propagation|
      |----|-----------------|----------|-----------------------|
      |1   |jib-maven-run    |successful|no                     |

## jib-maven P&Z pipelinerun: PIPELINES-17-TC02
Tags: e2e, linux/ppc64le, linux/s390x, linux/arm64, clustertasks, non-admin, jib-maven, sanity
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                          |
      |----|------------------------------------------------------|
      |1   |testdata/v1/clustertask/pipelines/jib-maven-pz.yaml   |
      |2   |testdata/pvc/pvc.yaml                                 |
      |3   |testdata/v1/clustertask/pipelineruns/jib-maven-pz.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name|status    |check_label_propagation|
      |----|-----------------|----------|-----------------------|
      |1   |jib-maven-pz-run |successful|no                     |

## kn-apply pipelinerun: PIPELINES-17-TC03
Tags: e2e, linux/amd64, clustertasks, non-admin, kn-apply
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                      |
      |----|--------------------------------------------------|
      |1   |testdata/v1/clustertask/pipelineruns/kn-apply.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name|status    |check_label_propagation|
      |----|-----------------|----------|-----------------------|
      |1   |kn-apply-run     |successful|no                     |

## kn-apply p&z pipelinerun: PIPELINES-17-TC04
Tags: e2e, linux/ppc64le, linux/s390x, clustertasks, non-admin, kn-apply
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                                |
      |----|------------------------------------------------------------|
      |1   |testdata/v1/clustertask/pipelineruns/kn-apply-multiarch.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name|status    |check_label_propagation|
      |----|-----------------|----------|-----------------------|
      |1   |kn-apply-pz-run  |successful|no                     |

## kn pipelinerun: PIPELINES-17-TC05
Tags: e2e, linux/amd64, clustertasks, non-admin, kn
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                |
      |----|--------------------------------------------|
      |1   |testdata/v1/clustertask/pipelineruns/kn.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name|status    |check_label_propagation|
      |----|-----------------|----------|-----------------------|
      |1   |kn-run           |successful|no                     |

## kn p&z pipelinerun: PIPELINES-17-TC06
Tags: e2e, linux/ppc64le, linux/s390x, clustertasks, non-admin, kn
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                   |
      |----|-----------------------------------------------|
      |1   |testdata/v1/clustertask/pipelineruns/kn-pz.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name|status    |check_label_propagation|
      |----|-----------------|----------|-----------------------|
      |1   |kn-pz-run        |successful|no                     |