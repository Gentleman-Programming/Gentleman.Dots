-- Configuración de formateo para Neovim
-- Objetivo: Usar markdownlint-cli2 como único formateador para Markdown
-- Mantiene Prettier para JS, TS, JSON y otros lenguajes

return {
  {
    "stevearc/conform.nvim",
    opts = {
      -- Magic here: Al especificar esto, eliminamos "prettier" SOLO para markdown
      -- El resto de lenguajes siguen usando lo que LazyVim tenga por defecto
      formatters_by_ft = {
        markdown = { "markdownlint-cli2" },
        ["markdown.mdx"] = { "markdownlint-cli2" },
      },

      formatters = {
        ["markdownlint-cli2"] = {
          -- Forzamos a que use tu archivo global y el modo --fix
          -- vim.fn.expand("~") asegura que encuentre la ruta correcta a tu home
          prepend_args = {
            "--config",
            vim.fn.expand("~/.config/GentlemanNvim/.markdownlint-cli2.yaml"),
            "--fix"
          },
        },
      },
    },
  },
}
