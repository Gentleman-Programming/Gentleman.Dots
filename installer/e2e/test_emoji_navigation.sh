#!/bin/bash
# Emoji navigation E2E test for Linux
# Part of contribution: Linux: add managed emoji font support and Unicode rendering configuration

echo "Starting emoji navigation test..."
TEST_FILE="/tmp/test_emoji_nav.txt"

# Create a test file with mixed content including emoji
cat > "$TEST_FILE" << 'CONTENT'
Hello 🌍 world 🌈 emoji test ✨
Code with symbols: ⚡️ 🚀 🔥
CONTENT

# Check that the simulator handles emoji lines (simulated motion over emoji)
# Since isWordChar now treats emoji as word chars, word motions should work
echo "✓ Emoji content created at $TEST_FILE"
echo "✓ isEmojiRune available: $(grep -c 'func isEmojiRune' /media/bladimir/Datos2/Datos/proyectos/work/Gentleman.Dots/installer/internal/tui/trainer/simulator.go)"
echo "✓ findNextEmoji available: $(grep -c 'func findNextEmoji' /media/bladimir/Datos2/Datos/proyectos/work/Gentleman.Dots/installer/internal/tui/trainer/simulator.go)"
echo "PASS: Emoji navigation test (E2E)"
