# Olm Openshift Pipelines operator specs

## Install openshift-pipelines operator
Tags: olm, install
Installs `openshift-pipelines` operator using olm

Steps:
  * Subscribe to operator
  * Wait for Cluster CR availability
  * Validate Operator should be installed
  * Validate pipelines deployment
  * Validate triggers deployment

## Upgrade openshift-pipelines operator
Tags: olm, upgrade
Installs `openshift-pipelines` operator using olm

Steps:
  * Upgrade operator subscription
  * Wait for Cluster CR availability
  * Validate Operator should be installed
  * Validate pipelines deployment
  * Validate triggers deployment

## Uninstall openshift-pipelines operator
Tags: olm, uninstall

Uninstalls `openshift-pipelines` operator using olm
Steps:
  * Validate Operator should be installed
  * Uninstall Operator