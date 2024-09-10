# Gentleman.Dots

## Description

This repository contains customized configurations for a complete development environment, including Neovim, Fish, Zsh, Tmux, Zellij, and terminal emulators like Alacritty, WezTerm, and Kitty. You can choose between automatic and manual installation methods depending on your preference and operating system. **Important:** Windows users must follow the manual installation instructions before running the script.

## Installation (Automatic Recommended!)

### The Easy Way! Test the automated process and let the script do all the work for you 😘

The **automatic installation script** is the quickest and easiest way to set up your development environment. This script handles all the heavy lifting, but remember that you **must install the font** mentioned below before running it. The script is designed for macOS, Linux, and WSL systems. If you’re on Windows, you’ll need to follow the manual steps first before attempting to run this script.

```bash
curl -O https://raw.githubusercontent.com/Gentleman-Programming/Gentleman.Dots/main/install-linux-mac.sh

sudo chmod +x install-linux-mac.sh
bash ./install-linux-mac.sh
```

## Manual Installation

### For Windows

**Important:** Windows users must follow these manual installation steps before running the automated script.

#### 1. Install WSL

WSL (Windows Subsystem for Linux) allows you to run Linux on Windows. Install it and set it to version 2:

```bash
wsl --install
wsl --set-default-version 2
```

#### 2. Install a Linux Distribution

Install a Linux distribution (e.g., Ubuntu) in WSL:

```bash
wsl --install -d Ubuntu
```

To list available distributions:

```bash
wsl --list --online
```

Install your preferred distribution:

```bash
wsl --install -d <distribution-name>
```

#### 3. Launch and Configure the Distribution

Open the installed distribution to complete setup. Update it with:

```bash
sudo apt-get update
sudo apt-get upgrade
```

#### 4. Install a Terminal Emulator

Choose and install one of the following terminal emulators:

