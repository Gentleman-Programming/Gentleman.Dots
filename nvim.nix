{ lib, ... }:
{
  home.activation.copyNvim = lib.hm.dag.entryAfter ["writeBoundary"] ''
    echo "Copying nvim configuration..."
    rm -rf "$HOME/.config/nvim"

    # Ensure ~/.config exists before copying
    if [ ! -d "$HOME/.config" ]; then
      mkdir -p "$HOME/.config"
    fi

    cp -r ${toString ./nvim} "$HOME/.config/nvim"
    chmod -R u+w "$HOME/.config/nvim"
  '';
}
