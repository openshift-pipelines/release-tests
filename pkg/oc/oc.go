package oc

import (
	"log"

	"github.com/openshift-pipelines/release-tests/pkg/cmd"
	resource "github.com/openshift-pipelines/release-tests/pkg/config"
)

// Create resources using oc command
func Create(path_dir, namespace string) {
	log.Printf("output: %s\n", cmd.MustSucceed("oc", "create", "-f", resource.Path(path_dir), "-n", namespace).Stdout())
}

// Delete resources using oc command
func Delete(path_dir, namespace string) {
	log.Printf("output: %s\n", cmd.MustSucceed("oc", "delete", "-f", resource.Path(path_dir), "-n", namespace).Stdout())
}

//CreateNewProject Helps you to create new project
func CreateNewProject(ns string) {
	log.Printf("output: %s\n", cmd.MustSucceed("oc", "new-project", ns).Stdout())
}

//DeleteProject Helps you to delete new project
func DeleteProject(ns string) {
	log.Printf("output: %s\n", cmd.MustSucceed("oc", "delete", "project", ns).Stdout())
}

func LinkSecretToPipelinesSA(secretname, namespace string) {
	log.Printf("output: %s\n", cmd.MustSucceed("oc", "secret", "link", "serviceaccount/pipeline", "secrets/"+secretname, "-n", namespace).Stdout())
}
