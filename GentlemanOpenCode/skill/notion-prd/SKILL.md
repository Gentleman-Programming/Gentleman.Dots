---
name: notion-prd
description: >
  Creates and manages PRDs (Product Requirement Documents) in Notion.
  Trigger: When user asks to create a PRD, document requirements, or prepare for engineering handoff.
license: Apache-2.0
metadata:
  author: prowler-cloud
  version: "7.1"
---

## When to Use

Use this skill when:
- Creating a PRD from a validated feature
- Documenting product requirements
- Preparing for engineering handoff
- Running PRD approval workflow

## Database Reference

**Database**: Docs (unified with RFCs, ADRs, and all engineering documentation)
**Collection ID**: `2e48202b-86e9-8089-b268-000b71c242df`
**URL**: https://www.notion.so/2e48202b86e980d79b0dfe7dcb7e7939

PRDs use `Category = PRD` in this unified database.

**Notion Template**: `2ef8202b-86e9-8082-9e48-cec4be4b16b3`
https://www.notion.so/prowlercom/PRD-Template-2ef8202b86e980829e48cec4be4b16b3

### Database Properties

| Property | Type | Values |
|----------|------|--------|
| Doc name | title | PRD title |
| Author | person | Product Owner |
| Status | status | Draft, Under Review, Published, Discarded, Outdated |
| Category | multi_select | **PRD** (also: RFC, ADR, Tech Spec, etc.) |
| Component | multi_select | Application, API, SDK, Platform, etc. |
| Summary | text | Brief one-liner |
| Reviewers | person | Stakeholders for approval |

### PRD-Specific Fields (in document body)

These fields go inside the document content, not as DB columns:

| Field | Description |
|-------|-------------|
| Feature Source | Link to source Feature |
| Design Link | Figma/prototype link |
| RFC Link | Link to RFC (after created by Engineering) |
| Epic Link | Link to Jira Epic (after sync) |
| Approval Date | When PRD was approved |

## Related Databases

| Entity | Database | Collection ID |
|--------|----------|---------------|
| Features (source) | Features | `927d90e5-adb9-4f3f-b3df-671bae57a514` |
| RFCs (created from PRD) | Docs | `2e48202b-86e9-8089-b268-000b71c242df` |

## Commands

```bash
/prd-create [feature-url]       # Generate PRD from validated feature
/prd-validate [prd-url]         # Run approval checklist
/prd-approve [prd-url]          # Mark as approved (Published)
```

## Workflow Position

```
Transcript --> Idea --> Feature --> PRD --> RFC --> Jira Epic
                                     ^
                               YOU ARE HERE
                         (Collaboration starts here!)
```

**IMPORTANT:** The PRD defines WHAT to build from a user/business perspective.
Engineering will create the RFC to define HOW to build it.

**DO NOT include technical implementation details** - that's Engineering's job in the RFC.

## Creating a PRD

```python
# 1. Fetch the source feature
feature = notion.fetch(feature_url)

# 2. Create PRD in the Docs database with Category = PRD
prd = notion.create_page(
    parent={"data_source_id": "2e48202b-86e9-8089-b268-000b71c242df"},
    properties={
        "Doc name": feature.title,  # No "PRD" suffix - Category field indicates type
        "Status": "Draft",
        "Category": "PRD",
        "Summary": f"Product requirements for {feature.title}"
    },
    content=prd_markdown  # Includes Feature Source link in body
)

# 3. Update the feature with PRD link
notion.update_page(
    page_id=feature_id,
    command="update_properties",
    properties={
        "Status": "PRD Created",
        "PRD Link": prd.url
    }
)
```

## PRD Title Convention

