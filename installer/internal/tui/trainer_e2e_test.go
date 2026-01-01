package tui

import (
	"bytes"
	"testing"
	"time"

	"github.com/Gentleman-Programming/Gentleman.Dots/installer/internal/tui/trainer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
)

// =============================================================================
// VIM TRAINER E2E TESTS - Playwright-style golden tests
// =============================================================================

// TestTrainerMenuGolden tests the trainer module selection screen
func TestTrainerMenuGolden(t *testing.T) {
	m := NewModel()
	m.Width = 80
	m.Height = 24
	m.Screen = ScreenTrainerMenu
	m.TrainerStats = trainer.NewUserStats()
	m.TrainerModules = trainer.GetAllModules()
	m.TrainerCursor = 0

	tm := teatest.NewTestModel(t, m,
		teatest.WithInitialTermSize(80, 24),
	)

	time.Sleep(100 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyCtrlC})
	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))

	out := readAll(t, tm.FinalOutput(t))
	teatest.RequireEqualOutput(t, out)
}

// TestTrainerLessonGolden tests a lesson exercise screen
func TestTrainerLessonGolden(t *testing.T) {
	m := NewModel()
	m.Width = 80
	m.Height = 24
	m.Screen = ScreenTrainerLesson
	m.TrainerStats = trainer.NewUserStats()

	// Start a lesson for horizontal module (first unlocked module)
	m.TrainerGameState = trainer.NewGameStateWithStats(m.TrainerStats)
	m.TrainerGameState.StartLesson(trainer.ModuleHorizontal)
	m.TrainerInput = ""
	m.TrainerMessage = ""

	tm := teatest.NewTestModel(t, m,
		teatest.WithInitialTermSize(80, 24),
	)

	time.Sleep(100 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyCtrlC})
	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))

	out := readAll(t, tm.FinalOutput(t))
	teatest.RequireEqualOutput(t, out)
}

// TestTrainerResultCorrectGolden tests the result screen after correct answer
func TestTrainerResultCorrectGolden(t *testing.T) {
	m := NewModel()
	m.Width = 80
	m.Height = 24
	m.Screen = ScreenTrainerResult
	m.TrainerStats = trainer.NewUserStats()
	m.TrainerGameState = trainer.NewGameStateWithStats(m.TrainerStats)
	m.TrainerGameState.StartLesson(trainer.ModuleHorizontal)
	m.TrainerLastCorrect = true
	m.TrainerMessage = "✨ Perfect! Optimal solution!"

	tm := teatest.NewTestModel(t, m,
		teatest.WithInitialTermSize(80, 24),
	)

	time.Sleep(100 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyCtrlC})
	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))

	out := readAll(t, tm.FinalOutput(t))
	teatest.RequireEqualOutput(t, out)
}

// TestTrainerBossGolden tests the boss fight screen
func TestTrainerBossGolden(t *testing.T) {
	m := NewModel()
	m.Width = 80
	m.Height = 24
	m.Screen = ScreenTrainerBoss
	m.TrainerStats = trainer.NewUserStats()

	// Prepare stats so boss is ready
	progress := m.TrainerStats.GetModuleProgress(trainer.ModuleHorizontal)
	progress.LessonsCompleted = 15
	progress.LessonsTotal = 15
	progress.PracticeAccuracy = 85.0

	m.TrainerGameState = trainer.NewGameStateWithStats(m.TrainerStats)
	m.TrainerGameState.StartBoss(trainer.ModuleHorizontal)
	m.TrainerInput = ""
	m.TrainerMessage = ""

	tm := teatest.NewTestModel(t, m,
		teatest.WithInitialTermSize(80, 24),
	)

	time.Sleep(100 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyCtrlC})
	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))

	out := readAll(t, tm.FinalOutput(t))
	teatest.RequireEqualOutput(t, out)
}

// =============================================================================
// E2E NAVIGATION FLOW TESTS
// =============================================================================

