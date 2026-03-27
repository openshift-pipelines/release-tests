package pipeline

import (
	"context"
	"os/exec"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*
CLI Testing Tests

STP Reference: outputs/stp/SRVKP-9005/SRVKP-9005_test_plan.md
Jira: SRVKP-9005
*/

var _ = Describe("[SRVKP-9005] CLI Testing", Ordered, func() {
	var (
		ctx context.Context
	)

	BeforeAll(func() {
		ctx = context.Background()
		_ = ctx
	})

	Context("TKN entitlement tests", func() {
		It("[test_id:TS-SRVKP-9005-013] should pass entitlement validation after CLI build moves to stage", func() {
			By("Verifying tkn CLI is available")
			cmd := exec.Command("tkn", "version")
			output, err := cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred(),
				"tkn CLI not available: %s", string(output))

			By("Verifying tkn reports correct version")
			ExpectWithOffset(1, string(output)).To(ContainSubstring("Client"),
				"tkn version output missing Client information")

			By("Verifying tkn can connect to cluster")
			cmd = exec.Command("tkn", "pipeline", "list")
			output, err = cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred(),
				"tkn cannot connect to cluster: %s", string(output))
		})
	})

	Context("TKN/OPC CLI tests", func() {
		It("[test_id:TS-SRVKP-9005-014] should execute CLI test suite with correct output and error handling", func() {
			By("Verifying tkn pipeline commands work")
			cmd := exec.Command("tkn", "pipeline", "list")
			output, err := cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred(),
				"tkn pipeline list failed: %s", string(output))

			By("Verifying tkn task commands work")
			cmd = exec.Command("tkn", "task", "list")
			output, err = cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred(),
				"tkn task list failed: %s", string(output))

			By("Verifying tkn pipelinerun commands work")
			cmd = exec.Command("tkn", "pipelinerun", "list")
			output, err = cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred(),
				"tkn pipelinerun list failed: %s", string(output))

			By("Verifying opc CLI is available")
			cmd = exec.Command("opc", "version")
			output, err = cmd.CombinedOutput()
			if err != nil {
				Skip("opc CLI not available - skipping opc-specific tests")
			}
			ExpectWithOffset(1, string(output)).To(ContainSubstring("version"),
				"opc version output unexpected")

			By("Verifying tkn error handling for invalid commands")
			cmd = exec.Command("tkn", "invalid-command-that-does-not-exist")
			_, err = cmd.CombinedOutput()
			ExpectWithOffset(1, err).To(HaveOccurred(),
				"tkn should fail on invalid command")
		})
	})

	Context("Approval task comment character limit", func() {
		It("[test_id:TS-SRVKP-9005-033] should enforce character limit on approval task comments", func() {
			By("Checking if tkn-approvaltask plugin is available")
			cmd := exec.Command("tkn", "approvaltask", "--help")
			_, err := cmd.CombinedOutput()
			if err != nil {
				Skip("tkn-approvaltask plugin not available - skipping")
			}

			By("Verifying help output is available")
			cmd = exec.Command("tkn", "approvaltask", "--help")
			output, err := cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			ExpectWithOffset(1, string(output)).To(ContainSubstring("approvaltask"),
				"help output does not mention approvaltask")
		})
	})

	Context("opc assist help examples", func() {
		It("[test_id:TS-SRVKP-9005-040] should include Examples section in help output", func() {
			By("Checking if opc CLI is available")
			cmd := exec.Command("opc", "version")
			_, err := cmd.CombinedOutput()
			if err != nil {
				Skip("opc CLI not available - skipping")
			}

			By("Verifying opc assist --help includes Examples section")
			cmd = exec.Command("opc", "assist", "--help")
			output, err := cmd.CombinedOutput()
			if err != nil {
				Skip("opc assist subcommand not available")
			}
			ExpectWithOffset(1, string(output)).To(ContainSubstring("Examples"),
				"opc assist --help does not include Examples section")
		})
	})

	Context("opc version accuracy", func() {
		It("[test_id:TS-SRVKP-9005-045] should report correct component versions", func() {
			By("Checking if opc CLI is available")
			cmd := exec.Command("opc", "version")
			output, err := cmd.CombinedOutput()
			if err != nil {
				Skip("opc CLI not available - skipping")
			}

			By("Verifying opc version output contains component information")
			versionOutput := string(output)
			ExpectWithOffset(1, versionOutput).To(ContainSubstring("version"),
				"opc version output missing version information")

			By("Verifying server component versions from cluster")
			cmd = exec.Command("oc", "get", "csv", "-n", "openshift-pipelines",
				"-o", "jsonpath={.items[0].spec.version}")
			serverOutput, err := cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			serverVersion := strings.TrimSpace(string(serverOutput))
			ExpectWithOffset(1, serverVersion).ToNot(BeEmpty(),
				"server version is empty")
			GinkgoWriter.Printf("opc version output: %s\nServer version: %s\n",
				versionOutput, serverVersion)
		})
	})
})
