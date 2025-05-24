{ lib, ... }:
{
  home.activation.copyTelevision = lib.hm.dag.entryAfter ["writeBoundary"] ''
    echo "Copying Television configuration..."
    rm -rf "$HOME/.config/television"

    # Check if the .config directory exists, if not, create it
    if [ ! -d "$HOME/.config/television" ]; then
      mkdir -p "$HOME/.config/television"
    fi

    cp -r ${toString ./television}/* "$HOME/.config/television/"
    chmod -R u+w "$HOME/.config/television"
  '';
}
