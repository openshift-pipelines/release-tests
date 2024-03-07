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
  * Define the tekton-hub-api variable
  * Verify namespace "openshift-pipelines" exists
  * Apply
    | S.NO | resource_dir                   |
    |------|--------------------------------|
    | 1    | testdata/hub/tektonhub.yaml    |
  * Create secrets for Tekton Results
  * Apply in namespace "openshift-pipelines"
    | S.NO | resource_dir                   |
    |------|--------------------------------|
    | 1    | testdata/pvc/tekton-logs.yaml  |
    | 2    | testdata/results/result.yaml   | 
  * Create Results route
  * Configure GitHub token for git resolver in TektonConfig
  * Configure the bundles resolver
  * Enable console plugin
  * Validate pipelines deployment
  * Validate triggers deployment
  * Validate PAC deployment
  * Validate chains deployment
  * Validate hub deployment
  * Validate tkn server cli deployment
  * Validate console plugin deployment
  * Ensure that Tekton Results is ready
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
  * Uninstall Operator