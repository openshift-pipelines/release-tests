package pac

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/clients"
)

func VerifyPipelinesAsCodeEnable(cs *clients.Clients, section, inputField, enable string) (string, error) {
	// Construct the JSON payload based on the 'enable' parameter
	payload := fmt.Sprintf(`{"spec":{"platforms":{"openshift":{"%s":{"%s": %s}}}}}`, inputField, section, enable)

	cmd := exec.Command("oc", "patch", "tektonconfigs.operator.tekton.dev", "config", "--type", "merge", "-p", payload)

	// Run the 'oc' command
	if err := cmd.Run(); err != nil {
		// Step failed - Use testsuit.T.Fail to fail the step and provide an error message
		testsuit.T.Errorf("Failed to set PipelinesAsCode enable status: %v", err)
		return "", err
	}

	// Return a message indicating the status change
	return fmt.Sprintf("PipelinesAsCode enable status has been set to %s", enable), nil
}

func VerifyInstallerSets(cs *clients.Clients, expectedStatus string) {
	// Sleep for 30 seconds
	time.Sleep(30 * time.Second)
	cmd := exec.Command("oc", "get", "tektoninstallersets", "-o", "custom-columns=NAME:.metadata.name")
	cmdOutput, err := cmd.CombinedOutput()

	if err != nil {
		// Step failed - Use testsuit.T.Fail to fail the step and provide an error message
		testsuit.T.Errorf("Failed to Verify InstallerSets status: %v", err)
		return
	}

	installerSets := strings.Split(string(cmdOutput), "\n")
	found := false
	for _, line := range installerSets {
		// Skip the header line
		if strings.Contains(line, "NAME") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) != 1 {
			continue
		}

		name := parts[0]

		if strings.HasPrefix(name, "openshiftpipelinesascode-") {
			found = true
			break
		}
	}

	if expectedStatus == "present" && !found {
		// Step failed - Use testsuit.T.Fail to fail the step
		testsuit.T.Fail(fmt.Errorf("InstallerSets related to PAC are not present"))
	} else if expectedStatus == "not present" && found {
		// Step failed - Use testsuit.T.Fail to fail the step
		testsuit.T.Fail(fmt.Errorf("InstallerSets related to PAC are present"))
	}
}

// VerifyPACPodsStatus checks the status of pods related to PAC in the specified namespace.
func VerifyPACPodsStatus(cs *clients.Clients, expectedStatus, namespace string) {
	// Sleep for 30 seconds
	time.Sleep(30 * time.Second)
	cmd := exec.Command("oc", "get", "pods", "-n", namespace, "-o", "custom-columns=NAME:.metadata.name")
	cmdOutput, err := cmd.CombinedOutput()

	if err != nil {
		if expectedStatus == "not present" {
			// Step succeeded - Use testsuit.T.Fail to fail the step and provide an error message
			testsuit.T.Errorf("Failed to get pod information: %v", err)
		}
		return
	}

	podNames := strings.Split(string(cmdOutput), "\n")
	found := false
	for _, line := range podNames {
		// Skip the header line
		if strings.Contains(line, "NAME") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) != 1 {
			continue
		}

		name := parts[0]

		if strings.HasPrefix(name, "pipelines-as-code-") {
			found = true
			break
		}
	}

	if expectedStatus == "present" && !found {
		// Step failed - Use testsuit.T.Fail to fail the step and provide an error message
		testsuit.T.Fail(fmt.Errorf("Pods related to PAC are not present"))
	} else if expectedStatus == "not present" && found {
		// Step failed - Use testsuit.T.Fail to fail the step and provide an error message
		testsuit.T.Fail(fmt.Errorf("Pods related to PAC are present"))
	}
}

func VerifyPACCustomResource(cs *clients.Clients, expectedStatus string) {
	// Sleep for 30 seconds
	time.Sleep(30 * time.Second)
	cmd := exec.Command("oc", "get", "crd", "-o", "custom-columns=NAME:.metadata.name")
	cmdOutput, err := cmd.CombinedOutput()

	if err != nil {
		if expectedStatus == "not present" {
			// Step succeeded - Use testsuit.T.Fail to fail the step and provide an error message
			testsuit.T.Fail(fmt.Errorf("Failed to get CRD information: %v", err))
		}
		return
	}

	crdNames := strings.Split(string(cmdOutput), "\n")
	found := false
	for _, name := range crdNames {
		if name == "repositories.pipelinesascode.tekton.dev" {
			found = true
			break
		}
	}

	if expectedStatus == "present" && !found {
		// Step failed - Use testsuit.T.Fail to fail the step and provide an error message
		testsuit.T.Fail(fmt.Errorf("CRD pipeline as code is not present"))
	} else if expectedStatus == "not present" && found {
		// Step failed - Use testsuit.T.Fail to fail the step and provide an error message
		testsuit.T.Fail(fmt.Errorf("CRD pipeline as code is present"))
	}
}
