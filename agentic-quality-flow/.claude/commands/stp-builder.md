---
name: stp-builder
description: Generate a Software Test Plan (STP) from a Jira ticket
argument-hint: <JIRA-ID or URL>
allowed-tools: Read, Write, Edit, Task, Glob, Grep, LSP, Skill, mcp__mcp-atlassian__jira_get_issue, mcp__mcp-atlassian__jira_search, mcp__github__pull_request_read, mcp__github__get_file_contents
---

# Generate STP for $ARGUMENTS

You are the STP Builder entry point. Initiate the STP generation workflow by activating the **stp-orchestrator** subagent.

## Input

The user has provided: `$ARGUMENTS`

This should be a Jira ticket ID (e.g., `CNV-12345`, `VIRTSTRAT-494`) or a Jira URL (e.g., `https://issues.redhat.com/browse/CNV-12345`).

## Workflow

### Step 0: Resolve Project

Use the Skill tool to invoke the project-resolver skill:

**Tool:** Skill
**Parameters:**
- skill: "project-resolver"
- args: "$ARGUMENTS"

This returns `project_context` containing:
- `project_id`, `display_name`, `jira_id`
- `config_dir` (path to project config files)
- `feature_toggles` (what capabilities are enabled)
- `stp_header`, `versioning`

**If project resolution fails:** Display the error and exit. Do not proceed.

**Check stp_generation toggle:**
If `project_context.feature_toggles.stp_generation` is false:
- Output: "STP generation is disabled for project {project_context.display_name} (stp_generation toggle is false)."
- Exit. Do not proceed.

### Step 1: Activate Orchestrator

Activate the **stp-orchestrator** subagent with the resolved Jira ticket ID **and** `project_context`.

Pass to orchestrator:
```yaml
jira_id: "{JIRA_ID}"
project_context: <from project-resolver>
```

The orchestrator will:
1. **Pre-processing Phase (Sequential Pipeline)**:
   - Launch **jira-collector** (cyan) to fetch Jira issue data and PR URLs
   - Launch **github-pr-fetcher** (green) with PR URLs from jira-collector
   - Launch **regression-analyzer** (yellow) with changed files from github-pr-fetcher

2. **Core Processing Phase (Sequential)**:
   - Pass aggregated data to **stp-generator** (purple)
   - Generate complete STP document with test scenarios

3. **Post-processing Phase (Sequential)**:
   - Pass document to **document-formatter** (orange)
   - Sanitize PII, validate structure, format tables
   - Save to `outputs/stp/{JIRA_ID}/{JIRA_ID}_test_plan.md`

## Expected Output

A complete Software Test Plan markdown file saved to `outputs/stp/{JIRA_ID}/{JIRA_ID}_test_plan.md`.

## Activation

1. Invoke the **project-resolver** skill with `$ARGUMENTS` to get `project_context`.
2. Activate the **stp-orchestrator** agent, passing both the Jira ticket ID and `project_context`.
