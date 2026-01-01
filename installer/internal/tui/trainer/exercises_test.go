package trainer

import (
	"testing"
)

// =============================================================================
// GET LESSONS
// =============================================================================

func TestGetLessons_Horizontal_ReturnsExercises(t *testing.T) {
	lessons := GetLessons(ModuleHorizontal)

	if len(lessons) == 0 {
		t.Error("GetLessons should return exercises for Horizontal module")
	}
}

func TestGetLessons_Horizontal_HasMinimum15(t *testing.T) {
	lessons := GetLessons(ModuleHorizontal)

	if len(lessons) < 15 {
		t.Errorf("Horizontal module should have at least 15 lessons, got %d", len(lessons))
	}
}

func TestGetLessons_Horizontal_AllHaveRequiredFields(t *testing.T) {
	lessons := GetLessons(ModuleHorizontal)

	for i, ex := range lessons {
		if ex.ID == "" {
			t.Errorf("Lesson %d: ID is empty", i)
		}
		if ex.Module != ModuleHorizontal {
			t.Errorf("Lesson %d: Module should be Horizontal, got %s", i, ex.Module)
		}
		if ex.Type != ExerciseLesson {
			t.Errorf("Lesson %d: Type should be Lesson, got %s", i, ex.Type)
		}
		if len(ex.Code) == 0 {
			t.Errorf("Lesson %d: Code is empty", i)
		}
		if ex.Mission == "" {
			t.Errorf("Lesson %d: Mission is empty", i)
		}
		if len(ex.Solutions) == 0 {
			t.Errorf("Lesson %d: Solutions is empty", i)
		}
		if ex.Optimal == "" {
			t.Errorf("Lesson %d: Optimal is empty", i)
		}
		if ex.Hint == "" {
			t.Errorf("Lesson %d: Hint is empty", i)
		}
		if ex.Explanation == "" {
			t.Errorf("Lesson %d: Explanation is empty", i)
		}
		if ex.Points <= 0 {
			t.Errorf("Lesson %d: Points should be positive, got %d", i, ex.Points)
		}
	}
}

