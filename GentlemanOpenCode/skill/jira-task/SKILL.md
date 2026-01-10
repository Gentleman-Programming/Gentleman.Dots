---
name: jira-task
description: >
  Creates Jira tasks following Prowler's standard format.
  Trigger: When user asks to create a Jira task, ticket, or issue.
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

## When to Use

Use this skill when creating Jira tasks for:
- Bug reports
- Feature requests
- Refactoring tasks
- Documentation tasks

## Multi-Component Work: Split into Multiple Tasks

**IMPORTANT:** When work requires changes in multiple components (API, UI, SDK), create **separate tasks for each component** instead of one big task.

### Why Split?
- Different developers can work in parallel
- Easier to review and test
- Better tracking of progress
- API needs to be done before UI (dependency)

### Bug vs Feature: Different Structures

#### For BUGS: Create separate sibling tasks
Bugs are typically urgent fixes, so create independent tasks per component:

**Task 1 - API:**
- Title: `[BUG] Add aws_region field to AWS provider secrets (API)`
- Must be done first (UI depends on it)

**Task 2 - UI:**
- Title: `[BUG] Add region selector to AWS provider connection form (UI)`
- Blocked by API task

#### For FEATURES: Create parent + child tasks
Features need business context for stakeholders, so use a parent-child structure:

**Parent Task (for PM/Stakeholders):**
- Title: `[FEATURE] AWS GovCloud support`
- Contains: Feature overview, user story, acceptance criteria from USER perspective
- NO technical details
- Links to child tasks

**Child Task 1 - API:**
- Title: `[FEATURE] AWS GovCloud support (API)`
- Contains: Technical details, affected files, API-specific acceptance criteria
- Links to parent

**Child Task 2 - UI:**
- Title: `[FEATURE] AWS GovCloud support (UI)`
- Contains: Technical details, component paths, UI-specific acceptance criteria
- Links to parent, blocked by API task

### Parent Task Template (Features Only)

```markdown
## Description

{User-facing description of the feature - what problem does it solve?}

## User Story

As a {user type}, I want to {action} so that {benefit}.

## Acceptance Criteria (User Perspective)

- [ ] User can {do something}
- [ ] User sees {something}
- [ ] {Behavior from user's point of view}

## Out of Scope

- {What this feature does NOT include}

## Design

- Figma: {link if available}
- Screenshots/mockups if available

## Child Tasks

- [ ] `[FEATURE] {Feature name} (API)` - Backend implementation
- [ ] `[FEATURE] {Feature name} (UI)` - Frontend implementation

## Priority

{High/Medium/Low} ({business justification})
```

### Child Task Template (Features Only)

```markdown
## Description

Technical implementation of {feature name} for {component}.

## Parent Task

`[FEATURE] {Feature name}`

## Acceptance Criteria (Technical)

- [ ] {Technical requirement 1}
- [ ] {Technical requirement 2}

## Technical Notes

- Affected files:
  - `{file path 1}`
  - `{file path 2}`
- {Implementation hints}

## Testing

- [ ] {Test case 1}
- [ ] {Test case 2}

## Related Tasks

- Parent: `[FEATURE] {Feature name}`
- Blocked by: {if any}
- Blocks: {if any}
```

### Linking Tasks

In each task description, add:
```markdown
## Related Tasks
- Parent: [Parent task title/link] (for child tasks)
- Blocked by: [API task title/link]
- Blocks: [UI task title/link]
```

## Task Template

```markdown
## Description

{Brief explanation of the problem or feature request}

**Current State:**
- {What's happening now / What's broken}
- {Impact on users}

**Expected State:**
- {What should happen}
- {Desired behavior}

## Acceptance Criteria

- [ ] {Specific, testable requirement}
- [ ] {Another requirement}
- [ ] {Include both API and UI tasks if applicable}

## Technical Notes

- {Implementation hints}
- {Affected files with full paths}
- {Dependencies or related components}

## Testing

- [ ] {Test case 1}
- [ ] {Test case 2}
- [ ] {Include regression tests}

## Priority

{High/Medium/Low} ({justification})
```

## Title Conventions

Format: `[TYPE] Brief description (components)`

**Types:**
- `[BUG]` - Something broken that worked before
- `[FEATURE]` - New functionality
- `[ENHANCEMENT]` - Improvement to existing feature
- `[REFACTOR]` - Code restructure without behavior change
- `[DOCS]` - Documentation only
- `[CHORE]` - Maintenance, dependencies, CI/CD

**Components (when multiple affected):**
- `(API)` - Backend only
- `(UI)` - Frontend only
- `(SDK)` - Prowler SDK only
- `(API + UI)` - Both backend and frontend
- `(SDK + API)` - SDK and backend
- `(Full Stack)` - All components

**Examples:**
- `[BUG] AWS GovCloud accounts cannot connect - STS region hardcoded (API + UI)`
- `[FEATURE] Add dark mode toggle (UI)`
- `[REFACTOR] Migrate E2E tests to Page Object Model (UI)`
- `[ENHANCEMENT] Improve scan performance for large accounts (SDK)`

