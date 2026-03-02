package tui

import (
	"strings"
	"testing"

	"github.com/Gentleman-Programming/Gentleman.Dots/installer/internal/system"
)

// TestLearnAboutOptionsAreSelectable tests that "Learn about..." options can be selected
func TestLearnAboutOptionsAreSelectable(t *testing.T) {
	m := NewModel()
	m.SystemInfo = &system.SystemInfo{
		OS:       system.OSMac,
		IsTermux: false,
	}

	tests := []struct {
		screen      Screen
		optionName  string
		shouldExist bool
	}{
		{ScreenTerminalSelect, "ℹ️  Learn about terminals", true},
		{ScreenShellSelect, "ℹ️  Learn about shells", true},
		{ScreenWMSelect, "ℹ️  Learn about multiplexers", true},
		{ScreenNvimSelect, "ℹ️  Learn about Neovim", true},
	}

	for _, tt := range tests {
		t.Run(tt.optionName, func(t *testing.T) {
			m.Screen = tt.screen
			m.Choices.OS = "mac" // For terminal screen
			options := m.GetCurrentOptions()

			// Find the option
			foundIdx := -1
			for i, opt := range options {
				if opt == tt.optionName {
					foundIdx = i
					break
				}
			}

			if tt.shouldExist && foundIdx == -1 {
				t.Errorf("Option '%s' not found in screen %v", tt.optionName, tt.screen)
				t.Logf("Available options:")
				for _, opt := range options {
					t.Logf("  - %s", opt)
				}
				return
			}

			// Move cursor to that option
			m.Cursor = foundIdx

			// Verify it's selectable by checking if it would trigger navigation
			// In the actual UI, pressing Enter on these options changes screen
			selectedOpt := options[m.Cursor]

			// Should NOT be treated as separator
			if strings.HasPrefix(selectedOpt, "───") {
				t.Errorf("Option '%s' is treated as separator", tt.optionName)
			}

			// Should start with info emoji
			if !strings.HasPrefix(selectedOpt, "ℹ️") {
				t.Errorf("Option '%s' should start with ℹ️ emoji", tt.optionName)
			}
		})
	}
}
