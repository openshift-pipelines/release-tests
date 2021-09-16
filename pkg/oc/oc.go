package oc

import (
	"log"
	"os"
	"strings"

	"github.com/getgauge-contrib/gauge-go/testsuit"
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

func LinkSecretToSA(secretname, sa, namespace string) {
	log.Printf("output: %s\n", cmd.MustSucceed("oc", "secret", "link", "serviceaccount/"+sa, "secrets/"+secretname, "-n", namespace).Stdout())
}

func CreateSecretWithSecretToken(secretname, namespace string) {
	log.Printf("output: %s\n", cmd.MustSucceed("oc", "create", "secret", "generic", secretname, "--from-literal=secretToken="+os.Getenv("SECRET_TOKEN"), "-n", namespace).Stdout())
}

func EnableTLSConfigForEventlisteners(namespace string) {
	log.Printf("output: %s\n", cmd.MustSucceed("oc", "label", "namespace", namespace, "operator.tekton.dev/enable-annotation=enabled").Stdout())
}

func UpdateTektonConfig(patch_data string) {
	log.Printf("output: %s\n", cmd.MustSucceed("oc", "patch", "tektonconfig", "config", "-p", patch_data, "--type=merge"))
}

func VerifyCronjobStatus(cronJobName, status, namespace string) {
	res := cmd.MustSucceed("oc", "get", "cronjob", "-n", namespace).Stdout()
	if status == "present" {
		if !strings.Contains(res, cronJobName) {
			testsuit.T.Errorf("Error: Expected a cronjob with name %v is present in namespace %v", cronJobName, namespace)
		} else {
			log.Printf("cronjob with name %v is present in namespace %v", cronJobName, namespace)
		}
	} else if status == "not present" {
		if strings.Contains(res, cronJobName) {
			testsuit.T.Errorf("Error: Expected a cronjob with name %v is not present in namespace %v", cronJobName, namespace)
		} else {
			log.Printf("cronjob with name %v is not present in namespace %v", cronJobName, namespace)
		}
	} else {
		testsuit.T.Errorf("Invalid input for status: %v", status)
	}
}
