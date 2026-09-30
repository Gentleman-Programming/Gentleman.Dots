# Proposal: Fix the Gentleman.Dots installer on Arch Linux / CachyOS

## Intent

On Arch-family systems the installer is unsafe and largely non-functional. OS detection is fine (`isArchLinux()` reads `/etc/arch-release`, which CachyOS ships), but four downstream defects were reproduced on v2.12.2:

1. `carapace` is in the `Arch:` package list for fish/zsh/nushell, but it is not in the official repos. pacman transactions are atomic, so the whole `pacman -S` aborts and zoxide/atuin/starship/neovim/zellij never install. The Homebrew fallback only runs when `hasBrew` is true, so an Arch box has no recovery path.
2. `--dry-run` sets `GENTLEMAN_DRY_RUN=1` and nothing reads it. A dry run really created a backup, really re-cloned 194MB and really ran `sudo pacman -S`.
3. `CopyDir` calls `CopyFile` (which uses `os.ReadFile`) on every non-directory, so one live Unix socket aborts the backup of every config.
4. `rm -rf Gentleman.Dots` runs on a **relative** path at two sites (`stepCloneRepo`, `stepCleanup`) with no git-status check and no confirmation. Running the installer from `~` destroyed uncommitted work.

Users on Arch/CachyOS currently get a failed install plus real data loss. That is the reason to act now.

## Scope

### In Scope
- Remove `carapace` from the three Arch package lists and update `platform_packages_test.go:111-131`.
- Skip non-regular files in `CopyDir` with a warning instead of aborting the backup.
- Guard both `rm -rf Gentleman.Dots` sites with a `git status --porcelain` dirty check plus confirmation.
- Make `--dry-run` real via a central `system.IsDryRun()` gate on the four `Run*` functions plus explicit guards in `CopyFile`, `CopyDir`, `CreateBackup` and the two delete sites.

### Out of Scope (explicit non-goals)
- **AUR / `carapace-bin` fast-follow**: detecting `paru`/`yay` and best-effort installing `carapace-bin` is a deliberate non-goal here. (`carapace-bin` does not `provides=carapace`, so it cannot rescue `pacman -S carapace`.) It is a separate later change.
- **Upstreaming to `Gentleman-Programming/Gentleman.Dots`**: not part of this change; it is a follow-up decision after local verification.
- No CachyOS-specific detection branch — detection already works.
- No full `dryRun` parameter threading through call signatures (disproportionate).
- No change to the three macOS-recorded golden snapshots (known baseline failures; run `go test ./... -skip Golden`).

## Capabilities

### New Capabilities
- `arch-package-install`: which packages the Arch/pacman path installs and how an unresolvable name must not abort the transaction.
- `config-backup`: backup MUST tolerate non-regular files (sockets, FIFOs, devices) and continue.
- `repo-clone-safety`: destructive removal of a `Gentleman.Dots` directory MUST be git-status-aware and confirmed.
- `dry-run-safety`: `--dry-run` MUST perform zero mutations of any kind.

### Modified Capabilities
- None (no existing specs in `openspec/specs/`).

## Approach

Four chained work units, in this order, each independently shippable and verifiable.

| # | Work unit | Boundary | Est. lines |
|---|-----------|----------|-----------|
| 1 | Unblock Arch installs | Drop `carapace` from the 3 `platformPackages.Arch` fields in `internal/tui/installer.go`; update the expected pacman string in `platform_packages_test.go` | 10–20 |
| 2 | Backup resilience | `info.Mode().IsRegular()` check in `CopyDir`'s Walk callback; skip + warn. Test uses `net.Listen("unix", ...)` | 15–25 |
| 3 | Destructive-delete guard | `git status --porcelain` guard at both call sites; handle dir-absent / not-a-git-repo / clean / dirty | 120–160 |
| 4 | Real `--dry-run` | `system.IsDryRun()` gate on the four `Run*` functions + explicit guards in `CopyFile`/`CopyDir`/`CreateBackup` and both delete sites | 150–250 |

