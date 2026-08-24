# Neutral/cyan Gentleman Fish colors.

set -l foreground F3F6F9 normal
set -l selection 263356 normal
set -l comment 8394A3 brblack
set -l red CB7C94 red
set -l orange DEBA87 orange
set -l yellow FFE066 yellow
set -l green B7CC85 green
set -l purple A3B5D6 purple
set -l cyan 7AA89F cyan
set -l pink FF8DD7 magenta

# Syntax highlighting colors.
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

# Completion pager colors.
set -g fish_pager_color_progress $comment
set -g fish_pager_color_prefix $cyan
set -g fish_pager_color_completion $foreground
set -g fish_pager_color_description $comment

# fzf colors.
set -gx FZF_DEFAULT_OPTS "--color=fg:#F3F6F9,bg:-1,gutter:-1,hl:#7AA89F,fg+:#F3F6F9,bg+:#263356,hl+:#7AA89F,info:#8394A3,prompt:#7AA89F,pointer:#FF8DD7,marker:#FFE066,spinner:#A3B5D6,header:#8394A3,border:#263356,separator:#263356,label:#8394A3"
