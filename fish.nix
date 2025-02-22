{ pkgs, lib, ... }: 
{
  home.activation.copyFish = lib.hm.dag.entryAfter ["writeBoundary"] ''
    echo "Copying fish configuration..."
    rm -rf "$HOME/.config/fish"

    cp -r ${toString ./fish} "$HOME/.config/fish"
    chmod -R u+w "$HOME/.config/fish"
  '';
}
