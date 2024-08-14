export ZSH="$HOME/.oh-my-zsh"

ZSH_THEME="robbyrussell"

# Si la sesión es interactiva
if [[ $- == *i* ]]; then
    # Commands to run in interactive sessions can go here
fi

## Para macOS
eval "$(/opt/homebrew/bin/brew shellenv)"

## Para Linux
#eval "$(/home/linuxbrew/.linuxbrew/bin/brew shellenv)"

export PROJECT_PATHS="/your/work/path/"

# Iniciar tmux si la sesión es interactiva y no estamos ya en tmux
if [[ $- == *i* ]] && [[ -z "$TMUX" ]]; then
    exec tmux
fi

# bun
export BUN_INSTALL="$HOME/.bun"
export PATH="$BUN_INSTALL/bin:$PATH"


#plugins
plugins=(
  pj 
  command-not-found     
  zsh-autosuggestions
)
source /opt/homebrew/share/zsh-autocomplete/zsh-autocomplete.plugin.zsh
source /opt/homebrew/share/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh

source $ZSH/oh-my-zsh.sh

# Inicializar Starship para zsh
eval "$(starship init zsh)"
