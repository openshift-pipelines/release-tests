PIPELINES-16
# Verify Clustertasks E2E spec

Pre condition:
  * Validate Operator should be installed

## buildah pipelinerun: PIPELINES-16-TC01
Tags: e2e, clustertasks, non-admin, buildah, sanity
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                          |
      |----|------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelines/buildah.yaml   |
      |2   |testdata/v1beta1/clustertask/pvc/pvc.yaml             |
      |3   |testdata/v1beta1/clustertask/pipelineruns/buildah.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name|status    |check_label_propagation|
      |----|-----------------|----------|-----------------------|
      |1   |buildah-run      |successful|no                     |

## buildah disconnected pipelinerun: PIPELINES-16-TC02
Tags: disconnected-e2e, clustertasks, non-admin, buildah
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                                       |
      |----|-------------------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelines/buildah.yaml                |
      |2   |testdata/v1beta1/clustertask/pvc/pvc.yaml                          |
      |3   |testdata/v1beta1/clustertask/pipelineruns/buildah-disconnected.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name       |status    |check_label_propagation|
      |----|------------------------|----------|-----------------------|
      |1   |buildah-disconnected-run|successful|no                     |

## git-cli pipelinerun: PIPELINES-16-TC03
Tags: e2e, clustertasks, non-admin, git-cli
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                          |
      |----|------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelines/git-cli.yaml   |
      |2   |testdata/v1beta1/clustertask/pvc/pvc.yaml             |
      |3   |testdata/v1beta1/clustertask/pipelineruns/git-cli.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name|status    |check_label_propagation|
      |----|-----------------|----------|-----------------------|
      |1   |git-cli-run      |successful|no                     |

## git-cli read private repo pipelinerun: PIPELINES-16-TC04
Tags: e2e, clustertasks, non-admin, git-cli
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                                    |
      |----|----------------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelines/git-cli-read-private.yaml|
      |2   |testdata/v1beta1/clustertask/pvc/pvc.yaml                       |
      |3   |testdata/v1beta1/clustertask/secrets/ssh-key.yaml               |
  * Link secret "ssh-key" to service account "pipeline"
  * Create
      |S.NO|resource_dir                                                       |
      |----|-------------------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelineruns/git-cli-read-private.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name       |status    |check_label_propagation|
      |----|------------------------|----------|-----------------------|
      |1   |git-cli-read-private-run|successful|no                     |

## git-cli read private repo using different service account pipelinerun: PIPELINES-16-TC05
Tags: e2e, clustertasks, non-admin, git-cli
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                                    |
      |----|----------------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelines/git-cli-read-private.yaml|
      |2   |testdata/v1beta1/clustertask/pvc/pvc.yaml                       |
      |3   |testdata/v1beta1/clustertask/secrets/ssh-key.yaml               |
      |4   |testdata/v1beta1/clustertask/serviceaccount/ssh-sa.yaml         |
      |5   |testdata/v1beta1/clustertask/rolebindings/ssh-sa-scc.yaml       |
  * Link secret "ssh-key" to service account "ssh-sa"
  * Create
      |S.NO|resource_dir                                                          |
      |----|----------------------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelineruns/git-cli-read-private-sa.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name          |status    |check_label_propagation|
      |----|---------------------------|----------|-----------------------|
      |1   |git-cli-read-private-sa-run|successful|no                     |
      
## maven pipelinerun: PIPELINES-16-TC06
Tags: e2e, clustertasks, non-admin, maven
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                               |
      |----|-----------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelines/maven.yaml          |
      |2   |testdata/v1beta1/clustertask/pvc/pvc.yaml                  |
      |3   |testdata/v1beta1/clustertask/configmaps/maven-settings.yaml|
      |4   |testdata/v1beta1/clustertask/pipelineruns/maven.yaml       |
  * Verify pipelinerun
      |S.NO|pipeline_run_name|status    |check_label_propagation|
      |----|-----------------|----------|-----------------------|
      |1   |maven-run        |successful|no                     |

