---
name: jira-task
description: >
  Manages Jira Stories and Sub-tasks following Prowler's standard format.
  Trigger: When user asks to create, update, or modify a Jira task, ticket, story, or issue.
license: Apache-2.0
metadata:
  author: prowler-cloud
  version: "4.0"
---

## When to Use

Use this skill when **creating or updating**:
- **Stories** under an Epic (user-facing functionality)
- **Sub-tasks** under a Story (implementation work)
- **Standalone Tasks** (bugs, chores not part of a feature)

### Common Scenarios

| Action | Example |
|--------|---------|
| Create new ticket | "Create a story for the filters feature" |
| Update existing ticket | "Update PROWLER-123 with implementation details" |
| Add codebase context | "Add technical notes to PROWLER-456 based on the codebase" |
| Enrich from RFC | "Update the story with info from the RFC" |
| Add sub-tasks | "Break down PROWLER-789 into sub-tasks" |

**Note:** This skill can be used standalone from any point in the workflow.

## CRITICAL: Two Modes of Operation

### Mode A: WITHOUT Codebase Access (Default)

When creating Stories/Sub-tasks **without** access to Prowler repository:

- ✅ Include: WHAT to do (requirements, acceptance criteria)
- ✅ Include: WHY (user value, business context)
- ❌ DO NOT include: HOW to do it (implementation details)
- ❌ DO NOT invent file paths, service names, or component names

**Stories describe WHAT the user needs, not HOW to build it.**

### Mode B: WITH Codebase Access

When creating Stories/Sub-tasks from **within the Prowler repository**:

- ✅ Include: WHAT + WHY
- ✅ Include: HOW (implementation hints, real file paths)
- ✅ Reference actual code patterns and locations

```bash
# Verify you're in Prowler repo
fd "viewset" --type f --extension py api/
fd "component" --type f --extension tsx ui/
```

## Jira Hierarchy (CRITICAL)

**Correct hierarchy:** Epic → Story → Sub-task

| Issue Type | Parent Type | Link Field | Use For |
|------------|-------------|------------|---------|
| Epic | None | N/A | Large features |
| Story | Epic | `customfield_10014` | User-facing functionality |
| Sub-task | Story | `parent` | Implementation work (API/UI) |

**NEVER create Tasks directly under Epics. Use Stories.**

### Example Structure
```
PROWLER-379 (Epic: Findings View)
├── PROWLER-XXX (Story: Groups API Endpoint)
│   ├── PROWLER-XXX (Sub-task: Create aggregation query) [Team: API]
│   └── PROWLER-XXX (Sub-task: Add pagination support) [Team: API]
└── PROWLER-XXX (Story: Hierarchical Tree UI)
    ├── PROWLER-XXX (Sub-task: Create TreeView component) [Team: UI]
    └── PROWLER-XXX (Sub-task: Implement lazy loading) [Team: UI]
```

---

## Chained PR Pattern for Sub-tasks

When structuring Sub-tasks for a Story, organize them as **chained PRs** that:
1. Can be reviewed independently
2. Allow API and UI teams to work in parallel
3. Merge in sequence to a **feature branch**

### Git Branch Strategy

```
main
 └── feature/PROWLER-XXX-story-name  ← Feature branch (base for all PRs)
      ├── PR #1: API contract     → merges to feature branch
      ├── PR #2: Endpoint A       → merges to feature branch
      ├── PR #3: Endpoint B       → merges to feature branch
      ├── PR #4: UI Components    → merges to feature branch
      ├── PR #5: Server Actions   → merges to feature branch
      ├── PR #6: Implement view   → merges to feature branch
      └── FINAL: feature branch   → merges to main (when all sub-tasks done)
```

**Key point:** All chained PRs target the **feature branch**, NOT main. The final "integration" is simply merging the feature branch to main.

### Sub-tasks Structure

```
Story: {User-facing feature}
├── 1. Define API contract/model (API) ← PREREQUISITE for both tracks
│
├── API Track (sequential)
│   ├── 2. Endpoint A (API) ← depends on 1
│   └── 3. Endpoint B (API) ← depends on 2
│
├── UI Track (parallel with API, uses mocks)
│   ├── 4. Reusable components (UI) ← depends on 1
│   ├── 5. Server Actions with mocks (UI) ← depends on 1
│   └── 6. Implement view + integrate real API (UI) ← depends on 3, 4, 5
│
└── [No separate integration sub-task - just merge feature branch to main]
```

