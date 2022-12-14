PIPELINES-14
# Verify Clustertasks E2E spec

Pre condition:
  * Validate Operator should be installed


## S2I nodejs pipelinerun: PIPELINES-14-TC01
Tags: e2e, clustertasks, non-admin, s2i, sanity, interop
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

## S2I dotnet pipelinerun: PIPELINES-14-TC02
Tags: e2e, clustertasks, non-admin, s2i, skip_linux/ppc64le
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                          |
      |----|------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelines/s2i-dotnet.yaml|
      |2   |testdata/v1beta1/clustertask/pvc/pvc.yaml             |
  * Get tags of the imagestream "dotnet" from namespace "openshift" and store to variable "dotnet-tags"
  * Start and verify dotnet pipeline "s2i-dotnet-pipeline" with values stored in variable "dotnet-tags" with workspace "name=source,claimName=shared-pvc"

## S2I golang pipelinerun: PIPELINES-14-TC03
Tags: e2e, clustertasks, non-admin, s2i
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                          |
      |----|------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelines/s2i-golang.yaml|
      |2   |testdata/v1beta1/clustertask/pvc/pvc.yaml             |
  * Get tags of the imagestream "golang" from namespace "openshift" and store to variable "golang-tags"
  * Start and verify pipeline "s2i-go-pipeline" with param "VERSION" with values stored in variable "golang-tags" with workspace "name=source,claimName=shared-pvc"

## S2I java pipelinerun: PIPELINES-14-TC04
Tags: e2e, clustertasks, non-admin, s2i, sanity, interop
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                          |
      |----|------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelines/s2i-java.yaml  |
      |2   |testdata/v1beta1/clustertask/pvc/pvc.yaml             |
  * Get tags of the imagestream "java" from namespace "openshift" and store to variable "java-tags"
  * Start and verify pipeline "s2i-java-pipeline" with param "VERSION" with values stored in variable "java-tags" with workspace "name=source,claimName=shared-pvc"

## S2I nodejs pipelinerun: PIPELINES-14-TC05
Tags: e2e, clustertasks, non-admin, s2i
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                          |
      |----|------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelines/s2i-nodejs.yaml|
      |2   |testdata/v1beta1/clustertask/pvc/pvc.yaml             |
  * Get tags of the imagestream "nodejs" from namespace "openshift" and store to variable "nodejs-tags"
  * Start and verify pipeline "s2i-nodejs-pipeline" with param "VERSION" with values stored in variable "nodejs-tags" with workspace "name=source,claimName=shared-pvc"

## S2I perl pipelinerun: PIPELINES-14-TC06
Tags: e2e, clustertasks, non-admin, s2i
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                          |
      |----|------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelines/s2i-perl.yaml  |
      |2   |testdata/v1beta1/clustertask/pvc/pvc.yaml             |
  * Get tags of the imagestream "perl" from namespace "openshift" and store to variable "perl-tags"
  * Start and verify pipeline "s2i-perl-pipeline" with param "VERSION" with values stored in variable "perl-tags" with workspace "name=source,claimName=shared-pvc"

## S2I php pipelinerun: PIPELINES-14-TC07
Tags: e2e, clustertasks, non-admin, s2i
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                       |
      |----|---------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelines/s2i-php.yaml|
      |2   |testdata/v1beta1/clustertask/pvc/pvc.yaml          |
  * Get tags of the imagestream "php" from namespace "openshift" and store to variable "php-tags"
  * Start and verify pipeline "s2i-php-pipeline" with param "VERSION" with values stored in variable "php-tags" with workspace "name=source,claimName=shared-pvc"

## S2I python pipelinerun: PIPELINES-14-TC08
Tags: e2e, clustertasks, non-admin, s2i
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                          |
      |----|------------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelines/s2i-python.yaml|
      |2   |testdata/v1beta1/clustertask/pvc/pvc.yaml             |
  * Get tags of the imagestream "python" from namespace "openshift" and store to variable "python-tags"
  * Start and verify pipeline "s2i-python-pipeline" with param "VERSION" with values stored in variable "python-tags" with workspace "name=source,claimName=shared-pvc"

## S2I ruby pipelinerun: PIPELINES-14-TC09
Tags: e2e, clustertasks, non-admin, s2i
Component: Pipelines
Level: Integration
Type: Functional
Importance: Critical

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create
      |S.NO|resource_dir                                        |
      |----|----------------------------------------------------|
      |1   |testdata/v1beta1/clustertask/pipelines/s2i-ruby.yaml|
      |2   |testdata/v1beta1/clustertask/pvc/pvc.yaml           |
  * Get tags of the imagestream "ruby" from namespace "openshift" and store to variable "ruby-tags"
  * Start and verify pipeline "s2i-ruby-pipeline" with param "VERSION" with values stored in variable "ruby-tags" with workspace "name=source,claimName=shared-pvc"