package trainer

// =============================================================================
// MACROS MODULE EXERCISES
// =============================================================================

func getMacrosLessons() []Exercise {
	return []Exercise{
		// Lesson 1: qa - start recording to register a
		{
			ID:     "macros_001",
			Module: ModuleMacros,
			Level:  1,
			Type:   ExerciseLesson,
			Code: []string{
				"item1: value1",
				"item2: value2",
				"item3: value3",
				"item4: value4",
			},
			CursorPos:   Position{Line: 0, Col: 0},
			Mission:     "Start recording a macro to register 'a' using qa",
			Solutions:   []string{"qa"},
			Optimal:     "qa",
			Hint:        "q followed by a register letter (a-z) starts recording",
			Explanation: "qa starts recording all your keystrokes to register 'a'. Macros are Vim's way of recording and playing back keystrokes - like a tape recorder for your editing!",
			TimeoutSecs: 30,
			Points:      10,
		},
		// Lesson 2: qa - more practice
		{
			ID:     "macros_002",
			Module: ModuleMacros,
			Level:  1,
			Type:   ExerciseLesson,
			Code: []string{
				"const a = 1;",
				"const b = 2;",
				"const c = 3;",
				"const d = 4;",
			},
			CursorPos:   Position{Line: 0, Col: 0},
			Mission:     "Start recording a macro to register 'a' (this will record edits)",
			Solutions:   []string{"qa"},
			Optimal:     "qa",
			Hint:        "Press q then a to begin recording to register a",
			Explanation: "Once you press qa, every keystroke you make will be recorded until you press q again. This is the foundation of automation in Vim.",
			TimeoutSecs: 30,
			Points:      10,
		},
		// Lesson 3: q - stop recording
		{
			ID:     "macros_003",
			Module: ModuleMacros,
			Level:  1,
			Type:   ExerciseLesson,
			Code: []string{
				"[RECORDING @a]",
				"const item = getValue();",
			},
			CursorPos:   Position{Line: 0, Col: 0},
			Mission:     "Stop the macro recording with q",
			Solutions:   []string{"q"},
			Optimal:     "q",
			Hint:        "Just press q to stop recording",
			Explanation: "Pressing q while recording stops the recording. Your macro is now saved in register 'a' and ready to be played back.",
			TimeoutSecs: 30,
			Points:      10,
		},
		// Lesson 4: q - stop recording (reinforcement)
		{
			ID:     "macros_004",
			Module: ModuleMacros,
			Level:  1,
			Type:   ExerciseLesson,
			Code: []string{
				"[RECORDING @b]",
				"line one",
				"line two",
			},
			CursorPos:   Position{Line: 0, Col: 0},
			Mission:     "You're recording to register 'b'. Stop the recording.",
			Solutions:   []string{"q"},
			Optimal:     "q",
			Hint:        "q stops any recording regardless of the register",
			Explanation: "Remember: qa starts recording to 'a', qb to 'b', etc. A single q stops any recording. The register letter is only needed to start.",
			TimeoutSecs: 30,
			Points:      10,
		},
		// Lesson 5: @a - play macro from register a
		{
			ID:     "macros_005",
			Module: ModuleMacros,
			Level:  2,
			Type:   ExerciseLesson,
			Code: []string{
				"TODO item 1",
				"TODO item 2",
				"TODO item 3",
			},
			CursorPos:   Position{Line: 0, Col: 0},
			Mission:     "Play the macro stored in register 'a' using @a",
			Solutions:   []string{"@a"},
			Optimal:     "@a",
			Hint:        "@ followed by register letter plays that macro",
			Explanation: "@ executes the macro stored in a register. If you recorded 'dwj' in register 'a', @a will delete a word and move down - exactly what you recorded!",
			TimeoutSecs: 30,
			Points:      15,
		},
		// Lesson 6: @a - play macro (reinforcement)
		{
			ID:     "macros_006",
			Module: ModuleMacros,
			Level:  2,
			Type:   ExerciseLesson,
			Code: []string{
				"old_name = value1",
				"old_name = value2",
				"old_name = value3",
			},
			CursorPos:   Position{Line: 0, Col: 0},
			Mission:     "Execute the macro in register 'a' to transform this line",
			Solutions:   []string{"@a"},
			Optimal:     "@a",
			Hint:        "@ + register letter executes the recorded keystrokes",
			Explanation: "Macros are perfect for repetitive edits. Record once, play many times. The key insight: make your macro work on ONE line, then repeat it.",
			TimeoutSecs: 30,
			Points:      15,
		},
		// Lesson 7: @@ - repeat last macro
		{
			ID:     "macros_007",
			Module: ModuleMacros,
			Level:  2,
			Type:   ExerciseLesson,
			Code: []string{
				"item1: done",
				"item2: pending",
				"item3: pending",
				"item4: pending",
			},
			CursorPos:   Position{Line: 1, Col: 0},
			Mission:     "Repeat the last macro you just played using @@",
			Solutions:   []string{"@@"},
			Optimal:     "@@",
			Hint:        "@@ repeats the most recently executed macro",
			Explanation: "@@ is a shortcut to repeat the last macro you executed. Much faster than typing @a again! Works with any register that was last used.",
			TimeoutSecs: 30,
			Points:      15,
		},
		// Lesson 8: @@ - repeat last macro (reinforcement)
		{
			ID:     "macros_008",
			Module: ModuleMacros,
			Level:  2,
			Type:   ExerciseLesson,
			Code: []string{
				"row1 = transform(data1);",
				"row2 = transform(data2);",
				"row3 = transform(data3);",
			},
			CursorPos:   Position{Line: 1, Col: 0},
			Mission:     "Apply the same macro again with @@",
			Solutions:   []string{"@@"},
			Optimal:     "@@",
			Hint:        "Double @ replays whatever macro you last executed",
			Explanation: "The @@ command is essential for rapid macro application. Pattern: @a to play once, then @@ @@ @@ for subsequent lines.",
			TimeoutSecs: 30,
			Points:      15,
		},
		// Lesson 9: 5@a - play macro 5 times
		{
			ID:     "macros_009",
			Module: ModuleMacros,
			Level:  3,
			Type:   ExerciseLesson,
			Code: []string{
				"line 1 - needs fix",
				"line 2 - needs fix",
				"line 3 - needs fix",
				"line 4 - needs fix",
				"line 5 - needs fix",
			},
			CursorPos:   Position{Line: 0, Col: 0},
			Mission:     "Play macro 'a' five times in one command using 5@a",
			Solutions:   []string{"5@a"},
			Optimal:     "5@a",
			Hint:        "Prefix @ with a count to repeat multiple times",
			Explanation: "Like most Vim commands, @ accepts a count. 5@a runs the macro 5 times. Perfect for applying the same edit to many lines at once!",
			TimeoutSecs: 30,
			Points:      20,
		},
		// Lesson 10: Count with macro (reinforcement)
		{
			ID:     "macros_010",
			Module: ModuleMacros,
			Level:  3,
			Type:   ExerciseLesson,
			Code: []string{
				"data[0] = input[0];",
				"data[1] = input[1];",
				"data[2] = input[2];",
				"data[3] = input[3];",
				"data[4] = input[4];",
				"data[5] = input[5];",
				"data[6] = input[6];",
				"data[7] = input[7];",
				"data[8] = input[8];",
				"data[9] = input[9];",
			},
			CursorPos:   Position{Line: 0, Col: 0},
			Mission:     "Apply macro 'a' to all 10 lines with 10@a",
			Solutions:   []string{"10@a"},
			Optimal:     "10@a",
			Hint:        "Use count before @ to run multiple times",
			Explanation: "Pro tip: if you have 100 lines, just use 100@a. If the macro fails partway (e.g., end of file), it stops safely. Overestimating is fine!",
			TimeoutSecs: 30,
			Points:      20,
		},
		// Lesson 11: "ay - yank to register a
		{
			ID:     "macros_011",
			Module: ModuleMacros,
			Level:  4,
			Type:   ExerciseLesson,
			Code: []string{
				"const importantValue = 'remember_this';",
				"",
				"// Need this value below",
			},
			CursorPos:   Position{Line: 0, Col: 6},
			Mission:     "Yank the word 'importantValue' to register 'a' using \"ayiw",
			Solutions:   []string{"\"ayiw"},
			Optimal:     "\"ayiw",
			Hint:        "\"a before yank stores it in register a instead of default",
			Explanation: "Registers aren't just for macros! \"a before any yank/delete stores that text in register 'a'. You have 26 named registers (a-z) to use.",
			TimeoutSecs: 30,
			Points:      25,
		},
		// Lesson 12: "ap - paste from register a
		{
			ID:     "macros_012",
			Module: ModuleMacros,
			Level:  4,
			Type:   ExerciseLesson,
			Code: []string{
				"// Register 'a' contains: myFunction",
				"const result = ();",
			},
			CursorPos:   Position{Line: 1, Col: 15},
			Mission:     "Paste the contents of register 'a' using \"ap",
			Solutions:   []string{"\"ap"},
			Optimal:     "\"ap",
			Hint:        "\"a before paste retrieves from register a",
			Explanation: "\"ap pastes from register 'a'. Unlike the default register which gets overwritten often, named registers (a-z) persist until you overwrite them.",
			TimeoutSecs: 30,
			Points:      25,
		},
		// Lesson 13: "+y - yank to system clipboard
		{
			ID:     "macros_013",
			Module: ModuleMacros,
			Level:  5,
			Type:   ExerciseLesson,
			Code: []string{
				"const API_KEY = 'abc123xyz789';",
				"",
				"// Copy the key to share externally",
			},
			CursorPos:   Position{Line: 0, Col: 16},
			Mission:     "Yank the string content to system clipboard using \"+yi'",
			Solutions:   []string{"\"+yi'", "\"+yi\""},
			Optimal:     "\"+yi'",
			Hint:        "\"+  is the system clipboard register",
			Explanation: "\"+ is the system clipboard! \"+y yanks to your OS clipboard so you can paste in other apps. \"+p pastes from clipboard into Vim.",
			TimeoutSecs: 30,
			Points:      30,
		},
		// Lesson 14: "+p - paste from system clipboard
		{
			ID:     "macros_014",
			Module: ModuleMacros,
			Level:  5,
			Type:   ExerciseLesson,
			Code: []string{
				"// Clipboard contains: external_data",
				"const value = '';",
			},
			CursorPos:   Position{Line: 1, Col: 14},
			Mission:     "Paste from system clipboard using \"+p",
			Solutions:   []string{"\"+p"},
			Optimal:     "\"+p",
			Hint:        "\"+ accesses system clipboard, p pastes",
			Explanation: "\"+p pastes whatever is in your system clipboard. This bridges Vim and other applications. Essential for copy-pasting from browsers, docs, etc.",
			TimeoutSecs: 30,
			Points:      30,
		},
		// Lesson 15: "0p - paste last yank (not delete)
		{
			ID:     "macros_015",
			Module: ModuleMacros,
			Level:  5,
			Type:   ExerciseLesson,
			Code: []string{
				"const original = 'keep_this';",
				"const deleted = 'gone';  // You just deleted this",
				"const target = '';       // Want to paste 'keep_this' here",
			},
			CursorPos:   Position{Line: 2, Col: 16},
			Mission:     "Paste the last YANKED text (not deleted) using \"0p",
			Solutions:   []string{"\"0p"},
			Optimal:     "\"0p",
			Hint:        "\"0 is the yank register - it only stores yanked text, never deleted",
			Explanation: "\"0 is the yank register! It always contains your last yank, unaffected by deletes. This solves the problem of 'I yanked something but then deleted and lost it'.",
			TimeoutSecs: 30,
			Points:      30,
		},
	}
}

