#!/bin/sh
# Interactive Docker test runner for Gentleman.Dots installer
# 
# Usage:
#   ./docker-test.sh              # Interactive mode
#   ./docker-test.sh e2e          # Run all E2E tests (non-interactive)
#   ./docker-test.sh e2e ubuntu   # Run E2E for specific image
#   ./docker-test.sh run debian   # Run tests for specific image
#   ./docker-test.sh shell alpine # Open shell in image

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
INSTALLER_DIR="$(dirname "$SCRIPT_DIR")"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m'

# Image configurations
IMAGES="alpine debian ubuntu termux"

get_dockerfile() {
    case "$1" in
        alpine) echo "Dockerfile.alpine" ;;
        debian) echo "Dockerfile.debian" ;;
        ubuntu) echo "Dockerfile.ubuntu" ;;
        termux) echo "Dockerfile.termux" ;;
    esac
}

get_description() {
    case "$1" in
        alpine) echo "Alpine (ash/sh, no bash)" ;;
        debian) echo "Debian (sh, no bash)" ;;
        ubuntu) echo "Ubuntu (bash, full e2e + backup tests)" ;;
        termux) echo "Termux-like (simulated pkg)" ;;
    esac
}

clear_screen() {
    printf "\033[2J\033[H"
}

print_header() {
    echo ""
    echo "${CYAN}╔══════════════════════════════════════════════════════════╗${NC}"
    echo "${CYAN}║${NC}  ${BOLD}$1${NC}"
    echo "${CYAN}╚══════════════════════════════════════════════════════════╝${NC}"
    echo ""
}

print_logo() {
    echo "${CYAN}"
    echo "   ██████╗ ███████╗███╗   ██╗████████╗██╗     ███████╗███╗   ███╗ █████╗ ███╗   ██╗"
    echo "  ██╔════╝ ██╔════╝████╗  ██║╚══██╔══╝██║     ██╔════╝████╗ ████║██╔══██╗████╗  ██║"
    echo "  ██║  ███╗█████╗  ██╔██╗ ██║   ██║   ██║     █████╗  ██╔████╔██║███████║██╔██╗ ██║"
    echo "  ██║   ██║██╔══╝  ██║╚██╗██║   ██║   ██║     ██╔══╝  ██║╚██╔╝██║██╔══██║██║╚██╗██║"
    echo "  ╚██████╔╝███████╗██║ ╚████║   ██║   ███████╗███████╗██║ ╚═╝ ██║██║  ██║██║ ╚████║"
    echo "   ╚═════╝ ╚══════╝╚═╝  ╚═══╝   ╚═╝   ╚══════╝╚══════╝╚═╝     ╚═╝╚═╝  ╚═╝╚═╝  ╚═══╝"
    echo "${NC}"
    echo "                        ${YELLOW}Docker Test Runner${NC}"
    echo ""
}

image_status() {
    if docker images -q "gentleman-test-$1" 2>/dev/null | grep -q .; then
        echo "${GREEN}●${NC}"
    else
        echo "${RED}○${NC}"
    fi
}

# Build the Go binary for Linux
build_binary() {
    echo "${BLUE}→ Building Linux AMD64 binary...${NC}"
    cd "$INSTALLER_DIR"
    GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o "$SCRIPT_DIR/gentleman-installer-linux-amd64" ./cmd/gentleman-installer
    echo "${GREEN}✓ Binary built${NC}"
}

build_image() {
    name="$1"
    platform="$2"
    dockerfile=$(get_dockerfile "$name")

    platform_flag=""
    tag_suffix=""
    if [ -n "$platform" ]; then
        platform_flag="--platform linux/$platform"
        tag_suffix="-$platform"
    fi

    # Ensure binary exists
    if [ ! -f "$SCRIPT_DIR/gentleman-installer-linux-amd64" ]; then
        build_binary
    fi

    echo "${BLUE}→ Building ${name}${tag_suffix}...${NC}"

    cd "$SCRIPT_DIR"
    # shellcheck disable=SC2086
    docker build $platform_flag \
        -f "$dockerfile" \
        -t "gentleman-test-${name}${tag_suffix}" \
        . 2>&1

    echo "${GREEN}✓ Built gentleman-test-${name}${tag_suffix}${NC}"
}

