# Tasks: Fix the Gentleman.Dots installer on Arch Linux / CachyOS

## Review Workload Forecast

| Field | Value |
|-------|-------|
| Estimated changed lines | 340-485 |
| 400-line budget risk | High |
| Chained PRs recommended | Yes |
| Suggested split | PR 1 (unit 1) -> PR 2 (units 2+3) -> PR 3 (unit 4) |
| Delivery strategy | auto-chain |
| Chain strategy | feature-branch-chain |

Decision needed before apply: No
Chained PRs recommended: Yes
Chain strategy: feature-branch-chain
400-line budget risk: High

### Suggested Work Units

| Unit | Goal | Likely PR | Focused test command | Runtime harness | Rollback boundary |
|------|------|-----------|----------------------|-----------------|-------------------|
| 1 | Arch pkg lists resolve | PR 1 (base: tracker) | `go test ./internal/tui/... -run TestPlatformPackages -skip Golden` | N/A (data-only, no exec) | Revert `installer.go:720,775,831` + test file |
| 2+3 | CopyDir skip + Notify seam + fail-closed delete guard | PR 2 (base: PR 1 branch) | `go test ./internal/system/... -skip Golden` | `net.Listen("unix",...)` fixture; real `git init` fixture in `t.TempDir()` | Revert `report.go`, `repostate.go`, `CopyDir` diff, both delete-site wiring |
| 4 | Real `--dry-run` gate | PR 3 (base: PR 2 branch) | `go test ./... -skip Golden` (installer/) | Full flow with `GENTLEMAN_DRY_RUN=1` on a dirty home dir | Revert `dryrun.go` + gate additions in `exec.go`/`installer.go`; units 1-3 remain valid |

## Phase 1: Unit 1 - Arch package lists (PR 1)

- [x] 1.1 RED: update `internal/tui/platform_packages_test.go:111-131` expected Arch strings to `"fish zoxide atuin starship"`, `"zsh zoxide atuin zsh-autosuggestions zsh-syntax-highlighting zsh-autocomplete"`, `"nushell zoxide atuin jq bash starship"` (arch-package-install: fish/zsh/nushell scenarios).
- [x] 1.2 GREEN: edit `internal/tui/installer.go:720,775,831` to remove `carapace` (fish, nushell) and `carapace`+`zsh-theme-powerlevel10k` (zsh).

## Phase 2: Unit 2 - CopyDir non-regular skip + Notify seam (PR 2)

- [x] 2.1 RED: `internal/system/exec_test.go` - dir with `net.Listen("unix",...)` socket + regular files backed up via `CreateBackup`; assert success, socket path in captured `Notify`, regular files copied (config-backup: socket scenario).
- [x] 2.2 RED: dir with socket + FIFO + regular files; assert both skipped with individual warnings, backup completes (config-backup: multiple non-regular scenario).
- [x] 2.3 RED: nominal all-regular-files backup emits zero skip warnings (config-backup: nominal scenario).
- [x] 2.4 GREEN: create `internal/system/report.go` with `var Notify = func(msg string){ fmt.Println(msg) }` and `report(logFunc func(string), msg string)`.
- [x] 2.5 GREEN: `CopyDir` Walk callback in `exec.go` checks `info.Mode().IsRegular()`; non-regular entries call `report(nil, "⚠ skipped non-regular file: "+path)` and continue instead of aborting.
- [x] 2.6 RED: TUI summary screen lists every skipped path from a run (config-backup: summary visibility scenario).
- [x] 2.7 GREEN: TUI aggregates skipped-path messages (captured via `system.Notify` override) into the final install summary.

## Phase 3: Unit 3 - Fail-closed delete guard (PR 2, same branch as unit 2)

