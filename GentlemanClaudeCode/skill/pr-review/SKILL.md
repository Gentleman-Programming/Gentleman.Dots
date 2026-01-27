---
name: pr-review
description: >
  Reviews GitHub PRs and leaves human, direct comments.
  Trigger: When user asks to review a PR, check a PR, or gives a PR URL.
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

## When to Use

Use this skill when:
- User shares a PR URL
- User asks to review a PR
- User says "review this", "check this PR", "what do you think of this PR"

## Philosophy

**Be a helpful colleague, not a robot.**

You're reviewing code like a senior dev who actually CARES about the codebase and the author's growth. Not a linter. Not a checklist machine.

## Review Process

### 1. Fetch PR Information

```bash
# Get PR metadata
gh pr view {number} --json title,body,files,additions,deletions,author,state,baseRefName,headRefName

# Get the diff (exclude lockfiles - they're noise)
gh pr diff {number} -- ':(exclude)**/pnpm-lock.yaml' ':(exclude)**/package-lock.json' ':(exclude)**/yarn.lock'
```

### 2. For Large Diffs

If the diff is huge, fetch specific files from the PR branch:

```bash
gh api "repos/{owner}/{repo}/contents/{path}?ref={branch}" --jq '.content' | base64 -d
```

### 3. Categorize Issues

**Only comment on things that MATTER:**

| Category | When to Comment |
|----------|-----------------|
| CRITICAL | Will break in production, security issue, data loss |
| NEEDS REVIEW | Might be wrong, need author to confirm intent |
| QUESTION | Curious about a choice, not necessarily wrong |

**DO NOT comment on:**
- Style preferences (that's what linters are for)
- "You could also do X" suggestions (unless X is significantly better)
- Minor nitpicks
- Things that are correct but you'd do differently

### 4. Verify Before Commenting

**NEVER assume.** If you're not sure about something:

```bash
# Check official docs
# Search for the API/convention being used
# Verify the claim before saying "this is wrong"
```

Example: If someone uses a new API signature, CHECK if that's the new standard before flagging it.

## Writing Comments

### Tone

Write like you're talking to a colleague on Slack. Not a formal code review template.

**BAD (robotic):**
```
Issue: The function export does not conform to Next.js 16 proxy.ts specification.
Recommendation: Modify the export to use named export pattern.
Severity: Critical
```

**GOOD (human):**
```
Hey! Looking at the proxy.ts export - according to Next.js 16 docs, the function needs to be named `proxy`, not an anonymous default. Should be something like:

export const proxy = auth((req) => { ... })

Can you double check this works correctly? I'm seeing conflicting info about whether next-auth handles this automatically.
```

### Structure

Keep it simple:
1. What you noticed
2. Why it might be a problem (with evidence/link if needed)
3. Suggested fix (if you have one)
4. Ask for confirmation if you're not 100% sure

### One Comment, Multiple Points

If you have several things to discuss, put them ALL in one comment. Don't spam the PR with 10 separate comments.

```
Hey! Did a review. Found a few things worth discussing:

**1. First thing**
{explanation}

**2. Second thing**
{explanation}

Let me know if any of these need clarification!
```

## Posting the Comment

```bash
gh pr comment {number} --body "Your comment here"
```

For inline comments on specific files/lines, use the review API:

```bash
gh api repos/{owner}/{repo}/pulls/{number}/comments -f body="comment" -f path="file.ts" -f line=42 -f side="RIGHT"
```

## What Makes a Good Review

1. **Helpful** - Focuses on real issues, not nitpicks
2. **Humble** - Asks questions when unsure instead of demanding changes
3. **Human** - Reads like a person wrote it, not a template
4. **Verified** - Claims are backed by docs/evidence
5. **Actionable** - Author knows exactly what to do (or what question to answer)

## Anti-Patterns to Avoid

- Starting with "LGTM" then listing 15 issues
- "This is wrong" without explaining why
- Suggesting rewrites of working code just because you'd do it differently
- Commenting on every file just to show you reviewed everything
- Being condescending ("Obviously you should...")
- Using formal review templates with severity ratings

## Example Output

After reviewing, summarize for the user:

```
Listo, dejé el comentario en el PR: {url}

Los puntos principales:
1. {brief summary of critical/review items}
2. {another point}

A ver qué responde el autor.
```

## Keywords

pr, pull request, review, github, code review, gh
