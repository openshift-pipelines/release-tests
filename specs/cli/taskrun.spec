# Verify tasks

Pre condition:
  * Operator should be installed

## Cancel taskrun
Tags: taskruns

Steps:
  * Run "tkn task create -f ../../testdata/tasks/python-script-task.yaml"
  * Start the task with command "tkn task start python-script-task"
  * List the task and get the taskrun name
  * Cancel the taskrun using command "tkn taskrun cancel <taskrun name>"
  * Verify the output message
  * Verify that the taskrun status is "Cancelled(TaskRunCancelled)"


## Cancel taskrun which does not exist
Tags: taskruns

Steps:
  * Run "tkn taskrun cancel which-doesnt-exist"
  * Verify that the error message is "Error: failed to find taskrun: which-doesnt-exist"


## tkn cancel help
Tags: taskruns

Steps:
  * Run "tkn taskrun cancel --help"
  * Verify help to run tkn taskrun cancel is shown


## Delete a taskrun force
Tags: taskruns

Steps:
  * Run "tkn task create -f ../../testdata/tasks/python-script-task.yaml"
  * Start the task with command "tkn task start python-script-task"
  * List the task and get the taskrun name
  * Delete the taskrun using command "tkn taskrun delete <taskrun name> -f"


## Delete a taskrun interactive
Tags: taskruns

Steps:
  * Run "tkn task create -f ../../testdata/tasks/python-script-task.yaml"
  * Start the task with command "tkn task start python-script-task"
  * List the task and get the taskrun name
  * Delete the taskrun using command "tkn taskrun delete <taskrun name>"
  * Check if it asks your permission to delete taskrun
  * Press Y
  * Verify if the taskrun is deleted


## Delete all taskrun with --all option
Tags: taskruns

Steps:
  * Run "tkn task create -f ../../testdata/tasks/python-script-task.yaml"
  * Start the task several times with command "tkn task start python-script-task"
  * Delete the taskrun using command "tkn taskrun delete --all -f"
  * Verify that all taskrun related to task created in step 1 is deleted


## Delete taskruns using --keep option
Tags: taskruns

Steps:
  * Run "tkn task create -f ../../testdata/tasks/python-script-task.yaml"
  * Start the task multiple times with command "tkn task start python-script-task"
  * Delete the taskrun using command "tkn taskrun delete --all -f 2"
  * Verify that all taskrun related to task created in step 1 except latest 2 taskruns


## Delete all taskruns of a particular task
Tags: taskruns

Steps:
  * Run "tkn task create -f ../../testdata/tasks/python-script-task.yaml"
  * Start the task several times with command "tkn task start python-script-task"
  * Delete the taskrun of task python-script-task using command "tkn taskrun delete -t python-script-task -f"
  * Verify that all taskrun related to python-script-task is deleted


## List all taskruns
Tags: taskruns

Steps:
  * Run "tkn task create -f ../../testdata/tasks/python-script-task.yaml"
  * Start the task multiple times with command "tkn task start python-script-task"
  * List the taskruns using the command "tkn taskrun list"
  * Verify that all the taskruns are list


## List taskruns of particular task
Tags: taskruns

Steps:
  * Run "tkn task create -f ../../testdata/tasks/python-script-task.yaml"
  * Start the task multiple times with command "tkn task start python-script-task"
  * List the taskruns using the command "tkn taskrun list python-script-task"
  * Verify that all the taskruns are list


## List taskruns with limit
Tags: taskruns

Steps:
  * Run "tkn task create -f ../../testdata/tasks/python-script-task.yaml"
  * Start the task multiple times with command "tkn task start python-script-task"
  * List the taskruns using the command "tkn taskrun list --limit 2"
  * Verify that 2 latest taskruns are listed


## Get logs of the taskrun
Tags: taskruns

Steps:
  * Run "tkn task create -f ../../testdata/tasks/python-script-task.yaml"
  * Start the task with command "tkn task start python-script-task"
  * list the taskrun and get the taskurn name
  * Wait for the task to start
  * Get the logs of the taskrun using the command "tkn taskrun logs <taskrun name>"
  * Verify the taskrun logs


## Get live logs of taskrun
Tags: taskruns

Steps:
  * Run "tkn task create -f ../../testdata/tasks/python-script-task.yaml"
  * Start the task with command "tkn task start python-script-task"
  * list the taskrun and get the taskurn name
  * Get the logs of the taskrun using the command "tkn taskrun logs <taskrun name> -f"
  * Verify the logs


## Get help for tkn taskrun logs command
Tags: taskruns

Steps:
  * Run "tkn taskrun logs --help" command
  * Verify the help output


## Get the logs for latest taskrun
Tags: taskruns

Steps:
  * Run "tkn task create -f ../../testdata/tasks/python-script-task.yaml"
  * Start the task multiple times with the command "tkn task start python-script-task"
  * Get the logs of the latest taskrun using the command "tkn taskrun logs --last"
  * Verify the logs


## List the latest taskruns to select before showing logs
Tags: taskruns

Steps:
  * Run "tkn task create -f ../../testdata/tasks/python-script-task.yaml"
  * Start the task multiple times with the command "tkn task start python-script-task"
  * Get the list of taskruns to select by running command "tkn taskrun logs --limit 3"
  * Verify latest 3 taskruns are listed
  * Select any one to get the logs
  * Verify the logs


## Get logs for specified step
Tags: taskruns

Steps:
  * Run "tkn task create -f ../../testdata/tasks/python-script-task.yaml"
  * Start the task multiple times with the command "tkn task start python-script-task"
  * list the taskruns using the command "tkn taskrun ls"
  * Get the logs for the particular step using command "tkn taskrun logs -s python-example-2 <taskrun name>"
  * Verify the logs


## Get description for taskrun
Tags: taskruns

Steps:
  * Run "tkn task create -f ../../testdata/tasks/python-script-task.yaml"
  * Start the task multiple times with the command "tkn task start python-script-task"
  * List the taskrun using the command"tkn taskrun ls"
  * Get the description for the taskrun by running the command "tkn taskrun describe <taskrun name>"
  * Verify the description