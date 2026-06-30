---
name: python-stub-generator
description: Generate Python/pytest test stubs with PSE docstrings from STD YAML (Phase 1 - design review)
model: claude-opus-4-6
---

# Python Stub Generator Skill

## Purpose

Generates **Python/pytest test stubs** with PSE docstrings for design review.

**Output:** Test stubs with `pass` body + `__test__ = False` (excluded from pytest collection)

**Key Principle:** The STD = the docstrings in the test files (no separate document needed for review).

---

## Input Required

- `jira_id`: Jira ticket ID (e.g., "CNV-66855")

**Prerequisites:**
- STD YAML file must exist at `outputs/std/{JIRA_ID}/{JIRA_ID}_test_description.yaml`

---

## Output

**Generated Files:**
```
outputs/std/{JIRA_ID}/python-tests/
├── test_{feature_name}_stubs.py           (stubs with PSE docstrings)
├── test_{another_feature}_stubs.py
└── ... (one file per feature group)
```

**File Characteristics:**
- **Language:** Python 3.x
- **Framework:** pytest
- **Size:** 50-150 lines per file (docstrings + pass)
- **Status:** Test stubs excluded from collection (`__test__ = False`)
- **Body:** `pass` only (no implementation)

---

## CRITICAL REQUIREMENT

**Generate ONE test stub per STD scenario. No exceptions.**

- CORRECT: 8 STD scenarios → 8 generated `def test_*()` functions
- WRONG: 8 STD scenarios → 3 test files (grouped without covering all scenarios)

**Pattern-based file grouping is allowed**, but **EVERY scenario must get a test stub**.

---

## PSE Docstring Format

**Note:** Import patterns and Polarion config come from `{project_context.config_dir}/tier2.yaml`. If the project has `polarion: false` in its config, skip Polarion marker references.

### Module Docstring (Required)
```python
"""
{Feature Name} Tests

STP Reference: {STP_URL}
Jira: {JIRA_ID}
"""
```

**Module docstring contains ONLY:**
- STP Reference (path to the approved STP)
- Jira issue ID

**Do NOT include in module docstring:**
- PR references or URLs
- Feature gate names
- VEP/enhancement references
- Implementation details

### Class Docstring (Shared Preconditions)
```python
class TestFeatureName:
    """
    Tests for {feature description}.

    Markers:
        - tier2
        - gating

    Parametrize:
        - storage_class: [ocs-storagecluster-ceph-rbd, hostpath-csi]

    Preconditions:
        - {Shared precondition 1 — resource creation}
        - {Shared precondition 2 — baseline data recorded}
    """
    __test__ = False
```

**Class-level Preconditions include ONLY test-specific setup:**
- Resource creation (VMs, NADs, peer VMs)
- Baseline data recording (MAC addresses, IP addresses, interface names)

**Do NOT include in Preconditions:**
- Test environment requirements (cluster version, node count, storage type, network infrastructure)
- Platform prerequisites (OCP version, CNV version, operator installations)
- Cluster configuration that the STP Test Environment section already describes

Tests assume the test environment described in the STP (Section II.3) is already in place.
Tests that share the same VM/NAD setup MUST be grouped in one class.

### Standalone Test Function (No Class Needed)

When a test stands alone without related tests, a class is not required:

```python
def test_specific_behavior():
    """
    Test that {specific ONE thing being verified}.

    Markers:
        - gating

    Parametrize:
        - os_image: [rhel9, fedora]

    Preconditions:
        - {Setup requirement}

    Steps:
        1. {Discrete action}

    Expected:
        - {Concrete, verifiable assertion}
    """

test_specific_behavior.__test__ = False
```

For standalone tests, `__test__ = False` goes AFTER the function definition.

### Test Docstring (PSE Format — Class Method)
```python
def test_specific_behavior(self):
    """
    Test that {specific ONE thing being verified}.

    Preconditions:  # Optional - only if test-specific beyond class
        - {Test-specific precondition}

    Steps:
        1. {Discrete action — patch, execute, wait}

    Expected:
        - {Concrete, verifiable assertion — e.g., "VM is Running"}
    """
```

**No fixture parameters in signature.** Write `def test_foo(self):` only.

### Parametrize Section

