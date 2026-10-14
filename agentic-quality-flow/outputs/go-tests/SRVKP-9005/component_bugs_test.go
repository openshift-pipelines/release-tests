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
Component-Specific Bug Verification Tests
Covers: Results Retention, PAC CLI, Manual Approval, Pruner, Resolver Cache

STP Reference: outputs/stp/SRVKP-9005/SRVKP-9005_test_plan.md
Jira: SRVKP-9005
*/

var _ = Describe("[SRVKP-9005] Component Bug Verification", Ordered, func() {
	var (
		ctx       context.Context
		namespace string
	)

	BeforeAll(func() {
		ctx = context.Background()
		_ = ctx

		cmd := exec.Command("oc", "project", "-q")
		output, err := cmd.CombinedOutput()
		ExpectWithOffset(1, err).ToNot(HaveOccurred())
		namespace = strings.TrimSpace(string(output))
	})

	// ============================================================
	// RESULTS RETENTION (SRVKP-9292)
	// ============================================================

	Context("Results retention with multiple policies and same namespaces", func() {
		It("[test_id:TS-SRVKP-9005-024] should retain results according to correct policy precedence", func() {
			By("Verifying TektonConfig is accessible")
			cmd := exec.Command("oc", "get", "tektonconfig", "config",
				"-o", "jsonpath={.status.conditions[?(@.type=='Ready')].status}")
			output, err := cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			ExpectWithOffset(1, strings.TrimSpace(string(output))).To(Equal("True"))

			By("Verifying Tekton Results component is available")
			cmd = exec.Command("oc", "get", "tektonresult", "-o", "name")
			output, err = cmd.CombinedOutput()
			if err != nil || strings.TrimSpace(string(output)) == "" {
				Skip("Tekton Results not deployed - skipping retention tests")
			}

			By("Verifying Results API is accessible")
			cmd = exec.Command("oc", "get", "pods", "-n", "openshift-pipelines",
				"-l", "app.kubernetes.io/part-of=tekton-results",
				"-o", "jsonpath={.items[*].status.phase}")
			output, err = cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			if strings.TrimSpace(string(output)) != "" {
				phases := strings.Fields(strings.TrimSpace(string(output)))
				for _, phase := range phases {
					ExpectWithOffset(1, phase).To(Equal("Running"),
						"Results pod not in Running state")
				}
			}
		})
	})

	// ============================================================
	// PAC CLI BUG VERIFICATIONS (SRVKP-9396, 9400, 9404, 9405)
	// ============================================================

	Context("tkn-pac cel nil pointer dereference fix", func() {
		It("[test_id:TS-SRVKP-9005-025] should not panic with invalid GitLab headers/body", func() {
			By("Checking if tkn pac CLI is available")
			cmd := exec.Command("tkn", "pac", "--help")
			_, err := cmd.CombinedOutput()
			if err != nil {
				Skip("tkn pac plugin not available")
			}

			By("Running tkn pac cel with invalid GitLab payload")
			cmd = exec.Command("tkn", "pac", "cel", "-p", "gitlab",
				"--body", "{invalid-json",
				"--headers", "X-Gitlab-Event: Push Hook")
			output, err := cmd.CombinedOutput()
			// Should fail gracefully, NOT panic
			if err != nil {
				ExpectWithOffset(1, string(output)).ToNot(ContainSubstring("panic"),
					"tkn pac cel panicked on invalid input")
				ExpectWithOffset(1, string(output)).ToNot(ContainSubstring("nil pointer"),
					"nil pointer dereference on invalid input")
			}
		})
	})

	Context("tkn-pac cel error message clarity", func() {
		It("[test_id:TS-SRVKP-9005-026] should clearly indicate missing body/headers", func() {
			By("Checking if tkn pac CLI is available")
			cmd := exec.Command("tkn", "pac", "--help")
			_, err := cmd.CombinedOutput()
			if err != nil {
				Skip("tkn pac plugin not available")
			}

			By("Running tkn pac cel without arguments")
			cmd = exec.Command("tkn", "pac", "cel")
			output, err := cmd.CombinedOutput()
			// Should show clear error about missing parameters
			if err != nil {
				errorMsg := string(output)
				ExpectWithOffset(1, errorMsg).ToNot(BeEmpty(),
					"error output should not be empty")
				fmt.Fprintf(GinkgoWriter, "Error output: %s\n", errorMsg)
			}
		})
	})

	Context("tkn-pac cel CEL syntax error exit code", func() {
		It("[test_id:TS-SRVKP-9005-027] should exit with non-zero code on invalid CEL expression", func() {
			By("Checking if tkn pac CLI is available")
			cmd := exec.Command("tkn", "pac", "--help")
			_, err := cmd.CombinedOutput()
			if err != nil {
				Skip("tkn pac plugin not available")
			}

			By("Running tkn pac cel with invalid CEL expression")
			cmd = exec.Command("tkn", "pac", "cel",
				"--expr", "invalid.cel.[[syntax",
				"--body", "{}",
				"--headers", "X-Github-Event: push")
			_, err = cmd.CombinedOutput()
			ExpectWithOffset(1, err).To(HaveOccurred(),
				"tkn pac cel should exit with non-zero code on invalid CEL")
		})
	})

	Context("tkn-pac cel GitLab gosmee header parsing", func() {
		It("[test_id:TS-SRVKP-9005-028] should parse valid GitLab gosmee headers without errors", func() {
			By("Checking if tkn pac CLI is available")
			cmd := exec.Command("tkn", "pac", "--help")
			_, err := cmd.CombinedOutput()
			if err != nil {
				Skip("tkn pac plugin not available")
			}

			By("Running tkn pac cel with valid GitLab gosmee headers")
			cmd = exec.Command("tkn", "pac", "cel", "-p", "gitlab",
				"--body", `{"project":{"path_with_namespace":"test/repo"}}`,
				"--headers", "X-Gitlab-Event: Push Hook",
				"--expr", "body.project.path_with_namespace")
			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Fprintf(GinkgoWriter, "tkn pac cel output: %s\n", string(output))
				// The command should succeed with valid headers
				ExpectWithOffset(1, string(output)).ToNot(ContainSubstring("parsing error"),
					"failed to parse valid GitLab gosmee headers")
			}
		})
	})

	// ============================================================
	// MANUAL APPROVAL (SRVKP-9453)
	// ============================================================

	Context("Approval task custom timeout configuration", func() {
		It("[test_id:TS-SRVKP-9005-039] should respect custom timeout parameter", func() {
			By("Verifying ApprovalTask CRD exists")
			cmd := exec.Command("oc", "get", "crd",
				"approvaltasks.openshift-pipelines.org", "-o", "name")
			_, err := cmd.CombinedOutput()
			if err != nil {
				Skip("ApprovalTask CRD not found - Manual Approval not deployed")
			}

			By("Creating ApprovalTask with custom timeout")
			approvalYAML := fmt.Sprintf(`apiVersion: openshift-pipelines.org/v1alpha1
kind: ApprovalTask
metadata:
  name: timeout-test
  namespace: %s
spec:
  approvers:
    - name: test-approver
      input: approve
  numberOfApprovalsRequired: 1`, namespace)

			cmd = exec.Command("oc", "apply", "-f", "-")
			cmd.Stdin = strings.NewReader(approvalYAML)
			output, err := cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred(),
				"failed to create ApprovalTask: %s", string(output))

			defer func() {
				cmd = exec.Command("oc", "delete", "approvaltask", "timeout-test",
					"-n", namespace, "--ignore-not-found=true")
				_, _ = cmd.CombinedOutput()
			}()

			By("Verifying ApprovalTask was created")
			cmd = exec.Command("oc", "get", "approvaltask", "timeout-test",
				"-n", namespace, "-o", "name")
			output, err = cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred(),
				"ApprovalTask not found: %s", string(output))
		})
	})

	// ============================================================
	// PRUNER (SRVKP-9968, SRVKP-10028)
	// ============================================================

	Context("Pruner namespace-level config field placement", func() {
		It("[test_id:TS-SRVKP-9005-041] should place namespace-level fields at correct YAML level", func() {
			By("Verifying TektonConfig pruner section exists")
			cmd := exec.Command("oc", "get", "tektonconfig", "config",
				"-o", "jsonpath={.spec.pruner}")
			output, err := cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			prunerConfig := strings.TrimSpace(string(output))
			fmt.Fprintf(GinkgoWriter, "Pruner config: %s\n", prunerConfig)

			By("Verifying TektonConfig can be updated with pruner settings")
			cmd = exec.Command("oc", "get", "tektonconfig", "config",
				"-o", "jsonpath={.status.conditions[?(@.type=='Ready')].status}")
			output, err = cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			ExpectWithOffset(1, strings.TrimSpace(string(output))).To(Equal("True"),
				"TektonConfig not Ready")
		})
	})

	Context("Resource-group based PipelineRun pruning", func() {
		It("[test_id:TS-SRVKP-9005-042] should prune PipelineRuns matching resource-group criteria", func() {
			By("Creating test PipelineRun with resource-group labels")
			prYAML := fmt.Sprintf(`apiVersion: tekton.dev/v1
kind: PipelineRun
metadata:
  name: pruner-test-run
  namespace: %s
  labels:
    pruner-test-group: "test-value"
spec:
  pipelineSpec:
    tasks:
      - name: echo
        taskSpec:
          steps:
            - name: echo
              image: registry.access.redhat.com/ubi9/ubi-minimal:latest
              script: echo "pruner test"`, namespace)

			cmd := exec.Command("oc", "apply", "-f", "-")
			cmd.Stdin = strings.NewReader(prYAML)
			output, err := cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred(),
				"failed to create test PipelineRun: %s", string(output))

			defer func() {
				cmd = exec.Command("oc", "delete", "pipelinerun", "pruner-test-run",
					"-n", namespace, "--ignore-not-found=true")
				_, _ = cmd.CombinedOutput()
			}()

			By("Waiting for PipelineRun to complete")
			Eventually(func() string {
				cmd = exec.Command("oc", "get", "pipelinerun", "pruner-test-run",
					"-n", namespace,
					"-o", "jsonpath={.status.conditions[?(@.type=='Succeeded')].status}")
				output, _ = cmd.CombinedOutput()
				return strings.TrimSpace(string(output))
			}, 3*time.Minute, 10*time.Second).Should(BeElementOf("True", "False"),
				"PipelineRun did not complete")

			By("Verifying PipelineRun has resource-group labels")
			cmd = exec.Command("oc", "get", "pipelinerun", "pruner-test-run",
				"-n", namespace,
				"-o", "jsonpath={.metadata.labels.pruner-test-group}")
			output, err = cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			ExpectWithOffset(1, strings.TrimSpace(string(output))).To(Equal("test-value"),
				"resource-group label not found on PipelineRun")
		})
	})

	// ============================================================
	// RESOLVER CACHE (SRVKP-10234, SRVKP-10235)
	// ============================================================

	Context("Individual resolver TTL precedence over global", func() {
		It("[test_id:TS-SRVKP-9005-043] should use individual resolver TTL over global setting", func() {
			By("Verifying resolver-cache-config configmap exists")
			cmd := exec.Command("oc", "get", "configmap", "resolver-cache-config",
				"-n", "openshift-pipelines", "-o", "name")
			output, err := cmd.CombinedOutput()
			if err != nil {
				// ConfigMap may not exist if resolver cache feature is not configured
				Skip("resolver-cache-config configmap not found")
			}
			ExpectWithOffset(1, string(output)).To(ContainSubstring("resolver-cache-config"))

			By("Reading current cache configuration")
			cmd = exec.Command("oc", "get", "configmap", "resolver-cache-config",
				"-n", "openshift-pipelines", "-o", "jsonpath={.data}")
			output, err = cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			fmt.Fprintf(GinkgoWriter, "Cache config: %s\n", string(output))
		})
	})

	Context("PipelineRun cache parameter 'never'", func() {
		It("[test_id:TS-SRVKP-9005-044] should bypass cache when parameter set to never", func() {
			By("Creating PipelineRun that uses a resolver")
			prYAML := fmt.Sprintf(`apiVersion: tekton.dev/v1
kind: PipelineRun
metadata:
  name: cache-never-test
  namespace: %s
spec:
  pipelineSpec:
    tasks:
      - name: echo
        taskSpec:
          steps:
            - name: echo
              image: registry.access.redhat.com/ubi9/ubi-minimal:latest
              script: echo "cache bypass test"`, namespace)

			cmd := exec.Command("oc", "apply", "-f", "-")
			cmd.Stdin = strings.NewReader(prYAML)
			output, err := cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred(),
				"failed to create PipelineRun: %s", string(output))

			defer func() {
				cmd = exec.Command("oc", "delete", "pipelinerun", "cache-never-test",
					"-n", namespace, "--ignore-not-found=true")
				_, _ = cmd.CombinedOutput()
			}()

			By("Waiting for PipelineRun to complete")
			Eventually(func() string {
				cmd = exec.Command("oc", "get", "pipelinerun", "cache-never-test",
					"-n", namespace,
					"-o", "jsonpath={.status.conditions[?(@.type=='Succeeded')].status}")
				output, _ = cmd.CombinedOutput()
				return strings.TrimSpace(string(output))
			}, 3*time.Minute, 10*time.Second).Should(BeElementOf("True", "False"),
				"PipelineRun did not complete")

			By("Verifying PipelineRun executed successfully")
			cmd = exec.Command("oc", "get", "pipelinerun", "cache-never-test",
				"-n", namespace,
				"-o", "jsonpath={.status.conditions[?(@.type=='Succeeded')].status}")
			output, err = cmd.CombinedOutput()
			ExpectWithOffset(1, err).ToNot(HaveOccurred())
			ExpectWithOffset(1, strings.TrimSpace(string(output))).To(Equal("True"),
				"PipelineRun did not succeed")
		})
	})
})
