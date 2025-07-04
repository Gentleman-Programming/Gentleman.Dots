-- This file contains the configuration overrides for specific Neovim plugins.

return {

  -- Change configuration for trouble.nvim

  {

    -- Plugin: trouble.nvim

    -- URL: https://github.com/folke/trouble.nvim

    -- Description: A pretty list for showing diagnostics, references, telescope results, quickfix and location lists.

    "folke/trouble.nvim",

    -- Options to be merged with the parent specification

    opts = { use_diagnostic_signs = true }, -- Use diagnostic signs for trouble.nvim
  },

  -- Add symbols-outline.nvim plugin

  {

    -- Plugin: symbols-outline.nvim

    -- URL: https://github.com/simrat39/symbols-outline.nvim

    -- Description: A tree like view for symbols in Neovim using the Language Server Protocol.

    "simrat39/symbols-outline.nvim",

    cmd = "SymbolsOutline", -- Command to open the symbols outline

    keys = { { "<leader>cs", "<cmd>SymbolsOutline<cr>", desc = "Symbols Outline" } }, -- Keybinding to open the symbols outline

    config = true, -- Use default configuration
  },

  -- Remove inlay hints from default configuration

  {

    -- Plugin: nvim-lspconfig

    -- URL: https://github.com/neovim/nvim-lspconfig

    -- Description: Quickstart configurations for the Neovim LSP client.

    "neovim/nvim-lspconfig",

    event = "VeryLazy", -- Load this plugin on the 'VeryLazy' event

    opts = {

      inlay_hints = { enabled = false }, -- Disable inlay hints

      servers = {

        angularls = {

          -- Configuration for Angular Language Server

          root_dir = function(fname)
            return require("lspconfig.util").root_pattern("angular.json", "project.json")(fname)
          end,
        },

        nil_ls = {

          -- Configuration for nil (Nix Language Server), already installed via nix

          cmd = { "nil" },

          autostart = true,

          mason = false, -- Explicitly disable mason management for nil_ls

          settings = {

            ["nil"] = {

              formatting = { command = { "nixpkgs-fmt" } },
            },
          },
        },
        -- Configuración para OmniSharp (C# Language Server)
        omnisharp = {
          settings = {
            Formatting = {
              -- Deshabilitamos el formateo de OmniSharp para que CSharpier tome el control
              Enable = true,
              EnableEditorConfigSupport = true, -- Igual dejas que respete .editorconfig
            },
            RoslynExtensionsOptions = {
              EnableDecompilationSupport = true, -- Soporte para descompilación
              EnableAnalyzersSupport = true, -- Habilita análisis de código Roslyn
              EnableCodeActions = true, -- Habilita las acciones de código
            },
          },
          -- Asegurarse de que OmniSharp detecte correctamente la raíz del proyecto
          -- Esto es crucial para proyectos .NET. Por defecto, lspconfig ya maneja .sln
          -- pero puedes personalizarlo si tus proyectos tienen una estructura particular.
          root_dir = require("lspconfig.util").root_pattern("*.sln", ".git", "Directory.Build.props", "omnisharp.json"),
        },
      },
    },
  },
}
