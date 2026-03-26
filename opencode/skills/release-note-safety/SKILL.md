---
name: release-note-safety
description: >
  Safe patterns for GitHub release notes and any shell command that embeds markdown,
  backticks, or multi-line text. Trigger: When creating or editing releases, passing
  markdown to `gh release create/edit`, or writing shell commands that include backticks.
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

## When to Use

- User asks to create or edit a GitHub release
- You use `gh release create` or `gh release edit`
- You pass markdown to any shell command
- The content includes backticks, fenced code blocks, or multi-line notes

## Critical Patterns

### Rule 1: Never pass markdown with backticks inside double-quoted shell strings

This is the failure mode:

```bash
gh release edit v1.2.3 --notes "Use `opus` and `sonnet`"
```

The shell executes the backticks. That is a command-substitution bug, not a markdown problem.

### Rule 2: Prefer single quotes for short literal note bodies

If the text is short and contains no single quotes, use single quotes:

```bash
gh release edit v1.2.3 --notes 'Use `opus` for orchestration and `sonnet` for implementation.'
```

### Rule 3: Prefer a heredoc for multi-line release notes

For anything longer than one short sentence, use a quoted heredoc and command substitution:

```bash
gh release edit v1.2.3 --notes "$(cat <<'EOF'
## Highlights
- Added per-phase Claude model assignments for SDD.
- Defaulted orchestration to `opus`, implementation-heavy phases to `sonnet`.

## Install
```bash
brew upgrade gentle-ai
```
EOF
)"
```

The `<<'EOF'` form is critical. The single-quoted heredoc marker prevents interpolation and command substitution inside the body.

### Rule 4: If using Bash tool, prefer `gh release edit --notes-file`

If the environment and workflow allow a temporary file, the safest option is to write notes to a file and pass `--notes-file`.

### Rule 5: Verify after editing

After `gh release create` or `gh release edit`, always verify the body with:

```bash
gh release view v1.2.3 --json body,url
```

Do not assume the release notes rendered correctly.

## Decision Table

| Situation | Safe approach |
|-----------|---------------|
| One short line, no single quotes | Single-quoted `--notes` |
| Multi-line markdown | Quoted heredoc |
| Complex reusable content | `--notes-file` |
| Content has backticks inside double quotes | Stop and rewrite |

## Commands

```bash
# Safe single-line edit
gh release edit v1.2.3 --notes 'Use `opus` for orchestration and `sonnet` for implementation.'

# Safe multi-line edit
gh release edit v1.2.3 --notes "$(cat <<'EOF'
## Highlights
- Added safer release handling.
EOF
)"

# Verify rendered result
gh release view v1.2.3 --json body,url
```
