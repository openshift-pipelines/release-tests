package pipeline

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPipelinesInstall(t *testing.T) {
	Convey("Given the Operator is installed", t, func() {
		Convey("It should have installed pipelines controllers as a deployment in target namespace (openshift-pipelines)", nil)
		Convey("It should have installed pipelines Webhooks as a deployment in target namespace (openshift-pipelines)", nil)
	})
}
