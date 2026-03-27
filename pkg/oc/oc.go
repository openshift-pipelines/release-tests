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

func CreateSecretInNamespace(secretData, secretName, namespace string) {
	cmd.MustSucceed("oc", "create", "secret", "generic", secretName, "--from-literal", "private-repo-token="+secretData, "-n", namespace)
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
	// Fetch package manifest JSON in a single call
	packageManifestJSON := cmd.MustSucceed("oc", "get", "packagemanifests", "openshift-pipelines-operator-rh", "-n", "openshift-marketplace", "-o", "json").Stdout()

	// Parse the JSON structure
	type Channel struct {
		Name           string `json:"name"`
		CurrentCSVDesc struct {
			Annotations map[string]string `json:"annotations"`
		} `json:"currentCSVDesc"`
	}

	type PackageManifest struct {
		Status struct {
			Channels []Channel `json:"channels"`
		} `json:"status"`
	}

	var packageManifest PackageManifest
	if err := json.Unmarshal([]byte(packageManifestJSON), &packageManifest); err != nil {
		return nil, fmt.Errorf("failed to parse package manifest JSON: %w", err)
	}

	// Build channel to skipRange mapping
	channelSkipRangeMap := make(map[string]string)
	for _, channel := range packageManifest.Status.Channels {
		skipRange, exists := channel.CurrentCSVDesc.Annotations["olm.skipRange"]
		if exists && skipRange != "" {
			channelSkipRangeMap[channel.Name] = skipRange
		}
	}

	if len(channelSkipRangeMap) == 0 {
		return nil, fmt.Errorf("no valid OLM Skip Ranges found")
	}
	return channelSkipRangeMap, nil
}

// extractMajorMinor extracts major.minor version from a full version string
// e.g., "1.19.2" -> "1.19", "1.18.1" -> "1.18"
func extractMajorMinor(version string) string {
	versionRegex := regexp.MustCompile(`^(\d+\.\d+)\.?\d*`)
	matches := versionRegex.FindStringSubmatch(version)
	if len(matches) >= 2 {
		return matches[1]
	}
	// Fallback: if regex fails, return the original version
	return version
}

