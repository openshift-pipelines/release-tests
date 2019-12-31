package pipelines

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPipelineRunWithResourceSpec(t *testing.T) {
	Convey("Given that Operator is installed", t, func() {
		Convey("When I logged in as non-admin user", func() {
			Convey("When I create a namespace", func() {
				Convey("I should Create a Task", nil)
				Convey("I should create a pipeline", nil)
				Convey("Then I should able to Run pipeline with embeded ResourceSpec", nil)
				Convey("Then I should validate the status of pipelineRun", nil)
			})
		})
	})
}
