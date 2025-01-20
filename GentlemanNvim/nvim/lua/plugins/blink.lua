return {
  "saghen/blink.cmp",
  version = false,
  build = "cargo +nightly build --release",
  dependencies = {
    { "saghen/blink.compat", lazy = true, version = false },
  },
  opts = {
    sources = {
      -- LazyVim as custom option copmpat to pass in external sources with blink.compat
      compat = { "obsidian", "obsidian_new", "obsidian_tags" },
    },
  },
}
