---
name: lsp-tracer
description: Trace call graphs using LSP operations to identify regression impact
model: claude-opus-4-6
---

# LSP Tracer Skill

**Phase:** Pre-Processing
**User-Invocable:** false

## Purpose

Trace call graphs using LSP operations to identify regression impact.

## When to Use

Invoked by the **regression-analyzer** subagent to trace code dependencies.

## Tools Required

- LSP (with operations: goToDefinition, findReferences, incomingCalls, outgoingCalls)

## Input

```yaml
# Repo path: Read `repositories.yaml` from `project_context.config_dir` to get the repo path.
symbols_to_trace:
  - name: HandleCPUHotplug
    file: <repo_path>/pkg/virt-controller/vm/vm.go
    line: 105
    character: 6
  - name: CPUHotplugSpec
    file: <repo_path>/api/v1/types.go
    line: 250
    character: 6
  - ...

# Optional: If symbols_to_trace is empty but feature_candidates is provided
feature_candidates:
  explicit_mentions:
    - VirtualMachine
    - multiarch
  component_hints:
    - component: virt-handler
      package_path: pkg/virt-handler/
    - component: node-labeller
      package_path: pkg/virt-handler/node-labeller/
  acceptance_criteria:
    - VM can run on ARM nodes
```

## Alternative Entry Point Discovery (when no PR data)

If `symbols_to_trace` is empty but `feature_candidates` is provided:

### 1. Discovery from explicit_mentions

For each candidate in `explicit_mentions`:
1. Use LSP `workspaceSymbol` to search for the symbol
2. If found, add to symbols_to_trace with `source: "jira_candidate"`

```yaml
# Example workspaceSymbol query
LSP operation="workspaceSymbol" query="VirtualMachine" filePath="<repo_path>/pkg/" line=1 character=1
```

### 2. Discovery from component_hints

For each component_hint:
1. Map to package path using the mapping table
2. Use Glob to find main files in the package
3. Use LSP `documentSymbol` to list exported functions
4. Add relevant functions (uppercase first letter = exported) to symbols_to_trace

**Component-to-Package Mapping:**

**Reference:** Read `{project_context.config_dir}/components.yaml` for the complete mapping.

### 3. Discovery from acceptance_criteria

Parse each acceptance criteria item for:
- Technical terms (capitalized words)
- API type names
- Action verbs that map to functions (e.g., "attach" → "Attach*", "migrate" → "Migrate*")

Use Grep to find matching symbols:
```
Grep pattern="func Migrate" path="<repo_path>/pkg/"
```

### 4. Output Discovered Entry Points

```yaml
discovered_entry_points:
  - name: ReconcileNode
    file: pkg/virt-handler/node-labeller/node_labeller.go
    line: 45
    character: 6
    source: component_hint
    component: node-labeller
  - name: IsARM
    file: pkg/virt-handler/node-labeller/util.go
    line: 20
    character: 6
    source: explicit_mention
    candidate: ARM
  - ...
```

## Output Format

```yaml
call_graph:
  - symbol: HandleCPUHotplug
    file: pkg/virt-controller/vm/vm.go
    line: 105

    incoming_calls:  # Who calls this function
      - caller: ReconcileVM
        file: pkg/virt-controller/vm/vm.go
        line: 45
        relationship: direct
      - caller: ProcessVMUpdate
        file: pkg/virt-controller/vm/update.go
        line: 89
        relationship: direct
      - ...

    outgoing_calls:  # What this function calls
      - callee: ValidateCPUChange
        file: pkg/virt-controller/vm/validation.go
        line: 120
        relationship: direct
      - callee: UpdateVMISpec
        file: pkg/virt-handler/vmi/vmi.go
        line: 230
        relationship: direct
      - ...

    references:  # All code that references this symbol
      - file: pkg/virt-controller/vm/vm_test.go
        line: 456
        context: test
      - file: tests/hotplug_test.go
        line: 78
        context: test
      - ...

  - symbol: CPUHotplugSpec
    file: api/v1/types.go
    line: 250

    usages:
      - file: pkg/virt-controller/vm/vm.go
        line: 110
        usage_type: field_access
      - file: pkg/virt-handler/vmi/vmi.go
        line: 340
        usage_type: type_assertion
      - ...

dependency_chains:
  - chain_name: CPU Hotplug → Migration
    path:
      - symbol: HandleCPUHotplug
        file: pkg/virt-controller/vm/vm.go
      - symbol: UpdateVMISpec
        file: pkg/virt-handler/vmi/vmi.go
      - symbol: PrepareMigration
        file: pkg/virt-handler/migration/migration.go
    impact: Migration may be affected by CPU hotplug changes
  - ...

summary:
  symbols_traced: 5
  total_callers: 12
  total_callees: 8
  total_references: 45
  max_chain_depth: 100
```

