package cmd

import (
	"testing"

	"github.com/openshift-pipelines/release-tests/pkg/config"
	"gotest.tools/v3/icmd"
)

type Cmd struct {
	Args     []string
	Expected icmd.Expected
}

func Run(cmd []string) *icmd.Result {
	return icmd.RunCmd(icmd.Cmd{Command: cmd, Timeout: config.Timeout})
}

// AssertOutput runs a command and verfies exit code (0)
func AssertOutput(cmd *Cmd) *icmd.Result {
	res := Run(cmd.Args)

	// TODO: fix this hack
	var t = &testing.T{}
	res.Assert(t, cmd.Expected)
	return res
}
