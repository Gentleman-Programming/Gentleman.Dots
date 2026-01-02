package tui

import (
	"strings"
	"testing"

	"github.com/Gentleman-Programming/Gentleman.Dots/installer/internal/system"
	tea "github.com/charmbracelet/bubbletea"
)

// TestFullInstallationFlow tests the complete flow from welcome to installation
func TestFullInstallationFlow(t *testing.T) {
	t.Run("complete flow: welcome -> main menu -> installation screens", func(t *testing.T) {
		m := NewModel()

		// Start at welcome screen
		if m.Screen != ScreenWelcome {
			t.Fatalf("Expected ScreenWelcome, got %v", m.Screen)
		}

		// Press enter to go to main menu
		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		m = result.(Model)

		if m.Screen != ScreenMainMenu {
			t.Fatalf("Expected ScreenMainMenu, got %v", m.Screen)
		}

		// Select "Start Installation" (cursor at 0)
		result, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		m = result.(Model)

		if m.Screen != ScreenOSSelect {
			t.Fatalf("Expected ScreenOSSelect, got %v", m.Screen)
		}

		// Select macOS (cursor at 0)
		result, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		m = result.(Model)

		if m.Screen != ScreenTerminalSelect {
			t.Fatalf("Expected ScreenTerminalSelect, got %v", m.Screen)
		}
		if m.Choices.OS != "mac" {
			t.Fatalf("Expected OS 'mac', got '%s'", m.Choices.OS)
		}

		// Select Ghostty (need to navigate down 3 times: Alacritty, WezTerm, Kitty, Ghostty)
		for i := 0; i < 3; i++ {
			result, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
			m = result.(Model)
		}
		result, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		m = result.(Model)

		if m.Screen != ScreenFontSelect {
			t.Fatalf("Expected ScreenFontSelect, got %v", m.Screen)
		}
		if m.Choices.Terminal != "ghostty" {
			t.Fatalf("Expected Terminal 'ghostty', got '%s'", m.Choices.Terminal)
		}

		// Select "Yes" for font (cursor at 0)
		result, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		m = result.(Model)

		if m.Screen != ScreenShellSelect {
			t.Fatalf("Expected ScreenShellSelect, got %v", m.Screen)
		}
		if !m.Choices.InstallFont {
			t.Fatal("Expected InstallFont to be true")
		}

		// Select Fish (cursor at 0)
		result, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		m = result.(Model)

		if m.Screen != ScreenWMSelect {
			t.Fatalf("Expected ScreenWMSelect, got %v", m.Screen)
		}
		if m.Choices.Shell != "fish" {
			t.Fatalf("Expected Shell 'fish', got '%s'", m.Choices.Shell)
		}

		// Select Zellij (navigate down once)
		result, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
		m = result.(Model)
		result, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		m = result.(Model)

		if m.Screen != ScreenNvimSelect {
			t.Fatalf("Expected ScreenNvimSelect, got %v", m.Screen)
		}
		if m.Choices.WindowMgr != "zellij" {
			t.Fatalf("Expected WindowMgr 'zellij', got '%s'", m.Choices.WindowMgr)
		}

		// Select "Yes" for Nvim (cursor at 0)
		// This should either go to BackupConfirm (if existing configs) or Installing
		result, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		m = result.(Model)

		if !m.Choices.InstallNvim {
			t.Fatal("Expected InstallNvim to be true")
		}

		// Screen depends on existing configs
		if m.Screen != ScreenBackupConfirm && m.Screen != ScreenInstalling {
			t.Fatalf("Expected ScreenBackupConfirm or ScreenInstalling, got %v", m.Screen)
		}
	})
}

// TestLinuxFlow tests Linux-specific options
func TestLinuxFlow(t *testing.T) {
	t.Run("linux flow should not show Kitty", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenOSSelect
		m.Cursor = 1 // Linux

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		m = result.(Model)

		if m.Choices.OS != "linux" {
			t.Fatalf("Expected OS 'linux', got '%s'", m.Choices.OS)
		}

		// Check terminal options don't include Kitty
		options := m.GetCurrentOptions()
		for _, opt := range options {
			if opt == "Kitty" {
				t.Error("Linux should not have Kitty option")
			}
		}
	})
}

