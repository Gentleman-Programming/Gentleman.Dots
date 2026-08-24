package trainer

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

// statsConfigPath is the directory for config files (can be overridden for testing)
var statsConfigPath = ""

const statsFileName = "stats.json"

// statsFileJSON represents the JSON structure for stats
type statsFileJSON struct {
	TotalScore       int                            `json:"totalScore"`
	CurrentStreak    int                            `json:"currentStreak"`
	BestStreak       int                            `json:"bestStreak"`
	TotalTimeSeconds int64                          `json:"totalTimeSeconds"`
	LastPlayed       string                         `json:"lastPlayed"`
	BossesDefeated   []string                       `json:"bossesDefeated"`
	Modules          map[string]*moduleProgressJSON `json:"modules"`
}

type moduleProgressJSON struct {
	LessonsCompleted    int                           `json:"lessonsCompleted"`
	LessonsTotal        int                           `json:"lessonsTotal"`
	PracticeAccuracy    float64                       `json:"practiceAccuracy"`
	PracticeAttempts    int                           `json:"practiceAttempts"`
	PracticeCorrect     int                           `json:"practiceCorrect"`
	BossDefeated        bool                          `json:"bossDefeated"`
	BossBestTimeSeconds int64                         `json:"bossBestTimeSeconds"`
	BossAttempts        int                           `json:"bossAttempts"`
	BossLivesLeft       int                           `json:"bossLivesLeft"`
	ExerciseStats       map[string]*exerciseStatsJSON `json:"exerciseStats"`
	WeakExercises       []string                      `json:"weakExercises"`
	LastPracticed       string                        `json:"lastPracticed"`
}

type exerciseStatsJSON struct {
	TotalAttempts    int   `json:"totalAttempts"`
	TotalCorrect     int   `json:"totalCorrect"`
	TotalWrong       int   `json:"totalWrong"`
	ConsecutiveRight int   `json:"consecutiveRight"`
	Mastered         bool  `json:"mastered"`
	LastAttempted    int64 `json:"lastAttempted"`
}

// GetStatsPath returns the full path to the stats file
func GetStatsPath() string {
	dir := getConfigDir()
	if dir == "" {
		return ""
	}
	return filepath.Join(dir, statsFileName)
}

