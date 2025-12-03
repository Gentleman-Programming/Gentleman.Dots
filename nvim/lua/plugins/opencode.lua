return {
  "NickvanDyke/opencode.nvim",
  dependencies = {
    { "folke/snacks.nvim", opts = { input = {}, picker = {}, terminal = {} } },
  },
  keys = {
    {
      "<leader>aa",
      function()
        require("opencode").toggle()
      end,
      mode = { "n" },
      desc = "Toggle OpenCode",
    },
    {
      "<leader>as",
      function()
        require("opencode").select({ submit = true })
      end,
      mode = { "n", "x" },
      desc = "OpenCode select",
    },
    {
      "<leader>ai",
      function()
        require("opencode").ask("", { submit = true })
      end,
      mode = { "n", "x" },
      desc = "OpenCode ask",
    },
    {
      "<leader>aI",
      function()
        require("opencode").ask("@this: ", { submit = true })
      end,
      mode = { "n", "x" },
      desc = "OpenCode ask with context",
    },
    {
      "<leader>ab",
      function()
        require("opencode").ask("@file ", { submit = true })
      end,
      mode = { "n", "x" },
      desc = "OpenCode ask about buffer",
    },
    {
      "<leader>ap",
      function()
        require("opencode").prompt("@this", { submit = true })
      end,
      mode = { "n", "x" },
      desc = "OpenCode prompt",
    },
    -- Built-in prompts
    {
      "<leader>ape",
      function()
        require("opencode").prompt("explain", { submit = true })
      end,
      mode = { "n", "x" },
      desc = "OpenCode explain",
    },
    {
      "<leader>apf",
      function()
        require("opencode").prompt("fix", { submit = true })
      end,
      mode = { "n", "x" },
      desc = "OpenCode fix",
    },
    {
      "<leader>apd",
      function()
        require("opencode").prompt("diagnose", { submit = true })
      end,
      mode = { "n", "x" },
      desc = "OpenCode diagnose",
    },
    {
      "<leader>apr",
      function()
        require("opencode").prompt("review", { submit = true })
      end,
      mode = { "n", "x" },
      desc = "OpenCode review",
    },
    {
      "<leader>apt",
      function()
        require("opencode").prompt("test", { submit = true })
      end,
      mode = { "n", "x" },
      desc = "OpenCode test",
    },
    {
      "<leader>apo",
      function()
        require("opencode").prompt("optimize", { submit = true })
      end,
      mode = { "n", "x" },
      desc = "OpenCode optimize",
    },
  },
  config = function()
    vim.g.opencode_opts = {
      provider = {
        snacks = {
          win = {
            position = "left",
          },
        },
      },
    }
    vim.o.autoread = true
  end,
}
