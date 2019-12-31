package tasks

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// TestDuplicatePodTaskRun creates 10 builds and checks that each of them has only one build pod.
func TestDuplicatePodTaskRun(t *testing.T) {
	Convey("Given that Operator is installed", t, func() {
		Convey("When I logged in as non-admin user", func() {
			Convey("When I create a namespace", func() {
				Convey("I should create a task", nil)
				Convey("And Then I should Run Task for (10 times) with same pod name", nil)
				Convey("Then I should validate status of TaskRun", nil)
				Convey("Then Verify Number of pods created should be equal to `1`", nil)
			})
		})
	})
}