**Delivery (auto-chain).** Combined 350–550 lines exceeds the 400-line review budget, so this ships as a Feature Branch Chain: PR #1 (unit 1) targets the tracker branch; each later PR targets the immediately previous PR branch. Units 1–3 are small enough to pair; unit 4 must be its own PR because it touches the most shared code.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `installer/internal/tui/installer.go` | Modified | Arch package lists (717-722, 771-777, 828-834); `stepCloneRepo` (~105-119); `stepCleanup` (~1195-1205) |
| `installer/internal/system/exec.go` | Modified | `CopyFile` (286-303), `CopyDir` (306-352), `CreateBackup` (442-493), the four `Run*` functions |
| `installer/cmd/gentleman-installer/main.go` | Modified | `GENTLEMAN_DRY_RUN` now actually consumed (70-72) |
| `installer/internal/tui/platform_packages_test.go` | Modified | Expected pacman string (111-131) |
| `installer/internal/system/exec_test.go` | Modified | New non-regular-file and dry-run coverage |
| `installer/internal/system/detect.go` | Unchanged | Detection already correct |

## Cross-Platform Impact

| Platform | Impact |
|----------|--------|
| Arch / CachyOS | Primary target — all four fixes |
| macOS | Units 2 and 4 touch shared `exec.go`; Homebrew path must stay byte-identical. Full golden suite must still pass on macOS |
| Debian / Ubuntu | Shared `exec.go` only; apt package lists untouched |
| Fedora | Shared `exec.go` only; dnf lists untouched |
| Alpine | Shared `exec.go` only; apk lists untouched |
| Termux | Shared `exec.go` and clone/cleanup path; pkg lists untouched |

Docker E2E matrix (ubuntu/debian/fedora/alpine/termux) is the regression net for units 2–4.

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Unit 4 regresses non-Arch platforms via shared `Run*` | Med | Default `IsDryRun()` false; dedicated tests per `Run*`; full Docker E2E matrix before merge |
| Unit 2 changes error semantics fatal → skip+warn, hiding real failures | Med | Warn loudly and surface skipped paths in the TUI summary; spec MUST state the visibility requirement |
| Unit 3 adds new failure modes (git absent, non-git dir) | Med | Explicit spec scenarios for all four states; fail closed (refuse to delete) when git status cannot be determined |
| Removing `carapace` silently drops a completion feature on Arch | Low | Note it in the install summary; AUR fast-follow tracked separately |
| Chained PRs show earlier slices in child diffs | Low | Rebase/retarget until each child diff is clean |

## Rollback Plan

Each unit is one commit on its own branch, so rollback is per-slice `git revert`.

- Unit 1: reverting restores the broken package list — safe, no state.
- Unit 2: reverting restores the fatal abort — safe, no state.
- Unit 3: reverting restores unguarded deletion — safe but reintroduces data loss; revert only with unit 4 also present or a full backup taken.
- Unit 4: highest blast radius; if any platform regresses, revert this slice alone — units 1–3 remain valid and independently useful.
- No migrations, no persisted state, no schema. Nothing to undo beyond source.

## Dependencies

- `git` present on PATH for unit 3 (must be handled as an explicit failure mode, not assumed).
- Verification runs `go test ./... -skip Golden` from `installer/`; `-skip Golden` is mandatory on Linux.

## Success Criteria

- [ ] `pacman -S` completes on CachyOS and zoxide, atuin, starship, neovim and zellij are installed.
- [ ] A backup of `~/.config` containing a live Unix socket completes, skipping the socket with a visible warning.
- [ ] `--dry-run` produces zero filesystem mutations, zero clones and zero `sudo` invocations, verified on a dirty home.
- [ ] Running the installer from a directory containing a dirty `Gentleman.Dots` checkout refuses to delete it without explicit confirmation.
- [ ] `go test ./... -skip Golden` passes from `installer/`; `go build ./cmd/gentleman-installer` succeeds; Docker E2E matrix green.
