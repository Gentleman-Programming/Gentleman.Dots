{ pkgs, lib, ... }:

{
  # Manual reinstall script for skhd (Homebrew-based)
  home.file."bin/install-skhd" = {
    text = ''
      #!/usr/bin/env bash
      set -e

      # ─── Locate Homebrew ───
      # brew is not guaranteed to be on PATH; probe the standard locations.
      BREW=""
      if [ -x /opt/homebrew/bin/brew ]; then
        BREW="/opt/homebrew/bin/brew"
      elif [ -x /usr/local/bin/brew ]; then
        BREW="/usr/local/bin/brew"
      fi

      if [ -z "$BREW" ]; then
        echo "⚠️  Homebrew no encontrado, instalá skhd manualmente con: brew install koekeishiya/formulae/skhd"
        exit 1
      fi

      echo "🚀 Installing skhd via Homebrew..."

      # Ensure the tap is present (idempotent)
      "$BREW" tap koekeishiya/formulae >/dev/null 2>&1 || true

      if "$BREW" list skhd >/dev/null 2>&1 || command -v skhd >/dev/null 2>&1; then
        echo "✅ skhd already installed"
      else
        "$BREW" install koekeishiya/formulae/skhd
        echo "✅ skhd installed"
      fi

      echo "🔧 Starting skhd service..."
      # brew services is unsupported by this formula (no #service); use the
      # binary's native service manager. Restart if running, else start.
      if pgrep -x skhd >/dev/null 2>&1; then
        skhd --restart-service || true
      else
        skhd --start-service || true
      fi

      echo ""
      echo "🎉 skhd setup complete!"
      echo "⚠️  Remember to grant Accessibility permission to skhd in System Settings."
    '';
    executable = true;
  };

  # Auto-install skhd on home-manager activation (Homebrew-based)
  home.activation.copySkhd = lib.hm.dag.entryAfter ["linkGeneration"] ''
    echo "🔧 Setting up skhd..."

    # ─── Copy declarative configuration ───
    echo "Copying Skhd configuration..."
    rm -rf "$HOME/.config/skhd"
    mkdir -p "$HOME/.config/skhd"

    cp -rf ${toString ./skhd}/* "$HOME/.config/skhd/" 2>/dev/null || true
    chmod -R u+w "$HOME/.config/skhd"
    echo "⚙️ Copied skhd config to $HOME/.config/skhd"

    # ─── Remove stale nix-era binary ───
    # Older nix generations installed skhd into ~/.local/bin, which takes PATH
    # priority over Homebrew. nix no longer manages it, so it lingers as an
    # orphan that both shadows the brew binary at runtime and makes the
    # `command -v skhd` check below skip the brew install. Drop it first.
    if [ -e "$HOME/.local/bin/skhd" ]; then
      rm -f "$HOME/.local/bin/skhd"
      echo "🧹 Removed stale ~/.local/bin/skhd (now provided by Homebrew)"
    fi

    # ─── Locate Homebrew ───
    # brew is not guaranteed to be on PATH during home-manager activation;
    # probe the standard locations (Apple Silicon first, Intel fallback).
    BREW=""
    if [ -x /opt/homebrew/bin/brew ]; then
      BREW="/opt/homebrew/bin/brew"
    elif [ -x /usr/local/bin/brew ]; then
      BREW="/usr/local/bin/brew"
    fi

    if [ -z "$BREW" ]; then
      echo "⚠️  Homebrew no encontrado, instalá skhd manualmente con: brew install koekeishiya/formulae/skhd"
    else
      # ─── Idempotent install ───
      "$BREW" tap koekeishiya/formulae >/dev/null 2>&1 || true

      if "$BREW" list skhd >/dev/null 2>&1 || command -v skhd >/dev/null 2>&1; then
        echo "✅ skhd already installed"
      else
        echo "🚀 Installing skhd via Homebrew..."
        "$BREW" install koekeishiya/formulae/skhd || echo "⚠️  skhd install failed, run: brew install koekeishiya/formulae/skhd"
      fi

      # ─── Start/refresh the daemon via the binary's native service manager ───
      # brew services is unsupported by this formula (no #service); use
      # skhd --start-service / --restart-service for idempotency.
      echo "🔧 Starting skhd service..."
      if pgrep -x skhd >/dev/null 2>&1; then
        skhd --restart-service >/dev/null 2>&1 || true
      else
        skhd --start-service >/dev/null 2>&1 || true
      fi
    fi

    echo "🎉 skhd setup complete!"
  '';
}
