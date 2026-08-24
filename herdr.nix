{ lib, ... }:

{
  # Herdr — agent multiplexer that lives in your terminal (https://herdr.dev)
  # Profile sources are immutable Home Manager files. The active config stays a
  # regular writable file because it is selected at runtime.
  #
  # There is no config validation gate: herdr 0.7.3 exposes no CLI subcommand
  # that parses config.toml, so a candidate can only be assembled and written
  # atomically. Bad edits surface in the repo sources, and `server reload-config`
  # reports anything the running server rejects.
  home.file = {
    ".config/herdr/config-base.toml".source = ./herdr/config-base.toml;
    ".config/herdr/profiles/gentleman.toml".source = ./herdr/profiles/gentleman.toml;
    ".config/herdr/profiles/gentleman-cute.toml".source = ./herdr/profiles/gentleman-cute.toml;

    ".local/bin/herdr-theme" = {
      executable = true;
      text = ''
        #!/usr/bin/env bash
        set -u

        usage() {
          printf 'Usage: herdr-theme <gentleman|gentleman-cute>\n' >&2
        }

        if [ "$#" -ne 1 ]; then
          usage
          exit 2
        fi

        theme="$1"
        case "$theme" in
          gentleman|gentleman-cute) ;;
          *)
            printf 'Unknown theme: %s\n' "$theme" >&2
            usage
            exit 2
            ;;
        esac

        config_dir="$HOME/.config/herdr"
        config_path="$config_dir/config.toml"
        base="$config_dir/config-base.toml"
        profile="$config_dir/profiles/$theme.toml"

        if [ ! -r "$base" ] || [ ! -r "$profile" ]; then
          printf 'Herdr theme sources are unavailable. Re-run Home Manager activation.\n' >&2
          exit 1
        fi

        mkdir -p "$config_dir"
        if ! candidate="$(mktemp "$config_dir/.config.toml.XXXXXX")"; then
          printf 'Unable to create a Herdr theme candidate.\n' >&2
          exit 1
        fi
        trap 'rm -f "$candidate"' EXIT HUP INT TERM

        if ! cat "$base" "$profile" > "$candidate"; then
          printf 'Unable to assemble Herdr theme candidate.\n' >&2
          exit 1
        fi

        if ! chmod u+w "$candidate"; then
          printf 'Unable to make the Herdr theme candidate writable.\n' >&2
          exit 1
        fi

        if ! mv -f "$candidate" "$config_path"; then
          printf 'Unable to activate the Herdr theme candidate.\n' >&2
          exit 1
        fi
        trap - EXIT HUP INT TERM

        printf 'Herdr theme selected: %s\n' "$theme"

        if ! command -v herdr >/dev/null 2>&1; then
          printf 'Herdr is not installed; the selection applies once it is.\n' >&2
          exit 0
        fi

        if ! herdr server reload-config; then
          printf 'Warning: Herdr could not reload its server; the selection remains active for the next start.\n' >&2
        fi
      '';
    };
  };

  # Auto-install Herdr on Home Manager activation if it is missing.
  # Guarded so a missing/failed brew never breaks the activation (same approach as engram.nix).
  home.activation.installHerdr = lib.hm.dag.entryAfter [ "linkGeneration" ] ''
    echo "🔧 Setting up Herdr..."

    if command -v herdr >/dev/null 2>&1; then
      echo "✅ Herdr already installed"
    elif command -v brew >/dev/null 2>&1; then
      echo "🚀 Installing Herdr via Homebrew..."
      brew install herdr || echo "❌ Herdr installation failed (run 'brew install herdr' manually)"
    else
      echo "⚠️  Homebrew not found — install Herdr manually: brew install herdr"
    fi

    config_dir="$HOME/.config/herdr"
    config_path="$config_dir/config.toml"
    mkdir -p "$config_dir"

    # Never reset a runtime selection. The first successful installation uses
    # Gentleman-Cute only when no active config exists yet.
    if [ -e "$config_path" ] || [ -L "$config_path" ]; then
      echo "📝 Preserving existing Herdr config"
    elif ! candidate="$(mktemp "$config_dir/.config.toml.XXXXXX")"; then
      echo "⚠️  Unable to create an initial Herdr config candidate"
    else
      trap 'rm -f "$candidate"' EXIT HUP INT TERM

      if cat "${./herdr/config-base.toml}" "${./herdr/profiles/gentleman-cute.toml}" > "$candidate" \
        && chmod u+w "$candidate" \
        && mv -f "$candidate" "$config_path"; then
        echo "📝 Installed Gentleman-Cute as the initial Herdr config"
      else
        rm -f "$candidate"
        echo "⚠️  Unable to write the initial Herdr config"
      fi

      trap - EXIT HUP INT TERM
    fi
  '';
}
