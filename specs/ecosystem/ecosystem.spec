PIPELINES-29
# Verify Ecosystem Tasks E2E spec

Pre condition:
  * Validate Operator should be installed

## buildah pipelinerun: PIPELINES-29-TC01
Tags: e2e, ecosystem, tasks, non-admin, buildah, sanity
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                     |
      |----|-------------------------------------------------|
      |1   |testdata/ecosystem/pipelines/buildah.yaml        |
      |2   |testdata/pvc/pvc.yaml                            |
      |3   |testdata/ecosystem/pipelineruns/buildah.yaml     |
  * Verify pipelinerun
      |S.NO|pipeline_run_name|status    |check_label_propagation|
      |----|-----------------|----------|-----------------------|
      |1   |buildah-run      |successful|no                     |

## buildah disconnected pipelinerun: PIPELINES-29-TC02
Tags: disconnected-e2e, ecosystem, tasks, non-admin, buildah
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                                  |
      |----|--------------------------------------------------------------|
      |1   |testdata/ecosystem/pipelines/buildah.yaml                     |
      |2   |testdata/pvc/pvc.yaml                                         |
      |3   |testdata/ecosystem/pipelineruns/buildah-disconnected.yaml     |
  * Verify pipelinerun
      |S.NO|pipeline_run_name       |status    |check_label_propagation|
      |----|------------------------|----------|-----------------------|
      |1   |buildah-disconnected-run|successful|no                     |

## git-cli pipelinerun: PIPELINES-29-TC03
Tags: e2e, ecosystem, tasks, non-admin, git-cli
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                     |
      |----|-------------------------------------------------|
      |1   |testdata/ecosystem/pipelines/git-cli.yaml        |
      |2   |testdata/pvc/pvc.yaml                            |
      |3   |testdata/ecosystem/pipelineruns/git-cli.yaml     |
  * Verify pipelinerun
      |S.NO|pipeline_run_name|status    |check_label_propagation|
      |----|-----------------|----------|-----------------------|
      |1   |git-cli-run      |successful|no                     |

## git-cli read private repo pipelinerun: PIPELINES-29-TC04
Tags: e2e, ecosystem, non-admin, git-cli
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                               |
      |----|-----------------------------------------------------------|
      |1   |testdata/ecosystem/pipelines/git-cli-read-private.yaml     |
      |2   |testdata/pvc/pvc.yaml                                      |
      |3   |testdata/ecosystem/secrets/ssh-key.yaml                    |
  * Link secret "ssh-key" to service account "pipeline"
  * Create
      |S.NO|resource_dir                                                  |
      |----|--------------------------------------------------------------|
      |1   |testdata/ecosystem/pipelineruns/git-cli-read-private.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name       |status    |check_label_propagation|
      |----|------------------------|----------|-----------------------|
      |1   |git-cli-read-private-run|successful|no                     |

## git-cli read private repo using different service account pipelinerun: PIPELINES-29-TC05
Tags: e2e, ecosystem, non-admin, git-cli
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                               |
      |----|-----------------------------------------------------------|
      |1   |testdata/ecosystem/pipelines/git-cli-read-private.yaml     |
      |2   |testdata/pvc/pvc.yaml                                      |
      |3   |testdata/ecosystem/secrets/ssh-key.yaml                    |
      |4   |testdata/ecosystem/serviceaccount/ssh-sa.yaml              |
      |5   |testdata/ecosystem/rolebindings/ssh-sa-scc.yaml            |
  * Link secret "ssh-key" to service account "ssh-sa"
  * Create
      |S.NO|resource_dir                                                     |
      |----|-----------------------------------------------------------------|
      |1   |testdata/ecosystem/pipelineruns/git-cli-read-private-sa.yaml     |
  * Verify pipelinerun
      |S.NO|pipeline_run_name          |status    |check_label_propagation|
      |----|---------------------------|----------|-----------------------|
      |1   |git-cli-read-private-sa-run|successful|no                     |
      
## git-clone read private repo taskrun PIPELINES-29-TC06
Tags: e2e, ecosystem, non-admin, git-clone, sanity
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical
CustomerScenario: yes

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      | S.NO | resource_dir                                                  |
      |------|---------------------------------------------------------------|
      | 1    | testdata/ecosystem/pipelines/git-clone-read-private.yaml      |
      | 2    | testdata/pvc/pvc.yaml                                         |
      | 3    | testdata/ecosystem/secrets/ssh-key.yaml                       |
  * Link secret "ssh-key" to service account "pipeline"
  * Create
      | S.NO | resource_dir                                                    |
      |------|-----------------------------------------------------------------|
      | 1    | testdata/ecosystem/pipelineruns/git-clone-read-private.yaml     |
  * Verify pipelinerun
      | S.NO | pipeline_run_name                   | status     | check_label_propagation |
      |------|-------------------------------------|------------|-------------------------|
      | 1    | git-clone-read-private-pipeline-run | successful | no                      |

