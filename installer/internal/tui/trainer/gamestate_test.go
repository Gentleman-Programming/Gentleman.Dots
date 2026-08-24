package trainer

import (
	"testing"
	"time"
)

// =============================================================================
// GAME STATE - Creation and Initialization
// =============================================================================

func TestNewGameState_Initialization(t *testing.T) {
	state := NewGameState()

	if state == nil {
		t.Fatal("NewGameState should not return nil")
	}
	if state.Stats == nil {
		t.Error("Stats should be initialized")
	}
	if state.CurrentStreak != 0 {
		t.Errorf("CurrentStreak should be 0, got %d", state.CurrentStreak)
	}
	if state.ComboMultiplier != 1 {
		t.Errorf("ComboMultiplier should be 1, got %d", state.ComboMultiplier)
	}
	if state.SessionScore != 0 {
		t.Errorf("SessionScore should be 0, got %d", state.SessionScore)
	}
}

func TestNewGameState_WithExistingStats(t *testing.T) {
	// Simulate loading existing stats
	stats := NewUserStats()
	stats.TotalScore = 1000
	stats.BestStreak = 10

	state := NewGameStateWithStats(stats)

	if state.Stats.TotalScore != 1000 {
		t.Errorf("Should preserve existing TotalScore, got %d", state.Stats.TotalScore)
	}
	if state.Stats.BestStreak != 10 {
		t.Errorf("Should preserve existing BestStreak, got %d", state.Stats.BestStreak)
	}
}

// =============================================================================
// GAME STATE - Starting Modes
// =============================================================================

func TestGameState_StartLesson(t *testing.T) {
	state := NewGameState()

	state.StartLesson(ModuleHorizontal)

	if state.CurrentModule != ModuleHorizontal {
		t.Errorf("CurrentModule should be Horizontal, got %s", state.CurrentModule)
	}
	if !state.IsLessonMode {
		t.Error("IsLessonMode should be true")
	}
	if state.IsPracticeMode {
		t.Error("IsPracticeMode should be false")
	}
	if state.IsBossMode {
		t.Error("IsBossMode should be false")
	}
	if state.ExerciseIndex != 0 {
		t.Errorf("ExerciseIndex should be 0, got %d", state.ExerciseIndex)
	}
	if len(state.Exercises) == 0 {
		t.Error("Exercises should be loaded")
	}
	if state.CurrentExercise == nil {
		t.Error("CurrentExercise should be set")
	}
}

func TestGameState_StartPractice(t *testing.T) {
	state := NewGameState()

	state.StartPractice(ModuleHorizontal)

	if state.CurrentModule != ModuleHorizontal {
		t.Errorf("CurrentModule should be Horizontal, got %s", state.CurrentModule)
	}
	if state.IsLessonMode {
		t.Error("IsLessonMode should be false")
	}
	if !state.IsPracticeMode {
		t.Error("IsPracticeMode should be true")
	}
	if state.IsBossMode {
		t.Error("IsBossMode should be false")
	}
	// In intelligent practice mode, we use random selection instead of loading all exercises
	// CurrentExercise should be set with a weighted random selection
	if state.CurrentExercise == nil {
		t.Error("CurrentExercise should be set for practice mode")
	}
}

func TestGameState_StartBoss(t *testing.T) {
	state := NewGameState()

	state.StartBoss(ModuleHorizontal)

	if !state.IsBossMode {
		t.Error("IsBossMode should be true")
	}
	if state.BossLives != 3 {
		t.Errorf("BossLives should be 3, got %d", state.BossLives)
	}
	if state.CurrentBoss == nil {
		t.Error("CurrentBoss should be set")
	}
	if state.BossStep != 0 {
		t.Errorf("BossStep should be 0, got %d", state.BossStep)
	}
}

// =============================================================================
// GAME STATE - Recording Answers
// =============================================================================

