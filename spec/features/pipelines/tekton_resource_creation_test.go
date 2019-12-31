package pipelines

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTektonResourceCreation(t *testing.T) {
	Convey("Given that Operator is installed", t, func() {
		Convey("When I logged in as non-admin user", func() {
			Convey("Create a namespace-1", func() {
				Convey("Then validate Creation of tekton Resources", func() {
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
			Convey("Create a namespace-2", func() {
				Convey("Then validate Creation of tekton Resources using TKN cli", func() {
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
	})
}