When a test should run with multiple parameter combinations, add a `Parametrize:` section:

```python
def test_online_disk_resize(self):
    """
    Test that a running VM's disk can be expanded.

    Parametrize:
        - storage_class: [ocs-storagecluster-ceph-rbd, hostpath-csi]

    Preconditions:
        - Storage class from parameter exists
        - Running VM with a DataVolume as boot disk

    Steps:
        1. Expand PVC by 1Gi

    Expected:
        - Disk size inside VM is greater than original size
    """
```

### Negative Test Indicator
```python
def test_failure_scenario(self):
    """
    [NEGATIVE] Test that {failure scenario description}.
    ...
    """
```

### Dependent Tests (Incremental)

When tests within a class depend on the execution order of previous tests,
use `@pytest.mark.incremental` marker in the class Markers section:

```python
class TestVMSomeFeature:
    """
    Tests for VM feature with ordered dependencies.

    Markers:
        - incremental

    Preconditions:
        - Running VM with feature configured
    """
    __test__ = False

    def test_vm_is_created(self):
        """Test that a VM with feature can be created."""

    def test_vm_migration(self):
        """Test that a VM with feature can be migrated."""
```

---

## Phase 1 Prohibitions

Phase 1 stubs are **design-only**. The following are implementation details
and MUST NOT appear in Phase 1 output:

- **No `@pytest.fixture` definitions** — Fixtures are Phase 2 implementation
- **No `@pytest.mark.*` decorators** — Use `Markers:` section in docstrings instead
- **No fixture parameters in test signatures** — Write `def test_foo(self):` not
  `def test_foo(self, bridge_nads, vm_with_secondary_interface):`
- **No `import pytest`** — Not needed when there are no decorators or fixtures
- **No PR references** — PRs are STP-level context, not STD
- **No block comments above tests** — All test information goes in the docstring
- **No fixture names in Preconditions** — Use descriptive requirements, not fixture names
  - GOOD: "Running Fedora virtual machine"
  - BAD: "Running Fedora VM (vm_to_restart fixture)"

**What Phase 1 stubs contain:**
- Module docstring (STP Reference + Jira only)
- Classes with docstring (shared Preconditions + Markers + optional Parametrize)
- `__test__ = False` on class level (for grouped tests) or after function definition (for standalone tests)
- Test methods with PSE docstrings (no body needed — method body is empty)
- Standalone test functions with `test_name.__test__ = False` after definition
- Nothing else

---

## Workflow

### Step 1: Read STD YAML

Load `outputs/std/{JIRA_ID}/{JIRA_ID}_test_description.yaml`

**Extract:**
- Total scenario count: `len(scenarios)`
- All tier2 scenarios (filter by `tier: "Tier 2"`)

### Step 2: Group Scenarios by Pattern

**CRITICAL: This step is for file organization only. ALL scenarios must still get stubs.**

For each scenario:
1. Detect basic patterns (localnet, NAD type, VM type) for grouping
2. Assign scenario to a file group

Result: Map of `{file_name: [scenario1, scenario2, ...]}`

### Step 3: For Each File Group

Generate the stub file with ALL scenarios in the group.

**For each scenario in the group:**

1. **Extract PSE Information from STD YAML**
   - Preconditions: from `specific_preconditions` and `test_steps.setup`
   - Steps: from `test_steps.test_execution`
   - Expected: from `test_objective.acceptance_criteria` or `assertions`

2. **Generate Test Function with PSE Docstring**
   ```python
   def test_{scenario_slug}(self):
       """
       Test that {test_objective.title}.

       Preconditions:
           - {precondition 1}
           - {precondition 2}

       Steps:
           1. {step 1}
           2. {step 2}

       Expected:
           - {expected outcome}
       """
       pass
   ```

**After all scenarios in group processed:**

3. **Assemble File Structure**
   ```python
   """
   {Feature} Tests

   STP Reference: {STP_URL}
   Jira: {JIRA_ID}
   """


   class TestFeatureName:
       """
       Tests for {feature description}.

       Markers:
           - tier2

       Preconditions:
           - {shared preconditions — ALL resource creation here}
           - {baseline data recording here}
       """
       __test__ = False

       # All test functions here (no fixture params, no decorators)
   ```

   **Note:** No `import pytest`, no `@pytest.fixture`, no `@pytest.mark` decorators.
   Markers are documented in docstring `Markers:` sections only.

