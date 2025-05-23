export TZ="America/Santo_Domingo"

# Enable Powerlevel10k instant prompt. Should stay close to the top of ~/.zshrc.
# Initialization code that may require console input (password prompts, [y/n]
# confirmations, etc.) must go above this block; everything else may go below.
if [[ -r "${XDG_CACHE_HOME:-$HOME/.cache}/p10k-instant-prompt-${(%):-%n}.zsh" ]]; then
  source "${XDG_CACHE_HOME:-$HOME/.cache}/p10k-instant-prompt-${(%):-%n}.zsh"
fi

export ZSH="$HOME/.oh-my-zsh"
export PATH="$HOME/.cargo/bin:$PATH"
#
#LOCAL BIN
export LOCAL_PATH="/home/deuri-vasquez/.local"
export PATH="$LOCAL_PATH/bin:$PATH"

if [[ $- == *i* ]]; then
    # Commands to run in interactive sessions can go here
fi

export LS_COLORS="di=38;5;67:ow=48;5;60:ex=38;5;132:ln=38;5;144:*.tar=38;5;180:*.zip=38;5;180:*.jpg=38;5;175:*.png=38;5;175:*.mp3=38;5;175:*.wav=38;5;175:*.txt=38;5;223:*.sh=38;5;132"
if [[ "$(uname)" == "Darwin" ]]; then
  alias ls='ls --color=auto'
else
  alias ls='gls --color=auto'
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
source $(dirname $BREW_BIN)/share/powerlevel10k/powerlevel10k.zsh-theme

export PROJECT_PATHS="/home/deuri-vasquez/projects"
export FZF_DEFAULT_COMMAND="fd --hidden --strip-cwd-prefix --exclude .git"
export FZF_DEFAULT_T_COMMAND="$FZF_DEFAULT_COMMAND"
export FZF_ALT_COMMAND="fd --type=d --hidden --strip-cwd-prefix --exlude .git"

# ZELLIJ initialization
WM_VAR="/$ZELLIJ"
WM_CMD="zellij"

function start_if_needed() {
    if [[ $- == *i* ]] && [[ -z "${WM_VAR#/}" ]] && [[ -t 1 ]]; then
        exec $WM_CMD
    fi
}

alias zz='zellij'

# Lazy load nvim
nvim() {
    unfunction nvim
    # Load nvim only when first called
    if command -v nvim &>/dev/null; then
        nvim "$@"
    else
        echo "Neovim is not installed"
        return 1
    fi
}

alias vim='nvim'
alias vi='nvim'

# alias
alias fzfbat='fzf --preview="bat --theme=gruvbox-dark --color=always {}"'
alias fzfnvim='nvim $(fzf --preview="bat --theme=gruvbox-dark --color=always {}")'

#plugins
plugins=(
  command-not-found
)

source $ZSH/oh-my-zsh.sh

export CARAPACE_BRIDGES='zsh,fish,bash,inshellisense'
zstyle ':completion:*' format $'\e[2;37mCompleting %d\e[m'
source <(carapace _carapace)

eval "$(fzf --zsh)"
eval "$(zoxide init zsh)"
eval "$(atuin init zsh)"

# To customize prompt, run `p10k configure` or edit ~/.p10k.zsh.
[[ ! -f ~/.p10k.zsh ]] || source ~/.p10k.zsh


# Start zellij if not already running
# start_if_needed

# proto
export PROTO_HOME="$HOME/.proto";
export PATH="$PROTO_HOME/shims:$PROTO_HOME/bin:$PATH";

# Docker
export COMPOSE_BAKE=true
alias lzd='lazydocker'

#DOTNET
# export DOTNET_ROOT='/home/linuxbrew/.linuxbrew/bin/dotnet'

#BUN
export PATH="/home/deuri-vasquez/.bun/bin:$PATH"

# AI config
export GEMINI_API_KEY='AIzaSyCd6OWpoAUKBMlQFP_yCjjVwG3uuAI0UnU'

