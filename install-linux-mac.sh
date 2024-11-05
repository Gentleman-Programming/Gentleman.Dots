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

sudo -v

while true; do
  sudo -n true
  sleep 60
  kill -0 "$$" || exit
done 2>/dev/null &

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

  # Check for the "templates" directory only if create_templates is true
  if [ "$create_templates" == "true" ]; then
    if [ ! -d "$dir_path/templates" ]; then
      echo -e "${YELLOW}Templates directory does not exist. Creating...${NC}"
      mkdir -p "$dir_path/templates"
      echo -e "${GREEN}Templates directory created at $dir_path/templates${NC}"
    else
      echo -e "${GREEN}Templates directory already exists at $dir_path/templates${NC}"
    fi
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
    run_command "sudo pacman -S rustup"
    run_command "rustup default stable"
    run_command ". $HOME/.cargo/env"
  else
    run_command "sudo apt-get update"
    run_command "sudo apt-get install -y build-essential curl file git"
    run_command "curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh"
    run_command ". $HOME/.cargo/env"
  fi
}

install_homebrew_with_progress() {
  local install_command="$1"

  echo -e "${YELLOW}Installing Homebrew...${NC}"

  if [ "$show_details" = "No" ]; then
    # Run installation in the background and show progress
    (eval "$install_command" &>/dev/null) &
    spinner
  else
    # Run installation normally
    eval "$install_command"
  fi
}

