package trainer

import (
	"testing"
	"time"
)

// =============================================================================
// MODULE PROGRESS
// =============================================================================

func TestModuleProgress_ZeroValue(t *testing.T) {
	var progress ModuleProgress

	if progress.LessonsCompleted != 0 {
		t.Errorf("LessonsCompleted should be 0, got %d", progress.LessonsCompleted)
	}
	if progress.PracticeAccuracy != 0.0 {
		t.Errorf("PracticeAccuracy should be 0.0, got %f", progress.PracticeAccuracy)
	}
	if progress.BossDefeated != false {
		t.Error("BossDefeated should be false")
	}
}

func TestModuleProgress_LessonsTracking(t *testing.T) {
	progress := ModuleProgress{
		LessonsCompleted: 10,
		LessonsTotal:     15,
	}

	if progress.LessonsCompleted != 10 {
		t.Errorf("LessonsCompleted: expected 10, got %d", progress.LessonsCompleted)
	}
	if progress.LessonsTotal != 15 {
		t.Errorf("LessonsTotal: expected 15, got %d", progress.LessonsTotal)
	}
}

func TestModuleProgress_PracticeAccuracy(t *testing.T) {
	progress := ModuleProgress{
		PracticeAttempts: 100,
		PracticeCorrect:  85,
		PracticeAccuracy: 0.85,
	}

	if progress.PracticeAccuracy != 0.85 {
		t.Errorf("PracticeAccuracy: expected 0.85, got %f", progress.PracticeAccuracy)
	}
}

func TestModuleProgress_BossStats(t *testing.T) {
	progress := ModuleProgress{
		BossDefeated:  true,
		BossBestTime:  28 * time.Second,
		BossAttempts:  3,
		BossLivesLeft: 2,
	}

	if !progress.BossDefeated {
		t.Error("BossDefeated should be true")
	}
	if progress.BossBestTime != 28*time.Second {
		t.Errorf("BossBestTime: expected 28s, got %v", progress.BossBestTime)
	}
	if progress.BossLivesLeft != 2 {
		t.Errorf("BossLivesLeft: expected 2, got %d", progress.BossLivesLeft)
	}
}

func TestModuleProgress_WeakExercises(t *testing.T) {
	progress := ModuleProgress{
		WeakExercises: []string{"horizontal_008", "horizontal_012"},
	}

	if len(progress.WeakExercises) != 2 {
		t.Errorf("WeakExercises: expected 2, got %d", len(progress.WeakExercises))
	}
	if progress.WeakExercises[0] != "horizontal_008" {
		t.Errorf("First weak exercise: expected 'horizontal_008', got %q", progress.WeakExercises[0])
	}
}

// =============================================================================
// USER STATS
// =============================================================================

func TestNewUserStats_CreatesEmptyStats(t *testing.T) {
	stats := NewUserStats()

	if stats == nil {
		t.Fatal("NewUserStats should not return nil")
	}
	if stats.TotalScore != 0 {
		t.Errorf("TotalScore should be 0, got %d", stats.TotalScore)
	}
	if stats.CurrentStreak != 0 {
		t.Errorf("CurrentStreak should be 0, got %d", stats.CurrentStreak)
	}
	if stats.BestStreak != 0 {
		t.Errorf("BestStreak should be 0, got %d", stats.BestStreak)
	}
	if stats.ModuleProgress == nil {
		t.Error("ModuleProgress map should be initialized")
	}
	if len(stats.BossesDefeated) != 0 {
		t.Errorf("BossesDefeated should be empty, got %d", len(stats.BossesDefeated))
	}
}

func TestUserStats_GetModuleProgress_CreatesIfNotExists(t *testing.T) {
	stats := NewUserStats()

	progress := stats.GetModuleProgress(ModuleHorizontal)

	if progress == nil {
		t.Fatal("GetModuleProgress should not return nil")
	}

	// Should return same instance on second call
	progress2 := stats.GetModuleProgress(ModuleHorizontal)
	if progress != progress2 {
		t.Error("GetModuleProgress should return same instance")
	}
}