## Priority Guidelines

| Priority | Criteria |
|----------|----------|
| **Critical** | Production down, data loss, security vulnerability |
| **High** | Blocks users, no workaround, affects paid features |
| **Medium** | Has workaround, affects subset of users |
| **Low** | Nice to have, cosmetic, internal tooling |

## Affected Files Section

Always include full paths when known:

```markdown
## Technical Notes

- Affected files:
  - `api/src/backend/api/v1/serializers.py`
  - `ui/components/providers/workflow/forms/aws-credentials-form.tsx`
  - `prowler/providers/aws/config.py`
```

## Component-Specific Sections

### API Tasks
Include:
- Serializer changes
- View/ViewSet changes
- Migration requirements
- API spec regeneration needs

### UI Tasks
Include:
- Component paths
- Form validation changes
- State management impact
- Responsive design considerations

### SDK Tasks
Include:
- Provider affected
- Service affected
- Check changes
- Config changes

## Checklist Before Submitting

1. ✅ Title follows `[TYPE] description (components)` format
2. ✅ Description has Current/Expected State
3. ✅ Acceptance Criteria are specific and testable
4. ✅ Technical Notes include file paths
5. ✅ Testing section covers happy path + edge cases
6. ✅ Priority has justification
7. ✅ **Multi-component work is split into separate tasks**
8. ✅ **Titles are recommended for all tasks**

## Output Format

### For BUGS (sibling tasks):

```markdown
## Recommended Tasks

### Task 1: [BUG] {Description} (API)
{Full task content}

---

### Task 2: [BUG] {Description} (UI)
{Full task content}
```

### For FEATURES (parent + children):

```markdown
## Recommended Tasks

### Parent Task: [FEATURE] {Feature name}
{User-facing content, no technical details}

---

### Child Task 1: [FEATURE] {Feature name} (API)
{Technical content for API team}

---

### Child Task 2: [FEATURE] {Feature name} (UI)
{Technical content for UI team}
```

## Formatting Rules

**CRITICAL:** All output MUST be in Markdown format, ready to paste into Jira.

- Use `##` for main sections (Description, Acceptance Criteria, etc.)
- Use `**bold**` for emphasis
- Use `- [ ]` for checkboxes
- Use ``` for code blocks with language hints
- Use `backticks` for file paths, commands, and code references
- Use tables where appropriate
- Use `---` to separate multiple tasks

## Jira MCP Integration

**CRITICAL:** When creating tasks via MCP, use these exact parameters:

### Required Fields

```json
{
  "project_key": "PROWLER",
  "summary": "[TYPE] Task title (component)",
  "issue_type": "Task",
  "additional_fields": {
    "parent": "PROWLER-XXX",
    "customfield_10359": {"value": "UI"}
  }
}
```

### Team Field (REQUIRED)

The `customfield_10359` (Team) field is **REQUIRED**. Options:
- `"UI"` - Frontend tasks
- `"API"` - Backend tasks
- `"SDK"` - Prowler SDK tasks

### Work Item Description Field

**IMPORTANT:** The project uses `customfield_10363` (Work Item Description) instead of the standard `description` field for display in the UI.

**CRITICAL:** Use **Jira Wiki markup**, NOT Markdown:
- `h2.` instead of `##`
- `*text*` for bold instead of `**text**`
- `* item` for bullets (same)
- `** subitem` for nested bullets

After creating the issue, update the description with:

```json
{
  "customfield_10363": "h2. Description\n\n{content}\n\n*Current State:*\n* {problem 1}\n* {problem 2}\n\n*Expected State:*\n* {solution 1}\n* {solution 2}\n\nh2. Acceptance Criteria\n\n* {criteria 1}\n* {criteria 2}\n\nh2. Technical Notes\n\nPR: [{pr_url}]\n\nAffected files:\n* {file 1}\n* {file 2}\n\nh2. Testing\n\n* [ ] PR - Local environment\n** {test case 1}\n** {test case 2}\n* [ ] After merge in prowler - dev\n** {test case 3}"
}
```

### Common Epics

| Epic | Key | Use For |
|------|-----|---------|
| UI - Bugs & Improvements | PROWLER-193 | UI bugs, enhancements |
| API - Bugs / Improvements | PROWLER-XXX | API bugs, enhancements |
| LightHouse AI | PROWLER-594 | AI features |
| Technical Debt - UI | PROWLER-502 | Refactoring |

### Workflow Transitions

```
Backlog (10037) → To Do (14) → In Progress (11) → Done (21)
                → Blocked (10)
```

### MCP Commands Sequence

1. **Create issue:**
```
mcp__mcp-atlassian__jira_create_issue
```

2. **Update Work Item Description:**
```
mcp__mcp-atlassian__jira_update_issue with customfield_10363
```

3. **Assign and transition:**
```
mcp__mcp-atlassian__jira_update_issue (assignee)
mcp__mcp-atlassian__jira_transition_issue (status)
```

## Keywords
jira, task, ticket, issue, bug, feature, prowler
