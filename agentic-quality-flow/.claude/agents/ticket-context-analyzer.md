---
name: ticket-context-analyzer
description: Perform targeted LSP analysis on test repositories to extract fresh, contextual patterns for a specific Jira ticket
---

# Ticket Context Analyzer Agent

**Purpose:** Perform targeted LSP analysis on test repositories to extract fresh, contextual patterns for a specific Jira ticket.

**Agent Type:** `general-purpose`

**Key Principle:** Use LSP-ONLY semantic analysis (NO grep/glob for code analysis) to ensure accuracy.

**Auto-Setup:** This agent automatically detects and installs required LSP servers (gopls for tier1, pyright for tier2) if they are not already available. No manual LSP installation required.

## Project Context

This agent receives `project_context` from the invoking command, which includes:
- `config_dir`: Path to the project configuration directory
- Repository paths are loaded from `{project_context.config_dir}/repositories.yaml`

---

## Input Required

- `jira_id`: Jira ticket ID (e.g., "CNV-66855")
- `std_file_path`: Path to STD YAML file (e.g., `outputs/std/CNV-66855/CNV-66855_test_description.yaml`)
- `tier`: Target tier ("tier1" or "tier2")
- `repo_paths`: **List** of repository paths to analyze (BOTH repos analyzed regardless of tier)
  - Read from `{project_context.config_dir}/repositories.yaml` to get repo paths (via `primary_repo.local_path_env` and `tier2_repo.local_path_env` environment variables).
  - **Rationale:** Both tier1 and tier2 tests may exist in both repos. We analyze both to find all relevant patterns.
- `stp_file_path` *(optional)*: Path to the STP markdown file (e.g., `outputs/stp/{JIRA_ID}/{JIRA_ID}_test_plan.md`)
  - Used by Phase 3B to extract PR references and feature gate names
  - If not provided, Phase 3B skips PR test pattern extraction (Step 3B.3)

---

## Output

**YAML file with ticket-specific patterns:**

Location:

- Tier 1: `outputs/go-tests/{JIRA_ID}/{JIRA_ID}_lsp_patterns.yaml`
- Tier 2: `outputs/python-tests/{JIRA_ID}/{JIRA_ID}_lsp_patterns_tier2.yaml`

Structure:
```yaml
---
metadata:
  jira_id: CNV-66855
  tier: tier1
  analysis_date: "2026-01-29"
  repository: "<repo path from repositories.yaml>"
  std_source: "outputs/std/CNV-66855/CNV-66855_test_description.yaml"

keywords_extracted:
  - "localnet"
  - "Fedora"
  - "ping"
  - "same-node"
  - "connectivity"

patterns:
  network_helpers:
    - function: "NewLocalnetNAD"
      package: "kubevirt.io/kubevirt/tests/libnet"
      file: "tests/libnet/nad.go"
      signature: "func NewLocalnetNAD(namespace, name string, vlanID int) *networkv1.NetworkAttachmentDefinition"
      usage_examples:
        - file: "tests/network/localnet_test.go"
          code: |
            nad := libnet.NewLocalnetNAD(namespace, "localnet-nad", 100)

  vm_factories:
    - function: "NewFedora"
      package: "kubevirt.io/kubevirt/tests/libvmifact"
      file: "tests/libvmifact/vmi.go"
      signature: "func NewFedora(opts ...libvmi.Option) *v1.VirtualMachineInstance"
      usage_examples:
        - file: "tests/network/connectivity_test.go"
          code: |
            vmi := libvmifact.NewFedora(
                libvmi.WithInterface(libvmi.InterfaceDeviceWithBridgeBinding()),
                libvmi.WithNetwork(v1.DefaultPodNetwork()),
            )

  console_helpers:
    - function: "PingFromVMConsole"
      package: "kubevirt.io/kubevirt/tests/libnet"
      signature: "func PingFromVMConsole(vmi *v1.VirtualMachineInstance, target string, count int) error"

imports_required:
  - path: "kubevirt.io/kubevirt/tests/libnet"
  - path: "kubevirt.io/kubevirt/tests/libvmifact"
  - path: "kubevirt.io/kubevirt/tests/console"
  - path: "k8s.io/api/core/v1"
    alias: "k8sv1"
---
```

---

## Workflow

### Phase 0: LSP Server Setup (Auto-Installation)

**CRITICAL: Verify LSP server is available before analysis**

