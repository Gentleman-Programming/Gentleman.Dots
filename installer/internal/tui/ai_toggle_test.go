package tui

import (
	"testing"

	"github.com/Gentleman-Programming/Gentleman.Dots/installer/internal/system"
)

// TestAIAssistantToggle tests that AI assistants can be selected and deselected
func TestAIAssistantToggle(t *testing.T) {
	m := NewModel()
	m.SystemInfo = &system.SystemInfo{
		OS:       system.OSMac,
		IsTermux: false,
	}
	m.Screen = ScreenAIAssistants
	m.Choices.InstallNvim = false // So Claude Code appears

	// Get options and find OpenCode
	options := m.GetCurrentOptions()

	// Find OpenCode in the list
	openCodeIdx := -1
	for i, opt := range options {
		if opt == "[ ] OpenCode" {
			openCodeIdx = i
			break
		}
	}

	if openCodeIdx == -1 {
		t.Fatal("OpenCode not found in options")
	}

	// Move cursor to OpenCode
	m.Cursor = openCodeIdx

	// Initial state - should NOT be selected
	if m.SelectedAIAssistants["opencode"] {
		t.Error("OpenCode should not be selected initially")
	}

	// Press space to SELECT
	updatedModel, _ := m.handleAIAssistantsKeys(" ")
	m = updatedModel.(Model)

	// Should now be selected
	if !m.SelectedAIAssistants["opencode"] {
		t.Error("OpenCode should be selected after pressing space")
	}

	// Verify the option now shows as checked
	options = m.GetCurrentOptions()
	if options[openCodeIdx] != "[✓] OpenCode" {
		t.Errorf("Expected '[✓] OpenCode', got '%s'", options[openCodeIdx])
	}

	// Press space again to DESELECT
	updatedModel, _ = m.handleAIAssistantsKeys(" ")
	m = updatedModel.(Model)

	// Should now be deselected
	if m.SelectedAIAssistants["opencode"] {
		t.Error("OpenCode should be deselected after pressing space again")
	}

	// Verify the option shows as unchecked
	options = m.GetCurrentOptions()
	if options[openCodeIdx] != "[ ] OpenCode" {
		t.Errorf("Expected '[ ] OpenCode', got '%s'", options[openCodeIdx])
	}
}
