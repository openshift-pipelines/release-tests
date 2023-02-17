PIPELINES-20
# Pipelines As Code tests

## Enable/Disable PAC: PIPELINES-20-TC01
Tags: pac, sanity, to-do
Component: PAC
Level: Integration
Type: Functional
Importance: Critical

This scenario tests enable/disable of pipelines as code from tektonconfig custom resource

Steps:
  * Set "enable" section under "pipelinesAsCode" to "false"
  * Verify the installersets related to PAC are "not present"
  * Verify that the pods related to PAC are "not present" from "openshift-pipelines" namespace
  * Verify that the custom resource "pipelines-as-code" of type "pac" is removed
  * Set "enable" section under "pipelinesAsCode" to "true"
  * Verify the installersets related to PAC are "present"
  * Verify that the pods related to PAC are "present" from "openshift-pipelines" namespace
  * Verify that the custom resource "pipelines-as-code" of type "pac" is removed

## Enable/Disable PAC: PIPELINES-20-TC02
Tags: pac, sanity, to-do
Component: PAC
Level: Integration
Type: Functional
Importance: Critical

This scenario tests if application name change is visible in github UI

Steps:
  * Change "application-name" to "Pipelines as Code test" in tektonconfig config CR
  * Configure PAC using github app
  * Create a repo in github
  * Create a repo CRD for the above created repo
  * Configure pac pipelines for the repo for "push" event
  * Push the changes to github
  * Verify that a pipelinerun is created in the namespace where pac is configured for the repo
  * Verify that the pipelinerun status is updated in github
  * Verify that the application name is shown as "Pipelines as Code test" in github UI

## Enable/Disable auto-configure-new-github-repo: PIPELINES-20-TC03
Tags: pac, sanity, to-do
Component: PAC
Level: Integration
Type: Functional
Importance: Critical

This scenario tests if a new repo cr is getting created after enabling auto-configure-new-github-repo

Steps:
  * Set "auto-configure-new-github-repo" section under "pipelinesAsCode" to "true"
  * Configure PAC using github app
  * Create a new repo in github
  * Verify that a new repo cr got created
  * Set "auto-configure-new-github-repo" section under "pipelinesAsCode" to "false"
  * Create a new repo in github
  * Verify that repo cr is not created

## Enable/Disable auto-configure-new-github-repo: PIPELINES-20-TC03
Tags: pac, sanity, to-do
Component: PAC
Level: Integration
Type: Functional
Importance: Critical

This scenario tests if a new repo cr is getting created after enabling auto-configure-new-github-repo

Steps:
  * Set "auto-configure-new-github-repo" section under "pipelinesAsCode" to "true"
  * Configure PAC using github app
  * Create a new repo in github
  * Verify that a new repo cr got created
  * Set "auto-configure-new-github-repo" section under "pipelinesAsCode" to "false"
  * Create a new repo in github
  * Verify that repo cr is not created

## Enable/Disable error-log-snippet: PIPELINES-20-TC04
Tags: pac, sanity, to-do
Component: PAC
Level: Integration
Type: Functional
Importance: Critical

This scenario tests is error log snippet is shown in github UI when error-log-snippet set to true

Steps:
  * Set "error-log-snippet" section under "pipelinesAsCode" to "false"
  * Configure PAC using github app
  * Create a repo in github
  * Create a repo CRD for the above created repo
  * Configure pac pipelines for the repo for "push" event
  * Create a pipelinerun which fails all the time
  * Create a push event for the github repo
  * Once the pipelinerun is completed verify that the error log message is not shown in the github UI
  * Set "error-log-snippet" section under "pipelinesAsCode" to "true"
  * Create a push event for the github repo
  * Once the pipelinerun is completed verify that the error log message is shown in the github UI
  
