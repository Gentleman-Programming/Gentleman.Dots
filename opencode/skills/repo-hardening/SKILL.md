---
name: repo-hardening
description: >
  Hardens any GitHub repository into a maintainer-grade setup with strong contribution
  gates, issue/PR templates, GitHub Actions policy enforcement, label conventions,
  stale strategy, and backlog hygiene — from first audit to durable process.
  Trigger: When hardening a repo, setting up maintainer workflow, tightening contribution
  gates, auditing repo health, adding issue/PR templates, or transforming a loose repo
  into a structured OSS-grade project.
  Key phrases: "harden repo", "repo setup", "contribution gates", "maintainer workflow",
  "PR templates", "issue templates", "label conventions", "stale strategy", "repo health".
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
allowed-tools: Read, Edit, Write, Glob, Grep, Bash, WebFetch
---

## When to Use

Use this skill when:
- Transforming a chaotic/fresh repo into a structured, maintainer-enforced one
- Adding or auditing contribution gates (CONTRIBUTING.md, templates, CI checks)
- Defining or formalizing the repo's maintainer philosophy and workflow
- Setting up label conventions, stale automation, or backlog triage cadence
- Doing a one-time hardening sprint or a rolling incremental hardening

**Don't use this skill when:**
- You only need a one-off backlog triage → use `backlog-triage` instead
- You only need to review PRs → use `pr-review` instead
- The repo already has a fully enforced workflow and you're just patching one file

---

## Operator vs. Repo Infrastructure

This is a critical distinction. Confusing them wastes effort.

| Layer | What it is | Where it lives |
|-------|-----------|----------------|
| **Operator tooling** | AI skills, agent rules, workflow shortcuts | `~/.config/opencode/skills/`, `AGENTS.md`, personal dotfiles |
| **Repo infrastructure** | Templates, workflows, docs, labels, rules | `.github/`, `CONTRIBUTING.md`, `README.md`, `CODEOWNERS` in the target repo |

**The skill you are reading now is operator tooling.**
**The files you create with it are repo infrastructure.**

Never confuse the two. Repo infrastructure must live in the target repo and survive without the agent. Operator tooling makes the agent smarter but has zero effect on the repo if missing.

---

## Maintainer Philosophy (customize per repo)

Before hardening, fill in this table from the repo's existing docs, past maintainer behavior, and stated goals. This drives every decision below.

| Principle | What it means in practice |
|-----------|--------------------------|
| **[Scope boundary]** | What is in scope? What is explicitly out of scope? |
| **[Quality bar]** | Minimum: tests required? Docs required? What merits a reject? |
| **[Contribution model]** | Issue-first? Direct PRs OK? Discussions before PRs? |
| **[Review speed SLA]** | How fast do PRs get reviewed? Stale after how many days? |
| **[Breaking change policy]** | Allowed? Requires RFC? Requires major bump? |
| **[Noise tolerance]** | Vague issues closed immediately? Given 7-day response window? |

### Example: minimal OSS utility

| Principle | Value |
|-----------|-------|
| Scope boundary | Narrow utility. No feature creep. New behaviors need an issue discussion first. |
| Quality bar | Tests required for logic changes. README update required for new behaviors. |
| Contribution model | Issue-first. Every PR must close an approved issue. |
| Review speed SLA | 7-day response target. Stale after 30 days, closed after 14 more. |
| Breaking change policy | Allowed with major bump. Must document migration path. |
| Noise tolerance | Vague issues get 7-day response request, then close. |

---

## Phased Hardening

Choose the tier that fits the repo's maturity and traffic level. **Don't over-engineer a low-traffic project.**

### Tier 1 — Minimal Hardening (any repo, day 1)

The absolute baseline. No friction for contributors, no workflow yet.

**Checklist:**
- [ ] `README.md` has a clear one-paragraph description + install/usage
- [ ] `LICENSE` file exists
- [ ] `CONTRIBUTING.md` exists with at minimum: how to open an issue, how to open a PR
- [ ] Default branch is `main` (not `master`)
- [ ] At least one issue template (bug report)
- [ ] At least one PR template