run_image() {
    name="$1"
    platform="$2"
    interactive="${3:-true}"

    tag_suffix=""
    platform_flag=""
    if [ -n "$platform" ]; then
        tag_suffix="-$platform"
        platform_flag="--platform linux/$platform"
    fi

    image_tag="gentleman-test-${name}${tag_suffix}"

    if ! docker images -q "$image_tag" 2>/dev/null | grep -q .; then
        build_image "$name" "$platform"
    fi

    if [ "$interactive" = "true" ]; then
        print_header "Running: $name"
    fi

    # shellcheck disable=SC2086
    if [ "$interactive" = "true" ]; then
        docker run --rm -it $platform_flag "$image_tag"
    else
        docker run --rm $platform_flag "$image_tag"
    fi
    local docker_exit=$?

    if [ "$interactive" = "true" ]; then
        echo ""
        echo "${GREEN}✓ Tests completed${NC}"
    fi
    
    return $docker_exit
}

shell_image() {
    name="$1"
    platform="$2"

    tag_suffix=""
    platform_flag=""
    if [ -n "$platform" ]; then
        tag_suffix="-$platform"
        platform_flag="--platform linux/$platform"
    fi

    image_tag="gentleman-test-${name}${tag_suffix}"

    if ! docker images -q "$image_tag" 2>/dev/null | grep -q .; then
        build_image "$name" "$platform"
    fi

    print_header "Interactive shell: $name"
    echo "${YELLOW}Type 'exit' to return to menu${NC}"
    echo ""

    # shellcheck disable=SC2086
    docker run --rm -it $platform_flag "$image_tag" /bin/sh
}

reset_image() {
    name="$1"
    echo "${YELLOW}→ Removing gentleman-test-$name...${NC}"
    docker rmi -f "gentleman-test-$name" 2>/dev/null || true
    docker rmi -f "gentleman-test-$name-arm64" 2>/dev/null || true
    docker rmi -f "gentleman-test-$name-amd64" 2>/dev/null || true
    echo "${GREEN}✓ Removed${NC}"
}

# Run all E2E tests (non-interactive mode for CI)
run_e2e_all() {
    target="${1:-all}"
    
    print_header "Gentleman.Dots E2E Test Suite"
    
    # Build binary first
    build_binary
    
    FAILED=0
    PASSED=0
    FAILED_IMAGES=""
    PASSED_IMAGES=""
    
    if [ "$target" = "all" ]; then
        images_to_test="$IMAGES"
    else
        images_to_test="$target"
    fi
    
    for img in $images_to_test; do
        echo ""
        echo "${CYAN}════════════════════════════════════════${NC}"
        echo "${CYAN}  Testing: $img${NC}"
        echo "${CYAN}════════════════════════════════════════${NC}"
        echo ""
        
        # Capture output to extract failure details
        test_output_file=$(mktemp)
        if run_image "$img" "" "false" 2>&1 | tee "$test_output_file"; then
            echo "${GREEN}✓ $img tests passed${NC}"
            PASSED=$((PASSED + 1))
            PASSED_IMAGES="$PASSED_IMAGES $img"
        else
            echo "${RED}✗ $img tests failed${NC}"
            FAILED=$((FAILED + 1))
            # Extract failed test names from output
            failed_tests=$(grep -E "^\[FAIL\]" "$test_output_file" | sed 's/\[FAIL\] //' | tr '\n' '; ' | sed 's/; $//')
            if [ -n "$failed_tests" ]; then
                FAILED_IMAGES="$FAILED_IMAGES|$img:$failed_tests"
            else
                FAILED_IMAGES="$FAILED_IMAGES|$img:unknown failure"
            fi
        fi
        rm -f "$test_output_file"
    done
    
    echo ""
    print_header "E2E Test Results"
    echo ""
    
    # Show passed images
    if [ -n "$PASSED_IMAGES" ]; then
        echo "  ${GREEN}✓ PASSED ($PASSED):${NC}"
        for img in $PASSED_IMAGES; do
            echo "    ${GREEN}•${NC} $img"
        done
        echo ""
    fi
    
    # Show failed images with details
    if [ -n "$FAILED_IMAGES" ]; then
        echo "  ${RED}✗ FAILED ($FAILED):${NC}"
        # Parse FAILED_IMAGES (format: |image1:reason1;reason2|image2:reason1)
        echo "$FAILED_IMAGES" | tr '|' '\n' | while read -r entry; do
            if [ -n "$entry" ]; then
                img_name=$(echo "$entry" | cut -d: -f1)
                failures=$(echo "$entry" | cut -d: -f2-)
                echo "    ${RED}•${NC} ${BOLD}$img_name${NC}"
                # Split failures by semicolon and show each
                echo "$failures" | tr ';' '\n' | while read -r fail; do
                    if [ -n "$fail" ]; then
                        echo "      ${RED}└─${NC} $fail"
                    fi
                done
            fi
        done
        echo ""
    fi
    
    # Summary line
    echo "  ─────────────────────────────────────────"
    printf "  Total: %d passed, %d failed\n" "$PASSED" "$FAILED"
    echo ""
    
    if [ $FAILED -gt 0 ]; then
        echo "${RED}══════════════════════════════════════════${NC}"
        echo "${RED}  SOME TESTS FAILED - SEE DETAILS ABOVE${NC}"
        echo "${RED}══════════════════════════════════════════${NC}"
        exit 1
    else
        echo "${GREEN}══════════════════════════════════════════${NC}"
        echo "${GREEN}  ALL TESTS PASSED${NC}"
        echo "${GREEN}══════════════════════════════════════════${NC}"
        exit 0
    fi
}

