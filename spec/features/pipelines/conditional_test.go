package pipelines

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// Conditions can declare input PipelineResources via the resources field to provide
// the Condition container step with data or context that is needed to perform the check.
func TestConditionalPipelineRun(t *testing.T) {
	Convey("Given that Operator is installed", t, func() {
		Convey("When I logged in as non-admin user", func() {
			Convey("When I create a namespace", func() {
				Convey("I should create a `conditional` Resource to perform some check", func() {
					Convey("Then I should create a pipeline that refers to `Conditional` Resource", nil)
					Convey("And Then I should Run Pipeline", nil)
					Convey("And If `conditional` check is success, And returns exit code 0", func() {
						Convey("Then It should resume pipeline execution", nil)
						Convey("And Then I should validate the status of pipelines", nil)
					})
					Convey("And If `conditional` check  failed", func() {
						Convey("Then It should stop the execution of PipelineRun", nil)
						Convey("And Then It should Mark PipelineRun status as `Failure` ", nil)
					})
				})
			})
		})
	})
}
