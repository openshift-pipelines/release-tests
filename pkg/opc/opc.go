package opc

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/cmd"
	"gotest.tools/v3/icmd"
)

type Cmd struct {
	// path to opc binary
	Path string
}

type PipelineRunList struct {
	Name   string
	Status string
}

type PacInfoInstall struct {
	PipelinesAsCode PipelinesAsCodeSection
}

type PipelinesAsCodeSection struct {
	InstallVersion   string
	InstallNamespace string
}

// New initializes Cmd
func New(opcPath string) Cmd {
	return Cmd{
		Path: opcPath,
	}
}

func GetOPCServerVersion(component string) string {
	var version string
	output := cmd.MustSucceed("opc", "version", "-s").Stdout()
	titleComp := strings.ToUpper(component[:1]) + component[1:]
	prefix := titleComp + " version:"

	for _, line := range strings.Split(output, "\n") {
		if strings.HasPrefix(line, prefix) {
			version = strings.TrimSpace(strings.TrimPrefix(line, prefix))
			break
		}
	}

	if strings.Contains(version, "unknown") {
		testsuit.T.Errorf("%s is not installed", titleComp)
	}
	return version
}

// Verify the versions of Openshift Pipelines components
func AssertComponentVersion(version string, component string) {
	var actualVersion string
	switch component {
	case "pipeline", "triggers", "operator", "chains":
		actualVersion = GetOPCServerVersion(component)
	case "OSP":
		actualVersion = cmd.MustSucceed("oc", "get", "tektonconfig", "config", "-o", "jsonpath={.status.version}").Stdout()
	case "pac":
		actualVersion = cmd.MustSucceed("oc", "get", "pac", "pipelines-as-code", "-o", "jsonpath={.status.version}").Stdout()
	case "hub":
		actualVersion = cmd.MustSucceed("oc", "get", "tektonhub", "hub", "-o", "jsonpath={.status.version}").Stdout()
	case "results":
		actualVersion = cmd.MustSucceed("oc", "get", "tektonresult", "result", "-o", "jsonpath={.status.version}").Stdout()
	case "manual-approval-gate":
		actualVersion = cmd.MustSucceed("oc", "get", "manualapprovalgate", "manual-approval-gate", "-o", "jsonpath={.status.version}").Stdout()
	default:
		testsuit.T.Errorf("Unknown component")
	}

	actualVersion = strings.Trim(actualVersion, "\n")
	if !strings.Contains(actualVersion, version) {
		testsuit.T.Errorf("The %s has an unexpected version: %s, expected: %s", component, actualVersion, version)
	}
}

func DownloadCLIFromCluster() {
	var architecture = strings.Trim(cmd.MustSucceed("uname").Stdout(), "\n") + " " + strings.Trim(cmd.MustSucceed("uname", "-m").Stdout(), "\n")
	var cliDownloadURL = cmd.MustSucceed("oc", "get", "consoleclidownloads", "tkn", "-o", "jsonpath={.spec.links[?(@.text==\"Download tkn and tkn-pac for "+architecture+"\")].href}").Stdout()
	result := cmd.MustSuccedIncreasedTimeout(time.Minute*10, "curl", "-o", "/tmp/tkn-binary.tar.gz", "-k", cliDownloadURL)
	if result.ExitCode != 0 {
		testsuit.T.Errorf("%s", result.Stderr())
	}
	cmd.MustSucceed("tar", "-xf", "/tmp/tkn-binary.tar.gz", "-C", "/tmp")
}

