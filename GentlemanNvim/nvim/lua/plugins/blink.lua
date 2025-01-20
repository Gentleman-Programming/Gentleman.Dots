return {
  "saghen/blink.cmp",
  version = "*",
  dependencies = {
    { "saghen/blink.compat", lazy = true, version = "*" },
  },
  opts = {
    sources = {
      -- LazyVim as custom option copmpat to pass in external sources with blink.compat
      compat = { "obsidian", "obsidian_new", "obsidian_tags" },
    },
  },
}
