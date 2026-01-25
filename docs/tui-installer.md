# TUI Installer

The Gentleman.Dots TUI Installer is a modern, interactive terminal application built with Go and [Bubbletea](https://github.com/charmbracelet/bubbletea) that guides you through the complete setup of your development environment.

## Table of Contents

- [Features](#features)
- [Quick Start](#quick-start)
- [Screens & Navigation](#screens--navigation)
- [Command Line Interface](#command-line-interface)
- [Backup & Restore](#backup--restore)
- [Learn Mode](#learn-mode)
- [Requirements](#requirements)
- [Troubleshooting](#troubleshooting)
- [Development](#development)

## Features

- **Interactive Navigation**: Arrow keys or Vim-style `j/k` bindings
- **Smart Detection**: Automatically detects your OS, existing configs, and installed tools
- **Backup & Restore**: Safely backup existing configurations before installation
- **Educational Content**: Learn about each tool before choosing (terminals, shells, multiplexers)
- **Neovim Keymaps Reference**: Built-in keymap browser organized by category
- **LazyVim Guide**: Comprehensive guide to LazyVim concepts and usage
- **Vim Trainer**: RPG-style interactive Vim learning with exercises and progression
- **Progress Tracking**: Real-time installation progress with detailed logs
- **Non-Interactive Mode**: CI/CD friendly installation via CLI flags

## Quick Start

### Option 1: Homebrew (Recommended)

```bash
brew install Gentleman-Programming/tap/gentleman-dots
gentleman-dots
```

### Option 2: Download Pre-built Binary

| Platform | Command |
|----------|---------|
| macOS Apple Silicon | `curl -sL https://github.com/Gentleman-Programming/Gentleman.Dots/releases/latest/download/gentleman-dots-darwin-arm64.tar.gz \| tar xz && ./gentleman-dots` |
| macOS Intel | `curl -sL https://github.com/Gentleman-Programming/Gentleman.Dots/releases/latest/download/gentleman-dots-darwin-amd64.tar.gz \| tar xz && ./gentleman-dots` |
| Linux x86_64 | `curl -sL https://github.com/Gentleman-Programming/Gentleman.Dots/releases/latest/download/gentleman-dots-linux-amd64.tar.gz \| tar xz && ./gentleman-dots` |
| Linux ARM64 | `curl -sL https://github.com/Gentleman-Programming/Gentleman.Dots/releases/latest/download/gentleman-dots-linux-arm64.tar.gz \| tar xz && ./gentleman-dots` |

### Option 3: Build from Source

```bash
git clone https://github.com/Gentleman-Programming/Gentleman.Dots.git
cd Gentleman.Dots/installer
go build -o gentleman-dots ./cmd/gentleman-installer
./gentleman-dots
```

## Screens & Navigation

### Main Menu

From the main menu you can access:

- **Start Installation**: Begin the guided setup process
- **Learn About Tools**: Explore terminals, shells, and multiplexers
- **Neovim Keymaps**: Browse all configured keybindings
- **LazyVim Guide**: Learn LazyVim fundamentals
- **Vim Trainer**: Practice Vim motions with interactive exercises
- **Restore from Backup**: Restore previous configurations (if backups exist)
- **Exit**: Quit the installer

### Installation Flow

1. **OS Selection**: Choose macOS, Linux, or Termux
2. **Terminal Emulator**: Select Ghostty, Kitty, WezTerm, Alacritty, or None
3. **Font Installation**: Iosevka Term Nerd Font (required for icons)
4. **Shell**: Choose Nushell, Fish, Zsh, or None
5. **Window Manager**: Select Tmux, Zellij, or None
6. **Neovim**: Configure LazyVim with LSP and AI assistants
7. **AI Coding Assistants**: Multi-select AI tools (OpenCode, Kilo Code, Continue.dev, Aider)
8. **Backup Confirmation**: Option to backup existing configs before overwriting
9. **Installation**: Watch real-time progress

**Skip Steps**: Each configuration screen (steps 2-7) includes a "⏭️  Skip this step" option. This is useful if you:
- Already have terminal/shell configurations you want to keep
- Only want to install specific components (e.g., just AI Assistants)
- Need to customize certain tools manually before running the installer

Skipped steps won't be installed or configured. For example, selecting "Skip this step" on the Terminal screen will keep your current terminal setup unchanged.

### Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `↑` / `k` | Move up |
| `↓` / `j` | Move down |
| `Enter` / `Space` | Select option |
| `Space` | Toggle checkbox (AI Assistants screen) |
| `Esc` | Go back |
| `q` | Quit (when not installing) |
| `d` | Toggle details (during installation) |
| `Ctrl+C` | Force quit |

## Command Line Interface

### Basic Flags

```bash
gentleman.dots [flags]
```

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--help` | `-h` | Show help message |
| `--version` | `-v` | Show version information |
| `--test` | `-t` | Run in test mode (uses temporary directory) |
| `--dry-run` | | Show what would be installed without doing it |
| `--non-interactive` | | Run without TUI, use CLI flags instead |

### Non-Interactive Mode

For CI/CD or scripted installations:

```bash
gentleman.dots --non-interactive --shell=<shell> [options]
```

| Flag | Values | Description |
|------|--------|-------------|
| `--shell` | `fish`, `zsh`, `nushell` | Shell to install (required) |
| `--terminal` | `alacritty`, `wezterm`, `kitty`, `ghostty`, `none` | Terminal emulator |
| `--wm` | `tmux`, `zellij`, `none` | Window manager |
| `--nvim` | | Install Neovim configuration |
| `--font` | | Install Nerd Font |
| `--backup` | `true`/`false` | Backup existing configs (default: true) |
| `--ai` | `opencode`, `kilocode`, `continue`, `aider` | AI assistants (comma-separated) |

### Examples

```bash
# Interactive TUI (default)
gentleman.dots

# Non-interactive with Fish + Zellij + Neovim
gentleman.dots --non-interactive --shell=fish --wm=zellij --nvim

# Non-interactive with OpenCode AI assistant
gentleman.dots --non-interactive --shell=zsh --nvim --ai=opencode

# Multiple AI assistants
gentleman.dots --non-interactive --shell=fish --ai=opencode,continue

# Test mode with Zsh + Tmux (no terminal, no nvim)
gentleman.dots --test --non-interactive --shell=zsh --wm=tmux

# Dry run to preview changes
gentleman.dots --dry-run

# Verbose output (shows all command logs)
GENTLEMAN_VERBOSE=1 gentleman.dots --non-interactive --shell=fish --nvim
```

## AI Coding Assistants

The installer supports multiple AI coding assistants that integrate with your development environment. You can select one or more assistants during installation.

### Available Assistants

| Assistant | Status | Description |
|-----------|--------|-------------|
| **OpenCode** | ✅ Available Now | Open-source AI coding assistant with context-aware completions |
| **Kilo Code** | 🚧 Coming Soon | Lightweight AI assistant optimized for performance |
| **Continue.dev** | 🚧 Coming Soon | Open-source autopilot for software development |
| **Aider** | 🚧 Coming Soon | AI pair programming in the terminal |

### How It Works

1. **Selection Screen**: Use `Space` to toggle checkboxes for each assistant
2. **Skills Installation**: Selected assistants will have their skills installed to `~/.{assistant}/skills/`
3. **Independent Installation**: AI assistants install independently from Neovim
4. **Configuration**: Each assistant uses its own configuration directory

### Interactive Mode

During installation, you'll see the AI Assistants screen:

```
Step 7: AI Coding Assistants
Select AI coding assistants (Space to toggle, Enter to confirm)

[ ] OpenCode
[ ] Kilo Code (Coming Soon)
[ ] Continue.dev (Coming Soon)
[ ] Aider (Coming Soon)
─────────────
✅ Confirm Selection
← Back
```

- **Navigation**: Use `↑`/`↓` or `j`/`k` to move
- **Toggle**: Press `Space` to select/deselect
- **Confirm**: Press `Enter` on "Confirm Selection"
- **Skip**: Leave all unchecked and confirm to skip

### Non-Interactive Mode

Use the `--ai` flag with comma-separated values:

```bash
# Single assistant
gentleman.dots --non-interactive --shell=fish --ai=opencode

# Multiple assistants
gentleman.dots --non-interactive --shell=fish --ai=opencode,continue

# No AI assistants (omit the flag)
gentleman.dots --non-interactive --shell=fish --nvim
```

### Skills Directory

Skills are installed to:

| Assistant | Skills Location |
|-----------|-----------------|
| OpenCode | `~/.opencode/skills/` |
| Kilo Code | `~/.kilocode/skills/` |
| Continue.dev | `~/.continue/skills/` |
| Aider | `~/.aider/skills/` |

### Manual Installation

If you want to install AI assistants later:

```bash
# OpenCode example
mkdir -p ~/.opencode/skills
cp -r ~/Gentleman.Dots/GentlemanClaude/skills/* ~/.opencode/skills/
```

See [Manual Installation](./manual-installation.md) for detailed steps.

## Backup & Restore

### Automatic Backup Detection

The installer automatically detects existing configurations for:

| Tool | Paths |
|------|-------|
| Neovim | `~/.config/nvim` |
| Fish | `~/.config/fish` |
| Zsh | `~/.zshrc`, `~/.oh-my-zsh` |
| Nushell | `~/.config/nushell`, `~/Library/Application Support/nushell` |
| Tmux | `~/.tmux.conf`, `~/.tmux` |
| Zellij | `~/.config/zellij` |
| Alacritty | `~/.config/alacritty` |
| WezTerm | `~/.config/wezterm`, `~/.wezterm.lua` |
| Kitty | `~/.config/kitty` |
| Ghostty | `~/.config/ghostty` |
| Starship | `~/.config/starship.toml` |
| OpenCode | `~/.opencode` |
| Kilo Code | `~/.kilocode` |
| Continue.dev | `~/.continue` |
| Aider | `~/.aider` |

### Backup Location

Backups are stored in your home directory with a timestamp:

```
~/.gentleman-backup-YYYYMMDD-HHMMSS/
```

### Restoring a Backup

1. Select "Restore from Backup" from the main menu
2. Choose the backup you want to restore
3. Confirm the restoration
4. Your previous configurations will be restored

## Learn Mode

The installer includes educational content to help you understand each tool:

### Terminals

| Terminal | Description |
|----------|-------------|
| Ghostty | GPU-accelerated, native, fast |
| Kitty | Feature-rich, GPU-based |
| WezTerm | Lua-configurable, cross-platform |
| Alacritty | Minimal, Rust-based |

### Shells

| Shell | Description |
|-------|-------------|
| Nushell | Structured data, modern syntax |
| Fish | User-friendly, great defaults |
| Zsh | Highly customizable, POSIX-compatible |

### Multiplexers

| Multiplexer | Description |
|-------------|-------------|
| Tmux | Battle-tested, widely used |
| Zellij | Modern, WebAssembly plugins |

### Neovim

- LazyVim configuration
- LSP setup
- Treesitter syntax highlighting
- Modern plugin ecosystem

### AI Assistants

| Assistant | Description |
|-----------|-------------|
| OpenCode | Context-aware code completions and generation |
| Kilo Code | Lightweight, performance-focused AI assistant |
| Continue.dev | Open-source autopilot for development |
| Aider | Terminal-based AI pair programming |

## Requirements

| Requirement | Details |
|-------------|---------|
| **macOS** | 10.15+ |
| **Linux** | Ubuntu 20.04+, Debian, Fedora/RHEL, Arch |
| **Termux** | Android terminal emulator |
| **Homebrew** | Will be installed if missing (macOS/Linux, except Fedora) |
| **Git** | For cloning the repository |
| **Internet** | For downloading packages |

## Troubleshooting

### Installation Fails

1. Press `d` during installation to view detailed logs
2. Ensure you have internet connectivity
3. Try running with `--test` flag first to verify detection
4. Check if Homebrew is properly installed: `brew --version`

### Backup Not Showing

Backups must be in your home directory with the format:

```
~/.gentleman-backup-*
```

### Font Not Displaying Correctly

1. Ensure the terminal is using "Iosevka Term Nerd Font"
2. Restart your terminal after font installation
3. On macOS, you may need to manually select the font in terminal preferences

## Development

The TUI installer is built with:

| Component | Description |
|-----------|-------------|
| **Go 1.25+** | Programming language |
| **Bubbletea** | Terminal UI framework |
| **Lipgloss** | Styling library |
| **Teatest** | Golden file testing |

### Running Tests

```bash
cd installer
go test ./... -v
```

### Updating Golden Files

```bash
cd installer
go test ./internal/tui/... -update
```

### Project Structure

```
installer/
├── cmd/
│   └── gentleman-installer/
│       └── main.go              # Entry point with CLI parsing
├── internal/
│   ├── system/
│   │   ├── detect.go            # OS/tool detection
│   │   └── exec.go              # Command execution, file ops, backups
│   └── tui/
│       ├── model.go             # App state, screens, choices
│       ├── update.go            # Event handlers
│       ├── view.go              # UI rendering
│       ├── installer.go         # Installation steps
│       ├── interactive.go       # TUI mode logic
│       ├── non_interactive.go   # CLI mode logic
│       ├── styles.go            # Gentleman theme colors
│       ├── tools_info.go        # Tool descriptions
│       ├── keymaps_*.go         # Keymap definitions
│       └── trainer/             # Vim Trainer RPG system
│           ├── types.go         # Exercise types, modules
│           ├── exercises.go     # Exercise definitions
│           ├── validation.go    # Input validation
│           ├── simulator.go     # Vim simulation
│           ├── stats.go         # Progress tracking
│           ├── gamestate.go     # Save/load game state
│           └── practice.go      # Practice mode
└── go.mod
```
