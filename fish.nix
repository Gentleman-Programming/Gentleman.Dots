{ config, pkgs, ... }:
{
  programs.fish.enable = true;

  home.file = {
    ".config/fish/config.fish" = {
      text = ''
if status is-interactive
end

if test (uname) = Darwin
    set BREW_BIN /opt/homebrew/bin/brew
else
    set BREW_BIN /home/linuxbrew/.linuxbrew/bin/brew
end

eval ($BREW_BIN shellenv)

if not set -q TMUX
    tmux
end

starship init fish | source
zoxide init fish | source
atuin init fish | source

set -x PATH $HOME/.cargo/bin $PATH
set -Ux CARAPACE_BRIDGES 'zsh,fish,bash,inshellisense'
mkdir -p ~/.config/fish/completions
carapace --list | awk '{print $1}' | xargs -I{} touch ~/.config/fish/completions/{}.fish
carapace _carapace | source

set -x LS_COLORS "di=38;5;67:ow=48;5;60:ex=38;5;132:ln=38;5;144:*.tar=38;5;180:*.zip=38;5;180:*.jpg=38;5;175:*.png=38;5;175:*.mp3=38;5;175:*.wav=38;5;175:*.txt=38;5;223:*.sh=38;5;132"
set -g fish_greeting ""

alias ls='gls --color=auto'
alias fzfbat='fzf --preview="bat --theme=gruvbox-dark --color=always {}"'
alias fzfnvim='nvim (fzf --preview="bat --theme=gruvbox-dark --color=always {}")'

set -l foreground C9C7CD
set -l selection 3B4252
set -l comment 4C566A

set -l red EA83A5
set -l orange F5A191
set -l yellow E6B99D
set -l green 90B99F
set -l purple 92A2D5
set -l cyan 85B5BA
set -l pink E29ECA

set -g fish_color_normal $foreground
set -g fish_color_command $cyan
set -g fish_color_keyword $pink
set -g fish_color_quote $yellow
set -g fish_color_redirection $foreground
set -g fish_color_end $orange
set -g fish_color_error $red
set -g fish_color_param $purple
set -g fish_color_comment $comment
set -g fish_color_selection --background=$selection
set -g fish_color_search_match --background=$selection
set -g fish_color_operator $green
set -g fish_color_escape $pink
set -g fish_color_autosuggestion $comment

set -g fish_pager_color_progress $comment
set -g fish_pager_color_prefix $cyan
set -g fish_pager_color_completion $foreground
set -g fish_pager_color_description $comment

clear
      '';
    };
  };
}
