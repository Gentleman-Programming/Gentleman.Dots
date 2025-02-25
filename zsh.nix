{ config, pkgs, ... }:

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

      # --------------------------
      # 7) LS configuration 
      # --------------------------
      export LS_COLORS="di=38;5;67:ow=48;5;60:ex=38;5;132:ln=38;5;144:*.tar=38;5;180:*.zip=38;5;180:*.jpg=38;5;175:*.png=38;5;175:*.mp3=38;5;175:*.wav=38;5;175:*.txt=38;5;223:*.sh=38;5;132"
      if [[ "$(uname)" == "Darwin" ]]; then
        alias ls='ls --color=auto'
      else
        alias ls='gls --color=auto'
      fi
    '';
  };
}
