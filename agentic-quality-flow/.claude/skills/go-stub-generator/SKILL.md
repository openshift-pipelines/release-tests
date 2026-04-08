---
name: go-stub-generator
description: Generate Go/Ginkgo test stubs with PSE comments from STD YAML (Phase 1 - design review)
model: claude-opus-4-6
---

# Go Stub Generator Skill

## Purpose

Generates **Go/Ginkgo test stubs** with PSE comments for design review.

**Output:** Test stubs with `PendingIt()` + `Skip()` + PSE comments (excluded from test execution)

**Key Principle:** The STD = the comments in the test files (no separate document needed for review).

---

## Input Required

- `jira_id`: Jira ticket ID (e.g., "CNV-66855")

**Prerequisites:**
- STD YAML file must exist at `outputs/std/{JIRA_ID}/{JIRA_ID}_test_description.yaml`

---

## Output

**Generated Files:**
```
outputs/std/{JIRA_ID}/go-tests/
├── {feature_name}_stubs_test.go           (stubs with PSE comments)
├── {another_feature}_stubs_test.go
└── ... (one file per feature group)
```

**File Characteristics:**
- **Language:** Go (Ginkgo v2 + Gomega)
- **Size:** 50-150 lines per file (PSE comments + PendingIt)
- **Status:** Test stubs excluded from execution (`PendingIt` + `Skip`)
- **Body:** `Skip("Phase 1: Design only")` only

---

## CRITICAL REQUIREMENT

**Generate ONE test stub per STD scenario. No exceptions.**

- CORRECT: 19 STD scenarios → 19 generated `PendingIt()` blocks
- WRONG: 19 STD scenarios → 7 test files (grouped without covering all scenarios)

**Pattern-based file grouping is allowed**, but **EVERY scenario must get a test stub**.

---

## Polarion Toggle

If `project_context.feature_toggles.polarion` is false, omit Polarion marker references from stubs.

## PSE Comment Format

### Package-Level Comment
```go
/*
{Feature Name} Tests

STP Reference: {STP_URL}
Jira: {JIRA_ID}
*/
```

### Describe Block Comments (Shared Preconditions)
```go
var _ = Describe("{Feature}", decorators.SigNetwork, Serial, func() {
    /*
    Markers:
        - tier1
        - gating

    Preconditions:
        - OpenShift cluster with CNV 4.22+
        - Multiple NetworkAttachmentDefinitions
    */
```

### It Block Comments (PSE Format)
```go
    /*
    [NEGATIVE] (if applicable)
    Preconditions:
        - VM running with original NAD

    Steps:
        1. Update VM spec to reference target NAD
        2. Wait for update to complete

    Expected:
        - VM is connected to target NAD network
    */
    PendingIt("[test_id:TS-CNV72329-001] should allow NAD swap while VM running", func() {
        Skip("Phase 1: Design only - awaiting implementation")
    })
```

---

## Workflow

### Step 1: Read STD YAML

Load `outputs/std/{JIRA_ID}/{JIRA_ID}_test_description.yaml`

**Extract:**
- Total scenario count: `len(scenarios)`
- All tier1 scenarios (filter by `tier: "Tier 1"`)

### Step 2: Group Scenarios by Pattern

**CRITICAL: This step is for file organization only. ALL scenarios must still get stubs.**

For each scenario:
1. Detect basic patterns (NAD type, OS type) for grouping
2. Assign scenario to a file group

Result: Map of `{file_name: [scenario1, scenario2, ...]}`

### Step 3: For Each File Group

Generate the stub file with ALL scenarios in the group.

**For each scenario in the group:**

1. **Extract PSE Information from STD YAML**
   - Preconditions: from `specific_preconditions` and `test_steps.setup`
   - Steps: from `test_steps.test_execution`
   - Expected: from `test_objective.acceptance_criteria` or `assertions`

2. **Generate PSE Comment Block**
   ```go
   /*
   Preconditions:
       - {precondition 1}
       - {precondition 2}

   Steps:
       1. {step 1}
       2. {step 2}

   Expected:
       - {expected outcome}
   */
   ```

3. **Generate PendingIt Block**
   ```go
   PendingIt("[test_id:{scenario.test_id}] {description}", func() {
       Skip("Phase 1: Design only - awaiting implementation")
   })
   ```

**After all scenarios in group processed:**

4. **Assemble File Structure**

   **Note:** The package name and import paths come from `{project_context.config_dir}/tier1.yaml` (e.g., `default_package`, `imports`). The decorators and SIG references also come from project config.

   ```go
   package network  // Read from tier1.yaml `default_package`

   import (
       . "github.com/onsi/ginkgo/v2"  // Read from tier1.yaml `imports`
   )

   /*
   {Feature} Tests

   STP Reference: {STP_URL}
   Jira: {JIRA_ID}
   */

   var _ = Describe("[{JIRA_ID}] {Feature}", decorators.SigNetwork, func() {
       /*
       Markers:
           - tier1

       Preconditions:
           - {shared preconditions}
       */

       Context("{context name}", func() {
           // All PendingIt blocks here
       })
   })
   ```

5. **Save File**
   - Derive filename from feature group (snake_case + _test.go)
   - Save to `outputs/std/{JIRA_ID}/go-tests/{feature_slug}_stubs_test.go`

