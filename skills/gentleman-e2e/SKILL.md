---
name: gentleman-e2e
description: >
  Docker-based E2E testing patterns for Gentleman.Dots installer.
  Trigger: When editing files in installer/e2e/, writing E2E tests, or adding platform support.
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

## When to Use

Use this skill when:
- Adding E2E tests for new features
- Creating Dockerfiles for new platforms
- Modifying the E2E test script
- Debugging installation failures
- Adding backup/restore test coverage

---

## Critical Patterns

### Pattern 1: Test Script Structure

All E2E tests in `e2e_test.sh` follow this pattern:

```bash
# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Logging functions
log_test() { printf "${YELLOW}[TEST]${NC} %s\n" "$1"; }
log_pass() { printf "${GREEN}[PASS]${NC} %s\n" "$1"; PASSED=$((PASSED + 1)); }
log_fail() { printf "${RED}[FAIL]${NC} %s\n" "$1"; FAILED=$((FAILED + 1)); }

# Test function pattern
test_feature_name() {
    log_test "Description of what we're testing"

    if some_condition; then
        log_pass "What passed"
    else
        log_fail "What failed"
    fi
}
```

### Pattern 2: Dockerfile Structure

Each platform has a Dockerfile in `e2e/`:

```dockerfile
FROM ubuntu:22.04

# Install base dependencies
RUN apt-get update && apt-get install -y \
    git curl sudo build-essential \
    && rm -rf /var/lib/apt/lists/*

# Create test user (non-root)
RUN useradd -m -s /bin/bash testuser && \
    echo "testuser ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers

# Copy installer binary
COPY gentleman-installer-linux-amd64 /usr/local/bin/gentleman-dots
RUN chmod +x /usr/local/bin/gentleman-dots

# Copy test script
COPY e2e/e2e_test.sh /home/testuser/e2e_test.sh
RUN chmod +x /home/testuser/e2e_test.sh

USER testuser
WORKDIR /home/testuser

# Run tests
CMD ["./e2e_test.sh"]
```

### Pattern 3: Non-Interactive Mode Testing

E2E tests use `--non-interactive` flag:

```bash
test_fish_tmux_nvim() {
    log_test "Install: Fish + Tmux + Nvim"

    if GENTLEMAN_VERBOSE=1 gentleman-dots --non-interactive \
        --shell=fish --wm=tmux --nvim --backup=false 2>&1; then

        # Verify fish config exists
        if [ -f "$HOME/.config/fish/config.fish" ]; then
            log_pass "Fish config was created"
        else
            log_fail "Fish config not found"
            return
        fi

        # Verify nvim config exists
        if [ -d "$HOME/.config/nvim" ]; then
            log_pass "Neovim config directory created"
        else
            log_fail "Neovim config directory not found"
        fi
    else
        log_fail "Installation failed"
    fi
}
```

### Pattern 4: Cleanup Between Tests

Always cleanup before each test:

```bash
test_something() {
    # Clean previous test artifacts
    rm -rf "$HOME/.config" "$HOME/.zshrc" 2>/dev/null || true
    mkdir -p "$HOME/.config"

    # Run test...
}
```

---

## Decision Tree

```
Adding new installation path test?
├── Create test_* function in e2e_test.sh
├── Use --non-interactive with all flags
├── Verify config files were created
├── Verify binaries are functional
└── Add to test execution at bottom

Adding new platform support?
├── Create Dockerfile.{platform} in e2e/
├── Install platform-specific dependencies
├── Create non-root test user with sudo
├── Copy binary and test script
└── Add to docker-test.sh matrix

Testing backup system?
├── Use setup_fake_configs() helper
├── Run with --backup=true or --backup=false
├── Check for .gentleman-backup-* directories
├── Verify backup contents
└── Test restore functionality
```

---

## Code Examples

### Example 1: Complete Installation Test

