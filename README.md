# Gentleman.Dots

<!--toc:start-->
- [Gentleman.Dots](#gentlemandots)
  - [Description](#description)
  - [Installation Steps](#installation-steps)
    - [Step 1: Clone the Repository](#step-1-clone-the-repository)
    - [For Windows](#for-windows)
    - [For macOS/Linux](#for-macoslinux)
    - [Shared Steps (for macOS, Linux, or WSL)](#shared-steps-for-macos-linux-or-wsl)
    - [Step 3: Shell Configuration (Fish and Zsh)](#step-3-shell-configuration-fish-and-zsh)
      - [Fish Configuration](#fish-configuration)
      - [Zsh Configuration](#zsh-configuration)
    - [Step 4: Additional Configurations](#step-4-additional-configurations)
      - [Neovim Configuration](#neovim-configuration)
      - [Tmux Configuration](#tmux-configuration)
      - [Zellij Configuration](#zellij-configuration)
      - [Starship Configuration](#starship-configuration)
    - [Note on Terminal Emulators](#note-on-terminal-emulators)
<!--toc:end-->

## Description

This repository contains customized configurations for the Neovim development environment, including specific plugins and keymaps to enhance productivity. It also includes configurations for both `fish` and `zsh` shells, allowing you to choose according to your preference. We utilize [LazyVim](https://github.com/LazyVim/LazyVim) as a preconfigured set of plugins and settings to simplify the use of Neovim.

## Installation Steps

### Step 1: Clone the Repository

Before proceeding with the configuration transfers, clone this repository and navigate into the cloned directory:

```bash
git clone https://github.com/Gentleman-Programming/Gentleman.Dots.git
cd Gentleman.Dots
```

All subsequent commands assume you are in the `Gentleman.Dots` directory.

### For Windows

1. **Install WSL**

   Windows Subsystem for Linux (WSL) is a compatibility layer for running Linux binary executables natively on Windows 10 and Windows Server 2019. This allows you to use a Linux environment on your Windows machine without the overhead of a virtual machine.

   To install WSL, follow these steps:

   ```bash
   wsl --install
   wsl --set-default-version 2
   ```

   For more detailed instructions, visit the [WSL installation guide](https://learn.microsoft.com/en-us/windows/wsl/install).

2. **Install a Terminal Emulator**

   You can choose between Kitty, WezTerm, or Alacritty as your terminal emulator. This repository provides configurations for all three, but it's recommended to use Alacritty.

   **Install Alacritty**

   1. Download the latest Alacritty release from the [Alacritty GitHub Releases page](https://github.com/alacritty/alacritty/releases).
   2. Extract the downloaded file and move the `alacritty.exe` to a folder in your PATH.

   **Install WezTerm**

   Download and install WezTerm from [this link](https://wezfurlong.org/wezterm/installation.html).

   **Install Kitty**

   Download and install Kitty from [this link](https://sw.kovidgoyal.net/kitty/#get-the-app).

3. **Configuration Transfer for Terminal Emulators**

   **Alacritty Configuration**

   ```powershell
   cp alacritty.toml %userprofile%\.config\alacritty\alacritty.toml
   ```

   **WezTerm Configuration**

   ```powershell
   cp .wezterm.lua %userprofile%

   Uncomment the lines under -- activate ONLY if windows --

   -- config.default_domain = 'WSL:Ubuntu'
   -- config.front_end = "WebGpu"
   -- config.max_fps = 120
   -- for _, gpu in ipairs(wezterm.gui.enumerate_gpus()) do
   -- if gpu.backend == "Vulkan" then
   --   config.webgpu_preferred_adapter = gpu
   --   break
   --  end
   -- end
   ```

   **Kitty Configuration**

   ```powershell
   cp -r GentlemanKitty/* %userprofile%\.config\kitty
   ```

### For macOS/Linux

1. **Install a Terminal Emulator**

   You can choose between Kitty, WezTerm, or Alacritty as your terminal emulator. This repository provides configurations for all three, but it is recommended to use Alacritty.

   **Install Alacritty**

   ```bash
   brew install --cask alacritty
   ```

   **Install WezTerm**

   Download and install WezTerm from [this link](https://wezfurlong.org/wezterm/installation.html).

   **Install Kitty**

   ```bash
   brew install --cask kitty
   ```

2. **Configuration Transfer for Terminal Emulators**

   **Alacritty Configuration**

   ```bash
   cp alacritty.toml ~/.config/alacritty/alacritty.toml
   ```

   **WezTerm Configuration**

   ```bash
   cp .wezterm.lua ~/.config/wezterm/wezterm.lua
   ```

   **Kitty Configuration**

   ```bash
   cp -r GentlemanKitty/* ~/.config/kitty
   ```

### Shared Steps (for macOS, Linux, or WSL)

### Step 3: Shell Configuration (Fish and Zsh)

Depending on your preference, you can configure either `fish` or `zsh` as your default shell.

#### Fish Configuration

1. **Install Homebrew (if not installed)**

   Install Homebrew by running the following command:

   ```bash
   /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
   ```

2. **Install Fish**

    ```bash
    brew install fish
    ```

3. **Set Fish as the Default Shell**

    ```bash
    sudo sh -c "echo $(which fish) >> /etc/shells"
    sudo chsh -s $(which fish)
    ```

4. **Install Fisher**

    ```bash
    curl -sL https://git.io/fisher | source && fisher install jorgebucaran/fisher
    ```

5. **Install PJ Plugin**

    ```bash
    fisher install oh-my-fish/plugin-pj
    ```

6. **Copy Fish Configuration**

    While in the `Gentleman.Dots` directory, copy the Fish configuration files:

    ```bash
    cp -r GentlemanFish/* ~/.config
    ```

7. **Set Project Paths**

    Modify the `PROJECT_PATHS` variable in `~/.config/fish/config.fish` to point to the directory where you store your projects. The default is:

    ```fish
    set PROJECT_PATHS /your/work/path/
    ```

    Replace `/your/work/path/` with the path to your preferred projects directory.

#### Zsh Configuration

1. **Install Homebrew (if not installed)**

   Install Homebrew by running the following command:

   ```bash
   /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
   ```

2. **Install Zsh**

    ```bash
    brew install zsh
    ```

3. **Install Oh My Zsh**

    ```bash
    sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"
    ```

4. **Install Required Plugins with Brew**

    Install the necessary plugins using Homebrew:

    ```bash
    brew install zsh-autosuggestions zsh-syntax-highlighting
    ```

5. **Copy Zsh Configuration**

    While in the `Gentleman.Dots` directory, copy the Zsh configuration file:

    ```bash
    cp .zshrc ~/
    ```

6. **Set Project Paths**

    Modify the `PROJECT_PATHS` variable in `~/.zshrc` to point to the directory where you store your projects. The default is:

    ```bash
    export PROJECT_PATHS="/your/work/path/"
    ```

    Replace `/your/work/path/` with the path to your preferred projects directory.

7. **Apply Zsh Configuration**

    To apply the configuration, reload your `.zshrc` file:

    ```bash
    source ~/.zshrc
    ```

8. **Additional Plugin Configuration**

    Ensure the following lines are in your `.zshrc`:

    ```bash
    source /opt/homebrew/share/zsh-autocomplete/zsh-autocomplete.plugin.zsh
    source /opt/homebrew/share/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh
    ```

### Step 4: Additional Configurations

#### Dependencies Install

1. **Install build-essentials for LINUX** (for Linux and WSL)

   ```bash
   sudo apt-get update
   sudo apt-get upgrade
   sudo apt-get install build-essential
   ```
2. **Install Starship**

   Starship is a cross-shell prompt that is fast, customizable, and easy to set up.

   **Install Starship**

   ```bash
   brew install starship
   ```
3. **Install NVIM**

   ```bash
   brew install nvim
   ```
4. **Install NODE & NPM**

   ```bash
   brew install node
   brew install npm
   ```
5. **Install GIT**

   ```bash
   brew install git
   ```
6. **Install the following dependencies**

    ```bash
    brew install gcc
    brew install fzf
    brew install fd
    brew install ripgrep
    brew install coreutils
    ```
7. **Install Iosevka Term Nerd Font**

    Download and install the Iosevka Term Nerd Font from [this link](https://github.com/ryanoasis/nerd-fonts/releases/download/v3.1.1/IosevkaTerm.zip).

#### Neovim Configuration

```bash
cp -r GentlemanNvim/nvim ~/.config
```

Restart Neovim to apply the changes.

#### Tmux Configuration

1. **Install Tmux**

   Tmux is a terminal multiplexer that allows you to run multiple terminal sessions within a single window.

   **Install Tmux**

   ```bash
   brew install tmux
   ```

2. **Install TPM (Tmux Plugin Manager)**

   TPM is a plugin manager for Tmux that allows you to easily manage and install Tmux plugins.

   **Install TPM**

   ```bash
   git clone https://github.com/tmux-plugins/tpm ~/.tmux/plugins/tpm
   ```

3. **Copy Tmux Configuration**

   While in the `Gentleman.Dots` directory, copy the Tmux configuration files:

   ```bash
   cp -r GentlemanTmux/.tmux ~/
   cp GentlemanTmux/.tmux.conf ~/
   ```

4. **Start Tmux and Load Configuration**

   ```bash
   tmux
   tmux source-file ~/.tmux.conf
   ```

5. **Install Tmux Plugins**

   Inside a Tmux session, press `Ctrl + a and then I` (capital I, as in Install) to fetch the plugins defined in your `.tmux.conf` file.

#### Zellij Configuration

1. **Install Zellij**

    ```bash
    brew install zellij
    ```

    If you find any issues with this method, use "Cargo" to install Zellij:

    ```bash
    // if installed with brew:
    brew uninstall zellij

    // Install Rust (needed for Cargo)
    curl https://sh.rustup.rs -sSf | sh

    // Install Zellij using cargo
    cargo install --locked zellij
    ```
2. **Copy Zellij Configuration**

While in the `Gentleman.Dots` directory, copy the Zellij configuration files:

```bash
cp -r GentlemanZellij/zellij ~/.config
```

3. **Choose the default Shell**

Go to ~/.config/zellij/config.kdl:

```bash
// uncomment the shell you want to use
default_shell "fish"
// default_shell "zsh"
``` 

#### Starship Configuration

While in the `Gentleman.Dots` directory, copy the starship configuration files:

```bash
cp starship.toml ~/.config
```

### Note on Terminal Emulators

You can choose between Kitty, WezTerm, or Alacritty as your terminal emulator. This repository provides configurations for all three

, but it is recommended to use Alacritty as it is preferred here.
