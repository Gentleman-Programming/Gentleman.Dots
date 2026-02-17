---
name: notion-to-jira
description: >
  Bridge between Notion (RFC) and Jira Epics/Tasks. Syncs approved RFCs to Jira.
  Trigger: When user asks to sync RFC to Jira, create epics from RFC, or bridge Notion to Jira.
license: Apache-2.0
metadata:
  author: prowler-cloud
  version: "4.0"
---

## When to Use

Use this skill when:
- Syncing an approved RFC to Jira
- Creating Epics and Tasks from technical design
- Updating Notion with Jira links

## Database Reference

**RFCs are in**: Docs database with `Category = RFC`
**Collection ID**: `2e48202b-86e9-8089-b268-000b71c242df`
**URL**: https://www.notion.so/2e48202b86e980d79b0dfe7dcb7e7939

## Commands

```bash
/rfc-to-jira [rfc-url]              # Full sync: create Epic + Tasks from RFC
/sync-epic [rfc-url]                # Create only the Epic
/sync-tasks [epic-key] [rfc-url]    # Create Tasks for existing Epic
```

## CRITICAL: Prerequisites

Before creating Epic/Tasks, you MUST:

### 1. Read the RFC
```python
rfc = notion.fetch(rfc_url)
# Extract: Summary, Problem, Proposed Solution, Sign Off status
```

### 2. Verify RFC is Published (Approved)
The RFC must be in **Published** status (this is "Approved" in Docs DB).
If not, stop and tell the user.

### 3. Explore the Codebase
**DO NOT invent file names, service names, or component names.**

```bash
# Find real file names and patterns
fd "viewset" --type f --extension py api/
fd "component" --type f --extension tsx ui/
rg "class.*Service" api/src/
rg "def.*task" api/src/
```

Use what you find to:
- Reference REAL file paths in tasks
- Use ACTUAL service/model names from the code
- Follow EXISTING patterns in the codebase

## Workflow Position

```
Transcript --> Feature --> PRD --> RFC --> Jira Epic
                                    ^           ^
                                 Source      Output
```

**Flow:**
- PRD = WHAT to build (Product owns)
- RFC = HOW to build (Engineering owns, in Docs DB with Category=RFC)
- Jira = Implementation tracking (based on RFC + actual codebase)

## Mapping RFC to Epic

| RFC Section | Epic Field |
|-------------|------------|
| Title | `[EPIC] {RFC Title without "RFC:"}` |
| Summary | Feature Overview |
| Problem/Motivation | Context section |
| Proposed Solution | Technical Considerations |

### Epic Template

```markdown
# {Feature Name}

**RFC:** {rfc_notion_url}
**PRD:** {prd_notion_url from RFC body}

## Feature Overview

{From RFC Summary}

## Context

{From RFC Problem/Motivation}

## Technical Approach

{From RFC Proposed Solution - HIGH LEVEL only}

## Implementation Checklist

{Tasks derived from RFC + codebase exploration}

- [ ] {Task based on REAL files/services found in codebase}
- [ ] {Task based on REAL files/services found in codebase}
```

## Mapping RFC to Tasks

**CRITICAL:** Tasks must reference REAL code, not invented names.

### Before Creating Tasks

1. **Explore the codebase** to find:
   - Actual file paths that need modification
   - Real model/service/component names
   - Existing patterns to follow

2. **Use RFC Proposed Solution** as the WHAT, codebase as the WHERE

### Task Template

```markdown
## Description

{Task description based on RFC + actual code location}

**RFC Reference:** {rfc_url}

## Acceptance Criteria

- [ ] {Criterion from RFC}
- [ ] {Criterion from RFC}

## Technical Notes

**Affected files (REAL paths from codebase):**
- `{actual/file/path.py}`
- `{actual/component/path.tsx}`

**Existing patterns to follow:**
- {Pattern found in codebase}

## Testing

- [ ] Unit tests
- [ ] Integration tests

## Related

- Epic: {epic_key}
- RFC: {rfc_url}
```

### Task Title Format

```
[FEATURE] {Description} ({Component})
```

Components: `API`, `UI`, `SDK`

## Jira Hierarchy (CRITICAL)

**Correct hierarchy:** Epic --> Story --> Sub-task

| Issue Type | Parent Type | Link Field | Use For |
|------------|-------------|------------|---------|
| Epic | None | N/A | Large features, RFCs |
| Story | Epic | `customfield_10014` | User-facing functionality |
| Sub-task | Story | `parent` | Implementation work (API/UI) |

