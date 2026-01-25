package tui

import (
	"strings"

	"github.com/Gentleman-Programming/Gentleman.Dots/installer/internal/system"
	tea "github.com/charmbracelet/bubbletea"
)

// handleAIAssistantsKeys handles keyboard input on the AI Assistants selection screen
func (m Model) handleAIAssistantsKeys(key string) (tea.Model, tea.Cmd) {
	options := m.GetCurrentOptions()

	switch key {
	case "up", "k":
		if m.Cursor > 0 {
			m.Cursor--
			// Skip separator lines
			if m.Cursor < len(options) && strings.HasPrefix(options[m.Cursor], "───") {
				if m.Cursor > 0 {
					m.Cursor--
				}
			}
		}
	case "down", "j":
		if m.Cursor < len(options)-1 {
			m.Cursor++
			// Skip separator lines
			if m.Cursor < len(options) && strings.HasPrefix(options[m.Cursor], "───") {
				if m.Cursor < len(options)-1 {
					m.Cursor++
				}
			}
		}
	case " ": // Space toggles selection
		if m.Cursor < len(m.AIAssistantsList) {
			ai := m.AIAssistantsList[m.Cursor]
			if ai.Available {
				m.SelectedAIAssistants[ai.ID] = !m.SelectedAIAssistants[ai.ID]
			}
		}
	case "enter":
		selected := options[m.Cursor]
		
		// Check if user selected "Skip this step"
		if strings.Contains(selected, "Skip this step") {
			m.SkippedSteps[ScreenAIAssistants] = true
			// No AI assistants selected
			m.Choices.AIAssistants = []string{}
		} else {
			// Enter confirms selection from any position (no need to navigate to "Confirm")
			// Convert selected map to slice for Choices
			m.Choices.AIAssistants = []string{}
			for id, isSelected := range m.SelectedAIAssistants {
				if isSelected {
					m.Choices.AIAssistants = append(m.Choices.AIAssistants, id)
				}
			}
			// CRITICAL: Clear skip flag when user confirms a selection (even if none selected)
			m.SkippedSteps[ScreenAIAssistants] = false
		}

		// Always show installation summary/confirmation screen
		// (previously called ScreenBackupConfirm, but now shows summary + backup if needed)
		m.ExistingConfigs = system.DetectExistingConfigs()
		m.Screen = ScreenBackupConfirm
		m.Cursor = 0
		return m, nil
	}

	return m, nil
}
