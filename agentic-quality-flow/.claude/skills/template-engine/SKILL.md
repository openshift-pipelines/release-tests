---
name: template-engine
description: Apply the official STP template structure to generated content
---

# Template Engine Skill

**Phase:** Utility
**User-Invocable:** false

## Purpose

Apply the official STP template structure to generated content.

## When to Use

Invoked by the **stp-generator** subagent to structure the final document.

## Template Location

Full template: Read STP template from `{project_context.config_dir}/templates/stp/stp-template.md` with fallback to `templates/stp-template.md` (relative to this skill).
Section requirements: `{project_context.config_dir}/templates/stp/section-requirements.md` with fallback to `references/section-requirements.md` (relative to this skill).

## Document Structure

The STP MUST contain sections in this EXACT order:

```
1. Document Header: `{project_context.stp_header}` (e.g., "# Openshift-virtualization-tests Test plan")
2. Feature Title: ## **[Title] - Quality Engineering Plan**
3. Metadata & Tracking (bullet list, 6 items)
4. Document Conventions (if applicable)
5. Feature Overview (2-4 sentence description)
6. ---
7. Section I.1 - Requirement & User Story Review Checklist (5 checkbox items)
8. Section I.2 - Known Limitations (free text / bullet list)
9. Section I.3 - Technology and Design Review (5 checkbox items)
10. Section II.1 - Scope of Testing + Testing Goals + Out of Scope (checkbox format)
11. Section II.2 - Test Strategy (categorized checkbox groups)
12. Section II.3 - Test Environment (bullet list, 10 items)
13. Section II.3.1 - Testing Tools & Frameworks (only NEW/SPECIAL tools)
14. Section II.4 - Entry Criteria (checkbox format)
15. Section II.5 - Risks (checkbox format with sub-items)
16. ---
17. Section III.1 - Requirements-to-Tests Mapping (bullet-based)
18. ---
19. Section IV - Sign-off
```

## Key Design Rules

### Section I is a Meta-Checklist

Section I items confirm the QE review PROCESS was followed. Each item is a checkbox
with the **verbatim standard guidance text from the upstream template** as the main label.
Feature-specific observations go in sub-item details beneath each checkbox, not in a
separate column.

**CRITICAL: The checkbox label text is NOT a feature-content field.** The text must
match the upstream template exactly:
- "Reviewed the relevant requirements."
- "Confirmed clear user stories and understood. Understand the value and customer use cases."
- "Confirmed requirements are **testable and unambiguous**."
- "Ensured acceptance criteria are **defined clearly**."
- "Confirmed coverage for NFRs."
- etc.

Do NOT replace checkbox labels with:
- Feature-specific acceptance criteria
- Technical requirement descriptions
- Feature-specific value propositions
- Detailed testability assessments
- PR references or implementation details

Feature-specific observations are added as indented sub-items below each checkbox.

### Section I.2 is Known Limitations

Known Limitations has moved from the old Section II.6 to Section I.2. This section
uses free text or a bullet list to describe known limitations and constraints
discovered during the motivation and requirements review.

### Section I.3 is Technology and Design Review

What was previously Section I.2 is now Section I.3. It uses checkbox format (not a table)
with 5 items. Each checkbox has the standard guidance text as its label, with
feature-specific observations as indented sub-items.

### Section II.1 Merges Scope, Goals, and Out of Scope

This is a single section containing:
1. A scope description paragraph
2. **Testing Goals** with priority levels (P0/P1/P2) using SMART criteria
3. **Out of Scope (Testing Scope Exclusions)** in checkbox format with rationale as sub-items

### Section II.2 Uses Categorized Checkboxes

The test strategy uses categorized checkbox groups instead of a single table.
Strategy items are organized by category (e.g., Core Testing, Extended Testing,
Integration & Operations) with each item as a checkbox. Sub-items provide the
description and applicability details.

There are 13 items total across four groups:

**Functional:** Functional Testing, Automation Testing, Regression Testing
**Non-Functional:** Performance Testing, Scale Testing, Security Testing, Usability Testing, Monitoring
**Integration & Compatibility:** Compatibility Testing, Upgrade Testing, Dependencies, Cross Integrations
**Infrastructure:** Cloud Testing

### Section II.3.1 Lists Only New/Special Tools

Only list tools that are **new** or **different** from standard testing infrastructure.
Standard tools (Ginkgo, pytest, Prow, kubectl, virtctl) should NOT be listed.
Leave empty if using only standard tools.

### No Related GitHub PRs Table

The upstream template does not include a Related GitHub Pull Requests table.
Do not add this section.

## Required Structure Counts

| Section | Format | Required Items |
|:--------|:-------|:---------------|
| Metadata | Bullet list | 6 (Enhancement, Feature Tracking, Epic Tracking, QE Owner, Owning SIG, Participating SIGs) |
| I.1 Requirement Review | Checkbox list | 5 (Review Requirements, Understand Value and Customer Use Cases, Testability, Acceptance Criteria, NFRs) |
| I.2 Known Limitations | Free text / bullets | At least 1 item or "None identified" |
| I.3 Technology Review | Checkbox list | 5 (Developer Handoff, Technology Challenges, Test Environment Needs, API Extensions, Topology) |
| II.1 Out of Scope | Checkbox list | 1+ items or "None" |
| II.2 Test Strategy | Categorized checkboxes | 13 items across categories |
| II.3 Test Environment | Bullet list | 10 (Cluster Topology, OCP Version, CPU Virtualization, Compute, Special Hardware, Storage, Network, Operators, Platform, Special Configs) |
| II.5 Risks | Checkbox with sub-items | 7 (Timeline, Coverage, Environment, Untestable, Resources, Dependencies, Other) |
| III.1 Requirements Mapping | Bullet-based | No minimum; comprehensive coverage |

