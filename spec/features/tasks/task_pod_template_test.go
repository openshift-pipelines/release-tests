package tasks

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

//Specifies a subset of PodSpec configuration that will be used as the basis for the Task pod.
// Eg: https://github.com/tektoncd/pipeline/blob/master/docs/taskruns.md
func TestTaskRunPodTemplate(t *testing.T) {
	Convey("Given that Operator is installed", t, func() {
		Convey("When I logged in as non-admin user", func() {
			Convey("When I create a namespace", func() {
				Convey("Then Create a task", func() {
					Convey("I should run task as non-root user", func() {
						Convey("Then I should run task with SecurityContext(`runAsNonRoot: true`) defined under PodTemplate", nil)
					})
				})
			})
		})
	})
}