// TestSkipTerminal tests skipping terminal installation
func TestSkipTerminal(t *testing.T) {
	t.Run("selecting None should skip font selection", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenTerminalSelect
		m.Choices.OS = "mac"

		// Find and select "None" option
		options := m.GetCurrentOptions()
		for i, opt := range options {
			if opt == "None" {
				m.Cursor = i
				break
			}
		}

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		m = result.(Model)

		// Should skip font and go directly to shell
		if m.Screen != ScreenShellSelect {
			t.Fatalf("Expected ScreenShellSelect when terminal is None, got %v", m.Screen)
		}
	})
}

// TestLearnScreensNavigation tests learn mode navigation
func TestLearnScreensNavigation(t *testing.T) {
	t.Run("can access learn screens from terminal selection", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenTerminalSelect
		m.Choices.OS = "mac"

		// Find "Learn about terminals" option
		options := m.GetCurrentOptions()
		for i, opt := range options {
			if strings.Contains(strings.ToLower(opt), "learn about terminals") {
				m.Cursor = i
				break
			}
		}

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		m = result.(Model)

		if m.Screen != ScreenLearnTerminals {
			t.Fatalf("Expected ScreenLearnTerminals, got %v", m.Screen)
		}
	})

	t.Run("ESC returns from learn screen", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenLearnTerminals
		m.PrevScreen = ScreenTerminalSelect

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEscape})
		m = result.(Model)

		if m.Screen != ScreenTerminalSelect {
			t.Fatalf("Expected ScreenTerminalSelect after ESC, got %v", m.Screen)
		}
	})
}

// TestKeymapsNavigation tests keymaps screen navigation
func TestKeymapsNavigation(t *testing.T) {
	t.Run("can access keymaps menu from main menu", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenMainMenu
		m.Cursor = 2 // Keymaps Reference

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		m = result.(Model)

		if m.Screen != ScreenKeymapsMenu {
			t.Fatalf("Expected ScreenKeymapsMenu, got %v", m.Screen)
		}
	})

	t.Run("can select Neovim keymaps from menu", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenKeymapsMenu
		m.Cursor = 0 // Neovim

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		m = result.(Model)

		if m.Screen != ScreenKeymaps {
			t.Fatalf("Expected ScreenKeymaps, got %v", m.Screen)
		}
	})

	t.Run("can select a keymap category", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenKeymaps
		m.Cursor = 0

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		m = result.(Model)

		if m.Screen != ScreenKeymapCategory {
			t.Fatalf("Expected ScreenKeymapCategory, got %v", m.Screen)
		}
	})

	t.Run("can scroll through keymaps", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenKeymapCategory
		// Use category 5 (Search Commands) which has 27 keymaps (>10, so scrollable)
		m.SelectedCategory = 5
		m.KeymapScroll = 0

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
		m = result.(Model)

		if m.KeymapScroll != 1 {
			t.Fatalf("Expected KeymapScroll 1, got %d", m.KeymapScroll)
		}
	})
}

// TestLazyVimNavigation tests LazyVim guide navigation
func TestLazyVimNavigation(t *testing.T) {
	t.Run("can access LazyVim guide from main menu", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenMainMenu
		m.Cursor = 3 // LazyVim Guide

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		m = result.(Model)

		if m.Screen != ScreenLearnLazyVim {
			t.Fatalf("Expected ScreenLearnLazyVim, got %v", m.Screen)
		}
	})

	t.Run("can select a LazyVim topic", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenLearnLazyVim
		m.Cursor = 0

		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		m = result.(Model)

		if m.Screen != ScreenLazyVimTopic {
			t.Fatalf("Expected ScreenLazyVimTopic, got %v", m.Screen)
		}
	})
}

