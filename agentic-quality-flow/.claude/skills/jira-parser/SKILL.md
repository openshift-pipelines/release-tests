---
name: jira-parser
description: Parse and normalize Jira issue fields into a structured format for STP generation
model: claude-opus-4-6
---

# Jira Parser Skill

**Phase:** Pre-Processing
**User-Invocable:** false

## Purpose

Parse and normalize Jira issue fields into a structured format for STP generation.

## When to Use

Invoked by the **jira-collector** subagent after fetching raw Jira issue data.

## Input

Raw Jira issue response containing:
- Standard fields (summary, description, status, priority, etc.)
- Custom fields (Feature Link, Git Pull Request, etc.)
- Issue links (outward and inward)
- Comments

## Output Format

```yaml
parsed_issue:
  key: CNV-12345
  summary: <issue summary>
  description: <full description text>
  status: <In Progress | Done | etc>
  issue_type: <Story | Bug | Enhancement | Epic | Task>
  priority: <Critical | Major | Minor | etc>
  labels:
    - label1
    - label2
  components:
    - component1
    - component2
  acceptance_criteria: <extracted from description or custom field>
  feature_link: <URL or null>
  git_pull_requests:
    - https://github.com/...
  reporter:
    name: <name>
    email: <email>
  assignee:
    name: <name>
    email: <email>
  created: <date>  # ISO 8601 format
  updated: <date>  # ISO 8601 format
  fix_version: <version string or null>
  custom_fields:
    - id: customfield_12345
      name: Feature Link
      value: <value>
    - id: customfield_67890
      name: Git Pull Request
      value: <value>
    - ...  # All non-null custom fields
```

## Field Extraction Rules

### Standard Fields

| Jira Field | Output Field | Notes |
|:-----------|:-------------|:------|
| `summary` | `summary` | Direct mapping |
| `description` | `description` | Full text, may contain markdown |
| `status.name` | `status` | Status display name |
| `issuetype.name` | `issue_type` | Issue type display name |
| `priority.name` | `priority` | Priority display name |
| `labels` | `labels` | Array of strings |
| `components[].name` | `components` | Array of component names |
| `fixVersions[0].name` | `fix_version` | First fix version or null |
| `reporter` | `reporter` | {name, email} object |
| `assignee` | `assignee` | {name, email} object |
| `created` | `created` | ISO 8601 timestamp |
| `updated` | `updated` | ISO 8601 timestamp |

### Custom Fields

Custom field IDs vary by Jira instance. Search for these patterns:

| Custom Field Name | Common Patterns | Output Field |
|:------------------|:----------------|:-------------|
| Feature Link | `customfield_*` containing "Feature", "Epic Link" | `feature_link` |
| Git Pull Request | `customfield_*` containing "Git", "Pull Request", "PR" | `git_pull_requests` |

### Custom Field Extraction Rules

Extract ALL non-null custom fields into the `custom_fields` array:

| Rule | Description |
|:-----|:------------|
| Include | All custom fields with non-null, non-empty values |
| Exclude | System fields (already mapped above) |
| Format | `{id, name, value}` for each field |
| Name Discovery | Use Jira API field metadata to get display names |
| Value Types | Strings, arrays, objects - preserve original structure |

**MANDATORY Custom Fields to Look For:**
1. Git Pull Request - Contains PR URLs (critical for PR collection)
2. Feature Link / Epic Link - Parent feature reference
3. Acceptance Criteria - Testing requirements
4. Story Points - Estimation data
5. Sprint - Sprint assignment
6. Target Version - Release targeting

### Acceptance Criteria Extraction

Look for acceptance criteria in:
1. Custom field named "Acceptance Criteria"
2. Description section starting with "Acceptance Criteria:" or "AC:"
3. Description section with `h3. Acceptance Criteria` (Jira wiki markup)

### Comment Processing

For each comment, extract:
```yaml
comments:
  - id: <comment id>
    author:
      name: <display name>
      email: <email if available>
    created: <timestamp>
    body: <comment body>
    pr_urls_found:
      - <any GitHub PR URLs found in body>
```

### GitHub PR URL Extraction

Scan text for patterns matching:
- `https://github.com/{owner}/{repo}/pull/{number}`
- `github.com/{owner}/{repo}/pull/{number}` (without https)

Return deduplicated list of full URLs.

## Normalization Rules

1. **Text cleaning**: Remove excessive whitespace, normalize line endings
2. **URL normalization**: Ensure all URLs have https:// prefix
3. **Empty handling**: Use `null` for missing optional fields, empty arrays for lists
4. **Date formatting**: ISO 8601 format (YYYY-MM-DDTHH:MM:SSZ)

## Example

Input (raw Jira response):
```json
{
  "key": "CNV-12345",
  "fields": {
    "summary": "Add hot-plug support for CPU",
    "description": "As a user, I want to...\n\nAcceptance Criteria:\n- CPU can be added\n- CPU can be removed",
    "status": {"name": "In Progress"},
    "issuetype": {"name": "Story"},
    "priority": {"name": "Major"},
    "labels": ["hot-plug", "cpu"],
    "components": [{"name": "virt-controller"}],
    "customfield_12345": "https://issues.redhat.com/browse/CNV-99000",
    "customfield_67890": "https://github.com/kubevirt/kubevirt/pull/1234"
  }
}
```

Output:
```yaml
parsed_issue:
  key: CNV-12345
  summary: Add hot-plug support for CPU
  description: "As a user, I want to...\n\nAcceptance Criteria:\n- CPU can be added\n- CPU can be removed"
  status: In Progress
  issue_type: Story
  priority: Major
  labels:
    - hot-plug
    - cpu
  components:
    - virt-controller
  acceptance_criteria: |
    - CPU can be added
    - CPU can be removed
  feature_link: https://issues.redhat.com/browse/CNV-99000
  git_pull_requests:
    - https://github.com/kubevirt/kubevirt/pull/1234
```
