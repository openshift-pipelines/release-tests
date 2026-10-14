---
name: jira-collector
description: Collect comprehensive Jira issue data including linked issues and PR URLs
---

# Jira Collector Subagent

**Phase:** Pre-Processing
**Purpose:** Collect comprehensive Jira issue data

## Tools Available

- mcp__mcp-atlassian__jira_get_issue
- mcp__mcp-atlassian__jira_search
- Read

## Required Skills

Must invoke these skills during execution:
1. **jira-parser** - Parse and normalize Jira fields
2. **link-resolver** - Build dependency graph from issue links

## Project Context

This agent receives `project_context` from the orchestrator, which includes:
- `config_dir`: Path to the project configuration directory
- `stp_header`: The expected STP document header
- `versioning`: Version derivation information

## Workflow

### Step 0: Load Project Jira Config

Read `{project_context.config_dir}/jira.yaml` to load:
- Jira project key prefix (e.g., "CNV")
- Custom field mappings (Feature Link, Git Pull Request, etc.)
- Issue type classifications
- Any project-specific Jira configuration

This config is used throughout subsequent steps when interpreting Jira fields.

### Step 1: Fetch Main Jira Issue

Use `mcp__mcp-atlassian__jira_get_issue` with:
- `issue_key`: The provided Jira ID
- `comment_limit`: 100 (to capture all comments)

Extract from the response:
- Summary, description, status
- Issue type and priority
- Labels and components
- Acceptance criteria (from description or custom field)
- "Feature Link" custom field (parent Feature/Epic link)
- Parent issue key and summary (from hierarchy links, e.g., VIRTSTRAT-xxx)
- "Git Pull Request" custom field (all PR links)
- All comments (scan for GitHub PR URLs)

### Step 2: Invoke jira-parser Skill

Invoke the **jira-parser** skill and apply it to normalize the fetched data.

The skill will extract and structure:
- Core issue metadata
- Custom fields (Feature Link, Git Pull Request)
- Comment content with PR URL extraction

### Step 3: Collect Issue Links and Subtasks

From the main issue response, collect:

**Outward Issue Links:**
- Links where the current issue is the SOURCE
- Examples: "this issue blocks CNV-456", "this issue relates to CNV-789"

**Inward Issue Links:**
- Links where the current issue is the TARGET
- Examples: "CNV-123 blocks this issue"

**Subtasks:**
- All subtask issues under the main issue

### Step 4: Fetch Linked Issues (1 Level Depth Only)

#### 4.1 Comprehensive Link Type Collection

Process ALL link types - do not filter by link type name:

| Link Type | Outward | Inward | Category |
|:----------|:--------|:-------|:---------|
| Blocks | blocks | is blocked by | blocking |
| Cloners | clones | is cloned by | clones |
| Duplicate | duplicates | is duplicated by | duplicates |
| Relates | relates to | relates to | relates |
| Parent/Child | is parent of | is child of | hierarchy |
| Epic-Story Link | has Epic | is Epic of | hierarchy |
| Implements | implements | is implemented by | implements |
| Dependency | depends on | is depended on by | dependency |
| Triggers | triggers | is triggered by | triggers |
| Split | split to | split from | split |
| Causes | causes | is caused by | causes |
| Problem/Incident | is problem of | is incident of | problem |

#### 4.2 Collect Full Metadata for Each Linked Issue

For EACH linked issue and subtask:
1. Fetch issue details using `mcp__mcp-atlassian__jira_get_issue` with `comment_limit`: 100
2. Extract complete metadata:
   - **Assignee**: `{name, email}`
   - **Reporter**: `{name, email}`
   - **Components**: Array of component names
   - **Labels**: Array of labels
   - **Fix Version**: Target release version
   - **Created**: ISO 8601 date
   - **Updated**: ISO 8601 date
   - **Description**: Full description text
   - **Acceptance Criteria**: From description or custom field

#### 4.3 **MANDATORY**: Extract PR URLs from Linked Issues

**CRITICAL:** PRs from linked issues MUST be collected:

1. Extract "Git Pull Request" custom field from each linked issue
2. Scan ALL comments for GitHub PR URLs (pattern: `https://github.com/.../pull/...`)
3. Track the source of each PR URL:
   - `source_issue`: The Jira key where the PR was found
   - `source_type`: `custom_field` or `comment`
   - `is_main_issue`: `true` if from main issue, `false` if from linked issue
4. Collect ALL PR URLs - do not deduplicate until aggregation step

**Important:** Do NOT recursively follow links from linked issues. Only process direct links from the main issue.

### Step 5: Invoke link-resolver Skill

Invoke the **link-resolver** skill and apply it to build the dependency graph.

The skill will:
- Categorize link types (blocks, relates to, implements)
- Build hierarchical relationship structure
- Identify key dependencies

### Step 5.5: Extract Feature Candidates for LSP Validation

