# TUI Installer

The Gentleman.Dots TUI Installer is a modern, interactive terminal application built with Go and [Bubbletea](https://github.com/charmbracelet/bubbletea) that guides you through the complete setup of your development environment.

## Features

- **Interactive Navigation**: Use arrow keys or Vim-style `j/k` bindings
- **Smart Detection**: Automatically detects your OS, existing configs, and installed tools
- **Backup & Restore**: Safely backup existing configurations before installation
- **Educational Content**: Learn about each tool before choosing (terminals, shells, multiplexers, Neovim)
- **Neovim Keymaps Reference**: Built-in keymap browser organized by category
- **LazyVim Guide**: Comprehensive guide to LazyVim concepts and usage
- **Progress Tracking**: Real-time installation progress with detailed logs

## Quick Start

### Option 1: Homebrew (Recommended)

```bash
brew install Gentleman-Programming/tap/gentleman-dots
gentleman.dots
```

### Option 2: Download Pre-built Binary

```bash
# macOS Apple Silicon (M1/M2/M3)
curl -sL https://github.com/Gentleman-Programming/Gentleman.Dots/releases/latest/download/gentleman-dots-darwin-arm64.tar.gz | tar xz
./gentleman-dots

# macOS Intel
curl -sL https://github.com/Gentleman-Programming/Gentleman.Dots/releases/latest/download/gentleman-dots-darwin-amd64.tar.gz | tar xz
./gentleman-dots

# Linux x86_64
curl -sL https://github.com/Gentleman-Programming/Gentleman.Dots/releases/latest/download/gentleman-dots-linux-amd64.tar.gz | tar xz
./gentleman-dots

# Linux ARM64
curl -sL https://github.com/Gentleman-Programming/Gentleman.Dots/releases/latest/download/gentleman-dots-linux-arm64.tar.gz | tar xz
./gentleman-dots
```

### Option 3: Build from Source

```bash
git clone https://github.com/Gentleman-Programming/Gentleman.Dots.git
cd Gentleman.Dots/installer
go build -o gentleman-dots ./cmd/gentleman-installer
./gentleman-dots
```

## Screens & Navigation

### Main Menu

From the main menu you can:

- **Start Installation**: Begin the guided setup process
- **Learn About Tools**: Explore terminals, shells, and multiplexers
- **Neovim Keymaps**: Browse all configured keybindings
- **LazyVim Guide**: Learn LazyVim fundamentals
- **Restore from Backup**: Restore previous configurations (if backups exist)
- **Exit**: Quit the installer

### Installation Flow

1. **OS Selection**: Choose macOS or Linux
2. **Terminal Emulator**: Select Ghostty, Kitty, WezTerm, Alacritty, or None
3. **Font Installation**: Iosevka Term Nerd Font (required for icons)
4. **Shell**: Choose Nushell, Fish, Zsh, or None
5. **Window Manager**: Select Tmux, Zellij, or None
6. **Neovim**: Configure LazyVim with LSP and AI assistants
7. **Backup Confirmation**: Option to backup existing configs before overwriting
8. **Installation**: Watch real-time progress

### Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `↑/k` | Move up |
| `↓/j` | Move down |
| `Enter/Space` | Select option |
| `Esc` | Go back |
| `q` | Quit (when not installing) |
| `d` | Toggle details (during installation) |
| `Ctrl+C` | Force quit |

## Backup & Restore

### Automatic Backup Detection

The installer automatically detects existing configurations for:
- Neovim (`~/.config/nvim`)
- Fish (`~/.config/fish`)
- Zsh (`~/.zshrc`, `~/.oh-my-zsh`)
- Nushell (`~/.config/nushell` or `~/Library/Application Support/nushell`)
- Tmux (`~/.tmux.conf`, `~/.tmux`)
- Zellij (`~/.config/zellij`)
- Alacritty (`~/.config/alacritty`)
- WezTerm (`~/.config/wezterm`, `~/.wezterm.lua`)
- Kitty (`~/.config/kitty`)
- Ghostty (`~/.config/ghostty`)
- Starship (`~/.config/starship.toml`)

### Backup Location

Backups are stored in your home directory:
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
- **Ghostty**: GPU-accelerated, native, fast
- **Kitty**: Feature-rich, GPU-based
- **WezTerm**: Lua-configurable, cross-platform
- **Alacritty**: Minimal, Rust-based

### Shells
- **Nushell**: Structured data, modern syntax
- **Fish**: User-friendly, great defaults
- **Zsh**: Highly customizable, POSIX-compatible

### Multiplexers
- **Tmux**: Battle-tested, widely used
- **Zellij**: Modern, WebAssembly plugins

### Neovim
- LazyVim configuration
- LSP setup
- AI assistants (OpenCode, Claude, Copilot, etc.)

## Command Line Flags

```bash
./gentleman-installer [flags]

Flags:
  --version     Show version information
  --help        Show help message
  --test        Run in test mode (no actual changes)
  --dry-run     Show what would be installed without making changes
```

## Requirements

- **macOS** 10.15+ or **Linux** (Ubuntu 20.04+, Arch, Debian)
- **Homebrew** (will be installed if missing)
- **Git** (for cloning the repository)
- **Internet connection** (for downloading packages)

## Troubleshooting

### Installation Fails

1. Check the detailed logs by pressing `d` during installation
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

- **Go 1.21+**
- **Bubbletea** - Terminal UI framework
- **Lipgloss** - Styling
- **Teatest** - Testing framework

### Running Tests

```bash
cd installer
go test ./... -v
```

### Updating Golden Files

```bash
go test ./internal/tui/... -update
```

### Project Structure

```
installer/
├── cmd/
│   └── gentleman-installer/
│       └── main.go           # Entry point
├── internal/
│   ├── system/
│   │   ├── detect.go         # OS/tool detection
│   │   └── exec.go           # Command execution, file ops, backups
│   └── tui/
│       ├── model.go          # App state, screens, choices
│       ├── update.go         # Event handlers
│       ├── view.go           # UI rendering
│       ├── installer.go      # Installation steps
│       ├── styles.go         # Gentleman theme colors
│       └── tools_info.go     # Tool descriptions, keymaps
└── go.mod
```
