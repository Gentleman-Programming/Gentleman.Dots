# Gentleman.Dots

> **Important Notice (January 2026)**: Anthropic has blocked third-party tools (OpenCode, Crush, etc.) from using Claude Max subscriptions. OAuth tokens are now restricted to Claude Code only. This config now uses **Claude Code as the primary AI assistant** in Neovim, and **OpenCode with API key** as secondary option.

<img width="1920" height="1080" alt="image" src="https://github.com/user-attachments/assets/9fcb1b15-89db-404e-b0f3-107801bd9115" />

---

## Description

This repository provides a complete, declarative development environment configuration using Nix Flakes and Home Manager. Everything is configured through local modules and automatically installs all dependencies.

### 🛠️ Development Tools & Languages

- **Languages**: Node.js, Bun, Cargo/Rust, Go, GCC
- **Package Managers**: Volta (Node.js), Cargo, Bun
- **Build Tools**: Nil (Nix LSP), Nixd (Nix language server)
- **Utilities**: jq, bash, fd, ripgrep, coreutils, unzip, bat, yazi

### 🐚 Shell Configurations

- **Fish Shell**: Complete configuration with 400+ command completions
- **Nushell**: Modern shell with custom config and environment setup
- **Zsh**: Traditional shell with modern enhancements
- **Starship**: Cross-shell prompt with custom configuration

### 🖥️ Terminal Emulators

- **Ghostty**: Modern GPU-accelerated terminal with custom themes
- **WezTerm**: Feature-rich terminal with Lua configuration
- **Tmux**: Terminal multiplexer with custom key bindings
- **Zellij**: Modern terminal workspace (optional - see customization)

### ⚡ Development Environment

- **Neovim**: Fully configured IDE with LazyVim, AI assistants, and 40+ plugins
- **Zed**: Modern code editor with custom themes and settings
- **Git & GitHub CLI**: Pre-configured version control
- **Lazy Git**: Terminal UI for Git operations

### 🤖 AI Integrations

- **Claude Code CLI**: Integrated AI coding assistant
- **OpenCode**: AI assistant integration
- **Gemini CLI**: Google's AI assistant (optional)
- **Multiple AI providers**: Support for various AI coding assistants

### 🔧 System Utilities

- **Television**: Modern file navigator and system monitor
- **Zoxide**: Smart directory jumping
- **Atuin**: Enhanced shell history
- **Carapace**: Universal shell completions
- **FZF**: Fuzzy finder integration
- **Raycast Scripts**: Custom automation scripts for display management

### 🪟 Window Management (macOS)

- **Yabai**: Tiling window manager with scripting support
- **Skhd**: Hotkey daemon for keyboard shortcuts
- **SketchyBar**: Customizable status bar with plugins
- **Aerospace**: Alternative tiling window manager (optional)

### 📝 Development Workflow

- **Oil.nvim**: Custom file navigation scripts
- **Custom Scripts**: Productivity-enhancing shell scripts
- **Nerd Fonts**: Iosevka Term for consistent typography
- **Declarative Configuration**: Everything version-controlled and reproducible

The flake automatically handles system-specific configurations, installs all dependencies, and sets up your complete development environment with a single command.

---

## Features Overview

### 🎯 What You Get

- **Zero Configuration**: Everything works out of the box
- **Declarative**: Version-controlled, reproducible environment
- **Modern Toolchain**: Latest development tools and utilities
- **AI-Enhanced**: Multiple AI coding assistants integrated
- **Shell Agnostic**: Fish, Nushell, and Zsh all configured
- **Terminal Flexibility**: Multiple terminal emulators supported
- **macOS Optimized**: Specifically tuned for macOS workflows

### 🔧 Technical Stack

| Category            | Tools                                     |
| ------------------- | ----------------------------------------- |
| **Package Manager** | Nix with Flakes + Home Manager            |
| **Shells**          | Fish, Nushell, Zsh with Starship prompt   |
| **Terminals**       | Ghostty, WezTerm, Tmux, Zellij (optional) |
| **Editor**          | Neovim (LazyVim) + Zed                    |
| **Languages**       | Node.js, Rust, Go, with Volta management  |
| **AI Tools**        | Claude Code, OpenCode, Gemini (opt.), multiple providers |
| **Navigation**      | Television, Yazi, Oil.nvim, Zoxide        |
| **Development**     | Git, GitHub CLI, Lazy Git                 |
| **Window Manager**  | Yabai + Skhd, SketchyBar, Aerospace (opt) |
| **Automation**      | Raycast Scripts                           |