// GetOlmSkipRange fetches OLM skipRange data and saves it to the specified file
// upgradeType: "pre-upgrade" or "post-upgrade"
// fieldName: unused parameter (kept for backward compatibility)
// fileName: path to JSON file where skipRange data will be stored (file must already exist)
func GetOlmSkipRange(upgradeType, fieldName, fileName string) {
	skipRangeMap, err := FetchOlmSkipRange()
	if err != nil {
		log.Printf("Error fetching OLM Skip Range: %v", err)
		testsuit.T.Fail(fmt.Errorf("failed to fetch OLM Skip Range: %w", err))
	}

	filePath := config.Path(fileName)

	// Read existing data from file
	file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		log.Printf("Error opening file %s: %v", fileName, err)
		testsuit.T.Fail(fmt.Errorf("failed to open file %s: %w", fileName, err))
	}
	defer file.Close()

	var existingData map[string]interface{}
	if err := json.NewDecoder(file).Decode(&existingData); err != nil {
		log.Printf("Error decoding existing data from file %s: %v", fileName, err)
		testsuit.T.Fail(fmt.Errorf("failed to decode JSON from file %s: %w", fileName, err))
	}

	// Store skipRange data based on upgrade type
	fieldKey := fmt.Sprintf("%s-olm-skip-range", upgradeType)
	existingData[fieldKey] = skipRangeMap
	upgradeTypeTitle := strings.ToUpper(upgradeType[:1]) + upgradeType[1:]
	log.Printf("%s OLM Skip Range stored: %+v", upgradeTypeTitle, skipRangeMap)

	// Write updated data back to file
	if _, err := file.Seek(0, 0); err != nil {
		log.Printf("Error seeking file %s: %v", fileName, err)
		testsuit.T.Fail(fmt.Errorf("failed to seek file %s: %w", fileName, err))
	}
	if err := file.Truncate(0); err != nil {
		log.Printf("Error truncating file %s: %v", fileName, err)
		testsuit.T.Fail(fmt.Errorf("failed to truncate file %s: %w", fileName, err))
		return
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(existingData); err != nil {
		log.Printf("Error writing data to file %s: %v", fileName, err)
		testsuit.T.Fail(fmt.Errorf("failed to write data to file %s: %w", fileName, err))
		return
	}

	log.Printf("OLM Skip Range for '%s' saved to file %s", upgradeType, fileName)
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
		// Extract major.minor from OSP_VERSION for channel matching
		// e.g., "1.19.2" -> "1.19" to match with "pipelines-1.19"
		ospMajorMinor := extractMajorMinor(ospVersion)
		log.Printf("Extracted major.minor '%s' from OSP_VERSION '%s' for channel matching", ospMajorMinor, ospVersion)
		for channel, skipRange := range skipRangeMap {
			if channel == "latest" {
				log.Printf("Skipping 'latest' channel as requested")
				continue
			}
			// Check if channel contains the major.minor version
			channelContainsVersion := strings.Contains(channel, ospMajorMinor)
			skipRangeContainsVersion := strings.Contains(skipRange, ospVersion)
			log.Printf("Channel: %s, SkipRange: %s", channel, skipRange)
			log.Printf("  - Channel contains major.minor '%s': %v", ospMajorMinor, channelContainsVersion)
			log.Printf("  - SkipRange contains OSP_VERSION '%s': %v", ospVersion, skipRangeContainsVersion)
			if channelContainsVersion && skipRangeContainsVersion {
				log.Printf("Success: OSP_VERSION '%s' found in channel '%s' (major.minor match) and its skipRange '%s'", ospVersion, channel, skipRange)
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
			ospMajorMinor := extractMajorMinor(ospVersion)
			testsuit.T.Fail(fmt.Errorf("Error: OSP_VERSION '%s' (major.minor: %s) not found in both channel name and skipRange for any non-latest channel", ospVersion, ospMajorMinor))
		}
	}
}

