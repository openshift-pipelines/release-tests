package uninstall

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUninstall(t *testing.T) {
	Convey("Given a new cluster with Operator installed", t, func() {
		Convey("When I logged in as admin user", func() {
			Convey("I should Delete instance (name: cluster) of config.operator.tekton.dev", func() {
			})
			Convey("I should Delete Clusterserviceversion", func() {
			})
			Convey("I should Delete install plan", func() {
			})
			Convey("I should delete subscription", func() {
			})
			Convey("Validate cluster should not have any CRDs/api-resources which contains `tekton`", func() {
			})
		})
	})
}
