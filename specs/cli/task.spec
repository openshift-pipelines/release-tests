# Verify tasks

Pre condition:
  * Operator should be installed

## Create task
Tags: tasks

Steps:
  * Run "tkn task create -f ../../testdata/tasks/shell-script-task.yaml"
  * Verify taks creation status is "successfull"


## Create duplicate task
Tags: tasks

Steps:
  * Run "tkn task create -f ../../testdata/tasks/shell-script-task.yaml"
  * Verify taks creation status is "successfull"
  * Create same task one more time with command "tkn task create -f ../../testdata/tasks/shell-script-task.yaml"
  * Verify failure message


## Create task using invalid syntax yaml file
Tags: tasks

Steps:
  * Run "tkn task create -f ../../testdata/tasks/invalid-task.yaml"
  * Verify failure message


## Delete task interactive
Tags: tasks

Steps:
  * Run "tkn task create -f ../../testdata/tasks/shell-script-task.yaml"
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

