{ lib, pkgs, ... }:
{
  # SketchyBar is installed via home.packages in flake.nix
  # This module only handles configuration files
  
  home.activation.copySketchybar = lib.hm.dag.entryAfter ["writeBoundary"] ''
    echo "Copying SketchyBar configuration..."
    
    # Create config directory
    SKETCHYBAR_DIR="$HOME/.config/sketchybar"
    mkdir -p "$SKETCHYBAR_DIR/plugins"
    
    # Copy main config
    cp "${toString ./sketchybar}/sketchybarrc" "$SKETCHYBAR_DIR/sketchybarrc"
    chmod +x "$SKETCHYBAR_DIR/sketchybarrc"
    
    # Copy plugins
    cp "${toString ./sketchybar}/plugins/"*.sh "$SKETCHYBAR_DIR/plugins/"
    chmod +x "$SKETCHYBAR_DIR/plugins/"*.sh
    
    echo "SketchyBar configuration copied to $SKETCHYBAR_DIR"
  '';
}
