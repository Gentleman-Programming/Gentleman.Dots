---
name: transcript-processor
description: >
  Processes meeting transcripts and generates structured output based on type.
  Trigger: When user asks to process a transcript, meeting notes, or recording.
license: Apache-2.0
metadata:
  author: prowler-cloud
  version: "5.0"
---

## When to Use

Use this skill when:
- Processing meeting transcripts
- Extracting insights from recordings
- Generating summaries and action items
- Creating features or ideas from discussions

## Database Reference

**Database**: Transcripts
**Collection ID**: `c4264969-980c-4ac7-9252-b03cec3b48d5`
**URL**: https://www.notion.so/83f6cfd8b7aa43d198209124af939ad4

### Database Properties

| Property | Type | Values |
|----------|------|--------|
| Transcript | title | Meeting title |
| Date | date | Meeting date |
| Type | select | Product Discovery, Technical, Design Review, Sprint Planning, Incident Review, Customer Call |
| Attendees | text | Comma-separated names |
| Status | select | Raw, Processing, Processed |
| Summary | text | Brief summary |
| Action Items | number | Count of action items |

## Related Databases

When creating entities FROM transcripts, use these:

| Entity | Database | Collection ID |
|--------|----------|---------------|
| Ideas | Ideas | `9164a9f9-5f77-4d5d-8a5e-6e521daf3fd4` |
| Features | Features | `927d90e5-adb9-4f3f-b3df-671bae57a514` |

## Commands

```bash
/process-transcript [notion-url]                    # Process existing transcript in Notion
/process-transcript [notion-url] --type=product     # Force specific type
/upload-transcript [file-path]                      # Upload local file to Notion first, then process
```

## IMPORTANT: Transcript Must Be in Notion First

The transcript MUST exist in the Notion Transcripts database before processing. This ensures:
- **Traceability**: All derived entities (Features, Ideas) link back to source
- **Audit trail**: We know where insights came from
- **Single source of truth**: No orphaned transcripts

### If user provides a local file path:

1. **Ask for metadata**: Meeting title, date, attendees, type
2. **Create entry in Transcripts DB** with Status = "Raw"
3. **Copy content** to the Notion page
4. **Then process** using the standard flow

```python
# Step 1: Create transcript in Notion
transcript = notion.create_page(
    parent={"data_source_id": "c4264969-980c-4ac7-9252-b03cec3b48d5"},
    properties={
        "Transcript": meeting_title,
        "Status": "Raw",
        "Type": detected_type,  # Product Discovery, Technical, etc.
        "date:Date:start": meeting_date,
        "Attendees": attendees
    },
    content=transcript_content
)

# Step 2: Process the transcript (standard flow)
# ... processing logic ...

# Step 3: Update status and summary
notion.update_page(
    transcript.id,
    command="update_properties",
    properties={
        "Status": "Processed",
        "Summary": executive_summary,
        "Action Items": action_item_count
    }
)
```

## CRITICAL: Privacy & Content Filtering

### What to INCLUDE
- Business discussions
- Product requirements
- User feedback and pain points
- Action items related to work
- Decisions made

### What to EXCLUDE
- Personal conversations (family, health, hobbies)
- Off-topic discussions (sports, weather, weekend plans)
- Private information (salaries, personal issues)
- Jokes, banter, casual chat

### Speaker Anonymization

1. **Named speakers**: Use names as-is (e.g., "Sarah", "Mike")
2. **Unnamed speakers**: Use generic labels: `Speaker 1`, `Speaker 2`
3. **Notion AI transcripts**: Normalize `[Speaker]` or timestamps

## Supported Types

| Type | Auto-detect Keywords | Output |
|------|---------------------|--------|
| `product` | user, customer, feature, pain point | Features, Ideas |
| `planning` | sprint, task, deadline, blocker | Action Items |
| `customer` | demo, feedback, pricing | Pain Points, Requests |

## Processing Flow

```
1. Verify transcript exists in Notion Transcripts DB
   - If local file provided: upload first (see above)
   - If not in DB: ERROR - must be uploaded first
2. Fetch transcript content from Notion
3. FILTER: Remove personal/off-topic content
4. ANONYMIZE: Handle unnamed speakers
5. Detect or use provided type
6. Generate type-specific output (business-focused)
7. Create related entities (Features, Ideas) WITH LINK to transcript
8. Update transcript status to "Processed"
9. Add links to created entities in transcript page
```

## Output Templates

### Product Discovery

```markdown
## Summary
{2-3 sentence executive summary}

## User Insights
| Insight | Quote/Evidence | Impact |
|---------|---------------|--------|
| {insight} | "{quote}" | High/Medium/Low |

## Pain Points Identified
1. **{Pain Point}**: {description}
   - Current workaround: {if any}
   - Frequency: {how often}

## Feature Ideas
| Idea | Priority | Related Pain Point |
|------|----------|-------------------|
| {idea} | High/Med/Low | {pain point} |

## Action Items
- [ ] {task} - {owner if mentioned}
```

### Planning Meeting

```markdown
## Summary
{Sprint goal and key focus areas}

## Commitments
| Task | Owner | Priority |
|------|-------|----------|
| {task} | {person} | P0/P1/P2 |

## Blockers Identified
| Blocker | Dependency | Owner |
|---------|------------|-------|
| {blocker} | {what's needed} | {who} |

## Action Items
- [ ] {task} - {owner}
```

### Customer Call

```markdown
## Summary
{Customer context, key takeaways}

## Customer Context
- **Use Case**: {how they use/want to use product}
- **Stage**: Prospect/Trial/Customer

## Pain Points Discussed
| Pain Point | Current Solution | Urgency |
|------------|-----------------|---------|
| {pain} | {workaround} | High/Med/Low |

## Feature Requests
| Request | Business Value | Priority |
|---------|---------------|----------|
| {feature} | {why they need it} | High/Med/Low |

## Next Steps
- [ ] {task} - {owner}
```

## Creating Entities from Transcript

### Creating an Idea

```python
# Always link back to the source transcript!
notion.create_page(
    parent={"data_source_id": "9164a9f9-5f77-4d5d-8a5e-6e521daf3fd4"},  # Ideas DB
    properties={
        "Idea": idea_title,
        "Status": "Raw",
        "Source": transcript_url  # Link back to transcript
    },
    content=idea_content
)
```

### Creating a Feature

```python
notion.create_page(
    parent={"data_source_id": "927d90e5-adb9-4f3f-b3df-671bae57a514"},  # Features DB
    properties={
        "Feature": feature_title,
        "Status": "Ideating",
        "Problem": pain_point
    },
    content=feature_content + f"\n\n## Source\n[{transcript_title}]({transcript_url})"
)
```

### Updating Transcript with Created Entities

After creating Features/Ideas, update the transcript page:

```python
notion.update_page(
    transcript_id,
    command="insert_content_after",
    selection="## Action Items...",  # or end of content
    new_str="""

## Entities Created

| Type | Title | Link |
|------|-------|------|
| Feature | {title} | [Open]({url}) |
| Idea | {title} | [Open]({url}) |
"""
)
```

## Output Format

```markdown
## Transcript Processed

**Type**: {detected/specified type}
**Content Filtered**: Yes (removed {N} off-topic segments)
**Speakers**: {Named / Anonymized}

### Summary
{Executive summary}

### Key Outputs
- Created {N} action items
- Created {N} features/ideas

### Entities Created
| Type | Title | Link |
|------|-------|------|
| Feature | {title} | [Link]({url}) |

### Next Steps
{What the user should do next}
```

## Keywords

transcript, meeting, notes, process, summary, action items, feature, privacy, filter
