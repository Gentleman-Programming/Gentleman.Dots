package tui

import (
	"strings"
	"testing"

	"github.com/Gentleman-Programming/Gentleman.Dots/installer/internal/system"
)

func TestNewModel(t *testing.T) {
	t.Run("should create model with initial state", func(t *testing.T) {
		m := NewModel()

		if m.Screen != ScreenWelcome {
			t.Errorf("Expected initial screen to be ScreenWelcome, got %v", m.Screen)
		}

		if m.Width != 80 {
			t.Errorf("Expected initial width 80, got %d", m.Width)
		}

		if m.Height != 24 {
			t.Errorf("Expected initial height 24, got %d", m.Height)
		}

		if m.SystemInfo == nil {
			t.Error("SystemInfo should not be nil")
		}

		if m.Cursor != 0 {
			t.Errorf("Expected cursor at 0, got %d", m.Cursor)
		}

		if m.ShowDetails {
			t.Error("ShowDetails should be false initially")
		}

		if m.Quitting {
			t.Error("Quitting should be false initially")
		}
	})
}

func TestGetCurrentOptions(t *testing.T) {
	m := NewModel()

	t.Run("should return OS options for OSSelect screen", func(t *testing.T) {
		m.Screen = ScreenOSSelect
		opts := m.GetCurrentOptions()

		if len(opts) != 2 {
			t.Errorf("Expected 2 OS options, got %d", len(opts))
		}
		// One of them should have "(detected)" based on current OS
		hasMac := strings.Contains(opts[0], "macOS")
		hasLinux := strings.Contains(opts[1], "Linux")
		if !hasMac || !hasLinux {
			t.Errorf("Unexpected OS options: %v", opts)
		}
	})

	t.Run("should return terminal options for mac", func(t *testing.T) {
		m.Screen = ScreenTerminalSelect
		m.Choices.OS = "mac"
		opts := m.GetCurrentOptions()

		// Should have: Alacritty, WezTerm, Kitty, Ghostty, None, separator, Learn
		if len(opts) != 7 {
			t.Errorf("Expected 7 terminal options for mac (including separator and learn), got %d", len(opts))
		}
		// Should include Kitty on mac
		hasKitty := false
		for _, opt := range opts {
			if opt == "Kitty" {
				hasKitty = true
				break
			}
		}
		if !hasKitty {
			t.Error("Mac should have Kitty option")
		}
	})

	t.Run("should return terminal options for linux without kitty", func(t *testing.T) {
		m.Screen = ScreenTerminalSelect
		m.Choices.OS = "linux"
		opts := m.GetCurrentOptions()

		// Should have: Alacritty, WezTerm, Ghostty, None, separator, Learn
		if len(opts) != 6 {
			t.Errorf("Expected 6 terminal options for linux (including separator and learn), got %d", len(opts))
		}
		// Should NOT include Kitty on linux
		for _, opt := range opts {
			if opt == "Kitty" {
				t.Error("Linux should not have Kitty option")
			}
		}
	})

	t.Run("should return shell options", func(t *testing.T) {
		m.Screen = ScreenShellSelect
		opts := m.GetCurrentOptions()

		// Should have: Fish, Zsh, Nushell, separator, Learn
		if len(opts) != 5 {
			t.Errorf("Expected 5 shell options (including separator and learn), got %d", len(opts))
		}
		expected := []string{"Fish", "Zsh", "Nushell"}
		for i, exp := range expected {
			if opts[i] != exp {
				t.Errorf("Expected %s at position %d, got %s", exp, i, opts[i])
			}
		}
	})

	t.Run("should return WM options", func(t *testing.T) {
		m.Screen = ScreenWMSelect
		opts := m.GetCurrentOptions()

		// Should have: Tmux, Zellij, None, separator, Learn
		if len(opts) != 5 {
			t.Errorf("Expected 5 WM options (including separator and learn), got %d", len(opts))
		}
	})

	t.Run("should return empty for non-selection screens", func(t *testing.T) {
		m.Screen = ScreenInstalling
		opts := m.GetCurrentOptions()

		if len(opts) != 0 {
			t.Errorf("Expected 0 options for installing screen, got %d", len(opts))
		}
	})
}

