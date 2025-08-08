{ lib, ... }:
{
  home.activation.copyZed = lib.hm.dag.entryAfter ["writeBoundary"] ''
    echo "Copying Zed configuration..."

    # Create temporary backup directory if existing configuration is present
    if [ -d "$HOME/.config/zed" ]; then
      echo "Backing up existing Zed configuration..."
      chmod -R u+w "$HOME/.config/zed" 2>/dev/null || true
      mv "$HOME/.config/zed" "$HOME/.config/zed.backup.$(date +%s)" 2>/dev/null || {
        echo "Cannot move existing config, forcing removal..."
        chmod -R u+w "$HOME/.config/zed" 2>/dev/null || true
        rm -rf "$HOME/.config/zed" 2>/dev/null || true
      }
    fi

    # Create fresh directory
    mkdir -p "$HOME/.config/zed"

    # Copy new configuration
    cp -r ${toString ./zed}/* "$HOME/.config/zed/"

    # Ensure write permissions
    chmod -R u+w "$HOME/.config/zed"

    echo "Zed configuration copied successfully"
  '';
}
