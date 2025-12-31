package tui

import (
	"fmt"

	"github.com/Gentleman-Programming/Gentleman.Dots/installer/internal/system"
	tea "github.com/charmbracelet/bubbletea"
)

// Screen represents the current screen being displayed
type Screen int

const (
	ScreenWelcome Screen = iota
	ScreenMainMenu
	ScreenOSSelect
	ScreenTerminalSelect
	ScreenFontSelect
	ScreenShellSelect
	ScreenWMSelect
	ScreenNvimSelect
	ScreenInstalling
	ScreenComplete
	ScreenError
	// Learn screens
	ScreenLearnTerminals
	ScreenLearnShells
	ScreenLearnWM
	ScreenLearnNvim
	// Keymaps screen
	ScreenKeymaps
	ScreenKeymapCategory
	// LazyVim learn screens
	ScreenLearnLazyVim
	ScreenLazyVimTopic
	// Backup screens
	ScreenBackupConfirm
	ScreenRestoreBackup
	ScreenRestoreConfirm
)

// InstallStep represents a single installation step
type InstallStep struct {
	ID          string
	Name        string
	Description string
	Status      StepStatus
	Progress    float64
	Error       error
}

type StepStatus int

const (
	StatusPending StepStatus = iota
	StatusRunning
	StatusDone
	StatusFailed
	StatusSkipped
)

// UserChoices stores all user selections
type UserChoices struct {
	OS           string // "mac", "linux"
	Terminal     string // "alacritty", "wezterm", "kitty", "ghostty", "none"
	InstallFont  bool
	Shell        string // "fish", "zsh", "nushell"
	WindowMgr    string // "tmux", "zellij", "none"
	InstallNvim  bool
	CreateBackup bool // Whether to backup existing configs
}

// Model is the main application state
type Model struct {
	Screen      Screen
	PrevScreen  Screen // For going back from learn/keymaps screens
	Width       int
	Height      int
	SystemInfo  *system.SystemInfo
	Choices     UserChoices
	Steps       []InstallStep
	CurrentStep int
	Cursor      int
	ErrorMsg    string
	ShowDetails bool
	LogLines    []string
	TotalTime   float64
	Quitting    bool
	// Program reference for sending messages during installation
	Program *tea.Program
	// Spinner animation
	SpinnerFrame int
	// Learn mode
	ViewingTool string // Current tool being viewed in learn mode
	// Keymaps mode
	KeymapCategories []KeymapCategory
	SelectedCategory int
	KeymapScroll     int // For scrolling through keymaps
	// LazyVim mode
	LazyVimTopics        []LazyVimTopic
	SelectedLazyVimTopic int
	LazyVimScroll        int // For scrolling through topic content
	// Backup mode
	ExistingConfigs  []string            // Configs that will be overwritten
	AvailableBackups []system.BackupInfo // Available backups for restore
	SelectedBackup   int                 // Selected backup index
	BackupDir        string              // Last backup directory created
}

// NewModel creates a new Model with initial state
func NewModel() Model {
	return Model{
		Screen:               ScreenWelcome,
		PrevScreen:           ScreenWelcome,
		Width:                80,
		Height:               24,
		SystemInfo:           system.Detect(),
		Choices:              UserChoices{},
		Steps:                []InstallStep{},
		CurrentStep:          0,
		Cursor:               0,
		ShowDetails:          false,
		LogLines:             []string{},
		SpinnerFrame:         0,
		KeymapCategories:     GetNvimKeymaps(),
		SelectedCategory:     0,
		KeymapScroll:         0,
		LazyVimTopics:        GetLazyVimTopics(),
		SelectedLazyVimTopic: 0,
		LazyVimScroll:        0,
		ExistingConfigs:      []string{},
		AvailableBackups:     []system.BackupInfo{},
		SelectedBackup:       0,
		BackupDir:            "",
		Program:              nil, // Will be set after tea.Program is created
	}
}

// SetProgram sets the tea.Program reference for sending messages during installation
func (m *Model) SetProgram(p *tea.Program) {
	m.Program = p
}

// globalProgram holds a reference to the tea.Program for sending logs during installation
var globalProgram *tea.Program

// SetGlobalProgram sets the global program reference
func SetGlobalProgram(p *tea.Program) {
	globalProgram = p
}

// SendLog sends a log message to the TUI during installation
func SendLog(stepID string, log string) {
	if globalProgram != nil {
		globalProgram.Send(stepProgressMsg{
			stepID: stepID,
			log:    log,
		})
	}
}

// SendLogLine is an alias for SendLog for compatibility
func (m *Model) SendLog(stepID string, log string) {
	SendLog(stepID, log)
}