// isValidOspVersionPatchUpdate validates that the post-upgrade skipRange represents a valid
// patch update for the current OSP_VERSION. The upper bound should increase (e.g., <1.15.3 -> <1.15.4)
// while the lower bound remains unchanged.
func isValidOspVersionPatchUpdate(preSkipRange, postSkipRange string) bool {
	ospVersion := os.Getenv("OSP_VERSION")
	if ospVersion == "" {
		log.Printf("OSP_VERSION not set, cannot validate patch update")
		return false
	}

	// Verify post-upgrade skipRange contains the OSP_VERSION
	if !skipRangeContainsVersion(postSkipRange, ospVersion) {
		log.Printf("Post-upgrade skip range '%s' does not contain OSP_VERSION '%s'", postSkipRange, ospVersion)
		return false
	}

	// Parse skipRange format: >=X.Y.Z <X.Y.Z
	rangeRegex := regexp.MustCompile(`>=(\d+\.\d+\.\d+)\s*<(\d+\.\d+\.\d+)`)
	preMatches := rangeRegex.FindStringSubmatch(preSkipRange)
	postMatches := rangeRegex.FindStringSubmatch(postSkipRange)

	if len(preMatches) != 3 || len(postMatches) != 3 {
		log.Printf("Invalid skipRange format: pre='%s', post='%s'", preSkipRange, postSkipRange)
		return false
	}

	preLower, preUpper := preMatches[1], preMatches[2]
	postLower, postUpper := postMatches[1], postMatches[2]

	// Lower bound must remain unchanged
	if preLower != postLower {
		log.Printf("Lower bound changed from '%s' to '%s' (should remain unchanged)", preLower, postLower)
		return false
	}

	// Parse version numbers to compare patch versions
	versionRegex := regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)$`)
	preUpperMatches := versionRegex.FindStringSubmatch(preUpper)
	postUpperMatches := versionRegex.FindStringSubmatch(postUpper)

	if len(preUpperMatches) != 4 || len(postUpperMatches) != 4 {
		log.Printf("Invalid version format: preUpper='%s', postUpper='%s'", preUpper, postUpper)
		return false
	}

	// Extract and compare patch versions
	var preUpperPatchInt, postUpperPatchInt int
	if _, err := fmt.Sscanf(preUpperMatches[3], "%d", &preUpperPatchInt); err != nil {
		log.Printf("Failed to parse pre-upgrade patch version: %v", err)
		return false
	}
	if _, err := fmt.Sscanf(postUpperMatches[3], "%d", &postUpperPatchInt); err != nil {
		log.Printf("Failed to parse post-upgrade patch version: %v", err)
		return false
	}

	// Patch version must increase
	if postUpperPatchInt <= preUpperPatchInt {
		log.Printf("Patch version did not increase: %s -> %s", preUpper, postUpper)
		return false
	}

	log.Printf("Valid patch update detected: %s -> %s (OSP_VERSION: %s)", preUpper, postUpper, ospVersion)
	return true
}

func skipRangeContainsVersion(skipRange, version string) bool {
	return strings.Contains(skipRange, version)
}

// ValidateChannelSkipRangeBounds validates that each channel's skipRange has correct bounds:
// - Lower bound should start from previous version (e.g., pipelines-1.14 should have >=1.13.0)
// - Upper bound should match current channel version (e.g., pipelines-1.14 should have <1.14.X)
func ValidateChannelSkipRangeBounds() {
	skipRangeMap, err := FetchOlmSkipRange()
	if err != nil {
		log.Printf("Error fetching OLM Skip Range: %v", err)
		testsuit.T.Fail(fmt.Errorf("Error fetching OLM Skip Range: %v", err))
		return
	}

	log.Printf("Validating channel skipRange bounds to ensure correct lower and upper bounds")
	log.Printf("Available channels and skipRanges:")
	for channel, skipRange := range skipRangeMap {
		log.Printf("  - Channel: %s, SkipRange: %s", channel, skipRange)
	}

	validationErrors := []string{}
	skipRangePattern := regexp.MustCompile(`>=(\d+\.\d+\.\d+)\s*<(\d+\.\d+\.\d+)`)
	channelVersionPattern := regexp.MustCompile(`pipelines-(\d+)\.(\d+)`)

	for channel, skipRange := range skipRangeMap {
		if channel == "latest" {
			log.Printf("Skipping 'latest' channel")
			continue
		}

		// Extract version from channel name (e.g., "pipelines-1.14" -> major=1, minor=14)
		channelMatches := channelVersionPattern.FindStringSubmatch(channel)
		if len(channelMatches) != 3 {
			log.Printf("Warning: Channel '%s' does not match expected pattern 'pipelines-X.Y', skipping", channel)
			continue
		}

		var major, minor int
		if _, err := fmt.Sscanf(channelMatches[1], "%d", &major); err != nil {
			validationErrors = append(validationErrors, fmt.Sprintf("Channel '%s' has invalid major version", channel))
			continue
		}
		if _, err := fmt.Sscanf(channelMatches[2], "%d", &minor); err != nil {
			validationErrors = append(validationErrors, fmt.Sprintf("Channel '%s' has invalid minor version", channel))
			continue
		}

		channelVersion := fmt.Sprintf("%d.%d", major, minor)

		// Parse skipRange format: >=X.Y.Z <X.Y.Z
		skipRangeMatches := skipRangePattern.FindStringSubmatch(skipRange)
		if len(skipRangeMatches) != 3 {
			validationErrors = append(validationErrors, fmt.Sprintf("Channel '%s' has invalid skipRange format: '%s' (expected format: '>=X.Y.Z <X.Y.Z')", channel, skipRange))
			continue
		}

		lowerBound := skipRangeMatches[1] // e.g., "1.13.0"
		upperBound := skipRangeMatches[2] // e.g., "1.14.5"

		lowerMajorMinor := extractMajorMinor(lowerBound)
		upperMajorMinor := extractMajorMinor(upperBound)

		log.Printf("\nValidating channel: %s (version: %s)", channel, channelVersion)
		log.Printf("  SkipRange: %s", skipRange)
		log.Printf("  Lower bound: %s (major.minor: %s)", lowerBound, lowerMajorMinor)
		log.Printf("  Upper bound: %s (major.minor: %s)", upperBound, upperMajorMinor)

		// Calculate previous version (e.g., 1.14 -> 1.13, 1.16 -> 1.15)
		prevMinor := minor - 1
		var prevVersion string
		if prevMinor < 0 {
			// If minor is 0, previous would be (major-1).X, but this case is unlikely for pipelines
			// For now, we'll skip validation for this edge case
			log.Printf("  ⚠️ Channel '%s' has minor version 0, skipping previous version validation", channel)
			continue
		}
		prevVersion = fmt.Sprintf("%d.%d", major, prevMinor)

		// Validate lower bound starts from previous version
		if lowerMajorMinor != prevVersion {
			validationErrors = append(validationErrors, fmt.Sprintf("Channel '%s' (version %s) has lower bound '%s' (major.minor: %s) that doesn't match previous version '%s'. Expected: >=%s.0", channel, channelVersion, lowerBound, lowerMajorMinor, prevVersion, prevVersion))
		} else {
			log.Printf("  ✅ Lower bound correctly starts from previous version: %s", prevVersion)
		}

		// Validate upper bound matches current channel version
		if upperMajorMinor != channelVersion {
			validationErrors = append(validationErrors, fmt.Sprintf("Channel '%s' (version %s) has upper bound '%s' (major.minor: %s) that doesn't match channel version. Expected: <%s.X", channel, channelVersion, upperBound, upperMajorMinor, channelVersion))
		} else {
			log.Printf("  ✅ Upper bound correctly matches channel version: %s", channelVersion)
		}
	}

	if len(validationErrors) > 0 {
		log.Printf("\n❌ Channel skipRange bounds validation failed with %d error(s):", len(validationErrors))
		for _, err := range validationErrors {
			log.Printf("  - %s", err)
		}
		testsuit.T.Fail(fmt.Errorf("Channel skipRange bounds validation failed: %v", strings.Join(validationErrors, "; ")))
	} else {
		log.Printf("\n✅ Success: All channel skipRanges have correct bounds - lower bound starts from previous version and upper bound matches channel version")
	}
}

// ValidateOlmSkipRangeDiff validates that skipRange changes between pre-upgrade and post-upgrade
// are valid. Only the channel matching the current OSP_VERSION should have its upper bound updated.
func ValidateOlmSkipRangeDiff(fileName string, preUpgradeSkipRange string, postUpgradeSkipRange string) {
	filePath := config.Path(fileName)
	file, err := os.Open(filePath)
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("failed to open file %s: %w", fileName, err))
		return
	}
	defer file.Close()

	var skipRangeData map[string]interface{}
	if err := json.NewDecoder(file).Decode(&skipRangeData); err != nil {
		testsuit.T.Fail(fmt.Errorf("failed to decode JSON from file %s: %w", fileName, err))
		return
	}

	// Extract pre-upgrade and post-upgrade data
	preUpgradeData, preExists := skipRangeData[preUpgradeSkipRange]
	postUpgradeData, postExists := skipRangeData[postUpgradeSkipRange]
	if !preExists || !postExists {
		testsuit.T.Fail(fmt.Errorf("missing skip range data: pre-upgrade exists=%v, post-upgrade exists=%v", preExists, postExists))
		return
	}

	preUpgradeMap, ok1 := preUpgradeData.(map[string]interface{})
	postUpgradeMap, ok2 := postUpgradeData.(map[string]interface{})
	if !ok1 || !ok2 {
		testsuit.T.Fail(fmt.Errorf("skip range data is not in expected map format"))
		return
	}

	log.Printf("Pre-Upgrade Skip Range: %+v", preUpgradeMap)
	log.Printf("Post-Upgrade Skip Range: %+v", postUpgradeMap)

	ospVersion := os.Getenv("OSP_VERSION")
	ospMajorMinor := ""
	if ospVersion != "" {
		ospMajorMinor = extractMajorMinor(ospVersion)
		log.Printf("Validating with OSP_VERSION: %s (major.minor: %s)", ospVersion, ospMajorMinor)
	} else {
		log.Printf("OSP_VERSION not set, validating that all channels remain unchanged")
	}

	validationErrors := []string{}
	skipRangePattern := regexp.MustCompile(`>=(\d+\.\d+\.\d+)\s*<(\d+\.\d+\.\d+)`)

	log.Printf("Validating channels (ignoring 'latest' channel)")

	// Validate each pre-upgrade channel
	for channel, preSkipRangeInterface := range preUpgradeMap {
		if channel == "latest" {
			log.Printf("Skipping 'latest' channel")
			continue
		}

		preSkipRange, ok := preSkipRangeInterface.(string)
		if !ok {
			validationErrors = append(validationErrors, fmt.Sprintf("Channel '%s' has invalid pre-upgrade skipRange format", channel))
			continue
		}

		postSkipRangeInterface, exists := postUpgradeMap[channel]
		if !exists {
			validationErrors = append(validationErrors, fmt.Sprintf("Channel '%s' missing in post-upgrade data", channel))
			continue
		}

		postSkipRange, ok := postSkipRangeInterface.(string)
		if !ok {
			validationErrors = append(validationErrors, fmt.Sprintf("Channel '%s' has invalid post-upgrade skipRange format", channel))
			continue
		}

		// Check if this channel matches the current OSP_VERSION (e.g., "pipelines-1.15" matches OSP_VERSION "1.15.4")
		channelMatchesOspVersion := ospVersion != "" && strings.Contains(channel, ospMajorMinor)

		// Check if skipRange changed
		if preSkipRange == postSkipRange {
			if channelMatchesOspVersion {
				// For channel matching OSP_VERSION, skipRange MUST change (upper bound should increase)
				// If unchanged, it's an error - the channel should have been updated for the new patch version
				validationErrors = append(validationErrors, fmt.Sprintf("Channel '%s' (matching OSP_VERSION %s) skipRange unchanged. Expected update from '%s' to include OSP_VERSION %s in upper bound, but got: %s", channel, ospVersion, preSkipRange, ospVersion, postSkipRange))
				continue
			}
			log.Printf("✅ Channel '%s': unchanged (%s)", channel, preSkipRange)
			continue
		}

		log.Printf("Channel '%s': changed from '%s' to '%s'", channel, preSkipRange, postSkipRange)

		// Parse skipRange format: >=X.Y.Z <X.Y.Z
		preMatches := skipRangePattern.FindStringSubmatch(preSkipRange)
		postMatches := skipRangePattern.FindStringSubmatch(postSkipRange)
		if len(preMatches) != 3 || len(postMatches) != 3 {
			validationErrors = append(validationErrors, fmt.Sprintf("Channel '%s' has invalid skipRange format", channel))
			continue
		}

		preLower, preUpper := preMatches[1], preMatches[2]
		postLower, postUpper := postMatches[1], postMatches[2]

		if channelMatchesOspVersion {
			// For the channel matching OSP_VERSION, validate patch update
			if preLower != postLower {
				validationErrors = append(validationErrors, fmt.Sprintf("Channel '%s' (matching OSP_VERSION %s) lower bound changed from '%s' to '%s' (should remain unchanged)", channel, ospVersion, preLower, postLower))
				continue
			}

			if isValidOspVersionPatchUpdate(preSkipRange, postSkipRange) {
				log.Printf("✅ Channel '%s': valid patch update for OSP_VERSION %s (%s -> %s)", channel, ospVersion, preUpper, postUpper)
			} else {
				validationErrors = append(validationErrors, fmt.Sprintf("Channel '%s' (matching OSP_VERSION %s) has invalid patch update: %s -> %s", channel, ospVersion, preSkipRange, postSkipRange))
			}
		} else {
			// For channels not matching OSP_VERSION, no changes allowed
			validationErrors = append(validationErrors, fmt.Sprintf("Channel '%s' skipRange changed from '%s' to '%s' (should remain unchanged, not matching OSP_VERSION %s)", channel, preSkipRange, postSkipRange, ospVersion))
		}
	}

	// Validate new channels in post-upgrade (for new major.minor releases)
	for postChannel, postSkipRangeInterface := range postUpgradeMap {
		if postChannel == "latest" {
			continue
		}
		if _, existedInPre := preUpgradeMap[postChannel]; !existedInPre {
			// This is a new channel - validate if it matches the current OSP_VERSION
			log.Printf("New channel '%s' found in post-upgrade data", postChannel)

			postSkipRange, ok := postSkipRangeInterface.(string)
			if !ok {
				validationErrors = append(validationErrors, fmt.Sprintf("New channel '%s' has invalid skipRange format", postChannel))
				continue
			}

			// Check if this new channel matches the current OSP_VERSION (new major.minor release)
			channelMatchesOspVersion := ospVersion != "" && strings.Contains(postChannel, ospMajorMinor)

			if channelMatchesOspVersion {
				// Validate the new channel's skipRange format and that it contains OSP_VERSION
				postMatches := skipRangePattern.FindStringSubmatch(postSkipRange)
				if len(postMatches) != 3 {
					validationErrors = append(validationErrors, fmt.Sprintf("New channel '%s' (matching OSP_VERSION %s) has invalid skipRange format: %s", postChannel, ospVersion, postSkipRange))
					continue
				}

				postUpper := postMatches[2]
				postUpperMajorMinor := extractMajorMinor(postUpper)

				// Validate that upper bound matches the channel version
				if postUpperMajorMinor != ospMajorMinor {
					validationErrors = append(validationErrors, fmt.Sprintf("New channel '%s' (matching OSP_VERSION %s) has upper bound '%s' (major.minor: %s) that doesn't match channel version", postChannel, ospVersion, postUpper, postUpperMajorMinor))
					continue
				}

				// Validate that skipRange contains the OSP_VERSION
				if !skipRangeContainsVersion(postSkipRange, ospVersion) {
					validationErrors = append(validationErrors, fmt.Sprintf("New channel '%s' (matching OSP_VERSION %s) skipRange '%s' does not contain OSP_VERSION", postChannel, ospVersion, postSkipRange))
					continue
				}

				log.Printf("✅ New channel '%s': valid for new major.minor release OSP_VERSION %s (skipRange: %s)", postChannel, ospVersion, postSkipRange)
			} else {
				// New channel that doesn't match OSP_VERSION - this shouldn't happen in normal upgrade scenarios
				log.Printf("⚠️ New channel '%s' found but doesn't match OSP_VERSION %s - this may indicate an unexpected new release", postChannel, ospVersion)
			}
		}
	}

	// Report results
	if len(validationErrors) > 0 {
		log.Printf("❌ OLM Skip Range validation failed with %d error(s):", len(validationErrors))
		for _, err := range validationErrors {
			log.Printf("  - %s", err)
		}
		testsuit.T.Fail(fmt.Errorf("OLM Skip Range validation failed: %v", strings.Join(validationErrors, "; ")))
	} else {
		log.Printf("✅ OLM Skip Range validation passed - all expected changes detected correctly")
	}
}
