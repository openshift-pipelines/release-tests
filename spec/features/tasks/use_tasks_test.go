package tasks

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUseClusterTasks(t *testing.T) {
	Convey("Given that Operator is installed", t, func() {

		Convey("When I am logged in as a non-admin user", func() {
			Convey("I can list all the cluster tasks", nil)
			Convey("I should be able to create a Pipeline that uses a cluster task", nil)
		})

	})
}
