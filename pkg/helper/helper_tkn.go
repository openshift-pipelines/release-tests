package helper

import (
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime/debug"
	"sync"

	"github.com/Netflix/go-expect"
	goexpect "github.com/Netflix/go-expect"
)

type TknRunner struct {
	// path to tkn binary
	Path string
}

// NewTknRunner initializes new OcRunner
func NewTknRunner(tknPath string) TknRunner {
	return TknRunner{
		Path: tknPath,
	}
}

// Run tkn with given arguments
func (tkn *TknRunner) Run(args ...string) string {
	output := CmdShouldPass(tkn.Path, args...)
	return output
}

// Prompt provides test utility for prompt test.
type Prompt struct {
	CmdArgs   []string
	Procedure func(*goexpect.Console) error
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
func (tkn *TknRunner) RunInteractiveTests(namespace string, ops *Prompt) *expect.Console {

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
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	if errStdout != nil || errStderr != nil {
		log.Fatal("failed to capture stdout or stderr\n")
	}

	return c
}
