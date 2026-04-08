---
name: refine-std
description: Iteratively refine an STD (YAML + test stubs) by running review, fixing findings, and re-reviewing until approved
argument-hint: <JIRA-ID>
allowed-tools: Read, Write, Edit, Glob, Grep, Skill
---

# Refine STD for $ARGUMENTS

You are the STD Refinement orchestrator. Automate the review-fix cycle for an existing
STD (YAML + test stubs) by iteratively reviewing, fixing findings, and re-reviewing
until the verdict reaches APPROVED or APPROVED_WITH_FINDINGS (0 critical findings).

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

- Output: "STD review is disabled for project {project_context.display_name} (std_review toggle is false). Cannot refine without review capability."
- Exit. Do not proceed.

### Step 1: Verify STD Artifacts Exist

Extract the Jira ID from `project_context.jira_id` (e.g., CNV-72329).

Check that the STD YAML file exists:

```text
outputs/std/{JIRA_ID}/{JIRA_ID}_test_description.yaml
```

**If STD YAML does NOT exist:**

- Inform the user: "STD file not found. Please run `/std-builder {JIRA_ID}` first."
- Exit — do not proceed.

**If STD YAML exists:**

- Read the full STD YAML content.
- Record the artifact path for editing.

Also locate stub files:

**Go stubs:**

- Use Glob to find `outputs/std/{JIRA_ID}/go-tests/*_stubs_test.go`
- Record paths of found files

**Python stubs:**

- Use Glob to find `outputs/std/{JIRA_ID}/python-tests/test_*_stubs.py`
- Record paths of found files

### Step 2: Check for Existing Review

Check if a review report already exists:

```text
outputs/reviews/{JIRA_ID}/{JIRA_ID}_std_review.md
```

**If review exists:**

- Read the review report.
- Parse findings by severity (critical, major, minor) and dimension.
- Extract the current verdict.
- If verdict is already APPROVED: inform user "STD already approved. No refinement needed." and exit.

**If review does NOT exist:**

- Run `/review-std {JIRA_ID}` by invoking the review-std command workflow:
  1. Read STD YAML and stub files
  2. Read source STP (for traceability)
  3. Resolve review rules via review-rules-extractor skill
  4. Load pattern library
  5. Invoke std-reviewer skill
  6. Save review report
- Parse the resulting review report.
- If verdict is APPROVED: inform user and exit.

### Step 3: Initialize Refinement Loop

Set up tracking variables:

```yaml
iteration: 0
max_iterations: 5
consecutive_no_improvement: 0
max_no_improvement: 2
findings_history: []
changes_log: []
```

Parse the review report to build a prioritized fix queue:

1. Group findings by dimension
2. Sort groups: CRITICAL findings first, then MAJOR
3. Each group becomes one iteration target

### Step 4: Iterative Fix Loop

For each iteration (up to `max_iterations`):

#### 4.1: Select Next Dimension to Fix

Pick the highest-priority unfixed dimension from the fix queue:

- CRITICAL findings take absolute priority
- Within same severity, process in dimension order (Dim 1 before Dim 2)
- Skip dimensions marked as PASS in the review

#### 4.2: Apply Targeted Edits

Read the current STD YAML (and stub files if relevant) and apply fixes for the
selected dimension only.

**Fix strategies by dimension:**

**Dimension 1 — STP-STD Traceability:**

- Read the source STP from `outputs/stp/{JIRA_ID}/{JIRA_ID}_test_plan.md`
- For missing forward coverage (STP scenario not in STD): add the missing scenario to the STD YAML with appropriate test structure, pattern metadata, and variables
- For orphan STD scenarios (in STD but not in STP): verify if they are valid additions; if the review says to remove them, remove them; otherwise add a comment noting they extend STP coverage
- Ensure `stp_requirement_id` fields in STD YAML match STP requirement IDs

**Dimension 2 — STD YAML Structure:**

- Fix v2.1-enhanced specification compliance issues
- Correct missing required fields (metadata, scenarios, test_cases)
- Fix incorrect field types or values
- Ensure proper nesting and YAML syntax

**Dimension 3 — Pattern Matching Correctness:**

- Read the pattern library from `{project_context.config_dir}/patterns/tier1_patterns.yaml`
- Fix incorrect primary pattern assignments
- Add missing helper patterns or decorators
- Correct pattern parameter references

**Dimension 4 — Test Step Quality:**

- Fix setup/execution/cleanup flow issues
- Improve assertion specificity
- Ensure each test case has clear setup, action, and verification steps
- Remove implementation details from test descriptions

