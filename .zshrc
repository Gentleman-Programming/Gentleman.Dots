export ZSH="$HOME/.oh-my-zsh"

ZSH_THEME="robbyrussell"

# Si la sesión es interactiva
if [[ $- == *i* ]]; then
    # Commands to run in interactive sessions can go here
fi

if [[ "$(uname)" == "Darwin" ]]; then
    # macOS
    BREW_BIN="/opt/homebrew/bin"
else
    # Linux
    BREW_BIN="/home/linuxbrew/.linuxbrew/bin"
fi

# Usar la variable BREW_BIN donde se necesite
eval "$($BREW_BIN/brew shellenv)"

source $(dirname $BREW_BIN)/share/zsh-autocomplete/zsh-autocomplete.plugin.zsh
source $(dirname $BREW_BIN)/share/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh
source $(dirname $BREW_BIN)/share/zsh-autosuggestions/zsh-autosuggestions.zsh

export PROJECT_PATHS="/home/alanbuscaglia/work"

# Run Tmux
if [[ $- == *i* ]] && [[ -z "$TMUX" ]]; then
    exec tmux
fi

# Run Zellij
#if [[ $- == *i* ]] && [[ -z "$ZELLIJ" ]]; then
#    exec zellij
#fi


# bun
export BUN_INSTALL="$HOME/.bun"
export PATH="$BUN_INSTALL/bin:$PATH"


#plugins
plugins=(
  pj 
  command-not-found     
)

source $ZSH/oh-my-zsh.sh

eval "$(fzf --zsh)"

# Inicializar Starship para zsh
eval "$(starship init zsh)"

