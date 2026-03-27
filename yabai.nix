{ lib, ... }:
{
  home.activation.copyYabai = lib.hm.dag.entryAfter ["writeBoundary"] ''
    # ─── macOS Mission Control settings ───
    # Disable automatic space reordering based on most recent use
    # This is CRITICAL for yabai/sketchybar to work correctly with numbered spaces
    echo "Configuring Mission Control settings..."
    /usr/bin/defaults write com.apple.dock mru-spaces -bool false
    killall Dock 2>/dev/null || true

    echo "Copying Yabai configuration..."
    rm -rf "$HOME/.config/yabai"

    if [ ! -d "$HOME/.config/yabai" ]; then
      mkdir -p "$HOME/.config/yabai"
    fi

    cp -r ${toString ./yabai}/* "$HOME/.config/yabai/"
    chmod -R u+w "$HOME/.config/yabai"
    chmod +x "$HOME/.config/yabai/yabairc"
    chmod +x "$HOME/.config/yabai/move-window-to-space.sh"

    # ─── Yabai sudoers setup for scripting addition ───
    # Required for space switching (yabai --load-sa needs passwordless sudo)
    YABAI_BIN="$(which yabai)"
    YABAI_HASH="$(shasum -a 256 "$YABAI_BIN" | awk '{print $1}')"
    EXPECTED_ENTRY="$USER ALL=(root) NOPASSWD: sha256:$YABAI_HASH $YABAI_BIN --load-sa"

    NEEDS_UPDATE=false
    if [ ! -f /private/etc/sudoers.d/yabai ]; then
      NEEDS_UPDATE=true
    elif ! grep -q "$YABAI_HASH" /private/etc/sudoers.d/yabai 2>/dev/null; then
      NEEDS_UPDATE=true
    fi

    if [ "$NEEDS_UPDATE" = true ]; then
      echo ""
      echo "══════════════════════════════════════════════════════════════"
      echo "  ⚠️  Yabai scripting addition needs sudoers update!"
      echo "  Run this command to enable space switching:"
      echo ""
      echo "  echo \"$EXPECTED_ENTRY\" | sudo tee /private/etc/sudoers.d/yabai"
      echo "  sudo yabai --load-sa"
      echo "══════════════════════════════════════════════════════════════"
      echo ""
    else
      echo "Yabai sudoers entry is up to date."
    fi
  '';
}
