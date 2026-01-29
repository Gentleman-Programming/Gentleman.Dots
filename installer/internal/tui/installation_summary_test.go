package tui

import (
	"strings"
	"testing"
)

// TestGetInstallationSummary_AllComponents tests when all components are selected
func TestGetInstallationSummary_AllComponents(t *testing.T) {
	m := NewModel()
	m.Choices.Terminal = "alacritty"
	m.Choices.InstallFont = true
	m.Choices.Shell = "fish"
	m.Choices.WindowMgr = "tmux"
	m.Choices.InstallNvim = true
	m.Choices.AIAssistants = []string{"opencode"}
	m.AIAssistantsList = GetAvailableAIAssistants()

	summary := m.GetInstallationSummary()

	expected := []string{
		"✓ Terminal: Alacritty",
		"  └─ Iosevka Nerd Font",
		"✓ Shell: Fish",
		"✓ Multiplexer: Tmux",
		"✓ Neovim: LazyVim configuration",
		"✓ AI Assistant: OpenCode",
	}

	if len(summary) != len(expected) {
		t.Errorf("Expected %d items, got %d", len(expected), len(summary))
		t.Logf("Summary: %v", summary)
		return
	}

	for i, exp := range expected {
		if summary[i] != exp {
			t.Errorf("Item %d: expected '%s', got '%s'", i, exp, summary[i])
		}
	}
}

// TestGetInstallationSummary_OnlyAIAssistant tests when only AI Assistant is selected
func TestGetInstallationSummary_OnlyAIAssistant(t *testing.T) {
	m := NewModel()
	m.SkippedSteps[ScreenTerminalSelect] = true
	m.SkippedSteps[ScreenShellSelect] = true
	m.SkippedSteps[ScreenWMSelect] = true
	m.SkippedSteps[ScreenNvimSelect] = true
	m.Choices.AIAssistants = []string{"opencode"}
	m.AIAssistantsList = GetAvailableAIAssistants()

	summary := m.GetInstallationSummary()

	expected := []string{
		"✗ Terminal (skipped)",
		"✗ Shell (skipped)",
		"✗ Multiplexer (skipped)",
		"✗ Neovim (skipped)",
		"✓ AI Assistant: OpenCode",
	}

	if len(summary) != len(expected) {
		t.Errorf("Expected %d items, got %d", len(expected), len(summary))
		t.Logf("Summary: %v", summary)
		return
	}

	for i, exp := range expected {
		if summary[i] != exp {
			t.Errorf("Item %d: expected '%s', got '%s'", i, exp, summary[i])
		}
	}
}

// TestGetInstallationSummary_TerminalAndAI tests Terminal + AI Assistant
func TestGetInstallationSummary_TerminalAndAI(t *testing.T) {
	m := NewModel()
	m.Choices.Terminal = "ghostty"
	m.Choices.InstallFont = false
	m.SkippedSteps[ScreenShellSelect] = true
	m.SkippedSteps[ScreenWMSelect] = true
	m.SkippedSteps[ScreenNvimSelect] = true
	m.Choices.AIAssistants = []string{"opencode"}
	m.AIAssistantsList = GetAvailableAIAssistants()

	summary := m.GetInstallationSummary()

	expected := []string{
		"✓ Terminal: Ghostty",
		"✗ Shell (skipped)",
		"✗ Multiplexer (skipped)",
		"✗ Neovim (skipped)",
		"✓ AI Assistant: OpenCode",
	}

	if len(summary) != len(expected) {
		t.Errorf("Expected %d items, got %d", len(expected), len(summary))
		t.Logf("Summary: %v", summary)
		return
	}

	for i, exp := range expected {
		if summary[i] != exp {
			t.Errorf("Item %d: expected '%s', got '%s'", i, exp, summary[i])
		}
	}
}

// TestGetInstallationSummary_ShellAndAI tests Shell + AI Assistant
func TestGetInstallationSummary_ShellAndAI(t *testing.T) {
	m := NewModel()
	m.SkippedSteps[ScreenTerminalSelect] = true
	m.Choices.Shell = "zsh"
	m.SkippedSteps[ScreenWMSelect] = true
	m.SkippedSteps[ScreenNvimSelect] = true
	m.Choices.AIAssistants = []string{"opencode"}
	m.AIAssistantsList = GetAvailableAIAssistants()

	summary := m.GetInstallationSummary()

	expected := []string{
		"✗ Terminal (skipped)",
		"✓ Shell: Zsh",
		"✗ Multiplexer (skipped)",
		"✗ Neovim (skipped)",
		"✓ AI Assistant: OpenCode",
	}

	if len(summary) != len(expected) {
		t.Errorf("Expected %d items, got %d", len(expected), len(summary))
		t.Logf("Summary: %v", summary)
		return
	}

	for i, exp := range expected {
		if summary[i] != exp {
			t.Errorf("Item %d: expected '%s', got '%s'", i, exp, summary[i])
		}
	}
}

