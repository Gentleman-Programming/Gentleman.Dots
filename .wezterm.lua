-- Pull in the wezterm API
local wezterm = require 'wezterm'

-- This table will hold the configuration.
local config = {}

-- In newer versions of wezterm, use the config_builder which will
-- help provide clearer error messages
if wezterm.config_builder then
  config = wezterm.config_builder()
end

-- Kanagawa colors with dragon background
config.color_scheme = 'Kanagawa (Gogh)'
config.window_background_gradient = {
  colors = {'#181616'}
}

-- Set opacity, you may need "WebGpu" as your front_end
config.window_background_opacity = 0.95
config.font = wezterm.font 'IosevkaTerm NFM'

-- Only show one tab if needed --

config.hide_tab_bar_if_only_one_tab = true

-- For Linux / OSX (uncomment if needed)

-- config.default_prog = { '/home/linuxbrew/.linuxbrew/bin/fish', '-l' }

-- For windows -- (comment if you want to use Linux / OSX)

-- Start directly in WSL
config.default_domain = 'WSL:Ubuntu'

-- Full usage of GPU for rendering
config.front_end = "WebGpu"

-- Change fps to your monitor ones --
config.max_fps = 120

-- Fix for nvidia, opacity not working
for _, gpu in ipairs(wezterm.gui.enumerate_gpus()) do
	if gpu.backend == "Vulkan" then
		config.webgpu_preferred_adapter = gpu
		break
	end
end

-- and finally, return the configuration to wezterm
return config
