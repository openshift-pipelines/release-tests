package pipelines

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

//Some times Pipeline run gets failed due to network issue, flaky task ,missing dependencies or upload problems etc.,
// Any of those issues results as False (corev1.ConditionFalse) within the PipelineRun Status
// retries attribute declaration, pipeline should be retried in case of failure
func TestRetriesOnPipelineRun(t *testing.T) {
	Convey("Given that Operator is installed", t, func() {
		Convey("When I logged in as non-admin user", func() {
			Convey("When I create a random namespace", func() {
				Convey("I should create a flaky Task", nil)
				Convey("Then I should create a pipeline, By specifying `retries` attribute value to `1`", nil)
				Convey("And Then I should Run Pipeline", nil)
				Convey("If status of pipelineRun, results as `False`", func() {
					Convey("Then Pipeline should Re-Run ", nil)
					Convey("Then validate status of Pipeline Run after Retry", nil)
				})
				Convey("If status of pipeline Run, results as `True`", func() {
					Convey("Then It should not Re-Run pipeline", nil)
					Convey("Then validate status of Pipeline Run", nil)
				})
			})
		})
	})
}
