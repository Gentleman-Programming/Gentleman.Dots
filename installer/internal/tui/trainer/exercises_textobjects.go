package trainer

// getTextObjectsLessons returns all lessons for the TextObjects module
func getTextObjectsLessons() []Exercise {
	return []Exercise{
		// Lesson 1: iw - inner word (basic)
		{
			ID:          "textobjects_001",
			Module:      ModuleTextObjects,
			Level:       1,
			Type:        ExerciseLesson,
			Code:        []string{"const userName = 'developer';"},
			CursorPos:   Position{Line: 0, Col: 6},
			Mission:     "Select the word 'userName' using inner word",
			Solutions:   []string{"iw"},
			Optimal:     "iw",
			Hint:        "iw selects the inner word under the cursor (without surrounding spaces)",
			Explanation: "The 'iw' text object selects the word under the cursor. It's 'inner' because it doesn't include surrounding whitespace. This is incredibly useful for quickly selecting variable names, function names, or any word in your code.",
			TimeoutSecs: 30,
			Points:      10,
		},
		// Lesson 2: iw - inner word (in context)
		{
			ID:          "textobjects_002",
			Module:      ModuleTextObjects,
			Level:       1,
			Type:        ExerciseLesson,
			Code:        []string{"function calculateTotal(items, taxRate) {"},
			CursorPos:   Position{Line: 0, Col: 12},
			Mission:     "Select 'calculateTotal' using inner word",
			Solutions:   []string{"iw"},
			Optimal:     "iw",
			Hint:        "Position doesn't matter - iw selects the entire word under cursor",
			Explanation: "Notice how 'iw' selects the entire word 'calculateTotal' even though the cursor is in the middle. Vim considers camelCase as a single word by default.",
			TimeoutSecs: 30,
			Points:      10,
		},
		// Lesson 3: aw - a word (includes space)
		{
			ID:          "textobjects_003",
			Module:      ModuleTextObjects,
			Level:       1,
			Type:        ExerciseLesson,
			Code:        []string{"const result = value + offset;"},
			CursorPos:   Position{Line: 0, Col: 17},
			Mission:     "Select 'value' AND trailing space using 'a word'",
			Solutions:   []string{"aw"},
			Optimal:     "aw",
			Hint:        "aw means 'a word' - includes the word plus one surrounding space",
			Explanation: "The 'aw' text object selects 'a word' including surrounding whitespace. This is perfect when you want to delete a word completely, leaving no extra spaces behind. Use 'daw' to delete a word cleanly.",
			TimeoutSecs: 30,
			Points:      15,
		},
		// Lesson 4: i" - inner double quotes
		{
			ID:          "textobjects_004",
			Module:      ModuleTextObjects,
			Level:       2,
			Type:        ExerciseLesson,
			Code:        []string{`const message = "Hello, World!";`},
			CursorPos:   Position{Line: 0, Col: 20},
			Mission:     "Select the text inside double quotes (not the quotes)",
			Solutions:   []string{`i"`},
			Optimal:     `i"`,
			Hint:        `i" selects everything INSIDE the double quotes`,
			Explanation: `The 'i"' text object selects text inside double quotes, excluding the quotes themselves. This is perfect for changing string contents while keeping the quotes: ci" lets you change the string content.`,
			TimeoutSecs: 30,
			Points:      15,
		},
		// Lesson 5: a" - a double quote (includes quotes)
		{
			ID:          "textobjects_005",
			Module:      ModuleTextObjects,
			Level:       2,
			Type:        ExerciseLesson,
			Code:        []string{`console.log("Processing data...");`},
			CursorPos:   Position{Line: 0, Col: 18},
			Mission:     "Select the entire quoted string INCLUDING the quotes",
			Solutions:   []string{`a"`},
			Optimal:     `a"`,
			Hint:        `a" selects the quotes AND everything inside them`,
			Explanation: `The 'a"' text object includes the quote characters themselves. Use this when you want to delete or replace an entire quoted string, including its delimiters. da" removes the whole string with quotes.`,
			TimeoutSecs: 30,
			Points:      15,
		},
		// Lesson 6: i' - inner single quotes
		{
			ID:          "textobjects_006",
			Module:      ModuleTextObjects,
			Level:       2,
			Type:        ExerciseLesson,
			Code:        []string{`const apiKey = 'sk-abc123xyz';`},
			CursorPos:   Position{Line: 0, Col: 20},
			Mission:     "Select the API key inside single quotes",
			Solutions:   []string{"i'"},
			Optimal:     "i'",
			Hint:        "i' works just like i\" but for single quotes",
			Explanation: "The i' text object works identically to i\" but for single-quoted strings. In JavaScript/TypeScript, both quote styles are common, so knowing both is essential.",
			TimeoutSecs: 30,
			Points:      15,
		},
		// Lesson 7: a' - a single quote
		{
			ID:          "textobjects_007",
			Module:      ModuleTextObjects,
			Level:       2,
			Type:        ExerciseLesson,
			Code:        []string{`import { useState } from 'react';`},
			CursorPos:   Position{Line: 0, Col: 28},
			Mission:     "Select 'react' INCLUDING the surrounding quotes",
			Solutions:   []string{"a'"},
			Optimal:     "a'",
			Hint:        "a' includes the quote characters in the selection",
			Explanation: "Using a' selects the entire quoted import path including quotes. This is useful when you need to replace an entire import path with a different one.",
			TimeoutSecs: 30,
			Points:      15,
		},
		// Lesson 8: i( - inner parentheses
		{
			ID:          "textobjects_008",
			Module:      ModuleTextObjects,
			Level:       3,
			Type:        ExerciseLesson,
			Code:        []string{"const sum = add(firstNum, secondNum);"},
			CursorPos:   Position{Line: 0, Col: 20},
			Mission:     "Select the function arguments inside parentheses",
			Solutions:   []string{"i(", "i)", "ib"},
			Optimal:     "i(",
			Hint:        "i( or i) or ib all select inside parentheses",
			Explanation: "The i( text object (also i) or ib for 'inner block') selects everything inside parentheses. This is incredibly powerful for changing function arguments. ci( lets you rewrite all arguments at once.",
			TimeoutSecs: 30,
			Points:      20,
		},
		// Lesson 9: a( - a parentheses
		{
			ID:          "textobjects_009",
			Module:      ModuleTextObjects,
			Level:       3,
			Type:        ExerciseLesson,
			Code:        []string{"const result = compute(x * y + z);"},
			CursorPos:   Position{Line: 0, Col: 25},
			Mission:     "Select the parentheses AND everything inside",
			Solutions:   []string{"a(", "a)", "ab"},
			Optimal:     "a(",
			Hint:        "a( includes the parentheses themselves",
			Explanation: "The a( text object includes the parentheses. Use da( to delete a function call's argument list entirely, or ca( to change the entire parenthesized expression.",
			TimeoutSecs: 30,
			Points:      20,
		},
		// Lesson 10: i{ - inner braces
		{
			ID:          "textobjects_010",
			Module:      ModuleTextObjects,
			Level:       3,
			Type:        ExerciseLesson,
			Code:        []string{"const config = { debug: true, timeout: 5000 };"},
			CursorPos:   Position{Line: 0, Col: 25},
			Mission:     "Select everything inside the curly braces",
			Solutions:   []string{"i{", "i}", "iB"},
			Optimal:     "i{",
			Hint:        "i{ or i} or iB selects inside curly braces",
			Explanation: "The i{ text object (also i} or iB for 'inner Block') selects content inside braces. Essential for working with objects, function bodies, and code blocks. ci{ rewrites an entire object's contents.",
			TimeoutSecs: 30,
			Points:      20,
		},
		// Lesson 11: a{ - a braces
		{
			ID:          "textobjects_011",
			Module:      ModuleTextObjects,
			Level:       3,
			Type:        ExerciseLesson,
			Code:        []string{"const user = { name: 'John', age: 30 };"},
			CursorPos:   Position{Line: 0, Col: 20},
			Mission:     "Select the entire object INCLUDING the braces",
			Solutions:   []string{"a{", "a}", "aB"},
			Optimal:     "a{",
			Hint:        "a{ includes the curly braces in the selection",
			Explanation: "The a{ text object includes the braces. Useful when replacing an entire object literal or when you need to wrap an object in something else.",
			TimeoutSecs: 30,
			Points:      20,
		},
		// Lesson 12: i[ - inner brackets
		{
			ID:          "textobjects_012",
			Module:      ModuleTextObjects,
			Level:       4,
			Type:        ExerciseLesson,
			Code:        []string{"const numbers = [1, 2, 3, 4, 5];"},
			CursorPos:   Position{Line: 0, Col: 20},
			Mission:     "Select the array elements inside the brackets",
			Solutions:   []string{"i[", "i]"},
			Optimal:     "i[",
			Hint:        "i[ or i] selects inside square brackets",
			Explanation: "The i[ text object selects content inside square brackets. Perfect for array manipulation - ci[ lets you replace all array elements, di[ empties the array while keeping the brackets.",
			TimeoutSecs: 30,
			Points:      20,
		},
		// Lesson 13: a[ - a brackets
		{
			ID:          "textobjects_013",
			Module:      ModuleTextObjects,
			Level:       4,
			Type:        ExerciseLesson,
			Code:        []string{"const item = inventory[selectedIndex];"},
			CursorPos:   Position{Line: 0, Col: 27},
			Mission:     "Select the bracket accessor INCLUDING brackets",
			Solutions:   []string{"a[", "a]"},
			Optimal:     "a[",
			Hint:        "a[ includes the square brackets",
			Explanation: "The a[ text object includes the brackets. Use da[ to remove an array index accessor entirely, converting inventory[selectedIndex] to just inventory.",
			TimeoutSecs: 30,
			Points:      20,
		},
		// Lesson 14: it - inner tag
		{
			ID:          "textobjects_014",
			Module:      ModuleTextObjects,
			Level:       5,
			Type:        ExerciseLesson,
			Code:        []string{"<button>Click me</button>"},
			CursorPos:   Position{Line: 0, Col: 12},
			Mission:     "Select 'Click me' inside the HTML tags",
			Solutions:   []string{"it"},
			Optimal:     "it",
			Hint:        "it selects inside XML/HTML tags",
			Explanation: "The 'it' text object selects content inside matching HTML/XML tags. This is a game-changer for web development - cit lets you change tag content while preserving the tags themselves.",
			TimeoutSecs: 30,
			Points:      25,
		},
		// Lesson 15: at - a tag
		{
			ID:          "textobjects_015",
			Module:      ModuleTextObjects,
			Level:       5,
			Type:        ExerciseLesson,
			Code:        []string{"<div><span>Hello</span></div>"},
			CursorPos:   Position{Line: 0, Col: 12},
			Mission:     "Select the entire <span> element including tags",
			Solutions:   []string{"at"},
			Optimal:     "at",
			Hint:        "at selects the entire tag including opening and closing tags",
			Explanation: "The 'at' text object selects an entire HTML/XML element including both opening and closing tags. Use dat to delete a complete element, or cat to replace it with something else.",
			TimeoutSecs: 30,
			Points:      30,
		},
	}
}

