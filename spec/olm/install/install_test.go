package install

import (
	. "github.com/onsi/ginkgo"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	"github.com/openshift-pipelines/release-tests/pkg/helper"
	"github.com/openshift-pipelines/release-tests/pkg/operator"
	flags "github.com/openshift-pipelines/release-tests/spec/flags"
)

var _ = Describe("Olm installation test", func() {

	Describe("New ocp  Cluster", func() {
		When("I subscribed to pipeline Operator on canary channel", func() {
			It("I should Wait for Cluster CR availability", func() {
				helper.WaitForClusterCR(flags.Clients, config.ClusterCRName)
			})

			It("I should validate SCC", func() {
				operator.ValidateSCC(flags.Clients)
			})

			It("installs Pipelines 0.9", func() {
				By("Validate pipelines deployment into target namespace (openshift-pipelines)")
				operator.ValidatePipelineDeployments(flags.Clients)

				By("Validate pipeline version")
				operator.VerifyPipelineVersion(flags.Clients, flags.PipelineVersion)

			})

			It("installs Triggers 0.1", func() {
				operator.ValidateTriggerDeployments(flags.Clients)
			})

			It("I should validate status of operator", func() {
				operator.ValidateOperatorInstalledStatus(flags.Clients)
			})
		})
	})
})
