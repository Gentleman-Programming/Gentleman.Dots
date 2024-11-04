# Gentleman.Dots

## Description

This repository contains customized configurations for a complete development environment, including Neovim, Fish, Zsh, **Nushell**, Tmux, Zellij, and terminal emulators like Alacritty, WezTerm, and Kitty. You can choose between automatic and manual installation methods depending on your preference and operating system. **Important:** Windows users must follow the manual installation instructions before running the script.

## Installation (Automatic Recommended!)

### The Easy Way! Test the automated process and let the script do all the work for you 😘

The **automatic installation script** is the quickest and easiest way to set up your development environment. This script handles all the heavy lifting, but remember that you **must install the font** mentioned below before running it. The script is designed for macOS, Linux, and WSL systems. If you’re on Windows, you’ll need to follow the manual steps first before attempting to run this script.

```bash
curl -O https://raw.githubusercontent.com/Gentleman-Programming/Gentleman.Dots/main/install-linux-mac.sh

sudo chmod +x install-linux-mac.sh
bash ./install-linux-mac.sh
```

## Manual Installation

Welcome to the Gentleman.Dots manual configuration guide! This document will walk you through the steps required to set up your development environment.

---

### For Windows

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

---

### For Linux, Arch Linux, macOS, and WSL

#### Prerequisites

- **macOS or Linux:** Ensure you have one of these operating systems.
- **Administrator privileges (sudo):** You'll need administrator access to install some tools.

#### Initial Dependencies Installation (Including Rust and Cargo)

