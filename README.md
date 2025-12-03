# Gentleman.Dots

<img width="2998" height="1649" alt="image" src="https://github.com/user-attachments/assets/c5a1eae2-69de-4ca1-8b4d-9d9b56e4cb5a" />

## Description

This repository contains customized configurations for a complete development environment, including:

- Neovim, with LSP, autocompletion, and integrated AI assistants (Claude, Gemini)
- Fish
- Zsh
- **Nushell**
- Tmux
- Zellij
- Terminal emulators:
  - Alacritty
  - WezTerm
  - Kitty
  - Ghostty

You can choose between automatic and manual installation methods depending on your preference and operating system.

**Want the bleeding edge?** If you want access to the latest experimental features that I use in my daily workflow (macOS only), check out the `nix-migration` branch: https://github.com/Gentleman-Programming/Gentleman.Dots/tree/nix-migration

This branch contains cutting-edge configurations that eventually make their way to `main` once they're stable and tested. If you can't wait for the stable release and want the absolute latest features, this branch is incredible - just keep in mind it's macOS-focused and may have occasional instabilities.

**Important:** Windows users **must** follow the manual installation instructions before running the script.

## Installation (Automatic Recommended!)

### The Easy Way! Test the automated process and let the script do all the work for you 😘

The **automatic installation script** is the quickest and easiest way to set up your development environment. This script handles all the heavy lifting, but remember that you **must install the font** mentioned below before running it. The script is designed for macOS, Linux, and WSL systems. If you’re on Windows, you’ll need to follow the manual steps first before attempting to run this script.

```bash
curl -O https://raw.githubusercontent.com/Gentleman-Programming/Gentleman.Dots/refs/heads/main/install-linux-mac.sh

sudo chmod +x install-linux-mac.sh
bash ./install-linux-mac.sh
```

## Manual Installation

Welcome to the Gentleman.Dots manual configuration guide! This document will walk you through the steps required to set up your development environment.

**_Clone the repo before continuing!!!_**

```bash
git clone git@github.com:Gentleman-Programming/Gentleman.Dots.git
cd Gentleman.Dots
```

---

<details>
<summary>For Windows</summary>

**Important:** Windows users must follow these manual installation steps before running the automated script.

#### 1. Install WSL

WSL (Windows Subsystem for Linux) allows you to run Linux on Windows. Install it and set it to version 2:

```powershell
wsl --install
wsl --set-default-version 2
```

#### 2. Install a Linux Distribution

Install a Linux distribution (e.g., Ubuntu) in WSL:

```powershell
wsl --install -d Ubuntu
```

To list available distributions:

```powershell
wsl --list --online
```

Install your preferred distribution:

```powershell
wsl --install -d <distribution-name>
```

#### 3. Installing the Iosevka Font

The Iosevka Term Nerd Font is required for terminal emulators in this setup. On Windows, this installation must be done manually.

