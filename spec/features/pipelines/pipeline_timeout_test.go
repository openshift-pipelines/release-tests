package pipelines

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPipelineRunTimeOut(t *testing.T) {
	Convey("Given that Operator is installed", t, func() {
		Convey("When I logged in as non-admin user", func() {
			Convey("When I create a namespace", func() {
				Convey("I should Create a task which runs for atleast 10 sec", nil)
				Convey("I should Create a pipeline that refers to above task", nil)
				Convey("Then run pipeline with a timeout of 5sec", nil)
				Convey("Then verify pipeline run status is, `false`", nil)
				Convey("Then verify reason for failure is, `TimeOut`", nil)
			})
		})
	})
}
