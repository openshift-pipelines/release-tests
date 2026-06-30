---
name: generate-python-tests
description: Generate working tier2 Python/pytest tests from STD YAML
argument-hint: <JIRA-ID>
allowed-tools: Read, Write, Edit, Task, Glob, Grep, LSP, Skill
---

# Generate Python Tests Command

Generates **full working Python/pytest test implementations** from STD YAML.

**Use this after design review is approved.** For test stubs (design phase), use `/std-builder` instead.

---

When the user runs `/generate-python-tests {JIRA_ID}`:

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

**Check tier2_tests toggle:**
If `project_context.feature_toggles.tier2_tests` is false:
- Output: "Project {project_context.display_name} does not support tier2 Python tests (tier2_tests is disabled)."
- Exit. Do not proceed.

## Step 1: Parse Jira ID

Extract the Jira ID from `project_context.jira_id` (e.g., CNV-66855, VIRTSTRAT-494).

## Step 2: Verify STD File Exists

**CRITICAL: Tier2 Python test generation requires an existing STD file.**

Check that the STD file exists:
```
outputs/std/{JIRA_ID}/{JIRA_ID}_test_description.yaml
```

**If STD file does NOT exist:**
- Inform the user: "STD file not found. Please run `/std-builder {JIRA_ID}` first."
- Exit - do not proceed

**If STD file exists:**
- Proceed to Step 3

## Step 3: Run LSP Pattern Analysis

**Toggle gate:** If `project_context.feature_toggles.lsp_analysis` is false, skip the ticket-context-analyzer. Proceed to Step 4 using the project pattern library fallback only (`{project_context.config_dir}/patterns/tier2_patterns.yaml`).

**CRITICAL: LSP analysis before code generation ensures accuracy**

Use the Task tool to spawn the ticket-context-analyzer agent:

**Tool:** Task
**Parameters:**
- subagent_type: "general-purpose"
- description: "LSP pattern analysis for {JIRA_ID}"
- prompt: |
    Read and follow the ticket-context-analyzer agent instructions.

    Analyze patterns for:
    - jira_id: "{JIRA_ID}"
    - std_file_path: "outputs/std/{JIRA_ID}/{JIRA_ID}_test_description.yaml"
    - tier: "tier2"
    - config_dir: "{project_context.config_dir}"

    Read repo_paths from `{project_context.config_dir}/repositories.yaml` to determine
    which repositories to analyze.

    Output:
    - outputs/python-tests/{JIRA_ID}/{JIRA_ID}_lsp_patterns_tier2.yaml (detailed patterns)
    - outputs/python-tests/{JIRA_ID}/{JIRA_ID}_lsp_summary.md (summary)

**This step:**
- Uses LSP (NOT grep) for semantic analysis
- Reads repository paths from project config (config_dir/repositories.yaml)
- Extracts current function signatures from configured repos
- Finds real usage examples from test files

## Step 4: Invoke python-test-generator Skill

Use the Skill tool to invoke the python-test-generator skill:

**Tool:** Skill
**Parameters:**
- skill: "python-test-generator"
- args: "{JIRA_ID}"

The python-test-generator skill will:
1. Read the STD file
2. **Read LSP patterns file** (outputs/python-tests/{JIRA_ID}/{JIRA_ID}_lsp_patterns_tier2.yaml)
3. Use fresh patterns for code generation (LSP-verified)
4. Generate **working** Python/pytest test files (full implementation)
5. Generate conftest.py if needed (shared fixtures)
6. Validate generated code (syntax + pytest collection)
7. Save to `outputs/python-tests/{JIRA_ID}/`

## Step 5: Report Results

Once the skill completes, show the user:

```
✅ Tier2 Python Test Generation Complete!

📄 Input: outputs/std/{JIRA_ID}/{JIRA_ID}_test_description.yaml

📊 Summary:
- Scenarios processed: {COUNT}
- Python test files generated: {COUNT}
- Conftest files generated: {COUNT}
- Total lines of code: {COUNT}

📁 Generated Files:
- outputs/python-tests/{JIRA_ID}/test_{file1}.py ({lines} lines)
- outputs/python-tests/{JIRA_ID}/test_{file2}.py ({lines} lines)
- outputs/python-tests/{JIRA_ID}/conftest.py ({lines} lines)
- ...

✅ Tests are ready to run!

📌 Next steps:
   cd outputs/python-tests/{JIRA_ID}
   pytest --collect-only   # Verify collection
   pytest -v               # Run tests

{Any errors or warnings}
```

---

## Example Usage

**Full Workflow:**
```
# Step 1: Generate STP
User: /stp-builder {JIRA_ID}
Output: outputs/stp/{JIRA_ID}/{JIRA_ID}_test_plan.md

# Step 2: Generate STD (stubs for design review)
User: /std-builder {JIRA_ID}
Output:
   - outputs/std/{JIRA_ID}/{JIRA_ID}_test_description.yaml
   - outputs/std/{JIRA_ID}/python-tests/test_*_stubs.py (stubs)

# Step 3: Team reviews stubs, approves design

# Step 4: Generate full implementation
User: /generate-python-tests {JIRA_ID}
Output:
   - outputs/python-tests/{JIRA_ID}/test_*.py (working code)
   - outputs/python-tests/{JIRA_ID}/conftest.py
```

---

## Error Handling

**If STD file not found:**
- Error: "STD file not found at outputs/std/{JIRA_ID}/{JIRA_ID}_test_description.yaml"
- Suggestion: "Please run `/std-builder {JIRA_ID}` first"
- Exit without invoking skill

**If skill fails:**
- Display error message from skill
- List any Python test files that were successfully generated before the failure
- Report which scenario IDs are missing test coverage
- Suggest reviewing errors and re-running for failed scenarios

---

## Prerequisites

Before running this command:
1. ✅ STP file must exist (run `/stp-builder {JIRA_ID}`)
2. ✅ STD file must exist (run `/std-builder {JIRA_ID}`)
3. ✅ Design review should be complete (stubs approved)
4. ✅ STD must contain tier2 test scenarios

---

## Workflow Overview

```
/stp-builder {JIRA_ID}
  ↓
outputs/stp/{JIRA_ID}/{JIRA_ID}_test_plan.md (STP)
  ↓
/std-builder {JIRA_ID}
  ↓
outputs/std/{JIRA_ID}/{JIRA_ID}_test_description.yaml (STD YAML)
outputs/std/{JIRA_ID}/python-tests/test_*_stubs.py (stubs for review)
  ↓
[Design Review & Approval]
  ↓
/generate-python-tests {JIRA_ID}
  ↓
Step 0: Resolve project → check tier2_tests toggle
  ↓
outputs/python-tests/{JIRA_ID}/test_*.py (full working implementation)
outputs/python-tests/{JIRA_ID}/conftest.py (shared fixtures)
```

---

## Notes

- **This command generates WORKING code** (not stubs)
- **For stubs:** Use `/std-builder` instead
- **Writes to separate directory from stubs:** Working implementations go to `outputs/python-tests/` while stubs remain in `outputs/std/`
- **LSP-verified:** Uses fresh LSP patterns for accuracy (generated per-ticket)
- **Fixture-based:** Leverages pytest fixture architecture
- **Context managers:** All resources use `with` statements for cleanup

---

**End of Generate Python Tests Command**
