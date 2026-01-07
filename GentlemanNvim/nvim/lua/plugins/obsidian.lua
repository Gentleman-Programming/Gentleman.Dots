return {
  "obsidian-nvim/obsidian.nvim",
  version = "*", -- recommended, use latest release instead of latest commit
  lazy = false,
  enabled = function()
    -- Disable Obsidian when running from Oil Simple (to avoid path issues in Zed context)
    return not vim.g.disable_obsidian
  end,
  dependencies = {
    -- Required.
    "nvim-lua/plenary.nvim",
  },
  opts = {
    workspaces = {
      {
        name = "GentlemanNotes", -- Name of the workspace
        path = os.getenv("HOME") .. "/.config/obsidian", -- Path to the notes directory
      },
    },
    completion = {
      cmp = true,
    },
    picker = {
      -- Set your preferred picker. Can be one of 'telescope.nvim', 'fzf-lua', 'mini.pick' or 'snacks.pick'.
      name = "snacks.pick",
    },
    -- Optional, define your own callbacks to further customize behavior.
    callbacks = {
      -- Runs anytime you enter the buffer for a note.
      -- FIX: Breaking change from obsidian-nvim/obsidian.nvim
      -- The callback now only receives 1 parameter (note) instead of 2 (client, note)
      enter_note = function(note)
        if not note then return end

        -- Setup keymaps for obsidian notes
        vim.keymap.set("n", "gf", function()
          return require("obsidian").util.gf_passthrough()
        end, { buffer = note.bufnr, expr = true, desc = "Obsidian follow link" })

        vim.keymap.set("n", "<leader>ch", function()
          return require("obsidian").util.toggle_checkbox()
        end, { buffer = note.bufnr, desc = "Toggle checkbox" })

        vim.keymap.set("n", "<cr>", function()
          return require("obsidian").util.smart_action()
        end, { buffer = note.bufnr, expr = true, desc = "Obsidian smart action" })
      end,
    },

    -- Settings for templates
    templates = {
      subdir = "templates", -- Subdirectory for templates
      date_format = "%Y-%m-%d-%a", -- Date format for templates
      gtime_format = "%H:%M", -- Time format for templates
      tags = "", -- Default tags for templates
    },
  },
}
