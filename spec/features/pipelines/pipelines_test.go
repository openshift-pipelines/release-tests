package pipeline

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPipelineRunTutorial(t *testing.T) {
	Convey("Given that Operator is installed", t, func() {
		Convey("I should be able to run Pipelines Tutorial as a non admin", nil)
	})
}

func TestPipelineResourceCreation(t *testing.T) {
	Convey("Given that Operator is installed", t, func() {
		Convey("When I create a namespace", func() {
			Convey("Then validate Creation of pipeline Resources", func() {
				Convey("I should be able to Create 'Task'", nil)
				Convey("I should be able to create 'pipeline'", nil)
				Convey("I should be able to Create 'pipelinesResources'", func() {
					Convey("I should be able to create 'Git' pipelineResource", nil)
					Convey("I should be able to create 'image' pipelineResource", nil)
					Convey("I should be able to create 'cluster' pipelineResource", nil)
					Convey("I should be able to create 'storage' pipelineResource", nil)
					Convey("I should be able to create 'pull request' pipelineResource", nil)
					Convey("I should be able to create 'cloud events' pipelineResource", nil)
				})
				Convey("I should be able to create 'Condition'", nil)
				Convey("I should be able to create 'pipelineRun'", nil)
			})
			Convey("Then validate status of created Resources", nil)
		})
	})
}
