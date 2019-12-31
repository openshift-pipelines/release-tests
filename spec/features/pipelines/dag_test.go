package pipelines

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// TestDAGPipelineRun creates a graph of arbitrary Tasks, then looks at the corresponding
// TaskRun start times to ensure they were run in the order intended, which is:
//                               |
//                        pipeline-task-1
//                       /               \
//   pipeline-task-2-parallel-1    pipeline-task-2-parallel-2
//                       \                /
//                        pipeline-task-3
//                               |
//                        pipeline-task-4
func TestDAGPipelineRun(t *testing.T) {
	Convey("Given that Operator is installed", t, func() {
		Convey("When I logged in as non-admin user", func() {
			Convey("Create a namespace", func() {
				Convey("Then I should Create multiple tasks", nil)
				Convey("And Then I should create a pipeline, By Specifing the order of execution to created tasks", nil)
				Convey("Then I should Run pipeline", nil)
				Convey("I should validate order of execution of tasks defined under pipeline", nil)
				Convey("Then Verify TaskRun start times to ensure they were run in the order intended", nil)
			})
		})
	})
}
