# Yabai + skhd Configuration

Tiling window manager setup for macOS with custom keybindings.

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

### 1. Install yabai and skhd

```bash
# Install via Homebrew
brew install koekeishiya/formulae/yabai
brew install koekeishiya/formulae/skhd
```

### 2. Create config directories

```bash
mkdir -p ~/.config/yabai
mkdir -p ~/.config/skhd
```

### 3. Copy config files

```bash
# From your dotfiles repo
cp yabairc ~/.config/yabai/yabairc
cp skhdrc ~/.config/skhd/skhdrc
cp move-window-to-space.sh ~/.config/yabai/move-window-to-space.sh

# Make helper script executable
chmod +x ~/.config/yabai/move-window-to-space.sh
chmod +x ~/.config/yabai/yabairc
```

### 4. Grant Accessibility Permissions

**IMPORTANT**: Both yabai and skhd need accessibility permissions to work.

1. Open **System Settings**
2. Go to **Privacy & Security** → **Accessibility**
3. Click the **+** button
4. Navigate to `/opt/homebrew/bin/` and add `yabai`
5. Click **+** again and add `skhd` from the same location
6. Enable checkboxes next to both

### 5. Create Mission Control Spaces

macOS spaces must be created manually via Mission Control:

1. Press `F3` or `Ctrl+Up Arrow` to open Mission Control
2. Move cursor to the top-right corner
3. Click **+** to add desktops
4. Create a total of 7 desktops

### 6. Enable Mission Control Keyboard Shortcuts

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

### 7. Start Services

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
5. stream
6. six
7. seven

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
```

## Notes

- Configuration works with **SIP enabled** (no need to disable System Integrity Protection)
- Uses Mission Control shortcuts as workaround for space switching
- Window tiling and management fully functional without scripting-addition
