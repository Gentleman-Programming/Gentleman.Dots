package trainer

import (
	"testing"
)

// =============================================================================
// MODULE UNLOCK LOGIC
// =============================================================================

func TestIsModuleUnlocked_HorizontalAlwaysUnlocked(t *testing.T) {
	stats := NewUserStats()

	if !stats.IsModuleUnlocked(ModuleHorizontal) {
		t.Error("Horizontal module should ALWAYS be unlocked")
	}
}

func TestIsModuleUnlocked_SecondModuleRequiresBossDefeat(t *testing.T) {
	stats := NewUserStats()

	// Vertical should be locked initially
	if stats.IsModuleUnlocked(ModuleVertical) {
		t.Error("Vertical should be locked without defeating Horizontal boss")
	}

	// Defeat Horizontal boss
	stats.BossesDefeated = append(stats.BossesDefeated, ModuleHorizontal)

	// Now Vertical should be unlocked
	if !stats.IsModuleUnlocked(ModuleVertical) {
		t.Error("Vertical should be unlocked after defeating Horizontal boss")
	}
}

func TestIsModuleUnlocked_AllModulesInOrder(t *testing.T) {
	stats := NewUserStats()

	moduleOrder := []ModuleID{
		ModuleHorizontal,
		ModuleVertical,
		ModuleTextObjects,
		ModuleChangeRepeat,
		ModuleSubstitution,
		ModuleRegex,
		ModuleMacros,
	}

	// Initially only first is unlocked
	for i, mod := range moduleOrder {
		unlocked := stats.IsModuleUnlocked(mod)
		if i == 0 && !unlocked {
			t.Errorf("Module %s (index 0) should be unlocked", mod)
		}
		if i > 0 && unlocked {
			t.Errorf("Module %s (index %d) should be locked initially", mod, i)
		}
	}

	// Unlock progressively by defeating bosses
	for i := 0; i < len(moduleOrder)-1; i++ {
		stats.BossesDefeated = append(stats.BossesDefeated, moduleOrder[i])

		// Next module should now be unlocked
		nextMod := moduleOrder[i+1]
		if !stats.IsModuleUnlocked(nextMod) {
			t.Errorf("Module %s should be unlocked after defeating %s", nextMod, moduleOrder[i])
		}

		// Module after next should still be locked (if exists)
		if i+2 < len(moduleOrder) {
			afterNext := moduleOrder[i+2]
			if stats.IsModuleUnlocked(afterNext) {
				t.Errorf("Module %s should still be locked", afterNext)
			}
		}
	}
}

func TestIsBossDefeated_ReturnsFalseInitially(t *testing.T) {
	stats := NewUserStats()

	for _, mod := range GetAllModules() {
		if stats.IsBossDefeated(mod.ID) {
			t.Errorf("Boss %s should not be defeated initially", mod.ID)
		}
	}
}

func TestIsBossDefeated_ReturnsTrueWhenInList(t *testing.T) {
	stats := NewUserStats()
	stats.BossesDefeated = []ModuleID{ModuleHorizontal, ModuleTextObjects}

	if !stats.IsBossDefeated(ModuleHorizontal) {
		t.Error("Horizontal boss should be marked as defeated")
	}
	if !stats.IsBossDefeated(ModuleTextObjects) {
		t.Error("TextObjects boss should be marked as defeated")
	}
	if stats.IsBossDefeated(ModuleVertical) {
		t.Error("Vertical boss should NOT be marked as defeated")
	}
}

// =============================================================================
// LESSON COMPLETION LOGIC
// =============================================================================

func TestIsLessonsComplete_FalseWhenNoLessons(t *testing.T) {
	stats := NewUserStats()

	if stats.IsLessonsComplete(ModuleHorizontal) {
		t.Error("Lessons should not be complete when total is 0")
	}
}

func TestIsLessonsComplete_FalseWhenPartial(t *testing.T) {
	stats := NewUserStats()
	progress := stats.GetModuleProgress(ModuleHorizontal)
	progress.LessonsTotal = 15
	progress.LessonsCompleted = 10

	if stats.IsLessonsComplete(ModuleHorizontal) {
		t.Error("Lessons should not be complete when only partially done")
	}
}

func TestIsLessonsComplete_TrueWhenAllDone(t *testing.T) {
	stats := NewUserStats()
	progress := stats.GetModuleProgress(ModuleHorizontal)
	progress.LessonsTotal = 15
	progress.LessonsCompleted = 15

	if !stats.IsLessonsComplete(ModuleHorizontal) {
		t.Error("Lessons should be complete when all done")
	}
}

func TestIsLessonsComplete_TrueWhenMoreThanTotal(t *testing.T) {
	// Edge case: completed > total (shouldn't happen, but be safe)
	stats := NewUserStats()
	progress := stats.GetModuleProgress(ModuleHorizontal)
	progress.LessonsTotal = 15
	progress.LessonsCompleted = 20

	if !stats.IsLessonsComplete(ModuleHorizontal) {
		t.Error("Lessons should be complete when completed >= total")
	}
}

// =============================================================================
// PRACTICE UNLOCK LOGIC
// =============================================================================

