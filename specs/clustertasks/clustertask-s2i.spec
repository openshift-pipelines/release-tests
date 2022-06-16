PIPELINES-14
# Verify Clustertasks E2E spec

Pre condition:
  * Validate Operator should be installed


## S2I nodejs pipelinerun: PIPELINES-14-TC01
Tags: e2e, integration, clustertasks, non-admin, s2i
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
      |1   |nodejs-ex-git-pr |successful |no                     |

## S2I dotnet pipelinerun dotnetcore-3.1 version: PIPELINES-14-TC02
Tags: e2e, integration, clustertasks, non-admin, s2i
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                                     |
      |----|-----------------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelines/s2i-dotnet.yaml           |
      |2   |testdata/v1beta1/clustertask/pvc/pvc.yaml                        |
      |3   |testdata/v1beta1/clustertask/pipelineruns/s2i-dotnet-31-ubi8.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name      |status     |check_lable_propagation|
      |----|-----------------------|-----------|-----------------------|
      |1   |s2i-dotnet-31-ubi8-run |successful |no                     |

## S2I dotnet pipelinerun dotnet-5.0 version: PIPELINES-14-TC03
Tags: e2e, integration, clustertasks, non-admin, s2i
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                                     |
      |----|-----------------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelines/s2i-dotnet.yaml           |
      |2   |testdata/v1beta1/clustertask/pvc/pvc.yaml                        |
      |3   |testdata/v1beta1/clustertask/pipelineruns/s2i-dotnet-50-ubi8.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name      |status     |check_lable_propagation|
      |----|-----------------------|-----------|-----------------------|
      |1   |s2i-dotnet-50-ubi8-run |successful |no                     |

## S2I dotnet pipelinerun dotnet-6.0 version: PIPELINES-14-TC04
Tags: e2e, integration, clustertasks, non-admin, s2i
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                                     |
      |----|-----------------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelines/s2i-dotnet.yaml           |
      |2   |testdata/v1beta1/clustertask/pvc/pvc.yaml                        |
      |3   |testdata/v1beta1/clustertask/pipelineruns/s2i-dotnet-60-ubi8.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name      |status     |check_lable_propagation|
      |----|-----------------------|-----------|-----------------------|
      |1   |s2i-dotnet-60-ubi8-run |successful |no                     |

## S2I go pipelinerun go-1.16.7-ubi7 version: PIPELINES-14-TC05
Tags: e2e, integration, clustertasks, non-admin, s2i
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                                   |
      |----|---------------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelines/s2i-go.yaml             |
      |2   |testdata/v1beta1/clustertask/pvc/pvc.yaml                      |
      |3   |testdata/v1beta1/clustertask/pipelineruns/s2i-go-1167-ubi7.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name    |status     |check_lable_propagation|
      |----|---------------------|-----------|-----------------------|
      |1   |s2i-go-1167-ubi7-run |successful |no                     |

## S2I go pipelinerun go-1.16.7-ubi8 version: PIPELINES-14-TC06
Tags: e2e, integration, clustertasks, non-admin, s2i
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                                   |
      |----|---------------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelines/s2i-go.yaml             |
      |2   |testdata/v1beta1/clustertask/pvc/pvc.yaml                      |
      |3   |testdata/v1beta1/clustertask/pipelineruns/s2i-go-1167-ubi8.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name    |status     |check_lable_propagation|
      |----|---------------------|-----------|-----------------------|
      |1   |s2i-go-1167-ubi8-run |successful |no                     |

## S2I go pipelinerun go-1.17-ubi9 version: PIPELINES-14-TC07
Tags: e2e, integration, clustertasks, non-admin, s2i
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                                  |
      |----|--------------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelines/s2i-go.yaml            |
      |2   |testdata/v1beta1/clustertask/pvc/pvc.yaml                     |
      |3   |testdata/v1beta1/clustertask/pipelineruns/s2i-go-117-ubi9.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name   |status     |check_lable_propagation|
      |----|--------------------|-----------|-----------------------|
      |1   |s2i-go-117-ubi9-run |successful |no                     |

## S2I java pipelinerun openjdk-11-el7 version: PIPELINES-14-TC08
Tags: e2e, integration, clustertasks, non-admin, s2i
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                                          |
      |----|----------------------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelines/s2i-java.yaml                  |
      |2   |testdata/v1beta1/clustertask/pvc/pvc.yaml                             |
      |3   |testdata/v1beta1/clustertask/pipelineruns/s2i-java-openjdk-11-el7.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name           |status     |check_lable_propagation|
      |----|----------------------------|-----------|-----------------------|
      |1   |s2i-java-openjdk-11-el7-run |successful |no                     |