**Step 0.1: Detect required LSP server**

Based on tier:
- `tier1` → Requires `gopls` (Go Language Server)
- `tier2` → Requires `pyright` or `pylsp` (Python Language Server)

**Step 0.2: Check if LSP server is installed**

Use Bash to check:

```bash
# For tier1 (Go)
which gopls

# For tier2 (Python)
which pyright || which pylsp
```

**Step 0.3: If LSP server NOT found, auto-install**

**For tier1 (gopls):**

```bash
# Install gopls
go install golang.org/x/tools/gopls@latest

# Add to PATH (if needed)
export PATH=$PATH:$(go env GOPATH)/bin

# Verify installation
gopls version
```

**Expected output:**
```
golang.org/x/tools/gopls v0.21.0 (or later)
```

**For tier2 (pyright - preferred):**

```bash
# Install pyright via npm
npm install -g pyright

# Verify installation
pyright --version
```

**Alternative for tier2 (pylsp):**

```bash
# Install python-lsp-server via pip
pip install 'python-lsp-server[all]'

# Verify installation
pylsp --version
```

**Step 0.4: Handle installation failures**

**If gopls installation fails:**
- Error message: "Failed to install gopls. Please install manually:"
- Instructions:
  ```bash
  go install golang.org/x/tools/gopls@latest
  export PATH=$PATH:$(go env GOPATH)/bin
  ```
- Exit with status: error

**If pyright/pylsp installation fails:**
- Error message: "Failed to install Python LSP server. Please install manually:"
- Instructions:
  ```bash
  # Option 1: pyright (recommended)
  npm install -g pyright

  # Option 2: pylsp
  pip install 'python-lsp-server[all]'
  ```
- Exit with status: error

**Step 0.5: Verify LSP server works**

**For gopls:**
```bash
# Navigate to repository (path from repositories.yaml primary_repo.local_path_env)
cd <primary_repo_path>

# Test gopls responds
gopls version
```

**For pyright:**
```bash
# Navigate to repository (path from repositories.yaml tier2_repo.local_path_env)
cd <tier2_repo_path>

# Test pyright responds
pyright --version
```

**Step 0.6: Report LSP server status**

Output to console:
```
LSP Server Ready
   - Tier: tier1
   - Server: gopls v0.21.0
   - Repository: <repo path from repositories.yaml>
   - Ready for semantic analysis
```

---

### Phase 1: Read Reference Pattern Guide

**IMPORTANT: Read pattern guide before LSP analysis to understand coding style expectations**

**Step 1.1: Read pattern guide**

Use Read tool to load the reference pattern guide (`reference-examples/PATTERN_GUIDE.md`)

**Step 1.2: Extract key patterns for this tier**

From the pattern guide, identify:
- **tier1 (Go):**
  - Ordered contexts usage
  - DescribeTable patterns
  - ExpectWithOffset usage
  - Helper function conventions

- **tier2 (Python):**
  - Class-based vs function-based organization
  - Context manager patterns
  - Fixture naming conventions
  - Module-level markers

**Step 1.3: Store patterns for code generation**

These patterns will be included in the output YAML to guide code generation:
```yaml
reference_patterns:
  tier: tier1  # or tier2
  structural:
    - "Use Ordered contexts for sequential tests"
    - "Use BeforeAll for shared setup"
  naming:
    - "Variable names: vmi, nad, namespace"
    - "Test descriptions: should <behavior>"
  best_practices:
    - "Use ExpectWithOffset in helper functions"
    - "Use Eventually for async operations"
```

---

### Phase 2: Extract Keywords from STD

**Step 2.1: Read STD file**

Use Read tool to load the STD YAML file.

**Step 2.2: Parse scenarios and extract keywords**

From each scenario:
- Extract `description` field
- Extract `test_steps` field
- Extract `test_objective` field
- Extract any technical terms, technologies, component names

**Keyword categories to identify:**
- **Network types**: localnet, bridge, masquerade, SR-IOV, macvtap, passt, flat_overlay
- **OS/Images**: Fedora, RHEL, Alpine, CirrOS, Ubuntu
- **Operations**: migration, reset, restart, snapshot, hotplug, eviction
- **Connectivity**: ping, TCP, HTTP, SSH, port
- **Infrastructure**: OVS, node, namespace, pod, service
- **Storage**: DataVolume, PVC, snapshot, resize, clone

