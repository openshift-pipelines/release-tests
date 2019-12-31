package tasks

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// To validate termination of pods, if we perform cancellation on TaskRun
// and Right TaskRun status
func TestTaskRunCancellation(t *testing.T) {
	Convey("Given that Operator is installed", t, func() {
		Convey("When I logged in as non-admin user", func() {
			Convey("When I create a namespace", func() {
				Convey("I should create Task", nil)
				Convey("And Then I should run task", nil)
				Convey("If status of TaskRun is still Running", func() {
					Convey("Then I should be able to cancel running Task, by updating status spec to `TaskRunCancelled`", nil)
					Convey("Then validate the termination of running Pods", nil)
				})
				Convey("If status of TaskRun is completed", func() {
					Convey("And When I try to cancel running Task, by updating status spec to `TaskRunCancelled`", nil)
					Convey("Then I should validate for error msg [Eg: TaskRun cannot be cancelled]", nil)
				})

			})
		})
	})
}
