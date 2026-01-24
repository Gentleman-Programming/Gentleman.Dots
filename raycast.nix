{ lib, ... }:
{
  home.activation.copyRaycast = lib.hm.dag.entryAfter ["writeBoundary"] ''
    echo "Copying Raycast scripts..."
    rm -rf "$HOME/Raycast Scripts"

    if [ ! -d "$HOME/Raycast Scripts" ]; then
      mkdir -p "$HOME/Raycast Scripts"
    fi

    cp -r ${toString ./raycast/scripts}/* "$HOME/Raycast Scripts/"
    chmod -R u+wx "$HOME/Raycast Scripts"
  '';
}