# Menu functions
select_image() {
    echo "${BOLD}Select image:${NC}"
    echo ""
    i=1
    for img in $IMAGES; do
        status=$(image_status "$img")
        desc=$(get_description "$img")
        printf "  ${CYAN}%d)${NC} %s %-10s - %s\n" "$i" "$status" "$img" "$desc"
        i=$((i + 1))
    done
    echo ""
    printf "  ${CYAN}0)${NC} ← Back"
    echo ""
    echo ""
    printf "${YELLOW}Enter choice:${NC} "
    read -r choice

    case "$choice" in
        1) SELECTED_IMAGE="alpine" ;;
        2) SELECTED_IMAGE="debian" ;;
        3) SELECTED_IMAGE="ubuntu" ;;
        4) SELECTED_IMAGE="termux" ;;
        0|"") SELECTED_IMAGE="" ;;
        *) SELECTED_IMAGE="" ;;
    esac
}

select_platform() {
    echo "${BOLD}Select platform:${NC}"
    echo ""
    echo "  ${CYAN}1)${NC} Native (default)"
    echo "  ${CYAN}2)${NC} ARM64 (Android/Termux)"
    echo "  ${CYAN}3)${NC} AMD64 (Intel/x86)"
    echo ""
    printf "  ${CYAN}0)${NC} ← Back"
    echo ""
    echo ""
    printf "${YELLOW}Enter choice [1]:${NC} "
    read -r choice

    case "$choice" in
        1|"") SELECTED_PLATFORM="" ;;
        2) SELECTED_PLATFORM="arm64" ;;
        3) SELECTED_PLATFORM="amd64" ;;
        0) SELECTED_PLATFORM="BACK" ;;
        *) SELECTED_PLATFORM="" ;;
    esac
}

menu_run() {
    clear_screen
    print_header "Run Tests"
    select_image
    [ -z "$SELECTED_IMAGE" ] && return

    clear_screen
    print_header "Run Tests - $SELECTED_IMAGE"
    select_platform
    [ "$SELECTED_PLATFORM" = "BACK" ] && return

    clear_screen
    run_image "$SELECTED_IMAGE" "$SELECTED_PLATFORM"
    echo ""
    printf "${YELLOW}Press Enter to continue...${NC}"
    read -r _
}

menu_e2e() {
    clear_screen
    print_header "Run ALL E2E Tests"
    echo "${YELLOW}This will run E2E tests on all images.${NC}"
    echo ""
    printf "Continue? [Y/n]: "
    read -r confirm

    if [ "$confirm" != "n" ] && [ "$confirm" != "N" ]; then
        clear_screen
        run_e2e_all "all"
        echo ""
        printf "${YELLOW}Press Enter to continue...${NC}"
        read -r _
    fi
}

menu_shell() {
    clear_screen
    print_header "Interactive Shell"
    select_image
    [ -z "$SELECTED_IMAGE" ] && return

    clear_screen
    print_header "Interactive Shell - $SELECTED_IMAGE"
    select_platform
    [ "$SELECTED_PLATFORM" = "BACK" ] && return

    clear_screen
    shell_image "$SELECTED_IMAGE" "$SELECTED_PLATFORM"
    echo ""
    printf "${YELLOW}Press Enter to continue...${NC}"
    read -r _
}

menu_reset() {
    clear_screen
    print_header "Reset Image"
    select_image
    [ -z "$SELECTED_IMAGE" ] && return

    clear_screen
    print_header "Reset - $SELECTED_IMAGE"
    reset_image "$SELECTED_IMAGE"
    echo ""
    printf "${YELLOW}Press Enter to continue...${NC}"
    read -r _
}

