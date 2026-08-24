package tui

// GetGhosttyKeymaps returns all Ghostty keymaps organized by category
func GetGhosttyKeymaps() []KeymapCategory {
	return []KeymapCategory{
		{
			Name:        "Splits",
			Description: "Create and navigate splits",
			Keymaps: []Keymap{
				{Keys: "Alt+v", Description: "New split to the right", Mode: ""},
				{Keys: "Alt+d", Description: "New split below", Mode: ""},
				{Keys: "Alt+h", Description: "Go to split on the left", Mode: ""},
				{Keys: "Alt+j", Description: "Go to split below", Mode: ""},
				{Keys: "Alt+k", Description: "Go to split above", Mode: ""},
				{Keys: "Alt+l", Description: "Go to split on the right", Mode: ""},
			},
		},
		{
			Name:        "Resize Splits",
			Description: "Resize split panes",
			Keymaps: []Keymap{
				{Keys: "Ctrl+Shift+h", Description: "Resize split left (10px)", Mode: ""},
				{Keys: "Ctrl+Shift+j", Description: "Resize split up (10px)", Mode: ""},
				{Keys: "Ctrl+Shift+k", Description: "Resize split down (10px)", Mode: ""},
				{Keys: "Ctrl+Shift+l", Description: "Resize split right (10px)", Mode: ""},
			},
		},
		{
			Name:        "Tabs",
			Description: "Tab management (native Ghostty)",
			Keymaps: []Keymap{
				{Keys: "Cmd+t", Description: "New tab", Mode: "macOS"},
				{Keys: "Cmd+w", Description: "Close tab", Mode: "macOS"},
				{Keys: "Cmd+Shift+[", Description: "Previous tab", Mode: "macOS"},
				{Keys: "Cmd+Shift+]", Description: "Next tab", Mode: "macOS"},
				{Keys: "Cmd+1-9", Description: "Go to tab number", Mode: "macOS"},
			},
		},
		{
			Name:        "General",
			Description: "Miscellaneous keybindings",
			Keymaps: []Keymap{
				{Keys: "Cmd+k", Description: "Clear screen", Mode: "macOS"},
				{Keys: "Shift+Enter", Description: "Send literal newline (\\x1b\\r)", Mode: ""},
				{Keys: "Alt+s", Description: "Write screen to file & paste path", Mode: ""},
			},
		},
		{
			Name:        "Copy & Paste",
			Description: "Clipboard operations",
			Keymaps: []Keymap{
				{Keys: "Cmd+c", Description: "Copy selection", Mode: "macOS"},
				{Keys: "Cmd+v", Description: "Paste from clipboard", Mode: "macOS"},
				{Keys: "Cmd+Shift+c", Description: "Copy (alternate)", Mode: "macOS"},
				{Keys: "Cmd+Shift+v", Description: "Paste (alternate)", Mode: "macOS"},
			},
		},
		{
			Name:        "Font & Zoom",
			Description: "Adjust font size",
			Keymaps: []Keymap{
				{Keys: "Cmd++", Description: "Increase font size", Mode: "macOS"},
				{Keys: "Cmd+-", Description: "Decrease font size", Mode: "macOS"},
				{Keys: "Cmd+0", Description: "Reset font size", Mode: "macOS"},
			},
		},
		{
			Name:        "Scrollback",
			Description: "Navigate scrollback buffer",
			Keymaps: []Keymap{
				{Keys: "Cmd+↑", Description: "Scroll up", Mode: "macOS"},
				{Keys: "Cmd+↓", Description: "Scroll down", Mode: "macOS"},
				{Keys: "Cmd+Home", Description: "Scroll to top", Mode: "macOS"},
				{Keys: "Cmd+End", Description: "Scroll to bottom", Mode: "macOS"},
				{Keys: "PageUp", Description: "Page up", Mode: ""},
				{Keys: "PageDown", Description: "Page down", Mode: ""},
			},
		},
		{
			Name:        "macOS Integration",
			Description: "macOS-specific features",
			Keymaps: []Keymap{
				{Keys: "Left Option", Description: "Acts as Alt key (macos-option-as-alt)", Mode: "macOS"},
				{Keys: "Cmd+,", Description: "Open preferences", Mode: "macOS"},
				{Keys: "Cmd+n", Description: "New window", Mode: "macOS"},
				{Keys: "Cmd+q", Description: "Quit Ghostty", Mode: "macOS"},
			},
		},
	}
}
