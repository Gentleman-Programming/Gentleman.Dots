#!/bin/bash

set -e

# Define colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

USER_HOME=$(eval echo ~${SUDO_USER})

echo -e "${GREEN}Welcome to the Gentleman.Dots installation and configuration guide!${NC}"

# Function to prompt user for input
prompt_user() {
  local prompt_message="$1"
  local default_answer="$2"
  read -p "$prompt_message [$default_answer] " user_input
  user_input="${user_input:-$default_answer}"
  echo "$user_input"
}

# Function to check if running on WSL
is_wsl() {
  grep -qEi "(Microsoft|WSL)" /proc/version &>/dev/null
  return $?
}

# Function to install basic dependencies
install_dependencies() {
  if [ "$os_choice" = "linux" ]; then
    sudo apt-get update
    sudo apt-get install -y build-essential curl file git
  fi
}

# Function to install Homebrew if not installed
install_homebrew() {
  if ! command -v brew &>/dev/null; then
    echo -e "${YELLOW}Homebrew is not installed. Installing Homebrew...${NC}"
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

    # Add Homebrew to PATH based on OS
    if [ "$os_choice" = "mac" ]; then
      (
        echo
        echo 'eval "$(/opt/homebrew/bin/brew shellenv)"'
      ) >>$USER_HOME/.zshrc
      (
        echo
        echo 'eval "$(/opt/homebrew/bin/brew shellenv)"'
      ) >>$USER_HOME/.bashrc
      mkdir -p $USER_HOME/.config/fish
      (
        echo
        echo 'eval "$(/opt/homebrew/bin/brew shellenv)"'
      ) >>$USER_HOME/.config/fish/config.fish
      eval "$(/opt/homebrew/bin/brew shellenv)"
    elif [ "$os_choice" = "linux" ]; then
      (
        echo
        echo 'eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"'
      ) >>$USER_HOME/.zshrc
      (
        echo
        echo 'eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"'
      ) >>$USER_HOME/.bashrc
      mkdir -p $USER_HOME/.config/fish
      (
        echo
        echo 'eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"'
      ) >>$USER_HOME/.config/fish/config.fish
      eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"
    fi
  else
    echo -e "${GREEN}Homebrew is already installed.${NC}"
  fi
}

# Ask for the operating system
os_choice=$(prompt_user "Which operating system are you using? (Options: mac, linux)" "none")

# Install basic dependencies
install_dependencies

# Prompt for project path and Obsidian path
PROJECT_PATHS=$(prompt_user "Enter the path for your projects" "/your/work/path/")
OBSIDIAN_PATH=$(prompt_user "Enter the path for your Obsidian vault" "/your/notes/path")

# Step 1: Clone the Repository
echo -e "${YELLOW}Step 1: Clone the Repository${NC}"
if [ -d "Gentleman.Dots" ]; then
  echo -e "${GREEN}Repository already cloned. Skipping...${NC}"
else
  git clone https://github.com/Gentleman-Programming/Gentleman.Dots.git
fi
cd Gentleman.Dots || exit

# Install Homebrew if not installed
install_homebrew

# Terminal Emulators Installation
echo -e "${YELLOW}Step 2: Choose and Install Terminal Emulator${NC}"
if is_wsl; then
  echo -e "${YELLOW}You are running WSL. Terminal emulators should be installed on Windows.${NC}"
