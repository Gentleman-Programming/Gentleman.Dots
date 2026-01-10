# E2E Tests for Gentleman.Dots Installer

End-to-end tests that verify the installer works across different environments.

## Quick Start

```bash
# Run all tests
./run_e2e.sh all

# Run specific environment
./run_e2e.sh debian   # Tests sh fallback (no bash)
./run_e2e.sh alpine   # Tests ash shell
./run_e2e.sh ubuntu   # Full E2E tests
./run_e2e.sh fedora   # Tests dnf package manager

# Test Termux compatibility
./test-termux.sh arm64  # Tests ARM64 (like real Termux)
./test-termux.sh amd64  # Tests AMD64
```

## Test Environments

| Environment | Shell | Purpose |
|-------------|-------|---------|
| Debian | sh (no bash) | Verifies installer works without bash |
| Alpine | ash | Verifies BusyBox shell compatibility |
| Ubuntu | bash | Full E2E with all features |
| Fedora | bash | Verifies dnf package manager support |
| Termux | sh (Alpine sim) | Simulates Termux/Android environment |

## Requirements

- Docker
- Go 1.21+

## Adding New Tests

Edit `e2e_test.sh` to add new test cases. Tests should:
1. Be POSIX-compliant (no bashisms)
2. Return non-zero on failure
3. Output clear success/failure messages