**Files to create/update:**
```
.github/
├── ISSUE_TEMPLATE/
│   └── bug_report.md
└── pull_request_template.md
CONTRIBUTING.md
```

**Command to verify:**
```bash
gh repo view <owner/repo> --json defaultBranchRef,hasIssuesEnabled,hasWikiEnabled
ls .github/
```

---

### Tier 2 — Standard Hardening (active project, regular contributors)

Enforced process. Reviewers enforced. Labels structured.

**Checklist (all of Tier 1 plus):**
- [ ] `CODEOWNERS` defined (at minimum one owner per area)
- [ ] Branch protection on `main`: require PR, require 1 review, dismiss stale reviews
- [ ] At least 3 issue templates: bug, feature request, docs
- [ ] PR template includes: linked issue field, checklist (tests, docs, changelog)
- [ ] Label taxonomy defined (see Label Conventions section)
- [ ] At least one GitHub Actions CI workflow (lint or test on push/PR)
- [ ] `stale.yml` action configured
- [ ] `CHANGELOG.md` or release notes policy documented

**Files to create/update:**
```
.github/
├── ISSUE_TEMPLATE/
│   ├── bug_report.md
│   ├── feature_request.md
│   └── docs_improvement.md
├── workflows/
│   ├── ci.yml
│   └── stale.yml
├── CODEOWNERS
└── pull_request_template.md
CONTRIBUTING.md           ← expand with issue-first policy
CHANGELOG.md
```

**Branch protection via CLI:**
```bash
# Require PR + 1 review + dismiss stale approvals on main
gh api repos/<owner>/<repo>/branches/main/protection \
  --method PUT \
  --field required_status_checks='{"strict":true,"contexts":["ci"]}' \
  --field enforce_admins=false \
  --field required_pull_request_reviews='{"required_approving_review_count":1,"dismiss_stale_reviews":true}' \
  --field restrictions=null
```

---

### Tier 3 — Strict Hardening (high-traffic OSS, multiple maintainers)

Full enforcement. Automated policy checks. Structured triage cadence.

**Checklist (all of Tier 2 plus):**
- [ ] Branch protection: require 2 reviews, require signed commits, require linear history
- [ ] `pr-policy.yml` workflow: check for linked issue in PR body, enforce label presence
- [ ] `SECURITY.md` with disclosure policy
- [ ] `GOVERNANCE.md` (or governance section in CONTRIBUTING) documenting merge authority
- [ ] Discussions enabled (redirect questions out of Issues)
- [ ] Release automation workflow (tag → changelog → GitHub Release)
- [ ] `dependabot.yml` configured for dependency updates
- [ ] Required labels enforced via policy workflow (see Label Conventions)
- [ ] Triage cadence documented (when/how maintainers review the backlog — e.g., weekly)

**Files to create/update:**
```
.github/
├── ISSUE_TEMPLATE/
│   ├── bug_report.md
│   ├── feature_request.md
│   ├── docs_improvement.md
│   └── config.yml          ← disable blank issues, add Discussions link
├── workflows/
│   ├── ci.yml
│   ├── stale.yml
│   ├── pr-policy.yml
│   ├── release.yml
│   └── dependabot.yml
├── CODEOWNERS
└── pull_request_template.md
CONTRIBUTING.md
CHANGELOG.md
SECURITY.md
GOVERNANCE.md
```

---

## Label Conventions

Use a structured namespace. Labels tell the story of an issue/PR lifecycle at a glance.

| Namespace | Examples | Purpose |
|-----------|---------|---------|
| `type:` | `type:bug`, `type:feature`, `type:docs`, `type:chore`, `type:security` | What kind of work |
| `status:` | `status:needs-triage`, `status:approved`, `status:in-progress`, `status:blocked`, `status:wontfix` | Where it is in the lifecycle |
| `priority:` | `priority:critical`, `priority:high`, `priority:low` | Urgency (use sparingly) |
| `area:` | `area:cli`, `area:api`, `area:docs`, `area:ci` | Subsystem (customize per repo) |
| `effort:` | `effort:small`, `effort:medium`, `effort:large` | Rough sizing for contributors |

