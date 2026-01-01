package trainer

import (
	"testing"
)

// =============================================================================
// ANSWER VALIDATION
// =============================================================================

func TestValidateAnswer_ExactMatch(t *testing.T) {
	exercise := &Exercise{
		ID:        "test_001",
		Solutions: []string{"w"},
	}

	if !ValidateAnswer(exercise, "w") {
		t.Error("Exact match should be valid")
	}
}

func TestValidateAnswer_NoMatch(t *testing.T) {
	// Without proper exercise code/position, any answer that doesn't match
	// predefined solutions AND doesn't simulate correctly will fail
	exercise := &Exercise{
		ID:        "test_001",
		Code:      []string{"const foo = 'bar';"},
		CursorPos: Position{Line: 0, Col: 0},
		Solutions: []string{"w"},
		Optimal:   "w",
	}

	// 'b' from position 0 stays at 0, while 'w' goes to col 6
	// So 'b' should not be valid
	if ValidateAnswer(exercise, "b") {
		t.Error("'b' should be invalid (doesn't reach same position as 'w')")
	}
}

func TestValidateAnswer_MultipleSolutions(t *testing.T) {
	exercise := &Exercise{
		ID:        "test_001",
		Solutions: []string{"w", "W", "fe"},
	}

	// All should be valid
	if !ValidateAnswer(exercise, "w") {
		t.Error("'w' should be valid")
	}
	if !ValidateAnswer(exercise, "W") {
		t.Error("'W' should be valid")
	}
	if !ValidateAnswer(exercise, "fe") {
		t.Error("'fe' should be valid")
	}

	// This should not
	if ValidateAnswer(exercise, "e") {
		t.Error("'e' should be invalid")
	}
}

func TestValidateAnswer_TrimsWhitespace(t *testing.T) {
	exercise := &Exercise{
		ID:        "test_001",
		Solutions: []string{"ciw"},
	}

	if !ValidateAnswer(exercise, " ciw ") {
		t.Error("Answer with whitespace should be trimmed and validated")
	}
	if !ValidateAnswer(exercise, "\tciw\n") {
		t.Error("Answer with tabs/newlines should be trimmed and validated")
	}
}

func TestValidateAnswer_EmptyAnswer(t *testing.T) {
	exercise := &Exercise{
		ID:        "test_001",
		Solutions: []string{"w"},
	}

	if ValidateAnswer(exercise, "") {
		t.Error("Empty answer should be invalid")
	}
	if ValidateAnswer(exercise, "   ") {
		t.Error("Whitespace-only answer should be invalid")
	}
}

func TestValidateAnswer_NilExercise(t *testing.T) {
	if ValidateAnswer(nil, "w") {
		t.Error("Nil exercise should return false")
	}
}

func TestValidateAnswer_EmptySolutions(t *testing.T) {
	exercise := &Exercise{
		ID:        "test_001",
		Solutions: []string{},
	}

	if ValidateAnswer(exercise, "anything") {
		t.Error("Exercise with no solutions should never validate")
	}
}

func TestValidateAnswer_CaseSensitive(t *testing.T) {
	exercise := &Exercise{
		ID:        "test_001",
		Solutions: []string{"w", "W"}, // w and W are different in Vim!
	}

	// w is valid, W is valid, but not mixed
	if !ValidateAnswer(exercise, "w") {
		t.Error("'w' should be valid")
	}
	if !ValidateAnswer(exercise, "W") {
		t.Error("'W' should be valid")
	}

	// But these are not in solutions
	exercise2 := &Exercise{
		ID:        "test_002",
		Solutions: []string{"w"}, // Only lowercase
	}
	if ValidateAnswer(exercise2, "W") {
		t.Error("'W' should be invalid when only 'w' is in solutions")
	}
}

func TestValidateAnswer_ComplexCommands(t *testing.T) {
	tests := []struct {
		name      string
		solutions []string
		input     string
		valid     bool
	}{
		{"find char", []string{"f="}, "f=", true},
		{"find char wrong", []string{"f="}, "f-", false},
		{"change inner", []string{"ci\""}, "ci\"", true},
		{"change around", []string{"ca{"}, "ca{", true},
		{"numbered motion", []string{"3w", "www"}, "3w", true},
		{"numbered motion alt", []string{"3w", "www"}, "www", true},
		{"search and change", []string{"cgn"}, "cgn", true},
		{"macro", []string{"@a"}, "@a", true},
		{"repeat", []string{"."}, ".", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exercise := &Exercise{
				ID:        "test",
				Solutions: tt.solutions,
			}
			result := ValidateAnswer(exercise, tt.input)
			if result != tt.valid {
				t.Errorf("ValidateAnswer(%q) = %v, want %v", tt.input, result, tt.valid)
			}
		})
	}
}

