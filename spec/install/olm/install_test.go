package olm

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFreshInstall(t *testing.T) {
	Convey("Given a new cluster", t, func() {
		Convey("When I subscribe to the Pipelines Operator", func() {

			Convey("It should install the following", func() {

				Convey("installs Pipelines 0.8", nil)
				Convey("installs Triggers 0.1", nil)

				Convey("installs the following cluster tasks", func() {
					Convey("s2i", nil)
					Convey("s2i-java-8", nil)
					Convey("s2i-java-11", nil)
					Convey("s2i-python-2", nil)
					Convey("s2i-python-3", nil)
					Convey("openshift-client", nil)
				})
			})

		})
	})
}
