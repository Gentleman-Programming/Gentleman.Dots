{ lib, ... }:
{
  home.activation.copyYabai = lib.hm.dag.entryAfter ["writeBoundary"] ''
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
