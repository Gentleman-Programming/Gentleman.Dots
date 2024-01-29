# Gentleman.Dots

## Description

This repository contains customized configurations for the Neovim development environment, including specific plugins and keymaps to enhance productivity. It makes use of [LazyVim](https://github.com/LazyVim/LazyVim) as a preconfigured set of plugins and settings to simplify the use of Neovim.

## Folder `GentlemanNvim`

### Configuration Transfer

```bash
git clone https://github.com/Gentleman-Programming/Gentleman.Dots
cp -r Gentleman.Dots/GentlemanNvim/* ~/.config
```

Restart Neovim to apply the changes.

### Folder `plugins`

#### File `codeium.lua`

This file configures the [codeium.vim](https://github.com/Exafunction/codeium.vim) plugin, providing keyboard shortcuts for accepting, completing, and clearing suggestions.

```lua
return {
  "Exafunction/codeium.vim",
  config = function()
    vim.keymap.set("i", "<C-g>", function()
      return vim.fn["codeium#Accept"]()
    end, { expr = true })

    vim.keymap.set("i", "<C-l>", function()
      return vim.fn["codeium#CycleCompletions"](1)
    end, { expr = true })

    vim.keymap.set("i", "<C-M>", function()
      return vim.fn["codeium#Complete"]()
    end, { expr = true })

    vim.keymap.set("i", "<C-x>", function()
      return vim.fn["codeium#Clear"]()
    end, { expr = true })
  end,
}
```

#### File `vim-tmux-navigation.lua`

This file configures the [nvim-tmux-navigation.vim](https://github.com/alexghergh/nvim-tmux-navigation) plugin, providing keyboard shortcuts for navigating between tmux and nvim in an optimal way.

```lua
return {
  "alexghergh/nvim-tmux-navigation",
}
```

#### File `colorscheme.lua`

This file configures the color scheme using the [nvim](https://github.com/catppuccin/nvim) plugin. The "kanagawa-dragon" theme with specific options is chosen. You can also choose catppucin or modus by changing the property `colorscheme = "kanagawa-dragon"`. If you want a transparent background, use `:TransparentEnable`, and NVIM will have extra opacity according to your terminal settings.

```lua
return {
  {
    "catppuccin/nvim",
    name = "catppuccin",
    lazy = false,
    opts = {
      transparent_background = true,
      flavour = "mocha",
    },
    integrations = {
      cmp = true,
      gitsigns = true,
      nvimtree = true,
      treesitter = true,
      notify = false,
      mini = {
        enabled = true,
        indentscope_color = "",
      },
      -- For more plugins integrations please scroll down (https://github.com/catppuccin/nvim#integrations)
    },
  },
  {
    "miikanissi/modus-themes.nvim",
    name = "modus",
    priority = 1000,
  },
  {
    "rebelot/kanagawa.nvim",
    name = "kanagawa",
    opts = {
      transparent_background = true,
    },
    priority = 1000,
  },
  {
    "xiyaowong/transparent.nvim",
  },
  {
    "LazyVim/LazyVim",
    opts = {
      colorscheme = "kanagawa-dragon",
    },
  },
}
```

#### File `editor.lua`

This file configures various plugins to enhance the editing experience, such as highlighting patterns in Markdown files and advanced search tools with Telescope.

```lua
1. **mini.hipatterns**:
   - Plugin: `echasnovski/mini.hipatterns`
   - Event: `BufReadPre`
   - Configuration:
     - Configures a highlighter for HSL colors. HSL colors in the code will be highlighted with a background color corresponding to the HSL color.

2. **git.nvim**:
   - Plugin: `dinhhuy258/git.nvim`
   - Event: `BufReadPre`
   - Configuration:
     - Configures keyboard shortcuts to open a blame window (`<Leader>gb`) and to open a file or folder in the Git repository (`<Leader>go`).

3. **telescope.nvim**:
   - Plugin: `telescope.nvim`
   - Dependencies: `nvim-telescope/telescope-fzf-native.nvim` and `nvim-telescope/telescope-file-browser.nvim`
   - Configuration:
     - Configures a series of keyboard shortcuts for various functionalities, such as searching for files, searching for a string in the current directory, listing open buffers, listing help tags, and resuming the previous telescope picker.
     - Also configured to open a file browser with the path of the current buffer with `<Leader>sf`.
   - Additional Configuration:
     - Configures results to wrap, layout strategy to be horizontal, prompt position at the top, and sorting strategy to be ascending.
     - Configures the diagnostics selector to have the "ivy" theme, initial mode to be "normal", and preview cutoff to be 9999.
     - Configures the file browser to have the "dropdown" theme, hijack netrw and use it instead, and have its own mappings.
```

#### File `harpoon.lua`

This file configures the [harpoon](https://github.com/ThePrimeagen/harpoon) plugin to facilitate navigation between marked files.

```lua
return {
  "ThePrimeagen/harpoon",
  lazy = false,
  dependencies = {
    "nvim-lua/plenary.nvim",
  },
  branch = "harpoon2",
  config = true,
}
```

#### File `telescope.lua`

This file configures the [Telescope](https://github.com/nvim-telescope/telescope.nvim) plugin for advanced searches in files and other resources.

```lua
return {
  "nvim-telescope/telescope.nvim",
  opts = {
    defaults = {
      layout_strategy = "vertical",
      layout_config = { preview_cutoff = 6 },
    },
  },
}
```

#### File `ui.lua`

This file configures various plugins to enhance the user interface, including notifications, animations, buffer lines, and status lines.

```lua
1. **Noice.nvim**:
   - Plugin: `folke/noice.nvim`
   - Configuration:


     - Adds a path to the noice configuration to filter notification messages with the text "No information available". These messages will be skipped.
     - Sets autocmds to detect when the Neovim window gains or loses focus. This is used to determine if the interface is focused or not.
     - Adds an additional path to display system notifications when Neovim loses focus.

2. **Nvim-notify**:
   - Plugin: `rcarriga/nvim-notify`
   - Configuration:
     - Sets the background color and timeout for notifications.

3. **Mini.animate**:
   - Plugin: `echasnovski/mini.animate`
   - Configuration:
     - Disables the scroll animation (`scroll`).

4. **Bufferline.nvim**:
   - Plugin: `akinsho/bufferline.nvim`
   - Configuration:
     - Defines keyboard shortcuts to switch between tabs.
     - Configures options to show icons and close tabs.

5. **Lualine.nvim**:
   - Plugin: `nvim-lualine/lualine.nvim`
   - Configuration:
     - Configures the "catppuccin" theme for the status line.

6. **Incline.nvim**:
   - Plugin: `b0o/incline.nvim`
   - Configuration:
     - Configures colors and visual options for highlighting file names in the status line.

7. **Zen-mode.nvim**:
   - Plugin: `folke/zen-mode.nvim`
   - Configuration:
     - Configures keyboard shortcuts to activate "Zen Mode," which hides UI elements to focus on text editing.

8. **Dashboard-nvim**:
   - Plugin: `nvimdev/dashboard-nvim`
   - Configuration:
     - Sets a custom logo for the Neovim startup dashboard.
```

#### File `keymaps.lua`

This file defines some custom keymaps to improve navigation, text manipulation in insert mode and plugin shortcuts.

```lua
-- Keymaps are automatically loaded on the VeryLazy event
-- Default keymaps that are always set: https://github.com/LazyVim/LazyVim/blob/main/lua/lazyvim/config/keymaps.lua Add any additional keymaps here

vim.keymap.set("i", "<C-d>", "<C-d>zz")
vim.keymap.set("i", "<C-u>", "<C-u>zz")
vim.keymap.set("i", "<C-b>", "<C-o>de")

----- Tmux Navigation ------
local nvim_tmux_nav = require("nvim-tmux-navigation")

vim.keymap.set("n", "<C-h>", nvim_tmux_nav.NvimTmuxNavigateLeft)
vim.keymap.set("n", "<C-j>", nvim_tmux_nav.NvimTmuxNavigateDown)
vim.keymap.set("n", "<C-k>", nvim_tmux_nav.NvimTmuxNavigateUp)
vim.keymap.set("n", "<C-l>", nvim_tmux_nav.NvimTmuxNavigateRight)
vim.keymap.set("n", "<C-\\>", nvim_tmux_nav.NvimTmuxNavigateLastActive)
vim.keymap.set("n", "<C-Space>", nvim_tmux_nav.NvimTmuxNavigateNext)

----- Harpoon 2 -----
local harpoon = require("harpoon")

-- REQUIRED
harpoon:setup()
-- REQUIRED

vim.keymap.set("n", "<leader>a", function()
  harpoon:list():append()
end, { desc = "Add harpoon mark" })

vim.keymap.set("n", "<C-e>", function()
  harpoon.ui:toggle_quick_menu(harpoon:list())
end)

vim.keymap.set("n", "<C-M-h>", function()
  harpoon:list():select(1)
end)

vim.keymap.set("n", "<C-M-j>", function()
  harpoon:list():select(2)
end)

vim.keymap.set("n", "<C-M-k>", function()
  harpoon:list():select(3)
end)

vim.keymap.set("n", "<C-M-l>", function()
  harpoon:list():select(4)
end)

-- Disable key mappings in insert mode
vim.api.nvim_set_keymap("i", "<A-j>", "<Nop>", { noremap = true, silent = true })
vim.api.nvim_set_keymap("i", "<A-k>", "<Nop>", { noremap = true, silent = true })

-- Disable key mappings in normal mode
vim.api.nvim_set_keymap("n", "<A-j>", "<Nop>", { noremap = true, silent = true })
vim.api.nvim_set_keymap("n", "<A-k>", "<Nop>", { noremap = true, silent = true })

-- Disable key mappings in visual block mode
vim.api.nvim_set_keymap("x", "<A-j>", "<Nop>", { noremap = true, silent = true })
vim.api.nvim_set_keymap("x", "<A-k>", "<Nop>", { noremap = true, silent = true })
vim.api.nvim_set_keymap("x", "J", "<Nop>", { noremap = true, silent = true })
vim.api.nvim_set_keymap("x", "K", "<Nop>", { noremap = true, silent = true })
```

## Folder `GentlemanKitty`

### File `kanagawa.nvim`

This file configures the Kanagawa theme in the Kitty terminal, providing visual adjustments and keyboard shortcuts for tab navigation.

```vim
# vim:fileencoding=utf-8:foldmethod=marker

#: Fonts {{{

font_family      IosevkaTerm Nerd Font
font_size 14.0

#: Foreground and background colors.

background_opacity 0.95
# background_blur 0

## name: Kanagawa
## license: MIT
## author: Tommaso Laurenzi
## upstream: https://github.com/rebelot/kanagawa.nvim/

background #0d0c0c
foreground #DCD7BA
selection_background #2D4F67
selection_foreground #C8C093
url_color #72A7BC
cursor #C8C093

# Tabs
active_tab_background #1F1F28
active_tab_foreground #C8C093
inactive_tab_background  #1F1F28
inactive_tab_foreground #727169
#tab_bar_background #15161E

# normal
color0 #16161D
color1 #C34043
color2 #76946A
color3 #C0A36E
color4 #7E9CD8
color5 #957FB8
color6 #6A9589
color7 #C8C093

# bright
color8  #727169
color9  #E82424
color10 #98BB6C
color11 #E6C384
color12 #7FB4CA
color13 #938AA9
color14 #7AA89F
color15 #DCD7BA


# extended colors
color16 #FFA066
color17 #FF5D62


map cmd+1 goto_tab 1
map cmd+2 goto_tab 2
map cmd+3 goto_tab 3
map cmd+4 goto_tab 4
map cmd+5 goto_tab 5
map cmd+6 goto_tab 6
map cmd+7 goto_tab 7
map cmd+8 goto_tab 8
map cmd+9 goto_tab 9


# make option key work for alt-f / alt-b
macos_option_as_alt yes
```

This file provides the configuration for the Kanagawa theme in Neovim, using the IosevkaTerm Nerd Font with a font size of 14.0. It defines a carefully selected color palette to enhance the coding experience. The settings include tab styles for active and inactive tabs, along with key mappings for quick navigation between tabs.

### Configuration Transfer

```bash
git clone https://github.com/Gentleman-Programming/Gentleman.Dots
cp -r Gentleman.Dots/GentlemanKitty/* ~/.config/kitty
```

**Theme Details:**

- **Name:** Kanagawa
- **Author:** Tommaso Laurenzi
- **License:** MIT
- **Upstream Repository:** [Kanagawa.nvim](https://github.com/rebelot/kanagawa.nvim/)

**Note:** The provided key mappings for navigating between tabs are configured as `cmd+1` to `cmd+9`.

## Folder `GentlemanFish`

### Fish Installation

#### HomeBrew (Recommended)

`brew install fish`

#### Ubuntu/Debian

```
sudo apt-get update
sudo apt-get install fish
```

#### Fedora

`sudo dnf install fish`

### Oh My Fish Installation

`curl https://raw.githubusercontent.com/oh-my-fish/oh-my-fish/master/bin/install | fish`

### Configuration Transfer

```bash
git clone https://github.com/Gentleman-Programming/Gentleman.Dots
cp -r Gentleman.Dots/GentlemanFish/* ~/.config
```

### Configure path for PJ plugin working folders in Oh My Fish

Go to the file `~/.config/fish/fish_variables` and change the following variable to the path to your working folder with your projects:

`SETUVAR --export PROJECT_PATHS: /YourWorkingPath`

### Choose Kanagawa theme for Fish

Run: `fish_config theme save Kanagawa`

When asked if you want to overwrite: `Y` and then press enter

## Folder `GentlemanTmux`

Contains configurations for the tmux environment. To install and use it, follow these steps:

### Tmux Installation

#### HomeBrew (Recommended)

`brew install tmux`

#### Ubuntu/Debian

```
sudo apt-get update
sudo apt-get install tmux
```

#### Fedora

`sudo dnf -y install tmux`

### Configuration Transfer
```bash
git clone https://github.com/Gentleman-

Programming/Gentleman.Dots
Uncompress and remove after the "Plugins.zip" inside the folder

cp -r Gentleman.Dots/GentlemanTmux/* ~/
```

### Start Tmux

#### Launch it

```bash
tmux
```

#### Load the configuration

```bash
tmux source-file ~/.tmux.conf
```

### Load Tmux plugins

```bash
<Ctrl-b> + I to load the plugins
```

### If you want Tmux to run by default when opening the terminal

#### Open `~/.config/fish/config.fish` and add the following line at the end:

```bash
if status is-interactive
    and not set -q TMUX
    exec tmux
end
```

### Configuration Explanation

1. **Default Shell Configuration:**

   ```bash
   set-option -g default-shell /usr/bin/fish
   ```

   Sets the default shell that Tmux will use as `/usr/bin/fish`.

2. **Plugin Configuration:**

- Other plugins used, such as the Tmux Plugin Manager (`tpm`) and default sensible plugins.

  ```bash
  set -g @plugin 'tmux-plugins/tpm'
  set -g @plugin 'tmux-plugins/tmux-sensible'
  set -g @plugin 'tmux-plugins/tmux-resurrect'
  set -g @plugin 'christoomey/vim-tmux-navigator'
  ```

- Note tmux-resurrect, which saves the session state so that it's not lost, is used by:

  ```bash
  <Ctrl-b> + <Ctrl-s> to save the state
  <Ctrl-b> + <Ctrl-r> to restore the state
  ```

- Note vim-tmux-navigator, which allows switching between splits in Vim and Tmux interchangeably using `<Ctrl-h/j/k/l>`:

  ```bash
  set -g @plugin 'christoomey/vim-tmux-navigator'
  ```

- Configuration of the default terminal type and additional settings for scrolling.
  ```bash
  set -g default-terminal "tmux-256color"
  set-option -ga terminal-overrides ",xterm*:Tc"
  ```

````

3. **Appearance Configuration for Windows and Panes:**
 ```bash
 set -g @catppuccin_window_left_separator "█"
 set -g @catppuccin_window_right_separator "█

 "
 set -g @catppuccin_window_number_position "right"
 set -g @catppuccin_window_middle_separator "  █"
 set -g @catppuccin_window_default_fill "number"
 set -g @catppuccin_window_current_fill "number"
 set -g @catppuccin_window_current_text "#{pane_current_path}"
````

- Configuration of the appearance of windows and panes, including separators, number position, and text of the current window.

4. **Status Bar Configuration:**

   ```bash
   set -g @catppuccin_status_modules "application session date_time"
   set -g @catppuccin_status_left_separator  ""
   set -g @catppuccin_status_right_separator " "
   set -g @catppuccin_status_right_separator_inverse "yes"
   set -g @catppuccin_status_fill "all"
   set -g @catppuccin_status_connect_separator "no"
   set -g @catppuccin_directory_text "#{pane_current_path}"
   ```

   - Configuration of modules and appearance of the status bar, including separators and text of the current directory.

5. **Initialization of Tmux Plugin Manager (TPM):**
   ```bash
   run '~/.tmux/plugins/tpm/tpm'
   ```
   - Initiates the Tmux Plugin Manager. This command should be kept at the end of the Tmux configuration file.

Enjoy your new Neovim development environment!
