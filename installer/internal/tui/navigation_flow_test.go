package tui

import (
	"testing"
)

// TestNavigationFlow_TerminalThenSkip simulates:
// 1. User selects Terminal (Ghostty)
// 2. User presses ESC to go back
// 3. User skips Terminal
// 4. User skips Shell, WM, Nvim
// 5. User selects AI Assistant
// Expected: Should show Terminal as skipped, not as selected
func TestNavigationFlow_TerminalThenSkip(t *testing.T) {
	m := NewModel()
	m.AIAssistantsList = GetAvailableAIAssistants()
	
	// Step 1: User selects Ghostty on Terminal screen
	m.Choices.Terminal = "ghostty"
	m.Choices.InstallFont = true
	
	// Step 2: User presses ESC and goes back
	// Step 3: User now selects "Skip this step" on Terminal
	m.SkippedSteps[ScreenTerminalSelect] = true
	// BUG: m.Choices.Terminal still has "ghostty" value!
	
	// User skips other steps
	m.SkippedSteps[ScreenShellSelect] = true
	m.SkippedSteps[ScreenWMSelect] = true
	m.SkippedSteps[ScreenNvimSelect] = true
	
	// User selects AI Assistant
	m.Choices.AIAssistants = []string{"opencode"}
	
	summary := m.GetInstallationSummary()
	
	// Terminal should show as skipped, NOT as selected
	hasTerminalSkipped := false
	hasTerminalSelected := false
	
	for _, item := range summary {
		if item == "✗ Terminal (skipped)" {
			hasTerminalSkipped = true
		}
		if item == "✓ Terminal: Ghostty" {
			hasTerminalSelected = true
		}
	}
	
	if !hasTerminalSkipped {
		t.Error("Terminal should show as skipped")
	}
	if hasTerminalSelected {
		t.Error("Terminal should NOT show as selected when it was skipped")
	}
	
	// Test configs to overwrite
	m.ExistingConfigs = []string{
		"ghostty: /home/user/.config/ghostty",
		"fish: /home/user/.config/fish",
	}
	
	configs := m.GetConfigsToOverwrite()
	
	// Should not include ghostty since terminal was skipped
	for _, cfg := range configs {
		if containsAny(cfg, []string{"ghostty", "terminal"}) {
			t.Errorf("Should not overwrite ghostty config when terminal was skipped, got: %s", cfg)
		}
	}
}

// TestNavigationFlow_ShellThenSkip simulates the same for Shell
func TestNavigationFlow_ShellThenSkip(t *testing.T) {
	m := NewModel()
	m.AIAssistantsList = GetAvailableAIAssistants()
	
	// User selects Zsh
	m.Choices.Shell = "zsh"
	
	// User presses ESC and goes back, then skips
	m.SkippedSteps[ScreenTerminalSelect] = true
	m.SkippedSteps[ScreenShellSelect] = true // NOW skipped
	m.SkippedSteps[ScreenWMSelect] = true
	m.SkippedSteps[ScreenNvimSelect] = true
	
	m.Choices.AIAssistants = []string{"opencode"}
	
	summary := m.GetInstallationSummary()
	
	// Shell should show as skipped
	hasShellSkipped := false
	hasShellSelected := false
	
	for _, item := range summary {
		if item == "✗ Shell (skipped)" {
			hasShellSkipped = true
		}
		if item == "✓ Shell: Zsh" {
			hasShellSelected = true
		}
	}
	
	if !hasShellSkipped {
		t.Error("Shell should show as skipped")
	}
	if hasShellSelected {
		t.Error("Shell should NOT show as selected when it was skipped")
	}
	
	// Test configs
	m.ExistingConfigs = []string{
		"zsh: /home/user/.zshrc",
		"oh-my-zsh: /home/user/.oh-my-zsh",
	}
	
	configs := m.GetConfigsToOverwrite()
	
	if len(configs) != 0 {
		t.Errorf("Should not overwrite any configs when shell was skipped, got: %v", configs)
	}
}

func containsAny(s string, substrs []string) bool {
	for _, substr := range substrs {
		if len(s) >= len(substr) && s[:len(substr)] == substr {
			return true
		}
	}
	return false
}
