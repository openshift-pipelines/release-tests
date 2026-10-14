---
name: std-generator
description: Generate comprehensive v2.1-ENHANCED STD YAML with pattern metadata, variables, test structure from ALL STP scenarios (single file)
model: claude-opus-4-6
---

# STD Generator Skill (v2.1-enhanced)

## Purpose

Transforms **all scenarios** from a Software Test Plan (STP) into **ONE comprehensive v2.1-ENHANCED** Software Test Description (STD) YAML file with:
- Shared metadata and common preconditions
- **code_generation_config** (NEW in v2.1): imports, context init, timeout mappings
- **variables section per scenario** (NEW in v2.1): closure-scoped variable declarations
- **test_structure section per scenario** (NEW in v2.1): decorator placement, SIG() wrapper
- Detailed specifications for each scenario
- **Pattern metadata** (patterns, helpers, decorators, code templates)
- **Fixed code templates** (v2.1): no variable shadowing, ExpectWithOffset, auto-generated cleanups
- **Production-ready** for code generation (compiles without errors)

**Key Features:**
- Generates ONE file for ALL scenarios (not one file per scenario)
- Automatically adds pattern metadata to all scenarios
- Infers helper libraries from matched patterns
- Generates code templates from pattern library
- Ready for downstream code generation

## Input Required

- `scenarios`: Array of ALL scenario rows from STP Section III
  - Each scenario has:
    - `scenario_id`: Scenario number (e.g., 1, 2, 3)
    - `tier`: Tier classification (e.g., "Tier 1", "Tier 2")
    - `priority`: Priority (e.g., "P0", "P1", "P2")
    - `description`: Scenario description text
    - `requirement_id`: Requirement ID (e.g., "OCPBUGS-59657")
- `stp_context`: Context from the STP document
  - `jira_issue`: Jira ticket ID and metadata
  - `feature_description`: Feature overview (from Feature Overview section)
  - `related_prs`: List of GitHub PRs (from Metadata)
  - `api_endpoints`: API endpoints (from Section I.3 API Extensions, if applicable)
  - `known_limitations`: Known limitations (from Section I.2)
  - `test_environment`: Test environment requirements (from Section II.3)
- `stp_file_path`: Path to source STP file (e.g., `outputs/stp/CNV-66855/CNV-66855_test_plan.md`)

## Output

**Single comprehensive STD YAML file:**
- Filename: `{JIRA_ID}_test_description.yaml`
- Example: `CNV-66855_test_description.yaml`
- Location: `outputs/std/{JIRA_ID}/{JIRA_ID}_test_description.yaml`
- Size: Variable (~100-200 lines per scenario + 100 lines shared metadata)
- Format: Valid YAML with document metadata + scenarios array

**Structure:**
```yaml
---
# Document Metadata (shared)
document_metadata: {...}
common_preconditions: {...}

# Scenarios Array (one entry per STP scenario)
scenarios:
  - scenario_001: {...}
  - scenario_002: {...}
  - scenario_003: {...}
  ...
---
```

---

## STD Structure (2 Main Sections + Scenarios Array)

### Section 1: document_metadata

**Purpose:** Shared metadata for the entire test suite

**Required fields:**
```yaml
document_metadata:
  std_version: "2.1-enhanced"
  generated_date: "YYYY-MM-DD"
  jira_issue: "{JIRA_ID}"
  jira_summary: "{Jira issue summary}"
  source_bugs: ["{OCPBUGS-XXXXX}", ...]  # If applicable
  stp_reference:
    file: "outputs/stp/{JIRA_ID}/{JIRA_ID}_test_plan.md"
    version: "v1"
    sections_covered: "Section III - Requirements-to-Tests Mapping"

  # related_prs is internal metadata for code generation context.
  # It MUST NOT be propagated to Phase 1 stub module docstrings.
  # Stub docstrings contain only STP Reference and Jira ID.
  related_prs:
    - repo: "{org/repo}"
      pr_number: {number}
      url: "{PR_URL}"
      title: "{PR title}"
      merged: true

  owning_sig: "{sig-name}"
  participating_sigs: ["{sig-1}", "{sig-2}"]

  total_scenarios: {count}
  tier_1_count: {count}
  tier_2_count: {count}
  p0_count: {count}
  p1_count: {count}
```

