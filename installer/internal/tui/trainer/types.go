// Package trainer implements a Vim mastery RPG-style trainer
package trainer

import "time"

// ModuleID identifies a training module
type ModuleID string

const (
	ModuleHorizontal   ModuleID = "horizontal"
	ModuleVertical     ModuleID = "vertical"
	ModuleTextObjects  ModuleID = "textobjects"
	ModuleChangeRepeat ModuleID = "cgn"
	ModuleSubstitution ModuleID = "substitution"
	ModuleRegex        ModuleID = "regex"
	ModuleMacros       ModuleID = "macros"
)

// ExerciseType defines the type of exercise
type ExerciseType string

const (
	ExerciseLesson   ExerciseType = "lesson"
	ExercisePractice ExerciseType = "practice"
	ExerciseBoss     ExerciseType = "boss"
)

// Position represents cursor position in code
type Position struct {
	Line int
	Col  int
}

// ExerciseStats tracks individual exercise performance for practice mode
type ExerciseStats struct {
	TotalAttempts    int   // Total times attempted
	TotalCorrect     int   // Total correct answers
	TotalWrong       int   // Total wrong answers
	ConsecutiveRight int   // Current streak of correct answers
	Mastered         bool  // True when mastered (removed from practice pool)
	LastAttempted    int64 // Unix timestamp
}

// Exercise represents a single training exercise
type Exercise struct {
	ID           string       // "horizontal_001"
	Module       ModuleID     // "horizontal", "textobjects", "cgn", etc.
	Level        int          // 1-10
	Type         ExerciseType // "lesson", "practice", "boss"
	Code         []string     // Lines of code to display
	CursorPos    Position     // Initial cursor position
	CursorTarget *Position    // Target cursor position (for movement exercises)
	Mission      string       // "Move cursor to the 'N' in 'Name'"
	Solutions    []string     // ["w", "W", "fe"] - all valid solutions
	Optimal      string       // "w" - the best/shortest solution
	Hint         string       // Hint shown after timeout
	Explanation  string       // Post-answer explanation
	TimeoutSecs  int          // Seconds before showing solution
	Points       int          // Base points for completion
}

// ModuleInfo contains display info for a module
type ModuleInfo struct {
	ID          ModuleID
	Name        string
	Icon        string
	Description string
	BossName    string
}

// BossStep represents a single step in a boss fight
type BossStep struct {
	Exercise  Exercise
	TimeLimit int // Seconds for this specific step
}

// BossExercise represents a boss fight challenge
type BossExercise struct {
	ID        string
	Module    ModuleID
	Name      string     // "The Line Walker"
	Lives     int        // 3
	Steps     []BossStep // Chain of missions
	BonusTime int        // Total time for bonus points
}

// ModuleProgress tracks progress within a module
type ModuleProgress struct {
	// Lessons
	LessonsCompleted int
	LessonsTotal     int

	// Practice
	PracticeAccuracy float64 // 0.0 - 1.0
	PracticeAttempts int
	PracticeCorrect  int

	// Boss
	BossDefeated  bool
	BossBestTime  time.Duration
	BossAttempts  int
	BossLivesLeft int // Lives remaining on best run

	// Per-exercise tracking for intelligent practice
	ExerciseStats map[string]*ExerciseStats

	// Legacy fields (kept for compatibility)
	WeakExercises []string // IDs of exercises that fail most
	LastPracticed time.Time
}

// UserStats contains all user statistics
type UserStats struct {
	TotalScore     int
	CurrentStreak  int
	BestStreak     int
	TotalTime      time.Duration
	ModuleProgress map[ModuleID]*ModuleProgress
	BossesDefeated []ModuleID
	LastPlayed     time.Time
}

// NewUserStats creates a new UserStats with defaults
func NewUserStats() *UserStats {
	return &UserStats{
		ModuleProgress: make(map[ModuleID]*ModuleProgress),
		BossesDefeated: []ModuleID{},
	}
}

// GetModuleProgress returns progress for a module, creating if needed
func (s *UserStats) GetModuleProgress(module ModuleID) *ModuleProgress {
	if s.ModuleProgress == nil {
		s.ModuleProgress = make(map[ModuleID]*ModuleProgress)
	}
	if _, ok := s.ModuleProgress[module]; !ok {
		s.ModuleProgress[module] = &ModuleProgress{}
	}
	return s.ModuleProgress[module]
}

