# Graphite-rose Gentleman-Cute Fish colors.

set -l surface "#241822"
set -l selection "#342230"
set -l text "#F6EFF3"
set -l muted "#A78E9B"
set -l dimmed "#76616B"
set -l accent "#F095C8"
set -l deep_rose "#C96AA2"
set -l active_rose "#FFB1DD"
set -l champagne "#E0C27A"
set -l pearl "#D2CBD0"
set -l soft_rose "#D7A0B8"
set -l powder_blue "#A9C7EE"
set -l warning "#F2B86D"
set -l error "#FF718F"

# Syntax highlighting colors.
set -g fish_color_normal $text
set -g fish_color_command --bold $accent
set -g fish_color_keyword --bold $active_rose
set -g fish_color_quote $champagne
set -g fish_color_redirection $powder_blue
set -g fish_color_end $warning
set -g fish_color_error --bold $error
set -g fish_color_param $soft_rose
set -g fish_color_comment --dim $dimmed
set -g fish_color_selection --background=$selection
set -g fish_color_search_match --background=$selection
set -g fish_color_operator $pearl
set -g fish_color_escape $deep_rose
set -g fish_color_option $soft_rose
set -g fish_color_autosuggestion --dim $muted

# Completion pager colors.
set -g fish_pager_color_progress --dim $dimmed
set -g fish_pager_color_prefix --bold $accent
set -g fish_pager_color_completion $text
set -g fish_pager_color_description --dim $muted
set -g fish_pager_color_selected_background --background=$selection
set -g fish_pager_color_secondary_background --background=$surface