**Derivation:**
- Extract from STP metadata table (Section I)
- Count scenarios by tier and priority
- List all related PRs from STP Section II.4

---

### Section 1.5: code_generation_config (NEW IN v2.1)

**Purpose:** Code generation configuration for downstream test file generation

**Note:** The imports, helper_library_imports, and timeout_constants below should be read from `{project_context.config_dir}/tier1.yaml` instead of being hardcoded. The following serves as documentation/example for the CNV project.

```yaml
code_generation_config:
  std_version: "2.1-enhanced"
  framework: "ginkgo-v2"
  assertion_library: "gomega"
  language: "go"
  package_name: "{INFER_FROM_SIG}"  # sig-network → "network", sig-compute → "compute"

  # Context initialization (injected at BeforeAll start)
  context_init:
    - statement: "ctx := context.Background()"
      variable: "ctx"
      type: "context.Context"
    - statement: "namespace := testsuite.GetTestNamespace(nil)"
      variable: "namespace"
      type: "string"

  # Global imports (read from {project_context.config_dir}/tier1.yaml)
  imports:
    dot_imports:
      - "github.com/onsi/ginkgo/v2"
      - "github.com/onsi/gomega"
    standard:
      - "context"
      - "time"
    k8s_core:
      - path: "k8s.io/api/core/v1"
        alias: "k8sv1"
      - path: "k8s.io/apimachinery/pkg/apis/meta/v1"
        alias: "metav1"
    kubevirt_base:
      - "kubevirt.io/kubevirt/tests/decorators"
      - "kubevirt.io/kubevirt/tests/framework/kubevirt"
      - "kubevirt.io/kubevirt/tests/testsuite"
      - "kubevirt.io/kubevirt/tests/libvmi"
    kubevirt_api:
      - path: "kubevirt.io/api/core/v1"
        alias: "v1"
    network:
      - path: "github.com/k8snetworkplumbingwg/network-attachment-definition-client/pkg/apis/k8s.cni.cncf.io/v1"
        alias: "networkv1"

  # Timeout constants mapping (read from {project_context.config_dir}/tier1.yaml)
  timeout_constants:
    tiny: "StartupTimeoutSecondsTiny"       # 30s
    small: "StartupTimeoutSecondsSmall"     # 60s
    medium: "StartupTimeoutSecondsMedium"   # 90s
    large: "StartupTimeoutSecondsLarge"     # 120s
    xlarge: "StartupTimeoutSecondsXLarge"   # 180s
    huge: "StartupTimeoutSecondsHuge"       # 240s
    xhuge: "StartupTimeoutSecondsXHuge"     # 300s

  migration_timeout: "MigrationWaitTime"  # 240s

  # Helper library import mappings (read from {project_context.config_dir}/tier1.yaml)
  helper_library_imports:
    libvmifact: "kubevirt.io/kubevirt/tests/libvmifact"
    libnet: "kubevirt.io/kubevirt/tests/libnet"
    libwait: "kubevirt.io/kubevirt/tests/libwait"
    libpod: "kubevirt.io/kubevirt/tests/libpod"
    libvmops: "kubevirt.io/kubevirt/tests/libvmops"
    libmigration: "kubevirt.io/kubevirt/tests/libmigration"
    libstorage: "kubevirt.io/kubevirt/tests/libstorage"
    console: "kubevirt.io/kubevirt/tests/console"
    matcher: "kubevirt.io/kubevirt/tests/framework/matcher"
```

