# Gentleman.Dots

<img width="2998" height="1649" alt="image" src="https://github.com/user-attachments/assets/c5a1eae2-69de-4ca1-8b4d-9d9b56e4cb5a" />

## What is this?

A complete development environment configuration including:

- **Neovim** with LSP, autocompletion, and AI assistants (Claude, Gemini, OpenCode)
- **Shells**: Fish, Zsh, Nushell
- **Terminal Multiplexers**: Tmux, Zellij
- **Terminal Emulators**: Alacritty, WezTerm, Kitty, Ghostty

## Quick Start

### Option 1: Homebrew (Recommended)

```bash
brew install Gentleman-Programming/tap/gentleman-dots
gentleman.dots
```

### Option 2: Direct Download

```bash
# macOS Apple Silicon
curl -fsSL https://github.com/Gentleman-Programming/Gentleman.Dots/releases/latest/download/gentleman.dots-darwin-arm64 -o gentleman.dots

# macOS Intel
curl -fsSL https://github.com/Gentleman-Programming/Gentleman.Dots/releases/latest/download/gentleman.dots-darwin-amd64 -o gentleman.dots

# Linux x86_64
curl -fsSL https://github.com/Gentleman-Programming/Gentleman.Dots/releases/latest/download/gentleman.dots-linux-amd64 -o gentleman.dots

# Then run
chmod +x gentleman.dots
./gentleman.dots
```

The TUI guides you through selecting your preferred tools and handles all the configuration automatically.

> **Windows users:** You must set up WSL first. See the [Manual Installation Guide](docs/manual-installation.md#windows-wsl).

---

## Documentation

### Installation

| Document | Description |
|----------|-------------|
| [**TUI Installer Guide**](docs/tui-installer.md) | Interactive installer features, navigation, backup/restore |
| [**Manual Installation**](docs/manual-installation.md) | Step-by-step manual setup for all platforms |

### Configuration

| Document | Description |
|----------|-------------|
| [**Neovim Keymaps**](docs/neovim-keymaps.md) | Complete reference of all keybindings |
| [**AI Configuration**](docs/ai-configuration.md) | OpenCode, Claude, Copilot, and other AI assistants |

---

## Tools Overview

### Terminal Emulators

| Tool | Description |
|------|-------------|
| **Ghostty** | GPU-accelerated, native, blazing fast |
| **Kitty** | Feature-rich, GPU-based rendering |
| **WezTerm** | Lua-configurable, cross-platform |
| **Alacritty** | Minimal, Rust-based, lightweight |

### Shells

| Tool | Description |
|------|-------------|
| **Nushell** | Structured data, modern syntax, pipelines |
| **Fish** | User-friendly, great defaults, no config needed |
| **Zsh** | Highly customizable, POSIX-compatible, Powerlevel10k |

### Multiplexers

| Tool | Description |
|------|-------------|
| **Tmux** | Battle-tested, widely used, lots of plugins |
| **Zellij** | Modern, WebAssembly plugins, floating panes |

### Editor

| Tool | Description |
|------|-------------|
| **Neovim** | LazyVim config with LSP, completions, AI |

---

## Bleeding Edge

Want the latest experimental features from my daily workflow (macOS only)?

Check out the [`nix-migration` branch](https://github.com/Gentleman-Programming/Gentleman.Dots/tree/nix-migration).

This branch contains cutting-edge configurations that eventually make their way to `main` once stable.

---

## Project Structure

```
Gentleman.Dots/
├── docs/                    # Documentation
│   ├── tui-installer.md     # TUI guide
│   ├── manual-installation.md
│   └── ai-configuration.md
├── installer/               # Go TUI installer
│   ├── cmd/                 # Entry point
│   └── internal/            # TUI and system packages
├── GentlemanNvim/           # Neovim configuration
├── GentlemanFish/           # Fish shell config
├── GentlemanZsh/            # Zsh + Oh-My-Zsh + P10k
├── GentlemanNushell/        # Nushell config
├── GentlemanTmux/           # Tmux config
├── GentlemanZellij/         # Zellij config
├── GentlemanGhostty/        # Ghostty terminal config
├── GentlemanKitty/          # Kitty terminal config
├── GentlemanOpenCode/       # OpenCode AI config
├── alacritty.toml           # Alacritty config
├── .wezterm.lua             # WezTerm config
└── starship.toml            # Starship prompt config
```

---

## Support

- **Issues**: [GitHub Issues](https://github.com/Gentleman-Programming/Gentleman.Dots/issues)
- **Discord**: [Gentleman Programming Community](https://discord.gg/gentleman-programming)
- **YouTube**: [@GentlemanProgramming](https://youtube.com/@GentlemanProgramming)
- **Twitch**: [GentlemanProgramming](https://twitch.tv/GentlemanProgramming)

---

## License

MIT License - feel free to use, modify, and share.

**Happy coding!** 🎩
