# Archive Report: fix-arch-installer

**Date**: 2026-08-29  
**Change Name**: fix-arch-installer  
**Project**: Gentleman.Dots-fix (installer/cmd/gentleman-installer)  
**Branch**: fix/arch-installer-cachyos  
**Status**: Complete and Archived

## Overview

The `fix-arch-installer` change addresses five distinct defects in the Gentleman.Dots installer's Arch Linux / CachyOS support path:

1. **Unit 1**: Arch package lists contained AUR-only package names
2. **Unit 2**: Backup failed silently on Unix domain sockets and FIFOs
3. **Unit 3**: Destructive directory removal was unguarded against data loss
4. **Unit 4**: The `--dry-run` flag was not implemented (no-op)
5. **Unit 5**: Design flaw in Unit 3's guard discovered during verification and fixed

All work is implemented, verified, and ready for delivery.

## Work Units Completed

### Unit 1: Arch Package Lists Sanitization
**Commits**: Included in 3801389  
**Tasks Completed**: 1.1, 1.2 (2/2)  

Removed AUR-only package names (`carapace`, `zsh-theme-powerlevel10k`) from Arch package lists for fish, zsh, and nushell profiles. Modified `/internal/tui/installer.go` (lines 720, 775, 831) to reflect corrected package dependencies that resolve exclusively in official Arch repositories (`core`, `extra`).

**Specs Updated**: `arch-package-install/spec.md` (5 scenarios)
- Requirement: Arch package lists must only contain names resolvable in official Arch repositories
- Requirement: A single unresolvable pacman package name must not block installation of remaining packages

### Unit 2: Backup Safety for Non-Regular Files
**Commits**: Included in 410a9a9  
**Tasks Completed**: 2.1–2.7 (7/7)  

Added `internal/system/report.go` with pluggable `Notify` function for user messaging. Modified `CopyDir` in `exec.go` to detect non-regular files (sockets, FIFOs) and skip them gracefully with individual warnings, enabling backup of mixed directories. Integrated skipped-file reporting into the TUI install summary screen.

**Specs Updated**: `config-backup/spec.md` (5 scenarios)
- Requirement: Backup must handle Unix domain sockets and FIFOs without aborting
- Requirement: Skipped paths must be visible in the install summary

### Unit 3: Fail-Closed Deletion Guard (Initial Implementation)
**Commits**: Included in 410a9a9  
**Tasks Completed**: 3.1–3.8 (7/7)  

Created `internal/system/repostate.go` with `RepoState` enum (RepoAbsent, RepoCleanCheckout, RepoDirtyCheckout, RepoNotAGitRepo, RepoUnknown) and `InspectRepoDir()` function to determine git state before deletion. Implemented `removeRepoDirIfSafe()` in `internal/tui/installer.go` and wired it into both `stepCloneRepo` and `stepCleanup` deletion sites, replacing unguarded `rm -rf`.

**Initial Design**: Refuse deletion on dirty/non-git/unknown states; proceed on absent/clean.

**Specs Updated**: `repo-clone-safety/spec.md` (7 scenarios, initial version)

### Unit 4: Real `--dry-run` Implementation
**Commits**: Included in 410a9a9  
**Tasks Completed**: 4.1–4.13 (13/13)  

Created `internal/system/dryrun.go` with `IsDryRun()` gate reading `GENTLEMAN_DRY_RUN` environment variable. Guarded all process spawning in `exec.go` (wrappers: RunSudo, RunBrew, RunPkg, RunPkgWithLogs, RunPkgInstall, RunBrewWithLogs, RunSudoWithLogs) to emit synthetic `ExecResult{ExitCode:0, Error:nil}` and print planned commands via `Notify`. Extended dry-run gates to `CreateBackup`, `CopyFile`, `CopyDir` (skip mutations, report plans), and `stepCloneRepo` (skip clone and post-clone `os.Stat` verification). Wired `system.Notify` callback into TUI at install start to capture all plan messages.

**Specs Updated**: `dry-run-safety/spec.md` (5 scenarios)
- Requirement: Dry-run must not spawn processes or mutate filesystem
- Requirement: Dry-run plans must be reported via log/notify for user visibility

### Unit 5: Deletion Guard Design Flaw Fix (POST-VERIFICATION)
**Commits**: b2564a2  
**Context**: Discovered during verification when installer destroyed the change's own commits  
**Root Cause**: Initial guard (`RepoCleanCheckout`) only checked `git status --porcelain` (no uncommitted changes), which proved nothing about unpushed commits or stashes. A clean working tree can hold commits that do not exist on any remote or stash entries. Clean status alone did not authorize deletion.

**Solution**:
- Extended `RepoState` enum with `RepoUnpublishedWork` state
- Modified deletion guard to additionally verify:
  - No stash entries (`git stash list` is empty)
  - `HEAD` is contained by at least one remote-tracking branch (`git branch -r --contains HEAD` finds at least one match)
- Deletion now permitted only when ALL conditions hold:
  - Working tree is clean (no staged/unstaged changes)
  - No untracked files
  - No stash entries
  - `HEAD` is published (reachable from a remote branch)
- Special case: Branch with no upstream is allowed if `HEAD` is already published elsewhere
- Non-interactive refusal remains deterministic (no TTY prompt)

