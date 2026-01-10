#!/bin/bash
# E2E Test Script for Gentleman.Dots Installer
# Runs in Docker container to avoid modifying host system

set -e

PASSED=0
FAILED=0
TESTS=()

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_test() {
    echo -e "${YELLOW}[TEST]${NC} $1"
}

log_pass() {
    echo -e "${GREEN}[PASS]${NC} $1"
    ((PASSED++))
    TESTS+=("PASS: $1")
}

log_fail() {
    echo -e "${RED}[FAIL]${NC} $1"
    ((FAILED++))
    TESTS+=("FAIL: $1")
}

# Test 1: Binary executes
test_binary_executes() {
    log_test "Binary executes with --help"
    if gentleman-dots --help > /dev/null 2>&1; then
        log_pass "Binary executes correctly"
    else
        log_fail "Binary failed to execute"
    fi
}

# Test 2: Binary shows version
test_binary_version() {
    log_test "Binary shows version"
    if gentleman-dots --version 2>&1 | grep -q "gentleman"; then
        log_pass "Version displays correctly"
    else
        log_fail "Version not displayed"
    fi
}

# Test 3: Test mode works
test_test_mode() {
    log_test "Test mode flag is recognized"
    if gentleman-dots --help 2>&1 | grep -q "test"; then
        log_pass "Test mode flag exists"
    else
        log_fail "Test mode flag not found"
    fi
}

# Test 4: Dry-run mode works
test_dry_run() {
    log_test "Dry-run flag is recognized"
    if gentleman-dots --help 2>&1 | grep -q "dry-run"; then
        log_pass "Dry-run flag exists"
    else
        log_fail "Dry-run flag not found"
    fi
}

# Test 5: Git is available (required for clone)
test_git_available() {
    log_test "Git is available"
    if command -v git &> /dev/null; then
        log_pass "Git is installed"
    else
        log_fail "Git is not installed"
    fi
}

# Test 6: Curl is available (required for homebrew)
test_curl_available() {
    log_test "Curl is available"
    if command -v curl &> /dev/null; then
        log_pass "Curl is installed"
    else
        log_fail "Curl is not installed"
    fi
}

# Test 7: Home directory exists
test_home_exists() {
    log_test "Home directory exists"
    if [ -d "$HOME" ]; then
        log_pass "Home directory exists: $HOME"
    else
        log_fail "Home directory not found"
    fi
}

# Test 8: Can create directories in home
test_can_create_dirs() {
    log_test "Can create directories in home"
    TEST_DIR="$HOME/.gentleman-test-$$"
    if mkdir -p "$TEST_DIR" && rmdir "$TEST_DIR"; then
        log_pass "Can create directories in home"
    else
        log_fail "Cannot create directories in home"
    fi
}

# Test 9: Can clone repository
test_can_clone_repo() {
    log_test "Can clone Gentleman.Dots repository"
    CLONE_DIR="$HOME/test-clone-$$"
    if git clone --depth 1 https://github.com/Gentleman-Programming/Gentleman.Dots.git "$CLONE_DIR" 2>/dev/null; then
        log_pass "Repository cloned successfully"
        rm -rf "$CLONE_DIR"
    else
        log_fail "Failed to clone repository"
    fi
}

# Test 10: Shell config files can be created
test_shell_config_creation() {
    log_test "Can create shell config files"
    TEST_ZSHRC="$HOME/.zshrc-test-$$"
    TEST_FISHCONF="$HOME/.config/fish-test-$$/config.fish"
    
    # Test zshrc
    echo "# Test" > "$TEST_ZSHRC"
    if [ -f "$TEST_ZSHRC" ]; then
        rm "$TEST_ZSHRC"
        log_pass "Can create .zshrc"
    else
        log_fail "Cannot create .zshrc"
    fi
    
    # Test fish config
    mkdir -p "$(dirname $TEST_FISHCONF)"
    echo "# Test" > "$TEST_FISHCONF"
    if [ -f "$TEST_FISHCONF" ]; then
        rm -rf "$HOME/.config/fish-test-$$"
        log_pass "Can create fish config"
    else
        log_fail "Cannot create fish config"
    fi
}

# Test 11: Check available shells
test_available_shells() {
    log_test "Check available shells"
    SHELLS_FOUND=""
    
    if command -v bash &> /dev/null; then
        SHELLS_FOUND="$SHELLS_FOUND bash"
    fi
    if command -v sh &> /dev/null; then
        SHELLS_FOUND="$SHELLS_FOUND sh"
    fi
    if command -v zsh &> /dev/null; then
        SHELLS_FOUND="$SHELLS_FOUND zsh"
    fi
    
    if [ -n "$SHELLS_FOUND" ]; then
        log_pass "Available shells:$SHELLS_FOUND"
    else
        log_fail "No shells found"
    fi
}

# Test 12: POSIX sh compatibility
test_posix_compatibility() {
    log_test "POSIX sh compatibility"
    
    # Test basic POSIX constructs
    RESULT=$(/bin/sh -c 'VAR=test; echo $VAR' 2>&1)
    if [ "$RESULT" = "test" ]; then
        log_pass "POSIX sh works correctly"
    else
        log_fail "POSIX sh failed"
    fi
}

# Run all tests
echo ""
echo "========================================"
echo "  Gentleman.Dots E2E Tests"
echo "========================================"
echo ""

test_binary_executes
test_binary_version
test_test_mode
test_dry_run
test_git_available
test_curl_available
test_home_exists
test_can_create_dirs
test_can_clone_repo
test_shell_config_creation
test_available_shells
test_posix_compatibility

# Summary
echo ""
echo "========================================"
echo "  Test Summary"
echo "========================================"
echo ""
for test in "${TESTS[@]}"; do
    echo "  $test"
done
echo ""
echo -e "  ${GREEN}Passed: $PASSED${NC}"
echo -e "  ${RED}Failed: $FAILED${NC}"
echo ""

if [ $FAILED -gt 0 ]; then
    echo -e "${RED}SOME TESTS FAILED${NC}"
    exit 1
else
    echo -e "${GREEN}ALL TESTS PASSED${NC}"
    exit 0
fi