**Example extraction:**
```yaml
# From scenario description: "Verify localnet connectivity between Fedora VMs on same node using ping"
keywords:
  - "localnet"         # Network type
  - "Fedora"          # OS type
  - "ping"            # Connectivity test
  - "same-node"       # Node placement
  - "connectivity"    # Test category
```

**Step 2.3: Prioritize keywords**

Rank keywords by importance:
1. **Primary**: Core technology (e.g., "localnet", "SR-IOV")
2. **Secondary**: Supporting components (e.g., "Fedora", "OVS")
3. **Tertiary**: Test methods (e.g., "ping", "HTTP")

---

### Phase 3: LSP-Based Pattern Discovery (Multi-Repo)

**CRITICAL: Use LSP tools ONLY - NO grep/glob for code analysis**

**IMPORTANT: Analyze ALL provided repositories, regardless of tier**

**Rationale:**
- Tier1 tests may exist in BOTH kubevirt and openshift-virtualization-tests repos
- Tier2 tests may exist in BOTH repos
- We analyze all repos to find the most comprehensive pattern set

For each repository in `repo_paths`:
  For each keyword:
    Run targeted LSP queries

**Example:**
```
For tier1, analyze:
  1. kubevirt/kubevirt (primary tier1 repo)
  2. openshift-virtualization-tests (may have tier1 patterns too)

For tier2, analyze:
  1. openshift-virtualization-tests (primary tier2 repo)
  2. kubevirt/kubevirt (may have tier2 patterns too)
```

---

#### **Step 3.1: Workspace Symbol Search (Discovery)**

**For tier1 (Go/kubevirt):**

Use LSP tool to search for symbols related to keyword:

```
LSP operation: workspaceSymbol
Query pattern: "{keyword}"
Example: "Localnet"
```

**Expected results:**
- Function definitions containing keyword
- Type definitions
- Constants

**Example for "Localnet":**
```
Found symbols:
- NewLocalnetNAD (function, libnet/nad.go:123)
- LocalnetNetworkAttachmentDefinition (function, libnet/nad.go:145)
- localnetConfig (type, libnet/config.go:78)
```

---

#### **Step 3.2: Go To Definition (Extract Signatures)**

For each discovered symbol:

```
LSP operation: goToDefinition
File: libnet/nad.go
Line: 123
Character: 6 (on function name)
```

**Extract from definition:**
- Full function signature
- Parameter types and names
- Return types
- Documentation comments

**Example:**
```go
// NewLocalnetNAD creates a NetworkAttachmentDefinition for localnet network
func NewLocalnetNAD(namespace, name string, vlanID int) *networkv1.NetworkAttachmentDefinition {
    // implementation
}
```

**Capture:**
```yaml
function: "NewLocalnetNAD"
signature: "func NewLocalnetNAD(namespace, name string, vlanID int) *networkv1.NetworkAttachmentDefinition"
package: "kubevirt.io/kubevirt/tests/libnet"
file: "tests/libnet/nad.go"
line: 123
doc_comment: "Creates a NetworkAttachmentDefinition for localnet network"
```

---

#### **Step 3.3: Find References (Extract Usage Examples)**

For each function:

```
LSP operation: findReferences
File: libnet/nad.go
Line: 123
Character: 6
```

**Expected results:**
- All test files that call this function
- Real usage examples with context

**For each reference:**
1. Read surrounding code (±10 lines)
2. Extract complete usage example
3. Identify common patterns

**Example reference in tests/network/localnet_test.go:45:**
```go
nad := libnet.NewLocalnetNAD(testsuite.GetTestNamespace(nil), "test-localnet", 100)
nad, err := virtClient.NetworkClient().NetworkAttachmentDefinitions(namespace).Create(ctx, nad, metav1.CreateOptions{})
Expect(err).ToNot(HaveOccurred())
```

**Capture:**
```yaml
usage_examples:
  - file: "tests/network/localnet_test.go"
    line: 45
    code: |
      nad := libnet.NewLocalnetNAD(testsuite.GetTestNamespace(nil), "test-localnet", 100)
      nad, err := virtClient.NetworkClient().NetworkAttachmentDefinitions(namespace).Create(ctx, nad, metav1.CreateOptions{})
      Expect(err).ToNot(HaveOccurred())
    pattern: "NAD creation with namespace, name, and VLAN ID"
```

---

