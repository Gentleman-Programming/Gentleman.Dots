package tui

import (
	"bytes"
	"io"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
)

// Helper to read all bytes from io.Reader
func readAll(t *testing.T, r io.Reader) []byte {
	t.Helper()
	bts, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}
	return bts
}

// TestWelcomeScreenGolden tests the welcome screen render against golden file
func TestWelcomeScreenGolden(t *testing.T) {
	m := NewModel()
	m.Width = 80
	m.Height = 24
	m.Screen = ScreenWelcome

	tm := teatest.NewTestModel(t, m,
		teatest.WithInitialTermSize(80, 24),
	)

	// Wait for initial render
	time.Sleep(100 * time.Millisecond)

	// Get final output
	tm.Send(tea.KeyMsg{Type: tea.KeyCtrlC})
	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))

	out := readAll(t, tm.FinalOutput(t))
	teatest.RequireEqualOutput(t, out)
}

// TestMainMenuGolden tests the main menu render against golden file
func TestMainMenuGolden(t *testing.T) {
	m := NewModel()
	m.Width = 80
	m.Height = 24
	m.Screen = ScreenMainMenu

	tm := teatest.NewTestModel(t, m,
		teatest.WithInitialTermSize(80, 24),
	)

	time.Sleep(100 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyCtrlC})
	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))

	out := readAll(t, tm.FinalOutput(t))
	teatest.RequireEqualOutput(t, out)
}

// TestOSSelectGolden tests OS selection screen against golden file
func TestOSSelectGolden(t *testing.T) {
	m := NewModel()
	m.Width = 80
	m.Height = 24
	m.Screen = ScreenOSSelect

	tm := teatest.NewTestModel(t, m,
		teatest.WithInitialTermSize(80, 24),
	)

	time.Sleep(100 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyCtrlC})
	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))

	out := readAll(t, tm.FinalOutput(t))
	teatest.RequireEqualOutput(t, out)
}

// TestNavigationFlowE2E tests navigating from welcome through menu like Playwright would
func TestNavigationFlowE2E(t *testing.T) {
	m := NewModel()
	m.Width = 80
	m.Height = 24

	tm := teatest.NewTestModel(t, m,
		teatest.WithInitialTermSize(80, 24),
	)

	// Start at welcome screen, press Enter to go to main menu
	time.Sleep(50 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})

	// Wait for main menu render and verify we can read output
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("Start Installation")) ||
			bytes.Contains(bts, []byte("Main Menu"))
	}, teatest.WithCheckInterval(50*time.Millisecond), teatest.WithDuration(2*time.Second))

	tm.Send(tea.KeyMsg{Type: tea.KeyCtrlC})
	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}

// TestInstallFlowE2E simulates a complete installation flow selection
func TestInstallFlowE2E(t *testing.T) {
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

	// Main Menu -> Start Installation (already cursor=0)
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(50 * time.Millisecond)

	// Should be at OS Select now
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("Operating System")) ||
			bytes.Contains(bts, []byte("macOS")) ||
			bytes.Contains(bts, []byte("Linux"))
	}, teatest.WithCheckInterval(50*time.Millisecond), teatest.WithDuration(2*time.Second))

	// Select macOS (cursor=0)
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(50 * time.Millisecond)

	// Should be at Terminal Select now
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("Terminal")) ||
			bytes.Contains(bts, []byte("Alacritty")) ||
			bytes.Contains(bts, []byte("WezTerm"))
	}, teatest.WithCheckInterval(50*time.Millisecond), teatest.WithDuration(2*time.Second))

	tm.Send(tea.KeyMsg{Type: tea.KeyCtrlC})
	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}

// TestKeymapsE2E tests navigating keymaps like a real user
func TestKeymapsE2E(t *testing.T) {
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

	// Main Menu -> Navigate down to Keymaps (index 2)
	tm.Send(tea.KeyMsg{Type: tea.KeyDown})
	time.Sleep(20 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyDown})
	time.Sleep(20 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(50 * time.Millisecond)

	// Should be at KeymapsMenu (tool selection: Neovim, Tmux, Zellij, Ghostty)
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("Neovim")) ||
			bytes.Contains(bts, []byte("Tmux")) ||
			bytes.Contains(bts, []byte("Zellij")) ||
			bytes.Contains(bts, []byte("Ghostty"))
	}, teatest.WithCheckInterval(50*time.Millisecond), teatest.WithDuration(2*time.Second))

	// Select Neovim (first option) to get to Neovim keymaps categories
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(50 * time.Millisecond)

	// Should be at Neovim Keymaps categories (Harpoon, Mini.files, etc.)
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("Harpoon")) ||
			bytes.Contains(bts, []byte("Mini.files"))
	}, teatest.WithCheckInterval(50*time.Millisecond), teatest.WithDuration(2*time.Second))

	// Select first category (Harpoon)
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(50 * time.Millisecond)

	// Should show keymaps now with leader key bindings
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("leader")) ||
			bytes.Contains(bts, []byte("Description")) ||
			bytes.Contains(bts, []byte("Keys"))
	}, teatest.WithCheckInterval(50*time.Millisecond), teatest.WithDuration(2*time.Second))

	tm.Send(tea.KeyMsg{Type: tea.KeyCtrlC})
	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}

