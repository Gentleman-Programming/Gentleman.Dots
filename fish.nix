{ pkgs, lib, ... }:
{
  programs.fish = {
    enable = true;

    interactiveShellInit = ''
      if status is-interactive
          # Install Fisher if not installed
          if not functions -q fisher
              curl -sL https://git.io/fisher | source
              fisher install jorgebucaran/fisher
          end

          # Set Catppuccin Mocha as default theme
          fish_config theme choose "Catppuccin Mocha"
      end

      # Set BREW_BIN based on OS
      if test (uname) = Darwin
          set BREW_BIN /opt/homebrew/bin/brew
      else
          set BREW_BIN /home/linuxbrew/.linuxbrew/bin/brew
      end

      # Set PATH
      set -x PATH $HOME/.volta/bin $HOME/.bun/bin $HOME/.nix-profile/bin /nix/var/nix/profiles/default/bin $PATH /usr/local/bin $HOME/.config $HOME/.cargo/bin /usr/local/lib/*

      eval ($BREW_BIN shellenv)

      starship init fish | source
      zoxide init fish | source
      atuin init fish | source
      fzf --fish | source

      set -Ux CARAPACE_BRIDGES 'zsh,fish,bash,inshellisense'

      if not test -d ~/.config/fish/completions
          mkdir -p ~/.config/fish/completions
      end

      if not test -f ~/.config/fish/completions/.initialized
          carapace --list | awk '{print $1}' | xargs -I{} touch ~/.config/fish/completions/{}.fish
          touch ~/.config/fish/completions/.initialized
      end

      carapace _carapace | source

      # set -x LS_COLORS "di=38;5;67:ow=48;5;60:ex=38;5;132:ln=38;5;144:*.tar=38;5;180:*.zip=38;5;180:*.jpg=38;5;175:*.png=38;5;175:*.mp3=38;5;175:*.wav=38;5;175:*.txt=38;5;223:*.sh=38;5;132"
      set -g fish_greeting ""

      # Enable vi mode
      fish_vi_key_bindings

      # Set visual indicator for vi mode
      function fish_mode_prompt
          switch $fish_bind_mode
              case default
                  set_color --bold red
                  echo '[N] '
              case insert
                  set_color --bold green
                  echo '[I] '
              case replace_one
                  set_color --bold green
                  echo '[R] '
              case visual
                  set_color --bold magenta
                  echo '[V] '
              case '*'
                  set_color --bold red
                  echo '[?] '
          end
          set_color normal
      end

      if test (uname) = Darwin
          alias ls='ls --color=auto'
      else
          alias ls='gls --color=auto'
      end

      alias fzfbat='fzf --preview="bat --theme=gruvbox-dark --color=always {}"'
      alias fzfnvim='nvim (fzf --preview="bat --theme=gruvbox-dark --color=always {}")'

      clear
    '';

    plugins = [
      {
        name = "fisher";
        src = pkgs.fetchFromGitHub {
          owner = "jorgebucaran";
          repo = "fisher";
          rev = "4.4.4";
          sha256 = "sha256-e8gIaVbuUzTwKtuMPNXBT5STeddYqQegduWBtURLT3M=";
        };
      }
      {
        name = "catppuccin";
        src = pkgs.fetchFromGitHub {
          owner = "catppuccin";
          repo = "fish";
          rev = "0ce27b518e8ead555dec34dd8be3df5bd75cff8e";
          sha256 = "sha256-Dc/zdxfzAUM5NX8PxzfljRbYvO9f9syuLO8yBr+R3qg=";
        };
      }
    ];
  };
}
