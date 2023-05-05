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

## Check server side components versions: PIPELINES-09-TC04
Tags: install, upgrade
Component: Operator
Level: Integration
Type: Functional
Importance: Critical
Steps:
  * Check version of component "pipeline"
  * Check version of component "triggers"
  * Check version of component "operator"
  * Check version of component "pipelines-as-code"
  * Check version of OSP

## Check client versions: PIPELINES-09-TC05
Tags: install, upgrade
Component: Operator
Level: Integration
Type: Functional
Importance: Critical
Steps: 
  * Download and extract CLI from cluster
  * Check "tkn" client version
  * Check "tkn-pac" version
  * Check "opc" client version