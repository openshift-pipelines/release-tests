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
  * Check version of component "chains"
  * Check version of component "pac"
  * Check version of component "hub"
  * Check version of component "results"
  * Check version of component "manual-approval-gate"
  * Check version of OSP
 
## Check client versions: PIPELINES-22-TC02
Tags: sanity
Component: Operator
Level: Integration
Type: Functional
Importance: High
Steps: 
  * Download and extract CLI from cluster
  * Check "tkn" client version
  * Check "tkn-pac" version
  * Check "opc" client version
  * Check "opc" server version

## Check OSP Version in OlmSkipRange : PIPELINES-22-TC03
Tags: e2e, sanity, olm
Component: Operator
Level: Integration
Type: Functional
Importance: High
Steps: 
  * Validate OSP Version in OlmSkipRange