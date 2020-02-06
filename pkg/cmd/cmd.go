package cmd

import (
	"fmt"

	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	"gotest.tools/assert"
	"gotest.tools/v3/icmd"
)

type Cmd struct {
	Args     []string
	Expected icmd.Expected
}

// testsuitAdaptor bridges the gap between testsuit.T and assert.TestingT
type testsuitAdaptor struct{}

// ensure testsuitAdaptor satisfies assert.TestingT interface
var _ assert.TestingT = (*testsuitAdaptor)(nil)

func (ta testsuitAdaptor) Fail() {
	testsuit.T.Fail(fmt.Errorf("Step failed execute"))
}

func (ta testsuitAdaptor) FailNow() {
	testsuit.T.Fail(fmt.Errorf("Step failed to execute"))
}

func (ta testsuitAdaptor) Log(args ...interface{}) {
	testsuit.T.Fail(fmt.Errorf("%v", args))
}

func Run(cmd []string) *icmd.Result {
	return icmd.RunCmd(icmd.Cmd{Command: cmd, Timeout: config.Timeout})
}

// AssertOutput runs a command and verfies exit code (0)
func AssertOutput(cmd *Cmd) *icmd.Result {
	res := Run(cmd.Args)

	t := &testsuitAdaptor{}
	res.Assert(t, cmd.Expected)
	return res
}
