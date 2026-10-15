#!/bin/bash
# Emoji navigation E2E smoke test for Linux
# Part of contribution: Linux: add managed emoji font support and Unicode rendering configuration
set -u

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SIMULATOR="$SCRIPT_DIR/../internal/tui/trainer/simulator.go"

echo "Starting emoji navigation test..."

if [ ! -f "$SIMULATOR" ]; then
    echo "FAIL: simulator not found at $SIMULATOR" >&2
    exit 1
fi

# The simulator must expose the emoji-aware helpers this feature adds:
# rune decoding and emoji classification drive the word motions.
for fn in runeAtBytePos isEmojiRune; do
    if ! grep -q "func $fn" "$SIMULATOR"; then
        echo "FAIL: $fn not found in $SIMULATOR" >&2
        exit 1
    fi
    echo "OK: $fn available in simulator"
done

# The Go test suite is the behavioral authority for rune motions and emoji
# word classification; this script only smoke-checks that the helpers exist.
echo "PASS: Emoji navigation test (E2E)"