func TestGetLessons_Horizontal_OptimalIsInSolutions(t *testing.T) {
	lessons := GetLessons(ModuleHorizontal)

	for i, ex := range lessons {
		found := false
		for _, sol := range ex.Solutions {
			if sol == ex.Optimal {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Lesson %d (%s): Optimal %q not in Solutions %v", i, ex.ID, ex.Optimal, ex.Solutions)
		}
	}
}

func TestGetLessons_Horizontal_UniqueIDs(t *testing.T) {
	lessons := GetLessons(ModuleHorizontal)
	seen := make(map[string]bool)

	for _, ex := range lessons {
		if seen[ex.ID] {
			t.Errorf("Duplicate lesson ID: %s", ex.ID)
		}
		seen[ex.ID] = true
	}
}

func TestGetLessons_Horizontal_ProgressiveLevels(t *testing.T) {
	lessons := GetLessons(ModuleHorizontal)

	// First lesson should be level 1
	if lessons[0].Level != 1 {
		t.Errorf("First lesson should be level 1, got %d", lessons[0].Level)
	}

	// Levels should be non-decreasing (can stay same or increase)
	prevLevel := 0
	for i, ex := range lessons {
		if ex.Level < prevLevel {
			t.Errorf("Lesson %d: Level %d is less than previous %d", i, ex.Level, prevLevel)
		}
		prevLevel = ex.Level
	}
}

func TestGetLessons_Horizontal_CoversBasicMotions(t *testing.T) {
	lessons := GetLessons(ModuleHorizontal)

	// Collect all solutions taught
	allSolutions := make(map[string]bool)
	for _, ex := range lessons {
		for _, sol := range ex.Solutions {
			allSolutions[sol] = true
		}
	}

	// Should cover these basic motions
	requiredMotions := []string{"w", "W", "e", "b", "$", "^", "0"}
	for _, motion := range requiredMotions {
		if !allSolutions[motion] {
			t.Errorf("Horizontal lessons should cover motion %q", motion)
		}
	}
}

func TestGetLessons_Horizontal_CoversFindMotions(t *testing.T) {
	lessons := GetLessons(ModuleHorizontal)

	// At least one lesson should use f{char} pattern
	hasFind := false
	for _, ex := range lessons {
		for _, sol := range ex.Solutions {
			if len(sol) >= 2 && sol[0] == 'f' {
				hasFind = true
				break
			}
		}
	}
	if !hasFind {
		t.Error("Horizontal lessons should include f{char} motion")
	}
}

func TestGetLessons_UnknownModule_ReturnsEmpty(t *testing.T) {
	lessons := GetLessons(ModuleID("unknown"))

	if len(lessons) != 0 {
		t.Errorf("Unknown module should return empty slice, got %d lessons", len(lessons))
	}
}

// =============================================================================
// GET PRACTICE EXERCISES
// =============================================================================

func TestGetPracticeExercises_Horizontal_ReturnsExercises(t *testing.T) {
	exercises := GetPracticeExercises(ModuleHorizontal)

	if len(exercises) == 0 {
		t.Error("GetPracticeExercises should return exercises")
	}
}

func TestGetPracticeExercises_Horizontal_AllPracticeType(t *testing.T) {
	exercises := GetPracticeExercises(ModuleHorizontal)

	for i, ex := range exercises {
		if ex.Type != ExercisePractice {
			t.Errorf("Practice exercise %d: Type should be Practice, got %s", i, ex.Type)
		}
	}
}

func TestGetPracticeExercises_Horizontal_HasMoreThanLessons(t *testing.T) {
	lessons := GetLessons(ModuleHorizontal)
	practice := GetPracticeExercises(ModuleHorizontal)

	// Practice should have at least as many exercises as lessons
	// (or more, since it includes variations)
	if len(practice) < len(lessons) {
		t.Errorf("Practice should have at least %d exercises, got %d", len(lessons), len(practice))
	}
}

func TestGetPracticeExercises_UnknownModule_ReturnsEmpty(t *testing.T) {
	exercises := GetPracticeExercises(ModuleID("unknown"))

	if len(exercises) != 0 {
		t.Errorf("Unknown module should return empty slice, got %d exercises", len(exercises))
	}
}

// =============================================================================
// GET BOSS
// =============================================================================

func TestGetBoss_Horizontal_ReturnsBoss(t *testing.T) {
	boss := GetBoss(ModuleHorizontal)

	if boss == nil {
		t.Fatal("GetBoss should return boss for Horizontal module")
	}
}

func TestGetBoss_Horizontal_HasCorrectName(t *testing.T) {
	boss := GetBoss(ModuleHorizontal)

	if boss.Name != "The Line Walker" {
		t.Errorf("Horizontal boss should be 'The Line Walker', got %q", boss.Name)
	}
}

func TestGetBoss_Horizontal_Has3Lives(t *testing.T) {
	boss := GetBoss(ModuleHorizontal)

	if boss.Lives != 3 {
		t.Errorf("Boss should have 3 lives, got %d", boss.Lives)
	}
}

func TestGetBoss_Horizontal_Has5Steps(t *testing.T) {
	boss := GetBoss(ModuleHorizontal)

	if len(boss.Steps) != 5 {
		t.Errorf("Boss should have 5 steps, got %d", len(boss.Steps))
	}
}

func TestGetBoss_Horizontal_StepsHaveTimeLimits(t *testing.T) {
	boss := GetBoss(ModuleHorizontal)

	for i, step := range boss.Steps {
		if step.TimeLimit <= 0 {
			t.Errorf("Boss step %d: TimeLimit should be positive, got %d", i, step.TimeLimit)
		}
	}
}

func TestGetBoss_Horizontal_StepsHaveExercises(t *testing.T) {
	boss := GetBoss(ModuleHorizontal)

	for i, step := range boss.Steps {
		if step.Exercise.ID == "" {
			t.Errorf("Boss step %d: Exercise ID is empty", i)
		}
		if step.Exercise.Mission == "" {
			t.Errorf("Boss step %d: Exercise Mission is empty", i)
		}
		if len(step.Exercise.Solutions) == 0 {
			t.Errorf("Boss step %d: Exercise Solutions is empty", i)
		}
	}
}

func TestGetBoss_UnknownModule_ReturnsNil(t *testing.T) {
	boss := GetBoss(ModuleID("unknown"))

	if boss != nil {
		t.Error("Unknown module should return nil boss")
	}
}

// =============================================================================
// EXERCISE QUALITY CHECKS
// =============================================================================

func TestHorizontalExercises_MissionsAreInEnglish(t *testing.T) {
	lessons := GetLessons(ModuleHorizontal)

	// Check that missions use English (contain words like "Move", "to", "the", "using", etc.)
	englishIndicators := []string{"Move", "to", "the", "using", "Reach", "Go"}
	hasEnglish := false

	for _, ex := range lessons {
		for _, indicator := range englishIndicators {
			if containsString(ex.Mission, indicator) {
				hasEnglish = true
				break
			}
		}
		if hasEnglish {
			break
		}
	}

	if !hasEnglish {
		t.Error("Lessons should have missions in English")
	}
}

func TestHorizontalExercises_HintsAreHelpful(t *testing.T) {
	lessons := GetLessons(ModuleHorizontal)

	for i, ex := range lessons {
		// Hint should be longer than just the solution
		if len(ex.Hint) <= len(ex.Optimal) {
			t.Errorf("Lesson %d: Hint should be more helpful than just the answer", i)
		}
	}
}

func TestHorizontalExercises_ExplanationsAreEducational(t *testing.T) {
	lessons := GetLessons(ModuleHorizontal)

	for i, ex := range lessons {
		// Explanation should be substantial (at least 20 characters)
		if len(ex.Explanation) < 20 {
			t.Errorf("Lesson %d: Explanation should be more educational (got %d chars)", i, len(ex.Explanation))
		}
	}
}

// Helper function
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && (s[:len(substr)] == substr || containsString(s[1:], substr)))
}