## S2I java pipelinerun openjdk-11-ubi8 version: PIPELINES-14-TC09
Tags: e2e, integration, clustertasks, non-admin, s2i
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                                           |
      |----|-----------------------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelines/s2i-java.yaml                   |
      |2   |testdata/v1beta1/clustertask/pvc/pvc.yaml                              |
      |3   |testdata/v1beta1/clustertask/pipelineruns/s2i-java-openjdk-11-ubi8.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name            |status     |check_lable_propagation|
      |----|-----------------------------|-----------|-----------------------|
      |1   |s2i-java-openjdk-11-ubi8-run |successful |no                     |

## S2I java pipelinerun openjdk-17-ubi8 version: PIPELINES-14-TC10
Tags: e2e, integration, clustertasks, non-admin, s2i
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                                           |
      |----|-----------------------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelines/s2i-java.yaml                   |
      |2   |testdata/v1beta1/clustertask/pvc/pvc.yaml                              |
      |3   |testdata/v1beta1/clustertask/pipelineruns/s2i-java-openjdk-17-ubi8.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name            |status     |check_lable_propagation|
      |----|-----------------------------|-----------|-----------------------|
      |1   |s2i-java-openjdk-17-ubi8-run |successful |no                     |

## S2I java pipelinerun openjdk-8-ubi8 version: PIPELINES-14-TC11
Tags: e2e, integration, clustertasks, non-admin, s2i
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                                          |
      |----|----------------------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelines/s2i-java.yaml                  |
      |2   |testdata/v1beta1/clustertask/pvc/pvc.yaml                             |
      |3   |testdata/v1beta1/clustertask/pipelineruns/s2i-java-openjdk-8-ubi8.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name          |status      |check_lable_propagation|
      |----|---------------------------|------------|-----------------------|
      |1   |s2i-java-openjdk-8-ubi8-run |successful |no                     |

## S2I nodejs pipelinerun nodejs-14-ubi7 version: PIPELINES-14-TC12
Tags: e2e, integration, clustertasks, non-admin, s2i
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                                     |
      |----|-----------------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelines/s2i-nodejs.yaml           |
      |2   |testdata/v1beta1/clustertask/pvc/pvc.yaml                        |
      |3   |testdata/v1beta1/clustertask/pipelineruns/s2i-nodejs-14-ubi7.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name     |status     |check_lable_propagation|
      |----|----------------------|-----------|-----------------------|
      |1   |s2i-nodejs-14-ubi7-run|successful |no                     |

## S2I nodejs pipelinerun nodejs-14-ubi8-minimal version: PIPELINES-14-TC13
Tags: e2e, integration, clustertasks, non-admin, s2i
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                                             |
      |----|-------------------------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelines/s2i-nodejs.yaml                   |
      |2   |testdata/v1beta1/clustertask/pvc/pvc.yaml                                |
      |3   |testdata/v1beta1/clustertask/pipelineruns/s2i-nodejs-14-ubi8-minimal.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name             |status     |check_lable_propagation|
      |----|------------------------------|-----------|-----------------------|
      |1   |s2i-nodejs-14-ubi8-minimal-run|successful |no                     |

## S2I nodejs pipelinerun nodejs-14-ubi8 version: PIPELINES-14-TC14
Tags: e2e, integration, clustertasks, non-admin, s2i
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                                     |
      |----|-----------------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelines/s2i-nodejs.yaml           |
      |2   |testdata/v1beta1/clustertask/pvc/pvc.yaml                        |
      |3   |testdata/v1beta1/clustertask/pipelineruns/s2i-nodejs-14-ubi8.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name     |status     |check_lable_propagation|
      |----|----------------------|-----------|-----------------------|
      |1   |s2i-nodejs-14-ubi8-run|successful |no                     |

## S2I nodejs pipelinerun nodejs-16-ubi8 version: PIPELINES-14-TC15
Tags: e2e, integration, clustertasks, non-admin, s2i
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                                     |
      |----|-----------------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelines/s2i-nodejs.yaml           |
      |2   |testdata/v1beta1/clustertask/pvc/pvc.yaml                        |
      |3   |testdata/v1beta1/clustertask/pipelineruns/s2i-nodejs-16-ubi8.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name      |status    |check_lable_propagation|
      |----|-----------------------|----------|-----------------------|
      |1   |s2i-nodejs-16-ubi8-run|successful |no                     |

