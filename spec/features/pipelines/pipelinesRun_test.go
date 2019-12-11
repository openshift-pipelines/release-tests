package pipelines

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	"github.com/openshift-pipelines/release-tests/pkg/helper"
	"github.com/openshift-pipelines/release-tests/pkg/pipelines"
	flags "github.com/openshift-pipelines/release-tests/spec/flags"
)

var _ = Describe("Run sample pipeline", func() {

	// Setup up state for each test spec
	// create new project (not set as active) and new context directory for each test spec
	// This is before every spec (It)
	BeforeEach(func() {
		SetDefaultEventuallyTimeout(config.Timeout)
		SetDefaultConsistentlyDuration(config.ConsistentlyDuration)
		flags.Clients, flags.Namespace, flags.Cleanup = helper.NewClientSet()
	})

	AfterEach(func() {
		defer flags.Cleanup()
	})

	Describe("I should create Random namespace", func() {
		When("I Create a sample Pipeline", func() {
			It("I should Run Pipeline", func() {
				By("Create simple pipeline under namespace " + flags.Namespace)
				pipelines.CreateSamplePipeline(flags.Clients, flags.Namespace)
				By("I should Run pipeline " + flags.Namespace)
				pipelines.RunSamplePipeline(flags.Clients, flags.Namespace)
			})
		})
	})
})
