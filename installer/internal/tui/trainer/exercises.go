package trainer

// GetLessons returns lesson exercises for a module
func GetLessons(module ModuleID) []Exercise {
	switch module {
	case ModuleHorizontal:
		return getHorizontalLessons()
	case ModuleVertical:
		return getVerticalLessons()
	case ModuleTextObjects:
		return getTextObjectsLessons()
	case ModuleChangeRepeat:
		return getChangeRepeatLessons()
	case ModuleSubstitution:
		return getSubstitutionLessons()
	case ModuleRegex:
		return getRegexLessons()
	case ModuleMacros:
		return getMacrosLessons()
	default:
		return []Exercise{}
	}
}

// GetPracticeExercises returns practice exercises for a module
func GetPracticeExercises(module ModuleID) []Exercise {
	switch module {
	case ModuleHorizontal:
		return getHorizontalPractice()
	case ModuleVertical:
		return getVerticalPractice()
	case ModuleTextObjects:
		return getTextObjectsPractice()
	case ModuleChangeRepeat:
		return getChangeRepeatPractice()
	case ModuleSubstitution:
		return getSubstitutionPractice()
	case ModuleRegex:
		return getRegexPractice()
	case ModuleMacros:
		return getMacrosPractice()
	default:
		return []Exercise{}
	}
}

// GetBoss returns the boss fight for a module
func GetBoss(module ModuleID) *BossExercise {
	switch module {
	case ModuleHorizontal:
		return getHorizontalBoss()
	case ModuleVertical:
		return getVerticalBoss()
	case ModuleTextObjects:
		return getTextObjectsBoss()
	case ModuleChangeRepeat:
		return getChangeRepeatBoss()
	case ModuleSubstitution:
		return getSubstitutionBoss()
	case ModuleRegex:
		return getRegexBoss()
	case ModuleMacros:
		return getMacrosBoss()
	default:
		return nil
	}
}

// =============================================================================
// HORIZONTAL MODULE EXERCISES
// =============================================================================