### 📁 Project Structure

```
.
├── flake.nix              # Main Nix flake configuration
├── README.md              # This file
├── AGENTS.md              # AI agent instructions
│
├── # ─── Shell Configurations ───
├── fish.nix               # Fish shell configuration
├── fish/                  # Fish completions (400+ commands)
├── nushell.nix            # Nushell configuration
├── nushell/               # Nushell config and env files
├── zsh.nix                # Zsh configuration
├── starship.nix           # Starship prompt configuration
│
├── # ─── Terminal Emulators ───
├── ghostty.nix            # Ghostty terminal configuration
├── ghostty/               # Ghostty config, themes, and shaders
├── wezterm.nix            # WezTerm configuration
├── tmux.nix               # Tmux configuration
├── zellij.nix             # Zellij terminal workspace (optional)
├── zellij/                # Zellij plugins (zjstatus, forgot)
│
├── # ─── Editors ───
├── nvim.nix               # Neovim configuration
├── nvim/                  # Neovim plugins and settings (LazyVim)
├── nvim-oil-minimal/      # Minimal Neovim config for Oil.nvim
├── zed.nix                # Zed editor configuration
├── zed/                   # Zed themes, keymaps, tasks, prompts
│
├── # ─── AI Tools ───
├── claude.nix             # Claude Code CLI configuration
├── claude/                # Claude settings, skills, themes, statusline
├── opencode.nix           # OpenCode AI configuration
├── opencode/              # OpenCode themes and skills
├── gemini.nix             # Gemini CLI configuration (optional)
│
├── # ─── Window Management (macOS) ───
├── yabai.nix              # Yabai window manager configuration
├── yabai/                 # Yabai scripts and config
├── skhd.nix               # Skhd hotkey daemon configuration
├── skhd/                  # Skhd keybindings (skhdrc)
├── sketchybar.nix         # SketchyBar status bar configuration
├── sketchybar/            # SketchyBar plugins and config
├── simple-bar.nix         # simple-bar for Übersicht (disabled)
├── simple-bar/            # simple-bar themes (disabled)
├── aerospace/             # Aerospace window manager config (optional)
│
├── # ─── Utilities ───
├── television.nix         # Television file navigator
├── television/            # Television config and channels
├── oil-scripts.nix        # Custom Oil.nvim scripts
├── scripts/               # Custom utility scripts (ocd, oil)
├── raycast.nix            # Raycast scripts configuration
└── raycast/               # Raycast automation scripts
```

---

## 🐧 Linux Users: Important Configuration

**If you're using Linux, you MUST disable macOS-specific modules before running the installation.**

The following modules are **macOS-only** and will fail on Linux:

| Module | Description | Why macOS-only |
|--------|-------------|----------------|
| `yabai.nix` | Tiling window manager | Uses macOS Accessibility API |
| `skhd.nix` | Hotkey daemon | Depends on macOS input system |
| `sketchybar.nix` | Status bar | macOS menu bar integration |
| `simple-bar.nix` | Status bar widget | Requires Übersicht (macOS app) |
| `raycast.nix` | Automation scripts | Raycast is macOS-only |

### How to Disable macOS Modules

Edit `flake.nix` and comment out these lines in the `modules` array (around line 40-57):

```nix
modules = [
  ./nushell.nix
  ./ghostty.nix
  ./zed.nix
  ./television.nix
  ./wezterm.nix
  # ./zellij.nix  # Optional - uncomment if you want Zellij
  ./tmux.nix
  ./fish.nix
  ./starship.nix
  ./nvim.nix
  ./zsh.nix
  ./oil-scripts.nix
  ./opencode.nix
  ./claude.nix
  # ─── macOS ONLY - Comment these on Linux ───
  # ./yabai.nix        # ← Comment this line
  # ./skhd.nix         # ← Comment this line
  # ./simple-bar.nix   # ← Comment this line
  {
    # ... rest of config
```

Also remove the macOS window manager packages from `home.packages` (around line 73-75):

