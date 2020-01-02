package olm

import (
	"testing"

	"github.com/openshift-pipelines/release-tests/pkg/config"
	"github.com/openshift-pipelines/release-tests/pkg/operator"
	. "github.com/smartystreets/goconvey/convey"
)

func TestFreshInstall(t *testing.T) {
	operator.SubscribeToChannel()
	defer operator.DeleteOperator(t, operator.Clt)
	Convey("Given a new cluster", t, func() {
		Convey("When I subscribe to the Pipelines Operator", func() {
			Convey("validate Cluster CR", func() {
				operator.CR = operator.ValidateClusterCR(operator.Clt)
				So(true, ShouldEqual, true)
			})
			Convey("validate SCC", func() {
				operator.ValidateSCC(operator.Clt, operator.CR.Spec.TargetNamespace, config.PipelineControllerName)
				So(true, ShouldEqual, true)
			})
			Convey("installs Pipelines 0.9", func() {
				Convey("Validate pipelines deployment into target namespace (openshift-pipelines)", func() {
					operator.ValidatePipelineAndTriggerSetup(operator.Clt, operator.CR, config.PipelineWebhookName, config.PipelineControllerName)
					So(true, ShouldEqual, true)
				})

				Convey("Validate pipelines version", func() {
					operator.VerifyPipelineVersion(operator.Clt, `v0.9`)
				})
			})
			Convey("installs Triggers 0.1", func() {
				operator.ValidatePipelineAndTriggerSetup(operator.Clt, operator.CR, config.TriggerWebhookName, config.TriggerControllerName)
				So(true, ShouldEqual, true)
			})
			SkipConvey("installs the following cluster tasks", func() {
				Convey("s2i", func() {
				})
				Convey("s2i-java-8", func() {
				})
				Convey("s2i-java-11", func() {
				})
				Convey("s2i-python-2", func() {
				})
				Convey("s2i-python-3", func() {
				})
				Convey("openshift-client", func() {
				})
				So(true, ShouldEqual, true)
			})
		})
	})
}