func getHorizontalLessons() []Exercise {
	return []Exercise{
		// Lesson 1: w - word forward
		{
			ID:          "horizontal_001",
			Module:      ModuleHorizontal,
			Level:       1,
			Type:        ExerciseLesson,
			Code:        []string{"const userName = 'gentleman';"},
			CursorPos:   Position{Line: 0, Col: 0},
			Mission:     "Move to the start of 'userName' using w (word)",
			Solutions:   []string{"w"},
			Optimal:     "w",
			Hint:        "w moves to the start of the next word",
			Explanation: "w (word) moves the cursor to the start of the next word. It's the most basic and useful horizontal motion in Vim.",
			TimeoutSecs: 30,
			Points:      10,
		},
		// Lesson 2: W - WORD forward
		{
			ID:          "horizontal_002",
			Module:      ModuleHorizontal,
			Level:       1,
			Type:        ExerciseLesson,
			Code:        []string{"user.data.name = 'value';"},
			CursorPos:   Position{Line: 0, Col: 0},
			Mission:     "Move to '=' using W (WORD, skips punctuation)",
			Solutions:   []string{"W"},
			Optimal:     "W",
			Hint:        "W moves to the start of the next space-separated block",
			Explanation: "W (WORD) treats everything separated by spaces as a single word. 'user.data.name' is one WORD, but multiple words with w.",
			TimeoutSecs: 30,
			Points:      10,
		},
		// Lesson 3: e - end of word
		{
			ID:          "horizontal_003",
			Module:      ModuleHorizontal,
			Level:       1,
			Type:        ExerciseLesson,
			Code:        []string{"const config = { name: 'app' };"},
			CursorPos:   Position{Line: 0, Col: 0},
			Mission:     "Move to the end of 'const' using e (end)",
			Solutions:   []string{"e"},
			Optimal:     "e",
			Hint:        "e moves to the end of the current or next word",
			Explanation: "e (end) moves the cursor to the last character of the current or next word. Very useful for positioning before adding text.",
			TimeoutSecs: 30,
			Points:      10,
		},
		// Lesson 4: b - back word
		{
			ID:          "horizontal_004",
			Module:      ModuleHorizontal,
			Level:       1,
			Type:        ExerciseLesson,
			Code:        []string{"const userName = 'gentleman';"},
			CursorPos:   Position{Line: 0, Col: 15},
			Mission:     "Go back to the start of 'userName' using b (back)",
			Solutions:   []string{"b"},
			Optimal:     "b",
			Hint:        "b moves to the start of the previous word",
			Explanation: "b (back) is the opposite of w. It moves to the start of the previous word. Essential for navigating backwards efficiently.",
			TimeoutSecs: 30,
			Points:      10,
		},
		// Lesson 5: $ - end of line
		{
			ID:          "horizontal_005",
			Module:      ModuleHorizontal,
			Level:       2,
			Type:        ExerciseLesson,
			Code:        []string{"    const config = getConfig();"},
			CursorPos:   Position{Line: 0, Col: 15},
			Mission:     "Move to the end of the line using $ (end of line)",
			Solutions:   []string{"$"},
			Optimal:     "$",
			Hint:        "$ moves to the last character of the line",
			Explanation: "$ moves to the end of the line. It's one of the most used motions along with 0 (start of line) and ^ (first non-blank).",
			TimeoutSecs: 30,
			Points:      10,
		},
		// Lesson 6: 0 - start of line
		{
			ID:          "horizontal_006",
			Module:      ModuleHorizontal,
			Level:       2,
			Type:        ExerciseLesson,
			Code:        []string{"    const config = getConfig();"},
			CursorPos:   Position{Line: 0, Col: 25},
			Mission:     "Move to the absolute start of the line using 0 (zero)",
			Solutions:   []string{"0"},
			Optimal:     "0",
			Hint:        "0 moves to column 0, the absolute start",
			Explanation: "0 (zero) moves to the absolute start of the line, including spaces. Different from ^ which goes to the first non-blank character.",
			TimeoutSecs: 30,
			Points:      10,
		},
		// Lesson 7: ^ - first non-blank
		{
			ID:          "horizontal_007",
			Module:      ModuleHorizontal,
			Level:       2,
			Type:        ExerciseLesson,
			Code:        []string{"    const config = getConfig();"},
			CursorPos:   Position{Line: 0, Col: 25},
			Mission:     "Move to the first non-blank character using ^ (caret)",
			Solutions:   []string{"^"},
			Optimal:     "^",
			Hint:        "^ moves to the first character that isn't a space",
			Explanation: "^ moves to the first non-blank character of the line. Very useful in indented code where 0 would take you to the spaces.",
			TimeoutSecs: 30,
			Points:      10,
		},
		// Lesson 8: f{char} - find character
		{
			ID:          "horizontal_008",
			Module:      ModuleHorizontal,
			Level:       3,
			Type:        ExerciseLesson,
			Code:        []string{"const userName = 'gentleman';"},
			CursorPos:   Position{Line: 0, Col: 0},
			Mission:     "Move to the '=' using f (find)",
			Solutions:   []string{"f="},
			Optimal:     "f=",
			Hint:        "f followed by a character takes you to that character",
			Explanation: "f{char} (find) moves the cursor to the next occurrence of the specified character (inclusive). It's one of the most powerful commands.",
			TimeoutSecs: 30,
			Points:      15,
		},
		// Lesson 9: t{char} - till character
		{
			ID:          "horizontal_009",
			Module:      ModuleHorizontal,
			Level:       3,
			Type:        ExerciseLesson,
			Code:        []string{"const userName = 'gentleman';"},
			CursorPos:   Position{Line: 0, Col: 0},
			Mission:     "Move UNTIL BEFORE the '=' using t (till)",
			Solutions:   []string{"t="},
			Optimal:     "t=",
			Hint:        "t is like f but stops ONE character before",
			Explanation: "t{char} (till) is like f but exclusive - it stops just before the target character. Useful for operations like ct= (change till =).",
			TimeoutSecs: 30,
			Points:      15,
		},
		// Lesson 10: F{char} - find backwards
		{
			ID:          "horizontal_010",
			Module:      ModuleHorizontal,
			Level:       3,
			Type:        ExerciseLesson,
			Code:        []string{"const userName = 'gentleman';"},
			CursorPos:   Position{Line: 0, Col: 27},
			Mission:     "Move backwards to the first quote using F'",
			Solutions:   []string{"F'"},
			Optimal:     "F'",
			Hint:        "F is like f but searches backwards (left)",
			Explanation: "F{char} searches for the character backwards (left). T{char} does the same but exclusive. They're the mirror of f and t.",
			TimeoutSecs: 30,
			Points:      15,
		},
		// Lesson 11: ; - repeat f/F/t/T
		{
			ID:          "horizontal_011",
			Module:      ModuleHorizontal,
			Level:       4,
			Type:        ExerciseLesson,
			Code:        []string{"name.first.middle.last"},
			CursorPos:   Position{Line: 0, Col: 0},
			Mission:     "Reach the third dot using f. then ;; (repeat twice)",
			Solutions:   []string{"f.;;", "3f."},
			Optimal:     "3f.",
			Hint:        "; repeats the last f/F/t/T in the same direction",
			Explanation: "; repeats the last f/F/t/T movement. You can also use {n}f. to go to the n-th directly. Both are valid.",
			TimeoutSecs: 30,
			Points:      20,
		},
		// Lesson 12: , - repeat opposite direction
		{
			ID:          "horizontal_012",
			Module:      ModuleHorizontal,
			Level:       4,
			Type:        ExerciseLesson,
			Code:        []string{"name.first.middle.last"},
			CursorPos:   Position{Line: 0, Col: 17},
			Mission:     "From the third dot, go back to the second using F. (find backward)",
			Solutions:   []string{"F."},
			Optimal:     "F.",
			Hint:        "F searches backward for a character, , repeats the last f/F in opposite direction",
			Explanation: "F{char} searches backward. You can also use , to repeat f/F/t/T in the opposite direction if you already did a search.",
			TimeoutSecs: 30,
			Points:      20,
		},
		// Lesson 13: ge - end of previous word
		{
			ID:          "horizontal_013",
			Module:      ModuleHorizontal,
			Level:       4,
			Type:        ExerciseLesson,
			Code:        []string{"const userName = 'gentleman';"},
			CursorPos:   Position{Line: 0, Col: 15},
			Mission:     "Move to the end of the previous word (the 'e' in 'userName') using ge",
			Solutions:   []string{"ge"},
			Optimal:     "ge",
			Hint:        "ge is like e but backwards",
			Explanation: "ge moves to the end of the previous word. It's the opposite of e, just like b is the opposite of w. gE does the same with WORDS.",
			TimeoutSecs: 30,
			Points:      20,
		},
		// Lesson 14: Numbered motions
		{
			ID:          "horizontal_014",
			Module:      ModuleHorizontal,
			Level:       5,
			Type:        ExerciseLesson,
			Code:        []string{"one two three four five six"},
			CursorPos:   Position{Line: 0, Col: 0},
			Mission:     "Reach 'four' using 3w (three words forward)",
			Solutions:   []string{"3w", "ff"},
			Optimal:     "3w",
			Hint:        "You can put a number before any motion to repeat it",
			Explanation: "{n}w jumps n words forward. 3w = www. This works with any motion: 3e, 3b, 3f{char}, etc.",
			TimeoutSecs: 30,
			Points:      25,
		},
		// Lesson 15: Combining for efficiency
		{
			ID:          "horizontal_015",
			Module:      ModuleHorizontal,
			Level:       5,
			Type:        ExerciseLesson,
			Code:        []string{"const getUserData = (userId) => fetchUser(userId);"},
			CursorPos:   Position{Line: 0, Col: 0},
			Mission:     "Reach the opening parenthesis of fetchUser efficiently",
			Solutions:   []string{"2f(", "f(;"},
			Optimal:     "2f(",
			Hint:        "There are two parentheses - how do you get to the second directly?",
			Explanation: "2f( jumps to the second '('. You can also use f( then ; to repeat. Mastering {n}f{char} is the key to being fast in Vim.",
			TimeoutSecs: 30,
			Points:      30,
		},
		// Lesson 16: E - end of WORD (skips punctuation)
		{
			ID:          "horizontal_016",
			Module:      ModuleHorizontal,
			Level:       5,
			Type:        ExerciseLesson,
			Code:        []string{"user.profile.settings = config.defaults;"},
			CursorPos:   Position{Line: 0, Col: 0},
			Mission:     "Move to the end of 'user.profile.settings' using E (end of WORD)",
			Solutions:   []string{"E"},
			Optimal:     "E",
			Hint:        "E moves to the end of the current WORD (space-separated)",
			Explanation: "E (end of WORD) moves to the end of the current space-separated block. While 'e' stops at 'user', 'E' goes all the way to 'settings'. Use E when dealing with dotted paths or hyphenated words.",
			TimeoutSecs: 30,
			Points:      15,
		},
		// Lesson 17: B - back WORD (skips punctuation)
		{
			ID:          "horizontal_017",
			Module:      ModuleHorizontal,
			Level:       5,
			Type:        ExerciseLesson,
			Code:        []string{"const result = user.profile.settings;"},
			CursorPos:   Position{Line: 0, Col: 36},
			Mission:     "Move back to 'user' using B (back WORD)",
			Solutions:   []string{"B"},
			Optimal:     "B",
			Hint:        "B moves to the start of the previous WORD (space-separated)",
			Explanation: "B (back WORD) is the opposite of W. It jumps back over entire space-separated blocks. From 'settings', B takes you to 'user' in one move, while 'b' would stop at each dot-separated part.",
			TimeoutSecs: 30,
			Points:      15,
		},
		// Lesson 18: T{char} - till backwards
		{
			ID:           "horizontal_018",
			Module:       ModuleHorizontal,
			Level:        6,
			Type:         ExerciseLesson,
			Code:         []string{"const value = 'hello world';"},
			CursorPos:    Position{Line: 0, Col: 26},
			CursorTarget: &Position{Line: 0, Col: 15},
			Mission:      "Move backwards to just AFTER the opening quote using T'",
			Solutions:    []string{"T'"},
			Optimal:      "T'",
			Hint:         "T is like F but stops one character AFTER the target",
			Explanation:  "T{char} (till backwards) searches backward but lands one character AFTER the target. F' would land ON the quote, T' lands just after it. Perfect for operations like dT' (delete backwards till quote).",
			TimeoutSecs:  30,
			Points:       20,
		},
		// Lesson 19: , - repeat f/F/t/T in opposite direction
		{
			ID:           "horizontal_019",
			Module:       ModuleHorizontal,
			Level:        6,
			Type:         ExerciseLesson,
			Code:         []string{"one.two.three.four.five"},
			CursorPos:    Position{Line: 0, Col: 0},
			CursorTarget: &Position{Line: 0, Col: 3},
			Mission:      "Use f. to jump to a dot, then ; to go forward, then , to go back",
			Solutions:    []string{"f.;,"},
			Optimal:      "f.;,",
			Hint:         ", repeats the last f/F/t/T but in the OPPOSITE direction",
			Explanation:  ", (comma) is the reverse of ; (semicolon). If you used f. to search forward, , goes backward to the previous dot. It's like having an undo for your character search - incredibly useful for overshooting!",
			TimeoutSecs:  30,
			Points:       20,
		},
	}
}

