# Contributing to <project>

Thanks for taking the time to contribute! This document explains how the contribution
process works and what gets merged.

---

## Before You Start

- **Search existing issues and PRs.** Your problem or idea may already be tracked.
- **No PR without an approved issue.** Every PR must close a `status:approved` issue.
  If you open a PR without one, it will be closed — not because we don't appreciate the
  work, but because untracked changes are hard to review and often break scope.
- **Questions go to Discussions.** The issue tracker is for bugs and approved features only.
  Use [Discussions](<link>) for questions, ideas, and feedback.

---

## Issue Process

1. **Use the right template.** Bug report, feature request, or docs improvement.
2. **Be specific.** Vague issues (e.g. "it doesn't work") will be closed.
3. **Wait for `status:approved`.** A maintainer will review your issue and add the label
   before a PR is welcome. This usually happens within **[X] days**.

---

## PR Process

1. **Fork the repo** and create a branch from `main`.
2. **Branch naming**: `type/short-description` — e.g. `fix/auth-crash`, `feat/export-json`.
3. **Write tests** for any changed or new behavior.
4. **Update docs** if the behavior is visible to users.
5. **Conventional commit title**: `fix: short description` / `feat: short description`.
6. **Reference the closing issue**: `Closes #N` in the PR description.
7. **Add exactly one `type:*` label** to the PR before requesting review.
8. **Ensure CI passes.** Failing CI = no review.

---

## What Gets Merged

| Contribution | Requirements |
|-------------|-------------|
| Bug fix | Linked approved issue + tests + CI green |
| New feature | Linked approved issue + tests + docs update + CI green |
| Docs improvement | Accurate, scoped, linked approved issue |
| Chore / refactor | Linked approved issue, no behavior changes |

---

## What Gets Closed

| Situation | Action |
|-----------|--------|
| PR without a linked approved issue | Closed immediately |
| PR without a `type:*` label | Requested changes |
| PR failing CI | Requested changes |
| Vague or duplicate issue | Closed with explanation |
| Question in Issues | Redirected to Discussions |
| Issue inactive > 30 days | Auto-stale, closed after 14 more days |

---

## Review SLA

- **New issues**: Maintainer response within **[X] days**.
- **Open PRs**: First review within **[X] days** of opening.
- **Stale**: Issues and PRs inactive for 30 days are marked stale. Active ones (approved or
  in-progress) are never auto-closed.

---

## Local Development

```bash
# Install dependencies
<install command>

# Run tests
<test command>

# Run linter
<lint command>
```

---

## Governance

Decisions about architecture, scope, and breaking changes are made by maintainers.
For significant proposals, open a feature request issue and discuss before opening a PR.