func TestGetScreenTitle(t *testing.T) {
	m := NewModel()

	tests := []struct {
		screen   Screen
		expected string
	}{
		{ScreenWelcome, "Welcome to Gentleman.Dots Installer"},
		{ScreenOSSelect, "Step 1: Select Your Operating System"},
		{ScreenTerminalSelect, "Step 2: Choose Terminal Emulator"},
		{ScreenShellSelect, "Step 4: Choose Your Shell"},
		{ScreenComplete, "Installation Complete!"},
		{ScreenError, "Error"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			m.Screen = tt.screen
			title := m.GetScreenTitle()
			if title != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, title)
			}
		})
	}
}

func TestSetupInstallSteps(t *testing.T) {
	t.Run("should setup minimal steps", func(t *testing.T) {
		m := NewModel()
		m.SystemInfo = &system.SystemInfo{
			OS:       system.OSMac,
			HasBrew:  true,
			HasXcode: true,
		}
		m.Choices = UserChoices{
			OS:          "mac",
			Terminal:    "none",
			InstallFont: false,
			Shell:       "fish",
			WindowMgr:   "none",
			InstallNvim: false,
		}

		m.SetupInstallSteps()

		// Should have: clone, shell, cleanup
		if len(m.Steps) < 3 {
			t.Errorf("Expected at least 3 steps, got %d", len(m.Steps))
		}

		// First step should be clone
		if m.Steps[0].ID != "clone" {
			t.Errorf("First step should be 'clone', got '%s'", m.Steps[0].ID)
		}

		// Last step should be cleanup
		lastStep := m.Steps[len(m.Steps)-1]
		if lastStep.ID != "cleanup" {
			t.Errorf("Last step should be 'cleanup', got '%s'", lastStep.ID)
		}
	})

	t.Run("should include homebrew step when not installed", func(t *testing.T) {
		m := NewModel()
		m.SystemInfo = &system.SystemInfo{
			OS:       system.OSMac,
			HasBrew:  false, // Not installed
			HasXcode: true,
		}
		m.Choices = UserChoices{
			OS:          "mac",
			Terminal:    "none",
			Shell:       "fish",
			WindowMgr:   "none",
			InstallNvim: false,
		}

		m.SetupInstallSteps()

		hasBrewStep := false
		for _, step := range m.Steps {
			if step.ID == "homebrew" {
				hasBrewStep = true
				break
			}
		}

		if !hasBrewStep {
			t.Error("Should include homebrew step when brew is not installed")
		}
	})

	t.Run("should include terminal step when selected", func(t *testing.T) {
		m := NewModel()
		m.SystemInfo = &system.SystemInfo{
			OS:       system.OSMac,
			HasBrew:  true,
			HasXcode: true,
		}
		m.Choices = UserChoices{
			OS:          "mac",
			Terminal:    "ghostty",
			Shell:       "fish",
			WindowMgr:   "none",
			InstallNvim: false,
		}

		m.SetupInstallSteps()

		hasTerminalStep := false
		for _, step := range m.Steps {
			if step.ID == "terminal" {
				hasTerminalStep = true
				break
			}
		}

		if !hasTerminalStep {
			t.Error("Should include terminal step when terminal is selected")
		}
	})

	t.Run("should include nvim step when selected", func(t *testing.T) {
		m := NewModel()
		m.SystemInfo = &system.SystemInfo{
			OS:       system.OSMac,
			HasBrew:  true,
			HasXcode: true,
		}
		m.Choices = UserChoices{
			OS:          "mac",
			Terminal:    "none",
			Shell:       "fish",
			WindowMgr:   "none",
			InstallNvim: true,
		}

		m.SetupInstallSteps()

		hasNvimStep := false
		for _, step := range m.Steps {
			if step.ID == "nvim" {
				hasNvimStep = true
				break
			}
		}

		if !hasNvimStep {
			t.Error("Should include nvim step when nvim is selected")
		}
	})

	t.Run("should include all steps for full installation", func(t *testing.T) {
		m := NewModel()
		m.SystemInfo = &system.SystemInfo{
			OS:       system.OSLinux,
			HasBrew:  false,
			HasXcode: false,
		}
		m.Choices = UserChoices{
			OS:          "linux",
			Terminal:    "ghostty",
			InstallFont: true,
			Shell:       "zsh",
			WindowMgr:   "tmux",
			InstallNvim: true,
		}

		m.SetupInstallSteps()

		expectedSteps := []string{"clone", "homebrew", "deps", "terminal", "font", "shell", "wm", "nvim", "setshell", "cleanup"}

		if len(m.Steps) != len(expectedSteps) {
			t.Errorf("Expected %d steps, got %d", len(expectedSteps), len(m.Steps))
		}

		for i, expected := range expectedSteps {
			if i < len(m.Steps) && m.Steps[i].ID != expected {
				t.Errorf("Step %d: expected '%s', got '%s'", i, expected, m.Steps[i].ID)
			}
		}
	})
}