func AssertClientVersion(binary string) {
	var commandResult, unexpectedVersion string

	switch binary {
	case "tkn-pac":
		commandResult = cmd.MustSucceed("/tmp/tkn-pac", "version").Stdout()
		expectedVersion := os.Getenv("PAC_VERSION")
		if !strings.Contains(commandResult, expectedVersion) {
			testsuit.T.Errorf("tkn-pac has an unexpected version: %s. Expected: %s", commandResult, expectedVersion)
		}

	case "tkn":
		expectedVersion := os.Getenv("TKN_CLIENT_VERSION")
		commandResult = cmd.MustSucceed("/tmp/tkn", "version").Stdout()
		var splittedCommandResult = strings.Split(commandResult, "\n")
		for i := range splittedCommandResult {
			if strings.Contains(splittedCommandResult[i], "Client") {
				if !strings.Contains(splittedCommandResult[i], expectedVersion) {
					unexpectedVersion = splittedCommandResult[i]
					testsuit.T.Errorf("tkn client has an unexpected version: %s. Expected: %s", unexpectedVersion, expectedVersion)
				}
			}
		}

	case "opc":
		commandResult = cmd.MustSucceed("/tmp/opc", "version").Stdout()
		components := [3]string{"OpenShift Pipelines Client", "Tekton CLI", "Pipelines as Code CLI"}
		expectedVersions := [3]string{os.Getenv("OSP_VERSION"), os.Getenv("TKN_CLIENT_VERSION"), os.Getenv("PAC_VERSION")}
		splittedCommandResult := strings.Split(commandResult, "\n")
		for i := 0; i < 3; i++ {
			if strings.Contains(splittedCommandResult[i], components[i]) {
				if !strings.Contains(splittedCommandResult[i], expectedVersions[i]) {
					unexpectedVersion = splittedCommandResult[i]
					testsuit.T.Errorf("%s has an unexpected version: %s. Expected: %s", components[i], unexpectedVersion, expectedVersions[i])
				}
			}
		}

	default:
		testsuit.T.Errorf("Unknown binary or client")
	}
}

func AssertServerVersion(binary string) {
	var commandResult, unexpectedVersion string

	switch binary {
	case "opc":
		commandResult = cmd.MustSucceed("/tmp/opc", "version", "--server").Stdout()
		components := [4]string{"Chains version", "Pipeline version", "Triggers version", "Operator version"}
		expectedVersions := [4]string{os.Getenv("CHAINS_VERSION"), os.Getenv("PIPELINE_VERSION"), os.Getenv("TRIGGERS_VERSION"), os.Getenv("OPERATOR_VERSION")}
		splittedCommandResult := strings.Split(commandResult, "\n")
		for i := 0; i < 4; i++ {
			if strings.Contains(splittedCommandResult[i], components[i]) {
				if !strings.Contains(splittedCommandResult[i], expectedVersions[i]) {
					unexpectedVersion = splittedCommandResult[i]
					testsuit.T.Errorf("%s has an unexpected version: %s. Expected: %s", components[i], unexpectedVersion, expectedVersions[i])
				}
			}
		}
	default:
		testsuit.T.Errorf("Unknown binary or client")
	}

}

func ValidateQuickstarts() {
	cmd.MustSucceed("oc", "get", "consolequickstart", "install-app-and-associate-pipeline").Stdout()
	cmd.MustSucceed("oc", "get", "consolequickstart", "configure-pipeline-metrics").Stdout()
}

// Run opc with given arguments
func (opc Cmd) MustSucceed(args ...string) string {
	return opc.Assert(icmd.Success, args...)
}

// Run opc with given arguments
func (opc Cmd) Assert(exp icmd.Expected, args ...string) string {
	run := append([]string{opc.Path}, args...)
	output := cmd.Assert(exp, run...)
	return output.Stdout()
}

type CapturingPassThroughWriter struct {
	m   sync.RWMutex
	buf bytes.Buffer
	w   io.Writer
}

// NewCapturingPassThroughWriter creates new CapturingPassThroughWriter
func NewCapturingPassThroughWriter(w io.Writer) *CapturingPassThroughWriter {
	return &CapturingPassThroughWriter{
		w: w,
	}
}

func (w *CapturingPassThroughWriter) Write(d []byte) (int, error) {
	w.m.Lock()
	defer w.m.Unlock()
	w.buf.Write(d)
	return w.w.Write(d)
}

// Bytes returns bytes written to the writer
func (w *CapturingPassThroughWriter) Bytes() []byte {
	w.m.RLock()
	defer w.m.RUnlock()
	return w.buf.Bytes()
}

