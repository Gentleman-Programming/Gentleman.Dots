#!/bin/bash
# ============================================================================
# Unit Tests for setup.sh
# ============================================================================
# Run: ./skills/setup_test.sh
# ============================================================================

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
SETUP_SCRIPT="$SCRIPT_DIR/setup.sh"

# Test counters
PASSED=0
FAILED=0

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Temp directory for tests
TEST_DIR=$(mktemp -d)
trap "rm -rf $TEST_DIR" EXIT

# ============================================================================
# Test Helpers
# ============================================================================

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

assert_file_exists() {
    if [ -f "$1" ]; then
        log_pass "File exists: $1"
    else
        log_fail "File missing: $1"
    fi
}

assert_file_contains() {
    if grep -q "$2" "$1" 2>/dev/null; then
        log_pass "File contains '$2': $1"
    else
        log_fail "File missing '$2': $1"
    fi
}

assert_dir_exists() {
    if [ -d "$1" ]; then
        log_pass "Directory exists: $1"
    else
        log_fail "Directory missing: $1"
    fi
}

# ============================================================================
# Setup Tests
# ============================================================================

test_setup_script_exists() {
    log_test "setup.sh exists and is executable"

    if [ -x "$SETUP_SCRIPT" ]; then
        log_pass "setup.sh is executable"
    else
        log_fail "setup.sh is not executable"
    fi
}

test_help_flag() {
    log_test "--help flag works"

    if $SETUP_SCRIPT --help 2>&1 | grep -q "Usage:"; then
        log_pass "--help shows usage"
    else
        log_fail "--help does not show usage"
    fi
}

test_agents_md_exists() {
    log_test "AGENTS.md exists in repo root"

    assert_file_exists "$REPO_ROOT/AGENTS.md"
}

test_skills_directory_exists() {
    log_test "skills/ directory exists"

    assert_dir_exists "$REPO_ROOT/skills"
}

# ============================================================================
# Generation Tests
# ============================================================================

test_generate_claude() {
    log_test "--claude generates CLAUDE.md"

    # Remove existing file
    rm -f "$REPO_ROOT/CLAUDE.md"

    # Generate
    $SETUP_SCRIPT --claude >/dev/null 2>&1

    assert_file_exists "$REPO_ROOT/CLAUDE.md"
    assert_file_contains "$REPO_ROOT/CLAUDE.md" "Auto-generated from AGENTS.md"
    assert_file_contains "$REPO_ROOT/CLAUDE.md" "Claude Code Instructions"
}

test_generate_gemini() {
    log_test "--gemini generates GEMINI.md"

    rm -f "$REPO_ROOT/GEMINI.md"

    $SETUP_SCRIPT --gemini >/dev/null 2>&1

    assert_file_exists "$REPO_ROOT/GEMINI.md"
    assert_file_contains "$REPO_ROOT/GEMINI.md" "Auto-generated from AGENTS.md"
    assert_file_contains "$REPO_ROOT/GEMINI.md" "Gemini CLI Instructions"
}

test_generate_codex() {
    log_test "--codex generates CODEX.md"

    rm -f "$REPO_ROOT/CODEX.md"

    $SETUP_SCRIPT --codex >/dev/null 2>&1

    assert_file_exists "$REPO_ROOT/CODEX.md"
    assert_file_contains "$REPO_ROOT/CODEX.md" "Auto-generated from AGENTS.md"
    assert_file_contains "$REPO_ROOT/CODEX.md" "OpenAI Codex Instructions"
}

test_generate_copilot() {
    log_test "--copilot generates .github/copilot-instructions.md"

    rm -f "$REPO_ROOT/.github/copilot-instructions.md"

    $SETUP_SCRIPT --copilot >/dev/null 2>&1

    assert_file_exists "$REPO_ROOT/.github/copilot-instructions.md"
    assert_file_contains "$REPO_ROOT/.github/copilot-instructions.md" "Auto-generated from AGENTS.md"
    assert_file_contains "$REPO_ROOT/.github/copilot-instructions.md" "GitHub Copilot Instructions"
}

test_generate_all() {
    log_test "--all generates all formats"

    rm -f "$REPO_ROOT/CLAUDE.md" "$REPO_ROOT/GEMINI.md" "$REPO_ROOT/CODEX.md"
    rm -f "$REPO_ROOT/.github/copilot-instructions.md"

    $SETUP_SCRIPT --all >/dev/null 2>&1

    assert_file_exists "$REPO_ROOT/CLAUDE.md"
    assert_file_exists "$REPO_ROOT/GEMINI.md"
    assert_file_exists "$REPO_ROOT/CODEX.md"
    assert_file_exists "$REPO_ROOT/.github/copilot-instructions.md"
}

# ============================================================================
# Content Tests
# ============================================================================

