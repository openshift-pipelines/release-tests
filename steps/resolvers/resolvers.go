package resolvers

import (
	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/openshift-pipelines/release-tests/pkg/resolvers"
)

var _ = gauge.Step("Check if project <project> exists", func(project string){
	resolvers.CheckProjectExists(project)
})