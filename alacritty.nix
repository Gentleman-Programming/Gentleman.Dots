{ config, pkgs, ... }:
{
  home.file = {
    ".config/alacritty/alacritty.toml" = {
      text = ''
[colors.primary]
background = "#000000" # bl: dark background (almost black)
foreground = "#C9C7CD" # na: main text (light gray)

[colors.cursor]
cursor = "#92A2D5" # ca: blue lavender (cursor)
text = "#C9C7CD"   # na: main text (light gray)

[colors.selection]
background = "#3B4252" # gr: dark gray (selection background)
text = "#C9C7CD"       # na: main text (light gray)

[colors.normal]
black = "#000000"   # bl: dark background (almost black)
red = "#EA83A5"     # ia: intense pink (errors)
green = "#90B99F"   # va: soft green (success)
yellow = "#E6B99D"  # ca: beige (warnings)
blue = "#85B5BA"    # va: light blue-green (information)
magenta = "#92A2D5" # ca: blue lavender (highlight)
cyan = "#85B5BA"    # va: light blue-green (links)
white = "#C9C7CD"   # na: main text (light gray)

[colors.bright]
black = "#4C566A"   # nb: medium gray (bright black)
red = "#EA83A5"     # ia: intense pink (bright red)
green = "#90B99F"   # va: soft green (bright green)
yellow = "#E6B99D"  # ca: beige (bright yellow)
blue = "#85B5BA"    # va: light blue-green (bright blue)
magenta = "#92A2D5" # ca: blue lavender (bright magenta)
cyan = "#85B5BA"    # va: light blue-green (bright cyan)
white = "#C9C7CD"   # na: main text (bright white)

[[colors.indexed_colors]]
index = 16
color = "#F5A191" # ca: light peach (orange)

[[colors.indexed_colors]]
index = 17
color = "#E29ECA" # ia: soft pink (pink)

[cursor]
style = "Block"
unfocused_hollow = true

[font]
size = 16

[font.normal]
family = "IosevkaTerm NF"

[window]
option_as_alt = "Both"
opacity = 0.96

[env]
TERM = "xterm-256color"
      '';
    };
  };
}