**Derivation:**
- **package_name**: Infer from `owning_sig`:
  - sig-network → "network"
  - sig-compute → "compute"
  - sig-storage → "storage"
  - sig-migration → "migration"
  - default → "tests"
- All other fields are STATIC (copy from template above)

---

### Section 2: common_preconditions

**Purpose:** Infrastructure and environment requirements shared by ALL scenarios

**Required fields:**
```yaml
common_preconditions:
  infrastructure:
    - name: "OpenShift cluster"
      requirement: "OCP {version}+ with OVN-Kubernetes"
      validation: "oc version"

    - name: "OpenShift Virtualization"
      requirement: "CNV {version}+"
      validation: "oc get csv -n openshift-cnv | grep kubevirt-hyperconverged"

    - name: "{Additional infrastructure}"
      requirement: "{From STP Section II.5}"
      validation: "{Validation command}"

  operators:
    - name: "{Operator name}"
      namespace: "{namespace}"
      validation: "{oc get csv command}"

  cluster_configuration:
    topology: "{Single-node|Multi-node}"
    cpu_virtualization: "{Standard|Nested}"
    storage: "{StorageClass requirement}"
    network: "{CNI requirement}"

  rbac_requirements:
    - permission: "{verb} on {resource}"
      scope: "{Cluster|Namespace: {namespace}}"
      validation: "oc auth can-i {verb} {resource}"
```

**Derivation:**
- Extract from STP Section II.5 (Test Environment)
- Extract from STP Section I.2 (Technology and Design Review)
- Infer RBAC from feature type (API operations, resource management)

---

### Section 3: scenarios

**Purpose:** Array of detailed scenario specifications

**Structure:** One entry per STP scenario

**Required fields for each scenario:**
```yaml
scenarios:
  - scenario_id: "{NUM}"
    test_id: "TS-{JIRA_ID}-{NUM:03d}"
    tier: "{Tier 1|Tier 2}"
    priority: "{P0|P1|P2}"
    mvp: {true|false}
    requirement_id: "{REQUIREMENT_ID}"

    # ===== PATTERN METADATA (AUTO-GENERATED) =====
    patterns:
      primary: "{matched_primary_pattern}"
      secondary:
        - "{matched_setup_pattern_1}"
        - "{matched_setup_pattern_2}"
        - "{matched_execution_pattern_1}"
      helpers_required:
        - name: "{helper_library_name}"
          functions: ["{function1}", "{function2}"]
          purpose: "{what_it_does}"
      decorators:
        - "{decorator_1}"
        - "{decorator_2}"

    # ===== VARIABLE DECLARATIONS (AUTO-GENERATED in v2.1) =====
    variables:
      closure_scope:
        - name: "{variable_name}"
          type: "{Go_type}"
          initialized_in: "{BeforeAll|It}"
          used_in: ["{BeforeAll}", "{It}", "{AfterEach}"]
          comment: "{Brief description}"
    # =========================================================

    # ===== TEST STRUCTURE (AUTO-GENERATED in v2.1) =====
    test_structure:
      type: "{single|table-driven}"

      describe:
        wrapper: "SIG"
        description: "{Feature description}"
        decorators:
          - "{SIG_decorator}"
          - "Serial"

      context:
        description: "{Scenario description}"
        decorators:
          - "Ordered"
          - "decorators.OncePerOrderedCleanup"

      it:
        description: "should {test_objective}"
        test_id_format: "[test_id:{test_id}]"
    # ===================================================

    code_structure: |
      Context("{scenario_description}", Ordered) {
        BeforeAll(func() {
          // Setup
        })
        It("[test_id:{test_id}]should {test_objective}", func() {
          // Test
        })
      }
    # =============================================

    test_objective:
      title: "{scenario.description}"
      what: |
        {Expand scenario description into 2-3 sentences explaining:
         - What functionality is being tested
         - What specific aspect/behavior is validated
         - What operations are performed}

      why: |
        {Explain business/technical rationale:
         - Why this test is important
         - What user need it addresses
         - What could break if this fails}

      acceptance_criteria:
        - "{Criterion 1: clear, measurable condition}"
        - "{Criterion 2: ...}"

    classification:
      test_type: "{Functional|Integration|E2E}"
      scope: "{Single-component|Multi-component}"
      automation_approach: "pytest with openshift-python-wrapper"

    specific_preconditions:
      # Scenario-specific requirements (beyond common_preconditions)
      - name: "{Specific requirement}"
        requirement: "{Details}"
        validation: "{Command}"

    test_data:
      # YAML definitions for this specific scenario
      resource_definitions:
        - name: "{resource_name}"
          type: "{VirtualMachine|Pod|NetworkAttachmentDefinition|etc}"
          yaml: |
            {Complete YAML definition}

      api_endpoints:
        # If applicable
        - operation: "{operation_name}"
          method: "{GET|POST|PUT|DELETE}"
          path: "{API path}"
          expected_status: {200|201|etc}

    test_steps:
      setup:
        - step_id: "SETUP-01"
          action: "{Setup action}"
          command: "{Command or API call}"
          validation: "{Expected result}"
          pattern_id: "{matched_pattern}"        # AUTO-ADDED
          code_template: |                       # AUTO-ADDED
            {code from pattern library}

      test_execution:
        - step_id: "TEST-01"
          action: "{Test action}"
          command: "{Command or API call}"
          validation: "{Expected result}"
          pattern_id: "{matched_pattern}"        # AUTO-ADDED
          code_template: |                       # AUTO-ADDED
            {code from pattern library}

      cleanup:
        - step_id: "CLEANUP-01"
          action: "{Cleanup action}"
          command: "{Command}"

    assertions:
      - assertion_id: "ASSERT-01"
        priority: "P0"
        description: "{What is being validated}"
        condition: "{Expected condition}"
        failure_impact: "{What failure means}"

    dependencies:
      kubernetes_resources:
        - "{Resource type}: {name}"

      external_tools:
        - "{Tool name} {version}+"

      scenario_specific_rbac:
        - "{permission description}"

```