**Dimension 4.5 — STD Content Policy:**

- Remove banned references (internal URLs, proprietary names)
- Remove implementation details that leaked into test descriptions
- Separate environment configuration from test logic

**Dimension 5 — PSE Docstring Quality (stub files):**

- Fix Preconditions/Steps/Expected format in stub files
- Ensure each stub has complete PSE documentation
- Fix vague or missing preconditions
- Make steps specific and numbered
- Make expected results observable and verifiable
- Edit both Go stubs (`*_stubs_test.go`) and Python stubs (`test_*_stubs.py`) as needed

**Dimension 6 — Code Generation Readiness:**

- Fix undefined or incorrectly typed variables
- Add missing import references
- Correct code structure issues that would prevent generation
- Ensure variable names match pattern library conventions

**General rules for all edits:**

- Do NOT delete content unless the finding explicitly says to remove it
- Do NOT modify sections the reviewer marked as PASS
- Default to rewriting, not removing
- For YAML edits, ensure the result is valid YAML after each edit
- Use the Edit tool for targeted modifications

#### 4.3: Structural Validation Guard

After applying edits, validate structural integrity:

**For STD YAML:** Verify the file is valid YAML by reading it back and confirming
it parses without error. Check that required top-level keys are present (`metadata`,
`scenarios` or `test_cases`).

**For stub files (if edited):** Verify syntax is reasonable:

- Go stubs: Check for balanced braces and valid function signatures
- Python stubs: Check for valid indentation and function definitions

**If validation fails:**

- Revert the edits that broke structure
- Log the failure
- Move to the next dimension in the queue

#### 4.4: Re-run Review

Run the full review again by executing the review-std workflow:

1. Read updated STD YAML and stub files
2. Read source STP (for traceability)
3. Resolve review rules
4. Invoke std-reviewer skill
5. Save updated review report to `outputs/reviews/{JIRA_ID}/{JIRA_ID}_std_review.md`

Parse the new review report. Extract updated finding counts.

#### 4.5: Measure Improvement

Compare findings before and after this iteration:

```yaml
before:
  critical: {count}
  major: {count}
  minor: {count}
after:
  critical: {count}
  major: {count}
  minor: {count}
delta:
  critical: {after - before}
  major: {after - before}
  minor: {after - before}
```

**If improvement (critical + major decreased):**

- Reset `consecutive_no_improvement` to 0
- Log the iteration as successful
- Record changes in `changes_log`

**If no improvement (critical + major same or increased):**

- Increment `consecutive_no_improvement`
- Log the iteration as no-improvement
- If `consecutive_no_improvement >= 2`: stop the loop and report to user

#### 4.6: Check Stopping Criteria

Stop the loop if ANY of these conditions are met:

1. **Verdict is APPROVED:** 0 critical, 0 major findings — fully approved
2. **Verdict is APPROVED_WITH_FINDINGS:** 0 critical findings — acceptable
3. **Max iterations reached:** `iteration >= max_iterations`
4. **Consecutive no-improvement:** `consecutive_no_improvement >= max_no_improvement`

If none met, increment `iteration` and return to Step 4.1.

### Step 5: Save Refinement Log

Generate and save the refinement log:

```text
outputs/reviews/{JIRA_ID}/{JIRA_ID}_std_refinement_log.md
```

Use the following format:

```markdown
# Refinement Log: {JIRA_ID}

**Artifact:** outputs/std/{JIRA_ID}/{JIRA_ID}_test_description.yaml
**Stub Files:** {list of stub file paths, or "None"}
**Date:** {YYYY-MM-DD}
**Iterations:** {count}

## Iteration Summary

| # | Dimension Addressed | Findings Before | Findings After | Delta |
|:--|:--------------------|:----------------|:---------------|:------|
| 1 | {dimension} | {X}C, {Y}M, {Z}m | {X}C, {Y}M, {Z}m | {delta} |
| ... | ... | ... | ... | ... |

## Final Verdict: {APPROVED | APPROVED_WITH_FINDINGS | NEEDS_REVISION}

## Changes Applied

### Iteration 1: {Dimension} — {Description}

- {specific edit 1}
- {specific edit 2}
- ...

### Iteration 2: {Dimension} — {Description}

- {specific edit 1}
- ...
```

### Step 6: Report to User

Once complete, show the user:

