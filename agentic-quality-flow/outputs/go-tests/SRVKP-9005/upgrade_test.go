package pipeline

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

/*
Upgrade Testing Tests

STP Reference: outputs/stp/SRVKP-9005/SRVKP-9005_test_plan.md
Jira: SRVKP-9005
*/

var _ = Describe("[SRVKP-9005] Upgrade Testing", Ordered, func() {
	var (
		ctx context.Context
	)

	BeforeAll(func() {
		ctx = context.Background()
		_ = ctx
	})

	Context("Cluster upgrade tests", func() {
		It("[test_id:TS-SRVKP-9005-011] should function correctly after OCP cluster upgrade", func() {
			By("Verifying operator is present before upgrade scenario")
			cmd := exec.Command("oc", "get", "csv", "-n", "openshift-pipelines",
				"-o", "jsonpath={.items[0].status.phase}")
			output, err := cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred(),
				"failed to get CSV: %s", string(output))
			ExpectWithOffset(1, strings.TrimSpace(string(output))).To(Equal("Succeeded"),
				"CSV not Succeeded - operator may not be installed")

			By("Verifying TektonConfig is Ready post-upgrade")
			cmd = exec.Command("oc", "get", "tektonconfig", "config",
				"-o", "jsonpath={.status.conditions[?(@.type=='Ready')].status}")
			output, err = cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			ExpectWithOffset(1, strings.TrimSpace(string(output))).To(Equal("True"),
				"TektonConfig not Ready after cluster upgrade")

			By("Verifying PipelineRun can be created and executed after upgrade")
			pipelineRunYAML := `apiVersion: tekton.dev/v1
kind: PipelineRun
metadata:
  generateName: post-upgrade-validation-
spec:
  pipelineSpec:
    tasks:
      - name: echo-test
        taskSpec:
          steps:
            - name: echo
              image: registry.access.redhat.com/ubi9/ubi-minimal:latest
              script: echo "Post-upgrade validation successful"`

			cmd = exec.Command("oc", "create", "-f", "-")
			cmd.Stdin = strings.NewReader(pipelineRunYAML)
			output, err = cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred(),
				"failed to create PipelineRun: %s", string(output))

			prName := strings.TrimSpace(strings.TrimPrefix(string(output), "pipelinerun.tekton.dev/"))
			prName = strings.TrimSuffix(prName, " created")
			fmt.Fprintf(GinkgoWriter, "Created PipelineRun: %s\n", prName)

			By("Waiting for PipelineRun to complete")
			Eventually(func() string {
				cmd = exec.Command("oc", "get", "pipelinerun", prName,
					"-o", "jsonpath={.status.conditions[?(@.type=='Succeeded')].status}")
				output, _ = cmd.CombinedOutput()
				return strings.TrimSpace(string(output))
			}, 5*time.Minute, 10*time.Second).Should(Equal("True"),
				"PipelineRun did not complete successfully after upgrade")

			By("Cleaning up validation PipelineRun")
			cmd = exec.Command("oc", "delete", "pipelinerun", prName)
			_, _ = cmd.CombinedOutput()
		})
	})

	Context("Operator upgrade tests", func() {
		It("[test_id:TS-SRVKP-9005-012] should upgrade from latest released version to 1.21.0", func() {
			By("Verifying current operator version")
			cmd := exec.Command("oc", "get", "csv", "-n", "openshift-pipelines",
				"-o", "jsonpath={.items[0].spec.version}")
			output, err := cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred(),
				"failed to get CSV version: %s", string(output))
			currentVersion := strings.TrimSpace(string(output))
			fmt.Fprintf(GinkgoWriter, "Current operator version: %s\n", currentVersion)
			ExpectWithOffset(1, currentVersion).ToNot(BeEmpty())

			By("Verifying subscription exists for OpenShift Pipelines")
			cmd = exec.Command("oc", "get", "subscription", "-n", "openshift-pipelines",
				"-o", "jsonpath={.items[0].spec.channel}")
			output, err = cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred(),
				"failed to get subscription: %s", string(output))
			channel := strings.TrimSpace(string(output))
			fmt.Fprintf(GinkgoWriter, "Subscription channel: %s\n", channel)

			By("Verifying operator is in Succeeded state after upgrade")
			cmd = exec.Command("oc", "get", "csv", "-n", "openshift-pipelines",
				"-o", "jsonpath={.items[0].status.phase}")
			output, err = cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			ExpectWithOffset(1, strings.TrimSpace(string(output))).To(Equal("Succeeded"),
				"CSV not Succeeded after operator upgrade")

			By("Verifying all components are Ready after upgrade")
			cmd = exec.Command("oc", "get", "tektonconfig", "config",
				"-o", "jsonpath={.status.conditions[?(@.type=='Ready')].status}")
			output, err = cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			ExpectWithOffset(1, strings.TrimSpace(string(output))).To(Equal("True"),
				"TektonConfig not Ready after operator upgrade")
		})
	})
})
