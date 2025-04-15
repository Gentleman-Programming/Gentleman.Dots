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
palette = "catppuccin_mocha"

[fill]
symbol = ' '

[palettes.catppuccin_mocha]
rosewater = "#f5e0dc"
flamingo = "#f2cdcd"
pink = "#f5c2e7"
mauve = "#cba6f7"
red = "#f38ba8"
maroon = "#eba0ac"
peach = "#fab387"
yellow = "#f9e2af"
green = "#a6e3a1"
teal = "#94e2d5"
sky = "#89dceb"
sapphire = "#74c7ec"
blue = "#89b4fa"
lavender = "#b4befe"
text = "#cdd6f4"
subtext1 = "#bac2de"
subtext0 = "#a6adc8"
overlay2 = "#9399b2"
overlay1 = "#7f849c"
overlay0 = "#6c7086"
surface2 = "#585b70"
surface1 = "#45475a"
surface0 = "#313244"
base = "#1e1e2e"
mantle = "#181825"
crust = "#11111b"

[character]
success_symbol = "[󱗞](fg:green)"
error_symbol   = "[󱗞](fg:red)"
vimcmd_symbol = "[N](bold red)"
vimcmd_replace_one_symbol = "[R](bold peach)"
vimcmd_visual_symbol = "[V](bold purple)"

[username]
style_user    = 'bold teal'
style_root    = 'bold red'
format        = '[󱗞 $user](fg:$style) '
disabled      = false
show_always   = true

[directory]
format                = "[$path](bold $style)[$read_only]($read_only_style) "
truncation_length     = 2
style                 = "fg:lavender"
read_only_style       = "fg:lavender"
before_repo_root_style= "fg:lavender"
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
style    = "bold fg:mauve"
symbol   = " "

[git_status]
format = '[$all_status$ahead_behind ]($style)'
style  = "fg:text bg:pink"

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
format         = '[[ $symbol ($version) ](fg:peach)]($style)'
version_format = "$raw"
disabled       = false

[bun]
version_format = "$raw"
format         = '[[ $symbol ($version) ](fg:text)]($style)'
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
