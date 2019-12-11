package pipelines

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/openshift-pipelines/release-tests/pkg/config"
	"github.com/openshift-pipelines/release-tests/pkg/helper"
	"github.com/openshift-pipelines/release-tests/pkg/pipelines"
	flags "github.com/openshift-pipelines/release-tests/spec/flags"
)

var _ = Describe("Validate PipelineRun, TaskRun status", func() {

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

	Describe("operator installed on cluster", func() {
		When("I create Task", func() {
			It("I should Run Task with `non-existence` SA", func() {
				By("Create Task in namespace " + flags.Namespace)
				pipelines.CreateTask(flags.Clients, flags.Namespace)
				By("I should Run Task")
				pipelines.CreateTaskRunWithSA(flags.Clients, flags.Namespace, "non-existence")
				By("validate TaskRun for failed status")
				pipelines.ValidateTaskRunForFailedStatus(flags.Clients, flags.Namespace)
			})
		})
		When("I create pipeline", func() {
			It("I should Run pipeline with `non-existence` SA", func() {
				By("Create Pipeline in namespace " + flags.Namespace)
				pipelines.CreatePipeline(flags.Clients, flags.Namespace)
				By("I should Run Pipeline")
				pipelines.CreatePipelineRunWithSA(flags.Clients, flags.Namespace, "non-existence")
				By("validate PipelineRun for failed status")
				pipelines.ValidatePipelineRunForFailedStatus(flags.Clients, flags.Namespace)
			})
		})
	})
})