// getTextObjectsPractice converts lessons into practice exercises
func getTextObjectsPractice() []Exercise {
	lessons := getTextObjectsLessons()
	practices := make([]Exercise, len(lessons))

	for i, lesson := range lessons {
		practice := lesson
		practice.ID = "textobjects_p" + lesson.ID[len("textobjects_"):]
		practice.Type = ExercisePractice
		practice.Hint = ""
		practice.TimeoutSecs = 20
		practice.Points = lesson.Points + 10
		practices[i] = practice
	}

	return practices
}

// getTextObjectsBoss returns the boss fight for TextObjects module
func getTextObjectsBoss() *BossExercise {
	return &BossExercise{
		ID:        "textobjects_boss",
		Module:    ModuleTextObjects,
		Name:      "The Text Surgeon",
		Lives:     3,
		BonusTime: 30,
		Steps: []BossStep{
			{
				TimeLimit: 6,
				Exercise: Exercise{
					ID:        "textobjects_boss_1",
					Module:    ModuleTextObjects,
					Level:     5,
					Type:      ExerciseBoss,
					Code:      []string{"const userAuthentication = validateCredentials(input);"},
					CursorPos: Position{Line: 0, Col: 10},
					Mission:   "Select 'userAuthentication' using inner word",
					Solutions: []string{"iw"},
					Optimal:   "iw",
					Points:    50,
				},
			},
			{
				TimeLimit: 6,
				Exercise: Exercise{
					ID:        "textobjects_boss_2",
					Module:    ModuleTextObjects,
					Level:     5,
					Type:      ExerciseBoss,
					Code:      []string{`const endpoint = "/api/v2/users/profile";`},
					CursorPos: Position{Line: 0, Col: 25},
					Mission:   "Select the URL path inside the quotes",
					Solutions: []string{`i"`},
					Optimal:   `i"`,
					Points:    50,
				},
			},
			{
				TimeLimit: 6,
				Exercise: Exercise{
					ID:        "textobjects_boss_3",
					Module:    ModuleTextObjects,
					Level:     5,
					Type:      ExerciseBoss,
					Code:      []string{"const filtered = items.filter(item => item.active && item.visible);"},
					CursorPos: Position{Line: 0, Col: 45},
					Mission:   "Select everything inside the filter parentheses",
					Solutions: []string{"i(", "i)", "ib"},
					Optimal:   "i(",
					Points:    50,
				},
			},
			{
				TimeLimit: 6,
				Exercise: Exercise{
					ID:        "textobjects_boss_4",
					Module:    ModuleTextObjects,
					Level:     5,
					Type:      ExerciseBoss,
					Code:      []string{"const settings = { theme: 'dark', fontSize: 14, autoSave: true };"},
					CursorPos: Position{Line: 0, Col: 30},
					Mission:   "Select all properties inside the braces",
					Solutions: []string{"i{", "i}", "iB"},
					Optimal:   "i{",
					Points:    50,
				},
			},
			{
				TimeLimit: 6,
				Exercise: Exercise{
					ID:        "textobjects_boss_5",
					Module:    ModuleTextObjects,
					Level:     5,
					Type:      ExerciseBoss,
					Code:      []string{"<section><article>Important content here</article></section>"},
					CursorPos: Position{Line: 0, Col: 25},
					Mission:   "Select the entire <article> element including its tags",
					Solutions: []string{"at"},
					Optimal:   "at",
					Points:    50,
				},
			},
		},
	}
}
