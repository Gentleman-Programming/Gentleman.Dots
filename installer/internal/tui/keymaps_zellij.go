package tui

// GetZellijKeymaps returns all Zellij keymaps organized by category
func GetZellijKeymaps() []KeymapCategory {
	return []KeymapCategory{
		{
			Name:        "Mode Switching",
			Description: "Zellij uses modes like Vim (default: locked mode)",
			Keymaps: []Keymap{
				{Keys: "Ctrl+g", Description: "Toggle locked/normal mode", Mode: ""},
				{Keys: "p", Description: "Enter Pane mode", Mode: "normal"},
				{Keys: "t", Description: "Enter Tab mode", Mode: "normal"},
				{Keys: "r", Description: "Enter Resize mode", Mode: "normal"},
				{Keys: "m", Description: "Enter Move mode", Mode: "normal"},
				{Keys: "s", Description: "Enter Scroll mode", Mode: "normal"},
				{Keys: "o", Description: "Enter Session mode", Mode: "normal"},
				{Keys: "Esc/Enter", Description: "Return to locked mode", Mode: "any"},
			},
		},
		{
			Name:        "Pane Mode (Ctrl+g p)",
			Description: "Create, navigate, and manage panes",
			Keymaps: []Keymap{
				{Keys: "n", Description: "New pane", Mode: "pane"},
				{Keys: "r", Description: "New pane to the right", Mode: "pane"},
				{Keys: "d", Description: "New pane below", Mode: "pane"},
				{Keys: "x", Description: "Close pane", Mode: "pane"},
				{Keys: "f", Description: "Toggle fullscreen", Mode: "pane"},
				{Keys: "w", Description: "Toggle floating panes", Mode: "pane"},
				{Keys: "e", Description: "Toggle embed/floating", Mode: "pane"},
				{Keys: "z", Description: "Toggle pane frames", Mode: "pane"},
				{Keys: "c", Description: "Rename pane", Mode: "pane"},
				{Keys: "h/j/k/l", Description: "Move focus left/down/up/right", Mode: "pane"},
				{Keys: "Tab", Description: "Switch focus", Mode: "pane"},
			},
		},
		{
			Name:        "Tab Mode (Ctrl+g t)",
			Description: "Create and manage tabs",
			Keymaps: []Keymap{
				{Keys: "n", Description: "New tab", Mode: "tab"},
				{Keys: "x", Description: "Close tab", Mode: "tab"},
				{Keys: "r", Description: "Rename tab", Mode: "tab"},
				{Keys: "h/l", Description: "Go to previous/next tab", Mode: "tab"},
				{Keys: "1-9", Description: "Go to tab number", Mode: "tab"},
				{Keys: "Tab", Description: "Toggle between tabs", Mode: "tab"},
				{Keys: "s", Description: "Sync tab (send input to all panes)", Mode: "tab"},
				{Keys: "b", Description: "Break pane to new tab", Mode: "tab"},
				{Keys: "[", Description: "Break pane left", Mode: "tab"},
				{Keys: "]", Description: "Break pane right", Mode: "tab"},
			},
		},
		{
			Name:        "Resize Mode (Ctrl+g r)",
			Description: "Resize panes",
			Keymaps: []Keymap{
				{Keys: "h/j/k/l", Description: "Increase size left/down/up/right", Mode: "resize"},
				{Keys: "H/J/K/L", Description: "Decrease size left/down/up/right", Mode: "resize"},
				{Keys: "+/=", Description: "Increase size", Mode: "resize"},
				{Keys: "-", Description: "Decrease size", Mode: "resize"},
			},
		},
		{
			Name:        "Move Mode (Ctrl+g m)",
			Description: "Move panes around",
			Keymaps: []Keymap{
				{Keys: "h/j/k/l", Description: "Move pane left/down/up/right", Mode: "move"},
				{Keys: "n/Tab", Description: "Move pane forward", Mode: "move"},
				{Keys: "p", Description: "Move pane backward", Mode: "move"},
			},
		},
		{
			Name:        "Scroll Mode (Ctrl+g s)",
			Description: "Scroll through output, search, edit scrollback",
			Keymaps: []Keymap{
				{Keys: "j/k", Description: "Scroll down/up", Mode: "scroll"},
				{Keys: "d/u", Description: "Half page down/up", Mode: "scroll"},
				{Keys: "Ctrl+f/b", Description: "Page down/up", Mode: "scroll"},
				{Keys: "f", Description: "Enter search mode", Mode: "scroll"},
				{Keys: "e", Description: "Edit scrollback in $EDITOR", Mode: "scroll"},
				{Keys: "Ctrl+c", Description: "Exit scroll mode", Mode: "scroll"},
			},
		},
		{
			Name:        "Search Mode",
			Description: "Search in scroll buffer",
			Keymaps: []Keymap{
				{Keys: "n/p", Description: "Next/previous match", Mode: "search"},
				{Keys: "c", Description: "Toggle case sensitivity", Mode: "search"},
				{Keys: "w", Description: "Toggle wrap", Mode: "search"},
				{Keys: "o", Description: "Toggle whole word", Mode: "search"},
			},
		},
		{
			Name:        "Session Mode (Ctrl+g o)",
			Description: "Manage sessions and plugins",
			Keymaps: []Keymap{
				{Keys: "d", Description: "Detach from session", Mode: "session"},
				{Keys: "w", Description: "Open session manager", Mode: "session"},
				{Keys: "c", Description: "Open configuration", Mode: "session"},
				{Keys: "p", Description: "Open plugin manager", Mode: "session"},
			},
		},
		{
			Name:        "Quick Actions (Always Available)",
			Description: "Work in both locked and normal modes",
			Keymaps: []Keymap{
				{Keys: "Alt+n", Description: "New pane", Mode: ""},
				{Keys: "Alt+h/j/k/l", Description: "Move focus (vim style)", Mode: ""},
				{Keys: "Alt+←/↓/↑/→", Description: "Move focus (arrows)", Mode: ""},
				{Keys: "Alt+f", Description: "Toggle floating panes", Mode: ""},
				{Keys: "Alt++/-", Description: "Increase/decrease pane size", Mode: ""},
				{Keys: "Alt+[/]", Description: "Previous/next layout", Mode: ""},
				{Keys: "Alt+i/o", Description: "Move tab left/right", Mode: ""},
				{Keys: "Ctrl+q", Description: "Quit Zellij", Mode: ""},
			},
		},
		{
			Name:        "Zellij Forgot (Alt+y)",
			Description: "Press Alt+y to open keybinding cheatsheet!",
			Keymaps: []Keymap{
				{Keys: "Alt+y", Description: "Open zellij_forgot plugin (shows all keybindings)", Mode: ""},
			},
		},
	}
}
