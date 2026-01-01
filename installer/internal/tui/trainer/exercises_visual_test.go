package trainer

import (
	"fmt"
	"testing"
)

// TestChangeRepeatExercisesVisualSelection tests ALL Change & Repeat exercises
// to ensure visual selection highlighting is pedagogically correct
func TestChangeRepeatExercisesVisualSelection(t *testing.T) {
	tests := []struct {
		name         string
		exerciseID   string
		code         []string
		cursorCol    int
		cursorLine   int
		input        string
		wantActive   bool
		wantStartCol int
		wantEndCol   int
		description  string
	}{
		// changerepeat_001: dw on 'badVariable'
		{
			name:         "dw deletes word plus space",
			exerciseID:   "changerepeat_001",
			code:         []string{"const badVariable = 'delete me';"},
			cursorLine:   0,
			cursorCol:    6,
			input:        "dw",
			wantActive:   true,
			wantStartCol: 6,
			wantEndCol:   17, // badVariable + space (not the '=')
			description:  "dw should select 'badVariable ' (word + trailing space)",
		},
		// changerepeat_001: de on 'badVariable'
		{
			name:         "de deletes word only",
			exerciseID:   "changerepeat_001",
			code:         []string{"const badVariable = 'delete me';"},
			cursorLine:   0,
			cursorCol:    6,
			input:        "de",
			wantActive:   true,
			wantStartCol: 6,
			wantEndCol:   16, // just badVariable, no space
			description:  "de should select 'badVariable' (word only, no space)",
		},
		// changerepeat_002: de on 'oldName'
		{
			name:         "de on oldName before parens",
			exerciseID:   "changerepeat_002",
			code:         []string{"function oldName() { return 'old implementation'; }"},
			cursorLine:   0,
			cursorCol:    9,
			input:        "de",
			wantActive:   true,
			wantStartCol: 9,
			wantEndCol:   15, // oldName
			description:  "de should select 'oldName' without the parenthesis",
		},
		// changerepeat_003: d$ delete to end
		{
			name:         "d$ deletes to end of line",
			exerciseID:   "changerepeat_003",
			code:         []string{"const result = oldValue + ' some extra stuff to remove';"},
			cursorLine:   0,
			cursorCol:    15,
			input:        "d$",
			wantActive:   true,
			wantStartCol: 15,
			wantEndCol:   55, // end of line
			description:  "d$ should select from cursor to end of line",
		},
		// changerepeat_004: dd deletes entire line
		{
			name:         "dd selects entire line",
			exerciseID:   "changerepeat_004",
			code:         []string{"const valid = true;", "const DELETE_THIS_LINE = 'garbage';", "const alsoValid = false;"},
			cursorLine:   1,
			cursorCol:    0,
			input:        "dd",
			wantActive:   true,
			wantStartCol: 0,
			wantEndCol:   34, // entire line
			description:  "dd should select the entire line",
		},
		// changerepeat_005: D deletes to end
		{
			name:         "D deletes to end of line",
			exerciseID:   "changerepeat_005",
			code:         []string{"const config = { name: 'app' /* DELETE: old comment */ };"},
			cursorLine:   0,
			cursorCol:    27,
			input:        "D",
			wantActive:   true,
			wantStartCol: 27,
			wantEndCol:   56,
			description:  "D should select from cursor to end of line",
		},
		// changerepeat_006: cw changes word
		{
			name:         "cw selects word to change",
			exerciseID:   "changerepeat_006",
			code:         []string{"const oldName = 'value';"},
			cursorLine:   0,
			cursorCol:    6,
			input:        "cw",
			wantActive:   true,
			wantStartCol: 6,
			wantEndCol:   13, // oldName + space before =
			description:  "cw should select 'oldName ' for change",
		},
		// changerepeat_007: ciw changes inner word
		// The exercise cursor should be on 'item' at col 28
		{
			name:         "ciw selects inner word (item)",
			exerciseID:   "changerepeat_007",
			code:         []string{"function calculateTotal(items) {", "  return items.reduce((sum, item) => sum + item.price, 0);", "}"},
			cursorLine:   1,
			cursorCol:    28, // 'i' of 'item' (cols 28-31)
			input:        "ciw",
			wantActive:   true,
			wantStartCol: 28, // start of 'item'
			wantEndCol:   31, // end of 'item'
			description:  "ciw should select 'item' (the word under cursor)",
		},
		// changerepeat_008: c$ changes to end
		{
			name:         "c$ selects to end of line",
			exerciseID:   "changerepeat_008",
			code:         []string{"const message = 'Hello, this needs to be completely rewritten';"},
			cursorLine:   0,
			cursorCol:    18,
			input:        "c$",
			wantActive:   true,
			wantStartCol: 18,
			wantEndCol:   62,
			description:  "c$ should select from cursor to end for change",
		},
		// changerepeat_009: cc changes entire line
		{
			name:         "cc selects entire line",
			exerciseID:   "changerepeat_009",
			code:         []string{"function oldFunction() {", "  // This entire line needs rewriting", "  return null;", "}"},
			cursorLine:   1,
			cursorCol:    5,
			input:        "cc",
			wantActive:   true,
			wantStartCol: 0,
			wantEndCol:   36, // len("  // This entire line needs rewriting") - 1 = 36
			description:  "cc should select the entire line for change",
		},
		// changerepeat_010: C changes to end
		{
			name:         "C selects to end of line",
			exerciseID:   "changerepeat_010",
			code:         []string{"const data = fetchData() // TODO: replace this implementation"},
			cursorLine:   0,
			cursorCol:    21,
			input:        "C",
			wantActive:   true,
			wantStartCol: 21,
			wantEndCol:   60,
			description:  "C should select from cursor to end for change",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := Position{Line: tt.cursorLine, Col: tt.cursorCol}
			result := SimulateMotionsWithSelection(start, tt.code, tt.input)

			// Log debug info
			line := tt.code[tt.cursorLine]
			t.Logf("Exercise: %s", tt.exerciseID)
			t.Logf("Line: %q", line)
			t.Logf("Cursor at col %d: '%c'", tt.cursorCol, line[tt.cursorCol])
			t.Logf("Input: %q", tt.input)
			t.Logf("Selection Active: %v", result.Selection.Active)
			t.Logf("Selection: [%d, %d]", result.Selection.StartCol, result.Selection.EndCol)
			if result.Selection.Active && result.Selection.EndCol < len(line) && result.Selection.StartCol >= 0 {
				selectedText := line[result.Selection.StartCol : result.Selection.EndCol+1]
				t.Logf("Selected text: %q", selectedText)
			}

			// Verify
			if result.Selection.Active != tt.wantActive {
				t.Errorf("Active = %v, want %v", result.Selection.Active, tt.wantActive)
			}
			if tt.wantActive {
				if result.Selection.StartCol != tt.wantStartCol {
					t.Errorf("StartCol = %d, want %d (%s)", result.Selection.StartCol, tt.wantStartCol, tt.description)
				}
				if result.Selection.EndCol != tt.wantEndCol {
					t.Errorf("EndCol = %d, want %d (%s)", result.Selection.EndCol, tt.wantEndCol, tt.description)
				}
			}
		})
	}
}