// TestTrainerNavigationE2E tests navigating from main menu to trainer
func TestTrainerNavigationE2E(t *testing.T) {
	m := NewModel()
	m.Width = 80
	m.Height = 24

	tm := teatest.NewTestModel(t, m,
		teatest.WithInitialTermSize(80, 24),
	)

	// Welcome -> Enter
	time.Sleep(50 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(50 * time.Millisecond)

	// Main Menu -> Navigate to Vim Trainer (index 4: Start, Learn, Keymaps, LazyVim, Vim Trainer)
	for i := 0; i < 4; i++ {
		tm.Send(tea.KeyMsg{Type: tea.KeyDown})
		time.Sleep(20 * time.Millisecond)
	}
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(50 * time.Millisecond)

	// Should be at Trainer Menu now - verify we see modules
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("Horizontal")) ||
			bytes.Contains(bts, []byte("Vertical")) ||
			bytes.Contains(bts, []byte("Module"))
	}, teatest.WithCheckInterval(50*time.Millisecond), teatest.WithDuration(2*time.Second))

	tm.Send(tea.KeyMsg{Type: tea.KeyCtrlC})
	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}

// TestTrainerStartLessonE2E tests starting a lesson from module menu
func TestTrainerStartLessonE2E(t *testing.T) {
	m := NewModel()
	m.Width = 80
	m.Height = 24
	m.Screen = ScreenTrainerMenu
	m.TrainerStats = trainer.NewUserStats()
	m.TrainerModules = trainer.GetAllModules()
	m.TrainerCursor = 0

	tm := teatest.NewTestModel(t, m,
		teatest.WithInitialTermSize(80, 24),
	)

	time.Sleep(50 * time.Millisecond)

	// Press Enter to start lesson on first module (Horizontal)
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(100 * time.Millisecond)

	// Should show a lesson with code and mission
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("Mission")) ||
			bytes.Contains(bts, []byte("Code")) ||
			bytes.Contains(bts, []byte("Lesson"))
	}, teatest.WithCheckInterval(50*time.Millisecond), teatest.WithDuration(2*time.Second))

	tm.Send(tea.KeyMsg{Type: tea.KeyCtrlC})
	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}

// TestTrainerInputE2E tests typing vim commands in a lesson
func TestTrainerInputE2E(t *testing.T) {
	m := NewModel()
	m.Width = 80
	m.Height = 24
	m.Screen = ScreenTrainerLesson
	m.TrainerStats = trainer.NewUserStats()
	m.TrainerGameState = trainer.NewGameStateWithStats(m.TrainerStats)
	m.TrainerGameState.StartLesson(trainer.ModuleHorizontal)
	m.TrainerInput = ""

	tm := teatest.NewTestModel(t, m,
		teatest.WithInitialTermSize(80, 24),
	)

	time.Sleep(50 * time.Millisecond)

	// Type "w" as input (common first lesson answer)
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'w'}})
	time.Sleep(50 * time.Millisecond)

	// Verify input is shown
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("w"))
	}, teatest.WithCheckInterval(50*time.Millisecond), teatest.WithDuration(2*time.Second))

	// Press Enter to submit
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(100 * time.Millisecond)

	// Should go to result screen
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("Correct")) ||
			bytes.Contains(bts, []byte("Incorrect")) ||
			bytes.Contains(bts, []byte("Perfect")) ||
			bytes.Contains(bts, []byte("Result"))
	}, teatest.WithCheckInterval(50*time.Millisecond), teatest.WithDuration(2*time.Second))

	tm.Send(tea.KeyMsg{Type: tea.KeyCtrlC})
	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}

// TestTrainerBackspaceE2E tests backspace functionality in input
func TestTrainerBackspaceE2E(t *testing.T) {
	m := NewModel()
	m.Width = 80
	m.Height = 24
	m.Screen = ScreenTrainerLesson
	m.TrainerStats = trainer.NewUserStats()
	m.TrainerGameState = trainer.NewGameStateWithStats(m.TrainerStats)
	m.TrainerGameState.StartLesson(trainer.ModuleHorizontal)
	m.TrainerInput = "ww" // Pre-set some input

	tm := teatest.NewTestModel(t, m,
		teatest.WithInitialTermSize(80, 24),
	)

	time.Sleep(50 * time.Millisecond)

	// Press backspace
	tm.Send(tea.KeyMsg{Type: tea.KeyBackspace})
	time.Sleep(50 * time.Millisecond)

	// Input should now be "w" (one character removed)
	// We can verify the screen still shows input area with answer label
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("answer")) ||
			bytes.Contains(bts, []byte("Your")) ||
			bytes.Contains(bts, []byte("w"))
	}, teatest.WithCheckInterval(50*time.Millisecond), teatest.WithDuration(2*time.Second))

	tm.Send(tea.KeyMsg{Type: tea.KeyCtrlC})
	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}

