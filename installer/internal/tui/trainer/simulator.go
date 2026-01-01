package trainer

import (
	"strings"
	"unicode"
)

// SimulatedPosition represents the cursor position after simulating vim motions
type SimulatedPosition struct {
	Line int
	Col  int
}

// lastFindCommand tracks the last f/F/t/T command for ; and , repeats
type lastFindCommand struct {
	cmd       byte // 'f', 'F', 't', or 'T'
	char      byte // the character that was searched for
	hasSearch bool // whether a search has been performed
}

// SimulateMotions takes a starting position, code lines, and a vim command string
// and returns the resulting cursor position after executing those motions
func SimulateMotions(start Position, code []string, input string) SimulatedPosition {
	if len(code) == 0 || input == "" {
		return SimulatedPosition{Line: start.Line, Col: start.Col}
	}

	pos := SimulatedPosition{Line: start.Line, Col: start.Col}
	lastFind := lastFindCommand{}

	i := 0
	for i < len(input) {
		// Handle '0' as start-of-line command (not count prefix)
		// '0' is a command when it's at position 0 or follows another command
		if input[i] == '0' {
			pos.Col = 0
			i++
			continue
		}

		// Parse count prefix (e.g., 3w, 2f.) - only 1-9 can start a count
		count := 0
		for i < len(input) && input[i] >= '1' && input[i] <= '9' {
			count = count*10 + int(input[i]-'0')
			i++
		}
		// After first digit, 0 can be part of count (e.g., 10j)
		for i < len(input) && input[i] >= '0' && input[i] <= '9' {
			count = count*10 + int(input[i]-'0')
			i++
		}
		if count == 0 {
			count = 1
		}

		if i >= len(input) {
			break
		}

		cmd := input[i]
		i++

		// Handle two-character commands (f, F, t, T, g)
		var char byte
		needsChar := cmd == 'f' || cmd == 'F' || cmd == 't' || cmd == 'T'
		isGCommand := cmd == 'g'

		if needsChar && i < len(input) {
			char = input[i]
			i++
			// Save this as the last find command
			lastFind.cmd = cmd
			lastFind.char = char
			lastFind.hasSearch = true
		} else if isGCommand && i < len(input) {
			// Handle g commands (ge, gE, gg, etc.)
			secondChar := input[i]
			i++
			for c := 0; c < count; c++ {
				pos = executeGCommand(pos, code, secondChar)
			}
			continue
		}

		// Execute the command count times
		for c := 0; c < count; c++ {
			switch cmd {
			case 'w':
				pos = moveWordForward(pos, code, false)
			case 'W':
				pos = moveWordForward(pos, code, true)
			case 'e':
				pos = moveEndOfWord(pos, code, false)
			case 'E':
				pos = moveEndOfWord(pos, code, true)
			case 'b':
				pos = moveWordBackward(pos, code, false)
			case 'B':
				pos = moveWordBackward(pos, code, true)
			case '^':
				pos = moveFirstNonBlank(pos, code)
			case '$':
				if pos.Line < len(code) {
					if len(code[pos.Line]) > 0 {
						pos.Col = len(code[pos.Line]) - 1
					} else {
						pos.Col = 0
					}
				}
			case 'f':
				pos = findChar(pos, code, char, true, true)
			case 'F':
				pos = findChar(pos, code, char, false, true)
			case 't':
				pos = findChar(pos, code, char, true, false)
			case 'T':
				pos = findChar(pos, code, char, false, false)
			case ';':
				// Repeat last f/F/t/T in the same direction
				if lastFind.hasSearch {
					switch lastFind.cmd {
					case 'f':
						pos = findChar(pos, code, lastFind.char, true, true)
					case 'F':
						pos = findChar(pos, code, lastFind.char, false, true)
					case 't':
						pos = findChar(pos, code, lastFind.char, true, false)
					case 'T':
						pos = findChar(pos, code, lastFind.char, false, false)
					}
				}
			case ',':
				// Repeat last f/F/t/T in the OPPOSITE direction
				if lastFind.hasSearch {
					switch lastFind.cmd {
					case 'f':
						pos = findChar(pos, code, lastFind.char, false, true) // opposite: backward
					case 'F':
						pos = findChar(pos, code, lastFind.char, true, true) // opposite: forward
					case 't':
						pos = findChar(pos, code, lastFind.char, false, false)
					case 'T':
						pos = findChar(pos, code, lastFind.char, true, false)
					}
				}
			case 'h':
				if pos.Col > 0 {
					pos.Col--
				}
			case 'l':
				if pos.Line < len(code) && pos.Col < len(code[pos.Line])-1 {
					pos.Col++
				}
			case 'j':
				if pos.Line < len(code)-1 {
					pos.Line++
					if pos.Line < len(code) && pos.Col >= len(code[pos.Line]) {
						pos.Col = max(0, len(code[pos.Line])-1)
					}
				}
			case 'k':
				if pos.Line > 0 {
					pos.Line--
					if pos.Line < len(code) && pos.Col >= len(code[pos.Line]) {
						pos.Col = max(0, len(code[pos.Line])-1)
					}
				}
			}
		}
	}

	// Clamp position to valid bounds
	if pos.Line >= len(code) {
		pos.Line = len(code) - 1
	}
	if pos.Line < 0 {
		pos.Line = 0
	}
	if pos.Line < len(code) {
		lineLen := len(code[pos.Line])
		if lineLen == 0 {
			pos.Col = 0
		} else if pos.Col >= lineLen {
			pos.Col = lineLen - 1
		}
		if pos.Col < 0 {
			pos.Col = 0
		}
	}

	return pos
}

