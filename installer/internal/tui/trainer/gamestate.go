package trainer

import (
	"time"
)

// GameState manages the current game session state
type GameState struct {
	// Stats & Progress (persisted)
	Stats *UserStats

	// Current session state (not persisted)
	CurrentModule   ModuleID
	CurrentExercise *Exercise
	CurrentBoss     *BossExercise
	Exercises       []Exercise
	ExerciseIndex   int

	// Mode flags
	IsLessonMode   bool
	IsPracticeMode bool
	IsBossMode     bool

	// Streak and scoring
	CurrentStreak   int
	ComboMultiplier int
	SessionScore    int

	// Boss state
	BossLives      int
	BossStep       int
	IsBossDefeated bool

	// Timing
	TimeElapsed time.Duration
}

// NewGameState creates a new game state with fresh stats
func NewGameState() *GameState {
	return &GameState{
		Stats:           NewUserStats(),
		ComboMultiplier: 1,
	}
}

// NewGameStateWithStats creates a game state with existing stats
func NewGameStateWithStats(stats *UserStats) *GameState {
	if stats == nil {
		stats = NewUserStats()
	}
	return &GameState{
		Stats:           stats,
		ComboMultiplier: 1,
	}
}

// StartLesson starts lesson mode for a module
func (g *GameState) StartLesson(module ModuleID) {
	g.CurrentModule = module
	g.IsLessonMode = true
	g.IsPracticeMode = false
	g.IsBossMode = false
	g.Exercises = GetLessons(module)
	g.ExerciseIndex = 0
	g.CurrentStreak = 0
	g.ComboMultiplier = 1

	if len(g.Exercises) > 0 {
		g.CurrentExercise = &g.Exercises[0]
	}

	// Initialize lesson total in stats
	progress := g.Stats.GetModuleProgress(module)
	progress.LessonsTotal = len(g.Exercises)
}

// StartPractice starts practice mode for a module using intelligent selection
func (g *GameState) StartPractice(module ModuleID) {
	g.CurrentModule = module
	g.IsLessonMode = false
	g.IsPracticeMode = true
	g.IsBossMode = false
	g.Exercises = nil // Not used in intelligent practice mode
	g.ExerciseIndex = 0
	g.CurrentStreak = 0
	g.ComboMultiplier = 1

	// Use weighted random selection for intelligent practice
	progress := g.Stats.GetModuleProgress(module)
	exercise := SelectRandomPracticeExercise(module, progress)
	g.CurrentExercise = exercise
}

// SetPracticeExercise sets a specific exercise for practice mode
func (g *GameState) SetPracticeExercise(exercise *Exercise) {
	g.CurrentExercise = exercise
}

// NextPracticeExercise selects the next random exercise for practice
// Returns false if practice is complete (all mastered)
func (g *GameState) NextPracticeExercise() bool {
	if !g.IsPracticeMode {
		return false
	}

	progress := g.Stats.GetModuleProgress(g.CurrentModule)
	exercise := SelectRandomPracticeExercise(g.CurrentModule, progress)

	if exercise == nil {
		// All exercises mastered - practice complete!
		return false
	}

	g.CurrentExercise = exercise
	return true
}

// StartBoss starts boss fight for a module
func (g *GameState) StartBoss(module ModuleID) {
	g.CurrentModule = module
	g.IsLessonMode = false
	g.IsPracticeMode = false
	g.IsBossMode = true
	g.CurrentBoss = GetBoss(module)
	g.BossStep = 0
	g.CurrentStreak = 0
	g.ComboMultiplier = 1
	g.IsBossDefeated = false

	if g.CurrentBoss != nil {
		g.BossLives = g.CurrentBoss.Lives
		if len(g.CurrentBoss.Steps) > 0 {
			g.CurrentExercise = &g.CurrentBoss.Steps[0].Exercise
		}
	}
}

