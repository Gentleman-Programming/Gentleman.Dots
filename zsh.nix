{
  programs.zsh = {
    enable = true;
    # Enable completions
    enableCompletion = false;

    # zplug configuration
    zplug = {
      enable = true;
      plugins = [
        { name = "zsh-users/zsh-autosuggestions"; }
        { name = "zsh-users/zsh-syntax-highlighting"; }
        { name = "marlonrichert/zsh-autocomplete"; }
        { name = "jeffreytse/zsh-vi-mode"; }
      ];
    };

    # Extra initialization
    initExtra = ''
      # --------------------------
      # 1) COMPINIT + CACHE
      # --------------------------
      autoload -Uz compinit
      # Use a directory in .cache or as you prefer
      compinit -d "$${XDG_CACHE_HOME:-$${HOME}/.cache}/zsh/zcompdump-$${ZSH_VERSION}"

      # --------------------------
      # 2) FZF
      # --------------------------
      export FZF_DEFAULT_COMMAND="fd --hidden --strip-cwd-prefix --exclude .git"
      export FZF_DEFAULT_T_COMMAND="$FZF_DEFAULT_COMMAND"
      export FZF_ALT_COMMAND="fd --type=d --hidden --strip-cwd-prefix --exclude .git"

      alias fzfbat='fzf --preview="bat --theme=gruvbox-dark --color=always {}"'
      alias fzfnvim='nvim $(fzf --preview="bat --theme=gruvbox-dark --color=always {})"'

      # If you really need this eval, leave it:
      # eval "$(fzf --zsh)"

      # --------------------------
      # 3) Carapace
      # --------------------------
      export CARAPACE_BRIDGES='zsh,fish,bash,inshellisense'
      zstyle ':completion:*' format $'\e[2;37mCompleting %d\e[m'
      source <(carapace _carapace)

      # --------------------------
      # 4) Tools initialization
      # --------------------------
      eval "$(zoxide init zsh)"
      eval "$(atuin init zsh)"
      eval "$(starship init zsh)"

      ya_zed() {
        tmp=$(mktemp -t "yazi-chooser.XXXXXXXXXX")
        yazi --chooser-file "$tmp" "$@"

        if [[ -s "$tmp" ]]; then
          opened_file=$(head -n 1 -- "$tmp")
          if [[ -n "$opened_file" ]]; then
            if [[ -d "$opened_file" ]]; then
              # Es una carpeta, la agregamos al workspace
              zed --add "$opened_file"
            else
              # Es un archivo, lo abrimos normalmente
              zed --add "$opened_file"
            fi
          fi
        fi

        rm -f -- "$tmp"
      }

      # --------------------------
      # 5) Final cleanup
      # --------------------------
      # Clear gives you that "fresh" feeling,
      # but if you prefer speed, you can comment it out.
      clear

      # --------------------------
      # 6) Login shell specific configuration
      # --------------------------
      if [[ -o login ]]; then
        # PATHS and Variables
        export PATH="$HOME/.cargo/bin:$HOME/.volta/bin:$HOME/.bun/bin:$HOME/.nix-profile/bin:/nix/var/nix/profiles/default/bin:$PATH:/usr/local/bin:$HOME/.config:$HOME/.cargo/bin:/usr/local/lib/*"

        # macOS vs Linux distinction
        if [[ "$(uname)" == "Darwin" ]]; then
          export BREW_BIN="/opt/homebrew/bin"
        else
          export BREW_BIN="/home/linuxbrew/.linuxbrew/bin"
        fi

        # Load brew
        if [ -x "$BREW_BIN/brew" ]; then
          eval "$($BREW_BIN/brew shellenv)"
        fi
      fi
    '';
  };
}