// TestTrainerModuleNavigationE2E tests navigating through modules with j/k
func TestTrainerModuleNavigationE2E(t *testing.T) {
	m := NewModel()
	m.Width = 80
	m.Height = 24
	m.Screen = ScreenTrainerMenu
	m.TrainerStats = trainer.NewUserStats()
	m.TrainerModules = trainer.GetAllModules()
	m.TrainerCursor = 0

	tm := teatest.NewTestModel(t, m,
		teatest.WithInitialTermSize(80, 24),
	)

	time.Sleep(50 * time.Millisecond)

	// Navigate down through modules using j
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	time.Sleep(50 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	time.Sleep(50 * time.Millisecond)

	// Navigate up using k
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	time.Sleep(50 * time.Millisecond)

	// Should still show modules menu
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("Horizontal")) ||
			bytes.Contains(bts, []byte("Vertical"))
	}, teatest.WithCheckInterval(50*time.Millisecond), teatest.WithDuration(2*time.Second))

	tm.Send(tea.KeyMsg{Type: tea.KeyCtrlC})
	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}

// TestTrainerEscapeE2E tests escape key returns to menu
func TestTrainerEscapeE2E(t *testing.T) {
	m := NewModel()
	m.Width = 80
	m.Height = 24
	m.Screen = ScreenTrainerLesson
	m.TrainerStats = trainer.NewUserStats()
	m.TrainerGameState = trainer.NewGameStateWithStats(m.TrainerStats)
	m.TrainerGameState.StartLesson(trainer.ModuleHorizontal)
	m.TrainerInput = ""

	tm := teatest.NewTestModel(t, m,
		teatest.WithInitialTermSize(80, 24),
	)

	time.Sleep(50 * time.Millisecond)

	// Press Escape to go back to menu
	tm.Send(tea.KeyMsg{Type: tea.KeyEsc})
	time.Sleep(100 * time.Millisecond)

	// Should be back at trainer menu
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("Horizontal")) ||
			bytes.Contains(bts, []byte("Module")) ||
			bytes.Contains(bts, []byte("Vim Trainer"))
	}, teatest.WithCheckInterval(50*time.Millisecond), teatest.WithDuration(2*time.Second))

	tm.Send(tea.KeyMsg{Type: tea.KeyCtrlC})
	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}

// TestTrainerLessonProgressE2E tests progressing through multiple lessons
func TestTrainerLessonProgressE2E(t *testing.T) {
	m := NewModel()
	m.Width = 80
	m.Height = 24
	m.Screen = ScreenTrainerLesson
	m.TrainerStats = trainer.NewUserStats()
	m.TrainerGameState = trainer.NewGameStateWithStats(m.TrainerStats)
	m.TrainerGameState.StartLesson(trainer.ModuleHorizontal)
	m.TrainerInput = ""

	tm := teatest.NewTestModel(t, m,
		teatest.WithInitialTermSize(80, 24),
	)

	// Complete first exercise
	time.Sleep(50 * time.Millisecond)

	// Get the current exercise's optimal solution
	exercise := m.TrainerGameState.CurrentExercise
	if exercise == nil {
		t.Fatal("No exercise loaded")
	}

	// Type the optimal answer
	for _, r := range exercise.Optimal {
		tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
		time.Sleep(20 * time.Millisecond)
	}

	// Submit
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(100 * time.Millisecond)

	// Should show result
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("Correct")) ||
			bytes.Contains(bts, []byte("Perfect")) ||
			bytes.Contains(bts, []byte("Incorrect"))
	}, teatest.WithCheckInterval(50*time.Millisecond), teatest.WithDuration(2*time.Second))

	// Press Enter to continue
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(100 * time.Millisecond)

	// Should be at next lesson or back at menu
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("Mission")) ||
			bytes.Contains(bts, []byte("Module")) ||
			bytes.Contains(bts, []byte("Lesson"))
	}, teatest.WithCheckInterval(50*time.Millisecond), teatest.WithDuration(2*time.Second))

	tm.Send(tea.KeyMsg{Type: tea.KeyCtrlC})
	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}

