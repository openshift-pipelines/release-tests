package resolvers

import (
	"strings"

	"github.com/openshift-pipelines/release-tests/pkg/cmd"
)

// This function checks if the project we need exists: if the project does not exist, then the function creates it
func CheckProjectExists(project string){
	commandResult := cmd.MustSucceed("oc", "projects").Stdout()
	splittedCommandResult := strings.Fields(commandResult)
	for i, _ := range splittedCommandResult{
		if strings.Contains(splittedCommandResult[i], project){
			cmd.MustSucceed("oc", "project", project)
			return
		}
	}
	cmd.MustSucceed("oc", "new-project", project)
}