func getHorizontalPractice() []Exercise {
	lessons := getHorizontalLessons()
	practice := make([]Exercise, len(lessons))

	// Convert lessons to practice type
	for i, ex := range lessons {
		practice[i] = ex
		practice[i].Type = ExercisePractice
		practice[i].ID = "horizontal_p" + ex.ID[len("horizontal_"):]
		practice[i].TimeoutSecs = 15 // Shorter timeout for practice
	}

	return practice
}

func getHorizontalBoss() *BossExercise {
	return &BossExercise{
		ID:        "horizontal_boss",
		Module:    ModuleHorizontal,
		Name:      "The Line Walker",
		Lives:     3,
		BonusTime: 30,
		Steps: []BossStep{
			{
				TimeLimit: 5,
				Exercise: Exercise{
					ID:        "horizontal_boss_1",
					Module:    ModuleHorizontal,
					Level:     5,
					Type:      ExerciseBoss,
					Code:      []string{"const getUserProfile = (userId) => fetchUserData(userId).then(processProfile);"},
					CursorPos: Position{Line: 0, Col: 0},
					Mission:   "Move to the 'e' in 'getUser'",
					Solutions: []string{"fe", "wl"},
					Optimal:   "fe",
					Points:    50,
				},
			},
			{
				TimeLimit: 5,
				Exercise: Exercise{
					ID:        "horizontal_boss_2",
					Module:    ModuleHorizontal,
					Level:     5,
					Type:      ExerciseBoss,
					Code:      []string{"const getUserProfile = (userId) => fetchUserData(userId).then(processProfile);"},
					CursorPos: Position{Line: 0, Col: 33},
					Mission:   "Move to the start of 'fetchUserData'",
					Solutions: []string{"ff", "ll"},
					Optimal:   "ff",
					Points:    50,
				},
			},
			{
				TimeLimit: 5,
				Exercise: Exercise{
					ID:        "horizontal_boss_3",
					Module:    ModuleHorizontal,
					Level:     5,
					Type:      ExerciseBoss,
					Code:      []string{"const getUserProfile = (userId) => fetchUserData(userId).then(processProfile);"},
					CursorPos: Position{Line: 0, Col: 36},
					Mission:   "Move to 'processProfile'",
					Solutions: []string{"fp", "2f(l"},
					Optimal:   "fp",
					Points:    50,
				},
			},
			{
				TimeLimit: 5,
				Exercise: Exercise{
					ID:        "horizontal_boss_4",
					Module:    ModuleHorizontal,
					Level:     5,
					Type:      ExerciseBoss,
					Code:      []string{"const getUserProfile = (userId) => fetchUserData(userId).then(processProfile);"},
					CursorPos: Position{Line: 0, Col: 62},
					Mission:   "Move to the end of the line",
					Solutions: []string{"$"},
					Optimal:   "$",
					Points:    50,
				},
			},
			{
				TimeLimit: 5,
				Exercise: Exercise{
					ID:        "horizontal_boss_5",
					Module:    ModuleHorizontal,
					Level:     5,
					Type:      ExerciseBoss,
					Code:      []string{"const getUserProfile = (userId) => fetchUserData(userId).then(processProfile);"},
					CursorPos: Position{Line: 0, Col: 77},
					Mission:   "Go back to the start of the line",
					Solutions: []string{"0", "^"},
					Optimal:   "0",
					Points:    50,
				},
			},
		},
	}
}

