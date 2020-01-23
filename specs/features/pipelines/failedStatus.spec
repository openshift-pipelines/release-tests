verify pipeline failure status
==============================

Specifications are defined via H1 tag, you could either use the syntax above or "# Pipeline failed status Spec"

* Operator should be installed

Run Task with non-existance SA
------------------------------
Tags: e2e, pipeline
 Creates a simple Task
 Validate for failure status 
  when we try to run Task with `non-existance` SA 

* Create Task
* Run Task with "non-existance" SA
* Validate TaskRun for failed status

Run Pipeline with non-existance SA
----------------------------------
Tags: e2e, pipeline
 Creates a simple pipeline
 Validate for failure status
  when we try to run pipeline with `non-existance` SA 

* Create pipeline
* Run pipeline with "non-existance" SA
* Validate pipelineRun for failed status