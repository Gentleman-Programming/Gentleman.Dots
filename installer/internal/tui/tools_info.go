package tui

// ToolInfo contains information about a tool
type ToolInfo struct {
	Name        string
	Description string
	Pros        []string
	Cons        []string
	Website     string
}

// KeymapCategory groups related keymaps
type KeymapCategory struct {
	Name        string
	Description string
	Keymaps     []Keymap
}

// Keymap represents a single keybinding
type Keymap struct {
	Keys        string
	Description string
	Mode        string // "n" normal, "v" visual, "i" insert
}

// GetTerminalInfo returns info about terminal emulators
func GetTerminalInfo() map[string]ToolInfo {
	return map[string]ToolInfo{
		"alacritty": {
			Name:        "Alacritty",
			Description: "GPU-accelerated terminal emulator focused on simplicity and performance",
			Pros: []string{
				"Fastest terminal emulator (GPU rendering)",
				"Very low latency",
				"Simple TOML configuration",
				"Cross-platform (Linux, macOS, Windows)",
				"Low memory usage",
			},
			Cons: []string{
				"No tabs/splits (use tmux/zellij)",
				"No ligatures support",
				"No built-in scrollback search (needs tmux)",
				"Minimal features by design",
			},
			Website: "https://alacritty.org",
		},
		"wezterm": {
			Name:        "WezTerm",
			Description: "GPU-accelerated terminal with built-in multiplexer, configured in Lua",
			Pros: []string{
				"Built-in tabs and splits (no tmux needed)",
				"Lua configuration (very flexible)",
				"Ligatures and font fallback support",
				"Image protocol support (sixel, iTerm2)",
				"SSH multiplexing built-in",
				"Excellent documentation",
			},
			Cons: []string{
				"Higher memory usage than Alacritty",
				"Lua config can be complex",
				"Slightly higher latency",
			},
			Website: "https://wezfurlong.org/wezterm/",
		},
		"kitty": {
			Name:        "Kitty",
			Description: "Fast, feature-rich terminal with graphics protocol support",
			Pros: []string{
				"Very fast (GPU accelerated)",
				"Native image display (kitty protocol)",
				"Ligatures support",
				"Built-in tabs and layouts",
				"Kitten extensions system",
				"Great for image-heavy workflows",
			},
			Cons: []string{
				"macOS only in this installer",
				"Custom config format",
				"Some apps need kitty-specific config",
			},
			Website: "https://sw.kovidgoyal.net/kitty/",
		},
		"ghostty": {
			Name:        "Ghostty",
			Description: "Native terminal by Mitchell Hashimoto (Hashicorp founder)",
			Pros: []string{
				"Native performance (not Electron)",
				"Zero config needed to start",
				"Native macOS/Linux look and feel",
				"Built-in splits and tabs",
				"Very fast rendering",
				"Modern codebase (Zig)",
			},
			Cons: []string{
				"Relatively new project",
				"Smaller community",
				"Fewer customization options (for now)",
			},
			Website: "https://ghostty.org",
		},
	}
}

// GetShellInfo returns info about shells
func GetShellInfo() map[string]ToolInfo {
	return map[string]ToolInfo{
		"fish": {
			Name:        "Fish",
			Description: "Friendly Interactive SHell - user-friendly with great defaults",
			Pros: []string{
				"Amazing autosuggestions out of the box",
				"Syntax highlighting by default",
				"Web-based configuration UI",
				"Great error messages",
				"No configuration needed to be productive",
				"Fast and responsive",
			},
			Cons: []string{
				"Not POSIX compliant (scripts differ)",
				"Can't run bash scripts directly",
				"Smaller plugin ecosystem than zsh",
			},
			Website: "https://fishshell.com",
		},
		"zsh": {
			Name:        "Zsh",
			Description: "Z Shell - powerful and highly customizable, POSIX-like",
			Pros: []string{
				"POSIX compatible (bash scripts work)",
				"Huge plugin ecosystem (oh-my-zsh)",
				"PowerLevel10k for amazing prompts",
				"Very mature and stable",
				"Great completion system",
				"Default on macOS",
			},
			Cons: []string{
				"Slow startup if misconfigured",
				"Needs plugins for good defaults",
				"Configuration can be complex",
			},
			Website: "https://www.zsh.org",
		},
		"nushell": {
			Name:        "Nushell",
			Description: "Modern shell with structured data - thinks in tables, not text",
			Pros: []string{
				"Data-first: output is structured (tables)",
				"Built-in support for JSON, YAML, TOML",
				"Pipeline operations like filter, select, sort",
				"Cross-platform consistency",
				"Modern, clean syntax",
				"Great for data manipulation",
			},
			Cons: []string{
				"Not POSIX compatible at all",
				"Steeper learning curve",
				"Smaller ecosystem",
				"Some tools need wrappers",
			},
			Website: "https://www.nushell.sh",
		},
	}
}

// GetWMInfo returns info about window managers/multiplexers
func GetWMInfo() map[string]ToolInfo {
	return map[string]ToolInfo{
		"tmux": {
			Name:        "Tmux",
			Description: "Terminal multiplexer - sessions, windows, and panes",
			Pros: []string{
				"Industry standard, everywhere",
				"Persistent sessions (survives disconnects)",
				"Huge plugin ecosystem (TPM)",
				"Remote pairing support",
				"Scriptable and automatable",
				"Very stable and mature",
			},
			Cons: []string{
				"Steep learning curve",
				"Default keybindings are awkward",
				"Configuration syntax is dated",
				"No native mouse support (needs config)",
			},
			Website: "https://github.com/tmux/tmux",
		},
		"zellij": {
			Name:        "Zellij",
			Description: "Modern terminal workspace - batteries included",
			Pros: []string{
				"Great UI out of the box",
				"Floating panes and tabs",
				"WebAssembly plugins",
				"Built-in session manager",
				"Discoverable keybindings (shows hints)",
				"Modern and actively developed",
			},
			Cons: []string{
				"Younger project than tmux",
				"Smaller plugin ecosystem",
				"Higher memory usage",
				"Less ubiquitous on servers",
			},
			Website: "https://zellij.dev",
		},
	}
}