// TestTextObjectsExercisesVisualSelection tests ALL Text Objects exercises
func TestTextObjectsExercisesVisualSelection(t *testing.T) {
	tests := []struct {
		name         string
		exerciseID   string
		code         []string
		cursorCol    int
		cursorLine   int
		input        string
		wantActive   bool
		wantStartCol int
		wantEndCol   int
		description  string
	}{
		// textobjects_001: iw on userName
		{
			name:         "iw selects inner word",
			exerciseID:   "textobjects_001",
			code:         []string{"const userName = 'developer';"},
			cursorLine:   0,
			cursorCol:    6,
			input:        "iw",
			wantActive:   true,
			wantStartCol: 6,
			wantEndCol:   13, // userName
			description:  "iw should select 'userName'",
		},
		// textobjects_002: iw on calculateTotal
		{
			name:         "iw on camelCase word",
			exerciseID:   "textobjects_002",
			code:         []string{"function calculateTotal(items, taxRate) {"},
			cursorLine:   0,
			cursorCol:    12, // middle of calculateTotal
			input:        "iw",
			wantActive:   true,
			wantStartCol: 9,
			wantEndCol:   22, // calculateTotal
			description:  "iw should select entire 'calculateTotal'",
		},
		// textobjects_003: aw on value (includes trailing space)
		{
			name:         "aw includes trailing space",
			exerciseID:   "textobjects_003",
			code:         []string{"const result = value + offset;"},
			cursorLine:   0,
			cursorCol:    17, // somewhere in 'value'
			input:        "aw",
			wantActive:   true,
			wantStartCol: 15,
			wantEndCol:   20, // 'value ' with trailing space
			description:  "aw should select 'value ' with trailing space",
		},
		// textobjects_004: i" on Hello, World!
		{
			name:         `i" selects inside quotes`,
			exerciseID:   "textobjects_004",
			code:         []string{`const message = "Hello, World!";`},
			cursorLine:   0,
			cursorCol:    20,
			input:        `i"`,
			wantActive:   true,
			wantStartCol: 17,
			wantEndCol:   29, // Hello, World!
			description:  `i" should select content inside double quotes`,
		},
		// textobjects_005: a" includes quotes
		{
			name:         `a" includes quotes`,
			exerciseID:   "textobjects_005",
			code:         []string{`console.log("Processing data...");`},
			cursorLine:   0,
			cursorCol:    18,
			input:        `a"`,
			wantActive:   true,
			wantStartCol: 12,
			wantEndCol:   31, // "Processing data..."
			description:  `a" should select quoted string including quotes`,
		},
		// textobjects_006: i' on single quoted string
		{
			name:         "i' selects inside single quotes",
			exerciseID:   "textobjects_006",
			code:         []string{`const apiKey = 'sk-abc123xyz';`},
			cursorLine:   0,
			cursorCol:    20,
			input:        "i'",
			wantActive:   true,
			wantStartCol: 16,
			wantEndCol:   27, // sk-abc123xyz
			description:  "i' should select content inside single quotes",
		},
		// textobjects_007: a' includes single quotes
		{
			name:         "a' includes single quotes",
			exerciseID:   "textobjects_007",
			code:         []string{`import { useState } from 'react';`},
			cursorLine:   0,
			cursorCol:    28,
			input:        "a'",
			wantActive:   true,
			wantStartCol: 25,
			wantEndCol:   31, // 'react'
			description:  "a' should select 'react' including quotes",
		},
		// textobjects_008: i( selects function args
		{
			name:         "i( selects inside parentheses",
			exerciseID:   "textobjects_008",
			code:         []string{"const sum = add(firstNum, secondNum);"},
			cursorLine:   0,
			cursorCol:    20,
			input:        "i(",
			wantActive:   true,
			wantStartCol: 16,
			wantEndCol:   34, // firstNum, secondNum
			description:  "i( should select function arguments",
		},
		// textobjects_009: a( includes parentheses
		{
			name:         "a( includes parentheses",
			exerciseID:   "textobjects_009",
			code:         []string{"const result = compute(x * y + z);"},
			cursorLine:   0,
			cursorCol:    25,
			input:        "a(",
			wantActive:   true,
			wantStartCol: 22,
			wantEndCol:   32, // (x * y + z)
			description:  "a( should select parens and content",
		},
		// textobjects_010: i{ selects inside braces
		{
			name:         "i{ selects inside braces",
			exerciseID:   "textobjects_010",
			code:         []string{"const config = { debug: true, timeout: 5000 };"},
			cursorLine:   0,
			cursorCol:    25,
			input:        "i{",
			wantActive:   true,
			wantStartCol: 16, // after '{'
			wantEndCol:   43, // before '}'
			description:  "i{ should select content inside braces (including inner whitespace)",
		},
		// textobjects_011: a{ includes braces
		{
			name:         "a{ includes braces",
			exerciseID:   "textobjects_011",
			code:         []string{"const user = { name: 'John', age: 30 };"},
			cursorLine:   0,
			cursorCol:    20,
			input:        "a{",
			wantActive:   true,
			wantStartCol: 13,
			wantEndCol:   37, // { name: 'John', age: 30 }
			description:  "a{ should select braces and content",
		},
		// textobjects_012: i[ selects inside brackets
		{
			name:         "i[ selects inside brackets",
			exerciseID:   "textobjects_012",
			code:         []string{"const numbers = [1, 2, 3, 4, 5];"},
			cursorLine:   0,
			cursorCol:    20,
			input:        "i[",
			wantActive:   true,
			wantStartCol: 17,
			wantEndCol:   29, // 1, 2, 3, 4, 5
			description:  "i[ should select array elements",
		},
		// textobjects_013: a[ includes brackets
		{
			name:         "a[ includes brackets",
			exerciseID:   "textobjects_013",
			code:         []string{"const item = inventory[selectedIndex];"},
			cursorLine:   0,
			cursorCol:    27,
			input:        "a[",
			wantActive:   true,
			wantStartCol: 22,
			wantEndCol:   36, // [selectedIndex]
			description:  "a[ should select brackets and content",
		},
		// textobjects_014: it selects inside tags
		{
			name:         "it selects inside tags",
			exerciseID:   "textobjects_014",
			code:         []string{"<button>Click me</button>"},
			cursorLine:   0,
			cursorCol:    12,
			input:        "it",
			wantActive:   true,
			wantStartCol: 8,
			wantEndCol:   15, // Click me
			description:  "it should select tag content",
		},
		// textobjects_015: at selects entire tag
		{
			name:         "at selects entire tag",
			exerciseID:   "textobjects_015",
			code:         []string{"<div><span>Hello</span></div>"},
			cursorLine:   0,
			cursorCol:    12,
			input:        "at",
			wantActive:   true,
			wantStartCol: 5,
			wantEndCol:   22, // <span>Hello</span>
			description:  "at should select entire span element",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := Position{Line: tt.cursorLine, Col: tt.cursorCol}
			result := SimulateMotionsWithSelection(start, tt.code, tt.input)

			line := tt.code[tt.cursorLine]
			t.Logf("Exercise: %s", tt.exerciseID)
			t.Logf("Line: %q", line)
			if tt.cursorCol < len(line) {
				t.Logf("Cursor at col %d: '%c'", tt.cursorCol, line[tt.cursorCol])
			}
			t.Logf("Input: %q", tt.input)
			t.Logf("Selection Active: %v", result.Selection.Active)
			t.Logf("Selection: [%d, %d]", result.Selection.StartCol, result.Selection.EndCol)
			if result.Selection.Active && result.Selection.EndCol < len(line) && result.Selection.StartCol >= 0 {
				selectedText := line[result.Selection.StartCol : result.Selection.EndCol+1]
				t.Logf("Selected text: %q", selectedText)
			}

			if result.Selection.Active != tt.wantActive {
				t.Errorf("Active = %v, want %v", result.Selection.Active, tt.wantActive)
			}
			if tt.wantActive {
				if result.Selection.StartCol != tt.wantStartCol {
					t.Errorf("StartCol = %d, want %d (%s)", result.Selection.StartCol, tt.wantStartCol, tt.description)
				}
				if result.Selection.EndCol != tt.wantEndCol {
					t.Errorf("EndCol = %d, want %d (%s)", result.Selection.EndCol, tt.wantEndCol, tt.description)
				}
			}
		})
	}
}

