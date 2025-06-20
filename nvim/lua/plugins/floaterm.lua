return {
  "voldikss/vim-floaterm",
  event = "VeryLazy",
  keys = {
    -- Toggle terminal
    { "<leader>ft", "<cmd>FloatermToggle<cr>", desc = "Toggle Terminal" },

    -- New terminals
    { "<leader>fT", "<cmd>FloatermNew<cr>", desc = "New Terminal" },

    -- Navigation
    { "<leader>fp", "<cmd>FloatermPrev<cr>", desc = "Prev Terminal" },
    { "<leader>fn", "<cmd>FloatermNext<cr>", desc = "Next Terminal" },
    { "<leader>fk", "<cmd>FloatermKill<cr>", desc = "Kill Terminal" },

    -- Special terminals
    { "<leader>fN", "<cmd>FloatermNew --wintype=normal --position=right<cr>", desc = "New Terminal (Right Split)" },

    -- Terminal mode mappings
    { "<Esc><Esc>", "<C-\\><C-n>", mode = "t", desc = "Switch to Normal mode" },
    { "<C-_>", "<cmd>FloatermToggle<cr>", mode = "t", desc = "Hide Terminal" },
    { "<M-f>", "<cmd>FloatermToggle<cr>", mode = "t", desc = "Toggle Terminal" },


    -- Normal mode mappings in terminal
    { "i", "i", mode = "n", buffer = true, desc = "Enter Insert mode" },
    --{ "q", "<cmd>FloatermToggle<cr>", mode = "n", buffer = true, desc = "Quit Terminal" },
  },
  config = function()
    -- Window appearance
    vim.g.floaterm_width = 0.5
    vim.g.floaterm_height = 0.5
    vim.g.floaterm_title = "Terminal ($1/$2)"
    vim.g.floaterm_position = "center"
    -- Border settings (rounded corners)
    vim.g.floaterm_borderchars = { "─", "│", "─", "│", "╭", "╮", "╯", "╰" }
    vim.g.floaterm_borderhighlight = "FloatermBorder"
    vim.g.floaterm_borderless = 0
    -- Terminal behavior
    vim.g.floaterm_autoclose = 1
    vim.g.floaterm_autohide = 1 -- Hide instead of close when using Esc
    -- Scrollback limit (critical for performance)
    vim.g.floaterm_scrollback = 1000 -- Lines to keep in buffer (default 10000)
    -- Shell configuration
    vim.g.floaterm_shell = vim.o.shell -- Use your default shell
    vim.g.floaterm_rootmarkers = { ".git", "package.json" } -- Auto-detect project root
    -- Keymaps for terminal job control
    vim.g.floaterm_keymap_toggle = "<leader>ft"
    vim.g.floaterm_keymap_kill = "<leader>fk"
    -- Highlight groups
    vim.api.nvim_set_hl(0, "FloatermBorder", { link = "FloatBorder" })
    vim.api.nvim_set_hl(0, "Floaterm", { bg = "#1e1e2e" }) -- Custom background
    -- Auto-commands
    vim.api.nvim_create_autocmd("TermOpen", {
      pattern = "floaterm",
      callback = function()
        vim.opt_local.number = false
        vim.opt_local.relativenumber = false
        vim.opt_local.signcolumn = "no"
        vim.cmd("startinsert")
      end,
    })
  end,
}