// getConfigDir returns the config directory
func getConfigDir() string {
	if statsConfigPath != "" {
		return statsConfigPath
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(homeDir, ".config", "gentleman-trainer")
}

// LoadStats loads user stats from file
func LoadStats() *UserStats {
	path := GetStatsPath()
	if path == "" {
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	var fileStats statsFileJSON
	if err := json.Unmarshal(data, &fileStats); err != nil {
		return nil
	}

	// Convert JSON struct to UserStats
	stats := &UserStats{
		TotalScore:     fileStats.TotalScore,
		CurrentStreak:  fileStats.CurrentStreak,
		BestStreak:     fileStats.BestStreak,
		TotalTime:      time.Duration(fileStats.TotalTimeSeconds) * time.Second,
		ModuleProgress: make(map[ModuleID]*ModuleProgress),
		BossesDefeated: make([]ModuleID, 0),
	}

	if fileStats.LastPlayed != "" {
		stats.LastPlayed, _ = time.Parse(time.RFC3339, fileStats.LastPlayed)
	}

	for _, boss := range fileStats.BossesDefeated {
		stats.BossesDefeated = append(stats.BossesDefeated, ModuleID(boss))
	}

	for modID, modProgress := range fileStats.Modules {
		mp := &ModuleProgress{
			LessonsCompleted: modProgress.LessonsCompleted,
			LessonsTotal:     modProgress.LessonsTotal,
			PracticeAccuracy: modProgress.PracticeAccuracy,
			PracticeAttempts: modProgress.PracticeAttempts,
			PracticeCorrect:  modProgress.PracticeCorrect,
			BossDefeated:     modProgress.BossDefeated,
			BossBestTime:     time.Duration(modProgress.BossBestTimeSeconds) * time.Second,
			BossAttempts:     modProgress.BossAttempts,
			BossLivesLeft:    modProgress.BossLivesLeft,
			ExerciseStats:    make(map[string]*ExerciseStats),
			WeakExercises:    modProgress.WeakExercises,
		}
		if modProgress.LastPracticed != "" {
			mp.LastPracticed, _ = time.Parse(time.RFC3339, modProgress.LastPracticed)
		}
		// Load exercise stats
		for exID, exStats := range modProgress.ExerciseStats {
			mp.ExerciseStats[exID] = &ExerciseStats{
				TotalAttempts:    exStats.TotalAttempts,
				TotalCorrect:     exStats.TotalCorrect,
				TotalWrong:       exStats.TotalWrong,
				ConsecutiveRight: exStats.ConsecutiveRight,
				Mastered:         exStats.Mastered,
				LastAttempted:    exStats.LastAttempted,
			}
		}
		stats.ModuleProgress[ModuleID(modID)] = mp
	}

	return stats
}

// SaveStats saves user stats to file
func SaveStats(stats *UserStats) error {
	if stats == nil {
		return errors.New("stats cannot be nil")
	}

	path := GetStatsPath()
	if path == "" {
		return errors.New("could not determine stats path")
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Convert to JSON struct
	lastPlayed := ""
	if !stats.LastPlayed.IsZero() {
		lastPlayed = stats.LastPlayed.Format(time.RFC3339)
	}

	fileStats := statsFileJSON{
		TotalScore:       stats.TotalScore,
		CurrentStreak:    stats.CurrentStreak,
		BestStreak:       stats.BestStreak,
		TotalTimeSeconds: int64(stats.TotalTime.Seconds()),
		LastPlayed:       lastPlayed,
		BossesDefeated:   make([]string, 0),
		Modules:          make(map[string]*moduleProgressJSON),
	}

	for _, boss := range stats.BossesDefeated {
		fileStats.BossesDefeated = append(fileStats.BossesDefeated, string(boss))
	}

	for modID, modProgress := range stats.ModuleProgress {
		lastPracticed := ""
		if !modProgress.LastPracticed.IsZero() {
			lastPracticed = modProgress.LastPracticed.Format(time.RFC3339)
		}
		weakExercises := modProgress.WeakExercises
		if weakExercises == nil {
			weakExercises = []string{}
		}

		// Convert exercise stats
		exerciseStats := make(map[string]*exerciseStatsJSON)
		for exID, exStats := range modProgress.ExerciseStats {
			exerciseStats[exID] = &exerciseStatsJSON{
				TotalAttempts:    exStats.TotalAttempts,
				TotalCorrect:     exStats.TotalCorrect,
				TotalWrong:       exStats.TotalWrong,
				ConsecutiveRight: exStats.ConsecutiveRight,
				Mastered:         exStats.Mastered,
				LastAttempted:    exStats.LastAttempted,
			}
		}

		fileStats.Modules[string(modID)] = &moduleProgressJSON{
			LessonsCompleted:    modProgress.LessonsCompleted,
			LessonsTotal:        modProgress.LessonsTotal,
			PracticeAccuracy:    modProgress.PracticeAccuracy,
			PracticeAttempts:    modProgress.PracticeAttempts,
			PracticeCorrect:     modProgress.PracticeCorrect,
			BossDefeated:        modProgress.BossDefeated,
			BossBestTimeSeconds: int64(modProgress.BossBestTime.Seconds()),
			BossAttempts:        modProgress.BossAttempts,
			BossLivesLeft:       modProgress.BossLivesLeft,
			ExerciseStats:       exerciseStats,
			WeakExercises:       weakExercises,
			LastPracticed:       lastPracticed,
		}
	}

	data, err := json.MarshalIndent(fileStats, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// ResetStats deletes the stats file
func ResetStats() error {
	path := GetStatsPath()
	if path == "" {
		return nil
	}
	err := os.Remove(path)
	if os.IsNotExist(err) {
		return nil // Not an error if file doesn't exist
	}
	return err
}
