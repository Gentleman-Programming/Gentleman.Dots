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
palette = "kanagawa"

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

[palettes.kanagawa]
text = "#DCD7BA"
red = "#E46876"
green = "#98BB6C"
yellow = "#FF9E3B"
blue = "#7E9CD8"
mauve = "#957FB8"
pink = "#D27E99"
teal = "#7AA89F"
peach = "#FFA066"
subtext0 = "#727169"
overlay0 = "#2D4F67"
rosewater = "#C0A36E"
flamingo = "#D27E99"
maroon = "#C34043"
lavender = "#b8b4d0"
subtext1 = "#C8C093"
overlay2 = "#54546D"
overlay1 = "#363646"
surface2 = "#223249"
surface1 = "#2A2A37"
surface0 = "#1F1F28"
base = "#1F1F28"
mantle = "#16161D"
crust = "#16161D"

[character]
success_symbol = "[󱗞](fg:green)"
error_symbol = "[󱗞](fg:red)"
vimcmd_symbol = "[N](bold red)"
vimcmd_replace_one_symbol = "[R](bold peach)"
vimcmd_visual_symbol = "[V](bold mauve)"

[username]
style_user = 'bold blue'
style_root = 'bold red'
format = '[󱗞 $user](fg:$style) '
disabled = false
show_always = true

[directory]
format = "[$path](bold $style)[$read_only]($read_only_style) "
truncation_length = 2
style = "fg:blue"
read_only_style = "fg:blue"
before_repo_root_style = "fg:blue"
truncation_symbol = "…/"
truncate_to_repo = true
read_only = "  "

[directory.substitutions]
"Documents" = "󰈙 "
"Downloads" = " "
"Music" = " "
"Pictures" = " "

[cmd_duration]
format = " took [ $duration]($style) "
style = "bold fg:yellow"
min_time = 500

[git_branch]
format = "-> [$symbol$branch]($style) "
style = "bold fg:mauve"
symbol = " "

[git_status]
format = '[$all_status$ahead_behind ]($style)'
style = "fg:text bg:pink"

[docker_context]
disabled = true
symbol = " "

[python]
disabled = false
format = "[$symbol$pyenv_prefix($version)( $virtualenv)](fg:peach)"
symbol = " "
version_format = "$raw"

[java]
format = '[[ $symbol ($version) ](fg:red)]($style)'
version_format = "$raw"
symbol = " "
disabled = false

[c]
format = '[[ $symbol ($version) ](fg:blue)]($style)'
symbol = " "
version_format = "$raw"
disabled = false

[zig]
format = '[[ $symbol ($version) ](fg:peach)]($style)'
version_format = "$raw"
disabled = false

[bun]
version_format = "$raw"
format = '[[ $symbol ($version) ](fg:text)]($style)'
disabled = false

[nodejs]
symbol = ""
format = '[[ $symbol ($version) ](fg:green)]($style)'

[rust]
symbol = ""
format = '[[ $symbol ($version) ](fg:red)]($style)'

[golang]
symbol = ""
format = '[[ $symbol ($version) ](fg:teal)]($style)'

[php]
symbol = ""
format = '[[ $symbol ($version) ](fg:peach)]($style)'

[time]
disabled = false
time_format = "%R"
format = '[[   $time ](fg:subtext0)]($style)'
      '';
    };
  };
}
