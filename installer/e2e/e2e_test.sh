#!/bin/sh
# E2E Test Script for Gentleman.Dots Installer
# Runs REAL installation tests in Docker containers
#
# Usage: Called by Dockerfiles, not directly

set -e

PASSED=0
FAILED=0

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_test() {
    printf "${YELLOW}[TEST]${NC} %s\n" "$1"
}

log_pass() {
    printf "${GREEN}[PASS]${NC} %s\n" "$1"
    PASSED=$((PASSED + 1))
}

log_fail() {
    printf "${RED}[FAIL]${NC} %s\n" "$1"
    FAILED=$((FAILED + 1))
}

log_section() {
    printf "\n${YELLOW}════════════════════════════════════════${NC}\n"
    printf "${YELLOW}  %s${NC}\n" "$1"
    printf "${YELLOW}════════════════════════════════════════${NC}\n\n"
}

# ============================================
# BASIC TESTS - Binary functionality
# ============================================

test_binary_runs() {
    log_test "Binary executes with --help"
    if gentleman-dots --help > /dev/null 2>&1; then
        log_pass "Binary executes correctly"
    else
        log_fail "Binary failed to execute"
    fi
}

test_version() {
    log_test "Binary shows version"
    if gentleman-dots --version 2>&1 | grep -q "gentleman"; then
        log_pass "Version displays correctly"
    else
        log_fail "Version not displayed"
    fi
}

test_non_interactive_flag() {
    log_test "Non-interactive flag exists"
    if gentleman-dots --help 2>&1 | grep -q "non-interactive"; then
        log_pass "Non-interactive mode available"
    else
        log_fail "Non-interactive mode not found"
    fi
}

# ============================================
# INSTALLATION TESTS - Real E2E
# ============================================

# Test: Zsh + Zellij (no nvim, no terminal)
test_zsh_zellij() {
    log_test "Install: Zsh + Zellij (no nvim)"
    
    # Run installation
    if GENTLEMAN_VERBOSE=1 gentleman-dots --test --non-interactive \
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

# Test: Fish + Tmux + Nvim
test_fish_tmux_nvim() {
    log_test "Install: Fish + Tmux + Nvim"
    
    # Clean previous test
    rm -rf "$HOME/.config" "$HOME/.zshrc" 2>/dev/null || true
    mkdir -p "$HOME/.config"
    
    if GENTLEMAN_VERBOSE=1 gentleman-dots --test --non-interactive \
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
        
        # Verify init.lua exists
        if [ -f "$HOME/.config/nvim/init.lua" ]; then
            log_pass "Neovim init.lua exists"
        else
            log_fail "Neovim init.lua not found"
        fi
        
        # Verify tmux config in fish (not zellij)
        if grep -qi "tmux" "$HOME/.config/fish/config.fish"; then
            log_pass "Fish config contains tmux"
        else
            log_fail "Fish config missing tmux"
        fi
    else
        log_fail "Installation failed"
    fi
}

# Test: Nushell + None WM
test_nushell_no_wm() {
    log_test "Install: Nushell + No WM"
    
    # Clean previous test
    rm -rf "$HOME/.config" "$HOME/.zshrc" 2>/dev/null || true
    mkdir -p "$HOME/.config"
    
    if GENTLEMAN_VERBOSE=1 gentleman-dots --test --non-interactive \
        --shell=nushell --wm=none --backup=false 2>&1; then
        
        # Verify nushell config exists
        if [ -d "$HOME/.config/nushell" ]; then
            log_pass "Nushell config directory created"
        else
            log_fail "Nushell config directory not found"
            return
        fi
        
        # Check for config.nu
        if [ -f "$HOME/.config/nushell/config.nu" ]; then
            log_pass "Nushell config.nu exists"
        else
            log_fail "Nushell config.nu not found"
        fi
    else
        log_fail "Installation failed"
    fi
}

# Test: Verify shell is installed and functional
test_shell_functional() {
    log_test "Installed shell is functional"
    
    # Check if fish runs
    if command -v fish >/dev/null 2>&1; then
        if fish -c "echo 'fish works'" 2>/dev/null | grep -q "fish works"; then
            log_pass "Fish shell is functional"
        else
            log_fail "Fish shell not working"
        fi
    fi
    
    # Check if zsh runs  
    if command -v zsh >/dev/null 2>&1; then
        if zsh -c "echo 'zsh works'" 2>/dev/null | grep -q "zsh works"; then
            log_pass "Zsh shell is functional"
        else
            log_fail "Zsh shell not working"
        fi
    fi
}

# Test: WM is installed
test_wm_installed() {
    log_test "Window manager is installed"
    
    if command -v tmux >/dev/null 2>&1; then
        log_pass "Tmux is installed"
    fi
    
    if command -v zellij >/dev/null 2>&1; then
        log_pass "Zellij is installed"
    fi
}

# Test: Nvim is installed and configured
test_nvim_configured() {
    log_test "Neovim is properly configured"
    
    if command -v nvim >/dev/null 2>&1; then
        log_pass "Neovim binary is installed"
        
        # Check if it can start (headless)
        if nvim --headless -c "q" 2>/dev/null; then
            log_pass "Neovim starts successfully"
        else
            log_fail "Neovim failed to start"
        fi
    else
        log_fail "Neovim not installed"
    fi
}

# ============================================
# RUN TESTS
# ============================================

log_section "Gentleman.Dots E2E Tests"

# Basic tests first
log_section "Basic Tests"
test_binary_runs
test_version
test_non_interactive_flag

# Installation tests (only if we have the full environment)
if [ "$RUN_FULL_E2E" = "1" ]; then
    log_section "Installation Tests"
    test_zsh_zellij
    test_fish_tmux_nvim
    test_nushell_no_wm
    
    log_section "Verification Tests"
    test_shell_functional
    test_wm_installed
    test_nvim_configured
fi

# Summary
log_section "Test Summary"
printf "  ${GREEN}Passed: %d${NC}\n" "$PASSED"
printf "  ${RED}Failed: %d${NC}\n" "$FAILED"

if [ $FAILED -gt 0 ]; then
    printf "\n${RED}SOME TESTS FAILED${NC}\n"
    exit 1
else
    printf "\n${GREEN}ALL TESTS PASSED${NC}\n"
    exit 0
fi
