---
name: pr-review
description: >
  Review GitHub PRs and Issues with structured analysis for open source projects.
  Trigger: When user wants to review PRs (even if first asking what's open), analyze issues, or audit PR/issue backlog.
  Key phrases: "pr review", "revisar pr", "qué PRs hay", "PRs pendientes", "issues abiertos", "sin atención", "hacer review".
license: MIT
metadata:
  author: gentleman-programming
  version: "1.2"
---

## When to Use

**ALWAYS use this skill when user mentions "pr review", "revisar PRs", or asks about pending PRs/issues** - even if they first ask what's pending. This skill handles the FULL flow: listing → analyzing → reviewing.

Specific triggers:
- User wants to review PRs (even if first asking what's open)
- Analyze issues or contributions
- Audit PR/issue backlog
- Check what needs attention

**Key phrases:** "pr review", "revisar", "qué hay pendiente", "sin atención", "PRs abiertos", "issues abiertos", "hacer review", "necesito revisar"

## Review Process

### Phase 1: Gather Information

```bash
# List all issues and PRs
gh issue list --state all --limit 20
gh pr list --state all --limit 20

# Get PR details (run in parallel)
gh pr view {number} --json title,body,files,additions,deletions,author
gh pr diff {number} --patch
```

### Phase 2: Load Project Skills (MANDATORY)

Before reviewing ANY code, check if the repo has project-specific skills that define conventions. These are your review criteria — not just generic best practices.

**How to find them:**
1. Check `AGENTS.md` at the repo root — it lists all available skills and auto-invoke rules
2. Check `skills/` directory for project-specific skill files
3. If the repo has an `AGENTS.md` with an Auto-invoke Skills table, read it

**For each PR, load the skills that match the changed files:**

| Files Changed | Skills to Load |
|---------------|----------------|
| `api/` (models, views, serializers) | project API skill + `django-drf` |
| `api/` (tests) | project API test skill + `pytest` |
| `ui/` (components, pages) | project UI skill + `react-19` + `nextjs-15` + `tailwind-4` |
| `ui/` (tests) | project UI test skill + `playwright` |
| `ui/` (schemas) | `zod-4` |
| `ui/` (stores) | `zustand-5` |
| Types/interfaces | `typescript` |

**Review against project conventions, not just general quality.** Check:
- Does the file structure match what the project skill defines?
- Are naming conventions followed? (serializer names, component placement, test structure)
- Are the right patterns used? (service layer vs serializer logic, Server Components vs Client)
- Do tests follow the project's test patterns? (fixtures, assertions, POM for E2E)

If no project skills exist, fall back to generic best practices.

### Phase 3: Read Current Codebase

Before reviewing diffs, **always read the current code** to understand context:
- Main entry points
- Files being modified
- Related modules

### Phase 4: Analyze Each PR

For each PR, evaluate these factors:

| Factor | What to Check |
|--------|---------------|
| **Project Conventions** | Does it follow the project skills? Structure, naming, patterns |
| **Code Quality** | Clean code, no duplication, proper error handling |
| **Tests** | Are there tests? Do they follow the project's test patterns? |
| **Breaking Changes** | Does it break existing functionality? |
| **Conflicts** | Will it conflict with other open PRs? |
| **Commit Hygiene** | Clean history, no test files, proper messages |
| **Documentation** | README updated if needed, comments where necessary |

## Critical Patterns

### Red Flags (DO NOT MERGE)

- [ ] Test/debug files committed (`test.js`, `console.log`, etc.)
- [ ] Variables declared but never used
- [ ] Code duplication instead of refactoring
- [ ] Broken indentation or syntax errors
- [ ] Config files with personal/local settings
- [ ] Hardcoded secrets or credentials
- [ ] Breaking changes without migration path

### Yellow Flags (Request Changes)

- [ ] Too many commits (should squash)
- [ ] Missing validation for dependencies (e.g., `jq`, `curl`)
- [ ] Potential conflicts with other PRs
- [ ] Incomplete feature (e.g., variable declared but not used)
- [ ] Fallback code with bugs (e.g., not escaping newlines)

### Green Flags (Good to Merge)

- [x] Small, focused changes
- [x] Tests included
- [x] Clean commit history (1-3 commits)
- [x] Documentation updated
- [x] No conflicts with other PRs
- [x] Solves a real issue

## Decision Matrix

```
Has red flags?           → DO NOT MERGE, request fixes
Has yellow flags only?   → Request changes, can merge after fixes
All green?               → MERGE
```

## Output Format

### For Issues

```markdown
## Issues Analysis

### Good Issues (Valid, should be addressed)
| # | Issue | Analysis |
|---|-------|----------|
| **#XX** | Title | Why it's valid |

### Questionable Issues
| # | Issue | Analysis |
|---|-------|----------|
| **#XX** | Title | Problems with this issue |

### Should Close
| # | Issue | Reason |
|---|-------|--------|
| **#XX** | Title | Why it should be closed |
```

### For PRs

```markdown
## PR Analysis

### Ready to Merge
| PR | Author | Why it's ready |
|----|--------|----------------|
| **#XX** | @user | Brief explanation |

### Needs Work
| PR | Author | What to fix |
|----|--------|-------------|
| **#XX** | @user | List of issues |

### Do Not Merge
| PR | Author | Critical problems |
|----|--------|-------------------|
| **#XX** | @user | Why it can't be merged |
```

## Review Comments

### Language Rules

**Reply in the same language the author used in their PR/issue:**
- PR written in Spanish → Reply in Spanish
- PR written in English → Reply in English

### Comment Style: Concise & Human

Write review comments like a senior engineer talking to a colleague — direct, clear, no fluff. NOT like a template.

**Rules:**
- Lead with the issues, numbered. No greetings, no "Hey {Name}!".
- Each issue: **bold the problem** in one phrase, then explain in 1-2 plain sentences. Include the concrete fix inline.
- End with 1-2 sentences acknowledging what's good. Don't force it — only if something genuinely stood out.
- No emojis in the review body. No `##` headings. No horizontal rules. Just numbered points and a closing line.
- No "Solution" sections — the fix goes inline with the issue description.
- Keep it short. If you can say it in one sentence, don't use two.

### Approve Format

One sentence — what's good, ship it. Optionally a follow-up note.

```
Clean refactor, all spec requirements covered, 28 tests. Ship it.
```

```
Well done. Service layer pattern, anti-enumeration, rate limiting, 32 tests. Synchronous email is fine for MVP.
```

```
Solid. Fire-and-forget with proper timeouts, 5 tests. One note: the spec still says single-field payload but code sends {type, data} — code is better, update the spec in a follow-up.
```

### Request Changes Format

```
Two things to address:

1. **UpdateModelMixin exposes PUT** — you only need PATCH here. Add `http_method_names = ["get", "patch", "head", "options"]` to the ViewSet so PUT isn't accidentally exposed.

2. **partner_id in refresh token** — `get_token()` adds partner_id to the refresh token, and the access inherits from it, so it ends up in both. The design doc says access only. Either move the claim injection to `validate()` on the access token, or update the design doc if you're ok with it being in both.

Everything else looks solid — sign-in guards correctly use 403, is_staff is in the serializer, tests are thorough. Nice work on the service layer separation.
```

```
One thing — in partner-kickoff/SKILL.md, the Step 5 heading is sitting between Step 9 and Step 11. Looks like it didn't get renumbered when the workflow was restructured. Move it to the right position or renumber.

The verify skill itself is well-structured — clear verdict rules, good fresh-context pattern.
```

### Anti-patterns to AVOID

- "Hey John! Thanks for the PR, the analysis is well done" — skip the greeting, get to the point
- "## Problem Category" / "## Solution" headings — too formal, use numbered list
- Long code blocks showing the fix — one line inline is enough
- "Great job! Just a few minor things..." — empty praise before criticism
- Emojis anywhere in the review body
- Repeating what the PR description already says

## Commands

```bash
# Merge a PR (after approval)
gh pr merge {number} --merge

# Leave a review comment
gh pr review {number} --comment --body "$(cat <<'EOF'
{comment content}
EOF
)"

# Request changes
gh pr review {number} --request-changes --body "..."

# Approve
gh pr review {number} --approve --body "..."

# Close issue as not planned
gh issue close {number} --reason "not planned" --comment "..."
```

## Conflict Detection

When reviewing multiple PRs, check for conflicts:

1. **Same files modified** - Check if PRs touch the same files
2. **Dependent features** - PR A adds feature, PR B extends it
3. **Version bumps** - Multiple PRs changing VERSION
4. **Provider patterns** - New providers need to be added to all switch/case statements

### Common Conflict Pattern: Provider Addition

When a PR adds timeout/wrapper logic with hardcoded providers:
```bash
case "$provider" in
  claude) ...
  gemini) ...
  *) echo "Unknown" ;;  # New providers will fail!
esac
```

**Flag this** - any PR adding new providers will conflict.

## Merge Order Strategy

When multiple PRs have dependencies:

1. **Independent, small PRs first** - Quick wins, no conflicts
2. **Infrastructure PRs second** - Timeout, error handling, etc.
3. **Feature PRs third** - New providers, modes, etc.
4. **Large refactors last** - Most likely to have conflicts

## Checklist Before Merging

- [ ] All tests pass (if CI exists)
- [ ] No red flags in code
- [ ] No conflicts with recently merged PRs
- [ ] Author's branch is up to date with main
- [ ] Review comments addressed (if any)
