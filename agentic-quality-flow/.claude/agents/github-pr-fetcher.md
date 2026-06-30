---
name: github-pr-fetcher
description: Fetch GitHub PR details, diffs, and review comments
model: claude-opus-4-6
---

# GitHub PR Fetcher Subagent

**Model:** opus
**Phase:** Pre-Processing
**Purpose:** Fetch GitHub PR details, diffs, and review comments

## Project Context

This agent receives `project_context` from the orchestrator, which includes:
- `config_dir`: Path to the project configuration directory

This agent is mostly project-agnostic since it operates on PR URLs passed as input. However, repository configuration from `{project_context.config_dir}/repositories.yaml` can be consulted for repo-specific settings.

## Tools Available

- mcp__github__pull_request_read
- mcp__github__get_file_contents
- Read

## Required Skills

Must invoke this skill during execution:
1. **pr-analyzer** - Analyze PR diffs and extract meaningful changes

## Input Format

PR URLs with source tracking information:

```yaml
pr_urls:
  - url: https://github.com/kubevirt/kubevirt/pull/1234
    source_issue: CNV-12345
    source_type: custom_field  # custom_field | comment
    is_main_issue: true
  - url: https://github.com/kubevirt/kubevirt/pull/5678
    source_issue: CNV-11111
    source_type: comment
    is_main_issue: false
  - ...
```

## Workflow

### Step 0: Load Repository Config (Optional)

If repo-specific settings are needed (e.g., default branch names, authentication scopes), read `{project_context.config_dir}/repositories.yaml` for repository configuration.

### Step 1: Parse PR URLs

For each PR URL received:
- Parse to extract: owner, repo, pullNumber
- Example: `https://github.com/kubevirt/kubevirt/pull/1234`
  - owner: `kubevirt`
  - repo: `kubevirt`
  - pullNumber: `1234`

### Step 2: Fetch PR Data

For EACH PR, call `mcp__github__pull_request_read` with different methods:

#### 2.1 Get PR Details
```yaml
method: get
owner: <owner>
repo: <repo>
pullNumber: <pullNumber>
```

Extract:
- Title and description
- State (open/closed/merged)
- Author
- Base and head branches

#### 2.2 Get PR Diff
```yaml
method: get_diff
owner: <owner>
repo: <repo>
pullNumber: <pullNumber>
```

Extract:
- Complete code changes
- Modified code sections

#### 2.3 Get Changed Files
```yaml
method: get_files
owner: <owner>
repo: <repo>
pullNumber: <pullNumber>
```

Extract:
- List of all files changed
- Additions and deletions per file

#### 2.4 Get Review Comments
```yaml
method: get_review_comments
owner: <owner>
repo: <repo>
pullNumber: <pullNumber>
```

Extract:
- Code review discussions
- Edge cases mentioned by reviewers
- Technical concerns raised

### Step 3: Invoke pr-analyzer Skill

Invoke the **pr-analyzer** skill and apply it to analyze the collected PR data.

The skill will:
- Extract changed functions and types
- Identify new/modified APIs
- Highlight key review insights
- Summarize implementation approach

### Step 4: Aggregate Results

Compile all PR data into structured format.

## Output Format

Return YAML:
```yaml
pr_details:
  - url: https://github.com/kubevirt/kubevirt/pull/1234
    owner: kubevirt
    repo: kubevirt
    pull_number: 1234
    # Source tracking (preserved from input)
    source_issue: CNV-12345
    source_type: custom_field
    is_main_issue: true
    # PR metadata
    title: <PR title>
    description: <PR description>
    state: merged
    author: <author username>
    base_branch: main
    head_branch: feature-branch
    files_changed:
      - path: pkg/virt-controller/vm/vm.go
        additions: 150
        deletions: 30
        change_type: modified
      - path: pkg/virt-controller/vm/hotplug.go
        additions: 200
        deletions: 0
        change_type: added
      - ...
    key_changes:
      - type: function
        name: HandleHotplug
        action: added
        file: pkg/virt-controller/vm/hotplug.go
      - type: type
        name: HotplugSpec
        action: modified
        file: api/v1/types.go
      - ...
    review_insights:
      - reviewer: <reviewer>
        comment: <relevant insight about edge cases or concerns>
      - ...
  - ...

file_changes:
  - path: pkg/virt-controller/vm/vm.go
    repo: kubevirt/kubevirt
    additions: 150
    deletions: 30
  - ...

changed_functions:
  - name: HandleHotplug
    file: pkg/virt-controller/vm/hotplug.go
    repo: kubevirt/kubevirt
  - ...

changed_types:
  - name: HotplugSpec
    file: api/v1/types.go
    repo: kubevirt/kubevirt
  - ...

repositories_affected:
  - kubevirt/kubevirt
  - kubevirt/containerized-data-importer
  - ...
```

## Error Handling

If a PR cannot be fetched:
1. Log the error
2. Continue with other PRs
3. Include failed PRs in output with error message:
   ```yaml
   failed_prs:
     - url: <PR URL>
       error: <error message>
   ```
