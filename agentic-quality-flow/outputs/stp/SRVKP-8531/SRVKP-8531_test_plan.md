# OpenShift Pipelines Test Plan

## **Migrate Chains from OpenCensus to OpenTelemetry - Quality Engineering Plan**

### **Metadata & Tracking**

- **Enhancement(s):** [tektoncd/pipeline#8969](https://github.com/tektoncd/pipeline/issues/8969)
- **Feature Tracking:** [SRVKP-8531](https://redhat.atlassian.net/browse/SRVKP-8531)
- **Epic Tracking:** [SRVKP-8531](https://redhat.atlassian.net/browse/SRVKP-8531)
- **QE Owner(s):** QE Team
- **Owning SIG:** chains
- **Participating SIGs:** pipelines, triggers, operator

**Document Conventions (if applicable):** This STP covers the migration of Tekton Chains from OpenCensus to OpenTelemetry for observability instrumentation. Related sibling tickets cover Triggers (SRVKP-8532) and Results (SRVKP-8533) migrations.

### **Feature Overview**

Tekton Chains currently uses OpenCensus for metrics and tracing instrumentation. OpenCensus has been deprecated in favor of OpenTelemetry (OTel), the CNCF standard for observability. This story migrates Chains from OpenCensus dependencies to OpenTelemetry, executed in two stages:

1. **Stage 1:** Catalogue existing OpenCensus metrics and map their usage in the Konflux dashboard
2. **Stage 2:** Upgrade Knative dependencies, replace OpenCensus with OpenTelemetry, publish documentation and release notes

The upstream reference PR ([tektoncd/pipeline#8933](https://github.com/tektoncd/pipeline/pull/8933)) demonstrates the migration pattern for the pipeline component, which serves as the blueprint for the Chains migration. Key changes include replacing OpenCensus metric views with OpenTelemetry metric instruments, updating config-observability configuration, and bumping Knative dependencies.

**Dependency:** This work is blocked by [SRVKP-9322](https://redhat.atlassian.net/browse/SRVKP-9322) ([Chains] bump knative.dev/pkg), which must be completed first to provide the OpenTelemetry-compatible Knative base.

---

### **I. Motivation and Requirements Review (QE Review Guidelines)**

This section documents the mandatory QE review process. The goal is to understand the feature's value,
technology, and testability before formal test planning.

#### **1. Requirement & User Story Review Checklist**

- [ ] **Review Requirements**
  - Reviewed the relevant requirements.
  - The Jira description references upstream issue [tektoncd/pipeline#8969](https://github.com/tektoncd/pipeline/issues/8969) as the driving requirement. OpenCensus is deprecated and must be replaced with OpenTelemetry across all Tekton components. The acceptance criteria in the Jira are placeholder templates and need to be filled in by the development team before test execution.
- [ ] **Understand Value and Customer Use Cases**
  - Confirmed clear user stories and understood.
  - Understand the difference between U/S and D/S requirements.
  - **What is the value of the feature for RH customers**: Customers relying on Tekton Chains metrics for supply chain security monitoring (via Konflux dashboards or custom Prometheus/Grafana setups) need continued metric availability after the migration. OpenTelemetry is the industry-standard successor to OpenCensus, ensuring long-term support and compatibility with modern observability backends.
  - Ensured requirements contain relevant **customer use cases**.
  - Primary use case: Platform engineers monitoring Chains signing, attestation, and storage operations through Prometheus metrics. The migration must preserve all existing metric names, labels, and semantics to avoid breaking dashboards and alerts.
- [ ] **Testability**
  - Confirmed requirements are **testable and unambiguous**.
  - The migration is testable by verifying that all previously exposed OpenCensus metrics remain available with identical names and label dimensions after the switch to OpenTelemetry. Metric endpoint scraping and Prometheus query validation provide concrete, automatable verification.
- [ ] **Acceptance Criteria**
  - Ensured acceptance criteria are **defined clearly** (clear user stories; D/S requirements clearly defined in Jira).
  - **Note:** The Jira acceptance criteria section contains only placeholder text. QE has derived testable acceptance criteria from the upstream issue and reference PR: (1) All existing Chains metrics must be preserved, (2) No OpenCensus imports remain in Chains codebase, (3) config-observability ConfigMap supports OTel configuration, (4) Knative dependency is bumped to OTel-compatible version.
- [ ] **Non-Functional Requirements (NFRs)**
  - Confirmed coverage for NFRs, including Performance, Security, Usability, Downtime, Connectivity, Monitoring (alerts/metrics), Scalability, Portability (e.g., cloud support), and Docs.
  - **Performance:** Metric collection overhead must not increase beyond OpenCensus baseline. **Monitoring:** This feature directly impacts metrics and monitoring -- all existing Prometheus scrape targets and metric names must remain functional. **Docs:** Downstream documentation and release notes are required per the Jira description.

#### **2. Known Limitations**

- The Jira acceptance criteria section is not filled in -- QE-derived criteria are used for this STP
- The upstream reference PR (tektoncd/pipeline#8933) is in CLOSED/WIP state, meaning the exact migration pattern may evolve
- The migration depends on SRVKP-9322 (Knative bump) being completed first
- Metric semantic equivalence between OpenCensus views and OpenTelemetry instruments may have subtle differences in histogram bucket boundaries or aggregation behavior
- Conflux dashboard compatibility is a cross-team dependency that cannot be fully validated in isolation

#### **3. Technology and Design Review**

- [ ] **Developer Handoff/QE Kickoff**
  - A meeting where Dev/Arch walked QE through the design, architecture, and implementation details. **Critical for identifying untestable aspects early.**
  - Kickoff should cover: the OpenCensus-to-OpenTelemetry migration pattern from the upstream pipeline PR, Chains-specific metrics (signing duration, attestation count, storage latency), and Knative dependency chain.
- [ ] **Technology Challenges**
  - Identified potential testing challenges related to the underlying technology.
  - OpenTelemetry uses a different SDK initialization pattern than OpenCensus. The Chains metrics are registered through Knative's metrics package, which abstracts the backend. Testing must verify both the Prometheus exporter path (production) and any OTLP exporter configuration (optional). The Knative metrics library version must match the OTel-compatible version from the bumped knative.dev/pkg.
- [ ] **Test Environment Needs**
  - Determined necessary **test environment setups and tools**.
  - Standard OCP cluster with OpenShift Pipelines operator installed. Prometheus must be accessible for metric scraping validation. No special hardware required. The `tkn` CLI and `oc` CLI are needed for pipeline and chains operations.
- [ ] **API Extensions**
  - Reviewed new or modified APIs and their impact on testing.
  - No new CRD APIs. The config-observability ConfigMap gains new OpenTelemetry-specific keys (e.g., `metrics.backend-destination`, `metrics.opencensus-address` replaced by `metrics.otel-collector-address`). The TektonConfig CR may expose new observability configuration fields through the operator.
- [ ] **Topology Considerations**
  - Evaluated multi-cluster, network topology, and architectural impacts.
  - Single-cluster topology is sufficient for functional validation. The metrics endpoint is local to the Chains controller pod. No multi-cluster or network topology concerns.

### **II. Software Test Plan (STP)**

This STP serves as the **overall roadmap for testing**, detailing the scope, approach, resources, and schedule.

#### **1. Scope of Testing**

This test plan covers the migration of Tekton Chains observability instrumentation from OpenCensus to OpenTelemetry within OpenShift Pipelines 1.21. Testing focuses on verifying that Chains metrics remain functional, accurate, and backward-compatible after the dependency migration.

**Testing Goals**

- Verify all existing Chains metrics are preserved with identical names and label dimensions after the OpenTelemetry migration
- Validate that the config-observability ConfigMap correctly configures OpenTelemetry backends
- Confirm Chains signing, attestation, and storage operations continue to emit metrics under load
- Ensure no OpenCensus dependency remnants cause runtime errors or metric gaps
- Validate that Prometheus scraping of the Chains controller metrics endpoint returns expected metric families
- Verify Knative dependency bump does not introduce regressions in Chains reconciliation behavior

**Out of Scope (Testing Scope Exclusions)**

- [ ] Konflux dashboard UI validation -- owned by the Konflux team
- [ ] Triggers migration to OpenTelemetry -- covered by SRVKP-8532
- [ ] Results migration to OpenTelemetry -- covered by SRVKP-8533
- [ ] OpenTelemetry Collector deployment and configuration -- infrastructure concern
- [ ] OTLP exporter testing -- Prometheus exporter is the primary production path
- [ ] Performance benchmarking of OpenTelemetry vs OpenCensus overhead -- informational only, not a gate

#### **2. Test Strategy**

**Functional**

- [ ] **Functional Testing** -- Validates that the feature works according to specified requirements and user stories
  - *Details:* Verify each Chains metric (signing duration, attestation count, storage operation latency, error counts) is emitted correctly after the OTel migration. Trigger Chains operations (TaskRun completion with signing configured) and scrape the metrics endpoint to confirm values.
- [ ] **Automation Testing** -- Confirms test automation plan is in place for CI and regression coverage (all tests are expected to be automated)
  - *Details:* All test scenarios will be implemented as Tier 1 Go/Ginkgo tests. Tests will use the Tekton Pipelines test framework to create TaskRuns, trigger Chains signing, and validate metrics via HTTP scraping of the controller pod's metrics port.
- [ ] **Regression Testing** -- Verifies that new changes do not break existing functionality
  - *Details:* Existing Chains functional tests (signing, attestation, storage) must continue to pass. The metrics migration must not alter Chains reconciliation behavior. Regression test suite should be run against both the pre-migration and post-migration builds.

**Non-Functional**

- [ ] **Performance Testing** -- Validates feature performance meets requirements (latency, throughput, resource usage)
  - *Details:* Measure metric collection overhead by comparing Chains controller CPU and memory usage before and after migration. OpenTelemetry SDK should not introduce measurable latency to the signing/attestation hot path.
- [ ] **Scale Testing** -- Validates feature behavior under increased load and at production-like scale (e.g., large number of TaskRuns or concurrent signing operations)
  - *Details:* Not a primary focus for this migration. Existing scale tests should be re-run to confirm no regression.
- [ ] **Security Testing** -- Verifies security requirements, RBAC, authentication, authorization, and vulnerability scanning
  - *Details:* Verify that the new OpenTelemetry dependencies do not introduce known CVEs. The metrics endpoint should remain accessible only within the cluster (no new network exposure).
- [ ] **Usability Testing** -- Validates user experience and accessibility requirements
  - *Details:* Not applicable for this infrastructure migration.
- [ ] **Monitoring** -- Does the feature require metrics and/or alerts?
  - *Details:* This feature IS a metrics migration. All existing Chains monitoring must continue to function. Verify that Prometheus ServiceMonitor or PodMonitor configurations continue to scrape the Chains controller correctly.

**Integration & Compatibility**

- [ ] **Compatibility Testing** -- Ensures feature works across supported platforms, versions, and configurations
  - *Details:* Verify on OCP 4.14+ with OpenShift Pipelines 1.21. Validate that the metrics endpoint works with both in-cluster Prometheus and user-workload monitoring.
- [ ] **Upgrade Testing** -- Validates upgrade paths from previous versions, data migration, and configuration preservation
  - *Details:* Upgrade from OpenShift Pipelines 1.20 (OpenCensus) to 1.21 (OpenTelemetry). Verify that existing config-observability ConfigMap settings are either preserved or gracefully migrated. Confirm that Prometheus queries using old metric names continue to work.
- [ ] **Dependencies** -- Blocked by deliverables from other components/products. Identify what we need from other teams before we can test.
  - *Details:* Blocked by SRVKP-9322 ([Chains] bump knative.dev/pkg). The Knative package bump must be completed and merged before the OTel migration can proceed. The operator team must update the TektonChain CR reconciler if observability config fields change.
- [ ] **Cross Integrations** -- Does the feature affect other features or require testing by other teams? Identify the impact we cause.
  - *Details:* The Triggers team (SRVKP-8532) and Results team (SRVKP-8533) are performing parallel migrations. The operator component must handle the new observability configuration. Konflux dashboard consumers must validate their Prometheus queries still work.

**Infrastructure**

- [ ] **Cloud Testing** -- Does the feature require multi-cloud platform testing? Consider cloud-specific features.
  - *Details:* No cloud-specific behavior. The metrics endpoint is standard Prometheus format regardless of cloud provider.

#### **3. Test Environment**

- **Cluster Topology:** Single or Multi-node (single-node sufficient for functional validation)
- **OCP & OpenShift Pipelines Version(s):** OCP 4.14+ with OpenShift Pipelines 1.21+
- **CPU Virtualization:** Standard (no special CPU features required)
- **Compute Resources:** 1+ worker nodes with default resource allocation
- **Special Hardware:** None required
- **Storage:** Default storage class (standard)
- **Network:** OVN-Kubernetes CNI (default)
- **Required Operators:** Red Hat OpenShift Pipelines (namespace: openshift-pipelines)
- **Platform:** OpenShift Container Platform (OCP)
- **Special Configurations:** Chains signing configured (Cosign or x509), Prometheus accessible for metrics scraping

#### **3.1. Testing Tools & Frameworks**

- **Test Framework:** Ginkgo v2 + Gomega (Go)
- **CI/CD:** Prow, OpenShift CI
- **Other Tools:** tkn CLI, oc CLI, curl/wget (metrics endpoint scraping), promtool (PromQL validation)

#### **4. Entry Criteria**

The following conditions must be met before testing can begin:

- [ ] Requirements and design documents are **approved and merged**
- [ ] Test environment can be **set up and configured** (see Section II.3 - Test Environment)
- [ ] SRVKP-9322 ([Chains] bump knative.dev/pkg) is completed and merged
- [ ] Chains OpenTelemetry migration code is merged and available in a build
- [ ] OpenShift Pipelines operator includes the updated Chains component
- [ ] Chains signing is configurable (Cosign or x509 key pair available)
- [ ] Prometheus is accessible in the test cluster for metric scraping validation

#### **5. Risks**

- [ ] **Timeline/Schedule**
  - Risk: The dependency on SRVKP-9322 (Knative bump) may delay the start of testing if it is not completed on time
  - Mitigation: Monitor the blocker ticket status; prepare test automation against the upstream pipeline migration pattern so tests are ready when the Chains migration lands
- [ ] **Test Coverage**
  - Risk: Incomplete catalogue of existing Chains metrics may result in missed regression
  - Mitigation: Stage 1 of the Jira (SRVKP-8560, now Closed) should have produced a comprehensive metric catalogue; use that as the test baseline
- [ ] **Test Environment**
  - Risk: Prometheus configuration in test clusters may differ from production Konflux environments
  - Mitigation: Use standard OpenShift user-workload monitoring as the reference Prometheus instance; document any cluster-specific configuration
- [ ] **Untestable Aspects**
  - Risk: Metric semantic equivalence (e.g., histogram bucket boundaries) between OpenCensus and OpenTelemetry may have subtle differences that are hard to detect in automated tests
  - Mitigation: Include explicit histogram bucket validation in test scenarios; compare bucket boundaries against the known OpenCensus defaults
- [ ] **Resource Constraints**
  - Risk: QE team bandwidth may be constrained by parallel Triggers (SRVKP-8532) and Results (SRVKP-8533) migrations
  - Mitigation: Share test patterns and helper utilities across the three migration STPs; use a common metrics validation library
- [ ] **Dependencies**
  - Risk: Operator component may not expose new OTel configuration fields in TektonConfig CR in time
  - Mitigation: Test with direct ConfigMap modification as a fallback; validate operator integration separately when available
- [ ] **Other**
  - Risk: The upstream reference PR (tektoncd/pipeline#8933) is in CLOSED/WIP state, and the final migration pattern may differ
  - Mitigation: Track upstream progress; update test scenarios if the migration approach changes before the Chains implementation is finalized

---

### **III. Test Scenarios & Traceability**

This section links requirements to test coverage, enabling reviewers to verify all requirements are tested.

#### **1. Requirements-to-Tests Mapping**

- **[SRVKP-8531]** -- Migrate Chains from OpenCensus to OpenTelemetry
  - TS-SRVKP-8531-001: Verify Chains controller metrics endpoint serves metrics in Prometheus format after OTel migration -- **Tier 1**
  - TS-SRVKP-8531-002: Verify all previously registered Chains OpenCensus metrics are present with identical names under OpenTelemetry -- **Tier 1**
  - TS-SRVKP-8531-003: Verify Chains signing duration metric is emitted with correct labels after a successful TaskRun signing operation -- **Tier 1**
  - TS-SRVKP-8531-004: Verify Chains attestation count metric increments correctly for each signed TaskRun -- **Tier 1**
  - TS-SRVKP-8531-005: Verify Chains storage operation latency metric is recorded when storing attestations to OCI registry -- **Tier 1**
  - TS-SRVKP-8531-006: Verify Chains error count metric increments on signing failure (e.g., invalid key configuration) -- **Tier 1**
  - TS-SRVKP-8531-007: Verify config-observability ConfigMap accepts OpenTelemetry configuration keys (metrics.backend-destination, request-metrics-backend-destination) -- **Tier 1**
  - TS-SRVKP-8531-008: Verify Chains controller starts successfully with default config-observability settings after OTel migration -- **Tier 1**
  - TS-SRVKP-8531-009: Verify no OpenCensus import paths remain in the Chains binary (dependency audit) -- **Tier 1**
  - TS-SRVKP-8531-010: Verify Prometheus ServiceMonitor/PodMonitor successfully scrapes the Chains controller metrics endpoint -- **Tier 1**
  - TS-SRVKP-8531-011: Verify histogram bucket boundaries for duration metrics match the pre-migration OpenCensus configuration -- **Tier 1**
  - TS-SRVKP-8531-012: Verify Chains reconciliation behavior is unaffected by the OTel migration (TaskRun signing completes normally) -- **Tier 1**

- **[SRVKP-9322]** -- [Chains] bump knative.dev/pkg (blocker dependency)
  - TS-SRVKP-8531-013: Verify Chains controller starts and reconciles correctly after knative.dev/pkg bump -- **Tier 1**
  - TS-SRVKP-8531-014: Verify Knative metrics package initialization uses OpenTelemetry backend after the bump -- **Tier 1**

- **[SRVKP-8531 - Upgrade Path]** -- Upgrade from OpenCensus to OpenTelemetry
  - TS-SRVKP-8531-015: Verify upgrade from OpenShift Pipelines 1.20 to 1.21 preserves existing config-observability ConfigMap settings -- **Tier 1**
  - TS-SRVKP-8531-016: Verify Prometheus queries using existing metric names return results after the upgrade -- **Tier 1**

- **[SRVKP-8531 - Negative/Edge Cases]** -- Error handling and edge cases
  - TS-SRVKP-8531-017: Verify Chains controller handles invalid config-observability values gracefully (does not crash, logs warning) -- **Tier 1**
  - TS-SRVKP-8531-018: Verify metrics are still emitted when the OTel collector endpoint is unreachable (fallback to Prometheus exporter) -- **Tier 1**

---

### **IV. Sign-off and Approval**

This Software Test Plan requires approval from the following stakeholders:

* **Reviewers:**
  - [Name / @github-username]
  - [Name / @github-username]
* **Approvers:**
  - [Name / @github-username]
  - [Name / @github-username]
