---
name: refine-stp
description: Iteratively refine an STP document by running review, fixing findings, and re-reviewing until approved
argument-hint: <JIRA-ID>
allowed-tools: Read, Write, Edit, Glob, Grep, Skill, mcp__mcp-atlassian__jira_get_issue, mcp__mcp-atlassian__jira_search
---

# Refine STP for $ARGUMENTS

You are the STP Refinement orchestrator. Automate the review-fix cycle for an existing
STP document by iteratively reviewing, fixing findings, and re-reviewing until the
verdict reaches APPROVED or APPROVED_WITH_FINDINGS (0 critical findings).

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

- Output: "STP review is disabled for project {project_context.display_name} (stp_review toggle is false). Cannot refine without review capability."
- Exit. Do not proceed.

### Step 1: Verify STP Artifact Exists

Extract the Jira ID from `project_context.jira_id` (e.g., CNV-72329).

Check that the STP file exists:

```text
outputs/stp/{JIRA_ID}/{JIRA_ID}_test_plan.md
```

**If STP file does NOT exist:**

- Inform the user: "STP file not found. Please run `/stp-builder {JIRA_ID}` first."
- Exit — do not proceed.

**If STP file exists:**

- Read the full STP file content.
- Record the artifact path for editing.

### Step 2: Check for Existing Review

Check if a review report already exists:

```text
outputs/reviews/{JIRA_ID}/{JIRA_ID}_stp_review.md
```

**If review exists:**

- Read the review report.
- Parse findings by severity (critical, major, minor) and dimension/rule.
- Extract the current verdict.
- If verdict is already APPROVED: inform user "STP already approved. No refinement needed." and exit.

**If review does NOT exist:**

- Run `/review-stp {JIRA_ID}` by invoking the review-stp command workflow:
  1. Fetch Jira source data
  2. Resolve review rules via review-rules-extractor skill
  3. Read STP template
  4. Invoke stp-reviewer skill
  5. Save review report
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

1. Group findings by dimension/rule
2. Sort groups: CRITICAL findings first, then MAJOR
3. Each group becomes one iteration target

### Step 4: Iterative Fix Loop

For each iteration (up to `max_iterations`):

#### 4.1: Select Next Dimension to Fix

Pick the highest-priority unfixed dimension/rule group from the fix queue:

- CRITICAL findings take absolute priority
- Within same severity, process in dimension order (Rule A before Rule B, Dim 1 before Dim 2)
- Skip dimensions marked as PASS in the review

#### 4.2: Apply Targeted Edits

Read the current STP content and apply fixes for the selected dimension only.

**Fix strategies by common finding types:**

**Rule A violations (Abstraction Level):**

- Read `internal_to_user_mappings` from project review rules (if available)
- Rewrite internal-mechanism language to user-observable language
- Example: "HandleCPUHotplug is invoked" -> "CPU can be hot-added to a running VM"
- Apply to all occurrences in the STP (scope, scenarios, strategy sections)

**Dimension 2 gaps (Requirement Coverage):**

- Identify uncovered acceptance criteria from the review findings
- For each uncovered criterion, generate test scenarios using the scenario-builder skill format:
  - Start with action verb (Verify, Test, Validate, Confirm)
  - 5-10 word description
  - Include positive and negative variants
- Add new scenarios to Section III under the appropriate requirement
- Assign tier and priority per existing conventions

**Rule B violations (Template Compliance):**

- Read the STP template from `{project_context.config_dir}/templates/stp/stp-template.md`
- Compare STP structure against template
- Add missing sections, fix section ordering, correct heading levels

**Rule J violations (One Tier per Row):**

- Find rows with multiple tiers listed
- Split each into separate entries, one per tier
- Preserve all other content in each split entry

**Rule K contradictions (Scope vs Section III):**

- Identify contradictions between scope/strategy sections and Section III
- Resolve in favor of Section III (requirements mapping is source of truth)
- Update scope/strategy sections to be consistent with Section III

**Dimension 3 (Scenario Quality):**

- Rewrite vague scenarios to be specific and action-oriented
- Ensure each scenario starts with an action verb
- Keep descriptions to 5-10 words
- Remove test steps or implementation details from scenario descriptions

**Dimension 4 (Risk & Limitation Accuracy):**

- Cross-reference risks/limitations against Jira data
- Add missing risks identified in the review
- Remove or correct inaccurate risk statements

**Dimension 5 (Scope Boundary Assessment):**

- Adjust scope boundaries to match the Jira feature description
- Add missing in-scope items or remove out-of-scope items per findings

**Dimension 6 (Test Strategy Appropriateness):**

- Fix incorrect Y/N/N/A classifications in the test strategy section
- Ensure classifications align with Section III scenarios

**Dimension 7 (Metadata Accuracy):**

- Correct metadata fields to match Jira source data
- Fix version numbers, component names, dates

