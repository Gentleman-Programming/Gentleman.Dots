# Manual Installation Guide

This guide walks you through manually setting up your development environment with Gentleman.Dots. Use this if you prefer full control over each step or if the automatic installer doesn't work for your setup.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Windows (WSL)](#windows-wsl)
  - [Install WSL](#1-install-wsl)
  - [Install a Linux Distribution](#2-install-a-linux-distribution)
  - [Install the Iosevka Font](#3-install-the-iosevka-font)
  - [Update the Distribution](#4-update-the-distribution)
  - [Install a Terminal Emulator](#5-install-a-terminal-emulator)
  - [Configure Terminal Emulator](#6-configure-terminal-emulator)
  - [Install Chocolatey and win32yank](#7-install-chocolatey-and-win32yank)
- [Linux, Arch Linux, macOS, and WSL](#linux-arch-linux-macos-and-wsl)
  - [Install Dependencies](#1-install-dependencies)
  - [Install Iosevka Term Nerd Font](#2-install-iosevka-term-nerd-font)
  - [Install Terminal Emulator](#3-install-terminal-emulator)
  - [Install a Shell](#4-install-a-shell)
  - [Install Window Manager](#5-install-window-manager)
  - [Install Neovim](#6-install-neovim)
  - [Set Default Shell](#7-set-default-shell)
  - [Restart](#8-restart)

---

## Prerequisites

**Clone the repository first!**

```bash
git clone git@github.com:Gentleman-Programming/Gentleman.Dots.git
cd Gentleman.Dots
```

---

## Windows (WSL)

> **Important:** Windows users must follow these manual installation steps. The automatic installer requires WSL.

### 1. Install WSL

WSL (Windows Subsystem for Linux) allows you to run Linux on Windows:

```powershell
wsl --install
wsl --set-default-version 2
```

### 2. Install a Linux Distribution

```powershell
wsl --install -d Ubuntu
```

To list available distributions:

```powershell
wsl --list --online
wsl --install -d <distribution-name>
```

### 3. Install the Iosevka Font

The Iosevka Term Nerd Font is required for terminal icons. On Windows, install manually:

1. Download from [Nerd Fonts GitHub](https://github.com/ryanoasis/nerd-fonts/releases) - look for `IosevkaTerm.zip`
2. Extract the archive
3. Right-click each font file and select **"Install for all users"**

### 4. Update the Distribution

Open the installed distribution and run:

```bash
sudo apt-get update
sudo apt-get upgrade
```

### 5. Install a Terminal Emulator

Choose one of the following:

| Terminal | Download | Notes |
|----------|----------|-------|
| **Alacritty** | [GitHub Releases](https://github.com/alacritty/alacritty/releases) | Lightweight, GPU-accelerated |
| **WezTerm** | [Official Site](https://wezfurlong.org/wezterm/installation.html) | Create `HOME` env var → `C:\Users\your-username` |
| **Kitty** | [Official Site](https://sw.kovidgoyal.net/kitty/#get-the-app) | Feature-rich, GPU-based |

### 6. Configure Terminal Emulator

Using PowerShell:

**Alacritty:**

```powershell
mkdir $env:APPDATA\alacritty
Copy-Item -Path alacritty.toml -Destination $env:APPDATA\alacritty\alacritty.toml

# In alacritty.toml, uncomment:
# [shell]
# program = "wsl.exe"
# args = ["--cd", "~"]
```

**WezTerm:**

```powershell
Copy-Item -Path .wezterm.lua -Destination $HOME

# In .wezterm.lua, uncomment for Windows:
# config.default_domain = 'WSL:Ubuntu'
```

**Transparency issues on Windows?** WezTerm has known issues with certain GPU drivers. Try different `front_end` values in order:
1. `config.front_end = "OpenGL"` (default, try first)
2. `config.front_end = "WebGpu"` (if OpenGL has issues)
3. `config.front_end = "Software"` (fallback, uses CPU)

For Windows Acrylic effect, set `window_background_opacity = 0` (not 0.95).

If WezTerm doesn't pick up the config:
- Create `C:\Users\your-username\.config\wezterm`
- Copy `.wezterm.lua` to `wezterm.lua` inside that directory
- Restart WezTerm

**Kitty:**

```powershell
mkdir $HOME\.config\kitty -Force
Copy-Item -Path GentlemanKitty\kitty.conf -Destination $HOME\.config\kitty\kitty.conf
```

### 7. Install Chocolatey and win32yank

**Chocolatey** (package manager for Windows):

```powershell
Set-ExecutionPolicy Bypass -Scope Process -Force; `
[System.Net.ServicePointManager]::SecurityProtocol = `
[System.Net.ServicePointManager]::SecurityProtocol -bor 3072; `
iwr https://community.chocolatey.org/install.ps1 -UseBasicParsing | iex
```

**win32yank** (for Neovim clipboard integration):

```powershell
choco install win32yank
```

---

## Linux, Arch Linux, macOS, and WSL

### 1. Install Dependencies

#### Arch Linux

```bash
sudo pacman -Syu --noconfirm
sudo pacman -S --needed --noconfirm base-devel curl file git wget
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
. $HOME/.cargo/env
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

#### Linux (Debian/Ubuntu)

```bash
sudo apt-get update
sudo apt-get install -y build-essential curl file git
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
. $HOME/.cargo/env
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

#### macOS

```bash
xcode-select --install
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
. $HOME/.cargo/env
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

### 2. Install Iosevka Term Nerd Font

#### Arch Linux / Linux

```bash
mkdir -p ~/.local/share/fonts
# Download latest version from https://github.com/ryanoasis/nerd-fonts/releases
wget -O ~/.local/share/fonts/Iosevka.zip https://github.com/ryanoasis/nerd-fonts/releases/latest/download/IosevkaTerm.zip
unzip ~/.local/share/fonts/Iosevka.zip -d ~/.local/share/fonts/
fc-cache -fv
```

#### macOS

```bash
brew install --cask font-iosevka-term-nerd-font
```

### 3. Install Terminal Emulator

#### Alacritty

```bash
# Arch Linux
sudo pacman -S --noconfirm alacritty

# macOS
brew install alacritty --cask

# Linux (Debian/Ubuntu)
sudo add-apt-repository ppa:aslatter/ppa
sudo apt update
sudo apt install alacritty

# All platforms - copy config
mkdir -p ~/.config/alacritty && cp alacritty.toml ~/.config/alacritty/alacritty.toml
```

#### WezTerm

```bash
# Arch Linux
sudo pacman -S --noconfirm wezterm

# macOS
brew install wezterm --cask

# Linux
brew tap wez/wezterm-linuxbrew
brew install wezterm

# All platforms - copy config (to home directory)
cp .wezterm.lua ~/.wezterm.lua
```

#### Ghostty

```bash
# Arch Linux
pacman -S ghostty

# macOS / Linux
brew install --cask ghostty

# All platforms - copy config
mkdir -p ~/.config/ghostty && cp -r GentlemanGhostty/* ~/.config/ghostty
```

#### Kitty

```bash
# macOS
brew install --cask kitty

# Copy config
mkdir -p ~/.config/kitty && cp GentlemanKitty/kitty.conf ~/.config/kitty/
```

**Note:** Reload Kitty config with `Ctrl+Shift+,` (Linux) or `Cmd+Shift+,` (macOS)

### 4. Install a Shell

#### Nushell

**Install dependencies:**

```bash
brew install nushell carapace zoxide atuin jq bash starship fzf
cp -rf bash-env-json ~/.config/
cp -rf bash-env.nu ~/.config/
cp -rf starship.toml ~/.config/
```

**Arch Linux / Linux:**

```bash
mkdir -p ~/.config/nushell
cp -rf GentlemanNushell/* ~/.config/nushell/
```

**macOS:**

```bash
mkdir -p ~/Library/Application\ Support/nushell

# Update brew path from Linux to macOS
# Edit GentlemanNushell/env.nu and replace:
#   /home/linuxbrew/.linuxbrew/bin  →  /opt/homebrew/bin

cp -rf GentlemanNushell/* ~/Library/Application\ Support/nushell/
```

#### Fish + Starship

```bash
brew install fish carapace zoxide atuin starship fzf
mkdir -p ~/.cache/starship
mkdir -p ~/.cache/carapace
mkdir -p ~/.local/share/atuin
cp -rf starship.toml ~/.config/
cp -rf GentlemanFish/fish ~/.config
```

#### Zsh + Powerlevel10k

```bash
brew install zsh carapace zoxide atuin fzf
brew install zsh-autosuggestions zsh-syntax-highlighting zsh-autocomplete
mkdir -p ~/.cache/carapace
mkdir -p ~/.local/share/atuin
cp -rf GentlemanZsh/.zshrc ~/
cp -rf GentlemanZsh/.p10k.zsh ~/
cp -rf GentlemanZsh/.oh-my-zsh ~/
brew install powerlevel10k
```

### 5. Install Window Manager

#### Tmux

```bash
brew install tmux
git clone https://github.com/tmux-plugins/tpm ~/.tmux/plugins/tpm
mkdir -p ~/.tmux
cp -r GentlemanTmux/.tmux/* ~/.tmux/
cp GentlemanTmux/.tmux.conf ~/
tmux new-session -d -s plugin-installation 'source ~/.tmux.conf; tmux run-shell ~/.tmux/plugins/tpm/bin/install_plugins'
tmux kill-session -t plugin-installation
```

#### Zellij

```bash
cargo install zellij
mkdir -p ~/.config/zellij
cp -r GentlemanZellij/zellij/* ~/.config/zellij/
```

**Update shell config for Zellij:**

If you chose Zellij instead of Tmux, update your shell configuration:

| Shell | Config File | Change |
|-------|-------------|--------|
| Zsh | `~/.zshrc` | Replace `TMUX` → `ZELLIJ` and `tmux` → `zellij` |
| Fish | `~/.config/fish/config.fish` | Replace `TMUX` → `ZELLIJ` and `tmux` → `zellij` |
| Nushell (Linux) | `~/.config/nushell/config.nu` | Replace `"tmux"` → `"zellij"` and `"TMUX"` → `"ZELLIJ"` |
| Nushell (macOS) | `~/Library/Application Support/nushell/config.nu` | Same as above |

### 6. Install Neovim

```bash
brew install nvim node npm git gcc fzf fd ripgrep coreutils bat curl lazygit tree-sitter
mkdir -p ~/.config/nvim
cp -r GentlemanNvim/nvim/* ~/.config/nvim/
```

**Configure Obsidian path (optional):**

If you use Obsidian for notes, configure the integration:

```bash
# 1. Create your notes directory
mkdir -p /path/to/your/notes/templates

# 2. Edit ~/.config/nvim/lua/plugins/obsidian.lua
# 3. Update the path variable:
#    path = "/path/to/your/notes",
```

### 7. Set Default Shell

```bash
# Get path to your preferred shell (zsh, fish, or nu)
shell_path=$(which zsh)

# Add to allowed shells if not present, then set as default
sudo sh -c "grep -Fxq \"$shell_path\" /etc/shells || echo \"$shell_path\" >> /etc/shells"
sudo chsh -s "$shell_path" "$USER"
```

> **Note:** Replace `zsh` with `fish` or `nu` depending on which shell you installed.

### 8. Restart

Close and reopen your terminal, or restart your computer/WSL instance for changes to take effect.

---

## Troubleshooting

### Copy/Paste Not Working in Zellij (Linux)

If you can't copy text from the terminal when using Zellij on Linux:

1. Edit `~/.config/zellij/config.kdl`
2. Uncomment the appropriate line for your system:
   - **X11**: `copy_command "xclip -selection clipboard"` AND `copy_clipboard "primary"`
   - **Wayland**: `copy_command "wl-copy"`

See: [Zellij FAQ](https://zellij.dev/documentation/faq.html#copy--paste-isnt-working-how-can-i-fix-this)

### Fish/Zsh Fails in WSL with "missing or unsuitable terminal"

When using WezTerm on Windows with WSL, Fish or Zsh may fail with:
```
missing or unsuitable terminal: wezterm
```

**Status**: ✅ Fixed automatically in the default `.wezterm.lua` config (v2.7.7+).

If you're using an older config, update your `.wezterm.lua` to include the auto-detection:
```lua
if wezterm.target_triple:find("windows") then
  config.term = "xterm-256color"
else
  config.term = "wezterm"
end
```

### Other Issues

If you encounter other problems:

1. Consult the official documentation of the specific tool
2. Open an issue on [GitHub](https://github.com/Gentleman-Programming/Gentleman.Dots/issues)

---

**Done!** You've manually configured your development environment. Enjoy! 🎉