func TestStepStatus(t *testing.T) {
	t.Run("status constants should have correct values", func(t *testing.T) {
		if StatusPending != 0 {
			t.Error("StatusPending should be 0")
		}
		if StatusRunning != 1 {
			t.Error("StatusRunning should be 1")
		}
		if StatusDone != 2 {
			t.Error("StatusDone should be 2")
		}
		if StatusFailed != 3 {
			t.Error("StatusFailed should be 3")
		}
		if StatusSkipped != 4 {
			t.Error("StatusSkipped should be 4")
		}
	})
}

func TestScreen(t *testing.T) {
	t.Run("screen constants should exist", func(t *testing.T) {
		// Just verify the screens exist and can be used
		screens := []Screen{
			ScreenWelcome,
			ScreenMainMenu,
			ScreenOSSelect,
			ScreenTerminalSelect,
			ScreenFontSelect,
			ScreenShellSelect,
			ScreenWMSelect,
			ScreenNvimSelect,
			ScreenInstalling,
			ScreenComplete,
			ScreenError,
			ScreenLearnTerminals,
			ScreenLearnShells,
			ScreenLearnWM,
			ScreenLearnNvim,
			ScreenKeymaps,
			ScreenKeymapCategory,
			ScreenLearnLazyVim,
			ScreenLazyVimTopic,
			ScreenBackupConfirm,
			ScreenRestoreBackup,
			ScreenRestoreConfirm,
		}

		// Verify we have all expected screens (including 3 new backup screens)
		if len(screens) != 22 {
			t.Errorf("Expected 22 screens, got %d", len(screens))
		}
	})
}

func TestBackupScreenOptions(t *testing.T) {
	t.Run("ScreenBackupConfirm should return correct options", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenBackupConfirm

		opts := m.GetCurrentOptions()

		if len(opts) != 3 {
			t.Errorf("Expected 3 options for BackupConfirm, got %d", len(opts))
		}

		// Check options contain expected text
		expectedOptions := []string{"Install with Backup", "Install without Backup", "Cancel"}
		for i, expected := range expectedOptions {
			found := false
			for _, opt := range opts {
				if containsString(opt, expected) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected option %d to contain '%s'", i, expected)
			}
		}
	})

	t.Run("ScreenRestoreBackup should return options based on available backups", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenRestoreBackup
		m.AvailableBackups = []system.BackupInfo{
			{Path: "/test/backup1", Files: []string{"file1", "file2"}},
			{Path: "/test/backup2", Files: []string{"file3"}},
		}

		opts := m.GetCurrentOptions()

		// Should have: 2 backups + separator + Back = 4 options
		if len(opts) != 4 {
			t.Errorf("Expected 4 options for RestoreBackup with 2 backups, got %d", len(opts))
		}

		// Last option should be Back
		if opts[len(opts)-1] != "← Back" {
			t.Errorf("Last option should be '← Back', got '%s'", opts[len(opts)-1])
		}
	})

	t.Run("ScreenRestoreConfirm should return correct options", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenRestoreConfirm

		opts := m.GetCurrentOptions()

		if len(opts) != 3 {
			t.Errorf("Expected 3 options for RestoreConfirm, got %d", len(opts))
		}
	})
}

func TestBackupScreenTitles(t *testing.T) {
	m := NewModel()

	tests := []struct {
		screen   Screen
		expected string
	}{
		{ScreenBackupConfirm, "Existing Configs Detected"},
		{ScreenRestoreBackup, "Restore from Backup"},
		{ScreenRestoreConfirm, "Confirm Restore"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			m.Screen = tt.screen
			title := m.GetScreenTitle()
			if !containsString(title, tt.expected) {
				t.Errorf("Expected title to contain '%s', got '%s'", tt.expected, title)
			}
		})
	}
}

