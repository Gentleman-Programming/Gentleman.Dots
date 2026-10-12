# Design: Fix the Gentleman.Dots installer on Arch Linux / CachyOS

## Technical Approach

Four chained work units, each shippable alone. Two are data-only edits (package lists, `CopyDir` mode check); two introduce small seams in `internal/system` that the TUI wires once. No call-signature threading, no new dependency, no change to `detect.go`. Every seam uses the package-level var-function pattern the codebase already relies on (`runPkgInstallWithLogs`, `runSudoWithLogs`, `runBrewWithLogs` at `internal/tui/installer.go:609-613`), so tests override and restore, and no interface or DI container is invented.

## Architecture Decisions

### Decision: dry-run gate lives at the two base exec functions, not at call sites

| Option | Tradeoff | Decision |
|---|---|---|
| Guard all four `Run*` funcs | `RunSudo` → `Run` and `RunSudoWithLogs` → `RunWithLogs`; guarding all four double-prints the `sudo ` command | Rejected |
| Guard only `Run` and `RunWithLogs` | Single chokepoint; wrappers pass the fully composed string (`sudo pacman -S …`), which is exactly the planned operation to print | **Chosen** |
| Thread a `dryRun bool` through signatures | Touches every step function; disproportionate (proposal non-goal) | Rejected |

`system.IsDryRun()` reads `GENTLEMAN_DRY_RUN` once via `sync.Once`-free direct lookup so tests can set/unset it. On a dry run the gate emits the planned operation and returns a synthetic success `*ExecResult{Command: command, ExitCode: 0, Error: nil}` — callers keep their existing success path.

### Decision: one reporter seam serves both "print planned operations" and "skipped file" warnings

Upstream #190 requires dry-run to **print** planned operations; `--test` is not a substitute. `Run` has no log callback, `RunWithLogs` does. So:

```go
// internal/system/report.go
var Notify = func(msg string) { fmt.Println(msg) } // overridable in tests and by the TUI

func report(logFunc func(string), msg string) { if logFunc != nil { logFunc(msg) } else { Notify(msg) } }
```

`RunWithLogs` routes planned operations to its own `logFunc` (so they land in the TUI step log); `Run`, `CopyFile`, `CopyDir`, `CreateBackup` and the delete helper route to `Notify`. Prefixes: `[dry-run] would run: <cmd>`, `[dry-run] would copy: <src> -> <dst>`, `⚠ skipped non-regular file: <path>`. The TUI sets `system.Notify` once at install start to forward into `SendLog`; non-interactive mode keeps the stdout default.

Rationale: one seam, one place to assert in tests, and skip-warnings become user-visible without changing `CopyDir`/`CreateBackup` signatures (every existing caller and test stays compiling).

### Decision: dry-run suppresses *effects*, never *derivations* — the state-correctness rule

This is the answer to the biggest stated risk. Skipping a mutation must not skip the value later steps read.

| Site | Effect skipped | Derivation preserved |
|---|---|---|
| `CreateBackup` (exec.go:442) | `EnsureDir`, `CopyDir`, `CopyFile`, the `RemoveAll` deferred cleanup | still computes `GetBackupDir()` and returns `(backupDir, nil)` → `installer.go:100` sets `m.BackupDir` to a real, non-empty path |
| `CopyFile` / `CopyDir` | `os.WriteFile`, `os.MkdirAll` | returns `nil`; `CopyDir` still walks to enumerate what it *would* copy so the printed plan is real |
| `stepCloneRepo` (installer.go:105-140) | `rm -rf`, `git clone` | the post-clone `os.Stat` verification at :132 MUST be skipped under dry-run, otherwise a dry run reports a false failure |
| `stepCleanup` (~1195) | `rm -rf` | already non-fatal; prints plan and returns nil |

Rule for tasks/tests: any `os.Stat`/existence check that verifies the result of a mutation the gate skipped must itself be gated. This is the only class of dry-run regression that unit tests can miss, so it gets an explicit RED test per site.

### Decision: destructive delete is fail-closed and shared by both sites

`docs/tui-installer.md:27` documents `--non-interactive` as CI/CD friendly, so a confirmation prompt is not acceptable — the guard decides on its own and refuses.

```go
// internal/system/repostate.go
type RepoState int
const (RepoAbsent RepoState = iota; RepoCleanCheckout; RepoDirtyCheckout; RepoNotAGitRepo; RepoUnknown)

var runGitStatus = func(dir string) (stdout string, err error) { /* git -C <dir> status --porcelain */ }

// InspectRepoDir resolves path with filepath.Abs first (the defect was a relative rm -rf),
// then classifies. git missing, non-zero exit, or unparseable output => RepoUnknown.
func InspectRepoDir(path string) RepoState
```