**This step runs ALWAYS to support LSP validation even when no PRs exist.**

From the parsed issue data, extract potential test features:

#### 5.5.1 From Summary and Description

Extract:
- Technical terms and feature names (capitalized terms, quoted identifiers)
- API types mentioned (VirtualMachine, VMI, DataVolume, VolumeSpec, etc.)
- Component names that map to packages

#### 5.5.2 Component-to-Package Mapping

Read `{project_context.config_dir}/components.yaml` for the component-to-package mapping.

Use the project-specific mapping to resolve Jira component names to package paths.

#### 5.5.3 From Acceptance Criteria

Each acceptance criteria item suggests a testable area. Extract as feature candidates.

#### 5.5.4 From Linked Issues

Extract:
- Related feature names from linked issue summaries
- Dependencies mentioned
- Integration points

#### 5.5.5 Output Feature Candidates

Build a structured list:
```yaml
feature_candidates:
  explicit_mentions:
    - <features/functions/components named in summary>
    - <API types mentioned: VirtualMachine, VMI, etc.>
  component_hints:
    - component: <component name>
      package_path: <mapped package path>
  acceptance_criteria:
    - <each AC item as potential test feature>
  integration_points:
    - <dependencies/integrations from linked issues>
```

### Step 6: Aggregate PR URLs

Compile all GitHub PR URLs from:
- Main issue "Git Pull Request" custom field
- Main issue comments
- Linked issues "Git Pull Request" custom fields
- Linked issues comments
- Subtask "Git Pull Request" custom fields
- Subtask comments

Deduplicate the list.

## Output Format

Return YAML:
```yaml
main_issue:
  key: {JIRA_ID}
  summary: <summary>
  description: <description>
  status: <status>
  issue_type: <type>
  priority: <priority>
  labels: [label1, label2]
  components: [comp1, comp2]
  acceptance_criteria: <criteria or null>
  feature_link: <feature link URL or null>
  parent_issue:
    key: <parent issue key or null>  # e.g., VIRTSTRAT-560
    summary: <parent issue summary or null>
  comments:
    - author: <author>
      created: <date>
      body: <body>
      pr_urls_found: [<URLs found in comment>]
    - ...

linked_issues:
  - key: <linked issue key>  # e.g., CNV-11111
    summary: <summary>
    description: <full description>
    status: <status>
    issue_type: <type>
    relationship: outward
    link_type: blocks
    link_category: blocking
    assignee:
      name: <name>
      email: <email>
    reporter:
      name: <name>
      email: <email>
    components: [comp1, comp2]
    labels: [label1, label2]
    fix_version: <version or null>
    created: <ISO date>
    updated: <ISO date>
    acceptance_criteria: <criteria or null>
    pr_urls:
      - url: https://github.com/.../pull/123
        source_type: custom_field
      - url: https://github.com/.../pull/456
        source_type: comment
  - key: <linked issue key>  # e.g., CNV-22222
    summary: <summary>
    description: <full description>
    status: <status>
    issue_type: <type>
    relationship: inward
    link_type: is blocked by
    link_category: blocking
    assignee: {...}
    reporter: {...}
    components: [...]
    labels: [...]
    fix_version: <version or null>
    created: <ISO date>
    updated: <ISO date>
    acceptance_criteria: <criteria or null>
    pr_urls: [...]
  - ...

subtasks:
  - key: <subtask key>  # e.g., {JIRA_ID}-1
    summary: <summary>
    status: <status>
    pr_urls: [<PR URLs from this subtask>]
  - ...

pr_urls:
  # Example entries - actual repos come from {project_context.config_dir}/repositories.yaml
  - url: https://github.com/<owner>/<repo>/pull/<number>
    source_issue: {JIRA_ID}
    source_type: custom_field
    is_main_issue: true
  - url: https://github.com/<owner>/<repo>/pull/<number>
    source_issue: <linked issue key>
    source_type: comment
    is_main_issue: false
  - ...

feature_candidates:
  explicit_mentions:
    - VirtualMachine
    - HotplugVolume
    - virt-controller
  component_hints:
    - component: virt-handler
      package_path: pkg/virt-handler/
    - component: storage
      package_path: pkg/storage/
  acceptance_criteria:
    - VM can attach volume while running
    - Volume remains attached after migration
  integration_points:
    - Live Migration (from linked issue CNV-11111)
    - Snapshot operations (mentioned in comments)

dependency_graph:
  blocking: [<issues this blocks>]
  blocked_by: [<issues blocking this>]
  related: [<related issues>]
```

## GitHub PR URL Pattern

Scan for URLs matching:
- `https://github.com/{owner}/{repo}/pull/{number}`

Parse each URL to extract:
- owner (e.g., `kubevirt` -- actual owners come from project config)
- repo (e.g., `kubevirt`, `containerized-data-importer` -- actual repos come from project config)
- pullNumber (e.g., `1234`)
