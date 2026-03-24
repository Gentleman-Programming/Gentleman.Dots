# Maintainer Voice — Copy-Ready Templates

Short prose paragraphs by default for GitHub. Fill in `[brackets]` before posting.
Target: 3–5 sentences. Bullets only when listing 3+ discrete items.
No em dash or en dash interruptions inside sentence prose.

---

## GitHub: Approve Issue

```
This is clear, reproducible, and in scope. Adding `status:approved`. A PR linking this issue is welcome.
```

---

## GitHub: Approve PR

```
Scoped, tested, linked to the approved issue. CI green. Merging.
```

---

## GitHub: Approve PR — With Minor Notes

```
Good to merge. Two minor notes, not blocking: [note 1]; [note 2]. Merging now.
```

---

## GitHub: Request Changes — Process Gap (issue-first)

```
[Optional one-phrase acknowledgment of genuine value.] I'm not going to merge this without a linked approved issue. Open or link one, get it approved, and add a `type:*` label. Then push an updated PR body and I'll re-review on the merits.
```

Real example (PR #101 style):
```
Thanks for the docs work. The content itself looks useful, but I am not going to merge a PR that skips the contribution contract. Open or link a real issue, get it approved, and add the correct type label. After that, push an updated PR body and we can review it on the merits instead of on preventable process failures.
```

---

## GitHub: Request Changes — Technical Gaps

```
[One sentence on what's solid.] Before this can merge: [gap 1], [gap 2]. [What to fix and how — one sentence.] Re-push and I'll take another look.
```

---

## GitHub: Close PR — Process Violation (idea still valid)

```
Closing this because it skips the issue-first process. No PR merges without a `status:approved` issue linked. The underlying idea is worth pursuing: open an issue, get it approved, then re-file the PR with the link. This isn't a rejection of the concept.
```

---

## GitHub: Close PR — Wrong Direction

```
Closing this. The direction conflicts with [principle] because [reason in one sentence]. The code quality is solid but the approach isn't where this project is going. [Optional: what the right direction looks like.]
```

---

## GitHub: Close Issue — Belongs in Discussions

```
This is more a discussion topic than a concrete bug or feature. Please continue in [Discussions](link). Re-open if you can reproduce it as a specific bug with steps.
```

---

## GitHub: Close Issue — Out of Scope

```
Closing this. [One sentence: why it conflicts with the project's scope.] Re-open if there's strong community demand with concrete evidence.
```

---

## GitHub: Close Issue — Duplicate

```
Duplicate of #[N]. Continuing there.
```

---

## GitHub: Needs Design

```
Good idea, but this touches [area] in a way that needs a design decision first. Propose an approach in the comments: what's the intended contract and behavior? We can get alignment before a PR makes sense.
```

---

## GitHub: Stale / Silence

```
This [PR/issue] has been waiting on contributor feedback since [date]. Closing to keep the backlog clean. Re-open or re-file any time.
```

---

## GitHub: Release Announcement

```
Released v[version].

- [What changed — bullet 1]
- [What changed — bullet 2]

Install: `[install command]`
Changelog: [link]
```

---

## GitHub: Redirect to Docs

```
This is covered in [docs section](link). [One sentence key answer.] Closing. Re-open with the specific scenario if the docs don't cover your case.
```

---

## GitHub: Ask for Evidence

```
Before this moves forward I need to see [specific evidence]. Without [reproduction steps / failing test / benchmark], there's no way to evaluate the claim. Drop that here and I'll take another look.
```

---

## GitHub: Explain Policy

```
Quick note on process: [policy in one sentence]. [Why it exists — one sentence.] [What to do now — one action.]
```

---

## GitHub: First-Time Contributor — Extra Context

Append to any close or request-changes message:

```
This is the normal process here, not a judgment on the work itself.
```

---

## Jira: Status Update

```
[Component] is [status]. [One sentence of current state.]

Blocked by: [PROJ-N or "nothing — continuing"]
Next: [one action and owner]
ETA: [date or "TBD pending X"]
```

---

## Jira: Scope Clarification

```
Narrowing scope to: [what's in]. Out of scope for this ticket: [what's out]. That belongs in [PROJ-N or a new ticket]. [Who owns the narrowed work and by when.]
```

---

## Slack: Quick Unblock

```
[Answer in 1–2 sentences.] Let me know if that unblocks you.
```

---

## Slack: Escalation or Blocker

```
Blocked on [X]. Needs [who] to [do what] by [when]. @[mention] can you confirm?
```

---

## Discord: Release Announcement

```
**Released v[version]**

- [Change 1]
- [Change 2]

```
[install command]
```

Changelog: [link]
```

---

## Discord: Process Reminder (contributor feedback)

```
Hey [name], quick note: [what was skipped]. [One sentence on the right path.] No rush, just flagging before you put more time in.
```
