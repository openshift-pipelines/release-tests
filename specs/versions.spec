PIPELINES-22
# Versions of OpenShift Pipelines upstream components and CLI

## Check server side components versions: PIPELINES-22-TC01
Tags: e2e, sanity
Component: Operator
Level: Integration
Type: Functional
Importance: High
Steps:
  * Check version of component "pipeline"
  * Check version of component "triggers"
  * Check version of component "operator"
  * Check version of component "pipelines-as-code" 
  * Check version of component "hub"
  * Check version of OSP
 
## Check client versions: PIPELINES-22-TC02
Tags: e2e, sanity
Component: Operator
Level: Integration
Type: Functional
Importance: High
Steps: 
  * Download and extract CLI from cluster
  * Check "tkn" client version
  * Check "tkn-pac" version
  * Check "opc" client version
