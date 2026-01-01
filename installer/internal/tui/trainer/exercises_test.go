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