// moduleUnlockOrder defines the order modules are unlocked
var moduleUnlockOrder = []ModuleID{
	ModuleHorizontal,
	ModuleVertical,
	ModuleTextObjects,
	ModuleChangeRepeat,
	ModuleSubstitution,
	ModuleRegex,
	ModuleMacros,
}

// IsModuleUnlocked checks if a module is unlocked
func (s *UserStats) IsModuleUnlocked(module ModuleID) bool {
	// First module is always unlocked
	if module == ModuleHorizontal {
		return true
	}

	// Find position of requested module
	var moduleIdx int = -1
	for i, m := range moduleUnlockOrder {
		if m == module {
			moduleIdx = i
			break
		}
	}

	// Unknown module
	if moduleIdx == -1 {
		return false
	}

	// Need to have defeated the previous boss
	if moduleIdx > 0 {
		prevModule := moduleUnlockOrder[moduleIdx-1]
		return s.IsBossDefeated(prevModule)
	}

	return false
}

// IsBossDefeated checks if a module's boss has been defeated
func (s *UserStats) IsBossDefeated(module ModuleID) bool {
	for _, defeated := range s.BossesDefeated {
		if defeated == module {
			return true
		}
	}
	return false
}

// IsLessonsComplete checks if lessons are 100% complete for a module
func (s *UserStats) IsLessonsComplete(module ModuleID) bool {
	progress := s.GetModuleProgress(module)
	return progress.LessonsTotal > 0 && progress.LessonsCompleted >= progress.LessonsTotal
}

// IsPracticeReady checks if practice mode is unlocked (lessons complete)
func (s *UserStats) IsPracticeReady(module ModuleID) bool {
	return s.IsModuleUnlocked(module) && s.IsLessonsComplete(module)
}

// IsBossReady checks if boss fight is unlocked (80% practice accuracy + minimum attempts)
func (s *UserStats) IsBossReady(module ModuleID) bool {
	if !s.IsPracticeReady(module) {
		return false
	}
	progress := s.GetModuleProgress(module)
	return progress.PracticeAccuracy >= 0.80 && progress.PracticeAttempts >= 10
}

// GetAllModules returns info for all modules in order
func GetAllModules() []ModuleInfo {
	return []ModuleInfo{
		{
			ID:          ModuleHorizontal,
			Name:        "Horizontal Motions",
			Icon:        "🏃",
			Description: "w, W, e, E, b, B, f, F, t, T, ;, ,, 0, $, ^",
			BossName:    "The Line Walker",
		},
		{
			ID:          ModuleVertical,
			Name:        "Vertical Motions",
			Icon:        "📐",
			Description: "j, k, gg, G, {, }, H, M, L, ctrl+d/u/f/b",
			BossName:    "The Code Tower",
		},
		{
			ID:          ModuleTextObjects,
			Name:        "Text Objects",
			Icon:        "🎯",
			Description: "viw, vaw, vi\", va\", vi{, diw, daw, ci\", di{, yiw, yi\"",
			BossName:    "The Bracket Demon",
		},
		{
			ID:          ModuleChangeRepeat,
			Name:        "Change & Repeat",
			Icon:        "🔁",
			Description: "d, c, dd, D, C, x, *, #, n, N, gn, cgn, dgn, .",
			BossName:    "The Clone Army",
		},
		{
			ID:          ModuleSubstitution,
			Name:        "Substitution",
			Icon:        "🔄",
			Description: "r, R, s, S, ~, gu, gU, J, :s, :%s, flags (g, c, i)",
			BossName:    "The Transformer",
		},
		{
			ID:          ModuleRegex,
			Name:        "Regex & Vimgrep",
			Icon:        "🔍",
			Description: "/, ?, n, N, *, #, \\v, :vimgrep, :copen, :cnext",
			BossName:    "The Pattern Master",
		},
		{
			ID:          ModuleMacros,
			Name:        "Macros",
			Icon:        "🎪",
			Description: "qa, q, @a, @@, :normal, :g/pattern/",
			BossName:    "The Automaton",
		},
	}
}