// =============================================================================
// VERTICAL MODULE EXERCISES
// =============================================================================

func getVerticalLessons() []Exercise {
	return []Exercise{
		// Lesson 1: j - move down
		{
			ID:     "vertical_001",
			Module: ModuleVertical,
			Level:  1,
			Type:   ExerciseLesson,
			Code: []string{
				"function greet() {",
				"  console.log('Hello');",
				"  return true;",
				"}",
			},
			CursorPos:   Position{Line: 0, Col: 0},
			Mission:     "Move down to the console.log line using j",
			Solutions:   []string{"j"},
			Optimal:     "j",
			Hint:        "j moves the cursor down one line",
			Explanation: "j is the most basic downward motion. Think of j as having a downward hook at the bottom of the letter.",
			TimeoutSecs: 30,
			Points:      10,
		},
		// Lesson 2: k - move up
		{
			ID:     "vertical_002",
			Module: ModuleVertical,
			Level:  1,
			Type:   ExerciseLesson,
			Code: []string{
				"function greet() {",
				"  console.log('Hello');",
				"  return true;",
				"}",
			},
			CursorPos:   Position{Line: 2, Col: 2},
			Mission:     "Move up to the console.log line using k",
			Solutions:   []string{"k"},
			Optimal:     "k",
			Hint:        "k moves the cursor up one line",
			Explanation: "k moves up. Together j and k are your bread and butter for vertical navigation. Think of k as pointing upward.",
			TimeoutSecs: 30,
			Points:      10,
		},
		// Lesson 3: Multiple j/k with count
		{
			ID:     "vertical_003",
			Module: ModuleVertical,
			Level:  1,
			Type:   ExerciseLesson,
			Code: []string{
				"const a = 1;",
				"const b = 2;",
				"const c = 3;",
				"const d = 4;",
				"const e = 5;",
			},
			CursorPos:   Position{Line: 0, Col: 0},
			Mission:     "Move down 3 lines to 'const d' using 3j",
			Solutions:   []string{"3j", "jjj"},
			Optimal:     "3j",
			Hint:        "You can prefix any motion with a count: {n}j",
			Explanation: "{count}j moves down count lines. 3j = jjj but faster to type. This works with any motion!",
			TimeoutSecs: 30,
			Points:      15,
		},
		// Lesson 4: gg - go to top
		{
			ID:     "vertical_004",
			Module: ModuleVertical,
			Level:  2,
			Type:   ExerciseLesson,
			Code: []string{
				"// File: utils.js",
				"function helper() {",
				"  return 42;",
				"}",
				"",
				"export default helper;",
			},
			CursorPos:   Position{Line: 4, Col: 0},
			Mission:     "Jump to the first line of the file using gg",
			Solutions:   []string{"gg"},
			Optimal:     "gg",
			Hint:        "gg goes to the first line of the file",
			Explanation: "gg instantly jumps to the top of the file. Essential for navigating large files. You can also use [n]gg to go to line n.",
			TimeoutSecs: 30,
			Points:      15,
		},
		// Lesson 5: G - go to bottom
		{
			ID:     "vertical_005",
			Module: ModuleVertical,
			Level:  2,
			Type:   ExerciseLesson,
			Code: []string{
				"// File: utils.js",
				"function helper() {",
				"  return 42;",
				"}",
				"",
				"export default helper;",
			},
			CursorPos:   Position{Line: 0, Col: 0},
			Mission:     "Jump to the last line of the file using G",
			Solutions:   []string{"G"},
			Optimal:     "G",
			Hint:        "G (capital) goes to the last line",
			Explanation: "G jumps to the end of the file. Use [n]G to go to a specific line number. gg and G are the vertical equivalent of 0 and $.",
			TimeoutSecs: 30,
			Points:      15,
		},
		// Lesson 6: [n]G - go to specific line
		{
			ID:     "vertical_006",
			Module: ModuleVertical,
			Level:  2,
			Type:   ExerciseLesson,
			Code: []string{
				"line 1",
				"line 2",
				"line 3",
				"line 4 - target",
				"line 5",
				"line 6",
			},
			CursorPos:   Position{Line: 0, Col: 0},
			Mission:     "Jump directly to line 4 using 4G",
			Solutions:   []string{"4G", "3j"},
			Optimal:     "4G",
			Hint:        "[n]G goes to line number n",
			Explanation: "[n]G is incredibly useful when you know the line number (from errors, linters, etc). Much faster than counting j presses!",
			TimeoutSecs: 30,
			Points:      20,
		},
		// Lesson 7: [n]gg - alternative to go to line
		{
			ID:     "vertical_007",
			Module: ModuleVertical,
			Level:  2,
			Type:   ExerciseLesson,
			Code: []string{
				"line 1",
				"line 2",
				"line 3 - target",
				"line 4",
				"line 5",
			},
			CursorPos:   Position{Line: 4, Col: 0},
			Mission:     "Jump to line 3 using 3gg",
			Solutions:   []string{"3gg", "2k"},
			Optimal:     "3gg",
			Hint:        "[n]gg also goes to line n, same as [n]G",
			Explanation: "Both [n]gg and [n]G go to line n. Some prefer gg because it's easier to type. Use whichever feels natural.",
			TimeoutSecs: 30,
			Points:      20,
		},
		// Lesson 8: } - next paragraph
		{
			ID:     "vertical_008",
			Module: ModuleVertical,
			Level:  3,
			Type:   ExerciseLesson,
			Code: []string{
				"function first() {",
				"  return 1;",
				"}",
				"",
				"function second() {",
				"  return 2;",
				"}",
			},
			CursorPos:   Position{Line: 0, Col: 0},
			Mission:     "Jump to the blank line between functions using }",
			Solutions:   []string{"}"},
			Optimal:     "}",
			Hint:        "} moves to the next blank line (paragraph boundary)",
			Explanation: "} jumps forward to the next blank line. In code, this usually means jumping between functions or blocks. Super useful!",
			TimeoutSecs: 30,
			Points:      20,
		},
		// Lesson 9: { - previous paragraph
		{
			ID:     "vertical_009",
			Module: ModuleVertical,
			Level:  3,
			Type:   ExerciseLesson,
			Code: []string{
				"function first() {",
				"  return 1;",
				"}",
				"",
				"function second() {",
				"  return 2;",
				"}",
			},
			CursorPos:   Position{Line: 5, Col: 2},
			Mission:     "Jump back to the blank line using {",
			Solutions:   []string{"{"},
			Optimal:     "{",
			Hint:        "{ moves to the previous blank line",
			Explanation: "{ jumps backward to the previous blank line. Combined with }, you can quickly navigate between code blocks.",
			TimeoutSecs: 30,
			Points:      20,
		},
		// Lesson 10: + move to next line first non-blank
		{
			ID:     "vertical_010",
			Module: ModuleVertical,
			Level:  3,
			Type:   ExerciseLesson,
			Code: []string{
				"if (condition) {",
				"    doSomething();",
				"    doMore();",
				"}",
			},
			CursorPos:   Position{Line: 0, Col: 0},
			Mission:     "Move to the first non-blank character of the next line using +",
			Solutions:   []string{"+", "j^"},
			Optimal:     "+",
			Hint:        "+ moves down and to the first non-blank character",
			Explanation: "+ is like j followed by ^. It moves to the next line's first non-blank character. Great for navigating indented code.",
			TimeoutSecs: 30,
			Points:      20,
		},
		// Lesson 11: - move to previous line first non-blank
		{
			ID:     "vertical_011",
			Module: ModuleVertical,
			Level:  3,
			Type:   ExerciseLesson,
			Code: []string{
				"if (condition) {",
				"    doSomething();",
				"    doMore();",
				"}",
			},
			CursorPos:   Position{Line: 2, Col: 4},
			Mission:     "Move to the first non-blank character of the previous line using -",
			Solutions:   []string{"-", "k^"},
			Optimal:     "-",
			Hint:        "- moves up and to the first non-blank character",
			Explanation: "- is like k followed by ^. Moves to previous line's first non-blank. Together + and - help you navigate indented blocks.",
			TimeoutSecs: 30,
			Points:      20,
		},
		// Lesson 12: Combining j/k with counts efficiently
		{
			ID:     "vertical_012",
			Module: ModuleVertical,
			Level:  4,
			Type:   ExerciseLesson,
			Code: []string{
				"function App() {",
				"  const [count, setCount] = useState(0);",
				"",
				"  return (",
				"    <button onClick={() => setCount(c => c + 1)}>",
				"      Count: {count}",
				"    </button>",
				"  );",
				"}",
			},
			CursorPos:    Position{Line: 0, Col: 0},
			CursorTarget: &Position{Line: 3, Col: 0},
			Mission:      "Jump to the 'return' statement (line 4) using 3j",
			Solutions:    []string{"3j"},
			Optimal:      "3j",
			Hint:         "Count the lines or use relative line numbers if enabled",
			Explanation:  "In real editing, enable relative line numbers (:set relativenumber) to easily see counts. Then [n]j becomes natural.",
			TimeoutSecs:  30,
			Points:       25,
		},
		// Lesson 13: Multiple paragraphs with count
		{
			ID:     "vertical_013",
			Module: ModuleVertical,
			Level:  4,
			Type:   ExerciseLesson,
			Code: []string{
				"// Block 1",
				"const a = 1;",
				"",
				"// Block 2",
				"const b = 2;",
				"",
				"// Block 3",
				"const c = 3;",
			},
			CursorPos:   Position{Line: 0, Col: 0},
			Mission:     "Jump forward two paragraph boundaries using 2}",
			Solutions:   []string{"2}"},
			Optimal:     "2}",
			Hint:        "You can use counts with } and { too",
			Explanation: "2} jumps two blank lines forward. This is faster than }} and works great in code with consistent spacing.",
			TimeoutSecs: 30,
			Points:      25,
		},
		// Lesson 14: Combining vertical and horizontal
		{
			ID:     "vertical_014",
			Module: ModuleVertical,
			Level:  5,
			Type:   ExerciseLesson,
			Code: []string{
				"const config = {",
				"  name: 'app',",
				"  version: '1.0.0',",
				"  debug: true,",
				"};",
			},
			CursorPos:   Position{Line: 0, Col: 0},
			Mission:     "Navigate to 'debug' using 3j then w",
			Solutions:   []string{"3jw"},
			Optimal:     "3jw",
			Hint:        "You can chain any motions together",
			Explanation: "Vim motions compose naturally. 3jw = go down 3 lines, then next word. This is the power of modal editing!",
			TimeoutSecs: 30,
			Points:      30,
		},
		// Lesson 15: Quick file navigation pattern
		{
			ID:     "vertical_015",
			Module: ModuleVertical,
			Level:  5,
			Type:   ExerciseLesson,
			Code: []string{
				"export function api() {",
				"  // implementation",
				"}",
				"",
				"export function helper() {",
				"  // implementation",
				"}",
				"",
				"export function main() {",
				"  // TARGET: this line",
				"}",
			},
			CursorPos:   Position{Line: 0, Col: 0},
			Mission:     "Reach the 'TARGET' comment using G then k",
			Solutions:   []string{"Gk"},
			Optimal:     "Gk",
			Hint:        "Sometimes going to the end and moving back is faster",
			Explanation: "Real Vim mastery is knowing multiple ways and choosing the fastest. Gk = end of file, one line up. Sometimes that's fastest!",
			TimeoutSecs: 30,
			Points:      30,
		},
		// Lesson 16: H - High (top of visible screen)
		{
			ID:     "vertical_016",
			Module: ModuleVertical,
			Level:  4,
			Type:   ExerciseLesson,
			Code: []string{
				"// TOP OF SCREEN - TARGET",
				"function first() {}",
				"function second() {}",
				"function third() {}",
				"function fourth() {}",
				"function fifth() {}",
				"// You are here",
			},
			CursorPos:   Position{Line: 6, Col: 0},
			Mission:     "Jump to the top of the visible screen using H (High)",
			Solutions:   []string{"H"},
			Optimal:     "H",
			Hint:        "H takes you to the Highest line on screen",
			Explanation: "H (High) jumps to the top of the visible screen. Unlike gg which goes to the file's first line, H goes to the first VISIBLE line. Combined with M and L, you can navigate the screen in thirds.",
			TimeoutSecs: 30,
			Points:      20,
		},
		// Lesson 17: M - Middle of visible screen
		{
			ID:     "vertical_017",
			Module: ModuleVertical,
			Level:  4,
			Type:   ExerciseLesson,
			Code: []string{
				"function first() {}",
				"function second() {}",
				"function third() {}",
				"// MIDDLE OF SCREEN - TARGET",
				"function fourth() {}",
				"function fifth() {}",
				"function sixth() {}",
			},
			CursorPos:   Position{Line: 0, Col: 0},
			Mission:     "Jump to the middle of the visible screen using M (Middle)",
			Solutions:   []string{"M"},
			Optimal:     "M",
			Hint:        "M takes you to the Middle line on screen",
			Explanation: "M (Middle) jumps to the middle of the visible screen. It's perfect for quickly getting to the center of what you're viewing without counting lines.",
			TimeoutSecs: 30,
			Points:      20,
		},
		// Lesson 18: L - Low (bottom of visible screen)
		{
			ID:     "vertical_018",
			Module: ModuleVertical,
			Level:  4,
			Type:   ExerciseLesson,
			Code: []string{
				"// You are here",
				"function first() {}",
				"function second() {}",
				"function third() {}",
				"function fourth() {}",
				"function fifth() {}",
				"// BOTTOM OF SCREEN - TARGET",
			},
			CursorPos:   Position{Line: 0, Col: 0},
			Mission:     "Jump to the bottom of the visible screen using L (Low)",
			Solutions:   []string{"L"},
			Optimal:     "L",
			Hint:        "L takes you to the Lowest line on screen",
			Explanation: "L (Low) jumps to the bottom of the visible screen. H, M, L together let you navigate any visible area in at most 2 keystrokes: one to get to the right third, then fine-tune with j/k.",
			TimeoutSecs: 30,
			Points:      20,
		},
		// Lesson 19: Ctrl+d - half page down
		{
			ID:     "vertical_019",
			Module: ModuleVertical,
			Level:  5,
			Type:   ExerciseLesson,
			Code: []string{
				"// Line 1",
				"// Line 2",
				"// Line 3",
				"// Line 4",
				"// Line 5",
				"// Line 6",
				"// Line 7",
				"// Line 8",
				"// TARGET - halfway down",
				"// Line 10",
			},
			CursorPos:   Position{Line: 0, Col: 0},
			Mission:     "Scroll half a page down using Ctrl+d",
			Solutions:   []string{"\x04"},
			Optimal:     "\x04",
			Hint:        "Ctrl+d scrolls Down half a page",
			Explanation: "Ctrl+d (down) scrolls half a page down, keeping your cursor in the middle of the screen. It's smoother than jumping full pages and helps you maintain context while navigating.",
			TimeoutSecs: 30,
			Points:      25,
		},
		// Lesson 20: Ctrl+u - half page up
		{
			ID:     "vertical_020",
			Module: ModuleVertical,
			Level:  5,
			Type:   ExerciseLesson,
			Code: []string{
				"// TARGET - halfway up",
				"// Line 2",
				"// Line 3",
				"// Line 4",
				"// Line 5",
				"// Line 6",
				"// Line 7",
				"// Line 8",
				"// Line 9",
				"// You are here",
			},
			CursorPos:   Position{Line: 9, Col: 0},
			Mission:     "Scroll half a page up using Ctrl+u",
			Solutions:   []string{"\x15"},
			Optimal:     "\x15",
			Hint:        "Ctrl+u scrolls Up half a page",
			Explanation: "Ctrl+u (up) scrolls half a page up. Together with Ctrl+d, these are the most practical scrolling commands - they move enough to make progress but not so much that you lose your place.",
			TimeoutSecs: 30,
			Points:      25,
		},
		// Lesson 21: Ctrl+f - full page forward
		{
			ID:     "vertical_021",
			Module: ModuleVertical,
			Level:  5,
			Type:   ExerciseLesson,
			Code: []string{
				"// Page 1 - Line 1",
				"// Page 1 - Line 2",
				"// Page 1 - Line 3",
				"// Page 1 - Line 4",
				"// --- PAGE BREAK ---",
				"// Page 2 - Line 1",
				"// Page 2 - Line 2",
				"// Page 2 - TARGET",
				"// Page 2 - Line 4",
				"// Page 2 - Line 5",
			},
			CursorPos:   Position{Line: 0, Col: 0},
			Mission:     "Scroll a full page forward using Ctrl+f",
			Solutions:   []string{"\x06"},
			Optimal:     "\x06",
			Hint:        "Ctrl+f scrolls Forward a full page",
			Explanation: "Ctrl+f (forward) scrolls a full page down. Use it when you need to move quickly through a long file. The cursor stays at the same screen position while the content scrolls.",
			TimeoutSecs: 30,
			Points:      25,
		},
		// Lesson 22: Ctrl+b - full page backward
		{
			ID:     "vertical_022",
			Module: ModuleVertical,
			Level:  5,
			Type:   ExerciseLesson,
			Code: []string{
				"// Page 1 - TARGET",
				"// Page 1 - Line 2",
				"// Page 1 - Line 3",
				"// Page 1 - Line 4",
				"// --- PAGE BREAK ---",
				"// Page 2 - Line 1",
				"// Page 2 - Line 2",
				"// Page 2 - Line 3",
				"// Page 2 - Line 4",
				"// You are here",
			},
			CursorPos:   Position{Line: 9, Col: 0},
			Mission:     "Scroll a full page backward using Ctrl+b",
			Solutions:   []string{"\x02"},
			Optimal:     "\x02",
			Hint:        "Ctrl+b scrolls Backward a full page",
			Explanation: "Ctrl+b (backward) scrolls a full page up. Summary: Ctrl+d/u for half pages (precise), Ctrl+f/b for full pages (fast). Learn to mix them based on how far you need to go!",
			TimeoutSecs: 30,
			Points:      25,
		},
	}
}

