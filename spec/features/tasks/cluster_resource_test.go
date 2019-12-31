package tasks

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestClusterResourceTask(t *testing.T) {
	Convey("Given that Operator is installed", t, func() {
		Convey("When I logged in as non-admin user", func() {
			Convey("When I create a namespace", func() {
				Convey("Then I should Create pipeline Resource of `Cluster` type", nil)
				Convey("Then Create a cluster Resource Task", func() {
					Convey("By defining, required Secrets and Configmaps to Task", nil)
				})
				Convey("And Then I should run Task successfully", nil)
			})

		})
	})
}