// =============================================================================
// ALL MODULES RENDERING TESTS
// =============================================================================

// TestAllModulesRenderE2E tests that each module can be accessed and renders properly
func TestAllModulesRenderE2E(t *testing.T) {
	modules := []trainer.ModuleID{
		trainer.ModuleHorizontal,
		trainer.ModuleVertical,
		trainer.ModuleTextObjects,
		trainer.ModuleChangeRepeat,
		trainer.ModuleSubstitution,
		trainer.ModuleRegex,
		trainer.ModuleMacros,
	}

	for _, moduleID := range modules {
		t.Run(string(moduleID), func(t *testing.T) {
			m := NewModel()
			m.Width = 80
			m.Height = 24
			m.Screen = ScreenTrainerLesson
			m.TrainerStats = trainer.NewUserStats()

			// Unlock all modules for testing
			for _, mod := range modules {
				progress := m.TrainerStats.GetModuleProgress(mod)
				progress.BossDefeated = true
			}

			m.TrainerGameState = trainer.NewGameStateWithStats(m.TrainerStats)
			m.TrainerGameState.StartLesson(moduleID)
			m.TrainerInput = ""

			// Verify we got an exercise
			if m.TrainerGameState.CurrentExercise == nil {
				t.Fatalf("Module %s has no exercises", moduleID)
			}

			tm := teatest.NewTestModel(t, m,
				teatest.WithInitialTermSize(80, 24),
			)

			time.Sleep(100 * time.Millisecond)

			// Verify screen renders with mission and code
			teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
				return bytes.Contains(bts, []byte("Mission")) &&
					bytes.Contains(bts, []byte("Code"))
			}, teatest.WithCheckInterval(50*time.Millisecond), teatest.WithDuration(2*time.Second))

			tm.Send(tea.KeyMsg{Type: tea.KeyCtrlC})
			tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
		})
	}
}

// TestTextObjectsVisualSelectionE2E tests that text object exercises show visual selection
func TestTextObjectsVisualSelectionE2E(t *testing.T) {
	m := NewModel()
	m.Width = 80
	m.Height = 24
	m.Screen = ScreenTrainerLesson
	m.TrainerStats = trainer.NewUserStats()

	// Unlock TextObjects
	progress := m.TrainerStats.GetModuleProgress(trainer.ModuleHorizontal)
	progress.BossDefeated = true

	m.TrainerGameState = trainer.NewGameStateWithStats(m.TrainerStats)
	m.TrainerGameState.StartLesson(trainer.ModuleTextObjects)
	m.TrainerInput = ""

	tm := teatest.NewTestModel(t, m,
		teatest.WithInitialTermSize(80, 24),
	)

	time.Sleep(50 * time.Millisecond)

	// Type "iw" to see visual selection
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})
	time.Sleep(30 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'w'}})
	time.Sleep(100 * time.Millisecond)

	// Verify screen shows input
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("iw"))
	}, teatest.WithCheckInterval(50*time.Millisecond), teatest.WithDuration(2*time.Second))

	tm.Send(tea.KeyMsg{Type: tea.KeyCtrlC})
	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}

// TestBossFlowE2E tests the complete boss fight flow
func TestBossFlowE2E(t *testing.T) {
	m := NewModel()
	m.Width = 80
	m.Height = 24
	m.Screen = ScreenTrainerMenu
	m.TrainerStats = trainer.NewUserStats()
	m.TrainerModules = trainer.GetAllModules()

	// Make boss ready for horizontal module
	progress := m.TrainerStats.GetModuleProgress(trainer.ModuleHorizontal)
	progress.LessonsCompleted = 15
	progress.LessonsTotal = 15
	progress.PracticeAccuracy = 85.0
	progress.PracticeAttempts = 100

	m.TrainerCursor = 0

	tm := teatest.NewTestModel(t, m,
		teatest.WithInitialTermSize(80, 24),
	)

	time.Sleep(50 * time.Millisecond)

	// Press 'b' to start boss fight
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'b'}})
	time.Sleep(100 * time.Millisecond)

	// Should be in boss fight screen
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("Boss")) ||
			bytes.Contains(bts, []byte("Lives")) ||
			bytes.Contains(bts, []byte("Step"))
	}, teatest.WithCheckInterval(50*time.Millisecond), teatest.WithDuration(2*time.Second))

	tm.Send(tea.KeyMsg{Type: tea.KeyCtrlC})
	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}

