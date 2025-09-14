return {
  "greggh/claude-code.nvim",
  dependencies = {
    "nvim-lua/plenary.nvim", -- Required for git operations
  },
  keys = {
    { "<leader>ac", "<cmd>ClaudeCode<cr>", desc = "Toggle Claude Code" },
    { "<leader>ar", "<cmd>ClaudeCodeResume<cr>", desc = "Resume conversation (picker)" },
    { "<leader>at", "<cmd>ClaudeCodeContinue<cr>", desc = "Continue recent conversation" },
    { "<leader>av", "<cmd>ClaudeCodeVerbose<cr>", desc = "Verbose logging" },
  },
  config = function()
    require("claude-code").setup()
  end,
}