## S2I nodejs pipelinerun nodejs-16-ubi9 version: PIPELINES-14-TC16
Tags: e2e, integration, clustertasks, non-admin, s2i
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                                     |
      |----|-----------------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelines/s2i-nodejs.yaml           |
      |2   |testdata/v1beta1/clustertask/pvc/pvc.yaml                        |
      |3   |testdata/v1beta1/clustertask/pipelineruns/s2i-nodejs-16-ubi8.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name      |status    |check_lable_propagation|
      |----|-----------------------|----------|-----------------------|
      |1   |s2i-nodejs-16-ubi9-run|successful |no                     |
## S2I perl pipelinerun perl-526-ubi8 version: PIPELINES-14-TC17
Tags: e2e, integration, clustertasks, non-admin, s2i
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                                    |
      |----|----------------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelines/s2i-perl.yaml            |
      |2   |testdata/v1beta1/clustertask/pvc/pvc.yaml                       |
      |3   |testdata/v1beta1/clustertask/pipelineruns/s2i-perl-526-ubi8.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name    |status     |check_lable_propagation|
      |----|---------------------|-----------|-----------------------|
      |1   |s2i-perl-526-ubi8-run|successful |no                     |

## S2I perl pipelinerun perl-530-ubi8 version: PIPELINES-14-TC18
Tags: e2e, integration, clustertasks, non-admin, s2i
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                                    |
      |----|----------------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelines/s2i-perl.yaml            |
      |2   |testdata/v1beta1/clustertask/pvc/pvc.yaml                       |
      |3   |testdata/v1beta1/clustertask/pipelineruns/s2i-perl-530-ubi8.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name    |status     |check_lable_propagation|
      |----|---------------------|-----------|-----------------------|
      |1   |s2i-perl-530-ubi8-run|successful |no                     |

## S2I perl pipelinerun perl-530 version: PIPELINES-14-TC19
Tags: e2e, integration, clustertasks, non-admin, s2i
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                               |
      |----|-----------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelines/s2i-perl.yaml       |
      |2   |testdata/v1beta1/clustertask/pvc/pvc.yaml                  |
      |3   |testdata/v1beta1/clustertask/pipelineruns/s2i-perl-530.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name|status     |check_lable_propagation|
      |----|-----------------|-----------|-----------------------|
      |1   |s2i-perl-530-run |successful |no                     |

## S2I perl pipelinerun perl-532-ubi9 version: PIPELINES-14-TC20
Tags: e2e, integration, clustertasks, non-admin, s2i
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                                    |
      |----|----------------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelines/s2i-perl.yaml            |
      |2   |testdata/v1beta1/clustertask/pvc/pvc.yaml                       |
      |3   |testdata/v1beta1/clustertask/pipelineruns/s2i-perl-532-ubi9.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name    |status     |check_lable_propagation|
      |----|---------------------|-----------|-----------------------|
      |1   |s2i-perl-532-ubi9-run|successful |no                     |

## S2I perl pipelinerun python-27-ubi8 version: PIPELINES-14-TC21
Tags: e2e, integration, clustertasks, non-admin, s2i
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                                     |
      |----|-----------------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelines/s2i-python.yaml           |
      |2   |testdata/v1beta1/clustertask/pvc/pvc.yaml                        |
      |3   |testdata/v1beta1/clustertask/pipelineruns/s2i-python-27-ubi8.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name     |status     |check_lable_propagation|
      |----|----------------------|-----------|-----------------------|
      |1   |s2i-python-27-ubi8-run|successful |no                     |

## S2I perl pipelinerun python-36-ubi8 version: PIPELINES-14-TC22
Tags: e2e, integration, clustertasks, non-admin, s2i
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                                     |
      |----|-----------------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelines/s2i-python.yaml           |
      |2   |testdata/v1beta1/clustertask/pvc/pvc.yaml                        |
      |3   |testdata/v1beta1/clustertask/pipelineruns/s2i-python-36-ubi8.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name     |status     |check_lable_propagation|
      |----|----------------------|-----------|-----------------------|
      |1   |s2i-python-36-ubi8-run|successful |no                     |

## S2I perl pipelinerun python-38-ubi7 version: PIPELINES-14-TC23
Tags: e2e, integration, clustertasks, non-admin, s2i
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                                     |
      |----|-----------------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelines/s2i-python.yaml           |
      |2   |testdata/v1beta1/clustertask/pvc/pvc.yaml                        |
      |3   |testdata/v1beta1/clustertask/pipelineruns/s2i-python-38-ubi7.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name     |status     |check_lable_propagation|
      |----|----------------------|-----------|-----------------------|
      |1   |s2i-python-38-ubi7-run|successful |no                     |

