---
name: notion-product-brain
description: >
  Manages Product's ideation space in Notion - ideas and features.
  Trigger: When user asks to create ideas, features, or validate product concepts.
license: Apache-2.0
metadata:
  author: prowler-cloud
  version: "4.0"
---

## When to Use

Use this skill when:
- Creating new product ideas
- Documenting features for validation
- Running validation checklists before PRD
- Managing the product ideation pipeline

## Database References

### Ideas Database
**Collection ID**: `9164a9f9-5f77-4d5d-8a5e-6e521daf3fd4`
**URL**: https://www.notion.so/a26784038c68488cb6db72ed8b756662

| Property | Type | Values |
|----------|------|--------|
| Idea | title | - |
| Author | rich_text | - |
| Source | url | Link to transcript/source |
| Status | select | Raw, Exploring, Validated, Discarded, Promoted to Feature |
| Priority | select | High, Medium, Low |

### Features Database
**Collection ID**: `927d90e5-adb9-4f3f-b3df-671bae57a514`
**URL**: https://www.notion.so/cc1d61460fac41478ab0d4cc8bdff126

| Property | Type | Values |
|----------|------|--------|
| Feature | title | - |
| Problem | rich_text | - |
| Hypothesis | rich_text | - |
| Priority | select | High, Medium, Low |
| Status | select | Ideating, Validating, Ready for PRD, PRD Created |
| Prototype Link | url | - |
| PRD Link | url | - |

## Commands

```bash
/product-idea [title]           # Create quick idea
/product-feature [title]        # Create structured feature
/product-validate [feature-url] # Run validation checklist
/product-list                   # List features by status
```

## Creating Items

### Creating an Idea

```python
notion.create_page(
    parent={"data_source_id": "9164a9f9-5f77-4d5d-8a5e-6e521daf3fd4"},
    properties={
        "Idea": idea_title,
        "Status": "Raw",
        "Author": author_name,
        "Source": source_url  # Optional: link to transcript
    },
    content=idea_markdown
)
```

### Creating a Feature

```python
notion.create_page(
    parent={"data_source_id": "927d90e5-adb9-4f3f-b3df-671bae57a514"},
    properties={
        "Feature": feature_title,
        "Status": "Ideating",
        "Problem": problem_summary,
        "Priority": "Medium"
    },
    content=feature_markdown
)
```

## Idea Template

Quick capture for raw ideas, brainstorming, early concepts.

```markdown
# {Idea Title}

## The Problem
{What user pain or opportunity does this address?}

## Proposed Solution
{High-level idea of how to solve it - from user perspective}

## Source
{Where did this idea come from? Customer feedback, competitor, internal, transcript}

## Initial Thoughts
- {Thought 1}
- {Thought 2}

## Questions to Answer
- [ ] {Question 1}
- [ ] {Question 2}
```

## Feature Template

More structured format for ideas that are being actively explored.

```markdown
# {Feature Title}

## The Problem
{Detailed description of the user pain point}

### Who experiences this?
- {User persona 1}
- {User persona 2}

### When do they experience it?
{Context and frequency}

### Current workarounds
- {Workaround 1}
- {Workaround 2}

## Proposed Solution
{How we plan to solve this - from user perspective}

### Key capabilities
- {Capability 1 - what user can do}
- {Capability 2 - what user can do}

## Hypothesis
**If** {we build X},
**Then** {users will Y},
**Because** {reason}.

### Success Metric
{How we'll know it worked - user behavior, not technical metrics}

## Prototypes
| Version | Link | Date | Feedback Summary |
|---------|------|------|------------------|
| v1 | {link} | {date} | {summary} |

## Validation Status
- [ ] Problem validated with users
- [ ] Solution concept tested
- [ ] Prototype feedback collected
- [ ] Scope defined (in/out)

## Collaboration Log
| Date | Person | Topic | Outcome |
|------|--------|-------|---------|
| {date} | {name} | {topic} | {what was decided} |
```

## Validation Checklist Command

When running `/product-validate [feature-url]`:

```markdown
## Feature Validation Checklist

**Feature**: {title}
**Current Status**: {status}

### Problem Validation
- [ ] Problem is clearly articulated
- [ ] Target users identified
- [ ] Pain point frequency/severity understood
- [ ] Current workarounds documented

### Solution Validation  
- [ ] Solution addresses the core problem
- [ ] Key capabilities defined (user perspective)
- [ ] Prototype created and tested
- [ ] User feedback collected and positive

### Scope Definition
- [ ] In-scope items listed
- [ ] Out-of-scope items explicitly stated
- [ ] MVP vs future phases clear

### Business Validation
- [ ] Aligns with company goals
- [ ] Priority justified
- [ ] Success metrics defined (user behavior)

---

## Validation Result

**Ready for PRD?** {Yes / No}

### Missing Items
{List of unchecked items that need attention}

### Recommended Next Steps
1. {Step 1}
2. {Step 2}
```

## Status Flow

```
Idea: Raw --> Exploring --> Validated --> Discarded
                               |
                               v
                    Promoted to Feature
                               |
                               v
Feature: Ideating --> Validating --> Ready for PRD --> PRD Created
                                            |
                                            v
PRD: Draft --> In Review --> Approved --> In Development --> Completed
```

## Best Practices

### For Ideas
- Capture quickly, refine later
- Link to source (transcript, customer call)
- Don't over-structure initially
- Review weekly to promote or discard

### For Features
- Be specific about the problem
- Include real user quotes if available
- Keep prototypes linked and updated
- Focus on USER outcomes, not technical solutions

### For Validation
- Don't skip the checklist
- Prototype before committing
- Define success metrics upfront (user behavior)

## Output Format

### After Creating Idea
```markdown
## Idea Created

**Title**: {title}
**Status**: Raw
**Link**: [Open in Notion]({url})

### Next Steps
1. Add more details to the problem statement
2. Identify source/evidence
3. Review in next product sync
```

### After Creating Feature
```markdown
## Feature Created

**Title**: {title}
**Status**: Ideating
**Link**: [Open in Notion]({url})

### Next Steps
1. Create hypothesis
2. Build prototype (Lovable/Figma)
3. Schedule validation sessions
4. Run `/product-validate` when ready
```

### After Validation
```markdown
## Validation Complete

**Feature**: {title}
**Result**: {Ready for PRD / Needs Work}

### Summary
- {X}/{Y} validation items passed
- Missing: {list of missing items}

### Next Steps
{If ready: Run `/prd-create {url}`}
{If not ready: List specific actions needed}
```

## Keywords

product, idea, feature, validate, notion, brainstorm, hypothesis, prototype