### Dependency Diagram

> **IMPORTANT: ASCII Box Alignment**
> The first line of the box (`┌───┐`) must have **9 spaces** of indentation to align with the text inside. This ensures proper rendering in Jira's monospace font.
>
> ```
> ❌ Wrong:    ┌─────────────┐
>              │ 885. Filter │
>
> ✅ Correct:         ┌─────────────┐
>                     │ 885. Filter │
> ```

```
                    ┌─────────────────┐
                    │  1. API Model   │  ← Define contract FIRST
                    └────────┬────────┘
                             │
         ┌───────────────────┼───────────────────┐
         │                   │                   │
         ▼                   ▼                   ▼
┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐
│ 2. Endpoint A   │ │ 4. Components   │ │ 5. Server       │
│    (API)        │ │    (UI)         │ │    Actions (UI) │
└────────┬────────┘ └────────┬────────┘ └────────┬────────┘
         │                   │                   │
         ▼                   └─────────┬─────────┘
┌─────────────────┐                    │
│ 3. Endpoint B   │                    ▼
│    (API)        │          ┌─────────────────────────┐
└────────┬────────┘          │ 6. Implement view +     │
         │                   │    integrate real API   │
         └───────────────────┴─────────────────────────┘
                             │
                    All PRs merge to feature branch
                             │
                             ▼
                    ┌─────────────────┐
                    │  Merge feature  │
                    │  branch → main  │
                    └─────────────────┘
```

### Key Principles

1. **Feature branch as base** - Create `feature/PROWLER-XXX` from main, all PRs target this branch
2. **API contract first** - Both teams need the interface defined upfront
3. **UI works with mocks** - Don't block UI on API completion
4. **Each sub-task = 1 PR** - Auto-contained, reviewable independently
5. **Clear dependencies** - Create Jira issue links (blocks/is blocked by) - NOT just text in description
6. **Final merge to main** - Not a sub-task, just merge the feature branch when done
7. **Include diagram in Story** - Add ASCII dependency diagram to Story description (see box alignment rules above)

### Real Example: PROWLER-774

**Branch:** `feature/PROWLER-774-hierarchical-view`

| # | Sub-task | Team | Depends On | PR Target |
|---|----------|------|------------|-----------|
| 1 | Define API contract for Groups/Resources | API | - | feature branch |
| 2 | GET /findings/groups endpoint | API | 1 | feature branch |
| 3 | GET /findings/groups/{id}/resources | API | 2 | feature branch |
| 4 | ExpandableRow reusable component | UI | 1 | feature branch |
| 5 | Server Actions for Findings Groups | UI | 1 | feature branch |
| 6 | Implement hierarchical view + integrate API | UI | 3, 4, 5 | feature branch |
| - | *Merge feature branch to main* | - | all done | main |

### Creating Dependency Links (IMPORTANT)

Dependencies must be **Jira issue links**, not just text in the description. Use `jira_create_issue_link`:

```json
{
  "link_type": "Blocks",
  "inward_issue_key": "PROWLER-876",   // This issue...
  "outward_issue_key": "PROWLER-877"   // ...blocks this one
}
```

**For PROWLER-774 example, create these links:**

| From (blocks) | To (is blocked by) |
|---------------|-------------------|
| 876 (API Model) | 877 (GET /groups) |
| 876 (API Model) | 879 (ExpandableRow) |
| 876 (API Model) | 880 (Server Actions) |
| 877 (GET /groups) | 878 (GET /resources) |
| 878 (GET /resources) | 882 (Integrate) |
| 879 (ExpandableRow) | 881 (Implement view) |
| 880 (Server Actions) | 881 (Implement view) |
| 881 (Implement view) | 882 (Integrate) |

This creates a proper dependency graph visible in Jira's "Dependencies" section.

---

## Story Template (Mode A - Without Codebase)

**Title:** `{User-facing functionality}`

```markdown
h2. User Story

*As a* {role},
*I want to* {what},
*So that* {why/value}.

----

h2. Context

{Background from PRD/RFC. Why this matters.}

----

h2. Acceptance Criteria

* [ ] {Testable requirement 1 - WHAT, not HOW}
* [ ] {Testable requirement 2}
* [ ] {Testable requirement 3}

----

h2. Design Reference

*Figma:* [Design Link|{url}]
*PRD Requirement:* R-001, R-002

----

h2. Out of Scope

* {What this story does NOT include}

----

h2. Technical Notes

{panel:title=For Engineering|borderColor=#0052CC|bgColor=#DEEBFF}
Implementation details to be determined by Engineering.
See RFC for architectural decisions.
{panel}

----

h2. Dependencies

* Blocked by: {story/task if any}
* Blocks: {story/task if any}
```

