package tui

// GetTmuxKeymaps returns all Tmux keymaps organized by category
func GetTmuxKeymaps() []KeymapCategory {
	return []KeymapCategory{
		{
			Name:        "Prefix & Basics",
			Description: "Essential Tmux keybindings (prefix is Ctrl+a)",
			Keymaps: []Keymap{
				{Keys: "Ctrl+a", Description: "Prefix key (replaces default Ctrl+b)", Mode: ""},
				{Keys: "Ctrl+a Ctrl+a", Description: "Send Ctrl+a to terminal", Mode: ""},
				{Keys: "Ctrl+a ?", Description: "Show all keybindings", Mode: ""},
				{Keys: "Ctrl+a :", Description: "Enter command mode", Mode: ""},
				{Keys: "Ctrl+a q", Description: "Show pane numbers", Mode: ""},
			},
		},
		{
			Name:        "Panes",
			Description: "Split, navigate, and manage panes",
			Keymaps: []Keymap{
				{Keys: "Ctrl+a v", Description: "Split pane vertically (right)", Mode: ""},
				{Keys: "Ctrl+a d", Description: "Split pane horizontally (down)", Mode: ""},
				{Keys: "Ctrl+a x", Description: "Close current pane", Mode: ""},
				{Keys: "Ctrl+h/j/k/l", Description: "Navigate panes (vim-tmux-navigator)", Mode: ""},
				{Keys: "Ctrl+a z", Description: "Toggle pane zoom (fullscreen)", Mode: ""},
				{Keys: "Ctrl+a {", Description: "Move pane left", Mode: ""},
				{Keys: "Ctrl+a }", Description: "Move pane right", Mode: ""},
				{Keys: "Ctrl+a Space", Description: "Cycle through pane layouts", Mode: ""},
				{Keys: "Ctrl+a !", Description: "Convert pane to window", Mode: ""},
			},
		},
		{
			Name:        "Windows (Tabs)",
			Description: "Create and manage windows",
			Keymaps: []Keymap{
				{Keys: "Ctrl+a c", Description: "Create new window", Mode: ""},
				{Keys: "Ctrl+a &", Description: "Close current window", Mode: ""},
				{Keys: "Ctrl+a n", Description: "Next window", Mode: ""},
				{Keys: "Ctrl+a p", Description: "Previous window", Mode: ""},
				{Keys: "Ctrl+a 0-9", Description: "Switch to window number", Mode: ""},
				{Keys: "Ctrl+a w", Description: "List all windows", Mode: ""},
				{Keys: "Ctrl+a ,", Description: "Rename current window", Mode: ""},
				{Keys: "Ctrl+a l", Description: "Toggle last active window", Mode: ""},
			},
		},
		{
			Name:        "Sessions",
			Description: "Manage Tmux sessions",
			Keymaps: []Keymap{
				{Keys: "Ctrl+a d", Description: "Detach from session", Mode: ""},
				{Keys: "Ctrl+a s", Description: "List sessions", Mode: ""},
				{Keys: "Ctrl+a $", Description: "Rename session", Mode: ""},
				{Keys: "Ctrl+a (", Description: "Switch to previous session", Mode: ""},
				{Keys: "Ctrl+a )", Description: "Switch to next session", Mode: ""},
				{Keys: "Ctrl+a K", Description: "Kill all other sessions", Mode: ""},
			},
		},
		{
			Name:        "Copy Mode (Vi)",
			Description: "Navigate and copy text (vi mode enabled)",
			Keymaps: []Keymap{
				{Keys: "Ctrl+a [", Description: "Enter copy mode", Mode: ""},
				{Keys: "v", Description: "Start selection (in copy mode)", Mode: "copy"},
				{Keys: "y", Description: "Copy selection (uses system clipboard)", Mode: "copy"},
				{Keys: "q", Description: "Exit copy mode", Mode: "copy"},
				{Keys: "Ctrl+u/d", Description: "Page up/down (in copy mode)", Mode: "copy"},
				{Keys: "/", Description: "Search forward (in copy mode)", Mode: "copy"},
				{Keys: "?", Description: "Search backward (in copy mode)", Mode: "copy"},
				{Keys: "n/N", Description: "Next/previous search result", Mode: "copy"},
			},
		},
		{
			Name:        "Floating & Special",
			Description: "Gentleman.Dots custom keybindings",
			Keymaps: []Keymap{
				{Keys: "Alt+g", Description: "Toggle floating scratch terminal", Mode: ""},
			},
		},
		{
			Name:        "Plugins (TPM)",
			Description: "Plugin management keybindings",
			Keymaps: []Keymap{
				{Keys: "Ctrl+a I", Description: "Install plugins", Mode: ""},
				{Keys: "Ctrl+a U", Description: "Update plugins", Mode: ""},
				{Keys: "Ctrl+a Alt+u", Description: "Uninstall unused plugins", Mode: ""},
			},
		},
		{
			Name:        "Resurrect (Sessions)",
			Description: "Save and restore sessions",
			Keymaps: []Keymap{
				{Keys: "Ctrl+a Ctrl+s", Description: "Save session", Mode: ""},
				{Keys: "Ctrl+a Ctrl+r", Description: "Restore session", Mode: ""},
			},
		},
	}
}
