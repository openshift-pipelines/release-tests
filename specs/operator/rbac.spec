PIPELINES-11
# Verify RBAC Resources and CA Bundle Configuration

Pre condition:
  * Validate Operator should be installed

## Disable RBAC resource creation: PIPELINES-11-TC01
Tags: e2e, rbac-disable, admin, sanity
Component: Operator
Level: Integration
Type: Functional
Importance: High

This scenario helps you to disable creation of RBAC resources at cluster level.

Steps:
  * Update TektonConfig CR to use param with name "createRbacResource" and value "true" to "enable" auto creation of "RBAC resources"
  * Verify RBAC resources are auto created successfully
  * Update TektonConfig CR to use param with name "createRbacResource" and value "false" to "disable" auto creation of "RBAC resources"
  * Verify RBAC resources disabled successfully
  * Update TektonConfig CR to use param with name "createRbacResource" and value "true" to "enable" auto creation of "RBAC resources"
  * Verify RBAC resources are auto created successfully

## Independent CA Bundle ConfigMap creation control: PIPELINES-11-TC02
Tags: e2e, cabundle-control, admin, sanity
Component: Operator
Level: Integration
Type: Functional
Importance: High

This scenario helps you to enable CA Bundle ConfigMap creation at cluster level.

Steps:
  * Update TektonConfig CR to use param with name "createCABundleConfigMaps" and value "true" to "enable" auto creation of "CA Bundle ConfigMaps"
  * Verify CA Bundle ConfigMaps are auto created successfully
  * Update TektonConfig CR to use param with name "createCABundleConfigMaps" and value "false" to "disable" auto creation of "CA Bundle ConfigMaps"
  * Verify CA Bundle ConfigMaps still exist

Teardown:
  * Update TektonConfig CR to use param with name "createRbacResource" and value "true" to "enable" auto creation of "RBAC resources"
  * Update TektonConfig CR to use param with name "createCABundleConfigMaps" and value "true" to "enable" auto creation of "CA Bundle ConfigMaps"