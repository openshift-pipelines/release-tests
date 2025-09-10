package oc

import (
	"encoding/json"
	"log"
	"slices"
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

// Create resources using remote path using oc command
func CreateRemote(remote_path, namespace string) {
	log.Printf("output: %s\n", cmd.MustSucceed("oc", "create", "-f", remote_path, "-n", namespace).Stdout())
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
	// Wait for project to be completely deleted if it's terminating
	waitForProjectDeletion(ns)
	log.Printf("output: %s\n", cmd.MustSucceed("oc", "new-project", ns).Stdout())
}

// waitForProjectDeletion waits for a project to be completely deleted
func waitForProjectDeletion(projectName string) {
	// Check if project exists and is terminating
	statusResult := cmd.Run("oc", "get", "project", projectName, "-o", "jsonpath={.status.phase}")
	if statusResult.ExitCode != 0 {
		// Project doesn't exist, no need to wait
		return
	}

	phase := strings.TrimSpace(statusResult.Stdout())
	if phase != "Terminating" {
		// Project is not terminating, no need to wait
		return
	}

	log.Printf("Project %s is terminating, waiting for deletion to complete...", projectName)

	checkInterval := 10 * time.Second
	maxAttempts := 30 // 30 attempts * 10 seconds = 5 minutes max wait
	attempts := 0

	for attempts < maxAttempts {
		checkResult := cmd.Run("oc", "get", "project", projectName)
		if checkResult.ExitCode != 0 {
			// Project no longer exists
			log.Printf("Project %s has been fully deleted", projectName)
			return
		}

		attempts++
		log.Printf("Still waiting for project %s to be deleted... (attempt %d/%d)", projectName, attempts, maxAttempts)
		time.Sleep(checkInterval)
	}

	log.Printf("Warning: Timed out waiting for project %s deletion after %d attempts", projectName, maxAttempts)
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

func DeleteResourceInNamespace(resourceType, name, namespace string) {
	log.Printf("output: %s\n", cmd.MustSucceed("oc", "delete", resourceType, name, "-n", namespace).Stdout())
}

func CheckProjectExists(projectName string) bool {
	// Check project status first to avoid issues with terminating projects
	statusResult := cmd.Run("oc", "get", "project", projectName, "-o", "jsonpath={.status.phase}")
	if statusResult.ExitCode != 0 {
		// Project doesn't exist
		return false
	}

	phase := strings.TrimSpace(statusResult.Stdout())
	// Don't consider Terminating projects as existing
	if phase == "Terminating" {
		log.Printf("Project %s is in Terminating state, treating as non-existent", projectName)
		return false
	}

	// Double-check with project command
	commandResult := cmd.Run("oc", "project", projectName)
	return commandResult.ExitCode == 0 && !strings.Contains(commandResult.String(), "error")
}

func SecretExists(secretName string, namespace string) bool {
	return !strings.Contains(cmd.Run("oc", "get", "secret", secretName, "-n", namespace).String(), "Error")
}

func CreateSecretForGitResolver(secretData string) {
	cmd.MustSucceed("oc", "create", "secret", "generic", "github-auth-secret", "--from-literal", "github-auth-key="+secretData, "-n", "openshift-pipelines")
}

func CreateSecretForWebhook(tokenSecretData, webhookSecretData, namespace string) {
	cmd.MustSucceed("oc", "create", "secret", "generic", "gitlab-webhook-config", "--from-literal", "provider.token="+tokenSecretData, "--from-literal", "webhook.secret="+webhookSecretData, "-n", namespace)
}

func EnableConsolePlugin() {
	json_output := cmd.MustSucceed("oc", "get", "consoles.operator.openshift.io", "cluster", "-o", "jsonpath={.spec.plugins}").Stdout()
	log.Printf("Already enabled console plugins: %s", json_output)
	var plugins []string

	if len(json_output) > 0 {
		err := json.Unmarshal([]byte(json_output), &plugins)

		if err != nil {
			testsuit.T.Errorf("Could not parse consoles.operator.openshift.io CR: %v", err)
		}

		if slices.Contains(plugins, config.ConsolePluginDeployment) {
			log.Printf("Pipelines console plugin is already enabled.")
			return
		}
	}

	plugins = append(plugins, config.ConsolePluginDeployment)

	patch_data := "{\"spec\":{\"plugins\":[\"" + strings.Join(plugins, "\",\"") + "\"]}}"
	cmd.MustSucceed("oc", "patch", "consoles.operator.openshift.io", "cluster", "-p", patch_data, "--type=merge").Stdout()
}

func GetSecretsData(secretName, namespace string) string {
	return cmd.MustSucceed("oc", "get", "secrets", secretName, "-n", namespace, "-o", "jsonpath=\"{.data}\"").Stdout()
}

func CreateChainsImageRegistrySecret(dockerConfig string) {
	cmd.MustSucceed("oc", "create", "secret", "generic", "chains-image-registry-credentials", "--from-literal=.dockerconfigjson="+dockerConfig, "--from-literal=config.json="+dockerConfig, "--type=kubernetes.io/dockerconfigjson")
}
