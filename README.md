# Gentleman.Dots

> 🤖 **NEW**: The AI development layer now lives in its own installer — [**AI Gentle Stack (gentle-ai)**](https://github.com/Gentleman-Programming/gentle-ai). It configures Claude Code, OpenCode, Gemini CLI, Cursor, and VS Code Copilot with persistent memory, SDD workflow, skills, and the Gentleman persona. Install Gentleman.Dots first, then run `gentle-ai` for the AI layer.

📄 Read this in: **English** | [Español](README.es.md)

## Table of Contents

- [What is this?](#what-is-this)
- [Quick Start](#quick-start)
- [Supported Platforms](#supported-platforms)
- [AI Development Layer](#-ai-development-layer)
- [Vim Mastery Trainer](#-vim-mastery-trainer)
- [Documentation](#documentation)
- [Tools Overview](#tools-overview)
- [Bleeding Edge](#bleeding-edge)
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

- **Neovim** with LSP, autocompletion, and AI integration
- **Shells**: Fish, Zsh, Nushell
- **Terminal Multiplexers**: Tmux, Zellij
- **Terminal Emulators**: Alacritty, WezTerm, Kitty, Ghostty
- **AI Agent Configs**: Source configurations for Claude Code and OpenCode (installed via [gentle-ai](https://github.com/Gentleman-Programming/gentle-ai))

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

Termux requires building locally. See the [Termux Installation Guide](docs/manual-installation.md#termux) for full instructions.

The TUI guides you through selecting your preferred tools and handles all the configuration automatically.

> **Tmux users:** After installation, open tmux and press `prefix + I` (capital I) to install plugins via TPM. This ensures the theme and all plugins load correctly.

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

## 🤖 AI Development Layer

Gentleman.Dots handles your **dev environment** (editor, shells, terminals). For the **AI development layer** (agents, memory, skills, workflow), use the companion project:

### [AI Gentle Stack (gentle-ai)](https://github.com/Gentleman-Programming/gentle-ai)

```bash
brew install Gentleman-Programming/tap/gentle-ai
gentle-ai
```

It configures your AI coding agents with everything they need:

| Component | What it does |
|-----------|-------------|
| **Engram** | Persistent memory across sessions (MCP server) |
| **SDD Workflow** | Spec-Driven Development with orchestrated sub-agents |
| **Skills** | 24 coding pattern libraries (React 19, Next.js 15, TypeScript, Tailwind 4, etc.) |
| **Context7** | Up-to-date library documentation via MCP |
| **Persona** | Gentleman teaching style for AI responses |
| **Permissions** | Security-first defaults for all agents |

### Supported Agents

| Agent | Single Agent | Multi Agent |
|-------|:----------:|:-----------:|
| **Claude Code** | ✅ | ✅ |
| **OpenCode** | ✅ | ✅ |
| **Gemini CLI** | ✅ | ✅ |
| **Cursor** | ✅ | — |
| **VS Code Copilot** | ✅ | — |

> **Single agent**: One orchestrator handles all SDD phases.
> **Multi agent**: Dedicated sub-agent per phase with individual model routing (e.g., Claude Opus for design, Gemini for specs, GPT for verification).

### What lives where

| | This repo (Gentleman.Dots) | gentle-ai |
|--|---------------------------|-----------|
| **Purpose** | Dev environment (editors, shells, terminals) | AI development layer (agents, memory, skills) |
| **Installs** | Neovim, Fish/Zsh, Tmux/Zellij, Ghostty | Configures Claude Code, OpenCode, Gemini CLI, Cursor, VS Code Copilot |
| **Source configs** | `GentlemanClaude/`, `GentlemanOpenCode/` | Reads from this repo + its own assets |

Install Gentleman.Dots first for your dev environment, then `gentle-ai` for the AI layer on top.

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
| [AI Gentle Stack](https://github.com/Gentleman-Programming/gentle-ai) | AI layer installer — Engram, SDD, Skills, Persona (separate repo) |
| [Vim Trainer Spec](docs/vim-trainer-spec.md) | Technical specification for the Vim Mastery Trainer |
| [Docker Testing](docs/docker-testing.md) | E2E testing with Docker containers |
| [Contributing](docs/contributing.md) | Development setup, skills system, E2E tests, release process |

---

## Tools Overview

- **Terminal Emulators**: Ghostty, Kitty, WezTerm, Alacritty
- **Shells**: Nushell, Fish, Zsh (+ Powerlevel10k)
- **Multiplexers**: Tmux, Zellij
- **Editor**: Neovim (LazyVim with LSP, completions, AI)
- **Prompt**: Starship

> See [Tools Reference](docs/tools.md) for detailed descriptions of each tool.

---

## Bleeding Edge

Want the latest experimental features from my daily workflow (macOS only)?

Check out the [`nix-migration` branch](https://github.com/Gentleman-Programming/Gentleman.Dots/tree/nix-migration).

This branch contains cutting-edge configurations that eventually make their way to `main` once stable.

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
