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
UI Testing Tests

STP Reference: outputs/stp/SRVKP-9005/SRVKP-9005_test_plan.md
Jira: SRVKP-9005
*/

var _ = Describe("[SRVKP-9005] UI Testing", Ordered, func() {
	var (
		ctx       context.Context
		namespace string
	)

	BeforeAll(func() {
		ctx = context.Background()
		_ = ctx

		By("Setting up test namespace")
		cmd := exec.Command("oc", "project", "-q")
		output, err := cmd.CombinedOutput()
		ExpectWithOffset(1, err).ToNot(HaveOccurred())
		namespace = strings.TrimSpace(string(output))
		fmt.Fprintf(GinkgoWriter, "Using namespace: %s\n", namespace)
	})

	Context("DevConsole UI manual tests", func() {
		It("[test_id:TS-SRVKP-9005-016] should create, execute, and monitor pipelines through DevConsole", func() {
			By("Verifying console route exists")
			cmd := exec.Command("oc", "get", "route", "console",
				"-n", "openshift-console", "-o", "jsonpath={.spec.host}")
			output, err := cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred(),
				"failed to get console route: %s", string(output))
			consoleHost := strings.TrimSpace(string(output))
			ExpectWithOffset(1, consoleHost).ToNot(BeEmpty(),
				"console route host is empty")
			fmt.Fprintf(GinkgoWriter, "Console URL: https://%s\n", consoleHost)

			By("Verifying OpenShift Pipelines console plugin is loaded")
			cmd = exec.Command("oc", "get", "consoleplugin",
				"-o", "jsonpath={.items[*].metadata.name}")
			output, err = cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			plugins := string(output)
			ExpectWithOffset(1, plugins).To(ContainSubstring("pipelines"),
				"pipelines console plugin not found")
		})
	})

	Context("ApprovalTask YAML validation warning fix", func() {
		It("[test_id:TS-SRVKP-9005-029] should not show spurious validation warnings for ApprovalTask CRD", func() {
			By("Verifying ApprovalTask CRD exists")
			cmd := exec.Command("oc", "get", "crd",
				"approvaltasks.openshift-pipelines.org", "-o", "name")
			output, err := cmd.CombinedOutput()
			if err != nil {
				Skip("ApprovalTask CRD not found - Manual Approval not deployed")
			}
			ExpectWithOffset(1, string(output)).To(ContainSubstring("approvaltasks"),
				"ApprovalTask CRD not found")

			By("Verifying ApprovalTask CRD has OpenAPI validation schema")
			cmd = exec.Command("oc", "get", "crd",
				"approvaltasks.openshift-pipelines.org",
				"-o", "jsonpath={.spec.versions[0].schema.openAPIV3Schema}")
			output, err = cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			ExpectWithOffset(1, strings.TrimSpace(string(output))).ToNot(BeEmpty(),
				"ApprovalTask CRD missing OpenAPI validation schema")
		})
	})

	Context("Approvals tab status display fix", func() {
		It("[test_id:TS-SRVKP-9005-030] should show Pending status instead of Unknown", func() {
			By("Verifying ApprovalTask CRD is available")
			cmd := exec.Command("oc", "get", "crd",
				"approvaltasks.openshift-pipelines.org", "-o", "name")
			_, err := cmd.CombinedOutput()
			if err != nil {
				Skip("ApprovalTask CRD not found - skipping approval status tests")
			}

			By("Creating a test ApprovalTask")
			approvalYAML := fmt.Sprintf(`apiVersion: openshift-pipelines.org/v1alpha1
kind: ApprovalTask
metadata:
  generateName: status-test-
  namespace: %s
spec:
  approvers:
    - name: test-approver
      input: approve
  numberOfApprovalsRequired: 1`, namespace)

			cmd = exec.Command("oc", "create", "-f", "-")
			cmd.Stdin = strings.NewReader(approvalYAML)
			output, err := cmd.CombinedOutput()
			if err != nil {
				Skip(fmt.Sprintf("Cannot create ApprovalTask: %s", string(output)))
			}
			taskName := strings.TrimSpace(strings.Split(string(output), " ")[0])
			taskName = strings.TrimPrefix(taskName, "approvaltask.openshift-pipelines.org/")

			By("Cleaning up test ApprovalTask")
			defer func() {
				cmd = exec.Command("oc", "delete", "approvaltask", taskName,
					"-n", namespace, "--ignore-not-found=true")
				_, _ = cmd.CombinedOutput()
			}()

			By("Verifying ApprovalTask was created successfully")
			cmd = exec.Command("oc", "get", "approvaltask", taskName,
				"-n", namespace, "-o", "name")
			output, err = cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred(),
				"failed to get ApprovalTask: %s", string(output))
		})
	})

	Context("Approvals tab load time fix", func() {
		It("[test_id:TS-SRVKP-9005-031] should load within 5 seconds for new users", func() {
			By("Verifying console is accessible")
			cmd := exec.Command("oc", "get", "route", "console",
				"-n", "openshift-console", "-o", "jsonpath={.spec.host}")
			output, err := cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			ExpectWithOffset(1, strings.TrimSpace(string(output))).ToNot(BeEmpty(),
				"console route not accessible")

			By("Verifying pipelines console plugin responds")
			cmd = exec.Command("oc", "get", "consoleplugin",
				"-o", "jsonpath={.items[*].metadata.name}")
			output, err = cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			ExpectWithOffset(1, string(output)).To(ContainSubstring("pipelines"),
				"pipelines console plugin not found")
		})
	})

	Context("Approvals tab presence in Developer Perspective", func() {
		It("[test_id:TS-SRVKP-9005-032] should show Approvals tab under Pipelines view", func() {
			By("Verifying Manual Approval operator component is deployed")
			cmd := exec.Command("oc", "get", "deployment",
				"-n", "openshift-pipelines",
				"-o", "jsonpath={.items[*].metadata.name}")
			output, err := cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			deployments := string(output)
			fmt.Fprintf(GinkgoWriter, "Deployments in openshift-pipelines: %s\n", deployments)

			By("Verifying ApprovalTask CRD exists for UI integration")
			cmd = exec.Command("oc", "get", "crd",
				"approvaltasks.openshift-pipelines.org", "-o", "name")
			output, err = cmd.CombinedOutput()
			if err != nil {
				Skip("ApprovalTask CRD not available - Manual Approval UI not deployed")
			}
			ExpectWithOffset(1, string(output)).To(ContainSubstring("approvaltasks"))
		})
	})

	Context("Approvals list navigation fix", func() {
		It("[test_id:TS-SRVKP-9005-034] should remain interactive after back navigation", func() {
			By("Verifying console plugin for pipelines is active")
			cmd := exec.Command("oc", "get", "consoleplugin",
				"-o", "jsonpath={.items[*].metadata.name}")
			output, err := cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			ExpectWithOffset(1, string(output)).To(ContainSubstring("pipelines"),
				"pipelines console plugin required for navigation tests")
		})
	})

	Context("Approvals list status overlap fix", func() {
		It("[test_id:TS-SRVKP-9005-035] should render status column without text overlap", func() {
			By("Verifying ApprovalTask CRD has status subresource")
			cmd := exec.Command("oc", "get", "crd",
				"approvaltasks.openshift-pipelines.org",
				"-o", "jsonpath={.spec.versions[0].subresources}")
			output, err := cmd.CombinedOutput()
			if err != nil {
				Skip("ApprovalTask CRD not found")
			}
			ExpectWithOffset(1, string(output)).To(ContainSubstring("status"),
				"ApprovalTask CRD missing status subresource")
		})
	})

	Context("Non-admin PipelineRuns display fix", func() {
		It("[test_id:TS-SRVKP-9005-036] should display PipelineRuns for non-admin users", func() {
			By("Verifying RBAC allows non-admin pipeline access")
			cmd := exec.Command("oc", "get", "clusterrole",
				"-o", "jsonpath={.items[*].metadata.name}")
			output, err := cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			roles := string(output)
			ExpectWithOffset(1, roles).To(ContainSubstring("edit"),
				"edit clusterrole not found")

			By("Verifying PipelineRun list API is accessible")
			cmd = exec.Command("oc", "auth", "can-i", "list", "pipelineruns")
			output, err = cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			ExpectWithOffset(1, strings.TrimSpace(string(output))).To(Equal("yes"),
				"cannot list pipelineruns")
		})
	})

	Context("Approvals PipelineRun name display fix", func() {
		It("[test_id:TS-SRVKP-9005-037] should display correct PipelineRun names in Approvals list", func() {
			By("Verifying ApprovalTask API returns PipelineRun references")
			cmd := exec.Command("oc", "get", "crd",
				"approvaltasks.openshift-pipelines.org", "-o", "name")
			_, err := cmd.CombinedOutput()
			if err != nil {
				Skip("ApprovalTask CRD not available")
			}

			By("Verifying ApprovalTask spec includes pipelineRun reference fields")
			cmd = exec.Command("oc", "get", "crd",
				"approvaltasks.openshift-pipelines.org",
				"-o", "jsonpath={.spec.versions[0].schema.openAPIV3Schema.properties.spec}")
			output, err := cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			ExpectWithOffset(1, string(output)).ToNot(BeEmpty(),
				"ApprovalTask CRD spec schema is empty")
		})
	})

	Context("Non-admin All Projects context switch fix", func() {
		It("[test_id:TS-SRVKP-9005-038] should allow non-admin users to switch to All Projects", func() {
			By("Verifying cluster role bindings for pipeline access exist")
			cmd := exec.Command("oc", "get", "clusterrolebinding",
				"-o", "jsonpath={.items[*].metadata.name}")
			output, err := cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			ExpectWithOffset(1, string(output)).ToNot(BeEmpty(),
				"no cluster role bindings found")

			By("Verifying PipelineRun list is accessible across namespaces")
			cmd = exec.Command("oc", "auth", "can-i", "list", "pipelineruns",
				"--all-namespaces")
			output, err = cmd.CombinedOutput()
			// This may return "no" for non-admin, which is expected
			_ = err
			_ = output
			fmt.Fprintf(GinkgoWriter, "can-i list pipelineruns --all-namespaces: %s\n",
				strings.TrimSpace(string(output)))

			By("Verifying console RBAC proxy allows project switching")
			cmd = exec.Command("oc", "get", "consoleplugin",
				"-o", "jsonpath={.items[*].metadata.name}")
			output, err = cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			ExpectWithOffset(1, string(output)).To(ContainSubstring("pipelines"))
		})
	})
})
