# Install Openshift Pipelines

## Install openshift-pipelines
Tags: e2e, integration, install

Installs `openshift-pipelines` operator using olm

Steps:
  * Wait for Cluster CR availability
  * Validate SCC
  * Validate installation of pipelines "v0.9.2"
  * Validate installation of triggers "v0.1"
  * Validate operator setup status
