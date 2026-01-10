# E2E Tests for Gentleman.Dots Installer

End-to-end tests that verify the installer works across different Linux environments using Docker.

## Table of Contents

- [Quick Start](#quick-start)
- [Test Environments](#test-environments)
- [Requirements](#requirements)
- [Test Scripts](#test-scripts)
- [Adding New Tests](#adding-new-tests)

## Quick Start

```bash
# Interactive mode (TUI menu)
./docker-test.sh

# Run all E2E tests (CI mode)
./docker-test.sh e2e

# Run specific environment
./docker-test.sh e2e debian
./docker-test.sh e2e alpine
./docker-test.sh e2e ubuntu
./docker-test.sh e2e fedora
./docker-test.sh e2e termux

# Run tests with specific platform
./docker-test.sh run debian         # Native platform
./docker-test.sh run debian arm64   # ARM64

# Open interactive shell in container
./docker-test.sh shell alpine
./docker-test.sh shell termux arm64

# Test Termux compatibility (standalone script)
./test-termux.sh arm64  # ARM64 (like real Termux/Android)
./test-termux.sh amd64  # AMD64
```

## Test Environments

| Environment | Shell | Package Manager | Purpose |
|-------------|-------|-----------------|---------|
| Alpine | ash (sh) | apk | BusyBox shell compatibility |
| Debian | sh (dash) | apt | POSIX sh fallback |
| Ubuntu | bash | apt | Full E2E with backup tests |
| Fedora | bash | dnf | RPM-based distro support |
| Termux | sh (simulated) | pkg (simulated) | Android/Termux environment |

## Requirements

| Requirement | Version |
|-------------|---------|
| Docker | Latest |
| Go | 1.21+ |

## Test Scripts

| File | Description |
|------|-------------|
| `docker-test.sh` | Main test runner (interactive + CLI modes) |
| `test-termux.sh` | Standalone Termux compatibility test |
| `e2e_test.sh` | Test cases run inside containers |
| `e2e_test_termux.sh` | Termux-specific test cases |

### Environment Variables

| Variable | Description |
|----------|-------------|
| `RUN_FULL_E2E=1` | Run installation tests (Ubuntu only) |
| `RUN_BACKUP_TESTS=1` | Run backup system tests |
| `GENTLEMAN_VERBOSE=1` | Enable verbose output |

## Adding New Tests

Edit `e2e_test.sh` to add new test cases. Tests must:

1. Be **POSIX-compliant** (no bashisms)
2. Return **non-zero on failure**
3. Output clear **[PASS]/[FAIL]** messages
4. Use provided helper functions: `log_test`, `log_pass`, `log_fail`

```sh
# Example test function
test_my_feature() {
    log_test "My feature works"
    if [ -f "$HOME/.config/myfile" ]; then
        log_pass "Config file exists"
    else
        log_fail "Config file not found"
    fi
}
```
