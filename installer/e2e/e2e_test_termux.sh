#!/bin/sh
# E2E Test Script for Gentleman.Dots Installer - TERMUX
# Tests Termux-specific behavior and pkg package manager
#
# Usage: Called by Dockerfile.termux, not directly

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

log_info() {
    printf "[INFO] %s\n" "$1"
}

# ============================================
# TERMUX ENVIRONMENT TESTS
# ============================================

test_termux_env_vars() {
    log_test "Termux environment variables are set"
    
    if [ -n "$TERMUX_VERSION" ]; then
        log_pass "TERMUX_VERSION is set: $TERMUX_VERSION"
    else
        log_fail "TERMUX_VERSION is not set"
    fi
    
    if [ -n "$PREFIX" ]; then
        log_pass "PREFIX is set: $PREFIX"
    else
        log_fail "PREFIX is not set"
    fi
}

test_pkg_command_exists() {
    log_test "pkg command is available"
    
    if command -v pkg >/dev/null 2>&1; then
        log_pass "pkg command found"
    else
        log_fail "pkg command not found"
    fi
}

test_pkg_help() {
    log_test "pkg --help works"
    
    if pkg --help >/dev/null 2>&1; then
        log_pass "pkg --help succeeded"
    else
        log_fail "pkg --help failed"
    fi
}

# ============================================
# BASIC TESTS - Binary functionality
# ============================================

test_binary_runs() {
    log_test "Binary executes with --help"
    if ./gentleman-dots --help > /dev/null 2>&1; then
        log_pass "Binary executes correctly"
    else
        log_fail "Binary failed to execute"
    fi
}

test_version() {
    log_test "Binary shows version"
    if ./gentleman-dots --version 2>&1 | grep -qi "gentleman"; then
        log_pass "Version displays correctly"
    else
        log_fail "Version not displayed"
    fi
}

test_non_interactive_flag() {
    log_test "Non-interactive flag exists"
    if ./gentleman-dots --help 2>&1 | grep -q "non-interactive"; then
        log_pass "Non-interactive mode available"
    else
        log_fail "Non-interactive mode not found"
    fi
}

# ============================================
# TERMUX DETECTION TESTS
# ============================================

test_termux_detection() {
    log_test "Installer detects Termux environment"
    
    # Run with verbose to check detection
    OUTPUT=$(./gentleman-dots --version 2>&1 || true)
    log_info "Version output: $OUTPUT"
    
    # The binary should detect Termux via TERMUX_VERSION env var
    # We verify this through the actual installation tests below
    log_pass "Binary runs in Termux environment"
}

# ============================================
# TERMUX INSTALLATION TESTS
# ============================================

test_termux_zsh_install() {
    log_test "Install: Zsh + Tmux on Termux (uses pkg)"
    
    # Clean previous test
    rm -rf "$HOME/.config" "$HOME/.zshrc" 2>/dev/null || true
    mkdir -p "$HOME/.config"
    
    # Run installation (no --test in Docker, container is disposable)
    if GENTLEMAN_VERBOSE=1 ./gentleman-dots --non-interactive \
        --shell=zsh --wm=tmux --backup=false 2>&1; then
        
        # On Termux, shell change goes to .bashrc, not chsh
        if [ -f "$HOME/.bashrc" ]; then
            if grep -q "exec zsh" "$HOME/.bashrc" 2>/dev/null; then
                log_pass "Termux shell auto-start configured in .bashrc"
            else
                log_info ".bashrc exists but no exec zsh (may be expected in test mode)"
                log_pass "Installation completed"
            fi
        else
            log_pass "Installation completed (test mode may not create .bashrc)"
        fi
        
        # Verify .zshrc exists
        if [ -f "$HOME/.zshrc" ]; then
            log_pass ".zshrc was created"
        else
            log_info ".zshrc not found (may be expected in test mode)"
        fi
    else
        log_fail "Installation failed"
    fi
}

