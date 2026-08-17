---
name: review-stp
description: Review a generated STP document for QE quality, rule compliance, and requirement coverage
argument-hint: <JIRA-ID>
allowed-tools: Read, Write, Edit, Glob, Grep, Skill, mcp__mcp-atlassian__jira_get_issue, mcp__mcp-atlassian__jira_search
---

# Review STP for $ARGUMENTS

You are the STP Reviewer entry point. Perform a comprehensive QE review of the generated
STP document by invoking the **stp-reviewer** skill.

## Input

The user has provided: `$ARGUMENTS`

This should be a Jira ticket ID (e.g., `CNV-12345`, `VIRTSTRAT-494`) for which an STP
has already been generated.

## Workflow

### Step 0: Resolve Project

Use the Skill tool to invoke the project-resolver skill:

**Tool:** Skill
**Parameters:**
- skill: "project-resolver"
- args: "$ARGUMENTS"

This returns `project_context` containing:
- `project_id`, `display_name`, `jira_id`
- `config_dir` (path to project config files)
- `feature_toggles` (what capabilities are enabled)
- `stp_header`, `versioning`

**If project resolution fails:** Display the error and exit. Do not proceed.

**Check stp_review toggle:**
If `project_context.feature_toggles.stp_review` is false:
- Output: "STP review is disabled for project {project_context.display_name} (stp_review toggle is false)."
- Exit. Do not proceed.

### Step 1: Parse the Jira ID

Extract the Jira ID from `project_context.jira_id` (e.g., CNV-72329).

### Step 2: Verify STP File Exists

Check that the STP file exists:
```
outputs/stp/{JIRA_ID}/{JIRA_ID}_test_plan.md
```

**If STP file does NOT exist:**
- Inform the user: "STP file not found. Please run `/stp-builder {JIRA_ID}` first."
- Exit — do not proceed with review.

**If STP file exists:**
- Read the full STP file content.
- Proceed to Step 3.

### Step 3: Fetch Jira Source Data

Fetch the Jira ticket data for comparison against the STP. This is essential for
Dimension 2 (Requirement Coverage) and Dimension 4 (Risk & Limitation Accuracy).

Use the Jira MCP tools to fetch:

1. **Main issue:** Use `jira_get_issue` with the Jira ID
   - Extract: summary, description, acceptance_criteria, status, fix_version, labels, components
2. **Linked issues:** From the main issue's links, fetch each linked issue
   - Extract: key, summary, relationship type, description

**If Jira fetch fails:**
- Log a warning: "Could not fetch Jira data. Review will be limited to content-only analysis."
- Set `jira_data_available: false`
- Continue with Steps 4-5 (content-only review)

**If Jira fetch succeeds:**
- Set `jira_data_available: true`
- Store the fetched data for passing to the reviewer skill

### Step 3.1: Extract PR Data for Fix-Scope Analysis

If the Jira issue type (from Step 3) is Bug, Customer Case, or Defect:

1. Check Jira issue data for PR URLs in custom fields or comments
2. If PR URLs are found, use GitHub MCP tools to fetch PR details and extract:
   - `key_changes`: list of modified functions/methods from the PR diff
   - `files_changed`: count of files changed in the PR
3. Construct `fix_scope` summary:
   ```yaml
   fix_scope:
     files_changed: <count>
     key_changes: <list from PR>
     issue_type: <bug|customer_case|defect>
   ```
4. Pass `fix_scope` to stp-reviewer skill context for Rule P evaluation

If no PR URLs are found or the issue type is Feature/Enhancement:
- Set `fix_scope: null` — Rule P will be skipped by the reviewer

### Step 3.5: Resolve Review Rules

Invoke the **review-rules-extractor** skill to produce project-specific review rules:

**Tool:** Skill
**Parameters:**

- skill: "review-rules-extractor"
- args: "{JIRA_ID}"

The skill reads project config files (`project.yaml`, `components.yaml`, `tier1.yaml`,
`tier2.yaml`, `patterns/tier1_patterns.yaml`) and optionally scans repositories to extract
review rules dynamically. If a static `review_rules.yaml` exists in the project's config
directory, its values take priority as overrides.

The output is a complete `review_rules` data structure identical to the format the
stp-reviewer skill expects. The `_extraction_metadata` section indicates where each
piece of data came from (config files, repo scans, static override, or defaults).

**If extraction fails entirely:** Log a warning and continue without project-specific
rules. The stp-reviewer skill's general rules still apply.

### Step 4: Read the STP Template

Read the project STP template for Rule B comparison:
```
{project_context.config_dir}/templates/stp/stp-template.md
```

**If template not found:**
- Log a warning: "STP template not found. Rule B check will be skipped."
- Set `template_available: false`
- Continue

