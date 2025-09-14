-- This file contains the configuration for the nvim-rip-substitute plugin in Neovim.

return {
  -- Plugin: nvim-rip-substitute
  -- URL: https://github.com/chrisgrieser/nvim-rip-substitute
  -- Description: A Neovim plugin for performing substitutions with ripgrep.
  "chrisgrieser/nvim-rip-substitute",

  cmd = "RipSubstitute", -- Command to trigger the plugin

  keys = {
    {
      -- Keybinding to perform a ripgrep substitution
      "<leader>fs",
      function()
        require("rip-substitute").sub() -- Call the substitution function from the plugin
      end,
      mode = { "n", "x" }, -- Enable the keybinding in normal and visual modes
      desc = " rip substitute", -- Description for the keybinding
    },
  },
}
