PIPELINES-11
# Verify RBAC Resources

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
  * Update TektonConfig CR to use param with name createRbacResource and value "true" to "enable" auto creation of RBAC resources
  * Verify RBAC resources are auto created successfully
  * Update TektonConfig CR to use param with name createRbacResource and value "false" to "disable" auto creation of RBAC resources
  * Verify RBAC resources disabled successfully
  * Update TektonConfig CR to use param with name createRbacResource and value "true" to "enable" auto creation of RBAC resources
  * Verify RBAC resources are auto created successfully
