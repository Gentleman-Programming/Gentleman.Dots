if status is-interactive
    # Commands to run in interactive sessions can go here
end

if test (uname) = Darwin
    # macOS
    set BREW_BIN /opt/homebrew/bin/brew
else
    # Linux
    set BREW_BIN /home/linuxbrew/.linuxbrew/bin/brew
end

eval ($BREW_BIN shellenv)

if not set -q TMUX
    tmux
end

#if not set -q ZELLIJ 
#  zellij
#end

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

## alias
alias ls='gls --color=auto'
alias fzfbat='fzf --preview="bat --theme=gruvbox-dark --color=always {}"'
alias fzfnvim='nvim (fzf --preview="bat --theme=gruvbox-dark --color=always {}")'

## everforest
#set -l foreground d3c6aa
#set -l selection 2d4f67
#set -l comment 859289
#set -l red e67e80
#set -l orange ff9e64
#set -l yellow dbbc7f
#set -l green a7c080
#set -l purple d699b6
#set -l cyan 7fbbb3
#set -l pink d699b6

## rose pine moon colors
#set -l foreground e0def4 # text - a soft white for main text
#set -l selection 2a2a37 # highlight_high - dark blue for selection
#set -l comment 6e6a86 # muted - gray for comments
#set -l red eb6f92 # love - soft red
#set -l orange f6c177 # gold - soft orange
#set -l yellow f6c177 # gold - warm yellow
#set -l green 9ccfd8 # pine - pastel green
#set -l purple c4a7e7 # iris - soft purple
#set -l cyan 9ccfd8 # foam - teal
#set -l pink eb6f92 # love - soft pink

## Sakura colors
#set -l foreground c5a3a9  # na: text (light pink)
#set -l selection 3f3b3e   # gr: dark gray (highlight)
#set -l comment 4e4044     # nb: dark brown (comments)
#set -l red c58ea7         # ia: intense pink (errors)
#set -l orange 9e97d0      # ca: soft purple (warnings)
#set -l yellow 9e97d0      # ca: soft purple (warnings)
#set -l green 878fb9       # va: light blue (success)
#set -l purple 9e97d0      # ca: soft purple (highlight)
#set -l cyan 878fb9        # va: light blue (information)
#set -l pink c58ea7        # ia: intense pink (highlight)

# --- Base colors ---
set -l foreground C9C7CD # na: main text (light gray)
set -l selection 3B4252 # gr: dark gray (highlight)
set -l comment 4C566A # nb: medium gray (comments)

# --- Accent colors ---
set -l red EA83A5 # ia: intense pink (errors)
set -l orange F5A191 # ca: light peach (warnings)
set -l yellow E6B99D # ca: beige (warnings)
set -l green 90B99F # va: soft green (success)
set -l purple 92A2D5 # ca: lavender blue (highlight)
set -l cyan 85B5BA # va: blue-green (information)
set -l pink E29ECA # ia: soft pink (highlight)

# Syntax Highlighting Colors
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

# Completion Pager Colors
set -g fish_pager_color_progress $comment
set -g fish_pager_color_prefix $cyan
set -g fish_pager_color_completion $foreground
set -g fish_pager_color_description $comment

clear
