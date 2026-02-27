package tui

import (
	"strings"
	"testing"

	"github.com/Gentleman-Programming/Gentleman.Dots/installer/internal/system"
)

// TestAIScreenWithNeovim tests that AI screen shows correct options when Neovim is selected
func TestAIScreenWithNeovim(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIAssistants
	m.Choices.InstallNvim = true

	opts := m.GetCurrentOptions()

	// When Neovim is installed, Claude Code, Gemini CLI, and Copilot CLI should be hidden
	// Expected: Info header + 3 bullets + blank + OpenCode + separator + skip + docs = 9 options
	expectedCount := 9
	if len(opts) != expectedCount {
		t.Errorf("Expected %d options when Neovim is installed, got %d", expectedCount, len(opts))
		t.Logf("Options:")
		for i, opt := range opts {
			t.Logf("  %d: %s", i, opt)
		}
	}

	// Verify Claude Code, Gemini CLI, and Copilot CLI are NOT selectable (checkbox format)
	for _, opt := range opts {
		if strings.Contains(opt, "[ ] Claude Code") || strings.Contains(opt, "[✓] Claude Code") {
			t.Error("Claude Code should not appear as selectable option when Neovim is installed")
		}
		if strings.Contains(opt, "[ ] Gemini CLI") || strings.Contains(opt, "[✓] Gemini CLI") {
			t.Error("Gemini CLI should not appear as selectable option when Neovim is installed")
		}
		if strings.Contains(opt, "[ ] GitHub Copilot CLI") || strings.Contains(opt, "[✓] GitHub Copilot CLI") {
			t.Error("GitHub Copilot CLI should not appear as selectable option when Neovim is installed")
		}
	}

	// Verify informational note is present
	foundNote := false
	for _, opt := range opts {
		if strings.HasPrefix(opt, "ℹ️  Note:") {
			foundNote = true
			break
		}
	}
	if !foundNote {
		t.Error("Should show informational note when Neovim is installed")
	}
}

// TestAIScreenWithoutNeovim tests that AI screen shows correct options when Neovim is skipped
func TestAIScreenWithoutNeovim(t *testing.T) {
	m := NewModel()
	m.SystemInfo = &system.SystemInfo{
		OS:       system.OSMac,
		IsTermux: false,
	}
	m.Screen = ScreenAIAssistants
	m.Choices.InstallNvim = false

	options := m.GetCurrentOptions()

	// Should have: Claude Code + Gemini + Copilot + OpenCode + separator + skip = 6
	expectedMinOptions := 6
	if len(options) < expectedMinOptions {
		t.Errorf("Expected at least %d options without Neovim, got %d", expectedMinOptions, len(options))
	}

	// Check that informational note is NOT present
	for _, opt := range options {
		if strings.HasPrefix(opt, "ℹ️  Note:") {
			t.Error("Informational note should not appear when Neovim is not installed")
		}
	}

	// Check that all AI assistants ARE in the selectable options
	foundClaudeCode := false
	foundGemini := false
	foundCopilot := false
	for _, opt := range options {
		if opt == "[ ] Claude Code" {
			foundClaudeCode = true
		}
		if opt == "[ ] Gemini CLI" {
			foundGemini = true
		}
		if opt == "[ ] GitHub Copilot CLI" {
			foundCopilot = true
		}
	}
	if !foundClaudeCode {
		t.Error("Claude Code should appear as selectable option when Neovim is not installed")
	}
	if !foundGemini {
		t.Error("Gemini CLI should appear as selectable option when Neovim is not installed")
	}
	if !foundCopilot {
		t.Error("GitHub Copilot CLI should appear as selectable option when Neovim is not installed")
	}

	// Check that "View AI Configuration Docs" link is NOT present
	for _, opt := range options {
		if opt == "📖 View AI Configuration Docs" {
			t.Error("'View AI Configuration Docs' link should not appear when Neovim is not installed")
		}
	}
}
