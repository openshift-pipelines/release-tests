---
name: std-orchestrator
description: Orchestrate STP → STD pipeline (generates comprehensive STD YAML only)
model: claude-opus-4-6
phase_support: true
default_phase: phase1
---

# STD Orchestrator Skill

## Purpose

Coordinates the Software Test Description (STD) generation workflow by:
1. Parsing STP Section III to extract test scenarios
2. Generating comprehensive STD YAML file for ALL scenarios (internal format)
3. **Routing to code generators with appropriate phase**
4. Validating output and generating summary report

**Output:**
- STD YAML file (internal format for automation)
- **Test stubs with PSE docstrings** (the actual deliverable per SOFTWARE_TEST_DESCRIPTION.md)

---

## Input Required

- `stp_file_path`: Path to the STP markdown file (e.g., `outputs/stp/CNV-66855/CNV-66855_test_plan.md`)
- `jira_id`: The Jira ticket ID (e.g., "CNV-66855")
- `output_dir`: Base directory for outputs (defaults to `outputs/std/{JIRA_ID}/`)
- `phase`: `phase1` (default) or `phase2`

---

## Phase Parameter

**Phase 1 (default):**
- Generate STD YAML (internal format)
- Call code generators with `phase=phase1`
- Output: Test stubs with PSE docstrings + `pass` body
- Tests excluded from collection (`__test__ = False` for Python, `PendingIt()` for Go)

**Phase 2:**
- Generate STD YAML (internal format)
- Call code generators with `phase=phase2`
- Output: Full working test implementations

---

## Workflow

Execute the following steps in order:

---

### Step 1: Parse STP Section III

**Read the STP file and extract all test scenarios from Section III (Test Scenarios & Traceability).**

**Expected bullet-based format:**
```markdown
- **[CNV-12345]** — As a user, I want to reset a VM
  - *Test Scenario:* [Tier 1] Verify basic reset operation succeeds
  - *Priority:* P0

- **[CNV-12345]** — As a user, I want to verify reset preserves state
  - *Test Scenario:* [Tier 1] Verify reset preserves pod UID
  - *Priority:* P0
```

**Parse and extract:**
- Requirement ID (Jira key from `**[ID]**`)
- Requirement summary (text after `— `)
- Tier classification from test scenario (Tier 1, Tier 2)
- Priority (P0, P1, P2)
- Scenario description (from `*Test Scenario:*` line)

**Store as:**
```yaml
scenarios:
  - scenario_id: 1
    tier: "Tier 1"
    priority: "P0"
    description: "Verify basic reset operation succeeds"
  - scenario_id: 2
    tier: "Tier 1"
    priority: "P0"
    description: "Verify reset preserves pod UID"
```

---

### Step 2: Generate Comprehensive STD YAML (Single File)

**Generate ONE comprehensive STD file for ALL scenarios:**

1. **Extract STP context** (needed by std-generator):
   - Jira issue metadata (from Metadata & Tracking)
   - Feature description (from Feature Overview)
   - Known limitations (from Section I.2)
   - API extensions (from Section I.3 API Extensions checkbox)
   - Test environment (from Section II.3)
   - Source bugs (if Closed Loop ticket)
   - Fix versions

2. **Call std-generator skill ONCE** with:
   - ALL scenarios array (from Step 1)
   - STP context
   - STP file path

3. **Output file:**
   - `outputs/std/{JIRA_ID}/{JIRA_ID}_test_description.yaml`
   - Example: `outputs/std/CNV-66855/CNV-66855_test_description.yaml`
   - Single comprehensive file with:
     - document_metadata (shared across all scenarios)
     - common_preconditions (shared infrastructure)
     - scenarios array (one entry per STP scenario)

4. **Validate STD output:**
   - File exists
   - Valid YAML syntax
   - All required sections populated:
     - document_metadata
     - common_preconditions
     - scenarios array (count matches STP scenarios)
   - Each scenario has required fields:
     - test_id, tier, priority
     - test_objective, test_steps, assertions

---

### Step 2.5: Route to Code Generators (Based on Phase)

**After STD YAML is validated, call appropriate code generators:**

**If phase == "phase1" (default):**
1. Detect tier split in STD (Tier 1 count, Tier 2 count)
2. **If Tier 1 scenarios exist:**
   - Call go-test-generator with `phase=phase1`
   - Output: Go/Ginkgo test stubs with PSE comments
3. **If Tier 2 scenarios exist:**
   - Call python-test-generator with `phase=phase1`
   - Output: Python/pytest test stubs with PSE docstrings

**If phase == "phase2":**
1. Same routing as phase1, but with `phase=phase2`
2. Output: Full working implementations

**Output files:**
```
outputs/std/{JIRA_ID}/
├── go-tests/           (if Tier 1 scenarios exist)
│   └── *_stubs_test.go (Phase 1: stubs, Phase 2: implementation)
└── python-tests/       (if Tier 2 scenarios exist)
    └── test_*_stubs.py (Phase 1: stubs, Phase 2: implementation)
```

---

### Step 3: Generate Summary Report

**Create a summary report with:**
- Total scenarios processed
- STD file generated
- **Phase indicator** (Phase 1 stubs or Phase 2 implementation)
- **Code files generated** (Go and/or Python)
- Validation results
- Execution time
- Any errors or warnings

