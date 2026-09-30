# Contribution: Linux: add managed emoji font support and Unicode rendering configuration

## Description

This contribution adds full Unicode and emoji support to the Vim simulator and TUI in Gentleman.Dots, specifically for Linux environments where input may contain multi-byte characters.

## Current Problem

The Vim simulator in `installer/internal/tui/trainer/simulator.go` uses `byte` for character processing:

```go
func isWordChar(ch byte, bigWord bool) bool {
    // Only works for ASCII
    return unicode.IsLetter(rune(ch)) || unicode.IsDigit(r) || ch == '_'
}
```

This fails with:

- Non-ASCII Unicode characters (accents, Chinese, Cyrillic, etc.)
- Emoji (😀, ❤️, etc.)
- Double-width characters (CJK, full-width emoji)

## Proposed Solution

### 1. Update `isWordChar` to handle runes

Change signature from `byte` to `rune` and use proper Unicode word detection:

```go
func isWordChar(r rune, bigWord bool) bool {
    if bigWord {
        // WORD: only spaces separate words (also consider control characters)
        return r != ' ' && r != '\t' && !unicode.IsControl(r)
    }
    // word: letters, digits, underscore (includes Unicode)
    return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}
```

### 2. Update `isWordCharForTextObj` internally

Align signature and convert byte to rune before calling `isWordChar`:

```go
func isWordCharForTextObj(ch byte, bigWord bool) bool {
    // ... convert to rune before calling isWordChar
}
```

### 3. UTF-8 input handling in simulator

Ensure all code lines respect rune boundaries, not byte boundaries:

```go
// Instead of len(line), use:
func runeCountInLine(line string) int {
    count := 0
    for range line {
        count++
    }
    return count
}
```

### 4. Emoji support in navigation

Add commands for moving by emoji:

```go
// New case for emoji in command switch
case 'M-<emoji>':
    // Navigate to next emoji
    pos = findNextEmoji(pos, code, true)
```

### 5. Emoji search function

```go
func findNextEmoji(pos SimulatedPosition, code []string, forward bool) SimulatedPosition {
    // Find next character that is an emoji (general category "So" - Symbol, other)
    for lineIdx := pos.Line; lineIdx < len(code); lineIdx++ {
        line := code[lineIdx]
        startCol := 0
        if lineIdx == pos.Line {
            startCol = pos.Col
        }
        
        for col := startCol; col < len(line); {
            // Get next rune
            if col >= len(line) {
                break
            }
            r, size := utf8.DecodeRuneInString(line[col:])
            if r == utf8.RuneError {
                col++ // Invalid UTF-8, advance one byte
                continue
            }
            // Check if emoji (So category or emoji ranges)
            if unicode.Is(unicode.Sym, r) || (r >= 0x1F300 && r <= 0x1F5FF) || (r >= 0x1F900 && r <= 0x1F9FF) {
                // Found
                pos.Line = lineIdx
                pos.Col = col
                return pos
            }
            col += size
        }
    }
    return pos
}
```

### 6. Selection operations support

- `dw` over emoji should select the full emoji
- `cw` to change emoji
- `yw` to copy emoji
- Search with `f`/`t` using Unicode code points

### 7. Linux testing

Add tests in `e2e/` to validate:

```bash
# Emoji navigation test
test_emoji_navigation() {
    # Create test file with emoji
    echo "Hello 🌍 World" > /tmp/test_emoji.txt
    # Run simulator and verify motions work
    ./gentleman-installer --train --file /tmp/test_emoji.txt
    # Verify w/e/b motions work over emoji
}
```

## Implementation Considerations

1. **Go `unicode/utf8` package**: Available in stdlib, use `utf8.DecodeRuneInString` for rune iteration.

2. **`strings` vs `[]rune`**: Convert string to `[]rune` when needed, but consider memory overhead for large files.

3. **Emoji width**: Emoji may have variable terminal width. Consider Go `width` library or manual Unicode category calculations.

4. **Backward compatibility**: Maintain ASCII compatibility—existing functions must continue working with ASCII while adding Unicode support.

## Change Checklist

- [ ] Update `isWordChar` to accept `rune` instead of `byte`
- [ ] Update internal `isWordCharForTextObj`
- [ ] Add `runeCountInLine` helper
- [ ] Add emoji search (`findNextEmoji`)
- [ ] Update `moveWordForward`/`moveWordBackward` for rune boundaries
- [ ] Add E2E tests for emoji navigation
- [ ] Update AGENTS.md documentation

## References

- [Go Unicode Documentation](https://go.dev/doc/unicode)
- [Rune basics in Go](https://go.dev/tour/basics/11)
- [Unicode emoji categories](https://unicode.org/emoji/charts/full-emoji-list.html)