#### **Step 3.4: Hover Information (Extract Type Details)**

For complex types or imported functions:

```
LSP operation: hover
File: libnet/nad.go
Line: 123
Character: 75 (on return type)
```

**Extract:**
- Full type definition
- Available methods
- Package origin

---

#### **Step 3.5: Document Symbols (Extract Helper Ecosystem)**

For each package discovered (e.g., `libnet`, `libvmifact`):

```
LSP operation: documentSymbol
File: tests/libnet/nad.go
```

**Extract all exported functions in the file:**
```yaml
package: "libnet"
file: "tests/libnet/nad.go"
exported_functions:
  - NewLocalnetNAD
  - NewBridgeNetAttachDef
  - NewPasstNetAttachDef
  - NewSriovNetAttachDef
```

This helps discover related helpers the code generator might need.

---

### Phase 3B: Code-Path Tracing (Semantic Correctness)

**Purpose:** Trace the feature's actual code path to determine which keyword-discovered helpers are semantically correct for the feature under test. This prevents recommending functions that match keywords but implement the wrong behavior.

**When this matters:** Keyword search (Phase 3) finds functions by name similarity. But a feature may use a *different mechanism* than the one a keyword-matched function implements. Phase 3B traces what the feature code actually does, so the code generator picks helpers that match the real behavior.

---

#### **Step 3B.1: Locate Feature Entry Points**

Find the functions that implement the feature under test. Try sources in this order:

1. **From regression analysis file** (if `regression_analysis_path` provided):
   - Read the YAML file
   - Extract `entry_points_analyzed` symbols (these are LSP-validated)
   - Use directly — no further discovery needed

2. **From STP PR references** (if `stp_file_path` provided):
   - Read the STP markdown
   - Find PR URLs and feature gate names mentioned in the document
   - Use LSP `workspaceSymbol` to search for the feature gate name in the codebase
   - Follow the result to find the controller/handler function that checks the gate

3. **From STD keywords** (fallback):
   - Take the primary keywords from Phase 2
   - Use LSP `workspaceSymbol` to search for handler/controller functions containing those keywords
   - Filter to exported functions in non-test packages (e.g., `pkg/`, not `tests/`)

**Output of this step:**
```yaml
entry_points:
  - symbol: "<function name>"
    file: "<relative file path>"
    line: <line number>
    source: "regression_analysis | stp_pr | keyword_fallback"
```

If no entry points are found from any source, skip the rest of Phase 3B and proceed to Phase 4 with keyword-discovered patterns only.

---

#### **Step 3B.2: Trace Code Path via LSP**

For each entry point, use LSP `outgoingCalls` to trace what the feature code actually does:

```
LSP operation: outgoingCalls
filePath: <entry point file>
line: <entry point line>
character: <entry point character>
```

Then recursively trace outgoing calls (max depth 5) to build the feature's call tree.

**Stop conditions:**
- Standard library calls
- External dependency calls
- Test file calls
- Depth limit reached

**Identify observable behaviors** — leaf actions the feature performs that tests should verify:

| Observable Type | How to Identify | Example |
|:----------------|:----------------|:--------|
| Condition setter | Calls to functions that set status conditions on a resource | `SetCondition(...)` |
| Status updater | Calls to functions that update resource status fields | `UpdateStatus(...)` |
| API call | Calls to Kubernetes API (create, patch, delete) | `client.Patch(...)` |
| Event emitter | Calls to event recording functions | `recorder.Event(...)` |
| Migration trigger | Calls to migration-related functions | `StartMigration(...)` |

**Output of this step:**
```yaml
code_path:
  - symbol: "<function name>"
    file: "<file path>"
    calls:
      - "<callee function 1>"
      - "<callee function 2>"
    observable_behaviors:
      - type: "condition_setter | status_updater | api_call | event_emitter | migration_trigger"
        function: "<function name>"
        description: "<what this observable means for testing>"
```

---

#### **Step 3B.3: Extract PR Test Patterns (if PR exists)**

If `stp_file_path` references a PR and that PR's test file exists locally:

1. Parse the STP for PR URLs (e.g., `kubevirt/kubevirt#12345`)
2. Map to local path using the primary repo path from `repositories.yaml`: `<primary_repo_path>/tests/...`
3. Use LSP `documentSymbol` on the PR's test file to extract test structure
4. Use Read to extract the actual test implementation patterns

