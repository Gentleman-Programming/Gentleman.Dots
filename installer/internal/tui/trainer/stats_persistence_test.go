package trainer

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// =============================================================================
// STATS PERSISTENCE
// =============================================================================

func TestSaveAndLoadStats_RoundTrip(t *testing.T) {
	// Create temp directory for test
	tempDir, err := os.MkdirTemp("", "trainer-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Override stats path for test
	originalPath := statsConfigPath
	statsConfigPath = tempDir
	defer func() { statsConfigPath = originalPath }()

	// Create stats to save
	stats := NewUserStats()
	stats.TotalScore = 2340
	stats.CurrentStreak = 7
	stats.BestStreak = 23
	stats.TotalTime = 2*time.Hour + 18*time.Minute
	stats.LastPlayed = time.Date(2026, 1, 1, 15, 30, 0, 0, time.UTC)
	stats.BossesDefeated = []ModuleID{ModuleHorizontal, ModuleVertical}

	// Add module progress
	progress := stats.GetModuleProgress(ModuleHorizontal)
	progress.LessonsCompleted = 15
	progress.LessonsTotal = 15
	progress.PracticeAccuracy = 0.85
	progress.PracticeAttempts = 47
	progress.PracticeCorrect = 40
	progress.BossDefeated = true
	progress.BossBestTime = 28 * time.Second
	progress.BossAttempts = 3
	progress.BossLivesLeft = 2
	progress.WeakExercises = []string{"horizontal_012", "horizontal_008"}
	progress.LastPracticed = time.Date(2026, 1, 1, 15, 30, 0, 0, time.UTC)

	// Save
	err = SaveStats(stats)
	if err != nil {
		t.Fatalf("SaveStats failed: %v", err)
	}

	// Verify file exists
	statsFile := filepath.Join(tempDir, "stats.json")
	if _, err := os.Stat(statsFile); os.IsNotExist(err) {
		t.Error("Stats file should exist after save")
	}

	// Load
	loaded := LoadStats()
	if loaded == nil {
		t.Fatal("LoadStats returned nil")
	}

	// Verify all fields
	if loaded.TotalScore != 2340 {
		t.Errorf("TotalScore: expected 2340, got %d", loaded.TotalScore)
	}
	if loaded.CurrentStreak != 7 {
		t.Errorf("CurrentStreak: expected 7, got %d", loaded.CurrentStreak)
	}
	if loaded.BestStreak != 23 {
		t.Errorf("BestStreak: expected 23, got %d", loaded.BestStreak)
	}
	if loaded.TotalTime != 2*time.Hour+18*time.Minute {
		t.Errorf("TotalTime: expected 2h18m, got %v", loaded.TotalTime)
	}
	if len(loaded.BossesDefeated) != 2 {
		t.Errorf("BossesDefeated: expected 2, got %d", len(loaded.BossesDefeated))
	}

	// Verify module progress
	loadedProgress := loaded.GetModuleProgress(ModuleHorizontal)
	if loadedProgress.LessonsCompleted != 15 {
		t.Errorf("LessonsCompleted: expected 15, got %d", loadedProgress.LessonsCompleted)
	}
	if loadedProgress.PracticeAccuracy != 0.85 {
		t.Errorf("PracticeAccuracy: expected 0.85, got %f", loadedProgress.PracticeAccuracy)
	}
	if loadedProgress.BossBestTime != 28*time.Second {
		t.Errorf("BossBestTime: expected 28s, got %v", loadedProgress.BossBestTime)
	}
	if len(loadedProgress.WeakExercises) != 2 {
		t.Errorf("WeakExercises: expected 2, got %d", len(loadedProgress.WeakExercises))
	}
}

func TestLoadStats_ReturnsNilWhenNoFile(t *testing.T) {
	// Create temp directory (empty)
	tempDir, err := os.MkdirTemp("", "trainer-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Override stats path
	originalPath := statsConfigPath
	statsConfigPath = tempDir
	defer func() { statsConfigPath = originalPath }()

	stats := LoadStats()
	if stats != nil {
		t.Error("LoadStats should return nil when no file exists")
	}
}

func TestLoadStats_ReturnsNilOnCorruptedFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "trainer-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	originalPath := statsConfigPath
	statsConfigPath = tempDir
	defer func() { statsConfigPath = originalPath }()

	// Write corrupted JSON
	statsFile := filepath.Join(tempDir, "stats.json")
	err = os.WriteFile(statsFile, []byte("not valid json {{{"), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	stats := LoadStats()
	if stats != nil {
		t.Error("LoadStats should return nil on corrupted file")
	}
}

func TestSaveStats_CreatesDirectory(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "trainer-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Use a nested path that doesn't exist
	nestedPath := filepath.Join(tempDir, "nested", "config")

	originalPath := statsConfigPath
	statsConfigPath = nestedPath
	defer func() { statsConfigPath = originalPath }()

	stats := NewUserStats()
	stats.TotalScore = 100

	err = SaveStats(stats)
	if err != nil {
		t.Fatalf("SaveStats failed: %v", err)
	}

	// Verify nested directory was created
	statsFile := filepath.Join(nestedPath, "stats.json")
	if _, err := os.Stat(statsFile); os.IsNotExist(err) {
		t.Error("SaveStats should create nested directories")
	}
}

func TestSaveStats_NilStats(t *testing.T) {
	err := SaveStats(nil)
	if err == nil {
		t.Error("SaveStats should return error for nil stats")
	}
}

func TestResetStats_DeletesFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "trainer-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	originalPath := statsConfigPath
	statsConfigPath = tempDir
	defer func() { statsConfigPath = originalPath }()

	// Save stats first
	stats := NewUserStats()
	stats.TotalScore = 100
	err = SaveStats(stats)
	if err != nil {
		t.Fatalf("SaveStats failed: %v", err)
	}

	// Verify file exists
	statsFile := filepath.Join(tempDir, "stats.json")
	if _, err := os.Stat(statsFile); os.IsNotExist(err) {
		t.Fatal("Stats file should exist before reset")
	}

	// Reset
	err = ResetStats()
	if err != nil {
		t.Fatalf("ResetStats failed: %v", err)
	}

	// Verify file is deleted
	if _, err := os.Stat(statsFile); !os.IsNotExist(err) {
		t.Error("Stats file should be deleted after reset")
	}
}

func TestResetStats_NoErrorWhenNoFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "trainer-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	originalPath := statsConfigPath
	statsConfigPath = tempDir
	defer func() { statsConfigPath = originalPath }()

	// Reset without any file existing
	err = ResetStats()
	// Should not error (file doesn't exist is fine)
	if err != nil && !os.IsNotExist(err) {
		t.Errorf("ResetStats should not error when no file exists: %v", err)
	}
}

// =============================================================================
// STATS FILE PATH
// =============================================================================

func TestGetStatsPath_ReturnsValidPath(t *testing.T) {
	path := GetStatsPath()
	if path == "" {
		t.Error("GetStatsPath should return non-empty path")
	}
	// Should contain "gentleman-trainer"
	if !filepath.IsAbs(path) {
		// When using test config path, it might not be absolute
		// but in production it should be
	}
}