### Step 5: Invoke stp-reviewer Skill

Use the Skill tool to invoke the stp-reviewer skill:

**Tool:** Skill
**Parameters:**
- skill: "stp-reviewer"
- args: "{JIRA_ID}"

The skill receives context from the conversation including:
- The STP file content (read in Step 2)
- The Jira source data (fetched in Step 3, if available)
- The project review rules (read in Step 3.5, if available)
- The STP template (read in Step 4, if available)
- The project_context (from Step 0)

Perform the review across all 7 dimensions:
1. **Rule Compliance (A-P)** — Check each rule against STP content
2. **Requirement Coverage** — Compare STP scenarios against Jira acceptance criteria (if Jira data available)
3. **Scenario Quality** — Evaluate specificity, user perspective, distribution
4. **Risk & Limitation Accuracy** — Compare against Jira data (if available)
5. **Scope Boundary Assessment** — Compare scope against Jira feature description
6. **Test Strategy Appropriateness** — Validate Y/N/N/A classifications
7. **Metadata Accuracy** — Compare metadata against Jira data

### Step 6: Generate Review Report

Generate the structured review report per the stp-reviewer skill's output format.

Determine the verdict:
- **APPROVED:** 0 critical findings, 0 major findings
- **APPROVED_WITH_FINDINGS:** 0 critical findings, 1+ major or minor findings
- **NEEDS_REVISION:** 1+ critical findings

Determine confidence level:
- **HIGH:** Jira data available, all sections present, template comparison done
- **MEDIUM:** Jira data incomplete or template unavailable
- **LOW:** Jira data unavailable (content-only review)

### Step 7: Save Review Report

Create the output directory and save the review report:

```
outputs/reviews/{JIRA_ID}/{JIRA_ID}_stp_review.md
```

Use the Write tool to save the report.

### Step 8: Report to User

Once complete, show the user:

```
STP Review Complete: {JIRA_ID}

Verdict: {APPROVED | APPROVED_WITH_FINDINGS | NEEDS_REVISION}
Confidence: {HIGH | MEDIUM | LOW}

Findings: {X} critical, {Y} major, {Z} minor

Reviewed: outputs/stp/{JIRA_ID}/{JIRA_ID}_test_plan.md
Report:   outputs/reviews/{JIRA_ID}/{JIRA_ID}_stp_review.md

{If NEEDS_REVISION:}
Critical findings require attention before proceeding to STD generation.
Review the report for details and recommendations.

{If APPROVED_WITH_FINDINGS:}
STP is usable but has findings worth addressing.
Review the report for improvement recommendations.

{If APPROVED:}
STP passes all review checks. Ready for STD generation.
Next step: /std-builder {JIRA_ID}
```

---

## Error Handling

**If STP file not found:**
- Error message: "STP file not found at outputs/stp/{JIRA_ID}/{JIRA_ID}_test_plan.md"
- Suggestion: "Please run `/stp-builder {JIRA_ID}` first to create the STP"
- Exit without proceeding

**If Jira data fetch fails:**
- Warning: "Jira data unavailable. Performing content-only review (reduced coverage)."
- Skip Dimensions 2 and 4 (require source comparison)
- Continue with remaining dimensions

**If template not found:**
- Warning: "STP template not found at {path}. Skipping Rule B check."
- Skip Rule B in Dimension 1
- Continue with remaining checks

**If review skill fails:**
- Display error message
- Suggest reviewing the STP manually
- Do not save a partial report

---

## Prerequisites

**Before running this command:**
1. STP file must exist (run `/stp-builder {JIRA_ID}` first)
2. Jira access is recommended but not required (review works without it at reduced confidence)

---

## Example Usage

```
User: /review-stp CNV-72329
Output: outputs/reviews/CNV-72329/CNV-72329_stp_review.md

User: /review-stp VIRTSTRAT-494
Output: outputs/reviews/VIRTSTRAT-494/VIRTSTRAT-494_stp_review.md
```

---

## Workflow Overview

```
User: /review-stp {JIRA_ID}
  |
  v
0. Resolve project: project-resolver -> project_context
  |
  v
1. Verify STP exists: outputs/stp/{JIRA_ID}/{JIRA_ID}_test_plan.md
  |
  v
2. Fetch Jira source data (for coverage comparison)
  |
  v
3. Resolve review rules: review-rules-extractor -> review_rules
  |
  v
4. Read STP template (for Rule B comparison)
  |
  v
5. Invoke stp-reviewer skill (7 dimensions, rules A-P)
  |
  v
6. Generate and save review report:
   -> outputs/reviews/{JIRA_ID}/{JIRA_ID}_stp_review.md
  |
  v
7. Report verdict to user
```

---

**End of Review STP Command**