**Derivation:**
- `test_objective.title`: Use scenario.description verbatim
- `test_objective.what`: Expand description with specifics
- `test_objective.why`: Infer from STP Section I.1 (Requirement Review)
- `acceptance_criteria`: Extract from scenario description and STP acceptance criteria
- `classification`: Infer from tier and scenario complexity
- `specific_preconditions`: Add scenario-specific requirements (e.g., external router for networking tests)
- `test_data`: Generate realistic YAML for VMs, pods, networks based on scenario
- `test_steps`: Expand scenario into 5-10 detailed steps (setup → execute → cleanup)
- `assertions`: Extract validation points from scenario description (2-5 per scenario)
- `dependencies`: List K8s resources, tools, and RBAC specific to this scenario

---

## PATTERN ENHANCEMENT (AUTO-GENERATION)

**CRITICAL:** All scenarios MUST include pattern metadata for production-ready STD

For each scenario, analyze the description and automatically add pattern metadata using the rules below.

### Pattern Matching Rules

Apply these rules to match scenarios to patterns from `{project_context.config_dir}/patterns/` directory:

#### 1. Keywords → Primary Pattern

Analyze the scenario description for these keywords:

- Contains **"connectivity"**, **"ping"**, **"reach"** → `network-connectivity-001`
- Contains **"migration"**, **"migrate"** → `migration-001`
- Contains **"hotplug"**, **"attach"** → `network-hotplug-001`
- Contains **"lifecycle"**, **"create"**, **"delete"** → `vm-lifecycle-001`
- Contains **"console"**, **"login"**, **"SSH"** → `console-001`
- Contains **"snapshot"**, **"restore"** → `snapshot-001`
- Contains **"clone"**, **"copy"** → `clone-001`

#### 2. Resources → Setup Patterns

