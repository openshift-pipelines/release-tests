package curatedtasks

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestClusterTasksValidation(t *testing.T) {
	Convey("Given the Operator is installed", t, func() {
		Convey("When I logged in as non-admin user", func() {
			Convey("Then I should be able to perform lint test on ClusterTask manifest yamls", nil)
		})
		Convey("When I logged in as admin user", func() {
			Convey("Then I should be able to perform lint test on ClusterTask manifest yamls", nil)
		})
	})
}
