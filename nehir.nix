{ lib, ... }:
{
  # Nehir (Niri-style scrolling-column window manager) configuration.
  # Config-only, same pattern as ghostty.nix: the binary is installed via
  # Homebrew (`brew tap guria/tap && brew install --cask nehir`), and Nehir
  # watches ~/.config/nehir for live changes (no reload/LaunchAgent needed).
  home.activation.copyNehir = lib.hm.dag.entryAfter ["writeBoundary"] ''
    echo "Copying Nehir configuration..."
    rm -rf "$HOME/.config/nehir"

    if [ ! -d "$HOME/.config/nehir" ]; then
      mkdir -p "$HOME/.config/nehir"
    fi

    cp -r ${toString ./nehir}/* "$HOME/.config/nehir/"
    chmod -R u+w "$HOME/.config/nehir"
  '';
}
