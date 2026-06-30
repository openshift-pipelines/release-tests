---
name: go-test-generator
description: Generate working tier1 Go/Ginkgo test code from STD YAML (full implementation)
model: claude-opus-4-6
---

# Go Test Generator Skill (Tier1)

## Purpose

Generates **working tier1 Go/Ginkgo test code** from STD YAML specifications.

**Output:** Working Go test files that compile with Bazel in kubevirt/kubevirt repository

**Note:** For test stubs (design review), use `go-stub-generator` instead.

---

## Input Required

- `jira_id`: Jira ticket ID (e.g., "CNV-66855")

**Prerequisites:**
- STD YAML file must exist at `outputs/std/{JIRA_ID}/{JIRA_ID}_test_description.yaml`
- STD must contain test scenarios

---

## Output

**Generated Files:**
```
outputs/go-tests/{JIRA_ID}/
├── {feature_name}_test.go           (working implementation)
├── {another_feature}_test.go
└── ... (one file per feature group)
```

**File Characteristics:**
- **Language:** Go (Ginkgo v2 + Gomega)
- **Size:** 200-500 lines per file
- **Status:** Working code (compiles with Bazel)
- **Format:** Follows kubevirt/kubevirt tier1 test patterns

---

## CRITICAL REQUIREMENT

**Generate ONE test case per STD scenario. No exceptions.**

- ✅ CORRECT: 19 STD scenarios → 19 generated `It()` or `PendingIt()` blocks
- ❌ WRONG: 19 STD scenarios → 7 test files (grouped without covering all scenarios)

**Pattern-based file grouping is allowed**, but **EVERY scenario must get a test case**.

---

## Polarion Toggle

If `project_context.feature_toggles.polarion` is false, omit Polarion test case ID markers from generated test code.

## Workflow

### Step 1: Read STD YAML

Load `outputs/std/{JIRA_ID}/{JIRA_ID}_test_description.yaml`

Expected structure:
```yaml
document_metadata:
  jira_id: CNV-66855
  title: "Feature title"
  tier: tier1

scenarios:
  - test_id: TS-CNV66855-001
    description: "Test scenario description"
    steps:
      - "Step 1"
      - "Step 2"
    assertions:
      - "Expected outcome"
```

**Extract:**
- Total scenario count: `len(scenarios)`
- All tier1 scenarios (filter by `tier: "Tier 1"`)

### Step 2: Load Pattern Rules

Read `{project_context.config_dir}/patterns/tier1_patterns.yaml` for:
- NAD type detection rules
- OS type detection rules
- Connectivity pattern rules
- Template selection logic

### Step 3: Group Scenarios by Pattern

**CRITICAL: This step is for file organization only. ALL scenarios must still get tests.**

For each scenario:
1. Detect patterns (NAD type, OS type, connectivity type)
2. Select template based on priority rules
3. Assign scenario to a file group

