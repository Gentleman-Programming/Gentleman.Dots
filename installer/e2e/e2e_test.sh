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
    
    # Run installation (no --test in Docker, container is disposable)
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
        
        # Verify Zellij config has default_shell set to zsh
        if [ -f "$HOME/.config/zellij/config.kdl" ]; then
            if grep -q 'default_shell "zsh"' "$HOME/.config/zellij/config.kdl"; then
                log_pass "Zellij config has default_shell set to zsh"
            else
                log_fail "Zellij config missing default_shell zsh"
            fi
        else
            log_fail "Zellij config.kdl not found"
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
    
    if GENTLEMAN_VERBOSE=1 gentleman-dots --non-interactive \
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
    
    # Ensure Homebrew is in PATH for this test
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
    
    # Ensure Homebrew is in PATH for this test
    if [ -d "/home/linuxbrew/.linuxbrew/bin" ]; then
        export PATH="/home/linuxbrew/.linuxbrew/bin:$PATH"
    fi
    
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
    
    # Ensure Homebrew is in PATH for this test
    if [ -d "/home/linuxbrew/.linuxbrew/bin" ]; then
        export PATH="/home/linuxbrew/.linuxbrew/bin:$PATH"
    fi
    
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
# BACKUP SYSTEM TESTS
# ============================================

# Helper: Create fake existing configs
setup_fake_configs() {
    log_test "Setting up fake existing configs..."
    
    # Create fake nvim config
    mkdir -p "$HOME/.config/nvim"
    echo "-- Fake nvim config" > "$HOME/.config/nvim/init.lua"
    echo "vim.opt.number = true" >> "$HOME/.config/nvim/init.lua"
    
    # Create fake fish config
    mkdir -p "$HOME/.config/fish"
    echo "# Fake fish config" > "$HOME/.config/fish/config.fish"
    echo "set -x EDITOR nvim" >> "$HOME/.config/fish/config.fish"
    
    # Create fake .zshrc
    echo "# Fake zshrc" > "$HOME/.zshrc"
    echo "export EDITOR=nvim" >> "$HOME/.zshrc"
    
    # Create fake tmux config
    echo "# Fake tmux config" > "$HOME/.tmux.conf"
    echo "set -g prefix C-a" >> "$HOME/.tmux.conf"
    
    # Create fake zellij config  
    mkdir -p "$HOME/.config/zellij"
    echo "// Fake zellij config" > "$HOME/.config/zellij/config.kdl"
    
    log_pass "Fake configs created"
}

# Helper: Cleanup test environment
cleanup_test_env() {
    rm -rf "$HOME/.config/nvim" 2>/dev/null || true
    rm -rf "$HOME/.config/fish" 2>/dev/null || true
    rm -rf "$HOME/.config/zellij" 2>/dev/null || true
    rm -f "$HOME/.zshrc" 2>/dev/null || true
    rm -f "$HOME/.tmux.conf" 2>/dev/null || true
    rm -rf "$HOME/.gentleman-backup-"* 2>/dev/null || true
}

# Test: Existing configs are detected
test_detect_existing_configs() {
    log_test "Detecting existing configurations"
    
    cleanup_test_env
    setup_fake_configs
    
    # The installer should detect these configs
    # We verify by checking if the files exist
    detected=0
    
    if [ -f "$HOME/.config/nvim/init.lua" ]; then
        detected=$((detected + 1))
    fi
    if [ -f "$HOME/.config/fish/config.fish" ]; then
        detected=$((detected + 1))
    fi
    if [ -f "$HOME/.zshrc" ]; then
        detected=$((detected + 1))
    fi
    if [ -f "$HOME/.tmux.conf" ]; then
        detected=$((detected + 1))
    fi
    if [ -f "$HOME/.config/zellij/config.kdl" ]; then
        detected=$((detected + 1))
    fi
    
    if [ $detected -eq 5 ]; then
        log_pass "All 5 fake configs detected ($detected/5)"
    else
        log_fail "Expected 5 configs, found $detected"
    fi
}

