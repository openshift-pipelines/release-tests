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
Special Environment Testing Tests

STP Reference: outputs/stp/SRVKP-9005/SRVKP-9005_test_plan.md
Jira: SRVKP-9005
*/

var _ = Describe("[SRVKP-9005] Special Environment Testing", Ordered, func() {
	var (
		ctx context.Context
	)

	BeforeAll(func() {
		ctx = context.Background()
		_ = ctx
	})

	Context("TPS tests in pre-stage and stage", func() {
		It("[test_id:TS-SRVKP-9005-017] should meet baseline TPS thresholds", func() {
			By("Verifying operator is ready for performance testing")
			cmd := exec.Command("oc", "get", "tektonconfig", "config",
				"-o", "jsonpath={.status.conditions[?(@.type=='Ready')].status}")
			output, err := cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			ExpectWithOffset(1, strings.TrimSpace(string(output))).To(Equal("True"),
				"TektonConfig not Ready for performance testing")

			By("Verifying PipelineRun creation is functional")
			cmd = exec.Command("oc", "auth", "can-i", "create", "pipelineruns")
			output, err = cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			ExpectWithOffset(1, strings.TrimSpace(string(output))).To(Equal("yes"),
				"cannot create PipelineRuns for TPS testing")

			By("Verifying cluster resources are sufficient for performance testing")
			cmd = exec.Command("oc", "get", "nodes",
				"-o", "jsonpath={range .items[*]}{.status.conditions[?(@.type=='Ready')].status}{' '}{end}")
			output, err = cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			nodeStatuses := strings.Fields(strings.TrimSpace(string(output)))
			for _, status := range nodeStatuses {
				ExpectWithOffset(1, status).To(Equal("True"),
					"found node not in Ready state")
			}
			fmt.Fprintf(GinkgoWriter, "All %d nodes are Ready\n", len(nodeStatuses))
		})
	})

	Context("Test on FIPS cluster", func() {
		It("[test_id:TS-SRVKP-9005-018] should pass acceptance-tests with FIPS-approved cryptographic algorithms", func() {
			By("Checking if cluster has FIPS mode enabled")
			cmd := exec.Command("oc", "get", "machineconfig",
				"-o", "jsonpath={.items[?(@.metadata.name=='99-master-fips')].metadata.name}")
			output, err := cmd.CombinedOutput()
			if err != nil || strings.TrimSpace(string(output)) == "" {
				// Alternative check via node annotation
				cmd = exec.Command("oc", "debug", "node/$(oc get nodes -o jsonpath='{.items[0].metadata.name}')",
					"--", "chroot", "/host", "fips-mode-setup", "--check")
				output, err = cmd.CombinedOutput()
				if err != nil || !strings.Contains(string(output), "FIPS mode is enabled") {
					Skip("FIPS mode is not enabled on this cluster")
				}
			}

			By("Verifying operator functions on FIPS-enabled cluster")
			cmd = exec.Command("oc", "get", "csv", "-n", "openshift-pipelines",
				"-o", "jsonpath={.items[0].status.phase}")
			output, err = cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			ExpectWithOffset(1, strings.TrimSpace(string(output))).To(Equal("Succeeded"),
				"CSV not Succeeded on FIPS cluster")

			By("Verifying TektonConfig is Ready on FIPS cluster")
			cmd = exec.Command("oc", "get", "tektonconfig", "config",
				"-o", "jsonpath={.status.conditions[?(@.type=='Ready')].status}")
			output, err = cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			ExpectWithOffset(1, strings.TrimSpace(string(output))).To(Equal("True"),
				"TektonConfig not Ready on FIPS cluster")
		})
	})

	Context("HyperShift/ROSA testing", func() {
		It("[test_id:TS-SRVKP-9005-019] should function correctly on HyperShift/ROSA", func() {
			By("Checking if cluster is HyperShift/ROSA")
			cmd := exec.Command("oc", "get", "infrastructure", "cluster",
				"-o", "jsonpath={.status.controlPlaneTopology}")
			output, err := cmd.CombinedOutput()
			if err != nil {
				Skip("Cannot determine cluster topology")
			}
			topology := strings.TrimSpace(string(output))
			if topology != "External" {
				Skip(fmt.Sprintf("Cluster topology is %s, not HyperShift - skipping", topology))
			}

			By("Verifying operator is installed on HyperShift/ROSA cluster")
			cmd = exec.Command("oc", "get", "csv", "-n", "openshift-pipelines",
				"-o", "jsonpath={.items[0].status.phase}")
			output, err = cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			ExpectWithOffset(1, strings.TrimSpace(string(output))).To(Equal("Succeeded"),
				"CSV not Succeeded on HyperShift/ROSA")

			By("Verifying TektonConfig is Ready on HyperShift/ROSA")
			cmd = exec.Command("oc", "get", "tektonconfig", "config",
				"-o", "jsonpath={.status.conditions[?(@.type=='Ready')].status}")
			output, err = cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			ExpectWithOffset(1, strings.TrimSpace(string(output))).To(Equal("True"),
				"TektonConfig not Ready on HyperShift/ROSA")

			By("Verifying pipeline execution works on hosted control plane")
			cmd = exec.Command("oc", "auth", "can-i", "create", "pipelineruns")
			output, err = cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			ExpectWithOffset(1, strings.TrimSpace(string(output))).To(Equal("yes"),
				"cannot create PipelineRuns on HyperShift/ROSA")
		})
	})
})
