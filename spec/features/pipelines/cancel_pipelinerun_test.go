package pipelines

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// When you perform cancellation on Pipeline Run by updating it's sepc to `cancelled`
// Related TaskRun instances should be marked as cancelled and running Pods should be deleted.
func TestPipelineRunCancellation(t *testing.T) {
	Convey("Given that Operator is installed", t, func() {
		Convey("When I logged in as non-admin user", func() {
			Convey("When I create a namespace", func() {
				Convey("I should create pipeline", nil)
				Convey("And Then I should run pipeline", nil)
				Convey("If status of pipeline run is still Running", func() {
					Convey("Then I should be able to cancel running pipeline, by updating status spec to `PipelineRunCancelled`", nil)
					Convey("Then validate status of Related TaskRun instances also marked as `cancelled`", nil)
					Convey("Then validate the termination of running Pods", nil)
				})
				Convey("If status of pipeline run is completed", func() {
					Convey("And When I try to cancel running pipelines, by updating status spec to `PipelineRunCancelled`", nil)
					Convey("Then I should validate for error msg [Eg: PipelineRun cannot be cancelled]", nil)
				})
			})
		})
	})
}
