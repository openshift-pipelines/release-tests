---
name: review-rules-extractor
description: Extract project-specific review rules dynamically from config files and optional repo scanning
---

# Review Rules Extractor Skill

**Phase:** Pre-Review (Step 3.5)
**User-Invocable:** false

## Purpose

Produces a complete `review_rules` data structure for the stp-reviewer and std-reviewer
skills by combining three sources in precedence order:

```
Hardcoded defaults (lowest) -> Dynamic extraction -> Static overrides (highest)
```

This enables Layer 2 (project-specific) review rules to activate automatically for any
project, without requiring manual creation of a `review_rules.yaml` file. Projects that
do have a static `review_rules.yaml` use it as an override layer on top of the dynamically
extracted rules.

## When to Use

Invoked by `/review-stp` and `/review-std` at **Step 3.5** to resolve review rules before
passing them to the reviewer skills. Replaces the previous direct read of `review_rules.yaml`.

## Tools Required

- Read
- Glob
- Grep

## Input

Expects `project_context` to already be in the conversation (from Step 0 project-resolver),
including:

```yaml
project_context:
  project_id: "cnv"
  config_dir: "config/projects/cnv"
  feature_toggles: { ... }
```

The skill receives the Jira ID as its argument (e.g., `CNV-72329`) but primarily uses
`project_context.config_dir` for file reads.

## Workflow

### Phase 0: Read Existing Config Files

Read the following files from `{config_dir}`. All reads are fast local file operations.
Missing files are not errors -- skip and note what was unavailable.

| File | Key data extracted |
|:-----|:-------------------|
| `project.yaml` | `decorator_mappings`, `scope_boundaries.cli_tools`, `domain_vocabulary` |
| `components.yaml` | `component_package_map` keys and feature names |
| `tier1.yaml` | `framework`, `helper_libraries`, `context_init`, `timeout_constants` |
| `tier2.yaml` | `framework` |
| `patterns/tier1_patterns.yaml` | `template_selection` rules (keyword-to-pattern mapping) |
| `repositories.yaml` | `primary_repo.local_path_env`, repo paths |
| `config/_defaults.yaml` | `test_id_format` |

Record which files were successfully loaded for the confidence annotation in the output.

### Phase 1: Build Review Rules from Config Data

Transform the config data into the `review_rules` structure. Each output key has a
defined source and transformation logic:

#### STP Rules

| Output key | Source | Logic |
|:-----------|:-------|:------|
| `stp_rules.abstraction.internal_components` | `components.yaml` keys | Extract component names from `component_package_map` keys: `["virt-controller", "virt-handler", ...]` |
| `stp_rules.testing_tools.standard_tools` | `project.yaml` cli_tools + tier1/tier2 frameworks | Combine: start with `cli_tools`, add framework display names. Example: `["virtctl", "oc", "kubectl", "Ginkgo", "pytest"]` |
| `stp_rules.testing_tools.standard_frameworks` | `tier1.yaml` framework + `tier2.yaml` framework | Transform to display names: `"ginkgo-v2"` becomes `"Ginkgo v2"`, `"pytest"` stays `"pytest"`. Also add known CI systems if identifiable from config (e.g., build_system `"bazel"` in repositories.yaml suggests Prow/OpenShift CI). |
| `stp_rules.upgrade.persistent_state_indicators` | `components.yaml` feature names | Scan feature names for CRD-like indicators. If features reference CRDs, configs, or stored state, include `["CRD", "stored config"]`. If VM lifecycle features exist, include `"running VM with feature-dependent data"`. |

#### STD Rules

