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
      - [Dependencies Install](#dependencies-install)
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

   This command will install WSL and set version 2 as the default.

1.2. **Install a Linux Distribution**

   Once WSL is installed, you need to install a Linux distribution. Common distributions include Ubuntu, Debian, and Kali Linux, among others.

   To install a Linux distribution (like Ubuntu), follow these steps:

   ```bash
   wsl --install -d Ubuntu
   ```

   This will install the latest version of Ubuntu on your system. You can replace "Ubuntu" with another available distribution if you prefer a different one (e.g., Debian, Kali-Linux, etc.).

   To see the list of available distributions, run:

   ```bash
   wsl --list --online
   ```

   Then, install the distribution of your choice with:

   ```bash
   wsl --install -d <distribution-name>
   ```

1.3. **Launch and Configure the Distribution**

   After the distribution is installed, you need to open it to complete the initial setup. To open the newly installed distribution, you can click on the shortcut created in the Start menu or run:

   ```bash
   wsl
   ```

   This will open your Linux distribution in a terminal window. The first time you open it, you'll be asked to set up a UNIX username and password for your Linux environment.

1.4. **Update the Distribution**

   After configuring the distribution, it's recommended to update the packages to the latest version:

   ```bash
   sudo apt-get update
   sudo apt-get upgrade
   ```

   This ensures your environment is up to date and ready to use.

1.5. **Set WSL Default Distribution**

   If you have installed multiple distributions, you can set which one will run by default when you start WSL without specifying a distribution:

   ```bash
   wsl --set-default Ubuntu
   ```

   This command sets Ubuntu as the default distribution. Replace "Ubuntu" with the name of your preferred distribution.

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
    brew install zsh-autosuggestions zsh-syntax-highlighting zsh-autocomplete
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

7. **Set Zsh as the Default Shell**

    ```bash
    sudo sh -c "echo $(which zsh) >> /etc/shells"
    sudo chsh -s $(which zsh)
    ```

8. **Apply Zsh Configuration**

    To apply the configuration, reload your `.zshrc` file:

    ```bash
    source ~/.zshrc
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

6. **Start Tmux by default**
   
For fish, go to ~/.config/fish/config.fish:
```bash
Uncomment Tmux Code
# Run TMUX
 if status is-interactive
     and not set -q TMUX
     exec tmux
 end

Comment Zellij Code
# Run Zellij
#if set -q ZELLIJ
#else
#   zellij
#end
```
For zsh, go to ~/.zshrc:
```bash
Uncomment Tmux Code
# Run Tmux
if [[ $- == *i* ]] && [[ -z "$TMUX" ]]; then
    exec tmux
fi

Comment Zellij Code
# Run Zellij
#if [[ $- == *i* ]] && [[ -z "$ZELLIJ" ]]; then
#    exec zellij
#fi
```


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

4. **Start Zellij by default**

For fish, go to ~/.config/fish/config.fish:
```bash
Comment Tmux Code
# Run TMUX
# if status is-interactive
#     and not set -q TMUX
#     exec tmux
# end

Uncomment Zellij Code
# Run Zellij
if set -q ZELLIJ
else
    zellij
end
```
For zsh, go to ~/.zshrc:
```bash
Comment Tmux Code
# Run Tmux
#if [[ $- == *i* ]] && [[ -z "$TMUX" ]]; then
#    exec tmux
#fi

Uncomment Zellij Code
# Run Zellij
if [[ $- == *i* ]] && [[ -z "$ZELLIJ" ]]; then
    exec zellij
fi
```

#### Starship Configuration

While in the `Gentleman.Dots` directory, copy the starship configuration files:

```bash
cp starship.toml ~/.config
```

### Note on Terminal Emulators

You can choose between Kitty, WezTerm, or Alacritty as your terminal emulator. This repository provides configurations for all three

, but it is recommended to use Alacritty as it is preferred here.
