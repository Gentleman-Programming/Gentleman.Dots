{ lib, ... }:
{
  home.activation.copySimpleBar = lib.hm.dag.entryAfter ["writeBoundary"] ''
    echo "Installing simple-bar for Übersicht..."
    
    WIDGETS_DIR="$HOME/Library/Application Support/Übersicht/widgets"
    SIMPLE_BAR_DIR="$WIDGETS_DIR/simple-bar"
    
    # Create widgets directory if it doesn't exist
    if [ ! -d "$WIDGETS_DIR" ]; then
      mkdir -p "$WIDGETS_DIR"
    fi
    
    # Clone or update simple-bar
    if [ ! -d "$SIMPLE_BAR_DIR" ]; then
      echo "Cloning simple-bar..."
      git clone --depth 1 https://github.com/Jean-Tinland/simple-bar "$SIMPLE_BAR_DIR"
    else
      echo "simple-bar already installed, updating..."
      cd "$SIMPLE_BAR_DIR" && git pull --ff-only || true
    fi
    
    # Copy settings if they exist
    if [ -f "${toString ./simple-bar}/settings.json" ]; then
      cp "${toString ./simple-bar}/settings.json" "$HOME/Library/Application Support/Übersicht/simple-bar-settings.json"
      chmod u+w "$HOME/Library/Application Support/Übersicht/simple-bar-settings.json"
    fi
  '';
}
