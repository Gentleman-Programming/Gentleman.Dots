{ lib, ... }:
{
  home.activation.copyGhostty = lib.hm.dag.entryAfter ["writeBoundary"] ''
    echo "Copying Ghostty configuration..."
    rm -rf "$HOME/.config/ghostty"

    # Check if the .config directory exists, if not, create it
    if [ ! -d "$HOME/.config/ghostty" ]; then
      mkdir -p "$HOME/.config/ghostty"
    fi

    cp -r ${toString ./ghostty}/* "$HOME/.config/ghostty/"
    chmod -R u+w "$HOME/.config/ghostty"
  '';
}

