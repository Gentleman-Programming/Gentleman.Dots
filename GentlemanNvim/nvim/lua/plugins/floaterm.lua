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

    -- Normal mode mappings in terminal
    { "i", "i", mode = "n", buffer = true, desc = "Enter Insert mode" },
    { "q", "<cmd>FloatermToggle<cr>", mode = "n", buffer = true, desc = "Quit Terminal" },
  },
  config = function()
    vim.g.floaterm_width = 0.6
    vim.g.floaterm_height = 0.6
    vim.g.floaterm_title = "Terminal ($1/$2)"
    vim.g.floaterm_autoclose = 1

    -- Auto-enter insert mode when terminal opens
    vim.api.nvim_create_autocmd("TermOpen", {
      pattern = "floaterm",
      callback = function()
        vim.cmd("startinsert")
      end,
    })
  end,
}
