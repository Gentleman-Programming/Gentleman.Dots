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
#   ./skills/setup.sh --opencode   # Sync to OpenCode config
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

# Sync skills to OpenCode config directory
sync_opencode() {
    local opencode_dir="$HOME/.config/opencode/skill"

    log_info "Syncing skills to OpenCode config..."

    if [ ! -d "$opencode_dir" ]; then
        log_warning "OpenCode skill directory not found: $opencode_dir"
        log_info "Creating directory..."
        mkdir -p "$opencode_dir"
    fi

    # Copy AGENTS.md as the main instruction file
    if [ -f "$REPO_ROOT/AGENTS.md" ]; then
        cp "$REPO_ROOT/AGENTS.md" "$opencode_dir/AGENTS.md"
        log_success "Copied AGENTS.md to OpenCode"
    fi

    # Copy individual skills
    local skills_dir="$REPO_ROOT/GentlemanClaude/skills"
    if [ -d "$skills_dir" ]; then
        for skill_dir in "$skills_dir"/*/; do
            if [ -f "$skill_dir/SKILL.md" ]; then
                skill_name=$(basename "$skill_dir")
                mkdir -p "$opencode_dir/$skill_name"
                cp "$skill_dir/SKILL.md" "$opencode_dir/$skill_name/SKILL.md"
                log_info "  → Copied $skill_name"
            fi
        done
        log_success "Synced all skills to OpenCode"
    fi
}

# Sync skills to Claude Code config directory
sync_claude_config() {
    local claude_dir="$HOME/.claude/skills"

    log_info "Syncing skills to Claude Code config..."

    if [ ! -d "$claude_dir" ]; then
        log_warning "Claude skills directory not found: $claude_dir"
        log_info "Creating directory..."
        mkdir -p "$claude_dir"
    fi

    # Copy individual skills
    local skills_dir="$REPO_ROOT/GentlemanClaude/skills"
    if [ -d "$skills_dir" ]; then
        for skill_dir in "$skills_dir"/*/; do
            if [ -f "$skill_dir/SKILL.md" ]; then
                skill_name=$(basename "$skill_dir")
                mkdir -p "$claude_dir/$skill_name"

                # Remove existing file if read-only, then copy
                local dest_file="$claude_dir/$skill_name/SKILL.md"
                if [ -f "$dest_file" ]; then
                    chmod u+w "$dest_file" 2>/dev/null || true
                fi
                cp -f "$skill_dir/SKILL.md" "$dest_file"

                # Copy assets if they exist
                if [ -d "$skill_dir/assets" ]; then
                    chmod -R u+w "$claude_dir/$skill_name/assets" 2>/dev/null || true
                    cp -rf "$skill_dir/assets" "$claude_dir/$skill_name/"
                fi

                log_info "  → Copied $skill_name"
            fi
        done
        log_success "Synced all skills to ~/.claude/skills/"
    fi
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
    echo "  ${CYAN}6)${NC} Sync to ~/.claude/skills/"
    echo "  ${CYAN}7)${NC} Sync to OpenCode config"
    echo "  ${CYAN}8)${NC} Sync to all user configs"
    echo ""
    echo "  ${CYAN}0)${NC} Exit"
    echo ""
    printf "Enter choice [0-8]: "
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
        6)
            sync_claude_config
            ;;
        7)
            sync_opencode
            ;;
        8)
            sync_claude_config
            sync_opencode
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
  --sync-claude Sync skills to ~/.claude/skills/
  --sync-opencode Sync skills to OpenCode config
  --sync-all    Sync skills to all user config directories
  --help        Show this help message

Examples:
  ./skills/setup.sh              # Interactive menu
  ./skills/setup.sh --all        # Generate all formats
  ./skills/setup.sh --claude     # Claude Code only
  ./skills/setup.sh --sync-all   # Sync to user configs
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
        --sync-claude)
            sync_claude_config
            ;;
        --sync-opencode)
            sync_opencode
            ;;
        --sync-all)
            sync_claude_config
            sync_opencode
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
