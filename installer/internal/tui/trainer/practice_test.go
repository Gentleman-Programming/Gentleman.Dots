package trainer

import (
	"testing"
)

// =============================================================================
// EXERCISE STATS
// =============================================================================

func TestGetExerciseStats_CreatesNew(t *testing.T) {
	mp := &ModuleProgress{}

	stats := mp.GetExerciseStats("test_001")

	if stats == nil {
		t.Fatal("GetExerciseStats should create new stats")
	}
	if stats.TotalAttempts != 0 {
		t.Errorf("New stats should have 0 attempts, got %d", stats.TotalAttempts)
	}
}

func TestGetExerciseStats_ReturnsSame(t *testing.T) {
	mp := &ModuleProgress{}

	stats1 := mp.GetExerciseStats("test_001")
	stats1.TotalAttempts = 5

	stats2 := mp.GetExerciseStats("test_001")

	if stats2.TotalAttempts != 5 {
		t.Errorf("Should return same stats, got %d attempts", stats2.TotalAttempts)
	}
}

// =============================================================================
// RECORD PRACTICE RESULT
// =============================================================================

func TestRecordPracticeResult_CorrectIncrementsCounters(t *testing.T) {
	mp := &ModuleProgress{}

	mp.RecordPracticeResult("test_001", true)

	stats := mp.GetExerciseStats("test_001")
	if stats.TotalAttempts != 1 {
		t.Errorf("TotalAttempts should be 1, got %d", stats.TotalAttempts)
	}
	if stats.TotalCorrect != 1 {
		t.Errorf("TotalCorrect should be 1, got %d", stats.TotalCorrect)
	}
	if stats.ConsecutiveRight != 1 {
		t.Errorf("ConsecutiveRight should be 1, got %d", stats.ConsecutiveRight)
	}
}

func TestRecordPracticeResult_WrongIncrementsCounters(t *testing.T) {
	mp := &ModuleProgress{}

	mp.RecordPracticeResult("test_001", false)

	stats := mp.GetExerciseStats("test_001")
	if stats.TotalAttempts != 1 {
		t.Errorf("TotalAttempts should be 1, got %d", stats.TotalAttempts)
	}
	if stats.TotalWrong != 1 {
		t.Errorf("TotalWrong should be 1, got %d", stats.TotalWrong)
	}
	if stats.ConsecutiveRight != 0 {
		t.Errorf("ConsecutiveRight should be 0, got %d", stats.ConsecutiveRight)
	}
}

func TestRecordPracticeResult_WrongResetsStreak(t *testing.T) {
	mp := &ModuleProgress{}

	// Build up a streak
	mp.RecordPracticeResult("test_001", true)
	mp.RecordPracticeResult("test_001", true)

	stats := mp.GetExerciseStats("test_001")
	if stats.ConsecutiveRight != 2 {
		t.Errorf("Should have 2 consecutive right, got %d", stats.ConsecutiveRight)
	}

	// Wrong answer resets streak
	mp.RecordPracticeResult("test_001", false)

	if stats.ConsecutiveRight != 0 {
		t.Errorf("Wrong should reset streak to 0, got %d", stats.ConsecutiveRight)
	}
}

func TestRecordPracticeResult_MasteryAfterThreeCorrect(t *testing.T) {
	mp := &ModuleProgress{}

	// Not mastered after 2
	mp.RecordPracticeResult("test_001", true)
	mp.RecordPracticeResult("test_001", true)

	stats := mp.GetExerciseStats("test_001")
	if stats.Mastered {
		t.Error("Should not be mastered after 2 correct")
	}

	// Mastered after 3
	mp.RecordPracticeResult("test_001", true)

	if !stats.Mastered {
		t.Error("Should be mastered after 3 consecutive correct")
	}
}

func TestRecordPracticeResult_WrongUnmasters(t *testing.T) {
	mp := &ModuleProgress{}

	// Master the exercise
	mp.RecordPracticeResult("test_001", true)
	mp.RecordPracticeResult("test_001", true)
	mp.RecordPracticeResult("test_001", true)

	stats := mp.GetExerciseStats("test_001")
	if !stats.Mastered {
		t.Fatal("Should be mastered")
	}

	// Wrong answer un-masters
	mp.RecordPracticeResult("test_001", false)

	if stats.Mastered {
		t.Error("Wrong answer should un-master")
	}
}

func TestRecordPracticeResult_UpdatesModuleProgress(t *testing.T) {
	mp := &ModuleProgress{}

	mp.RecordPracticeResult("test_001", true)
	mp.RecordPracticeResult("test_001", false)

	if mp.PracticeAttempts != 2 {
		t.Errorf("PracticeAttempts should be 2, got %d", mp.PracticeAttempts)
	}
	if mp.PracticeCorrect != 1 {
		t.Errorf("PracticeCorrect should be 1, got %d", mp.PracticeCorrect)
	}
	if mp.PracticeAccuracy != 0.5 {
		t.Errorf("PracticeAccuracy should be 0.5, got %f", mp.PracticeAccuracy)
	}
}

