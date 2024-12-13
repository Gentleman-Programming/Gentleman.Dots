return {
  { "saghen/blink.compat" },
  {
    "saghen/blink.cmp",
    dependencies = {
      { "epwalsh/obsidian.nvim" },
    },
    opts_extend = { "sources.completion.enabled_providers" },
    ---@module 'blink.cmp'
    ---@type blink.cmp.Config
    opts = {
      sources = {
        completion = {
          enabled_providers = {
            "lsp",
            "path",
            "snippets",
            "buffer",
            "obsidian",
            "obsidian_new",
            "obsidian_tags",
          },
        },
        providers = {
          obsidian = {
            name = "obsidian",
            module = "blink.compat.source",
          },
          obsidian_new = {
            name = "obsidian_new",
            module = "blink.compat.source",
          },
          obsidian_tags = {
            name = "obsidian_tags",
            module = "blink.compat.source",
          },
        },
      },
    },
  },
}