// Helper to print character positions for debugging
func printLinePositions(line string) {
	fmt.Println("Line:", line)
	fmt.Println("Positions:")
	for i, ch := range line {
		fmt.Printf("%2d: '%c'\n", i, ch)
	}
}

// TestSubstitutionExercisesVisualSelection tests visual selection for gu/gU commands
func TestSubstitutionExercisesVisualSelection(t *testing.T) {
	tests := []struct {
		name         string
		exerciseID   string
		code         []string
		cursorCol    int
		cursorLine   int
		input        string
		wantActive   bool
		wantStartCol int
		wantEndCol   int
		description  string
	}{
		// substitution_011: gu + motion - lowercase
		{
			name:         "guiW lowercases WORD",
			exerciseID:   "substitution_011",
			code:         []string{"const ERROR_MESSAGE = 'Not Found';"},
			cursorLine:   0,
			cursorCol:    6,
			input:        "guiW",
			wantActive:   true,
			wantStartCol: 6,
			wantEndCol:   18, // ERROR_MESSAGE
			description:  "guiW should select ERROR_MESSAGE for lowercase",
		},
		// substitution_012: guu - lowercase line
		{
			name:         "guu lowercases entire line",
			exerciseID:   "substitution_012",
			code:         []string{"const STATUS = 'ACTIVE';"},
			cursorLine:   0,
			cursorCol:    0,
			input:        "guu",
			wantActive:   true,
			wantStartCol: 0,
			wantEndCol:   23,
			description:  "guu should select entire line for lowercase",
		},
		// substitution_013: gU + motion - uppercase
		{
			name:         "gUiW uppercases WORD",
			exerciseID:   "substitution_013",
			code:         []string{"const api_key = 'secret123';"},
			cursorLine:   0,
			cursorCol:    6,
			input:        "gUiW",
			wantActive:   true,
			wantStartCol: 6,
			wantEndCol:   12, // api_key
			description:  "gUiW should select api_key for uppercase",
		},
		// substitution_014: gUiw - uppercase word
		{
			name:         "gUiw uppercases word",
			exerciseID:   "substitution_014",
			code:         []string{"function getUser() { return user; }"},
			cursorLine:   0,
			cursorCol:    9,
			input:        "gUiw",
			wantActive:   true,
			wantStartCol: 9,
			wantEndCol:   15, // getUser
			description:  "gUiw should select getUser for uppercase",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := Position{Line: tt.cursorLine, Col: tt.cursorCol}
			result := SimulateMotionsWithSelection(start, tt.code, tt.input)

			line := tt.code[tt.cursorLine]
			t.Logf("Exercise: %s", tt.exerciseID)
			t.Logf("Line: %q", line)
			if tt.cursorCol < len(line) {
				t.Logf("Cursor at col %d: '%c'", tt.cursorCol, line[tt.cursorCol])
			}
			t.Logf("Input: %q", tt.input)
			t.Logf("Selection Active: %v", result.Selection.Active)
			t.Logf("Selection: [%d, %d]", result.Selection.StartCol, result.Selection.EndCol)
			if result.Selection.Active && result.Selection.EndCol < len(line) && result.Selection.StartCol >= 0 {
				selectedText := line[result.Selection.StartCol : result.Selection.EndCol+1]
				t.Logf("Selected text: %q", selectedText)
			}

			if result.Selection.Active != tt.wantActive {
				t.Errorf("Active = %v, want %v", result.Selection.Active, tt.wantActive)
			}
			if tt.wantActive {
				if result.Selection.StartCol != tt.wantStartCol {
					t.Errorf("StartCol = %d, want %d (%s)", result.Selection.StartCol, tt.wantStartCol, tt.description)
				}
				if result.Selection.EndCol != tt.wantEndCol {
					t.Errorf("EndCol = %d, want %d (%s)", result.Selection.EndCol, tt.wantEndCol, tt.description)
				}
			}
		})
	}
}