**Create all labels via CLI:**
```bash
# Type labels
gh label create "type:bug"      --color "d73a4a" --repo <owner/repo>
gh label create "type:feature"  --color "a2eeef" --repo <owner/repo>
gh label create "type:docs"     --color "0075ca" --repo <owner/repo>
gh label create "type:chore"    --color "e4e669" --repo <owner/repo>
gh label create "type:security" --color "e11d48" --repo <owner/repo>

# Status labels
gh label create "status:needs-triage"  --color "ededed" --repo <owner/repo>
gh label create "status:approved"      --color "0e8a16" --repo <owner/repo>
gh label create "status:in-progress"   --color "fbca04" --repo <owner/repo>
gh label create "status:blocked"       --color "b60205" --repo <owner/repo>
gh label create "status:wontfix"       --color "ffffff" --repo <owner/repo>

# Priority labels
gh label create "priority:critical" --color "b60205" --repo <owner/repo>
gh label create "priority:high"     --color "e11d48" --repo <owner/repo>
gh label create "priority:low"      --color "c5def5" --repo <owner/repo>

# Effort labels
gh label create "effort:small"  --color "d4edda" --repo <owner/repo>
gh label create "effort:medium" --color "fff3cd" --repo <owner/repo>
gh label create "effort:large"  --color "f8d7da" --repo <owner/repo>
```

---

## GitHub Actions Policy Checks

These workflows enforce process without requiring a human to catch every violation.

### `pr-policy.yml` — Enforce PR hygiene automatically

```yaml
# .github/workflows/pr-policy.yml
name: PR Policy Check

on:
  pull_request:
    types: [opened, edited, synchronize, labeled, unlabeled]

jobs:
  check-linked-issue:
    name: Require linked issue
    runs-on: ubuntu-latest
    steps:
      - name: Check PR body for issue reference
        env:
          PR_BODY: ${{ github.event.pull_request.body }}
        run: |
          if ! echo "$PR_BODY" | grep -qiE '(closes|fixes|resolves) #[0-9]+'; then
            echo "❌ PR must reference a closing issue (e.g. 'Closes #42')"
            exit 1
          fi

  check-type-label:
    name: Require type label
    runs-on: ubuntu-latest
    steps:
      - name: Check for type label
        env:
          LABELS: ${{ toJSON(github.event.pull_request.labels.*.name) }}
        run: |
          if ! echo "$LABELS" | grep -q '"type:'; then
            echo "❌ PR must have exactly one type:* label"
            exit 1
          fi
```

### `stale.yml` — Auto-close stale issues and PRs

```yaml
# .github/workflows/stale.yml
name: Stale Issues and PRs

on:
  schedule:
    - cron: '0 9 * * 1'  # Every Monday at 9am UTC
  workflow_dispatch:

jobs:
  stale:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/stale@v9
        with:
          stale-issue-message: >
            This issue has been inactive for 30 days. If it's still relevant,
            please leave a comment with updated context. Otherwise it will be
            closed in 14 days.
          close-issue-message: >
            Closing due to inactivity. Feel free to reopen with updated info.
          stale-pr-message: >
            This PR has been inactive for 30 days. Please rebase and update it,
            or it will be closed in 14 days.
          close-pr-message: >
            Closing due to inactivity. Feel free to re-open when ready to continue.
          days-before-stale: 30
          days-before-close: 14
          stale-issue-label: 'status:stale'
          stale-pr-label: 'status:stale'
          exempt-issue-labels: 'status:approved,status:in-progress,priority:critical'
          exempt-pr-labels: 'status:in-progress'
```

---

## Issue & PR Templates

### Bug Report Template