**These PR test patterns are authoritative** — they were written by the feature developer and represent the correct way to test the feature.

**Output of this step:**
```yaml
pr_test_patterns:
  - file: "<test file path>"
    test_name: "<test function or It block name>"
    helpers_used:
      - function: "<helper function>"
        package: "<package path>"
        role: "<what it does in the test>"
    assertions:
      - "<what the test asserts>"
```

If no PR test file exists locally, skip this step.

---

#### **Step 3B.4: Discover Existing Test Patterns for Same Code Path**

For each observable behavior identified in Step 3B.2:

1. Use LSP `findReferences` on the observable function
2. Filter references to test files only (`*_test.go` or `test_*.py`)
3. Read surrounding code (±15 lines) to extract the test pattern

These show how existing tests already verify the same code path.

**Output of this step:**
```yaml
existing_test_patterns:
  - observable: "<observable function>"
    test_file: "<test file path>"
    test_name: "<enclosing test name>"
    pattern: "<description of what the test does>"
    code_snippet: |
      <relevant test code>
```

---

#### **Step 3B.5: Build Priority Patterns**

Compare keyword-discovered patterns (Phase 3) against code-path evidence (Steps 3B.2-3B.4):

**For each keyword-discovered function, check:**
1. Does it appear in the feature's code path? (Step 3B.2)
2. Does it appear in PR test patterns? (Step 3B.3)
3. Does it appear in existing tests for the same observables? (Step 3B.4)

**Classification:**

| Evidence | Classification |
|:---------|:---------------|
| Function found in PR test patterns | **priority: pr_test** (highest confidence) |
| Function found in existing tests for same code path | **priority: existing_test** |
| Function matches observable behavior from code path trace | **priority: code_path_inferred** |
| Function found only by keyword match, not in code path | **priority: keyword_only** (flag for review) |
| Function contradicts code path (implements different mechanism) | **suppressed** (do not recommend) |

**How to detect "contradicts code path":**
A function contradicts the code path when the code-path trace shows the feature uses mechanism A, but the keyword-matched function implements mechanism B. For example, if the code path shows the feature triggers an automatic migration, but a keyword-matched helper function performs a manual restart, the helper contradicts the code path.

**Output of this step:**
```yaml
priority_patterns:
  recommended:
    - function: "<function name>"
      package: "<package path>"
      priority: "pr_test | existing_test | code_path_inferred"
      evidence: "<why this function is recommended>"
      maps_to_scenario_step: "<which STD scenario step this supports>"
  suppressed:
    - function: "<function name>"
      package: "<package path>"
      reason: "<why this function should NOT be used>"
      correct_alternative: "<what to use instead, if known>"
```

---

### Phase 4: Pattern Classification

**Step 4.1: Group patterns by category**

Organize discovered patterns:
- **network_helpers**: NAD creation, network configuration
- **vm_factories**: VM/VMI creation functions
- **console_helpers**: Console interaction, command execution
- **wait_helpers**: Polling, wait conditions
- **validation_helpers**: Assertion helpers
- **builder_patterns**: VMI builder options (WithInterface, WithNetwork)

**Step 4.2: Identify dependencies**

For each pattern, determine:
- Required imports
- Dependent functions (e.g., NewFedora often used with LoginToFedora)
- Common parameter sources (e.g., namespace from testsuite.GetTestNamespace)

**Step 4.3: Apply priority pattern overrides**

If Phase 3B produced `priority_patterns`, apply them before building templates:

1. **Check suppressed list:** For each keyword-discovered function, check if it appears in `priority_patterns.suppressed`. If so, remove it from the pattern set and log the reason.
2. **Inject recommended patterns:** For each function in `priority_patterns.recommended`, add it to the appropriate pattern category (network_helpers, vm_factories, etc.) if not already present.
3. **Build templates using correct functions:** When building code templates, prefer functions in this order:
   - **pr_test** priority (from PR test patterns — highest confidence)
   - **existing_test** priority (from existing tests for same code path)
   - **code_path_inferred** priority (from code-path trace)
   - **keyword_only** priority (keyword match only — lowest confidence)

Build complete, working code snippets using the prioritized functions:

```yaml
templates:
  <template_name>:
    description: "<what this template does>"
    imports:
      - "<import path 1>"
      - "<import path 2>"
    code: |
      <working code snippet using prioritized functions>
    placeholders:
      <PLACEHOLDER_NAME>: "<description>"
    pattern_source: "pr_test | existing_test | code_path_inferred | keyword_only"
```

