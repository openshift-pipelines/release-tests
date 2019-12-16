package clusterconfig

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestClusterConfigReinstall(t *testing.T) {
	Convey("Given the Operator is installed", t, func() {
		Convey("Delete cluster config from cluster", nil)
		Convey("Verify pipeline, Triggers, ClusterTasks has been uninstalled", nil)
		Convey("Create new config name 'Cluster' with target namespace 'openshift-pipelines'", func() {
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
