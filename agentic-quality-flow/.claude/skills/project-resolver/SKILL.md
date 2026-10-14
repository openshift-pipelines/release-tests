---
name: project-resolver
description: Resolve Jira ID to project configuration and load project context
---

# Project Resolver Skill

**Phase:** Pre-Processing (Step 0)
**User-Invocable:** false

## Purpose

Central config loader for QualityFlow's multi-project architecture. Every command
invokes this skill as Step 0 to resolve the Jira ID to a project and load its
configuration.

## When to Use

Invoked as the **first step** of every command (`stp-builder`, `std-builder`,
`generate-go-tests`, `generate-python-tests`) before any other processing.

## Tools Required

- Read

## Input

```yaml
jira_input: "CNV-66855"  # or "https://issues.redhat.com/browse/CNV-66855"
```

## Workflow

### Step 1: Parse Jira ID

Extract the Jira ID from the input. Handle both formats:
- Direct ID: `CNV-12345` → extract prefix `CNV`, ID `CNV-12345`
- URL: `https://issues.redhat.com/browse/CNV-12345` → extract prefix `CNV`, ID `CNV-12345`

The prefix is the text before the first hyphen in the Jira key.

### Step 2: Read Routing Configuration

Read `config/routing.yaml` from the project root.

Extract the `routes` array and `default_project` value.

### Step 3: Resolve Project

Match the extracted prefix against `routes[].prefix`:

```
For each route in routes:
  if route.prefix == extracted_prefix:
    project_id = route.project
    break
```

If no match found:
- If `default_project` is not null: use `default_project`
- If `default_project` is null: **FAIL** with error:
  ```
  Unknown Jira prefix "{prefix}". No project configured for this prefix.
  Known prefixes: CNV, VIRTSTRAT, OCPBUGS, MTV
  To add a new project, create config/projects/{name}/ and add a route in config/routing.yaml.
  ```

### Step 4: Validate Project Directory

Check that `config/projects/{project_id}/` exists and contains the required files.

Read `config/_schema.yaml` to get the `required_files` list.

For each required file, verify it exists at `config/projects/{project_id}/{file}`.

If any required file is missing: **FAIL** with error:
```
Project "{project_id}" is missing required config file: {file}
Expected at: config/projects/{project_id}/{file}
```

### Step 5: Load Defaults

Read `config/_defaults.yaml` and extract the `feature_toggles` defaults.

### Step 6: Load Project Config

Read `config/projects/{project_id}/project.yaml` and extract:
- `project_id`
- `display_name`
- `feature_toggles` (project-specific overrides)
- `stp_document.header`
- `versioning`

### Step 7: Merge Feature Toggles

Deep-merge project toggles over defaults:

```
merged_toggles = defaults.feature_toggles
for key, value in project.feature_toggles:
  merged_toggles[key] = value
```

### Step 8: Validate Toggle Consistency

Read `config/_schema.yaml` `toggle_consistency` rules.

For each rule:
- If `merged_toggles[rule.toggle]` is true, verify `config/projects/{project_id}/{rule.requires_file}` exists
- If the required file is missing: **WARN** (not fail):
  ```
  Warning: {rule.toggle} is enabled but {rule.requires_file} not found.
  ```

### Step 9: Return Project Context

Return the resolved context:

```yaml
project_context:
  project_id: "{project_id}"
  display_name: "{display_name}"
  jira_id: "{JIRA_ID}"
  config_dir: "config/projects/{project_id}"
  feature_toggles:
    polarion: true/false
    unit_tests: true/false
    tier1_tests: true/false
    tier2_tests: true/false
    stp_generation: true/false
    std_generation: true/false
    lsp_analysis: true/false
    pii_sanitization: true/false
  stp_header: "{stp_document.header}"
  versioning:
    product_name: "{product_name}"
    platform_name: "{platform_name}"
    current_version: "{current_version}"
```

## Output Format

```yaml
project_context:
  project_id: "cnv"
  display_name: "OpenShift Virtualization (CNV)"
  jira_id: "CNV-66855"
  config_dir: "config/projects/cnv"
  feature_toggles:
    polarion: true
    unit_tests: false
    tier1_tests: true
    tier2_tests: true
    stp_generation: true
    std_generation: true
    lsp_analysis: true
    pii_sanitization: true
  stp_header: "Openshift-virtualization-tests Test plan"
  versioning:
    product_name: "OpenShift Virtualization"
    platform_name: "OCP"
    current_version: "4.22"
```

## Error Handling

**Unknown prefix:**
- Error: "Unknown Jira prefix. No project configured."
- Action: List known prefixes and suggest adding a route
- Exit command

**Missing project directory:**
- Error: "Project config directory not found"
- Action: Suggest creating the directory structure
- Exit command

**Missing required config file:**
- Error: "Required config file missing"
- Action: List the missing file and expected location
- Exit command

**Malformed YAML:**
- Error: "Cannot parse config file"
- Action: Show the file path and suggest checking YAML syntax
- Exit command

## Usage by Commands

Each command uses project_context differently:

| Command | Uses from project_context |
|:--------|:--------------------------|
| stp-builder | Passes to stp-orchestrator for all subagents |
| std-builder | Checks tier1_tests/tier2_tests to decide which stubs to generate |
| generate-go-tests | Checks tier1_tests; blocks if false |
| generate-python-tests | Checks tier2_tests; blocks if false |

## Usage by Agents

Each agent reads additional config files on-demand from `config_dir`:

| Agent | Reads from config_dir |
|:------|:----------------------|
| jira-collector | `jira.yaml`, `components.yaml` |
| github-pr-fetcher | `repositories.yaml` (optional) |
| regression-analyzer | `repositories.yaml`, `components.yaml` |
| stp-generator | `project.yaml`, `environment.yaml`, `tier1.yaml`, `tier2.yaml` |
| document-formatter | `pii_exceptions.yaml` |
| ticket-context-analyzer | `repositories.yaml` |

## Feature Toggle Notes

The `unit_tests` toggle is informational only. It signals whether unit tests are in scope for a project configuration, but no QualityFlow command or skill gates on it. All other toggles (`polarion`, `tier1_tests`, `tier2_tests`, `stp_generation`, `std_generation`, `lsp_analysis`, `pii_sanitization`) are actively gated by commands, agents, or skills.
