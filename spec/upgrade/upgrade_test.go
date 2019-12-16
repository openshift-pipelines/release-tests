package upgrade

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestManualApprovalStrategy(t *testing.T) {

	Convey("Given a new cluster with Operator installed with installation strategy 'manual'", t, func() {
		Convey("When I Logged as an admin user", func() {
			Convey("When new upgrade is pushed to channel", func() {
				Convey("When I Navigate to opeartor > installed operator and select namespace to 'openshift-operators'", nil)
				Convey("Then I should see 'openshift pipelines operator' status as 'Upgrade pending'", func() {
					Convey("Then click on 'openshift pipelines operator', it should show upgrade status as 'Upgrade requires approval'", nil)
					Convey("Click on 'requires approval' > select components ", func() {
						Convey("It should list all latest components csv, crd, SA, ClusterRole, ClusterBinding, etc.,", nil)
						Convey("Click on 'Approve' button it should get updated", nil)
					})
				})
			})
		})
	})

}
