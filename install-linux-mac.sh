#!/bin/bash

set -e

# Define colors for output using tput for better compatibility
PINK=$(tput setaf 204)
PURPLE=$(tput setaf 141)
GREEN=$(tput setaf 114)
ORANGE=$(tput setaf 208)
BLUE=$(tput setaf 75)
YELLOW=$(tput setaf 221)
RED=$(tput setaf 196)
NC=$(tput sgr0) # No Color

# Gentleman.Dots logo with pink color
logo='
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

# Display logo and title
echo -e "${PINK}${logo}${NC}"
echo -e "${PURPLE}Welcome to the Gentleman.Dots Auto Config!${NC}"

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

# Function to display a progress bar
progress_bar() {
  local duration=$1
  already_done() { for ((done = 0; done < $elapsed; done++)); do printf "▇"; done; }
  remaining() { for ((remain = $elapsed; remain < $duration; remain++)); do printf " "; done; }
  percentage() { printf "| %s%%" $(((($elapsed) * 100) / ($duration) * 100 / 100)); }
  for ((elapsed = 1; elapsed <= $duration; elapsed += 1)); do
    already_done
    remaining
    percentage
    printf "\r"
    sleep 0.1
  done
  printf "\n"
}

# Function to check and create directories if they do not exist
ensure_directory_exists() {
  local dir_path="$1"
  local create_templates="$2"
  if [ ! -d "$dir_path" ]; then
    echo -e "${YELLOW}Directory $dir_path does not exist. Creating...${NC}"
    mkdir -p "$dir_path"
    if [ "$create_templates" == "true" ]; then
      mkdir -p "$dir_path/templates"
      echo -e "${GREEN}Templates directory created at $dir_path/templates${NC}"
    fi
  else
    echo -e "${GREEN}Directory $dir_path already exists.${NC}"
  fi
}

# Function to check if running on WSL
is_wsl() {
  grep -qEi "(Microsoft|WSL)" /proc/version &>/dev/null
  return $?
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
  if [ -f /etc/arch-release ]; then
    return 0
  else
    return 1
  fi
}

# Function to install basic dependencies
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

    if [ "$show_details" = "No" ]; then
      # Show progress bar while installing Homebrew
      install_homebrew_progress &
      progress_bar 10
    else
      # Install Homebrew normally
      run_command "/bin/bash -c \"\$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\""
    fi

    # Add Homebrew to PATH based on OS
    if [ "$os_choice" = "mac" ]; then
      run_command "(echo 'eval \"\$(/opt/homebrew/bin/brew shellenv)\"' >> $USER_HOME/.zshrc)"
      run_command "(echo 'eval \"\$(/opt/homebrew/bin/brew shellenv)\"' >> $USER_HOME/.bashrc)"
      run_command "mkdir -p $USER_HOME/.config/fish"
      run_command "(echo 'eval \"\$(/opt/homebrew/bin/brew shellenv)\"' >> $USER_HOME/.config/fish/config.fish)"
      run_command "eval \"\$(/opt/homebrew/bin/brew shellenv)\""
    elif [ "$os_choice" = "linux" ]; then
      run_command "(echo 'eval \"\$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)\"' >> $USER_HOME/.zshrc)"
      run_command "(echo 'eval \"\$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)\"' >> $USER_HOME/.bashrc)"
      run_command "mkdir -p $USER_HOME/.config/fish"
      run_command "(echo 'eval \"\$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)\"' >> $USER_HOME/.config/fish/config.fish)"
      run_command "eval \"\$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)\""
    fi
  else
    echo -e "${GREEN}Homebrew is already installed.${NC}"
  fi
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

  if ! grep -Fxq "$shell_path" /etc/shells; then
    echo "$shell_path" | sudo tee -a /etc/shells
  fi

  sudo chsh -s "$shell_path" "$USER"
}

# Ask if the user wants to see detailed output
show_details=$(select_option "Do you want to see detailed output? " "No" "Yes")

# Ask for the operating system
os_choice=$(select_option "Which operating system are you using? " "mac" "linux")

# Install basic dependencies with progress bar
echo -e "${YELLOW}Installing basic dependencies...${NC}"
if [ "$show_details" = "No" ]; then
  install_dependencies &
  progress_bar 10
else
  install_dependencies
fi

# Prompt for project path and Obsidian path
PROJECT_PATHS=$(prompt_user "Enter the path for your projects, it will create the folders for you if they don't exist" "/your/work/path/")
ensure_directory_exists "$PROJECT_PATHS" "false"

OBSIDIAN_PATH=$(prompt_user "Enter the path for your Obsidian vault, it will create the folders for you if they don't exist" "/your/notes/path")
ensure_directory_exists "$OBSIDIAN_PATH" "true"

