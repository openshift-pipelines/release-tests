olm install Spec
===================

Specifications are defined via H1 tag, you could either use the syntax above or "# olm install Spec"

Any context step defined under spec section, gets executed before every scenario. Every unordered list is a step.
 (Nothing)

install openshift-pipelines operator 
------------------------------------
Tags: e2e, integration
 Installs `opesnshift-pipelines` operator using olm
 1. Waits for cluster config to be created.
 2. validates installation process.
 3. verifies the status of resources with right versions. 

* Wait for Cluster CR availability
* Validate SCC
* Validate installation of pipelines "v0.9.2"
* Validate installation of triggers "v0.1"
* Validate opeartor setup status