func TestIsPracticeReady_RequiresModuleUnlocked(t *testing.T) {
	stats := NewUserStats()

	// Vertical is locked, even with lessons complete, practice not ready
	progress := stats.GetModuleProgress(ModuleVertical)
	progress.LessonsTotal = 10
	progress.LessonsCompleted = 10

	if stats.IsPracticeReady(ModuleVertical) {
		t.Error("Practice should not be ready for locked module")
	}
}

func TestIsPracticeReady_RequiresLessonsComplete(t *testing.T) {
	stats := NewUserStats()

	// Horizontal is unlocked but lessons not complete
	progress := stats.GetModuleProgress(ModuleHorizontal)
	progress.LessonsTotal = 15
	progress.LessonsCompleted = 10

	if stats.IsPracticeReady(ModuleHorizontal) {
		t.Error("Practice should not be ready without completing lessons")
	}
}

func TestIsPracticeReady_TrueWhenConditionsMet(t *testing.T) {
	stats := NewUserStats()

	progress := stats.GetModuleProgress(ModuleHorizontal)
	progress.LessonsTotal = 15
	progress.LessonsCompleted = 15

	if !stats.IsPracticeReady(ModuleHorizontal) {
		t.Error("Practice should be ready for unlocked module with complete lessons")
	}
}

// =============================================================================
// BOSS UNLOCK LOGIC
// =============================================================================

func TestIsBossReady_RequiresPracticeReady(t *testing.T) {
	stats := NewUserStats()

	// Module unlocked, lessons incomplete
	progress := stats.GetModuleProgress(ModuleHorizontal)
	progress.LessonsTotal = 15
	progress.LessonsCompleted = 10

	if stats.IsBossReady(ModuleHorizontal) {
		t.Error("Boss should not be ready without practice being ready")
	}
}

func TestIsBossReady_RequiresMinimumAttempts(t *testing.T) {
	stats := NewUserStats()

	progress := stats.GetModuleProgress(ModuleHorizontal)
	progress.LessonsTotal = 15
	progress.LessonsCompleted = 15
	progress.PracticeAccuracy = 0.90 // High accuracy but...
	progress.PracticeAttempts = 5    // Not enough attempts

	if stats.IsBossReady(ModuleHorizontal) {
		t.Error("Boss should not be ready without minimum practice attempts (10)")
	}
}

func TestIsBossReady_Requires80PercentAccuracy(t *testing.T) {
	stats := NewUserStats()

	progress := stats.GetModuleProgress(ModuleHorizontal)
	progress.LessonsTotal = 15
	progress.LessonsCompleted = 15
	progress.PracticeAttempts = 20
	progress.PracticeAccuracy = 0.70 // Only 70%

	if stats.IsBossReady(ModuleHorizontal) {
		t.Error("Boss should not be ready without 80% accuracy")
	}
}

func TestIsBossReady_TrueWhenAllConditionsMet(t *testing.T) {
	stats := NewUserStats()

	progress := stats.GetModuleProgress(ModuleHorizontal)
	progress.LessonsTotal = 15
	progress.LessonsCompleted = 15
	progress.PracticeAttempts = 20
	progress.PracticeAccuracy = 0.85 // 85% >= 80%

	if !stats.IsBossReady(ModuleHorizontal) {
		t.Error("Boss should be ready with all conditions met")
	}
}

func TestIsBossReady_Exactly80PercentWorks(t *testing.T) {
	stats := NewUserStats()

	progress := stats.GetModuleProgress(ModuleHorizontal)
	progress.LessonsTotal = 15
	progress.LessonsCompleted = 15
	progress.PracticeAttempts = 10
	progress.PracticeAccuracy = 0.80 // Exactly 80%

	if !stats.IsBossReady(ModuleHorizontal) {
		t.Error("Boss should be ready with exactly 80% accuracy")
	}
}

func TestIsBossReady_Exactly10AttemptsWorks(t *testing.T) {
	stats := NewUserStats()

	progress := stats.GetModuleProgress(ModuleHorizontal)
	progress.LessonsTotal = 15
	progress.LessonsCompleted = 15
	progress.PracticeAttempts = 10 // Exactly minimum
	progress.PracticeAccuracy = 0.80

	if !stats.IsBossReady(ModuleHorizontal) {
		t.Error("Boss should be ready with exactly 10 attempts")
	}
}

func TestIsBossReady_9AttemptsNotEnough(t *testing.T) {
	stats := NewUserStats()

	progress := stats.GetModuleProgress(ModuleHorizontal)
	progress.LessonsTotal = 15
	progress.LessonsCompleted = 15
	progress.PracticeAttempts = 9 // One less than minimum
	progress.PracticeAccuracy = 0.95

	if stats.IsBossReady(ModuleHorizontal) {
		t.Error("Boss should NOT be ready with only 9 attempts")
	}
}

func TestIsBossReady_79PercentNotEnough(t *testing.T) {
	stats := NewUserStats()

	progress := stats.GetModuleProgress(ModuleHorizontal)
	progress.LessonsTotal = 15
	progress.LessonsCompleted = 15
	progress.PracticeAttempts = 100
	progress.PracticeAccuracy = 0.79 // Just under 80%

	if stats.IsBossReady(ModuleHorizontal) {
		t.Error("Boss should NOT be ready with 79% accuracy")
	}
}
