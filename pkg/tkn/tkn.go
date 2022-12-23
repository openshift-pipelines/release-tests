package tkn

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"runtime/debug"
	"strings"
	"sync"

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

func StartPipeline(pipelineName string, params map[string]string, workspaces map[string]string, args ...string) string {
	var commandArgs []string
	commandArgs = append(commandArgs, "tkn", "pipeline", "start", pipelineName, "-o", "name")
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

var Tkn_version_map map[string]string

func UpdateVersionsMap() {
	out, err := exec.Command("tkn", "version").Output()
	// tkn_version_map := map[string]string{}
	if err != nil {
		testsuit.T.Fail(fmt.Errorf("Error getting tkn version : %s", err))
	} else {
		temp_string_list := []string{}
		re := regexp.MustCompile(`([a-zA-Z ]*)\s*:\s*v*([0-9\.]*)`)
		op_temp := re.FindAllString(strings.ToLower(string(out)), -1)
		for _, s1 := range op_temp {
			temp_string_list = strings.Split(s1, `:`)
			Tkn_version_map[strings.Trim(strings.Split(temp_string_list[0], ` `)[0], ` `)] = strings.Trim(strings.ReplaceAll(temp_string_list[1], `v`, ``), ` `)
		}
	}
}