- **Alacritty**: [Download from GitHub Releases](https://github.com/alacritty/alacritty/releases) and place `alacritty.exe` in your `PATH`.
- **WezTerm**: [Download and Install](https://wezfurlong.org/wezterm/installation.html).
- **Kitty**: [Download and Install](https://sw.kovidgoyal.net/kitty/#get-the-app).

#### 5. Configuration Transfer for Terminal Emulators

Using Powershell:

**Alacritty Configuration**

```powershell
mkdir %userprofile%\AppData\Roaming\alacritty
cp alacritty.toml %userprofile%\AppData\Roaming\alacritty\alacritty.toml

# Uncomment in alacritty.toml

[shell]
program = "wsl.exe"
args = ["--cd","~"]
```

**WezTerm Configuration**

```powershell
cp .wezterm.lua %userprofile%

# Uncomment for Windows settings

# config.default_domain = 'WSL:Ubuntu'

# config.front_end = "WebGpu"
# config.max_fps = 120
```

**Kitty Configuration**

```powershell
cp -r GentlemanKitty/* %userprofile%\.config\kitty
```

---

# Gentleman.Dots - Manual Configuration

Welcome to the Gentleman.Dots manual configuration guide! This document will walk you through the steps required to set up your development environment by following the `install-linux-mac.txt` script.

## Prerequisites

* **macOS or Linux:** Ensure you have one of these operating systems.
* **Administrator privileges (sudo):** You'll need administrator access to install some tools.

## Initial Dependencies Installation

Before starting the configuration, you need to install some basic dependencies. Run the commands corresponding to your operating system:

### Arch Linux

```bash
sudo pacman -Syu --noconfirm
sudo pacman -S --needed --noconfirm base-devel curl file git
```

### Other Linux Systems

```bash
sudo apt-get update
sudo apt-get install -y build-essential curl file git
```

## Configuration Steps

### 1. Clone the Repository

* **If the repository is already cloned, you can skip this step.**
* Open your terminal and run the following command to clone the repository:

```bash
git clone https://github.com/Gentleman-Programming/Gentleman.Dots.git
```

* Then, navigate to the cloned directory:

```bash
cd Gentleman.Dots
```

### 2. Install Homebrew (if not installed)

* Homebrew is a package manager for macOS and Linux that makes it easy to install many tools.
* **If you already have Homebrew installed, you can skip this step.**
* Run the following command in your terminal to install Homebrew:

```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

* Follow the on-screen instructions to complete the installation.

### 3. Choose and Install a Terminal Emulator (Optional)

* **If you are using WSL, terminal emulators must be installed on Windows.**
* Choose one of the following terminal emulators and install it by following the instructions:

   * **Alacritty:**

      ```bash
      # macOS/Linux (not Arch)
      sudo add-apt-repository ppa:aslatter/ppa 
      sudo apt-get update 
      sudo apt-get install alacritty

      # Arch Linux
      sudo pacman -S --noconfirm alacritty

      # Configuration (all systems)
      mkdir -p ~/.config/alacritty 
      cp alacritty.toml ~/.config/alacritty/alacritty.toml
      ```

   * **WezTerm:**

     ```bash
     # macOS/Linux (not Arch)
     curl -fsSL https://apt.fury.io/wez/gpg.key | sudo gpg --yes --dearmor -o /usr/share/keyrings/wezterm-fury.gpg
     echo 'deb [signed-by=/usr/share/keyrings/wezterm-fury.gpg] https://apt.fury.io/wez/ * *' | sudo tee /etc/apt/sources.list.d/wezterm.list
     sudo apt update
     sudo apt install wezterm

     # Arch Linux
     sudo pacman -S --noconfirm wezterm

     # Configuration (all systems)
     mkdir -p ~/.config/wezterm 
     cp .wezterm.lua ~/.config/wezterm/wezterm.lua
     ```

   * **Kitty (macOS only):**

     ```bash
     brew install --cask kitty

     # Configuration
     mkdir -p ~/.config/kitty 
     cp -r GentlemanKitty/* ~/.config/kitty
     ```

   * **None:** If you do not want to install a terminal emulator, you can skip this step.

### 4. Choose and Install a Shell

* Choose one of the following shells and install it by following the instructions:

   * **Fish:**

     ```bash
     brew install fish

     # Configuration
     cp -r GentlemanFish/fish ~/.config

     # (Optional) Update your projects path in ~/.config/fish/config.fish
     set PROJECT_PATHS "/your/work/path/" 
     ```

   * **Zsh:**

     ```bash
     brew install zsh zsh-autosuggestions zsh-syntax-highlighting zsh-autocomplete

     # Install Oh My Zsh (if not installed)
     NO_INTERACTIVE=true sh -c "$(curl -fsSL https://raw.githubusercontent.com/subtlepseudonym/oh-my-zsh/feature/install-noninteractive/tools/install.sh)"

     # Configuration
     cp -r GentlemanZsh/.zshrc ~/

     # PowerLevel10K configuration
     brew install powerlevel10k
     cp -r GentlemanZsh/.p10k.zsh ~/

     # (Optional) Update your projects path in ~/.zshrc
     export PROJECT_PATHS="/your/work/path/" 
     ```

   * **None:** If you do not want to install a shell, you can skip this step.

### 5. Install Additional Dependencies

* **On Linux (not Arch):**

  ```bash
  sudo apt-get update && sudo apt-get upgrade -y
  ```

### 6. Choose and Install a Window Manager (Optional)

* Choose one of the following window managers and install it by following the instructions:

   * **Tmux:**

     ```bash
     brew install tmux

     # Configuration
     git clone https://github.com/tmux-plugins/tpm ~/.tmux/plugins/tpm 
     mkdir -p ~/.tmux
     cp -r GentlemanTmux/.tmux/* ~/.tmux/
     cp GentlemanTmux/.tmux.conf ~/

     # Install Tmux plugins
     tmux new-session -d -s plugin-installation 'source ~/.tmux.conf; tmux run-shell ~/.tmux/plugins/tpm/bin/install_plugins'
     # Wait for the plugin installation to finish. You can monitor the progress in another terminal with:
     # tmux attach-session -t plugin-installation
     tmux kill-session -t plugin-installation
     ```

   * **Zellij:**

     ```bash
     brew install zellij

     # Configuration
     mkdir -p ~/.config/zellij
     cp -r GentlemanZellij/zellij/* ~/.config/zellij/

     # Replace TMUX and tmux in the shell configuration (Zsh)
     sed -i '' 's/TMUX/WM_VAR="\/\$ZELLIJ/g' ~/.zshrc  # If you chose zsh
     sed -i '' 's/tmux/WM_CMD="zellij/g' ~/.zshrc       # If you chose zsh

     # Replace TMUX and tmux in the shell configuration (Fish)
     sed -i '' 's/TMUX/if not set -q ZELLIJ/g' ~/.config/fish/config.fish  # If you chose fish
     sed -i '' 's/tmux/zellij/g' ~/.config/fish/config.fish               # If you chose fish
     ```

   * **None:** If you do not want to install a window manager, you can skip this step.
   
     * If no window manager is chosen, delete the line that runs `tmux` or `zellij`:

     ```bash
     sed -i '' '/exec tmux/d' ~/.zshrc        # If you chose zsh
     sed -i '' '/exec zellij/d' ~/.zshrc      # If you chose zsh
     sed -i '' '/tmux/d' ~/.config/fish/config.fish  # If you chose fish
     sed -i '' '/zellij/d' ~/.config/fish/config.fish # If you chose fish
     ```

### 7. Configure Neovim (Optional)

* If you want to install and configure Neovim, follow these steps:

  ```bash
  brew install nvim node npm git gcc fzf fd ripgrep coreutils bat curl lazygit

  # Neovim configuration
  mkdir -p ~/.config/nvim
  cp -r GentlemanNvim/nvim/* ~/.config/nvim/

  # Obsidian configuration (optional)
  # Replace '/your/notes/path' with the actual path to your Obsidian vault
  sed -i '' 's/\/your\/notes\/path/your/actual/obsidian/vault/path/g' ~/.config/nvim/lua/plugins/obsidian.lua
  ```

### 8. Set the Default Shell

* Run the following command, replacing `<shell-name>` with the name of the shell you chose (fish or zsh):

```bash
chsh -s $(which <shell-name>)
```

### 9. Restart the Shell or Computer

* Close and reopen your terminal, or restart your computer or WSL instance for the changes to take effect.

You're done! You have manually configured your development environment following the Gentleman.Dots script. Enjoy your new environment!

**Note:** If you encounter any problems during configuration, consult the official documentation of the tools or seek help online.

**I hope this guide has been helpful!** 

If you have any further questions or need additional assistance, feel free to reach out to the Gentleman.Dots community or consult the project's documentation. 

Happy coding! 
