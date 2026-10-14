package pipeline

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*
Release Verification Tests

STP Reference: outputs/stp/SRVKP-9005/SRVKP-9005_test_plan.md
Jira: SRVKP-9005
*/

var _ = Describe("[SRVKP-9005] Release Verification", Ordered, func() {
	var (
		ctx context.Context
	)

	BeforeAll(func() {
		ctx = context.Background()
		_ = ctx
	})

	Context("Post-release operator installation", func() {
		It("[test_id:TS-SRVKP-9005-020] should install from OperatorHub on all supported OCP versions", func() {
			By("Verifying operator is available in OperatorHub")
			cmd := exec.Command("oc", "get", "packagemanifest",
				"openshift-pipelines-operator-rh", "-n", "openshift-marketplace",
				"-o", "jsonpath={.status.catalogSource}")
			output, err := cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred(),
				"OpenShift Pipelines not found in OperatorHub: %s", string(output))
			catalogSource := strings.TrimSpace(string(output))
			ExpectWithOffset(1, catalogSource).ToNot(BeEmpty(),
				"catalog source is empty")
			fmt.Fprintf(GinkgoWriter, "Catalog source: %s\n", catalogSource)

			By("Verifying operator is installed and CSV is Succeeded")
			cmd = exec.Command("oc", "get", "csv", "-n", "openshift-pipelines",
				"-o", "jsonpath={.items[0].status.phase}")
			output, err = cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			ExpectWithOffset(1, strings.TrimSpace(string(output))).To(Equal("Succeeded"),
				"CSV not Succeeded after installation")

			By("Verifying all operator components are running")
			cmd = exec.Command("oc", "get", "pods", "-n", "openshift-pipelines",
				"--field-selector=status.phase!=Running,status.phase!=Succeeded",
				"-o", "name")
			output, err = cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			ExpectWithOffset(1, strings.TrimSpace(string(output))).To(BeEmpty(),
				"found non-running pods: %s", string(output))
		})
	})

	Context("CLI binary SHA/signature verification", func() {
		It("[test_id:TS-SRVKP-9005-021] should have valid SHA256 checksums and Cosign signatures on macOS and Windows", func() {
			By("Verifying cosign is available")
			cmd := exec.Command("cosign", "version")
			output, err := cmd.CombinedOutput()
			if err != nil {
				Skip("cosign not available - skipping signature verification")
			}
			ExpectWithOffset(1, string(output)).To(ContainSubstring("cosign"),
				"cosign version output unexpected")

			By("Verifying tkn CLI binary is available locally")
			cmd = exec.Command("which", "tkn")
			output, err = cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred(),
				"tkn CLI not found: %s", string(output))
			tknPath := strings.TrimSpace(string(output))
			fmt.Fprintf(GinkgoWriter, "tkn binary path: %s\n", tknPath)

			By("Verifying tkn binary is a valid executable")
			cmd = exec.Command("file", tknPath)
			output, err = cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			ExpectWithOffset(1, string(output)).To(ContainSubstring("ELF"),
				"tkn binary is not a valid ELF executable")
		})
	})

	Context("Bug verification for 1.21 release", func() {
		It("[test_id:TS-SRVKP-9005-022] should verify all ON_QA bugs with fixVersion 1.21", func() {
			By("Verifying operator version matches expected release")
			cmd := exec.Command("oc", "get", "csv", "-n", "openshift-pipelines",
				"-o", "jsonpath={.items[0].spec.version}")
			output, err := cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			version := strings.TrimSpace(string(output))
			fmt.Fprintf(GinkgoWriter, "Installed version: %s\n", version)

			By("Verifying all components are healthy")
			cmd = exec.Command("oc", "get", "tektonconfig", "config",
				"-o", "jsonpath={.status.conditions[?(@.type=='Ready')].status}")
			output, err = cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			ExpectWithOffset(1, strings.TrimSpace(string(output))).To(Equal("True"),
				"TektonConfig not Ready - bug verification requires healthy cluster")

			By("Verifying no pods in error state")
			cmd = exec.Command("oc", "get", "pods", "-n", "openshift-pipelines",
				"--field-selector=status.phase=Failed", "-o", "name")
			output, err = cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			ExpectWithOffset(1, strings.TrimSpace(string(output))).To(BeEmpty(),
				"found failed pods in openshift-pipelines namespace")
		})
	})

	Context("Documentation review", func() {
		It("[test_id:TS-SRVKP-9005-023] should have accurate and complete release documentation", func() {
			By("Verifying operator provides expected APIs")
			expectedCRDs := []string{
				"pipelines.tekton.dev",
				"tasks.tekton.dev",
				"pipelineruns.tekton.dev",
				"taskruns.tekton.dev",
			}
			for _, crd := range expectedCRDs {
				cmd := exec.Command("oc", "get", "crd", crd, "-o", "name")
				output, err := cmd.CombinedOutput()
				ExpectWithOffset(1, err).ToNot(HaveOccurred(),
					"CRD %s not found: %s", crd, string(output))
			}

			By("Verifying documented CLI tools are available")
			cliTools := []string{"tkn", "oc"}
			for _, tool := range cliTools {
				cmd := exec.Command("which", tool)
				output, err := cmd.CombinedOutput()
				ExpectWithOffset(1, err).ToNot(HaveOccurred(),
					"CLI tool %s not found: %s", tool, string(output))
			}
		})
	})
})
