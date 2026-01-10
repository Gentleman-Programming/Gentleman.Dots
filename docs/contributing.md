# Contributing to Gentleman.Dots

Guide for contributors and developers working on Gentleman.Dots.

## Table of Contents

- [Development Setup](#development-setup)
- [Project Structure](#project-structure)
- [AI Skills System](#ai-skills-system)
- [E2E Testing](#e2e-testing)
- [Release Process](#release-process)

## Development Setup

### Requirements

| Tool | Version | Purpose |
|------|---------|---------|
| Go | 1.21+ | Build the installer |
| Docker | Latest | Run E2E tests |
| Git | Latest | Version control |

### Build from Source

```bash
git clone https://github.com/Gentleman-Programming/Gentleman.Dots.git
cd Gentleman.Dots/installer
go build -o gentleman-dots ./cmd/gentleman-installer
./gentleman-dots
```

### Run Tests

```bash
cd installer
go test ./... -v
```

## Project Structure

```
Gentleman.Dots/
├── installer/                    # Go TUI installer
│   ├── cmd/gentleman-installer/  # Entry point
│   ├── internal/
│   │   ├── system/               # OS detection, command execution
│   │   └── tui/                  # Bubbletea screens, views, installer
│   └── e2e/                      # Docker-based E2E tests
├── skills/                       # Repository-specific AI skills
│   ├── setup.sh                  # Sync script for AI assistants
│   └── */SKILL.md                # Individual skills
├── GentlemanClaude/              # Claude Code config + user skills
├── GentlemanNvim/                # Neovim configuration
├── GentlemanFish/                # Fish shell config
├── GentlemanZsh/                 # Zsh config
├── GentlemanNushell/             # Nushell config
├── GentlemanTmux/                # Tmux config
├── GentlemanZellij/              # Zellij config
├── docs/                         # Documentation
└── AGENTS.md                     # Single source of truth for AI skills
```

## AI Skills System

The repository uses a skills system to provide context to AI assistants (Claude, Gemini, Copilot, etc.).

### Single Source of Truth

`AGENTS.md` is the master file. All other AI instruction files are generated from it.

### Sync Skills

```bash
# Interactive menu
./skills/setup.sh

# Generate all formats (CLAUDE.md, GEMINI.md, etc.)
./skills/setup.sh --all

# Sync to user config directories
./skills/setup.sh --sync-all

# Individual targets
./skills/setup.sh --claude      # CLAUDE.md
./skills/setup.sh --gemini      # GEMINI.md
./skills/setup.sh --copilot     # .github/copilot-instructions.md
./skills/setup.sh --codex       # CODEX.md
./skills/setup.sh --sync-claude # ~/.claude/skills/
./skills/setup.sh --sync-opencode # ~/.config/opencode/skill/
```

### Skill Types

| Type | Location | Purpose |
|------|----------|---------|
| Repository skills | `skills/` | For this codebase (bubbletea, trainer, etc.) |
| User skills | `GentlemanClaude/skills/` | Installed to user's ~/.claude/skills/ |

### Creating a New Skill

1. Load the skill-creator skill: `mcp_skill("skill-creator")`
2. Create directory under appropriate location
3. Add `SKILL.md` following the template
4. Register in `AGENTS.md`
5. Run `./skills/setup.sh --all`

## E2E Testing

Docker-based tests verify the installer works across different environments.

### Quick Start

```bash
cd installer/e2e

# Interactive TUI menu
./docker-test.sh

# Run all E2E tests
./docker-test.sh e2e

# Run specific environment
./docker-test.sh e2e ubuntu
./docker-test.sh e2e debian
./docker-test.sh e2e fedora
./docker-test.sh e2e alpine
./docker-test.sh e2e termux

# Interactive shell for debugging
./docker-test.sh shell ubuntu
```

### Test Environments

| Environment | Shell | Package Manager | Tests |
|-------------|-------|-----------------|-------|
| Ubuntu | bash | apt + Homebrew | Full E2E + backup |
| Debian | dash | apt + Homebrew | Basic + shell detection |
| Fedora | bash | dnf | Full E2E |
| Alpine | ash | apk | POSIX compatibility |
| Termux | sh | pkg (simulated) | Android/Termux specific |

### Adding Tests

Edit `e2e_test.sh` or `e2e_test_termux.sh`:

```bash
test_my_feature() {
    log_test "My feature works"
    if [ -f "$HOME/.config/myfile" ]; then
        log_pass "Config file exists"
    else
        log_fail "Config file not found"
    fi
}
```

Tests must be POSIX-compliant (no bashisms).

## Release Process

### 1. Build Binaries

```bash
cd installer
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o gentleman-installer-darwin-amd64 ./cmd/gentleman-installer
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o gentleman-installer-darwin-arm64 ./cmd/gentleman-installer
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o gentleman-installer-linux-amd64 ./cmd/gentleman-installer
GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o gentleman-installer-linux-arm64 ./cmd/gentleman-installer
```

### 2. Create Tag and Release

```bash
git tag v{VERSION}
git push origin v{VERSION}

gh release create v{VERSION} \
  installer/gentleman-installer-darwin-amd64 \
  installer/gentleman-installer-darwin-arm64 \
  installer/gentleman-installer-linux-amd64 \
  installer/gentleman-installer-linux-arm64 \
  --title "v{VERSION}" \
  --notes "## Changes
- Feature/fix description"
```

### 3. Update Homebrew Formula

```bash
# Get SHA256 for each binary
shasum -a 256 installer/gentleman-installer-*

# Update homebrew-tap/Formula/gentleman-dots.rb with new version and hashes
# Commit to both repos
```

### Version Guidelines

| Change Type | Version Bump | Example |
|-------------|--------------|---------|
| New platform/major feature | Minor (x.Y.0) | v2.7.0 |
| Bug fixes, improvements | Patch (x.y.Z) | v2.6.2 |
| Breaking changes | Major (X.0.0) | v3.0.0 |

## Code Style

### Go

- Follow standard Go conventions
- Use `gofmt` for formatting
- Table-driven tests preferred
- Error wrapping with context

### Shell Scripts

- POSIX-compliant (no bashisms in tests)
- Use `shellcheck` for linting
- Quote all variables

### Documentation

- Use tables for structured data
- Include Table of Contents for long docs
- Code examples for every feature
