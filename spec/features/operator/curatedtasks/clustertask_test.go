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

func TestClusterTaskCreation(t *testing.T) {
	Convey("Given the Operator is installed", t, func() {
		Convey("When I logged in as non-admin user", func() {
			Convey("When I create Reusable ClusterTask", func() {
				Convey("Then validate existence of Reusable Task, clusterwide", nil)
				Convey("Create a random namespace", func() {
					Convey("I should able to run clusterTask", nil)
				})
			})
		})
		Convey("When I logged in as admin user", func() {
			Convey("When I create Reusable ClusterTask", func() {
				Convey("Then validate existence of Reusable Task, clusterwide", nil)
				Convey("Create a random namespace", func() {
					Convey("I should able to run clusterTask", nil)
				})
			})
		})
	})
}
