-- codeium with cmp
return {
  "nvim-cmp",
  dependencies = {
    -- codeium
    {
      "Exafunction/codeium.nvim",
      cmd = "Codeium",
      build = ":Codeium Auth",
      opts = {},
    },
  },
  ---@param opts cmp.ConfigSchema
  opts = function(_, opts)
    table.insert(opts.sources, 1, {
      name = "codeium",
      group_index = 1,
      priority = 100,
    })
  end,
}

-- codeium inline
-- return {
--   "Exafunction/codeium.vim",
--   config = function()
--     vim.keymap.set("i", "<C-g>", function()
--       return vim.fn["codeium#Accept"]()
--     end, { expr = true })
--     vim.keymap.set("i", "<C-l>", function()
--       return vim.fn["codeium#CycleCompletions"](1)
--     end, { expr = true })
--     vim.keymap.set("i", "<C-M>", function()
--       return vim.fn["codeium#Complete"]()
--     end, { expr = true })
--     vim.keymap.set("i", "<C-x>", function()
--       return vim.fn["codeium#Clear"]()
--     end, { expr = true })
--   end,
-- }
