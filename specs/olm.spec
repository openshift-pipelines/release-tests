PIPELINES-09
# Olm Openshift Pipelines operator specs

## Install openshift-pipelines operator: PIPELINES-09-TC01
Tags: install, admin, sanity
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

Installs `openshift-pipelines` operator using olm

Steps:
  * Subscribe to operator
  * Wait for TektonConfig CR availability
  * Validate pipelines deployment
  * Validate triggers deployment
  * Validate PAC deployment
  * Validate tkn server cli deployment
  * Verify TektonAddons Install status
  * Validate RBAC
  * Validate quickstarts
  * Validate default auto prune cronjob in target namespace
  * Validate tektoninstallersets status
  * Validate tektoninstallersets names
  * Change enable-api-fields to "beta"

## Upgrade openshift-pipelines operator: PIPELINES-09-TC02
Tags: upgrade, admin
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

Installs `openshift-pipelines` operator using olm

Steps:
  * Upgrade operator subscription
  * Wait for TektonConfig CR availability
  * Validate Operator should be installed
  * Validate RBAC
  * Validate quickstarts

## Uninstall openshift-pipelines operator: PIPELINES-09-TC03
Tags: uninstall, admin
Component: Operator
Level: Integration
Type: Functional
Importance: Critical

Uninstalls `openshift-pipelines` operator using olm
Steps:
  * Validate Operator should be installed
  * Uninstall Operator
