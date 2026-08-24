#!/bin/sh
# Test installer compatibility in Termux-like environment
# Usage: ./test-termux.sh [arm64|amd64]

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
INSTALLER_DIR="$(dirname "$SCRIPT_DIR")"

PLATFORM="${1:-arm64}"
IMAGE_NAME="gentleman-termux-test"

echo "══════════════════════════════════════════════════════════"
echo "  Testing Termux compatibility (platform: linux/$PLATFORM)"
echo "══════════════════════════════════════════════════════════"

# Build the test image
echo "\n→ Building test image..."
docker build \
    --platform "linux/$PLATFORM" \
    -f "$SCRIPT_DIR/Dockerfile.termux" \
    -t "$IMAGE_NAME:$PLATFORM" \
    "$INSTALLER_DIR"

# Run the tests
echo "\n→ Running POSIX compatibility tests..."
docker run --rm \
    --platform "linux/$PLATFORM" \
    "$IMAGE_NAME:$PLATFORM"

echo "\n══════════════════════════════════════════════════════════"
echo "  ✓ All tests passed on linux/$PLATFORM (no bash)"
echo "══════════════════════════════════════════════════════════"
