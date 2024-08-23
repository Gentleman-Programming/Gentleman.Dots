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
export FZF_DEFAULT_COMMAND="fd --hidden --strip-cwd-prefix --exclude .git"
export FZF_DEFAULT_T_COMMAND="$FZF_DEFAULT_COMMAND"
export FZF_ALT_COMMAND="fd --type=d --hidden --strip-cwd-prefix --exlude .git"

WM_VAR="\$TMUX" // change with ZELLIJ
WM_CMD="tmux" // change with zellij

if [[ $- == *i* ]] && [[ -z "$WM_VAR" ]]; then
    exec $WM_CMD
fi

# alias
alias fzfbat='fzf --preview="bat --theme=gruvbox-dark --color=always {}"'
alias fzfnvim='nvim $(fzf --preview="bat --theme=gruvbox-dark --color=always {}")'

#plugins
plugins=(
  pj 
  command-not-found     
)

source $ZSH/oh-my-zsh.sh

eval "$(fzf --zsh)"

# Inicializar Starship para zsh
eval "$(starship init zsh)"

