package pipelines

import (
	"log"

	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/openshift-pipelines/release-tests/pkg/cmd"
	"github.com/openshift-pipelines/release-tests/steps"
	"gotest.tools/v3/icmd"
)

var _ = gauge.Step("Start pipleine using tkn", func() {
	result := cmd.AssertOutput(&cmd.Cmd{
		Args: []string{
			steps.Tkn().Path, "pipeline", "start", "output-pipeline",
			"-r=source-repo=skaffold-git",
			"--showlog", "true",
			"-n", steps.Namespace()},
		Expected: icmd.Success,
	})

	log.Printf("output: %s", result.Stdout())
})
