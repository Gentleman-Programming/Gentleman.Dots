{ lib, ... }:
{
  home.activation.copySkhd = lib.hm.dag.entryAfter ["writeBoundary"] ''
    echo "Copying Skhd configuration..."
    rm -rf "$HOME/.config/skhd"

    if [ ! -d "$HOME/.config/skhd" ]; then
      mkdir -p "$HOME/.config/skhd"
    fi

    cp -r ${toString ./skhd}/* "$HOME/.config/skhd/"
    chmod -R u+w "$HOME/.config/skhd"
  '';
}