**Specs Updated**: `repo-clone-safety/spec.md` (expanded to 10 scenarios)
- Added Requirement: A clean working tree MUST NOT by itself authorise deletion
- Added Scenario: Clean checkout holding unpushed commits
- Added Scenario: Clean checkout with no remote configured
- Added Scenario: Clean checkout holding a stash
- Added Scenario: Branch with no upstream whose HEAD is already published

**Code Changes**:
- `internal/system/repostate.go`: `RepoUnpublishedWork` enum value, publication verification logic
- `internal/system/repostate_test.go`: New test scenarios for stash, unpushed commits, no-remote cases
- `internal/tui/installer.go`: Updated `removeRepoDirIfSafe()` to use enhanced guard

## Task Completion Status

**Total Tasks**: 33  
**Completed**: 33/33 (100%)

- Phase 1 (Unit 1): 2/2 completed ✓
- Phase 2 (Unit 2): 7/7 completed ✓
- Phase 3 (Unit 3): 7/7 completed ✓
- Phase 4 (Unit 4): 13/13 completed ✓
- Phase 5 (Unit 5): 3/3 completed ✓
- Phase 6 (Verification): 5.1–5.3 marked complete ✓

## Verification Status

**Test Results**:
- `go test ./... -skip Golden` from `installer/`: PASS
  - Only three documented macOS-recorded golden baseline failures remain (pre-existing, unrelated to feature changes)
  - All implementation tests pass on Linux (CachyOS/Arch)

**Build Status**:
- `go build ./cmd/gentleman-installer`: SUCCESS

**Code Quality**:
- `gofmt -l .`: Clean except pre-existing `internal/system/detect.go` formatting issue

**Docker E2E Matrix**:
- alpine: PASS ✓
- debian: PASS ✓
- ubuntu: PASS ✓
- fedora: PASS ✓
- termux: Reported PASS by harness, but inner suite had 12 failures (pre-existing on clean main at commit 0258450, not caused by this change; identified as defect in `docker-test.sh` exit-status handling)

## Specs Merged into Main Specs

Four new domain specs created and synced:

| Domain | Action | Requirements | Scenarios |
|--------|--------|---|---|
| `arch-package-install` | Created | 2 | 5 |
| `config-backup` | Created | 2 | 5 |
| `repo-clone-safety` | Created | 3 | 10 (includes Unit 5 addition) |
| `dry-run-safety` | Created | 2 | 5 |

**Total Requirements**: 9  
**Total Scenarios**: 25

All specs now live in `openspec/specs/{domain}/spec.md` and serve as source of truth for ongoing maintenance and future enhancements.

## Commits on Branch

| Hash | Message Summary |
|------|---|
| 3801389 | Unit 1–4 implementation (arch packages, backup guard, deletion guard, dry-run) |
| 410a9a9 | OpenSpec artifacts (proposal, design, specs, tasks) |
| b2564a2 | Unit 5: Deletion guard enhancement (unpublished-work check, stash detection, remote-branch publication verification) |

Nothing is pushed; all commits remain on `fix/arch-installer-cachyos` for orchestrator delivery.

## Issues and Defects Referenced

**Upstream Issues Addressed**:
- #182: AUR-only Arch packages (Unit 1)
- #174: AUR-only Arch packages (Unit 1)
- #190: dry-run is a no-op (Unit 4)
- #192: No Homebrew fallback on Arch (Unit 1)

**Additional Defects Identified (User Decision: NOT Reported Upstream)**:
1. **Backup dying on Unix sockets**: Fixed in Unit 2 (CopyDir now skips sockets gracefully)
2. **Destructive relative-path `rm -rf`**: Fixed in Unit 3–5 (deletion now guarded with publication verification)
3. **E2E harness false green (termux)**: `docker-test.sh` exits 0 even when per-image suites fail (pipe to `tee` masks exit status); identified as separate harness defect, not part of installer logic

User explicitly chose not to open upstream issues for the three additional defects discovered during implementation and verification.

## Archive Integrity

**Source Snapshot**: Captured before move operation  
**Archive Path**: `openspec/changes/archive/2026-08-29-fix-arch-installer/`  
**Verification**: `diff -r` confirms all files (proposal.md, design.md, tasks.md, specs/) copied intact  
**Archive Contents**:
- proposal.md ✓
- design.md ✓
- tasks.md ✓ (33/33 tasks marked complete)
- specs/ (4 domain subdirectories with spec.md files) ✓
- archive-report.md ✓

## Key Learnings

1. A clean git working tree is not sufficient proof of safety for deletion; stash entries and unpushed commits can hide work that would be destroyed.
2. Design flaws in safety guards can only be discovered when the guard is actually exercised; verification (running the installer in a realistic scenario) exposed the incomplete deletion criteria.
3. The terminal definition of "published" requires three checks: clean status, no stash, and HEAD contained by a remote branch; omitting any one opens a data-loss vulnerability.
4. E2E test harnesses must preserve per-platform exit status; masking it behind aggregators like `tee` can silently hide failures in specific environments.
5. Distinguishing pre-existing defects from new work during verification preserves change integrity; documenting non-reported defects and the reasons protects decision rationale.

## SDD Cycle Closure

**Change Status**: COMPLETE AND ARCHIVED  
**Date Archived**: 2026-08-29  
**Ready for**: Orchestrator delivery (apply and publish to main branch)

The change is fully specified, implemented, tested, verified, and archived. All artifacts (proposal, design, specs, tasks) are preserved in the archive for future reference. Specs are now integrated into the main `openspec/specs/` tree and will serve as baseline requirements for all related future work on Arch/CachyOS support and deletion safety.
