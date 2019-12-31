package pipelines

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// can use a bucket for temporary storage of artifacts shared between tasks
func TestStorageBucketPipelineRun(t *testing.T) {
	Convey("Given that Operator is installed", t, func() {
		Convey("When I logged in as non-admin user", func() {
			Convey("When I create a namespace", func() {
				Convey("I should create GCP Bucket Secret", nil)
				Convey("I should create GCP Bucket Task", nil)
				Convey("Then I should Run Bucket Task", nil)
				Convey("And Then I should update original ConfigMap data to newly created GCP Bucket data", nil)
				Convey("I should Create a Task1 which creates Artifact", nil)
				Convey("I should Create a Task2 Which uses Artifact created by Task1", nil)
				Convey("I Should create pipeline that refers to above Tasks", nil)
				Convey("I should Run Pipeline", func() {
					Convey("Then Verify artifact created by Task1 should be mounted to GCP Bucket", nil)
					Convey("And Then Verify artifact used by Task2 should be mounted to GCP Bucket", nil)
				})
				Convey("Then verify status of pipelineRun", nil)
			})

		})
	})
}
