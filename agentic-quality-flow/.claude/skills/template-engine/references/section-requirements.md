# STP Section Requirements Reference

This document specifies the exact requirements for each STP section,
aligned with the current upstream template structure.

## Section Overview

| Section | Required | Format | Notes |
|:--------|:---------|:-------|:------|
| Document Header | Yes | Heading | Must be exact text |
| Feature Title | Yes | Heading | Must include "Quality Engineering Plan" |
| Metadata & Tracking | Yes | Bullet list, 6 items | "Current Status" is removed |
| Document Conventions | Yes | 1 line | "N/A" if not applicable |
| Feature Overview | Yes | Prose | 2-4 sentence description |
| Section I.1 | Yes | Checkbox list, 5 items | Sub-bullets for details |
| Section I.2 | Yes | Bullet list | Known Limitations (moved from old II.6) |
| Section I.3 | Yes | Checkbox list, 5 items | Technology and Design Review (was I.2) |
| Section II.1 | Yes | Mixed | Scope + Testing Goals (SMART) + Out of Scope (checkbox) |
| Section II.2 | Yes | Checkbox list, grouped | 4 groups; "Scale Testing" added, "Backward Compatibility" removed |
| Section II.3 | Yes | Bullet list, 10 items | All components required |
| Section II.3.1 | Optional | Bullet list, 3 items | Only NEW/SPECIAL tools |
| Section II.4 | Yes | Checkbox list | 2 standard + feature-specific |
| Section II.5 | Yes | Checkbox list, 6 categories + Other | Sub-items for mitigation and impact |
| Section III | Yes | Bullet list | Requirements-to-Tests Mapping |
| Section IV | Yes | - | TBD reviewers/approvers |

## Sections NOT in Template

Do NOT include:

- Related GitHub Pull Requests table (not in upstream template)
- Appendix, Summary, Glossary, References
- Sub-tables (A, B) under Test Strategy
- Section II.6 (Known Limitations moved to I.2; old II.6 is removed)

## Detailed Requirements

### Metadata & Tracking (6 items, bullet list)

Bullet list format (not a table). Each item is a single bullet.

- **Enhancement(s)** -- Link or "None"
- **Feature Tracking** -- From Feature Link custom field
- **Epic Tracking** -- Epic key + Parent key if exists (e.g., `Epic: CNV-xxxxx, Parent: VIRTSTRAT-xxx`)
- **QE Owner(s)** -- Name or TBD
- **Owning SIG** -- SIG name
- **Participating SIGs** -- List or "None"

"Current Status" is **removed** from the upstream template.

### Feature Overview

A brief (2-4 sentences) description of the feature being tested.
Include: what it does, why it matters to customers, and key technical components.

### Section I.1 - Requirement & User Story Review Checklist (5 items, checkbox list)

Checkbox list with sub-bullets for details (not a table). 5 items:

1. **Review Requirements** -- `- [ ] Review Requirements`
   - Sub-bullet: guidance on reviewing relevant requirements
2. **Understand Value and Customer Use Cases** -- `- [ ] Understand Value and Customer Use Cases`
   - Sub-bullet: understand the difference between U/S and D/S requirements, the value for
     RH customers, and ensure requirements contain relevant customer use cases
3. **Testability** -- `- [ ] Testability`
   - Sub-bullet: confirm requirements are testable and unambiguous
4. **Acceptance Criteria** -- `- [ ] Acceptance Criteria`
   - Sub-bullet: ensure acceptance criteria are defined clearly (clear user stories;
     D/S requirements clearly defined in Jira)
5. **Non-Functional Requirements** -- `- [ ] Non-Functional Requirements`
   - Sub-bullet: confirm coverage for NFRs including Performance, Security, Usability,
     Downtime, Connectivity, Monitoring (alerts/metrics), Scalability, Portability
     (e.g., cloud support), and Docs

**Key change from previous template:** "Understand Value" and "Customer Use Cases" are
merged into a single checkbox item. The total count drops from 6 to 5.

Each checkbox item has:

- A `[ ]` or `[x]` checkbox on the main line
- One or more indented sub-bullets with standard guidance text
- Optional feature-specific comments as additional sub-bullets

### Section I.2 - Known Limitations (bullet list)

Bullet list of known limitations and constraints for the feature.
This section was previously located at II.6 and has moved here.

### Section I.3 - Technology and Design Review (5 items, checkbox list)

Checkbox list with sub-bullets for details (not a table). 5 items:

1. **Developer Handoff/QE Kickoff** -- `- [ ] Developer Handoff/QE Kickoff`
   - Sub-bullet: guidance on handoff and kickoff activities
2. **Technology Challenges** -- `- [ ] Technology Challenges`
   - Sub-bullet: guidance on identifying technology challenges
3. **Test Environment Needs** -- `- [ ] Test Environment Needs`
   - Sub-bullet: guidance on environment requirements
4. **API Extensions** -- `- [ ] API Extensions`
   - Sub-bullet: guidance on API changes and extensions
5. **Topology Considerations** -- `- [ ] Topology Considerations`
   - Sub-bullet: guidance on topology and deployment considerations

Each checkbox item follows the same format as Section I.1: checkbox line with
indented sub-bullets for standard guidance text and optional feature-specific comments.

**Key change from previous template:** This section was previously I.2. It is now I.3,
after the insertion of Known Limitations at I.2.