else
  if [ "$os_choice" = "linux" ]; then
    echo -e "${YELLOW}Note: Kitty is not available for Linux.${NC}"
    term_choice=$(prompt_user "Which terminal emulator do you want to install? (Options: alacritty, wezterm)" "none")
  else
    term_choice=$(prompt_user "Which terminal emulator do you want to install? (Options: alacritty, wezterm, kitty)" "none")
  fi

  case "$term_choice" in
  "alacritty")
    if ! command -v alacritty &>/dev/null; then
      if [ "$os_choice" = "mac" ]; then
        brew install --cask alacritty
      elif [ "$os_choice" = "linux" ]; then
        sudo add-apt-repository ppa:aslatter/ppa
        sudo apt-get update
        sudo apt-get install alacritty
      fi
    else
      echo -e "${GREEN}Alacritty is already installed.${NC}"
    fi
    echo -e "${YELLOW}Configuring Alacritty...${NC}"
    mkdir -p ~/.config/alacritty
    cp alacritty.toml ~/.config/alacritty/alacritty.toml
    ;;
  "wezterm")
    if ! command -v wezterm &>/dev/null; then
      if [ "$os_choice" = "mac" ]; then
        brew install --cask wezterm
      elif [ "$os_choice" = "linux" ]; then
        curl -fsSL https://apt.fury.io/wez/gpg.key | sudo gpg --yes --dearmor -o /usr/share/keyrings/wezterm-fury.gpg
        echo 'deb [signed-by=/usr/share/keyrings/wezterm-fury.gpg] https://apt.fury.io/wez/ * *' | sudo tee /etc/apt/sources.list.d/wezterm.list
        sudo apt update
        sudo apt install wezterm
      fi
    else
      echo -e "${GREEN}WezTerm is already installed.${NC}"
    fi
    echo -e "${YELLOW}Configuring WezTerm...${NC}"
    mkdir -p ~/.config/wezterm
    cp .wezterm.lua ~/.config/wezterm/wezterm.lua
    ;;
  "kitty")
    if [ "$os_choice" = "mac" ]; then
      if ! command -v kitty &>/dev/null; then
        brew install --cask kitty
      else
        echo -e "${GREEN}Kitty is already installed.${NC}"
      fi
      echo -e "${YELLOW}Configuring Kitty...${NC}"
      mkdir -p ~/.config/kitty
      cp -r GentlemanKitty/* ~/.config/kitty
    else
      echo -e "${YELLOW}Kitty installation is not available for Linux.${NC}"
    fi
    ;;
  *)
    echo -e "${YELLOW}No terminal emulator will be installed or configured.${NC}"
    ;;
  esac
fi

# Shared Steps (macOS, Linux, or WSL)

# Step 3: Shell Configuration (Fish and Zsh)
echo -e "${YELLOW}Step 3: Choose and Install Shell${NC}"
shell_choice=$(prompt_user "Which shell do you want to install? (Options: fish, zsh)" "none")

# Function to update or replace a line in a file
update_or_replace() {
  local file="$1"
  local search="$2"
  local replace="$3"

  if grep -q "$search" "$file"; then
    # Use awk to replace the line containing the search string
    awk -v search="$search" -v replace="$replace" '
    $0 ~ search {print replace; next}
    {print}
    ' "$file" >"${file}.tmp" && mv "${file}.tmp" "$file"
  else
    echo "$replace" >>"$file"
  fi
}

set_default_shell() {
  local shell_path="$1"

  if ! grep -Fxq "$shell_path" /etc/shells; then
    echo "$shell_path" | sudo tee -a /etc/shells
  fi

  sudo chsh -s "$shell_path" "$USER"
}

case "$shell_choice" in
"fish")
  if ! command -v fish &>/dev/null; then
    brew install fish
  else
    echo -e "${GREEN}Fish shell is already installed.${NC}"
  fi
  echo -e "${YELLOW}Configuring Fish shell...${NC}"
  mkdir -p ~/.config/fish
  cp -r GentlemanFish/* ~/.config/fish
  # Update or append the PROJECT_PATHS line
  update_or_replace ~/.config/fish/config.fish "set PROJECT_PATHS" "set PROJECT_PATHS $PROJECT_PATHS"

  # Set fish as the default shell
  set_default_shell "$(which fish)"

  # Install Fisher and plugins
  fish -c "curl -sL https://git.io/fisher | source && fisher install jorgebucaran/fisher"
  fish -c "fisher install oh-my-fish/plugin-pj"
  ;;
"zsh")
  if ! command -v zsh &>/dev/null; then
    brew install zsh
  else
    echo -e "${GREEN}Zsh is already installed.${NC}"
  fi
  brew install zsh-autosuggestions zsh-syntax-highlighting zsh-autocomplete

  if [ ! -d "$HOME/.oh-my-zsh" ]; then
    echo -e "${YELLOW}Installing Oh My Zsh...${NC}"
    echo -e "${YELLOW}After its done installing, just write exit and press enter to continue with the process${NC}"
    prompt_user "Press enter to continue"
    sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"
  else
    echo -e "${GREEN}Oh My Zsh is already installed.${NC}"
  fi

  echo -e "${YELLOW}Configuring Zsh...${NC}"
  cp .zshrc ~/
  # Update or append the PROJECT_PATHS line
  update_or_replace ~/.zshrc "export PROJECT_PATHS" "export PROJECT_PATHS=\"$PROJECT_PATHS\""

  # Set zsh as the default shell
  set_default_shell "$(which zsh)"
  ;;
*)
  echo -e "${YELLOW}No shell will be installed or configured.${NC}"
  ;;
esac

# Step 4: Additional Configurations

# Dependencies Install
echo -e "${YELLOW}Step 4: Installing Additional Dependencies...${NC}"

if [ "$os_choice" = "linux" ]; then
  sudo apt-get update
  sudo apt-get upgrade
fi

brew install node npm git gcc fzf fd ripgrep coreutils bat curl lazygit

# Neovim Configuration
echo -e "${YELLOW}Configuring Neovim...${NC}"
mkdir -p ~/.config/nvim
cp -r GentlemanNvim/nvim/* ~/.config/nvim/

# Starship Configuration
echo -e "${YELLOW}Configuring Starship...${NC}"
mkdir -p ~/.config
cp starship.toml ~/.config

# Obsidian Configuration
echo -e "${YELLOW}Configuring Obsidian...${NC}"
obsidian_config_file="$HOME/.config/nvim/lua/plugins/obsidian.lua"
if [ -f "$obsidian_config_file" ]; then
  # Replace the vault path in the existing configuration
  update_or_replace "$obsidian_config_file" "/your/notes/path" "path = '$OBSIDIAN_PATH'"
else
  echo -e "${RED}Obsidian configuration file not found at $obsidian_config_file. Please check your setup.${NC}"
fi

# Ask if they want to use Tmux or Zellij
wm_choice=$(prompt_user "Which window manager do you want to install? (Options: tmux, zellij)" "none")

case "$wm_choice" in
"tmux")
  if ! command -v tmux &>/dev/null; then
    brew install tmux
  else
    echo -e "${GREEN}Tmux is already installed.${NC}"
  fi
  echo -e "${YELLOW}Configuring Tmux...${NC}"
  if [ ! -d "$HOME/.tmux/plugins/tpm" ]; then
    git clone https://github.com/tmux-plugins/tpm ~/.tmux/plugins/tpm
  else
    echo -e "${GREEN}Tmux Plugin Manager is already installed.${NC}"
  fi
  mkdir -p ~/.tmux
  cp -r GentlemanTmux/.tmux/* ~/.tmux/
  cp GentlemanTmux/.tmux.conf ~/

  echo -e "${YELLOW}Please restart your computer to complete the Tmux installation after the script is done.${NC}"
  echo -e "${YELLOW}After restarting, press Ctrl + a followed by Shift + i to install TMUX plugins.${NC}"
  prompt_user "Press enter to continue"
  ;;
"zellij")
  if ! command -v zellij &>/dev/null; then
    brew install zellij
  else
    echo -e "${GREEN}Zellij is already installed.${NC}"
  fi
  echo -e "${YELLOW}Configuring Zellij...${NC}"
  mkdir -p ~/.config/zellij
  cp -r GentlemanZellij/zellij/* ~/.config/zellij/

  # Replace TMUX with ZELLIJ and tmux with zellij only in the selected shell configuration
  if [[ "$shell_choice" == "zsh" ]]; then
    update_or_replace ~/.zshrc "TMUX" "ZELLIJ"
    update_or_replace ~/.zshrc "tmux" "zellij"
  elif [[ "$shell_choice" == "fish" ]]; then
    update_or_replace ~/.config/fish/config.fish "TMUX" "ZELLIJ"
    update_or_replace ~/.config/fish/config.fish "tmux" "zellij"
  fi
  ;;
*)
  echo -e "${YELLOW}No window manager will be installed or configured.${NC}"
  ;;
esac

# Clean up: Remove the cloned repository
echo -e "${YELLOW}Cleaning up...${NC}"
cd ..
rm -rf Gentleman.Dots

echo -e "${YELLOW}After restarting, if you installed TMUX, remember to press Ctrl + a followed by Shift + i to install the plugins.${NC}"
echo -e "${GREEN}Installation and configuration complete! Please restart your computer to see the changes.${NC}"