Delete policy, identical at both sites via one TUI helper `removeRepoDirIfSafe(stepID, path string) error`:

| State | Action |
|---|---|
| `RepoAbsent` | no-op, proceed |
| `RepoCleanCheckout` | delete, proceed |
| `RepoDirtyCheckout` / `RepoNotAGitRepo` / `RepoUnknown` | **refuse**; abort the step with an actionable error naming the absolute path and the two exits (move the directory, or run the installer from another cwd) |

`stepCloneRepo` treats refusal as a hard `wrapStepError`. `stepCleanup` keeps its existing non-critical semantics: warn and `return nil`, leaving the directory in place — which is the safe outcome.

**Rejected alternative:** cloning into a suffixed directory. The name `Gentleman.Dots` is hardcoded at ~10 functional sites, including three shell-command string literals (`internal/tui/interactive.go:220,235,250`) and four `repoDir :=` assignments (`installer.go:286,703,905,1103`). Not a rename, and it would not protect the user's checkout anyway.

### Decision: package lists — remove two names, Arch fields only

Verified against the official Arch package database (archlinux.org JSON API), not against a CachyOS host whose extra `[cachyos]` repos masked one of them:

| Field | Line | Remove |
|---|---|---|
| fish `Arch` | installer.go:720 | `carapace` |
| zsh `Arch` | installer.go:775 | `carapace`, `zsh-theme-powerlevel10k` |
| nushell `Arch` | installer.go:831 | `carapace` |

