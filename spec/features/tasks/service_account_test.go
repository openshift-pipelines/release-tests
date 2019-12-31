package tasks

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTaskRunWithDefaultSA(t *testing.T) {
	Convey("Given that Operator is installed", t, func() {
		Convey("When I logged in as non-admin user", func() {
			Convey("When I create a namespace", func() {
				Convey("I should create Task", nil)
				Convey("And Then I should run Task with SA as `Default`", nil)
				Convey("Then validate status of TaskRun", nil)
			})
		})
	})
}