// TestGetInstallationSummary_AllSkipped tests when everything is skipped
func TestGetInstallationSummary_AllSkipped(t *testing.T) {
	m := NewModel()
	m.SkippedSteps[ScreenTerminalSelect] = true
	m.SkippedSteps[ScreenShellSelect] = true
	m.SkippedSteps[ScreenWMSelect] = true
	m.SkippedSteps[ScreenNvimSelect] = true
	m.SkippedSteps[ScreenAIAssistants] = true

	summary := m.GetInstallationSummary()

	expected := []string{
		"✗ Terminal (skipped)",
		"✗ Shell (skipped)",
		"✗ Multiplexer (skipped)",
		"✗ Neovim (skipped)",
		"✗ AI Assistants (skipped)",
	}

	if len(summary) != len(expected) {
		t.Errorf("Expected %d items, got %d", len(expected), len(summary))
		t.Logf("Summary: %v", summary)
		return
	}

	for i, exp := range expected {
		if summary[i] != exp {
			t.Errorf("Item %d: expected '%s', got '%s'", i, exp, summary[i])
		}
	}
}

// TestGetConfigsToOverwrite_OnlyRelevantConfigs tests config filtering
func TestGetConfigsToOverwrite_OnlyRelevantConfigs(t *testing.T) {
	m := NewModel()
	
	// Simulate existing configs
	m.ExistingConfigs = []string{
		"fish: /home/user/.config/fish",
		"zsh: /home/user/.zshrc",
		"oh-my-zsh: /home/user/.oh-my-zsh",
		"tmux: /home/user/.tmux.conf",
		"nvim: /home/user/.config/nvim",
	}
	
	// User only chose Fish and skipped everything else
	m.Choices.Shell = "fish"
	m.SkippedSteps[ScreenTerminalSelect] = true
	m.SkippedSteps[ScreenWMSelect] = true
	m.SkippedSteps[ScreenNvimSelect] = true
	m.SkippedSteps[ScreenAIAssistants] = true
	
	configs := m.GetConfigsToOverwrite()
	
	// Should only show fish config, not zsh, tmux, or nvim
	if len(configs) != 1 {
		t.Errorf("Expected 1 config to overwrite, got %d: %v", len(configs), configs)
		return
	}
	
	if !strings.Contains(configs[0], "fish") {
		t.Errorf("Expected fish config, got: %s", configs[0])
	}
}

// TestGetConfigsToOverwrite_ZshConfigs tests zsh-related configs
func TestGetConfigsToOverwrite_ZshConfigs(t *testing.T) {
	m := NewModel()
	
	// Simulate existing configs
	m.ExistingConfigs = []string{
		"fish: /home/user/.config/fish",
		"zsh: /home/user/.zshrc",
		"oh-my-zsh: /home/user/.oh-my-zsh",
		"zsh_p10k: /home/user/.p10k.zsh",
	}
	
	// User chose Zsh
	m.Choices.Shell = "zsh"
	m.SkippedSteps[ScreenTerminalSelect] = true
	m.SkippedSteps[ScreenWMSelect] = true
	m.SkippedSteps[ScreenNvimSelect] = true
	m.SkippedSteps[ScreenAIAssistants] = true
	
	configs := m.GetConfigsToOverwrite()
	
	// Should show zsh, oh-my-zsh, and zsh_p10k (but NOT fish)
	if len(configs) != 3 {
		t.Errorf("Expected 3 configs to overwrite, got %d: %v", len(configs), configs)
		return
	}
	
	hasZsh := false
	hasOhMyZsh := false
	hasP10k := false
	hasFish := false
	
	for _, cfg := range configs {
		if strings.Contains(cfg, "zsh:") {
			hasZsh = true
		}
		if strings.Contains(cfg, "oh-my-zsh") {
			hasOhMyZsh = true
		}
		if strings.Contains(cfg, "zsh_p10k") {
			hasP10k = true
		}
		if strings.Contains(cfg, "fish") {
			hasFish = true
		}
	}
	
	if !hasZsh || !hasOhMyZsh || !hasP10k {
		t.Error("Missing expected zsh-related configs")
	}
	if hasFish {
		t.Error("Should not include fish config when user chose zsh")
	}
}

// TestGetConfigsToOverwrite_NoConfigs tests when nothing will be overwritten
func TestGetConfigsToOverwrite_NoConfigs(t *testing.T) {
	m := NewModel()
	
	// User skipped everything
	m.SkippedSteps[ScreenTerminalSelect] = true
	m.SkippedSteps[ScreenShellSelect] = true
	m.SkippedSteps[ScreenWMSelect] = true
	m.SkippedSteps[ScreenNvimSelect] = true
	m.SkippedSteps[ScreenAIAssistants] = true
	
	// But has existing configs
	m.ExistingConfigs = []string{
		"fish: /home/user/.config/fish",
		"tmux: /home/user/.tmux.conf",
		"nvim: /home/user/.config/nvim",
	}
	
	configs := m.GetConfigsToOverwrite()
	
	// Should be empty since user skipped everything
	if len(configs) != 0 {
		t.Errorf("Expected 0 configs to overwrite, got %d: %v", len(configs), configs)
	}
}