## Bullet and Checkbox Formatting

- Metadata uses `- **Field:** Value` format
- Checkbox items use `- [ ] **Label** -- guidance text` with sub-items for details
- Risk items use `- [ ] **Category**` with indented sub-items for risk, mitigation, and status
- Test environment uses `- **Component:** configuration details`
- Section III uses bullet items with requirement ID, summary, scenarios, tier, and priority

## Section Headers

Use exact markdown header levels:
- `#` for document title
- `##` for feature title
- `###` for major sections (Metadata, Feature Overview, I, II, III, IV)
- `####` for subsections (1, 2, 3...)

## Horizontal Rules

Place `---` horizontal rules:
- After Feature Overview (before Section I)
- After Section II.5 (before Section III) -- note: II.6 no longer exists
- After Section III (before Section IV)

## Input

```yaml
content:
  metadata:
    enhancement: <link>
    feature_in_jira: <link>
    jira_tracking: <link>
    qe_owner: <name>
    owning_sig: <sig>
    participating_sigs: [...]
    document_conventions: <text or "N/A">

  feature_overview: <2-4 sentence description>

  section_i:
    requirement_review:
      - check: Review Requirements
        done: "[ ]"
        details: <feature-specific observations as sub-items>
      - check: Understand Value and Customer Use Cases
        done: "[ ]"
        details: <feature-specific observations as sub-items>
      - check: Testability
        done: "[ ]"
        details: <feature-specific observations as sub-items>
      - check: Acceptance Criteria
        done: "[ ]"
        details: <feature-specific observations as sub-items>
      - check: Non-Functional Requirements (NFRs)
        done: "[ ]"
        details: <feature-specific observations as sub-items>
    known_limitations:
      - <limitation text>
      - ...
    technology_review:
      - check: Developer Handoff
        done: "[ ]"
        details: <feature-specific observations as sub-items>
      - ...

  section_ii:
    scope: <text>
    testing_goals: <prioritized P0/P1/P2 list>
    out_of_scope:
      - item: <item>
        rationale: <rationale>
        agreement: "[ ] Name/Date"
      - ...
    test_strategy:
      - category: Core Testing
        items:
          - item: Functional Testing
            applicable: Y
            description: <description>
          - item: Automation Testing
            applicable: Y
            description: <description>
          - ...
      - category: Extended Testing
        items:
          - ...
      - category: Integration & Operations
        items:
          - ...
    environment:
      - component: Cluster Topology
        config: <configuration details>
      - ...
    tools:
      - category: Test Framework
        tools: <only NEW/SPECIAL tools or empty>
      - ...
    entry_criteria:
      - <extra feature-specific criteria>
      - ...
    risks:
      - category: Timeline
        specific_risk: <risk>
        mitigation: <mitigation>
        status: "[ ]"
      - ...

  section_iii:
    requirements_mapping:
      - requirement_id: <Jira issue key>
        requirement_summary: <summary>
        test_scenarios: <scenarios>
        tier: <Tier 1 or Tier 2>
        priority: <priority>
      - ...

  section_iv:
    reviewers: [...]
    approvers: [...]
```

## Output

Complete STP markdown document following the exact template structure.

## Validation Before Output

- [ ] Document starts with: `{project_context.stp_header}` (read from project config)
- [ ] Feature title format: `## **[Title] - Quality Engineering Plan**`
- [ ] Feature Overview section present (2-4 sentences)
- [ ] Document Conventions line present
- [ ] No Related GitHub PRs table
- [ ] Metadata is a bullet list with 6 items (no "Current Status" field)
- [ ] Section I.1 has 5 checkbox items (merged "Understand Value and Customer Use Cases")
- [ ] Section I.1 checkbox labels use **verbatim** standard template text (not feature-specific content)
- [ ] Section I.1 feature-specific observations are in sub-items (not in the checkbox label)
- [ ] Section I.2 is Known Limitations (not Technology Review)
- [ ] Section I.3 is Technology and Design Review with 5 checkbox items
- [ ] Section II.1 contains Scope + Testing Goals + Out of Scope (checkbox format)
- [ ] Test Strategy uses categorized checkbox groups with 13 items total
- [ ] Test Environment is a bullet list with 10 items
- [ ] Testing Tools lists only NEW/SPECIAL tools
- [ ] Entry Criteria uses checkbox format
- [ ] Risks use checkbox format with sub-items for risk, mitigation, and status
- [ ] Section II.6 does NOT exist (Known Limitations moved to I.2)
- [ ] Section III uses bullet-based format (not table)
- [ ] Each Section III item has exactly one tier
- [ ] Requirement IDs are Jira issue keys
- [ ] Horizontal rules in correct positions (after Overview, after II.5, after III)
- [ ] No extra sections added

## Prohibited Content

- Related GitHub Pull Requests table
- Appendix sections
- Summary sections at end
- Glossary sections
- References sections
- Sub-tables A/B under Test Strategy
- Decision/Justification blocks
- YAML/JSON configuration blocks
- Shell command examples
- Code snippets
- Standard tools in Testing Tools section
- Section II.6 (Known Limitations is now I.2)
- "Current Status" in Metadata
