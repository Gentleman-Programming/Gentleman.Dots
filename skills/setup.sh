#!/bin/bash
# ============================================================================
# Gentleman.Dots AI Skills Setup Script
# ============================================================================
# This script synchronizes AGENTS.md to tool-specific instruction files.
# AGENTS.md is the single source of truth - edits propagate to all tools.
#
# Usage:
#   ./skills/setup.sh              # Interactive menu
#   ./skills/setup.sh --all        # Generate all formats
#   ./skills/setup.sh --claude     # Generate CLAUDE.md only
#   ./skills/setup.sh --gemini     # Generate GEMINI.md only
#   ./skills/setup.sh --copilot    # Generate .github/copilot-instructions.md
#   ./skills/setup.sh --codex      # Generate CODEX.md only
#
# ============================================================================

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'
BOLD='\033[1m'

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# ============================================================================
# Utility Functions
# ============================================================================

log_info() {
    printf "${BLUE}[INFO]${NC} %s\n" "$1"
}

log_success() {
    printf "${GREEN}[SUCCESS]${NC} %s\n" "$1"
}

log_warning() {
    printf "${YELLOW}[WARNING]${NC} %s\n" "$1"
}

log_error() {
    printf "${RED}[ERROR]${NC} %s\n" "$1"
}

log_header() {
    printf "\n${CYAN}${BOLD}════════════════════════════════════════${NC}\n"
    printf "${CYAN}${BOLD}  %s${NC}\n" "$1"
    printf "${CYAN}${BOLD}════════════════════════════════════════${NC}\n\n"
}

# ============================================================================
# Generation Functions
# ============================================================================

# Find all AGENTS.md files in the repository
find_agents_files() {
    find "$REPO_ROOT" -name "AGENTS.md" -type f 2>/dev/null
}

# Generate CLAUDE.md from AGENTS.md
generate_claude() {
    local agents_file="$1"
    local dir=$(dirname "$agents_file")
    local claude_file="$dir/CLAUDE.md"

    log_info "Generating CLAUDE.md from $agents_file"

    # Add Claude-specific header
    cat > "$claude_file" << 'EOF'
# Claude Code Instructions

> **Auto-generated from AGENTS.md** - Do not edit directly.
> Run `./skills/setup.sh --claude` to regenerate.

EOF

    # Append AGENTS.md content
    cat "$agents_file" >> "$claude_file"

    log_success "Created $claude_file"
}

# Generate GEMINI.md from AGENTS.md
generate_gemini() {
    local agents_file="$1"
    local dir=$(dirname "$agents_file")
    local gemini_file="$dir/GEMINI.md"

    log_info "Generating GEMINI.md from $agents_file"

    cat > "$gemini_file" << 'EOF'
# Gemini CLI Instructions

> **Auto-generated from AGENTS.md** - Do not edit directly.
> Run `./skills/setup.sh --gemini` to regenerate.

EOF

    cat "$agents_file" >> "$gemini_file"

    log_success "Created $gemini_file"
}

# Generate .github/copilot-instructions.md from AGENTS.md
generate_copilot() {
    local agents_file="$1"
    local dir=$(dirname "$agents_file")
    local copilot_dir="$dir/.github"
    local copilot_file="$copilot_dir/copilot-instructions.md"

    log_info "Generating copilot-instructions.md from $agents_file"

    mkdir -p "$copilot_dir"

    cat > "$copilot_file" << 'EOF'
# GitHub Copilot Instructions

> **Auto-generated from AGENTS.md** - Do not edit directly.
> Run `./skills/setup.sh --copilot` to regenerate.

EOF

    cat "$agents_file" >> "$copilot_file"

    log_success "Created $copilot_file"
}

# Generate CODEX.md from AGENTS.md
generate_codex() {
    local agents_file="$1"
    local dir=$(dirname "$agents_file")
    local codex_file="$dir/CODEX.md"

    log_info "Generating CODEX.md from $agents_file"

    cat > "$codex_file" << 'EOF'
# OpenAI Codex Instructions

> **Auto-generated from AGENTS.md** - Do not edit directly.
> Run `./skills/setup.sh --codex` to regenerate.

EOF

    cat "$agents_file" >> "$codex_file"

    log_success "Created $codex_file"
}