**Output format:**
```yaml
---
status: success
component: std-orchestrator
jira_id: CNV-66855
phase: phase1  # or phase2
stp_file: outputs/stp/CNV-66855/CNV-66855_test_plan.md
output_dir: outputs/std/CNV-66855/

execution_summary:
  total_stp_scenarios: 12
  tier_1_scenarios: 9
  tier_2_scenarios: 3
  std_file_generated: "CNV-66855_test_description.yaml"
  scenarios_in_std: 12
  total_duration: "2 minutes"

code_generation:
  phase: phase1
  go_tests:
    file_count: 2
    test_count: 9
    status: "stubs_generated"  # or "implementation_generated" for phase2
  python_tests:
    file_count: 1
    test_count: 3
    status: "stubs_generated"

validation_results:
  std_file:
    file: CNV-66855_test_description.yaml
    status: valid
    yaml_syntax: passed
    required_sections: passed
    scenarios_count: 12

errors: []
warnings: []

notes:
  - "STD YAML generated as internal format"
  - "Use /generate-go-tests or /generate-python-tests for implementations"
---
```

---

## Output Structure

**Simple structure (STD YAML only):**
```
outputs/std/CNV-66855/
├── CNV-66855_test_description.yaml     (NEW - comprehensive STD for ALL scenarios)
└── std_generation_summary.yaml         (summary report)
```

**Key design:**
- **ONE comprehensive STD file** for all scenarios (not one file per scenario)
- **STD mirrors STP structure:** document_metadata + common_preconditions + scenarios array
- **No separate std/ folder** - single file at outputs/std/{JIRA_ID}/ level
- **No test stubs** - STD YAML is input for code generators
- **Downstream usage:** /generate-go-tests or /generate-python-tests read this STD file

---

## Skills Called

This orchestrator calls 1 specialized skill:

1. **std-generator** - Transforms STP scenarios → comprehensive STD YAML file

**Architecture:**
- STD YAML is the only output
- Code generation happens in separate commands (/generate-go-tests, /generate-python-tests)
- Clean separation: specification (STD) vs implementation (code)

---

## Error Handling

- **If STP parsing fails:**
  - Log error: "Cannot parse Section III from STP file"
  - Suggest: Check STP format, ensure Section III exists
  - Exit with status: error

- **If std-generator fails for a scenario:**
  - Log warning: "STD generation failed for scenario {num}"
  - Continue with other scenarios
  - Mark scenario as failed in summary

- **Set overall status to:**
  - `success`: All scenarios processed successfully
  - `partial`: Some scenarios failed, but >50% succeeded
  - `error`: >50% scenarios failed or critical error

---

## Success Criteria

The orchestration is complete when:
- ✅ All scenarios from STP Section III extracted
- ✅ Comprehensive STD YAML file created
- ✅ Valid YAML syntax
- ✅ All required sections populated
- ✅ Summary report generated

**Minimum acceptable outcome:**
- At least 80% of scenarios successfully included in STD
- All P0 scenarios successfully included
- Summary report explains any failures
- STD file passes validation

---

## Validation Checklist

Before marking orchestration as complete, validate:

- [ ] STD YAML file exists
- [ ] Valid YAML syntax (can be parsed)
- [ ] document_metadata section populated
- [ ] common_preconditions section populated
- [ ] scenarios array contains all STP scenarios
- [ ] Each scenario has required fields (test_id, tier, priority, test_objective, test_steps, assertions)
- [ ] No missing or null values in critical fields
- [ ] Test IDs are unique

---

## Usage Example

**User command (Phase 1 - default):**
```
Generate STD/PSE/Code for CNV-66855
```

**Orchestrator execution (Phase 1):**
```
1. Read outputs/stp/CNV-66855/CNV-66855_test_plan.md
2. Parse Section III → 12 scenarios found (9 Tier 1, 3 Tier 2)
3. Call std-generator ONCE → CNV-66855_test_description.yaml
4. Validate STD YAML
5. Call go-test-generator (phase=phase1) → 9 test stubs in go-tests/
6. Call python-test-generator (phase=phase1) → 3 test stubs in python-tests/
7. Generate summary → std_generation_summary.yaml
8. Report to user: "✅ Generated Phase 1 test stubs for 12 scenarios"
```

**Example output (Phase 1):**
```
✅ Phase 1 Test Stubs Generated!

📄 Input: outputs/stp/CNV-66855/CNV-66855_test_plan.md

📊 Summary:
- STP scenarios: 12 (9 Tier 1, 3 Tier 2)
- STD file: CNV-66855_test_description.yaml (internal format)
- Phase: 1 (Design stubs with PSE docstrings)

📁 Output:
- outputs/std/CNV-66855/CNV-66855_test_description.yaml
- outputs/std/CNV-66855/go-tests/ (9 test stubs with PSE comments)
- outputs/std/CNV-66855/python-tests/ (3 test stubs with PSE docstrings)

📋 Phase 1 Checklist:
- [ ] STP link in module docstring
- [ ] Tests grouped in class with shared preconditions
- [ ] Each test has: Preconditions, Steps, Expected
- [ ] Each test verifies ONE thing with ONE Expected
- [ ] Test bodies contain only 'pass'

✅ Ready for design review!

📌 Next steps:
1. Review the test stubs
2. Submit PR for design review
3. After approval, run:
   - /generate-go-tests CNV-66855     (Tier 1 implementation)
   - /generate-python-tests CNV-66855 (Tier 2 implementation)
```

---

## Notes

- **Output**: Single comprehensive STD YAML file only
- **No code generation**: Use /generate-go-tests or /generate-python-tests for code
- **Single comprehensive STD**: ONE file for ALL scenarios (mirrors STP structure)
- **STD replaces old multi-file approach**: More maintainable, less duplication
- **Clean separation**: Specification (STD) vs implementation (code generation)

---

**End of STD Orchestrator Skill**