**General rules for all edits:**

- Do NOT delete content unless the finding explicitly says to remove it
- Do NOT modify sections the reviewer marked as PASS
- Default to rewriting, not removing
- Use the Edit tool for targeted modifications

#### 4.3: Structural Validation Guard

After applying edits, validate structural integrity:

**Tool:** Skill
**Parameters:**

- skill: "output-validator"
- args: "{JIRA_ID}"

**If validation fails:**

- Revert the edits that broke structure
- Log the failure
- Move to the next dimension in the queue

#### 4.4: Re-run Review

Run the full review again by executing the review-stp workflow:

1. Fetch Jira source data (reuse from Step 2 if still in context)
2. Resolve review rules
3. Invoke stp-reviewer skill
4. Save updated review report to `outputs/reviews/{JIRA_ID}/{JIRA_ID}_stp_review.md`

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
outputs/reviews/{JIRA_ID}/{JIRA_ID}_stp_refinement_log.md
```

Use the following format:

```markdown
# Refinement Log: {JIRA_ID}

**Artifact:** outputs/stp/{JIRA_ID}/{JIRA_ID}_test_plan.md
**Date:** {YYYY-MM-DD}
**Iterations:** {count}

## Iteration Summary

| # | Dimension Addressed | Findings Before | Findings After | Delta |
|:--|:--------------------|:----------------|:---------------|:------|
| 1 | {dimension/rule} | {X}C, {Y}M, {Z}m | {X}C, {Y}M, {Z}m | {delta} |
| ... | ... | ... | ... | ... |

## Final Verdict: {APPROVED | APPROVED_WITH_FINDINGS | NEEDS_REVISION}

## Changes Applied

### Iteration 1: {Dimension/Rule} — {Description}

- {specific edit 1}
- {specific edit 2}
- ...

### Iteration 2: {Dimension/Rule} — {Description}

- {specific edit 1}
- ...
```

### Step 6: Report to User

Once complete, show the user:

```text
STP Refinement Complete: {JIRA_ID}

Initial Verdict:  {original verdict}
Final Verdict:    {final verdict}
Iterations:       {count}

Finding Progression:
  Start:  {X} critical, {Y} major, {Z} minor
  End:    {X} critical, {Y} major, {Z} minor

Artifact:  outputs/stp/{JIRA_ID}/{JIRA_ID}_test_plan.md
Review:    outputs/reviews/{JIRA_ID}/{JIRA_ID}_stp_review.md
Log:       outputs/reviews/{JIRA_ID}/{JIRA_ID}_stp_refinement_log.md

{If final verdict is APPROVED:}
STP is fully approved. Ready for STD generation.
Next step: /std-builder {JIRA_ID}

{If final verdict is APPROVED_WITH_FINDINGS:}
STP approved with minor findings. Ready for STD generation.
Review the refinement log for remaining minor items.
Next step: /std-builder {JIRA_ID}

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

**If STP file not found:**

- Error message: "STP file not found at outputs/stp/{JIRA_ID}/{JIRA_ID}_test_plan.md"
- Suggestion: "Please run `/stp-builder {JIRA_ID}` first to create the STP"
- Exit without proceeding

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

1. STP file must exist (run `/stp-builder {JIRA_ID}` first)
2. Jira access is recommended for full review coverage
3. A previous review report is optional (will be generated if missing)

---

## Example Usage

```text
User: /refine-stp CNV-72329
Output:
  - Updated STP: outputs/stp/CNV-72329/CNV-72329_test_plan.md
  - Updated review: outputs/reviews/CNV-72329/CNV-72329_stp_review.md
  - Refinement log: outputs/reviews/CNV-72329/CNV-72329_stp_refinement_log.md

User: /refine-stp VIRTSTRAT-494
Output:
  - Updated STP: outputs/stp/VIRTSTRAT-494/VIRTSTRAT-494_test_plan.md
  - Updated review: outputs/reviews/VIRTSTRAT-494/VIRTSTRAT-494_stp_review.md
  - Refinement log: outputs/reviews/VIRTSTRAT-494/VIRTSTRAT-494_stp_refinement_log.md
```

---

## Workflow Overview

```text
User: /refine-stp {JIRA_ID}
  |
  v
0. Resolve project: project-resolver -> project_context
  |
  v
1. Verify STP exists: outputs/stp/{JIRA_ID}/{JIRA_ID}_test_plan.md
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
   +-> 4.1 Select next dimension/rule group
   +-> 4.2 Apply targeted edits to STP
   +-> 4.3 Validate structure (output-validator)
   +-> 4.4 Re-run review (stp-reviewer)
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
   -> outputs/reviews/{JIRA_ID}/{JIRA_ID}_stp_refinement_log.md
  |
  v
6. Report results to user
```

---

**End of Refine STP Command**
