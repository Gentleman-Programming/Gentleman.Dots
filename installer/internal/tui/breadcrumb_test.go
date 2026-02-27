package tui

import (
	"testing"
)

func TestBreadcrumbWithSkippedSteps(t *testing.T) {
	t.Run("breadcrumb shows checkmarks for completed steps", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenAIAssistants
		m.Choices.OS = "mac"
		m.Choices.Terminal = "ghostty"
		m.Choices.Shell = "fish"

		progress := m.renderStepProgress()

		// OS, Terminal, Font, Shell, WM, Nvim should show as done (✓)
		// AI Assistants should show as active (●)
		if !contains(progress, "✓ OS") {
			t.Error("OS should be marked as done")
		}
		if !contains(progress, "✓ Terminal") {
			t.Error("Terminal should be marked as done")
		}
		if !contains(progress, "● AI Assistants") {
			t.Error("AI Assistants should be marked as active")
		}
	})

	t.Run("breadcrumb with skipped terminal step", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenShellSelect
		m.Choices.OS = "mac"
		m.SkippedSteps = make(map[Screen]bool)
		m.SkippedSteps[ScreenTerminalSelect] = true

		progress := m.renderStepProgress()

		// Even skipped steps show as done (✓) in the breadcrumb
		// because we've passed them
		if !contains(progress, "✓ OS") {
			t.Error("OS should be marked as done")
		}
		if !contains(progress, "✓ Terminal") {
			t.Error("Terminal should be marked as done even if skipped")
		}
		if !contains(progress, "● Shell") {
			t.Error("Shell should be marked as active")
		}
	})

	t.Run("breadcrumb navigation back doesn't affect checkmarks", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenNvimSelect
		m.Choices.OS = "mac"
		m.Choices.Terminal = "ghostty"

		// User goes back to Terminal screen
		m.Screen = ScreenTerminalSelect

		progress := m.renderStepProgress()

		// Only OS should be done, Terminal is active, rest are pending
		if !contains(progress, "✓ OS") {
			t.Error("OS should be marked as done")
		}
		if !contains(progress, "● Terminal") {
			t.Error("Terminal should be active when navigating back")
		}
		if !contains(progress, "○ Shell") {
			t.Error("Shell should be pending")
		}
	})
}
