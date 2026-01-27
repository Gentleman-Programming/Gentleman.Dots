-- ╔══════════════════════════════════════════════════════════════════════════════╗
-- ║                          GENTLEMAN DOTS - WEZTERM                            ║
-- ║                           Optimized for Neovim                               ║
-- ╚══════════════════════════════════════════════════════════════════════════════╝

local wezterm = require("wezterm")
local config = {}

-- ┌──────────────────────────────────────────────────────────────────────────────┐
-- │                                   FONT                                       │
-- └──────────────────────────────────────────────────────────────────────────────┘

config.font = wezterm.font("IosevkaTerm NF")
config.font_size = 14.0

-- ┌──────────────────────────────────────────────────────────────────────────────┐
-- │                                  WINDOW                                      │
-- └──────────────────────────────────────────────────────────────────────────────┘

config.window_background_opacity = 0.95
config.macos_window_background_blur = 20
config.win32_system_backdrop = "Acrylic"

config.window_padding = {
	top = 0,
	right = 0,
	left = 0,
	bottom = 0,
}

config.enable_scroll_bar = false
config.hide_tab_bar_if_only_one_tab = true

-- ┌──────────────────────────────────────────────────────────────────────────────┐
-- │                                  CURSOR                                      │
-- └──────────────────────────────────────────────────────────────────────────────┘

config.default_cursor_style = "SteadyBlock"
config.cursor_blink_rate = 500
config.cursor_blink_ease_in = "Constant"
config.cursor_blink_ease_out = "Constant"

-- ┌──────────────────────────────────────────────────────────────────────────────┐
-- │                            NEOVIM OPTIMIZATIONS                              │
-- └──────────────────────────────────────────────────────────────────────────────┘

-- Terminal & Colors
-- NOTE: WSL users may need "xterm-256color" if fish/zsh fails with "missing terminal"
-- See: https://github.com/Gentleman-Programming/Gentleman.Dots/issues/117
config.term = "wezterm"
config.enable_csi_u_key_encoding = true

-- Undercurl support (LSP diagnostics, spelling)
config.underline_thickness = 2
config.underline_position = -2

-- Scrollback
config.scrollback_lines = 10000

-- Performance
config.max_fps = 240

-- Image support
config.enable_kitty_graphics = true

-- Input handling
config.use_dead_keys = false
config.send_composed_key_when_left_alt_is_pressed = false
config.send_composed_key_when_right_alt_is_pressed = false

-- ┌──────────────────────────────────────────────────────────────────────────────┐
-- │                           GENTLEMAN THEME                                    │
-- └──────────────────────────────────────────────────────────────────────────────┘

config.colors = {
	-- Base Colors
	foreground = "#f3f6f9",
	background = "#06080f",

	-- Cursor
	cursor_bg = "#e0c15a",
	cursor_fg = "#06080f",
	cursor_border = "#e0c15a",

	-- Selection
	selection_fg = "#f3f6f9",
	selection_bg = "#263356",

	-- Normal Colors
	ansi = {
		"#06080f", -- black
		"#cb7c94", -- red
		"#b7cc85", -- green
		"#ffe066", -- yellow
		"#7fb4ca", -- blue
		"#ff8dd7", -- magenta
		"#7aa89f", -- cyan
		"#f3f6f9", -- white
	},

	-- Bright Colors
	brights = {
		"#8a8fa3", -- black
		"#de8fa8", -- red
		"#d1e8a9", -- green
		"#fff7b1", -- yellow
		"#a3d4d5", -- blue
		"#ffaeea", -- magenta
		"#7fb4ca", -- cyan
		"#f3f6f9", -- white
	},
}

-- ┌──────────────────────────────────────────────────────────────────────────────┐
-- │                            WINDOWS (WSL)                                     │
-- └──────────────────────────────────────────────────────────────────────────────┘

-- Uncomment for Windows/WSL:
-- config.default_domain = 'WSL:Ubuntu'
-- config.front_end = "OpenGL"

-- WSL terminal fix: uncomment if fish/zsh fails with "missing or unsuitable terminal"
-- if wezterm.target_triple:find("windows") then
--   config.term = "xterm-256color"
-- end

return config