```markdown
<!-- .github/ISSUE_TEMPLATE/bug_report.md -->
---
name: Bug report
about: Something is broken
labels: type:bug, status:needs-triage
---

## Description
<!-- What is broken? One clear sentence. -->

## Steps to reproduce
1.
2.
3.

## Expected behavior
<!-- What should happen -->

## Actual behavior
<!-- What happens instead -->

## Environment
- OS:
- Version:
- Relevant config:

## Logs / screenshots
<!-- Paste relevant output here -->
```

### Feature Request Template

```markdown
<!-- .github/ISSUE_TEMPLATE/feature_request.md -->
---
name: Feature request
about: Propose a new capability
labels: type:feature, status:needs-triage
---

## Problem statement
<!-- What problem does this solve? Be specific. -->

## Proposed solution
<!-- What should it do? What is the expected behavior? -->

## Alternatives considered
<!-- What else did you think of? Why is this better? -->

## Out of scope
<!-- What should this NOT do? -->
```

### PR Template

```markdown
<!-- .github/pull_request_template.md -->
## Summary
<!-- One sentence: what does this PR do? -->

## Closes
<!-- Required: reference the approved issue this PR closes -->
Closes #

## Type of change
<!-- Add exactly one type:* label to this PR -->
- [ ] Bug fix (`type:bug`)
- [ ] New feature (`type:feature`)
- [ ] Documentation (`type:docs`)
- [ ] Chore / internal (`type:chore`)

## Checklist
- [ ] Tests added or updated for changed behavior
- [ ] Documentation updated if behavior changed
- [ ] CHANGELOG.md entry added (if user-facing change)
- [ ] No unrelated changes included
```

---

## CONTRIBUTING.md Structure

A `CONTRIBUTING.md` that actually enforces behavior must cover:

```markdown
# Contributing

## Before You Start
- Check if the issue you want to fix already has status:approved
- Don't open a PR without a linked approved issue — it will be closed

## Issue Process
1. Search existing issues first
2. Use the correct template (bug, feature, or docs)
3. Wait for maintainer approval (status:approved label) before opening a PR

## PR Process
1. Fork the repo and create a branch: `type/short-description` (e.g. `fix/auth-crash`)
2. PR title must follow conventional commits: `fix: short description`
3. PR must close an approved issue (`Closes #N`)
4. PR must have exactly one `type:*` label
5. Ensure CI passes before requesting review

## What Gets Merged
- Bug fixes with tests
- Features with approved issue + tests + docs update
- Docs improvements that are accurate and scoped

## What Gets Closed
- PRs without an approved issue
- PRs without a `type:*` label
- PRs that are vague, scope-breaking, or fail CI
- Issues that are duplicates, vague, or belong in Discussions
```

---

## Stale Strategy Decision Table

| Scenario | Recommended action |
|----------|--------------------|
| Active project, regular contributors | 30-day stale, 14-day close |
| Low-traffic project | 60-day stale, 30-day close |
| High-traffic OSS | 14-day stale, 7-day close |
| Critical bugs | Never auto-close (`exempt-issue-labels: priority:critical`) |
| Approved issues pending PR | Never auto-close (`exempt-issue-labels: status:approved`) |
| PRs in review | Never auto-close (`exempt-pr-labels: status:in-progress`) |

---

## Backlog Triage Cadence

Hardening is not a one-off. Document the cadence explicitly.

| Triage type | Frequency | Who | What |
|-------------|-----------|-----|------|
| New issue triage | Within 2 days of creation | Maintainer | Label + approve/reject |
| PR review | Within 7 days of opening | Reviewer | Review + action |
| Backlog sweep | Weekly (Monday) | Maintainer | Close stale, prioritize approved |
| Label audit | Monthly | Maintainer | Remove orphan labels, check accuracy |
| Dependency updates | Weekly (Dependabot PRs) | Maintainer | Review + merge or defer |

**Put the cadence in CONTRIBUTING.md** so contributors know what to expect.

---

## Audit Checklist (State Assessment)

Run this before hardening to know what exists and what's missing.

```bash
# Check what exists in .github/
ls -la .github/ 2>/dev/null || echo "No .github/ directory"
ls -la .github/ISSUE_TEMPLATE/ 2>/dev/null || echo "No issue templates"
ls -la .github/workflows/ 2>/dev/null || echo "No workflows"