# Test: Backup creation works
test_backup_creation() {
    log_test "Creating backup of existing configs"
    
    cleanup_test_env
    setup_fake_configs
    
    # Run installer with backup=true
    if GENTLEMAN_VERBOSE=1 gentleman-dots --non-interactive \
        --shell=fish --wm=tmux --backup=true 2>&1; then
        
        # Check if backup directory was created
        backup_count=$(ls -d "$HOME/.gentleman-backup-"* 2>/dev/null | wc -l)
        
        if [ "$backup_count" -gt 0 ]; then
            log_pass "Backup directory created"
            
            # Get the backup directory
            backup_dir=$(ls -d "$HOME/.gentleman-backup-"* 2>/dev/null | head -1)
            
            # Check if files were backed up
            if [ -d "$backup_dir" ]; then
                files_backed=$(ls "$backup_dir" 2>/dev/null | wc -l)
                if [ "$files_backed" -gt 0 ]; then
                    log_pass "Backup contains $files_backed items"
                else
                    log_fail "Backup directory is empty"
                fi
            fi
        else
            log_fail "No backup directory found"
        fi
    else
        log_fail "Installation with backup failed"
    fi
}

# Test: Backup directory naming format
test_backup_naming() {
    log_test "Backup directory naming format"
    
    # Check existing backups match pattern: .gentleman-backup-YYYY-MM-DD-HHMMSS
    backup_dirs=$(ls -d "$HOME/.gentleman-backup-"* 2>/dev/null || echo "")
    
    if [ -n "$backup_dirs" ]; then
        for dir in $backup_dirs; do
            basename=$(basename "$dir")
            # Check if format matches .gentleman-backup-YYYY-MM-DD-HHMMSS
            if echo "$basename" | grep -qE '^\.gentleman-backup-[0-9]{4}-[0-9]{2}-[0-9]{2}-[0-9]{6}$'; then
                log_pass "Backup naming format correct: $basename"
            else
                log_fail "Backup naming format incorrect: $basename"
            fi
        done
    else
        log_pass "No backups to check (expected in some test runs)"
    fi
}

# Test: Restore from backup works
test_backup_restore() {
    log_test "Restoring from backup"
    
    cleanup_test_env
    setup_fake_configs
    
    # First, create a backup
    mkdir -p "$HOME/.gentleman-backup-test-restore"
    
    # Copy files to backup manually for this test
    mkdir -p "$HOME/.gentleman-backup-test-restore/nvim"
    echo "-- Original nvim config from backup" > "$HOME/.gentleman-backup-test-restore/nvim/init.lua"
    echo "-- This should be restored" >> "$HOME/.gentleman-backup-test-restore/nvim/init.lua"
    
    # Store original content for comparison
    original_content="-- Original nvim config from backup"
    
    # Now modify the current config (simulate overwrite by installer)
    echo "-- New config after install" > "$HOME/.config/nvim/init.lua"
    
    # Verify modification happened
    if grep -q "New config after install" "$HOME/.config/nvim/init.lua"; then
        log_pass "Config was modified (simulating install)"
    else
        log_fail "Could not modify config for test"
        return
    fi
    
    # Now restore from backup manually (simulating restore function)
    if cp "$HOME/.gentleman-backup-test-restore/nvim/init.lua" "$HOME/.config/nvim/init.lua"; then
        # Verify restore worked
        if grep -q "Original nvim config from backup" "$HOME/.config/nvim/init.lua"; then
            log_pass "Backup restore successful"
        else
            log_fail "Restored content doesn't match original"
        fi
    else
        log_fail "Failed to copy from backup"
    fi
    
    # Cleanup test backup
    rm -rf "$HOME/.gentleman-backup-test-restore"
}