test_generated_files_have_correct_header() {
    log_test "Generated files have correct header format"

    $SETUP_SCRIPT --all >/dev/null 2>&1

    # Check each file has the warning comment
    for file in CLAUDE.md GEMINI.md CODEX.md; do
        if grep -q "Do not edit directly" "$REPO_ROOT/$file" 2>/dev/null; then
            log_pass "$file has 'Do not edit directly' warning"
        else
            log_fail "$file missing warning"
        fi
    done
}

test_generated_files_contain_agents_content() {
    log_test "Generated files contain AGENTS.md content"

    $SETUP_SCRIPT --claude >/dev/null 2>&1

    # Check AGENTS.md content is present
    assert_file_contains "$REPO_ROOT/CLAUDE.md" "Gentleman.Dots AI Agent Skills"
    assert_file_contains "$REPO_ROOT/CLAUDE.md" "Auto-invoke Skills"
}

test_skill_references_are_correct() {
    log_test "Skill references point to correct locations"

    # Repository skills should be in skills/
    assert_file_contains "$REPO_ROOT/AGENTS.md" "skills/gentleman-bubbletea/SKILL.md"
    assert_file_contains "$REPO_ROOT/AGENTS.md" "skills/gentleman-trainer/SKILL.md"

    # User skills should be in GentlemanClaude/skills/
    assert_file_contains "$REPO_ROOT/AGENTS.md" "GentlemanClaude/skills/react-19"
}

# ============================================================================
# Repository Skills Tests
# ============================================================================

test_repo_skills_exist() {
    log_test "Repository skills exist in skills/"

    local skills=("gentleman-bubbletea" "gentleman-trainer" "gentleman-installer"
                  "gentleman-e2e" "gentleman-system" "go-testing")

    for skill in "${skills[@]}"; do
        assert_file_exists "$REPO_ROOT/skills/$skill/SKILL.md"
    done
}

test_user_skills_exist() {
    log_test "User skills exist in GentlemanClaude/skills/"

    local skills=("react-19" "nextjs-15" "typescript" "tailwind-4" "zod-4"
                  "zustand-5" "ai-sdk-5" "django-drf" "playwright" "pytest")

    for skill in "${skills[@]}"; do
        assert_file_exists "$REPO_ROOT/GentlemanClaude/skills/$skill/SKILL.md"
    done
}

# ============================================================================
# Idempotency Tests
# ============================================================================

test_multiple_runs_are_idempotent() {
    log_test "Multiple runs produce same result"

    $SETUP_SCRIPT --claude >/dev/null 2>&1
    local first_hash=$(md5 -q "$REPO_ROOT/CLAUDE.md" 2>/dev/null || md5sum "$REPO_ROOT/CLAUDE.md" | cut -d' ' -f1)

    $SETUP_SCRIPT --claude >/dev/null 2>&1
    local second_hash=$(md5 -q "$REPO_ROOT/CLAUDE.md" 2>/dev/null || md5sum "$REPO_ROOT/CLAUDE.md" | cut -d' ' -f1)

    if [ "$first_hash" = "$second_hash" ]; then
        log_pass "Multiple runs are idempotent"
    else
        log_fail "Multiple runs produce different results"
    fi
}

# ============================================================================
# Error Handling Tests
# ============================================================================

test_invalid_flag_shows_error() {
    log_test "Invalid flag shows error"

    if $SETUP_SCRIPT --invalid-flag 2>&1 | grep -q "Unknown option"; then
        log_pass "Invalid flag shows error"
    else
        log_fail "Invalid flag does not show error"
    fi
}

# ============================================================================
# Run Tests
# ============================================================================

echo ""
echo "════════════════════════════════════════"
echo "  Gentleman.Dots setup.sh Tests"
echo "════════════════════════════════════════"
echo ""

# Basic tests
test_setup_script_exists
test_help_flag
test_agents_md_exists
test_skills_directory_exists

# Generation tests
test_generate_claude
test_generate_gemini
test_generate_codex
test_generate_copilot
test_generate_all

# Content tests
test_generated_files_have_correct_header
test_generated_files_contain_agents_content
test_skill_references_are_correct

# Skills existence tests
test_repo_skills_exist
test_user_skills_exist

# Quality tests
test_multiple_runs_are_idempotent
test_invalid_flag_shows_error

# Summary
echo ""
echo "════════════════════════════════════════"
echo "  Test Summary"
echo "════════════════════════════════════════"
printf "  ${GREEN}Passed: %d${NC}\n" "$PASSED"
printf "  ${RED}Failed: %d${NC}\n" "$FAILED"
echo ""

if [ $FAILED -gt 0 ]; then
    printf "${RED}SOME TESTS FAILED${NC}\n"
    exit 1
else
    printf "${GREEN}ALL TESTS PASSED${NC}\n"
    exit 0
fi
