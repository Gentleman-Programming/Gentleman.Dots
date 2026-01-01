package trainer

// getSubstitutionLessons returns lessons for the Substitution module
// covering r, R, s, S, ~, gu, gU, and J commands
func getSubstitutionLessons() []Exercise {
	return []Exercise{
		// Lessons 1-2: r - replace single character
		{
			ID:          "substitution_001",
			Module:      ModuleSubstitution,
			Level:       1,
			Type:        ExerciseLesson,
			Code:        []string{"function getTaxRate() { return 0.o8; }"},
			CursorPos:   Position{Line: 0, Col: 33},
			Mission:     "Fix the typo: replace 'o' with '0' using r",
			Solutions:   []string{"r0"},
			Optimal:     "r0",
			Hint:        "r replaces single character without entering insert mode",
			Explanation: "The 'r' command replaces the character under the cursor with the next character you type. It's perfect for fixing single-character typos quickly without entering insert mode.",
			TimeoutSecs: 30,
			Points:      10,
		},
		{
			ID:          "substitution_002",
			Module:      ModuleSubstitution,
			Level:       1,
			Type:        ExerciseLesson,
			Code:        []string{"const API_ENDPOINT = 'htps://api.example.com';"},
			CursorPos:   Position{Line: 0, Col: 23},
			Mission:     "Fix 'htps' to 'https' by replacing 'p' with 't' using r, then move right and replace 's' with 'p'",
			Solutions:   []string{"rtlrp", "rt$F:hrp"},
			Optimal:     "rtlrp",
			Hint:        "Use r to replace, then l to move right, then r again",
			Explanation: "Multiple r commands can fix adjacent typos efficiently. After each r replacement, you stay in normal mode and can immediately move to the next character.",
			TimeoutSecs: 45,
			Points:      15,
		},
		// Lessons 3-4: R - replace mode (overwrite)
		{
			ID:          "substitution_003",
			Module:      ModuleSubstitution,
			Level:       2,
			Type:        ExerciseLesson,
			Code:        []string{"const VERSION = '1.0.0';"},
			CursorPos:   Position{Line: 0, Col: 18},
			Mission:     "Enter Replace mode with R (in real Vim: R then type new version)",
			Solutions:   []string{"R"},
			Optimal:     "R",
			Hint:        "R enters replace mode - each character you type overwrites the existing one",
			Explanation: "The 'R' command enters Replace mode, where each character you type replaces the character under the cursor, then advances. It's like the Insert key on a keyboard - perfect for overwriting text of the same length. In real Vim you'd type: R2.1.5<Esc>",
			TimeoutSecs: 30,
			Points:      15,
		},
		{
			ID:          "substitution_004",
			Module:      ModuleSubstitution,
			Level:       2,
			Type:        ExerciseLesson,
			Code:        []string{"const DATE = '2023-01-15';"},
			CursorPos:   Position{Line: 0, Col: 14},
			Mission:     "Enter Replace mode with R to overwrite the date (in real Vim: type the new date)",
			Solutions:   []string{"R"},
			Optimal:     "R",
			Hint:        "Position on the '2' and use R to start overwriting",
			Explanation: "Replace mode is ideal when you need to change text but keep the same length. The original text is overwritten character by character as you type. In real Vim: R2024-12-25<Esc>",
			TimeoutSecs: 30,
			Points:      15,
		},
		// Lessons 5-6: s - substitute (delete char + insert)
		{
			ID:          "substitution_005",
			Module:      ModuleSubstitution,
			Level:       3,
			Type:        ExerciseLesson,
			Code:        []string{"const pi = 3.14;"},
			CursorPos:   Position{Line: 0, Col: 15},
			Mission:     "Use s to substitute the ';' (in real Vim: s then type 159;)",
			Solutions:   []string{"s"},
			Optimal:     "s",
			Hint:        "s deletes the character under cursor and enters insert mode",
			Explanation: "The 's' command (substitute) deletes the character under the cursor and enters insert mode. It's equivalent to 'cl' - delete character and insert. Use it when you need to replace a character with multiple characters. In real Vim: s159;<Esc>",
			TimeoutSecs: 30,
			Points:      15,
		},
		{
			ID:          "substitution_006",
			Module:      ModuleSubstitution,
			Level:       3,
			Type:        ExerciseLesson,
			Code:        []string{"let x = getValue();"},
			CursorPos:   Position{Line: 0, Col: 4},
			Mission:     "Use s to substitute variable 'x' (in real Vim: s then type 'result')",
			Solutions:   []string{"s"},
			Optimal:     "s",
			Hint:        "s removes one character and enters insert mode for replacement",
			Explanation: "Using 's' on a single character variable lets you replace it with a longer, more descriptive name. This is more efficient than 'xi' when the character needs to be deleted. In real Vim: sresult<Esc>",
			TimeoutSecs: 30,
			Points:      15,
		},
		// Lessons 7-8: S - substitute line
		{
			ID:          "substitution_007",
			Module:      ModuleSubstitution,
			Level:       4,
			Type:        ExerciseLesson,
			Code:        []string{"// TODO: implement later"},
			CursorPos:   Position{Line: 0, Col: 0},
			Mission:     "Use S to clear the line for replacement (in real Vim: type new content)",
			Solutions:   []string{"S"},
			Optimal:     "S",
			Hint:        "S deletes the entire line content and enters insert mode",
			Explanation: "The 'S' command (substitute line) deletes the entire line's content (preserving indentation) and enters insert mode. It's equivalent to 'cc' - change the whole line. In real Vim: S// DONE: feature complete<Esc>",
			TimeoutSecs: 30,
			Points:      20,
		},
		{
			ID:          "substitution_008",
			Module:      ModuleSubstitution,
			Level:       4,
			Type:        ExerciseLesson,
			Code:        []string{"  const oldValue = 42;  // deprecated"},
			CursorPos:   Position{Line: 0, Col: 10},
			Mission:     "Use S to clear line content (preserves indent) for new code",
			Solutions:   []string{"S"},
			Optimal:     "S",
			Hint:        "S clears line content but keeps the indentation",
			Explanation: "When S is used, it preserves the leading whitespace/indentation of the line. This is helpful when rewriting code while maintaining proper formatting. In real Vim you'd type the new code after S.",
			TimeoutSecs: 30,
			Points:      20,
		},
		// Lessons 9-10: ~ - toggle case
		{
			ID:          "substitution_009",
			Module:      ModuleSubstitution,
			Level:       5,
			Type:        ExerciseLesson,
			Code:        []string{"const maxValue = 100;"},
			CursorPos:   Position{Line: 0, Col: 6},
			Mission:     "Toggle 'm' to 'M' using ~ to start making it MAX_VALUE style",
			Solutions:   []string{"~"},
			Optimal:     "~",
			Hint:        "~ toggles the case of the character under cursor and moves right",
			Explanation: "The '~' command toggles the case of the character under the cursor (lowercase becomes uppercase and vice versa), then moves the cursor one position to the right.",
			TimeoutSecs: 30,
			Points:      10,
		},
		{
			ID:          "substitution_010",
			Module:      ModuleSubstitution,
			Level:       5,
			Type:        ExerciseLesson,
			Code:        []string{"const api = 'REST';"},
			CursorPos:   Position{Line: 0, Col: 13},
			Mission:     "Toggle 'REST' to 'rest' by pressing ~ four times or using 4~",
			Solutions:   []string{"~~~~", "4~"},
			Optimal:     "4~",
			Hint:        "You can use a count before ~ to toggle multiple characters",
			Explanation: "Like most Vim commands, ~ accepts a count. '4~' toggles the case of 4 characters. This is faster than pressing ~ repeatedly for multiple characters.",
			TimeoutSecs: 45,
			Points:      15,
		},
		// Lessons 11-12: gu + motion - lowercase
		{
			ID:          "substitution_011",
			Module:      ModuleSubstitution,
			Level:       6,
			Type:        ExerciseLesson,
			Code:        []string{"const ERROR_MESSAGE = 'Not Found';"},
			CursorPos:   Position{Line: 0, Col: 6},
			Mission:     "Convert 'ERROR_MESSAGE' to 'error_message' using gu and a motion",
			Solutions:   []string{"guiW", "gue", "gu2e", "guE"},
			Optimal:     "guiW",
			Hint:        "gu followed by a motion lowercases the text covered by that motion",
			Explanation: "The 'gu' operator converts text to lowercase. Combined with motions like 'w' (word), 'iW' (inner WORD), or 'e' (end of word), it can lowercase any amount of text efficiently.",
			TimeoutSecs: 45,
			Points:      20,
		},
		{
			ID:          "substitution_012",
			Module:      ModuleSubstitution,
			Level:       6,
			Type:        ExerciseLesson,
			Code:        []string{"const STATUS = 'ACTIVE';"},
			CursorPos:   Position{Line: 0, Col: 0},
			Mission:     "Convert the entire line to lowercase using guu",
			Solutions:   []string{"guu", "gugu"},
			Optimal:     "guu",
			Hint:        "guu lowercases the entire current line",
			Explanation: "Doubling the motion operator (guu) applies it to the entire line - a common Vim pattern. You can also use 'gugu' which does the same thing.",
			TimeoutSecs: 45,
			Points:      20,
		},
		// Lessons 13-14: gU + motion - uppercase
		{
			ID:          "substitution_013",
			Module:      ModuleSubstitution,
			Level:       7,
			Type:        ExerciseLesson,
			Code:        []string{"const api_key = 'secret123';"},
			CursorPos:   Position{Line: 0, Col: 6},
			Mission:     "Convert 'api_key' to 'API_KEY' using gU and a motion",
			Solutions:   []string{"gUiW", "gUe", "gU2e", "gUE"},
			Optimal:     "gUiW",
			Hint:        "gU followed by a motion uppercases the text covered by that motion",
			Explanation: "The 'gU' operator converts text to uppercase. It's the counterpart to 'gu'. Combined with motions, it's perfect for converting variable names to CONSTANT_CASE.",
			TimeoutSecs: 45,
			Points:      20,
		},
		{
			ID:          "substitution_014",
			Module:      ModuleSubstitution,
			Level:       7,
			Type:        ExerciseLesson,
			Code:        []string{"function getUser() { return user; }"},
			CursorPos:   Position{Line: 0, Col: 9},
			Mission:     "Convert 'getUser' to 'GETUSER' using gUiw",
			Solutions:   []string{"gUiw", "gUe", "gUw"},
			Optimal:     "gUiw",
			Hint:        "gU with word motions uppercases entire words",
			Explanation: "Using gU with word motions (w, e, iw) is efficient for uppercasing identifiers. 'gUiw' uppercases the inner word regardless of cursor position within the word.",
			TimeoutSecs: 45,
			Points:      20,
		},
		// Lesson 15: J - join lines
		{
			ID:     "substitution_015",
			Module: ModuleSubstitution,
			Level:  8,
			Type:   ExerciseLesson,
			Code: []string{
				"const message = 'Hello'",
				"  + ' World';",
			},
			CursorPos:   Position{Line: 0, Col: 0},
			Mission:     "Join the two lines into one using J",
			Solutions:   []string{"J"},
			Optimal:     "J",
			Hint:        "J joins the current line with the next line, adding a space between them",
			Explanation: "The 'J' command joins the current line with the line below it, replacing the line break with a space. It also removes leading whitespace from the joined line. Use 'gJ' to join without adding a space.",
			TimeoutSecs: 30,
			Points:      15,
		},
	}
}

