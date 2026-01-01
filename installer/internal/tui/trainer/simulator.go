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

// Selection represents a visual selection range
type Selection struct {
	StartLine int
	StartCol  int
	EndLine   int
	EndCol    int
	Active    bool // Whether there's an active selection
}

// SimulationResult contains the cursor position and optional selection
type SimulationResult struct {
	Position  SimulatedPosition
	Selection Selection
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
	result := SimulateMotionsWithSelection(start, code, input)
	return result.Position
}

// SimulateMotionsWithSelection simulates vim motions and returns both position and selection
func SimulateMotionsWithSelection(start Position, code []string, input string) SimulationResult {
	result := SimulationResult{
		Position: SimulatedPosition{Line: start.Line, Col: start.Col},
	}

	if len(code) == 0 || input == "" {
		return result
	}

	pos := result.Position
	lastFind := lastFindCommand{}

	i := 0
	for i < len(input) {
		// Check for gu/gU + text object or motion (guiw, gUe, guu, etc.)
		if i+1 < len(input) && input[i] == 'g' && (input[i+1] == 'u' || input[i+1] == 'U') {
			caseOp := input[i+1]
			remaining := input[i+2:]
			if len(remaining) > 0 {
				// Check for gu + text object (guiw, guaw, gui", etc.)
				textObjResult, consumed := tryParseTextObject("v"+remaining, pos, code) // Use 'v' as dummy operator
				if consumed > 1 {                                                       // consumed includes the 'v' we added
					result.Selection = textObjResult
					i += 2 + (consumed - 1) // Skip gu + text object
					continue
				}
				// Check for gu + motion (gue, guw, gu$, etc.)
				sel, consumed := tryParseOperatorMotion(caseOp, remaining, pos, code)
				if consumed > 0 && sel.Active {
					result.Selection = sel
					i += 2 + consumed
					continue
				}
				// Check for guu/gUU (whole line)
				if remaining[0] == 'u' || remaining[0] == 'U' {
					// gu + u or gU + U = whole line
					line := code[pos.Line]
					result.Selection = Selection{
						StartLine: pos.Line,
						StartCol:  0,
						EndLine:   pos.Line,
						EndCol:    len(line) - 1,
						Active:    true,
					}
					i += 3
					continue
				}
			}
			// Just gu without complete motion
			i += 2
			continue
		}

		// Check for text object patterns: [operator] + i/a + object
		// Can be: viw, diw, ciw, yiw OR just iw, aw (implicit visual)
		textObjResult, consumed := tryParseTextObject(input[i:], pos, code)
		if consumed > 0 {
			result.Selection = textObjResult
			i += consumed
			continue
		}

		// Check for operator + motion (dw, de, d$, cw, yw, etc.)
		// Also handles dd, cc, yy (double operator = whole line)
		if input[i] == 'v' || input[i] == 'd' || input[i] == 'c' || input[i] == 'y' {
			operator := input[i]
			remaining := input[i+1:]
			if len(remaining) > 0 {
				// Check for dd, cc, yy (whole line operations)
				if remaining[0] == operator {
					line := code[pos.Line]
					result.Selection = Selection{
						StartLine: pos.Line,
						StartCol:  0,
						EndLine:   pos.Line,
						EndCol:    len(line) - 1,
						Active:    true,
					}
					i += 2
					continue
				}
				// Try to calculate selection for operator + motion
				sel, consumed := tryParseOperatorMotion(operator, remaining, pos, code)
				if consumed > 0 && sel.Active {
					result.Selection = sel
					i += 1 + consumed
					continue
				}
			}
			// Just an operator without complete motion, skip it
			i++
			continue
		}

		// Handle D (delete to end of line) and C (change to end of line)
		if input[i] == 'D' || input[i] == 'C' {
			line := code[pos.Line]
			if len(line) > 0 {
				result.Selection = Selection{
					StartLine: pos.Line,
					StartCol:  pos.Col,
					EndLine:   pos.Line,
					EndCol:    len(line) - 1,
					Active:    true,
				}
			}
			i++
			continue
		}

		// Handle '0' as start-of-line command (not count prefix)
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
			// For gg with count (like 3gg), pass count to go to specific line
			pos = executeGCommand(pos, code, secondChar, count)
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
			case 'G':
				// G - go to last line (or line N if count given)
				if count > 1 {
					// [count]G goes to line [count]
					pos.Line = count - 1
					if pos.Line >= len(code) {
						pos.Line = len(code) - 1
					}
				} else {
					pos.Line = len(code) - 1
				}
				pos = moveFirstNonBlank(pos, code)
			case '{':
				// { - move to previous paragraph (blank line)
				pos = moveParagraphBackward(pos, code)
			case '}':
				// } - move to next paragraph (blank line)
				pos = moveParagraphForward(pos, code)
			case '+':
				// + - move to first non-blank of next line
				if pos.Line < len(code)-1 {
					pos.Line++
					pos = moveFirstNonBlank(pos, code)
				}
			case '-':
				// - - move to first non-blank of previous line
				if pos.Line > 0 {
					pos.Line--
					pos = moveFirstNonBlank(pos, code)
				}
			case '_':
				// _ - move to first non-blank of current line (with count: N-1 lines down)
				if count > 1 {
					pos.Line += count - 1
					if pos.Line >= len(code) {
						pos.Line = len(code) - 1
					}
				}
				pos = moveFirstNonBlank(pos, code)
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

	result.Position = pos
	return result
}

func executeGCommand(pos SimulatedPosition, code []string, secondChar byte, count int) SimulatedPosition {
	switch secondChar {
	case 'e':
		// ge - end of previous word (count times)
		for c := 0; c < count; c++ {
			pos = moveEndOfPrevWord(pos, code, false)
		}
		return pos
	case 'E':
		// gE - end of previous WORD (count times)
		for c := 0; c < count; c++ {
			pos = moveEndOfPrevWord(pos, code, true)
		}
		return pos
	case 'g':
		// [n]gg - go to line n (or first line if no count)
		if count > 1 {
			pos.Line = count - 1 // Line numbers are 1-indexed
			if pos.Line >= len(code) {
				pos.Line = len(code) - 1
			}
		} else {
			pos.Line = 0
		}
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

	if bigWord {
		// For WORD motion: only spaces separate WORDs
		// Skip current WORD (non-space chars)
		for col < len(line) && line[col] != ' ' && line[col] != '\t' {
			col++
		}
		// Skip whitespace
		for col < len(line) && (line[col] == ' ' || line[col] == '\t') {
			col++
		}
	} else {
		// For word motion: words are alphanumeric OR punctuation sequences
		// In Vim, 'w' moves to the start of the next word, where a word is either:
		// - A sequence of word characters (letters, digits, underscore)
		// - A sequence of non-word, non-space characters (punctuation)

		if col < len(line) {
			currentChar := line[col]

			if isWordChar(currentChar, false) {
				// On a word char - skip the rest of this word
				for col < len(line) && isWordChar(line[col], false) {
					col++
				}
			} else if currentChar != ' ' && currentChar != '\t' {
				// On punctuation - skip the rest of this punctuation sequence
				for col < len(line) && !isWordChar(line[col], false) && line[col] != ' ' && line[col] != '\t' {
					col++
				}
			}

			// Now skip whitespace to get to the next word/punctuation
			for col < len(line) && (line[col] == ' ' || line[col] == '\t') {
				col++
			}
		}
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

func moveParagraphForward(pos SimulatedPosition, code []string) SimulatedPosition {
	// Move to next blank line (or end of file)
	line := pos.Line + 1
	for line < len(code) {
		if isBlankLine(code[line]) {
			pos.Line = line
			pos.Col = 0
			return pos
		}
		line++
	}
	// No blank line found, go to last line
	pos.Line = len(code) - 1
	pos.Col = 0
	return pos
}

func moveParagraphBackward(pos SimulatedPosition, code []string) SimulatedPosition {
	// Move to previous blank line (or start of file)
	line := pos.Line - 1
	for line >= 0 {
		if isBlankLine(code[line]) {
			pos.Line = line
			pos.Col = 0
			return pos
		}
		line--
	}
	// No blank line found, go to first line
	pos.Line = 0
	pos.Col = 0
	return pos
}

func isBlankLine(line string) bool {
	for _, ch := range line {
		if ch != ' ' && ch != '\t' {
			return false
		}
	}
	return true
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

// tryParseTextObject attempts to parse a text object pattern and returns the selection
// Patterns: [v/d/c/y] + i/a + object OR just i/a + object
// Returns the selection and number of characters consumed (0 if not a text object)
func tryParseTextObject(input string, pos SimulatedPosition, code []string) (Selection, int) {
	if len(input) < 2 {
		return Selection{}, 0
	}

	offset := 0

	// Check for optional operator prefix (v, d, c, y)
	if input[0] == 'v' || input[0] == 'd' || input[0] == 'c' || input[0] == 'y' {
		offset = 1
		if len(input) < 3 {
			return Selection{}, 0
		}
	}

	// Check for modifier (i = inner, a = around)
	if offset >= len(input) || (input[offset] != 'i' && input[offset] != 'a') {
		return Selection{}, 0
	}
	modifier := input[offset]
	offset++

	// Check for text object
	if offset >= len(input) {
		return Selection{}, 0
	}
	obj := input[offset]
	offset++

	// Validate it's a known text object
	validObjects := "wWsptb()[]{}\"'`<>B"
	if !strings.ContainsRune(validObjects, rune(obj)) {
		return Selection{}, 0
	}

	// Calculate the selection range based on text object type
	sel := calculateTextObjectSelection(pos, code, modifier, obj)
	return sel, offset
}

// tryParseOperatorMotion attempts to parse operator + motion (dw, de, d$, cw, etc.)
// and returns the selection range that would be affected
// Returns the selection and number of characters consumed (0 if not a valid operator+motion)
func tryParseOperatorMotion(operator byte, input string, pos SimulatedPosition, code []string) (Selection, int) {
	if len(input) == 0 || len(code) == 0 {
		return Selection{}, 0
	}

	line := code[pos.Line]
	if len(line) == 0 {
		return Selection{}, 0
	}

	// Parse optional count prefix
	i := 0
	count := 0
	for i < len(input) && input[i] >= '0' && input[i] <= '9' {
		count = count*10 + int(input[i]-'0')
		i++
	}
	if count == 0 {
		count = 1
	}

	if i >= len(input) {
		return Selection{}, 0
	}

	motion := input[i]
	consumed := i + 1

	// Calculate end position after motion
	startPos := pos
	endPos := pos

	// Handle motions that need a character argument (f, F, t, T)
	var motionChar byte
	if (motion == 'f' || motion == 'F' || motion == 't' || motion == 'T') && consumed < len(input) {
		motionChar = input[consumed]
		consumed++
	}

	// Execute motion to find end position
	for c := 0; c < count; c++ {
		switch motion {
		case 'w':
			endPos = moveWordForward(endPos, code, false)
		case 'W':
			endPos = moveWordForward(endPos, code, true)
		case 'e':
			endPos = moveEndOfWord(endPos, code, false)
		case 'E':
			endPos = moveEndOfWord(endPos, code, true)
		case 'b':
			endPos = moveWordBackward(endPos, code, false)
		case 'B':
			endPos = moveWordBackward(endPos, code, true)
		case '$':
			if endPos.Line < len(code) && len(code[endPos.Line]) > 0 {
				endPos.Col = len(code[endPos.Line]) - 1
			}
		case '0':
			endPos.Col = 0
		case '^':
			endPos = moveFirstNonBlank(endPos, code)
		case 'f':
			if motionChar != 0 {
				endPos = findChar(endPos, code, motionChar, true, true)
			}
		case 'F':
			if motionChar != 0 {
				endPos = findChar(endPos, code, motionChar, false, true)
			}
		case 't':
			if motionChar != 0 {
				endPos = findChar(endPos, code, motionChar, true, false)
			}
		case 'T':
			if motionChar != 0 {
				endPos = findChar(endPos, code, motionChar, false, false)
			}
		case 'h':
			if endPos.Col > 0 {
				endPos.Col--
			}
		case 'l':
			if endPos.Col < len(line)-1 {
				endPos.Col++
			}
		case 'j':
			if endPos.Line < len(code)-1 {
				endPos.Line++
			}
		case 'k':
			if endPos.Line > 0 {
				endPos.Line--
			}
		case 'g':
			// Handle gg, ge, gE
			if consumed < len(input) {
				secondChar := input[consumed]
				consumed++
				switch secondChar {
				case 'g':
					endPos.Line = 0
					endPos.Col = 0
				case 'e':
					endPos = moveEndOfPrevWord(endPos, code, false)
				case 'E':
					endPos = moveEndOfPrevWord(endPos, code, true)
				default:
					return Selection{}, 0 // Unknown g command
				}
			} else {
				return Selection{}, 0 // Incomplete g command
			}
		case 'G':
			endPos.Line = len(code) - 1
			endPos = moveFirstNonBlank(endPos, code)
		default:
			return Selection{}, 0 // Unknown motion
		}
	}

	// Create selection from start to end position
	sel := Selection{
		StartLine: startPos.Line,
		StartCol:  startPos.Col,
		EndLine:   endPos.Line,
		EndCol:    endPos.Col,
		Active:    true,
	}

	// For motions like 'w', the selection goes up to but not including the new position
	// For 'e' and '$', it includes the end position
	// Adjust based on motion type
	if motion == 'w' || motion == 'W' {
		// 'w' deletes up to the start of next word, so end is one before endPos
		if endPos.Col > 0 && endPos.Line == startPos.Line {
			sel.EndCol = endPos.Col - 1
		} else if endPos.Line > startPos.Line {
			// Multi-line: select to end of current line
			sel.EndCol = len(code[startPos.Line]) - 1
			sel.EndLine = startPos.Line
		}
	}

	// For backward motions, swap start and end
	if motion == 'b' || motion == 'B' || motion == 'F' || motion == 'T' || motion == 'h' || motion == 'k' {
		if sel.EndCol < sel.StartCol || sel.EndLine < sel.StartLine {
			sel.StartCol, sel.EndCol = sel.EndCol, sel.StartCol
			sel.StartLine, sel.EndLine = sel.EndLine, sel.StartLine
		}
	}

	// Validate selection
	if sel.StartLine >= len(code) || sel.EndLine >= len(code) {
		return Selection{}, 0
	}
	if sel.StartCol < 0 {
		sel.StartCol = 0
	}
	if sel.EndCol < 0 {
		sel.EndCol = 0
	}
	if sel.EndLine == sel.StartLine && sel.EndCol < sel.StartCol {
		// For same-line backward selections that weren't caught above
		sel.StartCol, sel.EndCol = sel.EndCol, sel.StartCol
	}

	return sel, consumed
}

// calculateTextObjectSelection calculates the selection range for a text object
func calculateTextObjectSelection(pos SimulatedPosition, code []string, modifier byte, obj byte) Selection {
	if pos.Line >= len(code) {
		return Selection{}
	}
	line := code[pos.Line]
	if len(line) == 0 {
		return Selection{}
	}

	sel := Selection{
		StartLine: pos.Line,
		EndLine:   pos.Line,
		Active:    false, // Start as false, set to true only if we find valid bounds
	}

	var startCol, endCol int

	switch obj {
	case 'w', 'W':
		// Inner/around word - always succeeds if cursor is on a word char
		bigWord := obj == 'W'
		startCol, endCol = findWordBounds(line, pos.Col, bigWord, modifier == 'a')
		sel.Active = true // Word bounds always return something valid
	case '"', '\'', '`':
		// Inner/around quotes
		startCol, endCol = findQuoteBounds(line, pos.Col, obj, modifier == 'a')
		if startCol >= 0 && endCol >= 0 {
			sel.Active = true
		}
	case '(', ')', 'b':
		// Inner/around parentheses
		startCol, endCol = findPairBounds(line, pos.Col, '(', ')', modifier == 'a')
		if startCol >= 0 && endCol >= 0 {
			sel.Active = true
		}
	case '[', ']':
		// Inner/around brackets
		startCol, endCol = findPairBounds(line, pos.Col, '[', ']', modifier == 'a')
		if startCol >= 0 && endCol >= 0 {
			sel.Active = true
		}
	case '{', '}', 'B':
		// Inner/around braces
		startCol, endCol = findPairBounds(line, pos.Col, '{', '}', modifier == 'a')
		if startCol >= 0 && endCol >= 0 {
			sel.Active = true
		}
	case '<', '>':
		// Inner/around angle brackets OR tag content
		// First try standard <...> pair
		startCol, endCol = findPairBounds(line, pos.Col, '<', '>', modifier == 'a')
		if startCol >= 0 && endCol >= 0 {
			sel.Active = true
		} else {
			// If not inside <...>, try to find tag content (between > and <)
			// This handles HTML like <button>Click me</button>
			startCol, endCol = findTagContentBounds(line, pos.Col, modifier == 'a')
			if startCol >= 0 && endCol >= 0 {
				sel.Active = true
			}
		}
	case 't':
		// Inner/around tag - find HTML tag with content
		startCol, endCol = findHTMLTagBounds(line, pos.Col, modifier == 'a')
		if startCol >= 0 && endCol >= 0 {
			sel.Active = true
		}
	case 's':
		// Inner/around sentence (simplified: treat as word for now)
		startCol, endCol = findWordBounds(line, pos.Col, false, modifier == 'a')
		sel.Active = true
	case 'p':
		// Inner/around paragraph (simplified: whole line for now)
		startCol = 0
		endCol = len(line) - 1
		sel.Active = true
	default:
		return Selection{} // Unknown object
	}

	if !sel.Active {
		return Selection{}
	}

	sel.StartCol = startCol
	sel.EndCol = endCol

	// Ensure valid bounds
	if sel.StartCol < 0 {
		sel.StartCol = 0
	}
	if sel.EndCol >= len(line) {
		sel.EndCol = len(line) - 1
	}
	// For inner selection, if start > end means empty content (like "")
	if sel.EndCol < sel.StartCol {
		sel.Active = false
	}

	return sel
}

// findWordBounds finds the start and end of a word at the given position
func findWordBounds(line string, col int, bigWord bool, around bool) (int, int) {
	if col >= len(line) {
		col = len(line) - 1
	}
	if col < 0 {
		return 0, 0
	}

	// Find start of word
	start := col
	for start > 0 && isWordCharForTextObj(line[start-1], bigWord) {
		start--
	}

	// Find end of word
	end := col
	for end < len(line)-1 && isWordCharForTextObj(line[end+1], bigWord) {
		end++
	}

	// If 'around', include trailing whitespace
	if around {
		for end < len(line)-1 && (line[end+1] == ' ' || line[end+1] == '\t') {
			end++
		}
	}

	return start, end
}

// isWordCharForTextObj checks if a character is part of a word for text object purposes
func isWordCharForTextObj(ch byte, bigWord bool) bool {
	if bigWord {
		return ch != ' ' && ch != '\t'
	}
	r := rune(ch)
	return unicode.IsLetter(r) || unicode.IsDigit(r) || ch == '_'
}

// findQuoteBounds finds the bounds of quoted text
// Returns (-1, -1) if no valid quote pair is found
func findQuoteBounds(line string, col int, quote byte, around bool) (int, int) {
	if col >= len(line) {
		col = len(line) - 1
	}
	if col < 0 {
		return -1, -1
	}

	// Count quotes before cursor to determine if we're inside quotes
	quotesBeforeCursor := 0
	firstQuoteIdx := -1
	for i := 0; i <= col; i++ {
		if line[i] == quote {
			quotesBeforeCursor++
			if firstQuoteIdx == -1 {
				firstQuoteIdx = i
			}
		}
	}

	// If odd number of quotes before (and including) cursor, we might be inside
	// Find the quote pair we're inside
	openIdx := -1
	closeIdx := -1

	if quotesBeforeCursor > 0 {
		// Find opening quote by searching backward
		for i := col; i >= 0; i-- {
			if line[i] == quote {
				openIdx = i
				break
			}
		}

		// Find closing quote by searching forward from after the opening
		if openIdx >= 0 {
			for i := openIdx + 1; i < len(line); i++ {
				if line[i] == quote {
					closeIdx = i
					break
				}
			}
		}

		// Check if cursor is actually between open and close
		if openIdx >= 0 && closeIdx >= 0 && col >= openIdx && col <= closeIdx {
			// We're inside this quote pair
		} else {
			// Not inside, reset
			openIdx = -1
			closeIdx = -1
		}
	}

	// If not found inside, try searching forward for a quote pair
	if openIdx == -1 || closeIdx == -1 {
		openIdx = -1
		closeIdx = -1
		for i := col; i < len(line); i++ {
			if line[i] == quote {
				if openIdx == -1 {
					openIdx = i
				} else {
					closeIdx = i
					break
				}
			}
		}
	}

	if openIdx == -1 || closeIdx == -1 || openIdx >= closeIdx {
		return -1, -1 // No valid quote pair found
	}

	if around {
		return openIdx, closeIdx
	}
	return openIdx + 1, closeIdx - 1
}

// findPairBounds finds the bounds of paired delimiters (parens, brackets, braces, angle brackets)
// Returns (-1, -1) if no valid pair is found
func findPairBounds(line string, col int, open, close byte, around bool) (int, int) {
	if col >= len(line) {
		col = len(line) - 1
	}
	if col < 0 || len(line) == 0 {
		return -1, -1
	}

	// Special case: if cursor is ON the closing delimiter, include it in the backward search
	// by treating it as being just inside the pair
	searchCol := col

	// Strategy: find the innermost pair that contains the cursor position
	// We search backward for an opening delimiter, tracking nesting depth
	openIdx := -1
	depth := 0

	// Search backward for the matching opening delimiter
	for i := searchCol; i >= 0; i-- {
		ch := line[i]
		if ch == close {
			// Found a closing delimiter
			// If this is the cursor position, we want to find the pair that contains it
			if i == col {
				// Cursor is ON closing delimiter - don't count it, look for its opening
				continue
			}
			depth++
		} else if ch == open {
			if depth == 0 {
				// Found our opening delimiter!
				openIdx = i
				break
			}
			// This open matches a close we already passed
			depth--
		}
	}

	// If we didn't find an opening, check if cursor is ON an open delimiter
	if openIdx == -1 && col < len(line) && line[col] == open {
		openIdx = col
	}

	if openIdx == -1 {
		return -1, -1 // No opening delimiter found
	}

	// Now search forward from openIdx for the matching close
	closeIdx := -1
	depth = 0

	for i := openIdx; i < len(line); i++ {
		ch := line[i]
		if ch == open {
			depth++
		} else if ch == close {
			depth--
			if depth == 0 {
				closeIdx = i
				break
			}
		}
	}

	if closeIdx == -1 || openIdx >= closeIdx {
		return -1, -1 // No valid pair found
	}

	// Verify cursor is actually inside this pair
	if col < openIdx || col > closeIdx {
		return -1, -1 // Cursor is outside the pair we found
	}

	if around {
		return openIdx, closeIdx
	}
	// Inner: exclude the delimiters themselves
	return openIdx + 1, closeIdx - 1
}

// findTagContentBounds finds the content between HTML tags (between > and <)
// For example, in "<button>Click me</button>", if cursor is on "Click",
// it returns the bounds of "Click me"
// Returns (-1, -1) if no valid tag content is found
func findTagContentBounds(line string, col int, around bool) (int, int) {
	if col >= len(line) {
		col = len(line) - 1
	}
	if col < 0 || len(line) == 0 {
		return -1, -1
	}

	// Find the > before cursor (end of opening tag)
	openTagEnd := -1
	for i := col; i >= 0; i-- {
		if line[i] == '>' {
			openTagEnd = i
			break
		}
		if line[i] == '<' {
			// We hit an opening < before finding >, cursor might be inside a tag
			break
		}
	}

	if openTagEnd == -1 {
		return -1, -1
	}

	// Find the < after cursor (start of closing tag)
	closeTagStart := -1
	for i := col; i < len(line); i++ {
		if line[i] == '<' {
			closeTagStart = i
			break
		}
		if line[i] == '>' {
			// We hit a > before finding <, cursor might be inside a tag
			break
		}
	}

	if closeTagStart == -1 || openTagEnd >= closeTagStart {
		return -1, -1
	}

	// Content is between openTagEnd+1 and closeTagStart-1
	contentStart := openTagEnd + 1
	contentEnd := closeTagStart - 1

	if contentEnd < contentStart {
		return -1, -1 // Empty content
	}

	if around {
		// "around" for tag content includes the tags themselves
		// Find the start of opening tag
		tagStart := openTagEnd
		for tagStart > 0 && line[tagStart-1] != '<' {
			tagStart--
		}
		if tagStart > 0 {
			tagStart-- // Include the <
		}

		// Find the end of closing tag
		tagEnd := closeTagStart
		for tagEnd < len(line)-1 && line[tagEnd] != '>' {
			tagEnd++
		}

		return tagStart, tagEnd
	}

	return contentStart, contentEnd
}

// findHTMLTagBounds finds an HTML tag element including its content
// For "at" (around tag): <span>Hello</span> - selects entire element
// For "it" (inner tag): <span>Hello</span> - selects just Hello
// Returns (-1, -1) if no valid HTML tag is found
func findHTMLTagBounds(line string, col int, around bool) (int, int) {
	if col >= len(line) {
		col = len(line) - 1
	}
	if col < 0 || len(line) == 0 {
		return -1, -1
	}

	// Strategy: Find the innermost tag pair that contains the cursor
	// 1. Search backward for < to find potential opening tag
	// 2. Match it with its closing tag

	// First, find if we're inside tag content (between > and <)
	// or inside a tag itself (between < and >)

	// Find boundaries around cursor
	prevClose := -1 // Previous >
	prevOpen := -1  // Previous <
	nextOpen := -1  // Next <
	nextClose := -1 // Next >

	for i := col; i >= 0; i-- {
		if line[i] == '>' && prevClose == -1 {
			prevClose = i
		}
		if line[i] == '<' && prevOpen == -1 {
			prevOpen = i
			break
		}
	}

	for i := col; i < len(line); i++ {
		if line[i] == '<' && nextOpen == -1 {
			nextOpen = i
		}
		if line[i] == '>' && nextClose == -1 {
			nextClose = i
			break
		}
	}

	// Case 1: Cursor is inside tag content (after > before <)
	if prevClose != -1 && (prevOpen == -1 || prevClose > prevOpen) {
		// We're in content between tags
		// Need to find the opening tag that contains us

		// Search backward for the opening tag start
		openTagStart := -1
		for i := prevClose; i >= 0; i-- {
			if line[i] == '<' && i+1 < len(line) && line[i+1] != '/' {
				openTagStart = i
				break
			}
		}

		if openTagStart == -1 {
			return -1, -1
		}

		// Extract tag name
		tagNameEnd := openTagStart + 1
		for tagNameEnd < len(line) && line[tagNameEnd] != '>' && line[tagNameEnd] != ' ' {
			tagNameEnd++
		}
		tagName := line[openTagStart+1 : tagNameEnd]

		// Find closing tag </tagname>
		closingTag := "</" + tagName + ">"
		closeTagStart := -1
		for i := prevClose + 1; i <= len(line)-len(closingTag); i++ {
			if line[i:i+len(closingTag)] == closingTag {
				closeTagStart = i
				break
			}
		}

		if closeTagStart == -1 {
			// Try to find any closing tag
			for i := prevClose + 1; i < len(line); i++ {
				if i+1 < len(line) && line[i] == '<' && line[i+1] == '/' {
					// Find the end of this closing tag
					for j := i; j < len(line); j++ {
						if line[j] == '>' {
							if around {
								return openTagStart, j
							}
							// Inner: content between tags
							return prevClose + 1, i - 1
						}
					}
				}
			}
			return -1, -1
		}

		closeTagEnd := closeTagStart + len(closingTag) - 1

		if around {
			return openTagStart, closeTagEnd
		}
		// Inner: just the content
		return prevClose + 1, closeTagStart - 1
	}

	// Case 2: Cursor might be inside a tag itself
	if prevOpen != -1 && (prevClose == -1 || prevOpen > prevClose) && nextClose != -1 {
		// We're inside a tag like <span> or </span>
		// Check if it's an opening or closing tag
		isClosing := prevOpen+1 < len(line) && line[prevOpen+1] == '/'

		if isClosing {
			// Inside closing tag - need to find matching opening
			// This is complex, for now return the closing tag itself
			if around {
				return prevOpen, nextClose
			}
			return -1, -1
		} else {
			// Inside opening tag - find the closing tag
			tagNameEnd := prevOpen + 1
			for tagNameEnd < len(line) && line[tagNameEnd] != '>' && line[tagNameEnd] != ' ' {
				tagNameEnd++
			}
			tagName := line[prevOpen+1 : tagNameEnd]

			// Find closing tag
			closingTag := "</" + tagName + ">"
			closeTagStart := -1
			for i := nextClose + 1; i <= len(line)-len(closingTag); i++ {
				if line[i:i+len(closingTag)] == closingTag {
					closeTagStart = i
					break
				}
			}

			if closeTagStart != -1 {
				closeTagEnd := closeTagStart + len(closingTag) - 1
				if around {
					return prevOpen, closeTagEnd
				}
				// Inner: content between > and <
				return nextClose + 1, closeTagStart - 1
			}
		}
	}

	return -1, -1
}
