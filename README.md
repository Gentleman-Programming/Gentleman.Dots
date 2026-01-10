# Gentleman.Dots

> ℹ️ **Update (January 2026)**: OpenCode now supports Claude Max/Pro subscriptions via the `opencode-anthropic-auth` plugin (included in this config). Both **Claude Code** and **OpenCode** work with your Claude subscription. *Note: This workaround is stable for now, but Anthropic could block it in the future.*

## TUI
<img width="1424" height="1536" alt="image" src="https://github.com/user-attachments/assets/1db56d3b-a8c0-4885-82aa-c5ec04af4ac0" />

## ShowCase
<img width="3840" height="2160" alt="image" src="https://github.com/user-attachments/assets/fff14c05-9676-4e04-b05e-dab5e3cf300a" />

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
gentleman.dots
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

```bash
# Install dependencies first
pkg install git curl

# Download the installer
curl -fsSL https://github.com/Gentleman-Programming/Gentleman.Dots/releases/latest/download/gentleman-installer-android-arm64 -o gentleman.dots

# Run it
chmod +x gentleman.dots
./gentleman.dots
```

> **Note:** Termux support includes Fish, Zsh, Nushell, Neovim, Tmux, and Zellij. Terminal emulators are not available on Android.

The TUI guides you through selecting your preferred tools and handles all the configuration automatically.

> **Windows users:** You must set up WSL first. See the [Manual Installation Guide](docs/manual-installation.md#windows-wsl).

---

## Supported Platforms

| Platform | Architecture | Install Method |
|----------|--------------|----------------|
| macOS | Apple Silicon (ARM64) | Homebrew, Direct Download |
| macOS | Intel (x86_64) | Homebrew, Direct Download |
| Linux | x86_64 | Homebrew, Direct Download |
| Linux | ARM64 | Homebrew, Direct Download |
| Windows | WSL | Direct Download (see docs) |
| Android | Termux (ARM64) | Direct Download |

---

## 🎮 Vim Mastery Trainer

Learn Vim the fun way! The installer includes an interactive RPG-style trainer with:

**7 Complete Modules:**
- 🔤 **Horizontal Movement** - `w`, `e`, `b`, `f`, `t`, `0`, `$`, `^`
- ↕️ **Vertical Movement** - `j`, `k`, `G`, `gg`, `{`, `}`
- 📦 **Text Objects** - `iw`, `aw`, `i"`, `a(`, `it`, `at`
- ✂️ **Change & Repeat** - `d`, `c`, `dd`, `cc`, `D`, `C`, `x`
- 🔄 **Substitution** - `r`, `R`, `s`, `S`, `~`, `gu`, `gU`, `J`
- 🎬 **Macros & Registers** - `qa`, `@a`, `@@`, `"ay`, `"+p`
- 🔍 **Regex/Search** - `/`, `?`, `n`, `N`, `*`, `#`, `\v`

**Each module includes:**
- 15 progressive lessons with explanations
- Practice mode with intelligent exercise selection
- Boss fight to test your skills
- XP and progress tracking

Launch it from the main menu: **Vim Mastery Trainer**

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
| [**AI Configuration**](docs/ai-configuration.md) | Claude Code, OpenCode, Copilot, and other AI assistants |

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
├── GentlemanClaude/         # Claude Code AI config (primary)
│   ├── tweakcc-theme.json   # Visual theme for Claude Code (via tweakcc)
│   └── ...
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
