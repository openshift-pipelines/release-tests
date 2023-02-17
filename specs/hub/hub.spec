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
  * Create resource "testdata/hub/tektonhub.yaml"
  * Verify the the hub cr is in True state
  * Verify that the hub UI deployment is up and running
  * Verify that the hub api deployment is up and running
  * Verify that the hub db deployment is up and running
  * Verify that the HUB UI is accessbile
  * Verify that the HUB UI contains tasks which are present in "https://github.com/VeereshAradhya/catalog" 
  * Update task in "https://github.com/VeereshAradhya/catalog"  repo
  * Verify the updates are visible in HUB UI