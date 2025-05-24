{ lib, ... }:
{
  home.activation.copyZed = lib.hm.dag.entryAfter ["writeBoundary"] ''
    echo "Copying Zed configuration..."
    rm -rf "$HOME/.config/zed"
    mkdir -p "$HOME/.config/zed"
    cp -r ${toString ./zed}/* "$HOME/.config/zed/"
    chmod -R u+w "$HOME/.config/zed"
  '';
}
