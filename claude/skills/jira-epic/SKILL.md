---
name: jira-epic
description: >
  Creates Jira epics for large features following Prowler's standard format.
  Trigger: When user asks to create an epic, large feature, or multi-task initiative.
license: Apache-2.0
metadata:
  author: prowler-cloud
  version: "3.1"
---

## When to Use

Use this skill when creating Jira epics for:
- Large features spanning multiple components
- Major refactoring initiatives
- Features requiring API + UI + SDK work

**Note:** This skill can be used standalone OR is used by `notion-to-jira` when syncing PRDs/RFCs.

## CRITICAL: Two Modes of Operation

### Mode A: WITHOUT Codebase Access (Default)

When creating an Epic from PRD/RFC **without** access to Prowler repository:

- ✅ Include: Feature Overview, Requirements, Decisions, User Stories, Links
- ✅ Mark technical sections as "Pending Engineering Review"
- ❌ DO NOT invent file paths, service names, or component names
- ❌ DO NOT create fake Implementation Checklist with guessed paths

**Use this mode when:**
- Running from `prowler-workflows` repo
- Only have PRD/RFC content
- Engineering hasn't added technical design yet

### Mode B: WITH Codebase Access

When creating an Epic from **within the Prowler repository**:

- ✅ Explore codebase to find real file paths
- ✅ Include Technical Considerations with REAL paths
- ✅ Create Implementation Checklist referencing actual code

**Use this mode when:**
- Running from Prowler repo (`api/`, `ui/` directories exist)
- Engineering has completed technical design in RFC

```bash
# Verify you're in Prowler repo
fd "viewset" --type f --extension py api/
fd "component" --type f --extension tsx ui/
```

## Epic Template (Mode A - Without Codebase)

```markdown
h2. Feature Overview

{2-3 paragraph description from PRD}

*Figma:* [Design Link|{url}]
*PRD:* [PRD Title|{url}]
*RFC:* [RFC Title|{url}]

----

h2. Key Architectural Decisions

||Decision||Description||
|{Decision 1}|{Description from RFC}|
|{Decision 2}|{Description from RFC}|

----

h2. Requirements Summary

h3. {Section 1 from PRD}
* R-001: {Requirement} (P0)
* R-002: {Requirement} (P0)

h3. {Section 2 from PRD}
* R-011: {Requirement} (P0)

----

h2. User Stories

h3. Story 1: {Title}
*As a* {role}, *I want to* {what}, *so that* {why}.

h3. Story 2: {Title}
*As a* {role}, *I want to* {what}, *so that* {why}.

----

h2. Out of Scope (v1)

* {Item 1 from PRD/RFC}
* {Item 2 from PRD/RFC}

----

h2. Technical Considerations

{panel:title=Pending Engineering Review|borderColor=#FFAB00|bgColor=#FFFAE6}
This epic needs technical design with actual:
* Endpoint specifications
* Database models/queries
* Frontend component architecture

See RFC for open questions.
{panel}

----

h2. Implementation Tasks

*Stories to be created after RFC technical design is complete:*

# {Story 1 - derived from requirements}
# {Story 2 - derived from requirements}
```

## Epic Template (Mode B - With Codebase)

```markdown
h2. Feature Overview

{Description}

*Figma:* [Link|{url}]
*PRD:* [Link|{url}]
*RFC:* [Link|{url}]

----

h2. Technical Considerations

h3. Affected Areas (from codebase exploration)

*API:*
* {{api/src/backend/api/v1/{real_file}.py}}
* {{api/src/backend/models/{real_model}.py}}

*UI:*
* {{ui/components/{real_path}/{component}.tsx}}
* {{ui/lib/{real_service}.ts}}

h3. Existing Patterns to Follow
* {Pattern found in similar features}

----

h2. Implementation Checklist

- [ ] {Task based on REAL files/services}
- [ ] {Task based on REAL files/services}
```

## Epic Title Conventions

Format: `{Feature Name}` (clear, descriptive - NO prefix needed, Jira shows issue type)

**Examples:**
- `Findings View - Hierarchical Grouping`
- `Multi-tenant Support`
- `GovCloud Provider Integration`

> **Note:** Don't add `[EPIC]` prefix - Jira already shows the issue type with color/icon.

## Jira Hierarchy (CRITICAL)

**Correct hierarchy:** Epic → Story → Sub-task

| Issue Type | Parent Type | Link Field | Use For |
|------------|-------------|------------|---------|
| Epic | None | N/A | Large features, RFCs |
| Story | Epic | `customfield_10014` | User-facing functionality |
| Sub-task | Story | `parent` | Implementation work (API/UI) |

