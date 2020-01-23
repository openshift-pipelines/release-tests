package pipelines

import (
	"log"

	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/openshift-pipelines/release-tests/pkg/helper"
	"github.com/openshift-pipelines/release-tests/steps"
	"gotest.tools/v3/icmd"
)

var _ = gauge.Step("Start pipleine using tkn", func() {
	log.Printf("output: %s", helper.RunCmd(&helper.TknCmd{
		Args: []string{steps.GetTknBinaryPath().Path, "pipeline", "start", "output-pipeline",
			"-r=source-repo=skaffold-git",
			"--showlog",
			"true",
			"-n", steps.GetNameSpace()},
		Expected: icmd.Expected{
			ExitCode: 0,
			Err:      icmd.None,
		},
	}).Stdout())
})