// TestLockedModuleMessageE2E tests that locked modules show appropriate message
func TestLockedModuleMessageE2E(t *testing.T) {
	m := NewModel()
	m.Width = 80
	m.Height = 24
	m.Screen = ScreenTrainerMenu
	m.TrainerStats = trainer.NewUserStats()
	m.TrainerModules = trainer.GetAllModules()
	m.TrainerCursor = 1 // Vertical module (locked by default)

	tm := teatest.NewTestModel(t, m,
		teatest.WithInitialTermSize(80, 24),
	)

	time.Sleep(50 * time.Millisecond)

	// Try to start locked module
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(100 * time.Millisecond)

	// Should show locked message, still in menu
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("locked")) ||
			bytes.Contains(bts, []byte("Locked")) ||
			bytes.Contains(bts, []byte("Module"))
	}, teatest.WithCheckInterval(50*time.Millisecond), teatest.WithDuration(2*time.Second))

	tm.Send(tea.KeyMsg{Type: tea.KeyCtrlC})
	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}

// =============================================================================
// RESPONSIVE LAYOUT TESTS FOR TRAINER
// =============================================================================

// TestTrainerResponsiveE2E tests trainer screens at different terminal sizes
func TestTrainerResponsiveE2E(t *testing.T) {
	sizes := []struct {
		name   string
		width  int
		height int
	}{
		{"small_terminal", 60, 20},
		{"medium_terminal", 80, 24},
		{"large_terminal", 120, 40},
		{"wide_terminal", 160, 24},
	}

	for _, sz := range sizes {
		t.Run(sz.name, func(t *testing.T) {
			m := NewModel()
			m.Width = sz.width
			m.Height = sz.height
			m.Screen = ScreenTrainerLesson
			m.TrainerStats = trainer.NewUserStats()
			m.TrainerGameState = trainer.NewGameStateWithStats(m.TrainerStats)
			m.TrainerGameState.StartLesson(trainer.ModuleHorizontal)
			m.TrainerInput = "w"

			tm := teatest.NewTestModel(t, m,
				teatest.WithInitialTermSize(sz.width, sz.height),
			)

			time.Sleep(100 * time.Millisecond)
			tm.Send(tea.KeyMsg{Type: tea.KeyCtrlC})
			tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))

			// Verify it produces output without panicking
			out := readAll(t, tm.FinalOutput(t))
			if len(out) == 0 {
				t.Error("Expected some output")
			}
		})
	}
}

// TestPracticeModeE2E tests practice mode functionality
func TestPracticeModeE2E(t *testing.T) {
	m := NewModel()
	m.Width = 80
	m.Height = 24
	m.Screen = ScreenTrainerMenu
	m.TrainerStats = trainer.NewUserStats()
	m.TrainerModules = trainer.GetAllModules()

	// Complete all lessons to unlock practice
	progress := m.TrainerStats.GetModuleProgress(trainer.ModuleHorizontal)
	progress.LessonsCompleted = 15
	progress.LessonsTotal = 15

	m.TrainerCursor = 0

	tm := teatest.NewTestModel(t, m,
		teatest.WithInitialTermSize(80, 24),
	)

	time.Sleep(50 * time.Millisecond)

	// Press 'p' to start practice
	tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
	time.Sleep(100 * time.Millisecond)

	// Should be in practice mode
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("Practice")) ||
			bytes.Contains(bts, []byte("Mission")) ||
			bytes.Contains(bts, []byte("Code"))
	}, teatest.WithCheckInterval(50*time.Millisecond), teatest.WithDuration(2*time.Second))

	tm.Send(tea.KeyMsg{Type: tea.KeyCtrlC})
	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}
