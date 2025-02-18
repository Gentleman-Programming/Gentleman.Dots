{ config, pkgs, ... }:
{
  home.file = {
    ".config/starship.toml" = {
      text = ''
format = """\
($directory)\
$os\
$git_branch\
$fill\
$nodejs\
$rust\
$golang\
$php\
$bun\
$java\
$c\
$conda\
$zig\
$cmd_duration\
$time\
\n$character\
"""

add_newline = true
command_timeout = 3600000
palette = "oldWorld"

[fill]
symbol = ' '

[palettes.oldWorld]
arrow = "#EA83A5"
os = "#85B5BA"
directory = "#92A2D5"
git = "#EA83A5"
duration = "#92A2D5"
text_color = "#C9C7CD"
text_light = "#000000"

[character]
success_symbol = "[󱗞](fg:#85B5BA)"
error_symbol   = "[󱗞](fg:red)"

[username]
style_user    = 'bold os'
style_root    = 'bold os_admin'
format        = '[󱗞 $user](fg:$style) '
disabled      = false
show_always   = true

[directory]
format                = "[$path](bold $style)[$read_only]($read_only_style) "
truncation_length     = 2
style                 = "fg:directory"
read_only_style       = "fg:directory"
before_repo_root_style= "fg:directory"
truncation_symbol     = "…/"
truncate_to_repo      = true
read_only             = "  "

[directory.substitutions]
"Documents" = "󰈙 "
"Downloads" = " "
"Music"     = " "
"Pictures"  = " "

[cmd_duration]
format    = " took [ $duration]($style) "
style     = "bold fg:duration"
min_time  = 500

[git_branch]
format   = "-> [$symbol$branch]($style) "
style    = "bold fg:git"
symbol   = " "

[git_status]
format = '[$all_status$ahead_behind ]($style)'
style  = "fg:text_color bg:git"

[docker_context]
disabled = true
symbol   = " "

# Python: usamos $symbol, $pyenv_prefix, $version, $virtualenv
[python]
disabled       = false
format         = "[$symbol$pyenv_prefix($version)( $virtualenv)](fg:#FF9E3B)($style)"
symbol         = " "
version_format = "$raw"

[java]
format         = '[[ $symbol ($version) ](fg:#FF5D62)]($style)'
version_format = "$raw"
symbol         = " "
disabled       = false

[c]
format         = '[[ $symbol ($version) ](fg:#7FB4CA)]($style)'
symbol         = " "
version_format = "$raw"
disabled       = false

[zig]
format         = '[[ $symbol ($version) ](fg:#FFA066)]($style)'
version_format = "$raw"
disabled       = false

[bun]
version_format = "$raw"
format         = '[[ $symbol ($version) ](fg:#DCD7BA)]($style)'
disabled       = false

[nodejs]
symbol = ""
format = '[[ $symbol ($version) ](fg:#87a987)]($style)'

[rust]
symbol = ""
format = '[[ $symbol ($version) ](fg:#FF5D62)]($style)'

[golang]
symbol = ""
format = '[[ $symbol ($version) ](fg:#7FB4CA)]($style)'

[php]
symbol = ""
format = '[[ $symbol ($version) ](fg:#FF9E3B)]($style)'

[time]
disabled    = false
time_format = "%R"
format      = '[[   $time ](fg:#8BA4B0)]($style)'
      '';
    };
  };
}