```nix
home.packages = with pkgs; [
  # ...
  
  # ─── Window management (macOS) ───
  # yabai   # ← Comment this line
  # skhd    # ← Comment this line
  
  # ...
];
```

### Additional Linux Changes

**1. Set your username** (around line 20):

```nix
# ─── User Configuration ───
# Change this to your Linux username
username = "YourUser";  # ← Replace with your username
```

**2. Change home directory path** (around line 67):

```nix
# Change from macOS path:
home.homeDirectory = "/Users/${username}";

# To Linux path:
home.homeDirectory = "/home/${username}";
```

**3. Add Linux system support** (around line 17):

```nix
# Change from:
supportedSystems = [ "x86_64-darwin" "aarch64-darwin" ];

# To (add your Linux architecture):
supportedSystems = [ "x86_64-darwin" "aarch64-darwin" "x86_64-linux" "aarch64-linux" ];
```

**4. Add Linux home configuration** (around line 140):

```nix
homeConfigurations = {
  # macOS system configurations
  "gentleman-macos-intel" = mkHomeConfiguration "x86_64-darwin";
  "gentleman-macos-arm" = mkHomeConfiguration "aarch64-darwin";
  
  # Linux system configurations (add these)
  "gentleman-linux" = mkHomeConfiguration "x86_64-linux";
  "gentleman-linux-arm" = mkHomeConfiguration "aarch64-linux";
  
  # Default to Apple Silicon
  "gentleman" = mkHomeConfiguration "aarch64-darwin";
};
```

**5. Run installation with Linux config:**

```bash
home-manager switch --flake .#gentleman-linux
```

### Linux Alternatives

For window management on Linux, consider:
- **i3/Sway** - Popular tiling window managers
- **Hyprland** - Modern Wayland compositor
- **bspwm** - Scriptable tiling window manager
- **Polybar/Waybar** - Status bars (replace simple-bar)

> **Note:** This configuration is primarily optimized for macOS. Linux support is possible but requires manual adjustment of these modules.

---

## 🖥️ SketchyBar Status Bar

This configuration includes a fully customized SketchyBar setup with:

- **Workspace indicators** with Yabai integration
- **System monitors** (CPU, memory, battery)
- **Media controls** and now playing info
- **Custom plugins** for various system information

### SketchyBar Plugins

| Plugin | Description |
|--------|-------------|
| `spaces.sh` | Workspace/space indicators with window counts |
| `front_app.sh` | Currently focused application |
| `media.sh` | Now playing media info |
| `battery.sh` | Battery status and percentage |
| `wifi.sh` | Network connection status |
| `clock.sh` | Date and time display |

### Starting SketchyBar

SketchyBar starts automatically via Nix. To manually control:

```bash
# Start
brew services start sketchybar

# Restart
brew services restart sketchybar

# Stop
brew services stop sketchybar
```

---

## ⚡ Raycast Scripts

Custom Raycast automation scripts for display and system management:

| Script | Description |
|--------|-------------|
| `set-4k.sh` | Set LG TV to 4K resolution with multi-monitor arrangement |
| `set-1080p.sh` | Set LG TV to 1080p resolution with multi-monitor arrangement |
| `reset-display-placement.sh` | Auto-detect resolution and reset display arrangement |
| `restart-sketchybar.sh` | Kill and restart SketchyBar |

### Setup

1. Open Raycast Settings → Extensions → Script Commands
2. Add `~/Raycast Scripts/` as a script directory
3. Scripts will be available in Raycast search

### Requirements

