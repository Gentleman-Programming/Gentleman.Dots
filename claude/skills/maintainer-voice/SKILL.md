---
name: maintainer-voice
description: >
  Write maintainer-grade async comments, PR reviews, issue responses, Jira updates,
  Slack messages, and Discord posts in a direct, evidence-based, contributor-friendly style.
  Trigger: When user asks to write a comment, reply, review, message, or update in
  any async communication channel — GitHub, Jira, Slack, Discord, or similar.
  Key phrases: "write a comment", "reply to this", "draft a message", "respond to",
  "leave feedback", "post an update", "announce", "close this", "approve this", "request changes".
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "2.1"
allowed-tools: Read, Edit, Write, Glob, Grep, Bash, WebFetch
---

## When to Use

Use this skill when writing any of these:
- GitHub issue/PR comments, review bodies, approval messages, close messages
- Jira ticket comments, status updates, blockers, or scope clarifications
- Slack replies in technical or project channels
- Discord announcements, changelogs, or contributor feedback
- Any async written message from a maintainer, tech lead, or senior contributor

---

## Default Format: GitHub PR/Issue Comments

**Short paragraph, not a checklist.** The canonical PR #101 example:

> "Thanks for the docs work. The content itself looks useful, but I am not going to merge a PR that skips the contribution contract. Open or link a real issue, get it approved, and add the correct type label. After that, push an updated PR body and we can review it on the merits instead of on preventable process failures."

What this demonstrates:
- Warm acknowledgment is **one phrase**, not a sentence of its own
- Verdict lands in sentence 1–2, not under a heading
- Required actions embedded as a short list **inside prose**, not as bullet checkboxes
- Ends with exactly one path forward
- Names the problem plainly ("preventable process failures") — no softening, no anger
- Total: 4 sentences, two logical beats, zero fluff
- **Zero em dash or en dash interruptions inside sentence prose**

**Length target for GitHub comments: 3–5 sentences, one or two short paragraphs.** Only go longer if the feedback is genuinely complex. When in doubt, cut.

---

## Core Style Rules

| Rule | Good | Bad |
|------|------|-----|
| **Compact by default** | 3–5 sentences, one paragraph | Bullet-heavy templates, multi-section review bodies |
| **Verdict in sentence 1–2** | "I'm not merging this until…" | Verdict buried after three sentences of context |
| **Evidence embedded in prose** | "skips the contribution contract" | "Missing: [ ] linked issue [ ] type label" |
| **One next action** | "Open an issue, get it approved, then update the PR body" | "You can either: A) … B) … C) …" |
| **Distinguish PR from idea** | "Closing the PR. The idea is still welcome." | "Closing" (ambiguous) |
| **Firm, not cold** | Name the problem accurately | Soften into vagueness |
| **No generic praise** | Skip it, or one short phrase tied to substance | "Thanks so much for contributing! We really appreciate the effort!" |
| **No corporate hedging** | "I am not going to merge" | "We may need to consider whether this aligns" |
| **No dash interruptions in prose** | Full sentences only, period-terminated | "The content looks useful — but this PR…" |

---

## Punctuation Rule: No Dash Interruptions in Prose

**Do not use em dashes, en dashes, or stylistic dash interruptions inside sentence prose.**

Bad:
> "The content looks useful — but this PR skips the contribution contract."
> "Closing this because it skips the process — no PR merges without a linked issue."

Good:
> "The content looks useful, but this PR skips the contribution contract."
> "Closing this because it skips the process. No PR merges without a linked issue."

The canonical PR #101 comment contains zero dash interruptions. That is the reference.

Use a dash only in two non-prose situations:
1. A changelog bullet list item: `- Improved startup time by 40%`
2. A label reference in backtick context: `` `status:approved` ``

**Never** use a dash to join two clauses inside a running sentence.

---

## Comment Structure

```
[Acknowledgment — max one phrase, tied to something specific, skip if nothing genuine to say]
[Verdict — what you're doing and the core reason, sentence 1 or 2]
[Evidence + required actions — embedded in prose, not a checklist]
[One next step — the single thing to do now]
```

For approvals, skip acknowledgment entirely. Get to "Merging." in two sentences.

---

## Intent Templates

### APPROVE (PR)

```
Scoped, tested, linked to the approved issue. CI green. Merging.
```

### APPROVE (issue)

```
This is clear, reproducible, and in scope. Adding `status:approved`. A PR linking this issue is welcome.
```

### REQUEST CHANGES — process gap

```
[Optional: one phrase acknowledging genuine value.] I'm not going to merge this until [specific gap]. [One sentence: what to do — open/link issue, add label, rebase, etc.] Once that's in, push an update and I'll re-review.
```

Example (issue-first):
```
The implementation looks solid, but this PR has no linked approved issue. Open or link one, get it approved, and add a type label. Then update the PR body and I can review the code on its merits.
```

### CLOSE PR — process violation (idea still valid)

```
Closing this because it skips the issue-first process. No PR merges without a `status:approved` issue linked. The underlying idea is worth pursuing: open an issue, get it approved, then re-file the PR with the link. This isn't a rejection of the concept.
```

### CLOSE PR — wrong direction

