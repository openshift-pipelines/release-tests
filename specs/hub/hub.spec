PIPELINES-21
# HUB tests tests

## Install HUB without authentication: PIPELINES-21-TC01
Tags: hub, sanity, to-do
Component: HUB
Level: Integration
Type: Functional
Importance: Critical

This scenario tests HUB installation without authentication

Steps:
  * Verify ServiceAccount "pipeline" exist
  * Create TektinHub CR
  * Verify that the hub CR is in True state
  * Verify that the TektonHub deployment is up and running
  * Verify that the HUB UI is accessbile
  * Verify that the TektonHub elements like Kind, Platform, Catalog, Category are available
  Verify that the HUB UI contains tasks which are present in "https://github.com/manojbison/catalog" 
  Update task in "https://github.com/manojbison/catalog" repo
  Verify the updates are visible in HUB UI