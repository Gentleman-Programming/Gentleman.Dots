{ lib, ... }:
{
  home.activation.copyGhostty = lib.hm.dag.entryAfter ["writeBoundary"] ''
    echo "Copying Zed configuration..."
    rm -rf "$HOME/.config/zed"

    # Check if the .config directory exists, if not, create it
    if [ ! -d "$HOME/.config/zed" ]; then
      mkdir -p "$HOME/.config/zed"
    fi

    cp -r ${toString ./zed}/* "$HOME/.config/zed/"
    chmod -R u+w "$HOME/.config/zed"
  '';
}

