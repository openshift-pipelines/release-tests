---
name: std-builder
description: Generate STD (YAML + test stubs with PSE docstrings) from an existing STP file
argument-hint: <JIRA-ID>
allowed-tools: Read, Write, Edit, Task, Glob, Grep, Skill
---

# STD Builder

Generates the complete Software Test Description (STD):
1. **STD YAML file** (internal format for automation)
2. **Test stubs with PSE docstrings** (the deliverable for human review)

Per [SOFTWARE_TEST_DESCRIPTION.md](https://github.com/rnetser/openshift-virtualization-tests/blob/std-template/docs/SOFTWARE_TEST_DESCRIPTION.md), the STD = docstrings in test files.

---

When the user runs this command with a Jira ID, you MUST:

## Step 0: Resolve Project

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

**Check std_generation toggle:**
If `project_context.feature_toggles.std_generation` is false:
- Output: "STD generation is disabled for project {project_context.display_name} (std_generation toggle is false)."
- Exit. Do not proceed.

## Step 1: Parse the Jira ID

Extract the Jira ID from `project_context.jira_id` (e.g., CNV-66855, VIRTSTRAT-494).

## Step 2: Verify STP File Exists

**CRITICAL: STD generation requires an existing STP file.**

Check that the STP file exists:
```
outputs/stp/{JIRA_ID}/{JIRA_ID}_test_plan.md
```

**If STP file does NOT exist:**
- Inform the user: "STP file not found. Please run `/stp-builder {JIRA_ID}` first."
- Exit - do not proceed with STD generation

**If STP file exists:**
- Proceed to Step 3

## Step 3: Generate STD YAML (Internal Format)

Use the Skill tool to invoke the std-orchestrator skill:

**Tool:** Skill
**Parameters:**
- skill: "std-orchestrator"
- args: "{JIRA_ID}"

The std-orchestrator skill will:
1. Read the STP file at `outputs/stp/{JIRA_ID}/{JIRA_ID}_test_plan.md`
2. Parse Section III (Requirements-to-Tests Mapping table)
3. Extract all test scenarios
4. Generate comprehensive STD YAML file:
   - Output: `outputs/std/{JIRA_ID}/{JIRA_ID}_test_description.yaml`
5. Validate STD YAML

## Step 4: Generate Test Stubs (The Actual STD)

After STD YAML is generated, generate test stubs with PSE docstrings.

**Check tier distribution in STD YAML:**
- Count Tier 1 scenarios
- Count Tier 2 scenarios

**Check feature toggles from project_context before generating stubs:**

**If Tier 1 scenarios exist AND `project_context.feature_toggles.tier1_tests` is true:**

Use the Skill tool to invoke go-stub-generator:

**Tool:** Skill
**Parameters:**
- skill: "go-stub-generator"
- args: "{JIRA_ID}"

This generates Go/Ginkgo test stubs with PSE comments:
- Output: `outputs/std/{JIRA_ID}/go-tests/*_stubs_test.go`
- Contains: `PendingIt()` blocks with PSE comments
- Excluded from test execution

**If Tier 1 scenarios exist BUT `project_context.feature_toggles.tier1_tests` is false:**
- Skip Go stub generation
- Log: "Skipping Go stub generation: tier1_tests is disabled for project {project_context.display_name}."

**If Tier 2 scenarios exist AND `project_context.feature_toggles.tier2_tests` is true:**

Use the Skill tool to invoke python-stub-generator:

**Tool:** Skill
**Parameters:**
- skill: "python-stub-generator"
- args: "{JIRA_ID}"

This generates Python/pytest test stubs with PSE docstrings:
- Output: `outputs/std/{JIRA_ID}/python-tests/test_*_stubs.py`
- Contains: `pass` body with comprehensive docstrings
- Marked with `__test__ = False` to exclude from collection

**If Tier 2 scenarios exist BUT `project_context.feature_toggles.tier2_tests` is false:**
- Skip Python stub generation
- Log: "Skipping Python stub generation: tier2_tests is disabled for project {project_context.display_name}."

## Step 5: Report to User

Once complete, show the user:

```
✅ STD Generation Complete!

📄 Input: outputs/stp/{JIRA_ID}/{JIRA_ID}_test_plan.md

📊 Summary:
- STP scenarios: {TOTAL_COUNT} ({TIER1_COUNT} Tier 1, {TIER2_COUNT} Tier 2)
- STD YAML: {JIRA_ID}_test_description.yaml (internal format)

📁 STD Output (for review):
- outputs/std/{JIRA_ID}/go-tests/*_stubs_test.go ({TIER1_COUNT} test stubs)
- outputs/std/{JIRA_ID}/python-tests/test_*_stubs.py ({TIER2_COUNT} test stubs)

📋 Phase 1 Checklist:
- [ ] STP link in module docstring
- [ ] Tests grouped in class with shared preconditions
- [ ] Each test has: Preconditions, Steps, Expected
- [ ] Each test verifies ONE thing with ONE Expected
- [ ] Test bodies contain only 'pass' / Skip()

✅ Ready for design review!

📌 Next steps:
1. Review the test stubs (the STD)
2. Submit PR for design review
3. After approval, run:
   - /generate-go-tests {JIRA_ID}     (Tier 1 implementation)
   - /generate-python-tests {JIRA_ID} (Tier 2 implementation)
```

---

## Output Structure

```
outputs/std/{JIRA_ID}/
├── {JIRA_ID}_test_description.yaml     (STD YAML - internal format)
├── go-tests/                           (Tier 1 STD - test stubs)
│   └── {feature}_stubs_test.go         (PendingIt + PSE comments)
└── python-tests/                       (Tier 2 STD - test stubs)
    └── test_{feature}_stubs.py         (__test__=False + PSE docstrings)
```

---

## Workflow Overview

```
User: /std-builder {JIRA_ID}
  ↓
0. Resolve project: project-resolver → project_context
  ↓
1. Verify STP exists: outputs/stp/{JIRA_ID}/{JIRA_ID}_test_plan.md
  ↓
2. Generate STD YAML (internal):
   → outputs/std/{JIRA_ID}/{JIRA_ID}_test_description.yaml
  ↓
3. Generate test stubs (respecting feature toggles):
   → outputs/std/{JIRA_ID}/go-tests/*_stubs_test.go (if tier1_tests enabled)
   → outputs/std/{JIRA_ID}/python-tests/test_*_stubs.py (if tier2_tests enabled)
  ↓
4. Report results:
   STD complete - ready for design review
```

---

## Error Handling

**If STP file not found:**
- Error message: "STP file not found at outputs/stp/{JIRA_ID}/{JIRA_ID}_test_plan.md"
- Suggestion: "Please run `/stp-builder {JIRA_ID}` first to create the STP"
- Exit without proceeding

**If STP Section III is empty:**
- Error message: "No test scenarios found in STP Section III"
- Suggestion: "Verify STP file is complete and contains Requirements-to-Tests Mapping table"
- Exit

**If std-orchestrator skill fails:**
- Display error message from skill
- Show partial results if any
- Suggest reviewing errors and re-running

**If stub generation fails:**
- Show which stubs were generated successfully
- Report which scenarios failed
- STD YAML is still available for manual review

---

## Prerequisites

**Before running this command:**
1. ✅ STP file must exist (run `/stp-builder {JIRA_ID}` first)
2. ✅ STP must contain Section III with test scenarios

---

## Example Usage

**Step 1: Generate STP**
```
User: /stp-builder {JIRA_ID}
Output: outputs/stp/{JIRA_ID}/{JIRA_ID}_test_plan.md
```

**Step 2: Generate STD (YAML + Stubs)**
```
User: /std-builder {JIRA_ID}
Output:
   - outputs/std/{JIRA_ID}/{JIRA_ID}_test_description.yaml (internal)
   - outputs/std/{JIRA_ID}/go-tests/*_stubs_test.go (if tier1_tests enabled)
   - outputs/std/{JIRA_ID}/python-tests/test_*_stubs.py (if tier2_tests enabled)
```

**Step 3: After Design Review - Generate Implementation**
```
User: /generate-go-tests {JIRA_ID}
User: /generate-python-tests {JIRA_ID}
Output: Full working test implementations
```

---

## Notes

- **STD = Test stubs with docstrings** (per SOFTWARE_TEST_DESCRIPTION.md)
- **STD YAML = Internal format** (for automation, not for review)
- **Two-phase workflow:**
  - Phase 1 (this command): Generate stubs for design review
  - Phase 2 (/generate-*-tests): Generate full implementation
- **Test stubs are excluded from execution:**
  - Go: `PendingIt()` with `Skip()`
  - Python: `__test__ = False` with `pass`

---

**End of STD Builder Command**