// RecordCorrectAnswer records a correct answer and updates stats
func (g *GameState) RecordCorrectAnswer(timeSeconds float64, isOptimal bool) {
	g.CurrentStreak++
	if g.CurrentStreak > g.Stats.BestStreak {
		g.Stats.BestStreak = g.CurrentStreak
	}
	g.Stats.CurrentStreak = g.CurrentStreak

	// Calculate and add points
	points := CalculatePoints(g.CurrentExercise, timeSeconds, isOptimal, g.ComboMultiplier)
	g.SessionScore += points
	g.Stats.TotalScore += points

	// Update practice stats
	if g.IsPracticeMode {
		progress := g.Stats.GetModuleProgress(g.CurrentModule)
		progress.PracticeAttempts++
		progress.PracticeCorrect++
		progress.PracticeAccuracy = float64(progress.PracticeCorrect) / float64(progress.PracticeAttempts)
		progress.LastPracticed = time.Now()
	}

	// Update lesson progress
	if g.IsLessonMode {
		progress := g.Stats.GetModuleProgress(g.CurrentModule)
		if g.ExerciseIndex+1 > progress.LessonsCompleted {
			progress.LessonsCompleted = g.ExerciseIndex + 1
		}
	}
}

// RecordIncorrectAnswer records an incorrect answer
func (g *GameState) RecordIncorrectAnswer() {
	g.CurrentStreak = 0
	g.Stats.CurrentStreak = 0
	g.ComboMultiplier = 1

	// Update practice stats
	if g.IsPracticeMode {
		progress := g.Stats.GetModuleProgress(g.CurrentModule)
		progress.PracticeAttempts++
		progress.PracticeAccuracy = float64(progress.PracticeCorrect) / float64(progress.PracticeAttempts)
	}

	// Boss mode: lose a life
	if g.IsBossMode {
		g.BossLives--
		if g.BossLives <= 0 {
			g.BossLives = 0
			g.IsBossDefeated = true
		}
	}
}

// NextExercise advances to the next exercise
func (g *GameState) NextExercise() bool {
	if g.IsBossMode {
		g.BossStep++
		if g.ComboMultiplier < 4 {
			g.ComboMultiplier++
		}

		if g.CurrentBoss == nil || g.BossStep >= len(g.CurrentBoss.Steps) {
			return false
		}
		g.CurrentExercise = &g.CurrentBoss.Steps[g.BossStep].Exercise
		return true
	}

	g.ExerciseIndex++

	if g.ExerciseIndex >= len(g.Exercises) {
		// Complete lessons if in lesson mode
		if g.IsLessonMode {
			progress := g.Stats.GetModuleProgress(g.CurrentModule)
			progress.LessonsCompleted = progress.LessonsTotal
		}
		return false
	}

	g.CurrentExercise = &g.Exercises[g.ExerciseIndex]
	return true
}

// RecordBossVictory records defeating a boss
func (g *GameState) RecordBossVictory() {
	// Add to defeated list if not already
	alreadyDefeated := false
	for _, boss := range g.Stats.BossesDefeated {
		if boss == g.CurrentModule {
			alreadyDefeated = true
			break
		}
	}
	if !alreadyDefeated {
		g.Stats.BossesDefeated = append(g.Stats.BossesDefeated, g.CurrentModule)
	}

	// Update module progress
	progress := g.Stats.GetModuleProgress(g.CurrentModule)
	progress.BossDefeated = true
	progress.BossAttempts++

	if progress.BossBestTime == 0 || g.TimeElapsed < progress.BossBestTime {
		progress.BossBestTime = g.TimeElapsed
	}
	progress.BossLivesLeft = g.BossLives

	// Boss victory bonus
	g.Stats.TotalScore += 500
	g.SessionScore += 500
}

// Reset resets the game state for a new session
func (g *GameState) Reset() {
	g.CurrentModule = ""
	g.CurrentExercise = nil
	g.CurrentBoss = nil
	g.Exercises = nil
	g.ExerciseIndex = 0

	g.IsLessonMode = false
	g.IsPracticeMode = false
	g.IsBossMode = false

	g.CurrentStreak = 0
	g.ComboMultiplier = 1
	g.SessionScore = 0

	g.BossLives = 0
	g.BossStep = 0
	g.IsBossDefeated = false

	g.TimeElapsed = 0
}
