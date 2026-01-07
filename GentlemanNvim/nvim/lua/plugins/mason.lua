-- Configuración de Mason para Neovim
-- Objetivo: Asegurar que markdownlint-cli2 esté instalado automáticamente

return {
  {
    "mason-org/mason.nvim",
    opts = function(_, opts)
      opts.ensure_installed = opts.ensure_installed or {}
      vim.list_extend(opts.ensure_installed, {
        "markdownlint-cli2",
      })
    end,
  },
}
