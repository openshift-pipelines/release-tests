package pipelines

import (
	"log"

	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/openshift-pipelines/release-tests/pkg/store"
)

var _ = gauge.Step("Start pipleine using tkn", func() {
	result := store.Tkn().MustSucceed(
		"pipeline", "start", "output-pipeline",
		"-r=source-repo=skaffold-git",
		"--showlog", "true",
		"-n", store.Namespace(),
	)

	log.Printf("output: %s", result)
})
