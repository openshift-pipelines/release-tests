PIPELINES-20
# Pipelines As Code tests

## Enable/Disable PAC: PIPELINES-20-TC01
Tags: pac, sanity, to-do
Component: PAC
Level: Integration
Type: Functional
Importance: Critical

This scenario tests enable/disable of pipelines as code from tektonconfig custom resource

Steps:
  * Set "enable" section under "pipelinesAsCode" to "false"
  * Verify the installersets related to PAC are "not present"
  * Verify that the pods related to PAC are "not present" from "openshift-pipelines" namespace
  * Verify that the custom resource "pipelines-as-code" of type "pac" is "not present"
  * Set "enable" section under "pipelinesAsCode" to "true"
  * Verify the installersets related to PAC are "present"
  * Verify that the pods related to PAC are "present" from "openshift-pipelines" namespace
  * Verify that the custom resource "pipelines-as-code" of type "pac" is "present"