| Output key | Source | Logic |
|:-----------|:-------|:------|
| `std_rules.patterns.sig_to_decorator` | `project.yaml` decorator_mappings | Strip `sig-` prefix from keys: `"sig-network": "decorators.SigNetwork"` becomes `network: "decorators.SigNetwork"` |
| `std_rules.patterns.closure_scope_required` | `tier1.yaml` context_init[].variable | Extract variable names: `["ctx", "namespace"]` |
| `std_rules.patterns.test_id_format` | `_defaults.yaml` test_id_format | Direct copy: `"TS-{JIRA_ID}-{NUM:03d}"` |
| `std_rules.patterns.ginkgo_structure` | `tier1.yaml` framework | If framework is `"ginkgo-v2"`: output `"Context -> BeforeAll -> It"`. Otherwise, derive from framework name. |
| `std_rules.patterns.keyword_to_pattern` | `patterns/tier1_patterns.yaml` template_selection | For each template_selection entry, extract keywords from `conditions[].match_any` and `conditions[].match_all`, map them to the entry `name` as the pattern ID. Example: `connectivity: "network-connectivity-001"` |
| `std_rules.patterns.pattern_to_helpers` | `tier1.yaml` helper_libraries + pattern keywords | Map pattern IDs to their required helper libraries based on keyword overlap. Example: patterns with "network" keywords need `["libvmifact", "libnet", "libwait"]` |
| `std_rules.timeouts` | `tier1.yaml` timeout_constants | Map operation types to timeout ranges. Example: `vm_startup: "medium to large"`, `api_call: "tiny to small"` |
| `std_rules.stub_conventions` | `tier1.yaml` + `tier2.yaml` framework | Derive from frameworks: if ginkgo-v2, `go_pending: "PendingIt()"`, `go_skip: "Skip()"`. If pytest, `python_pending: "pass"`, `python_test_disabled: "__test__ = False"`. If tier1 uses sig-based packages, `go_package_from_sig: true`. |

### Phase 2: Scan Repos if Locally Available (Optional)

Check `repositories.yaml` for `local_path_env` values. For each repo, check if the
environment variable is set and the path exists on disk.

**If primary repo is available:**

1. **Supplement `internal_components`:** Grep `{repo_path}/pkg/` for
   controller/operator/handler/reconciler/agent patterns in directory and file names.
   Add any component names not already in the list.

2. **Supplement `sig_to_decorator`:** Grep `{repo_path}/tests/**/*_test.go` for
   `decorators.Sig` patterns (e.g., `decorators.SigNetwork`, `decorators.SigCompute`).
   Extract the SIG suffix and add to the mapping if not already present.

3. **Confirm `ginkgo_structure`:** Grep `{repo_path}/tests/**/*_test.go` for `BeforeAll`
   vs `BeforeEach` usage counts. If `BeforeAll` is dominant, confirm
   `"Context -> BeforeAll -> It"`.

4. **Confirm `stub_conventions`:** Grep for `PendingIt` and `Skip(` patterns to validate
   the derived conventions.

**If tier2 repo is available:**

1. Grep `{repo_path}/tests/` for `pytest.mark.polarion` patterns to confirm Polarion
   marker format.

**If repos are not available:** Skip this phase entirely. Config-only extraction from
Phase 1 is sufficient for functional review rules.

Note which repos were scanned (if any) in the confidence annotation.

### Phase 3: Load Static Overrides

Attempt to read `{config_dir}/review_rules.yaml`.

**If found:**

Deep-merge the static file values over the dynamically extracted values:

- **Per-key override:** Static values replace dynamic values for the same key.
- **Array replacement:** Static arrays replace (not append to) dynamic arrays.
- **Missing keys preserved:** Keys present in the dynamic extraction but absent in the
  static file are preserved from the dynamic extraction.

Note in the confidence annotation that a static override file was applied.

**If not found:**

Continue with the dynamically extracted rules. This is the normal case for new projects.

### Phase 4: Fill Remaining Gaps with Hardcoded Defaults

After merging, check for any keys that are still empty/missing and fill them with
hardcoded defaults:

