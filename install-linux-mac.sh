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
      if is_wsl; then
        echo -e "${YELLOW}You are running WSL. Please install Alacritty from Windows.${NC}"
      else
        sudo apt-get install -y alacritty
      fi
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
      if is_wsl; then
        echo -e "${YELLOW}You are running WSL. Please install WezTerm from Windows.${NC}"
      else
        sudo apt-get install -y wezterm
      fi
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

# Shared Steps (macOS, Linux, or WSL)

# Step 3: Shell Configuration (Fish and Zsh)
echo -e "${YELLOW}Step 3: Choose and Install Shell${NC}"
shell_choice=$(prompt_user "Which shell do you want to install? (Options: fish, zsh)" "none")

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
  sed -i "s|/your/work/path/|$PROJECT_PATHS|g" ~/.config/fish/config.fish
  sudo sh -c "echo $(which fish) >> /etc/shells"
  chsh -s $(which fish)
  curl -sL https://git.io/fisher | source && fisher install jorgebucaran/fisher
  fisher install oh-my-fish/plugin-pj
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
    prompt_user "After it's done, write 'exit' to continue... press enter now"
    sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"
  else
    echo -e "${GREEN}Oh My Zsh is already installed.${NC}"
  fi

  echo -e "${YELLOW}Configuring Zsh...${NC}"
  cp .zshrc ~/
  sed -i "s|/your/work/path/|$PROJECT_PATHS|g" ~/.zshrc
  sudo sh -c "echo $(which zsh) >> /etc/shells"
  chsh -s $(which zsh)
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
  sudo apt-get install build-essential
fi

brew install starship nvim node npm git gcc fzf fd ripgrep coreutils

# Neovim Configuration
echo -e "${YELLOW}Configuring Neovim...${NC}"
cp -r GentlemanNvim/nvim ~/.config

# Starship Configuration
echo -e "${YELLOW}Configuring Starship...${NC}"
cp starship.toml ~/.config

# Obsidian Configuration
echo -e "${YELLOW}Configuring Obsidian...${NC}"
if [ -f ~/.config/nvim/lua/plugins/obsidian.lua ]; then
  sed -i "s|/your/notes/path|$OBSIDIAN_PATH|g" ~/.config/nvim/lua/plugins/obsidian.lua
else
  echo -e "${RED}Obsidian configuration file not found. Please check your setup.${NC}"
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
  cp -r GentlemanTmux/.tmux ~/
  cp GentlemanTmux/.tmux.conf ~/

  # Update .zshrc and config.fish to use Tmux
  if [ "$shell_choice" = "zsh" ]; then
    sed -i "s/exec tmux/exec tmux/g" ~/.zshrc
    sed -i "s/\$TMUX/\$TMUX/g" ~/.zshrc
  elif [ "$shell_choice" = "fish" ]; then
    sed -i "s/TMUX/TMUX/g" ~/.config/fish/config.fish
    sed -i "s/tmux/tmux/g" ~/.config/fish/config.fish

  fi

  echo -e "${YELLOW}Starting Tmux and Loading Configuration...${NC}"
  tmux new-session -d -s mysession "tmux source-file ~/.tmux.conf"
  tmux kill-session -t mysession
  ;;
"zellij")
  if ! command -v zellij &>/dev/null; then
    brew install zellij
  else
    echo -e "${GREEN}Zellij is already installed.${NC}"
  fi
  echo -e "${YELLOW}Configuring Zellij...${NC}"
  mkdir -p ~/.config/zellij
  cp -r GentlemanZellij/zellij ~/.config

  # Update the default shell in Zellij config and .zshrc or config.fish
  if [ "$shell_choice" = "fish" ]; then
    sed -i "s|default_shell \"fish\"|default_shell \"fish\"|g" ~/.config/zellij/config.kdl
    sed -i "s/TMUX/ZELLIJ/g" ~/.config/fish/config.fish
    sed -i "s/tmux/zellij/g" ~/.config/fish/config.fish
  elif [ "$shell_choice" = "zsh" ]; then
    sed -i "s|default_shell \"fish\"|default_shell \"zsh\"|g" ~/.config/zellij/config.kdl
    sed -i "s/exec tmux/exec zellij/g" ~/.zshrc
    sed -i "s/\$TMUX/\$ZELLIJ/g" ~/.zshrc
  fi
  zellij
  ;;
*)
  echo -e "${YELLOW}No window manager will be installed or configured.${NC}"
  ;;
esac

# Clean up: Remove the cloned repository
echo -e "${YELLOW}Cleaning up...${NC}"
cd ..
rm -rf Gentleman.Dots

echo -e "${GREEN}Installation and configuration complete! Please restart your terminal to see the changes.${NC}"
