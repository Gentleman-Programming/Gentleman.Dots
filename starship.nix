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
palette = "kanagawa" # Le decimos a Starship que use nuestra paleta personalizada de abajo

[fill]
symbol = ' '

# Paleta Catppuccin Mocha (la dejamos por si alguna vez querés cambiar rápido)
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

# MI PALETA KANAGAWA "WAVE" PERSONALIZADA
# Acá mapeamos los nombres de color estándar de Starship a los colores "wave"
[palettes.kanagawa]
text      = "#DCD7BA"  # fujiWhite (texto principal, legible)
red       = "#E46876"  # waveRed (para errores, Rust, Java - bien visible)
green     = "#98BB6C"  # springGreen (para éxito, NodeJS - fresco)
yellow    = "#FF9E3B"  # roninYellow (para cmd_duration - llamativo pero no agresivo)
blue      = "#7E9CD8"  # crystalBlue (para directorios, username, C - claro y distintivo)
mauve     = "#957FB8"  # oniViolet (para git_branch, vimcmd_visual - elegante)
pink      = "#D27E99"  # sakuraPink (para git_status bg - un toque de color)
teal      = "#7AA89F"  # waveAqua2 (para Go, y como un color secundario fachero)
peach     = "#FFA066"  # surimiOrange (para PHP, Python, Zig, vimcmd_replace_one - cálido)
subtext0  = "#727169"  # fujiGray (para la hora y texto menos importante - sutil)

# Colores adicionales de la paleta Kanagawa original o "wave" que pueden ser útiles
# o que Starship podría buscar si se usan nombres más específicos.
# Los mantengo para consistencia y por si los usás en algún lado.
overlay0  = "#2D4F67"  # waveBlue2
rosewater = "#C0A36E"  # boatYellow2 (un amarillo suave)
flamingo  = "#D27E99"  # sakuraPink (ya que 'pink' es este, y git_branch usará 'mauve')
maroon    = "#C34043"  # autumnRed (un rojo más oscuro, como alternativa)
lavender  = "#b8b4d0"  # oniViolet2 (un lavanda más claro)
subtext1  = "#C8C093"  # oldWhite
overlay2  = "#54546D"  # sumiInk6
overlay1  = "#363646"  # sumiInk5
surface2  = "#223249"  # waveBlue1
surface1  = "#2A2A37"  # sumiInk4
surface0  = "#1F1F28"  # sumiInk3
base      = "#1F1F28"  # sumiInk3 (fondo base)
mantle    = "#16161D"  # sumiInk0 (fondo más oscuro)
crust     = "#16161D"  # sumiInk0

[character]
success_symbol = "[󱗞](fg:green)" # Usará nuestro springGreen
error_symbol   = "[󱗞](fg:red)"   # Usará nuestro waveRed
vimcmd_symbol = "[N](bold red)" # Usará waveRed
vimcmd_replace_one_symbol = "[R](bold peach)" # Usará nuestro surimiOrange
vimcmd_visual_symbol = "[V](bold mauve)" # Usará nuestro oniViolet

[username]
style_user    = 'bold blue' # Cambiado de teal a blue (crystalBlue)
style_root    = 'bold red'  # Usará waveRed
format        = '[󱗞 $user](fg:$style) '
disabled      = false
show_always   = true

[directory]
format                = "[$path](bold $style)[$read_only]($read_only_style) "
truncation_length     = 2
style                 = "fg:blue" # Cambiado de teal a blue (crystalBlue)
read_only_style       = "fg:blue" # Consistente con el path
before_repo_root_style= "fg:blue" # Consistente
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
style     = "bold fg:yellow" # Usará nuestro roninYellow
min_time  = 500

[git_branch]
format   = "-> [$symbol$branch]($style) "
style    = "bold fg:mauve" # Cambiado de flamingo a mauve (oniViolet)
symbol   = " "

[git_status]
format = '[$all_status$ahead_behind ]($style)'
style  = "fg:text bg:pink" # Usará fujiWhite sobre sakuraPink

[docker_context]
disabled = true
symbol   = " "

[python]
disabled       = false
format         = "[$symbol$pyenv_prefix($version)( $virtualenv)](fg:peach)($style)" # Usará surimiOrange
symbol         = " "
version_format = "$raw"

[java]
format         = '[[ $symbol ($version) ](fg:red)]($style)' # Usará waveRed
version_format = "$raw"
symbol         = " "
disabled       = false

[c]
format         = '[[ $symbol ($version) ](fg:blue)]($style)' # Cambiado a blue (crystalBlue) para consistencia
symbol         = " "
version_format = "$raw"
disabled       = false

[zig]
format         = '[[ $symbol ($version) ](fg:peach)]($style)' # Usará surimiOrange
version_format = "$raw"
disabled       = false

[bun]
version_format = "$raw"
format         = '[[ $symbol ($version) ](fg:text)]($style)' # Usará fujiWhite
disabled       = false

[nodejs]
symbol = ""
format = '[[ $symbol ($version) ](fg:green)]($style)' # Usará springGreen

[rust]
symbol = ""
format = '[[ $symbol ($version) ](fg:red)]($style)' # Usará waveRed

[golang]
symbol = ""
format = '[[ $symbol ($version) ](fg:teal)]($style)' # Usará waveAqua2

[php]
symbol = ""
format = '[[ $symbol ($version) ](fg:peach)]($style)' # Usará surimiOrange

[time]
disabled    = false
time_format = "%R"
format      = '[[   $time ](fg:subtext0)]($style)' # Usará fujiGray
      '';
    };
  };
}
