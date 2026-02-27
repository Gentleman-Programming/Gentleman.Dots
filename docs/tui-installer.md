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
7. **AI Tools**: Multi-select Claude Code, OpenCode, Gemini CLI, GitHub Copilot
8. **AI Framework**: Choose preset or custom module selection (203 modules across 6 categories)
9. **Backup Confirmation**: Option to backup existing configs before overwriting
10. **Installation**: Watch real-time progress

> See [AI Tools & Framework Integration](ai-tools-integration.md) for detailed documentation on steps 7-8, including the category drill-down UI, viewport scrolling, preset reference, and SDD choice.

### Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `↑` / `k` | Move up |
| `↓` / `j` | Move down |
| `Enter` / `Space` | Select option |
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
| `--ai-tools` | `claude,opencode,gemini,copilot` | AI tools (comma-separated) |
| `--ai-framework` | | Install AI coding framework |
| `--ai-preset` | `minimal,frontend,backend,fullstack,data,complete` | Framework preset |
| `--ai-modules` | `hooks,commands,skills,agents,sdd,mcp` | Framework features (comma-separated) |
| `--agent-teams-lite` | | Install Agent Teams Lite SDD framework |

### Examples

```bash
# Interactive TUI (default)
gentleman.dots

# Non-interactive with Fish + Zellij + Neovim
gentleman.dots --non-interactive --shell=fish --wm=zellij --nvim

# Test mode with Zsh + Tmux (no terminal, no nvim)
gentleman.dots --test --non-interactive --shell=zsh --wm=tmux

# Dry run to preview changes
gentleman.dots --dry-run

# Full setup with AI tools and framework
gentleman.dots --non-interactive --shell=fish --nvim \
  --ai-tools=claude,opencode --ai-preset=fullstack

# Custom framework features with Agent Teams Lite
gentleman.dots --non-interactive --shell=zsh --ai-tools=claude --ai-framework \
  --ai-modules=hooks,skills --agent-teams-lite

# Verbose output (shows all command logs)
GENTLEMAN_VERBOSE=1 gentleman.dots --non-interactive --shell=fish --nvim
```

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
- AI assistants (OpenCode, Claude, Copilot, etc.)

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

### Copy/Paste Not Working in Zellij (Linux)

If you can't copy text from the terminal when using Zellij on Linux:

1. Edit `~/.config/zellij/config.kdl`
2. Uncomment the appropriate line for your system:
   - **X11**: `copy_command "xclip -selection clipboard"` AND `copy_clipboard "primary"`
   - **Wayland**: `copy_command "wl-copy"`

See: [Zellij FAQ](https://zellij.dev/documentation/faq.html#copy--paste-isnt-working-how-can-i-fix-this)

### Fish/Zsh Fails in WSL with "missing or unsuitable terminal"

When using WezTerm on Windows with WSL, Fish or Zsh may fail with:
```
missing or unsuitable terminal: wezterm
```

**Status**: ✅ Fixed automatically in the default `.wezterm.lua` config (v2.7.7+).

If you're using an older config, update your `.wezterm.lua` to include the auto-detection:
```lua
if wezterm.target_triple:find("windows") then
  config.term = "xterm-256color"
else
  config.term = "wezterm"
end
```

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
