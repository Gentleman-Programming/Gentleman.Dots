-- Import the wezterm API
local wezterm = require("wezterm")

-- Initialize an empty configuration table
local config = {}

-- Background image options
-- config.background = {
-- 	{
-- 		source = {
-- 			File = "C:/Users/alanb/Pictures/synth-kanagawa-blur-20.jpg", -- Path to the background image file
-- 		},
-- 		width = "100%", -- Set the background image width to 100% of the terminal window
-- 		height = "100%", -- Set the background image height to gg100% of the terminal window
-- 		opacity = 1, -- Set the opacity of the background image (0.0 - 1.0)
-- 		hsb = {
-- 			brightness = 0.04, -- Set the brightness of the background image (low value to darken the image)
-- 			saturation = 1, -- Set the saturation of the background image
-- 		},
-- 	},
-- }

config.color_scheme = "Kanagawa (Gogh)"
-- This is where you actually apply your config choices
config.window_padding = {
	top = 0,
	right = 0,
	left = 0,
}

-- Set the terminal font
config.font = wezterm.font("IosevkaTerm NF")

-- Hide the tab bar if only one tab is open
config.hide_tab_bar_if_only_one_tab = true
config.max_fps = 240 -- hack for smoothness
config.enable_kitty_graphics = true

-- Background with Transparency
config.window_background_opacity = 0.85 -- Adjust this value as needed
config.macos_window_background_blur = 20 -- Adjust this value as needed
config.win32_system_backdrop = "Acrylic" -- Only Works in Windows

-- Font Size
config.font_size = 16.0

-- Smooth hack
config.max_fps = 240

-- Enable Kitty Graphics
config.enable_kitty_graphics = true

-- Disable Scroll Bar
config.enable_scroll_bar = false

-- activate ONLY if windows --

-- config.default_domain = 'WSL:Ubuntu'
-- config.front_end = "OpenGL"
-- local gpus = wezterm.gui.enumerate_gpus()
-- if #gpus > 0 then
--   config.webgpu_preferred_adapter = gpus[1] -- only set if there's at least one GPU
-- else
--   -- fallback to default behavior or log a message
--   wezterm.log_info("No GPUs found, using default settings")
-- end

-- and finally, return the configuration to wezterm

return config