// =============================================================================
// VERTICAL MODULE TESTS
// =============================================================================

func TestGetLessons_Vertical_ReturnsExercises(t *testing.T) {
	lessons := GetLessons(ModuleVertical)

	if len(lessons) == 0 {
		t.Error("GetLessons should return exercises for Vertical module")
	}
}

func TestGetLessons_Vertical_HasMinimum15(t *testing.T) {
	lessons := GetLessons(ModuleVertical)

	if len(lessons) < 15 {
		t.Errorf("Vertical module should have at least 15 lessons, got %d", len(lessons))
	}
}

func TestGetLessons_Vertical_AllHaveRequiredFields(t *testing.T) {
	lessons := GetLessons(ModuleVertical)

	for i, ex := range lessons {
		if ex.ID == "" {
			t.Errorf("Lesson %d: ID is empty", i)
		}
		if ex.Module != ModuleVertical {
			t.Errorf("Lesson %d: Module should be Vertical, got %s", i, ex.Module)
		}
		if ex.Type != ExerciseLesson {
			t.Errorf("Lesson %d: Type should be Lesson, got %s", i, ex.Type)
		}
		if len(ex.Code) == 0 {
			t.Errorf("Lesson %d: Code is empty", i)
		}
		if ex.Mission == "" {
			t.Errorf("Lesson %d: Mission is empty", i)
		}
		if len(ex.Solutions) == 0 {
			t.Errorf("Lesson %d: Solutions is empty", i)
		}
		if ex.Optimal == "" {
			t.Errorf("Lesson %d: Optimal is empty", i)
		}
		if ex.Hint == "" {
			t.Errorf("Lesson %d: Hint is empty", i)
		}
		if ex.Points <= 0 {
			t.Errorf("Lesson %d: Points should be positive, got %d", i, ex.Points)
		}
	}
}

