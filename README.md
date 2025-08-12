# Gentleman.Dots

<img width="1920" alt="image" src="https://github.com/user-attachments/assets/746e6926-427c-4690-a0ba-07d8e2015d19" />

---

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

To enable the experimental features for flakes and the new `nix-command` (needed for our declarative setup), create/edit the configuration file:

```bash
# For daemon installation (default on macOS/Linux)
# The file may not exist, create it if needed
sudo mkdir -p /etc/nix
sudo nano /etc/nix/nix.conf
# Or: sudo vi /etc/nix/nix.conf
```

Add:

```
extra-experimental-features = flakes nix-command
build-users-group = nixbld
```

_(This is necessary because support for flakes and the new Nix command is still experimental, but it allows us to have a fully declarative and reproducible configuration.)_

### 2.1. Arch Linux Specific Instructions

#### 2.1.1. Install Nix on Arch Linux

```bash
sudo pacman -S nix
```

#### 2.1.2. Alternative installation method

Download the file with `curl --proto '=https' --tlsv1.2 -sSfL https://nixos.org/nix/install -o nix-install.sh`, view it: `less ./nix-install.sh`, and run the script `./nix-install.sh --daemon` to start Nix installation.

#### 2.1.3. Enable Nix daemon on Arch Linux

After configuring `/etc/nix/nix.conf` as shown in step 2 above, enable the Nix daemon:

```bash
sudo systemctl enable --now nix-daemon.service
```

### 3. Prepare Your System

**No need to edit `flake.nix` for system configuration!** The flake now automatically provides configurations for all systems.

You only need to update your username in `flake.nix`:
- Change `home.username = "anua";` to your actual username
- The home directory is automatically set based on your system

- Install your terminal emulator, configs will be already applied:
  - Wezterm: <https://wezterm.org/installation.html>
  - Ghostty: <https://ghostty.org/download> _Remember to reload Ghostty's Config inside the terminal_**(shift + command + ,)**

- If you want my `aerospace` tile windows manager configuration you can copy the one inside `./aerospace/.aerospace.toml` into your `$HOME` path

### 4. Run the Installation

Once you have cloned the repository and are **inside its directory**, run the following command.

**⚠️ Important:** You must be in the root of this project directory for the command to work, as it uses `.` to find the `flake.nix` file.

```bash
# Make sure you are in the Gentleman.Dots directory before running!
nix run github:nix-community/home-manager -- switch --flake .#gentleman -b backup
```

_(This command applies the configuration defined in the flake, installing all dependencies and applying the necessary settings.)_

### 5. Verify Installation

**For macOS users: PATH is configured automatically!**

**For WSL/Linux users: PATH is mostly automatic, but verify it works:**

```bash
hash -r  # Refresh command cache
which fish   # Should show path to fish
which nvim   # Should show path to nvim
which nu     # Should show path to nu
```

**If commands are not found on WSL/Linux, manually add to your shell config:**

- **Bash** (`~/.bashrc`):
  ```bash
  echo '. "$HOME/.nix-profile/etc/profile.d/nix.sh"' >> ~/.bashrc
  echo '[ -f "$HOME/.nix-profile/etc/profile.d/hm-session-vars.sh" ] && . "$HOME/.nix-profile/etc/profile.d/hm-session-vars.sh"' >> ~/.bashrc
  echo 'export PATH="$HOME/.local/state/nix/profiles/home-manager/bin:$HOME/.nix-profile/bin:$PATH"' >> ~/.bashrc
  source ~/.bashrc
  ```

- **Zsh** (`~/.zshrc`):
  ```bash
  echo '. "$HOME/.nix-profile/etc/profile.d/nix.sh"' >> ~/.zshrc
  echo '[ -f "$HOME/.nix-profile/etc/profile.d/hm-session-vars.sh" ] && . "$HOME/.nix-profile/etc/profile.d/hm-session-vars.sh"' >> ~/.zshrc
  echo 'export PATH="$HOME/.local/state/nix/profiles/home-manager/bin:$HOME/.nix-profile/bin:$PATH"' >> ~/.zshrc
  source ~/.zshrc
  ```

### 6. Default Shell

Now run the following script to add `Nushell`, `Fish` or `Zsh` to your list of available shells and select it as the default one:

**Fish:**