4. **Save File**
   - Derive filename: `test_{feature_slug}.py`
   - Save to `outputs/std/{JIRA_ID}/python-tests/test_{feature_slug}_stubs.py`

### Step 4: Validate Complete Coverage

**CRITICAL VALIDATION - This step is MANDATORY**

After all files generated:

1. **Count STD scenarios:** Count all Tier 2 scenarios in STD: `N_std`
2. **Count generated stubs:** Count all `def test_*()` functions: `N_stubs`
3. **Verify completeness:**
   - If `N_stubs < N_std`: ERROR + list missing scenario IDs
   - If `N_stubs == N_std`: SUCCESS

### Step 5: Report Results

Generate summary with:
- Scenarios processed
- Python stub files generated
- Total lines
- List of generated files
- Coverage validation result

---

## STD YAML to PSE Docstring Transformation

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
- Must be **concrete and verifiable** using assertion wording patterns:
  - GOOD: "MAC address equals pre-change value"
  - GOOD: "Ping succeeds with 0% packet loss"
  - GOOD: "VM is Running"
  - BAD: "Interfaces correctly configured" (missing the how)
  - BAD: "Everything works" (not verifiable)
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
| Patch/Update resource | **Steps** | "Patch VM spec to change NAD reference" |
| Wait for completion | **Steps** | "Wait for update to complete" |
| Execute command | **Steps** | "Execute ping from VM-A to VM-B" |
| Verify/Confirm outcome | **Expected** | "Ping succeeds with 0% packet loss" |
| Assert state | **Expected** | "VM is Running" |

### Field Mapping

| STD YAML Field | PSE Section | Transformation |
|----------------|-------------|----------------|
| `test_objective.title` | Brief description | "Test that {title}" |
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

**PSE Docstring Output:**
```python
def test_nad_swap_while_running(self):
    """
    Test that VM network interface can be swapped while running.

    Preconditions:
        - VM in Running state with original NAD

    Steps:
        1. Update VM spec to reference target NAD
        2. Wait for update to complete

    Expected:
        - VM is connected to target NAD network
    """
    pass
```

---

## Assertion Wording Patterns (for Expected section)

Use clear, natural language that maps directly to assertions:

| Wording Pattern | Maps To |
|-----------------|---------|
| `X equals Y` | `assert x == y` |
| `X does not equal Y` | `assert x != y` |
| `VM is "Running"` | `assert vm.status == Running` |
| `VM is not running` | `assert vm.status != Running` |
| `File exists` / `Resource x exists` | `assert exists(x)` |
| `File does not exist` / `Resource x does NOT exist` | `assert not exists(x)` |
| `X does not contain Y` | `assert y not in x` |
| `Ping succeeds` / `Operation succeeds` | `assert operation()` (no exception) |
| `Ping fails` / `Operation fails` | `assert` raises exception or returns failure |
| `X contains Y` | `assert y in x` |

---

## Success Criteria

Stub generation succeeds when:
- All STD Tier 2 scenarios have corresponding `def test_*()` functions
- Every scenario has PSE docstring (Preconditions/Steps/Expected)
- Each test verifies **ONE thing** with ONE Expected
- Related tests are grouped in classes with shared preconditions
- Standalone tests use `test_name.__test__ = False` after definition
- Classes with grouped tests have `__test__ = False` on the class
- Negative tests are marked with `[NEGATIVE]` in the description
- `Markers:` section used for pytest markers (not decorators)
- `Parametrize:` section used for parameter combinations (when applicable)
- No fixture names in Preconditions (use descriptive requirements)
- Valid Python syntax
- Files saved to `outputs/std/{JIRA_ID}/python-tests/`

---

## Error Handling

**If STD file not found:**
- Error: "STD file not found at outputs/std/{JIRA_ID}/{JIRA_ID}_test_description.yaml"
- Suggestion: "Run `/std-builder {JIRA_ID}` first"
- Exit

**If no Tier 2 scenarios found:**
- Warning: "No Tier 2 scenarios found in STD"
- Exit (no stubs to generate)

---

**End of Python Stub Generator Skill**