## openshift-client pipelinerun: PIPELINES-16-TC07
Tags: e2e, clustertasks, non-admin, openshift-client
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                                   |
      |----|---------------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelineruns/openshift-client.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name   |status    |check_label_propagation|
      |----|--------------------|----------|-----------------------|
      |1   |openshift-client-run|successful|no                     |

## skopeo-copy pipelinerun: PIPELINES-16-TC08
Tags: e2e, clustertasks, non-admin, skopeo-copy
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                              |
      |----|----------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelineruns/skopeo-copy.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name|status    |check_label_propagation|
      |----|-----------------|----------|-----------------------|
      |1   |skopeo-copy-run  |successful|no                     |

## tkn pipelinerun: PIPELINES-16-TC09
Tags: e2e, clustertasks, non-admin, tkn
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                      |
      |----|--------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelineruns/tkn.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name|status    |check_label_propagation|
      |----|-----------------|----------|-----------------------|
      |1   |tkn-run          |successful|no                     |

## tkn pac pipelinerun: PIPELINES-16-TC10
Tags: e2e, clustertasks, non-admin, tkn
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                          |
      |----|------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelineruns/tkn-pac.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name|status    |check_label_propagation|
      |----|-----------------|----------|-----------------------|
      |1   |tkn-pac-run      |successful|no                     |

## tkn version pipelinerun: PIPELINES-16-TC11
Tags: e2e, clustertasks, non-admin, tkn
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                              |
      |----|----------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelineruns/tkn-version.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name|status    |check_label_propagation|
      |----|-----------------|----------|-----------------------|
      |1   |tkn-version-run  |successful|no                     |
      
## git-clone read private repo taskrun PIPELINES-16-TC12
Tags: e2e, clustertasks, non-admin, git-clone, sanity
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical
CustomerScenario: yes

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      | S.NO | resource_dir                                                       |
      |------|--------------------------------------------------------------------|
      | 1    | testdata/v1beta1/clustertask/pipelines/git-clone-read-private.yaml |
      | 2    | testdata/v1beta1/clustertask/pvc/pvc.yaml                          |
      | 3    | testdata/v1beta1/clustertask/secrets/ssh-key.yaml                  |
  * Link secret "ssh-key" to service account "pipeline"
  * Create
      | S.NO | resource_dir                                                          |
      |------|-----------------------------------------------------------------------|
      | 1    | testdata/v1beta1/clustertask/pipelineruns/git-clone-read-private.yaml |
  * Verify pipelinerun
      | S.NO | pipeline_run_name                   | status     | check_label_propagation |
      |------|-------------------------------------|------------|-------------------------|
      | 1    | git-clone-read-private-pipeline-run | successful | no                      |

## git-clone read private repo using different service account taskrun PIPELINES-16-TC13
Tags: e2e, clustertasks, non-admin, git-clone
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      | S.NO | resource_dir                                                       |
      |------|--------------------------------------------------------------------|
      | 1    | testdata/v1beta1/clustertask/pipelines/git-clone-read-private.yaml |
      | 2    | testdata/v1beta1/clustertask/pvc/pvc.yaml                          |
      | 3    | testdata/v1beta1/clustertask/secrets/ssh-key.yaml                  |
      | 4    | testdata/v1beta1/clustertask/serviceaccount/ssh-sa.yaml            |
      | 5    | testdata/v1beta1/clustertask/rolebindings/ssh-sa-scc.yaml          |
  * Link secret "ssh-key" to service account "ssh-sa"
  * Create
      | S.NO | resource_dir                                                             |
      |------|--------------------------------------------------------------------------|
      | 1    | testdata/v1beta1/clustertask/pipelineruns/git-clone-read-private-sa.yaml |
  * Verify pipelinerun
      | S.NO | pipeline_run_name                      | status     | check_label_propagation |
      |------|----------------------------------------|------------|-------------------------|
      | 1    | git-clone-read-private-pipeline-sa-run | successful | no                      |
      