func StartPipeline(pipelineName string, params map[string]string, workspaces map[string]string, namespace string, args ...string) string {
	var commandArgs []string
	commandArgs = append(commandArgs, "opc", "pipeline", "start", pipelineName, "-o", "name", "-n", namespace)
	for key, value := range params {
		commandArgs = append(commandArgs, fmt.Sprintf("-p %s=%s", key, value))
	}
	for key, value := range workspaces {
		commandArgs = append(commandArgs, fmt.Sprintf("-w %s,%s", key, value))
	}
	commandArgs = append(commandArgs, args...)
	commandArgs = strings.Split(strings.Join(commandArgs, " "), " ")
	pipelineRunName := strings.Trim(cmd.MustSucceed(commandArgs...).Stdout(), "\n")
	log.Printf("Pipelinerun %s started", pipelineRunName)
	return pipelineRunName
}

// GetOpcPacInfoInstall fetches Pipelines as Code install information
func GetOpcPacInfoInstall() (*PacInfoInstall, error) {
	output := cmd.MustSucceed("opc", "pac", "info", "install").Stdout()
	lines := strings.Split(output, "\n")

	var pacInfo PacInfoInstall
	section := "" // current section: "pipelines"

	for _, rawLine := range lines {
		line := strings.TrimSpace(rawLine)
		if line == "" {
			continue
		}
		if line == "Pipelines as Code:" {
			section = "pipelines"
			continue
		}
		if section == "pipelines" {
			if strings.HasPrefix(line, "Install Version:") {
				pacInfo.PipelinesAsCode.InstallVersion = strings.TrimSpace(strings.TrimPrefix(line, "Install Version:"))
			} else if strings.HasPrefix(line, "Install Namespace:") {
				pacInfo.PipelinesAsCode.InstallNamespace = strings.TrimSpace(strings.TrimPrefix(line, "Install Namespace:"))
			}
		}
	}

	// Verify install version is not empty
	if pacInfo.PipelinesAsCode.InstallVersion == "" {
		return nil, fmt.Errorf("output of 'opc pac info install' is empty or missing Pipelines as Code information")
	}

	return &pacInfo, nil
}

// HubSearch performs an opc hub search for a resource
func HubSearch(resource string) error {
	output := cmd.MustSucceed("opc", "hub", "search", resource).Stdout()

	if !strings.Contains(output, resource) {
		log.Printf("Resource %q not found in opc hub search", resource)
		return fmt.Errorf("hub search failed for %s", resource)
	}
	return nil
}

// GetOpcPrList fetches pipeline run lists with status of each run
func GetOpcPrList(pipelineRunName, namespace string) ([]PipelineRunList, error) {
	result, err := VerifyResourceListMatchesName("pipelinerun", pipelineRunName, namespace)
	if err != nil {
		testsuit.T.Errorf("Failed to get pipelinerun list: %v", err)
	}
	output := strings.TrimSpace(result)
	lines := strings.Split(output, "\n")

	// Ensure output isn't empty
	if len(lines) < 2 {
		return nil, fmt.Errorf("unexpected pipelinerun output %s", output)
	}

	var runs []PipelineRunList
	for _, line := range lines[1:] { // Skip header
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 2 {
			log.Printf("Skipping malformed row: %s", line)
			continue
		}

		run := PipelineRunList{
			Name:   fields[0],
			Status: fields[len(fields)-1],
		}
		runs = append(runs, run)
	}

	return runs, nil
}

// resourceExists checks if a resource exists in output
func resourceExists(output, resourceName string) bool {
	trimmedOutput := strings.TrimSpace(output)
	if strings.HasPrefix(trimmedOutput, "No") {
		return false
	}

	lines := strings.Split(trimmedOutput, "\n")
	resourceLines := lines[1:] // Skip header line

	for _, line := range resourceLines {
		// Trim spaces from each line
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue // Skip empty lines
		}
		// Split the line into fields
		fields := strings.Fields(trimmed)
		if len(fields) > 0 && fields[0] == resourceName {
			return true
		}
	}
	return false
}

func VerifyResourceListMatchesName(resourceType, name, namespace string) (string, error) {
	output := cmd.MustSucceed("opc", resourceType, "list", "-n", namespace).Stdout()
	if !resourceExists(output, name) {
		return "", fmt.Errorf("%s %q not found in namespace %q", resourceType, name, namespace)
	}
	return output, nil
}
