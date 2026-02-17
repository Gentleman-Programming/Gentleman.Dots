---
name: notion-rfc
description: >
  Creates and manages RFCs (Request for Comments) for technical design review.
  Trigger: When user asks to create an RFC, technical design doc, or engineering proposal.
license: Apache-2.0
metadata:
  author: prowler-cloud
  version: "5.1"
---

## When to Use

Use this skill when:
- Creating an RFC from an approved PRD
- Documenting technical design decisions
- Proposing architectural changes
- Getting engineering review before implementation

## Database Reference

**Database**: Docs (unified with all engineering documentation)
**Collection ID**: `2e48202b-86e9-8089-b268-000b71c242df`
**URL**: https://www.notion.so/2e48202b86e980d79b0dfe7dcb7e7939

RFCs use `Category = RFC` in this unified database.

**Notion Template**: `2e78202b-86e9-801b-aa63-f94af9749f59`
https://www.notion.so/prowlercom/RFC-Template-2e78202b86e9801baa63f94af9749f59

### Database Properties

| Property | Type | Values |
|----------|------|--------|
| Doc name | title | RFC title |
| Author | person | Engineering lead |
| Status | status | Draft, Under Review, Published, Discarded, Outdated |
| Category | multi_select | **RFC** (also: ADR, Tech Spec, Guide, etc.) |
| Component | multi_select | Application, API, SDK, Platform, etc. |
| Tech | multi_select | Python, Elixir, AI, AWS, etc. |
| Summary | text | Brief one-liner |
| Reviewers | person | Technical reviewers |

### RFC-Specific Fields (in document body)

These fields go inside the document content, not as DB columns:

| Field | Description |
|-------|-------------|
| PRD Source | Link to source PRD |
| Epic Link | Link to Jira Epic (after sync) |
| Complexity | Low, Medium, High |
| Date | RFC creation date |

## Commands

```bash
/rfc-create [prd-url]           # Generate RFC from approved PRD
/rfc-validate [rfc-url]         # Run technical review checklist
/rfc-approve [rfc-url]          # Mark as approved by engineering
/rfc-to-jira [rfc-url]          # Create Epic/Tasks (uses notion-to-jira skill)
```

## Workflow Position

```
Transcript --> Feature --> PRD --> RFC --> Jira Epic
                                    ^
                              YOU ARE HERE
```

The RFC is the **engineering's response to the PRD**. It's where the technical team:
- Validates feasibility
- Proposes implementation approach
- Identifies technical risks
- Gets sign-off before implementation

## RFC Templates

## RFC Title Convention

