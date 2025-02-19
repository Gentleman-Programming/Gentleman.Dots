{ lib, ... }:
{
  home.activation.copyNvim = lib.hm.dag.entryAfter ["writeBoundary"] ''
    echo "Copiando configuración de nvim..."
    rm -rf "$HOME/.config/nvim"
    cp -r ${toString ./nvim} "$HOME/.config/nvim"
    chmod -R u+w "$HOME/.config/nvim"
  '';
}
