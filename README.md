# Gentleman.Dots

<img width="2998" height="1649" alt="image" src="https://github.com/user-attachments/assets/0ef4e8fb-e08c-4554-9028-43c36c79acde" />

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
- **Gemini CLI**: Google's AI assistant (optional - see customization)
- **Multiple AI providers**: Support for various AI coding assistants

### 🔧 System Utilities

- **Television**: Modern file navigator and system monitor
- **Zoxide**: Smart directory jumping
- **Atuin**: Enhanced shell history
- **Carapace**: Universal shell completions
- **FZF**: Fuzzy finder integration
- **Aerospace**: Tiling window manager configuration (optional)

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

### 📁 Project Structure

```
.
├── flake.nix              # Main Nix flake configuration
├── README.md               # This file
├── fish.nix               # Fish shell configuration
├── nushell.nix            # Nushell configuration
├── zsh.nix                # Zsh configuration
├── starship.nix           # Starship prompt configuration
├── nvim.nix               # Neovim configuration
├── ghostty.nix            # Ghostty terminal configuration
├── wezterm.nix            # WezTerm configuration
├── tmux.nix               # Tmux configuration
├── zed.nix                # Zed editor configuration
├── claude.nix             # Claude Code CLI configuration
├── opencode.nix           # OpenCode AI configuration
├── gemini.nix             # Gemini CLI configuration (optional)
├── television.nix         # Television file navigator
├── zellij.nix             # Zellij terminal workspace (optional)
├── oil-scripts.nix        # Custom Oil.nvim scripts
├── fish/                  # Fish completions and configs
├── nvim/                  # Neovim plugins and settings
├── ghostty/               # Ghostty themes and config
├── zed/                   # Zed themes and settings
├── scripts/               # Custom utility scripts
└── aerospace/             # Aerospace window manager config
```

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

### 3. Prepare Your System

**No need to edit `flake.nix` for system configuration!** The flake supports both Intel and Apple Silicon Macs.

You only need to update your username in `flake.nix`:

- Change `home.username = "YourUser";` to your actual username
- The home directory is automatically set to `/Users/YourUser`

### 4. Install Terminal Emulators (Optional)

Configurations are automatically applied. Choose your preferred terminal:

- **Ghostty** (Recommended): <https://ghostty.org/download>
  - Reload config with **Shift + Cmd + ,**
  - Modern GPU-accelerated with custom themes
  - Optimized for performance

- **WezTerm**: <https://wezterm.org/installation.html>
  - Feature-rich with Lua configuration
  - Cross-platform compatibility
  - Advanced customization options

### 5. Optional: Aerospace Window Manager

For tiling window management, copy the Aerospace configuration:

```bash
cp ./aerospace/.aerospace.toml ~/
```

Aerospace provides:

- Automatic window tiling
- Workspace management
- Keyboard-driven navigation
- macOS-native integration

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

| Tool           | Location                                 |
| -------------- | ---------------------------------------- |
| **Nushell**    | `~/Library/Application Support/nushell/` |
| **Fish**       | `~/.config/fish/`                        |
| **Ghostty**    | `~/.config/ghostty/`                     |
| **WezTerm**    | `~/.wezterm.lua`                         |
| **Neovim**     | `~/.config/nvim/`                        |
| **Zed**        | `~/Library/Application Support/Zed/`     |
| **Starship**   | `~/.config/starship.toml`                |
| **Tmux**       | `~/.config/tmux/`                        |
| **Zellij**     | `~/.config/zellij/` (optional)           |
| **Television** | `~/.config/television/`                  |

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
