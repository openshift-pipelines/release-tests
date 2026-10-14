---
name: table-generator
description: Generate properly formatted markdown tables for STP documents
model: claude-opus-4-6
---

# Table Generator Skill

**Phase:** Utility
**User-Invocable:** false

## Purpose

Generate properly formatted markdown tables for STP documents.

## When to Use

Invoked by the **document-formatter** subagent to ensure consistent table formatting.

## Input

```yaml
table:
  type: metadata | requirements | testing_types | environment | risks | custom
  headers:
    - Field
    - Details
  rows:
    - ["**Label**", "Value"]
    - ["**Label**", "Value"]
    - ...
```

## Output Format

```yaml
formatted_table: |
  | Field | Details |
  |:------|:--------|
  | **Label** | Value |
  | **Label** | Value |
  ...
```

## Table Formatting Rules

### 1. Header Row
```markdown
| Column 1 | Column 2 | Column 3 |
```
- Single pipe at start and end
- Single space padding around content
- No trailing whitespace

### 2. Alignment Row
```markdown
|:-------|:-------|:-------|
```
- Use `|:---` for left alignment (default)
- Use `|---:` for right alignment
- Use `|:---:` for center alignment
- Minimum 3 dashes per column

### 3. Data Rows
```markdown
| Content | More content | Even more |
```
- Single pipe at start and end
- Single space padding around content
- Empty cells use appropriate placeholder or leave empty

## Sections Now Using Bullet/Checkbox Format (Not Tables)

The upstream STP template has moved most sections from table format to bullet/checkbox
format. The table-generator skill is **not** invoked for these sections:

- **Metadata & Tracking** — bullet list (6 items)
- **Section I.1 Requirement Review** — checkbox list (5 items)
- **Section I.2 Known Limitations** — bullet list
- **Section I.3 Technology Review** — checkbox list (5 items)
- **Section II.1 Out of Scope** — checkbox list with rationale
- **Section II.2 Test Strategy** — grouped checkbox list (13 items across 4 groups)
- **Section II.3 Test Environment** — bullet list (10 items)
- **Section II.4 Entry Criteria** — checkbox list
- **Section II.5 Risks** — checkbox list with sub-items (7 categories)
- **Section III Requirements Mapping** — bullet-based format

## Standard Table Templates

Tables are still used for custom/ad-hoc data presentation. Below are reference
templates for any remaining table needs.

### Generic Two-Column Table

```markdown
| Field | Details |
|:------|:--------|
| **Label** | Value |
| **Label** | Value |
```

### Generic Multi-Column Table

```markdown
| Column 1 | Column 2 | Column 3 |
|:---------|:---------|:---------|
| Content | Content | Content |
```

## Cell Content Rules

### Bold Headers
Use `**text**` for row headers in first column where appropriate

### Checkboxes
Use `[ ]` for unchecked, `[x]` for checked

### Empty Cells
- Use `N/A` with justification for explicitly not applicable
- Leave empty only if to be filled later
- Never leave required cells empty

### Links
Use markdown link format: `[Text](URL)`

### Lists in Cells
Use `<br>` for line breaks within cells, or separate with commas

## Validation

Before outputting, verify:
- [ ] All rows have same number of columns as header
- [ ] Alignment row present after header
- [ ] No trailing whitespace
- [ ] Proper pipe alignment
- [ ] Required row count met
