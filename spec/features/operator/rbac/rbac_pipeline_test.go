package rbac

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPipelineSA(t *testing.T) {
	Convey("Given the Operator is installed", t, func() {
		Convey("When I create a new namespace", func() {
			Convey("It should create a service-account named pipeline", nil)
			Convey("The pipeline serviceaccount must have edit role", nil)
		})
	})
}

func TestPipelineRunSA(t *testing.T) {
	Convey("Given the Operator is installed", t, func() {
		Convey("When I create a new namespace", func() {
			Convey("I should be able to run pipeline, it should have SA as 'pipelines', by default", func() {
				Convey("I should be able to Create pipeline", nil)
				Convey("Then invoke pipelineRun without specifying SA", nil)
			})
		})
	})
}