```
Closing this. The direction conflicts with [principle] because [one sentence why]. The code quality is solid but the approach isn't where this project is going. [Optional: what the right direction looks like in one sentence.]
```

### CLOSE ISSUE — belongs in Discussions

```
This is more a discussion topic than a concrete bug or feature. Please continue in [Discussions](link). Re-open if you can reproduce it as a specific bug with steps.
```

### CLOSE ISSUE — out of scope

```
Closing this. [One sentence: why it conflicts with project scope.] Re-open if there's strong community demand with concrete evidence.
```

### CLOSE ISSUE — duplicate

```
Duplicate of #N. Continuing there.
```

### NEEDS DESIGN

```
Good idea, but this touches [area] in a way that needs a design decision first. Propose an approach in the comments: what's the intended contract and behavior? We can get alignment before a PR makes sense.
```

### EXPLAIN POLICY

```
Quick note on process: [policy in one sentence]. [Why it exists — one sentence.] [What to do now — one action.]
```

### ASK FOR EVIDENCE

```
Before this moves forward I need to see [specific evidence]. Without [reproduction steps / failing test / benchmark], there's no way to evaluate the claim. Drop that here and I'll take another look.
```

### STALE / SILENCE

```
This [PR/issue] has been waiting on contributor feedback since [date]. Closing to keep the backlog clean. Re-open or re-file any time.
```

### APPROVE WITH MINOR NOTES

```
Good to merge. Two minor notes, not blocking: [note 1]; [note 2]. Merging now.
```

### ANNOUNCE RELEASE

```
Released v[version].

- [What changed — bullet 1]
- [What changed — bullet 2]

Install: `[command]`
Changelog: [link]
```

### REDIRECT TO DOCS

```
This is covered in [docs section](link). [One sentence key answer.] Closing. Re-open with the specific scenario if the docs don't cover your case.
```

---

## Channel Adaptations

### GitHub (issues + PRs) — DEFAULT
- Short prose paragraphs first. Only use bullets if listing 3+ discrete items that won't read well in a sentence.
- Reference issues with `#N`, labels with backticks
- End with a clear action verb: "Closing.", "Merging.", "Re-open after X."
- Target 3–5 sentences. Go to 8 only for genuinely complex technical feedback.

### Jira
- Skip Markdown bullets; use plain sentences
- State the ticket transition explicitly: "Moving to `In Progress`." / "Blocking on PROJ-42."
- `@mention` only when an action is required from that person
- 3–6 lines max

### Slack
- One paragraph, one action. No walls of text.
- Reply in thread. `@mention` only for required action.
- Use backticks for code. Never code blocks in main channel.

### Discord
- Announcements: heading, then short bullets, then install command
- Contributor feedback: DM over public callout for negative feedback
- Keep it casual, full sentences, no typos

---

## Anti-Patterns

| Anti-pattern | Why it's wrong |
|---|---|
| Bullet checklist for a 2-item request | Prose is faster to read and write. Bullets for 3+ items only. |
| "Thanks so much for contributing!" (opener) | Preamble. Lead with substance or skip it. |
| "This looks good but..." | Vague approval undermines the feedback that follows. |
| "We'll consider this" with no timeline or owner | Non-answer. Approve, reject, or schedule explicitly. |
| Closing without distinguishing PR from idea | Confusing. Say which one you're rejecting. |
| "Needs improvement" without specifics | Non-actionable. Name the exact gap. |
| "LGTM" alone | Adds no signal. State what you verified. |
| Multiple options at the end | One next step. Not "you can A, or B, or C". |
| Em dash or en dash inside sentence prose | Kills the clean sentence flow. Use a comma, period, or rewrite. |

---

## Tone

**Warm about the person, exact about the requirement.**

- Firm does not mean cold. You can close a PR warmly.
- Direct does not mean rude. A clear reason is more respectful than a soft hedge.
- Name problems accurately: "preventable process failure", "no linked issue", "CI failing on rebase" — not "some concerns".
- First-time contributors get one sentence of extra context: "This is the normal process here, not a judgment on the work itself."

---

## Decision Tree

```
What's the intent?
├── Approving something → APPROVE
├── Merging with minor notes → APPROVE WITH MINOR NOTES
├── Blocking until fixed → REQUEST CHANGES
├── Closing (process skip, idea valid) → CLOSE PR (process violation)
├── Closing (wrong direction) → CLOSE PR (wrong direction)
├── Closing stale item → STALE / SILENCE
├── Closing issue (out of scope) → CLOSE ISSUE (out of scope)
├── Redirecting channel → CLOSE ISSUE (belongs in Discussions)
├── Idea valid, design unclear → NEEDS DESIGN
├── Explaining policy → EXPLAIN POLICY
├── Need proof → ASK FOR EVIDENCE
├── Shipping → ANNOUNCE RELEASE
└── Pointing to docs → REDIRECT TO DOCS
```

---

## Before Writing: Verify

1. **Read the full thread.** Never respond to just the title.
2. **Verify technical claims** before citing them.
3. **Check if addressed elsewhere** — docs, another issue, merged PR — before closing.
4. **Know the distinction**: closing the PR vs. rejecting the idea vs. both.

---

## Resources

- **Templates**: See [assets/](assets/) for copy-ready comment bodies per intent