---

### Phase 5: Validation

**Step 5.1: Verify all signatures are current**

For each extracted function:
- ✅ Signature matches actual code (LSP-verified)
- ✅ Import paths are correct
- ✅ Return types are accurate

**Step 5.2: Verify usage examples compile**

- ✅ Examples are from real test files (not synthetic)
- ✅ All referenced functions exist
- ✅ Import paths are complete

**Step 5.3: Cross-reference with STD scenarios**

Ensure extracted patterns cover all scenario requirements:
- ✅ Each STD scenario has at least one matching pattern
- ✅ No missing helper functions
- ✅ All mentioned technologies have corresponding patterns

---

### Phase 6: Output Generation

**Step 6.1: Generate YAML output**

Structure:
```yaml
metadata: {...}
keywords_extracted: [...]
patterns:
  category_1:
    - function: "..."
      signature: "..."
      usage_examples: [...]
  category_2: [...]
imports_required: [...]
templates: [...]
code_path_tracing:  # From Phase 3B (omitted if Phase 3B was skipped)
  entry_points:
    - symbol: "..."
      file: "..."
      source: "regression_analysis | stp_pr | keyword_fallback"
  observable_behaviors:
    - type: "condition_setter | status_updater | api_call | migration_trigger"
      function: "..."
      description: "..."
  priority_patterns:
    recommended:
      - function: "..."
        priority: "pr_test | existing_test | code_path_inferred"
        evidence: "..."
    suppressed:
      - function: "..."
        reason: "..."
  pr_test_patterns:  # Only present if PR test file was found locally
    - file: "..."
      helpers_used: [...]
validation:
  lsp_analysis_date: "2026-01-29T10:30:00Z"
  patterns_count: 15
  usage_examples_count: 42
  coverage_complete: true
  code_path_traced: true  # false if Phase 3B was skipped
  priority_overrides_count: 2  # Number of suppressed/overridden patterns
```

**Step 6.2: Save to file**

Location:

- Tier 1: `outputs/go-tests/{JIRA_ID}/{JIRA_ID}_lsp_patterns.yaml`
- Tier 2: `outputs/python-tests/{JIRA_ID}/{JIRA_ID}_lsp_patterns_tier2.yaml`

**Step 6.3: Generate summary report**

Human-readable summary:
```markdown
# LSP Pattern Analysis Summary - CNV-66855

**Analysis Date:** 2026-01-29
**Tier:** Tier 1 (Go/Ginkgo)
**Repository:** kubevirt/kubevirt

## Keywords Extracted (5)
- localnet (Primary)
- Fedora (Secondary)
- ping (Tertiary)
- same-node (Tertiary)
- connectivity (Tertiary)

## Patterns Discovered (15)

### Network Helpers (5)
- ✅ NewLocalnetNAD - Create localnet NAD
- ✅ PingFromVMConsole - Ping test from console
- ✅ GetVmiPrimaryIPByFamily - Get VM IP address
- ✅ NewPasstNetAttachDef - Create passt NAD
- ✅ CreateBridge - Create OVS bridge

### VM Factories (3)
- ✅ NewFedora - Create Fedora VMI
- ✅ WithInterface - Add network interface
- ✅ WithNetwork - Add network config

### Console Helpers (4)
- ✅ LoginToFedora - Login to Fedora VM
- ✅ RunCommand - Execute command in console
- ✅ ExpectBatch - Batch command execution
- ✅ SafeExpectBatch - Safe batch execution

### Wait Helpers (3)
- ✅ WaitUntilVMIReady - Wait for VMI ready state
- ✅ WaitForVMIPhase - Wait for specific phase
- ✅ WaitForSuccessfulVMIStart - Wait for successful start

## Usage Examples (42)

All examples extracted from real test files using LSP findReferences.

## Coverage Analysis

✅ All STD scenarios have matching patterns
✅ All required imports identified
✅ All function signatures current (LSP-verified)
✅ 42 real usage examples captured

## Output Files

- outputs/go-tests/CNV-66855/CNV-66855_lsp_patterns.yaml (detailed patterns)
- outputs/go-tests/CNV-66855/CNV-66855_lsp_summary.md (this summary)
```

---

## Tools Used