// TestLearnToolsE2E tests learn about tools navigation
func TestLearnToolsE2E(t *testing.T) {
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

	// Main Menu -> Navigate to Learn About Tools (index 1)
	tm.Send(tea.KeyMsg{Type: tea.KeyDown})
	time.Sleep(20 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(50 * time.Millisecond)

	// Should show tool categories
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("Terminal")) ||
			bytes.Contains(bts, []byte("Shell")) ||
			bytes.Contains(bts, []byte("Multiplexer"))
	}, teatest.WithCheckInterval(50*time.Millisecond), teatest.WithDuration(2*time.Second))

	tm.Send(tea.KeyMsg{Type: tea.KeyCtrlC})
	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}

// TestBackupScreenGolden tests the backup confirmation screen
func TestBackupScreenGolden(t *testing.T) {
	m := NewModel()
	m.Width = 80
	m.Height = 24
	m.Screen = ScreenBackupConfirm
	m.ExistingConfigs = []string{".config/nvim", ".zshrc", ".tmux.conf"}

	tm := teatest.NewTestModel(t, m,
		teatest.WithInitialTermSize(80, 24),
	)

	time.Sleep(100 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyCtrlC})
	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))

	out := readAll(t, tm.FinalOutput(t))
	teatest.RequireEqualOutput(t, out)
}

// TestErrorScreenGolden tests the error screen render
func TestErrorScreenGolden(t *testing.T) {
	m := NewModel()
	m.Width = 80
	m.Height = 24
	m.Screen = ScreenError
	m.ErrorMsg = "Test error: something went wrong during installation"

	tm := teatest.NewTestModel(t, m,
		teatest.WithInitialTermSize(80, 24),
	)

	time.Sleep(100 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyCtrlC})
	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))

	out := readAll(t, tm.FinalOutput(t))
	teatest.RequireEqualOutput(t, out)
}

// TestCompleteScreenGolden tests the completion screen render
func TestCompleteScreenGolden(t *testing.T) {
	m := NewModel()
	m.Width = 80
	m.Height = 24
	m.Screen = ScreenComplete
	m.Choices = UserChoices{
		OS:          "mac",
		Terminal:    "ghostty",
		Shell:       "fish",
		WindowMgr:   "tmux",
		InstallNvim: true,
	}

	tm := teatest.NewTestModel(t, m,
		teatest.WithInitialTermSize(80, 24),
	)

	time.Sleep(100 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyCtrlC})
	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))

	out := readAll(t, tm.FinalOutput(t))
	teatest.RequireEqualOutput(t, out)
}

// TestKeyboardNavigationE2E tests various keyboard interactions
func TestKeyboardNavigationE2E(t *testing.T) {
	t.Run("j/k navigation works like vim", func(t *testing.T) {
		m := NewModel()
		m.Width = 80
		m.Height = 24
		m.Screen = ScreenMainMenu

		tm := teatest.NewTestModel(t, m,
			teatest.WithInitialTermSize(80, 24),
		)

		// Use j to move down
		tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
		time.Sleep(50 * time.Millisecond)
		tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
		time.Sleep(50 * time.Millisecond)

		// Use k to move up
		tm.Send(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
		time.Sleep(50 * time.Millisecond)

		tm.Send(tea.KeyMsg{Type: tea.KeyCtrlC})
		tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
	})

	t.Run("escape goes back from learn screens", func(t *testing.T) {
		m := NewModel()
		m.Width = 80
		m.Height = 24
		m.Screen = ScreenLearnTerminals
		m.PrevScreen = ScreenMainMenu

		tm := teatest.NewTestModel(t, m,
			teatest.WithInitialTermSize(80, 24),
		)

		tm.Send(tea.KeyMsg{Type: tea.KeyEsc})
		time.Sleep(100 * time.Millisecond)

		// Should be back at main menu (learn screens support escape)
		teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
			return bytes.Contains(bts, []byte("Start Installation")) ||
				bytes.Contains(bts, []byte("Main Menu"))
		}, teatest.WithCheckInterval(50*time.Millisecond), teatest.WithDuration(2*time.Second))

		tm.Send(tea.KeyMsg{Type: tea.KeyCtrlC})
		tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
	})
}

// TestResponsiveLayoutE2E tests different terminal sizes
func TestResponsiveLayoutE2E(t *testing.T) {
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
			m.Screen = ScreenMainMenu

			tm := teatest.NewTestModel(t, m,
				teatest.WithInitialTermSize(sz.width, sz.height),
			)

			time.Sleep(100 * time.Millisecond)
			tm.Send(tea.KeyMsg{Type: tea.KeyCtrlC})
			tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))

			// Just verify it doesn't panic and produces output
			out := readAll(t, tm.FinalOutput(t))
			if len(out) == 0 {
				t.Error("Expected some output")
			}
		})
	}
}

// TestLazyVimGuideE2E tests LazyVim guide navigation flow
func TestLazyVimGuideE2E(t *testing.T) {
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

	// Main Menu -> Navigate to LazyVim Guide (index 3)
	tm.Send(tea.KeyMsg{Type: tea.KeyDown})
	time.Sleep(20 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyDown})
	time.Sleep(20 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyDown})
	time.Sleep(20 * time.Millisecond)
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(50 * time.Millisecond)

	// Should be at LazyVim guide screen
	teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
		return bytes.Contains(bts, []byte("LazyVim")) ||
			bytes.Contains(bts, []byte("lazy")) ||
			bytes.Contains(bts, []byte("plugin"))
	}, teatest.WithCheckInterval(50*time.Millisecond), teatest.WithDuration(2*time.Second))

	tm.Send(tea.KeyMsg{Type: tea.KeyCtrlC})
	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second))
}