menu_reset_all() {
    clear_screen
    print_header "Reset ALL Images"
    echo "${YELLOW}This will remove all test images.${NC}"
    echo ""
    printf "Are you sure? [y/N]: "
    read -r confirm

    if [ "$confirm" = "y" ] || [ "$confirm" = "Y" ]; then
        for img in $IMAGES; do
            reset_image "$img"
        done
        echo ""
        echo "${GREEN}✓ All images reset${NC}"
    else
        echo "${YELLOW}Cancelled${NC}"
    fi
    echo ""
    printf "${YELLOW}Press Enter to continue...${NC}"
    read -r _
}

menu_status() {
    clear_screen
    print_header "Image Status"
    echo ""
    printf "  ${BOLD}%-12s %-35s %s${NC}\n" "IMAGE" "DESCRIPTION" "STATUS"
    echo "  ─────────────────────────────────────────────────────────"
    for img in $IMAGES; do
        status=$(image_status "$img")
        desc=$(get_description "$img")

        # Check for platform-specific images too
        arm_status=""
        amd_status=""
        if docker images -q "gentleman-test-$img-arm64" 2>/dev/null | grep -q .; then
            arm_status=" ${GREEN}arm64${NC}"
        fi
        if docker images -q "gentleman-test-$img-amd64" 2>/dev/null | grep -q .; then
            amd_status=" ${GREEN}amd64${NC}"
        fi

        printf "  %-12s %-35s %b%b%b\n" "$img" "$desc" "$status" "$arm_status" "$amd_status"
    done
    echo ""
    echo "  ${GREEN}●${NC} = built    ${RED}○${NC} = not built"
    echo ""
    printf "${YELLOW}Press Enter to continue...${NC}"
    read -r _
}

main_menu() {
    while true; do
        clear_screen
        print_logo

        echo "${BOLD}What do you want to do?${NC}"
        echo ""
        echo "  ${CYAN}1)${NC} 🚀 Run tests (single image)"
        echo "  ${CYAN}2)${NC} 🧪 Run ALL E2E tests"
        echo "  ${CYAN}3)${NC} 🐚 Interactive shell"
        echo "  ${CYAN}4)${NC} 🔄 Reset image"
        echo "  ${CYAN}5)${NC} 💣 Reset ALL images"
        echo "  ${CYAN}6)${NC} 📊 View status"
        echo ""
        echo "  ${CYAN}q)${NC} Exit"
        echo ""
        printf "${YELLOW}Enter choice:${NC} "
        read -r choice

        case "$choice" in
            1) menu_run ;;
            2) menu_e2e ;;
            3) menu_shell ;;
            4) menu_reset ;;
            5) menu_reset_all ;;
            6) menu_status ;;
            q|Q|0) clear_screen; echo "${GREEN}Bye!${NC}"; exit 0 ;;
            *) ;;
        esac
    done
}

# Show usage
usage() {
    echo "Usage: $0 [command] [image] [platform]"
    echo ""
    echo "Commands:"
    echo "  (none)          Interactive mode"
    echo "  e2e [image]     Run E2E tests (all images or specific one)"
    echo "  run <image>     Run tests for specific image"
    echo "  shell <image>   Open shell in image"
    echo "  reset <image>   Reset specific image"
    echo "  status          Show image status"
    echo ""
    echo "Images: alpine, debian, ubuntu, termux"
    echo "Platforms: arm64, amd64 (optional)"
    echo ""
    echo "Examples:"
    echo "  $0                    # Interactive mode"
    echo "  $0 e2e                # Run all E2E tests"
    echo "  $0 e2e ubuntu         # Run E2E for Ubuntu only"
    echo "  $0 run debian         # Run Debian tests"
    echo "  $0 shell alpine arm64 # Shell in Alpine ARM64"
}

# Entry point - check for direct commands or run interactive
case "${1:-}" in
    e2e)
        run_e2e_all "${2:-all}"
        ;;
    run)
        if [ -z "$2" ]; then
            echo "${RED}Error: image name required${NC}"
            usage
            exit 1
        fi
        run_image "$2" "$3"
        ;;
    shell)
        if [ -z "$2" ]; then
            echo "${RED}Error: image name required${NC}"
            usage
            exit 1
        fi
        shell_image "$2" "$3"
        ;;
    reset)
        if [ -z "$2" ]; then
            echo "${RED}Error: image name required${NC}"
            usage
            exit 1
        fi
        reset_image "$2"
        ;;
    status)
        menu_status
        ;;
    help|--help|-h)
        usage
        ;;
    "")
        # Interactive mode
        main_menu
        ;;
    *)
        echo "${RED}Unknown command: $1${NC}"
        usage
        exit 1
        ;;
esac