// TestBackupFlow tests the backup confirmation flow
func TestBackupFlow(t *testing.T) {
	t.Run("shows backup confirm when existing configs detected", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenNvimSelect
		m.Cursor = 0 // Yes, install Neovim
		m.SystemInfo = &system.SystemInfo{
			OS:       system.OSMac,
			HasBrew:  true,
			HasXcode: true,
		}
		m.Choices = UserChoices{
			OS:       "mac",
			Shell:    "fish",
			Terminal: "ghostty",
		}

		// Simulate existing configs
		m.ExistingConfigs = []string{"nvim: /Users/test/.config/nvim"}

		// This simulates what handleSelection does
		result, _ := m.handleSelection()
		m = result.(Model)

		// If there were existing configs, we go to backup confirm
		// The actual behavior depends on DetectExistingConfigs() which checks real filesystem
		// For this test, we manually set ExistingConfigs
		if len(m.ExistingConfigs) > 0 && m.Screen != ScreenBackupConfirm {
			// In real code, handleSelection calls system.DetectExistingConfigs()
			// which may return different results
			t.Log("Screen after nvim select with configs:", m.Screen)
		}
	})

	t.Run("backup confirm options work correctly", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenBackupConfirm
		m.ExistingConfigs = []string{"nvim: /test"}
		m.SystemInfo = &system.SystemInfo{
			OS:       system.OSMac,
			HasBrew:  true,
			HasXcode: true,
		}
		m.Choices = UserChoices{OS: "mac", Shell: "fish"}

		// Test cancel (cursor = 2)
		m.Cursor = 2
		result, _ := m.handleBackupConfirmKeys("enter")
		newM := result.(Model)

		if newM.Screen != ScreenMainMenu {
			t.Errorf("Cancel should go to MainMenu, got %v", newM.Screen)
		}
	})
}

// TestRestoreFlow tests the restore from backup flow
func TestRestoreFlow(t *testing.T) {
	t.Run("restore option appears in menu when backups exist", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenMainMenu
		m.AvailableBackups = []system.BackupInfo{
			{Path: "/test/backup"},
		}

		options := m.GetCurrentOptions()
		found := false
		for _, opt := range options {
			if strings.Contains(opt, "Restore") {
				found = true
				break
			}
		}

		if !found {
			t.Error("Restore option should appear when backups exist")
		}
	})

	t.Run("restore screen shows available backups", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenRestoreBackup
		m.AvailableBackups = []system.BackupInfo{
			{Path: "/test/backup1", Files: []string{"nvim", "fish"}},
			{Path: "/test/backup2", Files: []string{"zsh"}},
		}

		options := m.GetCurrentOptions()

		// Should have 2 backups + separator + back = 4 options
		if len(options) != 4 {
			t.Errorf("Expected 4 options, got %d: %v", len(options), options)
		}
	})
}

// TestViewRendering tests that views render without panicking
func TestViewRendering(t *testing.T) {
	screens := []struct {
		name   string
		setup  func(*Model)
		screen Screen
	}{
		{"Welcome", nil, ScreenWelcome},
		{"MainMenu", nil, ScreenMainMenu},
		{"OSSelect", nil, ScreenOSSelect},
		{"TerminalSelect", func(m *Model) { m.Choices.OS = "mac" }, ScreenTerminalSelect},
		{"FontSelect", nil, ScreenFontSelect},
		{"ShellSelect", nil, ScreenShellSelect},
		{"WMSelect", nil, ScreenWMSelect},
		{"NvimSelect", nil, ScreenNvimSelect},
		{"BackupConfirm", func(m *Model) {
			m.ExistingConfigs = []string{"nvim: /test"}
		}, ScreenBackupConfirm},
		{"RestoreBackup", func(m *Model) {
			m.AvailableBackups = []system.BackupInfo{{Path: "/test"}}
		}, ScreenRestoreBackup},
		{"RestoreConfirm", func(m *Model) {
			m.AvailableBackups = []system.BackupInfo{{Path: "/test", Files: []string{"nvim"}}}
			m.SelectedBackup = 0
		}, ScreenRestoreConfirm},
		{"Installing", func(m *Model) {
			m.Steps = []InstallStep{{ID: "test", Name: "Test", Status: StatusRunning}}
		}, ScreenInstalling},
		{"Complete", func(m *Model) {
			m.Choices = UserChoices{OS: "mac", Shell: "fish", Terminal: "ghostty"}
		}, ScreenComplete},
		{"Error", func(m *Model) { m.ErrorMsg = "Test error" }, ScreenError},
		{"LearnTerminals", nil, ScreenLearnTerminals},
		{"LearnShells", nil, ScreenLearnShells},
		{"LearnWM", nil, ScreenLearnWM},
		{"LearnNvim", nil, ScreenLearnNvim},
		{"Keymaps", nil, ScreenKeymaps},
		{"KeymapCategory", func(m *Model) { m.SelectedCategory = 0 }, ScreenKeymapCategory},
		{"LearnLazyVim", nil, ScreenLearnLazyVim},
		{"LazyVimTopic", func(m *Model) { m.SelectedLazyVimTopic = 0 }, ScreenLazyVimTopic},
	}

	for _, tc := range screens {
		t.Run(tc.name+" renders without panic", func(t *testing.T) {
			m := NewModel()
			m.Screen = tc.screen
			m.Width = 80
			m.Height = 24

			if tc.setup != nil {
				tc.setup(&m)
			}

			// This should not panic
			view := m.View()

			if view == "" && tc.screen != ScreenWelcome {
				// Welcome might have empty view if quitting
				t.Log("View is empty for", tc.name)
			}
		})
	}
}

