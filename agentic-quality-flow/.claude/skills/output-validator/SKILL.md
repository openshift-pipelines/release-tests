---
name: output-validator
description: Validate STP document structure and content completeness
model: claude-opus-4-6
---

# Output Validator Skill

**Phase:** Post-Processing
**User-Invocable:** false

## Purpose

Validate STP document structure and content completeness.

## When to Use

Invoked by the **document-formatter** subagent to verify the STP before saving.

## Input

```yaml
document: |
  # Openshift-virtualization-tests Test plan

  ## **CPU Hot-Plug Feature - Quality Engineering Plan**

  ### **Metadata & Tracking**
  ...
```

## Output Format

```yaml
validation_results:
  valid: true

  structure:
    document_header: pass
    feature_title: pass
    all_sections_present: pass
    horizontal_rules: pass

  list_items:
    metadata: pass  # 6 bullet items
    section_i_1: pass  # 5 checkbox items
    section_i_3: pass  # 5 checkbox items
    section_ii_1_out_of_scope: pass  # At least 1 checkbox item
    section_ii_2: pass  # 13 checkbox items across 4 categories
    section_ii_3: pass  # 10 bullet items
    section_ii_5: pass  # 7 checkbox categories with sub-items
    section_iii: pass  # No minimum - comprehensive coverage

  content:
    no_code_blocks: pass
    valid_test_tiers: pass
    unique_requirement_summaries: pass
    no_generic_scenarios: pass

  errors: []
  warnings:
    - "Section II.5 'Other' risk category is empty - consider adding risks if applicable"

total_checks: 18
passed: 17
failed: 0
warnings: 1
```

## Validation Rules

### Document Structure

#### 1. Document Header

Read `project_context.stp_header` for the expected document header.

```
Example: # Openshift-virtualization-tests Test plan
```
Check: First line matches the configured `stp_header` value.

#### 2. Feature Title

Read `project_context.feature_title_format` for the expected format if configured.

```
Example: ## **[Any Title] - Quality Engineering Plan**
```
Check: Second non-empty line matches pattern

#### 3. Required Sections

Check presence of all sections in order:
- Metadata & Tracking
- Feature Overview
- I. Motivation and Requirements Review (QE Review Guidelines)
- Section I.1 - Requirement & User Story Review Checklist
- Section I.2 - Known Limitations
- Section I.3 - Technology and Design Review
- II. Software Test Plan (STP)
- Section II.1 - Scope of Testing (including Testing Goals and Out of Scope)
- Section II.2 - Test Strategy (including Functional, Non-Functional, Integration & Compatibility, Infrastructure)
- Section II.3 - Test Environment
- Section II.3.1 - Testing Tools & Frameworks
- Section II.4 - Entry Criteria
- Section II.5 - Risks
- III. Test Scenarios & Traceability
- Section III.1 - Requirements-to-Tests Mapping
- Section IV - Sign-off and Approval

#### 4. Horizontal Rules

```
Expected: --- after Feature Overview, after Section I.3 (before Section II), after Section III (before Section IV)
```

### List Item Counts

| Section | Format | Required Items | Check Method |
|:--------|:-------|:---------------|:-------------|
| Metadata | Bullet list (`- **Field:**`) | 6 | Count lines matching `- **` pattern |
| I.1 | Checkbox list (`- [ ]`) | 5 | Count lines matching `- [ ]` with bold label |
| I.3 | Checkbox list (`- [ ]`) | 5 | Count lines matching `- [ ]` with bold label |
| II.1 Out of Scope | Checkbox list (`- [ ]`) | At least 1 | Count `- [ ]` lines under Out of Scope heading |
| II.2 Test Strategy | Grouped checkbox list | 13 total | Count across 4 category groups (see below) |
| II.3 Test Environment | Bullet list (`- **Label:**`) | 10 | Count lines matching `- **` pattern |
| II.5 Risks | Checkbox list with sub-items | 7 categories | Count `- [ ]` top-level items |
| III.1 | Bullet list (`- **[Jira-ID]**`) | No minimum | Count requirement entries (no minimum enforced) |

#### Test Strategy Category Breakdown (Section II.2)

| Category | Expected Items |
|:---------|:---------------|
| Functional | 4 (Functional, Automation, Regression, Upgrade) |
| Non-Functional | 5 (Performance, Scale, Security, Usability, Monitoring) |
| Integration & Compatibility | 3 (Compatibility, Dependencies, Cross Integrations) |
| Infrastructure | 1 (Cloud Testing) |

Check: Each category heading is present as a bold sub-heading, followed by its checkbox items.

### Content Validation

#### 1. No Code Blocks
Check: Document does not contain:
- ` ``` ` code fence blocks
- YAML/JSON content blocks
- Shell command examples with flags
- Kubernetes manifest examples

#### 2. Valid Test Tiers
Check: All tier references in Section III.1 use inline format:
- `[Tier 1]`
- `[Tier 2]`

Invalid values:
- `Tier 1 (Functional)` (old column format)
- `Tier 2 (End-to-End)` (old column format)
- `Functional`, `API`, `Integration`, `Security`, `Upgrade` (not valid tier names)
- `Unit Tests` (not used in Section III mapping)

#### 3. Unique Requirement Summaries
Check: In Section III.1, each requirement entry is unique
- No repeated generic summaries
- Each `- **[Jira-ID]**` line describes a specific capability

#### 4. No Generic/Meta Scenarios
Check: Test scenarios do not include:
- "Verify automated tests pass in CI"
- "All tests should pass"
- "Ensure test coverage is complete"
- "Validate CI pipeline runs successfully"

#### 5. Section III.1 Format
Check: Each requirement entry follows the bullet-based format:
- `- **[Jira-123]** -- As a user...` (requirement line with Jira ID in bold)
- Indented `*Test Scenario:*` sub-item describing the test
- Indented `*Priority:*` sub-item with priority value

### Prohibited Content

Check document does NOT contain:
- Appendix sections
- Summary sections at end
- Glossary sections
- References sections
- Old-style section numbering (II.4.A, II.4.B, II.4.C, II.4.D, II.6, II.7, II.8)
- "Decision:" or "Justification:" blocks
- Real IP addresses (except RFC 5737)
- Real email addresses
- Third-party vendor names
- "Current Status" metadata field (removed from template)

## Error Severity

| Severity | Handling |
|:---------|:---------|
| **Error** | Must be fixed before save |
| **Warning** | Document can be saved, but should be addressed |
| **Info** | Optional improvement suggestion |

### Error Examples
- Missing required section
- Wrong tier reference format
- Code block present
- Section missing expected list items
- Old-style section numbering used (e.g., II.4.A, II.7, II.8)

### Warning Examples
- Optional section empty
- No negative test scenarios
- Fewer than expected test scenarios
- Known Limitations section (I.2) is empty

### Info Examples
- Consider adding more detail to section X
- Entry criteria could be more specific

## Auto-Fix Capabilities

Some issues can be auto-fixed:

| Issue | Auto-Fix |
|:------|:---------|
| Wrong tier format | Convert "Tier 1 (Functional)" to "[Tier 1]" |
| Missing horizontal rule | Insert `---` at expected positions |
| Trailing whitespace | Remove |
| Old metadata field name | Rename "Feature in Jira" to "Feature Tracking", "Jira Tracking" to "Epic Tracking" |

Issues that CANNOT be auto-fixed:
- Missing sections (need content)
- Invalid test scenarios (need rewriting)
- Code blocks (need conceptual replacement)
- Platform-level tests (need rejection)
