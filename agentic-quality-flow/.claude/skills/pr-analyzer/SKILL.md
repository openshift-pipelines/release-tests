---
name: pr-analyzer
description: Analyze GitHub PR diffs and extract meaningful changes for STP generation
model: claude-opus-4-6
---

# PR Analyzer Skill

**Phase:** Pre-Processing
**User-Invocable:** false

## Purpose

Analyze GitHub PR diffs and extract meaningful changes for STP generation.

## When to Use

Invoked by the **github-pr-fetcher** subagent after fetching PR details, diffs, and review comments.

## Input

```yaml
pr_data:
  url: https://github.com/kubevirt/kubevirt/pull/1234
  owner: kubevirt
  repo: kubevirt
  pull_number: 1234
  title: Add CPU hot-plug support
  description: |
    This PR implements CPU hot-plug functionality...
  state: merged
  author: developer
  base_branch: main
  head_branch: feature/cpu-hotplug
  diff: |
    diff --git a/pkg/virt-controller/vm/vm.go b/pkg/virt-controller/vm/vm.go
    index abc123..def456 100644
    --- a/pkg/virt-controller/vm/vm.go
    +++ b/pkg/virt-controller/vm/vm.go
    @@ -100,6 +100,20 @@ func (c *VMController) Reconcile() {
    +func (c *VMController) HandleCPUHotplug(vm *v1.VirtualMachine) error {
    ...
  files:
    - filename: pkg/virt-controller/vm/vm.go
      status: modified
      additions: 50
      deletions: 10
    - filename: pkg/virt-controller/vm/hotplug.go
      status: added
      additions: 200
      deletions: 0
    - ...
  review_comments:
    - user: reviewer1
      body: "Consider edge case when VM is migrating"
      path: pkg/virt-controller/vm/hotplug.go
      line: 45
    - ...
```

## Output Format

```yaml
analysis:
  pr_url: https://github.com/kubevirt/kubevirt/pull/1234
  summary: Implements CPU hot-plug functionality for running VMs

  key_changes:
    functions:
      - name: HandleCPUHotplug
        file: pkg/virt-controller/vm/vm.go
        action: added
        purpose: Main entry point for CPU hot-plug operations
      - name: ValidateCPUChange
        file: pkg/virt-controller/vm/hotplug.go
        action: added
        purpose: Validates CPU changes before applying
      - ...

    types:
      - name: CPUHotplugSpec
        file: api/v1/types.go
        action: added
        fields_changed:
          - MaxSockets
          - CurrentSockets
      - ...

    apis:
      - endpoint: /virtualmachines/{name}/cpu
        method: PATCH
        action: added
        purpose: Hot-plug CPU to running VM
      - ...

    configurations:
      - name: EnableCPUHotplug
        type: feature_gate
        location: HyperConverged CR
        default: false
      - ...

  files_by_category:
    controllers:
      - pkg/virt-controller/vm/vm.go
      - pkg/virt-controller/vm/hotplug.go
    handlers:
      - pkg/virt-handler/hotplug/cpu.go
    api:
      - api/v1/types.go
      - api/v1/types_swagger_generated.go
    tests:
      - tests/hotplug_test.go
    other:
      - ...

  review_insights:
    edge_cases:
      - "VM migration during hot-plug needs handling"
      - "Consider maximum CPU limit validation"
    concerns:
      - "Performance impact of frequent hot-plug operations"
    suggestions:
      - "Add metrics for hot-plug success/failure rate"

  impact_assessment:
    components_affected:
      - virt-controller
      - virt-handler
      - virt-api
    features_potentially_impacted:
      - VM lifecycle
      - Live migration
      - Resource quotas
    breaking_changes: false
    api_changes: true
    config_changes: true
```

## Analysis Rules

### Function Detection

Parse diff for:
- `func (receiver) FunctionName(` - Go methods
- `func FunctionName(` - Go functions
- Added/Modified/Deleted based on diff markers (+/-)

### Type Detection

Parse diff for:
- `type TypeName struct {`
- `type TypeName interface {`
- Field additions/removals within structs

### API Detection

Look for:
- Route registrations (e.g., `router.Handle`, `http.HandleFunc`)
- OpenAPI/Swagger annotations
- CRD changes (in `api/` or `config/` directories)

### Configuration Detection

Look for:
- Feature gates
- Environment variables
- ConfigMap references
- HyperConverged CR fields

### Review Insight Extraction

From review comments, extract:
- **Edge cases**: Comments mentioning "edge case", "corner case", "what if"
- **Concerns**: Comments with "concern", "worry", "problem", "issue"
- **Suggestions**: Comments with "suggest", "should", "consider", "might want"

## File Categorization

Read `{project_context.config_dir}/components.yaml` for project-specific file categorization rules.

The following is an example (CNV) of directory-to-category mapping:

| Directory Pattern | Category |
|:------------------|:---------|
| `pkg/virt-controller/` | controllers |
| `pkg/virt-handler/` | handlers |
| `pkg/virt-api/` | api |
| `pkg/virt-launcher/` | launcher |
| `api/`, `staging/` | api |
| `tests/`, `*_test.go` | tests |
| `config/`, `deploy/` | config |
| `cmd/` | cmd |
| `pkg/util/`, `pkg/util*/` | util |

## Usage Notes

1. **Focus on Behavioral Changes**: Identify what the PR changes functionally
2. **Ignore Noise**: Skip formatting-only changes, comment updates
3. **Highlight Test Implications**: Note what new tests should cover
4. **Extract Edge Cases**: Review comments often reveal test scenarios