The Arch **zsh** list is broken too, not only fish/nushell (upstream #174 reports this independently). `Brew`, `Fedora`, `Debian` and `Termux` fields are untouched. `platform_packages_test.go:117,125` expectations update to match. `system.EnsureDir(~/.cache/carapace)` at installer.go:711 stays — harmless and unrelated.

Not fixed here (proposal non-goal): `runNativeWithBrewFallback` (installer.go:635-642) only falls back when `hasBrew`, so Arch still has no recovery path if pacman fails for another reason (upstream #192). Removing the unresolvable names is what makes the atomic transaction succeed.

## Data Flow

```
--dry-run ─→ main.go:70 GENTLEMAN_DRY_RUN=1
                          │
      ┌───────────────────┴────────────────────┐
      ▼                                        ▼
system.IsDryRun()                       tui sets system.Notify = SendLog
      │
      ├─ Run / RunWithLogs ─→ report("[dry-run] would run: …") ─→ ExecResult{ok}
      ├─ CopyFile/CopyDir ──→ report("[dry-run] would copy: …") ─→ nil
      └─ CreateBackup ──────→ report(plan) ─→ (GetBackupDir(), nil) ─→ m.BackupDir set

stepCloneRepo ─┐                                  ┌─ Absent/Clean ─→ rm -rf (or dry-run print)
               ├─→ removeRepoDirIfSafe(absPath) ──┤
stepCleanup  ──┘        │ InspectRepoDir          └─ Dirty/NotRepo/Unknown ─→ refuse + abort
                        └─ runGitStatus (var seam, overridable in tests)
```

## File Changes

| File | Action | Description |
|---|---|---|
| `internal/system/report.go` | Create | `Notify` var seam + `report(logFunc, msg)` |
| `internal/system/dryrun.go` | Create | `IsDryRun()` and the synthetic-result helper |
| `internal/system/repostate.go` | Create | `RepoState`, `InspectRepoDir`, `runGitStatus` var seam |
| `internal/system/exec.go` | Modify | gate in `Run`/`RunWithLogs`; `IsRegular()` skip in `CopyDir` walk; dry-run guards in `CopyFile`, `CopyDir`, `CreateBackup` |
| `internal/tui/installer.go` | Modify | Arch lists 720/775/831; `removeRepoDirIfSafe` helper; both delete sites; gated post-clone verify; wire `system.Notify` |
| `internal/tui/platform_packages_test.go` | Modify | expected pacman string (111-131) |
| `internal/system/*_test.go` | Create/Modify | new dry-run, skip-file and repo-state tests |
| `internal/system/detect.go` | Unchanged | Arch detection already covers CachyOS (detect.go:93-96) |

## Work Units, Ordering and Budget

| # | Unit | Ordering rationale | Est. lines |
|---|---|---|---|
| 1 | Arch package lists | Zero shared code, no dependency, unblocks installs immediately; ship first so a broken later unit never blocks the actual fix | 15–25 |
| 2 | `CopyDir` skips non-regular files + `Notify` seam | Introduces the reporter seam unit 4 reuses, so it must precede unit 4 | 45–70 |
| 3 | Fail-closed delete guard | Collapses two delete sites into one helper, so unit 4 adds a dry-run branch in **one** place instead of two | 130–170 |
| 4 | Real `--dry-run` | Largest blast radius (shared `exec.go`); last so a revert leaves 1–3 intact and useful | 150–220 |

**Decision needed before apply: Yes** — combined 340–485 lines. **Chained PRs recommended: Yes.** **400-line budget risk: High.** Feature Branch Chain: PR#1 = unit 1, PR#2 = units 2+3 targeting PR#1's branch, PR#3 = unit 4 targeting PR#2's branch. Unit 4 is never bundled.

## Cross-Platform Blast Radius

| Platform | Units touching it | Guard |
|---|---|---|
| Arch / CachyOS | 1,2,3,4 | primary target |
| macOS | 2,4 (`exec.go`) | `IsDryRun()` defaults false → Homebrew path byte-identical; full golden suite must pass on macos-latest |
| Debian / Ubuntu, Fedora, Alpine | 2,4 (`exec.go`) | package lists untouched; Docker E2E matrix |
| Termux | 2,3,4 (`exec.go` + clone/cleanup) | `Run`'s Termux direct-exec branch sits *after* the gate, so dry-run never reaches `parseCommand`; E2E termux image |

## Testing Strategy (TDD — RED first, `go test ./... -skip Golden` from `installer/`)

| Layer | What | Approach |
|---|---|---|
| Unit | `IsDryRun` true → `Run`/`RunWithLogs` execute nothing and return `Error == nil` | `t.Setenv`, override `Notify`, assert emitted plan text |
| Unit | dry-run `CreateBackup` returns a non-empty dir and creates nothing on disk | assert `os.Stat(dir)` is `IsNotExist` and return value non-empty |
| Unit | `CopyDir` over a dir containing `net.Listen("unix", …)` completes and reports the skip | assert regular files copied, socket path in captured `Notify` messages |
| Unit | `InspectRepoDir` → absent / clean / dirty / not-a-repo / git-error | override `runGitStatus`; real `git init` fixture in `t.TempDir()` for the clean/dirty pair |
| Unit | relative path is absolutized before git/delete | assert `runGitStatus` received an absolute path |
| Integration | `stepCloneRepo` refuses on dirty checkout; `stepCleanup` warns and returns nil | existing TUI step tests + fixture cwd |
| Integration | Arch pacman string no longer contains `carapace` / `zsh-theme-powerlevel10k` | `platform_packages_test.go` var-seam capture |
| E2E | Docker matrix ubuntu/debian/fedora/alpine/termux green | `e2e/docker-test.sh` |

Note: `-skip Golden` is mandatory on Linux; the three macOS-recorded snapshots are known baseline failures (`openspec/config.yaml` `known_baseline_failures`) and are never regressions.

## Threat Matrix

| Boundary | Applicability | Design response | Planned RED tests |
|---|---|---|---|
| Documentation-like paths | N/A — no file-classification or executable-content boundary; `CopyDir` classifies by `Mode().IsRegular()` only, and never executes what it copies | — | — |
| Git repository selection (`git -C`, relative vs absolute) | **Applicable** — the defect was `rm -rf` on a relative path; the guard now runs `git -C` | `filepath.Abs` before both `runGitStatus` and any delete; no shell interpolation of the path into a command string | abs-path assertion; path containing spaces; symlinked cwd |
| Commit state (staged, untracked, empty index) | **Applicable** — `git status --porcelain` decides deletion | any non-empty porcelain output = dirty (staged, unstaged and untracked all block); empty = clean | dirty-unstaged, dirty-staged, untracked-only, clean |
| Push state | N/A — this change performs no push | — | — |
| PR commands | N/A — no PR automation in the installer | — | — |
| Subprocess composition (`rm -rf`, `pacman -S`) | **Applicable** — commands are built as shell strings and run via `GetShell() -c` | delete stays a fixed literal target; the dry-run gate intercepts before `exec.CommandContext`, so a dry run spawns no process at all | assert zero process spawn under dry-run for `Run`, `RunWithLogs` and their sudo wrappers |

## Migration / Rollout

No migration, no persisted state, no schema. Rollback is a per-slice `git revert`. Reverting unit 3 alone reintroduces the data-loss path — revert it only together with unit 4 or after a full backup.

## Open Questions

- [ ] `carapace` remains in the `Fedora` lists; not verified against Fedora repos in this change. Tracked as a follow-up, not a blocker.
- [ ] Whether `system.Notify` should also feed the TUI final summary (aggregated skipped-path list) or only the per-step log. Spec's visibility requirement decides; per-step log is the minimum.
