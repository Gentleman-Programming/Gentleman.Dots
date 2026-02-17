---
name: chained-pr
description: >
  Creates GitHub PRs following the Chained PRs workflow pattern.
  Trigger: When user asks to create a PR for a feature with sub-tasks, chained PR, or stacked PR workflow.
license: Apache-2.0
metadata:
  author: prowler-cloud
  version: "3.0"
---

## When to Use

Use this skill when creating GitHub Pull Requests that are part of a **Chained PRs** or **Stacked PRs** workflow:

| Scenario | Approach | PR Type | Target Branch |
|----------|----------|---------|---------------|
| Sub-task of a feature (tightly coupled API+UI) | Feature Branch | Chained PR | Feature branch |
| All sub-tasks done, ready to ship | Feature Branch | Main PR | master |
| Backend endpoint (can ship independently) | Stacked to Master | Stacked PR | master or previous PR |
| Frontend that depends on backend PRs | Stacked to Master | Dependent PR | master |
| Simple fix/hotfix (no chain) | **Don't use this skill** | Regular PR | master |
| Regular PR to master | **Don't use this skill** | Regular PR | master |

### Choosing the Right Approach

| Scenario | Feature Branch | Stacked to Master |
|----------|----------------|-------------------|
| API & UI tightly coupled, need atomic rollback | ✅ Better | ❌ |
| API can ship independently | ❌ | ✅ Better |
| Long-running feature (2+ weeks) | ❌ Conflict risk | ✅ Better |
| Multiple devs on same track | ✅ Better | ⚠️ Coordination needed |
| Want fast incremental reviews | ⚠️ Slower | ✅ Better |

### When NOT to Use

- **Regular PRs** that go directly to master (no chain or stack)
- **Hotfixes** or **bug fixes**
- **Very small features** (50 lines, single PR is fine)
- **Solo work** without parallel tracks

### Common Triggers

- "Create a PR for this sub-task"
- "Open a chained PR"
- "Create the main PR for the feature branch"
- "Stacked PR for PROWLER-XXX"
- "Create a stacked PR to master"

---

# Approach 1: Feature Branch (Chained PRs)

## Overview

A series of Pull Requests where each one:
1. Builds on top of the previous one
2. Is small enough to review independently
3. Merges to a **feature branch** (not master)
4. Allows parallel work between API and UI teams

### The Pattern

```
master
 └── feat/PROWLER-XXX-feature-name  ← Feature branch (base for all PRs)
      ├── PR #1: API contract     → merges to feature branch
      ├── PR #2: Endpoint A       → merges to feature branch
      ├── PR #3: Endpoint B       → merges to feature branch
      ├── PR #4: UI Components    → merges to feature branch
      ├── PR #5: Server Actions   → merges to feature branch
      ├── PR #6: Implement view   → merges to feature branch
      └── FINAL: feature branch   → merges to master (Main PR)
```

**Key Rule:** All chained PRs target the **feature branch**, NOT master.

### What Generates a PR vs What Doesn't

| Task Type | Generates PR? | Example |
|-----------|---------------|---------|
| **Implementation** | ✅ Yes | Create endpoint, Build component |
| **Analysis/Design** | ❌ No | Define API contract, Design data model |
| **Documentation** | ❌ No | Write RFC, Update specs |
| **Investigation** | ❌ No | Research options, Spike |

**Analysis tasks** (like "Define API model/contract") are **prerequisites** that happen BEFORE coding. They produce documentation/decisions, not code PRs.

---

## PR Types

### 1. Chained PR (Sub-task → Feature Branch)

Individual PRs that implement one sub-task. They merge to the feature branch.

### 2. Main PR (Feature Branch → Master)

The final PR that brings all the work to master. Created when all sub-tasks are complete.

---

## Chained PR Template

**Use when:** Creating a PR for a sub-task that merges to the feature branch.

### Title Format

```
[CHAIN] {conventional commit type}({scope}): {description}
```

**Examples:**
- `[CHAIN] feat(api): define contract for Findings Groups`
- `[CHAIN] feat(api): create GET /findings/groups endpoint`
- `[CHAIN] feat(ui): add TreeView component for hierarchical data`

### Description Template

> **IMPORTANT: Public Repository**
> Do NOT include links to private resources (Jira, Notion, Figma) in PR descriptions.
> Only reference other GitHub PRs and public information.

````markdown
## 🔗 Part of Chained PRs

