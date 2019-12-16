package triggers

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTriggersInstall(t *testing.T) {
	Convey("Given the Operator is installed", t, func() {
		Convey("It should have installed Trigger controllers as a deployment to target namespace (openshift-pipelines)", nil)
		Convey("It should have installed Trigger Webhooks as a deployment to target namespace (openshift-pipelines)", nil)
	})
}
