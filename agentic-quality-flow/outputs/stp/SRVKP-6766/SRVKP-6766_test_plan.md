# OpenShift Pipelines Test Plan

## **Show "All" Namespaces in Pipelines Overview for Non-Admins - Quality Engineering Plan**

### **Metadata & Tracking**

- **Enhancement(s):** [SRVKP-6766](https://redhat.atlassian.net/browse/SRVKP-6766)
- **Feature Tracking:** [SRVKP-6766 - Show "All" namespaces in Pipelines Overview for non-admins](https://redhat.atlassian.net/browse/SRVKP-6766)
- **Epic Tracking:** N/A
- **QE Owner(s):** TBD
- **Components:** Tekton Results, UI
- **Labels:** docs-pending, groomable, release-notes-pending, test-req, tests-pending
- **Product Version:** OpenShift Pipelines 1.21
- **Platform:** OCP 4.14+

**Document Conventions (if applicable):** This STP covers the feature request to enable non-admin users to select "All" namespaces in the Pipelines Overview page, showing an aggregated view of all namespaces the user has access to. The feature spans UI (console plugin) and backend (Tekton Results API) components.

### **Feature Overview**

Currently, on the Pipelines Overview page in the OpenShift Console, only cluster administrators can select "All" in the namespace dropdown to view an overview across all namespaces. Non-admin users (e.g., namespace admins or developers with access to multiple namespaces) are restricted to viewing one namespace at a time.

This feature extends the "All" namespace option to non-admin users. When a non-admin selects "All," the Pipelines Overview page should display an aggregated overview of all namespaces that the user has access to, as determined by their RBAC permissions. This eliminates the need for developers working across multiple namespaces to manually switch between each namespace to monitor their pipelines.

**Business Value:** Development teams operating across multiple namespaces gain consolidated visibility into pipeline status, reducing context-switching overhead and improving pipeline adoption.

---

### **I. Motivation and Requirements Review (QE Review Guidelines)**

This section documents the mandatory QE review process. The goal is to understand the feature's value, technology, and testability before formal test planning.

#### **1. Requirement & User Story Review Checklist**

- [ ] **Review Requirements**
  - Reviewed the relevant requirements.
  - The feature request (SRVKP-6766) describes the need for non-admin users to see an "All" namespaces option in the Pipelines Overview page. Requirements are derived from a customer case (support case tracking). The issue type is Feature with Normal priority.
- [ ] **Understand Value and Customer Use Cases**
  - Confirmed clear user stories and understood.
  - Understand the difference between U/S and D/S requirements.
  - **What is the value of the feature for RH customers**: Customers with development teams spanning multiple namespaces currently must check each namespace individually in the Pipelines Overview. This creates friction and reduces Pipelines adoption. The "All" namespace view consolidates monitoring across namespaces the user can access.
  - Ensured requirements contain relevant **customer use cases**.
  - Use case: A namespace admin managing pipelines in 5+ namespaces can view all pipeline activity from a single overview page instead of switching between namespaces individually.
- [ ] **Testability**
  - Confirmed requirements are **testable and unambiguous**.
  - The core requirement is testable: verify that the "All" option appears for non-admin users and that it shows data from exactly the namespaces the user has access to. Boundary conditions include users with access to zero, one, or many namespaces.
- [ ] **Acceptance Criteria**
  - Ensured acceptance criteria are **defined clearly** (clear user stories; D/S requirements clearly defined in Jira).
  - No explicit acceptance criteria defined in the Jira ticket. Derived acceptance criteria:
    1. Non-admin users see the "All" option in the namespace dropdown on the Pipelines Overview page
    2. Selecting "All" displays pipeline data from all namespaces the user has RBAC access to
    3. The overview does not expose data from namespaces the user lacks access to
    4. Admin users retain existing behavior (cluster-wide "All" view)
    5. Performance remains acceptable when aggregating across multiple namespaces
- [ ] **Non-Functional Requirements (NFRs)**
  - Confirmed coverage for NFRs, including Performance, Security, Usability, Downtime, Connectivity, Monitoring (alerts/metrics), Scalability, Portability (e.g., cloud support), and Docs.
  - **Performance:** Aggregating data across many namespaces must not cause unacceptable latency in the UI. Pagination or lazy loading may be needed.
  - **Security:** RBAC boundaries must be strictly enforced -- the user must never see pipeline data from namespaces they do not have access to.
  - **Usability:** The "All" namespace option must be discoverable and behave consistently with the existing admin experience.

#### **2. Known Limitations**

- The Jira ticket is in "New" status with no PRs or implementation details available yet. Test scenarios are based on the feature description and derived acceptance criteria.
- The ticket comments indicate this is primarily a UI-side change, but may also require backend (Tekton Results) changes to support cross-namespace queries scoped to user permissions.
- UI resource availability has been flagged as a constraint in the comments.
- No formal design document or architecture review has been shared.

#### **3. Technology and Design Review**

- [ ] **Developer Handoff/QE Kickoff**
  - A meeting where Dev/Arch walked QE through the design, architecture, and implementation details. **Critical for identifying untestable aspects early.**
  - Not yet completed. The feature is in refinement. A kickoff should cover: how the UI determines which namespaces the user has access to, whether the backend (Tekton Results) needs API changes, and how pagination is handled for cross-namespace queries.
- [ ] **Technology Challenges**
  - Identified potential testing challenges related to the underlying technology.
  - Testing requires setting up multiple users with varying RBAC roles across multiple namespaces. The test environment must support OpenShift RBAC configuration (creating users, roles, role bindings). The UI testing may require browser automation or API-level validation of the console plugin behavior.
- [ ] **Test Environment Needs**
  - Determined necessary **test environment setups and tools**.
  - Need an OCP cluster with multiple namespaces, multiple non-admin users with different namespace access patterns, and the OpenShift Pipelines operator installed. The tkn CLI and oc CLI are needed for creating test data (pipelines, pipeline runs) and RBAC configuration.
- [ ] **API Extensions**
  - Reviewed new or modified APIs and their impact on testing.
  - Potential API changes in Tekton Results to support listing results across namespaces scoped to user permissions. The console plugin API may be extended to support the "All" namespace query for non-admin users. Exact API changes TBD pending design review.
- [ ] **Topology Considerations**
  - Evaluated multi-cluster, network topology, and architectural impacts.
  - Single-cluster feature. No multi-cluster or network topology impacts expected. The feature operates within the OpenShift Console and Tekton Results within a single cluster.

### **II. Software Test Plan (STP)**

This STP serves as the **overall roadmap for testing**, detailing the scope, approach, resources, and schedule.

#### **1. Scope of Testing**

This test plan covers the "All" namespaces feature for non-admin users on the Pipelines Overview page. Testing will validate that non-admin users can view aggregated pipeline data across their accessible namespaces, that RBAC boundaries are enforced, and that existing admin functionality is preserved.

**Testing Goals**

- Verify the "All" namespace option is available to non-admin users on the Pipelines Overview page
- Validate that the aggregated view shows data only from namespaces the user has access to
- Confirm RBAC boundary enforcement -- no data leakage from inaccessible namespaces
- Validate existing admin "All" namespace functionality is not regressed
- Verify performance and usability of the aggregated view across multiple namespaces
- Validate upgrade path from previous versions where non-admins did not have this option

**Out of Scope (Testing Scope Exclusions)**

- [ ] Console UI framework internals (React/PatternFly rendering, general console plugin lifecycle)
- [ ] Tekton Results storage backend internals (database schema, gRPC internals)
- [ ] Pipelines pages other than the Overview page (Pipeline details, TaskRun logs, etc.)
- [ ] Multi-cluster federation or cross-cluster namespace aggregation
- [ ] Non-OpenShift Kubernetes distributions

#### **2. Test Strategy**

**Functional**

- [ ] **Functional Testing** -- Validates that the feature works according to specified requirements and user stories
  - *Details:* Test the "All" namespace option for non-admin users with various RBAC configurations. Validate data correctness in the aggregated overview (PipelineRun counts, status summaries, namespace grouping). Test with different numbers of accessible namespaces (1, few, many). Verify the dropdown selection and page rendering.
- [ ] **Automation Testing** -- Confirms test automation plan is in place for CI and regression coverage (all tests are expected to be automated)
  - *Details:* Tier 1 functional tests will be implemented in Go/Ginkgo targeting the Tekton Results API and RBAC validation. UI-level validation may require integration with the console plugin test framework. All tests will be automated for CI inclusion.
- [ ] **Regression Testing** -- Verifies that new changes do not break existing functionality
  - *Details:* Verify that admin "All" namespace behavior is unchanged. Verify that single-namespace selection for non-admins continues to work. Verify that the Pipelines Overview page loads correctly for users with no special permissions.

**Non-Functional**

- [ ] **Performance Testing** -- Validates feature performance meets requirements (latency, throughput, resource usage)
  - *Details:* Measure page load time for the "All" namespace view when a non-admin has access to 5, 10, and 50+ namespaces. Compare against single-namespace page load time. Verify no excessive API calls or memory usage in the browser.
- [ ] **Scale Testing** -- Validates feature behavior under increased load and at production-like scale
  - *Details:* Test with a user who has access to a large number of namespaces (50+) with significant pipeline history. Validate pagination or result limiting behavior under scale.
- [ ] **Security Testing** -- Verifies security requirements, RBAC, authentication, authorization, and vulnerability scanning
  - *Details:* Critical for this feature. Validate that the aggregated view strictly respects RBAC boundaries. Test with users who have mixed permissions (view in some namespaces, edit in others, none in others). Test for information disclosure vulnerabilities -- ensure namespace names and pipeline data from inaccessible namespaces are never returned. Test token scope and impersonation scenarios.
- [ ] **Usability Testing** -- Validates user experience and accessibility requirements
  - *Details:* Verify the "All" namespace option is clearly labeled and behaves consistently with the admin experience. Validate that the UI communicates which namespaces are included in the aggregated view.
- [ ] **Monitoring** -- Does the feature require metrics and/or alerts?
  - *Details:* No new metrics or alerts expected for this feature. Existing Tekton Results and console plugin metrics should continue to function.

**Integration & Compatibility**

- [ ] **Compatibility Testing** -- Ensures feature works across supported platforms, versions, and configurations
  - *Details:* Test on OCP 4.14+ with OpenShift Pipelines 1.21. Validate with OVN-Kubernetes CNI. Test with different authentication providers (htpasswd, LDAP, OIDC).
- [ ] **Upgrade Testing** -- Validates upgrade paths from previous versions
  - *Details:* Upgrade from OpenShift Pipelines 1.20 (where non-admins lack the "All" option) to 1.21 and verify the feature becomes available without manual configuration. Verify no data migration is needed.
- [ ] **Dependencies** -- Blocked by deliverables from other components/products
  - *Details:* Depends on OpenShift Console plugin framework supporting the namespace aggregation query. May depend on Tekton Results API changes for cross-namespace listing with RBAC scoping.
- [ ] **Cross Integrations** -- Does the feature affect other features or require testing by other teams?
  - *Details:* The Pipelines Overview page integrates with Tekton Results for data retrieval. Changes to namespace filtering may affect how Tekton Results processes queries. The Console team should be consulted on the UI plugin changes.

**Infrastructure**

- [ ] **Cloud Testing** -- Does the feature require multi-cloud platform testing?
  - *Details:* Not cloud-specific. Standard OCP cluster is sufficient. The feature behavior should be consistent across bare-metal, AWS, Azure, and GCP deployments.

#### **3. Test Environment**

- **Cluster Topology:** Single or Multi-node (minimum 1 worker node)
- **OCP & OpenShift Pipelines Version(s):** OCP 4.14+ / OpenShift Pipelines 1.21+
- **CPU Virtualization:** Standard (no special requirements)
- **Compute Resources:** Standard cluster sizing
- **Special Hardware:** None required
- **Storage:** Default storage class (standard)
- **Network:** OVN-Kubernetes
- **Required Operators:** Red Hat OpenShift Pipelines (namespace: openshift-pipelines)
- **Platform:** OpenShift Container Platform
- **Special Configurations:** Multiple non-admin users with varying RBAC role bindings across multiple namespaces

#### **3.1. Testing Tools & Frameworks**

- **Test Framework:** Ginkgo v2 + Gomega (Go, Tier 1)
- **CI/CD:** Prow / OpenShift CI
- **Other Tools:** tkn CLI, oc CLI, kubectl

#### **4. Entry Criteria**

The following conditions must be met before testing can begin:

- [ ] Requirements and design documents are **approved and merged**
- [ ] Test environment can be **set up and configured** (see Section II.3 - Test Environment)
- [ ] Developer handoff/QE kickoff meeting completed
- [ ] Feature implementation PRs are merged (or available for testing in a staging environment)
- [ ] OpenShift Pipelines operator with this feature is deployable on the target OCP cluster
- [ ] Non-admin user RBAC test fixtures are documented and reproducible

#### **5. Risks**

- [ ] **Timeline/Schedule**
  - Risk: Feature is in "New" status with no PRs. UI resource constraints have been flagged. Implementation timeline is uncertain.
  - Mitigation: Begin test planning early (this STP). Engage with the UI team to align on timelines. Prepare test fixtures and automation framework in advance.
- [ ] **Test Coverage**
  - Risk: Without a detailed design document, test scenarios may miss implementation-specific edge cases (e.g., how namespace list is retrieved, caching behavior).
  - Mitigation: Schedule a QE kickoff meeting with the development team. Update test scenarios after the design review.
- [ ] **Test Environment**
  - Risk: Setting up multiple users with specific RBAC configurations across many namespaces requires careful orchestration.
  - Mitigation: Create reusable test fixtures (scripts/manifests) for RBAC setup. Document the environment configuration steps.
- [ ] **Untestable Aspects**
  - Risk: UI rendering and browser-specific behavior may be difficult to validate in automated Tier 1 tests.
  - Mitigation: Focus Tier 1 tests on API-level validation (Tekton Results queries with RBAC). Supplement with manual UI verification during development.
- [ ] **Resource Constraints**
  - Risk: Limited UI development resources may delay feature delivery and testing.
  - Mitigation: Prioritize API-level test automation that can proceed independently of UI completion.
- [ ] **Dependencies**
  - Risk: The feature may require changes in both the console plugin and Tekton Results backend. Misalignment between teams could delay integration testing.
  - Mitigation: Coordinate with both the Console/UI team and the Tekton Results team. Establish integration test milestones.

---

### **III. Test Scenarios & Traceability**

This section links requirements to test coverage, enabling reviewers to verify all requirements are tested.

#### **1. Requirements-to-Tests Mapping**

- **[SRVKP-6766]** -- Non-admin users can select "All" namespaces in the Pipelines Overview page to view aggregated pipeline data from all accessible namespaces
  - TS-01: Verify "All" namespace option is visible for non-admin users on the Pipelines Overview page -- **Tier 1**
  - TS-02: Verify non-admin "All" namespace view shows PipelineRun data from all accessible namespaces -- **Tier 1**
  - TS-03: Verify non-admin "All" namespace view excludes data from namespaces the user lacks access to (RBAC enforcement) -- **Tier 1**
  - TS-04: Verify admin "All" namespace view behavior is unchanged (regression) -- **Tier 1**
  - TS-05: Verify non-admin user with access to a single namespace sees the "All" option (shows same data as single namespace view) -- **Tier 1**
  - TS-06: Verify non-admin user with access to zero namespaces sees the "All" option with an empty overview or appropriate message -- **Tier 1**
  - TS-07: Verify namespace dropdown correctly reflects user permissions when switching between "All" and specific namespaces -- **Tier 1**
  - TS-08: Verify PipelineRun status counts (succeeded, failed, running) are correctly aggregated across namespaces in the "All" view -- **Tier 1**
  - TS-09: Verify that revoking namespace access dynamically updates the "All" namespace view (user no longer sees data from revoked namespace) -- **Tier 1**
  - TS-10: Verify performance of "All" namespace view for non-admin user with access to 10+ namespaces (page load within acceptable threshold) -- **Tier 1**
  - TS-11: Verify the "All" namespace view after upgrading from OpenShift Pipelines 1.20 to 1.21 (feature becomes available without manual steps) -- **Tier 1**
  - TS-12: Verify RBAC boundary enforcement -- non-admin user with view-only access to some namespaces and edit access to others sees correct data in the aggregated view -- **Tier 1**
  - TS-13: Verify that namespace names from inaccessible namespaces are not disclosed in API responses or error messages (security) -- **Tier 1**

---

### **IV. Sign-off and Approval**

This Software Test Plan requires approval from the following stakeholders:

* **Reviewers:**
  - [TBD]
  - [TBD]
* **Approvers:**
  - [TBD]
  - [TBD]
