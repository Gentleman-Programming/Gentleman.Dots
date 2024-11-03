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

#### 3. Installing the Iosevka Font

The Iosevka Term Nerd Font is required for terminal emulators in this setup. On Windows, this installation must be done manually.

1. Download the Iosevka font from its official site or from [Nerd Fonts GitHub](https://github.com/ryanoasis/nerd-fonts).
2. Extract the archive and locate the font files (`.ttf` or `.otf`).
3. Right-click each font file and select "Install for all users" to install the font system-wide.

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
mkdir %userprofile%\AppData\Roaming\alacritty
cp alacritty.toml %userprofile%\AppData\Roaming\alacritty\alacritty.toml

# Uncomment in alacritty.toml

#[shell]
#program = "wsl.exe"
#args = ["--cd","~"]
```

**WezTerm Configuration**

```powershell
cp .wezterm.lua %userprofile%

# Uncomment for Windows settings

# config.default_domain = 'WSL:Ubuntu'

# config.front_end = "WebGpu"
# config.max_fps = 120
```

**If WezTerm doesn't take the initial configuration**

- Create a `wezterm` folder in `C:\User\your-username\.config`
- Copy `.wezterm.lua` into a `wezterm.lua` file inside the previous directory
- Start WezTerm again

**Kitty Configuration**

```powershell
cp -r GentlemanKitty/* %userprofile%\.config\kitty
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

### For Linux, Arch Linux, and WSL

#### Prerequisites

- **macOS or Linux:** Ensure you have one of these operating systems.
- **Administrator privileges (sudo):** You'll need administrator access to install some tools.

#### Initial Dependencies Installation

Before starting the configuration, you need to install some basic dependencies. Run the commands corresponding to your operating system:

##### Arch Linux

```bash
sudo pacman -Syu --noconfirm
sudo pacman -S --needed --noconfirm base-devel curl file git
```

##### Other Linux Systems

```bash
sudo apt-get update
sudo apt-get install -y build-essential curl file git
```

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
  brew install rsteube/carapace/carapace-bin
  ```

- **Zoxide**: A smarter `cd` command that learns your habits.

  ```bash
  brew install zoxide
  ```

- **Atuin**: A magical shell history.

  ```bash
  bash <(curl https://raw.githubusercontent.com/atuinsh/atuin/main/install.sh)
  ```

Additionally, make sure to create the necessary directories if they do not already exist:

```bash
mkdir -p ~/.cache/starship
mkdir -p ~/.cache/carapace
mkdir -p ~/.local/share/atuin
```

**Note:** Configuration steps for these tools are included in the repository files and will be applied automatically when you copy the configuration files.

#### 3. Install Homebrew (if not installed)

- Homebrew is a package manager for macOS and Linux that makes it easy to install many tools.
- **If you already have Homebrew installed, you can skip this step.**
- Run the following command in your terminal to install Homebrew:

```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

- Follow the on-screen instructions to complete the installation.

#### 4. Install Iosevka Term Nerd Font

Iosevka Term Nerd Font is required for the terminal emulator. Follow these steps to install it:

```bash
brew tap homebrew/cask-fonts
brew install --cask font-iosevka-term-nerd-font
```

#### 5. Choose and Install a Terminal Emulator (Optional)

- **If you are using WSL, terminal emulators must be installed on Windows.**
- Choose one of the following terminal emulators and install it by following the instructions:

  - **Alacritty:**

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

  - **WezTerm:**

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

  - **Kitty (macOS only):**

    ```bash
    brew install --cask kitty

    # Configuration
    mkdir -p ~/.config/kitty
    cp -r GentlemanKitty/* ~/.config/kitty
    ```

  - **None:** If you do not want to install a terminal emulator, you can skip this step.

#### 6. Choose and Install a Shell

- Choose one of the following shells and install it by following the instructions:

  - **Fish:**

    ```bash
    brew install fish

    # Configuration
    cp -r GentlemanFish/fish ~/.config
    ```

  - **Zsh:**

    ```bash
    brew install zsh zsh-autosuggestions zsh-syntax-highlighting zsh-autocomplete

    # Install Oh My Zsh (if not installed)
    NO_INTERACTIVE=true sh -c "$(curl -fsSL https://raw.githubusercontent.com/subtlepseudonym/oh-my-zsh/feature/install-noninteractive/tools/install.sh)"

    # Configuration
    cp -r GentlemanZsh/.zshrc ~/

    # PowerLevel10K configuration
    brew install powerlevel10k
    cp -r GentlemanZsh/.p10k.zsh ~/
    ```

  - **Nushell:**

    ```bash
    # Install Nushell
    brew install nushell

    # Configuration

    // Linux
    mkdir -p ~/.config/nushell
    cp -r GentlemanNushell/* ~/.config/nushell/

    // Mac
    mkdir -p ~/Library/Application\ Support/nushell
    cp -r GentlemanNushell/* ~/Library/Application\ Support/nushell/
    ```

  - **None:** If you do not want to install a shell, you can skip this step.

#### 7. Install Additional Dependencies

- **On Linux (not Arch):**

  ```bash
  sudo apt-get update && sudo apt-get upgrade -y
  ```

#### 8. Choose and Install a Window Manager (Optional)

- Choose one of the following window managers and install it by following the instructions:

  - **Tmux:**

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

  - **Zellij:**

    ```bash
    brew install zellij

    # Configuration
    mkdir -p ~/.config/zellij
    cp -r GentlemanZellij/zellij/* ~/.config/zellij/
    ```

  - **None:** If you do not want to install a window manager, you can skip this step.

    - If no window manager is chosen, you may need to adjust your shell configuration files to prevent automatic launching of a window manager.

#### 9. Configure Neovim

- To install and configure Neovim, follow these steps:

  ```bash
  brew install nvim node npm git gcc fzf fd ripgrep coreutils bat curl lazygit

  # Neovim configuration
  mkdir -p ~/.config/nvim
  cp -r GentlemanNvim/nvim/* ~/.config/nvim/
  ```

- **Manual Obsidian Configuration**

  To set up your Obsidian vault path in Neovim, follow these steps:

  1. **Create the directory for your Obsidian notes** (replace `/path/to/your/notes` with the actual path you’d like to use):

     ```bash
     mkdir -p /path/to/your/notes
     ```

  2. **Inside the notes directory, create a `templates` folder** to store any templates you plan to use:

     ```bash
     mkdir -p /path/to/your/notes/templates
     ```

  3. **Open the `obsidian.lua` file** to configure the vault path:

     ```bash
     nano ~/.config/nvim/lua/plugins/obsidian.lua
     ```

     Alternatively, use Neovim:

     ```bash
     nvim ~/.config/nvim/lua/plugins/obsidian.lua
     ```

  4. **Find the line** that contains:

     ```lua
     path = "/your/notes/path",
     ```

  5. **Replace `"/your/notes/path"`** with the path you created for your notes, for example:

     ```lua
     path = "/path/to/your/notes",
     ```

  6. **Save and close** the file:
     - In `nano`, press `Ctrl + O` to save and `Ctrl + X` to exit.
     - In `nvim`, type `:wq` and press `Enter` to save and close.

#### 10. Set the Default Shell

Before setting your default shell, make sure it’s registered in the `/etc/shells` file. This allows the system to recognize it as a valid default shell option.

1. Check if your chosen shell (e.g., `fish`, `zsh`, or `nu`) is already registered in `/etc/shells` by running:

   ```bash
   cat /etc/shells
   ```

2. If your shell isn’t listed, add it manually. For example, to add `fish`, use:

   ```bash
   echo "$(which fish)" | sudo tee -a /etc/shells
   ```

   For `nushell`, you might need to add:

   ```bash
   echo "$(which nu)" | sudo tee -a /etc/shells
   ```

3. Now, set the default shell by replacing `<shell-name>` with the name of your shell (like `fish`, `zsh`, or `nu`):

   ```bash
   chsh -s $(which <shell-name>)
   ```

#### 11. Restart the Shell or Computer

- Close and reopen your terminal, or restart your computer or WSL instance for the changes to take effect.

You're done! You have manually configured your development environment following the Gentleman.Dots script. Enjoy your new environment!

**Note:** If you encounter any problems during configuration, consult the official documentation of the tools or seek help online.

**I hope this guide has been helpful!**

If you have any further questions or need additional assistance, feel free to reach out to the Gentleman.Dots community or consult the project's documentation.

Happy coding!