### Step 4: Validate Complete Coverage

**CRITICAL VALIDATION - This step is MANDATORY**

After all files generated:

1. **Count STD scenarios:** Count all Tier 1 scenarios in STD: `N_std`
2. **Count generated stubs:** Count all `PendingIt()` blocks: `N_stubs`
3. **Verify completeness:**
   - If `N_stubs < N_std`: ERROR + list missing scenario IDs
   - If `N_stubs == N_std`: SUCCESS

### Step 5: Report Results

Generate summary with:
- Scenarios processed
- Go stub files generated
- Total lines
- List of generated files
- Coverage validation result

---

## STD YAML to PSE Comment Transformation

### PSE Boundary Rules

**These rules are mandatory and override the field mapping when there is a conflict.**

#### What goes in Preconditions (setup — before the test runs)
- ALL resource creation: VMs, NADs, pods, peer VMs, storage
- ALL baseline data recording: "Record MAC address", "Save IP address"
- Baseline state verification: confirming the starting state before the test action
- Any action that establishes the starting state for the test
- **Never** test environment requirements (cluster version, node count, storage class, operator version)

#### What goes in Steps (actions — during the test)
- The test action itself: "Patch VM spec", "Execute ping", "Wait for completion"
- ONLY actions that are part of the test execution
- **Never** resource creation (that's a Precondition)
- **Never** verification statements (that's Expected)

#### What goes in Expected (assertions — what the test verifies)
- Outcome verification: checking the RESULT of the test action
- Any sentence starting with "Verify", "Confirm", "Check", "Ensure", "Assert"
  that checks the OUTCOME of the test action belongs here, NOT in Steps
- Must be **concrete and verifiable**:
  - GOOD: "MAC address equals pre-change value"
  - GOOD: "Ping succeeds with 0% packet loss"
  - BAD: "Interfaces correctly configured" (missing the how)
- Must describe the **observable outcome**, not the internal mechanism

**Baseline vs Outcome verification:**
- Baseline verification (confirming starting state BEFORE the test action) → **Preconditions**
- Outcome verification (confirming result AFTER the test action) → **Expected**

If a baseline check fails, the test cannot run (setup error). If an outcome check
fails, the test ran but produced a wrong result (test failure).

#### Quick Reference

| Action | PSE Section | Example |
|--------|-------------|---------|
| Create VM/NAD/Pod | **Preconditions** | "Running VM with secondary interface" |
| Record baseline data | **Preconditions** | "MAC address and interface name recorded" |
| Verify baseline state | **Preconditions** | "VM is in Running state before test action" |
| Patch/Update resource | **Steps** | "Update VM spec to reference target NAD" |
| Wait for completion | **Steps** | "Wait for operation to complete" |
| Execute command | **Steps** | "Ping from VM-A to VM-B" |
| Verify/Confirm outcome | **Expected** | "Ping succeeds with 0% packet loss" |
| Assert state | **Expected** | "VM is Running after operation" |

### Field Mapping

| STD YAML Field | PSE Section | Transformation |
|----------------|-------------|----------------|
| `specific_preconditions[*].requirement` | Preconditions | Bullet list |
| `test_steps.setup` | Preconditions (context) | Add to preconditions |
| `test_steps.test_execution[*].action` | Steps | Numbered list — **filter out "Verify/Confirm" actions → move to Expected** |
| `test_objective.acceptance_criteria[0]` | Expected | Natural language |
| `assertions[*].description` | Expected (fallback) | Natural language |
| Title contains "fail"/"error"/"negative" | [NEGATIVE] marker | Prefix |

### Example Transformation

**STD YAML Input:**
```yaml
test_objective:
  title: "VM network interface can be swapped while running"
  acceptance_criteria:
    - "Swap completes without VM restart"
specific_preconditions:
  - requirement: "VM in Running state with original NAD"
test_steps:
  test_execution:
    - action: "Update VM spec to reference target NAD"
    - action: "Wait for update to complete"
```

**PSE Comment Output:**
```go
/*
Preconditions:
    - VM in Running state with original NAD

Steps:
    1. Update VM spec to reference target NAD
    2. Wait for update to complete

Expected:
    - VM is connected to target NAD network
*/
PendingIt("[test_id:TS-CNV72329-001] should allow NAD swap while VM running", func() {
    Skip("Phase 1: Design only - awaiting implementation")
})
```

---

## Success Criteria

Stub generation succeeds when:
- All STD Tier 1 scenarios have corresponding `PendingIt()` blocks
- Every scenario ID appears in generated code with `[test_id:TS-XXX]` label
- Valid Go syntax (proper imports, package declaration)
- Files saved to `outputs/std/{JIRA_ID}/go-tests/`

---

## Error Handling

**If STD file not found:**
- Error: "STD file not found at outputs/std/{JIRA_ID}/{JIRA_ID}_test_description.yaml"
- Suggestion: "Run `/std-builder {JIRA_ID}` first"
- Exit

**If no Tier 1 scenarios found:**
- Warning: "No Tier 1 scenarios found in STD"
- Exit (no stubs to generate)

---

**End of Go Stub Generator Skill**