func TestGameState_RecordCorrectAnswer(t *testing.T) {
	state := NewGameState()
	state.StartLesson(ModuleHorizontal)
	initialScore := state.SessionScore

	state.RecordCorrectAnswer(5.0, true) // 5 seconds, optimal answer

	if state.CurrentStreak != 1 {
		t.Errorf("CurrentStreak should be 1 after correct answer, got %d", state.CurrentStreak)
	}
	if state.SessionScore <= initialScore {
		t.Error("SessionScore should increase after correct answer")
	}
}

func TestGameState_RecordCorrectAnswer_UpdatesStreak(t *testing.T) {
	state := NewGameState()
	state.StartLesson(ModuleHorizontal)

	// Answer correctly 5 times
	for i := 0; i < 5; i++ {
		state.RecordCorrectAnswer(3.0, false)
	}

	if state.CurrentStreak != 5 {
		t.Errorf("CurrentStreak should be 5, got %d", state.CurrentStreak)
	}
}

func TestGameState_RecordCorrectAnswer_UpdatesBestStreak(t *testing.T) {
	state := NewGameState()
	state.StartLesson(ModuleHorizontal)

	// Answer correctly 5 times
	for i := 0; i < 5; i++ {
		state.RecordCorrectAnswer(3.0, false)
	}

	if state.Stats.BestStreak < 5 {
		t.Errorf("BestStreak should be at least 5, got %d", state.Stats.BestStreak)
	}
}

func TestGameState_RecordIncorrectAnswer(t *testing.T) {
	state := NewGameState()
	state.StartLesson(ModuleHorizontal)
	state.CurrentStreak = 5

	state.RecordIncorrectAnswer()

	if state.CurrentStreak != 0 {
		t.Errorf("CurrentStreak should be 0 after incorrect answer, got %d", state.CurrentStreak)
	}
}

func TestGameState_RecordIncorrectAnswer_PreservesBestStreak(t *testing.T) {
	state := NewGameState()
	state.Stats.BestStreak = 10
	state.CurrentStreak = 5

	state.RecordIncorrectAnswer()

	if state.Stats.BestStreak != 10 {
		t.Errorf("BestStreak should be preserved at 10, got %d", state.Stats.BestStreak)
	}
}

func TestGameState_RecordIncorrectAnswer_InBossMode(t *testing.T) {
	state := NewGameState()
	state.StartBoss(ModuleHorizontal)
	initialLives := state.BossLives

	state.RecordIncorrectAnswer()

	if state.BossLives != initialLives-1 {
		t.Errorf("BossLives should decrease by 1, expected %d, got %d", initialLives-1, state.BossLives)
	}
}

func TestGameState_RecordIncorrectAnswer_BossDefeatsOnZeroLives(t *testing.T) {
	state := NewGameState()
	state.StartBoss(ModuleHorizontal)
	state.BossLives = 1

	state.RecordIncorrectAnswer()

	if state.BossLives != 0 {
		t.Errorf("BossLives should be 0, got %d", state.BossLives)
	}
	if !state.IsBossDefeated {
		t.Error("IsBossDefeated should be true when lives reach 0")
	}
}

// =============================================================================
// GAME STATE - Exercise Progression
// =============================================================================

func TestGameState_NextExercise_InLessonMode(t *testing.T) {
	state := NewGameState()
	state.StartLesson(ModuleHorizontal)

	hasMore := state.NextExercise()

	if state.ExerciseIndex != 1 {
		t.Errorf("ExerciseIndex should be 1, got %d", state.ExerciseIndex)
	}
	if !hasMore && state.ExerciseIndex < len(state.Exercises) {
		t.Error("NextExercise should return true when more exercises exist")
	}
}

