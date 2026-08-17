---
name: tier-classifier
description: Classify test scenarios as Unit Tests, Tier 1, or Tier 2
model: claude-opus-4-6
---

# Tier Classifier Skill

**Phase:** Core Processing
**User-Invocable:** false

## Purpose

Classify test scenarios as Unit Tests, Tier 1 (Functional), or Tier 2 (End-to-End).

## When to Use

Invoked by the **stp-generator** subagent for each test scenario.

## Input

```yaml
scenario:
  requirement_id: CNV-12345
  requirement_summary: CPU can be hot-added to running VM
  test_description: Verify CPU hot-add to running VM
  type: positive
  priority: P0
  fix_scope:                          # optional, from github_data
    files_changed: 1
    functions_changed: ["validateCPUModel"]
    packages_changed: ["pkg/virt-controller/cpu"]
    requires_cluster_interaction: false
    issue_type: bug                   # bug vs feature
```

## Output Format

```yaml
classification:
  requirement_id: CNV-12345
  test_description: Verify CPU hot-add to running VM
  test_type: Tier 1 (Functional)
  reasoning: Tests single feature (hot-plug) in isolation
```

## Valid Test Types

**ONLY these three values are valid:**

| Test Type | Description |
|:----------|:------------|
| `Unit Tests` | Isolated components with mocks; validates individual functions/modules. **Note:** Unit tests are classified for tracking in the STP but are developer-responsibility -- no auto-generation pipeline exists for this tier. |
| `Tier 1 (Functional)` | Single feature in real cluster; API contracts; basic workflows |
| `Tier 2 (End-to-End)` | Complete user workflows; multi-feature integrations; **user-scenario focused** |

## Key Principle: User-Scenario Focus for Tier 2

**Tier 2 tests are strictly user-scenario focused.** They validate what end users experience and interact with, not internal system behavior, implementation details, or diagnostic information.

**Key Principle:** Tests should only verify observable user outcomes, not internal system state or logs.

## Decision Matrix

| Question | Unit | Tier 1 | Tier 2 |
|:---------|:-----|:-------|:-------|
| Tests isolated functions with mocks? | YES | no | no |
| Tests single feature in real cluster? | no | YES | no |
| Requires multiple features working together? | no | no | YES |
| Tests basic API or component functionality? | no | YES | no |
| Validates complete user workflow? | no | no | YES |
| Can run without cluster (mocked dependencies)? | YES | no | no |
| Requires minimal test cluster? | no | YES | no |
| Requires production-like environment? | no | no | YES |
| Tests upgrade or migration paths? | no | no | YES |
| Tests at scale (100+ resources)? | no | no | YES |
| Involves multiple VMs interacting? | no | no | YES |
| Tests data persistence across operations? | no | no | YES |

## Classification Flow (Updated)

```
0. Fix-Scope Demotion Check (optional)
   SKIP if fix_scope is absent OR issue_type is feature/enhancement.
   ONLY activate when fix_scope is present AND issue_type is bug/customer_case/defect.

   a. Single function changed AND requires_cluster_interaction is false?
      YES -> Unit Tests
             reasoning: "Fix modifies single function {name} with no cluster
             interaction. Unit test provides equivalent coverage at lower cost."
      NO  -> Continue

   b. Single package changed AND single resource type?
      YES -> Tier 1 (Functional)
             reasoning: "Fix is scoped to {package}, single resource operation.
             Tier 1 provides equivalent coverage."
      NO  -> Continue to Step 1 (no demotion)

1. Does it require a cluster?
   NO  -> Unit Tests
   YES -> Continue

2. Check Tier 2 PROMOTION triggers first (see below)
   ANY trigger matches -> Tier 2 (End-to-End)
   NO triggers match -> Continue

3. Does it test a single feature in isolation?
   YES -> Tier 1 (Functional)
   NO  -> Tier 2 (End-to-End)
```

**IMPORTANT:** Check Tier 2 triggers BEFORE defaulting to Tier 1.

## Tier 2 Promotion Triggers

**Qualifying Rule:** A trigger matches only when the **test itself** exercises that
workflow as its primary action — not when the feature merely uses that mechanism
internally. Classify based on what the **test** does, not what the **feature** does
under the hood.

Example: A feature that uses live migration internally to apply a NAD change does NOT
make a test "Live migration with workload validation" — unless the test's primary
action is to perform and validate a migration. If the test patches a spec field and
checks a VM condition, it's Tier 1 regardless of whether migration happens behind
the scenes.

**If ANY of these are true for what the test exercises, classify as Tier 2:**