Identify required resources and add setup patterns:

- Mentions **"VM"**, **"VMI"**, **"VirtualMachineInstance"** → Add `factory-001`
- Mentions **"Pod"** → Add `factory-pod-001`
- Mentions **"NetworkAttachmentDefinition"**, **"NAD"** → Add `network-nad-001`
- **Any VMI creation** → Also add `wait-002` (always wait for VMI ready)

#### 3. Actions → Execution Patterns

Identify test actions and add execution patterns:

- Action: **"ping"**, **"connectivity test"** → Add `network-ping-001`
- Action: **"migrate"** → Add `migration-execute-001`
- Action: **"console"**, **"login"** → Add `console-001`
- Action: **"hotplug"**, **"attach"** → Add `network-hotplug-execute-001`

#### 4. Infer Helpers from Patterns

Based on matched patterns, automatically infer required helper libraries:

**Pattern → Helper Mapping:**
- `network-connectivity-001` requires: **libvmifact**, **libnet**, **libwait**
- `factory-001` requires: **libvmifact**, **libvmi**
- `migration-001` requires: **libvmifact**, **libmigration**
- `console-001` requires: **console**
- `factory-pod-001` requires: **libpod**
- `network-nad-001` requires: **libnet**
- `wait-002` requires: **libwait**

**Helper Library Functions (Common):**
- **libvmifact**: `NewFedora`, `NewAlpineWithTestTooling`, `NewCirros`, `NewAlpine`
- **libvmi**: `WithInterface`, `WithNetwork`, `WithMasqueradeNetworking`
- **libnet**: `PingFromVMConsole`, `GetVmiPrimaryIPByFamily`, `CreateNetworkAttachmentDefinition`
- **libwait**: `WaitUntilVMIReady`, `WaitForVMIPhase`
- **console**: `LoginToFedora`, `LoginToAlpine`, `RunCommand`, `SafeExpectBatch`
- **libmigration**: `MigrateVMI`, `ConfirmVMIPostMigration`
- **libpod**: `RenderPod`, `CreatePodFromDefinition`

#### 5. Add Decorators

Add test decorators based on tier and domain:

**Tier-based:**
- Tier 1 → `decorators.Tier1`
- Tier 2 → `decorators.Tier2`

**Domain-based (from scenario description):**
- Network-related → `decorators.SigNetwork`
- Migration-related → `decorators.SigCompute`
- Storage-related → `decorators.SigStorage`
- Compute-related → `decorators.SigCompute`

**Always add:**
- `Ordered` (for proper test execution order)
- `decorators.OncePerOrderedCleanup` (for cleanup after ordered tests)

#### 6. Generate Code Templates

For each matched pattern:

1. **Read pattern definition** from `{project_context.config_dir}/patterns/tier1_patterns.yaml`
2. **Extract the `template` field** for that pattern
3. **Add as `code_template`** to the corresponding test step
4. **Add `pattern_id`** to link step to pattern

**Example:**
```yaml
test_steps:
  setup:
    - step_id: "SETUP-01"
      action: "Create Fedora VMI"
      pattern_id: "factory-001"           # Added
      code_template: |                    # Added from pattern library
        vmi := libvmifact.NewFedora(
            libvmi.WithInterface(iface),
            libvmi.WithNetwork(network),
        )
```

#### 7. Generate Code Structure

For each scenario, generate a Ginkgo test structure hint:

```go
Context("{scenario_description}", Ordered) {
  BeforeAll(func() {
    // Setup from test_steps.setup
  })
  It("[test_id:{test_id}]should {test_objective}", func() {
    // Test execution from test_steps.test_execution
  })
}
```

Replace placeholders:
- `{scenario_description}`: Brief description of scenario
- `{test_id}`: The test_id field (e.g., TS-CNV66855-001)
- `{test_objective}`: The test_objective.title field

---

### Pattern Library Reference

