package tasks

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUseClusterTasks(t *testing.T) {
	Convey("Given that Operator is installed", t, func() {
		Convey("When I am logged in as a non-admin user", func() {
			Convey("I should be able to list all the cluster tasks", nil)
		})
		Convey("When I am logged in as a admin user", func() {
			Convey("I should be able to list all the cluster tasks", nil)
		})
	})
}

func TestPipelineRunForClusterTasks(t *testing.T) {
	Convey("Given that Operator is installed", t, func() {
		Convey("When I am logged in as a non-admin user", func() {
			Convey("When I create a pipeline using ClusterTask 'buildah-v0-8-0'", func() {
				Convey("Then I should be able to run pipeline successfully", nil)
			})
			Convey("When I create a pipeline using ClusterTask 'openshift-client-v0-8-0'", func() {
				Convey("Then I should be able to run pipeline successfully", nil)
			})
			Convey("When I create a pipeline using ClusterTask 's2i-go-v0-8-0'", func() {
				Convey("Then I should be able to run pipeline successfully", nil)
			})
			Convey("When I create a pipeline using ClusterTask 's2i-java-11-v0-8-0'", func() {
				Convey("Then I should be able to run pipeline successfully", nil)
			})
			Convey("When I create a pipeline using ClusterTask 's2i-java-8-v0-8-0'", func() {
				Convey("Then I should be able to run pipeline successfully", nil)
			})
			Convey("When I create a pipeline using ClusterTask 's2i-nodejs-v0-8-0'", func() {
				Convey("Then I should be able to run pipeline successfully", nil)
			})
			Convey("When I create a pipeline using ClusterTask 's2i-python-3-v0-8-0'", func() {
				Convey("Then I should be able to run pipeline successfully", nil)
			})
			Convey("When I create a pipeline using ClusterTask 's2i-v0-8-0'", func() {
				Convey("Then I should be able to run pipeline successfully", nil)
			})
		})
	})
}

func TestTaskRunForClusterTasks(t *testing.T) {
	Convey("Given that Operator is installed", t, func() {
		Convey("When I am logged in as a non-admin user", func() {
			Convey("Then I should create TaskRun using ClusterTask 'buildah-v0-8-0'", nil)
			Convey("Then I should create TaskRun using ClusterTask 'openshift-client-v0-8-0'", nil)
			Convey("Then I should create TaskRun using ClusterTask 's2i-go-v0-8-0'", nil)
			Convey("Then I should create TaskRun using ClusterTask 's2i-java-11-v0-8-0'", nil)
			Convey("Then I should create TaskRun using ClusterTask 's2i-java-8-v0-8-0'", nil)
			Convey("Then I should create TaskRun using ClusterTask 's2i-nodejs-v0-8-0'", nil)
			Convey("Then I should create TaskRun using ClusterTask 's2i-python-3-v0-8-0'", nil)
			Convey("Then I should create TaskRun using ClusterTask 's2i-v0-8-0'", nil)
		})
	})
}
