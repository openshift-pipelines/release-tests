package oc

import (
	"log"
	"os"
	"strings"

	"fmt"
	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/assert"
	"github.com/openshift-pipelines/release-tests/pkg/clients"
	"github.com/openshift-pipelines/release-tests/pkg/cmd"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	resource "github.com/openshift-pipelines/release-tests/pkg/config"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
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

func VerifyKubernetesEventsForEventListener(namespace string) {
	result := cmd.Run("oc", "-n", namespace, "get", "events")
	startedEvent := strings.Contains(result.String(), "dev.tekton.event.triggers.started.v1")
	successfulEvent := strings.Contains(result.String(), "dev.tekton.event.triggers.successful.v1")
	doneEvent := strings.Contains(result.String(), "dev.tekton.event.triggers.done.v1")
	if !startedEvent || !successfulEvent || !doneEvent {
		testsuit.T.Errorf("No events for successful, done and started")
	}
}

func UpdateTektonConfig(patch_data string) {
	log.Printf("output: %s\n", cmd.MustSucceed("oc", "patch", "tektonconfig", "config", "-p", patch_data, "--type=merge").Stdout())
}

func UpdateTektonConfigwithInvalidData(patch_data, errorMessage string) {
	result := cmd.Run("oc", "patch", "tektonconfig", "config", "-p", patch_data, "--type=merge")
	log.Printf("Output: %s\n", result.Stdout())
	if result.ExitCode != 1 {
		testsuit.T.Errorf("Expected exit code 1 but got %v", result.ExitCode)
	}
	if !strings.Contains(result.Stderr(), errorMessage) {
		testsuit.T.Errorf("Expected error message substring %v in %v", errorMessage, result.Stderr())
	}
}

func AssertCronjobPresent(c *clients.Clients, cronJobName, namespace string) {
	err := wait.Poll(config.APIRetry, config.ResourceTimeout, func() (bool, error) {
		log.Printf("Verifying if cronjob with prefix %v is present in namespace %v", cronJobName, namespace)
		cjs, err := c.KubeClient.Kube.BatchV1().CronJobs(namespace).List(c.Ctx, v1.ListOptions{})
		if err != nil {
			return false, err
		}
		for _, cj := range cjs.Items {
			if strings.Contains(cj.Name, cronJobName) {
				return true, nil
			}
		}
		return false, nil
	})
	if err != nil {
		assert.FailOnError(fmt.Errorf("Expected: cronjob with prefix %v to be present in namespace %v, Actual:cronjob with prefix %v to be not present in namespace %v", cronJobName, namespace, cronJobName, namespace))
	}
	fmt.Printf("Cronjob with prefix %v is present in namespace %v", cronJobName, namespace)
}

func AssertCronjobNotPresent(c *clients.Clients, cronJobName, namespace string) {
	err := wait.Poll(config.APIRetry, config.ResourceTimeout, func() (bool, error) {
		log.Printf("Verifying if cronjob with prefix %v is present in namespace %v", cronJobName, namespace)
		cjs, err := c.KubeClient.Kube.BatchV1().CronJobs(namespace).List(c.Ctx, v1.ListOptions{})
		if err != nil {
			return false, err
		}
		for _, cj := range cjs.Items {
			if strings.Contains(cj.Name, cronJobName) {
				return false, nil
			}
		}
		return true, nil
	})
	if err != nil {
		assert.FailOnError(fmt.Errorf("Expected: cronjob with prefix %v to be present in namespace %v, Actual:cronjob with prefix %v to be not present in namespace %v", cronJobName, namespace, cronJobName, namespace))
	}
	fmt.Printf("Cronjob with prefix %v is present in namespace %v", cronJobName, namespace)
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

func AnnotateNamespace(namespace, annotation string) {
	log.Printf("output: %s\n", cmd.MustSucceed("oc", "annotate", "namespace", namespace, annotation).Stdout())
}

func RemovePrunerConfig() {
	cmd.Run("oc", "patch", "tektonconfig", "config", "-p", "[{ \"op\": \"remove\", \"path\": \"/spec/pruner\" }]", "--type=json")
}

func LabelNamespace(namespace, label string) {
	log.Printf("output: %s\n", cmd.MustSucceed("oc", "label", "namespace", namespace, label).Stdout())
}
