# Tools Reference

Detailed overview of all tools configured by Gentleman.Dots.

## Terminal Emulators

| Tool | Description |
|------|-------------|
| **Ghostty** | GPU-accelerated, native, blazing fast |
| **Kitty** | Feature-rich, GPU-based rendering |
| **WezTerm** | Lua-configurable, cross-platform |
| **Alacritty** | Minimal, Rust-based, lightweight |

## Shells

| Tool | Description |
|------|-------------|
| **Nushell** | Structured data, modern syntax, pipelines |
| **Fish** | User-friendly, great defaults, no config needed |
| **Zsh** | Highly customizable, POSIX-compatible, Powerlevel10k |

Fish's managed `gentleman` and `gentleman-cute` profiles also export listing
colors when loaded at startup or through `fish-theme gentleman` /
`fish-theme gentleman-cute`. GNU `ls` uses `LS_COLORS` with the profile's RGB
palette; BSD/macOS `ls` uses `LSCOLORS` with ANSI approximations from the terminal
palette. Cute directories use rose; Gentleman directories use cyan. This does
not replace `ls`, change its options, or force colors when disabled. Apply the
Home Manager configuration before reloading the installed profile; repository
edits alone do not update an existing shell.

## Multiplexers

| Tool | Description |
|------|-------------|
| **Tmux** | Battle-tested, widely used, lots of plugins |
| **Zellij** | Modern, WebAssembly plugins, floating panes |
| **Herdr** | Agent-focused terminal multiplexer for AI coding workflows |

## Editor

| Tool | Description |
|------|-------------|
| **Neovim** | LazyVim config with LSP, completions, AI integration |

## Prompts

| Tool | Description |
|------|-------------|
| **Starship** | Cross-shell prompt with Git integration |