| Trigger | Example |
|:--------|:--------|
| Involves multiple VMs interacting | Multi-tier app deployment |
| Tests complete user story/workflow | Create VM -> Run workload -> Migrate -> Verify |
| Resources must survive across operations | VM state preserved through migration |
| Validates data/state persistence across operations | Snapshot -> Restore -> Verify data |
| Tests upgrade or version compatibility | Upgrade from 4.18 to 4.19 |
| Requires external systems | External router, load balancer |
| Simulates production deployment | Full application stack |
| Tests disaster recovery or failover | Node failure recovery |
| RBAC across multiple resources/operations | User permissions through VM lifecycle |
| Storage lifecycle with multiple steps | Provision -> Attach -> Snapshot -> Restore |
| Live migration with workload validation | Migrate while workload running, verify continuity |

## What's NOT in Tier 1

**Classify as Tier 2 (not Tier 1) if the scenario involves:**

- Multi-feature integration scenarios
- Complex end-to-end user workflows and user stories
- Performance and scale testing
- Upgrade scenarios
- Disaster recovery scenarios
- Multi-step workflows (create -> operate -> verify persistence)
- Cross-component interactions

## What Tier 2 Does NOT Test

**Do NOT classify as Tier 2 if testing:**

- Internal debug logs validation (not user-facing)
- Internal component implementation details
- Code-level unit behaviors
- Low-level API internals not exposed to users
- Developer debugging workflows
- Kubernetes/OpenShift platform features (not virtualization)
- System metrics users don't interact with
- Internal error messages or stack traces

**Note:** Tests may verify user-observable Kubernetes Events (user-facing API) but should not parse internal pod logs.

## Tier 1 (Functional) Indicators

Classify as Tier 1 if:
- Tests a single feature in isolation
- Validates API contracts
- Basic CRUD operations
- Single resource lifecycle
- Error handling for single feature
- Basic configuration validation
- **AND** no Tier 2 promotion triggers apply

**Examples:**
- Create VM with DataVolume
- Attach network interface via API
- Create VM snapshot (single operation)
- Stop running VM
- Hot-plug single disk

## Tier 2 (End-to-End) Indicators

Classify as Tier 2 if:
- Requires multiple features working together
- Tests complete user workflow
- Involves cross-component interaction
- Requires production-like environment
- Tests upgrade/migration paths
- Tests at scale (100+ resources)
- Involves multi-step scenarios with state verification

**Examples:**
- Deploy multi-tier app with VMs
- Live migration with workload validation
- Create -> Snapshot -> Restore -> Verify workflow
- Upgrade from version X to Y
- RBAC workflow across VM lifecycle
- Storage lifecycle (provision -> attach -> snapshot -> restore)
- Multi-VM network communication
- CPU hotplug followed by migration and state verification

## Unit Test Indicators

Classify as Unit Tests if:
- Tests individual function/method
- Uses mocks for dependencies
- No cluster required
- Developer responsibility typically

**Examples:**
- Validate input parsing function
- Test error message formatting
- Test configuration parsing

## Common Misclassifications

| Scenario | Wrong | Correct | Reason |
|:---------|:------|:--------|:-------|
| Deploy 3-tier app | Tier 1 | Tier 2 | Multi-VM workflow |
| VM migration (single) | Tier 2 | Tier 1 | Single feature operation |
| API validation | Unit | Tier 1 | Requires cluster |
| Upgrade with running VMs | Tier 1 | Tier 2 | Multi-step, cross-version |
| Hot-plug single disk | Tier 2 | Tier 1 | Single feature |
| Migrate then verify workload | Tier 1 | Tier 2 | Multi-step with state verification |
| Snapshot and restore | Tier 1 | Tier 2 | Multi-step workflow |
| VM survives node drain | Tier 1 | Tier 2 | Cross-component, DR scenario |
| Scale test with 100 VMs | Tier 1 | Tier 2 | Scale testing |
| CPU hotplug + migration | Tier 1 | Tier 2 | Multi-feature integration |

## Priority Influence

Priority doesn't determine tier:
- P0 can be Tier 1 or Tier 2
- P2 can be Tier 1 or Tier 2

Tier is based on **scope and complexity**, not importance.

## Output Examples

Input:
```yaml
test_description: Verify CPU hot-add to running VM
```

Output:
```yaml
test_type: Tier 1 (Functional)
reasoning: Tests single feature (CPU hot-plug) in real cluster, no multi-step workflow
```

Input:
```yaml
test_description: Verify VM state preserved through snapshot and restore
```

Output:
```yaml
test_type: Tier 2 (End-to-End)
reasoning: Multi-step workflow (create -> snapshot -> restore -> verify state)
```

Input:
```yaml
test_description: Verify upgrade preserves VM configuration
```

Output:
```yaml
test_type: Tier 2 (End-to-End)
reasoning: Cross-version testing, requires upgrade scenario
```

Input:
```yaml
test_description: Verify CPU hotplug followed by live migration preserves CPU count
```

Output:
```yaml
test_type: Tier 2 (End-to-End)
reasoning: Multi-feature integration (hotplug + migration), state verification across operations
```

Input:
```yaml
test_description: Verify VM can be created with 216 cores
```

Output:
```yaml
test_type: Tier 1 (Functional)
reasoning: Single feature (VM creation), single operation, no multi-step workflow
```
