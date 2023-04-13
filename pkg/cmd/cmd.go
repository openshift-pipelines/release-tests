package cmd

import (
	"fmt"
	"time"

	"github.com/getgauge-contrib/gauge-go/testsuit"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/icmd"
)

type FooCmd struct {
	Command  []string
	Expected icmd.Expected
}

// testsuitAdaptor bridges the gap between testsuit.T and assert.TestingT as
// testsuit.T does not implement assert.TestingT interface
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

func Run(cmd ...string) *icmd.Result {
	return icmd.RunCmd(icmd.Cmd{Command: cmd, Timeout: config.CLITimeout})
}

// MustSucceed asserts that the command ran with 0 exit code
func MustSucceed(args ...string) *icmd.Result {
	return Assert(icmd.Success, args...)
}

// Assert runs a command and verifies exit code (0)
func Assert(exp icmd.Expected, args ...string) *icmd.Result {
	res := Run(args...)
	t := &testsuitAdaptor{}
	res.Assert(t, exp)
	return res
}

func MustSuccedIncreasedTimeout(timeout time.Duration, args ...string) *icmd.Result {
    return AssertIncreasedTimeout(icmd.Success, timeout, args...)
}

func AssertIncreasedTimeout(exp icmd.Expected, timeout time.Duration, args ...string) *icmd.Result {
    res := RunIncreasedTimeout(timeout, args...)
    t := &testsuitAdaptor{}
    res.Assert(t, exp)
    return res
}

func RunIncreasedTimeout(timeout time.Duration, cmd ...string) *icmd.Result {
    return icmd.RunCmd(icmd.Cmd{Command: cmd, Timeout: timeout})
}


