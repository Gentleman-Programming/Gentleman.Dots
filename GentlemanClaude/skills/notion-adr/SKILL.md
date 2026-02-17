---
name: notion-adr
description: >
  Creates ADRs (Architecture Decision Records) from RFCs or technical discussions.
  Trigger: When user asks to create an ADR, document a technical decision, or formalize an architecture choice.
license: Apache-2.0
metadata:
  author: prowler-cloud
  version: "2.0"
---

## When to Use

Use this skill when:
- Formalizing technical decisions from an RFC
- Documenting architecture choices made in discussions
- Recording why a specific approach was chosen over alternatives

## Database Reference

**Database**: Docs (unified with all engineering documentation)
**Collection ID**: `2e48202b-86e9-8089-b268-000b71c242df`
**URL**: https://www.notion.so/2e48202b86e980d79b0dfe7dcb7e7939

ADRs use `Category = ADR` in this unified database.

**Notion Template**: `2ef8202b-86e9-80d1-91d3-c92476495cc3`
https://www.notion.so/prowlercom/ADR-Template-2ef8202b86e980d191d3c92476495cc3

### Database Properties

| Property | Type | Values |
|----------|------|--------|
| Doc name | title | ADR title |
| Author | person | Decision maker |
| Status | status | Draft, Under Review, Published, Discarded, Outdated |
| Category | multi_select | **ADR** (also: RFC, Tech Spec, Guide, etc.) |
| Component | multi_select | Application, API, SDK, Platform, etc. |
| Tech | multi_select | Python, Elixir, AI, AWS, etc. |
| Summary | text | One-line decision summary |
| Reviewers | person | Technical reviewers |

### ADR-Specific Fields (in document body)

These fields go inside the document content, not as DB columns:

| Field | Description |
|-------|-------------|
| Date | Decision date |
| RFC Source | Link to source RFC (if applicable) |
| Context | Why this decision was needed |
| Decision | The actual decision made |

## Commands

```bash
/adr-create [rfc-url]           # Extract decisions from RFC and create ADRs
/adr-create [title]             # Create standalone ADR
/adr-list                       # List all ADRs by status
```

## CRITICAL: Run from Prowler Repo

This skill should be run from within the Prowler codebase to:
- Reference actual code affected by the decision
- Link to real file paths
- Follow existing patterns

## ADR Template (Prowler Format)

```markdown
# ADR: {Decision Title}

| Field | Value |
|-------|-------|
| **Status** | Draft / Accepted / Deprecated / Superseded |
| **Date** | {YYYY-MM-DD} |
| **Author** | {Name} |
| **RFC Source** | [Link to RFC]({rfc-url}) _(if applicable)_ |

---

## Context

{What is the issue we're facing? Why do we need to make this decision?}

---

## Decision

{What is the decision we made? Be specific and concrete.}

---

## Alternatives Considered

### Alternative 1: {Name}
- **Pros:** {advantages}
- **Cons:** {disadvantages}
- **Why rejected:** {reason}

### Alternative 2: {Name}
- **Pros:** {advantages}
- **Cons:** {disadvantages}
- **Why rejected:** {reason}

---

## Consequences

### Positive
- {Positive consequence 1}
- {Positive consequence 2}

### Negative
- {Negative consequence 1}
- {Trade-off we accept}

---

## Implementation Notes

{Where in the codebase this affects - REAL paths from codebase exploration}

- Affected files: `{actual/path/to/file.py}`
- Related components: {component names from codebase}
```

## Extracting ADRs from RFC

When running `/adr-create [rfc-url]`:

1. **Read the RFC** - Look for the "Alternative Approaches Considered" section
2. **Identify decisions** - Each alternative comparison = potential ADR
3. **List candidates** - Ask user which to formalize

### Example

From RFC "Unified Findings Management":
```
Alternative Approaches Considered:
- GraphQL instead of REST --> Decision: REST for v1
- Client-side filtering --> Decision: Server-side filtering
```

**Extracted ADR candidates:**
1. ADR: REST API over GraphQL for Findings v1
2. ADR: Server-side Filtering for Large Datasets

## Creating ADR

```python
# Create ADR in Docs database with Category = ADR
adr = notion.create_page(
    parent={"data_source_id": "2e48202b-86e9-8089-b268-000b71c242df"},
    properties={
        "Doc name": f"ADR: {decision_title}",
        "Status": "Draft",
        "Category": "ADR",
        "Summary": decision_summary,
        "Component": component  # e.g., "API", "Application"
    },
    content=adr_markdown  # Includes RFC Source link in body
)
```

## Status Flow

```
Draft --> Under Review --> Published --> Outdated
                                            |
                                            v
                                      (Superseded by new ADR)
```

Note: 
- "Published" = Accepted decision
- "Outdated" = Superseded by newer decision

## Output Format

### After Creating ADR

```markdown
## ADR Created

**Title:** ADR: {title}
**Status:** Draft
**Category:** ADR
**Link:** [Open in Notion]({url})

### Source
- RFC: [{rfc_title}]({rfc_url}) _(if applicable)_

### Next Steps
1. Review with engineering team
2. Update status to "Under Review"
3. After approval, change to "Published"
4. Link from RFC's "Decision - ADR" section
```

### After Extracting from RFC

```markdown
## ADR Candidates from RFC

**RFC:** [{title}]({url})

### Decisions Found

| # | Decision | Alternatives | Recommended ADR Title |
|---|----------|--------------|----------------------|
| 1 | REST over GraphQL | GraphQL | ADR: REST API for Findings v1 |
| 2 | Server-side filtering | Client-side | ADR: Server-side Filtering Architecture |

### Actions
- Reply with numbers to create (e.g., "1, 2" or "all")
- Or "skip" to not create ADRs
```

## Best Practices

1. **One decision per ADR** - Keep them focused
2. **Document alternatives** - Show what you didn't choose and why
3. **Link to RFC** - Maintain traceability in the body
4. **Update when superseded** - Don't delete old ADRs, mark as Outdated
5. **Reference real code** - When running from repo, include actual file paths
6. **Use Component tag** - Tag which component this decision affects

## Keywords

adr, architecture, decision, record, technical, choice, alternative, rfc
