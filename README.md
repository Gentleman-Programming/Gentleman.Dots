# Gentleman.Dots

## Description

This repository contains customized configurations for the Neovim development environment, including specific plugins and keymaps to enhance productivity. It makes use of [LazyVim](https://github.com/LazyVim/LazyVim) as a preconfigured set of plugins and settings to simplify the use of Neovim.

## Installation Steps

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

   After installing and configuring your terminal to run WSL, you can execute the rest of the commands from the WSL terminal.

3. **Configuration Transfer for Terminal Emulators**

   **Alacritty Configuration**

   ```powershell
   git clone https://github.com/Gentleman-Programming/Gentleman.Dots
   cp Gentleman.Dots/alacritty.toml %userprofile%\.config\alacritty\alacritty.toml
   ```

   **WezTerm Configuration**

   ```powershell
   git clone https://github.com/Gentleman-Programming/Gentleman.Dots
   cp Gentleman.Dots/.wezterm.lua %userprofile%

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
   git clone https://github.com/Gentleman-Programming/Gentleman.Dots
   cp -r Gentleman.Dots/GentlemanKitty/* %userprofile%\.config\kitty
   ```

### For macOS/Linux

1. **Install a Terminal Emulator**

   You can choose between Kitty, WezTerm, or Alacritty as your terminal emulator. This repository provides configurations for all three, but it's recommended to use Alacritty.

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
   git clone https://github.com/Gentleman-Programming/Gentleman.Dots
   cp Gentleman.Dots/alacritty.toml ~/.config/alacritty/alacritty.toml
   ```

   **WezTerm Configuration**

   ```bash
   git clone https://github.com/Gentleman-Programming/Gentleman.Dots
   cp Gentleman.Dots/.wezterm.lua ~/.config/wezterm/wezterm.lua
   ```

   **Kitty Configuration**

   ```bash
   git clone https://github.com/Gentleman-Programming/Gentleman.Dots
   cp -r Gentleman.Dots/GentlemanKitty/* ~/.config/kitty
   ```

### Shared Steps (for macOS, Linux, or WSL)

3. **Install HomeBrew**

   ```bash
   /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
   ```

4. **Include HomeBrew Path**

   ```bash
   Change 'YourUserName' with the device username

   (echo; echo 'eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"') >> /home/YourUserName/.bashrc
   eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"
   ```

5. **Install build-essentials for LINUX**

   ```bash
   sudo apt-get update
   sudo apt-get upgrade
   sudo apt-get install build-essential
   ```

6. **Install Starship**

   Starship is a cross-shell prompt that is fast, customizable, and easy to set up.

   **Install Starship**

   ```bash
   brew install starship
   ```

   **Configure Starship**

   Add the following line to your `~/.config/fish/config.fish` for Fish shell:

   ```bash
   starship init fish | source
   ```

   For other shells, refer to the [Starship installation guide](https://starship.rs/guide/#%F0%9F%9A%80-installation).

7. **Install NVIM**

   ```bash
   brew install nvim
   ```

8. **Install NODE & NPM**

   ```bash
   brew install node
   brew install npm
   ```

9. **Install GIT**

   ```bash
   brew install git
   ```

10. **Install FISH**

    ```bash
    brew install fish

    // set as default:

    which fish
    // this will return a path, let‘s call it whichFishResultingPath

    // add it as an available shell
    echo whichFishResultingPath | sudo tee -a /etc/shells

    // set it as default
    sudo chsh -s whichFishResultingPath
    ```

11. **Install Fisher**

    ```bash
    curl -sL https://git.io/fisher | source && fisher install jorgebucaran/fisher
    ```

12. **Install the following dependencies**

    ```bash
    brew install gcc
    brew install fzf
    brew install fd
    brew install ripgrep
    brew install coreutils
    ```

13. **Install Zellij**

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

14. **Install Iosevka Term Nerd Font**

    Download and install the Iosevka Term Nerd Font from [this link](https://github.com/ryanoasis/nerd-fonts/releases/download/v3.1.1/IosevkaTerm.zip).

### Configuration Transfer

#### Neovim Configuration

```bash
git clone https://github.com/Gentleman-Programming/Gentleman.Dots
cp -r Gentleman.Dots/GentlemanNvim/nvim ~/.config
```

Restart Neovim to apply the changes.

#### Tmux Configuration

```bash
git clone https://github.com/Gentleman-Programming/Gentleman.Dots
cp -r Gentleman.Dots/GentlemanTmux/.tmux ~/
cp Gentleman.Dots/GentlemanTmux/.tmux.conf ~/
```

Start Tmux and load the configuration:

```bash
tmux
tmux source-file ~/.tmux.conf
```

#### Zellij Configuration

```bash
git clone https://github.com/Gentleman-Programming/Gentleman.Dots
cp -r Gentleman.Dots/GentlemanZellij/zellij ~/.config
```

#### Starship Configuration

```bash
git clone https://github.com/Gentleman-Programming/Gentleman.Dots
cp Gentleman.Dots/starship.toml ~/.config
```

### Additional Configurations

#### Fish Configuration

```bash
git clone https://github.com/Gentleman-Programming/Gentleman.Dots
cp -r Gentleman.Dots/GentlemanFish/* ~/.config
```

Run `fisher install oh-my-fish/plugin-pj` to install the PJ plugin. Then go to the file `~/.config/fish/fish_variables` and change the following variable to the path to your working folder with your projects:

```bash
SETUVAR --export PROJECT_PATHS: /YourWorkingPath
```

### Note on Terminal Emulators

You can choose between Kitty, WezTerm, or Alacritty as your terminal emulator. This repository provides configurations for all three, but it's recommended to use Alacritty as it is preferred here.
