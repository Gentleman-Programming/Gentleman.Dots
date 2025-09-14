# Enable Powerlevel10k instant prompt. Should stay close to the top of ~/.zshrc.
# Initialization code that may require console input (password prompts, [y/n]
# confirmations, etc.) must go above this block; everything else may go below.
if [[ -r "${XDG_CACHE_HOME:-$HOME/.cache}/p10k-instant-prompt-${(%):-%n}.zsh" ]]; then
  source "${XDG_CACHE_HOME:-$HOME/.cache}/p10k-instant-prompt-${(%):-%n}.zsh"
fi


typeset -g POWERLEVEL9K_INSTANT_PROMPT=quiet

#--------------------------------------------------------------------------------------------------------------------

export ZSH="$HOME/.oh-my-zsh"
export PATH="$HOME/.volta/bin:$HOME/.cargo/bin:$PATH"

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

export NVM_DIR="$HOME/.nvm"
[ -s "/home/linuxbrew/.linuxbrew/opt/nvm/nvm.sh" ] && \. "/home/linuxbrew/.linuxbrew/opt/nvm/nvm.sh"                                       # This loads nvm
[ -s "/home/linuxbrew/.linuxbrew/opt/nvm/etc/bash_completion.d/nvm" ] && \. "/home/linuxbrew/.linuxbrew/opt/nvm/etc/bash_completion.d/nvm" # This loads nvm bash_completion

export HOMEBREW_NO_AUTO_UPDATE=1
export HOMEBREW_NO_ENV_HINTS=1
export HOMEBREW_NO_ANALYTICS=1

export PROJECT_PATHS="/home/deuri-vasquez/projects"
export FZF_DEFAULT_COMMAND="fd --hidden --strip-cwd-prefix --exclude .git"
export FZF_DEFAULT_T_COMMAND="$FZF_DEFAULT_COMMAND"
export FZF_ALT_COMMAND="fd --type=d --hidden --strip-cwd-prefix --exlude .git"

export PATH="/home/linuxbrew/.linuxbrew/opt/coreutils/libexec/gnubin:$PATH"
#---------------------------------------------------------------------------------------------------------------------
WM_VAR="/$ZELLIJ"
# change with ZELLIJ
WM_CMD="zellij"
# change with zellij

function start_if_needed() {
    if [[ $- == *i* ]] && [[ -z "${WM_VAR#/}" ]] && [[ -t 1 ]]; then
        exec $WM_CMD
    fi
}

# ~/.zshrc

# ---
# Lógica para auto-renombrar los tabs de Zellij en Zsh.
# Se ejecuta solo si estamos dentro de una sesión de Zellij.
# ---

# Array de comandos que queremos ignorar.
# Estos procesos son de corta duración y no deberían cambiar el nombre del tab principal.
typeset -gA ignored_commands
ignored_commands=(
    git 1
    ls 1
    grep 1
    clear 1
    docker 1
    npm 1
    yarn 1
    pnpm 1
    npx 1
    exit 1
    z 1
    cd 1
)

# Esta función se ejecuta justo antes de cada comando que escribas.
# Chequea si el comando está en la lista de ignorados antes de renombrar el tab.
function zellij_rename_preexec() {
    # Evita que se ejecute si no estamos en un Zellij
    if [[ -z "$ZELLIJ" ]]; then
        return
    fi

    # Extrae el nombre del comando a ejecutar
    local command_name="${1%% *}"

    # Si el comando es "z" (un alias común de zellij) o está en la lista de ignorados,
    # no hacemos nada para evitar cambios innecesarios.
    if [[ "$command_name" = "z" || -n "${ignored_commands[$command_name]}" ]]; then
        return
    fi
    
    # Renombra el tab con el nombre del comando.
    # El "&" lo ejecuta en segundo plano para no bloquear el shell.
    nohup zellij action rename-tab "$command_name" >/dev/null 2>&1 
}

# Esta función se ejecuta justo antes de mostrar el prompt.
# Se usa para limpiar el nombre del tab, volviéndolo al nombre del directorio actual.
function zellij_rename_precmd() {
    # Evita que se ejecute si no estamos en un Zellij
    if [[ -z "$ZELLIJ" ]]; then
        return
    fi

    # Renombra el tab con el nombre del directorio actual.
  nohup zellij action rename-tab "$(basename "$PWD")" >/dev/null 2>&1 
}

# Registra las funciones a los hooks de Zsh.
# `preexec` se ejecuta antes de cada comando.
# `precmd` se ejecuta antes de mostrar el prompt.
autoload -Uz add-zsh-hook
add-zsh-hook preexec zellij_rename_preexec
add-zsh-hook precmd zellij_rename_precmd

#---------------------------------------------------------------------------------------------------------------------
# alias
alias fzfbat='fzf --preview="bat --theme=gruvbox-dark --color=always {}"'
alias fzfnvim='nvim $(fzf --preview="bat --theme=gruvbox-dark --color=always {}")'
alias lzd="lazydocker"
alias zz="zellij"
alias vi="nvim"
alias vim="nvim"
alias zz="zellij"

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

start_if_needed
