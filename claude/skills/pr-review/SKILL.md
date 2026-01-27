---
name: pr-review
description: >
  Review GitHub PRs and Issues with structured analysis for open source projects.
  Trigger: When reviewing PRs, analyzing issues, or evaluating contributions to repositories.
license: MIT
metadata:
  author: gentleman-programming
  version: "1.1"
---

## When to Use

- Reviewing open PRs in a repository
- Analyzing GitHub issues for validity/priority
- Evaluating external contributions before merge
- Auditing PR quality across a project

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

### Phase 2: Read Current Codebase

Before reviewing diffs, **always read the current code** to understand context:
- Main entry points
- Files being modified
- Related modules

### Phase 3: Analyze Each PR

For each PR, evaluate these factors:

| Factor | What to Check |
|--------|---------------|
| **Code Quality** | Clean code, no duplication, proper error handling |
| **Tests** | Are there tests? Do they cover the changes? |
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

### Comment Structure

```markdown
Hey {Name}! / ¡Hola {Name}! {Positive feedback about the PR}

{Brief context if needed}

## {Problem Category}

{Explanation of the issue with code example}

## Solution

```bash
{Concrete fix}
```

---

{Closing - what happens after they fix it}
```

**Examples:**
- Spanish PR → "¡Hola Juan! Gracias por el PR, el análisis está muy bien hecho 👏"
- English PR → "Hey John! Thanks for the PR, the analysis is well done 👏"

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
