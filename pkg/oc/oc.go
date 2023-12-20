package oc

import (
	"log"
	"strings"
	"time"

	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/cmd"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	resource "github.com/openshift-pipelines/release-tests/pkg/config"
	"github.com/openshift-pipelines/release-tests/pkg/store"
)

// Create resources using oc command
func Create(path_dir, namespace string) {
	log.Printf("output: %s\n", cmd.MustSucceed("oc", "create", "-f", resource.Path(path_dir), "-n", namespace).Stdout())
}

func Apply(path_dir, namespace string) {
	log.Printf("output: %s\n", cmd.MustSucceed("oc", "apply", "-f", resource.Path(path_dir), "-n", namespace).Stdout())
}

// Delete resources using oc command
func Delete(path_dir, namespace string) {
	log.Printf("output: %s\n", cmd.MustSucceed("oc", "delete", "-f", resource.Path(path_dir), "-n", namespace).Stdout())
}

// CreateNewProject Helps you to create new project
func CreateNewProject(ns string) {
	log.Printf("output: %s\n", cmd.MustSucceed("oc", "new-project", ns).Stdout())
}

// DeleteProject Helps you to delete new project
func DeleteProject(ns string) {
	log.Printf("output: %s\n", cmd.MustSucceed("oc", "delete", "project", ns).Stdout())
}

func DeleteProjectIgnoreErors(ns string) {
	log.Printf("output: %s\n", cmd.Run("oc", "delete", "project", ns).Stdout())
}

func LinkSecretToSA(secretname, sa, namespace string) {
	log.Printf("output: %s\n", cmd.MustSucceed("oc", "secret", "link", "serviceaccount/"+sa, "secrets/"+secretname, "-n", namespace).Stdout())
}

func CreateSecretWithSecretToken(secretname, namespace string) {
	log.Printf("output: %s\n", cmd.MustSucceed("oc", "create", "secret", "generic", secretname, "--from-literal=secretToken="+config.TriggersSecretToken, "-n", namespace).Stdout())
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

func AnnotateNamespace(namespace, annotation string) {
	log.Printf("output: %s\n", cmd.MustSucceed("oc", "annotate", "namespace", namespace, annotation).Stdout())
}

func AnnotateNamespaceIgnoreErrors(namespace, annotation string) {
	log.Printf("output: %s\n", cmd.Run("oc", "annotate", "namespace", namespace, annotation).Stdout())
}

func RemovePrunerConfig() {
	cmd.Run("oc", "patch", "tektonconfig", "config", "-p", "[{ \"op\": \"remove\", \"path\": \"/spec/pruner\" }]", "--type=json")
}

func LabelNamespace(namespace, label string) {
	log.Printf("output: %s\n", cmd.MustSucceed("oc", "label", "namespace", namespace, label).Stdout())
}

func DeleteResource(resourceType, name string) {
	log.Printf("output: %s\n", cmd.MustSucceed("oc", "delete", resourceType, name, "-n", store.Namespace()).Stdout())
}

func CheckProjectExists(projectName string) bool {
	return !strings.Contains(cmd.Run("oc", "project", projectName).String(), "error")
}

func SecretExists(secretName string, namespace string) bool {
	return !strings.Contains(cmd.Run("oc", "get", "secret", secretName, "-n", namespace).String(), "Error")
}

func CreateSecretForGitResolver(secretData string) {
	cmd.MustSucceed("oc", "create", "secret", "generic", "github-auth-secret", "--from-literal", "github-auth-key="+secretData, "-n", "openshift-pipelines")
}

func CreateSecretsForTektonResults(){
	var password string = cmd.MustSucceed("openssl", "rand", "-base64", "20").Stdout()
	password = strings.Replace(password, "\n", "", -1)
	cmd.MustSucceed("oc", "create", "secret", "-n", "openshift-pipelines", "generic", "tekton-results-postgres", "--from-literal=POSTGRES_USER=result", "--from-literal=POSTGRES_PASSWORD="+password)
	//generating tls certifiacte
	cmd.MustSucceed("openssl", "req", "-x509", "-newkey", "rsa:4096", "-keyout", "key.pem", "-out", "cert.pem", "-days", "365", "-nodes", "-subj", "/CN=tekton-results-api-service.openshift-pipelines.svc.cluster.local", "-addext", "subjectAltName=DNS:tekton-results-api-service.openshift-pipelines.svc.cluster.local")
	//creating secret with generated certificate
	cmd.MustSucceed("oc", "create", "secret", "tls", "-n", "openshift-pipelines", "tekton-results-tls", "--cert=cert.pem", "--key=key.pem")
}

func EnsureResutsReady(){
	cmd.MustSucceed("oc", "wait", "--for=condition=Ready", "tektoninstallerset", "-l", "operator.tekton.dev/type=result")
}

func GetLogsAndAnnotaions(resorceType string) {
	if resorceType == "pr" {
		log.Printf(cmd.MustSuccedIncreasedTimeout(time.Minute*10, "tkn", "pipeline", "start", "pipeline-results", "--showlog").Stdout())
	} 
	if resorceType == "tr" {
		log.Printf(cmd.MustSuccedIncreasedTimeout(time.Minute*10, "tkn", "tr", "logs", "-f", "--last").Stdout())
	}
	log.Print(cmd.MustSucceed("tkn", resorceType, "describe", "--last", "-o", "jsonpath={.metadata.annotations}").Stdout())
}