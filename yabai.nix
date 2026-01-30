{ lib, ... }:
{
  home.activation.copyYabai = lib.hm.dag.entryAfter ["writeBoundary"] ''
    # ─── macOS Mission Control settings ───
    # Disable automatic space reordering based on most recent use
    # This is CRITICAL for yabai/sketchybar to work correctly with numbered spaces
    echo "Configuring Mission Control settings..."
    defaults write com.apple.dock mru-spaces -bool false
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
  '';
}