**WRONG:** Epic --> Task (Tasks shouldn't be direct children of Epics)

**RIGHT:** Epic --> Story --> Sub-task

## Jira Integration

### Create Epic

```python
epic = jira.create_issue(
    project_key="PROWLER",
    summary=f"[EPIC] {rfc_title.replace('RFC: ', '')}",
    issue_type="Epic",
    additional_fields={
        "customfield_10359": {"value": component},  # Team field
        "customfield_10363": epic_description_wiki  # Work Item Description
    }
)
```

### Create Stories under Epic

```python
story = jira.create_issue(
    project_key="PROWLER",
    summary=f"[STORY] {story_title}",
    issue_type="Story",
    additional_fields={
        "customfield_10014": epic.key,  # Epic Link
        "customfield_10359": {"value": component},
        "customfield_10363": story_description_wiki
    }
)
```

### Create Sub-tasks under Story

```python
subtask = jira.create_issue(
    project_key="PROWLER",
    summary=f"[{component}] {task_title}",
    issue_type="Sub-task",
    additional_fields={
        "parent": story.key,  # Parent Story
        "customfield_10359": {"value": component},
        "customfield_10363": task_description_wiki
    }
)
```

### Update Notion RFC

After creating Epic, update the RFC's Epic Link in the document body:

```python
# RFC doesn't have Epic Link as a column - it's in the body
# Find and update the Epic Link field in the document content
notion.update_page(
    page_id=rfc_id,
    command="replace_content_range",
    selection_with_ellipsis="| **Epic Link** |..._ |",
    new_str=f"| **Epic Link** | [PROWLER-{epic.key}](https://prowlerpro.atlassian.net/browse/{epic.key}) |"
)
```

## Jira Field References

```yaml
jira:
  project_key: "PROWLER"
  epic_type: "Epic"
  story_type: "Story"
  subtask_type: "Sub-task"
  team_field: "customfield_10359"       # REQUIRED: UI, API, or SDK
  epic_link_field: "customfield_10014"  # For linking Story to Epic
  description_field: "customfield_10363" # Work Item Description (visible in UI)
```

### Jira Wiki Format

**CRITICAL:** Use Jira Wiki markup, NOT Markdown:
- `h2.` instead of `##`
- `*text*` for bold instead of `**text**`
- `* item` for bullets

## Error Handling

### RFC Not Published (Approved)
```
Error: RFC must be in Published status before syncing to Jira.
Current Status: {status}

Next Steps:
1. Get RFC approved by reviewers
2. Change status to "Published"
3. Then run /rfc-to-jira {url}
```

### Can't Find Codebase
```
Error: Need access to Prowler codebase to create accurate tasks.

This skill must be run from within the Prowler repository so I can:
- Find actual file paths
- Identify existing patterns
- Reference real service/model names
```

## Output Format

```markdown
## RFC Synced to Jira

**RFC:** [{title}]({rfc_url})
**Status:** Synced

### Epic Created
- **Key:** PROWLER-XXX
- **Title:** [EPIC] {title}
- **Link:** [Open in Jira]({jira_url})

### Stories Created ({count})

| Key | Title | Component |
|-----|-------|-----------|
| PROWLER-001 | [STORY] {title} | API |
| PROWLER-002 | [STORY] {title} | UI |

### Sub-tasks Created ({count})

| Key | Title | Parent Story | Based On |
|-----|-------|--------------|----------|
| PROWLER-003 | [API] {title} | PROWLER-001 | {file found in codebase} |
| PROWLER-004 | [UI] {title} | PROWLER-002 | {component found in codebase} |

### Notion Updated
- RFC Epic Link updated in document body

### Next Steps
1. Review created Epic, Stories, and Sub-tasks in Jira
2. Assign tasks to team members
3. Add to sprint
```

## Best Practices

1. **RFC Must Be Published**: Never sync draft RFCs
2. **Explore Codebase First**: Don't invent file/service names
3. **Real Paths Only**: Tasks must reference actual code locations
4. **Follow Existing Patterns**: Look at how similar features were implemented
5. **Keep Traceability**: Link RFC --> Epic --> Story --> Sub-task
6. **Use Correct Hierarchy**: Epic --> Story --> Sub-task (NOT Epic --> Task)

## Keywords

notion, jira, sync, epic, story, task, rfc, bridge, codebase
