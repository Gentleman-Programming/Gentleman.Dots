package tui

import (
	"testing"

	"github.com/Gentleman-Programming/Gentleman.Dots/installer/internal/system"
	tea "github.com/charmbracelet/bubbletea"
)

// mockSystemInfo returns a basic system info for testing
func mockSystemInfo() *system.SystemInfo {
	return &system.SystemInfo{
		OS:     system.OSMac,
		OSName: "macOS",
	}
}

func TestLeaderNavigationFlow(t *testing.T) {
	t.Run("Full Nvim Configuration Flow", func(t *testing.T) {
		m := NewModel()
		m.SystemInfo = mockSystemInfo()
		m.Screen = ScreenNvimSelect
		m.Cursor = 0 // "Yes" option

		// 1. Press Enter on Nvim Select (Yes)
		model, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		m = model.(Model)

		if m.Screen != ScreenExperienceSelect {
			t.Errorf("Expected ScreenExperienceSelect, got %v", m.Screen)
		}

		// 2. Select "Advanced" (Cursor 2) and Press Enter
		m.Cursor = 2
		model, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		m = model.(Model)

		if m.Choices.Experience != "advanced" {
			t.Errorf("Expected experience 'advanced', got %s", m.Choices.Experience)
		}
		if m.Screen != ScreenLeaderKeySelect {
			t.Errorf("Expected ScreenLeaderKeySelect, got %v", m.Screen)
		}

		// 3. Select "Comma" (Cursor 1) and Press Enter
		m.Cursor = 1
		model, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		m = model.(Model)

		if m.Choices.LeaderKey != "comma" {
			t.Errorf("Expected leader 'comma', got %s", m.Choices.LeaderKey)
		}
		// Should have proceeded to either Backup or Installing
		if m.Screen != ScreenBackupConfirm && m.Screen != ScreenInstalling {
			t.Errorf("Expected ScreenBackupConfirm or ScreenInstalling, got %v", m.Screen)
		}
	})

	t.Run("Skip Nvim Flow", func(t *testing.T) {
		m := NewModel()
		m.SystemInfo = mockSystemInfo()
		m.Screen = ScreenNvimSelect
		m.Cursor = 1 // "No" option

		// Press Enter on Nvim Select (No)
		model, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
		m = model.(Model)

		// Should skip Experience and Leader screens
		if m.Screen == ScreenExperienceSelect || m.Screen == ScreenLeaderKeySelect {
			t.Errorf("Should have skipped nvim sub-screens, got %v", m.Screen)
		}
	})

	t.Run("Go Back Navigation", func(t *testing.T) {
		m := NewModel()
		m.SystemInfo = mockSystemInfo()
		m.Screen = ScreenLeaderKeySelect

		// Press Esc on Leader Key screen
		model, _ := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
		m = model.(Model)

		if m.Screen != ScreenExperienceSelect {
			t.Errorf("Expected ScreenExperienceSelect after back from Leader, got %v", m.Screen)
		}

		// Press Esc on Experience screen
		model, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
		m = model.(Model)

		if m.Screen != ScreenNvimSelect {
			t.Errorf("Expected ScreenNvimSelect after back from Experience, got %v", m.Screen)
		}
	})
}