# Function to install Homebrew if not installed
install_homebrew() {
  if ! command -v brew &>/dev/null; then
    echo -e "${YELLOW}Homebrew is not installed. Installing Homebrew...${NC}"

    if [ "$show_details" = "No" ]; then
      # Show progress bar while installing Homebrew
      install_homebrew_with_progress "/bin/bash -c \"\$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\""
      spinner
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
      run_command "(echo 'eval \"\$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)\"' >> ~/.zshrc)"
      run_command "(echo 'eval \"\$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)\"' >> ~/.bashrc)"
      run_command "mkdir -p ~/.config/fish"
      run_command "(echo 'eval \"\$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)\"' >> ~/.config/fish/config.fish)"
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

# Ask if the user wants to see detailed output
show_details=$(select_option "Do you want to see detailed output? " "No" "Yes")

# Ask for the operating system
os_choice=$(select_option "Which operating system are you using? " "mac" "linux")

if [ "$os_choice" != "mac" ]; then
  # Install basic dependencies with progress bar
  echo -e "${YELLOW}Installing basic dependencies...${NC}"
  if [ "$show_details" = "No" ]; then
    install_dependencies &
    spinner
  else
    install_dependencies
  fi
else
  if xcode-select -p &>/dev/null; then
    echo -e "${GREEN}Xcode is already installed.${NC}"
  else
    run_command "xcode-select --install"
  fi
  run_command "curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh"
  run_command ". $HOME/.cargo/env"
fi

# Function to clone repository with progress bar
clone_repository_with_progress() {
  local repo_url="$1"
  local clone_dir="$2"
  local progress_duration=$3

  echo -e "${YELLOW}Cloning repository...${NC}"

  if [ "$show_details" = "No" ]; then
    # Run clone command in the background and show progress
    (git clone "$repo_url" "$clone_dir" &>/dev/null) &
    spinner "$progress_duration"
  else
    # Run clone command normally
    git clone "$repo_url" "$clone_dir"
  fi
}

# Step 1: Clone the Repository
echo -e "${YELLOW}Step 1: Clone the Repository${NC}"
if [ -d "Gentleman.Dots" ]; then
  echo -e "${GREEN}Repository already cloned. Overwriting...${NC}"
  rm -rf "Gentleman.Dots"
fi
clone_repository_with_progress "https://github.com/Gentleman-Programming/Gentleman.Dots.git" "Gentleman.Dots" 20
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
    spinner
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
      term_choice=$(select_option "Which terminal emulator do you want to install? " "alacritty" "wezterm" "none")
    else
      echo -e "${YELLOW}Note: Kitty is not available for Linux.${NC}"
      term_choice=$(select_option "Which terminal emulator do you want to install? " "alacritty" "wezterm" "none")
    fi
  else
    term_choice=$(select_option "Which terminal emulator do you want to install? " "alacritty" "wezterm" "kitty" "none")
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

  if [ "$term_choice" != "none" ]; then
    echo -e "${YELLOW}Iosevka Term Nerd Font is required for the terminal emulator.${NC}"
    font_installed=$(select_option "Do you already have Iosevka Term Nerd Font installed?" "Yes" "No")
    if [ "$font_installed" = "No" ]; then
      echo -e "${YELLOW}Installing Iosevka Term Nerd Font...${NC}"
      if [ "$os_choice" = "linux" ]; then
        mkdir -p ~/.local/share/fonts
        wget -O ~/.local/share/fonts/Iosevka.zip https://github.com/ryanoasis/nerd-fonts/releases/download/v2.1.0/Iosevka.zip
        unzip ~/.local/share/fonts/Iosevka.zip -d ~/.local/share/fonts/
        fc-cache -fv
      elif [ "$os_choice" = "mac" ]; then
        brew install --cask font-iosevka-term-nerd-font
      fi
      echo -e "${GREEN}Iosevka Term Nerd Font installed.${NC}"
    else
      echo -e "${GREEN}Iosevka Term Nerd Font is already installed.${NC}"
    fi
  fi
fi

# Shared Steps (macOS, Linux, or WSL)

# Function to install shell or plugins with progress bar
install_shell_with_progress() {
  local name="$1"
  local install_command="$2"

  echo -e "${YELLOW}Installing $name...${NC}"
  if [ "$show_details" = "No" ]; then
    (eval "$install_command" &>/dev/null) &
    spinner
  else
    eval "$install_command"
  fi
}

set_as_default_shell() {
  local name="$1"

  echo -e "${YELLOW}Setting default shell to $name...${NC}"
  local shell_path
  shell_path=$(which "$name") # Obtener el camino completo del shell

  if [ -n "$shell_path" ]; then
    sudo sh -c "grep -Fxq \"$shell_path\" /etc/shells || echo \"$shell_path\" >> /etc/shells"

    sudo chsh -s "$shell_path" "$USER"

    if [ "$SHELL" != "$shell_path" ]; then
      echo -e "${RED}Error: Shell did not change. Please check manually.${NC}"
      echo -e "${GREEN}Command: sudo chsh -s $shell_path \$USER ${NC}"
    else
      echo -e "${GREEN}Shell changed to $shell_path successfully.${NC}"
    fi
  else
    echo -e "${RED}Shell $name not found.${NC}"
  fi
}

echo -e "${YELLOW}Step 3: Choose and Install Shell${NC}"
shell_choice=$(select_option "Which shell do you want to install? " "fish" "zsh" "nushell")

# Case for shell choice
case "$shell_choice" in
"nushell")
  run_command "cp -rf bash-env-json ~/.config/"
  run_command "cp -rf bash-env.nu ~/.config/"

  if ! command -v nu &>/dev/null; then
    install_shell_with_progress "nushell" "brew install nushell carapace zoxide atuin jq bash"
  else
    echo -e "${GREEN}Nushell shell is already installed.${NC}"
  fi
  if ! command -v starship &>/dev/null; then
    install_shell_with_progress "starship" "brew install starship"
  else
    echo -e "${GREEN}starship is already installed.${NC}"
  fi

  [ ! -d ~/.cache/starship ] && mkdir ~/.cache/starship
  [ ! -d ~/.cache/carapace ] && mkdir ~/.cache/carapace
  [ ! -d ~/.local/share/atuin ] && mkdir ~/.local/share/atuin

  echo -e "${YELLOW}Configuring Nushell...${NC}"

  if [[ "$OSTYPE" == "darwin"* ]]; then
    mkdir -p ~/Library/Application\ Support/nushell
    update_or_replace "GentlemanNushell/env.nu" "/home/linuxbrew/.linuxbrew/bin" "    | prepend '/opt/homebrew/bin'"
    run_command "cp -rf GentlemanNushell/* ~/Library/Application\ Support/nushell/"
  else
    mkdir -p ~/.config/nushell
    run_command "cp -rf GentlemanNushell/* ~/.config/nushell/"
  fi

  ;;
"fish")
  if ! command -v fish &>/dev/null; then
    install_shell_with_progress "fish" "brew install fish carapace zoxide atuin"
  else
    echo -e "${GREEN}Fish shell is already installed.${NC}"
  fi
  if ! command -v starship &>/dev/null; then
    install_shell_with_progress "starship" "brew install starship"
  else
    echo -e "${GREEN}starship is already installed.${NC}"
  fi

  [ ! -d ~/.cache/starship ] && mkdir ~/.cache/starship
  [ ! -d ~/.cache/carapace ] && mkdir ~/.cache/carapace
  [ ! -d ~/.local/share/atuin ] && mkdir ~/.local/share/atuin

  echo -e "${YELLOW}Configuring Fish...${NC}"
  run_command "cp -rf GentlemanFish/fish ~/.config"
  ;;
