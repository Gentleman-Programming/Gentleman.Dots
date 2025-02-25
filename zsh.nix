{ config, pkgs, ... }:

{
  programs.zsh = {
    enable = true;
    oh-my-zsh = {
      enable = true;
      plugins = [ "command-not-found" ];
    };
    initExtra = ''
      export ZSH="$HOME/.oh-my-zsh"
      export PATH="$HOME/.cargo/bin:$PATH"

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

      source ${pkgs.zsh-vi-mode}/share/zsh-vi-mode/zsh-vi-mode.plugin.zsh
      source ${pkgs.zsh-autocomplete}/share/zsh-autocomplete/zsh-autocomplete.plugin.zsh
      source ${pkgs.zsh-autosuggestions}/share/zsh-autosuggestions/zsh-autosuggestions.zsh
      source ${pkgs.zsh-syntax-highlighting}/share/zsh-syntax-highlighting/zsh-syntax-highlighting.zsh

      export PATH="$HOME/.volta/bin:$HOME/.bun/bin:$HOME/.nix-profile/bin:/nix/var/nix/profiles/default/bin:$PATH:/usr/local/bin:$HOME/.config:$HOME/.cargo/bin:/usr/local/lib/*"
      export FZF_DEFAULT_COMMAND="fd --hidden --strip-cwd-prefix --exclude .git"
      export FZF_DEFAULT_T_COMMAND="$FZF_DEFAULT_COMMAND"
      export FZF_ALT_COMMAND="fd --type=d --hidden --strip-cwd-prefix --exlude .git"

      # alias
      alias fzfbat='fzf --preview="bat --theme=gruvbox-dark --color=always {}"'
      alias fzfnvim='nvim $(fzf --preview="bat --theme=gruvbox-dark --color=always {}")'

      export CARAPACE_BRIDGES='zsh,fish,bash,inshellisense'
      zstyle ':completion:*' format $'\e[2;37mCompleting %d\e[m'
      source <(carapace _carapace)

      eval "$(fzf --zsh)"
      eval "$(zoxide init zsh)"
      eval "$(atuin init zsh)"
      eval "$(starship init zsh)"

      clear
    '';
  };
}