func executeGCommand(pos SimulatedPosition, code []string, secondChar byte) SimulatedPosition {
	switch secondChar {
	case 'e':
		// ge - end of previous word
		return moveEndOfPrevWord(pos, code, false)
	case 'E':
		// gE - end of previous WORD
		return moveEndOfPrevWord(pos, code, true)
	case 'g':
		// gg - go to first line
		pos.Line = 0
		pos.Col = 0
		return moveFirstNonBlank(pos, code)
	}
	return pos
}

func moveWordForward(pos SimulatedPosition, code []string, bigWord bool) SimulatedPosition {
	if pos.Line >= len(code) {
		return pos
	}
	line := code[pos.Line]
	col := pos.Col

	// Skip current word
	for col < len(line) && isWordChar(line[col], bigWord) {
		col++
	}
	// Skip non-word characters (punctuation/spaces)
	for col < len(line) && !isWordChar(line[col], bigWord) {
		col++
	}

	// If we reached end of line, try next line
	if col >= len(line) && pos.Line < len(code)-1 {
		pos.Line++
		pos.Col = 0
		// Skip leading spaces on new line
		line = code[pos.Line]
		for pos.Col < len(line) && line[pos.Col] == ' ' {
			pos.Col++
		}
		return pos
	}

	pos.Col = col
	return pos
}

func moveEndOfWord(pos SimulatedPosition, code []string, bigWord bool) SimulatedPosition {
	if pos.Line >= len(code) {
		return pos
	}
	line := code[pos.Line]
	col := pos.Col

	// Move at least one character
	if col < len(line)-1 {
		col++
	}

	// Skip spaces
	for col < len(line) && line[col] == ' ' {
		col++
	}

	// Skip to end of word
	for col < len(line)-1 && isWordChar(line[col+1], bigWord) {
		col++
	}

	pos.Col = col
	return pos
}

func moveWordBackward(pos SimulatedPosition, code []string, bigWord bool) SimulatedPosition {
	if pos.Line >= len(code) {
		return pos
	}
	line := code[pos.Line]
	col := pos.Col

	// Move back at least one
	if col > 0 {
		col--
	}

	// Skip spaces backwards
	for col > 0 && line[col] == ' ' {
		col--
	}

	// If we're at start of line, we're done
	if col == 0 {
		pos.Col = col
		return pos
	}

	// Determine what type of character we're on
	if isWordChar(line[col], bigWord) {
		// On a word char - skip word chars backwards to find start of word
		for col > 0 && isWordChar(line[col-1], bigWord) {
			col--
		}
	} else if line[col] != ' ' {
		// On punctuation - skip punctuation backwards
		for col > 0 && !isWordChar(line[col-1], bigWord) && line[col-1] != ' ' {
			col--
		}
	}

	pos.Col = col
	return pos
}

func moveEndOfPrevWord(pos SimulatedPosition, code []string, bigWord bool) SimulatedPosition {
	if pos.Line >= len(code) {
		return pos
	}
	line := code[pos.Line]
	col := pos.Col

	// Move back at least one
	if col > 0 {
		col--
	}

	// Skip spaces backwards
	for col > 0 && line[col] == ' ' {
		col--
	}

	// We're now at the end of previous word (or at a word char)
	// If we're in the middle of a word, find its end
	// ge goes to end of PREVIOUS word, so we need to go back more if we're in a word

	pos.Col = col
	return pos
}

func moveFirstNonBlank(pos SimulatedPosition, code []string) SimulatedPosition {
	if pos.Line >= len(code) {
		return pos
	}
	line := code[pos.Line]
	for i, ch := range line {
		if ch != ' ' && ch != '\t' {
			pos.Col = i
			return pos
		}
	}
	pos.Col = 0
	return pos
}

func findChar(pos SimulatedPosition, code []string, char byte, forward bool, inclusive bool) SimulatedPosition {
	if pos.Line >= len(code) {
		return pos
	}
	line := code[pos.Line]

	if forward {
		for i := pos.Col + 1; i < len(line); i++ {
			if line[i] == char {
				if inclusive {
					pos.Col = i
				} else {
					pos.Col = i - 1
				}
				return pos
			}
		}
	} else {
		for i := pos.Col - 1; i >= 0; i-- {
			if line[i] == char {
				if inclusive {
					pos.Col = i
				} else {
					pos.Col = i + 1
				}
				return pos
			}
		}
	}
	return pos
}

func isWordChar(ch byte, bigWord bool) bool {
	if bigWord {
		// WORD: only spaces separate words
		return ch != ' ' && ch != '\t'
	}
	// word: letters, digits, underscore
	r := rune(ch)
	return unicode.IsLetter(r) || unicode.IsDigit(r) || ch == '_'
}

// IsValidInput checks if the input so far could be a valid vim motion
func IsValidInput(input string) bool {
	if input == "" {
		return true
	}

	// Valid starting characters for motions
	validStarts := "wWeEbB0^$fFtThljkgG;,"

	// Check if first non-digit char is valid
	i := 0
	for i < len(input) && input[i] >= '0' && input[i] <= '9' {
		i++
	}

	if i >= len(input) {
		// Just digits - could be a count prefix
		return true
	}

	firstCmd := input[i]
	if !strings.ContainsRune(validStarts, rune(firstCmd)) {
		return false
	}

	return true
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