## Story Template (Mode B - With Codebase)

**Title:** `{User-facing functionality}`

```markdown
h2. User Story

*As a* {role},
*I want to* {what},
*So that* {why/value}.

----

h2. Acceptance Criteria

* [ ] {Testable requirement 1}
* [ ] {Testable requirement 2}

----

h2. Technical Notes

h3. Affected Files
* {{api/src/backend/api/v1/{real_file}.py}}
* {{ui/components/{real_path}/{component}.tsx}}

h3. Implementation Hints
* Follow pattern in {{api/src/backend/api/v1/resources.py}}
* Use existing {{TreeView}} component from {{ui/components/ui/tree-view.tsx}}

h3. Suggested Sub-tasks
# [API] {Specific backend task}
# [UI] {Specific frontend task}
```

---

## Sub-task Template (Mode A - Without Codebase)

**Title:** `{What to implement}`

```markdown
h2. Description

{What this sub-task delivers - NOT how to implement it}

----

h2. Acceptance Criteria

* [ ] {Verifiable outcome 1}
* [ ] {Verifiable outcome 2}

----

h2. Parent Story Context

This sub-task is part of: [STORY-XXX|{url}]

Story Goal: {Brief reminder of what the parent story achieves}

----

h2. Validation

* [ ] {How to verify this works}
* [ ] {Test scenario}
```

## Sub-task Template (Mode B - With Codebase)

**Title:** `{What to implement}`

```markdown
h2. Description

{What this sub-task delivers}

----

h2. Acceptance Criteria

* [ ] {Verifiable outcome 1}
* [ ] {Verifiable outcome 2}

----

h2. Implementation

h3. Files to Modify
* {{api/src/backend/api/v1/{file}.py}} - {what to change}
* {{api/src/backend/models/{model}.py}} - {what to add}

h3. Code Pattern
Follow the pattern in:
{code:python}
# Example from existing code
class ExistingViewSet(viewsets.ModelViewSet):
    ...
{code}

h3. Steps
# {Step 1 with specific file reference}
# {Step 2}
# {Step 3}

----

h2. Testing

* [ ] Unit test: {{api/tests/test_{feature}.py}}
* [ ] E2E test: {{ui/e2e/{feature}.spec.ts}}
```

---

## Title Conventions

> **Note:** Don't add `[STORY]`, `[EPIC]`, `[BUG]` prefixes - Jira already shows the issue type with color/icon. Use the **Team field** (`customfield_10359`) to indicate API/UI/SDK.

### Stories
Format: `{User-facing functionality}` (clear, descriptive)

**Examples:**
- `Groups API - Aggregate findings by check`
- `Hierarchical Tree View for Findings`
- `Resource Detail Panel with Remediation`

### Sub-tasks
Format: `{What to implement}` (set Team field for API/UI/SDK)

**Examples:**
- `Create findings aggregation endpoint` (Team: API)
- `Add pagination to groups endpoint` (Team: API)
- `Create TreeView component` (Team: UI)
- `Implement lazy loading for resources` (Team: UI)

### Standalone Tasks (Outside Stories)

Not all work fits in the Epic → Story → Sub-task hierarchy. Use standalone Tasks for:

| Task Type | When to Use | Example |
|-----------|-------------|---------|
| **Investigation/Spike** | Research before committing to implementation | "Investigate aggregation strategies for findings" |
| **Technical Debt** | Refactoring, cleanup, optimization | "Refactor FindingsTable to use TanStack Table" |
| **Performance** | Optimization work | "Optimize findings query with database indexes" |
| **UI Polish** | Visual improvements not tied to features | "Improve loading states across tables" |
| **DevOps/Infra** | CI/CD, deployment, tooling | "Add Sentry error tracking to UI" |
| **Bug Fix** | Isolated bugs not part of a feature | "Fix pagination reset on filter change" |

**Issue Type:** Task (or Bug)
**Parent:** None (link to Epic with `customfield_10014` if loosely related)

#### Standalone Task Template

```jira
h2. Description

{What needs to be done and why}

----

h2. Acceptance Criteria

* [ ] {Verifiable outcome 1}
* [ ] {Verifiable outcome 2}

----

h2. Context

*Why now:* {Why this is needed at this time}
*Related to:* [PROWLER-XXX|{url}] (if any)

----

h2. Notes

{Any technical considerations or constraints}
```

