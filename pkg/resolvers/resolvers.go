package resolvers

import (
	"strings"

	"github.com/openshift-pipelines/release-tests/pkg/cmd"
)

// This function checks if the project we need exists: if the project does not exist, then the function creates it
func CheckProjectExists(project string){
	availableProjects := cmd.MustSucceed("oc", "projects").Stdout()
	splittedAvailableProjects := strings.Fields(availableProjects)
	for i, _ := range splittedAvailableProjects{
		if strings.Contains(splittedAvailableProjects[i], project){
			cmd.MustSucceed("oc", "project", project)
			return
		}
	}
	cmd.MustSucceed("oc", "new-project", project)
}