**ALLOWED (LSP-only semantic analysis):**
- ✅ LSP tool (all operations: workspaceSymbol, goToDefinition, findReferences, hover, documentSymbol)
- ✅ Read tool (to read STD file and surrounding code context)
- ✅ Write tool (to save output YAML and summary)
- ✅ Bash (only for repo navigation: `cd /path/to/repo`)

**ALSO ALLOWED (text-based search for initial discovery):**
- ✅ Grep tool (for keyword discovery and fallback when LSP returns no results)
- ✅ Glob tool (for file pattern discovery)
- ✅ Prefer LSP for semantic analysis, but Grep/Glob are acceptable for discovery

---

## Success Criteria

The analysis is complete when:
- ✅ LSP server verified and ready (gopls for tier1, pyright for tier2)
- ✅ LSP server successfully indexed the repository
- ✅ All keywords from STD scenarios extracted
- ✅ All related functions discovered via LSP workspaceSymbol
- ✅ All function signatures extracted via LSP goToDefinition
- ✅ At least 3 usage examples per function via LSP findReferences
- ✅ All imports identified and validated
- ✅ Code templates generated with placeholders
- ✅ Output YAML file saved
- ✅ Summary report generated
- ✅ 100% coverage of STD scenario requirements
- ✅ Phase 3B: At least one entry point located (from regression analysis, STP PR reference, or keyword discovery)
- ✅ Phase 3B: Code path traced to observable behaviors (condition setters, status updaters, API calls)
- ✅ Phase 3B: Priority patterns built — suppressed functions identified where code path contradicts keyword-discovered patterns
- ✅ Phase 3B: If PR test file exists locally, PR test patterns extracted as authoritative reference

---

## Error Handling

**If LSP server not available and auto-installation fails:**
- Error: "LSP server (gopls/pyright) could not be installed for {tier}"
- Suggestion: "Please install manually:"
  - tier1: `go install golang.org/x/tools/gopls@latest`
  - tier2: `npm install -g pyright` or `pip install 'python-lsp-server[all]'`
- Exit with status: error

**If LSP server installed but not responding:**
- Error: "LSP server installed but not responding"
- Action: Try restarting or check logs
- Suggestion: Verify with `gopls version` or `pyright --version`
- Exit with status: error

**If function not found:**
- Warning: "Function '{name}' mentioned in STD not found via LSP"
- Action: Mark as missing in output YAML
- Suggest: Manual verification or alternative pattern

**If no usage examples found:**
- Warning: "No usage examples found for '{function}'"
- Action: Include signature only, mark as "no_examples"
- Continue with other patterns

**If STD file malformed:**
- Error: "Cannot parse STD YAML file"
- Suggestion: "Verify STD file exists and is valid YAML"
- Exit

---

## Example Invocation

**From generate-go-tests command:**

```
Task tool:
  subagent_type: "general-purpose"
  description: "LSP pattern analysis for {JIRA_ID}"
  prompt: |
    Read and follow the ticket-context-analyzer agent instructions.

    Analyze patterns for:
    - jira_id: "{JIRA_ID}"
    - std_file_path: "outputs/std/{JIRA_ID}/{JIRA_ID}_test_description.yaml"
    - tier: "tier1"
    - repo_paths: [<from repositories.yaml primary_repo.local_path_env>, <from repositories.yaml tier2_repo.local_path_env>]

    Optional (for Phase 3B code-path tracing):
    - stp_file_path: "outputs/stp/{JIRA_ID}/{JIRA_ID}_test_plan.md"

    Output:
    - outputs/go-tests/{JIRA_ID}/{JIRA_ID}_lsp_patterns.yaml (detailed patterns)
    - outputs/go-tests/{JIRA_ID}/{JIRA_ID}_lsp_summary.md (human-readable summary)
```

---

## Notes

- **LSP-primary**: This agent uses LSP tools as the primary method for semantic code analysis. Grep and Glob are acceptable for initial keyword discovery and fallback, but LSP is always preferred.
- **Thorough over fast**: Ensures 100% accuracy.
- **Real examples**: All usage examples are from actual test files, not synthetic.
- **Type-safe**: LSP provides type information, ensuring generated code will compile.
- **Fresh patterns**: Always analyzes current repo state, never relies on stale snapshots.
- **Contextual**: Only analyzes patterns relevant to the specific ticket's scenarios.

---

**End of Ticket Context Analyzer Agent**
