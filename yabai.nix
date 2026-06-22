{ pkgs, lib, ... }:

{
  # Manual reinstall script for yabai (Homebrew-based)
  home.file."bin/install-yabai" = {
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
        echo "⚠️  Homebrew no encontrado, instalá yabai manualmente con: brew install asmvik/formulae/yabai"
        exit 1
      fi

      echo "🚀 Installing yabai via Homebrew..."

      # Ensure the tap is present (idempotent)
      "$BREW" tap asmvik/formulae >/dev/null 2>&1 || true

      if "$BREW" list yabai >/dev/null 2>&1 || command -v yabai >/dev/null 2>&1; then
        echo "✅ yabai already installed"
      else
        "$BREW" install asmvik/formulae/yabai
        echo "✅ yabai installed"
      fi

      echo "🔧 Starting yabai service..."
      # brew services is unsupported by this formula (no #service); use the
      # binary's native service manager. Restart if running, else start.
      if pgrep -x yabai >/dev/null 2>&1; then
        yabai --restart-service || true
      else
        yabai --start-service || true
      fi

      echo ""
      echo "🎉 yabai setup complete!"
      echo "⚠️  Remember to grant Accessibility permission to yabai in System Settings."
    '';
    executable = true;
  };

  # Auto-install yabai on home-manager activation (Homebrew-based)
  home.activation.copyYabai = lib.hm.dag.entryAfter ["linkGeneration"] ''
    echo "🔧 Setting up yabai..."

    # ─── macOS Mission Control settings ───
    # Disable automatic space reordering based on most recent use.
    # This is CRITICAL for yabai/sketchybar to work correctly with numbered spaces.
    echo "Configuring Mission Control settings..."
    /usr/bin/defaults write com.apple.dock mru-spaces -bool false
    killall Dock 2>/dev/null || true

    # ─── Copy declarative configuration ───
    echo "Copying Yabai configuration..."
    rm -rf "$HOME/.config/yabai"
    mkdir -p "$HOME/.config/yabai"

    cp -rf ${toString ./yabai}/* "$HOME/.config/yabai/" 2>/dev/null || true
    chmod -R u+w "$HOME/.config/yabai"
    chmod +x "$HOME/.config/yabai/yabairc" 2>/dev/null || true
    chmod +x "$HOME/.config/yabai"/*.sh 2>/dev/null || true
    echo "⚙️ Copied yabai config to $HOME/.config/yabai"

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
      echo "⚠️  Homebrew no encontrado, instalá yabai manualmente con: brew install asmvik/formulae/yabai"
    else
      # ─── Idempotent install ───
      "$BREW" tap asmvik/formulae >/dev/null 2>&1 || true

      if "$BREW" list yabai >/dev/null 2>&1 || command -v yabai >/dev/null 2>&1; then
        echo "✅ yabai already installed"
      else
        echo "🚀 Installing yabai via Homebrew..."
        "$BREW" install asmvik/formulae/yabai || echo "⚠️  yabai install failed, run: brew install asmvik/formulae/yabai"
      fi

      # ─── Start/refresh the daemon via the binary's native service manager ───
      # brew services is unsupported by this formula (no #service); use
      # yabai --start-service / --restart-service for idempotency.
      echo "🔧 Starting yabai service..."
      if pgrep -x yabai >/dev/null 2>&1; then
        yabai --restart-service >/dev/null 2>&1 || true
      else
        yabai --start-service >/dev/null 2>&1 || true
      fi

      # ─── Yabai sudoers setup for scripting addition ───
      # Required for space switching (yabai --load-sa needs passwordless sudo).
      # brew upgrade yabai changes the binary, which invalidates the sha256 in
      # the sudoers entry — so we recompute the hash on every activation and, if
      # it differs from what's installed, print the command for the user to run.
      YABAI_PATH="$(command -v yabai 2>/dev/null || true)"
      if [ -z "$YABAI_PATH" ] && [ -x /opt/homebrew/bin/yabai ]; then
        YABAI_PATH="/opt/homebrew/bin/yabai"
      elif [ -z "$YABAI_PATH" ] && [ -x /usr/local/bin/yabai ]; then
        YABAI_PATH="/usr/local/bin/yabai"
      fi

      if [ -n "$YABAI_PATH" ] && [ -x "$YABAI_PATH" ]; then
        YABAI_HASH=$(/usr/bin/shasum -a 256 "$YABAI_PATH" | /usr/bin/awk '{print $1}')
        EXPECTED_ENTRY="$USER ALL=(root) NOPASSWD: sha256:$YABAI_HASH $YABAI_PATH --load-sa"

        NEEDS_UPDATE=false
        if [ ! -f /private/etc/sudoers.d/yabai ]; then
          NEEDS_UPDATE=true
        elif ! /usr/bin/grep -q "$YABAI_HASH" /private/etc/sudoers.d/yabai 2>/dev/null; then
          NEEDS_UPDATE=true
        fi

        if [ "$NEEDS_UPDATE" = true ]; then
          echo ""
          echo "══════════════════════════════════════════════════════════════"
          echo "  ⚠️  Yabai scripting addition needs sudoers update!"
          echo "  Run this command to enable space switching:"
          echo ""
          echo "  echo \"$EXPECTED_ENTRY\" | sudo tee /private/etc/sudoers.d/yabai"
          echo "  sudo $YABAI_PATH --load-sa"
          echo "══════════════════════════════════════════════════════════════"
          echo ""
        else
          echo "✅ Yabai sudoers entry is up to date."
        fi
      else
        echo "⚠️  yabai binary not found — skipping sudoers setup."
      fi
    fi

    echo "🎉 yabai setup complete!"
  '';
}