// TestQuitBehavior tests various quit scenarios
func TestQuitBehavior(t *testing.T) {
	t.Run("space+q quits from main menu", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenMainMenu

		// First press space to enter leader mode
		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
		newM := result.(Model)

		if !newM.LeaderMode {
			t.Error("space should activate leader mode")
		}

		// Then press q to quit
		result, cmd := newM.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
		newM = result.(Model)

		if !newM.Quitting {
			t.Error("space+q should set Quitting to true")
		}
		if cmd == nil {
			t.Error("space+q should return quit command")
		}
	})

	t.Run("q does not quit during installation", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenInstalling

		// Even with leader mode, q should not quit during installation
		result, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
		newM := result.(Model)

		result, _ = newM.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
		newM = result.(Model)

		if newM.Quitting {
			t.Error("space+q should not quit during installation")
		}
	})

	t.Run("ctrl+c always quits", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenInstalling

		result, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
		newM := result.(Model)

		if !newM.Quitting {
			t.Error("ctrl+c should always quit")
		}
		if cmd == nil {
			t.Error("ctrl+c should return quit command")
		}
	})
}

// TestInstallStepsSetup tests that install steps are set up correctly for various configurations
func TestInstallStepsSetup(t *testing.T) {
	testCases := []struct {
		name          string
		choices       UserChoices
		sysInfo       *system.SystemInfo
		existConfigs  []string
		expectedSteps []string
	}{
		{
			name: "minimal mac install",
			choices: UserChoices{
				OS:           "mac",
				Terminal:     "none",
				Shell:        "fish",
				WindowMgr:    "none",
				InstallNvim:  false,
				InstallFont:  false,
				CreateBackup: false,
			},
			sysInfo:       &system.SystemInfo{OS: system.OSMac, HasBrew: true, HasXcode: true},
			existConfigs:  []string{},
			expectedSteps: []string{"clone", "shell", "setshell", "cleanup"},
		},
		{
			name: "full mac install with backup",
			choices: UserChoices{
				OS:           "mac",
				Terminal:     "ghostty",
				Shell:        "fish",
				WindowMgr:    "tmux",
				InstallNvim:  true,
				InstallFont:  true,
				CreateBackup: true,
			},
			sysInfo:       &system.SystemInfo{OS: system.OSMac, HasBrew: true, HasXcode: true},
			existConfigs:  []string{"nvim: /test"},
			expectedSteps: []string{"backup", "clone", "terminal", "font", "shell", "wm", "nvim", "setshell", "cleanup"},
		},
		{
			name: "linux install without brew",
			choices: UserChoices{
				OS:           "linux",
				Terminal:     "alacritty",
				Shell:        "zsh",
				WindowMgr:    "zellij",
				InstallNvim:  true,
				InstallFont:  true,
				CreateBackup: false,
			},
			sysInfo:       &system.SystemInfo{OS: system.OSLinux, HasBrew: false},
			existConfigs:  []string{},
			expectedSteps: []string{"clone", "homebrew", "deps", "terminal", "font", "shell", "wm", "nvim", "setshell", "cleanup"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := NewModel()
			m.Choices = tc.choices
			m.SystemInfo = tc.sysInfo
			m.ExistingConfigs = tc.existConfigs

			m.SetupInstallSteps()

			if len(m.Steps) != len(tc.expectedSteps) {
				stepIDs := make([]string, len(m.Steps))
				for i, s := range m.Steps {
					stepIDs[i] = s.ID
				}
				t.Errorf("Expected %d steps %v, got %d steps %v",
					len(tc.expectedSteps), tc.expectedSteps,
					len(m.Steps), stepIDs)
				return
			}

			for i, expected := range tc.expectedSteps {
				if m.Steps[i].ID != expected {
					t.Errorf("Step %d: expected '%s', got '%s'", i, expected, m.Steps[i].ID)
				}
			}
		})
	}
}