// GetCurrentOptions returns the options for the current screen
func (m Model) GetCurrentOptions() []string {
	switch m.Screen {
	case ScreenMainMenu:
		opts := []string{
			"🚀 Start Installation",
			"📚 Learn About Tools",
			"⌨️  Neovim Keymaps Reference",
			"📖 LazyVim Guide",
		}
		// Add restore option if backups exist
		if len(m.AvailableBackups) > 0 {
			opts = append(opts, "🔄 Restore from Backup")
		}
		opts = append(opts, "❌ Exit")
		return opts
	case ScreenOSSelect:
		return []string{"macOS", "Linux"}
	case ScreenTerminalSelect:
		if m.Choices.OS == "mac" {
			return []string{"Alacritty", "WezTerm", "Kitty", "Ghostty", "None", "─────────────", "ℹ️  Learn about terminals"}
		}
		return []string{"Alacritty", "WezTerm", "Ghostty", "None", "─────────────", "ℹ️  Learn about terminals"}
	case ScreenFontSelect:
		return []string{"Yes, install Iosevka Term Nerd Font", "No, I already have it"}
	case ScreenShellSelect:
		return []string{"Fish", "Zsh", "Nushell", "─────────────", "ℹ️  Learn about shells"}
	case ScreenWMSelect:
		return []string{"Tmux", "Zellij", "None", "─────────────", "ℹ️  Learn about multiplexers"}
	case ScreenNvimSelect:
		return []string{"Yes, install Neovim with config", "No, skip Neovim", "─────────────", "ℹ️  Learn about Neovim", "⌨️  View Keymaps", "📖 LazyVim Guide"}
	case ScreenBackupConfirm:
		return []string{
			"✅ Install with Backup (recommended)",
			"⚠️  Install without Backup",
			"❌ Cancel",
		}
	case ScreenRestoreBackup:
		opts := make([]string, len(m.AvailableBackups)+2)
		for i, backup := range m.AvailableBackups {
			// Format: timestamp + file count
			opts[i] = fmt.Sprintf("%s (%d items)", backup.Timestamp.Format("2006-01-02 15:04:05"), len(backup.Files))
		}
		opts[len(m.AvailableBackups)] = "─────────────"
		opts[len(m.AvailableBackups)+1] = "← Back"
		return opts
	case ScreenRestoreConfirm:
		return []string{
			"✅ Yes, restore this backup",
			"🗑️  Delete this backup",
			"❌ Cancel",
		}
	case ScreenLearnTerminals:
		return []string{"Alacritty", "WezTerm", "Kitty", "Ghostty", "─────────────", "← Back"}
	case ScreenLearnShells:
		return []string{"Fish", "Zsh", "Nushell", "─────────────", "← Back"}
	case ScreenLearnWM:
		return []string{"Tmux", "Zellij", "─────────────", "← Back"}
	case ScreenLearnNvim:
		return []string{"View Features", "View Keymaps", "📖 LazyVim Guide", "─────────────", "← Back"}
	case ScreenKeymaps:
		categories := make([]string, len(m.KeymapCategories)+2)
		for i, cat := range m.KeymapCategories {
			categories[i] = cat.Name
		}
		categories[len(m.KeymapCategories)] = "─────────────"
		categories[len(m.KeymapCategories)+1] = "← Back"
		return categories
	case ScreenLearnLazyVim:
		titles := GetLazyVimTopicTitles()
		result := make([]string, len(titles)+2)
		copy(result, titles)
		result[len(titles)] = "─────────────"
		result[len(titles)+1] = "← Back"
		return result
	default:
		return []string{}
	}
}

// GetScreenTitle returns the title for the current screen
func (m Model) GetScreenTitle() string {
	switch m.Screen {
	case ScreenWelcome:
		return "Welcome to Gentleman.Dots Installer"
	case ScreenMainMenu:
		return "Main Menu"
	case ScreenOSSelect:
		return "Step 1: Select Your Operating System"
	case ScreenTerminalSelect:
		return "Step 2: Choose Terminal Emulator"
	case ScreenFontSelect:
		return "Step 3: Nerd Font Installation"
	case ScreenShellSelect:
		return "Step 4: Choose Your Shell"
	case ScreenWMSelect:
		return "Step 5: Choose Window Manager"
	case ScreenNvimSelect:
		return "Step 6: Neovim Configuration"
	case ScreenBackupConfirm:
		return "⚠️  Existing Configs Detected"
	case ScreenRestoreBackup:
		return "🔄 Restore from Backup"
	case ScreenRestoreConfirm:
		return "🔄 Confirm Restore"
	case ScreenInstalling:
		return "Installing..."
	case ScreenComplete:
		return "Installation Complete!"
	case ScreenError:
		return "Error"
	case ScreenLearnTerminals:
		return "📚 Learn: Terminal Emulators"
	case ScreenLearnShells:
		return "📚 Learn: Shells"
	case ScreenLearnWM:
		return "📚 Learn: Window Managers"
	case ScreenLearnNvim:
		return "📚 Learn: Neovim"
	case ScreenKeymaps:
		return "⌨️  Neovim Keymaps Reference"
	case ScreenKeymapCategory:
		if m.SelectedCategory < len(m.KeymapCategories) {
			return "⌨️  " + m.KeymapCategories[m.SelectedCategory].Name
		}
		return "⌨️  Keymaps"
	case ScreenLearnLazyVim:
		return "📖 LazyVim Guide"
	case ScreenLazyVimTopic:
		if m.SelectedLazyVimTopic < len(m.LazyVimTopics) {
			return "📖 " + m.LazyVimTopics[m.SelectedLazyVimTopic].Title
		}
		return "📖 LazyVim"
	default:
		return ""
	}
}