**Examples:**
- `Investigate aggregation strategies for findings` (Type: Task, Team: API)
- `AWS GovCloud accounts cannot connect` (Type: Bug, Team: API)
- `Improve loading states across tables` (Type: Task, Team: UI)

---

## Jira MCP Integration

### Creating a Story under Epic

```json
{
  "project_key": "PROWLER",
  "summary": "Story title - clear description",
  "issue_type": "Story",
  "additional_fields": {
    "customfield_10014": "PROWLER-379",
    "customfield_10359": {"value": "API"},
    "customfield_10363": "h2. User Story\n\n*As a* security engineer,\n*I want to* see findings grouped by check,\n*So that* I can prioritize remediation.\n\n----\n\nh2. Acceptance Criteria\n\n* [ ] Groups are aggregated by check_id\n* [ ] Each group shows highest severity\n* [ ] Pagination with 50 items default"
  }
}
```

### Creating a Sub-task under Story

```json
{
  "project_key": "PROWLER",
  "summary": "Create aggregation endpoint",
  "issue_type": "Sub-task",
  "additional_fields": {
    "parent": "PROWLER-XXX",
    "customfield_10359": {"value": "API"},
    "customfield_10363": "h2. Description\n\nCreate API endpoint that aggregates findings by check_id.\n\n----\n\nh2. Acceptance Criteria\n\n* [ ] Endpoint returns grouped findings\n* [ ] Includes provider breakdown per group\n* [ ] Supports pagination"
  }
}
```

### Updating an Existing Ticket

Use `jira_update_issue` to modify any field on an existing ticket.

#### Update Description (Most Common)

```json
{
  "issue_key": "PROWLER-123",
  "fields": {
    "customfield_10363": "h2. User Story\n\n*As a* security engineer...\n\n----\n\nh2. Technical Notes\n\n{Added after codebase exploration}"
  }
}
```

#### Update Summary (Title)

```json
{
  "issue_key": "PROWLER-123",
  "fields": {
    "summary": "New title for the ticket"
  }
}
```

#### Update Multiple Fields

```json
{
  "issue_key": "PROWLER-123",
  "fields": {
    "summary": "Updated title",
    "customfield_10363": "h2. Updated Description\n\n...",
    "customfield_10359": {"value": "UI"}
  }
}
```

#### Common Update Scenarios

| Scenario | What to Update |
|----------|----------------|
| Add technical notes from codebase | `customfield_10363` - append Technical Notes section |
| Add implementation details from RFC | `customfield_10363` - append relevant RFC decisions |
| Break down into sub-tasks | Create sub-tasks with `parent: "PROWLER-XXX"` |
| Change team assignment | `customfield_10359` |
| Clarify scope | `customfield_10363` - update Acceptance Criteria |

#### Workflow: Enriching a Ticket

1. **Fetch current ticket** with `jira_get_issue` to see existing content
2. **Explore codebase** (Mode B) to find relevant files/patterns
3. **Read RFC/PRD** if available for context
4. **Update ticket** preserving existing content + adding new sections

```python
# Pseudo-workflow
current = jira_get_issue("PROWLER-123")
existing_desc = current.fields.customfield_10363

new_desc = existing_desc + """

----

h2. Technical Notes

h3. Affected Files
* {{api/src/backend/api/v1/findings.py}}

h3. Implementation Hints
* Follow pattern in existing filters
"""

jira_update_issue("PROWLER-123", {"customfield_10363": new_desc})
```

### Key Fields Reference

| Field | Custom Field ID | Usage |
|-------|-----------------|-------|
| Epic Link | `customfield_10014` | String: `"PROWLER-379"` (for Stories under Epic) |
| Parent (Sub-task) | `parent` | String: `"PROWLER-XXX"` (for Sub-tasks under Story) |
| Team | `customfield_10359` | Object: `{"value": "API"}` or `{"value": "UI"}` |
| Work Item Description | `customfield_10363` | String with Jira Wiki markup |

### Team Field (REQUIRED)

`customfield_10359` options:
- `{"value": "UI"}` - Frontend work
- `{"value": "API"}` - Backend work
- `{"value": "SDK"}` - Prowler SDK work

### Work Item Description Field (CRITICAL)

**ALWAYS use `customfield_10363` in `additional_fields`**

The `description` parameter goes to a hidden field. Only `customfield_10363` shows in Jira UI.

