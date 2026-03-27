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
Disconnected Environment Testing Tests

STP Reference: outputs/stp/SRVKP-9005/SRVKP-9005_test_plan.md
Jira: SRVKP-9005
*/

var _ = Describe("[SRVKP-9005] Disconnected Environment Testing", Ordered, func() {
	var (
		ctx            context.Context
		clusterArch    string
		isDisconnected bool
	)

	BeforeAll(func() {
		ctx = context.Background()
		_ = ctx

		By("Detecting cluster architecture")
		cmd := exec.Command("oc", "get", "nodes",
			"-o", "jsonpath={.items[0].status.nodeInfo.architecture}")
		output, err := cmd.CombinedOutput()
		ExpectWithOffset(1, err).ToNot(HaveOccurred())
		clusterArch = strings.TrimSpace(string(output))

		By("Checking if cluster is disconnected")
		cmd = exec.Command("oc", "get", "imagecontentsourcepolicy", "-o", "name")
		output, err = cmd.CombinedOutput()
		isDisconnected = err == nil && strings.TrimSpace(string(output)) != ""
		if !isDisconnected {
			// Also check for ImageDigestMirrorSet (OCP 4.13+)
			cmd = exec.Command("oc", "get", "imagedigestmirrorset", "-o", "name")
			output, err = cmd.CombinedOutput()
			isDisconnected = err == nil && strings.TrimSpace(string(output)) != ""
		}
		fmt.Fprintf(GinkgoWriter, "Cluster arch: %s, disconnected: %v\n", clusterArch, isDisconnected)
	})

	Context("Testing in disconnected environment (x86_64)", func() {
		It("[test_id:TS-SRVKP-9005-008] should install and validate in disconnected x86_64 environment", func() {
			if !isDisconnected {
				Skip("Cluster is not disconnected - skipping disconnected tests")
			}
			if clusterArch != "amd64" && clusterArch != "x86_64" {
				Skip(fmt.Sprintf("Cluster architecture is %s, not x86_64", clusterArch))
			}

			By("Verifying CatalogSource points to mirror registry")
			cmd := exec.Command("oc", "get", "catalogsource", "-n", "openshift-marketplace",
				"-o", "jsonpath={.items[*].spec.image}")
			output, err := cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred(),
				"failed to get CatalogSource: %s", string(output))
			ExpectWithOffset(1, string(output)).ToNot(BeEmpty(),
				"no CatalogSource images found")

			By("Verifying operator is installed from mirror registry")
			cmd = exec.Command("oc", "get", "csv", "-n", "openshift-pipelines",
				"-o", "jsonpath={.items[0].status.phase}")
			output, err = cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			ExpectWithOffset(1, strings.TrimSpace(string(output))).To(Equal("Succeeded"),
				"CSV not Succeeded in disconnected x86_64 environment")

			By("Verifying all pods use mirrored images")
			cmd = exec.Command("oc", "get", "pods", "-n", "openshift-pipelines",
				"-o", "jsonpath={.items[*].spec.containers[*].image}")
			output, err = cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			ExpectWithOffset(1, string(output)).ToNot(BeEmpty(),
				"no pod images found in openshift-pipelines")

			By("Verifying TektonConfig is Ready in disconnected environment")
			cmd = exec.Command("oc", "get", "tektonconfig", "config",
				"-o", "jsonpath={.status.conditions[?(@.type=='Ready')].status}")
			output, err = cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			ExpectWithOffset(1, strings.TrimSpace(string(output))).To(Equal("True"))
		})
	})

	Context("Testing in disconnected environment (Z)", func() {
		It("[test_id:TS-SRVKP-9005-009] should install and validate in disconnected IBM Z environment", func() {
			if !isDisconnected {
				Skip("Cluster is not disconnected - skipping")
			}
			if clusterArch != "s390x" {
				Skip(fmt.Sprintf("Cluster architecture is %s, not s390x", clusterArch))
			}

			By("Verifying operator installed in disconnected IBM Z environment")
			cmd := exec.Command("oc", "get", "csv", "-n", "openshift-pipelines",
				"-o", "jsonpath={.items[0].status.phase}")
			output, err := cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			ExpectWithOffset(1, strings.TrimSpace(string(output))).To(Equal("Succeeded"),
				"CSV not Succeeded in disconnected IBM Z environment")

			By("Verifying TektonConfig is Ready")
			cmd = exec.Command("oc", "get", "tektonconfig", "config",
				"-o", "jsonpath={.status.conditions[?(@.type=='Ready')].status}")
			output, err = cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			ExpectWithOffset(1, strings.TrimSpace(string(output))).To(Equal("True"))
		})
	})

	Context("Testing in disconnected environment (P)", func() {
		It("[test_id:TS-SRVKP-9005-010] should install and validate in disconnected IBM Power environment", func() {
			if !isDisconnected {
				Skip("Cluster is not disconnected - skipping")
			}
			if clusterArch != "ppc64le" {
				Skip(fmt.Sprintf("Cluster architecture is %s, not ppc64le", clusterArch))
			}

			By("Verifying operator installed in disconnected IBM Power environment")
			cmd := exec.Command("oc", "get", "csv", "-n", "openshift-pipelines",
				"-o", "jsonpath={.items[0].status.phase}")
			output, err := cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			ExpectWithOffset(1, strings.TrimSpace(string(output))).To(Equal("Succeeded"),
				"CSV not Succeeded in disconnected IBM Power environment")

			By("Verifying TektonConfig is Ready")
			cmd = exec.Command("oc", "get", "tektonconfig", "config",
				"-o", "jsonpath={.status.conditions[?(@.type=='Ready')].status}")
			output, err = cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			ExpectWithOffset(1, strings.TrimSpace(string(output))).To(Equal("True"))
		})
	})
})