**WRONG:** Epic → Task (Tasks shouldn't be direct children of Epics)

**RIGHT:** Epic → Story → Sub-task

### Example Structure
```
PROWLER-379 (Epic: Findings View)
├── PROWLER-XXX (Story: Groups API Endpoint)
│   ├── Sub-task: [API] Create aggregation query
│   └── Sub-task: [API] Add pagination
└── PROWLER-XXX (Story: Hierarchical Tree UI)
    ├── Sub-task: [UI] Create TreeView component
    └── Sub-task: [UI] Implement lazy loading
```

## Jira MCP Integration

### Creating an Epic (Mode A - Without Codebase)

```json
{
  "project_key": "PROWLER",
  "summary": "Feature Name - Brief Description",
  "issue_type": "Epic",
  "additional_fields": {
    "customfield_10359": {"value": "API"},
    "customfield_10363": "h2. Feature Overview\n\n{content from PRD}\n\n*PRD:* [link]\n*RFC:* [link]\n\n----\n\nh2. Key Architectural Decisions\n\n||Decision||Description||\n|{dec1}|{desc1}|\n\n----\n\nh2. Technical Considerations\n\n{panel:title=Pending Engineering Review|borderColor=#FFAB00|bgColor=#FFFAE6}\nThis epic needs technical design.\n{panel}"
  }
}
```

### Key Fields Reference

| Field | Custom Field ID | Usage |
|-------|-----------------|-------|
| Epic Link | `customfield_10014` | String: `"PROWLER-741"` (for Stories under Epic) |
| Parent (Sub-task) | `parent` | String: `"PROWLER-787"` (for Sub-tasks under Story) |
| Team | `customfield_10359` | Object: `{"value": "API"}` or `{"value": "UI"}` |
| Work Item Description | `customfield_10363` | String with Jira Wiki markup |

### Team Field (REQUIRED)

`customfield_10359` options:
- `{"value": "UI"}` - Frontend epics
- `{"value": "API"}` - Backend epics
- `{"value": "SDK"}` - Prowler SDK epics

### Work Item Description Field (CRITICAL)

**IMPORTANT:** The `description` parameter in `jira_create_issue` populates a DIFFERENT field than what shows in the Jira UI!

- `description` parameter → Goes to the hidden "Description" field
- `customfield_10363` → Goes to **"Work Item Description"** field (what users see!)

**ALWAYS use `customfield_10363` in `additional_fields`**

**Use Jira Wiki markup (NOT Markdown):**
- `h2.` instead of `##`
- `*text*` for bold instead of `**text**`
- `* item` or `- item` for bullets
- `{code:python}...{code}` for code blocks
- `{{monospace}}` for inline code
- `{panel:title=X}...{panel}` for callouts
- `||header||header||` and `|cell|cell|` for tables

## Content Mapping: PRD/RFC → Epic

| PRD/RFC Section | Epic Section |
|-----------------|--------------|
| Summary | Feature Overview |
| Problems | Feature Overview (context) |
| Goals | Feature Overview |
| Requirements | Requirements Summary |
| User Stories | User Stories |
| Non-Goals | Out of Scope |
| Architectural Decisions (RFC) | Key Architectural Decisions |
| Design Links | Figma link in header |
| Open Questions (RFC) | Technical Considerations panel |

## Checklist Before Creating

### Mode A (Without Codebase)
1. ✅ Read PRD and RFC
2. ✅ Title is clear and descriptive (no `[EPIC]` prefix needed)
3. ✅ Feature Overview from PRD Summary + Problems
4. ✅ Key Decisions from RFC
5. ✅ Requirements Summary from PRD
6. ✅ User Stories from PRD
7. ✅ Out of Scope from PRD/RFC Non-Goals
8. ✅ Technical Considerations marked "Pending Engineering Review"
9. ✅ Links to PRD, RFC, Figma included
10. ✅ Using `customfield_10363` for description

### Mode B (With Codebase)
1. ✅ All Mode A items
2. ✅ Explored codebase for real file paths
3. ✅ Technical Considerations has REAL paths
4. ✅ Implementation Checklist references actual code

## Error Handling

### No Codebase Access (Expected for Mode A)
```
Note: Creating Epic in Mode A (without codebase access).

Technical sections marked as "Pending Engineering Review".
Engineering should update with real file paths and implementation details
after completing technical design in the RFC.
```

## Keywords

jira, epic, feature, initiative, prowler, prd, rfc, codebase