Result: Map of `{file_name: [scenario1, scenario2, ...]}

Example:
```
{
  "localnet_connectivity_test.go": [TS-001, TS-002, TS-005, TS-006],
  "openflow_validation_test.go": [TS-003, TS-004],
  "route_localnet_test.go": [TS-007, TS-008, TS-009],
  "vm_operations_test.go": [TS-010, TS-011]
}
```

### Step 4: For Each File Group

Generate the test file with ALL scenarios in the group.

**For each scenario in the group:**

1. **Detect Patterns** (from `{project_context.config_dir}/patterns/tier1_patterns.yaml`)
   - NAD type: localnet, bridge, SR-IOV, macvtap
   - OS type: Fedora, Alpine, CirrOS
   - Connectivity: ping, TCP, HTTP
   - Special: IPv4/IPv6, migration

2. **Generate Test Case Code**
   - Create an `It()` block with `[test_id:{scenario.test_id}]` label
   - Use scenario description as test description
   - Populate test body based on detected patterns
   - Include all steps from STD scenario
   - Include all assertions from STD scenario

3. **Generate Connectivity Test Code** (if needed)
   - Based on detected connectivity patterns (ping, TCP, HTTP)
   - Generate appropriate test code snippets

**After all scenarios in group processed:**

4. **Select Template** (priority-based from `{project_context.config_dir}/patterns/tier1_patterns.yaml`)
   - Priority 1: `{project_context.config_dir}/templates/tier1/parametric_ipv4_ipv6_test.go.template` (if dual stack)
   - Priority 2: `{project_context.config_dir}/templates/tier1/migration_test.go.template` (if migration + connectivity)
   - Priority 3: `{project_context.config_dir}/templates/tier1/network_connectivity_test.go.template` (if NAD + connectivity)
   - Priority 4: `{project_context.config_dir}/templates/tier1/basic_vmi_test.go.template` (default)

5. **Read Template**
   - Load selected template from `{project_context.config_dir}/templates/tier1/` directory
   - Reference tests are available at `{project_context.config_dir}/reference/tier1/`

6. **Populate File-Level Placeholders**
   - `{{PACKAGE_NAME}}` → "network"
   - `{{TEST_SUITE_NAME}}` → derived from file group name
   - `{{IMPORTS}}` → all required imports (collected from all scenarios)
   - `{{SETUP_CODE}}` → shared BeforeEach/BeforeAll setup
   - `{{TEST_CASES}}` → ALL generated It() blocks from step 2

7. **Validate Generated Code**
   - Check all imports correct
   - Verify `WaitUntilVMIReady` has login function parameter
   - Ensure no syntax errors (pattern-based validation)
   - **CRITICAL:** Verify number of `It()` blocks equals number of scenarios in group

8. **Save File**
   - Derive filename from file group name (snake_case + _test.go)
   - Save to `outputs/go-tests/{JIRA_ID}/{feature_slug}_test.go`

### Step 5: Validate Complete Coverage

**CRITICAL VALIDATION - This step is MANDATORY**

After all files generated:

1. **Count STD scenarios:**
   - Count all Tier 1 scenarios in STD: `N_std`

2. **Count generated test cases:**
   - Count all `It()` blocks with `[test_id:` in generated files: `N_tests`

3. **Verify completeness:**
   - If `N_tests < N_std`:
     - ERROR: "Incomplete coverage - {N_std - N_tests} scenarios missing"
     - List missing scenario IDs
     - FAIL generation
   - If `N_tests == N_std`:
     - SUCCESS: "Complete coverage - all {N_std} scenarios have tests"

4. **Scenario ID mapping:**
   - For each STD scenario ID, verify it appears in generated code
   - Report any missing IDs

### Step 6: Report Results

Generate summary report with:
- Scenarios processed count
- Go test files generated count
- Total lines of code
- List of generated files with line counts
- Any errors or warnings

---

## Pattern Detection Examples

**NAD Type Detection:**
```
"localnet" OR "passt" → libnet.NewPasstNetAttachDef()
"bridge network" → libnet.NewBridgeNetAttachDef()
"SR-IOV" OR "sriov" → libnet.NewSriovNetAttachDef()
```

**OS Type Detection:**
```
"Fedora" → libvmifact.NewFedora() + console.LoginToFedora
"Alpine" → libvmifact.NewAlpine() + console.LoginToAlpine
"CirrOS" → libvmifact.NewCirros() + console.LoginToCirros
```

**Template Selection:**
```
IF "IPv4 and IPv6" → parametric_ipv4_ipv6_test.go.template
ELIF "migration" AND "connectivity" → migration_test.go.template
ELIF "NAD" AND "connectivity" → network_connectivity_test.go.template
ELSE → basic_vmi_test.go.template
```

---

## Connectivity Test Code Generation

**Ping Test:**
```go
pingCmd := []string{"ping", "-c", "5", vmiIP}
output, err := exec.ExecuteCommandOnPod(testPod, testPod.Spec.Containers[0].Name, pingCmd)
ExpectWithOffset(1, err).ToNot(HaveOccurred())
ExpectWithOffset(1, output).To(ContainSubstring("0% packet loss"))
```

**TCP Test:**
```go
ncCmd := []string{"nc", "-zv", vmiIP, "22"}
output, err := exec.ExecuteCommandOnPod(testPod, testPod.Spec.Containers[0].Name, ncCmd)
ExpectWithOffset(1, err).ToNot(HaveOccurred())
ExpectWithOffset(1, output).To(ContainSubstring("succeeded"))
```

---

## Critical Validation Rules

**ALWAYS validate:**
1. ✅ `WaitUntilVMIReady` includes login function parameter
   - ❌ WRONG: `libwait.WaitUntilVMIReady(vmi)`
   - ✅ CORRECT: `libwait.WaitUntilVMIReady(vmi, console.LoginToFedora)`

2. ✅ All imports are correct and minimal
3. ✅ Package declaration is `package network`
4. ✅ Test structure follows Ginkgo v2 patterns

---

## Success Criteria

Code generation succeeds when:
- ✅ **All STD scenarios processed (1:1 mapping: scenario → test case)**
- ✅ **Count of `It()` blocks equals count of STD Tier 1 scenarios**
- ✅ **Every scenario ID appears in generated code with `[test_id:TS-XXX]` label**
- ✅ Valid Go files generated (proper syntax)
- ✅ All imports correct and minimal
- ✅ WaitUntilVMIReady always has login function
- ✅ Files saved to `outputs/go-tests/{JIRA_ID}/`
- ✅ Summary report generated with coverage validation

---

## Error Handling

**If STD file not found:**
- Error: "STD file not found at outputs/std/{JIRA_ID}/{JIRA_ID}_test_description.yaml"
- Suggestion: "Run `/std-builder {JIRA_ID}` first"
- Exit

**If scenario pattern not recognized:**
- Warning: "Could not detect pattern for scenario {id}"
- Fallback: Use basic_vmi_test.go.template
- Continue with next scenario

**If validation fails:**
- Error: "Generated code has validation errors"
- Action: Save to `.go.invalid` file for review
- Show validation error details
- Continue with remaining scenarios

---

## Files Structure

Pattern rules, templates, and reference tests are loaded from project config:

```
{project_context.config_dir}/
├── patterns/
│   └── tier1_patterns.yaml                       # Pattern detection rules
├── templates/
│   └── tier1/
│       ├── network_connectivity_test.go.template  # Network connectivity test
│       ├── basic_vmi_test.go.template             # Basic VMI lifecycle test
│       ├── parametric_ipv4_ipv6_test.go.template  # Dual stack parametrized test
│       └── migration_test.go.template              # Migration with connectivity
└── reference/
    └── tier1/                                      # Reference test implementations
```

---

## Example: Multiple Scenarios in One File

**STD Input (4 scenarios):**
```yaml
scenarios:
  - test_id: "TS-CNV66855-001"
    tier: "Tier 1"
    test_objective: "Verify ICMP connectivity same node same subnet"

  - test_id: "TS-CNV66855-002"
    tier: "Tier 1"
    test_objective: "Verify TCP connectivity via external router"

  - test_id: "TS-CNV66855-005"
    tier: "Tier 1"
    test_objective: "Validate connectivity from worker node to VM"

  - test_id: "TS-CNV66855-006"
    tier: "Tier 1"
    test_objective: "Validate VM-initiated connections to pods"
```

**Generated Go File (ALL 4 scenarios included):**
```go
package network

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	// ... other imports
)