# Generate all formats for a single AGENTS.md
generate_all_for_file() {
    local agents_file="$1"

    generate_claude "$agents_file"
    generate_gemini "$agents_file"
    generate_copilot "$agents_file"
    generate_codex "$agents_file"
}

# Generate all formats for all AGENTS.md files
generate_all() {
    log_header "Generating All Formats"

    local agents_files=$(find_agents_files)

    if [ -z "$agents_files" ]; then
        log_error "No AGENTS.md files found in repository"
        exit 1
    fi

    for agents_file in $agents_files; do
        log_info "Processing: $agents_file"
        generate_all_for_file "$agents_file"
        echo ""
    done

    log_success "All formats generated!"
}

# ============================================================================
# Interactive Menu
# ============================================================================

show_menu() {
    log_header "Gentleman.Dots AI Skills Setup"

    echo "This script synchronizes AGENTS.md to tool-specific formats."
    echo "AGENTS.md is the single source of truth for all AI assistants."
    echo ""
    echo "Select which assistants to configure:"
    echo ""
    echo "  ${CYAN}1)${NC} Claude Code      (CLAUDE.md)"
    echo "  ${CYAN}2)${NC} Gemini CLI       (GEMINI.md)"
    echo "  ${CYAN}3)${NC} GitHub Copilot   (.github/copilot-instructions.md)"
    echo "  ${CYAN}4)${NC} OpenAI Codex     (CODEX.md)"
    echo "  ${CYAN}5)${NC} All of the above"
    echo ""
    echo "  ${CYAN}0)${NC} Exit"
    echo ""
    echo "  ${YELLOW}Note:${NC} AI tool configs (Claude Code, OpenCode) are now managed by"
    echo "  gentle-ai: https://github.com/Gentleman-Programming/gentle-ai"
    echo ""
    printf "Enter choice [0-5]: "
}

handle_menu_choice() {
    local choice="$1"
    local agents_file="$REPO_ROOT/AGENTS.md"

    if [ ! -f "$agents_file" ]; then
        log_error "AGENTS.md not found at $agents_file"
        exit 1
    fi

    case $choice in
        1)
            generate_claude "$agents_file"
            ;;
        2)
            generate_gemini "$agents_file"
            ;;
        3)
            generate_copilot "$agents_file"
            ;;
        4)
            generate_codex "$agents_file"
            ;;
        5)
            generate_all_for_file "$agents_file"
            ;;
        0)
            log_info "Exiting..."
            exit 0
            ;;
        *)
            log_error "Invalid choice: $choice"
            exit 1
            ;;
    esac
}

interactive_menu() {
    show_menu
    read -r choice
    handle_menu_choice "$choice"
}

# ============================================================================
# CLI Argument Parsing
# ============================================================================

show_help() {
    cat << EOF
Gentleman.Dots AI Skills Setup

Usage: ./skills/setup.sh [OPTIONS]

Options:
  --claude      Generate CLAUDE.md from AGENTS.md
  --gemini      Generate GEMINI.md from AGENTS.md
  --copilot     Generate .github/copilot-instructions.md
  --codex       Generate CODEX.md from AGENTS.md
  --all         Generate all format-specific files
  --help        Show this help message

Note: AI tool configs (Claude Code, OpenCode) are now managed by
  gentle-ai: https://github.com/Gentleman-Programming/gentle-ai

Examples:
  ./skills/setup.sh              # Interactive menu
  ./skills/setup.sh --all        # Generate all formats
  ./skills/setup.sh --claude     # Claude Code only
EOF
}

parse_args() {
    local agents_file="$REPO_ROOT/AGENTS.md"

    if [ ! -f "$agents_file" ]; then
        log_error "AGENTS.md not found at $agents_file"
        exit 1
    fi

    case "$1" in
        --claude)
            generate_claude "$agents_file"
            ;;
        --gemini)
            generate_gemini "$agents_file"
            ;;
        --copilot)
            generate_copilot "$agents_file"
            ;;
        --codex)
            generate_codex "$agents_file"
            ;;
        --all)
            generate_all_for_file "$agents_file"
            ;;
        --help|-h)
            show_help
            exit 0
            ;;
        *)
            log_error "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
}

# ============================================================================
# Main
# ============================================================================

main() {
    cd "$REPO_ROOT"

    if [ $# -eq 0 ]; then
        interactive_menu
    else
        parse_args "$@"
    fi

    echo ""
    log_success "Done!"
}

main "$@"