**Location**: `{project_context.config_dir}/patterns/tier1_patterns.yaml`

**Available Patterns:**
- `network-connectivity-001` - Network connectivity tests
- `factory-001` - VMI creation with factory
- `wait-002` - Wait for VMI ready
- `console-001` - Console login
- `migration-001` - Live migration
- `network-nad-001` - NetworkAttachmentDefinition creation
- `factory-pod-001` - Pod creation
- `table-001` - Table-driven tests
- `network-hotplug-001` - Network interface hotplug
- `snapshot-001` - VM snapshot operations
- `clone-001` - VM cloning operations

Each pattern provides:
- **keywords**: Trigger words for matching
- **resources**: Applicable K8s resources
- **actions**: Test actions
- **helpers**: Required helper libraries
- **template**: Ready-to-use code template

---

### Pattern Enhancement Validation

Before generating the STD file, validate pattern metadata:

- [ ] All scenarios have `patterns.primary` field
- [ ] All scenarios have `patterns.helpers_required` array
- [ ] All scenarios have `patterns.decorators` array
- [ ] All scenarios have `code_structure` field
- [ ] All test steps have `pattern_id` where applicable
- [ ] All test steps have `code_template` where applicable
- [ ] Pattern IDs reference actual patterns in `{project_context.config_dir}/patterns/`

---

## LLM Prompt for Comprehensive STD Generation

**System Prompt:**
```
You are an expert QE engineer generating a comprehensive Software Test Description (STD) from a Software Test Plan (STP).

Your task:
1. Read ALL scenarios from the STP Section III (Requirements-to-Tests Mapping table)
2. Extract shared metadata from STP Sections I and II
3. Generate ONE comprehensive STD YAML file with:
   - document_metadata (shared across all scenarios)
   - common_preconditions (shared infrastructure/environment)
   - scenarios array (detailed spec for each scenario)

Guidelines:
- Generate ONE file for ALL scenarios (not one file per scenario)
- Extract common preconditions to avoid duplication
- Be specific and detailed in scenario specifications
- Use realistic KubeVirt/OpenShift patterns
- Include complete YAML definitions for test resources
- Link scenarios to requirements (Jira, GitHub PRs)
- Prioritize assertions (P0 = critical, P1 = nice to have)

CRITICAL - Pattern Enhancement (AUTO-GENERATED):
- For EACH scenario, analyze the description and automatically add pattern metadata
- Apply pattern matching rules (keywords → patterns, resources → setup patterns, etc.)
- Infer helper libraries from matched patterns
- Add decorators based on tier and domain
- Generate code templates from pattern library (`{project_context.config_dir}/patterns/tier1_patterns.yaml`)
- Add code_structure hint for each scenario
- Add pattern_id and code_template to each test step
- This is NOT optional - ALL scenarios MUST have pattern metadata

Output only valid YAML. Do not include explanations outside the YAML structure.
```

**User Prompt Template:**
```
Generate a comprehensive Software Test Description (STD) YAML file for ALL scenarios in the following STP:

STP FILE: {stp_file_path}

DOCUMENT METADATA:
  Jira Issue: {jira_id} - {jira_summary}
  Source Bugs: {source_bugs}
  Related PRs: {related_prs}
  Owning SIG: {owning_sig}
  Total Scenarios: {total_scenarios}

STP CONTEXT:
  Feature Description: {feature_description}
  Non-Goals: {non_goals}
  Test Environment: {test_environment}
  Fix Versions: {fix_versions}

ALL SCENARIOS (from STP Section III):
{scenarios_array}

Generate ONE comprehensive STD YAML file with:
1. document_metadata (shared metadata for entire test suite)
2. common_preconditions (shared infrastructure/environment requirements)
3. scenarios (array with detailed spec for each scenario)

Output filename: {JIRA_ID}_test_description.yaml
Output only valid YAML.
```

---

## Final Validation Checklist

Before outputting the STD YAML, validate ALL of the following:

