package tasks

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTaskRunTimeOut(t *testing.T) {
	Convey("Given that Operator is installed", t, func() {
		Convey("When I logged in as non-admin user", func() {
			Convey("When I create a namespace", func() {
				Convey("I should Create a task which runs for atleast 10 sec", nil)
				Convey("Then run above task with a timeout of 5sec", nil)
				Convey("Then verify taskRun status is, `false`", nil)
				Convey("Then verify reason for failure is, `TimeOut`", nil)
			})
		})
	})
}