// =============================================================================
// PRACTICE WEIGHT
// =============================================================================

func TestGetPracticeWeight_MasteredIsZero(t *testing.T) {
	stats := &ExerciseStats{Mastered: true}

	if stats.GetPracticeWeight() != 0 {
		t.Error("Mastered exercise should have weight 0")
	}
}

func TestGetPracticeWeight_NeverAttemptedHasMediumWeight(t *testing.T) {
	stats := &ExerciseStats{}

	weight := stats.GetPracticeWeight()
	if weight != 15 {
		t.Errorf("Never attempted should have weight 15, got %d", weight)
	}
}

func TestGetPracticeWeight_MoreErrorsHigherWeight(t *testing.T) {
	statsLowErrors := &ExerciseStats{TotalAttempts: 5, TotalWrong: 1}
	statsHighErrors := &ExerciseStats{TotalAttempts: 5, TotalWrong: 5}

	if statsHighErrors.GetPracticeWeight() <= statsLowErrors.GetPracticeWeight() {
		t.Error("More errors should mean higher weight")
	}
}

func TestGetPracticeWeight_ConsecutiveCorrectReducesWeight(t *testing.T) {
	statsNoStreak := &ExerciseStats{TotalAttempts: 5, ConsecutiveRight: 0}
	statsStreak := &ExerciseStats{TotalAttempts: 5, ConsecutiveRight: 2}

	if statsStreak.GetPracticeWeight() >= statsNoStreak.GetPracticeWeight() {
		t.Error("Consecutive correct should reduce weight")
	}
}

// =============================================================================
// WEIGHTED PRACTICE EXERCISES
// =============================================================================

func TestGetWeightedPracticeExercises_ReturnsExercises(t *testing.T) {
	mp := &ModuleProgress{}

	exercises := GetWeightedPracticeExercises(ModuleHorizontal, mp)

	if len(exercises) == 0 {
		t.Error("Should return exercises for horizontal module")
	}
}

func TestGetWeightedPracticeExercises_ExcludesMastered(t *testing.T) {
	mp := &ModuleProgress{}

	lessons := GetLessons(ModuleHorizontal)
	if len(lessons) == 0 {
		t.Skip("No lessons available")
	}

	// Master the first exercise
	mp.RecordPracticeResult(lessons[0].ID, true)
	mp.RecordPracticeResult(lessons[0].ID, true)
	mp.RecordPracticeResult(lessons[0].ID, true)

	exercises := GetWeightedPracticeExercises(ModuleHorizontal, mp)

	// Should have one less than total
	if len(exercises) != len(lessons)-1 {
		t.Errorf("Should have %d exercises (excluding mastered), got %d", len(lessons)-1, len(exercises))
	}

	// First exercise should not be in list
	for _, ex := range exercises {
		if ex.ID == lessons[0].ID {
			t.Error("Mastered exercise should not be in practice list")
		}
	}
}

func TestGetWeightedPracticeExercises_EmptyWhenAllMastered(t *testing.T) {
	mp := &ModuleProgress{}

	lessons := GetLessons(ModuleHorizontal)
	if len(lessons) == 0 {
		t.Skip("No lessons available")
	}

	// Master all exercises
	for _, lesson := range lessons {
		mp.RecordPracticeResult(lesson.ID, true)
		mp.RecordPracticeResult(lesson.ID, true)
		mp.RecordPracticeResult(lesson.ID, true)
	}

	exercises := GetWeightedPracticeExercises(ModuleHorizontal, mp)

	if len(exercises) != 0 {
		t.Errorf("All mastered should return empty list, got %d", len(exercises))
	}
}

func TestGetWeightedPracticeExercises_ArePracticeType(t *testing.T) {
	mp := &ModuleProgress{}

	exercises := GetWeightedPracticeExercises(ModuleHorizontal, mp)

	for _, ex := range exercises {
		if ex.Type != ExercisePractice {
			t.Errorf("Exercise should be Practice type, got %s", ex.Type)
		}
	}
}

// =============================================================================
// SELECT RANDOM PRACTICE EXERCISE
// =============================================================================

func TestSelectRandomPracticeExercise_ReturnsExercise(t *testing.T) {
	mp := &ModuleProgress{}

	ex := SelectRandomPracticeExercise(ModuleHorizontal, mp)

	if ex == nil {
		t.Error("Should return an exercise")
	}
}

func TestSelectRandomPracticeExercise_NilWhenAllMastered(t *testing.T) {
	mp := &ModuleProgress{}

	lessons := GetLessons(ModuleHorizontal)
	if len(lessons) == 0 {
		t.Skip("No lessons available")
	}

	// Master all
	for _, lesson := range lessons {
		mp.RecordPracticeResult(lesson.ID, true)
		mp.RecordPracticeResult(lesson.ID, true)
		mp.RecordPracticeResult(lesson.ID, true)
	}

	ex := SelectRandomPracticeExercise(ModuleHorizontal, mp)

	if ex != nil {
		t.Error("Should return nil when all mastered")
	}
}

