---
name: jira-task
description: >
  Creates Jira tasks following Prowler's standard format.
  Trigger: When user asks to create a Jira task, ticket, or issue.
license: Apache-2.0
metadata:
  author: prowler-cloud
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

### How to Split

When you identify multi-component work:

1. **Create a parent/epic task** (optional) for tracking the overall feature
2. **Create individual tasks** for each component
3. **Always recommend titles** for all tasks

### Example: GovCloud Support

Instead of one task `[BUG] AWS GovCloud cannot connect (API + UI)`, create:

**Task 1 - API:**
- Title: `[BUG] Add aws_region field to AWS provider secrets (API)`
- Must be done first (UI depends on it)

**Task 2 - UI:**
- Title: `[BUG] Add region selector to AWS provider connection form (UI)`
- Blocked by API task

### Linking Tasks

In each task description, add:
```markdown
## Related Tasks
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

When creating tasks, always output:

```markdown
## Recommended Tasks

### Task 1: [Full title here]
{Full task content}

---

### Task 2: [Full title here]
{Full task content}

---

(repeat for each task)
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

## Keywords
jira, task, ticket, issue, bug, feature, prowler