- [displayplacer](https://github.com/jakehilborn/displayplacer): `brew install displayplacer`

> **Note:** Display scripts use hardcoded monitor IDs. Run `displayplacer list` to get your monitor IDs and update the scripts accordingly.

---

## Installation Steps (for macOS)

### 1. Install the Nix Package Manager

```bash
sh <(curl -L https://nixos.org/nix/install)
```

### 2. Configure Nix to Use Extra Experimental Features

To enable the experimental features for flakes and the new `nix-command` (needed for our declarative setup), create/edit the configuration file:

```bash
# For daemon installation (default on macOS)
# The file may not exist, create it if needed
sudo mkdir -p /etc/nix
sudo nano /etc/nix/nix.conf
# Or: sudo vi /etc/nix/nix.conf
```

Add:

```
extra-experimental-features = flakes nix-command
build-users-group = nixbld
```

_(This is necessary because support for flakes and the new Nix command is still experimental, but it allows us to have a fully declarative and reproducible configuration.)_

### 3. Configure Your Username

**No need to edit `flake.nix` for system configuration!** The flake supports both Intel and Apple Silicon Macs.

You only need to update the `username` variable at the top of `flake.nix` (around line 20):

```nix
# ─── User Configuration ───
# Change this to your macOS username
username = "YourUser";  # ← Replace with your username
```

This single variable is used for both `home.username` and `home.homeDirectory`, so you only need to change it in one place.

### 4. Install Terminal Emulators (Optional)

Configurations are automatically applied. Choose your preferred terminal:

- **Ghostty** (Recommended): <https://ghostty.org/download>
  - Reload config with **Shift + Cmd + ,**
  - Modern GPU-accelerated with custom themes
  - Optimized for performance
  - **50+ custom shaders included** (CRT effects, cursor trails, matrix, etc.)

- **WezTerm**: <https://wezterm.org/installation.html>
  - Feature-rich with Lua configuration
  - Cross-platform compatibility
  - Advanced customization options

#### Ghostty Shaders

The configuration includes 50+ GLSL shaders for visual effects:

| Category | Examples |
|----------|----------|
| **CRT Effects** | `crt.glsl`, `bettercrt.glsl`, `retro-terminal.glsl` |
| **Cursor Trails** | `cursor_blaze.glsl`, `cursor_smear.glsl`, `cursor_frozen.glsl` |
| **Backgrounds** | `starfield.glsl`, `matrix-hallway.glsl`, `galaxy.glsl` |
| **Effects** | `bloom.glsl`, `glitchy.glsl`, `underwater.glsl` |

To enable a shader, edit `~/.config/ghostty/config` and add:
```
custom-shader = ~/.config/ghostty/shaders/cursor_smear_gentleman.glsl
```

### 5. Window Management Options

#### Option A: Yabai + Skhd (Recommended)

Yabai and Skhd are automatically configured via the flake. They provide:

- **Yabai**: Tiling window manager with BSP layout
- **Skhd**: Hotkey daemon for keyboard shortcuts
- **SketchyBar**: Status bar with workspace indicators

**Key bindings (configured in `skhd/skhdrc`):**

| Shortcut | Action |
|----------|--------|
| `alt + h/j/k/l` | Focus window (vim-style) |
| `shift + alt + h/j/k/l` | Move window |
| `alt + 1-9` | Switch to space |
| `shift + alt + 1-9` | Move window to space |
| `alt + f` | Toggle fullscreen |
| `alt + t` | Toggle float |

> **Note:** Yabai requires accessibility permissions and SIP configuration. See `yabai/README.md` for details.

#### Option B: Aerospace (Alternative)

For a simpler setup without SIP changes, copy the Aerospace configuration:

```bash
cp ./aerospace/.aerospace.toml ~/
```

Aerospace provides:

- Automatic window tiling
- Workspace management
- Keyboard-driven navigation
- macOS-native integration (no SIP changes needed)

### 6. Install Home Manager

Before running the flake configuration, you need to set up Home Manager:

```bash
# Add home-manager channel
nix-channel --add https://github.com/nix-community/home-manager/archive/master.tar.gz home-manager

# Update channels
nix-channel --update

# Install home-manager
nix-shell '<home-manager>' -A install
```

### 7. Run the Installation

Once you have cloned the repository and are **inside its directory**, run the command for your system.

**⚠️ Important:** You must be in the root of this project directory for the command to work, as it uses `.` to find the `flake.nix` file.

**For any Mac (the flake auto-detects your system):**

```bash
home-manager switch --flake .#gentleman
```

**Alternative: Specific system configurations:**

- **Apple Silicon Macs (M1/M2/M3/M4)**:
  ```bash
  home-manager switch --flake .#gentleman-macos-arm
  ```

- **Intel Macs**:
  ```bash
  home-manager switch --flake .#gentleman-macos-intel
  ```

_(These commands apply the configuration defined in the flake, installing all dependencies and applying the necessary settings.)_

### 8. Verify Installation

**PATH is configured automatically on macOS!**

### 9. Default Shell

Now run the following script to add `Nushell`, `Fish` or `Zsh` to your list of available shells and select it as the default one:

**Fish (Recommended):**

```bash
shellPath="$HOME/.local/state/nix/profiles/home-manager/home-path/bin/fish" && sudo sh -c "grep -Fxq '$shellPath' /etc/shells || echo '$shellPath' >> /etc/shells" && sudo chsh -s "$shellPath" "$USER"
```

**Nushell:**

```bash
shellPath="$HOME/.local/state/nix/profiles/home-manager/home-path/bin/nu" && sudo sh -c "grep -Fxq '$shellPath' /etc/shells || echo '$shellPath' >> /etc/shells" && sudo chsh -s "$shellPath" "$USER"
```

**Zsh:**

```bash
shellPath="$HOME/.local/state/nix/profiles/home-manager/home-path/bin/zsh" && sudo sh -c "grep -Fxq '$shellPath' /etc/shells || echo '$shellPath' >> /etc/shells" && sudo chsh -s "$shellPath" "$USER"
```

---

## Configuration Details

### 🔧 How It Works

- **Declarative Setup**: All configurations are defined in Nix modules
- **Automatic Deployment**: Files are copied to correct macOS locations
- **Dependency Management**: All tools and dependencies installed automatically
- **Version Pinning**: Reproducible builds with locked versions
- **System Integration**: Proper PATH configuration and shell integration

### 📍 File Locations

Configurations are automatically deployed to:

| Tool              | Location                                 |
| ----------------- | ---------------------------------------- |
| **Nushell**       | `~/Library/Application Support/nushell/` |
| **Fish**          | `~/.config/fish/`                        |
| **Ghostty**       | `~/.config/ghostty/`                     |
| **WezTerm**       | `~/.wezterm.lua`                         |
| **Neovim**        | `~/.config/nvim/`                        |
| **Zed**           | `~/Library/Application Support/Zed/`     |
| **Starship**      | `~/.config/starship.toml`                |
| **Tmux**          | `~/.config/tmux/`                        |
| **Zellij**        | `~/.config/zellij/` (optional)           |
| **Television**    | `~/.config/television/`                  |
| **Claude Code**   | `~/.claude/`                             |
| **OpenCode**      | `~/.config/opencode/`                    |
| **Yabai**         | `~/.config/yabai/`                       |
| **Skhd**          | `~/.config/skhd/`                        |
| **SketchyBar**    | `~/.config/sketchybar/`                  |
| **Raycast**       | `~/Raycast Scripts/`                     |

### 🚀 Performance Features

- **Shell Completions**: 400+ Fish completions for better productivity
- **Smart History**: Atuin for enhanced command history across shells
- **Fuzzy Finding**: FZF integration for quick file/command finding
- **Directory Navigation**: Zoxide for intelligent directory jumping
- **File Management**: Yazi and Television for modern file browsing
- **Git Workflow**: Lazy Git for streamlined version control

### 🤖 AI Development Features

- **Claude Code Integration**: Native AI coding assistant
- **Multiple AI Providers**: Support for various AI services
- **Context-Aware**: AI tools integrated with your development workflow
- **Productivity Focused**: AI assistants configured for maximum productivity

### 🎨 Theming & Customization

- **Consistent Themes**: Catppuccin and custom themes across all tools
- **Nerd Font Support**: Iosevka Term for perfect icon rendering
- **GPU Acceleration**: Modern terminals with hardware acceleration
- **Custom Key Bindings**: Vim-like navigation across all tools

## Troubleshooting

### Common Issues

**Command not found after installation:**

```bash
hash -r  # Refresh command cache
source ~/.zshrc  # or ~/.bashrc
```

**`nix: command not found` when running home-manager:**

This happens when the Nix binaries aren't in your PATH. Source the Nix daemon profile first:

```bash
# Source Nix profile
source /nix/var/nix/profiles/default/etc/profile.d/nix-daemon.sh

# Then run home-manager
home-manager switch --flake .#gentleman
```

**`Permission denied` when copying configs (sketchybar, nvim, etc):**

Some config files may be set as read-only. Fix permissions before running home-manager:

```bash
# For sketchybar:
chmod -R u+w ~/.config/sketchybar/

# For nvim:
chmod -R u+w ~/.config/nvim/

# Then re-run home-manager
home-manager switch --flake .#gentleman
```

**Nix installation issues:**

- Ensure `/etc/nix/nix.conf` has experimental features enabled
- Restart terminal after Nix installation
- Check that you're in the project directory when running commands

**Terminal not picking up themes:**

- For Ghostty: Use **Shift + Cmd + ,** to reload config
- For WezTerm: Restart the terminal
- Verify config files are in correct locations

### Customization

**Adding your own configurations:**

1. Edit the relevant `.nix` files
2. Run the installation command again
3. Changes are applied automatically

**Managing versions:**

- Update `flake.lock` with: `nix flake update`
- Pin specific package versions in the flake

**Enabling Optional Configurations:**

Some configurations are commented out by default. To enable them:

1. **Zellij Terminal Workspace:**
   ```bash
   # Edit flake.nix and uncomment the Zellij line
   sed -i '' 's|# ./zellij.nix|./zellij.nix|' flake.nix
   
   # Re-run the installation
   nix run github:nix-community/home-manager -- switch --flake .#gentleman-macos-arm -b backup
   ```
   
   Features:
   - Modern terminal multiplexer alternative to tmux
   - Vim-like keybindings with custom themes
   - Plugin system with status bar and session management
   - Custom layouts and workspace management

   **Additional Configuration Required:**
   After enabling Zellij, you need to update shell configurations to use Zellij instead of tmux:

   **Fish Shell (`~/.config/fish/config.fish`):**
   ```fish
   # Change line ~31 from:
   if not set -q TMUX; and not set -q ZED_TERMINAL
       tmux
   
   # To:
   if not set -q ZELLIJ; and not set -q ZED_TERMINAL
       zellij
   ```

   **Zsh Shell (`~/.zshrc`):**
   ```bash
   # Change lines ~100-102 from:
   WM_VAR="/$TMUX"
   WM_CMD="tmux"
   
   # To:
   WM_VAR="/$ZELLIJ"
   WM_CMD="zellij"
   ```

   **Nushell (`~/.config/nushell/config.nu`):**
   ```nu
   # Change lines ~1015-1016 from:
   let MULTIPLEXER = "tmux"
   let MULTIPLEXER_ENV_PREFIX = "TMUX"
   
   # To:
   let MULTIPLEXER = "zellij"
   let MULTIPLEXER_ENV_PREFIX = "ZELLIJ"
   ```

2. **Gemini CLI Integration:**
   ```bash
   # Edit flake.nix and add Gemini module
   # Add './gemini.nix' to the modules list in flake.nix
   
   # Re-run the installation
   nix run github:nix-community/home-manager -- switch --flake .#gentleman-macos-arm -b backup
   ```
   
   Features:
   - Google's AI assistant CLI tool
   - Integrated via Bun package manager
   - Direct access with `gemini` command
   - Perfect for AI-powered development workflows

## 🤖 Claude Code CLI Configuration

This configuration includes a complete Claude Code CLI setup with custom skills, themes, and output styles.

### Claude Code Features

| Feature | Description |
|---------|-------------|
| **CLAUDE.md** | Custom system instructions and personality |
| **Skills** | Framework-specific coding patterns (React 19, Next.js 15, etc.) |
| **Output Styles** | Custom response formatting (Gentleman style) |
| **Themes** | Custom TweakCC theme for terminal |
| **Statusline** | Custom status bar script |
| **MCP Servers** | Template for Model Context Protocol servers |

### Included Skills

The following skills auto-load based on context:

| Skill | Trigger |
|-------|---------|
| `react-19` | React components, hooks, JSX |
| `nextjs-15` | App Router, Server Components |
| `typescript` | Types, interfaces, generics |
| `tailwind-4` | Tailwind CSS styling |
| `zod-4` | Schema validation |
| `zustand-5` | State management |
| `ai-sdk-5` | Vercel AI SDK |
| `django-drf` | Django REST Framework |
| `playwright` | E2E testing |
| `pytest` | Python testing |

### OpenCode Configuration

OpenCode also includes similar skills plus:

| Skill | Description |
|-------|-------------|
| `jira-epic` | Create Jira epics following standard format |
| `jira-task` | Create Jira tasks following standard format |

---

## AI Configuration for Neovim

This configuration includes support for the following AI tools:

- **Avante.nvim** - AI-powered coding assistant
- **CopilotChat.nvim** - GitHub Copilot chat interface
- **OpenCode.nvim** - OpenCode AI integration
- **CodeCompanion.nvim** - Multi-AI provider support
- **Claude Code.nvim** - Claude AI integration *(enabled by default)*
- **Gemini.nvim** - Google Gemini integration

### How to Switch AI Plugins

**Claude Code is already enabled by default.** To switch to a different AI assistant:

1. **Navigate to the disabled plugins file:**
   ```bash
   nvim ~/.config/nvim/lua/plugins/disabled.lua
   ```

2. **Disable Claude Code** by changing `enabled = true` to `enabled = false`:
   ```lua
   {
     "greggh/claude-code.nvim",
     enabled = false,  -- Disable Claude Code
   },
   ```

3. **Enable your preferred AI assistant** by changing `enabled = false` to `enabled = true`:

   ```lua
   {
     "yetone/avante.nvim",
     enabled = true,  -- Change to true to enable
   },
   ```

4. **Save the file** and restart Neovim.

### Important Notes

- **Only enable ONE AI plugin at a time** to avoid conflicts and keybinding issues
- **Required CLI tools** are automatically installed by the script:
  - Claude Code CLI (`brew install --cask claude-code`)
  - OpenCode CLI (`curl -fsSL https://opencode.ai/install | bash`)
  - Gemini CLI (`brew install gemini-cli`)
- **API keys may be required** for some services - check each plugin's documentation
- **Node.js 18+** is required for most AI plugins (automatically handled by the configuration)

### Switching Between AI Assistants

To switch from one AI assistant to another:

1. Set your current AI plugin to `enabled = false`
2. Set your desired AI plugin to `enabled = true`
3. Restart Neovim

### Recommended AI Assistants

- **For beginners:** Start with **CodeCompanion.nvim** - supports multiple AI providers
- **For Claude users:** Use **Claude Code.nvim** with the Claude Code CLI
- **For GitHub Copilot users:** Use **CopilotChat.nvim**
- **For Google Gemini users:** Use **Gemini.nvim** with the Gemini CLI

### OpenCode Configuration

> **Why Claude Code is the default:** In January 2026, Anthropic restricted their OAuth tokens to only work with the official Claude Code CLI. Third-party tools like OpenCode, Crush, etc. were blocked from using Claude Max/Pro subscriptions.

**Recommended: Use OpenCode with API Key**

The safest way to use OpenCode is with an Anthropic API key:

```json
{
  "provider": "anthropic",
  "model": "anthropic/claude-sonnet-4-20250514"
}
```

Then set your API key:
```bash
export ANTHROPIC_API_KEY="your-api-key-here"
```

**Location:** `~/.config/opencode/opencode.json`

---

<details>
<summary>⚠️ <strong>RISKY ALTERNATIVE: opencode-anthropic-auth plugin</strong> (click to expand)</summary>

There's a community plugin (`opencode-anthropic-auth`) that enables OAuth authentication with Claude Max/Pro subscriptions. **However, this violates Anthropic's Terms of Service and can result in account suspension.**

**What the plugin does:**
- Spoofs the Claude Code CLI's OAuth client ID
- Fakes the user-agent to impersonate Claude Code
- Replaces "OpenCode" strings to evade server detection

**From Anthropic's Usage Policy:**
> "Do Not Abuse our Platform: Intentionally bypass capabilities, restrictions, or guardrails established within our products without prior authorization from Anthropic"

**Risks:**
- 🚫 **Account ban** - Your Claude account can be suspended
- ⚠️ You can request an unban via X/Twitter, but it's not guaranteed
- 🔄 Anthropic could block this workaround at any time

**If you still want to use it (at your own risk):**

```json
{
  "plugin": ["opencode-anthropic-auth"],
  "model": "anthropic/claude-sonnet-4-20250514"
}
```

**We do NOT recommend this approach.** Use Claude Code CLI or OpenCode with an API key instead.

</details>

## Contributing

Contributions welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Test your changes thoroughly
4. Submit a pull request

For questions or issues, open a GitHub issue.

---

**Happy coding!** 🚀

— Gentleman Programming
