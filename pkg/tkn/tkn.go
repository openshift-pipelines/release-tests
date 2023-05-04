package tkn

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/Netflix/go-expect"
	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/cmd"
	"gotest.tools/v3/icmd"
)

type Cmd struct {
	// path to tkn binary
	Path string
}

// New initializes Cmd
func New(tknPath string) Cmd {
	return Cmd{
		Path: tknPath,
	}
}

// Verify the versions of Openshift Pipelines components
func AssertComponentVersion(version string, component string) {
	var commandResult string
	switch component {
	case "pipeline", "triggers", "operator":
		commandResult = cmd.MustSucceed("tkn", "version", "--component", component).Stdout()
	case "OSP":
		commandResult = cmd.MustSucceed("oc","get","tektonconfig","config", "-o", "jsonpath={.status.version}").Stdout()
	case "pipelines-as-code":
		commandResult = cmd.MustSucceed("oc", "get", "pac", "pipelines-as-code", "-o", "jsonpath={.status.version}").Stdout()
	default:
		testsuit.T.Errorf("Unknown component")
	}
	fmt.Println(commandResult)
	if !strings.Contains(commandResult, version) {
		testsuit.T.Errorf(component + " has an unexpected version: " + commandResult + " expected version is: " + version)
	}
}

func DownloadCLIFromCluster() {
	var architecture = strings.Trim(cmd.MustSucceed("uname").Stdout(), "\n") + " " + strings.Trim(cmd.MustSucceed("uname", "-m").Stdout(), "\n")
	var cliDownloadURL = cmd.MustSucceed("oc", "get", "consoleclidownloads", "tkn", "-o", "jsonpath={.spec.links[?(@.text==\"Download tkn and tkn-pac for "+architecture+"\")].href}").Stdout()
	cmd.MustSuccedIncreasedTimeout(time.Minute*10, "curl", "-o", "/tmp/tkn-binary.tar.gz", "-k", cliDownloadURL)
	cmd.MustSucceed("tar", "-xf", "/tmp/tkn-binary.tar.gz", "-C", "/tmp")
}

func AssertClientVersion(binary string) {
	var commandResult string
	var unexpectedVersion string
	if binary == "tkn-pac" {
		commandResult = cmd.MustSucceed("/tmp/tkn-pac", "version").Stdout()
		expectedVersion := os.Getenv("PAC_VERSION")
		if !strings.Contains(commandResult, expectedVersion) {
			testsuit.T.Errorf("tkn-pac has an unexpected version: " + commandResult + ". Expected: " + expectedVersion)
		}
	} else if binary == "tkn" {
		expectedVersion := os.Getenv("TKN_CLIENT_VERSION")
		commandResult = cmd.MustSucceed("/tmp/tkn", "version").Stdout()
		var splittedCommandResult = strings.Split(commandResult, "\n")
		for i, _ := range splittedCommandResult {
			if strings.Contains(splittedCommandResult[i], "Client") {
				if !strings.Contains(splittedCommandResult[i], expectedVersion) {
					unexpectedVersion = splittedCommandResult[i]
					testsuit.T.Errorf("tkn client has an unexpected version: " + unexpectedVersion + " Expeced: " + expectedVersion)
				}
			}
		}
	} else if binary == "opc" {
		commandResult = cmd.MustSucceed("/tmp/opc", "version").Stdout()
		components := [3]string{"OpenShift Pipelines Client", "Tekton CLI", "Pipelines as Code CLI"}
		expectedVersions := [3]string{os.Getenv("OSP_VERSION"), os.Getenv("TKN_CLIENT_VERSION"), os.Getenv("PAC_VERSION")}
		var splittedCommandResult = strings.Split(commandResult, "\n")
		for i := 0; i < 3; i++ {
			if strings.Contains(splittedCommandResult[i], components[i]) {
				if !strings.Contains(splittedCommandResult[i], expectedVersions[i]) {
					unexpectedVersion = splittedCommandResult[i]
					testsuit.T.Errorf(components[i] + " has an unexpected version: \"" + unexpectedVersion + "\". Expected: " + expectedVersions[i])
				}
			}
		}
	} else {
		testsuit.T.Errorf("Unknown binary or client")
	}
}

// Run tkn with given arguments
func (tkn Cmd) MustSucceed(args ...string) string {
	return tkn.Assert(icmd.Success, args...)
}

// Run tkn with given arguments
func (tkn Cmd) Assert(exp icmd.Expected, args ...string) string {
	run := append([]string{tkn.Path}, args...)
	output := cmd.Assert(exp, run...)
	return output.Stdout()
}

// Prompt provides test utility for prompt test.
type Prompt struct {
	CmdArgs   []string
	Procedure func(*expect.Console) error
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

func testCloser(closer io.Closer) {
	if err := closer.Close(); err != nil {
		log.Fatalf("Close failed: %s", err)
		debug.PrintStack()
	}
}

// Helps to Run Interactive Session
func (tkn *Cmd) RunInteractiveTests(namespace string, ops *Prompt) *expect.Console {

	c, err := expect.NewConsole(expect.WithStdout(os.Stdout))
	if err != nil {
		log.Fatal(err)
	}
	defer testCloser(c)

	cmd := exec.Command(tkn.Path, ops.CmdArgs[0:len(ops.CmdArgs)]...) //nolint:gosec
	cmd.Stdin = c.Tty()
	var errStdout, errStderr error
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()

	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := ops.Procedure(c); err != nil {
			log.Printf("procedure failed: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		_, errStdout = io.Copy(NewCapturingPassThroughWriter(c.Tty()), stdoutIn)
		_, errStderr = io.Copy(NewCapturingPassThroughWriter(c.Tty()), stderrIn)

	}()

	wg.Wait()

	err = cmd.Wait()
	if err != nil {
		testsuit.T.Errorf("cmd.Run() failed with %s\n", err)
	}
	if errStdout != nil || errStderr != nil {
		testsuit.T.Errorf("failed to capture stdout or stderr\n")
	}

	return c
}

func StartPipeline(pipelineName string, params map[string]string, workspaces map[string]string, namespace string, args ...string) string {
	var commandArgs []string
	commandArgs = append(commandArgs, "tkn", "pipeline", "start", pipelineName, "-o", "name", "-n", namespace)
	for key, value := range params {
		commandArgs = append(commandArgs, fmt.Sprintf("-p %s=%s", key, value))
	}
	for key, value := range workspaces {
		commandArgs = append(commandArgs, fmt.Sprintf("-w %s,%s", key, value))
	}
	for _, arg := range args {
		commandArgs = append(commandArgs, arg)
	}
	commandArgs = strings.Split(strings.Join(commandArgs, " "), " ")
	pipelineRunName := strings.Trim(cmd.MustSucceed(commandArgs...).Stdout(), "\n")
	log.Printf("Pipelinerun %s started", pipelineRunName)
	return pipelineRunName
}