// GetScreenDescription returns a description for the current screen
func (m Model) GetScreenDescription() string {
	switch m.Screen {
	case ScreenOSSelect:
		detected := m.SystemInfo.OSName
		if m.SystemInfo.IsWSL {
			detected += " (WSL)"
		}
		return "Detected: " + detected
	case ScreenTerminalSelect:
		if m.SystemInfo.IsWSL {
			return "Note: Terminal emulators should be installed on Windows for WSL"
		}
		return "Select your preferred terminal emulator"
	case ScreenFontSelect:
		return "Iosevka Term Nerd Font is required for icons and glyphs"
	case ScreenShellSelect:
		return "Current shell: " + m.SystemInfo.UserShell
	case ScreenWMSelect:
		return "Terminal multiplexer for managing sessions"
	case ScreenNvimSelect:
		return "Includes LSP, TreeSitter, and Gentleman config"
	default:
		return ""
	}
}

// SetupInstallSteps creates the installation steps based on user choices
func (m *Model) SetupInstallSteps() {
	m.Steps = []InstallStep{}

	// Backup step if user chose to backup
	if m.Choices.CreateBackup && len(m.ExistingConfigs) > 0 {
		m.Steps = append(m.Steps, InstallStep{
			ID:          "backup",
			Name:        "Backup Existing Configs",
			Description: "Creating backup of your current configuration",
			Status:      StatusPending,
		})
	}

	// Always clone repo first
	m.Steps = append(m.Steps, InstallStep{
		ID:          "clone",
		Name:        "Clone Repository",
		Description: "Downloading Gentleman.Dots",
		Status:      StatusPending,
	})

	// Homebrew
	if !m.SystemInfo.HasBrew {
		m.Steps = append(m.Steps, InstallStep{
			ID:          "homebrew",
			Name:        "Install Homebrew",
			Description: "Package manager",
			Status:      StatusPending,
		})
	}

	// Dependencies based on OS
	if m.Choices.OS == "linux" {
		m.Steps = append(m.Steps, InstallStep{
			ID:          "deps",
			Name:        "Install Dependencies",
			Description: "Base packages",
			Status:      StatusPending,
		})
	} else if m.Choices.OS == "mac" && !m.SystemInfo.HasXcode {
		m.Steps = append(m.Steps, InstallStep{
			ID:          "xcode",
			Name:        "Install Xcode CLI",
			Description: "Developer tools",
			Status:      StatusPending,
		})
	}

	// Terminal
	if m.Choices.Terminal != "none" && m.Choices.Terminal != "" {
		m.Steps = append(m.Steps, InstallStep{
			ID:          "terminal",
			Name:        "Install " + m.Choices.Terminal,
			Description: "Terminal emulator",
			Status:      StatusPending,
		})
	}

	// Font
	if m.Choices.InstallFont {
		m.Steps = append(m.Steps, InstallStep{
			ID:          "font",
			Name:        "Install Iosevka Nerd Font",
			Description: "Nerd font with icons",
			Status:      StatusPending,
		})
	}

	// Shell
	m.Steps = append(m.Steps, InstallStep{
		ID:          "shell",
		Name:        "Install " + m.Choices.Shell,
		Description: "Shell and plugins",
		Status:      StatusPending,
	})

	// Window manager
	if m.Choices.WindowMgr != "none" && m.Choices.WindowMgr != "" {
		m.Steps = append(m.Steps, InstallStep{
			ID:          "wm",
			Name:        "Install " + m.Choices.WindowMgr,
			Description: "Terminal multiplexer",
			Status:      StatusPending,
		})
	}

	// Neovim
	if m.Choices.InstallNvim {
		m.Steps = append(m.Steps, InstallStep{
			ID:          "nvim",
			Name:        "Install Neovim",
			Description: "Editor with config",
			Status:      StatusPending,
		})
	}

	// Cleanup
	m.Steps = append(m.Steps, InstallStep{
		ID:          "cleanup",
		Name:        "Cleanup",
		Description: "Removing temporary files",
		Status:      StatusPending,
	})
}
