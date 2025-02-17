# Gentleman.Dots

![Screenshot](https://github.com/user-attachments/assets/3f6c4f62-23d7-41d7-b7b1-42c7e0c32336)

## Description

This repository contains customized configurations for a complete development environment, including:

- Neovim
- Fish
- Zsh
- **Nushell**
- Tmux
- Zellij
- Terminal emulators:
  - Alacritty
  - **WezTerm** (default)

You can now automatically set up your environment using our new Nix Flake approach with Home Manager. This method is fully declarative and reproducible, and it lets you easily override default options. In our flake, all configurations are defined inline in local modules (e.g., `fish.nix`, `nushell.nix`, etc.), and the flake also:

1. Installs all the required dependencies (git, curl, rustc, cargo, tmux, zellij, neovim, etc.).
2. Automatically performs placeholder replacements in configuration files according to your chosen multiplexer (Zellij vs. Tmux).

---

## Previous Steps

### Installing Nix and Home Manager

Before running the automated installation commands, make sure Nix is installed:

- **macOS and Linux:**

  ```bash
  curl -L https://nixos.org/nix/install | sh
  . ~/.nix-profile/etc/profile.d/nix.sh
  ```

- **WSL (Windows Subsystem for Linux):**

  Open your WSL terminal and run:

  ```bash
  curl -L https://nixos.org/nix/install | sh
  . ~/.nix-profile/etc/profile.d/nix.sh
  ```

Next, install Home Manager (see the official instructions or use our flake method):

```bash
nix-channel --add https://github.com/nix-community/home-manager/archive/master.tar.gz home-manager
nix-channel --update
nix-shell '<home-manager>' -A install
```

---

## Automatic Installation (Recommended for Linux, macOS, and WSL)

### The Easy Way – Let Nix Do the Heavy Lifting

1. **Clone the Repository:**

   ```bash
   git clone https://github.com/Gentleman-Programming/Gentleman.Dots.git
   cd Gentleman.Dots
   ```

2. **Run the Installation Commands:**

   ```bash
   nix profile install .#gentleman-dots
   home-manager switch --flake .#zellij-nushell-starship
   ```

This flake uses local modules (e.g., `fish.nix`, `nushell.nix`, `kitty.nix`, etc.) that contain inline configurations for each tool. It also installs all necessary dependencies based on your options and automatically performs placeholder replacements in configuration files—e.g., if you choose Zellij as your multiplexer, it will update placeholders in your .zshrc, config.fish, and config.nu accordingly. Files are deployed to system-specific locations; for example, the Nushell configuration is copied to:

- **macOS:** `~/Library/Application Support/nushell`
- **Linux/WSL:** `~/.config/nushell`

---

## Available Configurations

The default options are defined in the flake as follows:

```nix
gentlemanOptions = {
  terminal = "wezterm";      # Default terminal is WezTerm
  shell = "nushell";         # Default shell is Nushell
  windowManager = "zellij";  # Default multiplexer is Zellij
  installNeovim = true;
  osType = if system == "x86_64-darwin" then "mac" else "linux";
  starship = true;           # Starship is enabled by default
  powerlevel10k = false;
  useTmux = false;
};
```

You can override these defaults by editing the `gentlemanOptions` block in `flake.nix` or by selecting one of the preset configurations listed below.

### Preset Configurations

To activate one, run the corresponding command:

- **Zellij with Fish and Starship:**

  ```bash
  home-manager switch --flake .#zellij-fish-starship
  ```

- **Zellij with Nushell and Starship:**

  ```bash
  home-manager switch --flake .#zellij-nushell-starship
  ```

- **Zellij with Zsh and Powerlevel10k:**

  ```bash
  home-manager switch --flake .#zellij-zsh-power10k
  ```

- **Tmux with Fish and Starship:**

  ```bash
  home-manager switch --flake .#tmux-fish-starship
  ```

- **Tmux with Nushell and Starship:**

  ```bash
  home-manager switch --flake .#tmux-nushell-starship
  ```

- **Tmux with Zsh and Powerlevel10k:**

  ```bash
  home-manager switch --flake .#tmux-zsh-power10k
  ```

### Overriding the Terminal Emulator

If you want to use a different terminal than the default WezTerm, you can override the `terminal` option. For example, to use Alacritty instead of WezTerm with the Zellij with Nushell and Starship preset, run:

```bash
home-manager switch --flake .#zellij-nushell-starship --override 'gentlemanOptions.terminal="alacritty"'
```

Similarly, for Kitty:

```bash
home-manager switch --flake .#zellij-nushell-starship --override 'gentlemanOptions.terminal="kitty"'
```

_Note:_ These presets are defined in the flake. If you wish to create additional variants or adjust the options, modify the `gentlemanOptions` block or add new homeConfigurations.

---

## Manual Installation for Windows

### Clone the Repository

```bash
git clone git@github.com:Gentleman-Programming/Gentleman.Dots.git
cd Gentleman.Dots
```

#### 1. Install WSL

WSL (Windows Subsystem for Linux) lets you run Linux on Windows. Install and set it to version 2:

```powershell
wsl --install
wsl --set-default-version 2
```

#### 2. Install a Linux Distribution

For example, install Ubuntu:

```powershell
wsl --install -d Ubuntu
```

List available distributions:

```powershell
wsl --list --online
```

Then install your preferred distribution:

```powershell
wsl --install -d <distribution-name>
```

#### 3. Install the Iosevka Term Nerd Font

This font is required by our terminal emulators. Download it from the [Nerd Fonts GitHub](https://github.com/ryanoasis/nerd-fonts) or its official site. Then extract and install the font files (right-click each file and select **"Install for all users"**).

#### 4. Install a Terminal Emulator

Choose and install one of the following:

- **Alacritty:** [Download from GitHub Releases](https://github.com/alacritty/alacritty/releases). Make sure `alacritty.exe` is in your `PATH`.
- **WezTerm:** [Download and Install](https://wezfurlong.org/wezterm/installation.html). Also, set the `HOME` environment variable to point to `C:\Users\your-username`.

#### 5. Transfer Emulator Configurations

Using PowerShell:

**Alacritty:**

```powershell
mkdir $env:APPDATA\alacritty
Copy-Item -Path alacritty.toml -Destination $env:APPDATA\alacritty\alacritty.toml

# In alacritty.toml, uncomment and set:
#[shell]
#program = "wsl.exe"
#args = ["--cd", "~"]
```

**WezTerm:**

```powershell
Copy-Item -Path .wezterm.lua -Destination $HOME
```

_If WezTerm doesn’t pick up the configuration, create a folder `C:\Users\your-username\.config\wezterm` and place `.wezterm.lua` there._

#### 6. Install Chocolatey and win32yank

**Chocolatey** is a Windows package manager.

**Install Chocolatey:**

Open PowerShell as Administrator and run:

```powershell
Set-ExecutionPolicy Bypass -Scope Process -Force; `
[System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; `
iwr https://community.chocolatey.org/install.ps1 -UseBasicParsing | iex
```

**Install win32yank:**

```powershell
choco install win32yank
```

_win32yank is needed for clipboard integration in Neovim when using WSL._

#### 7. Launch and Update Your Linux Distribution

Open your installed Linux distribution (WSL) and run the appropriate update commands:

- **For Ubuntu/Debian:**

  ```bash
  sudo apt-get update
  sudo apt-get upgrade
  ```

- **For Arch Linux:**

  ```bash
  sudo pacman -Syu
  ```

- **For Fedora:**

  ```bash
  sudo dnf upgrade --refresh
  ```

---

## Summary

With the new Nix Flake method, you can automatically install your complete development environment with a single, declarative configuration. Key points:

- **Defaults:**

  - Terminal: WezTerm
  - Shell: Nushell
  - Multiplexer: Zellij

- **Overriding Options:**  
  Modify the `gentlemanOptions` block in `flake.nix` or use Nix overrides at build time.

- **Local Configuration Files:**  
  All configurations are defined inline in local modules (e.g., `fish.nix`, `nushell.nix`, etc.) and are deployed automatically to system-specific locations. For example, the Nushell configuration is copied to:

  - **macOS:** `~/Library/Application Support/nushell`
  - **Linux/WSL:** `~/.config/nushell`

- **Dependencies & Automatic Replacements:**  
  The flake installs all necessary dependencies (git, curl, rustc, cargo, tmux, etc.) and performs placeholder replacements in configuration files (e.g., replacing “tmux” with “zellij” when applicable).

- **Windows Users:**  
  Must install and configure WSL, manually install the Iosevka Term Nerd Font, set up Alacritty or WezTerm, install Chocolatey with win32yank, and finally launch and update the Linux distribution.

For any questions or further customizations, please open an issue or submit a pull request.

**Happy coding!**

— Gentleman
