# Repo Clone Safety Specification

## Purpose

Defines safety requirements around destructive removal of the `Gentleman.Dots` working directory at `stepCloneRepo` (before cloning) and `stepCleanup` (after installation), preventing data loss on a dirty checkout while remaining scriptable for CI/CD (`--non-interactive`).

## Requirements

### Requirement: Destructive removal MUST be git-status-aware and MUST fail closed

Before removing an existing `Gentleman.Dots` directory at either `stepCloneRepo` or `stepCleanup`, the installer MUST determine that directory's git status. It MUST refuse to delete and MUST abort the step with an actionable error message (naming the path and instructing the user to move, rename, or commit it) whenever the directory is a git repository with uncommitted changes or untracked files, or whenever its git status cannot be determined for any reason (git binary missing, not a git repository, or a git command error). The installer MUST NOT prompt interactively for confirmation to override this refusal, because `--non-interactive` mode MUST remain fully scriptable.

#### Scenario: Directory does not exist

- GIVEN no `Gentleman.Dots` directory exists at the relevant relative path
- WHEN `stepCloneRepo` or `stepCleanup` runs
- THEN the step MUST proceed without any git-status check or deletion attempt

#### Scenario: Directory exists but is not a git repository

- GIVEN a `Gentleman.Dots` directory exists at the relevant path
- AND it is not a git repository (no `.git`)
- WHEN the step attempts removal
- THEN the installer MUST refuse to delete it
- AND MUST abort the step with a message instructing the user to move, rename, or commit the directory

#### Scenario: Directory exists and is a clean git repository

- GIVEN a `Gentleman.Dots` directory exists at the relevant path
- AND `git status --porcelain` reports no uncommitted changes and no untracked files
- WHEN the step attempts removal
- THEN the installer MUST proceed to delete the directory without prompting

#### Scenario: Directory exists and is a dirty git repository

- GIVEN a `Gentleman.Dots` directory exists at the relevant path
- AND `git status --porcelain` reports uncommitted changes or untracked files
- WHEN the step attempts removal
- THEN the installer MUST refuse to delete it
- AND MUST abort the step with an actionable message
- AND MUST NOT prompt for confirmation to override

#### Scenario: Git is unavailable

- GIVEN a `Gentleman.Dots` directory exists at the relevant path
- AND the `git` binary is not present on PATH, or invoking `git status` returns an error
- WHEN the step attempts removal
- THEN the installer MUST treat the status as undeterminable
- AND MUST refuse to delete the directory
- AND MUST abort the step with an actionable message

### Requirement: The refusal path MUST remain non-interactive-compatible

The installer MUST NOT require a TTY prompt to resolve a dirty-checkout refusal. Refusal MUST be a deterministic abort (non-zero exit / step failure) so `--non-interactive` runs behave predictably in CI/CD.

#### Scenario: Non-interactive run hits a dirty checkout

- GIVEN the installer runs with `--non-interactive`
- AND the `Gentleman.Dots` directory is a dirty git repository
- WHEN `stepCloneRepo` or `stepCleanup` runs
- THEN the installer MUST abort that step deterministically without blocking on user input