## LSP Operation Guide

### Finding Incoming Calls (Who Calls This)

```yaml
operation: incomingCalls
filePath: <absolute path to file>
line: <1-based line number>
character: <1-based column>
```

Returns list of functions that call the target function.

### Finding Outgoing Calls (What This Calls)

```yaml
operation: outgoingCalls
filePath: <absolute path to file>
line: <1-based line number>
character: <1-based column>
```

Returns list of functions that the target function calls.

### Finding All References

```yaml
operation: findReferences
filePath: <absolute path to file>
line: <1-based line number>
character: <1-based column>
```

Returns all locations where the symbol is referenced.

### Going to Definition

```yaml
operation: goToDefinition
filePath: <absolute path to file>
line: <1-based line number>
character: <1-based column>
```

Returns the location where the symbol is defined.

## Tracing Strategy

### Level 1: Direct Dependencies
- All incoming calls to changed functions
- All outgoing calls from changed functions
- All usages of changed types

### Level 2: Indirect Dependencies
- Callers of the direct callers
- Callees of the direct callees
- Types that embed changed types

### Level 3: Transitive (Optional)
- Only trace if Level 2 shows high-risk patterns
- Stop at 100 levels to avoid explosion

## Path Normalization

**Repository Base:** Read `repositories.yaml` from `project_context.config_dir` to get the repo path.

When reporting paths, use relative paths from repo root:
- Absolute: `<repo_path>/pkg/virt-controller/vm/vm.go`
- Relative: `pkg/virt-controller/vm/vm.go`

## Feature Mapping

Map code locations to features by reading `{project_context.config_dir}/components.yaml` `path_to_feature` mapping.

The following is an example (CNV) of path-to-feature mapping:

| Path Pattern | Feature |
|:-------------|:--------|
| `pkg/virt-controller/vm/` | VM Lifecycle |
| `pkg/virt-handler/migration/` | Live Migration |
| `pkg/virt-controller/migration/` | Live Migration |
| `pkg/*/hotplug/` | Hot-plug |
| `pkg/virt-controller/snapshot/` | Snapshots |
| `pkg/virt-controller/clone/` | Clone |
| `pkg/network/` | Networking |
| `pkg/storage/` | Storage |
| `pkg/instancetype/` | Instance Types |

## Depth Limits

- **Maximum Call Chain Depth:** 100 levels
- **Maximum References:** 50 per symbol
- **Stop Conditions:**
  - Reaching test files (unless analyzing test impact)
  - Reaching standard library
  - Reaching external dependencies

## Example Trace

Input:
```yaml
# Repo path from `repositories.yaml` in `project_context.config_dir`
symbols_to_trace:
  - name: HandleCPUHotplug
    file: <repo_path>/pkg/virt-controller/vm/vm.go
    line: 105
    character: 6
```

LSP Calls:
1. `incomingCalls` on HandleCPUHotplug → finds ReconcileVM
2. `outgoingCalls` on HandleCPUHotplug → finds ValidateCPUChange, UpdateVMISpec
3. `incomingCalls` on UpdateVMISpec → finds PrepareMigration (Level 2)
4. Stop - reached migration feature (dependency identified)

Output:
```yaml
dependency_chains:
  - chain_name: CPU Hotplug → VMI Update → Migration
    path:
      - HandleCPUHotplug (vm.go:105)
      - UpdateVMISpec (vmi.go:230)
      - PrepareMigration (migration.go:45)
    impact: Live migration depends on VMI spec updates
    recommended_test: Verify migration works after CPU hotplug
```
