package tasks

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// Verifies the usage of configMaps under defined Tasks
func TestConfigMapTasks(t *testing.T) {
	Convey("Given that Operator is installed", t, func() {
		Convey("When I am logged in as a non-admin user", func() {
			Convey("Create a Git Checkout Task", func() {
				Convey("By configure configMap to Task, with Revision Data", nil)
			})
			Convey("And Then I should run Task", nil)
			Convey("Then Validate the status of TaskRun", nil)
		})
	})
}
