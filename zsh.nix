{
  # Zsh configuration via home.file (not programs.zsh to avoid recursive .zshenv bug)
  
  home.file.".zshenv" = {
    text = ''
      # Source Home Manager session variables if available
      if [ -e "$HOME/.nix-profile/etc/profile.d/hm-session-vars.sh" ]; then
        . "$HOME/.nix-profile/etc/profile.d/hm-session-vars.sh"
      fi
    '';
  };

  home.file.".zshrc" = {
    text = ''
      # --------------------------
      # 1) ZPLUG
      # --------------------------
      if [[ ! -d ~/.zplug ]]; then
        git clone https://github.com/zplug/zplug ~/.zplug
      fi
      source ~/.zplug/init.zsh

      zplug "zsh-users/zsh-autosuggestions"
      zplug "zsh-users/zsh-syntax-highlighting"
      zplug "marlonrichert/zsh-autocomplete"
      zplug "jeffreytse/zsh-vi-mode"

      if ! zplug check; then
        zplug install
      fi
      zplug load

      # --------------------------
      # 2) COMPINIT + CACHE
      # --------------------------
      autoload -Uz compinit
      compinit -d "''${XDG_CACHE_HOME:-''${HOME}/.cache}/zsh/zcompdump-''${ZSH_VERSION}"

      # --------------------------
      # 3) EDITOR
      # --------------------------
      export EDITOR="nvim"
      export VISUAL="nvim"

      # --------------------------
      # 4) FZF
      # --------------------------
      export FZF_DEFAULT_COMMAND="fd --hidden --strip-cwd-prefix --exclude .git"
      export FZF_DEFAULT_T_COMMAND="$FZF_DEFAULT_COMMAND"
      export FZF_ALT_COMMAND="fd --type=d --hidden --strip-cwd-prefix --exclude .git"

      alias fzfbat='fzf --preview="bat --theme=gruvbox-dark --color=always {}"'
      alias fzfnvim='nvim $(fzf --preview="bat --theme=gruvbox-dark --color=always {}")'

      # --------------------------
      # 5) Carapace
      # --------------------------
      export CARAPACE_BRIDGES='zsh,fish,bash,inshellisense'
      zstyle ':completion:*' format $'\e[2;37mCompleting %d\e[m'
      source <(carapace _carapace)

      # --------------------------
      # 6) Tools initialization
      # --------------------------
      eval "$(zoxide init zsh)"
      eval "$(atuin init zsh)"
      eval "$(starship init zsh)"

      # --------------------------
      # 7) Yazi + Zed function
      # --------------------------
      ya_zed() {
        tmp=$(mktemp -t "yazi-chooser.XXXXXXXXXX")
        yazi --chooser-file "$tmp" "$@"

        if [[ -s "$tmp" ]]; then
          opened_file=$(head -n 1 -- "$tmp")
          if [[ -n "$opened_file" ]]; then
            zed --add "$opened_file"
          fi
        fi

        rm -f -- "$tmp"
      }

      # --------------------------
      # 7.1) Oil aliases
      # --------------------------
      alias o='oil'
      alias oo='oil .'
      alias of='oil-float'
      alias oz='oil-zed'

      # --------------------------
      # 8) PATH and Brew
      # --------------------------
      export PATH="$HOME/.local/bin:$HOME/.local/state/nix/profiles/home-manager/home-path/bin:$HOME/.opencode/bin:$HOME/.cargo/bin:$HOME/.volta/bin:$HOME/.bun/bin:$HOME/.nix-profile/bin:/nix/var/nix/profiles/default/bin:/usr/local/bin:$PATH"

      if [[ "$(uname)" == "Darwin" ]]; then
        export BREW_BIN="/opt/homebrew/bin"
      else
        export BREW_BIN="/home/linuxbrew/.linuxbrew/bin"
      fi

      if [ -x "$BREW_BIN/brew" ]; then
        eval "$($BREW_BIN/brew shellenv)"
      fi

      # --------------------------
      # 9) TMUX auto-start
      # --------------------------
      WM_VAR="/$TMUX"
      WM_CMD="tmux"

      function start_if_needed() {
        if [[ $- == *i* ]] && [[ -z "''${WM_VAR#/}" ]] && [[ -t 1 ]] && [[ -z "$ZED_TERMINAL" ]]; then
          exec $WM_CMD
        fi
      }
      start_if_needed

      # --------------------------
      # 10) Clear screen
      # --------------------------
      clear
    '';
  };
}