Before starting the configuration, you need to install some essential dependencies, including **Rust** (and Cargo, Rust's package manager) for managing Rust-based tools. Run the commands corresponding to your operating system.

##### Arch Linux

1. **Update the system and install base development tools:**

   ```bash
   sudo pacman -Syu --noconfirm
   sudo pacman -S --needed --noconfirm base-devel curl file git
   ```

2. **Install Rustup (the Rust toolchain installer):**

   ```bash
   sudo pacman -S rustup
   ```

3. **Initialize Rustup and install the stable toolchain:**

   ```bash
   rustup default stable
   ```

4. **Ensure Cargo's bin directory is in your `PATH`:**

   ```bash
   source $HOME/.cargo/env
   ```

5. **Verify the installation:**

   ```bash
   rustc --version
   cargo --version
   ```

##### Other Linux Systems (e.g., Ubuntu)

1. **Update the system and install base development tools:**

   ```bash
   sudo apt-get update
   sudo apt-get install -y build-essential curl file git
   ```

2. **Install Rustup (the Rust toolchain installer):**

   ```bash
   curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
   ```

3. **Follow the on-screen instructions to complete the installation.**

4. **Ensure Cargo's bin directory is in your `PATH`:**

   ```bash
   source $HOME/.cargo/env
   ```

5. **Verify the installation:**

   ```bash
   rustc --version
   cargo --version
   ```

##### macOS

1. **Install Xcode command line tools:**

   ```bash
   xcode-select --install
   ```

2. **Install Homebrew (if not already installed):**

   ```bash
   /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
   ```

3. **Install Rustup (the Rust toolchain installer):**

   ```bash
   curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
   ```

4. **Follow the on-screen instructions to complete the installation.**

5. **Ensure Cargo's bin directory is in your `PATH`:**

   ```bash
   source $HOME/.cargo/env
   ```

6. **Verify the installation:**

   ```bash
   rustc --version
   cargo --version
   ```

---

#### 1. Clone the Repository

- **If the repository is already cloned, you can skip this step.**
- Open your terminal and run the following command to clone the repository:
```bash
git clone https://github.com/Gentleman-Programming/Gentleman.Dots.git
```

- Then, navigate to the cloned directory:

```bash
cd Gentleman.Dots
```

#### 2. Install Shell Enhancements

To enhance your shell experience, we recommend installing the following tools, regardless of which shell you are using:

- **Carapace**: A universal completion generator for your shell.

  ```bash
  cargo install carapace-bin
  ```

- **Zoxide**: A smarter `cd` command that learns your habits.

  ```bash
  cargo install zoxide
  ```

- **Atuin**: A magical shell history.

  ```bash
  bash <(curl https://raw.githubusercontent.com/atuinsh/atuin/main/install.sh)
  ```

Ensure the required directories exist:

```bash
mkdir -p ~/.cache/starship
mkdir -p ~/.cache/carapace
mkdir -p ~/.local/share/atuin
```

**Note:** Configuration steps for these tools are included in the repository files and will be applied automatically when you copy the configuration files.

#### 3. Install Homebrew (if not installed on macOS)

- **If you already have Homebrew installed, you can skip this step.**
- Run the following command in your terminal to install Homebrew:

```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

- Follow the on-screen instructions to complete the installation.

#### 4. Install Iosevka Term Nerd Font

Iosevka Term Nerd Font is required for the terminal emulator.

**On macOS:**

```bash
brew tap homebrew/cask-fonts
brew install --cask font-iosevka-term-nerd-font
```

**On Linux:**

- Download the font from [Nerd Fonts GitHub](https://github.com/ryanoasis/nerd-fonts/releases).
- Extract the font files and move them to your local fonts directory:

  ```bash
  mkdir -p ~/.local/share/fonts
  cp path_to_downloaded_fonts/*.ttf ~/.local/share/fonts/
  fc-cache -fv
  ```

#### 5. Choose and Install a Terminal Emulator (Optional)

- **If you are using WSL, terminal emulators must be installed on Windows.**
- Choose one of the following terminal emulators and install it by following the instructions:

  - **Alacritty:**

    **macOS:**

    ```bash
    brew install --cask alacritty
    ```

    **Linux (Debian/Ubuntu):**

    ```bash
    sudo add-apt-repository ppa:aslatter/ppa
    sudo apt-get update
    sudo apt-get install alacritty
    ```

    **Arch Linux:**

    ```bash
    sudo pacman -S --noconfirm alacritty
    ```

    **Configuration (all systems):**

    ```bash
    mkdir -p ~/.config/alacritty
    cp alacritty.toml ~/.config/alacritty/alacritty.toml
    ```

  - **WezTerm:**

    **macOS:**

    ```bash
    brew install --cask wezterm
    ```

    **Linux (Debian/Ubuntu):**

    ```bash
    curl -LO https://github.com/wez/wezterm/releases/latest/download/WezTerm-debian-ubuntu.latest.deb
    sudo apt install ./WezTerm-debian-ubuntu.latest.deb
    ```

    **Arch Linux:**

    ```bash
    sudo pacman -S --noconfirm wezterm
    ```

    **Configuration (all systems):**

    ```bash
    mkdir -p ~/.config/wezterm
    cp .wezterm.lua ~/.config/wezterm/wezterm.lua
    ```

  - **Kitty:**

    **macOS:**

    ```bash
    brew install --cask kitty
    ```

    **Linux:**

    ```bash
    # Download the installer script and run it
    curl -L https://sw.kovidgoyal.net/kitty/installer.sh | sh /dev/stdin
    ```

    **Configuration:**

    ```bash
    mkdir -p ~/.config/kitty
    cp -r GentlemanKitty/* ~/.config/kitty
    ```

  - **None:** If you do not want to install a terminal emulator, you can skip this step.

#### 6. Choose and Install a Shell

Choose one of the following shells and install it:

- **Fish:**

  **macOS:**

  ```bash
  brew install fish
  ```

  **Linux (Debian/Ubuntu):**

  ```bash
  sudo apt-get install fish
  ```

  **Arch Linux:**

  ```bash
  sudo pacman -S fish
  ```

  **Configuration:**

  ```bash
  mkdir -p ~/.config/fish
  cp -r GentlemanFish/fish/* ~/.config/fish/
  ```

- **Zsh:**

  **macOS:**

  ```bash
  brew install zsh zsh-autosuggestions zsh-syntax-highlighting
  ```

  **Linux (Debian/Ubuntu):**

  ```bash
  sudo apt-get install zsh
  ```

  **Arch Linux:**

  ```bash
  sudo pacman -S zsh zsh-autosuggestions zsh-syntax-highlighting
  ```

  **Install Oh My Zsh:**

  ```bash
  sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"
  ```

  **Configuration:**

  ```bash
  cp -r GentlemanZsh/.zshrc ~/
  ```

  **PowerLevel10K Theme:**

  ```bash
  git clone --depth=1 https://github.com/romkatv/powerlevel10k.git ${ZSH_CUSTOM:-$HOME/.oh-my-zsh/custom}/themes/powerlevel10k
  cp -r GentlemanZsh/.p10k.zsh ~/
  ```

- **Nushell:**

  **Install Nushell using Cargo:**

  ```bash
  cargo install nu
  ```

  **Configuration:**

  - **Linux:**

    ```bash
    mkdir -p ~/.config/nushell
    cp -rf bash-env-json ~/.config/
    cp -rf bash-env.nu ~/.config/
    cp -r GentlemanNushell/* ~/.config/nushell/
    ```

  - **macOS:**

    ```bash
    mkdir -p ~/Library/Application\ Support/nushell
    cp -r GentlemanNushell/* ~/Library/Application\ Support/nushell/
    ```

- **None:** If you do not want to install a shell, you can skip this step.

#### 7. Install Additional Dependencies

- **On Linux (Debian/Ubuntu):**

  ```bash
  sudo apt-get update && sudo apt-get upgrade -y
  sudo apt-get install -y fzf fd-find ripgrep bat exa git gcc curl lazygit jq bash
  ```

- **On macOS:**

  ```bash
  brew install fzf fd ripgrep bat exa git gcc curl lazygit jq bash
  ```

- **On Arch Linux:**

  ```bash
  sudo pacman -S --noconfirm fzf fd ripgrep bat exa git gcc curl lazygit
  ```

#### 8. Choose and Install a Window Manager (Optional)

- **Tmux:**

  **macOS:**

  ```bash
  brew install tmux
  ```

  **Linux (Debian/Ubuntu):**

  ```bash
  sudo apt-get install tmux
  ```

  **Arch Linux:**

  ```bash
  sudo pacman -S tmux
  ```

  **Configuration:**

  ```bash
  git clone https://github.com/tmux-plugins/tpm ~/.tmux/plugins/tpm
  mkdir -p ~/.tmux
  cp -r GentlemanTmux/.tmux/* ~/.tmux/
  cp GentlemanTmux/.tmux.conf ~/
  ```

  **Install Tmux plugins:**

  ```bash
  ~/.tmux/plugins/tpm/bin/install_plugins
  ```

- **Zellij:**

  **Install Zellij using Cargo:**

  ```bash
  cargo install --locked zellij
  ```

  **Configuration:**

  ```bash
  mkdir -p ~/.config/zellij
  cp -r GentlemanZellij/zellij/* ~/.config/zellij/
  ```

- **None:** If you do not want to install a window manager, you can skip this step.

#### 9. Configure Neovim

- **Install Neovim and dependencies:**

  **macOS:**

  ```bash
  brew install neovim node npm git gcc fzf fd ripgrep coreutils bat curl lazygit
  ```

  **Linux (Debian/Ubuntu):**

  ```bash
  sudo apt-get install -y neovim nodejs npm git gcc fzf fd-find ripgrep coreutils bat curl
  ```

  **Arch Linux:**

  ```bash
  sudo pacman -S neovim nodejs npm git gcc fzf fd ripgrep coreutils bat curl lazygit
  ```

- **Neovim configuration:**

  ```bash
  mkdir -p ~/.config/nvim
  cp -r GentlemanNvim/nvim/* ~/.config/nvim/
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

#### 10. Set the Default Shell

Before setting your default shell, ensure it's listed in `/etc/shells`.

1. **Check if your shell is listed:**

   ```bash
   cat /etc/shells
   ```

2. **Add your shell to `/etc/shells` if necessary:**

   - For Fish:

     ```bash
     echo "$(which fish)" | sudo tee -a /etc/shells
     ```

   - For Nushell:

     ```bash
     echo "$(which nu)" | sudo tee -a /etc/shells
     ```

3. **Change your default shell:**

   ```bash
   chsh -s $(which <shell-name>)
   ```

   Replace `<shell-name>` with `fish`, `zsh`, or `nu`.

#### 11. Restart the Shell or Computer

- **Close and reopen your terminal**, or **restart your computer** or **WSL instance** for the changes to take effect.

---

You're done! You have manually configured your development environment following the Gentleman.Dots guide. Enjoy your new setup!

**Note:** If you encounter any problems during configuration, consult the official documentation of the tools or seek help online.

**Happy coding!**