# Check key docs
ls CONTRIBUTING.md LICENSE README.md CHANGELOG.md SECURITY.md 2>/dev/null

# Check branch protection
gh api repos/<owner>/<repo>/branches/main/protection 2>/dev/null || echo "No branch protection"

# Check existing labels
gh label list --repo <owner/repo>

# Check Discussions enabled
gh repo view <owner/repo> --json hasDiscussionsEnabled

# Quick health summary
gh repo view <owner/repo> --json name,description,defaultBranchRef,isPrivate,hasIssuesEnabled,hasWikiEnabled,hasProjectsEnabled
```

---

## Decision Table: What to Harden When

| Signal | Action |
|--------|--------|
| No issue templates | Tier 1 minimum → add bug report template immediately |
| PRs merged without review | Add branch protection (1 required review) |
| No labels / chaotic labels | Define label taxonomy, bulk-apply `status:needs-triage` |
| PRs without tests getting merged | Add CI workflow + PR checklist |
| Issues with no response for 30+ days | Enable stale action + triage sweep |
| Vague PRs with no linked issue | Add `pr-policy.yml` + CONTRIBUTING process gate |
| No CODEOWNERS | Add CODEOWNERS immediately for any multi-maintainer repo |
| Contributors confused about process | Expand CONTRIBUTING.md — it's your best async documentation |
| High PR volume, slow reviews | Enable Discussions for questions, reduce issue noise |
| Dependabot PRs piling up | Add `dependabot.yml` + weekly merge cadence |

---

## Anti-Patterns to Avoid

| Anti-pattern | Why it fails | Better approach |
|--------------|-------------|----------------|
| Adding every possible template | Template fatigue → contributors skip them | 2-3 focused templates > 8 half-used ones |
| Enforcing everything on day 1 | New contributors bounce off | Tier 1 first, earn stricter gates |
| Labels no one maintains | Drift → labels become noise | Only define labels you will actually apply |
| Stale closes everything | Closes real issues, erodes trust | Exempt approved + critical issues |
| Complex branch naming | Nobody remembers it | Simple: `type/description`, document it |
| CODEOWNERS for everything | Approval bottleneck | CODEOWNERS only for sensitive paths |
| PR policy failing silently | Contributors don't know why CI failed | Policy checks must have clear error messages |

---

## Reusable Prompt Template

Copy this to apply the hardening workflow to any repository:

```
You are a repo-hardening agent. Your job is to transform <owner/repo> from its current
state into a maintainer-grade project with durable contribution gates and backlog hygiene.

## Repo Context
- Owner/repo: <owner/repo>
- Project type: <CLI tool / web app / library / OSS utility / internal tool>
- Current state: <fresh / messy / some structure / already has CI>
- Target tier: <Tier 1 / Tier 2 / Tier 3>

## Maintainer Philosophy
<Fill in the philosophy table from the skill — customize for this repo>

## Steps
1. Run the Audit Checklist to assess current state
2. Identify which Tier to target based on project maturity
3. Create all missing files using the templates from the skill
4. Set up label taxonomy using the gh CLI commands
5. Configure branch protection appropriate to the tier
6. Add GitHub Actions workflows (pr-policy.yml, stale.yml, ci.yml as needed)
7. Update CONTRIBUTING.md with the full contribution process
8. Document the triage cadence in CONTRIBUTING.md
9. Produce a hardening report: what was done, what remains, what tier was reached

## Constraints
- Do not over-engineer a low-traffic project — match the tier to the reality
- All workflows must have clear error messages so contributors understand failures
- Stale action must exempt approved and priority:critical issues
- PR policy must be enforced by CI, not by convention alone
- CONTRIBUTING.md is the source of truth — keep it honest and short
```

---

## Resources

- **Templates**: See [assets/](assets/) for ready-to-use file templates
- **Related skills**: `backlog-triage` (one-off audit), `pr-review` (reviewing individual PRs)
