{ lib, ... }:
{
  home.activation.copyAlacritty = lib.hm.dag.entryAfter ["writeBoundary"] ''
    echo "Copying Alacritty configuration..."
    rm -rf "$HOME/.config/alacritty"

    # Check if the .config directory exists, if not, create it
    if [ ! -d "$HOME/.config/alacritty" ]; then
      mkdir -p "$HOME/.config/alacritty"
    fi

    cp -r ${toString ./alacritty}/* "$HOME/.config/alacritty/"
    chmod -R u+w "$HOME/.config/alacritty"
  '';
}
