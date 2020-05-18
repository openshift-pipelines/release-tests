# Olm Openshift Pipelines spec

## Install openshift-pipelines
Tags: olm, install

Installs `openshift-pipelines` operator using olm

Steps:
  * Subscribe to operator
  * Wait for Cluster CR availability
  * Validate installation of pipelines "v0.11.3"
  * Validate installation of triggers
  * Validate operator setup status


## Uninstall openshift-pipelines
Tags: olm, uninstall

Uninstalls `openshift-pipelines` operator using olm
Steps:
  * Operator should be installed
  * Uninstall Operator
