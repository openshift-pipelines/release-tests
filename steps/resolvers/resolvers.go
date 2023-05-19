package resolvers

import (
	"github.com/getgauge-contrib/gauge-go/gauge"
	"github.com/openshift-pipelines/release-tests/pkg/oc"
)



var _ = gauge.Step("Delete projects", func(){
	oc.DeleteProject("resolver-test-tasks")
	oc.DeleteProject("resolver-test-pipelines")
	oc.DeleteProject("resolver-test-pipelineruns")
})