- [x] 3.1 RED: `internal/system/repostate_test.go` - `InspectRepoDir` returns `RepoAbsent`/`RepoCleanCheckout`/`RepoDirtyCheckout`/`RepoNotAGitRepo`/`RepoUnknown` for absent dir, real `git init` clean fixture, dirty (unstaged/staged/untracked) fixture, non-git dir, and `runGitStatus` override returning an error (repo-clone-safety: all five scenarios).
- [x] 3.2 GREEN: create `internal/system/repostate.go` with `RepoState` enum, `var runGitStatus`, and `InspectRepoDir(path string) RepoState`.
- [x] 3.3 RED: assert `runGitStatus` receives an absolute path when `InspectRepoDir` is called with a relative one.
- [x] 3.4 GREEN: `InspectRepoDir` calls `filepath.Abs` before invoking `runGitStatus`.
- [x] 3.5 RED: `stepCloneRepo` refuses (hard `wrapStepError`, actionable message naming absolute path) on dirty/not-a-repo/unknown; proceeds silently on absent/clean (repo-clone-safety: directory-exists scenarios).
- [x] 3.6 RED: `stepCleanup` warns and returns nil (non-critical) on the same refusal states, deletes on absent/clean.
- [x] 3.7 RED: non-interactive run hitting a dirty checkout aborts deterministically with no prompt (repo-clone-safety: non-interactive scenario).
- [x] 3.8 GREEN: implement `removeRepoDirIfSafe(stepID, path string) error` in `internal/tui/installer.go`; wire into both `stepCloneRepo` and `stepCleanup` delete sites, replacing the unguarded `rm -rf`.

## Phase 4: Unit 4 - Real `--dry-run` (PR 3)

- [x] 4.1 RED: `internal/system/dryrun_test.go` - with `GENTLEMAN_DRY_RUN=1` (`t.Setenv`), `Run` and `RunWithLogs` execute nothing (zero process spawn), return `ExecResult{ExitCode:0, Error:nil}`, and emit `[dry-run] would run: <cmd>` via `Notify`/`logFunc` (dry-run-safety: central-gate scenario).
- [x] 4.2 RED: same zero-spawn assertion individually for all 7 wrappers: `RunSudo`, `RunBrew`, `RunPkg`, `RunPkgWithLogs`, `RunPkgInstall`, `RunBrewWithLogs`, `RunSudoWithLogs` (threat-matrix: subprocess composition scenario).
- [x] 4.3 GREEN: create `internal/system/dryrun.go` with `IsDryRun()` and a synthetic-`ExecResult` helper; gate `Run` and `RunWithLogs` in `exec.go` (covers all 7 wrappers by delegation).
- [x] 4.4 RED: dry-run `CreateBackup` creates nothing on disk (`os.Stat` on backup dir is `IsNotExist`) yet returns a non-empty `backupDir` (dry-run-safety: filesystem-mutation scenario).
- [x] 4.5 RED: dry-run `CopyFile`/`CopyDir` perform no `os.WriteFile`/`os.MkdirAll`, `CopyDir` still walks and prints the plan via `Notify`.
- [x] 4.6 GREEN: add `IsDryRun()` guards in `CopyFile`, `CopyDir`, `CreateBackup` (`exec.go`) that skip mutations, still compute `GetBackupDir()`, and `report()` the plan.
- [x] 4.7 RED: dry-run `stepCloneRepo` performs no `git clone`, no directory create/remove, and the post-clone `os.Stat` verification at `installer.go:132` is skipped (does not report false failure) (dry-run-safety: no-clone scenario).
- [x] 4.8 GREEN: gate the `installer.go:132` post-clone `os.Stat` check behind `IsDryRun()`.
- [x] 4.9 RED: dry-run delete sites (`stepCloneRepo`, `stepCleanup`) print the planned removal instead of calling `rm -rf`.
- [x] 4.10 GREEN: wire `IsDryRun()` branch into `removeRepoDirIfSafe` to print-and-skip instead of deleting.
- [x] 4.11 RED: `--test` mode (not `--dry-run`) still executes pacman/sudo/network operations (dry-run-safety: not-mistaken-for-test scenario).
- [x] 4.12 GREEN: confirm `IsDryRun()` reads only `GENTLEMAN_DRY_RUN`, never `--test`'s flag/env.
- [x] 4.13 GREEN: wire `system.Notify = SendLog` once at TUI install start in `internal/tui/installer.go`; keep stdout default for non-interactive mode.

## Phase 5: Verification

- [x] 5.1 Run `go test ./... -skip Golden` from `installer/`; confirm only the three documented macOS-golden baseline failures remain (never treated as regressions).
- [x] 5.2 Run `go build ./cmd/gentleman-installer` from `installer/`.
- [x] 5.3 Run the Docker E2E matrix (`installer/e2e/docker-test.sh`, ubuntu/debian/fedora/alpine/termux) to confirm unit 4's shared `exec.go` change is byte-identical off Arch.
