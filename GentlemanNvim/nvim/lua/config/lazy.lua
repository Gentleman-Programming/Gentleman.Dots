-- This file contains the configuration for setting up the lazy.nvim plugin manager in Neovim.

-- Node.js configuration - always use latest stable version
vim.g.node_host_prog = vim.fn.exepath("node") or "/usr/local/bin/node"
-- Ensure we're using a recent Node version for LSPs and plugins
if vim.fn.executable("node") == 1 then
  local node_version = vim.fn.system("node --version"):gsub("\n", "")
  print("Using Node.js version: " .. node_version)
end

-- Spell-checking
vim.opt.spell = true -- activa spell checker
vim.opt.spelllang = { "en" }

-- Define the path to the lazy.nvim plugin
local lazypath = vim.fn.stdpath("data") .. "/lazy/lazy.nvim"

-- Check if the lazy.nvim plugin is not already installed
if not vim.loop.fs_stat(lazypath) then
    -- Bootstrap lazy.nvim by cloning the repository
    -- stylua: ignore
    vim.fn.system({ "git", "clone", "--filter=blob:none", "https://github.com/folke/lazy.nvim.git", "--branch=stable",
        lazypath })
end

-- Prepend the lazy.nvim path to the runtime path
vim.opt.rtp:prepend(vim.env.LAZY or lazypath)

-- Fix copy and paste in WSL (Windows Subsystem for Linux)
if vim.fn.has("wsl") == 1 then
  vim.g.clipboard = {
    name = "win32yank", -- Use win32yank for clipboard operations
    copy = {
      ["+"] = "win32yank.exe -i --crlf", -- Command to copy to the system clipboard
      ["*"] = "win32yank.exe -i --crlf", -- Command to copy to the primary clipboard
    },
    paste = {
      ["+"] = "win32yank.exe -o --lf", -- Command to paste from the system clipboard
      ["*"] = "win32yank.exe -o --lf", -- Command to paste from the primary clipboard
    },
    cache_enabled = false, -- Disable clipboard caching
  }
end

-- FIX for Angular inline template highlighting

-- helper to get the smallest named TSNode at (row, col)
-- local function get_node_at_pos(bufnr, row, col)
--   return vim.treesitter.get_node({ bufnr = bufnr, pos = { row, col } })
-- end
--
-- -- override highlight handler
-- local orig = vim.lsp.handlers["textDocument/documentHighlight"]
-- vim.lsp.handlers["textDocument/documentHighlight"] = function(err, result, ctx, cfg)
--   if not result or vim.tbl_isempty(result) then
--     return orig(err, result, ctx, cfg)
--   end
--
--   local bufnr = ctx.bufnr
--   local filtered = {}
--   for _, h in ipairs(result) do
--     local s = h.range.start
--     local node = get_node_at_pos(bufnr, s.line, s.character)
--     local in_tpl = false
--     while node do
--       if node:type() == "template_string" then
--         in_tpl = true
--         break
--       end
--       node = node:parent()
--     end
--     if not in_tpl then
--       table.insert(filtered, h)
--     end
--   end
--   return orig(err, filtered, ctx, cfg)
-- end

-- Setup lazy.nvim with the specified configuration
require("lazy").setup({
  spec = {
    -- Add LazyVim and import its plugins
    { "LazyVim/LazyVim", import = "lazyvim.plugins" },
    -- Import any extra modules here
    -- Editor plugins
    { import = "lazyvim.plugins.extras.editor.harpoon2" },
    { import = "lazyvim.plugins.extras.editor.mini-files" },
    -- { import = "lazyvim.plugins.extras.editor.snacks_explorer" },
    { import = "lazyvim.plugins.extras.editor.snacks_picker" },

    -- Formatting plugins
    { import = "lazyvim.plugins.extras.formatting.biome" },
    { import = "lazyvim.plugins.extras.formatting.prettier" },

    -- Linting plugins
    { import = "lazyvim.plugins.extras.linting.eslint" },

    -- Language support plugins
    { import = "lazyvim.plugins.extras.lang.json" },
    { import = "lazyvim.plugins.extras.lang.markdown" },
    { import = "lazyvim.plugins.extras.lang.typescript" },
    { import = "lazyvim.plugins.extras.lang.angular" },
    { import = "lazyvim.plugins.extras.lang.astro" },
    { import = "lazyvim.plugins.extras.lang.go" },
    { import = "lazyvim.plugins.extras.lang.nix" },
    { import = "lazyvim.plugins.extras.lang.toml" },

    -- Coding plugins
    { import = "lazyvim.plugins.extras.coding.mini-surround" },
    { import = "lazyvim.plugins.extras.editor.mini-diff" },
    { import = "lazyvim.plugins.extras.coding.blink" },

    -- Utility plugins
    { import = "lazyvim.plugins.extras.util.mini-hipatterns" },

    -- AI plugins
    { import = "lazyvim.plugins.extras.ai.copilot" },
    { import = "lazyvim.plugins.extras.ai.copilot-chat" },

    -- Import/override with your plugins
    { import = "plugins" },
  },
  defaults = {
    -- By default, only LazyVim plugins will be lazy-loaded. Your custom plugins will load during startup.
    -- If you know what you're doing, you can set this to `true` to have all your custom plugins lazy-loaded by default.
    lazy = false,
    -- It's recommended to leave version=false for now, since a lot of the plugins that support versioning
    -- have outdated releases, which may break your Neovim install.
    version = false, -- Always use the latest git commit
    -- version = "*", -- Try installing the latest stable version for plugins that support semver
  },
  install = { colorscheme = { "tokyonight", "habamax" } }, -- Specify colorschemes to install
  checker = { enabled = true }, -- Automatically check for plugin updates
  performance = {
    rtp = {
      -- Disable some runtime path plugins to improve performance
      disabled_plugins = {
        "gzip",
        -- "matchit",
        -- "matchparen",
        -- "netrwPlugin",
        "tarPlugin",
        "tohtml",
        "tutor",
        "zipPlugin",
      },
    },
  },
})