func TestGetLessons_Vertical_OptimalIsInSolutions(t *testing.T) {
	lessons := GetLessons(ModuleVertical)

	for i, ex := range lessons {
		found := false
		for _, sol := range ex.Solutions {
			if sol == ex.Optimal {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Lesson %d (%s): Optimal %q not in Solutions %v", i, ex.ID, ex.Optimal, ex.Solutions)
		}
	}
}

func TestGetLessons_Vertical_UniqueIDs(t *testing.T) {
	lessons := GetLessons(ModuleVertical)
	seen := make(map[string]bool)

	for _, ex := range lessons {
		if seen[ex.ID] {
			t.Errorf("Duplicate lesson ID: %s", ex.ID)
		}
		seen[ex.ID] = true
	}
}

func TestGetLessons_Vertical_CoversBasicMotions(t *testing.T) {
	lessons := GetLessons(ModuleVertical)

	allSolutions := make(map[string]bool)
	for _, ex := range lessons {
		for _, sol := range ex.Solutions {
			allSolutions[sol] = true
		}
	}

	// Should cover these basic vertical motions
	requiredMotions := []string{"j", "k", "gg", "G", "{", "}"}
	for _, motion := range requiredMotions {
		if !allSolutions[motion] {
			t.Errorf("Vertical lessons should cover motion %q", motion)
		}
	}
}

func TestGetLessons_Vertical_HasMultiLineCode(t *testing.T) {
	lessons := GetLessons(ModuleVertical)

	// All vertical lessons should have multi-line code
	for i, ex := range lessons {
		if len(ex.Code) < 2 {
			t.Errorf("Lesson %d: Vertical exercises should have multi-line code, got %d lines", i, len(ex.Code))
		}
	}
}

func TestGetPracticeExercises_Vertical_ReturnsExercises(t *testing.T) {
	exercises := GetPracticeExercises(ModuleVertical)

	if len(exercises) == 0 {
		t.Error("GetPracticeExercises should return exercises for Vertical")
	}
}

func TestGetPracticeExercises_Vertical_AllPracticeType(t *testing.T) {
	exercises := GetPracticeExercises(ModuleVertical)

	for i, ex := range exercises {
		if ex.Type != ExercisePractice {
			t.Errorf("Practice exercise %d: Type should be Practice, got %s", i, ex.Type)
		}
	}
}

func TestGetBoss_Vertical_ReturnsBoss(t *testing.T) {
	boss := GetBoss(ModuleVertical)

	if boss == nil {
		t.Error("GetBoss should return a boss for Vertical module")
	}
}

func TestGetBoss_Vertical_HasCorrectName(t *testing.T) {
	boss := GetBoss(ModuleVertical)

	if boss.Name != "The Code Climber" {
		t.Errorf("Vertical boss name should be 'The Code Climber', got %q", boss.Name)
	}
}

func TestGetBoss_Vertical_Has3Lives(t *testing.T) {
	boss := GetBoss(ModuleVertical)

	if boss.Lives != 3 {
		t.Errorf("Boss should have 3 lives, got %d", boss.Lives)
	}
}

func TestGetBoss_Vertical_Has5Steps(t *testing.T) {
	boss := GetBoss(ModuleVertical)

	if len(boss.Steps) != 5 {
		t.Errorf("Boss should have 5 steps, got %d", len(boss.Steps))
	}
}

func TestGetBoss_Vertical_StepsHaveExercises(t *testing.T) {
	boss := GetBoss(ModuleVertical)

	for i, step := range boss.Steps {
		if step.Exercise.ID == "" {
			t.Errorf("Boss step %d: Exercise ID is empty", i)
		}
		if step.Exercise.Mission == "" {
			t.Errorf("Boss step %d: Exercise Mission is empty", i)
		}
		if len(step.Exercise.Solutions) == 0 {
			t.Errorf("Boss step %d: Exercise Solutions is empty", i)
		}
	}
}

func TestVerticalExercises_SolutionsReachCorrectPosition(t *testing.T) {
	lessons := GetLessons(ModuleVertical)

	for i, ex := range lessons {
		// Simulate optimal solution
		targetPos := SimulateMotions(ex.CursorPos, ex.Code, ex.Optimal)

		// Check all solutions reach the same position
		for _, sol := range ex.Solutions {
			actualPos := SimulateMotions(ex.CursorPos, ex.Code, sol)
			if actualPos.Line != targetPos.Line || actualPos.Col != targetPos.Col {
				t.Errorf("Lesson %d (%s): Solution %q reaches (%d,%d) but optimal %q reaches (%d,%d)",
					i, ex.ID, sol, actualPos.Line, actualPos.Col, ex.Optimal, targetPos.Line, targetPos.Col)
			}
		}
	}
}

// =============================================================================
// TEXT OBJECTS MODULE TESTS
// =============================================================================

func TestGetLessons_TextObjects_ReturnsExercises(t *testing.T) {
	lessons := GetLessons(ModuleTextObjects)
	if len(lessons) == 0 {
		t.Error("GetLessons should return exercises for TextObjects module")
	}
}

func TestGetLessons_TextObjects_HasMinimum15(t *testing.T) {
	lessons := GetLessons(ModuleTextObjects)
	if len(lessons) < 15 {
		t.Errorf("TextObjects module should have at least 15 lessons, got %d", len(lessons))
	}
}

func TestGetLessons_TextObjects_AllHaveRequiredFields(t *testing.T) {
	lessons := GetLessons(ModuleTextObjects)
	for i, ex := range lessons {
		if ex.ID == "" {
			t.Errorf("Lesson %d has empty ID", i)
		}
		if ex.Module != ModuleTextObjects {
			t.Errorf("Lesson %d has wrong module: %s", i, ex.Module)
		}
		if len(ex.Solutions) == 0 {
			t.Errorf("Lesson %d (%s) has no solutions", i, ex.ID)
		}
		if ex.Optimal == "" {
			t.Errorf("Lesson %d (%s) has no optimal solution", i, ex.ID)
		}
	}
}

func TestGetBoss_TextObjects_ReturnsBoss(t *testing.T) {
	boss := GetBoss(ModuleTextObjects)
	if boss == nil {
		t.Error("GetBoss should return a boss for TextObjects module")
	}
}

func TestGetBoss_TextObjects_HasCorrectName(t *testing.T) {
	boss := GetBoss(ModuleTextObjects)
	if boss.Name != "The Text Surgeon" {
		t.Errorf("Boss name should be 'The Text Surgeon', got %s", boss.Name)
	}
}

// =============================================================================
// CHANGE/REPEAT MODULE TESTS
// =============================================================================

func TestGetLessons_ChangeRepeat_ReturnsExercises(t *testing.T) {
	lessons := GetLessons(ModuleChangeRepeat)
	if len(lessons) == 0 {
		t.Error("GetLessons should return exercises for ChangeRepeat module")
	}
}

func TestGetLessons_ChangeRepeat_HasMinimum15(t *testing.T) {
	lessons := GetLessons(ModuleChangeRepeat)
	if len(lessons) < 15 {
		t.Errorf("ChangeRepeat module should have at least 15 lessons, got %d", len(lessons))
	}
}

func TestGetLessons_ChangeRepeat_AllHaveRequiredFields(t *testing.T) {
	lessons := GetLessons(ModuleChangeRepeat)
	for i, ex := range lessons {
		if ex.ID == "" {
			t.Errorf("Lesson %d has empty ID", i)
		}
		if ex.Module != ModuleChangeRepeat {
			t.Errorf("Lesson %d has wrong module: %s", i, ex.Module)
		}
		if len(ex.Solutions) == 0 {
			t.Errorf("Lesson %d (%s) has no solutions", i, ex.ID)
		}
	}
}

func TestGetBoss_ChangeRepeat_ReturnsBoss(t *testing.T) {
	boss := GetBoss(ModuleChangeRepeat)
	if boss == nil {
		t.Error("GetBoss should return a boss for ChangeRepeat module")
	}
}

func TestGetBoss_ChangeRepeat_HasCorrectName(t *testing.T) {
	boss := GetBoss(ModuleChangeRepeat)
	if boss.Name != "The Change Master" {
		t.Errorf("Boss name should be 'The Change Master', got %s", boss.Name)
	}
}

// =============================================================================
// SUBSTITUTION MODULE TESTS
// =============================================================================

func TestGetLessons_Substitution_ReturnsExercises(t *testing.T) {
	lessons := GetLessons(ModuleSubstitution)
	if len(lessons) == 0 {
		t.Error("GetLessons should return exercises for Substitution module")
	}
}

func TestGetLessons_Substitution_HasMinimum15(t *testing.T) {
	lessons := GetLessons(ModuleSubstitution)
	if len(lessons) < 15 {
		t.Errorf("Substitution module should have at least 15 lessons, got %d", len(lessons))
	}
}

func TestGetLessons_Substitution_AllHaveRequiredFields(t *testing.T) {
	lessons := GetLessons(ModuleSubstitution)
	for i, ex := range lessons {
		if ex.ID == "" {
			t.Errorf("Lesson %d has empty ID", i)
		}
		if ex.Module != ModuleSubstitution {
			t.Errorf("Lesson %d has wrong module: %s", i, ex.Module)
		}
		if len(ex.Solutions) == 0 {
			t.Errorf("Lesson %d (%s) has no solutions", i, ex.ID)
		}
	}
}

func TestGetBoss_Substitution_ReturnsBoss(t *testing.T) {
	boss := GetBoss(ModuleSubstitution)
	if boss == nil {
		t.Error("GetBoss should return a boss for Substitution module")
	}
}

func TestGetBoss_Substitution_HasCorrectName(t *testing.T) {
	boss := GetBoss(ModuleSubstitution)
	if boss.Name != "The Transformer" {
		t.Errorf("Boss name should be 'The Transformer', got %s", boss.Name)
	}
}

// =============================================================================
// REGEX MODULE TESTS
// =============================================================================

func TestGetLessons_Regex_ReturnsExercises(t *testing.T) {
	lessons := GetLessons(ModuleRegex)
	if len(lessons) == 0 {
		t.Error("GetLessons should return exercises for Regex module")
	}
}

func TestGetLessons_Regex_HasMinimum15(t *testing.T) {
	lessons := GetLessons(ModuleRegex)
	if len(lessons) < 15 {
		t.Errorf("Regex module should have at least 15 lessons, got %d", len(lessons))
	}
}

func TestGetLessons_Regex_AllHaveRequiredFields(t *testing.T) {
	lessons := GetLessons(ModuleRegex)
	for i, ex := range lessons {
		if ex.ID == "" {
			t.Errorf("Lesson %d has empty ID", i)
		}
		if ex.Module != ModuleRegex {
			t.Errorf("Lesson %d has wrong module: %s", i, ex.Module)
		}
		if len(ex.Solutions) == 0 {
			t.Errorf("Lesson %d (%s) has no solutions", i, ex.ID)
		}
	}
}

func TestGetBoss_Regex_ReturnsBoss(t *testing.T) {
	boss := GetBoss(ModuleRegex)
	if boss == nil {
		t.Error("GetBoss should return a boss for Regex module")
	}
}

func TestGetBoss_Regex_HasCorrectName(t *testing.T) {
	boss := GetBoss(ModuleRegex)
	if boss.Name != "The Pattern Hunter" {
		t.Errorf("Boss name should be 'The Pattern Hunter', got %s", boss.Name)
	}
}

// =============================================================================
// MACROS MODULE TESTS
// =============================================================================

func TestGetLessons_Macros_ReturnsExercises(t *testing.T) {
	lessons := GetLessons(ModuleMacros)
	if len(lessons) == 0 {
		t.Error("GetLessons should return exercises for Macros module")
	}
}

func TestGetLessons_Macros_HasMinimum15(t *testing.T) {
	lessons := GetLessons(ModuleMacros)
	if len(lessons) < 15 {
		t.Errorf("Macros module should have at least 15 lessons, got %d", len(lessons))
	}
}

func TestGetLessons_Macros_AllHaveRequiredFields(t *testing.T) {
	lessons := GetLessons(ModuleMacros)
	for i, ex := range lessons {
		if ex.ID == "" {
			t.Errorf("Lesson %d has empty ID", i)
		}
		if ex.Module != ModuleMacros {
			t.Errorf("Lesson %d has wrong module: %s", i, ex.Module)
		}
		if len(ex.Solutions) == 0 {
			t.Errorf("Lesson %d (%s) has no solutions", i, ex.ID)
		}
	}
}

func TestGetBoss_Macros_ReturnsBoss(t *testing.T) {
	boss := GetBoss(ModuleMacros)
	if boss == nil {
		t.Error("GetBoss should return a boss for Macros module")
	}
}

func TestGetBoss_Macros_HasCorrectName(t *testing.T) {
	boss := GetBoss(ModuleMacros)
	if boss.Name != "The Automation Wizard" {
		t.Errorf("Boss name should be 'The Automation Wizard', got %s", boss.Name)
	}
}

// =============================================================================
// ALL MODULES - CROSS-MODULE TESTS
// =============================================================================

func TestAllModules_HaveUniqueExerciseIDs(t *testing.T) {
	modules := []ModuleID{
		ModuleHorizontal,
		ModuleVertical,
		ModuleTextObjects,
		ModuleChangeRepeat,
		ModuleSubstitution,
		ModuleRegex,
		ModuleMacros,
	}

	allIDs := make(map[string]bool)
	for _, mod := range modules {
		lessons := GetLessons(mod)
		for _, ex := range lessons {
			if allIDs[ex.ID] {
				t.Errorf("Duplicate exercise ID found: %s", ex.ID)
			}
			allIDs[ex.ID] = true
		}
	}
}

func TestAllModules_BossesHave5Steps(t *testing.T) {
	modules := []ModuleID{
		ModuleHorizontal,
		ModuleVertical,
		ModuleTextObjects,
		ModuleChangeRepeat,
		ModuleSubstitution,
		ModuleRegex,
		ModuleMacros,
	}

	for _, mod := range modules {
		boss := GetBoss(mod)
		if boss == nil {
			t.Errorf("Module %s has no boss", mod)
			continue
		}
		if len(boss.Steps) != 5 {
			t.Errorf("Module %s boss should have 5 steps, got %d", mod, len(boss.Steps))
		}
	}
}

func TestAllModules_BossesHave3Lives(t *testing.T) {
	modules := []ModuleID{
		ModuleHorizontal,
		ModuleVertical,
		ModuleTextObjects,
		ModuleChangeRepeat,
		ModuleSubstitution,
		ModuleRegex,
		ModuleMacros,
	}

	for _, mod := range modules {
		boss := GetBoss(mod)
		if boss == nil {
			continue
		}
		if boss.Lives != 3 {
			t.Errorf("Module %s boss should have 3 lives, got %d", mod, boss.Lives)
		}
	}
}