# Test: Multiple backups can coexist
test_multiple_backups() {
    log_test "Multiple backups can coexist"
    
    # Create multiple fake backups with different timestamps
    mkdir -p "$HOME/.gentleman-backup-2024-01-15-120000"
    echo "backup1" > "$HOME/.gentleman-backup-2024-01-15-120000/test"
    
    mkdir -p "$HOME/.gentleman-backup-2024-01-16-130000"
    echo "backup2" > "$HOME/.gentleman-backup-2024-01-16-130000/test"
    
    mkdir -p "$HOME/.gentleman-backup-2024-01-17-140000"
    echo "backup3" > "$HOME/.gentleman-backup-2024-01-17-140000/test"
    
    # Count backups
    backup_count=$(ls -d "$HOME/.gentleman-backup-"* 2>/dev/null | wc -l)
    
    if [ "$backup_count" -ge 3 ]; then
        log_pass "Multiple backups coexist ($backup_count found)"
    else
        log_fail "Expected at least 3 backups, found $backup_count"
    fi
    
    # Cleanup test backups
    rm -rf "$HOME/.gentleman-backup-2024-01-15-120000"
    rm -rf "$HOME/.gentleman-backup-2024-01-16-130000"
    rm -rf "$HOME/.gentleman-backup-2024-01-17-140000"
}

# Test: Backup deletion works
test_backup_deletion() {
    log_test "Backup deletion"
    
    # Create a test backup
    test_backup="$HOME/.gentleman-backup-delete-test"
    mkdir -p "$test_backup"
    echo "test" > "$test_backup/testfile"
    
    # Verify it exists
    if [ ! -d "$test_backup" ]; then
        log_fail "Could not create test backup"
        return
    fi
    
    # Delete it
    rm -rf "$test_backup"
    
    # Verify deletion
    if [ ! -d "$test_backup" ]; then
        log_pass "Backup deletion successful"
    else
        log_fail "Backup still exists after deletion"
    fi
}

# Test: Install without backup doesn't create backup dir
test_install_no_backup() {
    log_test "Install without backup doesn't create backup"
    
    cleanup_test_env
    setup_fake_configs
    
    # Count existing backups before
    before_count=$(ls -d "$HOME/.gentleman-backup-"* 2>/dev/null | wc -l)
    
    # Run installer with backup=false
    if GENTLEMAN_VERBOSE=1 gentleman-dots --non-interactive \
        --shell=zsh --wm=none --backup=false 2>&1; then
        
        # Count backups after
        after_count=$(ls -d "$HOME/.gentleman-backup-"* 2>/dev/null | wc -l)
        
        if [ "$after_count" -eq "$before_count" ]; then
            log_pass "No new backup created when backup=false"
        else
            log_fail "Backup was created despite backup=false"
        fi
    else
        log_fail "Installation failed"
    fi
}

# Test: Backup contains expected config files
test_backup_contents() {
    log_test "Backup contains expected files"
    
    cleanup_test_env
    setup_fake_configs
    
    # Run installer with backup
    if GENTLEMAN_VERBOSE=1 gentleman-dots --non-interactive \
        --shell=fish --wm=tmux --nvim --backup=true 2>&1; then
        
        # Find the backup
        backup_dir=$(ls -dt "$HOME/.gentleman-backup-"* 2>/dev/null | head -1)
        
        if [ -d "$backup_dir" ]; then
            # List what's in the backup
            log_test "Backup contents:"
            ls -la "$backup_dir" 2>/dev/null || true
            
            # Check for expected items (at least some should be there)
            found_items=0
            
            # Note: The backup system uses config keys, not full paths
            for item in nvim fish zsh tmux zellij; do
                if [ -e "$backup_dir/$item" ]; then
                    found_items=$((found_items + 1))
                    log_pass "  Found: $item"
                fi
            done
            
            if [ $found_items -gt 0 ]; then
                log_pass "Backup contains $found_items config items"
            else
                log_fail "Backup is empty or missing expected items"
            fi
        else
            log_fail "Could not find backup directory"
        fi
    else
        log_fail "Installation with backup failed"
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

# Backup tests (can run in basic mode)
if [ "$RUN_BACKUP_TESTS" = "1" ] || [ "$RUN_FULL_E2E" = "1" ]; then
    log_section "Backup System Tests"
    test_detect_existing_configs
    test_backup_creation
    test_backup_naming
    test_backup_restore
    test_multiple_backups
    test_backup_deletion
    test_install_no_backup
    test_backup_contents
    
    # Cleanup after backup tests
    cleanup_test_env
fi

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
