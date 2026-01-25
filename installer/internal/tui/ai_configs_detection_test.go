package tui

import (
	"testing"
)

// TestGetConfigsToOverwrite_OpenCodeConfig tests OpenCode config detection
func TestGetConfigsToOverwrite_OpenCodeConfig(t *testing.T) {
	m := NewModel()

	// User selects ONLY OpenCode AI Assistant
	m.Choices.OS = "darwin"
	m.Choices.AIAssistants = []string{"opencode"}
	m.SkippedSteps = make(map[Screen]bool)
	m.SkippedSteps[ScreenTerminalSelect] = true
	m.SkippedSteps[ScreenShellSelect] = true
	m.SkippedSteps[ScreenWMSelect] = true
	m.SkippedSteps[ScreenNvimSelect] = true
	m.SkippedSteps[ScreenAIAssistants] = false // User selected OpenCode

	// Simulate existing OpenCode config
	m.ExistingConfigs = []string{
		"opencode: /Users/test/.config/opencode",
	}

	configs := m.GetConfigsToOverwrite()

	// Should include OpenCode config since user selected it
	if len(configs) != 1 {
		t.Fatalf("Expected 1 config to overwrite, got %d: %v", len(configs), configs)
	}

	if configs[0] != "opencode: /Users/test/.config/opencode" {
		t.Errorf("Expected OpenCode config, got: %s", configs[0])
	}
}

// TestGetConfigsToOverwrite_AIAssistantSkipped tests when AI is skipped
func TestGetConfigsToOverwrite_AIAssistantSkipped(t *testing.T) {
	m := NewModel()

	// User SKIPS AI Assistants
	m.Choices.OS = "darwin"
	m.Choices.Shell = "zsh"
	m.Choices.AIAssistants = []string{} // No AI selected
	m.SkippedSteps = make(map[Screen]bool)
	m.SkippedSteps[ScreenTerminalSelect] = true
	m.SkippedSteps[ScreenShellSelect] = false // Selected Zsh
	m.SkippedSteps[ScreenWMSelect] = true
	m.SkippedSteps[ScreenNvimSelect] = true
	m.SkippedSteps[ScreenAIAssistants] = true // Skipped AI

	// Simulate existing OpenCode config AND Zsh config
	m.ExistingConfigs = []string{
		"opencode: /Users/test/.config/opencode",
		"zsh: /Users/test/.zshrc",
	}

	configs := m.GetConfigsToOverwrite()

	// Should include ONLY Zsh, NOT OpenCode (because AI was skipped)
	if len(configs) != 1 {
		t.Fatalf("Expected 1 config to overwrite, got %d: %v", len(configs), configs)
	}

	if configs[0] != "zsh: /Users/test/.zshrc" {
		t.Errorf("Expected Zsh config only, got: %s", configs[0])
	}
}

// TestGetConfigsToOverwrite_MultipleAIAssistants tests multiple AI selection
func TestGetConfigsToOverwrite_MultipleAIAssistants(t *testing.T) {
	m := NewModel()

	// User selects BOTH OpenCode and Continue.dev
	m.Choices.OS = "darwin"
	m.Choices.AIAssistants = []string{"opencode", "continue"}
	m.SkippedSteps = make(map[Screen]bool)
	m.SkippedSteps[ScreenTerminalSelect] = true
	m.SkippedSteps[ScreenShellSelect] = true
	m.SkippedSteps[ScreenWMSelect] = true
	m.SkippedSteps[ScreenNvimSelect] = true
	m.SkippedSteps[ScreenAIAssistants] = false

	// Simulate existing configs for BOTH
	m.ExistingConfigs = []string{
		"opencode: /Users/test/.config/opencode",
		"continue: /Users/test/.continue",
	}

	configs := m.GetConfigsToOverwrite()

	// Should include BOTH configs
	if len(configs) != 2 {
		t.Fatalf("Expected 2 configs to overwrite, got %d: %v", len(configs), configs)
	}

	// Check both are present (order doesn't matter)
	hasOpenCode := false
	hasContinue := false
	for _, config := range configs {
		if config == "opencode: /Users/test/.config/opencode" {
			hasOpenCode = true
		}
		if config == "continue: /Users/test/.continue" {
			hasContinue = true
		}
	}

	if !hasOpenCode {
		t.Error("Expected OpenCode config in list")
	}
	if !hasContinue {
		t.Error("Expected Continue.dev config in list")
	}
}

// TestGetConfigsToOverwrite_AIAndShellMix tests AI + Shell selection
func TestGetConfigsToOverwrite_AIAndShellMix(t *testing.T) {
	m := NewModel()

	// User selects Zsh AND OpenCode
	m.Choices.OS = "darwin"
	m.Choices.Shell = "zsh"
	m.Choices.AIAssistants = []string{"opencode"}
	m.SkippedSteps = make(map[Screen]bool)
	m.SkippedSteps[ScreenTerminalSelect] = true
	m.SkippedSteps[ScreenShellSelect] = false // Zsh selected
	m.SkippedSteps[ScreenWMSelect] = true
	m.SkippedSteps[ScreenNvimSelect] = true
	m.SkippedSteps[ScreenAIAssistants] = false // OpenCode selected

	// Simulate existing configs
	m.ExistingConfigs = []string{
		"opencode: /Users/test/.config/opencode",
		"zsh: /Users/test/.zshrc",
		"fish: /Users/test/.config/fish", // User has Fish BUT chose Zsh
	}

	configs := m.GetConfigsToOverwrite()

	// Should include OpenCode + Zsh, but NOT Fish
	if len(configs) != 2 {
		t.Fatalf("Expected 2 configs to overwrite, got %d: %v", len(configs), configs)
	}

	hasOpenCode := false
	hasZsh := false
	hasFish := false
	for _, config := range configs {
		if config == "opencode: /Users/test/.config/opencode" {
			hasOpenCode = true
		}
		if config == "zsh: /Users/test/.zshrc" {
			hasZsh = true
		}
		if config == "fish: /Users/test/.config/fish" {
			hasFish = true
		}
	}

	if !hasOpenCode {
		t.Error("Expected OpenCode config in list")
	}
	if !hasZsh {
		t.Error("Expected Zsh config in list")
	}
	if hasFish {
		t.Error("Did NOT expect Fish config (user chose Zsh, not Fish)")
	}
}
