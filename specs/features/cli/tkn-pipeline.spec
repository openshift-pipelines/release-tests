Verify tkn pipeline features
============================

* Operator should be installed

Create pipeline using Tkn binary 
--------------------------------
Tags: e2e, integration, tkn
 Create a pipeline from file using Tkn cli
 
* Create pipeline from file "../testdata/pipeline.yaml"

Create pipeline using Tkn binary (NEGATIVE) 
-------------------------------------------
Tags: e2e, integration, tkn
 (Negative Scenario) validate Error creation pipeline from file using Tkn cli 

* Create pipeline from file "../testdata/pipeline.yaml" - In Non-existance namespace
* Create pipeline from file "../testdata/pipeline.pdf" - with unsupported file format
* Create pipeline from file "../testdata/pipelinerun.yaml" - with mismatched Resource kind

Start pipeline using Tkn
------------------------
Tags: e2e, integration, tkn
 Start pipeline intarcatively using Tkn binary
* Create sample pipeline
* Start pipleine using tkn