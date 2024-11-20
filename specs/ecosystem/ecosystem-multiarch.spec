PIPELINES-32
# Verify ecosystem E2E spec

Pre condition:
  * Validate Operator should be installed

## jib-maven pipelinerun: PIPELINES-32-TC01
Tags: linux/amd64, ecosystem, non-admin, jib-maven, sanity
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                          |
      |----|------------------------------------------------------|
      |1   |testdata/ecosystem/pipelines/jib-maven.yaml           |
      |2   |testdata/pvc/pvc.yaml                                 |
      |3   |testdata/ecosystem/pipelineruns/jib-maven.yaml        |
  * Verify pipelinerun
      |S.NO|pipeline_run_name|status    |check_label_propagation|
      |----|-----------------|----------|-----------------------|
      |1   |jib-maven-run    |successful|no                     |

## jib-maven P&Z pipelinerun: PIPELINES-32-TC02
Tags: linux/ppc64le, linux/s390x, linux/arm64, ecosystem, non-admin, jib-maven, sanity
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                          |
      |----|------------------------------------------------------|
      |1   |testdata/ecosystem/pipelines/jib-maven-pz.yaml        |
      |2   |testdata/pvc/pvc.yaml                                 |
      |3   |testdata/ecosystem/pipelineruns/jib-maven-pz.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name|status    |check_label_propagation|
      |----|-----------------|----------|-----------------------|
      |1   |jib-maven-pz-run |successful|no                     |

## kn-apply pipelinerun: PIPELINES-32-TC03
Tags: e2e, linux/amd64, ecosystem, non-admin, kn-apply
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                      |
      |----|--------------------------------------------------|
      |1   |testdata/ecosystem/pipelineruns/kn-apply.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name|status    |check_label_propagation|
      |----|-----------------|----------|-----------------------|
      |1   |kn-apply-run     |successful|no                     |

## kn-apply p&z pipelinerun: PIPELINES-32-TC04
Tags: e2e, linux/ppc64le, linux/s390x, ecosystem, non-admin, kn-apply
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                                |
      |----|------------------------------------------------------------|
      |1   |testdata/ecosystem/pipelineruns/kn-apply-multiarch.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name|status    |check_label_propagation|
      |----|-----------------|----------|-----------------------|
      |1   |kn-apply-pz-run  |successful|no                     |

## kn pipelinerun: PIPELINES-32-TC05
Tags: e2e, linux/amd64, ecosystem, non-admin, kn
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                |
      |----|--------------------------------------------|
      |1   |testdata/ecosystem/pipelineruns/kn.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name|status    |check_label_propagation|
      |----|-----------------|----------|-----------------------|
      |1   |kn-run           |successful|no                     |

## kn p&z pipelinerun: PIPELINES-32-TC06
Tags: e2e, linux/ppc64le, linux/s390x, ecosystem, non-admin, kn
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                   |
      |----|-----------------------------------------------|
      |1   |testdata/ecosystem/pipelineruns/kn-pz.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name|status    |check_label_propagation|
      |----|-----------------|----------|-----------------------|
      |1   |kn-pz-run        |successful|no                     |