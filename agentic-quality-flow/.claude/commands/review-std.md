---
name: review-std
description: Review a generated STD (YAML + test stubs) for traceability, pattern correctness, and code readiness
argument-hint: <JIRA-ID>
allowed-tools: Read, Write, Edit, Glob, Grep, Skill
---

# Review STD for $ARGUMENTS

You are the STD Reviewer entry point. Perform a comprehensive QE review of the generated
STD (YAML + test stubs) by invoking the **std-reviewer** skill.

## Input

The user has provided: `$ARGUMENTS`

This should be a Jira ticket ID (e.g., `CNV-12345`, `VIRTSTRAT-494`) for which an STD
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

**Check std_review toggle:**
If `project_context.feature_toggles.std_review` is false:
- Output: "STD review is disabled for project {project_context.display_name} (std_review toggle is false)."
- Exit. Do not proceed.

### Step 1: Parse the Jira ID

Extract the Jira ID from `project_context.jira_id` (e.g., CNV-72329).

### Step 2: Verify STD Files Exist

Check that the STD YAML file exists:
```
outputs/std/{JIRA_ID}/{JIRA_ID}_test_description.yaml
```

**If STD YAML does NOT exist:**
- Inform the user: "STD file not found. Please run `/std-builder {JIRA_ID}` first."
- Exit — do not proceed with review.

**If STD YAML exists:**
- Read the full STD YAML content.

Also check for stub files (these may or may not exist depending on feature toggles):

**Go stubs:**
```
outputs/std/{JIRA_ID}/go-tests/
```
- Use Glob to find `*_stubs_test.go` files
- If found, read each stub file

**Python stubs:**
```
outputs/std/{JIRA_ID}/python-tests/
```
- Use Glob to find `test_*_stubs.py` files
- If found, read each stub file

Record which artifacts are available:
```yaml
artifacts:
  std_yaml: true
  go_stubs: true/false
  go_stub_files: [list of file paths]
  python_stubs: true/false
  python_stub_files: [list of file paths]
```

### Step 3: Read Source STP

The STD was generated from an STP. Read the source STP for traceability checking:
```
outputs/stp/{JIRA_ID}/{JIRA_ID}_test_plan.md
```

**If STP file does NOT exist:**
- Log a warning: "STP file not found. Traceability review (Dimension 1) will be skipped."
- Set `stp_available: false`
- Continue with remaining dimensions

**If STP file exists:**
- Read the full STP content.
- Parse Section III (Requirements-to-Tests Mapping table) to extract the scenario list.
- Set `stp_available: true`

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
std-reviewer skill expects. The `_extraction_metadata` section indicates where each
piece of data came from (config files, repo scans, static override, or defaults).

**If extraction fails entirely:** Log a warning and continue without project-specific
rules. The std-reviewer skill's general rules still apply.

### Step 4: Load Pattern Library

Read the pattern library for pattern validation (Dimension 3):
```
{project_context.config_dir}/patterns/tier1_patterns.yaml
```

**If pattern library not found:**
- Log a warning: "Pattern library not found. Pattern library validation will use rule-based checks only."
- Set `pattern_library_available: false`
- Continue

### Step 5: Invoke std-reviewer Skill

Use the Skill tool to invoke the std-reviewer skill:

**Tool:** Skill
**Parameters:**
- skill: "std-reviewer"
- args: "{JIRA_ID}"

The skill receives context from the conversation including:
- The STD YAML content (read in Step 2)
- The stub file contents (read in Step 2, if available)
- The STP content (read in Step 3, if available)
- The project review rules (read in Step 3.5, if available)
- The pattern library (read in Step 4, if available)
- The project_context (from Step 0)

Perform the review across all 7 dimensions:
1. **STP-STD Traceability** — Forward and reverse coverage mapping (if STP available)
2. **STD YAML Structure** — v2.1-enhanced specification compliance
3. **Pattern Matching Correctness** — Primary patterns, helpers, decorators
4. **Test Step Quality** — Setup/execution/cleanup flow and assertion quality
4.5. **STD Content Policy** — No banned references, no implementation details, environment separation
5. **PSE Docstring Quality** — Preconditions/Steps/Expected in stub files (if stubs available)
6. **Code Generation Readiness** — Variables, imports, code structure validity

