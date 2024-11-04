-- This file contains the configuration for various Neovim plugins related to the editor.

return {
  {
    -- Plugin: goto-preview
    -- URL: https://github.com/rmagatti/goto-preview
    -- Description: Provides preview functionality for definitions, declarations, implementations, type definitions, and references.
    "rmagatti/goto-preview",
    event = "BufEnter", -- Load the plugin when a buffer is entered
    config = true, -- Enable default configuration
    keys = {
      {
        "gpd",
        "<cmd>lua require('goto-preview').goto_preview_definition()<CR>",
        noremap = true, -- Do not allow remapping
        desc = "goto preview definition", -- Description for the keybinding
      },
      {
        "gpD",
        "<cmd>lua require('goto-preview').goto_preview_declaration()<CR>",
        noremap = true,
        desc = "goto preview declaration",
      },
      {
        "gpi",
        "<cmd>lua require('goto-preview').goto_preview_implementation()<CR>",
        noremap = true,
        desc = "goto preview implementation",
      },
      {
        "gpy",
        "<cmd>lua require('goto-preview').goto_preview_type_definition()<CR>",
        noremap = true,
        desc = "goto preview type definition",
      },
      {
        "gpr",
        "<cmd>lua require('goto-preview').goto_preview_references()<CR>",
        noremap = true,
        desc = "goto preview references",
      },
      {
        "gP",
        "<cmd>lua require('goto-preview').close_all_win()<CR>",
        noremap = true,
        desc = "close all preview windows",
      },
    },
  },
  {
    -- Plugin: mini.hipatterns
    -- URL: https://github.com/echasnovski/mini.hipatterns
    -- Description: Provides highlighter patterns for various text patterns.
    "echasnovski/mini.hipatterns",
    event = "BufReadPre", -- Load the plugin before reading a buffer
    opts = {
      highlighters = {
        hsl_color = {
          pattern = "hsl%(%d+,? %d+,? %d+%)", -- Pattern to match HSL color values
          group = function(_, match)
            local utils = require("config.gentleman.utils")
            local h, s, l = match:match("hsl%((%d+),? (%d+),? (%d+)%)")
            h, s, l = tonumber(h), tonumber(s), tonumber(l)
            local hex_color = utils.hslToHex(h, s, l)
            return MiniHipatterns.compute_hex_color_group(hex_color, "bg")
          end,
        },
      },
    },
  },
  {
    -- Plugin: git.nvim
    -- URL: https://github.com/dinhhuy258/git.nvim
    -- Description: Provides Git integration for Neovim.
    "dinhhuy258/git.nvim",
    event = "BufReadPre", -- Load the plugin before reading a buffer
    opts = {
      keymaps = {
        blame = "<Leader>gb", -- Keybinding to open blame window
        browse = "<Leader>go", -- Keybinding to open file/folder in git repository
      },
    },
  },
  {
    -- Plugin: telescope.nvim
    -- URL: https://github.com/nvim-telescope/telescope.nvim
    -- Description: A highly extendable fuzzy finder over lists.
    "nvim-telescope/telescope.nvim",
    opts = function(_, opts)
      local actions = require("telescope.actions")

      opts.defaults = {
        path_display = { "smart" }, -- Display paths smartly
        file_ignore_patterns = {
          "node_modules",
          "package-lock.json",
          "yarn.lock",
          "bun.lockb",
        },
        prompt_prefix = "> ", -- Set the prompt to just ">"
        layout_strategy = "horizontal", -- Use horizontal layout
        sorting_strategy = "ascending", -- Sort results in ascending order
        winblend = 0, -- No transparency
        results_title = false, -- Remove the "Results" title
        borderchars = {
          prompt = { "─", " ", " ", " ", " ", " ", " ", " " }, -- Top border for the prompt only
          results = { " ", " ", " ", " ", " ", " ", " ", " " }, -- No borders for results
          preview = { "─", "│", " ", "│", "╭", "╮", "", "" }, -- Borders for the preview (top and sides)
        },
        mappings = {
          i = {
            ["<C-Down>"] = actions.cycle_history_next, -- Cycle to next history item
            ["<C-Up>"] = actions.cycle_history_prev, -- Cycle to previous history item
            ["<C-f>"] = actions.preview_scrolling_down, -- Scroll preview down
            ["<C-b>"] = actions.preview_scrolling_up, -- Scroll preview up
          },
          n = {
            ["q"] = actions.close, -- Close the telescope window
          },
        },
      }

      -- Load the fzf extension for fast searches
      require("telescope").load_extension("fzf")

      -- Add hidden files and no-ignore options to file search and live_grep
      opts.pickers = {
        find_files = {
          find_command = { "rg", "--files", "--hidden", "--no-ignore", "--iglob", "!.git/" },
        },
        live_grep = {
          additional_args = function()
            return { "--hidden", "--no-ignore" }
          end,
        },
      }
      return opts
    end,

    dependencies = {
      {
        -- Plugin: telescope-live-grep-args.nvim
        -- URL: https://github.com/nvim-telescope/telescope-live-grep-args.nvim
        -- Description: Adds live grep arguments to Telescope.
        "nvim-telescope/telescope-live-grep-args.nvim",
        version = "^1.0.0",
        config = function()
          require("telescope").load_extension("live_grep_args")
        end,
      },
      {
        -- Plugin: telescope-fzf-native.nvim
        -- URL: https://github.com/nvim-telescope/telescope-fzf-native.nvim
        -- Description: FZF sorter for Telescope written in C.
        "nvim-telescope/telescope-fzf-native.nvim",
        build = "make", -- Build the plugin using make
        config = function()
          require("telescope").load_extension("fzf")
        end,
      },
    },
    config = function(_, opts)
      require("telescope").setup(opts)

      -- Keybinding to open live grep with arguments
      vim.keymap.set("n", "<leader>fg", ":lua require('telescope').extensions.live_grep_args.live_grep_args()<CR>")
    end,
  },
}
