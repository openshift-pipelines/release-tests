# Verify tkn pipeline sub-command

Pre Condition:
  * Operator should be installed

## Create pipeline using tkn binary

Tags: e2e, integration, tkn

  * Create pipeline from "../testdata/pipeline.yaml"

## Creating invalid pipeline using tkn binary fails

Tags: e2e, integration, tkn, negative

(Negative Scenario) validate Error creation pipeline from file using `tkn` cli

  * Create pipeline from  "../testdata/pipeline.yaml" - In Non-existance namespace
  * Create pipeline from  "../testdata/pipeline.pdf"  - with unsupported file format
  * Create pipeline from  "../testdata/pipelinerun.yaml" - with mismatched Resource kind

## Start pipeline using Tkn
Tags: e2e, integration, tkn

Start pipeline interactively using `tkn` binary
  * Create sample pipeline
  * Start pipeline using tkn
