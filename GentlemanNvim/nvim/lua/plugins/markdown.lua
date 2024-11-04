-- This file contains the configuration for the markview.nvim plugin in Neovim.

return {
  {
    -- Plugin: markview.nvim
    -- URL: https://github.com/OXY2DEV/markview.nvim
    -- Description: A Neovim plugin for previewing markdown files.
    "OXY2DEV/markview.nvim",
    lazy = false, -- Load this plugin immediately (recommended)
    -- ft = "markdown" -- Uncomment this line if you decide to lazy-load the plugin for markdown files only

    dependencies = {
      -- Dependency: nvim-treesitter
      -- URL: https://github.com/nvim-treesitter/nvim-treesitter
      -- Description: Neovim Treesitter configurations and abstraction layer.
      -- Note: You will not need this if you installed the parsers manually or if the parsers are in your $RUNTIMEPATH.
      "nvim-treesitter/nvim-treesitter",

      -- Dependency: nvim-web-devicons
      -- URL: https://github.com/nvim-tree/nvim-web-devicons
      -- Description: A Lua fork of vim-web-devicons for Neovim.
      "nvim-tree/nvim-web-devicons",
    },
  },
}
