package pipelines

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPipelineRunWithPipelineSpec(t *testing.T) {
	Convey("Given that Operator is installed", t, func() {
		Convey("When I logged in as non-admin user", func() {
			Convey("When I create a namespace", func() {
				Convey("I should Create a Task", nil)
				Convey("I should be able to run pipeline with embeded pipelinespec", nil)
				Convey("Then I should validate the status of pipelineRun", nil)
			})
		})
	})
}
