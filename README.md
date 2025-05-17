# Gentleman.Dots

<img width="1920" alt="image" src="https://github.com/user-attachments/assets/746e6926-427c-4690-a0ba-07d8e2015d19" />

---

# Gentleman.Dots

## Description

This repository contains customized configurations for a complete development environment, including:

- Neovim
- **Nushell**
- **Fish**
- **Zsh**
- Terminal emulators:
  - **WezTerm**
  - **Ghostty**

You can now automatically set up your environment using our new Nix Flake approach with Home Manager. This method is fully declarative and reproducible, and it lets you easily override default options. In our flake, all configurations are defined inline in local modules (e.g., `nushell.nix`, etc.), and the flake also installs all the required dependencies (git, curl, rustc, cargo, neovim, etc.).

---

## Previous Steps

## Installation Steps (for macOS/Linux/WSL)

### 1. Install the Nix Package Manager

```bash
sh <(curl -L https://nixos.org/nix/install)
```

### 2. Configure Nix to Use Extra Experimental Features

To enable the experimental features for flakes and the new `nix-command` (needed for our declarative setup), open the configuration file and add the following line:

```bash
sudo nvim /etc/nix/nix.conf
```

Add:

```
extra-experimental-features = flakes nix-command
```

_(This is necessary because support for flakes and the new Nix command is still experimental, but it allows us to have a fully declarative and reproducible configuration.)_

### 1.1. Archalinux

### 1.1. Install Nix and enable flakes on Arch Linux 
   ```bash
   sudo pacman -S nix
```
#### 1.1.1 If it doesn't work:

Download the file with `curl --proto '=https' --tlsv1.2 -sSfL https://nixos.org/nix/install -o nix-install.sh`, view it: `less ./nix-install.sh`, and run the script `./nix-install.sh --daemon` to start Nix installation. 

### 1.2 Configure Nix to Use Extra Experimental Features

To enable the experimental features for flakes and the new `nix-command` (needed for our declarative setup), open the configuration file and add the following lin
e:
```bash
# In /etc/nix/nix.conf:
   experimental-features = nix-command flakes
   sudo systemctl enable --now nix-daemon.servi
   ce
```

### 3. Prepare the `flake.nix` File

Before running the next command, you need to make a few changes in the `flake.nix` file to match your environment. **Make sure to update the `system` field with your system's value**. Below are the possible values:

- **"aarch64-darwin"**  
  _Architecture: Apple Silicon (M1, M2, M3, etc.)._  
  _Operating System: macOS._

- **"x86_64-darwin"**  
  _Architecture: Intel (older Macs)._  
  _Operating System: macOS._

- **"x86_64-linux"**  
  _Architecture: 64-bit Intel/AMD._  
  _Operating System: Linux._

- **"aarch64-linux"**  
  _Architecture: 64-bit ARM (e.g., Raspberry Pi 4, ARM servers)._  
  _Operating System: Linux._

- **"i686-linux"**  
  _Architecture: 32-bit Intel/AMD (obsolete)._  
  _Operating System: Linux._

- **"riscv64-linux"**  
  _Architecture: 64-bit RISC-V._  
  _Operating System: Linux (experimental)._

- **"x86_64-freebsd"**  
  _Architecture: 64-bit Intel/AMD._  
  _Operating System: FreeBSD._

Modify the parameters in your `flake.nix` file as follows:

- Change the line `home.username = "YourUser";` to reflect your machine's username.

- Install your terminal emulator, configs will we already applied:

  - Wezterm: <https://wezterm.org/installation.html>
  - Ghostty: <https://ghostty.org/download> _Remember to reload Ghostty's Config inside the terminal_**(shift + command + ,)**

- Modify `home.homeDirectory` accordingly:

  - On macOS: `/Users/YourUser`
  - On Linux: `/home/YourUser`

- If you want my `aerospace` tile windows manager configuration you can copy the one inside `./aerospace/.aerospace.toml` into your `$HOME` path

- **Don't forget to update the `system` field** (currently set to `"aarch64-darwin"`) with the appropriate value from the list above.

### 4. Run the Installation

Once you're in the repo directory and have made the above changes, run:

```bash
nix run github:nix-community/home-manager -- switch --flake .#gentleman -b backup
```

_(This command applies the configuration defined in the flake, installing all dependencies and applying the necessary settings.)_

### 5. Default Shell

Now run the following script to add `Nushell`, `Fish` or `Zsh` to your list of available shells and select it as the default one:

**Fish:**

```bash
shellPath=$(which fish)

sudo sh -c "grep -Fxq '$shellPath' /etc/shells || echo '$shellPath' >> /etc/shells"
sudo chsh -s "$shellPath" "$USER"
```

**Nushell:**

```bash
shellPath=$(which nu)

sudo sh -c "grep -Fxq '$shellPath' /etc/shells || echo '$shellPath' >> /etc/shells"
sudo chsh -s "$shellPath" "$USER"
```

**Zsh:**

```bash
shellPath=$(which zsh)

sudo sh -c "grep -Fxq '$shellPath' /etc/shells || echo '$shellPath' >> /etc/shells"
sudo chsh -s "$shellPath" "$USER"
```

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

## Summary

- **Local Configuration Files:**  
  All configurations are defined inline in local modules (e.g., `fish.nix`, `nushell.nix`, etc.) and are deployed automatically to system-specific locations. For example, the Nushell configuration is copied to:

  - **macOS:** `~/Library/Application Support/nushell`
  - **Linux/WSL:** `~/.config/nushell`

- **Dependencies & Automatic Replacements:**  
  The flake installs all necessary dependencies (git, curl, rustc, cargo, tmux, etc.) and performs placeholder replacements in configuration files.

- **Windows Users:**  
  Must install and configure WSL, manually install the Iosevka Term Nerd Font, set up Alacritty or WezTerm, install Chocolatey with win32yank, and finally launch and update the Linux distribution.

For any questions or further customizations, please open an issue or submit a pull request.

**Happy coding!**

— Gentleman
