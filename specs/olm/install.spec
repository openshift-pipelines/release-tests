# olm install Spec

## install openshift-pipelines
Tags: e2e, integration

Installs `opesnshift-pipelines` operator using olm

1. Waits for cluster config to be created.
2. validates installation process.
3. verifies the status of resources with right versions.

Steps:
  * Wait for Cluster CR availability
  * Validate SCC
  * Validate installation of pipelines "v0.9.2"
  * Validate installation of triggers "v0.1"
  * Validate operator setup status