## git-clone read private repo using different service account taskrun PIPELINES-29-TC07
Tags: e2e, ecosystem, non-admin, git-clone
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      | S.NO | resource_dir                                                  |
      |------|---------------------------------------------------------------|
      | 1    | testdata/ecosystem/pipelines/git-clone-read-private.yaml      |
      | 2    | testdata/pvc/pvc.yaml                                         |
      | 3    | testdata/ecosystem/secrets/ssh-key.yaml                       |
      | 4    | testdata/ecosystem/serviceaccount/ssh-sa.yaml                 |
      | 5    | testdata/ecosystem/rolebindings/ssh-sa-scc.yaml               |
  * Link secret "ssh-key" to service account "ssh-sa"
  * Create
      | S.NO | resource_dir                                                       |
      |------|--------------------------------------------------------------------|
      | 1    | testdata/ecosystem/pipelineruns/git-clone-read-private-sa.yaml|
  * Verify pipelinerun
      | S.NO | pipeline_run_name                      | status     | check_label_propagation |
      |------|----------------------------------------|------------|-------------------------|
      | 1    | git-clone-read-private-pipeline-sa-run | successful | no                      |

## openshift-client pipelinerun: PIPELINES-29-TC08
Tags: e2e, ecosystem, tasks, non-admin, openshift-client
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                              |
      |----|----------------------------------------------------------|
      |1   |testdata/ecosystem/pipelineruns/openshift-client.yaml     |
  * Verify pipelinerun
      |S.NO|pipeline_run_name   |status    |check_label_propagation|
      |----|--------------------|----------|-----------------------|
      |1   |openshift-client-run|successful|no                     |

## skopeo-copy pipelinerun: PIPELINES-29-TC09
Tags: e2e, ecosystem, tasks, non-admin, skopeo-copy
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                         |
      |----|-----------------------------------------------------|
      |1   |testdata/ecosystem/pipelineruns/skopeo-copy.yaml     |
  * Verify pipelinerun
      |S.NO|pipeline_run_name|status    |check_label_propagation|
      |----|-----------------|----------|-----------------------|
      |1   |skopeo-copy-run  |successful|no                     |

## tkn pipelinerun: PIPELINES-29-TC10
Tags: e2e, ecosystem, tasks, non-admin, tkn
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                 |
      |----|---------------------------------------------|
      |1   |testdata/ecosystem/pipelineruns/tkn.yaml     |
  * Verify pipelinerun
      |S.NO|pipeline_run_name|status    |check_label_propagation|
      |----|-----------------|----------|-----------------------|
      |1   |tkn-run          |successful|no                     |

## tkn pac pipelinerun: PIPELINES-29-TC11
Tags: e2e, ecosystem, tasks, non-admin, tkn
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                     |
      |----|-------------------------------------------------|
      |1   |testdata/ecosystem/pipelineruns/tkn-pac.yaml     |
  * Verify pipelinerun
      |S.NO|pipeline_run_name|status    |check_label_propagation|
      |----|-----------------|----------|-----------------------|
      |1   |tkn-pac-run      |successful|no                     |
  * Verify "tkn-pac" version from the pipelinerun logs

## tkn version pipelinerun: PIPELINES-29-TC12
Tags: e2e, ecosystem, tasks, non-admin, tkn
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                         |
      |----|-----------------------------------------------------|
      |1   |testdata/ecosystem/pipelineruns/tkn-version.yaml     |
  * Verify pipelinerun
      |S.NO|pipeline_run_name|status    |check_label_propagation|
      |----|-----------------|----------|-----------------------|
      |1   |tkn-version-run  |successful|no                     |
  * Verify "tkn" version from the pipelinerun logs

## maven pipelinerun: PIPELINES-29-TC13
Tags: e2e, ecosystem, tasks, non-admin, maven
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                          |
      |----|------------------------------------------------------|
      |1   |testdata/ecosystem/pipelines/maven.yaml               |
      |2   |testdata/pvc/pvc.yaml                                 |
      |3   |testdata/ecosystem/configmaps/maven-settings.yaml     |
      |4   |testdata/ecosystem/pipelineruns/maven.yaml            |
  * Verify pipelinerun
      |S.NO|pipeline_run_name|status    |check_label_propagation|
      |----|-----------------|----------|-----------------------|
      |1   |maven-run        |successful|no                     |

## Test the functionality of step action resolvers: PIPELINES-29-TC14
Tags: e2e, sanity, ecosystem, non-admin
Component: Resolvers
Level: Integration
Type: Functional
Importance: High

Steps:
    * Verify ServiceAccount "pipeline" exist
    * Create 
      |S.NO|resource_dir                                                     |
      |----|-----------------------------------------------------------------|
      |1   |testdata/ecosystem/tasks/git-clone-stepaction.yaml               |
      |2   |testdata/pvc/pvc.yaml                                            |
      |3   |testdata/ecosystem/pipelineruns/git-clone-stepaction.yaml        |
    * Verify pipelinerun
      |S.NO|pipeline_run_name                  |status      |check_label_propagation  |
      |----|-----------------------------------|--------------------------------------|
      |1   |git-clone-stepaction-run           |successful  |no                       |
