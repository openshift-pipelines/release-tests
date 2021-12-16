PIPELINES-13

# Verify HPA (Horizontal Pod Autosclaer)

Pre condition:
  * Validate Operator should be installed

## Test HPA for tekton-pipelines-webhook deployment: PIPELINES-13-TC01
Tags: hpa, admin, to-do
Component: Operator
Level: Integration
Type: Functional
Importance: Critical
CustomerScenario: yes

This scenario tests HPA for tekton-pipelines-webhook deployment

Steps:
  * Run "kubectl -n openshift-pipelines scale --replicas=3 deployment/tekton-pipelines-webhook"
  * Sleep for "30" seconds
  * Assert if "3" pods related to "tekton-pipelines-webhook" are present and running in "openshift-pipelinse" namespace