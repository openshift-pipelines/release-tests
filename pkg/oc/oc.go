package oc

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
	"time"

	"regexp"

	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/cmd"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	"github.com/openshift-pipelines/release-tests/pkg/store"
)

// Create resources using oc command
func Create(path_dir, namespace string) {
	log.Printf("output: %s\n", cmd.MustSucceed("oc", "create", "-f", config.Path(path_dir), "-n", namespace).Stdout())
}

// Create resources using remote path using oc command
func CreateRemote(remote_path, namespace string) {
	log.Printf("output: %s\n", cmd.MustSucceed("oc", "create", "-f", remote_path, "-n", namespace).Stdout())
}

func Apply(path_dir, namespace string) {
	log.Printf("output: %s\n", cmd.MustSucceed("oc", "apply", "-f", config.Path(path_dir), "-n", namespace).Stdout())
}

// Delete resources using oc command
func Delete(path_dir, namespace string) {
	// Tekton Results sets a finalizer that prevent resource removal for some time
	// see parameters "store_deadline" and "forward_buffer"
	// by default, it waits at least 150 seconds
	log.Printf("output: %s\n", cmd.MustSuccedIncreasedTimeout(time.Second*300, "oc", "delete", "-f", config.Path(path_dir), "-n", namespace).Stdout())
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

func FetchOlmSkipRange() (map[string]string, error) {
	skipRangesJson := cmd.MustSucceed("bash", "-c", `oc get packagemanifests openshift-pipelines-operator-rh -n openshift-marketplace -o json | jq -r '.status.channels[].currentCSVDesc.annotations["olm.skipRange"]'`).Stdout()
	skipRanges := strings.Split(strings.TrimSpace(skipRangesJson), "\n")
	channelsJson := cmd.MustSucceed("bash", "-c", `oc get packagemanifests openshift-pipelines-operator-rh -n openshift-marketplace -o json | jq -r '.status.channels[].name'`).Stdout()
	channels := strings.Split(strings.TrimSpace(channelsJson), "\n")

	if len(channels) != len(skipRanges) {
		return nil, fmt.Errorf("mismatch between number of channels (%d) and skipRanges (%d)", len(channels), len(skipRanges))
	}

	channelSkipRangeMap := make(map[string]string)
	for i, channel := range channels {
		if skipRanges[i] != "null" && skipRanges[i] != "" {
			channelSkipRangeMap[channel] = skipRanges[i]
		}
	}

	if len(channelSkipRangeMap) == 0 {
		return nil, fmt.Errorf("no valid OLM Skip Ranges found")
	}
	return channelSkipRangeMap, nil
}

func GetOlmSkipRange(upgradeType, fieldName, fileName string) {
	skipRangeMap, err := FetchOlmSkipRange()
	if err != nil {
		log.Printf("Error fetching OLM Skip Range: %v", err)
		return
	}
	file, err := os.OpenFile(config.Path(fileName), os.O_RDWR, 0644)
	if err != nil {
		log.Printf("Error opening file %s: %v", fileName, err)
		return
	}
	defer file.Close()
	var existingData map[string]interface{}
	if err := json.NewDecoder(file).Decode(&existingData); err != nil {
		log.Printf("Error decoding existing data from file %s: %v", fileName, err)
		return
	}
	switch upgradeType {
	case "pre-upgrade":
		existingData["pre-upgrade-olm-skip-range"] = skipRangeMap
		log.Printf("Pre-upgrade OLM Skip Range is stored as: %+v", skipRangeMap)
	case "post-upgrade":
		existingData["post-upgrade-olm-skip-range"] = skipRangeMap
		log.Printf("Post-upgrade OLM Skip Range is stored as: %+v", skipRangeMap)
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
	skipRangeMap, err := FetchOlmSkipRange()
	if err != nil {
		log.Printf("Error fetching OLM Skip Range: %v", err)
		return
	}

	ospVersion := os.Getenv("OSP_VERSION")
	log.Printf("Validating OSP_VERSION: %s", ospVersion)
	found := false

	if ospVersion == "5.0.5" {
		log.Printf("Detected nightly build (OSP_VERSION=5.0.5), validating only skipRange, not channel name")
		for channel, skipRange := range skipRangeMap {
			if channel == "latest" {
				log.Printf("Skipping 'latest' channel as requested")
				continue
			}
			skipRangeContainsVersion := strings.Contains(skipRange, ospVersion)
			log.Printf("Channel: %s, SkipRange: %s", channel, skipRange)
			log.Printf("  - SkipRange contains OSP_VERSION '%s': %v", ospVersion, skipRangeContainsVersion)

			if skipRangeContainsVersion {
				log.Printf("Success: OSP_VERSION '%s' found in skipRange for channel '%s': '%s'", ospVersion, channel, skipRange)
				found = true
				break
			}
		}
	} else {
		log.Printf("Regular release build, validating both channel name and skipRange")
		for channel, skipRange := range skipRangeMap {
			if channel == "latest" {
				log.Printf("Skipping 'latest' channel as requested")
				continue
			}
			channelContainsVersion := strings.Contains(channel, ospVersion)
			skipRangeContainsVersion := strings.Contains(skipRange, ospVersion)
			log.Printf("Channel: %s, SkipRange: %s", channel, skipRange)
			log.Printf("  - Channel contains OSP_VERSION '%s': %v", ospVersion, channelContainsVersion)
			log.Printf("  - SkipRange contains OSP_VERSION '%s': %v", ospVersion, skipRangeContainsVersion)
			if channelContainsVersion && skipRangeContainsVersion {
				log.Printf("Success: OSP_VERSION '%s' found in both channel '%s' and its skipRange '%s'", ospVersion, channel, skipRange)
				found = true
				break
			}
		}
	}

	if !found {
		log.Printf("Available channels and their skipRanges:")
		for channel, skipRange := range skipRangeMap {
			if channel != "latest" {
				log.Printf("  - Channel: %s, SkipRange: %s", channel, skipRange)
			}
		}

		if ospVersion == "5.0.5" {
			testsuit.T.Fail(fmt.Errorf("Error: OSP_VERSION '%s' not found in skipRange for any non-latest channel", ospVersion))
		} else {
			testsuit.T.Fail(fmt.Errorf("Error: OSP_VERSION '%s' not found in both channel name and skipRange for any non-latest channel", ospVersion))
		}
	}
}

func isValidOspVersionPatchUpdate(preSkipRange, postSkipRange string) bool {
	ospVersion := os.Getenv("OSP_VERSION")

	if !skipRangeContainsVersion(postSkipRange, ospVersion) {
		log.Printf("Post-upgrade skip range '%s' does not contain OSP_VERSION '%s'", postSkipRange, ospVersion)
		return false
	}

	rangeRegex := regexp.MustCompile(`>=(\d+\.\d+\.\d+)\s*<(\d+\.\d+\.\d+)`)
	preMatches := rangeRegex.FindStringSubmatch(preSkipRange)
	postMatches := rangeRegex.FindStringSubmatch(postSkipRange)

	preUpperBound := preMatches[2]   // Upper bound from pre-upgrade
	postUpperBound := postMatches[2] // Upper bound from post-upgrade

	versionRegex := regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)$`)
	preUpperMatches := versionRegex.FindStringSubmatch(preUpperBound)
	postUpperMatches := versionRegex.FindStringSubmatch(postUpperBound)

	preUpperPatch := preUpperMatches[3]
	postUpperPatch := postUpperMatches[3]

	var preUpperPatchInt, postUpperPatchInt int
	fmt.Sscanf(preUpperPatch, "%d", &preUpperPatchInt)
	fmt.Sscanf(postUpperPatch, "%d", &postUpperPatchInt)

	if postUpperPatchInt <= preUpperPatchInt {
		log.Printf("Version did not increase: %s -> %s", preUpperBound, postUpperBound)
		return false
	}

	log.Printf("Valid OSP version-based patch update detected: %s -> %s (OSP_VERSION: %s)",
		preUpperBound, postUpperBound, ospVersion)
	return true
}

func skipRangeContainsVersion(skipRange, version string) bool {
	return strings.Contains(skipRange, version)
}

func ValidateOlmSkipRangeDiff(fileName string, preUpgradeSkipRange string, postUpgradeSkipRange string) {
	file, err := os.Open(config.Path(fileName))
	if err != nil {
		log.Printf("Error opening file %s: %v", fileName, err)
		testsuit.T.Fail(fmt.Errorf("Error opening file %s: %v", fileName, err))
		return
	}
	defer file.Close()
	var skipRangeData map[string]interface{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&skipRangeData); err != nil {
		log.Printf("Error decoding JSON from file %s: %v", fileName, err)
		testsuit.T.Fail(fmt.Errorf("Error decoding JSON from file %s: %v", fileName, err))
		return
	}
	preUpgradeData, preExists := skipRangeData[preUpgradeSkipRange]
	postUpgradeData, postExists := skipRangeData[postUpgradeSkipRange]
	if !preExists || !postExists {
		log.Printf("Error: One of the skip ranges is missing. Pre-Upgrade exists: %v, Post-Upgrade exists: %v", preExists, postExists)
		testsuit.T.Fail(fmt.Errorf("One of the skip ranges is missing. Pre-Upgrade exists: %v, Post-Upgrade exists: %v", preExists, postExists))
		return
	}

	preUpgradeMap, ok1 := preUpgradeData.(map[string]interface{})
	postUpgradeMap, ok2 := postUpgradeData.(map[string]interface{})

	if !ok1 || !ok2 {
		log.Printf("Error: Skip range data is not in expected map format")
		testsuit.T.Fail(fmt.Errorf("Skip range data is not in expected map format"))
		return
	}

	log.Printf("Pre-Upgrade Skip Range: %+v", preUpgradeMap)
	log.Printf("Post-Upgrade Skip Range: %+v", postUpgradeMap)

	validationErrors := []string{}

	log.Printf("Validating that all pre-upgrade channels are preserved in post-upgrade (ignoring 'latest' channel)")

	// Check each channel from pre-upgrade data
	for preChannel, preSkipRange := range preUpgradeMap {
		// Skip 'latest' channel from pre-upgrade validation
		if preChannel == "latest" {
			log.Printf("Skipping 'latest' channel from pre-upgrade data as requested")
			continue
		}

		if postSkipRange, exists := postUpgradeMap[preChannel]; exists {
			if preSkipRange == postSkipRange {
				log.Printf("✅ Success: Channel '%s' preserved with skipRange: %v", preChannel, preSkipRange)
			} else {
				// There's a skipRange mismatch - check if it's related to current OSP_VERSION
				preSkipRangeStr, ok1 := preSkipRange.(string)
				postSkipRangeStr, ok2 := postSkipRange.(string)

				if ok1 && ok2 {
					ospVersion := os.Getenv("OSP_VERSION")

					if ospVersion != "" && (strings.Contains(preSkipRangeStr, ospVersion) || strings.Contains(postSkipRangeStr, ospVersion)) {
						log.Printf("ℹ️ SkipRange mismatch involves current OSP_VERSION '%s', applying OSP_VERSION-based validation", ospVersion)

						if isValidOspVersionPatchUpdate(preSkipRangeStr, postSkipRangeStr) {
							log.Printf("✅ Success: Channel '%s' updated with valid OSP_VERSION-based patch release: %v -> %v", preChannel, preSkipRange, postSkipRange)
						} else {
							validationErrors = append(validationErrors, fmt.Sprintf("Channel '%s' skipRange changed from '%v' to '%v' (not a valid OSP_VERSION-based patch update for OSP_VERSION '%s')", preChannel, preSkipRange, postSkipRange, ospVersion))
						}
					} else {
						if ospVersion == "" {
							log.Printf("ℹ️ OSP_VERSION not set, applying standard validation (no changes allowed)")
						} else {
							log.Printf("ℹ️ SkipRange mismatch does not involve current OSP_VERSION '%s', applying standard validation (no changes allowed)", ospVersion)
						}
						validationErrors = append(validationErrors, fmt.Sprintf("Channel '%s' skipRange changed from '%v' to '%v' (should remain unchanged)", preChannel, preSkipRange, postSkipRange))
					}
				} else {
					validationErrors = append(validationErrors, fmt.Sprintf("Channel '%s' skipRange changed from '%v' to '%v' (invalid format)", preChannel, preSkipRange, postSkipRange))
				}
			}
		} else {
			validationErrors = append(validationErrors, fmt.Sprintf("Pre-upgrade channel '%s' is missing in post-upgrade data", preChannel))
		}
	}

	log.Printf("Additional channels found in post-upgrade data:")
	for postChannel, postSkipRange := range postUpgradeMap {
		if _, existedInPre := preUpgradeMap[postChannel]; existedInPre {
			continue
		}

		if postChannel == "latest" {
			log.Printf("  - Ignoring 'latest' channel in post-upgrade as requested")
			continue
		}

		log.Printf("  - New channel '%s' with skipRange: %v", postChannel, postSkipRange)
	}

	if len(validationErrors) > 0 {
		log.Printf("❌ OLM Skip Range validation failed with errors:")
		for _, err := range validationErrors {
			log.Printf("  - %s", err)
		}
		testsuit.T.Fail(fmt.Errorf("OLM Skip Range validation failed: %v", strings.Join(validationErrors, "; ")))
	} else {
		log.Printf("✅ Success: OLM Skip Range validation passed - all expected changes detected correctly")
	}
}
