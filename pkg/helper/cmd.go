package helper

import (
	"testing"

	"github.com/openshift-pipelines/release-tests/pkg/config"
	"gotest.tools/v3/icmd"
)

var t = &testing.T{}

type TknCmd struct {
	Args     []string
	Expected icmd.Expected
}

func RunQuiet(cmd []string) *icmd.Result {
	return icmd.RunCmd(icmd.Cmd{Command: cmd, Timeout: config.Timeout})
}

// RunCmd runs a command and verfies exit code (0)
func RunCmd(cmd *TknCmd) *icmd.Result {
	res := RunQuiet(cmd.Args)
	AssertCmd(res, cmd.Expected)
	return res
}

func AssertCmd(res *icmd.Result, expected icmd.Expected) {
	res.Assert(t, expected)
}