func TestGameState_NextExercise_CompletesLessons(t *testing.T) {
	state := NewGameState()
	state.StartLesson(ModuleHorizontal)

	// Advance through all exercises
	for state.NextExercise() {
		// Keep going
	}

	// Check lessons are marked complete in stats
	progress := state.Stats.GetModuleProgress(ModuleHorizontal)
	if progress.LessonsCompleted != progress.LessonsTotal {
		t.Errorf("LessonsCompleted should equal LessonsTotal after completing all")
	}
}

func TestGameState_NextExercise_InBossMode(t *testing.T) {
	state := NewGameState()
	state.StartBoss(ModuleHorizontal)
	initialStep := state.BossStep

	hasMore := state.NextExercise()

	if state.BossStep != initialStep+1 {
		t.Errorf("BossStep should increase by 1")
	}
	if !hasMore && state.BossStep < len(state.CurrentBoss.Steps) {
		t.Error("NextExercise should return true when more boss steps exist")
	}
}

func TestGameState_NextExercise_IncreasesCombo(t *testing.T) {
	state := NewGameState()
	state.StartBoss(ModuleHorizontal)
	state.ComboMultiplier = 2

	state.NextExercise()

	if state.ComboMultiplier != 3 {
		t.Errorf("ComboMultiplier should increase to 3, got %d", state.ComboMultiplier)
	}
}

func TestGameState_NextExercise_ComboMaxAt4(t *testing.T) {
	state := NewGameState()
	state.StartBoss(ModuleHorizontal)
	state.ComboMultiplier = 4

	state.NextExercise()

	if state.ComboMultiplier != 4 {
		t.Errorf("ComboMultiplier should max at 4, got %d", state.ComboMultiplier)
	}
}

// =============================================================================
// GAME STATE - Stats Updates
// =============================================================================

func TestGameState_UpdatePracticeStats(t *testing.T) {
	state := NewGameState()
	state.StartPractice(ModuleHorizontal)

	// Record some practice attempts
	state.RecordCorrectAnswer(3.0, false)
	state.RecordCorrectAnswer(4.0, true)
	state.RecordIncorrectAnswer()

	progress := state.Stats.GetModuleProgress(ModuleHorizontal)

	if progress.PracticeAttempts != 3 {
		t.Errorf("PracticeAttempts should be 3, got %d", progress.PracticeAttempts)
	}
	if progress.PracticeCorrect != 2 {
		t.Errorf("PracticeCorrect should be 2, got %d", progress.PracticeCorrect)
	}
}

func TestGameState_PracticeAccuracyCalculation(t *testing.T) {
	state := NewGameState()
	state.StartPractice(ModuleHorizontal)

	// 8 correct, 2 incorrect = 80% accuracy
	for i := 0; i < 8; i++ {
		state.RecordCorrectAnswer(3.0, false)
	}
	for i := 0; i < 2; i++ {
		state.RecordIncorrectAnswer()
	}

	progress := state.Stats.GetModuleProgress(ModuleHorizontal)
	expectedAccuracy := 0.80

	if progress.PracticeAccuracy < expectedAccuracy-0.01 || progress.PracticeAccuracy > expectedAccuracy+0.01 {
		t.Errorf("PracticeAccuracy should be ~0.80, got %f", progress.PracticeAccuracy)
	}
}

func TestGameState_RecordBossVictory(t *testing.T) {
	state := NewGameState()
	state.StartBoss(ModuleHorizontal)
	state.BossLives = 2
	state.TimeElapsed = 25 * time.Second

	state.RecordBossVictory()

	if !state.Stats.IsBossDefeated(ModuleHorizontal) {
		t.Error("Boss should be marked as defeated")
	}

	progress := state.Stats.GetModuleProgress(ModuleHorizontal)
	if !progress.BossDefeated {
		t.Error("Module progress should show boss defeated")
	}
	if progress.BossBestTime != 25*time.Second {
		t.Errorf("BossBestTime should be 25s, got %v", progress.BossBestTime)
	}
}

