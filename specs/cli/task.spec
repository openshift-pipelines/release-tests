# Verify tasks

Pre condition:
  * Operator should be installed

## Create task using cli
Tags: tasks

Steps:
  * Cretae task from tkn "../../testdata/tasks/shell-script-task.yaml"
  * Verify taks creation status "successfull"


## Create duplicate task using cli
Tags: tasks

Steps:
  * Create task from tkn "../../testdata/tasks/shell-script-task.yaml"
  * Verify taks creation status "successfull"
  * Create task from tkn "../../testdata/tasks/shell-script-task.yaml"
  * Verify task creation status "unsuccessfull"


## Create task using invalid syntax yaml file
Tags: tasks

Steps:
  * create task from tkn "../../testdata/tasks/invalid-task.yaml"
  * Verify task creation status "successfull"

## Delete task interactive
Tags: tasks

Steps:
  * Create task from tkn "../../testdata/tasks/shell-script-task.yaml"
  * Verify taks creation status is "successfull"
  * Run "tkn task delete shell-script-task"
  * Verify whether it asks permission whether to delete the task or not
  * Answer y
  * Verify whether task got deleted or not


## Delete task interactive negative
Tags: tasks

Steps:
  * Run "tkn task create -f ../../testdata/tasks/shell-script-task.yaml"
  * Verify taks creation status is "successfull"
  * Run "tkn task delete shell-script-task"
  * Verify whether it asks permission whether to delete the task or not
  * Answer n
  * Verify the cancel message


## Delete task by force
Tags: tasks

Steps:
  * Run "tkn task create -f ../../testdata/tasks/shell-script-task.yaml"
  * Verify taks creation status is "successfull"
  * Run "tkn task delete shell-script-task -f"
  * Verify whether task got deleted or not


## Delete task which does not exist
Tags: tasks

Steps:
  * Run "tkn task delete non-existing-task"
  * Verify the error message


## Delete task with taskrun (--all)

Steps:  
  * Run "tkn task create -f ../../testdata/tasks/shell-script-task.yaml"
  * Verify taks creation status is "successfull"
  * Start the task using command "tkn task start shell-script-task
  * Run "tkn task delete shell-script-task --all -f"
  * Verify deletion of both task and taskrun


## Delete help

Steps:
  * Run "tkn task delete --help"
  * Help to use tkn task delete should be shown