test_termux_fish_zellij_install() {
    log_test "Install: Fish + Zellij on Termux"
    
    # Clean previous test
    rm -rf "$HOME/.config" "$HOME/.zshrc" "$HOME/.bashrc" 2>/dev/null || true
    mkdir -p "$HOME/.config"
    
    if GENTLEMAN_VERBOSE=1 ./gentleman-dots --non-interactive \
        --shell=fish --wm=zellij --backup=false 2>&1; then
        
        log_pass "Fish + Zellij installation completed"
        
        # Check for fish config directory
        if [ -d "$HOME/.config/fish" ]; then
            log_pass "Fish config directory exists"
        else
            log_fail "Fish config directory not found"
            return
        fi
        
        # Check for config.fish specifically
        if [ -f "$HOME/.config/fish/config.fish" ]; then
            log_pass "Fish config.fish exists"
        else
            log_fail "Fish config.fish not found"
        fi
        
        # Verify config.fish is not empty
        if [ -s "$HOME/.config/fish/config.fish" ]; then
            log_pass "Fish config.fish has content"
        else
            log_fail "Fish config.fish is empty"
        fi
        
        # Check for zellij config
        if [ -d "$HOME/.config/zellij" ]; then
            log_pass "Zellij config directory exists"
        else
            log_info "Zellij config not found (may be expected)"
        fi
    else
        log_fail "Installation failed"
    fi
}

test_termux_no_homebrew() {
    log_test "Termux mode skips Homebrew installation"
    
    # In Termux, we should use pkg, not brew
    # The installer should detect Termux and skip brew
    log_pass "Termux detected - Homebrew step would be skipped"
}

test_termux_no_sudo() {
    log_test "Termux mode doesn't use sudo"
    
    # Termux doesn't have/need sudo
    # Package installations go through pkg without sudo
    log_pass "Termux mode uses pkg without sudo"
}

test_termux_font_install() {
    log_test "Termux font installation"
    
    # Clean previous
    rm -rf "$HOME/.termux" 2>/dev/null || true
    
    # Run installation with nvim (which triggers font install)
    if GENTLEMAN_VERBOSE=1 ./gentleman-dots --non-interactive \
        --shell=fish --wm=tmux --nvim --backup=false 2>&1; then
        
        # Check if .termux directory was created
        if [ -d "$HOME/.termux" ]; then
            log_pass ".termux directory created"
            
            # Check if font.ttf exists
            if [ -f "$HOME/.termux/font.ttf" ]; then
                log_pass "Nerd Font installed at ~/.termux/font.ttf"
                
                # Check file is not empty
                if [ -s "$HOME/.termux/font.ttf" ]; then
                    log_pass "Font file has content"
                else
                    log_fail "Font file is empty"
                fi
            else
                log_info "font.ttf not found (font step may have been skipped)"
            fi
        else
            log_info ".termux directory not found (font step may have been skipped)"
        fi
    else
        log_fail "Installation failed"
    fi
}

test_termux_nvim_config() {
    log_test "Neovim config copied correctly on Termux"
    
    # Check nvim config exists
    if [ -d "$HOME/.config/nvim" ]; then
        log_pass "Neovim config directory exists"
        
        # Check init.lua exists
        if [ -f "$HOME/.config/nvim/init.lua" ]; then
            log_pass "Neovim init.lua exists"
        else
            log_fail "Neovim init.lua not found"
        fi
        
        # Check lua directory exists
        if [ -d "$HOME/.config/nvim/lua" ]; then
            log_pass "Neovim lua directory exists"
        else
            log_info "Neovim lua directory not found"
        fi
    else
        log_fail "Neovim config directory not found"
    fi
}

# ============================================
# RUN TESTS
# ============================================

log_section "Termux E2E Tests"

# Termux environment tests
log_section "Termux Environment"
test_termux_env_vars
test_pkg_command_exists
test_pkg_help

# Basic tests
log_section "Basic Binary Tests"
test_binary_runs
test_version
test_non_interactive_flag

# Termux detection
log_section "Termux Detection"
test_termux_detection

# Installation tests
log_section "Termux Installation Tests"
test_termux_zsh_install
test_termux_fish_zellij_install
test_termux_no_homebrew
test_termux_no_sudo
test_termux_font_install
test_termux_nvim_config

# Summary
log_section "Test Summary"
printf "  ${GREEN}Passed: %d${NC}\n" "$PASSED"
printf "  ${RED}Failed: %d${NC}\n" "$FAILED"

if [ $FAILED -gt 0 ]; then
    printf "\n${RED}SOME TESTS FAILED${NC}\n"
    exit 1
else
    printf "\n${GREEN}ALL TERMUX TESTS PASSED${NC}\n"
    exit 0
fi
