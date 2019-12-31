package tasks

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// Task requires input resources or output resources, they must be provided to run the Task.

// They can be provided via references to existing PipelineResources or by embedding Resource spec

// Eg: resourceSpec:
// type: git
// params:
//   - name: url
// 	value: https://github.com/pivotal-nader-ziada/gohelloworld
func TestEmbedResourceSpec(t *testing.T) {
	Convey("Given that Operator is installed", t, func() {
		Convey("When I logged in as non-admin user", func() {
			Convey("I should Create a Task", func() {
				Convey("By embedding the specs of (i/p or o/p) resources to task", nil)
			})
			Convey("And Then I should Run Task", nil)
			Convey("Then I should validate status of TaskRun", nil)
		})
	})
}