func TestGameState_RecordBossVictory_AddsToDefeatedList(t *testing.T) {
	state := NewGameState()
	state.StartBoss(ModuleHorizontal)

	state.RecordBossVictory()

	found := false
	for _, boss := range state.Stats.BossesDefeated {
		if boss == ModuleHorizontal {
			found = true
			break
		}
	}
	if !found {
		t.Error("Horizontal should be in BossesDefeated list")
	}
}

func TestGameState_RecordBossVictory_DoesntDuplicateInList(t *testing.T) {
	state := NewGameState()
	state.Stats.BossesDefeated = []ModuleID{ModuleHorizontal} // Already defeated
	state.StartBoss(ModuleHorizontal)

	state.RecordBossVictory()

	count := 0
	for _, boss := range state.Stats.BossesDefeated {
		if boss == ModuleHorizontal {
			count++
		}
	}
	if count != 1 {
		t.Errorf("Horizontal should only appear once in BossesDefeated, got %d", count)
	}
}

// =============================================================================
// GAME STATE - Session Management
// =============================================================================

func TestGameState_Reset(t *testing.T) {
	state := NewGameState()
	state.StartLesson(ModuleHorizontal)
	state.CurrentStreak = 10
	state.SessionScore = 500
	state.ComboMultiplier = 3

	state.Reset()

	if state.CurrentStreak != 0 {
		t.Errorf("CurrentStreak should be 0 after reset, got %d", state.CurrentStreak)
	}
	if state.SessionScore != 0 {
		t.Errorf("SessionScore should be 0 after reset, got %d", state.SessionScore)
	}
	if state.ComboMultiplier != 1 {
		t.Errorf("ComboMultiplier should be 1 after reset, got %d", state.ComboMultiplier)
	}
	if state.IsLessonMode || state.IsPracticeMode || state.IsBossMode {
		t.Error("All modes should be false after reset")
	}
}

// =============================================================================
// GAME STATE - Intelligent Practice Mode
// =============================================================================

func TestGameState_SetPracticeExercise(t *testing.T) {
	state := NewGameState()
	exercise := &Exercise{
		ID:      "test-exercise",
		Mission: "Test mission",
	}

	state.SetPracticeExercise(exercise)

	if state.CurrentExercise != exercise {
		t.Error("SetPracticeExercise should set the current exercise")
	}
	if state.CurrentExercise.ID != "test-exercise" {
		t.Errorf("Expected exercise ID 'test-exercise', got '%s'", state.CurrentExercise.ID)
	}
}

func TestGameState_NextPracticeExercise_ReturnsTrue(t *testing.T) {
	state := NewGameState()
	state.StartPractice(ModuleHorizontal)

	// Should return true when there are unmastered exercises
	hasNext := state.NextPracticeExercise()

	if !hasNext {
		t.Error("NextPracticeExercise should return true when exercises remain")
	}
	if state.CurrentExercise == nil {
		t.Error("NextPracticeExercise should set a new exercise")
	}
}

func TestGameState_NextPracticeExercise_ReturnsFalseWhenNotPracticeMode(t *testing.T) {
	state := NewGameState()
	state.StartLesson(ModuleHorizontal)

	hasNext := state.NextPracticeExercise()

	if hasNext {
		t.Error("NextPracticeExercise should return false when not in practice mode")
	}
}

func TestGameState_NextPracticeExercise_ReturnsFalseWhenAllMastered(t *testing.T) {
	state := NewGameState()
	state.StartPractice(ModuleHorizontal)

	// Mark all exercises as mastered
	progress := state.Stats.GetModuleProgress(ModuleHorizontal)
	lessons := GetLessons(ModuleHorizontal)
	for _, lesson := range lessons {
		exStats := progress.GetExerciseStats(lesson.ID)
		exStats.Mastered = true
	}

	hasNext := state.NextPracticeExercise()

	if hasNext {
		t.Error("NextPracticeExercise should return false when all exercises are mastered")
	}
}