func getVerticalPractice() []Exercise {
	lessons := getVerticalLessons()
	practice := make([]Exercise, len(lessons))

	for i, ex := range lessons {
		practice[i] = ex
		practice[i].Type = ExercisePractice
		practice[i].ID = "vertical_p" + ex.ID[len("vertical_"):]
		practice[i].TimeoutSecs = 15
	}

	return practice
}

func getVerticalBoss() *BossExercise {
	return &BossExercise{
		ID:        "vertical_boss",
		Module:    ModuleVertical,
		Name:      "The Code Climber",
		Lives:     3,
		BonusTime: 30,
		Steps: []BossStep{
			{
				TimeLimit: 5,
				Exercise: Exercise{
					ID:     "vertical_boss_1",
					Module: ModuleVertical,
					Level:  5,
					Type:   ExerciseBoss,
					Code: []string{
						"import express from 'express';",
						"import cors from 'cors';",
						"import helmet from 'helmet';",
						"",
						"const app = express();",
						"",
						"app.use(cors());",
						"app.use(helmet());",
						"",
						"app.listen(3000);",
					},
					CursorPos: Position{Line: 0, Col: 0},
					Mission:   "Jump to the blank line after imports",
					Solutions: []string{"}", "3j"},
					Optimal:   "}",
					Points:    50,
				},
			},
			{
				TimeLimit: 5,
				Exercise: Exercise{
					ID:     "vertical_boss_2",
					Module: ModuleVertical,
					Level:  5,
					Type:   ExerciseBoss,
					Code: []string{
						"import express from 'express';",
						"import cors from 'cors';",
						"import helmet from 'helmet';",
						"",
						"const app = express();",
						"",
						"app.use(cors());",
						"app.use(helmet());",
						"",
						"app.listen(3000);",
					},
					CursorPos: Position{Line: 3, Col: 0},
					Mission:   "Jump to the last line (app.listen)",
					Solutions: []string{"G"},
					Optimal:   "G",
					Points:    50,
				},
			},
			{
				TimeLimit: 5,
				Exercise: Exercise{
					ID:     "vertical_boss_3",
					Module: ModuleVertical,
					Level:  5,
					Type:   ExerciseBoss,
					Code: []string{
						"import express from 'express';",
						"import cors from 'cors';",
						"import helmet from 'helmet';",
						"",
						"const app = express();",
						"",
						"app.use(cors());",
						"app.use(helmet());",
						"",
						"app.listen(3000);",
					},
					CursorPos: Position{Line: 9, Col: 0},
					Mission:   "Jump back to the very first line",
					Solutions: []string{"gg"},
					Optimal:   "gg",
					Points:    50,
				},
			},
			{
				TimeLimit: 6,
				Exercise: Exercise{
					ID:     "vertical_boss_4",
					Module: ModuleVertical,
					Level:  5,
					Type:   ExerciseBoss,
					Code: []string{
						"import express from 'express';",
						"import cors from 'cors';",
						"import helmet from 'helmet';",
						"",
						"const app = express();",
						"",
						"app.use(cors());",
						"app.use(helmet());",
						"",
						"app.listen(3000);",
					},
					CursorPos: Position{Line: 0, Col: 0},
					Mission:   "Jump to line 7 (app.use(cors))",
					Solutions: []string{"7G", "6j", "2}"},
					Optimal:   "7G",
					Points:    50,
				},
			},
			{
				TimeLimit: 6,
				Exercise: Exercise{
					ID:     "vertical_boss_5",
					Module: ModuleVertical,
					Level:  5,
					Type:   ExerciseBoss,
					Code: []string{
						"import express from 'express';",
						"import cors from 'cors';",
						"import helmet from 'helmet';",
						"",
						"const app = express();",
						"",
						"app.use(cors());",
						"app.use(helmet());",
						"",
						"app.listen(3000);",
					},
					CursorPos: Position{Line: 6, Col: 0},
					Mission:   "Jump back to 'const app' line using {",
					Solutions: []string{"{k", "2k"},
					Optimal:   "{k",
					Points:    50,
				},
			},
		},
	}
}