| Key | Default | Rationale |
|:----|:--------|:----------|
| `stp_rules.abstraction.internal_to_user_mappings` | `{}` | Requires domain knowledge; Rule A still flags generically |
| `stp_rules.abstraction.acceptable_locations` | `["Technology Challenges (I.3 sub-items)", "Risks (II.5)", "Checkbox sub-items", "Known Limitations (I.2)"]` | Same defaults as stp-reviewer SKILL.md |
| `stp_rules.dependencies.infrastructure_not_dependency` | `[]` | General Rule D heuristic still applies |
| `stp_rules.dependencies.dependency_examples` | `["Another team must merge a prerequisite PR", "External service must deploy new API version", "Platform team must release feature gate"]` | Generic examples applicable to any project |
| `stp_rules.strategy.always_y` | `["Functional Testing", "Automation Testing"]` | Universal for QualityFlow projects |
| `stp_rules.strategy.requires_justification_for_y` | `{"Performance Testing": "latency/throughput requirements", "Security Testing": "RBAC, auth, or security boundary changes", "Usability Testing": "UI component"}` | Standard criteria |
| `stp_rules.metadata.version_source` | `"fix_version"` | Standard Jira field |
| `stp_rules.scope.layered_product` | `null` | Null means skip layered product check |

## Output

The skill produces a complete `review_rules` data structure in the conversation context,
identical in format to the static `review_rules.yaml`. The reviewer skills
(stp-reviewer, std-reviewer) receive this structure unchanged.

### Output Format

```yaml
review_rules:
  project_id: "{project_id}"

  stp_rules:
    abstraction:
      internal_components: [...]
      internal_to_user_mappings: { ... }
      acceptable_locations: [...]
    dependencies:
      infrastructure_not_dependency: [...]
      dependency_examples: [...]
    upgrade:
      persistent_state_indicators: [...]
    testing_tools:
      standard_tools: [...]
      standard_frameworks: [...]
    strategy:
      always_y: [...]
      requires_justification_for_y: { ... }
    metadata:
      version_source: "fix_version"
    scope:
      layered_product: null | { name, platform, platform_teams, ownership_note }

  std_rules:
    patterns:
      keyword_to_pattern: { ... }
      pattern_to_helpers: { ... }
      sig_to_decorator: { ... }
      closure_scope_required: [...]
      test_id_format: "..."
      ginkgo_structure: "..."
    timeouts: { ... }
    stub_conventions:
      go_pending: "PendingIt()"
      go_skip: "Skip()"
      python_pending: "pass"
      python_test_disabled: "__test__ = False"
      go_package_from_sig: true/false

  _extraction_metadata:
    sources_used:
      config_files: ["project.yaml", "components.yaml", ...]
      repo_scans: ["kubevirt/kubevirt"] | []
      static_override: true/false
    keys_from_static: ["stp_rules.scope.layered_product", ...]
    keys_from_defaults: ["stp_rules.abstraction.internal_to_user_mappings", ...]
```

The `_extraction_metadata` section is informational only. It tells the reviewer skills
where each piece of data came from, which helps set review confidence:

- **All keys from config + repo scans:** HIGH confidence for project-specific rules
- **Most keys from config, some defaults:** MEDIUM confidence
- **Mostly defaults:** LOW confidence (equivalent to Layer 1 only)

## Error Handling

**Config file read failure (individual file):**
- Log warning: "Could not read {file}. Extraction will proceed without this data."
- Continue with remaining files. Do not fail the entire extraction.

**All config files fail:**
- Log warning: "No config files could be read. Returning hardcoded defaults only."
- Return Phase 4 defaults. The reviewer skills will operate in Layer 1 mode.

**Repo scan failure:**
- Log warning: "Repo scan failed for {repo}. Continuing with config-only extraction."
- Phase 2 results are always supplementary; their absence is not an error.

**Static file parse failure:**
- Log warning: "Could not parse review_rules.yaml. Using dynamic extraction only."
- Skip Phase 3. Continue with Phase 1 + Phase 4 results.

## Precedence Summary

```
Phase 4 defaults  <  Phase 1 config extraction  <  Phase 2 repo scan  <  Phase 3 static override
(lowest)                                                                  (highest)
```

Each phase only fills keys not already set by a higher-precedence phase. Phase 3 (static
override) always wins when present.
