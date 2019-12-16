package uninstall

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUninstall(t *testing.T) {
	Convey("Given a new cluster with Operator installed", t, func() {
		Convey("Delete instance of openshift pipelines config", nil) // $ oc delete config.operator.tekton.dev cluster
		//   config.operator.tekton.dev "cluster" deleted
		Convey("Uninstall openshift pipeline Operator from console(UI)", func() {
			Convey("Login as an admin user", nil)
			Convey("Navigate to operators > operatorHub and search for 'openshift pipelines operator'", nil)
			Convey("Then click on 'openshift pipelines operator' and click on uninstall button", nil)
		})
	})
}
