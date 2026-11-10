# OmniWM migration (Home Manager)

The macOS flake now uses OmniWM instead of the imported yabai, skhd, and SketchyBar modules. Review the changes and run your normal Home Manager switch only when you choose to apply them; this repository change does not itself switch your machine.

## Start at login

1. Open your existing OmniWM app and enable **Start at Login** in its settings.
2. Verify it in macOS **System Settings → General → Login Items & Extensions**. This is a native macOS registration, not an OmniWM TOML setting; Home Manager does not create another launch agent or enable it for you.
3. If OmniWM requests Accessibility access, grant it yourself in **System Settings → Privacy & Security → Accessibility**. No security protections or permissions are changed automatically.

## What activation does

- Requires Apple Silicon and macOS 26 or later. The platform check runs before Home Manager's write boundary. The existing Intel flake target cannot activate this module unchanged.
- Preserves an existing `~/Applications/OmniWM.app` or `/Applications/OmniWM.app`. If neither exists, it uses `brew install --cask omniwm`; a missing Homebrew, failed install, or missing app afterward stops before changing legacy launchd services.
- Leaves `~/.config/omniwm/settings.toml` GUI-managed and writable. No settings, nativeSpaces, backups, or host-specific window and display IDs are imported, copied, or replaced. When no settings exist, OmniWM may create its own on first launch.
- Once an app is present, disables these current-user launchd labels and unloads them if running: `com.asmvik.yabai`, `com.koekeishiya.yabai`, `com.koekeishiya.skhd`, `homebrew.mxcl.sketchybar`, and `com.felixkratz.sketchybar`. Missing services are tolerated; other labels are not targeted. If a disable, inspection, or unload fails, activation reports the label and stops; earlier labels may already have changed.
- Does not uninstall the old packages, erase their configs, or start the legacy services. Old modules and Raycast scripts remain as inactive recovery material; invoking a legacy restart script or re-importing a legacy module can start a service again. Home Manager may unlink files that those old modules previously managed, such as their installer helpers.

The application check and migration run after Home Manager's link step. A brew failure leaves OmniWM settings and legacy service state unchanged, but is not a rollback of unrelated Home Manager work already performed. A missing GUI launchd domain is reported rather than silently claiming the migration succeeded.

## Test without touching services

Run `python3 scripts/omniwm/test_omniwm.py` from the repository root. The tests use a temporary HOME, fake app paths, and mock `brew`, `launchctl`, and platform commands. One wiring test additionally uses `nix-instantiate` when available; no real install or service operation is performed.