```bash
test_zsh_zellij() {
    log_test "Install: Zsh + Zellij (no nvim)"

    # Run installation
    if GENTLEMAN_VERBOSE=1 gentleman-dots --non-interactive \
        --shell=zsh --wm=zellij --backup=false 2>&1; then

        # Verify .zshrc exists
        if [ -f "$HOME/.zshrc" ]; then
            log_pass ".zshrc was created"
        else
            log_fail ".zshrc not found"
            return
        fi

        # Verify Zellij config in .zshrc (not tmux!)
        if grep -q "ZELLIJ" "$HOME/.zshrc"; then
            log_pass ".zshrc contains ZELLIJ config"
        else
            log_fail ".zshrc missing ZELLIJ config"
        fi

        # Verify NO tmux in .zshrc
        if grep -q 'WM_CMD="tmux"' "$HOME/.zshrc"; then
            log_fail ".zshrc still has tmux (should be zellij)"
        else
            log_pass ".zshrc correctly has no tmux"
        fi
    else
        log_fail "Installation failed"
    fi
}
```

### Example 2: Backup Test

```bash
setup_fake_configs() {
    # Create fake nvim config
    mkdir -p "$HOME/.config/nvim"
    echo "-- Fake nvim config" > "$HOME/.config/nvim/init.lua"

    # Create fake fish config
    mkdir -p "$HOME/.config/fish"
    echo "# Fake fish config" > "$HOME/.config/fish/config.fish"

    # Create fake .zshrc
    echo "# Fake zshrc" > "$HOME/.zshrc"
}

test_backup_creation() {
    log_test "Creating backup of existing configs"

    cleanup_test_env
    setup_fake_configs

    if GENTLEMAN_VERBOSE=1 gentleman-dots --non-interactive \
        --shell=fish --wm=tmux --backup=true 2>&1; then

        backup_count=$(ls -d "$HOME/.gentleman-backup-"* 2>/dev/null | wc -l)

        if [ "$backup_count" -gt 0 ]; then
            log_pass "Backup directory created"
        else
            log_fail "No backup directory found"
        fi
    else
        log_fail "Installation with backup failed"
    fi
}
```

### Example 3: Functional Verification

```bash
test_shell_functional() {
    log_test "Installed shell is functional"

    # Ensure Homebrew is in PATH
    if [ -d "/home/linuxbrew/.linuxbrew/bin" ]; then
        export PATH="/home/linuxbrew/.linuxbrew/bin:$PATH"
    fi

    # Check if fish runs
    if command -v fish >/dev/null 2>&1; then
        if fish -c "echo 'fish works'" 2>/dev/null | grep -q "fish works"; then
            log_pass "Fish shell is functional"
        else
            log_fail "Fish shell not working"
        fi
    fi
}
```

---

## Docker Commands

```bash
# Build and run specific platform
docker build -f e2e/Dockerfile.ubuntu -t gentleman-e2e-ubuntu .
docker run --rm gentleman-e2e-ubuntu

# Run with full E2E tests
docker run --rm -e RUN_FULL_E2E=1 gentleman-e2e-ubuntu

# Run with backup tests only
docker run --rm -e RUN_BACKUP_TESTS=1 gentleman-e2e-ubuntu

# Interactive debugging
docker run --rm -it gentleman-e2e-ubuntu /bin/bash
```

---

## Test Categories

| Variable | Tests Run |
|----------|-----------|
| (default) | Basic binary tests only |
| `RUN_BACKUP_TESTS=1` | Backup system tests |
| `RUN_FULL_E2E=1` | Full installation tests |

---

## Commands

```bash
cd installer/e2e && ./docker-test.sh            # Run all platforms
docker build -f e2e/Dockerfile.alpine -t test . # Build specific
docker run --rm -e RUN_FULL_E2E=1 test          # Run full suite
```

---

## Resources

- **Test Script**: See `installer/e2e/e2e_test.sh` for test patterns
- **Dockerfiles**: See `installer/e2e/Dockerfile.*` for platform configs
- **Non-interactive**: See `installer/internal/tui/non_interactive.go` for CLI flags
- **Documentation**: See `docs/docker-testing.md` for full guide