// TestMacrosExercisesVisualSelection tests visual selection for yank/paste with text objects
func TestMacrosExercisesVisualSelection(t *testing.T) {
	tests := []struct {
		name         string
		exerciseID   string
		code         []string
		cursorCol    int
		cursorLine   int
		input        string
		wantActive   bool
		wantStartCol int
		wantEndCol   int
		description  string
	}{
		// macros_011: "ayiw - yank inner word to register
		{
			name:         "yiw selects inner word for yank",
			exerciseID:   "macros_011",
			code:         []string{"const importantValue = 'remember_this';"},
			cursorLine:   0,
			cursorCol:    6,
			input:        "yiw", // We test without "a prefix since that's register, the selection is same
			wantActive:   true,
			wantStartCol: 6,
			wantEndCol:   19, // importantValue
			description:  "yiw should select importantValue",
		},
		// macros_013: "+yi' - yank inside quotes to clipboard
		{
			name:         "yi' selects inside quotes",
			exerciseID:   "macros_013",
			code:         []string{"const API_KEY = 'abc123xyz789';"},
			cursorLine:   0,
			cursorCol:    20,
			input:        "yi'",
			wantActive:   true,
			wantStartCol: 17,
			wantEndCol:   28, // abc123xyz789
			description:  "yi' should select content inside quotes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := Position{Line: tt.cursorLine, Col: tt.cursorCol}
			result := SimulateMotionsWithSelection(start, tt.code, tt.input)

			line := tt.code[tt.cursorLine]
			t.Logf("Exercise: %s", tt.exerciseID)
			t.Logf("Line: %q", line)
			if tt.cursorCol < len(line) {
				t.Logf("Cursor at col %d: '%c'", tt.cursorCol, line[tt.cursorCol])
			}
			t.Logf("Input: %q", tt.input)
			t.Logf("Selection Active: %v", result.Selection.Active)
			t.Logf("Selection: [%d, %d]", result.Selection.StartCol, result.Selection.EndCol)
			if result.Selection.Active && result.Selection.EndCol < len(line) && result.Selection.StartCol >= 0 {
				selectedText := line[result.Selection.StartCol : result.Selection.EndCol+1]
				t.Logf("Selected text: %q", selectedText)
			}

			if result.Selection.Active != tt.wantActive {
				t.Errorf("Active = %v, want %v", result.Selection.Active, tt.wantActive)
			}
			if tt.wantActive {
				if result.Selection.StartCol != tt.wantStartCol {
					t.Errorf("StartCol = %d, want %d (%s)", result.Selection.StartCol, tt.wantStartCol, tt.description)
				}
				if result.Selection.EndCol != tt.wantEndCol {
					t.Errorf("EndCol = %d, want %d (%s)", result.Selection.EndCol, tt.wantEndCol, tt.description)
				}
			}
		})
	}
}
