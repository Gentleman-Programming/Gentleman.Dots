# Gentleman.Dots

📄 Read this in: **English** | [Español](README.es.md)

Project description in English...

> ℹ️ **Update (January 2026)**: OpenCode now supports Claude Max/Pro subscriptions via the `opencode-anthropic-auth` plugin (included in this config). Both **Claude Code** and **OpenCode** work with your Claude subscription. *Note: This workaround is stable for now, but Anthropic could block it in the future.*

## Table of Contents

- [What is this?](#what-is-this)
- [Quick Start](#quick-start)
- [Supported Platforms](#supported-platforms)
- [Vim Mastery Trainer](#-vim-mastery-trainer)
- [Documentation](#documentation)
- [Tools Overview](#tools-overview)
- [Bleeding Edge](#bleeding-edge)
- [Project Structure](#project-structure)
- [Support](#support)

---

## Preview

### TUI Installer

<img width="1424" height="1536" alt="TUI Installer" src="https://github.com/user-attachments/assets/1db56d3b-a8c0-4885-82aa-c5ec04af4ac0" />

### Showcase

<img width="3840" height="2160" alt="Development Environment Showcase" src="https://github.com/user-attachments/assets/fff14c05-9676-4e04-b05e-dab5e3cf300a" />

---

## What is this?

A complete development environment configuration including:

- **Neovim** with LSP, autocompletion, and AI assistants (Claude Code, Gemini, OpenCode)
- **Shells**: Fish, Zsh, Nushell
- **Terminal Multiplexers**: Tmux, Zellij
- **Terminal Emulators**: Alacritty, WezTerm, Kitty, Ghostty

## Quick Start

### Option 1: Homebrew (Recommended)

```bash
brew install Gentleman-Programming/tap/gentleman-dots
gentleman-dots
```

### Option 2: Direct Download

```bash
# macOS Apple Silicon
curl -fsSL https://github.com/Gentleman-Programming/Gentleman.Dots/releases/latest/download/gentleman-installer-darwin-arm64 -o gentleman.dots

# macOS Intel
curl -fsSL https://github.com/Gentleman-Programming/Gentleman.Dots/releases/latest/download/gentleman-installer-darwin-amd64 -o gentleman.dots

# Linux x86_64
curl -fsSL https://github.com/Gentleman-Programming/Gentleman.Dots/releases/latest/download/gentleman-installer-linux-amd64 -o gentleman.dots

# Linux ARM64 (Raspberry Pi, etc.)
curl -fsSL https://github.com/Gentleman-Programming/Gentleman.Dots/releases/latest/download/gentleman-installer-linux-arm64 -o gentleman.dots

# Then run
chmod +x gentleman.dots
./gentleman.dots
```

### Option 3: Termux (Android)

Termux requires building the installer locally (Go cross-compilation to Android has limitations).

```bash
# 1. Install dependencies
pkg update && pkg upgrade
pkg install git golang

# 2. Clone the repository
git clone https://github.com/Gentleman-Programming/Gentleman.Dots.git
cd Gentleman.Dots/installer

# 3. Build and run
go build -o ~/gentleman-installer ./cmd/gentleman-installer
cd ~
./gentleman-installer
```

| Termux Support | Status |
|----------------|--------|
| Shells (Fish, Zsh, Nushell) | ✅ Available |
| Multiplexers (Tmux, Zellij) | ✅ Available |
| Neovim with full config | ✅ Available |
| Nerd Fonts | ✅ Auto-installed to `~/.termux/font.ttf` |
| Terminal emulators | ❌ Not applicable |
| Homebrew | ❌ Uses `pkg` instead |

> **Tip:** After installation, restart Termux to apply the font, then run `tmux` or `zellij` to start your configured environment.

The TUI guides you through selecting your preferred tools and handles all the configuration automatically.

> **Windows users:** You must set up WSL first. See the [Manual Installation Guide](docs/manual-installation.md#windows-wsl).

---

## Supported Platforms

| Platform | Architecture | Install Method | Package Manager |
|----------|--------------|----------------|-----------------|
| macOS | Apple Silicon (ARM64) | Homebrew, Direct Download | Homebrew |
| macOS | Intel (x86_64) | Homebrew, Direct Download | Homebrew |
| Linux (Ubuntu/Debian) | x86_64, ARM64 | Homebrew, Direct Download | Homebrew |
| Linux (Fedora/RHEL) | x86_64, ARM64 | Direct Download | dnf |
| Linux (Arch) | x86_64 | Homebrew, Direct Download | Homebrew |
| Windows | WSL | Direct Download (see docs) | Homebrew |
| Android | Termux (ARM64) | Build locally (see above) | pkg |

---

## 🎮 Vim Mastery Trainer

Learn Vim the fun way! The installer includes an interactive RPG-style trainer with:

| Module | Keys Covered |
|--------|--------------|
| 🔤 Horizontal Movement | `w`, `e`, `b`, `f`, `t`, `0`, `$`, `^` |
| ↕️ Vertical Movement | `j`, `k`, `G`, `gg`, `{`, `}` |
| 📦 Text Objects | `iw`, `aw`, `i"`, `a(`, `it`, `at` |
| ✂️ Change & Repeat | `d`, `c`, `dd`, `cc`, `D`, `C`, `x` |
| 🔄 Substitution | `r`, `R`, `s`, `S`, `~`, `gu`, `gU`, `J` |
| 🎬 Macros & Registers | `qa`, `@a`, `@@`, `"ay`, `"+p` |
| 🔍 Regex/Search | `/`, `?`, `n`, `N`, `*`, `#`, `\v` |

Each module includes 15 progressive lessons, practice mode with intelligent exercise selection, boss fights, and XP tracking.

Launch it from the main menu: **Vim Mastery Trainer**

---

## Documentation

| Document | Description |
|----------|-------------|
| [TUI Installer Guide](docs/tui-installer.md) | Interactive installer features, navigation, backup/restore |
| [Manual Installation](docs/manual-installation.md) | Step-by-step manual setup for all platforms |
| [Neovim Keymaps](docs/neovim-keymaps.md) | Complete reference of all keybindings |
| [AI Configuration](docs/ai-configuration.md) | Claude Code, OpenCode, Copilot, and other AI assistants |
| [Vim Trainer Spec](docs/vim-trainer-spec.md) | Technical specification for the Vim Mastery Trainer |
| [Docker Testing](docs/docker-testing.md) | E2E testing with Docker containers |
| [Contributing](docs/contributing.md) | Development setup, skills system, E2E tests, release process |

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

### Prompts

| Tool | Description |
|------|-------------|
| **Starship** | Cross-shell prompt with Git integration |

---

## Bleeding Edge

Want the latest experimental features from my daily workflow (macOS only)?

Check out the [`nix-migration` branch](https://github.com/Gentleman-Programming/Gentleman.Dots/tree/nix-migration).

This branch contains cutting-edge configurations that eventually make their way to `main` once stable.

---

## Project Structure

```
Gentleman.Dots/
├── installer/               # Go TUI installer source
│   ├── cmd/                 # Entry point
│   ├── internal/            # TUI, system, and trainer packages
│   └── e2e/                 # Docker-based E2E tests
├── docs/                    # Documentation
├── skills/                  # AI agent skills (repo-specific)
│
├── GentlemanNvim/           # Neovim configuration (LazyVim)
├── GentlemanClaude/         # Claude Code config + user skills
│   └── skills/              # Installable skills (React, Next.js, etc.)
├── GentlemanOpenCode/       # OpenCode AI config
│
├── GentlemanFish/           # Fish shell config
├── GentlemanZsh/            # Zsh + Oh-My-Zsh + Powerlevel10k
├── GentlemanNushell/        # Nushell config
├── GentlemanTmux/           # Tmux config
├── GentlemanZellij/         # Zellij config
│
├── GentlemanGhostty/        # Ghostty terminal config
├── GentlemanKitty/          # Kitty terminal config
├── alacritty.toml           # Alacritty config
├── .wezterm.lua             # WezTerm config
│
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