func TestMainMenuWithBackups(t *testing.T) {
	t.Run("should include restore option when backups exist", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenMainMenu
		m.AvailableBackups = []system.BackupInfo{
			{Path: "/test/backup"},
		}

		opts := m.GetCurrentOptions()

		hasRestoreOption := false
		for _, opt := range opts {
			if containsString(opt, "Restore from Backup") {
				hasRestoreOption = true
				break
			}
		}

		if !hasRestoreOption {
			t.Error("Main menu should include 'Restore from Backup' when backups exist")
		}
	})

	t.Run("should not include restore option when no backups", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenMainMenu
		m.AvailableBackups = []system.BackupInfo{}

		opts := m.GetCurrentOptions()

		for _, opt := range opts {
			if containsString(opt, "Restore from Backup") {
				t.Error("Main menu should not include 'Restore from Backup' when no backups exist")
			}
		}
	})
}

func TestSetupInstallStepsWithBackup(t *testing.T) {
	t.Run("should include backup step when CreateBackup is true", func(t *testing.T) {
		m := NewModel()
		m.SystemInfo = &system.SystemInfo{
			OS:       system.OSMac,
			HasBrew:  true,
			HasXcode: true,
		}
		m.Choices = UserChoices{
			OS:           "mac",
			Terminal:     "none",
			Shell:        "fish",
			WindowMgr:    "none",
			InstallNvim:  false,
			CreateBackup: true,
		}
		m.ExistingConfigs = []string{"nvim: /home/user/.config/nvim"}

		m.SetupInstallSteps()

		// First step should be backup
		if m.Steps[0].ID != "backup" {
			t.Errorf("First step should be 'backup' when CreateBackup is true, got '%s'", m.Steps[0].ID)
		}
	})

	t.Run("should not include backup step when CreateBackup is false", func(t *testing.T) {
		m := NewModel()
		m.SystemInfo = &system.SystemInfo{
			OS:       system.OSMac,
			HasBrew:  true,
			HasXcode: true,
		}
		m.Choices = UserChoices{
			OS:           "mac",
			Terminal:     "none",
			Shell:        "fish",
			WindowMgr:    "none",
			InstallNvim:  false,
			CreateBackup: false,
		}
		m.ExistingConfigs = []string{"nvim: /home/user/.config/nvim"}

		m.SetupInstallSteps()

		// First step should NOT be backup
		if m.Steps[0].ID == "backup" {
			t.Error("First step should not be 'backup' when CreateBackup is false")
		}
	})

	t.Run("should not include backup step when no existing configs", func(t *testing.T) {
		m := NewModel()
		m.SystemInfo = &system.SystemInfo{
			OS:       system.OSMac,
			HasBrew:  true,
			HasXcode: true,
		}
		m.Choices = UserChoices{
			OS:           "mac",
			Terminal:     "none",
			Shell:        "fish",
			WindowMgr:    "none",
			InstallNvim:  false,
			CreateBackup: true, // Even if true
		}
		m.ExistingConfigs = []string{} // Empty

		m.SetupInstallSteps()

		// First step should be clone, not backup
		if m.Steps[0].ID == "backup" {
			t.Error("Should not include backup step when ExistingConfigs is empty")
		}
	})
}

func TestUserChoicesBackupField(t *testing.T) {
	t.Run("UserChoices should have CreateBackup field", func(t *testing.T) {
		choices := UserChoices{
			OS:           "mac",
			Terminal:     "ghostty",
			InstallFont:  true,
			Shell:        "fish",
			WindowMgr:    "tmux",
			InstallNvim:  true,
			CreateBackup: true,
		}

		if !choices.CreateBackup {
			t.Error("CreateBackup should be true")
		}

		choices.CreateBackup = false
		if choices.CreateBackup {
			t.Error("CreateBackup should be false")
		}
	})
}

func TestModelBackupFields(t *testing.T) {
	t.Run("Model should have backup-related fields initialized", func(t *testing.T) {
		m := NewModel()

		if m.ExistingConfigs == nil {
			t.Error("ExistingConfigs should not be nil")
		}

		if m.AvailableBackups == nil {
			t.Error("AvailableBackups should not be nil")
		}

		if m.SelectedBackup != 0 {
			t.Errorf("SelectedBackup should be 0, got %d", m.SelectedBackup)
		}

		if m.BackupDir != "" {
			t.Errorf("BackupDir should be empty, got '%s'", m.BackupDir)
		}
	})
}

// Helper function
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