| Field | Value |
|-------|-------|
| **Feature Branch** | `feat/PROWLER-XXX-{name}` |
| **Main PR** | [#NNNN](link-to-main-pr) |
| **Chain Position** | {N} of {total} |

### Chain Overview

```
         ┌─────────────────┐
         │ #NNN Sub-task 1 │
         └────────┬────────┘
                  │
         ┌───────────────────┐
         │                   │
         ▼                   ▼
┌─────────────────┐ ┌─────────────────┐
│ 📍 #NNN This PR │ │ #NNN Sub-task 3 │
└────────┬────────┘ └────────┬────────┘
         │                   │
         └─────────┬─────────┘
                   ▼
         ┌─────────────────┐
         │ #NNN Main PR    │
         │   → master      │
         └─────────────────┘
```

---

### Context

{Explain WHY this sub-task exists and how it fits into the larger feature.}

Part of chained PRs for {Feature Name}. This sub-task implements {specific piece}.

---

### Description

{Summary of what this PR delivers - the WHAT, not the HOW}

**Changes:**
- {Change 1}
- {Change 2}
- {Change 3}

---

### Steps to review

1. {Step 1 - what to check first}
2. {Step 2 - how to verify functionality}
3. {Step 3 - edge cases to test}

---

### Checklist

<details>

<summary><b>Community Checklist</b></summary>

- [x] This feature/issue is listed in [here](https://github.com/prowler-cloud/prowler/issues?q=sort%3Aupdated-desc+is%3Aissue+is%3Aopen) or roadmap.prowler.com
- [x] Is it assigned to me, if not, request it via the issue/feature in [here](https://github.com/prowler-cloud/prowler/issues?q=sort%3Aupdated-desc+is%3Aissue+is%3Aopen) or [Prowler Community Slack](goto.prowler.com/slack)

</details>

- [ ] Review if the code is being covered by tests.
- [ ] Review if code is being documented following this specification https://github.com/google/styleguide/blob/gh-pages/pyguide.md#38-comments-and-docstrings
- [ ] Review if backport is needed.
- [ ] Review if is needed to change the [Readme.md](https://github.com/prowler-cloud/prowler/blob/master/README.md)
- [ ] Ensure new entries are added to [CHANGELOG.md](https://github.com/prowler-cloud/prowler/blob/master/prowler/CHANGELOG.md), if applicable.

#### SDK/CLI
- Are there new checks included in this PR? Yes / No
    - If so, do we need to update permissions for the provider? Please review this carefully.

#### UI
- [ ] All issue/task requirements work as expected on the UI
- [ ] Screenshots/Video of the functionality flow (if applicable) - Mobile (X < 640px)
- [ ] Screenshots/Video of the functionality flow (if applicable) - Table (640px > X < 1024px)
- [ ] Screenshots/Video of the functionality flow (if applicable) - Desktop (X > 1024px)
- [ ] Ensure new entries are added to [CHANGELOG.md](https://github.com/prowler-cloud/prowler/blob/master/ui/CHANGELOG.md), if applicable.

#### API
- [ ] All issue/task requirements work as expected on the API
- [ ] Endpoint response output (if applicable)
- [ ] EXPLAIN ANALYZE output for new/modified queries or indexes (if applicable)
- [ ] Performance test results (if applicable)
- [ ] Any other relevant evidence of the implementation (if applicable)
- [ ] Verify if API specs need to be regenerated.
- [ ] Check if version updates are required (e.g., specs, Poetry, etc.).
- [ ] Ensure new entries are added to [CHANGELOG.md](https://github.com/prowler-cloud/prowler/blob/master/api/CHANGELOG.md), if applicable.

### License

By submitting this pull request, I confirm that my contribution is made under the terms of the Apache 2.0 license.
````

---

## Main PR Template

**Use when:** Creating the Main PR that will eventually merge the feature branch to master.

> **IMPORTANT:** Create the Main PR early with the `no-merge` label. This PR serves as:
> - Documentation of the full feature
> - Central place to track all chained PRs
> - Final PR to merge when all sub-tasks are done
>
> Remove the `no-merge` label only when ALL chained PRs are merged to the feature branch.

### Label

```bash
gh pr edit <number> --add-label "no-merge"
```

### Title Format

```
{conventional commit type}({scope}): {feature description}
```

**Examples:**
- `feat(ui): Findings Hierarchical Tree View (Check → Resources)`
- `feat(api): Multi-tenant Support`
- `feat: GovCloud Provider Integration`

### Description Template

> **IMPORTANT: Public Repository**
> Do NOT include links to private resources (Jira, Notion, Figma) in PR descriptions.
> Only reference other GitHub PRs and public information.

````markdown
## 🚀 Feature Complete - Chained PRs

This PR merges the `feat/PROWLER-XXX-{name}` feature branch.

### Included PRs

| # | PR | Sub-task | Status |
|---|-----|----------|--------|
| 1 | [#NNN](link) | {Title} | 🟢 Merged |
| 2 | [#NNN](link) | {Title} | 🟢 Merged |
| 3 | [#NNN](link) | {Title} | 🟡 Open |

**Legend:** 🟢 Merged | 🟡 Open | 🔴 Changes Requested | ⚪ Pending

### Dependency Diagram

```
         ┌─────────────────┐
         │ #NNN Sub-task 1 │
         └────────┬────────┘
                  │
         ┌───────────────────┐
         │                   │
         ▼                   ▼
┌─────────────────┐ ┌─────────────────┐
│ #NNN Sub-task 2 │ │ #NNN Sub-task 3 │
└────────┬────────┘ └────────┬────────┘
         │                   │
         └─────────┬─────────┘
                   ▼
         ┌─────────────────┐
         │  📍 THIS PR     │
         │   → master      │
         └─────────────────┘
```

---

### Context

{Explain the feature, its value, and why it was built this way.}

This PR merges the `feat/PROWLER-XXX-{name}` feature branch which implements {feature description}.

---

### Description

{Summary of all changes included in this feature}

**Key Changes:**

#### API
- {API change 1}
- {API change 2}

#### UI
- {UI change 1}
- {UI change 2}

---

### Steps to review

1. Review individual chained PRs first (linked above)
2. {Integration testing steps}
3. {E2E verification steps}

---

### Checklist

<details>

<summary><b>Community Checklist</b></summary>

- [x] This feature/issue is listed in [here](https://github.com/prowler-cloud/prowler/issues?q=sort%3Aupdated-desc+is%3Aissue+is%3Aopen) or roadmap.prowler.com
- [x] Is it assigned to me, if not, request it via the issue/feature in [here](https://github.com/prowler-cloud/prowler/issues?q=sort%3Aupdated-desc+is%3Aissue+is%3Aopen) or [Prowler Community Slack](goto.prowler.com/slack)

</details>

- [ ] All chained PRs merged to feature branch
- [ ] Feature tested end-to-end on feature branch
- [ ] No unresolved conflicts with master
- [ ] Review if the code is being covered by tests.
- [ ] Review if code is being documented following this specification https://github.com/google/styleguide/blob/gh-pages/pyguide.md#38-comments-and-docstrings
- [ ] Review if backport is needed.
- [ ] Review if is needed to change the [Readme.md](https://github.com/prowler-cloud/prowler/blob/master/README.md)
- [ ] Ensure new entries are added to [CHANGELOG.md](https://github.com/prowler-cloud/prowler/blob/master/prowler/CHANGELOG.md), if applicable.

#### SDK/CLI
- Are there new checks included in this PR? Yes / No
    - If so, do we need to update permissions for the provider? Please review this carefully.

#### UI
- [ ] All issue/task requirements work as expected on the UI
- [ ] Screenshots/Video of the functionality flow (if applicable) - Mobile (X < 640px)
- [ ] Screenshots/Video of the functionality flow (if applicable) - Table (640px > X < 1024px)
- [ ] Screenshots/Video of the functionality flow (if applicable) - Desktop (X > 1024px)
- [ ] Ensure new entries are added to [CHANGELOG.md](https://github.com/prowler-cloud/prowler/blob/master/ui/CHANGELOG.md), if applicable.

#### API
- [ ] All issue/task requirements work as expected on the API
- [ ] Endpoint response output (if applicable)
- [ ] EXPLAIN ANALYZE output for new/modified queries or indexes (if applicable)
- [ ] Performance test results (if applicable)
- [ ] Any other relevant evidence of the implementation (if applicable)
- [ ] Verify if API specs need to be regenerated.
- [ ] Check if version updates are required (e.g., specs, Poetry, etc.).
- [ ] Ensure new entries are added to [CHANGELOG.md](https://github.com/prowler-cloud/prowler/blob/master/api/CHANGELOG.md), if applicable.

### License

By submitting this pull request, I confirm that my contribution is made under the terms of the Apache 2.0 license.
````

---

## Git Commands: Feature Branch

### Setting Up the Chain

```bash
# 1. Create feature branch from master
git checkout master
git pull origin master
git checkout -b feat/PROWLER-XXX-feature-name

# 2. IMPORTANT: Create an empty commit (required for PR creation)
#    Without this, GitHub won't allow creating PRs to this branch
git commit --allow-empty -m "chore: init feature branch for PROWLER-XXX"

# 3. Push feature branch
git push -u origin feat/PROWLER-XXX-feature-name
```

> **Why the empty commit?** GitHub requires at least one commit difference to create a PR. The empty commit establishes the feature branch as a valid PR target.

### Creating a Chained PR

```bash
# 1. Create sub-task branch FROM feature branch
git checkout feat/PROWLER-XXX-feature-name
git pull origin feat/PROWLER-XXX-feature-name
git checkout -b feat/PROWLER-YYY-subtask-name

# 2. Make changes, commit
git add .
git commit -m "feat(scope): implement X for PROWLER-YYY"

# 3. Push and create PR targeting FEATURE BRANCH
git push -u origin feat/PROWLER-YYY-subtask-name

# 4. Create PR with gh CLI (targets feature branch!)
gh pr create \
  --base feat/PROWLER-XXX-feature-name \
  --title "[CHAIN] feat(scope): sub-task title" \
  --body "$(cat <<'EOF'
{PR description from template}
EOF
)"
```

### After a Chained PR is Merged

```bash
# Update feature branch
git checkout feat/PROWLER-XXX-feature-name
git pull origin feat/PROWLER-XXX-feature-name

# Start next sub-task
git checkout -b feat/PROWLER-ZZZ-next-subtask
```

### Creating the Main PR

```bash
# Ensure feature branch is up to date
git checkout feat/PROWLER-XXX-feature-name
git pull origin feat/PROWLER-XXX-feature-name

# Rebase on master to resolve conflicts early
git fetch origin master
git rebase origin/master

# Push and create Main PR
git push origin feat/PROWLER-XXX-feature-name

gh pr create \
  --base master \
  --title "feat(scope): Feature title" \
  --body "$(cat <<'EOF'
{Main PR description from template}
EOF
)"
```

---

# Approach 2: Stacked PRs Direct to Master

## Overview

Instead of a feature branch, use **stacked PRs** that go directly to master. The key insight: it's not a branch dependency, it's a **merge order dependency**.

**Goal:** Code to master as fast as possible, in reviewable chunks, without breaking anything.

- Small PRs → fast reviews → quick merge → fewer conflicts
- Order dependencies, not branch dependencies → each PR is independent against master
- QA post-merge in staging → deploy to prod when validated

### The Pattern

```
master ←── PR1: Endpoint 1 (base: master)
              └── PR2: Endpoint 2 (base: PR1)
                    ↓ PR1 merges → PR2 retargets to master
master ←── PR2: Endpoint 2 (base: master, after rebase)

master ←── PR: UI (base: master, blocked by PR1 & PR2)
```

**Key Rule:** Backend stacks on itself, frontend waits on master. Merge order matters.

---

## Stacked PR Template

**Use when:** Creating a PR that is part of a stack going directly to master.

### Title Format

```
{conventional commit type}({scope}): {description}
```

**Examples:**
- `feat(api): create GET /findings/groups endpoint`
- `feat(api): add POST /findings/groups/export endpoint`
- `feat(ui): implement Findings Groups view`

### Description Template

> **IMPORTANT: Public Repository**
> Do NOT include links to private resources (Jira, Notion, Figma) in PR descriptions.
> Only reference other GitHub PRs and public information.

````markdown
## Description

{What this PR does}

---

## Dependencies

- [ ] Blocked by #123 (endpoint-1)
- [ ] Blocked by #124 (endpoint-2)

⚠️ Do not merge until dependencies are in master

---

## Type

- [ ] Part of a stack (change base after previous PR merges)
- [ ] Final PR in stack (ready to merge when dependencies are done)

---

## Notes for reviewer

This PR is part of a stacked PR workflow. Please review only the changes in this PR, not the dependencies.

---

### Steps to review

1. {Step 1 - what to check first}
2. {Step 2 - how to verify functionality}
3. {Step 3 - edge cases to test}

---

### Checklist

<details>

<summary><b>Community Checklist</b></summary>

- [x] This feature/issue is listed in [here](https://github.com/prowler-cloud/prowler/issues?q=sort%3Aupdated-desc+is%3Aissue+is%3Aopen) or roadmap.prowler.com
- [x] Is it assigned to me, if not, request it via the issue/feature in [here](https://github.com/prowler-cloud/prowler/issues?q=sort%3Aupdated-desc+is%3Aissue+is%3Aopen) or [Prowler Community Slack](goto.prowler.com/slack)

</details>

- [ ] Review if the code is being covered by tests.
- [ ] Review if code is being documented following this specification https://github.com/google/styleguide/blob/gh-pages/pyguide.md#38-comments-and-docstrings
- [ ] Review if backport is needed.
- [ ] Review if is needed to change the [Readme.md](https://github.com/prowler-cloud/prowler/blob/master/README.md)
- [ ] Ensure new entries are added to [CHANGELOG.md](https://github.com/prowler-cloud/prowler/blob/master/prowler/CHANGELOG.md), if applicable.

#### SDK/CLI
- Are there new checks included in this PR? Yes / No
    - If so, do we need to update permissions for the provider? Please review this carefully.

#### UI
- [ ] All issue/task requirements work as expected on the UI
- [ ] Screenshots/Video of the functionality flow (if applicable) - Mobile (X < 640px)
- [ ] Screenshots/Video of the functionality flow (if applicable) - Table (640px > X < 1024px)
- [ ] Screenshots/Video of the functionality flow (if applicable) - Desktop (X > 1024px)
- [ ] Ensure new entries are added to [CHANGELOG.md](https://github.com/prowler-cloud/prowler/blob/master/ui/CHANGELOG.md), if applicable.

#### API
- [ ] All issue/task requirements work as expected on the API
- [ ] Endpoint response output (if applicable)
- [ ] EXPLAIN ANALYZE output for new/modified queries or indexes (if applicable)
- [ ] Performance test results (if applicable)
- [ ] Any other relevant evidence of the implementation (if applicable)
- [ ] Verify if API specs need to be regenerated.
- [ ] Check if version updates are required (e.g., specs, Poetry, etc.).
- [ ] Ensure new entries are added to [CHANGELOG.md](https://github.com/prowler-cloud/prowler/blob/master/api/CHANGELOG.md), if applicable.

### License

By submitting this pull request, I confirm that my contribution is made under the terms of the Apache 2.0 license.
````

---

## Git Commands: Stacked to Master

### Creating the Backend Stack

```bash
# PR1: First endpoint
git checkout master
git pull origin master
git checkout -b feat/PROWLER-XXX-endpoint-1
# ... work ...
git push origin feat/PROWLER-XXX-endpoint-1
# Open PR in GitHub → base: master
```

```bash
# PR2: Second endpoint (stacked on PR1)
git checkout feat/PROWLER-XXX-endpoint-1
git checkout -b feat/PROWLER-XXX-endpoint-2
# ... work ...
git push origin feat/PROWLER-XXX-endpoint-2
# Open PR in GitHub → base: feat/PROWLER-XXX-endpoint-1
```

### Creating the Frontend PR

```bash
git checkout master
git pull origin master
git checkout -b feat/PROWLER-XXX-ui
# ... work ...
git push origin feat/PROWLER-XXX-ui
# Open PR in GitHub → base: master
```

### Sequential Backend Merge

1. **PR1 approved** → Merge to master

2. **PR2: Change base manually in GitHub**
   - Go to PR2 → Edit → Change base → `master`
   - GitHub will show only PR2 changes

3. **PR2 probably needs rebase:**
```bash
git checkout feat/PROWLER-XXX-endpoint-2
git fetch origin
git rebase origin/master
# If conflicts, resolve them
git push --force-with-lease
```

4. **PR2 approved** → Merge to master

### Frontend Merge

```bash
git checkout feat/PROWLER-XXX-ui
git fetch origin
git rebase origin/master
# Now you have the backend
git push --force-with-lease
```

Approved → Merge to master

---

## Corner Cases (Stacked PRs)

### PR1 needs changes after review

```bash
git checkout feat/PROWLER-XXX-endpoint-1
# Make changes
git add .
git commit -m "fix: address review comments"
git push origin feat/PROWLER-XXX-endpoint-1
```

PR2 is now outdated:

```bash
git checkout feat/PROWLER-XXX-endpoint-2
git rebase feat/PROWLER-XXX-endpoint-1
git push --force-with-lease
```

Frontend is not affected, still on master.

### PR2 needs changes that affect PR1

If PR1 is not merged:

```bash
git checkout feat/PROWLER-XXX-endpoint-1
# Make the fix
git add .
git commit -m "fix: add missing field"
git push origin feat/PROWLER-XXX-endpoint-1

# Rebase PR2
git checkout feat/PROWLER-XXX-endpoint-2
git rebase feat/PROWLER-XXX-endpoint-1
git push --force-with-lease
```

If PR1 is already merged:

```bash
# Option A: Fix in PR2
git checkout feat/PROWLER-XXX-endpoint-2
# Add the fix there

# Option B: New hotfix PR
git checkout master
git pull origin master
git checkout -b fix/PROWLER-XXX-endpoint-1-missing-field
# Quick fix, PR direct to master
```

### Conflicts in PR2 after merging PR1

```bash
git checkout feat/PROWLER-XXX-endpoint-2
git fetch origin
git rebase origin/master
# Resolve file by file
git add .
git rebase --continue
git push --force-with-lease
```

Then in GitHub: Edit PR → Change base → `master`

### Frontend has conflicts when backend merges

```bash
git checkout feat/PROWLER-XXX-ui
git fetch origin
git rebase origin/master
# Resolve conflicts
git add .
git rebase --continue
git push --force-with-lease
```

### Need to revert backend after merging

```bash
git checkout master
git pull origin master
git revert <commit-hash-from-pr1>
git push origin master
```

- If PR2 not merged: decide to rebase on the revert or wait for the fix
- Frontend stays blocked, didn't reach production (advantage of this approach)

### Forgot to change PR2 base in GitHub

After merging PR1, PR2 shows inflated diff (PR1 + PR2 commits).

Fix: GitHub → PR2 → Edit → Change base to `master`

### Force push was rejected

```bash
git push --force-with-lease
# ! [rejected] (stale info)
```

Someone else pushed to that branch:

```bash
git fetch origin
git rebase origin/feat/PROWLER-XXX-endpoint-2
# Resolve if there are conflicts
git push --force-with-lease
```

### See PR2 diff only (without PR1) before merging PR1

```bash
git log feat/PROWLER-XXX-endpoint-1..feat/PROWLER-XXX-endpoint-2 --oneline
git diff feat/PROWLER-XXX-endpoint-1..feat/PROWLER-XXX-endpoint-2
```

Or wait for PR1 to merge and change base to master.

---

## Workflow Checklists

### Feature Branch Approach

#### Before Starting
- [ ] Feature branch exists: `feat/PROWLER-XXX-{name}`
- [ ] Feature branch has empty init commit (required for PRs)
- [ ] Jira Story has sub-tasks with dependencies defined
- [ ] Team aligned on which sub-tasks can run in parallel

#### For Each Chained PR
- [ ] Branch created FROM feature branch
- [ ] PR targets feature branch (NOT master)
- [ ] Title has `[CHAIN]` prefix
- [ ] Description includes chain position and dependencies
- [ ] Linked to Jira sub-task

#### Before Main PR
- [ ] All chained PRs merged to feature branch
- [ ] Feature branch rebased on latest master
- [ ] End-to-end testing done on feature branch
- [ ] All Jira sub-tasks marked Done

#### Main PR
- [ ] PR targets master
- [ ] Description lists all included PRs
- [ ] Has screenshots/demo
- [ ] All checklists completed

### Stacked to Master Approach

#### Before Merging Each Stacked PR
- [ ] Previous PR already merged?
- [ ] Changed base to `master` in GitHub?
- [ ] Rebased on master?
- [ ] CI passes?
- [ ] Diff shows only YOUR changes?

#### For Dependent PRs (e.g., Frontend)
- [ ] All dependency PRs merged to master?
- [ ] Rebased on latest master (now includes dependencies)?
- [ ] Integration works with real backend?
- [ ] CI passes?

---

## Useful Git Commands

```bash
# See your branches status vs master
git log master..feat/PROWLER-XXX-endpoint-1 --oneline
git log master..feat/PROWLER-XXX-endpoint-2 --oneline

# See what commits PR2 has that PR1 doesn't
git log feat/PROWLER-XXX-endpoint-1..feat/PROWLER-XXX-endpoint-2 --oneline

# Interactive rebase if you need to clean up commits
git rebase -i feat/PROWLER-XXX-endpoint-1

# Abort a rebase that got complicated
git rebase --abort

# See branches already merged
git branch --merged master
```

---

## Related Skills

- **jira-task** - Create Jira sub-tasks with proper dependencies
- **jira-epic** - Create the parent Epic for the feature

---

## Keywords

pr, pull request, chained, stacked, feature branch, github, merge, chain, parallel, api, ui, stacked prs, direct to master