**Format:** `{Feature Name}` (NO "RFC:" prefix - Category field already indicates it's an RFC)

**Examples:**
- `Findings View - Hierarchical Grouping`
- `Multi-tenant Support`
- `GovCloud Provider Integration`

### Template A: WITHOUT Repo Access (Architectural Decisions Only)

Use this when creating RFC from PRD/meeting notes WITHOUT access to Prowler codebase.

```markdown
# {Feature Name}

<callout icon="⚠️" color="yellow_bg">
  **Status: Needs Technical Design** - This RFC contains architectural decisions 
  from the discovery meeting. Technical implementation details (endpoints, models, 
  queries) to be added by Engineering with access to the Prowler codebase.
</callout>

| Field | Value |
|-------|-------|
| **Author** | {Name} |
| **Date** | {YYYY-MM-DD} |
| **Status** | Draft - Needs Technical Design |
| **PRD Source** | [Link to PRD]({prd-url}) |
| **Complexity** | Low / Medium / High |
| **Reviewers** | TBD (Backend Lead, Frontend Lead) |
| **Epic Link** | _{to be filled after Jira sync}_ |

---

## Summary

{What we're building and key architectural decisions - NO invented code}

---

## Architectural Decisions (from Discovery Meeting)

### Decision 1: {Decision Name}
**Decision**: {What was decided}
**Rationale**: {Why this approach}

### Decision 2: {Decision Name}
**Decision**: {What was decided}
**Rationale**: {Why this approach}

{Use ASCII diagrams for visual structure - NO fake JSON/code}

---

## Out of Scope (v1)

{What was explicitly excluded and why}

---

## Open Questions for Engineering

- [ ] {Question that needs codebase access to answer}
- [ ] {Implementation detail to be decided}

---

## Next Steps

1. Engineering review of architectural decisions
2. Technical design with actual endpoint specs, models, queries
3. Sign-off before implementation

---

## Sign Off

| Reviewer | Role | Status |
|----------|------|--------|
| TBD | Backend Lead | Pending |
| TBD | Frontend Lead | Pending |
| {Product Owner} | Product | Pending |
```

### Template B: WITH Repo Access (Full Technical Design)

Use this when you have access to Prowler codebase and can reference real code.

```markdown
# {Feature Name}

| Field | Value |
|-------|-------|
| **Author** | {Name} |
| **Date** | {YYYY-MM-DD} |
| **Status** | Draft |
| **PRD Source** | [Link to PRD]({prd-url}) |
| **Complexity** | Low / Medium / High |
| **Reviewers** | {Name 1}, {Name 2} |
| **Epic Link** | _{to be filled after Jira sync}_ |

---

## Summary

{Executive summary: 2-3 paragraphs explaining what we're building, why this approach, and high-level overview}

---

## Problem / Motivation

{Detailed description of:
- Current state and its problems (reference actual code)
- Pain points
- Business impact}

---

## Proposed Solution(s)

{Technical proposal including:
- Architecture approach
- Key components (reference existing patterns)
- Data flow (if relevant)
- API design (actual endpoints)
- Database changes (actual models)}

### API Changes

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/v1/{resource}` | GET | {description} |

### Database Changes

```sql
-- Actual schema changes based on existing models
```

### Component Diagram

```
[Component A] --> [Component B] --> [Database]
```

---

## Alternative Approaches Considered

### Option 1: {Name}
- **Pros:** {advantages}
- **Cons:** {disadvantages}
- **Why rejected:** {reason}

---

## Open Questions

- [ ] {Open question 1}
- [ ] {Open question 2}

---

## Sign Off

| Reviewer | Role | Status | Date |
|----------|------|--------|------|
| {Name 1} | {Role} | Pending | |
| {Name 2} | {Role} | Pending | |

---

## Decision - ADR

{Link to ADR if a formal architecture decision is needed, or "N/A"}
```

## Creating RFC from PRD

```python
# 1. Fetch PRD details
prd = notion.fetch(prd_url)

# 2. Create RFC page in Docs database
rfc = notion.create_page(
    parent={"data_source_id": "2e48202b-86e9-8089-b268-000b71c242df"},
    properties={
        "Doc name": prd.title,  # No "RFC:" prefix - Category field indicates type
        "Status": "Draft",
        "Category": "RFC",
        "Summary": f"Technical design for {feature_name}"
    },
    content=rfc_markdown  # Includes PRD Source link in body
)

# 3. Update PRD with RFC link (in content, not property)
# PRD doesn't have RFC Link column - update the content instead
```

## Status Flow

```
Draft --> Under Review --> Published --> (to Jira)
               |
               v
          Discarded
```

Note: Docs database uses "Published" instead of "Approved"

## Validation Checklist

When running `/rfc-validate [rfc-url]`:

```markdown
## RFC Validation Checklist

**RFC**: {title}
**Status**: {status}

### Completeness
- [ ] Summary clearly explains what and why
- [ ] PRD Source link is present
- [ ] Problem/Motivation section is complete
- [ ] Proposed Solution is concrete (not vague)
- [ ] API/Database changes documented (if applicable)
- [ ] Alternative approaches documented
- [ ] Open questions listed

### Quality
- [ ] Technical approach is sound
- [ ] Complexity assessment is realistic
- [ ] No major risks unaddressed

### Review Readiness
- [ ] Reviewers assigned in Sign Off table
- [ ] All open questions can be resolved in review

---

## Validation Result

**Ready for Review?** {Yes / No}

### Missing Items
{List unchecked items}

### Next Steps
{Change status to "Under Review" and notify reviewers}
```

## Cross-Reference Updates

| Action | Update Location |
|--------|-----------------|
| Create RFC from PRD | RFC body contains PRD Source link |
| Approve RFC | Change Status to "Published" |
| Create Jira Epic | Update RFC body with Epic Link |

## Output Format

### After Creating RFC

```markdown
## RFC Created

**Title:** {feature_name}
**Status:** Draft
**Category:** RFC (no prefix needed in title)
**Link:** [Open in Notion]({url})

### Source PRD
- [{prd_title}]({prd_url})

### Next Steps
1. Complete technical design in Proposed Solution
2. Add reviewers to Sign Off table
3. Change status to "Under Review"
4. Run `/rfc-validate` when ready for approval
5. After approval: Run `/rfc-to-jira` to create Epic
```

## CRITICAL: Technical Content Rules

### When to Include Technical Implementation Details

| Context | What to Include | What to AVOID |
|---------|-----------------|---------------|
| **WITH access to Prowler repo** | Real endpoints, actual models, existing patterns from codebase | N/A - use real code |
| **WITHOUT access to Prowler repo** | Only architectural decisions from PRD/meetings | Invented endpoints, fake models, made-up code |

### Without Repo Access - RFC Contains ONLY:

1. **Architectural decisions** (from meetings/PRD)
   - Grouping strategy (e.g., "hierarchical by check_id")
   - Data flow decisions (e.g., "two endpoints, lazy loading")
   - What's in/out of scope

2. **High-level technical approach**
   - Visual diagrams (ASCII)
   - Decision rationale
   - Tradeoffs discussed

3. **Open questions for Engineering**
   - Questions that need codebase access to answer
   - Implementation details to be decided

4. **Callout warning**:
   ```markdown
   <callout icon="⚠️" color="yellow_bg">
     **Status: Needs Technical Design** - This RFC contains architectural decisions 
     from the discovery meeting. Technical implementation details (endpoints, models, 
     queries) to be added by Engineering with access to the Prowler codebase.
   </callout>
   ```

### DO NOT INVENT:
- ❌ Specific endpoint paths (`/api/v1/findings/groups`)
- ❌ JSON response structures
- ❌ Django models or fields
- ❌ SQL queries
- ❌ Python/code snippets
- ❌ Database schema changes

### With Repo Access - RFC Can Include:

- ✅ Actual existing endpoints to modify
- ✅ Real model names and relationships
- ✅ Existing patterns from codebase
- ✅ Concrete implementation proposal

## Best Practices

1. **Keep it simple** - The template is intentionally lightweight
2. **Be concrete ONLY with repo access** - Without repo, stay at architecture level
3. **Document alternatives** - Show what you didn't choose and why
4. **Resolve questions** - All questions should be answered before approval
5. **Get sign-off** - Don't skip the review process
6. **Link bidirectionally** - PRD links to RFC, RFC links to PRD
7. **Mark status clearly** - "Needs Technical Design" if created without repo access

## Keywords

rfc, technical design, engineering, architecture, implementation, proposal, review
