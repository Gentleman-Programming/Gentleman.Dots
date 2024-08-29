#!/bin/bash

set -e

# Define colors for output using tput for better compatibility
define_colors() {
  PINK=$(tput setaf 204)
  PURPLE=$(tput setaf 141)
  GREEN=$(tput setaf 114)
  ORANGE=$(tput setaf 208)
  BLUE=$(tput setaf 75)
  YELLOW=$(tput setaf 221)
  RED=$(tput setaf 196)
  NC=$(tput sgr0) # No Color
}

define_colors

# Gentleman.Dots logo with pink color
display_logo() {
  local logo='
                      ░░░░░░      ░░░░░░                      
                    ░░░░░░░░░░  ░░░░░░░░░░                    
                  ░░░░░░░░░░░░░░░░░░░░░░░░░░                  
                ░░░░░░░░░░▒▒▒▒░░▒▒▒▒░░░░░░░░░░                
              ░░░░░░▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒░░░░░░              
  ▒▒        ░░░░░░▒▒▒▒▒▒▒▒▒▒██▒▒██▒▒▒▒▒▒▒▒▒▒░░░░░░        ▒▒  
▒▒░░    ░░░░░░░░▒▒▒▒▒▒▒▒▒▒████▒▒████▒▒▒▒▒▒▒▒▒▒░░░░░░░░    ░░▒▒
▒▒▒▒░░░░░░▒▒▒▒▒▒▒▒▒▒▒▒▒▒██████▒▒██████▒▒▒▒▒▒▒▒▒▒▒▒▒▒░░░░░░▒▒▒▒
██▒▒▒▒▒▒▒▒▒▒▒▒▒▒██▒▒▒▒██████▓▓██▒▒██████▒▒▓▓██▒▒▒▒▒▒▒▒▒▒▒▒▒▒██
  ████▒▒▒▒▒▒████▒▒▒▒██████████  ██████████▒▒▒▒████▒▒▒▒▒▒▒▒██  
    ████████████████████████      ████████████████████████    
      ██████████████████              ██████████████████      
          ██████████                      ██████████          
'
  echo -e "${PINK}${logo}${NC}"
  echo -e "${PURPLE}Welcome to the Gentleman.Dots Auto Config!${NC}"
}

display_logo

# Keep sudo session alive
keep_sudo_alive() {
  sudo -v
  while true; do
    sudo -n true
    sleep 60
    kill -0 "$$" || exit
  done 2>/dev/null &
}

keep_sudo_alive

# Function to prompt user for input with a select menu
select_option() {
  local prompt_message="$1"
  shift
  local options=("$@")
  PS3="${ORANGE}$prompt_message${NC} "
  select opt in "${options[@]}"; do
    if [ -n "$opt" ]; then
      echo "$opt"
      break
    else
      echo -e "${RED}Invalid option. Please try again.${NC}"
    fi
  done
}

# Function to prompt user for input with a default option
prompt_user() {
  local prompt_message="$1"
  local default_answer="$2"
  read -p "$(echo -e ${BLUE}$prompt_message [$default_answer]${NC}) " user_input
  user_input="${user_input:-$default_answer}"
  echo "$user_input"
}