### Section II.1 - Scope of Testing

Contains three parts in order:

1. **Scope description** -- brief paragraph of what will be tested
2. **Testing Goals** -- prioritized list using P0/P1/P2 format with SMART criteria.
   Goals should be organized into three groups when applicable:
   - **Functional Goals** -- specific functional verification goals
   - **Quality Goals** -- performance, reliability, quality attribute goals
   - **Integration Goals** -- compatibility, upgrade, integration goals
3. **Out of Scope (Testing Scope Exclusions)** -- checkbox format with rationale and
   PM/Lead sign-off per item: `- [ ] {exclusion} -- *Rationale:* {reason} -- *PM/Lead Agreement:* {name/date}`

**Key change from previous template:** Out of Scope uses checkbox format with
rationale and PM/Lead sign-off instead of a table.

### Section II.2 - Test Strategy (checkbox list, grouped)

Checkbox list grouped into four categories (not a table). Each item uses
`- [ ] {strategy}` format.

**Functional:**

- Functional Testing
- Automation Testing
- Regression Testing

**Non-Functional:**

- Performance Testing
- Scale Testing
- Security Testing
- Usability Testing
- Monitoring

**Integration & Compatibility:**

- Compatibility Testing
- Upgrade Testing
- Dependencies
- Cross Integrations

**Infrastructure:**

- Cloud Testing

**Key changes from previous template:**

- "Backward Compatibility Testing" is **removed** (merged into Compatibility Testing)
- "Scale Testing" is **added** under Non-Functional
- Format changed from a table with Applicable/Comments columns to a checkbox list
- Items are grouped by category instead of a flat numbered list

### Section II.3 - Test Environment (10 items, bullet list)

Bullet list format (not a table). 10 items:

1. Cluster Topology
2. OCP & OpenShift Virtualization Version(s)
3. CPU Virtualization
4. Compute Resources
5. Special Hardware
6. Storage
7. Network
8. Required Operators
9. Platform
10. Special Configurations

"N/A" means explicitly not applicable. Cannot leave items empty.

### Section II.3.1 - Testing Tools & Frameworks (3 items, bullet list)

Bullet list format (not a table). 3 items.

Only list tools that are **new** or **different** from standard testing infrastructure.
Standard tools (Ginkgo v2, Gomega, pytest, Prow, kubectl, oc, virtctl) should NOT be listed.
Leave items empty if using only standard tools.

### Section II.4 - Entry Criteria (checkbox list)

Checkbox format with two standard items plus feature-specific items:

- `- [ ] Requirements and design documents are **approved and merged**`
- `- [ ] Test environment can be **set up and configured** (see Section II.3 - Test Environment)`
- `- [ ] [Add feature-specific entry criteria as needed]`

### Section II.5 - Risks (checkbox list, 6 categories + Other)

Checkbox list with sub-items per risk category (not a table). 6 named categories
plus an Other category:

1. **Timeline/Schedule** -- `- [ ] Timeline/Schedule`
   - Sub-item: mitigation strategy
   - Sub-item: estimated impact
2. **Test Coverage** -- `- [ ] Test Coverage`
   - Sub-item: mitigation strategy
   - Sub-item: estimated impact
3. **Test Environment** -- `- [ ] Test Environment`
   - Sub-item: mitigation strategy
   - Sub-item: estimated impact
4. **Untestable Aspects** -- `- [ ] Untestable Aspects`
   - Sub-item: mitigation strategy
   - Sub-item: estimated impact
5. **Resource Constraints** -- `- [ ] Resource Constraints`
   - Sub-item: mitigation strategy
   - Sub-item: estimated impact
6. **Dependencies** -- `- [ ] Dependencies`
   - Sub-item: mitigation strategy
   - Sub-item: estimated impact
7. **Other** -- `- [ ] Other`
   - Sub-item: mitigation strategy
   - Sub-item: estimated impact

Each risk checkbox has:

- A `[ ]` or `[x]` checkbox on the main line with the risk description (or "N/A" with justification)
- Indented sub-item for mitigation strategy
- Indented sub-item for estimated impact

Do not duplicate information already in Test Environment section.

### Section III - Requirements-to-Tests Mapping (bullet list)

**No minimum item requirement.** Generate comprehensive coverage based on feature complexity.

Bullet-based format (not a table). Each requirement is a top-level bullet:

```markdown
- **[CNV-72329]** -- As a user, I want to hotplug a network interface to a running VM
  - *Test Scenario:* Verify hotplug attaches interface and traffic flows
  - *Priority:* P0
```

Each item must have:

- **Requirement ID** in bold brackets: `**[Jira-123]**`
- **Requirement Summary** after the dash: user story format preferred
- Indented **Test Scenario:** in italics: `*Test Scenario:*` followed by a brief phrase
- Indented **Priority:** in italics: `*Priority:*` followed by P0, P1, or P2

## Horizontal Rules

Place `---` in these locations:

1. After Feature Overview (before Section I)
2. After Section II.5 Risks (before Section III)
3. After Section III (before Section IV)

## Prohibited Sections

Do NOT add:

- Related GitHub Pull Requests table
- Appendix
- Summary
- Glossary
- References
- Any sections after Section IV
- Section II.6 (removed; Known Limitations is now at I.2)
