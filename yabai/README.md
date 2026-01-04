# Yabai + skhd + simple-bar Configuration

Tiling window manager setup for macOS with custom keybindings and a beautiful status bar.

## Prerequisites

### Update Command Line Tools
If you encounter errors about Command Line Tools not supporting macOS 26:

```bash
# Remove old Command Line Tools
sudo rm -rf /Library/Developer/CommandLineTools

# Install new Command Line Tools
sudo xcode-select --install

# Set developer directory
sudo xcode-select --switch /Library/Developer/CommandLineTools
```

Wait for the installation dialog to complete before proceeding.

## Installation

### 1. Install yabai, skhd, and Übersicht

```bash
# Install via Homebrew
brew install koekeishiya/formulae/yabai
brew install koekeishiya/formulae/skhd
brew install --cask ubersicht
```

### 2. Create config directories

```bash
mkdir -p ~/.config/yabai
mkdir -p ~/.config/skhd
```

### 3. Copy config files

```bash
# From your dotfiles repo
cp yabai/yabairc ~/.config/yabai/yabairc
cp skhd/skhdrc ~/.config/skhd/skhdrc
cp yabai/move-window-to-space.sh ~/.config/yabai/move-window-to-space.sh

# Make helper script executable
chmod +x ~/.config/yabai/move-window-to-space.sh
chmod +x ~/.config/yabai/yabairc
```

### 4. Install simple-bar

```bash
# Open Übersicht once to create the widgets folder
open -a "Übersicht"

# Clone simple-bar
git clone --depth 1 https://github.com/Jean-Tinland/simple-bar "$HOME/Library/Application Support/Übersicht/widgets/simple-bar"

# Copy Gentleman theme
cp simple-bar/themes/gentleman.js "$HOME/Library/Application Support/Übersicht/widgets/simple-bar/lib/styles/themes/"

# Copy settings
cp simple-bar/settings.json "$HOME/Library/Application Support/Übersicht/simple-bar-settings.json"
```

### 5. Register Gentleman theme in simple-bar

Edit `$HOME/Library/Application Support/Übersicht/widgets/simple-bar/lib/styles/themes.js`:

Add at the top with other imports:
```javascript
import * as Gentleman from "./themes/gentleman";
```

Add at the bottom of the collection object:
```javascript
  Gentleman: Gentleman.theme,
```

### 6. Install simple-bar-server (for fast updates)

```bash
# Clone simple-bar-server
git clone https://github.com/Jean-Tinland/simple-bar-server.git ~/.config/simple-bar-server
cd ~/.config/simple-bar-server
npm install

# Copy launchd service
cp ../Gentleman.Dots2/simple-bar/com.simple-bar-server.plist ~/Library/LaunchAgents/

# Load service (starts automatically on login)
launchctl load -w ~/Library/LaunchAgents/com.simple-bar-server.plist
```

### 7. Configure simple-bar

1. Click on simple-bar (the status bar at the top)
2. Press `Cmd + ,` to open settings
3. Set **yabai path** to `/opt/homebrew/bin/yabai`
4. Enable **"Enable server"** option
5. In Themes section, select **"Gentleman"** as dark theme

### 8. Configure macOS menu bar and Dock

**Menu Bar:**
1. Open **System Settings**
2. Go to **Menu Bar** (in the left sidebar)
3. Set **"Automatically hide and show the menu bar"** to **"Always"**
4. Enable **"Show menu bar background"** (toggle ON)

This allows you to access the native macOS menu bar by moving your mouse to the top of the screen.

**Dock:**
1. Open **System Settings**
2. Go to **Desktop & Dock**
3. Enable **"Automatically hide and show the Dock"**
4. Disable **"Show suggested and recent apps in Dock"** (optional)

### 9. Grant Accessibility Permissions

**IMPORTANT**: Both yabai and skhd need accessibility permissions to work.

1. Open **System Settings**
2. Go to **Privacy & Security** → **Accessibility**
3. Click the **+** button
4. Navigate to `/opt/homebrew/bin/` and add `yabai`
5. Click **+** again and add `skhd` from the same location
6. Enable checkboxes next to both

### 10. Create Mission Control Spaces

macOS spaces must be created manually via Mission Control:

1. Press `F3` or `Ctrl+Up Arrow` to open Mission Control
2. Move cursor to the top-right corner
3. Click **+** to add desktops
4. Create a total of 7 desktops

### 11. Enable Mission Control Keyboard Shortcuts

1. Open **System Settings**
2. Go to **Keyboard** → **Keyboard Shortcuts** → **Mission Control**
3. Enable and set these shortcuts:
   - ☑ Switch to Desktop 1: `Ctrl+1` (^1)
   - ☑ Switch to Desktop 2: `Ctrl+2` (^2)
   - ☑ Switch to Desktop 3: `Ctrl+3` (^3)
   - ☑ Switch to Desktop 4: `Ctrl+4` (^4)
   - ☑ Switch to Desktop 5: `Ctrl+5` (^5)
   - ☑ Switch to Desktop 6: `Ctrl+6` (^6)
   - ☑ Switch to Desktop 7: `Ctrl+7` (^7)

### 12. Start Services

```bash
# Start services (will auto-start at login)
yabai --start-service
skhd --start-service
```

This will:
- Create launchd service files in `~/Library/LaunchAgents/`
- Start both services immediately
- Configure them to start automatically at login

## Keybindings

### Window Management
- `alt+ctrl+h/j/k/l` - Focus window (left/down/up/right)
- `alt+shift+h/j/k/l` - Move/swap window
- `alt+shift+f` - Toggle fullscreen

### Resizing
- `alt+arrow keys` - Resize windows in all directions
- `alt+-` / `alt+=` - Resize horizontally

### Layout
- `alt+r` - Rotate layout 90° (convert left/right to top/bottom)
- `alt+x` - Mirror horizontally
- `alt+y` - Mirror vertically
- `alt+e` - Toggle split orientation for focused window

### Workspaces/Spaces
- `alt+1-7` - Switch to workspace
- `alt+shift+1-7` - Move window to workspace and follow
- `alt+tab` - Switch to recent workspace

### Special
- `ctrl+alt+cmd+r` - Reload configs
- `ctrl+alt+cmd+b` - Balance layout
- `ctrl+alt+cmd+f` - Toggle floating
- `ctrl+alt+cmd+h/j/k/l` - Warp windows

## Space Labels
1. social
2. work
3. development
4. others
5. five
6. teleprompter
7. arzopa

## Service Management

```bash
# Restart services
yabai --restart-service
skhd --restart-service

# Stop services
yabai --stop-service
skhd --stop-service

# Start services
yabai --start-service
skhd --start-service

# Restart simple-bar-server
launchctl unload ~/Library/LaunchAgents/com.simple-bar-server.plist
launchctl load -w ~/Library/LaunchAgents/com.simple-bar-server.plist
```

## Gentleman Theme

The Gentleman theme for simple-bar matches the Ghostty terminal theme with colors:

- Background: `#06080f`
- Foreground: `#f3f6f9`
- Red: `#cb7c94`
- Green: `#b7cc85`
- Yellow: `#ffe066`
- Blue: `#7fb4ca`
- Magenta: `#ff8dd7`
- Cyan: `#7aa89f`

## Notes

- Configuration works with **SIP partially disabled** for full yabai features
- Uses simple-bar-server for efficient widget updates via curl
- Window tiling and management fully functional
- macOS menu bar accessible by moving mouse to top of screen
