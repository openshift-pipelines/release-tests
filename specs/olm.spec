# Olm Openshift Pipelines operator specs

## Install openshift-pipelines operator
Tags: install
Installs `openshift-pipelines` operator using olm

Steps:
  * Subscribe to operator
  * Wait for TektonConfig CR availability
  * Validate pipelines deployment
  * Validate triggers deployment
  * Verify TektonAddons Install status
  * Validate RBAC

## Upgrade openshift-pipelines operator
Tags: upgrade
Installs `openshift-pipelines` operator using olm

Steps:
  * Upgrade operator subscription
  * Wait for TektonConfig CR availability
  * Validate Operator should be installed
  * Validate RBAC

## Uninstall openshift-pipelines operator
Tags: uninstall

Uninstalls `openshift-pipelines` operator using olm
Steps:
  * Validate Operator should be installed
  * Uninstall Operator