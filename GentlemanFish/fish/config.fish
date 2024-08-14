if status is-interactive
    # Commands to run in interactive sessions can go here
end

set PROJECT_PATHS /your/work/path/

if test (uname) = "Darwin"
    # macOS
    set BREW_BIN /opt/homebrew/bin/brew
else
    # Linux
    set BREW_BIN /home/linuxbrew/.linuxbrew/bin/brew
end

eval ($BREW_BIN shellenv)

# Run TMUX
# if status is-interactive
#     and not set -q TMUX
#     exec tmux
# end

# Run Zellij
if set -q ZELLIJ
else
    zellij
end

starship init fish | source

set -x LS_COLORS "di=38;5;73:ow=48;5;22:ex=38;5;131:ln=38;5;179:*.tar=38;5;209:*.zip=38;5;209:*.jpg=38;5;176:*.png=38;5;176:*.mp3=38;5;176:*.wav=38;5;176:*.txt=38;5;223:*.sh=38;5;131"
alias ls='gls --color=auto'

## kanagawa
#set -l foreground DCD7BA
#set -l selection 2D4F67
#set -l comment 727169
#set -l red C34043
#set -l orange FF9E64
#set -l yellow C0A36E
#set -l green 76946A
#set -l purple 957FB8
#set -l cyan 7AA89F
#set -l pink D27E99

## everforest
set -l foreground d3c6aa
set -l selection 2d4f67
set -l comment 859289
set -l red e67e80
set -l orange ff9e64
set -l yellow dbbc7f
set -l green a7c080
set -l purple d699b6
set -l cyan 7fbbb3
set -l pink d699b6

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

# bun
set --export BUN_INSTALL "$HOME/.bun"
set --export PATH $BUN_INSTALL/bin $PATH
