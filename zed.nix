{ lib, ... }:
{
  home.activation.copyZed = lib.hm.dag.entryAfter ["writeBoundary"] ''
    echo "Copying Zed configuration..."

    # Crear directorio de respaldo temporal si existe configuración previa
    if [ -d "$HOME/.config/zed" ]; then
      echo "Backing up existing Zed configuration..."
      chmod -R u+w "$HOME/.config/zed" 2>/dev/null || true
      mv "$HOME/.config/zed" "$HOME/.config/zed.backup.$(date +%s)" 2>/dev/null || {
        echo "Cannot move existing config, forcing removal..."
        chmod -R u+w "$HOME/.config/zed" 2>/dev/null || true
        rm -rf "$HOME/.config/zed" 2>/dev/null || true
      }
    fi

    # Crear directorio limpio
    mkdir -p "$HOME/.config/zed"

    # Copiar nueva configuración
    cp -r ${toString ./zed}/* "$HOME/.config/zed/"

    # Asegurar permisos de escritura
    chmod -R u+w "$HOME/.config/zed"

    echo "Zed configuration copied successfully"
  '';
}
