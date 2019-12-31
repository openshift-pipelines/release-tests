package tasks

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// TestSideCarTaskRun injects side car to pod as local registry where we can push build image to local registry
func TestSideCarTaskRun(t *testing.T) {
	Convey("Given that Operator is installed", t, func() {
		Convey("When I logged in as non-admin user", func() {
			Convey("Then I should create a Task", func() {
				Convey("By configuring side car to Task", nil)
			})
			Convey("And Then I should run Task", nil)
			Convey("Then validate status of TaskRun", nil)
		})
	})
}