// =============================================================================
// IS OPTIMAL CHECK
// =============================================================================

func TestIsOptimalAnswer_ReturnsTrue(t *testing.T) {
	exercise := &Exercise{
		ID:        "test_001",
		Solutions: []string{"w", "W", "fe"},
		Optimal:   "w",
	}

	if !IsOptimalAnswer(exercise, "w") {
		t.Error("'w' should be optimal")
	}
}

func TestIsOptimalAnswer_ReturnsFalse(t *testing.T) {
	exercise := &Exercise{
		ID:        "test_001",
		Solutions: []string{"w", "W", "fe"},
		Optimal:   "w",
	}

	if IsOptimalAnswer(exercise, "W") {
		t.Error("'W' should not be optimal (valid but not optimal)")
	}
	if IsOptimalAnswer(exercise, "fe") {
		t.Error("'fe' should not be optimal")
	}
}

func TestIsOptimalAnswer_TrimsWhitespace(t *testing.T) {
	exercise := &Exercise{
		ID:      "test_001",
		Optimal: "w",
	}

	if !IsOptimalAnswer(exercise, " w ") {
		t.Error("Whitespace should be trimmed")
	}
}

func TestIsOptimalAnswer_NilExercise(t *testing.T) {
	if IsOptimalAnswer(nil, "w") {
		t.Error("Nil exercise should return false")
	}
}

// =============================================================================
// CALCULATE POINTS
// =============================================================================

func TestCalculatePoints_BasePoints(t *testing.T) {
	exercise := &Exercise{
		ID:     "test",
		Points: 100,
	}

	points := CalculatePoints(exercise, 5.0, false, 1)
	if points != 100 {
		t.Errorf("Base points should be 100, got %d", points)
	}
}

func TestCalculatePoints_OptimalBonus(t *testing.T) {
	exercise := &Exercise{
		ID:     "test",
		Points: 100,
	}

	// Optimal answer gives 50% bonus
	points := CalculatePoints(exercise, 5.0, true, 1)
	if points != 150 {
		t.Errorf("Optimal answer should give 150 points (100 + 50%%), got %d", points)
	}
}

func TestCalculatePoints_SpeedBonus(t *testing.T) {
	exercise := &Exercise{
		ID:          "test",
		Points:      100,
		TimeoutSecs: 30,
	}

	// Very fast answer (under 2 seconds) gets speed bonus
	points := CalculatePoints(exercise, 1.5, false, 1)
	if points <= 100 {
		t.Errorf("Fast answer should get speed bonus, got %d", points)
	}
}

func TestCalculatePoints_ComboMultiplier(t *testing.T) {
	exercise := &Exercise{
		ID:     "test",
		Points: 100,
	}

	// x2 combo
	points := CalculatePoints(exercise, 5.0, false, 2)
	if points != 200 {
		t.Errorf("x2 combo should give 200 points, got %d", points)
	}

	// x4 combo
	points = CalculatePoints(exercise, 5.0, false, 4)
	if points != 400 {
		t.Errorf("x4 combo should give 400 points, got %d", points)
	}
}

func TestCalculatePoints_AllBonusesCombined(t *testing.T) {
	exercise := &Exercise{
		ID:          "test",
		Points:      100,
		TimeoutSecs: 30,
	}

	// Optimal + fast + x2 combo
	points := CalculatePoints(exercise, 1.0, true, 2)
	// Base: 100, Optimal: +50%, Speed: +25% (under 2s), Combo: x2
	// (100 * 1.5 * 1.25) * 2 = 375
	if points < 300 {
		t.Errorf("Combined bonuses should give significant points, got %d", points)
	}
}

func TestCalculatePoints_NilExercise(t *testing.T) {
	points := CalculatePoints(nil, 5.0, false, 1)
	if points != 0 {
		t.Errorf("Nil exercise should give 0 points, got %d", points)
	}
}

