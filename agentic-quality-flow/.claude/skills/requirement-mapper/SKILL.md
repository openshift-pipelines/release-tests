---
name: requirement-mapper
description: Map Jira requirements to testable scenarios with validation gates
---

# Requirement Mapper Skill

**Phase:** Core Processing
**User-Invocable:** false

## Purpose

Map Jira requirements to testable scenarios, applying validation gates and KubeVirt scope boundaries.

## When to Use

Invoked by the **stp-generator** subagent to transform regression analysis into validated requirements.

## Input

```yaml
jira_data:
  main_issue:
    key: CNV-12345
    summary: Add CPU hot-plug support
    description: ...
    acceptance_criteria: ...
  linked_issues: [...]

regression_data:
  impacted_features:
    - feature_name: Live Migration
      relationship: Direct caller
      why_might_break: ...
    - ...
  recommended_tests:
    - requirement: Live migration works with CPU changes
      test_scenario: Verify VM migration succeeds after CPU hot-plug
      priority: P1
    - ...
```

## Output Format

```yaml
validated_requirements:
  - requirement_id: CNV-12345  # Jira issue key — NEVER invent IDs
    requirement_summary: Live migration completes successfully after CPU hot-plug
    source: regression_analysis
    evidence: MigrateVMI calls UpdateVMISpec which was modified
    validation_passed: true
    test_scenario: Verify VM migration succeeds after CPU hot-plug
    priority: P1
  - requirement_id: ""  # Leave blank for subsequent rows under the same epic
    requirement_summary: CPU can be hot-added to running VM
    source: regression_analysis
    evidence: HandleCPUHotplug is new entry point
    validation_passed: true
    test_scenario: Verify CPU addition to running VM
    priority: P0
  - ...

rejected_requirements:
  - requirement_summary: Kubernetes scheduler places VM pods correctly
    reason: Platform-level test - Kubernetes scheduler is tested by platform team
    gate_failed: Requirement Level Validation
  - requirement_summary: PVC binds to PV correctly
    reason: Platform-level test - CSI/storage tested by storage team
    gate_failed: Requirement Level Validation
  - ...

coverage_summary:
  total_from_regression: 15
  validated: 12
  rejected: 3
  tier1_count: 8
  tier2_count: 4
```

## Requirement Level Validation Gate

### Step 1: Identify Testing Level

| Level | Description | Action |
|:------|:------------|:-------|
| Kubernetes Platform | Core K8s (scheduling, storage primitives, RBAC engine) | REJECT |
| OpenShift Platform | OCP features (routes, image streams, OAuth) | REJECT |
| KubeVirt/CNV | Virtualization (VMs, migration, snapshots, CDI) | ACCEPT |
| Integration | KubeVirt using platform capabilities | ACCEPT |

### Step 2: "Who Tests This?" Question

| Answer | Action |
|:-------|:-------|
| Kubernetes upstream QE | REJECT |
| OpenShift Platform QE | REJECT |
| Storage/Network/Security platform team | REJECT |
| KubeVirt/OpenShift Virtualization QE | ACCEPT |

### Step 3: Project Scope Context Check

Read `{project_context.config_dir}/project.yaml` `scope_boundaries` for in-scope and out-of-scope resources.

The following is an example (CNV) of scope boundaries:

**Accept if involves:**
- VirtualMachine, VirtualMachineInstance
- DataVolume, VirtualMachineSnapshot
- VirtualMachineInstanceMigration
- VirtualMachineClusterInstancetype
- HyperConverged, KubeVirt, CDI CRs
- virt-launcher, virt-handler, virt-controller

**Reject if involves only:**
- Pod, Deployment, StatefulSet (raw)
- PersistentVolumeClaim (raw, not DataVolume)
- Node, Namespace (raw)
- ConfigMap, Secret (raw, not VM cloud-init)
- Service, Ingress, Route (raw)

### Step 4: Final Check

Read the `scope_boundaries.validation_gate` question from `{project_context.config_dir}/project.yaml`. For example: "Would removing KubeVirt make this test meaningless?"
- YES → ACCEPT
- NO → REJECT

## Requirement ID Rules

### Jira Issue Keys Only

Requirement IDs MUST be Jira issue keys (e.g., `CNV-72329`). Never invent IDs
like `REQ-xxx-001`, `REQ-NAD-001`, or any other synthetic ID format.

- Use the **epic key** for the first row under that epic
- Leave the Requirement ID **blank** for subsequent rows under the same epic
  (avoids redundant repetition of the same key)
- If a linked sub-issue has its own Jira key, use that key instead

| BAD (Invented) | GOOD (Jira Key) |
|:----------------|:-----------------|
| REQ-NAD-001 | CNV-72329 |
| REQ-CPU-001 | CNV-12345 |
| REQ-MIG-001 | CNV-67890 |

## Requirement Quality Rules

### STP Level Requirements

Requirements must be HIGH-LEVEL capabilities:

| BAD (Too Low-Level) | GOOD (STP Level) |
|:--------------------|:-----------------|
| Create VM with 2 CPUs, start it, add 2 more | CPU can be hot-added to running VM |
| Run `oc get vm` and check status | VM status is accurately reported via API |
| Create PVC, attach, write file, verify | Data persists across disk attach/detach |

### Avoid Trivial Atomic Requirements

Consolidate into feature capabilities:

| BAD (Fragmented) | GOOD (Consolidated) |
|:-----------------|:--------------------|
| Start VM, Stop VM, Restart VM | VM lifecycle operations function correctly |
| Create disk, Attach, Detach, Delete | Disk hot-plug operations complete successfully |
| Add CPU, Remove CPU, Add memory | Resource hot-plug preserves VM stability |

### Target Count

**5-15 high-level requirements per feature** - not 30-50 trivial operations.

## Source Priority

**EXCLUSIVE source for test scenarios:** Regression Analysis

DO NOT derive test scenarios from:
- Jira ticket descriptions
- Acceptance criteria
- PR descriptions or review comments
- Web search results
- General assumptions

## Negative Scenario Coverage

Include negative test scenarios for:
- Invalid input handling
- Resource constraints
- Permission denied
- Invalid state
- Conflict handling
- Recovery/interruption
- Boundary conditions
- Missing dependencies

## Example Mapping

Input (from regression analysis):
```yaml
impacted_features:
  - feature_name: Live Migration
    relationship: Direct caller
    why_might_break: Migration calls VMI update which was modified
```

Output:
```yaml
validated_requirements:
  - requirement_id: CNV-12345
    requirement_summary: Live migration completes successfully after CPU configuration changes
    source: regression_analysis
    evidence: MigrateVMI → UpdateVMISpec (modified)
    validation_passed: true
    test_scenario: Verify VM migration succeeds with modified CPU config
    priority: P1
```

## Coverage Checklist

Before finalizing, verify:
- [ ] All operations covered (every action the feature supports)
- [ ] All configuration options covered
- [ ] All API fields covered
- [ ] All states covered
- [ ] All integration points covered
- [ ] Positive AND negative scenarios included
- [ ] No gaps between regression findings and test scenarios
