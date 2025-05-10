 { config, pkgs, ... }:
 {
   home.file = {
     ".config/wezter/wezterm.lua" = {
       text = ''
-- Import the wezterm API
local wezterm = require("wezterm")

-- Initialize an empty configuration table
local config = {}

-- OldWorld Theme
-- config.colors = {
--foreground = "--C9C7CD",
--background = "--000000",
--
--cursor_bg = "--92A2D5",
--cursor_fg = "--C9C7CD",
--cursor_border = "--92A2D5",
--
--selection_fg = "--C9C7CD",
--selection_bg = "--3B4252",
--
--scrollbar_thumb = "--4C566A",
--split = "--4C566A",
--
--ansi = {
--	"--000000",
--	"--EA83A5",
--	"--90B99F",
--	"--E6B99D",
--	"--85B5BA",
--	"--92A2D5",
--	"--85B5BA",
--	"--C9C7CD",
--},
--
--brights = {
--	"--4C566A",
--	"--EA83A5",
--	"--90B99F",
--	"--E6B99D",
--	"--85B5BA",
--	"--92A2D5",
--	"--85B5BA",
--	"--C9C7CD",
--},
--
--indexed = {
--	[16] = "--F5A191",
--	[17] = "--E29ECA",
 	},
 }
config.color_scheme = "Catppuccin Mocha"
config.window_padding = {
	top = 0,
	right = 0,
	left = 0,
}

config.font = wezterm.font("IosevkaTerm NF")
config.hide_tab_bar_if_only_one_tab = true
config.max_fps = 240
config.enable_kitty_graphics = true

config.window_background_opacity = 0.85
config.macos_window_background_blur = 20
config.win32_system_backdrop = "Acrylic"

config.font_size = 16.0
config.enable_scroll_bar = false

return config
       '';
     };
   };
 }