1. **Download the Iosevka font** from its official site or from [Nerd Fonts GitHub](https://github.com/ryanoasis/nerd-fonts).
2. **Extract the archive** and locate the font files (`.ttf` or `.otf`).
3. **Install the fonts**:
   - Right-click each font file and select **"Install for all users"** to install the font system-wide.

#### 4. Launch and Configure the Distribution

Open the installed distribution to complete setup. Update it with:

```bash
sudo apt-get update
sudo apt-get upgrade
```

#### 5. Install a Terminal Emulator

Choose and install one of the following terminal emulators:

- **Alacritty**: [Download from GitHub Releases](https://github.com/alacritty/alacritty/releases) and place `alacritty.exe` in your `PATH`.
- **WezTerm**: [Download and Install](https://wezfurlong.org/wezterm/installation.html) and create an environment variable called `HOME` that resolves to `C:\Users\your-username`.
- **Kitty**: [Download and Install](https://sw.kovidgoyal.net/kitty/#get-the-app).

#### 6. Configuration Transfer for Terminal Emulators

Using PowerShell:

**Alacritty Configuration**

```powershell
mkdir $env:APPDATA\alacritty
Copy-Item -Path alacritty.toml -Destination $env:APPDATA\alacritty\alacritty.toml

# In alacritty.toml, uncomment and set the shell program to WSL:

#[shell]
#program = "wsl.exe"
#args = ["--cd", "~"]
```

**WezTerm Configuration**

```powershell
Copy-Item -Path .wezterm.lua -Destination $HOME

# Uncomment for Windows settings in .wezterm.lua:

# config.default_domain = 'WSL:Ubuntu'
# config.front_end = "WebGpu"
# config.max_fps = 120
```

If WezTerm doesn't take the initial configuration:

- Create a `wezterm` folder in `C:\Users\your-username\.config`
- Copy `.wezterm.lua` into `wezterm.lua` inside that directory
- Restart WezTerm

**Kitty Configuration**

```powershell
Copy-Item -Path GentlemanKitty\* -Destination $HOME\.config\kitty -Recurse
```

#### 7. Install Chocolatey and win32yank

**Chocolatey** is a package manager for Windows that simplifies the installation of software.

**To install Chocolatey:**

- Open **PowerShell** as an administrator.
- Run the following command:

```powershell
Set-ExecutionPolicy Bypass -Scope Process -Force; `
[System.Net.ServicePointManager]::SecurityProtocol = `
[System.Net.ServicePointManager]::SecurityProtocol -bor 3072; `
iwr https://community.chocolatey.org/install.ps1 -UseBasicParsing | iex
```

**To install win32yank:**

- After installing Chocolatey, run:

```powershell
choco install win32yank
```

**Note:** `win32yank` is required for clipboard integration in Neovim when using WSL.
</details>

---

<details>
<summary>For Linux, Arch Linux, macOS, and WSL</summary>

#### Prerequisites

- **macOS or Linux:** Ensure you have one of these operating systems.
- **Administrator privileges (sudo):** You'll need administrator access to install some tools.

#### 1. Install Dependencies

###### Arch Linux

```bash
sudo pacman -Syu --noconfirm
sudo pacman -S --needed --noconfirm base-devel curl file git wget
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
. $HOME/.cargo/env
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

###### Linux

```bash
sudo apt-get update
sudo apt-get install -y build-essential curl file git
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
. $HOME/.cargo/env
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

###### Mac

```bash
xcode-select --install
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
. $HOME/.cargo/env
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

#### 2. Install Iosevka Term Nerd Font (icons and font style)

###### Arch Linux / Linux

```bash
mkdir -p ~/.local/share/fonts
wget -O ~/.local/share/fonts/Iosevka.zip https://github.com/ryanoasis/nerd-fonts/releases/download/v3.3.0/IosevkaTerm.zip
unzip ~/.local/share/fonts/Iosevka.zip -d ~/.local/share/fonts/
fc-cache -fv
```

###### Mac

```bash
brew tap homebrew/cask-fonts
brew install --cask font-iosevka-term-nerd-font
```

#### 3. Choose and Install Terminal Emulator

##### Alacritty

###### Arch Linux

```bash
sudo pacman -S --noconfirm alacritty
mkdir -p ~/.config/alacritty && cp alacritty.toml ~/.config/alacritty/alacritty.toml
```

###### Mac

```bash
brew install alacritty --cask
mkdir -p ~/.config/alacritty && cp alacritty.toml ~/.config/alacritty/alacritty.toml
```

###### Linux

```bash
sudo add-apt-repository ppa:aslatter/ppa; sudo apt update; sudo apt install alacritty
mkdir -p ~/.config/alacritty && cp alacritty.toml ~/.config/alacritty/alacritty.toml
```

##### WezTerm

###### Arch Linux

```bash
sudo pacman -S --noconfirm wezterm
mkdir -p ~/.config/wezterm && cp .wezterm.lua ~/.config/wezterm/wezterm.lua
```

###### Mac

```bash
brew install wezterm --cask
mkdir -p ~/.config/wezterm && cp .wezterm.lua ~/.config/wezterm/wezterm.lua
```

###### Linux

```bash
brew tap wez/wezterm-linuxbrew; brew install wezterm
mkdir -p ~/.config/wezterm && cp .wezterm.lua ~/.config/wezterm/wezterm.lua
```

##### Ghostty

###### Arch Linux

```bash
pacman -S ghostty
mkdir -p ~/.config/ghostty && cp -r GentlemanGhostty/* ~/.config/ghostty
```

###### Mac

```bash
brew install --cask ghostty
mkdir -p ~/.config/ghostty && cp -r GentlemanGhostty/* ~/.config/ghostty
```

###### Linux

```bash
brew install --cask ghostty
mkdir -p ~/.config/ghostty && cp -r GentlemanGhostty/* ~/.config/ghostty
```

##### Kitty

###### Mac

```bash
brew install --cask kitty
mkdir -p ~/.config/kitty && cp -r GentlemanKitty/* ~/.config/kitty
```

**Reload the config after install doing `ctrl+shift+,` | `cmd+shift+,`**

#### 4. Choose and Install a Shell

##### Nushell

###### 1. Step

```bash
cp -rf bash-env-json ~/.config/
cp -rf bash-env.nu ~/.config/
brew install nushell carapace zoxide atuin jq bash starship fzf
cp -rf starship.toml ~/.config/
```

###### 2. Step

**_Arch Linux / Linux_**

```bash
mkdir -p ~/.config/nushell
run_command "cp -rf GentlemanNushell/* ~/.config/nushell/"
```

**_Mac_**

```bash
mkdir -p ~/Library/Application\ Support/nushell

## udpate config to use mac
if grep -q "/home/linuxbrew/.linuxbrew/bin" GentlemanNushell/env.nu; then
  awk -v search="/home/linuxbrew/.linuxbrew/bin" -v replace="    | prepend '/opt/homebrew/bin'" '
  $0 ~ search {print replace; next}
  {print}
  ' GentlemanNushell/env.nu > GentlemanNushell/env.nu.tmp && mv GentlemanNushell/env.nu.tmp GentlemanNushell/env.nu
else
  echo "    | prepend '/opt/homebrew/bin'" >> GentlemanNushell/env.nu
fi

cp -rf GentlemanNushell/* ~/Library/Application\ Support/nushell/
```

###### Fish + Starship

```bash
brew install fish carapace zoxide atuin starship fzf
mkdir -p ~/.cache/starship
mkdir -p ~/.cache/carapace
mkdir -p ~/.local/share/atuin
cp -rf starship.toml ~/.config/
cp -rf GentlemanFish/fish ~/.config
```

###### Zsh + Power10k\*\*

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

#### 5. Choose and Install Window Manager

##### Tmux

```bash
brew install tmux
git clone https://github.com/tmux-plugins/tpm ~/.tmux/plugins/tpm
mkdir -p ~/.tmux
cp -r GentlemanTmux/.tmux/* ~/.tmux/
cp GentlemanTmux/.tmux.conf ~/
tmux new-session -d -s plugin-installation 'source ~/.tmux.conf; tmux run-shell ~/.tmux/plugins/tpm/bin/install_plugins'
tmux kill-session -t plugin-installation
```

##### Zellij

###### 1. Step => Install Zellij

```bash
cargo install zellij
mkdir -p ~/.config/zellij
cp -r GentlemanZellij/zellij/* ~/.config/zellij/
```

###### 2. Step => If you use ZSH

```bash
# update or replace TMUX
if grep -q "TMUX" ~/.zshrc; then
  awk -v search="TMUX" -v replace='WM_VAR="/$ZELLIJ"' '
  $0 ~ search {print replace; next}
  {print}
  ' ~/.zshrc > ~/.zshrc.tmp && mv ~/.zshrc.tmp ~/.zshrc
else
  echo 'WM_VAR="/$ZELLIJ"' >> ~/.zshrc
fi

# update or replace tmux
if grep -q "tmux" ~/.zshrc; then
  awk -v search="tmux" -v replace='WM_CMD="zellij"' '
  $0 ~ search {print replace; next}
  {print}
  ' ~/.zshrc > ~/.zshrc.tmp && mv ~/.zshrc.tmp ~/.zshrc
else
  echo 'WM_CMD="zellij"' >> ~/.zshrc
fi
```

###### 3. Step => If you use Fish

```bash
if grep -q "TMUX" ~/.config/fish/config.fish; then
  awk -v search="TMUX" -v replace="if not set -q ZELLIJ" '
  $0 ~ search {print replace; next}
  {print}
  ' ~/.config/fish/config.fish > ~/.config/fish/config.fish.tmp && mv ~/.config/fish/config.fish.tmp ~/.config/fish/config.fish
else
  echo "if not set -q ZELLIJ" >> ~/.config/fish/config.fish
fi

# update or replace tmux
if grep -q "tmux" ~/.config/fish/config.fish; then
  awk -v search="tmux" -v replace="zellij" '
  $0 ~ search {print replace; next}
  {print}
  ' ~/.config/fish/config.fish > ~/.config/fish/config.fish.tmp && mv ~/.config/fish/config.fish.tmp ~/.config/fish/config.fish
else
  echo "zellij" >> ~/.config/fish/config.fish
fi
```

###### 3. Step => If you use Nushell

**_Mac_**

```bash
# update or replace "tmux"
if grep -q '"tmux"' GentlemanNushell/config.nu; then
  awk -v search='"tmux"' -v replace='let MULTIPLEXER = "zellij"' '
  $0 ~ search {print replace; next}
  {print}
  ' GentlemanNushell/config.nu > GentlemanNushell/config.nu.tmp && mv GentlemanNushell/config.nu.tmp GentlemanNushell/config.nu
else
  echo 'let MULTIPLEXER = "zellij"' >> GentlemanNushell/config.nu
fi

# update or replace "TMUX"
if grep -q '"TMUX"' GentlemanNushell/config.nu; then
  awk -v search='"TMUX"' -v replace='let MULTIPLEXER_ENV_PREFIX = "ZELLIJ"' '
  $0 ~ search {print replace; next}
  {print}
  ' GentlemanNushell/config.nu > GentlemanNushell/config.nu.tmp && mv GentlemanNushell/config.nu.tmp GentlemanNushell/config.nu
else
  echo 'let MULTIPLEXER_ENV_PREFIX = "ZELLIJ"' >> GentlemanNushell/config.nu
fi

# copy files to nushell support directory
cp -rf GentlemanNushell/* ~/Library/Application\ Support/nushell/
```

**_Arch Linux / Linux_**

```bash
if grep -q '"tmux"' ~/.config/nushell/config.nu; then
  awk -v search='"tmux"' -v replace='let MULTIPLEXER = "zellij"' '
  $0 ~ search {print replace; next}
  {print}
  ' ~/.config/nushell/config.nu > ~/.config/nushell/config.nu.tmp && mv ~/.config/nushell/config.nu.tmp ~/.config/nushell/config.nu
else
  echo 'let MULTIPLEXER = "zellij"' >> ~/.config/nushell/config.nu
fi

# update or replace "TMUX"
if grep -q '"TMUX"' ~/.config/nushell/config.nu; then
  awk -v search='"TMUX"' -v replace='let MULTIPLEXER_ENV_PREFIX = "ZELLIJ"' '
  $0 ~ search {print replace; next}
  {print}
  ' ~/.config/nushell/config.nu > ~/.config/nushell/config.nu.tmp && mv ~/.config/nushell/config.nu.tmp ~/.config/nushell/config.nu
else
  echo 'let MULTIPLEXER_ENV_PREFIX = "ZELLIJ"' >> ~/.config/nushell/config.nu
fi
```

#### 6. Install NVIM

```bash
brew install nvim node npm git gcc fzf fd ripgrep coreutils bat curl lazygit tree-sitter
mkdir -p ~/.config/nvim
cp -r GentlemanNvim/nvim/* ~/.config/nvim/
# update or replace /your/notes/path
if grep -q "/your/notes/path" "$HOME/.config/nvim/lua/plugins/obsidian.lua"; then
  awk -v search="/your/notes/path" -v replace="path = '$OBSIDIAN_PATH'" '
  $0 ~ search {print replace; next}
  {print}
  ' "$obsidian_config_file" > "${obsidian_config_file}.tmp" && mv "${obsidian_config_file}.tmp" "$obsidian_config_file"
else
  echo "path = '$OBSIDIAN_PATH'" >> "$obsidian_config_file"
fi
```

- **Manual Obsidian Configuration**

  To set up your Obsidian vault path in Neovim:

  1. **Create the directory for your Obsidian notes** (replace `/path/to/your/notes` with your desired path):

     ```bash
     mkdir -p /path/to/your/notes
     ```

  2. **Create a `templates` folder** inside your notes directory:

     ```bash
     mkdir -p /path/to/your/notes/templates
     ```

  3. **Edit the `obsidian.lua` file** to configure the vault path:

     ```bash
     nvim ~/.config/nvim/lua/plugins/obsidian.lua
     ```

  4. **Update the `path` setting** in `obsidian.lua`:

     ```lua
     path = "/path/to/your/notes",
     ```

  5. **Save and close** the file.

#### 7. Set Default Shell

###### 1. Step

**_ZSH_**

```bash
shell_path=$(which "zsh")
```

**_Fish_**

```bash
shell_path=$(which "fish")
```

**_Nushell_**

```bash
shell_path=$(which "nu")
```

###### 2. Step => Execute to Replace Default Shell

```bash
if [ -n "$shell_path" ]; then
  # Add shell to /etc/shells if not already present
  sudo sh -c "grep -Fxq \"$shell_path\" /etc/shells || echo \"$shell_path\" >> /etc/shells"

  # Change the default shell for the user
  sudo chsh -s "$shell_path" "$USER"

  # Verify if the shell has been changed
  if [ "$SHELL" != "$shell_path" ]; then
    echo -e "${RED}Error: Shell did not change. Please check manually.${NC}"
    echo -e "${GREEN}Command: sudo chsh -s $shell_path \$USER ${NC}"
  else
    echo -e "${GREEN}Shell changed to $shell_path successfully.${NC}"
  fi
else
  echo -e "${RED}Shell $shell_choice not found.${NC}"
fi

# Execute the chosen shell
exec $shell_choice
```

#### 8. Restart the Shell or Computer

- **Close and reopen your terminal**, or **restart your computer** or **WSL instance** for the changes to take effect.

---

You're done! You have manually configured your development environment following the Gentleman.Dots guide. Enjoy your new setup!

**Note:** If you encounter any problems during configuration, consult the official documentation of the tools or seek help online.

**Happy coding!**
</details>

---

## AI Configuration for Neovim

This configuration includes several AI assistants integrated with Neovim. By default, **OpenCode is enabled** as the primary AI assistant with the custom Gentleman agent, while all other AI plugins are disabled.

### Available AI Assistants

The configuration includes support for the following AI tools:

| Plugin | Description | Status |
|--------|-------------|--------|
| **OpenCode.nvim** | OpenCode AI integration with Gentleman agent | ✅ Enabled by default |
| **Avante.nvim** | AI-powered coding assistant | Disabled |
| **CopilotChat.nvim** | GitHub Copilot chat interface | Disabled |
| **CodeCompanion.nvim** | Multi-AI provider support | Disabled |
| **Claude Code.nvim** | Claude AI integration | Disabled |
| **Gemini.nvim** | Google Gemini integration | Disabled |

### How to Switch AI Plugins

**OpenCode.nvim is enabled by default** with the Gentleman agent configuration. All plugin states are managed in a single file:

```bash
nvim ~/.config/nvim/lua/plugins/disabled.lua
```

This file shows all available AI plugins with their current state. To switch:

1. **Find the plugin you want to disable** and set `enabled = false`
2. **Find the plugin you want to enable** and set `enabled = true`
3. **Save and restart Neovim**

Example - switching from OpenCode to Claude Code:

```lua
{
  "NickvanDyke/opencode.nvim",
  enabled = false,  -- Disable OpenCode
},
{
  "coder/claudecode.nvim",
  enabled = true,   -- Enable Claude Code
},
```

### Important Notes

- **Only enable ONE AI plugin at a time** to avoid conflicts and keybinding issues
- **Required CLI tools** are automatically installed by the script:
  - OpenCode CLI (`curl -fsSL https://opencode.ai/install | bash`)
  - Claude Code CLI (`curl -fsSL https://claude.ai/install.sh | bash`)
  - Gemini CLI (`brew install gemini-cli`)
- **API keys may be required** for some services - check each plugin's documentation
- **Node.js 18+** is required for most AI plugins (automatically handled by the configuration)

### Optional: mgrep - AI-Powered Semantic Search

For enhanced code search capabilities using AI, you can install mgrep:

```bash
# Install globally (choose your package manager)
npm install -g @mixedbread/mgrep    # or pnpm / bun

# Sign in once (required for first use)
mgrep login
```

Features:
- Semantic code search powered by AI
- Understands code context and meaning
- Fast and intelligent search results
- Perfect for large codebases

### Recommended AI Assistants

- **For the full Gentleman experience:** Use **OpenCode.nvim** (default) - comes with the custom Gentleman agent that verifies, challenges, and proposes alternatives
- **For Claude users:** Use **Claude Code.nvim** with the Claude Code CLI
- **For GitHub Copilot users:** Use **CopilotChat.nvim**
- **For multi-provider flexibility:** Use **CodeCompanion.nvim** - supports multiple AI providers
- **For Google Gemini users:** Use **Gemini.nvim** with the Gemini CLI

---

## OpenCode Configuration

OpenCode is installed automatically with a custom **Gentleman** agent and theme. This agent is designed to be a collaborative partner, not a yes-man.

### Gentleman Agent Philosophy

The Gentleman agent is a Senior Architect persona with 15+ years of experience. Key characteristics:

- **Never a Yes-Man**: Won't say "you're right" without verifying first. Uses "let's check that" or "dejame verificar eso" instead.
- **Collaborative Partner**: Like Jarvis to Tony Stark - provides data, alternatives, and pushes back when needed.
- **Proposes Alternatives**: Always presents options with tradeoffs: "Option A does X, Option B does Y..."
- **Verifies Before Agreeing**: When you challenge a suggestion, it investigates first using available tools.
- **Bilingual**: Responds in Rioplatense Spanish (Argentina/Uruguay slang) if you write in Spanish, or confrontational English if you write in English.

### Using the Gentleman Agent

To use the Gentleman agent in OpenCode:

1. Open OpenCode in your terminal:
   ```bash
   opencode
   ```

2. Type `/agent` and press Enter

3. Select **gentleman** from the list

That's it! The agent will now respond with the Gentleman personality.

### Default Model

The current default model is determined by OpenCode's configuration. To change it, edit your `opencode.json`:

```bash
nvim ~/.config/opencode/opencode.json
```

Add or modify the `model` field:

```json
{
  "$schema": "https://opencode.ai/config.json",
  "theme": "gentleman",
  "model": "anthropic/claude-sonnet-4-20250514",
  "agent": {
    "gentleman": {
      ...
    }
  }
}
```

#### Available Models

Some popular model options:

| Provider | Model ID |
|----------|----------|
| Anthropic | `anthropic/claude-sonnet-4-20250514` |
| Anthropic | `anthropic/claude-haiku-4-20250514` |
| OpenAI | `openai/gpt-4o` |
| OpenAI | `openai/gpt-4o-mini` |
| Google | `google/gemini-2.0-flash` |
| Google | `google/gemini-2.5-pro-preview-06-05` |

You can also set a specific model per agent:

```json
{
  "agent": {
    "gentleman": {
      "model": "anthropic/claude-sonnet-4-20250514",
      ...
    }
  }
}
```

### OpenCode Theme

The configuration includes a custom **Gentleman** theme with a dark background and Kanagawa-inspired colors. The theme is automatically applied when you run OpenCode.

### MCP Integrations

The Gentleman OpenCode config includes two MCP (Model Context Protocol) integrations:

- **Context7**: Remote MCP for fetching up-to-date documentation
- **mgrep**: Local semantic search tool for your codebase

These are enabled by default and enhance the agent's ability to verify information and search your code.