func TestUserStats_GetModuleProgress_InitializesMapIfNil(t *testing.T) {
	stats := &UserStats{
		ModuleProgress: nil, // Explicitly nil
	}

	progress := stats.GetModuleProgress(ModuleHorizontal)

	if progress == nil {
		t.Fatal("Should create progress even with nil map")
	}
	if stats.ModuleProgress == nil {
		t.Error("Should initialize ModuleProgress map")
	}
}

func TestUserStats_TrackScore(t *testing.T) {
	stats := NewUserStats()

	stats.TotalScore = 2340
	stats.CurrentStreak = 7
	stats.BestStreak = 23

	if stats.TotalScore != 2340 {
		t.Errorf("TotalScore: expected 2340, got %d", stats.TotalScore)
	}
	if stats.CurrentStreak != 7 {
		t.Errorf("CurrentStreak: expected 7, got %d", stats.CurrentStreak)
	}
	if stats.BestStreak != 23 {
		t.Errorf("BestStreak: expected 23, got %d", stats.BestStreak)
	}
}

func TestUserStats_TrackTotalTime(t *testing.T) {
	stats := NewUserStats()

	stats.TotalTime = 2*time.Hour + 18*time.Minute

	if stats.TotalTime != 2*time.Hour+18*time.Minute {
		t.Errorf("TotalTime: expected 2h18m, got %v", stats.TotalTime)
	}
}

func TestUserStats_TrackBossesDefeated(t *testing.T) {
	stats := NewUserStats()

	stats.BossesDefeated = []ModuleID{ModuleHorizontal, ModuleVertical}

	if len(stats.BossesDefeated) != 2 {
		t.Errorf("BossesDefeated: expected 2, got %d", len(stats.BossesDefeated))
	}
	if stats.BossesDefeated[0] != ModuleHorizontal {
		t.Errorf("First boss should be Horizontal, got %s", stats.BossesDefeated[0])
	}
}

func TestUserStats_LastPlayed(t *testing.T) {
	stats := NewUserStats()

	now := time.Now()
	stats.LastPlayed = now

	if stats.LastPlayed != now {
		t.Errorf("LastPlayed mismatch")
	}
}

// =============================================================================
// BOSS EXERCISE
// =============================================================================

func TestBossExercise_Creation(t *testing.T) {
	boss := BossExercise{
		ID:        "horizontal_boss",
		Module:    ModuleHorizontal,
		Name:      "The Line Walker",
		Lives:     3,
		BonusTime: 30,
		Steps:     []BossStep{},
	}

	if boss.ID != "horizontal_boss" {
		t.Errorf("Boss ID: expected 'horizontal_boss', got %q", boss.ID)
	}
	if boss.Name != "The Line Walker" {
		t.Errorf("Boss Name: expected 'The Line Walker', got %q", boss.Name)
	}
	if boss.Lives != 3 {
		t.Errorf("Boss Lives: expected 3, got %d", boss.Lives)
	}
}

func TestBossStep_Creation(t *testing.T) {
	step := BossStep{
		TimeLimit: 5,
		Exercise: Exercise{
			ID:      "boss_step_1",
			Mission: "Move to end of line",
		},
	}

	if step.TimeLimit != 5 {
		t.Errorf("TimeLimit: expected 5, got %d", step.TimeLimit)
	}
	if step.Exercise.ID != "boss_step_1" {
		t.Errorf("Exercise ID: expected 'boss_step_1', got %q", step.Exercise.ID)
	}
}

func TestBossExercise_WithMultipleSteps(t *testing.T) {
	boss := BossExercise{
		ID:     "test_boss",
		Module: ModuleHorizontal,
		Name:   "Test Boss",
		Lives:  3,
		Steps: []BossStep{
			{TimeLimit: 5, Exercise: Exercise{ID: "step1"}},
			{TimeLimit: 4, Exercise: Exercise{ID: "step2"}},
			{TimeLimit: 3, Exercise: Exercise{ID: "step3"}},
			{TimeLimit: 5, Exercise: Exercise{ID: "step4"}},
			{TimeLimit: 5, Exercise: Exercise{ID: "step5"}},
		},
	}

	if len(boss.Steps) != 5 {
		t.Errorf("Boss should have 5 steps, got %d", len(boss.Steps))
	}
}
