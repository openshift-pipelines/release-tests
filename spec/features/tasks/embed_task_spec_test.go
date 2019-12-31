package tasks

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// Embed the spec of `task` directly into `TaskRun`
// `Eg: spec:
//   taskSpec:
//    inputs:
// 	  	resources:
// 	  	- name: workspace
// 		  type: git
//    steps:
// 	  - name: build-and-push
// 	    image: gcr.io/kaniko-project/executor:v0.9.0
// `
func TestEmbedTaskSpec(t *testing.T) {
	Convey("Given that Operator is installed", t, func() {
		Convey("When I logged in as non-admin user", func() {
			Convey("When I create a namespace", func() {
				Convey("Then Embed the task spec into `TaskRun`", nil)
				Convey("And Then I should run TaskRun", nil)
				Convey("Then I should validate the status of TaskRun", nil)
			})
		})
	})
}
