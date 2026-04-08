# OpenShift Pipelines Test Plan

## **Test OpenShift Pipelines 1.21.0 - Quality Engineering Plan**

### **Metadata & Tracking**

- **Enhancement(s):** N/A (Minor release testing epic)
- **Feature Tracking:** N/A
- **Epic Tracking:** [SRVKP-9005](https://redhat.atlassian.net/browse/SRVKP-9005)
- **QE Owner(s):** OpenShift Pipelines QE Team
- **Owning Component:** QA
- **Participating Components:** Tekton Pipelines, Tekton Triggers, Tekton Chains, Pipelines as Code, Tekton CLI, Tekton Results, Pruner, Manual Approval, UI
- **Product Version:** OpenShift Pipelines 1.21.0
- **Platform:** OCP 4.14+

**Document Conventions:** This STP covers the complete release testing scope for OpenShift Pipelines 1.21.0, including feature testing, regression testing, multi-architecture validation, disconnected environment testing, upgrade testing, and bug verification. Test scenarios are derived from the 46 subtasks under the release testing epic.

### **Feature Overview**

OpenShift Pipelines 1.21.0 is a minor release of the Tekton-based CI/CD pipeline engine for OpenShift. This release testing epic (SRVKP-9005) encompasses comprehensive validation of all product components across multiple environments, architectures, and upgrade paths. The testing scope includes new feature validation, regression testing, multi-arch support (ARM64, IBM Z, IBM Power), disconnected environment testing, operator and cluster upgrade paths, CLI verification, UI validation, and bug verification for all fixed issues in this release.

---

### **I. Motivation and Requirements Review (QE Review Guidelines)**

This section documents the mandatory QE review process. The goal is to understand the release scope, technology changes, and testability before formal test planning.

#### **1. Requirement & User Story Review Checklist**

- [ ] **Review Requirements**
  - Reviewed the release testing epic SRVKP-9005 and all 46 subtasks
  - Release scope covers feature testing, regression, multi-arch, disconnected, upgrades, CLI, UI, and bug verification
- [ ] **Understand Value and Customer Use Cases**
  - OpenShift Pipelines 1.21.0 delivers CI/CD pipeline capabilities on OpenShift
  - Key customer value: reliable Tekton-based pipelines with supply chain security (Chains), event-driven triggers, and Pipelines as Code integration
  - Release testing ensures production readiness across all supported configurations
- [ ] **Testability**
  - All subtasks define concrete testing activities with clear pass/fail criteria
  - Automated test suites (acceptance-tests, release-tests) cover core functionality
  - Manual test procedures documented for Hub, DevConsole UI, and CLI verification
- [ ] **Acceptance Criteria**
  - Epic acceptance criteria: "All subtasks are done"
  - Each subtask has component-specific acceptance criteria
- [ ] **Non-Functional Requirements (NFRs)**
  - Performance: Covered through TPS (Transactions Per Second) tests in pre-stage and stage environments
  - Security: CLI binary SHA/signature verification on Mac and Windows; FIPS cluster testing
  - Usability: DevConsole UI manual tests; CLI usability validation
  - Monitoring: Tekton Results retention policy validation
  - Scalability: Multi-arch testing across ARM64, IBM Z, IBM Power

#### **2. Known Limitations**

- Tier 2 (Python/pytest) tests are not generated for this project; E2E tests use Go/Gauge (managed outside QualityFlow)
- Some subtasks reference internal GitLab CI pipelines (plumbing repo) that are outside the scope of this STP
- UI bugs (SRVKP-9435, SRVKP-9437, SRVKP-9438, SRVKP-9440, SRVKP-9443, SRVKP-9448, SRVKP-9451, SRVKP-9452) are in "To Do" status and may require resolution before release
- Resolver cache bugs (SRVKP-10234) is in "Code Review" status

#### **3. Technology and Design Review**

- [ ] **Developer Handoff/QE Kickoff**
  - Release testing follows established procedures documented in the release-tests repository
  - CI configuration updates required in ci-config.yaml for version 1.21
- [ ] **Technology Challenges**
  - Multi-architecture support requires ARM64, IBM Z (s390x), and IBM Power (ppc64le) clusters
  - Disconnected environment testing requires mirror registries and air-gapped network configurations
  - FIPS compliance testing requires FIPS-enabled OCP clusters
- [ ] **Test Environment Needs**
  - Connected OCP clusters (x86_64, ARM64, IBM Z, IBM Power)
  - Disconnected OCP clusters (x86_64, Z, P environments)
  - FIPS-enabled OCP cluster
  - HyperShift/ROSA environment for interop testing
  - Pre-stage and stage environments for TPS testing
- [ ] **API Extensions**
  - ApprovalTask CRD (openshift-pipelines.org/v1alpha1) - Manual approval workflow
  - Resolver cache configuration (resolver-cache-config) - Caching for pipeline resolvers
  - TektonPruner configuration - Resource retention and cleanup policies
- [ ] **Topology Considerations**
  - Single-node and multi-node cluster topologies
  - HyperShift hosted control plane topology (ROSA)
  - Disconnected/air-gapped network topology

### **II. Software Test Plan (STP)**

This STP serves as the overall roadmap for testing OpenShift Pipelines 1.21.0, detailing the scope, approach, resources, and schedule.

#### **1. Scope of Testing**

The testing scope for OpenShift Pipelines 1.21.0 covers the following areas derived from the epic subtasks:

- Feature testing for new capabilities in this release ([SRVKP-9006](https://redhat.atlassian.net/browse/SRVKP-9006))
- Automated regression testing via acceptance-tests and release-tests suites
- Multi-architecture validation: ARM64 ([SRVKP-9011](https://redhat.atlassian.net/browse/SRVKP-9011)), IBM Z ([SRVKP-9015](https://redhat.atlassian.net/browse/SRVKP-9015)), IBM Power ([SRVKP-9026](https://redhat.atlassian.net/browse/SRVKP-9026))
- Disconnected environment testing: x86_64 ([SRVKP-9016](https://redhat.atlassian.net/browse/SRVKP-9016)), Z ([SRVKP-9007](https://redhat.atlassian.net/browse/SRVKP-9007)), P ([SRVKP-9012](https://redhat.atlassian.net/browse/SRVKP-9012))
- Operator upgrade testing ([SRVKP-9020](https://redhat.atlassian.net/browse/SRVKP-9020)) and cluster upgrade testing ([SRVKP-9018](https://redhat.atlassian.net/browse/SRVKP-9018))
- CLI (tkn/opc) testing ([SRVKP-9025](https://redhat.atlassian.net/browse/SRVKP-9025)) and entitlement tests ([SRVKP-9017](https://redhat.atlassian.net/browse/SRVKP-9017))
- Tekton Hub manual tests ([SRVKP-9010](https://redhat.atlassian.net/browse/SRVKP-9010))
- DevConsole UI manual tests ([SRVKP-9019](https://redhat.atlassian.net/browse/SRVKP-9019))
- TPS (performance) testing in pre-stage and stage ([SRVKP-9021](https://redhat.atlassian.net/browse/SRVKP-9021))
- FIPS cluster testing ([SRVKP-9027](https://redhat.atlassian.net/browse/SRVKP-9027))
- HyperShift/ROSA interop testing ([SRVKP-9013](https://redhat.atlassian.net/browse/SRVKP-9013))
- Bug verification for all bugs with fixVersion 1.21 ([SRVKP-9024](https://redhat.atlassian.net/browse/SRVKP-9024))
- Binary verification on mirror.openshift.com ([SRVKP-9009](https://redhat.atlassian.net/browse/SRVKP-9009))
- CLI binary SHA/signature verification ([SRVKP-9023](https://redhat.atlassian.net/browse/SRVKP-9023))
- Documentation review ([SRVKP-9028](https://redhat.atlassian.net/browse/SRVKP-9028))
- Post-release operator installation verification ([SRVKP-9022](https://redhat.atlassian.net/browse/SRVKP-9022))

**Testing Goals**

- Validate all new features in OpenShift Pipelines 1.21.0 function correctly
- Confirm no regressions in existing Pipeline, Trigger, Chain, PAC, CLI, and Results functionality
- Verify multi-architecture support across ARM64, IBM Z, and IBM Power
- Validate disconnected/air-gapped installation and operation
- Confirm operator and cluster upgrade paths from previous versions
- Verify all fixed bugs in this release
- Validate UI (DevConsole) integration and Manual Approval workflow
- Confirm FIPS compliance on FIPS-enabled clusters
- Validate HyperShift/ROSA interoperability

**Out of Scope (Testing Scope Exclusions)**

- [ ] Upstream Tekton unit tests (covered by upstream CI)
- [ ] Performance benchmarking beyond TPS tests (no load/stress testing in scope)
- [ ] Third-party Git provider integrations beyond what is covered by PAC tests
- [ ] OCP platform-level testing (handled by OCP QE)
- [ ] Tier 2 Python/pytest generation (project uses Go/Gauge for E2E)

#### **2. Test Strategy**

**Functional**

- [ ] **Functional Testing** -- Validates that all OpenShift Pipelines 1.21.0 features work according to specifications
  - *Details:* Feature testing ([SRVKP-9006](https://redhat.atlassian.net/browse/SRVKP-9006)) covers new capabilities. Automated acceptance-tests and release-tests suites run on connected clusters. Manual testing covers Hub ([SRVKP-9010](https://redhat.atlassian.net/browse/SRVKP-9010)) and DevConsole UI ([SRVKP-9019](https://redhat.atlassian.net/browse/SRVKP-9019)).
- [ ] **Automation Testing** -- Confirms test automation plan is in place for CI and regression coverage
  - *Details:* Release branch release-v1.21 created in openshift-pipelines/release-tests ([SRVKP-9014](https://redhat.atlassian.net/browse/SRVKP-9014)). CI configuration updated in ci-config.yaml ([SRVKP-9008](https://redhat.atlassian.net/browse/SRVKP-9008)). Automated suites run via Prow CI.
- [ ] **Regression Testing** -- Verifies that new changes do not break existing functionality
  - *Details:* Full regression suite execution across all components. Bug verification for all ON_QA bugs ([SRVKP-9024](https://redhat.atlassian.net/browse/SRVKP-9024)). Release-testing bugs tracked with label `release-testing-bug`.

**Non-Functional**

- [ ] **Performance Testing** -- Validates feature performance meets requirements
  - *Details:* TPS tests executed in pre-stage and stage environments ([SRVKP-9021](https://redhat.atlassian.net/browse/SRVKP-9021)). Measures transactions per second for pipeline execution.
- [ ] **Scale Testing** -- Validates feature behavior under increased load
  - *Details:* Not explicitly scoped for this release. Covered indirectly by TPS tests and multi-node cluster testing.
- [ ] **Security Testing** -- Verifies security requirements, RBAC, and supply chain security
  - *Details:* FIPS cluster testing ([SRVKP-9027](https://redhat.atlassian.net/browse/SRVKP-9027)). CLI binary SHA/signature verification ([SRVKP-9023](https://redhat.atlassian.net/browse/SRVKP-9023)). Tekton Chains supply chain security validation. RBAC testing for Manual Approval workflows.
- [ ] **Usability Testing** -- Validates user experience and accessibility
  - *Details:* DevConsole UI manual tests ([SRVKP-9019](https://redhat.atlassian.net/browse/SRVKP-9019)). CLI (tkn/opc) usability validation ([SRVKP-9025](https://redhat.atlassian.net/browse/SRVKP-9025)). Manual Approval UI workflow testing.
- [ ] **Monitoring** -- Validates metrics and operational visibility
  - *Details:* Tekton Results retention policy validation ([SRVKP-9292](https://redhat.atlassian.net/browse/SRVKP-9292)). Pruner configuration validation ([SRVKP-9968](https://redhat.atlassian.net/browse/SRVKP-9968), [SRVKP-10028](https://redhat.atlassian.net/browse/SRVKP-10028)).

**Integration & Compatibility**

- [ ] **Compatibility Testing** -- Ensures compatibility across supported platforms and architectures
  - *Details:* Multi-arch testing on ARM64 ([SRVKP-9011](https://redhat.atlassian.net/browse/SRVKP-9011)), IBM Z ([SRVKP-9015](https://redhat.atlassian.net/browse/SRVKP-9015)), IBM Power ([SRVKP-9026](https://redhat.atlassian.net/browse/SRVKP-9026)). OCP 4.14+ compatibility. Disconnected environment support.
- [ ] **Upgrade Testing** -- Validates upgrade paths from previous versions
  - *Details:* Operator upgrade from latest released version to 1.21.0 on three OCP versions ([SRVKP-9020](https://redhat.atlassian.net/browse/SRVKP-9020)). Cluster upgrade testing ([SRVKP-9018](https://redhat.atlassian.net/browse/SRVKP-9018)). Automatic execution via CI with manual verification.
- [ ] **Dependencies** -- External dependencies required for testing
  - *Details:* OCP platform availability (4.14+). Mirror registry for disconnected testing. ROSA/HyperShift environment from interop team ([SRVKP-9013](https://redhat.atlassian.net/browse/SRVKP-9013)). CLI advisory state for binary signing.
- [ ] **Cross Integrations** -- Impact on other features or teams
  - *Details:* DevConsole UI integration (OpenShift Console team). ROSA interop testing (ROSA team notification). Tekton Hub integration. PAC integration with Git providers (GitHub, GitLab).

**Infrastructure**

- [ ] **Cloud Testing** -- Multi-cloud platform testing requirements
  - *Details:* HyperShift/ROSA testing on AWS ([SRVKP-9013](https://redhat.atlassian.net/browse/SRVKP-9013)). Standard OCP testing on supported cloud platforms.

#### **3. Test Environment**

- **Cluster Topology:** Single or Multi-node (Standard OCP clusters, HyperShift hosted control plane)
- **OCP & OpenShift Pipelines Version(s):** OCP 4.14+ with OpenShift Pipelines 1.21.0
- **CPU Architecture:** x86_64 (primary), ARM64, IBM Z (s390x), IBM Power (ppc64le)
- **Compute Resources:** Standard worker nodes (minimum 1 worker node)
- **Special Hardware:** IBM Z and IBM Power systems for multi-arch testing
- **Storage:** Standard storage class (no block or shared storage requirements)
- **Network:** OVN-Kubernetes CNI; Disconnected/air-gapped configurations for isolated testing
- **Required Operators:** Red Hat OpenShift Pipelines (namespace: openshift-pipelines)
- **Platform:** OpenShift Container Platform (OCP) 4.14+
- **Special Configurations:** FIPS-enabled cluster; Disconnected/air-gapped environments; HyperShift/ROSA

#### **3.1. Testing Tools & Frameworks**

- **Test Framework:** Ginkgo v2 + Gomega (Go), Gauge (E2E)
- **CI/CD:** Prow, GitLab CI (plumbing repo)
- **CLI Tools:** tkn, opc, oc, kubectl
- **Other Tools:** Cosign (signature verification), mirror registry tooling

#### **4. Entry Criteria**

The following conditions must be met before testing can begin:

- [ ] Requirements and design documents are **approved and merged**
- [ ] Test environment can be **set up and configured** (see Section II.3 - Test Environment)
- [ ] Release branch release-v1.21 created in openshift-pipelines/release-tests ([SRVKP-9014](https://redhat.atlassian.net/browse/SRVKP-9014))
- [ ] CI configuration updated in ci-config.yaml for 1.21 ([SRVKP-9008](https://redhat.atlassian.net/browse/SRVKP-9008))
- [ ] Component versions updated in env/default/default.properties
- [ ] OpenShift Pipelines 1.21.0 operator available for installation
- [ ] Multi-arch clusters (ARM64, IBM Z, IBM Power) provisioned and accessible
- [ ] Disconnected environments (x86_64, Z, P) configured with mirror registries
- [ ] FIPS-enabled OCP cluster available

#### **5. Risks**

- [ ] **Timeline/Schedule**
  - Risk: Multi-arch and disconnected environment provisioning may delay testing start
  - Mitigation: Initiate environment provisioning early; parallelize testing across architectures
- [ ] **Test Coverage**
  - Risk: Open UI bugs (9 issues in "To Do" status) may block UI testing scenarios
  - Mitigation: Prioritize critical UI bugs (SRVKP-9450, SRVKP-9452); defer minor UI issues if not fixed before release
- [ ] **Test Environment**
  - Risk: Disconnected environments require specialized mirror registry setup and may be unstable
  - Mitigation: Use established disconnected testing procedures; maintain fallback connected environments
- [ ] **Untestable Aspects**
  - Risk: HyperShift/ROSA testing depends on external team (interop team) availability
  - Mitigation: Notify ROSA interop team early ([SRVKP-9013](https://redhat.atlassian.net/browse/SRVKP-9013)); track as dependency
- [ ] **Resource Constraints**
  - Risk: Testing across 4 architectures and 3 disconnected environments requires significant infrastructure
  - Mitigation: Leverage CI automation for connected tests; prioritize manual testing for disconnected and multi-arch
- [ ] **Dependencies**
  - Risk: CLI binary signing depends on advisory state transition; resolver cache fix (SRVKP-10234) in "Code Review"
  - Mitigation: Track advisory state; verify signing after state transition; monitor SRVKP-10234 for merge

---

### **III. Test Scenarios & Traceability**

This section links requirements to test coverage, enabling reviewers to verify all requirements are tested.

#### **1. Requirements-to-Tests Mapping**

- **[SRVKP-9006](https://redhat.atlassian.net/browse/SRVKP-9006)** -- Feature testing for new capabilities in OpenShift Pipelines 1.21.0
  - Validate all new features introduced in this release function correctly according to specifications
  - Priority: **P1** | Tier: **Tier 1**

- **[SRVKP-9008](https://redhat.atlassian.net/browse/SRVKP-9008)** -- Update ci-config.yaml for release 1.21
  - Verify CI configuration is updated and automated test jobs execute successfully against 1.21 builds
  - Priority: **P1** | Tier: **Tier 1**

- **[SRVKP-9009](https://redhat.atlassian.net/browse/SRVKP-9009)** -- Verify binaries on mirror.openshift.com
  - Run pipeline verify-binaries-and-256sum with the 1.21 release URL and confirm all binaries are present and checksums match
  - Priority: **P2** | Tier: **Tier 1**

- **[SRVKP-9010](https://redhat.atlassian.net/browse/SRVKP-9010)** -- Hub manual tests
  - Verify Tekton Hub installation, task/pipeline catalog browsing, and task installation from Hub
  - Priority: **P2** | Tier: **Tier 1**

- **[SRVKP-9011](https://redhat.atlassian.net/browse/SRVKP-9011)** -- Multiarch testing on ARM64
  - Execute acceptance-tests suite on ARM64 architecture cluster and verify all tests pass
  - Priority: **P1** | Tier: **Tier 1**

- **[SRVKP-9012](https://redhat.atlassian.net/browse/SRVKP-9012)** -- Testing in disconnected environment (P)
  - Install and validate OpenShift Pipelines 1.21.0 in a disconnected IBM Power environment with mirror registry
  - Priority: **P1** | Tier: **Tier 1**

- **[SRVKP-9013](https://redhat.atlassian.net/browse/SRVKP-9013)** -- HyperShift/ROSA testing
  - Notify ROSA interop team and verify OpenShift Pipelines 1.21.0 functions correctly on HyperShift/ROSA
  - Priority: **P2** | Tier: **Tier 1**

- **[SRVKP-9014](https://redhat.atlassian.net/browse/SRVKP-9014)** -- Create release branch in release-tests for 1.21
  - Verify release-v1.21 branch is created, component versions updated in default.properties, and CI jobs target the correct branch
  - Priority: **P1** | Tier: **Tier 1**

- **[SRVKP-9015](https://redhat.atlassian.net/browse/SRVKP-9015)** -- Multiarch testing on IBM Z
  - Execute acceptance-tests suite on IBM Z (s390x) architecture cluster and verify all tests pass
  - Priority: **P1** | Tier: **Tier 1**

- **[SRVKP-9016](https://redhat.atlassian.net/browse/SRVKP-9016)** -- Testing in disconnected environment (x86_64)
  - Install and validate OpenShift Pipelines 1.21.0 in a disconnected x86_64 environment with mirror registry
  - Priority: **P1** | Tier: **Tier 1**

- **[SRVKP-9007](https://redhat.atlassian.net/browse/SRVKP-9007)** -- Testing in disconnected environment (Z)
  - Install and validate OpenShift Pipelines 1.21.0 in a disconnected IBM Z environment with mirror registry
  - Priority: **P1** | Tier: **Tier 1**

- **[SRVKP-9017](https://redhat.atlassian.net/browse/SRVKP-9017)** -- TKN entitlement tests
  - Run entitlement-tests pipeline after CLI build moves to stage and verify entitlement validation passes
  - Priority: **P2** | Tier: **Tier 1**

- **[SRVKP-9018](https://redhat.atlassian.net/browse/SRVKP-9018)** -- Cluster upgrade tests
  - Verify OpenShift Pipelines functions correctly after OCP cluster upgrade with the operator installed
  - Priority: **P1** | Tier: **Tier 1**

- **[SRVKP-9019](https://redhat.atlassian.net/browse/SRVKP-9019)** -- DevConsole UI manual tests
  - Execute manual UI test scenarios for pipeline creation, execution, and monitoring through the OpenShift DevConsole
  - Priority: **P2** | Tier: **Tier 1**

- **[SRVKP-9020](https://redhat.atlassian.net/browse/SRVKP-9020)** -- Operator upgrade tests
  - Verify operator upgrade from the latest released version to 1.21.0 on three OCP versions; confirm automatic CI execution and results
  - Priority: **P1** | Tier: **Tier 1**

- **[SRVKP-9021](https://redhat.atlassian.net/browse/SRVKP-9021)** -- TPS tests in pre-stage and stage
  - Execute TPS (transactions per second) tests in pre-stage and stage environments and verify performance meets baseline thresholds
  - Priority: **P2** | Tier: **Tier 1**

- **[SRVKP-9022](https://redhat.atlassian.net/browse/SRVKP-9022)** -- Post-release verify operator installation
  - After release, verify that the OpenShift Pipelines operator can be installed on all supported OCP versions from OperatorHub
  - Priority: **P1** | Tier: **Tier 1**

- **[SRVKP-9023](https://redhat.atlassian.net/browse/SRVKP-9023)** -- Verification of CLI binary SHA/signature on Mac and Windows
  - Download CLI binaries from the release server; verify SHA256 checksums and Cosign signatures on macOS and Windows platforms
  - Priority: **P2** | Tier: **Tier 1**

- **[SRVKP-9024](https://redhat.atlassian.net/browse/SRVKP-9024)** -- Bug verification
  - Verify all bugs with fixVersion of 1.21 in ON_QA state are correctly fixed and pass acceptance criteria
  - Priority: **P1** | Tier: **Tier 1**

- **[SRVKP-9025](https://redhat.atlassian.net/browse/SRVKP-9025)** -- TKN/OPC CLI tests
  - Execute CLI test suite for tkn and opc commands; verify command output, error handling, and plugin functionality
  - Priority: **P2** | Tier: **Tier 1**

- **[SRVKP-9026](https://redhat.atlassian.net/browse/SRVKP-9026)** -- Multiarch testing on IBM Power
  - Execute acceptance-tests suite on IBM Power (ppc64le) architecture cluster and verify all tests pass
  - Priority: **P1** | Tier: **Tier 1**

- **[SRVKP-9027](https://redhat.atlassian.net/browse/SRVKP-9027)** -- Test on FIPS cluster
  - Run acceptance-tests on a FIPS-enabled OCP cluster and verify all cryptographic operations use FIPS-approved algorithms
  - Priority: **P1** | Tier: **Tier 1**

- **[SRVKP-9028](https://redhat.atlassian.net/browse/SRVKP-9028)** -- Documentation review
  - Review release documentation for accuracy and completeness; verify documented features match implemented behavior
  - Priority: **P3** | Tier: **Tier 1**

- **[SRVKP-9292](https://redhat.atlassian.net/browse/SRVKP-9292)** -- Results retention with multiple policies and same namespaces
  - Configure TektonConfig with two retention policies targeting the same namespace; verify results are retained according to the correct policy precedence
  - Priority: **P1** | Tier: **Tier 1**

- **[SRVKP-9396](https://redhat.atlassian.net/browse/SRVKP-9396)** -- tkn-pac cel nil pointer dereference with invalid GitLab headers/body
  - Run `tkn pac cel -p gitlab` with invalid/malformed GitLab payload and headers; verify the command does not panic and returns a meaningful error
  - Priority: **P1** | Tier: **Tier 1**

- **[SRVKP-9400](https://redhat.atlassian.net/browse/SRVKP-9400)** -- tkn pac cel unclear error messages for missing body/headers
  - Run `tkn pac cel` without arguments or with only headers/body; verify error messages clearly indicate what is missing
  - Priority: **P2** | Tier: **Tier 1**

- **[SRVKP-9404](https://redhat.atlassian.net/browse/SRVKP-9404)** -- tkn pac cel CEL syntax errors exit with code 0
  - Run `tkn pac cel` with an invalid CEL expression; verify the command exits with a non-zero exit code
  - Priority: **P2** | Tier: **Tier 1**

- **[SRVKP-9405](https://redhat.atlassian.net/browse/SRVKP-9405)** -- tkn pac cel fails to parse valid GitLab gosmee headers
  - Run `tkn pac cel` with valid X-Gitlab-Event headers from gosmee save scripts; verify the command parses them correctly without errors
  - Priority: **P2** | Tier: **Tier 1**

- **[SRVKP-9435](https://redhat.atlassian.net/browse/SRVKP-9435)** -- UI YAML validation warning for ApprovalTask CRD
  - Edit an existing ApprovalTask from Console UI with Tekton/OpenShift YAML schema enabled; verify no spurious validation warnings are displayed
  - Priority: **P3** | Tier: **Tier 1**

- **[SRVKP-9437](https://redhat.atlassian.net/browse/SRVKP-9437)** -- UI Approvals tab status shown as Unknown for pending tasks
  - Create a manual approval task for a PipelineRun; verify the Approvals tab correctly shows "Pending" status (not "Unknown")
  - Priority: **P2** | Tier: **Tier 1**

- **[SRVKP-9438](https://redhat.atlassian.net/browse/SRVKP-9438)** -- UI Approvals tab takes longer to load
  - Log in as a new user; navigate to the Approvals tab; verify the tab loads within an acceptable time (under 5 seconds)
  - Priority: **P2** | Tier: **Tier 1**

- **[SRVKP-9440](https://redhat.atlassian.net/browse/SRVKP-9440)** -- UI Approvals tab missing from Pipelines view in Developer Perspective
  - Navigate to the Developer Perspective in the OpenShift Console; verify the Approvals tab is present under the Pipelines view
  - Priority: **P1** | Tier: **Tier 1**

- **[SRVKP-9441](https://redhat.atlassian.net/browse/SRVKP-9441)** -- CLI set character limit on Approval Task comments
  - Run `tkn-approvaltask` with an excessively long comment; verify the CLI enforces a character limit or handles gracefully
  - Priority: **P3** | Tier: **Tier 1**

- **[SRVKP-9443](https://redhat.atlassian.net/browse/SRVKP-9443)** -- UI navigating back to Approvals list prevents interaction
  - Modify filters on the Approvals list; navigate to a detail view and back; verify the list remains interactive and filters are preserved
  - Priority: **P2** | Tier: **Tier 1**

- **[SRVKP-9448](https://redhat.atlassian.net/browse/SRVKP-9448)** -- UI current status overlaps in Approvals list view
  - View the Approvals tab with multiple approval tasks in various states; verify the Current Status column renders without text overlap
  - Priority: **P3** | Tier: **Tier 1**

- **[SRVKP-9450](https://redhat.atlassian.net/browse/SRVKP-9450)** -- UI PipelineRuns not showing up as non-admin users
  - Log in as a non-admin user; navigate to PipelineRuns; verify the list loads and displays runs without hanging on a loading spinner
  - Priority: **P1** | Tier: **Tier 1**

- **[SRVKP-9451](https://redhat.atlassian.net/browse/SRVKP-9451)** -- UI PipelineRun names missing in Approvals list view
  - View the Approvals tab; verify the PipelineRun name column displays the correct pipeline run name for all entries
  - Priority: **P3** | Tier: **Tier 1**

- **[SRVKP-9452](https://redhat.atlassian.net/browse/SRVKP-9452)** -- UI non-admin users unable to switch context to All Projects
  - Log in as a non-admin user; view PipelineRuns/Approvals under All Projects; switch context and verify resources remain visible
  - Priority: **P1** | Tier: **Tier 1**

- **[SRVKP-9453](https://redhat.atlassian.net/browse/SRVKP-9453)** -- Manual approval task timeout not configurable via pipeline
  - Configure an ApprovalTask with a custom timeout parameter; verify the controller respects the configured timeout value
  - Priority: **P1** | Tier: **Tier 1**

- **[SRVKP-9460](https://redhat.atlassian.net/browse/SRVKP-9460)** -- tkn-assist show example in opc CLI usage
  - Run `opc assist --help`; verify the help output includes an Examples section with usage examples
  - Priority: **P3** | Tier: **Tier 1**

- **[SRVKP-9968](https://redhat.atlassian.net/browse/SRVKP-9968)** -- Pruner namespace-level config field placement
  - Configure TektonPruner with namespace-level settings for PipelineRuns, TaskRuns, and enforcedConfigLevel; verify fields are placed at the correct level (not inside the namespaces block)
  - Priority: **P1** | Tier: **Tier 1**

- **[SRVKP-10028](https://redhat.atlassian.net/browse/SRVKP-10028)** -- Resource-groups based PipelineRuns not pruned per configmap
  - Configure resource-group (label and annotation) based pruning in the configmap; verify PipelineRuns matching the resource-group criteria are pruned according to the specified configuration
  - Priority: **P1** | Tier: **Tier 1**

- **[SRVKP-10234](https://redhat.atlassian.net/browse/SRVKP-10234)** -- Global resolver-cache-config TTL overrides individual resolver TTL
  - Configure a global cache TTL and an individual resolver TTL; verify the individual resolver TTL takes precedence over the global setting
  - Priority: **P1** | Tier: **Tier 1**

- **[SRVKP-10235](https://redhat.atlassian.net/browse/SRVKP-10235)** -- PipelineRun cache parameter ignored when set to never
  - Create a PipelineRun with cache parameter set to "never"; verify the resolver bypasses the cache and fetches the task/pipeline definition fresh
  - Priority: **P1** | Tier: **Tier 1**

- **[SRVKP-10488](https://redhat.atlassian.net/browse/SRVKP-10488)** -- CLI version mismatch in opc
  - Run `opc version`; verify the reported versions for Pipelines as Code and Tekton Results match the actual server component versions
  - Priority: **P3** | Tier: **Tier 1**

---

### **IV. Sign-off and Approval**

This Software Test Plan requires approval from the following stakeholders:

* **Reviewers:**
  - [OpenShift Pipelines QE Lead]
  - [OpenShift Pipelines Dev Lead]
* **Approvers:**
  - [OpenShift Pipelines QE Manager]
  - [OpenShift Pipelines Engineering Manager]