func getMacrosPractice() []Exercise {
	lessons := getMacrosLessons()
	practice := make([]Exercise, len(lessons))

	for i, ex := range lessons {
		practice[i] = ex
		practice[i].Type = ExercisePractice
		practice[i].ID = "macros_p" + ex.ID[len("macros_"):]
		practice[i].TimeoutSecs = 15 // Shorter timeout for practice
	}

	return practice
}

func getMacrosBoss() *BossExercise {
	return &BossExercise{
		ID:        "macros_boss",
		Module:    ModuleMacros,
		Name:      "The Automation Wizard",
		Lives:     3,
		BonusTime: 30,
		Steps: []BossStep{
			{
				TimeLimit: 5,
				Exercise: Exercise{
					ID:     "macros_boss_1",
					Module: ModuleMacros,
					Level:  5,
					Type:   ExerciseBoss,
					Code: []string{
						"user1: active",
						"user2: active",
						"user3: active",
						"user4: active",
						"user5: active",
					},
					CursorPos: Position{Line: 0, Col: 0},
					Mission:   "Start recording a macro to register 'a'",
					Solutions: []string{"qa"},
					Optimal:   "qa",
					Points:    50,
				},
			},
			{
				TimeLimit: 5,
				Exercise: Exercise{
					ID:     "macros_boss_2",
					Module: ModuleMacros,
					Level:  5,
					Type:   ExerciseBoss,
					Code: []string{
						"[RECORDING @a]",
						"Edits made: f:r|j",
						"Ready to stop",
					},
					CursorPos: Position{Line: 0, Col: 0},
					Mission:   "Stop the macro recording",
					Solutions: []string{"q"},
					Optimal:   "q",
					Points:    50,
				},
			},
			{
				TimeLimit: 5,
				Exercise: Exercise{
					ID:     "macros_boss_3",
					Module: ModuleMacros,
					Level:  5,
					Type:   ExerciseBoss,
					Code: []string{
						"row1: pending",
						"row2: pending",
						"row3: pending",
						"row4: pending",
						"row5: pending",
					},
					CursorPos: Position{Line: 1, Col: 0},
					Mission:   "Execute the macro in register 'a'",
					Solutions: []string{"@a"},
					Optimal:   "@a",
					Points:    50,
				},
			},
			{
				TimeLimit: 5,
				Exercise: Exercise{
					ID:     "macros_boss_4",
					Module: ModuleMacros,
					Level:  5,
					Type:   ExerciseBoss,
					Code: []string{
						"item1: done",
						"item2: done",
						"item3: pending",
						"item4: pending",
						"item5: pending",
					},
					CursorPos: Position{Line: 2, Col: 0},
					Mission:   "Apply macro 'a' to remaining 3 lines with a count",
					Solutions: []string{"3@a"},
					Optimal:   "3@a",
					Points:    50,
				},
			},
			{
				TimeLimit: 6,
				Exercise: Exercise{
					ID:     "macros_boss_5",
					Module: ModuleMacros,
					Level:  5,
					Type:   ExerciseBoss,
					Code: []string{
						"// Yank register contains: 'savedValue'",
						"// You just deleted something else",
						"const x = '';  // paste yanked value here",
					},
					CursorPos: Position{Line: 2, Col: 11},
					Mission:   "Paste from the yank register (\"0)",
					Solutions: []string{"\"0p"},
					Optimal:   "\"0p",
					Points:    50,
				},
			},
		},
	}
}
