# Manual Installation Guide

This guide walks you through manually setting up your development environment with Gentleman.Dots. Use this if you prefer full control over each step or if the automatic installer doesn't work for your setup.

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

Choose one:

- **Alacritty**: [Download from GitHub Releases](https://github.com/alacritty/alacritty/releases)
- **WezTerm**: [Download and Install](https://wezfurlong.org/wezterm/installation.html) - Create a `HOME` environment variable pointing to `C:\Users\your-username`
- **Kitty**: [Download and Install](https://sw.kovidgoyal.net/kitty/#get-the-app)

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
# config.front_end = "WebGpu"
# config.max_fps = 120
```

If WezTerm doesn't pick up the config:
- Create `C:\Users\your-username\.config\wezterm`
- Copy `.wezterm.lua` to `wezterm.lua` inside that directory
- Restart WezTerm

**Kitty:**

```powershell
Copy-Item -Path GentlemanKitty\* -Destination $HOME\.config\kitty -Recurse
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
wget -O ~/.local/share/fonts/Iosevka.zip https://github.com/ryanoasis/nerd-fonts/releases/download/v3.3.0/IosevkaTerm.zip
unzip ~/.local/share/fonts/Iosevka.zip -d ~/.local/share/fonts/
fc-cache -fv
```

#### macOS

```bash
brew tap homebrew/cask-fonts
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

# All platforms - copy config
mkdir -p ~/.config/wezterm && cp .wezterm.lua ~/.config/wezterm/wezterm.lua
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
mkdir -p ~/.config/kitty && cp -r GentlemanKitty/* ~/.config/kitty
```

**Note:** Reload Kitty config with `Ctrl+Shift+,` (Linux) or `Cmd+Shift+,` (macOS)

### 4. Install a Shell

#### Nushell

```bash
# Install dependencies
cp -rf bash-env-json ~/.config/
cp -rf bash-env.nu ~/.config/
brew install nushell carapace zoxide atuin jq bash starship fzf
cp -rf starship.toml ~/.config/

# Arch Linux / Linux
mkdir -p ~/.config/nushell
cp -rf GentlemanNushell/* ~/.config/nushell/

# macOS
mkdir -p ~/Library/Application\ Support/nushell

# Update config for macOS brew path
if grep -q "/home/linuxbrew/.linuxbrew/bin" GentlemanNushell/env.nu; then
  awk -v search="/home/linuxbrew/.linuxbrew/bin" -v replace="    | prepend '/opt/homebrew/bin'" '
  $0 ~ search {print replace; next}
  {print}
  ' GentlemanNushell/env.nu > GentlemanNushell/env.nu.tmp && mv GentlemanNushell/env.nu.tmp GentlemanNushell/env.nu
fi

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

**If using Zsh, update config:**

```bash
# Replace TMUX with ZELLIJ
if grep -q "TMUX" ~/.zshrc; then
  sed -i 's/TMUX/ZELLIJ/g' ~/.zshrc
  sed -i 's/tmux/zellij/g' ~/.zshrc
fi
```

**If using Fish, update config:**

```bash
# Replace TMUX with ZELLIJ in fish config
if grep -q "TMUX" ~/.config/fish/config.fish; then
  sed -i 's/TMUX/ZELLIJ/g' ~/.config/fish/config.fish
  sed -i 's/tmux/zellij/g' ~/.config/fish/config.fish
fi
```

**If using Nushell, update config:**

```bash
# For macOS
if grep -q '"tmux"' GentlemanNushell/config.nu; then
  sed -i 's/"tmux"/"zellij"/g' GentlemanNushell/config.nu
  sed -i 's/"TMUX"/"ZELLIJ"/g' GentlemanNushell/config.nu
fi
cp -rf GentlemanNushell/* ~/Library/Application\ Support/nushell/

# For Linux
if grep -q '"tmux"' ~/.config/nushell/config.nu; then
  sed -i 's/"tmux"/"zellij"/g' ~/.config/nushell/config.nu
  sed -i 's/"TMUX"/"ZELLIJ"/g' ~/.config/nushell/config.nu
fi
```

### 6. Install Neovim

```bash
brew install nvim node npm git gcc fzf fd ripgrep coreutils bat curl lazygit tree-sitter
mkdir -p ~/.config/nvim
cp -r GentlemanNvim/nvim/* ~/.config/nvim/
```

**Configure Obsidian path (optional):**

1. Create your notes directory:
   ```bash
   mkdir -p /path/to/your/notes/templates
   ```

2. Edit the Obsidian config:
   ```bash
   nvim ~/.config/nvim/lua/plugins/obsidian.lua
   ```

3. Update the path:
   ```lua
   path = "/path/to/your/notes",
   ```

### 7. Set Default Shell

```bash
# Choose your shell
shell_path=$(which zsh)   # or fish, or nu

# Add to /etc/shells and set as default
if [ -n "$shell_path" ]; then
  sudo sh -c "grep -Fxq \"$shell_path\" /etc/shells || echo \"$shell_path\" >> /etc/shells"
  sudo chsh -s "$shell_path" "$USER"
fi
```

### 8. Restart

Close and reopen your terminal, or restart your computer/WSL instance for changes to take effect.

---

**Done!** You've manually configured your development environment. Enjoy! 🎉

If you encounter problems, consult the official documentation of each tool or open an issue on GitHub.