# Function to clone repository with progress bar
clone_repository_with_progress() {
  local repo_url="$1"
  local clone_dir="$2"
  local progress_duration=$3

  echo -e "${YELLOW}Cloning repository...${NC}"

  if [ "$show_details" = "No" ]; then
    # Run clone command in the background and show progress
    (git clone "$repo_url" "$clone_dir" &>/dev/null) &
    progress_bar "$progress_duration"
  else
    # Run clone command normally
    git clone "$repo_url" "$clone_dir"
  fi
}

# Step 1: Clone the Repository
echo -e "${YELLOW}Step 1: Clone the Repository${NC}"
if [ -d "Gentleman.Dots" ]; then
  echo -e "${GREEN}Repository already cloned. Skipping...${NC}"
else
  clone_repository_with_progress "https://github.com/Gentleman-Programming/Gentleman.Dots.git" "Gentleman.Dots" 20
fi
cd Gentleman.Dots || exit

# Install Homebrew if not installed
install_homebrew

# Function to install a terminal emulator with progress
install_terminal_with_progress() {
  local term_name="$1"
  local install_command="$2"
  local config_command="$3"

  echo -e "${YELLOW}Installing $term_name...${NC}"

  if [ "$show_details" = "No" ]; then
    # Run installation in the background and show progress
    (eval "$install_command" &>/dev/null) &
    progress_bar 10
  else
    # Run installation normally
    eval "$install_command"
  fi

  echo -e "${YELLOW}Configuring $term_name...${NC}"
  eval "$config_command"
}

echo -e "${YELLOW}Step 2: Choose and Install Terminal Emulator${NC}"
if is_wsl; then
  echo -e "${YELLOW}You are running WSL. Terminal emulators should be installed on Windows.${NC}"
else
  if [ "$os_choice" = "linux" ]; then
    if is_arch; then
      term_choice=$(select_option "Which terminal emulator do you want to install? " "alacritty" "wezterm")
    else
      echo -e "${YELLOW}Note: Kitty is not available for Linux.${NC}"
      term_choice=$(select_option "Which terminal emulator do you want to install? " "alacritty" "wezterm")
    fi
  else
    term_choice=$(select_option "Which terminal emulator do you want to install? " "alacritty" "wezterm" "kitty")
  fi

  case "$term_choice" in
  "alacritty")
    if ! command -v alacritty &>/dev/null; then
      if is_arch; then
        install_terminal_with_progress "Alacritty" "sudo pacman -S --noconfirm alacritty" "mkdir -p ~/.config/alacritty && cp alacritty.toml ~/.config/alacritty/alacritty.toml"
      else
        install_terminal_with_progress "Alacritty" "sudo add-apt-repository ppa:aslatter/ppa && sudo apt-get update && sudo apt-get install alacritty" "mkdir -p ~/.config/alacritty && cp alacritty.toml ~/.config/alacritty/alacritty.toml"
      fi
    else
      echo -e "${GREEN}Alacritty is already installed.${NC}"
    fi
    ;;
  "wezterm")
    if ! command -v wezterm &>/dev/null; then
      if is_arch; then
        install_terminal_with_progress "WezTerm" "sudo pacman -S --noconfirm wezterm" "mkdir -p ~/.config/wezterm && cp .wezterm.lua ~/.config/wezterm/wezterm.lua"
      else
        install_terminal_with_progress "WezTerm" "curl -fsSL https://apt.fury.io/wez/gpg.key | sudo gpg --yes --dearmor -o /usr/share/keyrings/wezterm-fury.gpg && echo 'deb [signed-by=/usr/share/keyrings/wezterm-fury.gpg] https://apt.fury.io/wez/ * *' | sudo tee /etc/apt/sources.list.d/wezterm.list && sudo apt update && sudo apt install wezterm" "mkdir -p ~/.config/wezterm && cp .wezterm.lua ~/.config/wezterm/wezterm.lua"
      fi
    else
      echo -e "${GREEN}WezTerm is already installed.${NC}"
    fi
    ;;
  "kitty")
    if [ "$os_choice" = "mac" ]; then
      if ! command -v kitty &>/dev/null; then
        install_terminal_with_progress "Kitty" "brew install --cask kitty" "mkdir -p ~/.config/kitty && cp -r GentlemanKitty/* ~/.config/kitty"
      else
        echo -e "${GREEN}Kitty is already installed.${NC}"
      fi
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

# Function to install shell or plugins with progress bar
install_shell_with_progress() {
  local name="$1"
  local install_command="$2"
  local post_install_command="$3"
  local set_default_command="$4"

  echo -e "${YELLOW}Installing $name...${NC}"
  if [ "$show_details" = "No" ]; then
    (eval "$install_command" &>/dev/null) &
    progress_bar 20
  else
    eval "$install_command"
  fi

  if [ -n "$post_install_command" ]; then
    echo -e "${YELLOW}Post-install configuration for $name...${NC}"
    eval "$post_install_command"
  fi

  if [ -n "$set_default_command" ]; then
    echo -e "${YELLOW}Setting default shell to $name...${NC}"
    local shell_path=$(which $name)  # Obtener el camino completo del shell
    set_default_shell "$shell_path"
  fi
}

