{ lib, ... }:
{
  home.activation.copyAerospace = lib.hm.dag.entryAfter ["writeBoundary"] ''
    echo "Copying Aerospace configuration..."

    cp -r ${toString ./aerospace/.aerospace.toml} "$HOME"
  '';
}