"zsh")
  if ! command -v zsh &>/dev/null; then
    install_shell_with_progress "zsh" "brew install zsh carapace zoxide atuin"
  else
    echo -e "${GREEN}zsh is already installed.${NC}"
  fi

  if ! command -v zsh-autosuggestions &>/dev/null || ! command -v zsh-syntax-highlighting &>/dev/null || ! command -v zsh-autocomplete &>/dev/null; then
    install_shell_with_progress "zsh" "brew install zsh-autosuggestions zsh-syntax-highlighting zsh-autocomplete" ""
  fi

  if [ ! -d "$HOME/.oh-my-zsh" ]; then
    echo -e "${YELLOW}Installing Oh My Zsh...${NC}"
    if [ "$show_details" = "No" ]; then
      (
        NO_INTERACTIVE=true sh -c "$(curl -fsSL https://raw.githubusercontent.com/subtlepseudonym/oh-my-zsh/feature/install-noninteractive/tools/install.sh)" &>/dev/null
      ) &
      spinner
    else
      NO_INTERACTIVE=true sh -c "$(curl -fsSL https://raw.githubusercontent.com/subtlepseudonym/oh-my-zsh/feature/install-noninteractive/tools/install.sh)"
    fi
  else
    echo -e "${GREEN}Oh My Zsh is already installed.${NC}"
  fi

  echo -e "${YELLOW}Configuring Zsh...${NC}"
  run_command "cp -rf GentlemanZsh/.zshrc ~/"

  [ ! -d ~/.cache/starship ] && mkdir ~/.cache/starship
  [ ! -d ~/.cache/carapace ] && mkdir ~/.cache/carapace
  [ ! -d ~/.local/share/atuin ] && mkdir ~/.local/share/atuin

  # PowerLevel10K Configuration
  echo -e "${YELLOW}Configuring PowerLevel10K...${NC}"
  run_command "brew install powerlevel10k"
  run_command "cp -rf GentlemanZsh/.p10k.zsh ~/"
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
    spinner
  else
    # Run installation normally
    eval "$install_command"
  fi
}

# Step 5: Additional Configurations

# Dependencies Install
echo -e "${YELLOW}Step 4: Installing Additional Dependencies...${NC}"

if [ "$os_choice" = "linux" ]; then
  if ! is_arch; then
    # Combine the update and upgrade commands for progress (only if not Arch Linux)
    install_dependencies_with_progress "sudo apt-get update && sudo apt-get upgrade -y"
  fi
fi

# Function to install window manager with progress bar
install_window_manager_with_progress() {
  local install_command="$1"

  echo -e "${YELLOW}Installing window manager...${NC}"

  if [ "$show_details" = "No" ]; then
    # Run installation in the background and show progress
    (eval "$install_command" &>/dev/null) &
    spinner
  else
    # Run installation normally
    eval "$install_command"
  fi
}

# Ask if they want to use Tmux or Zellij, or none
echo -e "${YELLOW}Step 4: Choose and Install Window Manager${NC}"
wm_choice=$(select_option "Which window manager do you want to install? " "tmux" "zellij" "none")

case "$wm_choice" in
"tmux")
  if ! command -v tmux &>/dev/null; then
    if [ "$show_details" = "Yes" ]; then
      install_window_manager_with_progress "brew install tmux"
    else
      run_command "brew install tmux"
    fi
  else
    echo -e "${GREEN}Tmux is already installed.${NC}"
  fi

  echo -e "${YELLOW}Configuring Tmux...${NC}"
  if [ ! -d "$HOME/.tmux/plugins/tpm" ]; then
    if [ "$show_details" = "Yes" ]; then
      install_window_manager_with_progress "git clone https://github.com/tmux-plugins/tpm ~/.tmux/plugins/tpm"
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
  if [ "$show_details" = "Yes" ]; then
    install_window_manager_with_progress "brew uninstall zellij | cargo install zellij"
  else
    run_command "cargo install zellij"
  fi

  echo -e "${YELLOW}Configuring Zellij...${NC}"
  run_command "mkdir -p ~/.config/zellij"
  run_command "cp -r GentlemanZellij/zellij/* ~/.config/zellij/"

  # Replace TMUX with ZELLIJ and tmux with zellij only in the selected shell configuration
  if [[ "$shell_choice" == "zsh" ]]; then
    update_or_replace ~/.zshrc "TMUX" 'WM_VAR="/$ZELLIJ"'
    update_or_replace ~/.zshrc "tmux" 'WM_CMD="zellij"'
  elif [[ "$shell_choice" == "fish" ]]; then
    update_or_replace ~/.config/fish/config.fish "TMUX" "if not set -q ZELLIJ"
    update_or_replace ~/.config/fish/config.fish "tmux" "zellij"
  elif [[ "$shell_choice" == "nushell" ]]; then
    os_type=$(uname)

    if [[ "$os_type" == "Darwin" ]]; then
      update_or_replace "GentlemanNushell/config.nu" '"tmux"' 'let MULTIPLEXER = "zellij"'
      update_or_replace "GentlemanNushell/config.nu" '"TMUX"' 'let MULTIPLEXER_ENV_PREFIX = "ZELLIJ"'
      run_command "cp -rf GentlemanNushell/* ~/Library/Application\ Support/nushell/"
    else
      update_or_replace ~/.config/nushell/config.nu '"tmux"' 'let MULTIPLEXER = "zellij"'
      update_or_replace ~/.config/nushell/config.nu '"TMUX"' 'let MULTIPLEXER_ENV_PREFIX = "ZELLIJ"'
    fi
  fi
  ;;
