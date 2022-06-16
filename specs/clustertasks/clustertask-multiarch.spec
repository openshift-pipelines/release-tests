PIPELINES-14
# Verify Clustertasks E2E spec

Pre condition:
  * Validate Operator should be installed

## jib-maven P&Z pipelinerun: PIPELINES-17-TC01
Tags: linux/ppc64le, linux/s390x, integration, clustertasks, non-admin, jib-maven
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                               |
      |----|-----------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelines/jib-maven-pz.yaml   |
      |2   |testdata/v1beta1/clustertask/pvc/pvc.yaml                  |
      |3   |testdata/v1beta1/clustertask/pipelineruns/jib-maven-pz.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name|status    |check_label_propagation|
      |----|-----------------|----------|-----------------------|
      |1   |jib-maven-pz-run |successful|no                     |

## kn-apply pipelinerun: PIPELINES-17-TC02
Tags: linux/amd64, integration, clustertasks, non-admin, kn-apply
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Create
      |S.NO|resource_dir                                                        |
      |----|--------------------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/serviceaccount/kn-deployer-account.yaml|
      |2   |testdata/v1beta1/clustertask/rolebindings/kn-deployer.yaml          |
  * Create
      |S.NO|resource_dir                                           |
      |----|-------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelineruns/kn-apply.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name|status    |check_label_propagation|
      |----|-----------------|----------|-----------------------|
      |1   |kn-apply-run     |successful|no                     |

## kn-apply p&z pipelinerun: PIPELINES-17-TC03
Tags: linux/ppc64le, linux/s390x, integration, clustertasks, non-admin, kn-apply
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Create
      |S.NO|resource_dir                                                        |
      |----|--------------------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/serviceaccount/kn-deployer-account.yaml|
      |2   |testdata/v1beta1/clustertask/rolebindings/kn-deployer.yaml          |
  * Create
      |S.NO|resource_dir                                              |
      |----|----------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelineruns/kn-apply-pz.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name|status    |check_label_propagation|
      |----|-----------------|----------|-----------------------|
      |1   |kn-apply-run     |successful|no                     |

## kn pipelinerun: PIPELINES-17-TC04
Tags: linux/amd64, integration, clustertasks, non-admin, kn-apply
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Create
      |S.NO|resource_dir                                     |
      |----|-------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelineruns/kn.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name|status    |check_label_propagation|
      |----|-----------------|----------|-----------------------|
      |1   |kn-run           |successful|no                     |

## kn p&z pipelinerun: PIPELINES-17-TC05
Tags: linux/ppc64le, linux/s390x, integration, clustertasks, non-admin, kn-apply
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Create
      |S.NO|resource_dir                                        |
      |----|----------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelineruns/kn-pz.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name|status    |check_label_propagation|
      |----|-----------------|----------|-----------------------|
      |1   |kn-pz-run        |successful|no                     |