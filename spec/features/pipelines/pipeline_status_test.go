package pipelines

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// verify a very simple "hello world"  PipelineRun failure
// execution lead to the correct TaskRun status.
func TestPipelineRunStatus(t *testing.T) {

	Convey("Given that Operator is installed", t, func() {
		Convey("When I logged in as non-admin user", func() {
			Convey("When I create a namespace", func() {
				Convey("I should create a task", nil)
				Convey("I should create a pipeline that refers to above task", nil)
				Convey("Then Run Pipeline with `Nonexistance` service Account", nil)
				Convey("Then validate correct status/Reason for PipelineRun failure", nil)
			})
		})
	})
}