```text
STD Refinement Complete: {JIRA_ID}

Initial Verdict:  {original verdict}
Final Verdict:    {final verdict}
Iterations:       {count}

Finding Progression:
  Start:  {X} critical, {Y} major, {Z} minor
  End:    {X} critical, {Y} major, {Z} minor

Artifacts:
  STD YAML:      outputs/std/{JIRA_ID}/{JIRA_ID}_test_description.yaml
  Go Stubs:      {count} files (or "N/A")
  Python Stubs:  {count} files (or "N/A")

Review:    outputs/reviews/{JIRA_ID}/{JIRA_ID}_std_review.md
Log:       outputs/reviews/{JIRA_ID}/{JIRA_ID}_std_refinement_log.md

{If final verdict is APPROVED:}
STD is fully approved. Ready for test generation.
Next steps:
  /generate-go-tests {JIRA_ID}      (Tier 1 implementation)
  /generate-python-tests {JIRA_ID}  (Tier 2 implementation)

{If final verdict is APPROVED_WITH_FINDINGS:}
STD approved with minor findings. Ready for test generation.
Review the refinement log for remaining minor items.
Next steps:
  /generate-go-tests {JIRA_ID}      (Tier 1 implementation)
  /generate-python-tests {JIRA_ID}  (Tier 2 implementation)

{If final verdict is NEEDS_REVISION:}
Refinement reached maximum iterations or stalled.
Remaining critical/major findings require manual attention.
Review the refinement log for details.

{If stopped due to consecutive no-improvement:}
Refinement stalled after 2 consecutive iterations with no improvement.
Remaining findings may require manual review or regeneration.
```

---

## Error Handling

**If STD YAML not found:**

- Error message: "STD file not found at outputs/std/{JIRA_ID}/{JIRA_ID}_test_description.yaml"
- Suggestion: "Please run `/std-builder {JIRA_ID}` first to create the STD"
- Exit without proceeding

**If STD YAML is invalid YAML after edit:**

- Revert the edit
- Log the parse error
- Move to the next dimension

**If review command fails during loop:**

- Log the failure for that iteration
- Attempt to continue with the next dimension
- If review fails twice consecutively, stop and report

**If structural validation fails after edit:**

- Log which edits broke validation
- Skip that dimension and move to the next
- Do not count as a no-improvement iteration

---

## Prerequisites

**Before running this command:**

1. STD must exist (run `/std-builder {JIRA_ID}` first)
2. STP should exist for traceability review (run `/stp-builder {JIRA_ID}` before STD)
3. A previous review report is optional (will be generated if missing)

---

## Example Usage

```text
User: /refine-std CNV-72329
Output:
  - Updated STD YAML: outputs/std/CNV-72329/CNV-72329_test_description.yaml
  - Updated stubs: outputs/std/CNV-72329/go-tests/, outputs/std/CNV-72329/python-tests/
  - Updated review: outputs/reviews/CNV-72329/CNV-72329_std_review.md
  - Refinement log: outputs/reviews/CNV-72329/CNV-72329_std_refinement_log.md

User: /refine-std VIRTSTRAT-494
Output:
  - Updated STD YAML: outputs/std/VIRTSTRAT-494/VIRTSTRAT-494_test_description.yaml
  - Updated review: outputs/reviews/VIRTSTRAT-494/VIRTSTRAT-494_std_review.md
  - Refinement log: outputs/reviews/VIRTSTRAT-494/VIRTSTRAT-494_std_refinement_log.md
```

---

## Workflow Overview

```text
User: /refine-std {JIRA_ID}
  |
  v
0. Resolve project: project-resolver -> project_context
  |
  v
1. Verify STD exists: outputs/std/{JIRA_ID}/{JIRA_ID}_test_description.yaml
  |
  v
2. Run or read existing review -> parse findings
  |
  v
3. Build prioritized fix queue (CRITICAL first, then MAJOR)
  |
  v
4. Iterative fix loop (max 5 iterations):
   |
   +-> 4.1 Select next dimension
   +-> 4.2 Apply targeted edits to STD YAML / stubs
   +-> 4.3 Validate structure (YAML parse, stub syntax)
   +-> 4.4 Re-run review (std-reviewer)
   +-> 4.5 Measure improvement (delta)
   +-> 4.6 Check stopping criteria
   |        |
   |        +-> APPROVED or APPROVED_WITH_FINDINGS -> stop
   |        +-> Max iterations reached -> stop
   |        +-> 2 consecutive no-improvement -> stop
   |        +-> Otherwise -> next iteration
   |
   v
5. Save refinement log:
   -> outputs/reviews/{JIRA_ID}/{JIRA_ID}_std_refinement_log.md
  |
  v
6. Report results to user
```

---

**End of Refine STD Command**
