{ config, pkgs, ... }:
{
  programs.zsh.enable = true;
  home.file = {
    ".zshrc" = {
      text = ''
if [[ -r "${XDG_CACHE_HOME:-$HOME/.cache}/p10k-instant-prompt-${(%):-%n}.zsh" ]]; then
  source "${XDG_CACHE_HOME:-$HOME/.cache}/p10k-instant-prompt-${(%):-%n}.zsh"
fi

export ZSH="$HOME/.oh-my-zsh"
export PATH="$HOME/.cargo/bin:$PATH"

if [[ $- == *i* ]]; then
    # Commands to run in interactive sessions can go here
fi

if [[ "$(uname)" == "Darwin" ]]; then
    BREW_BIN="/opt/homebrew/bin"
else
    BREW_BIN="/home/linuxbrew/.linuxbrew/bin"
fi

eval "$($BREW_BIN/brew shellenv)"

source $(dirname $BREW_BIN)/share/zsh-autocomplete/zsh-autocomplete.plugin.zsh
source $(dirname $BREW_BIN)/share/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh
source $(dirname $BREW_BIN)/share/zsh-autosuggestions/zsh-autosuggestions.zsh
source $(dirname $BREW_BIN)/share/powerlevel10k/powerlevel10k.zsh-theme

export PROJECT_PATHS="/home/alanbuscaglia/work"
export FZF_DEFAULT_COMMAND="fd --hidden --strip-cwd-prefix --exclude .git"
export FZF_DEFAULT_T_COMMAND="$FZF_DEFAULT_COMMAND"
export FZF_ALT_COMMAND="fd --type=d --hidden --strip-cwd-prefix --exlude .git"

WM_VAR="/$TMUX"
WM_CMD="tmux"

function start_if_needed() {
    if [[ $- == *i* ]] && [[ -z "${WM_VAR#/}" ]] && [[ -t 1 ]]; then
        exec $WM_CMD
    fi
}

alias fzfbat='fzf --preview="bat --theme=gruvbox-dark --color=always {}"'
alias fzfnvim='nvim $(fzf --preview="bat --theme=gruvbox-dark --color=always {}")'

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

[[ ! -f ~/.p10k.zsh ]] || source ~/.p10k.zsh

start_if_needed
      '';
    };
  };
}