// GetNvimKeymaps returns all Neovim keymaps organized by category
func GetNvimKeymaps() []KeymapCategory {
	return []KeymapCategory{
		{
			Name:        "Harpoon (Quick File Navigation)",
			Description: "Mark and jump to files instantly - by ThePrimeagen",
			Keymaps: []Keymap{
				{Keys: "<leader>H", Description: "Add file to Harpoon list", Mode: "n"},
				{Keys: "<leader>h", Description: "Toggle Harpoon quick menu", Mode: "n"},
				{Keys: "<leader>1", Description: "Jump to Harpoon file 1", Mode: "n"},
				{Keys: "<leader>2", Description: "Jump to Harpoon file 2", Mode: "n"},
				{Keys: "<leader>3", Description: "Jump to Harpoon file 3", Mode: "n"},
				{Keys: "<leader>4", Description: "Jump to Harpoon file 4", Mode: "n"},
				{Keys: "<leader>5", Description: "Jump to Harpoon file 5", Mode: "n"},
			},
		},
		{
			Name:        "Mini.files (File Explorer)",
			Description: "Edit filesystem like a buffer",
			Keymaps: []Keymap{
				{Keys: "<leader>fm", Description: "Open mini.files (current file dir)", Mode: "n"},
				{Keys: "<leader>fM", Description: "Open mini.files (cwd)", Mode: "n"},
				{Keys: "g.", Description: "Toggle hidden files", Mode: "n"},
				{Keys: "gc", Description: "Set current directory as cwd", Mode: "n"},
				{Keys: "<C-w>s", Description: "Open in horizontal split", Mode: "n"},
				{Keys: "<C-w>v", Description: "Open in vertical split", Mode: "n"},
				{Keys: "<C-w>S", Description: "Open in horizontal split and close", Mode: "n"},
				{Keys: "<C-w>V", Description: "Open in vertical split and close", Mode: "n"},
			},
		},
		{
			Name:        "Oil.nvim (Directory Editor)",
			Description: "Edit your filesystem like a buffer",
			Keymaps: []Keymap{
				{Keys: "-", Description: "Open Oil (parent dir)", Mode: "n"},
				{Keys: "<leader>E", Description: "Open Oil (floating)", Mode: "n"},
				{Keys: "<leader>-", Description: "Open Oil in current file's directory", Mode: "n"},
				{Keys: "g?", Description: "Show help", Mode: "n"},
				{Keys: "<CR>", Description: "Select file/directory", Mode: "n"},
				{Keys: "<C-s>", Description: "Open in vertical split", Mode: "n"},
				{Keys: "<C-v>", Description: "Open in horizontal split", Mode: "n"},
				{Keys: "<C-t>", Description: "Open in new tab", Mode: "n"},
				{Keys: "<C-p>", Description: "Preview file", Mode: "n"},
				{Keys: "<C-c>", Description: "Close Oil", Mode: "n"},
				{Keys: "<C-r>", Description: "Refresh", Mode: "n"},
				{Keys: "_", Description: "Open cwd", Mode: "n"},
				{Keys: "`", Description: "cd to current directory", Mode: "n"},
				{Keys: "~", Description: ":tcd to current directory", Mode: "n"},
				{Keys: "gs", Description: "Change sort", Mode: "n"},
				{Keys: "gx", Description: "Open external", Mode: "n"},
				{Keys: "g.", Description: "Toggle hidden files", Mode: "n"},
				{Keys: "g\\", Description: "Toggle trash", Mode: "n"},
				{Keys: "q", Description: "Close Oil", Mode: "n"},
			},
		},
		{
			Name:        "Mini.surround (Text Manipulation)",
			Description: "Add, delete, replace surroundings (brackets, quotes, etc)",
			Keymaps: []Keymap{
				{Keys: "sa", Description: "Add surrounding (e.g., saiw\" = surround word with \")", Mode: "n,v"},
				{Keys: "sd", Description: "Delete surrounding (e.g., sd\" = delete surrounding \")", Mode: "n"},
				{Keys: "sr", Description: "Replace surrounding (e.g., sr\"' = replace \" with ')", Mode: "n"},
				{Keys: "sf", Description: "Find surrounding (to the right)", Mode: "n"},
				{Keys: "sF", Description: "Find surrounding (to the left)", Mode: "n"},
				{Keys: "sh", Description: "Highlight surrounding", Mode: "n"},
				{Keys: "sn", Description: "Update number of search lines", Mode: "n"},
			},
		},
		{
			Name:        "Snacks Picker (Fuzzy Finding)",
			Description: "Find files, grep text, search everything",
			Keymaps: []Keymap{
				{Keys: "<leader><space>", Description: "Find Files (Root Dir)", Mode: "n"},
				{Keys: "<leader>,", Description: "Switch Buffer", Mode: "n"},
				{Keys: "<leader>/", Description: "Grep (Root Dir)", Mode: "n"},
				{Keys: "<leader>:", Description: "Command History", Mode: "n"},
				{Keys: "<leader>fb", Description: "Buffers", Mode: "n"},
				{Keys: "<leader>fB", Description: "Buffers (all, including hidden)", Mode: "n"},
				{Keys: "<leader>fc", Description: "Find Config File", Mode: "n"},
				{Keys: "<leader>ff", Description: "Find Files (Root Dir)", Mode: "n"},
				{Keys: "<leader>fF", Description: "Find Files (cwd)", Mode: "n"},
				{Keys: "<leader>fg", Description: "Find Files (git-files)", Mode: "n"},
				{Keys: "<leader>fr", Description: "Recent Files", Mode: "n"},
				{Keys: "<leader>fR", Description: "Recent Files (cwd)", Mode: "n"},
				{Keys: "<leader>fp", Description: "Projects", Mode: "n"},
				{Keys: "<leader>n", Description: "Notification History", Mode: "n"},
			},
		},
		{
			Name:        "Search Commands (Snacks)",
			Description: "Search registers, marks, diagnostics, and more",
			Keymaps: []Keymap{
				{Keys: "<leader>s\"", Description: "Registers", Mode: "n"},
				{Keys: "<leader>s/", Description: "Search History", Mode: "n"},
				{Keys: "<leader>sa", Description: "Auto Commands", Mode: "n"},
				{Keys: "<leader>sb", Description: "Buffer Lines", Mode: "n"},
				{Keys: "<leader>sB", Description: "Grep Open Buffers", Mode: "n"},
				{Keys: "<leader>sc", Description: "Command History", Mode: "n"},
				{Keys: "<leader>sC", Description: "Commands", Mode: "n"},
				{Keys: "<leader>sd", Description: "Document Diagnostics", Mode: "n"},
				{Keys: "<leader>sD", Description: "Workspace Diagnostics", Mode: "n"},
				{Keys: "<leader>sg", Description: "Grep (Root Dir)", Mode: "n"},
				{Keys: "<leader>sG", Description: "Grep (cwd)", Mode: "n"},
				{Keys: "<leader>sh", Description: "Help Pages", Mode: "n"},
				{Keys: "<leader>sH", Description: "Highlights", Mode: "n"},
				{Keys: "<leader>si", Description: "Icons", Mode: "n"},
				{Keys: "<leader>sj", Description: "Jumps", Mode: "n"},
				{Keys: "<leader>sk", Description: "Keymaps", Mode: "n"},
				{Keys: "<leader>sl", Description: "Location List", Mode: "n"},
				{Keys: "<leader>sm", Description: "Marks", Mode: "n"},
				{Keys: "<leader>sM", Description: "Man Pages", Mode: "n"},
				{Keys: "<leader>sp", Description: "Search Plugin Spec", Mode: "n"},
				{Keys: "<leader>sq", Description: "Quickfix List", Mode: "n"},
				{Keys: "<leader>sR", Description: "Resume Last Search", Mode: "n"},
				{Keys: "<leader>ss", Description: "LSP Symbols (Document)", Mode: "n"},
				{Keys: "<leader>sS", Description: "LSP Symbols (Workspace)", Mode: "n"},
				{Keys: "<leader>su", Description: "Undotree", Mode: "n"},
				{Keys: "<leader>sw", Description: "Word under cursor (Root Dir)", Mode: "n,v"},
				{Keys: "<leader>sW", Description: "Word under cursor (cwd)", Mode: "n,v"},
			},
		},
		{
			Name:        "LSP & Code Navigation",
			Description: "Go-to definitions, references, and previews",
			Keymaps: []Keymap{
				{Keys: "gd", Description: "Go to Definition", Mode: "n"},
				{Keys: "gr", Description: "Go to References", Mode: "n"},
				{Keys: "gI", Description: "Go to Implementation", Mode: "n"},
				{Keys: "gy", Description: "Go to Type Definition", Mode: "n"},
				{Keys: "gD", Description: "Go to Declaration", Mode: "n"},
				{Keys: "K", Description: "Hover Documentation", Mode: "n"},
				{Keys: "gK", Description: "Signature Help", Mode: "n"},
				{Keys: "<leader>ca", Description: "Code Actions", Mode: "n,v"},
				{Keys: "<leader>cc", Description: "Run Codelens", Mode: "n"},
				{Keys: "<leader>cC", Description: "Refresh Codelens", Mode: "n"},
				{Keys: "<leader>cr", Description: "Rename Symbol", Mode: "n"},
				{Keys: "<leader>cR", Description: "Rename File", Mode: "n"},
				{Keys: "<leader>cs", Description: "Symbols Outline", Mode: "n"},
				{Keys: "gpd", Description: "Preview Definition (popup)", Mode: "n"},
				{Keys: "gpD", Description: "Preview Declaration (popup)", Mode: "n"},
				{Keys: "gpi", Description: "Preview Implementation (popup)", Mode: "n"},
				{Keys: "gpy", Description: "Preview Type Definition (popup)", Mode: "n"},
				{Keys: "gpr", Description: "Preview References (popup)", Mode: "n"},
				{Keys: "gP", Description: "Close all preview windows", Mode: "n"},
			},
		},
		{
			Name:        "Git Integration",
			Description: "Git operations from within Neovim",
			Keymaps: []Keymap{
				{Keys: "<leader>gb", Description: "Git Blame (show who changed line)", Mode: "n"},
				{Keys: "<leader>go", Description: "Open file/folder in git repo", Mode: "n"},
				{Keys: "<leader>gB", Description: "Git Browse", Mode: "n"},
				{Keys: "<leader>gc", Description: "Git Commits", Mode: "n"},
				{Keys: "<leader>gd", Description: "Git Diff (hunks)", Mode: "n"},
				{Keys: "<leader>gf", Description: "Git File History (current)", Mode: "n"},
				{Keys: "<leader>gg", Description: "Open Lazygit", Mode: "n"},
				{Keys: "<leader>gl", Description: "Git Log", Mode: "n"},
				{Keys: "<leader>gL", Description: "Git Log (line)", Mode: "n"},
				{Keys: "<leader>gs", Description: "Git Status", Mode: "n"},
				{Keys: "<leader>gS", Description: "Git Stash", Mode: "n"},
				{Keys: "]h", Description: "Next Hunk", Mode: "n"},
				{Keys: "[h", Description: "Previous Hunk", Mode: "n"},
				{Keys: "<leader>ghp", Description: "Preview Hunk", Mode: "n"},
				{Keys: "<leader>ghs", Description: "Stage Hunk", Mode: "n"},
				{Keys: "<leader>ghr", Description: "Reset Hunk", Mode: "n"},
				{Keys: "<leader>ghS", Description: "Stage Buffer", Mode: "n"},
				{Keys: "<leader>ghR", Description: "Reset Buffer", Mode: "n"},
			},
		},
		{
			Name:        "AI Assistants (Copilot + OpenCode)",
			Description: "AI-powered coding assistance",
			Keymaps: []Keymap{
				{Keys: "<Tab>", Description: "Accept Copilot suggestion", Mode: "i"},
				{Keys: "<C-]>", Description: "Dismiss Copilot suggestion", Mode: "i"},
				{Keys: "<M-]>", Description: "Next Copilot suggestion", Mode: "i"},
				{Keys: "<M-[>", Description: "Previous Copilot suggestion", Mode: "i"},
				{Keys: "<leader>aa", Description: "Toggle OpenCode", Mode: "n"},
				{Keys: "<leader>as", Description: "OpenCode select", Mode: "n,x"},
				{Keys: "<leader>ai", Description: "OpenCode ask", Mode: "n,x"},
				{Keys: "<leader>aI", Description: "OpenCode ask with context", Mode: "n,x"},
				{Keys: "<leader>ab", Description: "OpenCode ask about buffer", Mode: "n,x"},
				{Keys: "<leader>ap", Description: "OpenCode prompt", Mode: "n,x"},
				{Keys: "<leader>ape", Description: "OpenCode explain", Mode: "n,x"},
				{Keys: "<leader>apf", Description: "OpenCode fix", Mode: "n,x"},
				{Keys: "<leader>apd", Description: "OpenCode diagnose", Mode: "n,x"},
				{Keys: "<leader>apr", Description: "OpenCode review", Mode: "n,x"},
				{Keys: "<leader>apt", Description: "OpenCode test", Mode: "n,x"},
				{Keys: "<leader>apo", Description: "OpenCode optimize", Mode: "n,x"},
			},
		},
		{
			Name:        "Tmux Navigation",
			Description: "Seamless navigation between Neovim and Tmux panes",
			Keymaps: []Keymap{
				{Keys: "<C-h>", Description: "Navigate to the left pane", Mode: "n"},
				{Keys: "<C-j>", Description: "Navigate to the bottom pane", Mode: "n"},
				{Keys: "<C-k>", Description: "Navigate to the top pane", Mode: "n"},
				{Keys: "<C-l>", Description: "Navigate to the right pane", Mode: "n"},
				{Keys: "<C-\\>", Description: "Navigate to last active pane", Mode: "n"},
				{Keys: "<C-Space>", Description: "Navigate to next pane", Mode: "n"},
			},
		},
		{
			Name:        "Obsidian (Notes)",
			Description: "Note-taking and knowledge management",
			Keymaps: []Keymap{
				{Keys: "<leader>oc", Description: "Obsidian Check Checkbox", Mode: "n"},
				{Keys: "<leader>ot", Description: "Insert Obsidian Template", Mode: "n"},
				{Keys: "<leader>oo", Description: "Open in Obsidian App", Mode: "n"},
				{Keys: "<leader>ob", Description: "Show Obsidian Backlinks", Mode: "n"},
				{Keys: "<leader>ol", Description: "Show Obsidian Links", Mode: "n"},
				{Keys: "<leader>on", Description: "Create New Note", Mode: "n"},
				{Keys: "<leader>os", Description: "Search Obsidian", Mode: "n"},
				{Keys: "<leader>oq", Description: "Quick Switch", Mode: "n"},
				{Keys: "gf", Description: "Follow link (in note buffer)", Mode: "n"},
				{Keys: "<leader>ch", Description: "Toggle checkbox (in note buffer)", Mode: "n"},
				{Keys: "<CR>", Description: "Smart action (in note buffer)", Mode: "n"},
			},
		},
		{
			Name:        "Buffers & Tabs",
			Description: "Buffer and tab management",
			Keymaps: []Keymap{
				{Keys: "<leader>bb", Description: "Switch to Other Buffer", Mode: "n"},
				{Keys: "<leader>bd", Description: "Delete Buffer", Mode: "n"},
				{Keys: "<leader>bD", Description: "Delete Buffer and Window", Mode: "n"},
				{Keys: "<leader>bo", Description: "Delete Other Buffers", Mode: "n"},
				{Keys: "<leader>bp", Description: "Toggle Buffer Pin", Mode: "n"},
				{Keys: "<leader>bP", Description: "Delete Non-Pinned Buffers", Mode: "n"},
				{Keys: "[b", Description: "Previous Buffer", Mode: "n"},
				{Keys: "]b", Description: "Next Buffer", Mode: "n"},
				{Keys: "<leader><tab>l", Description: "Last Tab", Mode: "n"},
				{Keys: "<leader><tab>f", Description: "First Tab", Mode: "n"},
				{Keys: "<leader><tab><tab>", Description: "New Tab", Mode: "n"},
				{Keys: "<leader><tab>d", Description: "Close Tab", Mode: "n"},
				{Keys: "<leader><tab>]", Description: "Next Tab", Mode: "n"},
				{Keys: "<leader><tab>[", Description: "Previous Tab", Mode: "n"},
			},
		},
		{
			Name:        "Windows & Splits",
			Description: "Window navigation and management",
			Keymaps: []Keymap{
				{Keys: "<C-h>", Description: "Go to Left Window", Mode: "n"},
				{Keys: "<C-j>", Description: "Go to Lower Window", Mode: "n"},
				{Keys: "<C-k>", Description: "Go to Upper Window", Mode: "n"},
				{Keys: "<C-l>", Description: "Go to Right Window", Mode: "n"},
				{Keys: "<leader>w", Description: "Windows menu (which-key)", Mode: "n"},
				{Keys: "<leader>wd", Description: "Delete Window", Mode: "n"},
				{Keys: "<leader>wm", Description: "Maximize Window", Mode: "n"},
				{Keys: "<leader>-", Description: "Split Below", Mode: "n"},
				{Keys: "<leader>|", Description: "Split Right", Mode: "n"},
				{Keys: "<C-Up>", Description: "Increase Height", Mode: "n"},
				{Keys: "<C-Down>", Description: "Decrease Height", Mode: "n"},
				{Keys: "<C-Left>", Description: "Decrease Width", Mode: "n"},
				{Keys: "<C-Right>", Description: "Increase Width", Mode: "n"},
			},
		},
		{
			Name:        "Debugging (DAP)",
			Description: "Debug Adapter Protocol integration",
			Keymaps: []Keymap{
				{Keys: "<leader>dB", Description: "Breakpoint Condition", Mode: "n"},
				{Keys: "<leader>db", Description: "Toggle Breakpoint", Mode: "n"},
				{Keys: "<leader>dc", Description: "Continue", Mode: "n"},
				{Keys: "<leader>da", Description: "Run with Args", Mode: "n"},
				{Keys: "<leader>dC", Description: "Run to Cursor", Mode: "n"},
				{Keys: "<leader>dg", Description: "Go to Line (No Execute)", Mode: "n"},
				{Keys: "<leader>di", Description: "Step Into", Mode: "n"},
				{Keys: "<leader>dj", Description: "Down", Mode: "n"},
				{Keys: "<leader>dk", Description: "Up", Mode: "n"},
				{Keys: "<leader>dl", Description: "Run Last", Mode: "n"},
				{Keys: "<leader>do", Description: "Step Out", Mode: "n"},
				{Keys: "<leader>dO", Description: "Step Over", Mode: "n"},
				{Keys: "<leader>dp", Description: "Pause", Mode: "n"},
				{Keys: "<leader>dr", Description: "Toggle REPL", Mode: "n"},
				{Keys: "<leader>ds", Description: "Session", Mode: "n"},
				{Keys: "<leader>dt", Description: "Terminate", Mode: "n"},
				{Keys: "<leader>dw", Description: "Widgets", Mode: "n"},
			},
		},
		{
			Name:        "UI Toggles",
			Description: "Toggle UI elements and settings",
			Keymaps: []Keymap{
				{Keys: "<leader>ub", Description: "Toggle Background (dark/light)", Mode: "n"},
				{Keys: "<leader>uC", Description: "Colorscheme Picker", Mode: "n"},
				{Keys: "<leader>ud", Description: "Toggle Diagnostics", Mode: "n"},
				{Keys: "<leader>uf", Description: "Toggle Auto Format (global)", Mode: "n"},
				{Keys: "<leader>uF", Description: "Toggle Auto Format (buffer)", Mode: "n"},
				{Keys: "<leader>ug", Description: "Toggle Indent Guides", Mode: "n"},
				{Keys: "<leader>uh", Description: "Toggle Inlay Hints", Mode: "n"},
				{Keys: "<leader>uI", Description: "Inspect Pos (Treesitter)", Mode: "n"},
				{Keys: "<leader>uk", Description: "Toggle Screenkey", Mode: "n"},
				{Keys: "<leader>ul", Description: "Toggle Line Numbers", Mode: "n"},
				{Keys: "<leader>uL", Description: "Toggle Relative Numbers", Mode: "n"},
				{Keys: "<leader>un", Description: "Dismiss Notifications", Mode: "n"},
				{Keys: "<leader>us", Description: "Toggle Spelling", Mode: "n"},
				{Keys: "<leader>uT", Description: "Toggle Treesitter Highlight", Mode: "n"},
				{Keys: "<leader>uw", Description: "Toggle Word Wrap", Mode: "n"},
				{Keys: "<leader>uz", Description: "Toggle Zen Mode", Mode: "n"},
				{Keys: "<leader>uZ", Description: "Toggle Zoom", Mode: "n"},
				{Keys: "<leader>z", Description: "Zen Mode", Mode: "n"},
			},
		},
		{
			Name:        "Diagnostics & Quickfix",
			Description: "Navigate errors, warnings, and quickfix",
			Keymaps: []Keymap{
				{Keys: "<leader>xx", Description: "Document Diagnostics (Trouble)", Mode: "n"},
				{Keys: "<leader>xX", Description: "Workspace Diagnostics (Trouble)", Mode: "n"},
				{Keys: "<leader>xL", Description: "Location List (Trouble)", Mode: "n"},
				{Keys: "<leader>xQ", Description: "Quickfix List (Trouble)", Mode: "n"},
				{Keys: "[d", Description: "Previous Diagnostic", Mode: "n"},
				{Keys: "]d", Description: "Next Diagnostic", Mode: "n"},
				{Keys: "[e", Description: "Previous Error", Mode: "n"},
				{Keys: "]e", Description: "Next Error", Mode: "n"},
				{Keys: "[w", Description: "Previous Warning", Mode: "n"},
				{Keys: "]w", Description: "Next Warning", Mode: "n"},
				{Keys: "[q", Description: "Previous Quickfix", Mode: "n"},
				{Keys: "]q", Description: "Next Quickfix", Mode: "n"},
			},
		},
		{
			Name:        "Session & Quit",
			Description: "Session management and exiting",
			Keymaps: []Keymap{
				{Keys: "<leader>qq", Description: "Quit All", Mode: "n"},
				{Keys: "<leader>qs", Description: "Restore Session", Mode: "n"},
				{Keys: "<leader>qS", Description: "Select Session", Mode: "n"},
				{Keys: "<leader>ql", Description: "Restore Last Session", Mode: "n"},
				{Keys: "<leader>qd", Description: "Don't Save Current Session", Mode: "n"},
			},
		},
		{
			Name:        "Custom Keymaps",
			Description: "Gentleman.Dots custom keybindings",
			Keymaps: []Keymap{
				{Keys: "<C-b>", Description: "Delete to end of word (insert)", Mode: "i"},
				{Keys: "<C-c>", Description: "Escape from any mode", Mode: "i,n,v"},
				{Keys: "<leader>bq", Description: "Delete other buffers but current", Mode: "n"},
				{Keys: "<C-s>", Description: "Save file", Mode: "n"},
				{Keys: "<leader>sg", Description: "Grep Selected Text", Mode: "v"},
				{Keys: "<leader>sG", Description: "Grep Selected Text (Root Dir)", Mode: "v"},
				{Keys: "<leader>md", Description: "Delete all marks", Mode: "n"},
				{Keys: "<leader>fs", Description: "Rip Substitute (find & replace)", Mode: "n,x"},
				{Keys: "ff", Description: "FFF Find files (fast fuzzy)", Mode: "n"},
			},
		},
		{
			Name:        "General & Which-Key",
			Description: "Essential keybindings every Neovim user needs",
			Keymaps: []Keymap{
				{Keys: "<leader>?", Description: "Show Buffer Local Keymaps", Mode: "n"},
				{Keys: "<leader>e", Description: "Toggle Explorer (neo-tree)", Mode: "n"},
				{Keys: "<leader>E", Description: "Explorer (neo-tree cwd)", Mode: "n"},
				{Keys: "<leader>l", Description: "Lazy (plugin manager)", Mode: "n"},
				{Keys: "<leader>L", Description: "LazyVim Changelog", Mode: "n"},
				{Keys: "<Esc>", Description: "Escape and Clear hlsearch", Mode: "n,i"},
				{Keys: "<leader>ur", Description: "Redraw / Clear hlsearch", Mode: "n"},
				{Keys: "gx", Description: "Open with system app", Mode: "n"},
				{Keys: "j", Description: "Down (respects wrapped lines)", Mode: "n,x"},
				{Keys: "k", Description: "Up (respects wrapped lines)", Mode: "n,x"},
			},
		},
	}
}

