package trainer

import (
	"strings"
)

// ValidationResult contains detailed validation information
type ValidationResult struct {
	IsCorrect       bool     // Answer reaches the correct position
	IsInSolutions   bool     // Answer is in the predefined solutions list
	IsOptimal       bool     // Answer is the optimal solution
	TargetPosition  Position // Where the answer should end up
	ActualPosition  Position // Where the answer actually ended up
	OptimalSolution string   // The best solution
	AllSolutions    []string // All predefined valid solutions
}

// ValidateAnswerDetailed performs comprehensive validation using the simulator
func ValidateAnswerDetailed(exercise *Exercise, answer string) ValidationResult {
	result := ValidationResult{
		OptimalSolution: exercise.Optimal,
		AllSolutions:    exercise.Solutions,
	}

	if exercise == nil || answer == "" {
		return result
	}

	answer = strings.TrimSpace(answer)
	if answer == "" {
		return result
	}

	// Check if it's in the predefined solutions (normalize both for comparison)
	for _, sol := range exercise.Solutions {
		if answer == strings.TrimSpace(sol) {
			result.IsInSolutions = true
			break
		}
	}

	// Check if it's optimal (normalize for comparison)
	result.IsOptimal = answer == strings.TrimSpace(exercise.Optimal)

	// Detect if this is an exercise that shouldn't use simulator
	// 1. Ex commands (start with : / ?)
	// 2. Substitution module (r, R, s, S, ~, etc. are edit commands)
	// 3. Macros module (q, @, :normal, :g/)
	// 4. Regex module (/, ?, :vimgrep, etc.)
	isExCommand := len(exercise.Solutions) > 0 && len(exercise.Solutions[0]) > 0 &&
		(exercise.Solutions[0][0] == ':' || exercise.Solutions[0][0] == '/' || exercise.Solutions[0][0] == '?')
	isNonMotionModule := exercise.Module == ModuleSubstitution ||
		exercise.Module == ModuleMacros ||
		exercise.Module == ModuleRegex
	skipSimulation := isExCommand || isNonMotionModule

	if skipSimulation {
		// For non-motion exercises, correct if it matches any predefined solution
		result.IsCorrect = result.IsInSolutions
		result.TargetPosition = exercise.CursorPos
		result.ActualPosition = exercise.CursorPos
		return result
	}

	// Use simulator to check if answer reaches the correct position
	// First, find the target position by simulating the optimal solution
	targetPos := SimulateMotions(exercise.CursorPos, exercise.Code, exercise.Optimal)
	result.TargetPosition = Position{Line: targetPos.Line, Col: targetPos.Col}

	// Now simulate the user's answer
	actualPos := SimulateMotions(exercise.CursorPos, exercise.Code, answer)
	result.ActualPosition = Position{Line: actualPos.Line, Col: actualPos.Col}

	// Answer is correct if it reaches the same position as the optimal solution
	result.IsCorrect = (actualPos.Line == targetPos.Line && actualPos.Col == targetPos.Col)

	return result
}

// ValidateAnswer checks if an answer is valid for an exercise
// Now uses simulator to accept any answer that reaches the correct position
func ValidateAnswer(exercise *Exercise, answer string) bool {
	if exercise == nil {
		return false
	}

	answer = strings.TrimSpace(answer)
	if answer == "" {
		return false
	}

	// First check predefined solutions (fast path) - normalize both for comparison
	for _, sol := range exercise.Solutions {
		if answer == strings.TrimSpace(sol) {
			return true
		}
	}

	// Detect if this exercise shouldn't use simulator
	isExCommand := len(exercise.Solutions) > 0 && len(exercise.Solutions[0]) > 0 &&
		(exercise.Solutions[0][0] == ':' || exercise.Solutions[0][0] == '/' || exercise.Solutions[0][0] == '?')
	isNonMotionModule := exercise.Module == ModuleSubstitution ||
		exercise.Module == ModuleMacros ||
		exercise.Module == ModuleRegex
	skipSimulation := isExCommand || isNonMotionModule

	if skipSimulation {
		// For non-motion exercises, only predefined solutions are valid
		return false
	}

	// Use simulator to check if answer reaches correct position
	// This allows creative solutions not in the predefined list
	targetPos := SimulateMotions(exercise.CursorPos, exercise.Code, exercise.Optimal)
	actualPos := SimulateMotions(exercise.CursorPos, exercise.Code, answer)

	return actualPos.Line == targetPos.Line && actualPos.Col == targetPos.Col
}

// IsOptimalAnswer checks if the answer is the optimal solution
func IsOptimalAnswer(exercise *Exercise, answer string) bool {
	if exercise == nil {
		return false
	}

	answer = strings.TrimSpace(answer)
	return answer == exercise.Optimal
}

// IsInSolutions checks if the answer is in the predefined solutions list
func IsInSolutions(exercise *Exercise, answer string) bool {
	if exercise == nil {
		return false
	}

	answer = strings.TrimSpace(answer)
	for _, sol := range exercise.Solutions {
		if answer == strings.TrimSpace(sol) {
			return true
		}
	}
	return false
}

// GetAlternativeSolutions returns all valid solutions except the one provided
func GetAlternativeSolutions(exercise *Exercise, usedAnswer string) []string {
	if exercise == nil {
		return nil
	}

	usedAnswer = strings.TrimSpace(usedAnswer)
	alternatives := make([]string, 0, len(exercise.Solutions))

	for _, sol := range exercise.Solutions {
		if sol != usedAnswer {
			alternatives = append(alternatives, sol)
		}
	}

	return alternatives
}

// FormatSolutionsHint returns a formatted string with all solutions
func FormatSolutionsHint(exercise *Exercise) string {
	if exercise == nil || len(exercise.Solutions) == 0 {
		return ""
	}

	if len(exercise.Solutions) == 1 {
		return exercise.Solutions[0]
	}

	// Format as "optimal (or alt1, alt2)"
	var alternatives []string
	for _, sol := range exercise.Solutions {
		if sol != exercise.Optimal {
			alternatives = append(alternatives, sol)
		}
	}

	if len(alternatives) == 0 {
		return exercise.Optimal
	}

	return exercise.Optimal + " (or " + strings.Join(alternatives, ", ") + ")"
}

// CalculatePoints calculates points earned for an exercise
func CalculatePoints(exercise *Exercise, timeSeconds float64, isOptimal bool, comboMultiplier int) int {
	if exercise == nil {
		return 0
	}

	// Ensure combo is at least 1
	if comboMultiplier < 1 {
		comboMultiplier = 1
	}

	points := float64(exercise.Points)

	// Optimal bonus: 50%
	if isOptimal {
		points *= 1.5
	}

	// Speed bonus: up to 25% for very fast answers (under 2 seconds)
	if exercise.TimeoutSecs > 0 && timeSeconds < 2.0 {
		points *= 1.25
	}

	// Apply combo multiplier
	points *= float64(comboMultiplier)

	return int(points)
}