# Function to display a spinner
spinner() {
  local pid=$!
  local delay=0.1
  local spinstr='|/-\'
  while [ "$(ps a | awk '{print $1}' | grep $pid)" ]; do
    local temp=${spinstr#?}
    printf " [%c]  " "$spinstr"
    spinstr=$temp${spinstr%"$temp"}
    sleep $delay
    printf "\b\b\b\b\b\b"
  done
  printf "    \b\b\b\b"
}

# Function to check and create directories if they do not exist
ensure_directory_exists() {
  local dir_path="$1"
  local create_templates="$2"

  if [ ! -d "$dir_path" ]; then
    echo -e "${YELLOW}Directory $dir_path does not exist. Creating...${NC}"
    mkdir -p "$dir_path"
  else
    echo -e "${GREEN}Directory $dir_path already exists.${NC}"
  fi

  if [ "$create_templates" == "true" ]; then
    ensure_directory_exists "$dir_path/templates" "false"
  fi
}

# Function to check if running on WSL
is_wsl() {
  grep -qEi "(Microsoft|WSL)" /proc/version &>/dev/null
}

# Function to run commands with optional suppression of output
run_command() {
  local command=$1
  if [ "$show_details" = "Yes" ]; then
    eval $command
  else
    eval $command &>/dev/null
  fi
}

# Function to detect if the system is Arch Linux
is_arch() {
  [ -f /etc/arch-release ]
}

# Function to install dependencies
install_dependencies() {
  if is_arch; then
    run_command "sudo pacman -Syu --noconfirm"
    run_command "sudo pacman -S --needed --noconfirm base-devel curl file git"
  else
    run_command "sudo apt-get update"
    run_command "sudo apt-get install -y build-essential curl file git"
  fi
}

# Function to install Homebrew if not installed
install_homebrew() {
  if ! command -v brew &>/dev/null; then
    echo -e "${YELLOW}Homebrew is not installed. Installing Homebrew...${NC}"
    run_command "/bin/bash -c \"\$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\""

    local brew_prefix="/opt/homebrew/bin/brew"
    [ "$os_choice" = "linux" ] && brew_prefix="/home/linuxbrew/.linuxbrew/bin/brew"
    add_to_shell_path "$brew_prefix"
  else
    echo -e "${GREEN}Homebrew is already installed.${NC}"
  fi
}

# Function to add Homebrew to shell path
add_to_shell_path() {
  local brew_path="$1"
  run_command "(echo 'eval \"\$(\"$brew_path\" shellenv)\"' >> $USER_HOME/.zshrc)"
  run_command "(echo 'eval \"\$(\"$brew_path\" shellenv)\"' >> $USER_HOME/.bashrc)"
  run_command "mkdir -p $USER_HOME/.config/fish"
  run_command "(echo 'eval \"\$(\"$brew_path\" shellenv)\"' >> $USER_HOME/.config/fish/config.fish)"
  run_command "eval \"\$(\"$brew_path\" shellenv)\""
}

# Function to update or replace a line in a file
update_or_replace() {
  local file="$1"
  local search="$2"
  local replace="$3"

  if grep -q "$search" "$file"; then
    awk -v search="$search" -v replace="$replace" '
    $0 ~ search {print replace; next}
    {print}
    ' "$file" >"${file}.tmp" && mv "${file}.tmp" "$file"
  else
    echo "$replace" >>"$file"
  fi
}

# Function to set the default shell
set_default_shell() {
  local shell_path="$1"
  sudo sh -c "grep -Fxq \"$shell_path\" /etc/shells || echo \"$shell_path\" >> /etc/shells"
  sudo chsh -s "$shell_path" "$USER"
}

# Function to install and configure a terminal emulator
install_terminal() {
  local term_name="$1"
  local install_command="$2"
  local config_command="$3"

  echo -e "${YELLOW}Installing $term_name...${NC}"
  run_command "$install_command"
  echo -e "${YELLOW}Configuring $term_name...${NC}"
  run_command "$config_command"
}

# Function to install and configure a shell
install_shell() {
  local shell_name="$1"
  local install_command="$2"
  local post_install_command="$3"

  echo -e "${YELLOW}Installing $shell_name...${NC}"
  run_command "$install_command"
  run_command "$post_install_command"
  set_default_shell "$(which $shell_name)"
}

# Function to install a window manager
install_window_manager() {
  local wm_name="$1"
  local install_command="$2"
  local config_command="$3"

  echo -e "${YELLOW}Installing $wm_name...${NC}"
  run_command "$install_command"
  echo -e "${YELLOW}Configuring $wm_name...${NC}"
  run_command "$config_command"
}

# Main script execution
show_details=$(select_option "Do you want to see detailed output? " "No" "Yes")
os_choice=$(select_option "Which operating system are you using? " "mac" "linux")

if [ "$os_choice" != "mac" ]; then
  echo -e "${YELLOW}Installing basic dependencies...${NC}"
  install_dependencies &
  spinner
fi

PROJECT_PATHS=$(prompt_user "Enter the path for your projects, it will create the folders for you if they don't exist" "/your/work/path/")
ensure_directory_exists "$PROJECT_PATHS" "false"

OBSIDIAN_PATH=$(prompt_user "Enter the path for your Obsidian vault, it will create the folders for you if they don't exist" "/your/notes/path")
ensure_directory_exists "$OBSIDIAN_PATH" "true"

clone_repository_with_progress "https://github.com/Gentleman-Programming/Gentleman.Dots.git" "Gentleman.Dots" 20
cd Gentleman.Dots || exit

install_homebrew

echo -e "${YELLOW}Step 2: Choose and Install Terminal Emulator${NC}"
if is_wsl; then
  echo -e "${YELLOW}You are running WSL. Terminal emulators should be installed on Windows.${NC}"
else
  term_choice=$(select_option "Which terminal emulator do you want to install? " "alacritty" "wezterm" "kitty")
  case "$term_choice" in
    "alacritty")
      install_terminal "Alacritty" "sudo apt-get install -y alacritty" "mkdir -p ~/.config/alacritty && cp alacritty.toml ~/.config/alacritty/"
      ;;
    "wezterm")
      install_terminal "WezTerm" "sudo apt-get install -y wezterm" "mkdir -p ~/.config/wezterm && cp .wezterm.lua ~/.config/wezterm/wezterm.lua"
      ;;
    "kitty")
      install_terminal "Kitty" "brew install --cask kitty" "mkdir -p ~/.config/kitty && cp -r GentlemanKitty/* ~/.config/kitty"
      ;;
    *)
      echo -e "${YELLOW}No terminal emulator will be installed or configured.${NC}"
      ;;
  esac