## S2I perl pipelinerun python-38-ubi8 version: PIPELINES-14-TC24
Tags: e2e, integration, clustertasks, non-admin, s2i
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                                     |
      |----|-----------------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelines/s2i-python.yaml           |
      |2   |testdata/v1beta1/clustertask/pvc/pvc.yaml                        |
      |3   |testdata/v1beta1/clustertask/pipelineruns/s2i-python-38-ubi8.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name     |status     |check_lable_propagation|
      |----|----------------------|-----------|-----------------------|
      |1   |s2i-python-38-ubi8-run|successful |no                     |

## S2I perl pipelinerun python-38 version: PIPELINES-14-TC25
Tags: e2e, integration, clustertasks, non-admin, s2i
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                                |
      |----|------------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelines/s2i-python.yaml      |
      |2   |testdata/v1beta1/clustertask/pvc/pvc.yaml                   |
      |3   |testdata/v1beta1/clustertask/pipelineruns/s2i-python-38.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name|status     |check_lable_propagation|
      |----|-----------------|-----------|-----------------------|
      |1   |s2i-python-38-run|successful |no                     |

## S2I perl pipelinerun python-39-ubi8 version: PIPELINES-14-TC26
Tags: e2e, integration, clustertasks, non-admin, s2i
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                                     |
      |----|-----------------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelines/s2i-python.yaml           |
      |2   |testdata/v1beta1/clustertask/pvc/pvc.yaml                        |
      |3   |testdata/v1beta1/clustertask/pipelineruns/s2i-python-39-ubi8.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name     |status     |check_lable_propagation|
      |----|----------------------|-----------|-----------------------|
      |1   |s2i-python-39-ubi8-run|successful |no                     |

## S2I perl pipelinerun python-39-ubi9 version: PIPELINES-14-TC27
Tags: e2e, integration, clustertasks, non-admin, s2i
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                                     |
      |----|-----------------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelines/s2i-python.yaml           |
      |2   |testdata/v1beta1/clustertask/pvc/pvc.yaml                        |
      |3   |testdata/v1beta1/clustertask/pipelineruns/s2i-python-39-ubi9.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name     |status     |check_lable_propagation|
      |----|----------------------|-----------|-----------------------|
      |1   |s2i-python-39-ubi9-run|successful |no                     |

## S2I perl pipelinerun php-73-ubi7 version: PIPELINES-14-TC28
Tags: e2e, integration, clustertasks, non-admin, s2i
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                                  |
      |----|--------------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelines/s2i-php.yaml           |
      |2   |testdata/v1beta1/clustertask/pvc/pvc.yaml                     |
      |3   |testdata/v1beta1/clustertask/pipelineruns/s2i-php-73-ubi7.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name  |status     |check_lable_propagation|
      |----|-------------------|-----------|-----------------------|
      |1   |s2i-php-73-ubi7-run|successful |no                     |

## S2I perl pipelinerun php-73 version: PIPELINES-14-TC29
Tags: e2e, integration, clustertasks, non-admin, s2i
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                             |
      |----|---------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelines/s2i-php.yaml      |
      |2   |testdata/v1beta1/clustertask/pvc/pvc.yaml                |
      |3   |testdata/v1beta1/clustertask/pipelineruns/s2i-php-73.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name|status     |check_lable_propagation|
      |----|-----------------|-----------|-----------------------|
      |1   |s2i-php-73-run   |successful |no                     |

## S2I perl pipelinerun php-74-ubi8 version: PIPELINES-14-TC30
Tags: e2e, integration, clustertasks, non-admin, s2i
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                                  |
      |----|--------------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelines/s2i-php.yaml           |
      |2   |testdata/v1beta1/clustertask/pvc/pvc.yaml                     |
      |3   |testdata/v1beta1/clustertask/pipelineruns/s2i-php-74-ubi8.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name  |status     |check_lable_propagation|
      |----|-------------------|-----------|-----------------------|
      |1   |s2i-php-74-ubi8-run|successful |no                     |

## S2I perl pipelinerun php-80-ubi9 version: PIPELINES-14-TC31
Tags: e2e, integration, clustertasks, non-admin, s2i
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                                  |
      |----|--------------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelines/s2i-php.yaml           |
      |2   |testdata/v1beta1/clustertask/pvc/pvc.yaml                     |
      |3   |testdata/v1beta1/clustertask/pipelineruns/s2i-php-80-ubi9.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name  |status     |check_lable_propagation|
      |----|-------------------|-----------|-----------------------|
      |1   |s2i-php-80-ubi9-run|successful |no                     |