**Base STD Structure:**
- [ ] Valid YAML syntax (parse with YAML parser)
- [ ] document_metadata section complete
- [ ] document_metadata.std_version is "2.1-enhanced"
- [ ] common_preconditions section complete
- [ ] scenarios array has entries for ALL STP scenarios
- [ ] Each scenario has required fields:
  - [ ] scenario_id, test_id, tier, priority
  - [ ] test_objective (title, what, why, acceptance_criteria)
  - [ ] test_steps (setup, test_execution, cleanup)
  - [ ] assertions (at least 1 per scenario)
- [ ] No "TODO" or placeholder values
- [ ] All scenario test_ids follow format: TS-{JIRA_ID}-{NUM:03d}

**Pattern Enhancement:**
- [ ] ALL scenarios have `patterns` section with `primary` field
- [ ] ALL scenarios have `patterns.helpers_required` array
- [ ] ALL scenarios have `patterns.decorators` array
- [ ] ALL scenarios have `code_structure` field
- [ ] ALL test steps have `pattern_id` where applicable
- [ ] ALL test steps have `code_template` where applicable
- [ ] Pattern IDs match patterns in `{project_context.config_dir}/patterns/tier1_patterns.yaml`

**v2.1 Enhancement:**
- [ ] `code_generation_config` section exists at document level
- [ ] `code_generation_config.std_version` is "2.1-enhanced"
- [ ] `code_generation_config.package_name` is inferred from owning_sig
- [ ] ALL scenarios have `variables` section
- [ ] ALL scenarios have `test_structure` section
- [ ] ALL `variables.closure_scope` includes at minimum: ctx, namespace, err
- [ ] ALL `test_structure.context.decorators` includes: Ordered, decorators.OncePerOrderedCleanup
- [ ] ALL code_templates use `=` (not `:=`) for closure variables
- [ ] ALL `Expect(err)` calls use `ExpectWithOffset(1, err)`
- [ ] ALL scenarios with setup steps have corresponding cleanup templates

**If ANY validation fails:**
- Log error with scenario_id and specific failure
- Do NOT output incomplete STD
- Return error report to user

---

## Success Criteria

STD generation is successful when:
- ✅ Valid YAML file created
- ✅ All 3 main sections populated (document_metadata, common_preconditions, scenarios)
- ✅ Scenarios array has {total_scenarios} entries
- ✅ No duplicate scenario IDs
- ✅ File size appropriate (~150 lines per scenario + 100 lines metadata)
- ✅ Traceability complete (Jira, PRs, STP reference)

---

## Error Handling

- **If scenario description is vague:**
  - Log warning: "Scenario {num} lacks detail - generating best-effort spec"
  - Generate spec with inferred details
  - Mark for human review in validation report

- **If STP context is incomplete:**
  - Use defaults (e.g., if no PRs → empty related_prs array)
  - Generate STD with available information
  - Log warning about missing context

- **If YAML generation fails:**
  - Return error message with LLM output
  - Suggest manual review and correction
  - Provide partial output if possible

---

## Output Location

**Primary output:**
- `outputs/std/{JIRA_ID}/{JIRA_ID}_test_description.yaml`

**Example:**
- `outputs/std/CNV-66855/CNV-66855_test_description.yaml`

**Note:** This comprehensive STD YAML is the single source of truth for all test scenarios. It is used by downstream generators (go-test-generator, python-test-generator) to produce working test code.

---

## Usage Example

See `supplemental.md` (in this skill directory) for detailed input/output examples.

---

## v2.1 Enhancements

v2.1 adds auto-generated `code_generation_config`, `variables`, and `test_structure` sections,
plus code template transformations (variable shadowing fixes, ExpectWithOffset, auto-cleanup).

For detailed algorithms, inference rules, transformation pseudocode, and the v2.1 changelog,
see `supplemental.md` (in this skill directory).

---

**End of STD Generator Skill**