**Format:** `{Feature Name}` (NO "PRD" suffix - Category field already indicates it's a PRD)

**Examples:**
- `Findings View - Hierarchical Grouping`
- `Multi-tenant Support`
- `GovCloud Provider Integration`

## PRD Template

```markdown
# {Feature Name}

| Field | Value |
|-------|-------|
| **Version** | 1.0 |
| **Date** | {YYYY-MM-DD} |
| **Author** | {Product Owner} |
| **Status** | Draft |
| **Feature Source** | [Link to Feature]({feature-url}) |
| **Design Link** | _{Figma/prototype URL}_ |
| **RFC Link** | _{to be filled after RFC created}_ |
| **Epic Link** | _{to be filled after Jira sync}_ |
| **Approval Date** | _{to be filled when approved}_ |

---

## Summary

{2-3 paragraphs explaining:
- What we're building (user perspective)
- Why we're building it (business value)
- Who benefits from this}

---

## Problems

### Problem 1: {Title}
{Description of the problem from user perspective}
- **Who experiences this**: {user persona}
- **Impact**: {How it affects them}
- **Frequency**: {How often this occurs}
- **Current Workaround**: {What users do today}

### Problem 2: {Title}
{Description}
- **Who experiences this**: {user persona}
- **Impact**: {impact}
- **Frequency**: {frequency}
- **Current Workaround**: {workaround}

---

## Goals

### Primary Goals
1. **{Goal 1}**
   - {Measurable user outcome}
   
2. **{Goal 2}**
   - {Measurable user outcome}

### Secondary Goals
1. {Secondary goal 1}
2. {Secondary goal 2}

---

## Non-Goals

The following are explicitly **out of scope** for this initiative:

1. **{Non-goal 1}**
   - Reason: {why not now}
   
2. **{Non-goal 2}**
   - Reason: {why not now}

---

## Requirements

### {Functional Area 1}

| ID | Requirement | Description | Priority |
|----|-------------|-------------|----------|
| R-001 | {What user can do} | {Detailed description} | P0 |
| R-002 | {What user can do} | {Detailed description} | P1 |

### {Functional Area 2}

| ID | Requirement | Description | Priority |
|----|-------------|-------------|----------|
| R-003 | {What user can do} | {Detailed description} | P0 |
| R-004 | {What user can do} | {Detailed description} | P2 |

---

## User Stories

### Story 1: {Story Name}

**As a** {user type},
**I want to** {action},
**So that** {benefit}.

#### Acceptance Criteria
- [ ] {Criterion 1 - user perspective}
- [ ] {Criterion 2 - user perspective}
- [ ] {Criterion 3 - user perspective}

---

## Design

### Prototype Links
| Type | Link | Status |
|------|------|--------|
| Prototype | {Lovable/Figma URL} | Validated |

### Key Flows
{Description of main user flows}

---

## Risks

| Risk | Probability | Impact | Mitigation |
|------|-------------|--------|------------|
| {Business/User risk} | High/Med/Low | High/Med/Low | {mitigation} |

---

## Approvals

| Role | Name | Status | Date |
|------|------|--------|------|
| Product Owner | {name} | Pending | |
| Stakeholder | {name} | Pending | |
```

## PRD Validation Checklist

When running `/prd-validate [prd-url]`:

```markdown
## PRD Validation Checklist

**PRD**: {title}
**Status**: {status}

### Completeness
- [ ] Summary is clear and user-focused
- [ ] Problems are well-defined with user impact
- [ ] Goals are measurable (user outcomes)
- [ ] Non-goals explicitly stated
- [ ] Requirements have IDs and priorities
- [ ] User stories have acceptance criteria
- [ ] Design/prototype links are present
- [ ] Risks identified

### Quality
- [ ] No technical implementation details (that's for RFC)
- [ ] Requirements are from user perspective
- [ ] Priorities are justified
- [ ] Scope is appropriate for timeline

### Stakeholder Readiness
- [ ] Product Owner reviewed
- [ ] No open questions blocking approval

---

## Validation Result

**Ready for Approval?** {Yes / No}

### Missing Items
{List unchecked items}

### Next Step After Approval
Engineering creates RFC --> `/rfc-create {prd-url}`
```

## Status Flow

```
Draft --> Under Review --> Published --> (Engineering creates RFC)
               |
               v
          (Revisions)
               |
               v
            Draft
```

Note: "Published" = Approved in Docs database

## Best Practices

1. **User Focus**: Everything should be from user perspective
2. **No Technical Details**: Don't specify HOW - that's Engineering's job in RFC
3. **Be Specific**: Vague requirements lead to scope creep
4. **Prioritize Ruthlessly**: Not everything is P0
5. **Define Non-Goals**: What you won't do is as important as what you will
6. **Measurable Success**: If you can't measure it, you can't know if it worked

## What NOT to Include

- Database schemas
- API endpoint definitions
- Code examples
- Architecture decisions
- Implementation approaches
- Technical dependencies

These belong in the RFC, created by Engineering.

## Output Format

### After Creating PRD
```markdown
## PRD Created

**Title**: {title}
**Category**: PRD (no suffix needed in title)
**Status**: Draft
**Link**: [Open in Notion]({url})

### Generated From
- Feature: [{feature_title}]({feature_url})

### Next Steps
1. Review generated content
2. Add design/prototype links
3. Run `/prd-validate` when ready
4. Get approval from stakeholders (change to "Under Review")
5. **After approval** (Published): Engineering runs `/rfc-create {prd-url}`
```

### After Approval
```markdown
## PRD Approved

**Title**: {title}
**Status**: Published
**Approval Date**: {date}

### Next Steps
**Engineering**: Create RFC with `/rfc-create {url}`

Note: Jira Epic is created AFTER RFC approval, not after PRD approval.
```

## Keywords

prd, requirements, document, specification, product, handoff, approval