## S2I perl pipelinerun ruby-25-ubi8 version: PIPELINES-14-TC32
Tags: e2e, integration, clustertasks, non-admin, s2i
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                                   |
      |----|---------------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelines/s2i-ruby.yaml           |
      |2   |testdata/v1beta1/clustertask/pvc/pvc.yaml                      |
      |3   |testdata/v1beta1/clustertask/pipelineruns/s2i-ruby-25-ubi8.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name   |status     |check_lable_propagation|
      |----|--------------------|-----------|-----------------------|
      |1   |s2i-ruby-25-ubi8-run|successful |no                     |

## S2I perl pipelinerun ruby-26 version: PIPELINES-14-TC33
Tags: e2e, integration, clustertasks, non-admin, s2i, latest
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                              |
      |----|----------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelines/s2i-ruby.yaml      |
      |2   |testdata/v1beta1/clustertask/pvc/pvc.yaml                 |
      |3   |testdata/v1beta1/clustertask/pipelineruns/s2i-ruby-26.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name   |status     |check_lable_propagation|
      |----|--------------------|-----------|-----------------------|
      |1   |s2i-ruby-26-run     |successful |no                     |

## S2I perl pipelinerun ruby-27-ubi7 version: PIPELINES-14-TC34
Tags: e2e, integration, clustertasks, non-admin, s2i
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                                   |
      |----|---------------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelines/s2i-ruby.yaml           |
      |2   |testdata/v1beta1/clustertask/pvc/pvc.yaml                      |
      |3   |testdata/v1beta1/clustertask/pipelineruns/s2i-ruby-27-ubi7.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name   |status     |check_lable_propagation|
      |----|--------------------|-----------|-----------------------|
      |1   |s2i-ruby-27-ubi7-run|successful |no                     |

## S2I perl pipelinerun ruby-27-ubi8 version: PIPELINES-14-TC35
Tags: e2e, integration, clustertasks, non-admin, s2i
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                                   |
      |----|---------------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelines/s2i-ruby.yaml           |
      |2   |testdata/v1beta1/clustertask/pvc/pvc.yaml                      |
      |3   |testdata/v1beta1/clustertask/pipelineruns/s2i-ruby-27-ubi8.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name   |status     |check_lable_propagation|
      |----|--------------------|-----------|-----------------------|
      |1   |s2i-ruby-27-ubi8-run|successful |no                     |

## S2I perl pipelinerun ruby-30-ubi7 version: PIPELINES-14-TC36
Tags: e2e, integration, clustertasks, non-admin, s2i
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                              |
      |----|----------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelines/s2i-ruby.yaml      |
      |2   |testdata/v1beta1/clustertask/pvc/pvc.yaml                 |
      |3   |testdata/v1beta1/clustertask/pipelineruns/s2i-ruby-27.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name   |status     |check_lable_propagation|
      |----|--------------------|-----------|-----------------------|
      |1   |s2i-ruby-27-run     |successful |no                     |

## S2I perl pipelinerun ruby-30-ubi7 version: PIPELINES-14-TC37
Tags: e2e, integration, clustertasks, non-admin, s2i
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                                   |
      |----|---------------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelines/s2i-ruby.yaml           |
      |2   |testdata/v1beta1/clustertask/pvc/pvc.yaml                      |
      |3   |testdata/v1beta1/clustertask/pipelineruns/s2i-ruby-30-ubi7.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name   |status     |check_lable_propagation|
      |----|--------------------|-----------|-----------------------|
      |1   |s2i-ruby-30-ubi7-run|successful |no                     |

## S2I perl pipelinerun ruby-30-ubi8 version: PIPELINES-14-TC38
Tags: e2e, integration, clustertasks, non-admin, s2i
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                                   |
      |----|---------------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelines/s2i-ruby.yaml           |
      |2   |testdata/v1beta1/clustertask/pvc/pvc.yaml                      |
      |3   |testdata/v1beta1/clustertask/pipelineruns/s2i-ruby-30-ubi8.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name   |status     |check_lable_propagation|
      |----|--------------------|-----------|-----------------------|
      |1   |s2i-ruby-30-ubi8-run|successful |no                     |

## S2I perl pipelinerun ruby-30-ubi9 version: PIPELINES-14-TC39
Tags: e2e, integration, clustertasks, non-admin, s2i
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                                   |
      |----|---------------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelines/s2i-ruby.yaml           |
      |2   |testdata/v1beta1/clustertask/pvc/pvc.yaml                      |
      |3   |testdata/v1beta1/clustertask/pipelineruns/s2i-ruby-30-ubi9.yaml|
  * Verify pipelinerun
      |S.NO|pipeline_run_name   |status     |check_lable_propagation|
      |----|--------------------|-----------|-----------------------|
      |1   |s2i-ruby-30-ubi9-run|successful |no                     |