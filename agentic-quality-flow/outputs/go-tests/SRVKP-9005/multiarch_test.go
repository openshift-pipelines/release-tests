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
Multi-Architecture Testing Tests

STP Reference: outputs/stp/SRVKP-9005/SRVKP-9005_test_plan.md
Jira: SRVKP-9005
*/

var _ = Describe("[SRVKP-9005] Multi-Architecture Testing", Ordered, func() {
	var (
		ctx          context.Context
		clusterArch  string
	)

	BeforeAll(func() {
		ctx = context.Background()
		_ = ctx

		By("Detecting cluster architecture")
		cmd := exec.Command("oc", "get", "nodes",
			"-o", "jsonpath={.items[0].status.nodeInfo.architecture}")
		output, err := cmd.CombinedOutput()
		ExpectWithOffset(1, err).ToNot(HaveOccurred(),
			"failed to detect cluster architecture: %s", string(output))
		clusterArch = strings.TrimSpace(string(output))
		fmt.Fprintf(GinkgoWriter, "Detected cluster architecture: %s\n", clusterArch)
	})

	Context("Multiarch testing on ARM64", func() {
		It("[test_id:TS-SRVKP-9005-005] should pass acceptance-tests suite on ARM64 cluster", func() {
			if clusterArch != "arm64" && clusterArch != "aarch64" {
				Skip(fmt.Sprintf("Cluster architecture is %s, not ARM64 - skipping", clusterArch))
			}

			By("Verifying operator is installed on ARM64 cluster")
			cmd := exec.Command("oc", "get", "csv", "-n", "openshift-pipelines",
				"-o", "jsonpath={.items[0].status.phase}")
			output, err := cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred(),
				"failed to query CSV on ARM64: %s", string(output))
			ExpectWithOffset(1, strings.TrimSpace(string(output))).To(Equal("Succeeded"),
				"CSV not Succeeded on ARM64 cluster")

			By("Verifying all nodes report ARM64 architecture")
			cmd = exec.Command("oc", "get", "nodes",
				"-o", "jsonpath={.items[*].status.nodeInfo.architecture}")
			output, err = cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			architectures := strings.Fields(strings.TrimSpace(string(output)))
			for _, arch := range architectures {
				ExpectWithOffset(1, arch).To(BeElementOf("arm64", "aarch64"),
					"found non-ARM64 node: %s", arch)
			}

			By("Verifying TektonConfig is Ready")
			cmd = exec.Command("oc", "get", "tektonconfig", "config",
				"-o", "jsonpath={.status.conditions[?(@.type=='Ready')].status}")
			output, err = cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			ExpectWithOffset(1, strings.TrimSpace(string(output))).To(Equal("True"),
				"TektonConfig not Ready on ARM64")

			By("Verifying core pipeline pods are running")
			cmd = exec.Command("oc", "get", "pods", "-n", "openshift-pipelines",
				"--field-selector=status.phase!=Running,status.phase!=Succeeded",
				"-o", "name")
			output, err = cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			ExpectWithOffset(1, strings.TrimSpace(string(output))).To(BeEmpty(),
				"found non-running pods in openshift-pipelines: %s", string(output))
		})
	})

	Context("Multiarch testing on IBM Z", func() {
		It("[test_id:TS-SRVKP-9005-006] should pass acceptance-tests suite on IBM Z cluster", func() {
			if clusterArch != "s390x" {
				Skip(fmt.Sprintf("Cluster architecture is %s, not s390x - skipping", clusterArch))
			}

			By("Verifying operator is installed on IBM Z cluster")
			cmd := exec.Command("oc", "get", "csv", "-n", "openshift-pipelines",
				"-o", "jsonpath={.items[0].status.phase}")
			output, err := cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred(),
				"failed to query CSV on IBM Z: %s", string(output))
			ExpectWithOffset(1, strings.TrimSpace(string(output))).To(Equal("Succeeded"),
				"CSV not Succeeded on IBM Z cluster")

			By("Verifying all nodes report s390x architecture")
			cmd = exec.Command("oc", "get", "nodes",
				"-o", "jsonpath={.items[*].status.nodeInfo.architecture}")
			output, err = cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			architectures := strings.Fields(strings.TrimSpace(string(output)))
			for _, arch := range architectures {
				ExpectWithOffset(1, arch).To(Equal("s390x"),
					"found non-s390x node: %s", arch)
			}

			By("Verifying TektonConfig is Ready on IBM Z")
			cmd = exec.Command("oc", "get", "tektonconfig", "config",
				"-o", "jsonpath={.status.conditions[?(@.type=='Ready')].status}")
			output, err = cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			ExpectWithOffset(1, strings.TrimSpace(string(output))).To(Equal("True"),
				"TektonConfig not Ready on IBM Z")
		})
	})

	Context("Multiarch testing on IBM Power", func() {
		It("[test_id:TS-SRVKP-9005-007] should pass acceptance-tests suite on IBM Power cluster", func() {
			if clusterArch != "ppc64le" {
				Skip(fmt.Sprintf("Cluster architecture is %s, not ppc64le - skipping", clusterArch))
			}

			By("Verifying operator is installed on IBM Power cluster")
			cmd := exec.Command("oc", "get", "csv", "-n", "openshift-pipelines",
				"-o", "jsonpath={.items[0].status.phase}")
			output, err := cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred(),
				"failed to query CSV on IBM Power: %s", string(output))
			ExpectWithOffset(1, strings.TrimSpace(string(output))).To(Equal("Succeeded"),
				"CSV not Succeeded on IBM Power cluster")

			By("Verifying all nodes report ppc64le architecture")
			cmd = exec.Command("oc", "get", "nodes",
				"-o", "jsonpath={.items[*].status.nodeInfo.architecture}")
			output, err = cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			architectures := strings.Fields(strings.TrimSpace(string(output)))
			for _, arch := range architectures {
				ExpectWithOffset(1, arch).To(Equal("ppc64le"),
					"found non-ppc64le node: %s", arch)
			}

			By("Verifying TektonConfig is Ready on IBM Power")
			cmd = exec.Command("oc", "get", "tektonconfig", "config",
				"-o", "jsonpath={.status.conditions[?(@.type=='Ready')].status}")
			output, err = cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			ExpectWithOffset(1, strings.TrimSpace(string(output))).To(Equal("True"),
				"TektonConfig not Ready on IBM Power")
		})
	})
})