// getSubstitutionPractice converts substitution lessons to practice exercises
func getSubstitutionPractice() []Exercise {
	lessons := getSubstitutionLessons()
	practices := make([]Exercise, len(lessons))

	for i, lesson := range lessons {
		practice := lesson
		practice.ID = "substitution_p" + lesson.ID[len("substitution_"):]
		practice.Type = ExercisePractice
		practice.Hint = ""
		practice.Explanation = ""
		practice.TimeoutSecs = lesson.TimeoutSecs - 10
		if practice.TimeoutSecs < 15 {
			practice.TimeoutSecs = 15
		}
		practice.Points = lesson.Points + 10
		practices[i] = practice
	}

	return practices
}

// getSubstitutionBoss returns the boss challenge for the Substitution module
func getSubstitutionBoss() *BossExercise {
	return &BossExercise{
		ID:        "substitution_boss",
		Module:    ModuleSubstitution,
		Name:      "The Transformer",
		Lives:     3,
		BonusTime: 45,
		Steps: []BossStep{
			{
				TimeLimit: 8,
				Exercise: Exercise{
					ID:     "substitution_boss_1",
					Module: ModuleSubstitution,
					Level:  8,
					Type:   ExerciseBoss,
					Code: []string{
						"// todo: fix this code",
						"const apiEndpoint = 'htps://api.example.com';",
					},
					CursorPos: Position{Line: 0, Col: 3},
					Mission:   "Change 'todo' to 'TODO' using gU",
					Solutions: []string{"gUw", "gUe", "gUiw"},
					Optimal:   "gUiw",
					Points:    50,
				},
			},
			{
				TimeLimit: 10,
				Exercise: Exercise{
					ID:     "substitution_boss_2",
					Module: ModuleSubstitution,
					Level:  8,
					Type:   ExerciseBoss,
					Code: []string{
						"// TODO: fix this code",
						"const apiEndpoint = 'htps://api.example.com';",
					},
					CursorPos: Position{Line: 1, Col: 22},
					Mission:   "Fix 'htps' to 'https' using r twice (replace 'p' with 't', then 's' with 'p')",
					Solutions: []string{"rtlrp", "rt$F:hrp"},
					Optimal:   "rtlrp",
					Points:    50,
				},
			},
			{
				TimeLimit: 8,
				Exercise: Exercise{
					ID:     "substitution_boss_3",
					Module: ModuleSubstitution,
					Level:  8,
					Type:   ExerciseBoss,
					Code: []string{
						"const maxRetries = 3;",
						"const timeout = 5000;",
					},
					CursorPos: Position{Line: 0, Col: 6},
					Mission:   "Change 'maxRetries' to 'MAX_RETRIES' using gU",
					Solutions: []string{"gUiW", "gUe", "gU2e"},
					Optimal:   "gUiW",
					Points:    50,
				},
			},
			{
				TimeLimit: 6,
				Exercise: Exercise{
					ID:     "substitution_boss_4",
					Module: ModuleSubstitution,
					Level:  8,
					Type:   ExerciseBoss,
					Code: []string{
						"const errorCode",
						"  = 500;",
					},
					CursorPos: Position{Line: 0, Col: 0},
					Mission:   "Join the two lines into one using J",
					Solutions: []string{"J"},
					Optimal:   "J",
					Points:    50,
				},
			},
			{
				TimeLimit: 10,
				Exercise: Exercise{
					ID:        "substitution_boss_5",
					Module:    ModuleSubstitution,
					Level:     8,
					Type:      ExerciseBoss,
					Code:      []string{"const VERSION = '1.0.0';"},
					CursorPos: Position{Line: 0, Col: 18},
					Mission:   "Enter Replace mode with R to change version (in real Vim: type 2.5.0)",
					Solutions: []string{"R"},
					Optimal:   "R",
					Points:    50,
				},
			},
		},
	}
}