fi

echo -e "${YELLOW}Step 3: Choose and Install Shell${NC}"
shell_choice=$(select_option "Which shell do you want to install? " "fish" "zsh")
case "$shell_choice" in
  "fish")
    install_shell "fish" "brew install fish" "mkdir -p ~/.config/fish && cp -r GentlemanFish/* ~/.config/fish"
    ;;
  "zsh")
    install_shell "zsh" "brew install zsh" "brew install zsh-autosuggestions zsh-syntax-highlighting zsh-autocomplete"
    run_command "cp -r GentlemanZsh/.zshrc ~/"
    run_command "brew install powerlevel10k && cp -r GentlemanZsh/.p10k.zsh ~/"
    update_or_replace ~/.zshrc "export PROJECT_PATHS" "export PROJECT_PATHS=\"$PROJECT_PATHS\""
    ;;
  *)
    echo -e "${YELLOW}No shell will be installed or configured.${NC}"
    ;;
esac

echo -e "${YELLOW}Step 4: Installing Additional Dependencies...${NC}"
install_nvim=$(select_option "Do you want to install Neovim?" "Yes" "No")
if [ "$install_nvim" = "Yes" ]; then
  install_dependencies_with_progress "brew install nvim node npm git gcc fzf fd ripgrep coreutils bat curl lazygit"
  echo -e "${YELLOW}Configuring Neovim...${NC}"
  run_command "mkdir -p ~/.config/nvim"
  run_command "cp -r GentlemanNvim/nvim/* ~/.config/nvim/"
  obsidian_config_file="$HOME/.config/nvim/lua/plugins/obsidian.lua"
  if [ -f "$obsidian_config_file" ]; then
    update_or_replace "$obsidian_config_file" "/your/notes/path" "path = '$OBSIDIAN_PATH'"
  else
    echo -e "${RED}Obsidian configuration file not found at $obsidian_config_file. Please check your setup.${NC}"
  fi
fi

wm_choice=$(select_option "Which window manager do you want to install? " "tmux" "zellij" "none")
case "$wm_choice" in
  "tmux")
    install_window_manager "Tmux" "brew install tmux" "mkdir -p ~/.tmux && cp -r GentlemanTmux/.tmux/* ~/.tmux/ && cp GentlemanTmux/.tmux.conf ~/"
    ;;
  "zellij")
    install_window_manager "Zellij" "brew install zellij" "mkdir -p ~/.config/zellij && cp -r GentlemanZellij/zellij/* ~/.config/zellij/"
    ;;
  *)
    echo -e "${YELLOW}No window manager will be installed or configured.${NC}"
    ;;
esac

# Clean up: Remove the cloned repository
sudo chown -R $(whoami) $(brew --prefix)/*
echo -e "${YELLOW}Cleaning up...${NC}"
cd ..
run_command "rm -rf Gentleman.Dots"

echo -e "${GREEN}Configuration complete. Restarting shell...${NC}"
echo -e "${GREEN}If it doesn't restart, restart your computer or WSL instance😘${NC}"
exec $(which $SHELL)