echo -e "${YELLOW}Step 3: Choose and Install Shell${NC}"
shell_choice=$(select_option "Which shell do you want to install? " "fish" "zsh")

# Case for shell choice
case "$shell_choice" in
"fish")
  if ! command -v fish &>/dev/null; then
    install_shell_with_progress "Fish shell" "brew install fish" "mkdir -p ~/.config/fish && cp -r GentlemanFish/* ~/.config" "set_default_shell \"$(which fish)\""
  else
    echo -e "${GREEN}Fish shell is already installed.${NC}"
  fi
  ;;
"zsh")
  if ! command -v zsh &>/dev/null; then
    install_shell_with_progress "Zsh" "brew install zsh" "" "set_default_shell \"$(which zsh)\""
  else
    echo -e "${GREEN}Zsh is already installed.${NC}"
  fi

  if ! command -v zsh-autosuggestions &>/dev/null || ! command -v zsh-syntax-highlighting &>/dev/null || ! command -v zsh-autocomplete &>/dev/null; then
    install_shell_with_progress "Zsh plugins" "brew install zsh-autosuggestions zsh-syntax-highlighting zsh-autocomplete" "" ""
  fi

  if [ ! -d "$HOME/.oh-my-zsh" ]; then
    echo -e "${YELLOW}Installing Oh My Zsh...${NC}"
    if [ "$show_details" = "No" ]; then
      (
        NO_INTERACTIVE=true sh -c "$(curl -fsSL https://raw.githubusercontent.com/subtlepseudonym/oh-my-zsh/feature/install-noninteractive/tools/install.sh)" &>/dev/null
      ) &
      progress_bar 30
    else
      NO_INTERACTIVE=true sh -c "$(curl -fsSL https://raw.githubusercontent.com/subtlepseudonym/oh-my-zsh/feature/install-noninteractive/tools/install.sh)"
    fi
  else
    echo -e "${GREEN}Oh My Zsh is already installed.${NC}"
  fi

  echo -e "${YELLOW}Configuring Zsh...${NC}"
  run_command "cp .zshrc ~/"

  # Update or append the PROJECT_PATHS line
  update_or_replace ~/.zshrc "export PROJECT_PATHS" "export PROJECT_PATHS=\"$PROJECT_PATHS\""
  ;;
*)
  echo -e "${YELLOW}No shell will be installed or configured.${NC}"
  ;;
esac

# Function to install dependencies with progress bar
install_dependencies_with_progress() {
  local install_command="$1"

  echo -e "${YELLOW}Installing dependencies...${NC}"

  if [ "$show_details" = "No" ]; then
    # Run installation in the background and show progress
    (eval "$install_command" &>/dev/null) &
    progress_bar 10
  else
    # Run installation normally
    eval "$install_command"
  fi
}

# Step 4: Additional Configurations

# Dependencies Install
echo -e "${YELLOW}Step 4: Installing Additional Dependencies...${NC}"

if [ "$os_choice" = "linux" ]; then
  if ! is_arch; then
    # Combine the update and upgrade commands for progress (only if not Arch Linux)
    install_dependencies_with_progress "sudo apt-get update && sudo apt-get upgrade -y"
  fi
fi

# Install additional packages with progress
install_dependencies_with_progress "brew install nvim starship node npm git gcc fzf fd ripgrep coreutils bat curl lazygit"

# Neovim Configuration
echo -e "${YELLOW}Configuring Neovim...${NC}"
run_command "mkdir -p ~/.config/nvim"
run_command "cp -r GentlemanNvim/nvim/* ~/.config/nvim/"

# Starship Configuration
echo -e "${YELLOW}Configuring Starship...${NC}"
run_command "mkdir -p ~/.config"
run_command "cp starship.toml ~/.config"

# Obsidian Configuration
echo -e "${YELLOW}Configuring Obsidian...${NC}"
obsidian_config_file="$HOME/.config/nvim/lua/plugins/obsidian.lua"
if [ -f "$obsidian_config_file" ]; then
  # Replace the vault path in the existing configuration
  update_or_replace "$obsidian_config_file" "/your/notes/path" "path = '$OBSIDIAN_PATH'"
else
  echo -e "${RED}Obsidian configuration file not found at $obsidian_config_file. Please check your setup.${NC}"
fi

# Function to install window manager with progress bar
install_window_manager_with_progress() {
  local install_command="$1"
  local progress_duration=$2

  echo -e "${YELLOW}Installing window manager...${NC}"

  if [ "$show_details" = "No" ]; then
    # Run installation in the background and show progress
    (eval "$install_command" &>/dev/null) &
    progress_bar "$progress_duration"
  else
    # Run installation normally
    eval "$install_command"
  fi
}