"none")
  echo -e "${YELLOW}No window manager will be installed or configured.${NC}"
  # If no window manager is chosen, remove the line that executes tmux or zellij

  # Determine the OS type
  OS_TYPE=$(uname)

  # Function to run sed with appropriate options based on OS
  run_sed() {
    local file=$1
    local pattern=$2

    if [ "$OS_TYPE" = "Darwin" ]; then
      # macOS
      sed -i '' "$pattern" "$file"
    else
      # Linux and other Unix-like systems
      sed -i "$pattern" "$file"
    fi
  }

  # Check and modify ~/.zshrc if it exists
  if [ -f ~/.zshrc ]; then
    run_sed ~/.zshrc '/exec $WM_CMD/d'
    update_or_replace ~/.zshrc "exec $WM_CMD" "true"
  fi

  # Check and modify ~/.config/fish/config.fish if it exists
  if [ -f ~/.config/fish/config.fish ]; then
    update_or_replace ~/.config/fish/config.fish "tmux" "true"
    update_or_replace ~/.config/fish/config.fish "zellij" "true"
  fi

  # Check and modify ~/.config/nushell/config.nu if it exists
  if [ -f ~/.config/fish/config.fish ] || [ -d ~/Library/Application\ Support/nushell ]; then
    if [[ "$os_type" == "Darwin" ]]; then
      update_or_replace "GentlemanNushell/config.nu" "/run-external $MULTIPLEXER" "true"
    else
      update_or_replace ~/.config/nushell/config.nu "/run-external $MULTIPLEXER" "true"
    fi
  fi
  ;;
*)
  echo -e "${YELLOW}Invalid option. No window manager will be installed or configured.${NC}"
  ;;
esac

# Neovim Configuration
echo -e "${YELLOW}Step 5: Choose and Install NVIM${NC}"
install_nvim=$(select_option "Do you want to install Neovim?" "Yes" "No")

if [ "$install_nvim" = "Yes" ]; then
  OBSIDIAN_PATH=$(prompt_user "Enter the path for your Obsidian vault, it will create the folders for you if they don't exist" "/your/notes/path")
  ensure_directory_exists "$OBSIDIAN_PATH" "true"

  # Install additional packages with Neovim
  install_dependencies_with_progress "brew install nvim node npm git gcc fzf fd ripgrep coreutils bat curl lazygit"

  # Neovim Configuration
  echo -e "${YELLOW}Configuring Neovim...${NC}"
  run_command "mkdir -p ~/.config/nvim"
  run_command "cp -r GentlemanNvim/nvim/* ~/.config/nvim/"
  # Obsidian Configuration
  echo -e "${YELLOW}Configuring Obsidian...${NC}"
  obsidian_config_file="$HOME/.config/nvim/lua/plugins/obsidian.lua"

  if [ -f "$obsidian_config_file" ]; then
    # Replace the vault path in the existing configuration
    update_or_replace "$obsidian_config_file" "/your/notes/path" "path = '$OBSIDIAN_PATH'"
  else
    echo -e "${RED}Obsidian configuration file not found at $obsidian_config_file. Please check your setup.${NC}"
  fi
fi

# Clean up: Remove the cloned repository
sudo chown -R $(whoami) $(brew --prefix)/*
echo -e "${YELLOW}Cleaning up...${NC}"
cd ..
run_command "rm -rf Gentleman.Dots"

if [ "$shell_choice" = "nushell" ]; then
  shell_choice="nu"
fi

set_as_default_shell "$shell_choice"

echo -e "${GREEN}Configuration complete. Restarting shell...${NC}"
echo -e "${GREEN}If it doesn't restart, restart your computer or WSL instance??${NC}"

exec $shell_choice