func TestCalculatePoints_ZeroCombo(t *testing.T) {
	exercise := &Exercise{
		ID:     "test",
		Points: 100,
	}

	// Combo 0 should be treated as 1
	points := CalculatePoints(exercise, 5.0, false, 0)
	if points != 100 {
		t.Errorf("Zero combo should be treated as 1, got %d", points)
	}
}

// =============================================================================
// SIMULATOR-BASED VALIDATION
// =============================================================================

func TestValidateAnswer_AcceptsCreativeSolutions(t *testing.T) {
	// An exercise where 'w' is the optimal solution to reach col 6
	exercise := &Exercise{
		ID:        "test_001",
		Code:      []string{"const userName = 'value';"},
		CursorPos: Position{Line: 0, Col: 0},
		Solutions: []string{"w"},
		Optimal:   "w",
	}

	// 'w' should work (predefined)
	if !ValidateAnswer(exercise, "w") {
		t.Error("Predefined solution 'w' should be valid")
	}

	// 'llllll' (6 times l) should also reach col 6 - creative solution!
	if !ValidateAnswer(exercise, "llllll") {
		t.Error("Creative solution 'llllll' should be valid (reaches same position)")
	}
}

func TestValidateAnswerDetailed_ReturnsFullInfo(t *testing.T) {
	exercise := &Exercise{
		ID:        "test_001",
		Code:      []string{"const userName = 'value';"},
		CursorPos: Position{Line: 0, Col: 0},
		Solutions: []string{"w", "fe"},
		Optimal:   "w",
	}

	// Test optimal answer
	result := ValidateAnswerDetailed(exercise, "w")
	if !result.IsCorrect {
		t.Error("Should be correct")
	}
	if !result.IsOptimal {
		t.Error("'w' should be optimal")
	}
	if !result.IsInSolutions {
		t.Error("'w' should be in solutions list")
	}

	// Test valid but not optimal
	result = ValidateAnswerDetailed(exercise, "llllll")
	if !result.IsCorrect {
		t.Error("Should be correct (reaches same position)")
	}
	if result.IsOptimal {
		t.Error("'llllll' should not be optimal")
	}
	if result.IsInSolutions {
		t.Error("'llllll' should not be in predefined solutions")
	}

	// Test incorrect answer
	result = ValidateAnswerDetailed(exercise, "b")
	if result.IsCorrect {
		t.Error("'b' should not be correct (doesn't reach target)")
	}
}

func TestIsInSolutions(t *testing.T) {
	exercise := &Exercise{
		ID:        "test",
		Solutions: []string{"w", "W", "fe"},
	}

	if !IsInSolutions(exercise, "w") {
		t.Error("'w' should be in solutions")
	}
	if !IsInSolutions(exercise, "fe") {
		t.Error("'fe' should be in solutions")
	}
	if IsInSolutions(exercise, "llllll") {
		t.Error("'llllll' should not be in solutions")
	}
}

func TestGetAlternativeSolutions(t *testing.T) {
	exercise := &Exercise{
		ID:        "test",
		Solutions: []string{"w", "W", "fe"},
	}

	// When user uses 'w', alternatives should be 'W' and 'fe'
	alts := GetAlternativeSolutions(exercise, "w")
	if len(alts) != 2 {
		t.Errorf("Expected 2 alternatives, got %d", len(alts))
	}

	// Check that 'w' is not in alternatives
	for _, alt := range alts {
		if alt == "w" {
			t.Error("Used answer should not be in alternatives")
		}
	}
}

func TestFormatSolutionsHint(t *testing.T) {
	tests := []struct {
		name      string
		solutions []string
		optimal   string
		expected  string
	}{
		{
			name:      "single solution",
			solutions: []string{"w"},
			optimal:   "w",
			expected:  "w",
		},
		{
			name:      "multiple solutions",
			solutions: []string{"w", "fe", "W"},
			optimal:   "w",
			expected:  "w (or fe, W)",
		},
		{
			name:      "optimal is only solution",
			solutions: []string{"$"},
			optimal:   "$",
			expected:  "$",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exercise := &Exercise{
				Solutions: tt.solutions,
				Optimal:   tt.optimal,
			}
			result := FormatSolutionsHint(exercise)
			if result != tt.expected {
				t.Errorf("FormatSolutionsHint() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestFormatSolutionsHint_NilExercise(t *testing.T) {
	result := FormatSolutionsHint(nil)
	if result != "" {
		t.Errorf("Nil exercise should return empty string, got %q", result)
	}
}
