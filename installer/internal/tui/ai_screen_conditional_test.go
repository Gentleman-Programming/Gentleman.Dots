package tui

import (
	"testing"

	"github.com/Gentleman-Programming/Gentleman.Dots/installer/internal/system"
)

// TestAIScreenWithNeovim tests that AI screen shows correct options when Neovim is selected
func TestAIScreenWithNeovim(t *testing.T) {
	m := NewModel()
	m.SystemInfo = &system.SystemInfo{
		OS:       system.OSMac,
		IsTermux: false,
	}
	m.Screen = ScreenAIAssistants
	m.Choices.InstallNvim = true

	options := m.GetCurrentOptions()

	// Should have: info note (2 lines) + blank line + OpenCode + 3 unavailable + separator + skip + docs link
	expectedMinOptions := 9
	if len(options) < expectedMinOptions {
		t.Errorf("Expected at least %d options with Neovim, got %d", expectedMinOptions, len(options))
	}

	// Check that informational note is present
	foundInfo := false
	for _, opt := range options {
		if opt == "ℹ️  Note: Claude Code is installed automatically with Neovim" {
			foundInfo = true
			break
		}
	}
	if !foundInfo {
		t.Error("Expected informational note about Claude Code, not found")
	}

	// Check that Claude Code is NOT in the selectable options
	for _, opt := range options {
		if opt == "[ ] Claude Code" || opt == "[✓] Claude Code" {
			t.Error("Claude Code should not appear as selectable option when Neovim is installed")
		}
	}

	// Check that "View AI Configuration Docs" link is present
	foundDocs := false
	for _, opt := range options {
		if opt == "📖 View AI Configuration Docs" {
			foundDocs = true
			break
		}
	}
	if !foundDocs {
		t.Error("Expected 'View AI Configuration Docs' link, not found")
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

	// Should have: Claude Code + OpenCode + 3 unavailable + separator + skip (no info note, no docs link)
	expectedMinOptions := 7
	if len(options) < expectedMinOptions {
		t.Errorf("Expected at least %d options without Neovim, got %d", expectedMinOptions, len(options))
	}

	// Check that informational note is NOT present
	for _, opt := range options {
		if opt == "ℹ️  Note: Claude Code is installed automatically with Neovim" {
			t.Error("Informational note should not appear when Neovim is not installed")
		}
	}

	// Check that Claude Code IS in the selectable options
	foundClaudeCode := false
	for _, opt := range options {
		if opt == "[ ] Claude Code" {
			foundClaudeCode = true
			break
		}
	}
	if !foundClaudeCode {
		t.Error("Claude Code should appear as selectable option when Neovim is not installed")
	}

	// Check that "View AI Configuration Docs" link is NOT present
	for _, opt := range options {
		if opt == "📖 View AI Configuration Docs" {
			t.Error("'View AI Configuration Docs' link should not appear when Neovim is not installed")
		}
	}
}
