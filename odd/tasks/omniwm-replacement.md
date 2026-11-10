# OmniWM replacement

Objective: replace the legacy macOS window-management startup imports with a safe OmniWM Home Manager integration.
Branch: `feat/omniwm-replacement` · Runner: mocked shell tests and offline Nix parsing/evaluation where available.

## Specs

- **S1 — Requested outcome:** “ok otra cosa que te pido sacame yabai y skhb y sketchybar del arranque de la mac, y cambia mi Gentleman.Dots para agregar esta nueva app como reemplazo a las otras”. The human confirmed the target repository with “si” to “¿Confirmas que trabajemos en Gentleman.Dots2?”.
- **S2 — Replacement and preservation:** Import OmniWM instead of the yabai, skhd and sketchybar modules; keep their files as inactive recovery material. Do not touch Herdr, existing OmniWM GUI-managed settings, nativeSpaces, backups, old packages, SIP, or the user's other services. Keep existing non-target imports unchanged.
- **S3 — Safe install and startup:** Preserve an existing OmniWM app in `~/Applications` or `/Applications`; otherwise install the official Homebrew cask only on Apple Silicon and macOS 26+. Reject unsupported hosts before side effects. Do not create a duplicate launch agent, bypass Gatekeeper, grant permissions, or claim native Start at Login is enabled.
- **S4 — Legacy services:** If activation handles migration, disable/unload only the five specified legacy current-user launchd labels, idempotently and safely when absent; never start the old services. Prevent install failures from changing configuration or services.
- **S5 — Evidence and delivery:** Mock commands and fake HOME for deterministic behavior tests, including imports, exact output/status/effects, existing configuration preservation, unsupported platforms, and failed installation. Parse both Nix files and attempt offline evaluation when cached without fetching or changing the lock. No live activation, install, download, service changes, commit, or remote operations. Independent verification and commit remain parent-owned.

## Tasks

- [x] **T1** — Replace the imports, implement the guarded integration, document user-managed startup and migration, and add focused mocked tests. The unit exceeds the approximately 400-line forecast; tests and documentation remain with the behavior.
- [x] **T2** — Independent verification passed for S1–S5; parent committed T1/T2 behavior and proof in `3cb3025`. This documentation-only follow-up awaits a separate parent-owned commit.

## Log

- **L1:** Original request: “ok otra cosa que te pido sacame yabai y skhb y sketchybar del arranque de la mac, y cambia mi Gentleman.Dots para agregar esta nueva app como reemplazo a las otras”. Repository correction and confirmation: “¿Confirmas que trabajemos en Gentleman.Dots2?” → “si”.
- **L2:** Existing `fish.nix` and `kitty.nix` are unrelated dirty files and must remain untouched. Existing OmniWM settings are GUI-managed; native Start at Login is a separate manual action. All five legacy labels are already disabled/unloaded on the current host according to the parent.
- **L3:** Engram mirror saved as observation `26845`, but observation readback is unavailable to this worker. Mirror/readback remains pending; the parent must verify or mirror this document before source work proceeds.
- **L4:** Blocked before source changes because the required observation readback cannot be performed. `nix-instantiate --parse flake.nix` succeeded; parsing `omniwm.nix` failed because it has not yet been created. T1 and T2 remain pending.
- **L5:** Parent confirmed the full observation `26845` matches this document; pre-source readback is fulfilled. T1 is in progress and the parent owns future mirror readback.
- **L6:** T1 self-check: initial behavior RED had 13 failures; the final mocked runner passes 13 tests. Both Nix files parse, and an isolated offline evaluation using cached Home Manager/Nixpkgs DAG libraries yields `checkOmniWM`, `writeBoundary`, `linkGeneration`, `setupOmniWM` in order. No live activation, install, or service operation was performed. Full flake evaluation and real macOS launchctl behavior remain unverified; independent verification is required.
- **L7:** Parent reports an earlier verifier attempt was denied before any runtime probes. A separately authorized bounded verifier then passed S1–S5 with no product blocker: both Nix parses and the cached DAG order passed, and read-only `launchctl` inspection returned 113 for all five legacy labels. This new evidence supersedes L6's launchctl uncertainty. Full flake evaluation was skipped because the verifier was permitted only four files; no Home Manager switch occurred.
- **L8:** The independently tested regression file was copied byte-for-byte. `python3 -B scripts/omniwm/test_verifier_regressions.py` passed 2 tests and `python3 -B scripts/omniwm/test_omniwm.py` passed 13 tests. Parent also reports its 13-test rerun and `git diff --check` passed. Native OmniWM Start at Login remains a manual action and is not enabled by this module.
- **L9:** The verifier reported a transient untracked `.claude` path of unknown origin; the parent reports its subsequent status did not show it. It was not inspected or deleted and is excluded from staging. Existing dirty `fish.nix` and `kitty.nix` remain untouched.
- **L10:** Parent committed T1/T2 behavior and proof as `3cb3025` (`feat(macos): replace legacy window managers with OmniWM`; 9 files, 584 additions, 71 deletions) and reports 13 original tests, 2 independent tests, and diff check passed. No push occurred. This documentation-only checkpoint is pending a separate parent commit. Full flake evaluation remains skipped; no Home Manager switch occurred, and native Start at Login remains manual. Unrelated `fish.nix` and `kitty.nix` remain untouched.