func TestSelectRandomPracticeExercise_IsPracticeType(t *testing.T) {
	mp := &ModuleProgress{}

	ex := SelectRandomPracticeExercise(ModuleHorizontal, mp)

	if ex != nil && ex.Type != ExercisePractice {
		t.Errorf("Should be Practice type, got %s", ex.Type)
	}
}

// =============================================================================
// PRACTICE STATS
// =============================================================================

func TestGetPracticeStatsForModule_CountsCorrectly(t *testing.T) {
	mp := &ModuleProgress{}

	lessons := GetLessons(ModuleHorizontal)
	if len(lessons) < 3 {
		t.Skip("Need at least 3 lessons")
	}

	// Master 2 exercises
	for i := 0; i < 2; i++ {
		mp.RecordPracticeResult(lessons[i].ID, true)
		mp.RecordPracticeResult(lessons[i].ID, true)
		mp.RecordPracticeResult(lessons[i].ID, true)
	}

	stats := GetPracticeStatsForModule(ModuleHorizontal, mp)

	if stats.TotalExercises != len(lessons) {
		t.Errorf("TotalExercises should be %d, got %d", len(lessons), stats.TotalExercises)
	}
	if stats.MasteredCount != 2 {
		t.Errorf("MasteredCount should be 2, got %d", stats.MasteredCount)
	}
	if stats.RemainingCount != len(lessons)-2 {
		t.Errorf("RemainingCount should be %d, got %d", len(lessons)-2, stats.RemainingCount)
	}
}

func TestGetPracticeStatsForModule_TracksWeakest(t *testing.T) {
	mp := &ModuleProgress{}

	lessons := GetLessons(ModuleHorizontal)
	if len(lessons) < 3 {
		t.Skip("Need at least 3 lessons")
	}

	// Make first exercise have most errors
	mp.RecordPracticeResult(lessons[0].ID, false)
	mp.RecordPracticeResult(lessons[0].ID, false)
	mp.RecordPracticeResult(lessons[0].ID, false)

	// Second exercise has some errors
	mp.RecordPracticeResult(lessons[1].ID, false)

	stats := GetPracticeStatsForModule(ModuleHorizontal, mp)

	if len(stats.WeakestExercises) == 0 {
		t.Fatal("Should have weakest exercises")
	}
	if stats.WeakestExercises[0] != lessons[0].ID {
		t.Errorf("First exercise should be weakest, got %s", stats.WeakestExercises[0])
	}
}

func TestGetPracticeStatsForModule_PracticeComplete(t *testing.T) {
	mp := &ModuleProgress{}

	lessons := GetLessons(ModuleHorizontal)
	if len(lessons) == 0 {
		t.Skip("No lessons")
	}

	// Master all
	for _, lesson := range lessons {
		mp.RecordPracticeResult(lesson.ID, true)
		mp.RecordPracticeResult(lesson.ID, true)
		mp.RecordPracticeResult(lesson.ID, true)
	}

	stats := GetPracticeStatsForModule(ModuleHorizontal, mp)

	if !stats.PracticeComplete {
		t.Error("PracticeComplete should be true when all mastered")
	}
}

// =============================================================================
// RESET MODULE PRACTICE
// =============================================================================

func TestResetModulePractice_ClearsAll(t *testing.T) {
	mp := &ModuleProgress{
		PracticeAttempts: 100,
		PracticeCorrect:  50,
		PracticeAccuracy: 0.5,
	}

	// Add some exercise stats
	mp.RecordPracticeResult("test_001", true)
	mp.RecordPracticeResult("test_001", true)

	mp.ResetModulePractice()

	if mp.PracticeAttempts != 0 {
		t.Errorf("PracticeAttempts should be 0, got %d", mp.PracticeAttempts)
	}
	if mp.PracticeCorrect != 0 {
		t.Errorf("PracticeCorrect should be 0, got %d", mp.PracticeCorrect)
	}
	if mp.PracticeAccuracy != 0 {
		t.Errorf("PracticeAccuracy should be 0, got %f", mp.PracticeAccuracy)
	}
	if len(mp.ExerciseStats) != 0 {
		t.Errorf("ExerciseStats should be empty, got %d", len(mp.ExerciseStats))
	}
}

// =============================================================================
// IS PRACTICE COMPLETE
// =============================================================================

func TestIsPracticeComplete_FalseWhenNotAllMastered(t *testing.T) {
	mp := &ModuleProgress{}

	if mp.IsPracticeComplete(ModuleHorizontal) {
		t.Error("Should not be complete with no practice")
	}
}

func TestIsPracticeComplete_TrueWhenAllMastered(t *testing.T) {
	mp := &ModuleProgress{}

	lessons := GetLessons(ModuleHorizontal)
	if len(lessons) == 0 {
		t.Skip("No lessons")
	}

	// Master all
	for _, lesson := range lessons {
		mp.RecordPracticeResult(lesson.ID, true)
		mp.RecordPracticeResult(lesson.ID, true)
		mp.RecordPracticeResult(lesson.ID, true)
	}

	if !mp.IsPracticeComplete(ModuleHorizontal) {
		t.Error("Should be complete when all mastered")
	}
}