**Use Jira Wiki markup (NOT Markdown):**
- `h2.` instead of `##`
- `*text*` for bold instead of `**text**`
- `* item` or `- item` for bullets
- `* [ ]` for checkboxes
- `{code:python}...{code}` for code blocks
- `{{monospace}}` for inline code
- `{panel:title=X}...{panel}` for callouts
- `----` for horizontal rule

---

## Using PRD/RFC to Create or Enrich Tickets

PRDs and RFCs are valuable sources for both **creating new tickets** and **enriching existing ones**.

### Creating New: From PRD Requirements → Stories

| PRD Section | Story Derivation |
|-------------|------------------|
| Requirements R-001 to R-010 (Main View) | Story: "Findings List with Grouping" |
| Requirements R-011 to R-016 (Drill-Down) | Story: "Resource Drill-Down" |
| Requirements R-017 to R-021 (Detail Panel) | Story: "Resource Detail Panel" |

### Creating New: From RFC Decisions → Sub-tasks

| RFC Decision | Sub-task Derivation |
|--------------|---------------------|
| "Two endpoints architecture" | [API] Groups endpoint, [API] Resources endpoint |
| "3-level tree view" | [UI] TreeView component |
| "Lazy loading" | [UI] Implement lazy loading |

### Enriching Existing Tickets from PRD/RFC

When updating an existing ticket with PRD/RFC context:

| Source | What to Extract | Where to Add |
|--------|-----------------|--------------|
| PRD Requirements | Acceptance criteria, user value | h2. Acceptance Criteria |
| PRD Figma links | Design references | h2. Design Reference |
| RFC Architecture decisions | Technical approach | h2. Technical Notes |
| RFC API contracts | Endpoint specs, models | h2. Implementation (Mode B) |
| RFC Dependencies | Blocking/blocked relationships | Create issue links |

**Example: Enriching PROWLER-123 from RFC**

```
# Before (ticket created without RFC context)
h2. User Story
As a user, I want to filter findings...

# After (enriched with RFC decisions)
h2. User Story
As a user, I want to filter findings...

----

h2. Technical Notes

*From RFC-005:*
* Filter state managed via URL params (not local state)
* Debounced API calls (300ms)
* Server-side filtering for performance

h2. Dependencies
* Requires: Filter API endpoint (PROWLER-456)
```

---

## Checklists

### Creating Stories (Mode A - Without Codebase)
1. ✅ Title is clear and descriptive (no `[STORY]` prefix needed)
2. ✅ User Story has role, what, why
3. ✅ Acceptance Criteria are WHAT (testable outcomes)
4. ✅ NO implementation details (no file paths, no "how")
5. ✅ Links to Figma, PRD requirements
6. ✅ Out of Scope defined
7. ✅ Epic link set (`customfield_10014`)
8. ✅ Team field set (`customfield_10359`)

### Creating Stories (Mode B - With Codebase)
1. ✅ All Mode A items
2. ✅ Technical Notes with REAL file paths
3. ✅ Implementation hints from actual codebase
4. ✅ Suggested sub-tasks based on code exploration

### Creating Sub-tasks (Mode A - Without Codebase)
1. ✅ Title is clear and descriptive (Team field indicates API/UI/SDK)
2. ✅ Description says WHAT, not HOW
3. ✅ Acceptance Criteria are verifiable
4. ✅ Parent story link set (`parent`)
5. ✅ Team field matches component

### Creating Sub-tasks (Mode B - With Codebase)
1. ✅ All Mode A items
2. ✅ Files to Modify with REAL paths
3. ✅ Code patterns from actual codebase
4. ✅ Testing file paths

### Updating Existing Tickets
1. ✅ Fetched current ticket content first (`jira_get_issue`)
2. ✅ Preserved existing content (don't overwrite, append)
3. ✅ Used same Jira Wiki markup format
4. ✅ Added clear section headers for new content
5. ✅ If adding technical notes: verified paths exist in codebase
6. ✅ If adding from RFC: cited RFC reference
7. ✅ Created issue links for dependencies (not just text)

---

## Error Handling

### No Codebase Access (Expected for Mode A)
```
Note: Creating/updating in Mode A (without codebase access).

Content describes WHAT to deliver, not HOW to implement.
Engineering should add implementation details based on
actual codebase patterns.
```

## Keywords

jira, task, story, subtask, ticket, issue, bug, feature, codebase, user story, update, modify, enrich