# Ask if they want to use Tmux or Zellij, or none
wm_choice=$(select_option "Which window manager do you want to install? " "tmux" "zellij" "none")

case "$wm_choice" in
"tmux")
  if ! command -v tmux &>/dev/null; then
    if [ "$show_details" = "Yes" ]; then
      install_window_manager_with_progress "brew install tmux" 10
    else
      run_command "brew install tmux"
    fi
  else
    echo -e "${GREEN}Tmux is already installed.${NC}"
  fi

  echo -e "${YELLOW}Configuring Tmux...${NC}"
  if [ ! -d "$HOME/.tmux/plugins/tpm" ]; then
    if [ "$show_details" = "Yes" ]; then
      install_window_manager_with_progress "git clone https://github.com/tmux-plugins/tpm ~/.tmux/plugins/tpm" 10
    else
      run_command "git clone https://github.com/tmux-plugins/tpm ~/.tmux/plugins/tpm"
    fi
  else
    echo -e "${GREEN}Tmux Plugin Manager is already installed.${NC}"
  fi

  run_command "mkdir -p ~/.tmux"
  run_command "cp -r GentlemanTmux/.tmux/* ~/.tmux/"
  run_command "cp GentlemanTmux/.tmux.conf ~/"

  echo -e "${YELLOW}Installing Tmux plugins...${NC}"
  SESSION_NAME="plugin-installation"

  # Check if session already exists and kill it if necessary
  if tmux has-session -t $SESSION_NAME 2>/dev/null; then
    echo -e "${YELLOW}Session $SESSION_NAME already exists. Killing it...${NC}"
    tmux kill-session -t $SESSION_NAME
  fi

  # Create a new session in detached mode with the specified name
  tmux new-session -d -s $SESSION_NAME 'source ~/.tmux.conf; tmux run-shell ~/.tmux/plugins/tpm/bin/install_plugins'

  # Check if the user wants to see details
  if [ "$show_details" = "Yes" ]; then
    # Use a loop to show progress (adjust as needed)
    while tmux has-session -t $SESSION_NAME 2>/dev/null; do
      echo -n "."
      sleep 1
    done
    echo -e "\n${GREEN}Tmux plugins installation complete!${NC}"
  else
    # Wait for a few seconds to ensure the installation completes
    while tmux has-session -t $SESSION_NAME 2>/dev/null; do
      sleep 1
    done

    echo -e "${GREEN}Tmux plugins installation complete!${NC}"
  fi

  # Ensure the tmux session is killed
  if tmux has-session -t $SESSION_NAME 2>/dev/null; then
    tmux kill-session -t $SESSION_NAME
  fi
  ;;
"zellij")
  if ! command -v zellij &>/dev/null; then
    install_window_manager_with_progress "brew install zellij" 10
  else
    echo -e "${GREEN}Zellij is already installed.${NC}"
  fi
  echo -e "${YELLOW}Configuring Zellij...${NC}"
  run_command "mkdir -p ~/.config/zellij"
  run_command "cp -r GentlemanZellij/zellij/* ~/.config/zellij/"

  # Replace TMUX with ZELLIJ and tmux with zellij only in the selected shell configuration
  if [[ "$shell_choice" == "zsh" ]]; then
    update_or_replace ~/.zshrc "TMUX" 'if [[ $- == *i* ]] && [[ -z "\$ZELLIJ" ]]; then'
    update_or_replace ~/.zshrc "exec tmux" "exec zellij"
  elif [[ "$shell_choice" == "fish" ]]; then
    update_or_replace ~/.config/fish/config.fish "TMUX" "if not set -q ZELLIJ"
    update_or_replace ~/.config/fish/config.fish "tmux" "zellij"
  fi
  ;;
"none")
  echo -e "${YELLOW}No window manager will be installed or configured.${NC}"
  # If no window manager is chosen, remove the line that executes tmux or zellij
  sed -i '' '/exec tmux/d' ~/.zshrc
  sed -i '' '/exec zellij/d' ~/.zshrc
  sed -i '' '/tmux/d' ~/.config/fish/config.fish
  sed -i '' '/zellij/d' ~/.config/fish/config.fish
  ;;
*)
  echo -e "${YELLOW}Invalid option. No window manager will be installed or configured.${NC}"
  ;;
esac

# Clean up: Remove the cloned repository
echo -e "${YELLOW}Cleaning up...${NC}"
cd ..
run_command "rm -rf Gentleman.Dots"

echo -e "${GREEN}Configuration complete. Restarting shell...${NC}"
echo -e "${GREEN}If it doesn't work, restart your computer or WSL instance😘${NC}"
# Para Bash o Zsh
exec $SHELL
