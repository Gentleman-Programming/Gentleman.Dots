package trainer

// GetLessons returns lesson exercises for a module
func GetLessons(module ModuleID) []Exercise {
	switch module {
	case ModuleHorizontal:
		return getHorizontalLessons()
	default:
		return []Exercise{}
	}
}

// GetPracticeExercises returns practice exercises for a module
func GetPracticeExercises(module ModuleID) []Exercise {
	switch module {
	case ModuleHorizontal:
		return getHorizontalPractice()
	default:
		return []Exercise{}
	}
}

// GetBoss returns the boss fight for a module
func GetBoss(module ModuleID) *BossExercise {
	switch module {
	case ModuleHorizontal:
		return getHorizontalBoss()
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
			CursorPos:   Position{Line: 0, Col: 28},
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
