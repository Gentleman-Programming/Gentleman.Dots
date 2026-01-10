#!/bin/bash
# Run E2E tests in Docker containers
# Usage: ./run_e2e.sh [debian|alpine|ubuntu|all]

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
INSTALLER_DIR="$(dirname "$SCRIPT_DIR")"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_section() {
    echo ""
    echo -e "${YELLOW}========================================${NC}"
    echo -e "${YELLOW}  $1${NC}"
    echo -e "${YELLOW}========================================${NC}"
    echo ""
}

# Build binary if needed
build_binary() {
    log_info "Building Linux AMD64 binary..."
    cd "$INSTALLER_DIR"
    GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o e2e/gentleman-installer-linux-amd64 ./cmd/gentleman-installer
    log_success "Binary built successfully"
}

# Run tests on Debian (no bash initially)
test_debian() {
    log_section "Testing on Debian (sh fallback)"
    
    cd "$SCRIPT_DIR"
    
    docker build -f Dockerfile.debian -t gentleman-e2e-debian .
    
    if docker run --rm gentleman-e2e-debian; then
        log_success "Debian tests passed"
        return 0
    else
        log_error "Debian tests failed"
        return 1
    fi
}

# Run tests on Alpine (ash shell, no bash)
test_alpine() {
    log_section "Testing on Alpine (ash/sh only)"
    
    cd "$SCRIPT_DIR"
    
    docker build -f Dockerfile.alpine -t gentleman-e2e-alpine .
    
    if docker run --rm gentleman-e2e-alpine; then
        log_success "Alpine tests passed"
        return 0
    else
        log_error "Alpine tests failed"
        return 1
    fi
}

# Run full E2E tests on Ubuntu
test_ubuntu() {
    log_section "Testing on Ubuntu (full E2E)"
    
    cd "$SCRIPT_DIR"
    
    docker build -f Dockerfile.ubuntu -t gentleman-e2e-ubuntu .
    
    if docker run --rm gentleman-e2e-ubuntu; then
        log_success "Ubuntu E2E tests passed"
        return 0
    else
        log_error "Ubuntu E2E tests failed"
        return 1
    fi
}

# Run all tests
test_all() {
    FAILED=0
    
    test_debian || ((FAILED++))
    test_alpine || ((FAILED++))
    test_ubuntu || ((FAILED++))
    
    log_section "Final Results"
    
    if [ $FAILED -eq 0 ]; then
        log_success "All E2E tests passed!"
        return 0
    else
        log_error "$FAILED test suite(s) failed"
        return 1
    fi
}

# Cleanup
cleanup() {
    log_info "Cleaning up..."
    rm -f "$SCRIPT_DIR/gentleman-installer-linux-amd64"
    docker rmi gentleman-e2e-debian gentleman-e2e-alpine gentleman-e2e-ubuntu 2>/dev/null || true
}

# Main
main() {
    TARGET="${1:-all}"
    
    log_section "Gentleman.Dots E2E Test Suite"
    log_info "Target: $TARGET"
    
    # Build binary first
    build_binary
    
    case "$TARGET" in
        debian)
            test_debian
            ;;
        alpine)
            test_alpine
            ;;
        ubuntu)
            test_ubuntu
            ;;
        all)
            test_all
            ;;
        clean)
            cleanup
            ;;
        *)
            echo "Usage: $0 [debian|alpine|ubuntu|all|clean]"
            exit 1
            ;;
    esac
}

main "$@"