```bash
shellPath="$HOME/.nix-profile/bin/fish"

sudo sh -c "grep -Fxq '$shellPath' /etc/shells || echo '$shellPath' >> /etc/shells"
sudo chsh -s "$shellPath" "$USER"
```

**Nushell:**

```bash
shellPath="$HOME/.nix-profile/bin/nu"

sudo sh -c "grep -Fxq '$shellPath' /etc/shells || echo '$shellPath' >> /etc/shells"
sudo chsh -s "$shellPath" "$USER"
```

**Zsh:**

```bash
shellPath="$HOME/.nix-profile/bin/zsh"

sudo sh -c "grep -Fxq '$shellPath' /etc/shells || echo '$shellPath' >> /etc/shells"
sudo chsh -s "$shellPath" "$USER"
```

---

## Manual Installation for Windows

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

#### 5. Clone the Repository

Clone the repository to get the configuration files. Open PowerShell and run:

```powershell
git clone https://github.com/Gentleman-Programming/Gentleman.Dots.git
cd Gentleman.Dots
```

#### 6. Transfer Emulator Configurations

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

#### 7. Install Chocolatey and win32yank

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

#### 8. Launch and Update Your Linux Distribution

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

#### 9. Configure Nix and Home Manager in WSL

**⚠️ IMPORTANT FOR WSL USERS: Follow these steps to use Nix with Home Manager in WSL:**

1. **Install Nix in WSL:**

   ```bash
   sh <(curl -L https://nixos.org/nix/install) --no-daemon
   ```

2. **Enable flakes in Nix:**

   ```bash
   # For non-daemon installation (WSL uses --no-daemon)
   # Create the config directory and file
   mkdir -p ~/.config/nix
   echo "extra-experimental-features = flakes nix-command" >> ~/.config/nix/nix.conf
   ```

3. **Source Nix in your current shell:**

   ```bash
   . ~/.nix-profile/etc/profile.d/nix.sh
   ```

4. **Add Nix to your shell configuration permanently:**

   For Bash (`~/.bashrc`):

   ```bash
   echo '. ~/.nix-profile/etc/profile.d/nix.sh' >> ~/.bashrc
   ```

   For Zsh (`~/.zshrc`):

   ```bash
   echo '. ~/.nix-profile/etc/profile.d/nix.sh' >> ~/.zshrc
   ```

5. **Clone and configure the repository:**

   ```bash
   git clone https://github.com/Gentleman-Programming/Gentleman.Dots.git
   cd Gentleman.Dots
   ```

6. **Edit `flake.nix` with your WSL configuration:**
   - Change `home.username = "YourUsername"`
   - The home directory will be automatically set to `/home/YourUsername`

7. **Run Home Manager with the Linux configuration:**

   **⚠️ Important:** Make sure you are in the `Gentleman.Dots` directory (cloned in step 5) before running this command.

   ```bash
   nix run github:nix-community/home-manager -- switch --flake .#gentleman-linux-x64 -b backup
   ```

8. **Configure PATH for the installed programs:**

   **For WSL, the PATH configuration is now automatic!** But if needed, add these to your shell config:

   ```bash
   # For Bash (~/.bashrc)
   echo '. "$HOME/.nix-profile/etc/profile.d/nix.sh"' >> ~/.bashrc
   echo '[ -f "$HOME/.nix-profile/etc/profile.d/hm-session-vars.sh" ] && . "$HOME/.nix-profile/etc/profile.d/hm-session-vars.sh"' >> ~/.bashrc
   echo 'export PATH="$HOME/.local/state/nix/profiles/home-manager/bin:$HOME/.nix-profile/bin:$PATH"' >> ~/.bashrc
   source ~/.bashrc
   ```

9. **Verify the installation:**

   ```bash
   hash -r  # Refresh command cache
   which fish   # Should show path to fish
   which nvim   # Should show path to nvim
   which nu     # Should show path to nu
   ```

   **Note:** In WSL, binaries might be in `~/.local/state/nix/profiles/home-manager/bin/` or `~/.nix-profile/bin/`

10. **Set your default shell (optional):**

    ```bash
    # For Fish
    sudo sh -c "echo '$HOME/.nix-profile/bin/fish' >> /etc/shells"
    chsh -s "$HOME/.nix-profile/bin/fish"

    # For Nushell
    sudo sh -c "echo '$HOME/.nix-profile/bin/nu' >> /etc/shells"
    chsh -s "$HOME/.nix-profile/bin/nu"
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
