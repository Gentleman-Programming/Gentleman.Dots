{ lib, pkgs, ... }:
{
  # Install node for simple-bar-server
  home.packages = with pkgs; [
    nodejs
  ];

  home.activation.installSimpleBarComplete = lib.hm.dag.entryAfter ["writeBoundary"] ''
    echo "═══════════════════════════════════════════════════════════════"
    echo "           Installing simple-bar complete setup..."
    echo "═══════════════════════════════════════════════════════════════"
    
    # ─────────────────────────────────────────────────────────────────
    # 1. Install Übersicht via Homebrew if not installed
    # ─────────────────────────────────────────────────────────────────
    if [ ! -d "/Applications/Übersicht.app" ]; then
      echo "📦 Installing Übersicht..."
      brew install --cask ubersicht
    else
      echo "✓ Übersicht already installed"
    fi
    
    # ─────────────────────────────────────────────────────────────────
    # 2. Create Übersicht widgets directory
    # ─────────────────────────────────────────────────────────────────
    WIDGETS_DIR="$HOME/Library/Application Support/Übersicht/widgets"
    if [ ! -d "$WIDGETS_DIR" ]; then
      echo "📁 Creating Übersicht widgets directory..."
      mkdir -p "$WIDGETS_DIR"
    fi
    
    # ─────────────────────────────────────────────────────────────────
    # 3. Install simple-bar
    # ─────────────────────────────────────────────────────────────────
    SIMPLE_BAR_DIR="$WIDGETS_DIR/simple-bar"
    if [ ! -d "$SIMPLE_BAR_DIR" ]; then
      echo "📦 Cloning simple-bar..."
      git clone --depth 1 https://github.com/Jean-Tinland/simple-bar "$SIMPLE_BAR_DIR"
    else
      echo "✓ simple-bar already installed, updating..."
      cd "$SIMPLE_BAR_DIR" && git pull --ff-only 2>/dev/null || true
    fi
    
    # ─────────────────────────────────────────────────────────────────
    # 4. Install Gentleman theme
    # ─────────────────────────────────────────────────────────────────
    echo "🎨 Installing Gentleman theme..."
    cp "${toString ./simple-bar}/themes/gentleman.js" "$SIMPLE_BAR_DIR/lib/styles/themes/"
    
    # Register Gentleman theme in themes.js if not already registered
    THEMES_FILE="$SIMPLE_BAR_DIR/lib/styles/themes.js"
    if ! grep -q "import \* as Gentleman" "$THEMES_FILE"; then
      echo "   Registering Gentleman theme in themes.js..."
      
      # Add import line after RosePineDawn import
      sed -i "" 's|import \* as RosePineDawn from "./themes/rose-pine-dawn";|import * as RosePineDawn from "./themes/rose-pine-dawn";\
import * as Gentleman from "./themes/gentleman";|' "$THEMES_FILE"
      
      # Add to collection object
      sed -i "" 's|RosePineDawn: RosePineDawn.theme,|RosePineDawn: RosePineDawn.theme,\
  Gentleman: Gentleman.theme,|' "$THEMES_FILE"
    else
      echo "   ✓ Gentleman theme already registered"
    fi
    
    # ─────────────────────────────────────────────────────────────────
    # 5. Copy simple-bar settings
    # ─────────────────────────────────────────────────────────────────
    echo "⚙️  Copying simple-bar settings..."
    cp "${toString ./simple-bar}/settings.json" "$HOME/Library/Application Support/Übersicht/simple-bar-settings.json"
    chmod u+w "$HOME/Library/Application Support/Übersicht/simple-bar-settings.json"
    
    # ─────────────────────────────────────────────────────────────────
    # 6. Install simple-bar-server
    # ─────────────────────────────────────────────────────────────────
    SERVER_DIR="$HOME/.config/simple-bar-server"
    if [ ! -d "$SERVER_DIR" ]; then
      echo "📦 Cloning simple-bar-server..."
      git clone --depth 1 https://github.com/Jean-Tinland/simple-bar-server.git "$SERVER_DIR"
      echo "   Installing npm dependencies..."
      cd "$SERVER_DIR" && npm install --silent
    else
      echo "✓ simple-bar-server already installed"
    fi
    
    # ─────────────────────────────────────────────────────────────────
    # 7. Install and load launchd service
    # ─────────────────────────────────────────────────────────────────
    echo "🚀 Setting up simple-bar-server launchd service..."
    mkdir -p "$HOME/Library/LaunchAgents"
    cp "${toString ./simple-bar}/com.simple-bar-server.plist" "$HOME/Library/LaunchAgents/"
    
    # Unload if exists, then load fresh
    launchctl unload "$HOME/Library/LaunchAgents/com.simple-bar-server.plist" 2>/dev/null || true
    launchctl load -w "$HOME/Library/LaunchAgents/com.simple-bar-server.plist" 2>/dev/null || true
    
    # ─────────────────────────────────────────────────────────────────
    # 8. Configure macOS settings via defaults
    # ─────────────────────────────────────────────────────────────────
    echo "🍎 Configuring macOS settings..."
    
    # Auto-hide dock
    defaults write com.apple.dock autohide -bool true
    
    # Auto-hide menu bar (Always)
    defaults write NSGlobalDomain _HIHideMenuBar -bool false
    defaults write NSGlobalDomain AppleMenuBarVisibleInFullscreen -bool true
    
    # Restart Dock to apply changes
    killall Dock 2>/dev/null || true
    
    # ─────────────────────────────────────────────────────────────────
    # 9. Start Übersicht if not running
    # ─────────────────────────────────────────────────────────────────
    if ! pgrep -q "Übersicht"; then
      echo "🚀 Starting Übersicht..."
      open -a "Übersicht"
    else
      echo "✓ Übersicht already running, refreshing..."
      osascript -e 'tell application id "tracesOf.Uebersicht" to refresh' 2>/dev/null || true
    fi
    
    echo ""
    echo "═══════════════════════════════════════════════════════════════"
    echo "           ✅ simple-bar installation complete!"
    echo "═══════════════════════════════════════════════════════════════"
    echo ""
    echo "📋 Manual steps required:"
    echo "   1. System Settings → Menu Bar → 'Automatically hide' → Always"
    echo "   2. System Settings → Menu Bar → 'Show menu bar background' → ON"
    echo "   3. Add Übersicht to Login Items (System Settings → General → Login Items)"
    echo ""
  '';
}
