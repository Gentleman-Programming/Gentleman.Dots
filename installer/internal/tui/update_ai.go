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
			// Skip non-selectable items (separator, info notes, blank lines)
			for m.Cursor >= 0 && m.Cursor < len(options) {
				opt := options[m.Cursor]
				if strings.HasPrefix(opt, "───") ||
					strings.HasPrefix(opt, "ℹ️") ||
					strings.HasPrefix(opt, "        ") || // Info note continuation
					opt == "" { // Blank line
					if m.Cursor > 0 {
						m.Cursor--
					} else {
						break
					}
				} else {
					break
				}
			}
		}
	case "down", "j":
		if m.Cursor < len(options)-1 {
			m.Cursor++
			// Skip non-selectable items (separator, info notes, blank lines)
			for m.Cursor < len(options) {
				opt := options[m.Cursor]
				if strings.HasPrefix(opt, "───") ||
					strings.HasPrefix(opt, "ℹ️") ||
					strings.HasPrefix(opt, "        ") || // Info note continuation
					opt == "" { // Blank line
					if m.Cursor < len(options)-1 {
						m.Cursor++
					} else {
						break
					}
				} else {
					break
				}
			}
		}
	case " ": // Space toggles selection
		selected := options[m.Cursor]

		// Only toggle if this is an AI assistant option (starts with checkbox)
		if strings.HasPrefix(selected, "[ ]") || strings.HasPrefix(selected, "[✓]") {
			// Extract the AI name from the option (format: "[ ] Name" or "[✓] Name")
			// Remove checkbox prefix - handle both "[ ] " and "[✓] "
			optionText := selected
			if strings.HasPrefix(optionText, "[ ] ") {
				optionText = strings.TrimPrefix(optionText, "[ ] ")
			} else if strings.HasPrefix(optionText, "[✓] ") {
				optionText = strings.TrimPrefix(optionText, "[✓] ")
			}
			optionText = strings.TrimSpace(optionText)

			// Remove " (Coming Soon)" suffix if present
			optionText = strings.TrimSuffix(optionText, " (Coming Soon)")

			// Find the AI assistant by name
			for _, ai := range m.AIAssistantsList {
				if ai.Name == optionText && ai.Available {
					m.SelectedAIAssistants[ai.ID] = !m.SelectedAIAssistants[ai.ID]
					break
				}
			}
		}
	case "enter":
		selected := options[m.Cursor]

		// Check if user selected "View AI Configuration Docs"
		if strings.Contains(selected, "View AI Configuration Docs") {
			// Open docs/ai-configuration.md in browser
			// Note: During installation, the repo is cloned to ~/Gentleman.Dots
			docsPath := "https://github.com/Gentleman-Programming/Gentleman.Dots/blob/main/docs/ai-configuration.md"
			if err := openURL(docsPath); err != nil {
				// Don't fail, just log - user can access docs later
				m.ErrorMsg = "Could not open docs (available after installation)"
			}
			// Stay on same screen, don't navigate
			return m, nil
		}

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
