package oc

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
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
	// Tekton Results sets a finalizer that prevent resource removal for some time
	// see parameters "store_deadline" and "forward_buffer"
	// by default, it waits at least 150 seconds
	log.Printf("output: %s\n", cmd.MustSuccedIncreasedTimeout(time.Second*300, "oc", "delete", "-f", resource.Path(path_dir), "-n", namespace).Stdout())
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
	// Tekton Results sets a finalizer that prevent resource removal for some time
	// see parameters "store_deadline" and "forward_buffer"
	// by default, it waits at least 150 seconds
	log.Printf("output: %s\n", cmd.MustSuccedIncreasedTimeout(time.Second*300, "oc", "delete", resourceType, name, "-n", store.Namespace()).Stdout())
}

func DeleteResourceInNamespace(resourceType, name, namespace string) {
	log.Printf("output: %s\n", cmd.MustSucceed("oc", "delete", resourceType, name, "-n", namespace).Stdout())
}

func CheckProjectExists(projectName string) bool {
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

func CopySecret(secretName string, sourceNamespace string, destNamespace string) {
	secretJson := cmd.MustSucceed("oc", "get", "secret", secretName, "-n", sourceNamespace, "-o", "json").Stdout()
	cmdOutput := cmd.MustSucceed("bash", "-c", fmt.Sprintf(`echo '%s' | jq 'del(.metadata["namespace", "creationTimestamp", "resourceVersion", "selfLink", "uid", "annotations"]) | .data |= with_entries(if .key == "github-auth-key" then .key = "token" else . end)'`, secretJson)).Stdout()
	cmd.MustSucceed("bash", "-c", fmt.Sprintf(`echo '%s' | kubectl apply -n %s -f -`, cmdOutput, destNamespace))
	log.Printf("Successfully copied secret %s from %s to %s", secretName, sourceNamespace, destNamespace)
}

func FetchOlmSkipRange() (string, error) {
	olmManifestJson := cmd.MustSucceed("oc", "get", "packagemanifests", "openshift-pipelines-operator-rh", "-n", "openshift-marketplace", "-o", "json").Stdout()
	SkipRange := cmd.MustSucceed("bash", "-c", fmt.Sprintf(`echo '%s' | jq -r '.status.channels[].currentCSVDesc.annotations["olm.skipRange"]'`, olmManifestJson)).Stdout()
	if SkipRange == "" {
		return "", fmt.Errorf("OLM Skip Range is empty")
	}
	return SkipRange, nil
}

func GetOlmSkipRange(upgradeType, fieldName, fileName string) {
	SkipRange, err := FetchOlmSkipRange()
	if err != nil {
		log.Printf("Error fetching OLM Skip Range: %v", err)
		return
	}
	file, err := os.OpenFile(resource.Path(fileName), os.O_RDWR, 0644)
	if err != nil {
		log.Printf("Error opening file %s: %v", fileName, err)
		return
	}
	defer file.Close()
	var existingData map[string]string
	if err := json.NewDecoder(file).Decode(&existingData); err != nil {
		log.Printf("Error decoding existing data from file %s: %v", fileName, err)
		return
	}
	if upgradeType == "pre-upgrade" {
		existingData["pre-upgrade-olm-skip-range"] = SkipRange
		log.Printf("Pre-upgrade OLM Skip Range is stored as: %s", SkipRange)
	} else if upgradeType == "post-upgrade" {
		existingData["post-upgrade-olm-skip-range"] = SkipRange
		log.Printf("Post-upgrade OLM Skip Range is stored as: %s", SkipRange)
	}
	if _, err := file.Seek(0, 0); err != nil {
		log.Printf("Error seeking file %s: %v", fileName, err)
		return
	}
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty-print the JSON output
	if err := encoder.Encode(existingData); err != nil {
		log.Printf("Error writing updated data to file %s: %v", fileName, err)
		return
	}
	log.Printf("OLM Skip Range for '%s' has been saved to file %s", fieldName, fileName)
}

func ValidateOlmSkipRange() {
	SkipRange, err := FetchOlmSkipRange()
	if err != nil {
		log.Printf("Error fetching OLM Skip Range: %v", err)
		return
	}
	lines := strings.Split(SkipRange, "\n")
	if len(lines) == 0 {
		log.Printf("Error: No lines found in OLM Skip Range")
		return
	}
	firstLine := lines[0]
	pipelineVersion := os.Getenv("OPERATOR_VERSION")
	if strings.Contains(firstLine, pipelineVersion) {
		log.Printf("Success: OPERATOR_VERSION '%s' matches the first line of OLM Skip Range: '%s'", pipelineVersion, firstLine)
	} else {
		testsuit.T.Fail(fmt.Errorf("Error: OPERATOR_VERSION '%s' does not match the first line of OLM Skip Range: '%s'", pipelineVersion, firstLine))
	}
}

func ValidateOlmSkipRangeDiff(fileName string, preUpgradeSkipRange string, postUpgradeSkipRange string) {
	file, err := os.Open(resource.Path(fileName))
	if err != nil {
		log.Printf("Error opening file %s: %v", fileName, err)
		testsuit.T.Fail(fmt.Errorf("Error opening file %s: %v", fileName, err))
		return
	}
	defer file.Close()
	var skipRangeData map[string]string
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&skipRangeData); err != nil {
		log.Printf("Error decoding JSON from file %s: %v", fileName, err)
		testsuit.T.Fail(fmt.Errorf("Error decoding JSON from file %s: %v", fileName, err))
		return
	}
	preUpgradeSkipRange, preExists := skipRangeData[preUpgradeSkipRange]
	postUpgradeSkipRange, postExists := skipRangeData[postUpgradeSkipRange]
	if !preExists || !postExists || preUpgradeSkipRange == "" || postUpgradeSkipRange == "" {
		log.Printf("Error: One of the skip ranges is missing or empty. Pre-Upgrade: %v, Post-Upgrade: %v", preUpgradeSkipRange, postUpgradeSkipRange)
		testsuit.T.Fail(fmt.Errorf("One of the skip ranges is missing or empty. Pre-Upgrade: %v, Post-Upgrade: %v", preUpgradeSkipRange, postUpgradeSkipRange))
		return
	}
	log.Printf("Pre-Upgrade Skip Range: %v", preUpgradeSkipRange)
	log.Printf("Post-Upgrade Skip Range: %v", postUpgradeSkipRange)
	if preUpgradeSkipRange == postUpgradeSkipRange {
		log.Printf("OLM Skip Range before and after upgrade are the same. Pre-Upgrade: %s, Post-Upgrade: %s", preUpgradeSkipRange, postUpgradeSkipRange)
	} else {
		log.Printf("OLM Skip Range mismatch detected! Pre-Upgrade: '%s', Post-Upgrade: '%s'", preUpgradeSkipRange, postUpgradeSkipRange)
		testsuit.T.Fail(fmt.Errorf("OLM Skip Range mismatch detected! Pre-Upgrade: '%s', Post-Upgrade: '%s'", preUpgradeSkipRange, postUpgradeSkipRange))
	}
}