// GetNvimInfo returns info about Neovim and its features
func GetNvimInfo() ToolInfo {
	return ToolInfo{
		Name:        "Neovim + LazyVim + Gentleman Config",
		Description: "Hyperextensible Vim-based text editor with LazyVim distribution and Gentleman customizations",
		Pros: []string{
			"Blazing fast startup (~50ms)",
			"LazyVim: Pre-configured, sane defaults",
			"LSP support for 100+ languages",
			"TreeSitter for better syntax highlighting",
			"AI integration (Copilot, Avante)",
			"Git integration (Lazygit, Gitsigns)",
			"Fuzzy finding (Snacks picker)",
			"File navigation (Oil, Mini.files, Harpoon)",
			"Note-taking with Obsidian.nvim",
			"Fully keyboard-driven workflow",
		},
		Cons: []string{
			"Steep learning curve for Vim motions",
			"Requires terminal with good font support",
			"Plugin ecosystem can be overwhelming",
			"Configuration complexity (mitigated by LazyVim)",
		},
		Website: "https://lazyvim.org",
	}
}

// LazyVimTopic represents a learning topic about LazyVim
type LazyVimTopic struct {
	Title       string
	Description string
	Content     []string
	CodeExample string
	Tips        []string
}

// GetLazyVimTopics returns all LazyVim learning topics
func GetLazyVimTopics() []LazyVimTopic {
	return []LazyVimTopic{
		{
			Title:       "What is LazyVim?",
			Description: "LazyVim is a Neovim setup powered by lazy.nvim",
			Content: []string{
				"LazyVim is NOT a plugin, it's a complete Neovim distribution.",
				"It provides sane defaults, pre-configured plugins, and a",
				"modular architecture that makes customization easy.",
				"",
				"Key concepts:",
				"• lazy.nvim: The plugin manager (handles loading/updates)",
				"• LazyVim: The distribution (pre-configured setup)",
				"• Extras: Optional modules you can enable/disable",
				"• Your config: Overrides in ~/.config/nvim/lua/plugins/",
			},
			CodeExample: `-- Your config structure:
~/.config/nvim/
├── lua/
│   ├── config/
│   │   ├── lazy.lua     -- Plugin manager setup
│   │   ├── keymaps.lua  -- Your custom keymaps
│   │   ├── options.lua  -- Neovim options
│   │   └── autocmds.lua -- Auto commands
│   └── plugins/
│       ├── example.lua  -- Your custom plugins
│       └── overrides.lua -- Override LazyVim defaults
└── init.lua`,
			Tips: []string{
				"Press <leader>l to open Lazy (plugin manager UI)",
				"Press <leader>L to see LazyVim changelog",
				"All your customizations go in lua/plugins/",
			},
		},
		{
			Title:       "LazyExtras: Enable Features",
			Description: "The EASIEST way to add languages, tools & features",
			Content: []string{
				"LazyExtras is your ONE-STOP SHOP for adding functionality!",
				"",
				"★ THE EASY WAY (recommended):",
				"  1. Open Neovim",
				"  2. Run :LazyExtras",
				"  3. Browse categories, press 'x' to toggle any extra",
				"  4. Restart Neovim - DONE!",
				"",
				"Each extra includes EVERYTHING you need:",
				"• LSP server (auto-completion, go-to-definition)",
				"• TreeSitter parser (syntax highlighting)",
				"• Formatter & Linter integration",
				"• Debugging support (DAP)",
				"• Testing support (where applicable)",
				"",
				"Categories available in :LazyExtras:",
				"• lang.*       - Languages (typescript, go, rust, python...)",
				"• editor.*     - Editor features (harpoon, mini-files...)",
				"• coding.*     - Coding helpers (copilot, snippets...)",
				"• formatting.* - Formatters (prettier, biome...)",
				"• linting.*    - Linters (eslint...)",
				"• ai.*         - AI assistants (copilot, copilot-chat...)",
				"• test.*       - Testing frameworks (core, coverage...)",
				"• dap.*        - Debugging adapters",
			},
			CodeExample: `-- :LazyExtras UI is the easiest way!
-- But if you prefer code, add to lua/config/lazy.lua:

spec = {
  { "LazyVim/LazyVim", import = "lazyvim.plugins" },
  
  -- These get added automatically when you use :LazyExtras!
  { import = "lazyvim.plugins.extras.lang.typescript" },
  { import = "lazyvim.plugins.extras.lang.go" },
  { import = "lazyvim.plugins.extras.test.core" },
  { import = "lazyvim.plugins.extras.ai.copilot" },
  
  { import = "plugins" },
}`,
			Tips: []string{
				":LazyExtras - Interactive UI to enable/disable extras",
				"Press 'x' on any extra to toggle it",
				"Changes apply after restarting Neovim",
				"No manual config needed - extras handle everything!",
			},
		},
		{
			Title:       "Adding a New Language",
			Description: "Use :LazyExtras - it's that simple!",
			Content: []string{
				"★ THE EASY WAY (99% of cases):",
				"  1. Open Neovim",
				"  2. Run :LazyExtras",
				"  3. Search for 'lang.' to see all languages",
				"  4. Press 'x' on the language you want",
				"  5. Restart Neovim - DONE!",
				"",
				"Each lang.* extra gives you EVERYTHING:",
				"• LSP server (auto-completion, diagnostics)",
				"• TreeSitter (syntax highlighting, text objects)",
				"• Formatter (auto-format on save)",
				"• Linter (code quality checks)",
				"• DAP debugger (where available)",
				"• Test runner integration (where available)",
				"",
				"Popular lang extras:",
				"• lang.typescript - TS/JS with tsserver",
				"• lang.python - Python with pyright/ruff",
				"• lang.go - Go with gopls",
				"• lang.rust - Rust with rust-analyzer",
				"• lang.java - Java with jdtls",
				"• lang.docker - Dockerfile support",
				"• lang.yaml - YAML with schemas",
				"• lang.json - JSON with schemas",
			},
			CodeExample: `-- RECOMMENDED: Just use :LazyExtras UI!
-- It adds the import automatically to your config.

-- ADVANCED: Manual config (only if no extra exists)
-- Create lua/plugins/my-language.lua:
return {
  -- TreeSitter parser
  {
    "nvim-treesitter/nvim-treesitter",
    opts = { ensure_installed = { "my_lang" } },
  },
  -- LSP server  
  {
    "neovim/nvim-lspconfig",
    opts = {
      servers = {
        my_lsp = {},
      },
    },
  },
}`,
			Tips: []string{
				":LazyExtras → search 'lang' → press 'x' → restart",
				"That's it! No config files to edit!",
				":LspInfo to verify LSP is working",
				":Mason to see/manage installed servers",
			},
		},
		{
			Title:       "Installing Custom Plugins",
			Description: "How to add plugins not included in LazyVim",
			Content: []string{
				"Adding a plugin is simple - create a file in lua/plugins/",
				"Each file should return a table (or list of tables).",
				"",
				"Plugin spec options:",
				"• First element: 'owner/repo' (GitHub)",
				"• dependencies: Other plugins it needs",
				"• event: When to load (VeryLazy, BufRead, etc.)",
				"• cmd: Commands that trigger loading",
				"• keys: Keymaps (also trigger loading)",
				"• opts: Plugin options (passed to setup())",
				"• config: Custom configuration function",
			},
			CodeExample: `-- lua/plugins/my-plugins.lua
return {
  -- Simple plugin
  { "tpope/vim-sleuth" }, -- Auto-detect indent
  
  -- Plugin with options
  {
    "folke/todo-comments.nvim",
    opts = {
      signs = true,
      keywords = {
        TODO = { icon = " ", color = "info" },
        HACK = { icon = " ", color = "warning" },
      },
    },
  },
  
  -- Plugin with lazy loading
  {
    "ThePrimeagen/vim-be-good",
    cmd = "VimBeGood", -- Only load when running :VimBeGood
  },
  
  -- Plugin with keymaps
  {
    "folke/zen-mode.nvim",
    keys = {
      { "<leader>z", "<cmd>ZenMode<cr>", desc = "Zen Mode" },
    },
    opts = {
      window = { width = 90 },
    },
  },
}`,
			Tips: []string{
				"Use :Lazy to manage plugins (update, clean, profile)",
				"Press <leader>l then 'p' to profile startup time",
				"Lazy load plugins to keep startup fast!",
			},
		},
		{
			Title:       "Overriding LazyVim Defaults",
			Description: "How to customize built-in plugins",
			Content: []string{
				"You can override any LazyVim plugin configuration.",
				"Just use the same plugin name - lazy.nvim merges specs.",
				"",
				"Common overrides:",
				"• Change plugin options (opts)",
				"• Add/change keymaps (keys)",
				"• Disable a plugin entirely (enabled = false)",
				"• Change when it loads (event, cmd)",
			},
			CodeExample: `-- lua/plugins/overrides.lua
return {
  -- Change colorscheme
  {
    "LazyVim/LazyVim",
    opts = {
      colorscheme = "catppuccin",
    },
  },
  
  -- Modify telescope options
  {
    "nvim-telescope/telescope.nvim",
    opts = {
      defaults = {
        layout_strategy = "vertical",
      },
    },
  },
  
  -- Disable a plugin completely
  { "folke/flash.nvim", enabled = false },
  
  -- Add keys to existing plugin
  {
    "folke/trouble.nvim",
    keys = {
      { "<leader>tt", "<cmd>Trouble<cr>", desc = "Trouble" },
    },
  },
  
  -- Use a function for complex opts
  {
    "hrsh7th/nvim-cmp",
    opts = function(_, opts)
      local cmp = require("cmp")
      opts.mapping["<C-y>"] = cmp.mapping.confirm({ select = true })
      return opts
    end,
  },
}`,
			Tips: []string{
				"opts can be a table OR a function(_, opts)",
				"Function lets you modify existing opts",
				"Check LazyVim source for default configurations",
			},
		},
		{
			Title:       "Custom Keymaps",
			Description: "How to add your own keybindings",
			Content: []string{
				"There are two ways to add keymaps:",
				"1. In lua/config/keymaps.lua (always loaded)",
				"2. In plugin specs using 'keys' (lazy loaded)",
				"",
				"LazyVim uses <leader> = <Space> by default.",
				"",
				"Keymap groups (which-key):",
				"• <leader>f = file/find",
				"• <leader>s = search",
				"• <leader>g = git",
				"• <leader>c = code",
				"• <leader>b = buffer",
				"• <leader>w = window",
				"• <leader>u = ui toggles",
				"• <leader>x = diagnostics",
			},
			CodeExample: "-- lua/config/keymaps.lua\n" +
				"local map = vim.keymap.set\n\n" +
				"-- Basic mappings\n" +
				"map(\"n\", \"<leader>w\", \"<cmd>w<cr>\", { desc = \"Save file\" })\n" +
				"map(\"n\", \"<leader>q\", \"<cmd>q<cr>\", { desc = \"Quit\" })\n\n" +
				"-- Better navigation\n" +
				"map(\"n\", \"J\", \"mzJ`z\")  -- Join lines, keep cursor\n" +
				"map(\"n\", \"<C-d>\", \"<C-d>zz\") -- Center after scroll\n" +
				"map(\"n\", \"<C-u>\", \"<C-u>zz\")\n" +
				"map(\"n\", \"n\", \"nzzzv\")  -- Center after search\n\n" +
				"-- Move lines in visual mode\n" +
				"map(\"v\", \"J\", \":m '>+1<CR>gv=gv\")\n" +
				"map(\"v\", \"K\", \":m '<-2<CR>gv=gv\")\n\n" +
				"-- Quick escape\n" +
				"map(\"i\", \"jk\", \"<Esc>\")\n\n" +
				"-- Delete to void (don't override register)\n" +
				"map({\"n\", \"v\"}, \"<leader>d\", '\"_d')",
			Tips: []string{
				"Always add 'desc' for which-key integration",
				"Use :map to see all current mappings",
				"<leader>sk to search keymaps with telescope",
			},
		},
		{
			Title:       "LSP Configuration",
			Description: "Understanding and customizing LSP servers",
			Content: []string{
				"LSP (Language Server Protocol) provides:",
				"• Auto-completion",
				"• Go to definition/references",
				"• Hover documentation",
				"• Rename symbol",
				"• Code actions",
				"• Diagnostics (errors/warnings)",
				"",
				"Mason manages LSP installation automatically.",
				"LazyVim configures common LSPs out of the box.",
			},
			CodeExample: `-- lua/plugins/lsp.lua
return {
  {
    "neovim/nvim-lspconfig",
    opts = {
      -- Add or override LSP servers
      servers = {
        -- TypeScript
        tsserver = {
          settings = {
            typescript = {
              inlayHints = {
                includeInlayParameterNameHints = "all",
              },
            },
          },
        },
        -- Lua (for Neovim config)
        lua_ls = {
          settings = {
            Lua = {
              workspace = { checkThirdParty = false },
              completion = { callSnippet = "Replace" },
            },
          },
        },
        -- Disable a server
        jsonls = { enabled = false },
      },
    },
  },
}

-- Key LSP commands:
-- :LspInfo     - See attached servers
-- :LspLog      - View LSP logs
-- :LspRestart  - Restart LSP servers
-- :Mason       - Manage LSP installations`,
			Tips: []string{
				"Use :LspInfo to debug LSP issues",
				"Check :Mason for available servers",
				"Most lang extras configure LSP for you",
			},
		},
		{
			Title:       "Useful Commands",
			Description: "Essential LazyVim commands to know",
			Content: []string{
				"★ MOST IMPORTANT COMMAND:",
				"• :LazyExtras - Add languages, tools, features!",
				"  (This is your go-to for adding anything)",
				"",
				"Plugin Management:",
				"• :Lazy - Plugin manager UI (update, clean, profile)",
				"• :LazyHealth - Check plugin health",
				"",
				"LSP & Mason:",
				"• :Mason - LSP/formatter/linter installer UI",
				"• :LspInfo - See attached LSP servers",
				"• :LspRestart - Restart LSP servers",
				"",
				"Formatting & Linting:",
				"• :ConformInfo - Formatter status",
				"• :Format - Format current buffer",
				"",
				"Treesitter:",
				"• :TSInstall <lang> - Install parser",
				"• :TSUpdate - Update all parsers",
				"• :InspectTree - View syntax tree",
			},
			CodeExample: `-- Quick reference cheatsheet

-- File navigation
<leader><space>  Find files
<leader>,        Switch buffer
<leader>ff       Find files
<leader>fr       Recent files
<leader>fg       Git files

-- Search
<leader>/        Grep in project
<leader>sg       Live grep
<leader>sw       Search word under cursor
<leader>ss       LSP symbols

-- Code
gd               Go to definition
gr               Go to references
K                Hover docs
<leader>ca       Code actions
<leader>cr       Rename

-- Git
<leader>gg       Lazygit
<leader>gs       Git status
]h / [h          Next/prev hunk

-- UI
<leader>e        File explorer
<leader>l        Lazy (plugins)
<leader>?        Keybindings help`,
			Tips: []string{
				"<leader>? shows context-aware keybindings",
				"<leader>sk to fuzzy search all keymaps",
				"Most commands work with telescope/snacks",
			},
		},
	}
}

// GetLazyVimTopicTitles returns just the titles for menu display
func GetLazyVimTopicTitles() []string {
	topics := GetLazyVimTopics()
	titles := make([]string, len(topics))
	for i, t := range topics {
		titles[i] = t.Title
	}
	return titles
}