### Step 6: Generate Review Report

Generate the structured review report per the std-reviewer skill's output format.

Determine the verdict:
- **APPROVED:** 0 critical findings, 0 major findings
- **APPROVED_WITH_FINDINGS:** 0 critical findings, 1+ major or minor findings
- **NEEDS_REVISION:** 1+ critical findings

Determine confidence level:
- **HIGH:** STD YAML valid, STP available, stub files present, pattern library available
- **MEDIUM:** STD YAML valid and STP available, but stubs or pattern library missing
- **LOW:** STD YAML valid but STP unavailable

### Step 7: Save Review Report

Create the output directory and save the review report:

```
outputs/reviews/{JIRA_ID}/{JIRA_ID}_std_review.md
```

Use the Write tool to save the report.

### Step 8: Report to User

Once complete, show the user:

```
STD Review Complete: {JIRA_ID}

Verdict: {APPROVED | APPROVED_WITH_FINDINGS | NEEDS_REVISION}
Confidence: {HIGH | MEDIUM | LOW}

Findings: {X} critical, {Y} major, {Z} minor

Reviewed:
  STD YAML:      outputs/std/{JIRA_ID}/{JIRA_ID}_test_description.yaml
  Go Stubs:      {count} files (or "N/A")
  Python Stubs:  {count} files (or "N/A")

Report: outputs/reviews/{JIRA_ID}/{JIRA_ID}_std_review.md

Traceability:
  STP scenarios:       {count}
  STD scenarios:       {count}
  Forward coverage:    {X}/{Y} ({percent}%)
  Orphan STD scenarios: {count}

{If NEEDS_REVISION:}
Critical findings require attention before proceeding to test generation.
Review the report for details and recommendations.

{If APPROVED_WITH_FINDINGS:}
STD is usable but has findings worth addressing.
Review the report for improvement recommendations.

{If APPROVED:}
STD passes all review checks. Ready for test generation.
Next steps:
  /generate-go-tests {JIRA_ID}      (Tier 1 implementation)
  /generate-python-tests {JIRA_ID}  (Tier 2 implementation)
```

---

## Error Handling

**If STD YAML not found:**
- Error message: "STD file not found at outputs/std/{JIRA_ID}/{JIRA_ID}_test_description.yaml"
- Suggestion: "Please run `/std-builder {JIRA_ID}` first to create the STD"
- Exit without proceeding

**If STD YAML is invalid YAML:**
- Report as CRITICAL finding
- Attempt partial review of structure
- Verdict: NEEDS_REVISION

**If STP file not found:**
- Warning: "STP file not found. Skipping traceability review (Dimension 1)."
- Continue with remaining dimensions at reduced confidence

**If stub files not found:**
- Warning: "No stub files found. Skipping PSE docstring review (Dimension 5)."
- Continue with remaining dimensions

**If review skill fails:**
- Display error message
- Suggest reviewing the STD manually
- Do not save a partial report

---

## Prerequisites

**Before running this command:**
1. STD must exist (run `/std-builder {JIRA_ID}` first)
2. STP should exist for full traceability review (run `/stp-builder {JIRA_ID}` before STD)

---

## Example Usage

```
User: /review-std CNV-72329
Output: outputs/reviews/CNV-72329/CNV-72329_std_review.md

User: /review-std VIRTSTRAT-494
Output: outputs/reviews/VIRTSTRAT-494/VIRTSTRAT-494_std_review.md
```

---

## Workflow Overview

```
User: /review-std {JIRA_ID}
  |
  v
0. Resolve project: project-resolver -> project_context
  |
  v
1. Verify STD exists: outputs/std/{JIRA_ID}/{JIRA_ID}_test_description.yaml
  |
  v
2. Read stub files (Go and/or Python, if available)
  |
  v
3. Read source STP (for traceability, if available)
  |
  v
3.5. Resolve review rules: review-rules-extractor -> review_rules
  |
  v
4. Load pattern library (for pattern validation, if available)
  |
  v
5. Invoke std-reviewer skill (7 dimensions)
  |
  v
6. Generate and save review report:
   -> outputs/reviews/{JIRA_ID}/{JIRA_ID}_std_review.md
  |
  v
7. Report verdict to user
```

---

**End of Review STD Command**
