# OpenShift Pipelines Test Plan

## **CVE-2026-33211: Information Disclosure via Path Traversal in Git Resolver - Quality Engineering Plan**

### **Metadata & Tracking**

- **Enhancement(s):** N/A (Security vulnerability fix)
- **Feature Tracking:** [SRVKP-11236](https://redhat.atlassian.net/browse/SRVKP-11236)
- **Epic Tracking:** [SRVKP-11236](https://redhat.atlassian.net/browse/SRVKP-11236)
- **QE Owner(s):** Anitha Natarajan
- **Owning SIG:** pipelines
- **Participating SIGs:** pipelines, operator

**Document Conventions (if applicable):** This STP covers a security vulnerability (CVE-2026-33211) requiring verification that the path traversal fix in the git resolver is effective and does not introduce regressions.

### **Feature Overview**

CVE-2026-33211 identifies an information disclosure vulnerability in the Tekton Pipelines git resolver. The `pathInRepo` parameter in `ResolutionRequest` resources does not properly validate input, allowing path traversal sequences (e.g., `../`) to read arbitrary files from the resolver pod's filesystem. This includes sensitive files such as ServiceAccount tokens. The file contents are returned base64-encoded in `resolutionrequest.status.data`.

**Affected component:** `openshift-pipelines/pipelines-operator-proxy-rhel9`

**Attack vector:** A tenant with permission to create `ResolutionRequests` (e.g., by creating `TaskRuns` or `PipelineRuns` that use the git resolver) can exploit the `pathInRepo` parameter to traverse directories and read arbitrary files.

**Fix versions (upstream):** Tekton Pipelines 1.0.1, 1.3.3, 1.6.1, 1.9.2, 1.10.2

**Downstream target:** OpenShift Pipelines 1.21 (pipelines-1.20 stream)

---

### **I. Motivation and Requirements Review (QE Review Guidelines)**

This section documents the mandatory QE review process. The goal is to understand the feature's value, technology, and testability before formal test planning.

#### **1. Requirement & User Story Review Checklist**

- [ ] **Review Requirements**
  - Reviewed the relevant requirements.
  - CVE-2026-33211 describes a path traversal vulnerability in the git resolver's `pathInRepo` parameter. The fix must sanitize path inputs to prevent directory traversal beyond the repository root.
- [ ] **Understand Value and Customer Use Cases**
  - Confirmed clear user stories and understood.
  - Understand the difference between U/S and D/S requirements.
  - **What is the value of the feature for RH customers**: Prevents information disclosure where a malicious tenant could read arbitrary files (including ServiceAccount tokens) from the git resolver pod, protecting multi-tenant cluster security.
  - Ensured requirements contain relevant **customer use cases**.
  - Customer use case: Multi-tenant OpenShift clusters where different teams create PipelineRuns/TaskRuns using the git resolver must be protected from cross-tenant information leakage.
- [ ] **Testability**
  - Confirmed requirements are **testable and unambiguous**.
  - The vulnerability is testable by attempting path traversal sequences in the `pathInRepo` parameter and verifying they are rejected or sanitized.
- [ ] **Acceptance Criteria**
  - Ensured acceptance criteria are **defined clearly** (clear user stories; D/S requirements clearly defined in Jira).
  - AC1: `ResolutionRequests` with path traversal sequences (`../`, `..%2f`, etc.) in `pathInRepo` must be rejected with an appropriate error.
  - AC2: Legitimate `pathInRepo` values (e.g., `tasks/my-task.yaml`) must continue to resolve correctly.
  - AC3: The resolver pod's filesystem (including `/var/run/secrets/`) must not be accessible via `pathInRepo`.
- [ ] **Non-Functional Requirements (NFRs)**
  - Confirmed coverage for NFRs, including Performance, Security, Usability, Downtime, Connectivity, Monitoring (alerts/metrics), Scalability, Portability (e.g., cloud support), and Docs.
  - Primary NFR is **Security**: Input validation must block all known path traversal variants. No performance regression expected from path validation.

#### **2. Known Limitations**

- The upstream fix covers Tekton Pipelines versions 1.0.1, 1.3.3, 1.6.1, 1.9.2, and 1.10.2. The downstream OpenShift Pipelines 1.21 build must incorporate the appropriate upstream patch.
- Testing is limited to the git resolver component; other resolvers (bundle, hub, cluster) are out of scope unless regression testing reveals issues.
- The CVE tracker indicates this affects the `pipelines-operator-proxy-rhel9` container image specifically.

#### **3. Technology and Design Review**

- [ ] **Developer Handoff/QE Kickoff**
  - A meeting where Dev/Arch walked QE through the design, architecture, and implementation details. **Critical for identifying untestable aspects early.**
  - Pending: QE kickoff meeting to review the path validation implementation approach in the git resolver.
- [ ] **Technology Challenges**
  - Identified potential testing challenges related to the underlying technology.
  - Testing requires creating `ResolutionRequests` directly or via `TaskRuns`/`PipelineRuns` that reference git-hosted resources. Must verify both direct API access and indirect access via pipeline execution.
- [ ] **Test Environment Needs**
  - Determined necessary **test environment setups and tools**.
  - Requires an OpenShift cluster with OpenShift Pipelines operator installed, a git repository for resolver targets, and RBAC configuration for multi-tenant scenarios.
- [ ] **API Extensions**
  - Reviewed new or modified APIs and their impact on testing.
  - The `ResolutionRequest` API's `pathInRepo` field behavior changes: path traversal sequences are now rejected. No new API fields are introduced.
- [ ] **Topology Considerations**
  - Evaluated multi-cluster, network topology, and architectural impacts.
  - Single-cluster topology is sufficient. The vulnerability is in the resolver pod's file handling, not in network-level access.

### **II. Software Test Plan (STP)**

This STP serves as the **overall roadmap for testing**, detailing the scope, approach, resources, and schedule.

#### **1. Scope of Testing**

This test plan covers verification of the CVE-2026-33211 fix in the Tekton Pipelines git resolver, ensuring that path traversal via the `pathInRepo` parameter is properly blocked while preserving legitimate git resolution functionality.

**Testing Goals**

- Verify that path traversal sequences in `pathInRepo` are rejected by the git resolver
- Verify that legitimate file paths in git repositories continue to resolve correctly
- Verify that ServiceAccount tokens and other sensitive pod filesystem contents are not accessible
- Verify that the fix is effective across multiple encoding variants of traversal sequences
- Verify no regression in PipelineRun/TaskRun execution using the git resolver

**Out of Scope (Testing Scope Exclusions)**

- [ ] Bundle resolver, hub resolver, and cluster resolver (not affected by this CVE)
- [ ] Tekton Chains signing and attestation functionality
- [ ] Tekton Triggers and EventListener functionality
- [ ] Performance benchmarking of the git resolver (no expected impact from input validation)
- [ ] Upstream Tekton release testing (covered by upstream CI)

#### **2. Test Strategy**

**Functional**

- [ ] **Functional Testing** -- Validates that the feature works according to specified requirements and user stories
  - *Details:* Verify path traversal blocking in the git resolver by testing various malicious `pathInRepo` inputs. Verify continued functionality with legitimate paths.
- [ ] **Automation Testing** -- Confirms test automation plan is in place for CI and regression coverage (all tests are expected to be automated)
  - *Details:* All test scenarios will be automated in Go/Ginkgo (Tier 1) and integrated into the OpenShift Pipelines CI pipeline for continuous regression detection.
- [ ] **Regression Testing** -- Verifies that new changes do not break existing functionality
  - *Details:* Existing git resolver tests must continue to pass. PipelineRuns and TaskRuns using git-referenced tasks/pipelines must execute successfully with legitimate paths.

**Non-Functional**

- [ ] **Performance Testing** -- Validates feature performance meets requirements (latency, throughput, resource usage)
  - *Details:* Not applicable. Path validation adds negligible overhead to resolution requests.
- [ ] **Scale Testing** -- Validates feature behavior under increased load and at production-like scale (e.g., large number of PipelineRuns or concurrent ResolutionRequests)
  - *Details:* Not applicable for this security fix.
- [ ] **Security Testing** -- Verifies security requirements, RBAC, authentication, authorization, and vulnerability scanning
  - *Details:* Primary focus of this STP. Verify that the path traversal vulnerability is fully remediated. Test RBAC boundaries to confirm that tenants cannot access files outside their git repository scope. Verify that encoded path traversal variants (URL-encoded, double-encoded) are also blocked.
- [ ] **Usability Testing** -- Validates user experience and accessibility requirements
  - *Details:* Verify that error messages for rejected path traversal attempts are clear and actionable (not exposing internal filesystem details).
- [ ] **Monitoring** -- Does the feature require metrics and/or alerts?
  - *Details:* Consider whether rejected path traversal attempts should generate audit events or metrics for security monitoring.

**Integration & Compatibility**

- [ ] **Compatibility Testing** -- Ensures feature works across supported platforms, versions, and configurations
  - *Details:* Verify fix works on OCP 4.14+ with OpenShift Pipelines 1.21.
- [ ] **Upgrade Testing** -- Validates upgrade paths from previous versions, data migration, and configuration preservation
  - *Details:* Verify that upgrading from a vulnerable OpenShift Pipelines version to the fixed version properly applies the path validation. Existing `ResolutionRequests` with legitimate paths must not be disrupted.
- [ ] **Dependencies** -- Blocked by deliverables from other components/products. Identify what we need from other teams before we can test.
  - *Details:* Depends on the upstream Tekton Pipelines patch being incorporated into the downstream OpenShift Pipelines operator build.
- [ ] **Cross Integrations** -- Does the feature affect other features or require testing by other teams? Identify the impact we cause.
  - *Details:* The git resolver is used by PipelineRuns and TaskRuns that reference git-hosted task/pipeline definitions. Any team using git-referenced resources should verify their workflows are not affected.

**Infrastructure**

- [ ] **Cloud Testing** -- Does the feature require multi-cloud platform testing? Consider cloud-specific features.
  - *Details:* Not specifically required. The vulnerability is in the resolver pod's file handling, which is platform-agnostic.

#### **3. Test Environment**

- **Cluster Topology:** Single or Multi-node (Single-node sufficient for functional testing)
- **OCP & OpenShift Pipelines Version(s):** OCP 4.14+ with OpenShift Pipelines 1.21+
- **CPU Virtualization:** Standard (no special CPU requirements)
- **Compute Resources:** Standard worker nodes (1+ worker nodes)
- **Special Hardware:** None required
- **Storage:** Default storage class (standard)
- **Network:** OVN-Kubernetes CNI (default)
- **Required Operators:** Red Hat OpenShift Pipelines (namespace: openshift-pipelines)
- **Platform:** OpenShift Container Platform (OCP)
- **Special Configurations:** Git repository accessible from the cluster for resolver testing; RBAC configuration for multi-tenant test scenarios

#### **3.1. Testing Tools & Frameworks**

- **Test Framework:** Ginkgo v2 + Gomega (Go)
- **CI/CD:** OpenShift Pipelines CI (Prow)
- **Other Tools:** `tkn` CLI, `oc` CLI, `kubectl`

#### **4. Entry Criteria**

The following conditions must be met before testing can begin:

- [ ] Requirements and design documents are **approved and merged**
- [ ] Test environment can be **set up and configured** (see Section II.3 - Test Environment)
- [ ] The upstream Tekton Pipelines patch for CVE-2026-33211 is incorporated into the downstream OpenShift Pipelines operator build
- [ ] The patched `pipelines-operator-proxy-rhel9` container image is available for testing
- [ ] A git repository with test task/pipeline definitions is accessible from the test cluster

#### **5. Risks**

- [ ] **Timeline/Schedule**
  - Risk: CVE remediation timelines may be compressed, limiting testing window
  - Mitigation: Prioritize critical path traversal tests; automate all scenarios for rapid execution
- [ ] **Test Coverage**
  - Risk: Novel path traversal encoding variants may not be covered by initial test scenarios
  - Mitigation: Include comprehensive encoding variants (URL-encoded, double-encoded, null-byte injection, Unicode normalization) based on OWASP path traversal cheat sheet
- [ ] **Test Environment**
  - Risk: Git repository access from cluster may be restricted by network policies
  - Mitigation: Prepare a local git server or use cluster-internal git repository for testing
- [ ] **Untestable Aspects**
  - Risk: Verifying that all filesystem paths are inaccessible may not be exhaustive
  - Mitigation: Focus on known sensitive paths (`/var/run/secrets/`, `/etc/`, `/proc/`) and verify the validation logic blocks traversal at the input level
- [ ] **Resource Constraints**
  - Risk: Security testing expertise may be needed for comprehensive path traversal variant coverage
  - Mitigation: Leverage OWASP testing guides and security team consultation
- [ ] **Dependencies**
  - Risk: Upstream patch availability may delay downstream testing
  - Mitigation: Begin test development against the known vulnerability description; execute once the patched build is available

---

### **III. Test Scenarios & Traceability**

This section links requirements to test coverage, enabling reviewers to verify all requirements are tested.

#### **1. Requirements-to-Tests Mapping**

- **[SRVKP-11236]** -- Verify path traversal via `pathInRepo` is blocked in the git resolver
  - Verify that a `ResolutionRequest` with `pathInRepo` containing `../` sequences returns an error and does not resolve filesystem content -- **Tier 1 (P1)**
  - Verify that a `TaskRun` using the git resolver with a `pathInRepo` containing `../../etc/passwd` does not return file contents -- **Tier 1 (P1)**
  - Verify that a `PipelineRun` using the git resolver with path traversal in `pathInRepo` is rejected before execution -- **Tier 1 (P1)**
  - Verify that URL-encoded path traversal (`..%2f`, `%2e%2e%2f`) in `pathInRepo` is detected and rejected -- **Tier 1 (P1)**
  - Verify that double-encoded path traversal (`..%252f`) in `pathInRepo` is detected and rejected -- **Tier 1 (P1)**
  - Verify that null-byte injection (`../../../etc/passwd%00.yaml`) in `pathInRepo` is detected and rejected -- **Tier 1 (P1)**
  - Verify that ServiceAccount token path (`/var/run/secrets/kubernetes.io/serviceaccount/token`) cannot be accessed via `pathInRepo` -- **Tier 1 (P1)**

- **[SRVKP-11236]** -- Verify legitimate git resolver functionality is preserved (regression)
  - Verify that a `ResolutionRequest` with a valid `pathInRepo` (e.g., `tasks/my-task.yaml`) resolves the correct file content from the git repository -- **Tier 1 (P1)**
  - Verify that a `TaskRun` referencing a git-hosted task definition via the git resolver executes successfully -- **Tier 1 (P1)**
  - Verify that a `PipelineRun` referencing git-hosted pipeline and task definitions via the git resolver executes successfully -- **Tier 1 (P2)**
  - Verify that nested valid paths (e.g., `dir1/dir2/task.yaml`) resolve correctly without being falsely flagged as traversal -- **Tier 1 (P2)**

- **[SRVKP-11236]** -- Verify error handling and observability for rejected requests
  - Verify that rejected path traversal attempts return a clear error message indicating invalid path without exposing internal filesystem details -- **Tier 1 (P2)**
  - Verify that the `ResolutionRequest` status reflects the rejection with an appropriate condition and reason -- **Tier 1 (P2)**
  - Verify that rejected path traversal attempts are logged in the resolver pod logs for security auditing -- **Tier 1 (P3)**

- **[SRVKP-11236]** -- Verify upgrade path from vulnerable version
  - Verify that upgrading OpenShift Pipelines from a pre-fix version to the patched version applies the path validation correctly -- **Tier 1 (P2)**
  - Verify that existing `PipelineRuns` using legitimate git resolver paths continue to function after the upgrade -- **Tier 1 (P2)**

- **[SRVKP-11236]** -- Verify RBAC boundaries with git resolver
  - Verify that a tenant with `ResolutionRequest` create permissions cannot use path traversal to access files outside the specified git repository -- **Tier 1 (P1)**
  - Verify that a tenant without `ResolutionRequest` create permissions cannot trigger git resolution at all -- **Tier 1 (P2)**

---

### **IV. Sign-off and Approval**

This Software Test Plan requires approval from the following stakeholders:

* **Reviewers:**
  - [Name / @github-username]
  - [Name / @github-username]
* **Approvers:**
  - [Name / @github-username]
  - [Name / @github-username]
