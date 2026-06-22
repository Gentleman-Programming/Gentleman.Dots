{ pkgs, lib, ... }:

{
  # Manual reinstall script for SketchyBar (Homebrew-based)
  home.file."bin/install-sketchybar" = {
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
        echo "⚠️  Homebrew no encontrado, instalá SketchyBar manualmente con: brew install FelixKratz/formulae/sketchybar"
        exit 1
      fi

      echo "🚀 Installing SketchyBar via Homebrew..."

      # Ensure the tap is present (idempotent)
      "$BREW" tap FelixKratz/formulae >/dev/null 2>&1 || true

      if "$BREW" list --formula sketchybar >/dev/null 2>&1; then
        echo "✅ SketchyBar already installed"
      else
        "$BREW" install FelixKratz/formulae/sketchybar
        echo "✅ SketchyBar installed"
      fi

      echo "🔧 Starting SketchyBar service..."
      # SketchyBar is managed through Homebrew services; unlike yabai/skhd, it
      # does not provide native --start-service / --restart-service flags.
      if pgrep -x sketchybar >/dev/null 2>&1; then
        "$BREW" services restart sketchybar || {
          echo "⚠️  Failed to restart SketchyBar service. If Homebrew reports an untrusted tap, run:"
          echo "   $BREW trust --formula FelixKratz/formulae/sketchybar"
        }
      else
        "$BREW" services start sketchybar || {
          echo "⚠️  Failed to start SketchyBar service. If Homebrew reports an untrusted tap, run:"
          echo "   $BREW trust --formula FelixKratz/formulae/sketchybar"
        }
      fi

      echo ""
      echo "🎉 SketchyBar setup complete!"
      echo "⚠️  Remember to grant Accessibility permission to SketchyBar in System Settings if plugins need it."
    '';
    executable = true;
  };

  # Auto-install SketchyBar on home-manager activation (Homebrew-based)
  home.activation.copySketchybar = lib.hm.dag.entryAfter ["linkGeneration"] ''
    echo "🔧 Setting up SketchyBar..."

    # ─── Copy declarative configuration ───
    echo "Copying SketchyBar configuration..."
    SKETCHYBAR_DIR="$HOME/.config/sketchybar"
    rm -rf "$SKETCHYBAR_DIR"
    mkdir -p "$SKETCHYBAR_DIR"

    cp -pR ${toString ./sketchybar}/. "$SKETCHYBAR_DIR/" 2>/dev/null || true
    chmod -R u+w "$SKETCHYBAR_DIR"
    chmod +x "$SKETCHYBAR_DIR/sketchybarrc" 2>/dev/null || true
    /usr/bin/find "$SKETCHYBAR_DIR/plugins" -type f -name "*.sh" -exec chmod +x {} \; 2>/dev/null || true
    echo "⚙️ Copied SketchyBar config to $SKETCHYBAR_DIR"

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
      echo "⚠️  Homebrew no encontrado, instalá SketchyBar manualmente con: brew install FelixKratz/formulae/sketchybar"
    else
      # ─── Idempotent install ───
      "$BREW" tap FelixKratz/formulae >/dev/null 2>&1 || true

      if "$BREW" list --formula sketchybar >/dev/null 2>&1; then
        echo "✅ SketchyBar already installed"
      else
        echo "🚀 Installing SketchyBar via Homebrew..."
        "$BREW" install FelixKratz/formulae/sketchybar || echo "⚠️  SketchyBar install failed, run: brew install FelixKratz/formulae/sketchybar"
      fi

      # ─── Start/refresh the daemon ───
      # SketchyBar is managed through Homebrew services; unlike yabai/skhd, it
      # does not provide native --start-service / --restart-service flags.
      echo "🔧 Starting SketchyBar service..."
      if pgrep -x sketchybar >/dev/null 2>&1; then
        "$BREW" services restart sketchybar >/dev/null 2>&1 || {
          echo "⚠️  Failed to restart SketchyBar service. If Homebrew reports an untrusted tap, run:"
          echo "   $BREW trust --formula FelixKratz/formulae/sketchybar"
        }
      else
        "$BREW" services start sketchybar >/dev/null 2>&1 || {
          echo "⚠️  Failed to start SketchyBar service. If Homebrew reports an untrusted tap, run:"
          echo "   $BREW trust --formula FelixKratz/formulae/sketchybar"
        }
      fi
    fi

    echo "🎉 SketchyBar setup complete!"
  '';
}
