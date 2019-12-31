package tasks

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// verify a very simple "hello world" TaskRun failure
// execution lead to the correct TaskRun status.
func TestTaskRunStatus(t *testing.T) {
	Convey("Given that Operator is installed", t, func() {
		Convey("When I logged in as non-admin user", func() {
			Convey("When I create a namespace", func() {
				Convey("I should create a task", nil)
				Convey("Then Run Task with `Nonexistance` service Account", nil)
				Convey("Then validate correct status/Reason for TaskRun failure", nil)
			})
		})
	})

}
