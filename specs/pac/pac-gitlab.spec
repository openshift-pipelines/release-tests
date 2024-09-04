PIPELINES-29
# Pipelines As Code tests

## Enable/Disable PAC: PIPELINES-29-TC01
Tags: pac, sanity
Component: PAC
Level: Integration
Type: Functional
Importance: Critical

This scenario tests enable/disable of pipelines as code from tektonconfig custom resource

Steps:
  * Create
     | S.NO | resource_dir                          |
     |------|---------------------------------------|
     | 1    | testdata/pac/eventlistener.yaml       |
     | 2    | testdata/pac/trigger-binding.yaml     |
     | 3    | testdata/pac/trigger-template.yaml    |
     | 4    | testdata/pac/pipeline.yaml            |
  * Create Smee Deployment
  * Configure Gitlab repo