var _ = Describe("[CNV-66855] Localnet connectivity", decorators.SigNetwork, func() {
	var ctx context.Context
	var namespace string
	// ... shared variables

	BeforeEach(func() {
		ctx = context.Background()
		namespace = testsuite.GetTestNamespace(nil)
	})

	Context("Same-node connectivity tests", func() {
		// ✅ Scenario 1: MUST be included
		It("[test_id:TS-CNV66855-001] should allow ICMP from pod to VM same subnet", func() {
			// Test implementation for TS-CNV66855-001
		})

		// ✅ Scenario 2: MUST be included
		It("[test_id:TS-CNV66855-002] should allow TCP via external router", func() {
			// Test implementation for TS-CNV66855-002
		})

		// ✅ Scenario 5: MUST be included (was previously missing!)
		It("[test_id:TS-CNV66855-005] should allow connectivity from worker node to VM", func() {
			// Test implementation for TS-CNV66855-005
		})

		// ✅ Scenario 6: MUST be included (was previously missing!)
		It("[test_id:TS-CNV66855-006] should allow VM-initiated connections to pods", func() {
			// Test implementation for TS-CNV66855-006
		})
	})
})
```

**Validation:**
```bash
# Count scenarios in STD
yq '.scenarios[] | select(.tier == "Tier 1") | .test_id' CNV-66855_test_description.yaml | wc -l
# Output: 4

# Count test cases in generated code
grep -c '\[test_id:TS-CNV66855-' localnet_connectivity_test.go
# Output: 4

# ✅ PASS: 4 scenarios → 4 test cases (100% coverage)
```

**Summary Report:**
```
✅ Tier1 Go Test Generation Complete!

📊 Coverage Validation:
- STD Tier 1 scenarios: 4
- Generated test cases: 4
- Coverage: 100% ✅

📁 Generated Files:
- outputs/go-tests/CNV-66855/localnet_connectivity_test.go (4 test cases, 287 lines)

✅ All scenarios covered - validation passed
```

---

**End of Go Test Generator Skill**
