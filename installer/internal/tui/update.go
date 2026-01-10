package tui

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/Gentleman-Programming/Gentleman.Dots/installer/internal/system"
	"github.com/Gentleman-Programming/Gentleman.Dots/installer/internal/tui/trainer"
	tea "github.com/charmbracelet/bubbletea"
)

// Messages
type (
	// tickMsg is sent periodically for animations
	tickMsg time.Time

	// installStartMsg signals to start installation
	installStartMsg struct{}

	// stepCompleteMsg signals a step completed
	stepCompleteMsg struct {
		stepID string
		err    error
	}

	// stepProgressMsg updates progress of current step
	stepProgressMsg struct {
		stepID   string
		progress float64
		log      string
	}

	// installCompleteMsg signals all installation is done
	installCompleteMsg struct {
		totalTime float64
	}

	// loadBackupsMsg signals to load available backups
	loadBackupsMsg struct {
		backups []system.BackupInfo
	}

	// execFinishedMsg signals an interactive process finished
	execFinishedMsg struct {
		stepID string
		err    error
	}
)

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		tea.SetWindowTitle("Gentleman.Dots Installer"),
		tickCmd(),
		loadBackupsCmd(),
	)
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func loadBackupsCmd() tea.Cmd {
	return func() tea.Msg {
		backups := system.ListBackups()
		return loadBackupsMsg{backups: backups}
	}
}

// Update implements tea.Model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case tickMsg:
		// Animate spinner during installation
		if m.Screen == ScreenInstalling {
			m.SpinnerFrame++
		}
		// Continue ticking for animations
		return m, tickCmd()

	case installStartMsg:
		// Start the installation process
		return m, m.runNextStep()

	case stepProgressMsg:
		// Update progress
		for i := range m.Steps {
			if m.Steps[i].ID == msg.stepID {
				m.Steps[i].Progress = msg.progress
				break
			}
		}
		if msg.log != "" {
			m.LogLines = append(m.LogLines, msg.log)
			// Keep only last 20 lines
			if len(m.LogLines) > 20 {
				m.LogLines = m.LogLines[len(m.LogLines)-20:]
			}
		}
		return m, nil

	case stepCompleteMsg:
		// Mark step as complete
		for i := range m.Steps {
			if m.Steps[i].ID == msg.stepID {
				if msg.err != nil {
					m.Steps[i].Status = StatusFailed
					m.Steps[i].Error = msg.err
					m.Screen = ScreenError
					// Include step name in error message for clarity
					m.ErrorMsg = fmt.Sprintf("Step '%s' failed:\n%s", m.Steps[i].Name, msg.err.Error())
					return m, nil
				}
				m.Steps[i].Status = StatusDone
				m.Steps[i].Progress = 1.0
				break
			}
		}
		m.CurrentStep++
		return m, m.runNextStep()

	case installCompleteMsg:
		m.TotalTime = msg.totalTime
		m.Screen = ScreenComplete
		return m, nil

	case loadBackupsMsg:
		m.AvailableBackups = msg.backups
		return m, nil

	case execFinishedMsg:
		// Interactive process finished (sudo commands, chsh, etc)
		for i := range m.Steps {
			if m.Steps[i].ID == msg.stepID {
				if msg.err != nil {
					m.Steps[i].Status = StatusFailed
					m.Steps[i].Error = msg.err
					m.Screen = ScreenError
					// Include step name in error message for clarity
					m.ErrorMsg = fmt.Sprintf("Step '%s' failed:\n%s", m.Steps[i].Name, msg.err.Error())
					return m, nil
				}
				m.Steps[i].Status = StatusDone
				m.Steps[i].Progress = 1.0
				break
			}
		}
		m.CurrentStep++
		return m, m.runNextStep()

	case needsExecProcessMsg:
		// This step needs to run with tea.ExecProcess for interactive input
		return m, tea.ExecProcess(msg.cmd, func(err error) tea.Msg {
			return execFinishedMsg{stepID: msg.stepID, err: err}
		})
	}

	return m, nil
}

// execInteractiveCmd creates a tea.Cmd that runs an interactive process
// This suspends the TUI and gives full terminal control to the process
func execInteractiveCmd(stepID string, name string, args ...string) tea.Cmd {
	c := exec.Command(name, args...)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return execFinishedMsg{stepID: stepID, err: err}
	})
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	// ctrl+c always quits immediately (no leader needed)
	if key == "ctrl+c" {
		m.Quitting = true
		return m, tea.Quit
	}

	// Leader key mode: <space> activates, next key executes command
	// Commands: <space>q = quit, <space>d = toggle details
	if m.LeaderMode {
		m.LeaderMode = false // Reset leader mode
		switch key {
		case "q":
			// Quit application
			if m.Screen != ScreenInstalling {
				m.Quitting = true
				return m, tea.Quit
			}
			return m, nil
		case "d":
			// Toggle details during installation
			if m.Screen == ScreenInstalling {
				m.ShowDetails = !m.ShowDetails
			}
			return m, nil
		default:
			// Unknown leader command, ignore
			return m, nil
		}
	}

	// <space> activates leader mode EXCEPT in screens that need space for input
	// (Trainer screens use space in commands, Welcome screen uses space to continue)
	if key == " " {
		// Screens where space should NOT activate leader mode
		switch m.Screen {
		case ScreenWelcome:
			// Welcome screen: space continues to main menu
			m.Screen = ScreenMainMenu
			m.Cursor = 0
			return m, nil
		case ScreenComplete, ScreenError:
			// Complete/Error screens: space quits the app
			m.Quitting = true
			return m, tea.Quit
		case ScreenTrainerLesson, ScreenTrainerPractice, ScreenTrainerBoss:
			// Trainer input screens: space is part of the input, pass through
			// (handled below in screen-specific handlers)
		default:
			// All other screens: activate leader mode
			m.LeaderMode = true
			return m, nil
		}
	}

	// ESC goes back from content/learn screens (and cancels leader mode implicitly)
	if key == "esc" {
		return m.handleEscape()
	}

	// Screen-specific keys
	switch m.Screen {
	case ScreenWelcome:
		switch key {
		case "enter":
			m.Screen = ScreenMainMenu
			m.Cursor = 0
		}

	case ScreenMainMenu:
		return m.handleMainMenuKeys(key)

	case ScreenOSSelect, ScreenTerminalSelect, ScreenFontSelect, ScreenShellSelect, ScreenWMSelect, ScreenNvimSelect, ScreenGhosttyWarning:
		return m.handleSelectionKeys(key)

	case ScreenLearnTerminals, ScreenLearnShells, ScreenLearnWM, ScreenLearnNvim:
		return m.handleLearnMenuKeys(key)

	case ScreenKeymaps:
		return m.handleKeymapsMenuKeys(key)

	case ScreenKeymapCategory:
		return m.handleKeymapCategoryKeys(key)

	case ScreenKeymapsMenu:
		return m.handleToolKeymapsMenuKeys(key)

	case ScreenKeymapsTmux:
		return m.handleTmuxKeymapsMenuKeys(key)

	case ScreenKeymapsTmuxCat:
		return m.handleTmuxKeymapCategoryKeys(key)

	case ScreenKeymapsZellij:
		return m.handleZellijKeymapsMenuKeys(key)

	case ScreenKeymapsZellijCat:
		return m.handleZellijKeymapCategoryKeys(key)

	case ScreenKeymapsGhostty:
		return m.handleGhosttyKeymapsMenuKeys(key)

	case ScreenKeymapsGhosttyCat:
		return m.handleGhosttyKeymapCategoryKeys(key)

	case ScreenLearnLazyVim:
		return m.handleLazyVimMenuKeys(key)

	case ScreenLazyVimTopic:
		return m.handleLazyVimTopicKeys(key)

	case ScreenBackupConfirm:
		return m.handleBackupConfirmKeys(key)

	case ScreenRestoreBackup:
		return m.handleRestoreBackupKeys(key)

	case ScreenRestoreConfirm:
		return m.handleRestoreConfirmKeys(key)

	// Trainer screens
	case ScreenTrainerMenu:
		return m.handleTrainerMenuKeys(key)

	case ScreenTrainerLesson, ScreenTrainerPractice:
		return m.handleTrainerExerciseKeys(key)

	case ScreenTrainerBoss:
		return m.handleTrainerBossKeys(key)

	case ScreenTrainerResult:
		return m.handleTrainerResultKeys(key)

	case ScreenTrainerBossResult:
		return m.handleTrainerBossResultKeys(key)

	case ScreenComplete:
		switch key {
		case "enter", " ":
			m.Quitting = true
			return m, tea.Quit
		}

	case ScreenError:
		switch key {
		case "enter", " ":
			m.Quitting = true
			return m, tea.Quit
		case "r":
			// Retry - go back to beginning
			m.Screen = ScreenWelcome
			m.ErrorMsg = ""
		}
	}

	return m, nil
}

func (m Model) handleEscape() (tea.Model, tea.Cmd) {
	switch m.Screen {
	// Installation wizard screens - go back through the flow
	case ScreenOSSelect, ScreenTerminalSelect, ScreenFontSelect, ScreenShellSelect, ScreenWMSelect, ScreenNvimSelect:
		return m.goBackInstallStep()
	case ScreenGhosttyWarning:
		// Go back to terminal selection
		m.Screen = ScreenTerminalSelect
		m.Cursor = 0
	case ScreenBackupConfirm:
		// Go back to Nvim selection (not abort)
		m.Screen = ScreenNvimSelect
		m.Cursor = 0
	// Content/Learn screens
	case ScreenKeymapCategory:
		m.Screen = ScreenKeymaps
		m.KeymapScroll = 0
	case ScreenKeymapsTmuxCat:
		m.Screen = ScreenKeymapsTmux
		m.TmuxKeymapScroll = 0
	case ScreenKeymapsZellijCat:
		m.Screen = ScreenKeymapsZellij
		m.ZellijKeymapScroll = 0
	case ScreenKeymapsGhosttyCat:
		m.Screen = ScreenKeymapsGhostty
		m.GhosttyKeymapScroll = 0
	case ScreenLazyVimTopic:
		m.Screen = ScreenLearnLazyVim
		m.LazyVimScroll = 0
	case ScreenLearnTerminals, ScreenLearnShells, ScreenLearnWM, ScreenLearnNvim:
		m.Screen = m.PrevScreen
		m.Cursor = 0
		m.ViewingTool = ""
	case ScreenKeymaps:
		m.Screen = ScreenKeymapsMenu
		m.Cursor = 0
	case ScreenKeymapsTmux, ScreenKeymapsZellij, ScreenKeymapsGhostty:
		m.Screen = ScreenKeymapsMenu
		m.Cursor = 0
	case ScreenKeymapsMenu, ScreenLearnLazyVim:
		m.Screen = m.PrevScreen
		m.Cursor = 0
	// Restore screens
	case ScreenRestoreBackup, ScreenRestoreConfirm:
		m.Screen = ScreenMainMenu
		m.Cursor = 0
	// Trainer screens
	case ScreenTrainerMenu:
		// Save stats and return to main menu
		if m.TrainerStats != nil {
			trainer.SaveStats(m.TrainerStats)
		}
		m.Screen = ScreenMainMenu
		m.Cursor = 0
	case ScreenTrainerLesson, ScreenTrainerPractice, ScreenTrainerBoss:
		// Return to trainer menu (stats saved in handlers)
		m.Screen = ScreenTrainerMenu
		m.TrainerMessage = ""
	case ScreenTrainerResult, ScreenTrainerBossResult:
		// Return to trainer menu
		if m.TrainerStats != nil {
			trainer.SaveStats(m.TrainerStats)
		}
		m.Screen = ScreenTrainerMenu
		m.TrainerMessage = ""
	// Main menu - quit
	case ScreenMainMenu:
		m.Quitting = true
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) handleMainMenuKeys(key string) (tea.Model, tea.Cmd) {
	options := m.GetCurrentOptions()
	hasRestoreOption := len(m.AvailableBackups) > 0

	switch key {
	case "up", "k":
		if m.Cursor > 0 {
			m.Cursor--
		}
	case "down", "j":
		if m.Cursor < len(options)-1 {
			m.Cursor++
		}
	case "enter", " ":
		selected := options[m.Cursor]
		switch {
		case strings.Contains(selected, "Start Installation"):
			m.Screen = ScreenOSSelect
			// Pre-select detected OS
			if m.SystemInfo.OS == system.OSLinux {
				m.Cursor = 1 // Linux is second option
			} else {
				m.Cursor = 0 // macOS is first option (default)
			}
		case strings.Contains(selected, "Learn About Tools"):
			m.Screen = ScreenLearnTerminals
			m.PrevScreen = ScreenMainMenu
			m.Cursor = 0
		case strings.Contains(selected, "Keymaps Reference"):
			m.Screen = ScreenKeymapsMenu
			m.PrevScreen = ScreenMainMenu
			m.Cursor = 0
		case strings.Contains(selected, "LazyVim Guide"):
			m.Screen = ScreenLearnLazyVim
			m.PrevScreen = ScreenMainMenu
			m.Cursor = 0
		case strings.Contains(selected, "Vim Trainer"):
			// Load user stats when entering trainer
			stats := trainer.LoadStats()
			if stats == nil {
				stats = trainer.NewUserStats()
			}
			m.TrainerStats = stats
			m.TrainerGameState = nil
			m.TrainerCursor = 0
			m.TrainerInput = ""
			m.Screen = ScreenTrainerMenu
			m.PrevScreen = ScreenMainMenu
		case strings.Contains(selected, "Restore from Backup") && hasRestoreOption:
			m.Screen = ScreenRestoreBackup
			m.Cursor = 0
		case strings.Contains(selected, "Exit"):
			m.Quitting = true
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) handleSelectionKeys(key string) (tea.Model, tea.Cmd) {
	options := m.GetCurrentOptions()

	switch key {
	case "up", "k":
		if m.Cursor > 0 {
			m.Cursor--
			// Skip separator lines
			if strings.HasPrefix(options[m.Cursor], "───") {
				if m.Cursor > 0 {
					m.Cursor--
				}
			}
		}

	case "down", "j":
		if m.Cursor < len(options)-1 {
			m.Cursor++
			// Skip separator lines
			if strings.HasPrefix(options[m.Cursor], "───") {
				if m.Cursor < len(options)-1 {
					m.Cursor++
				}
			}
		}

	case "esc", "backspace":
		// Go back to previous installation step
		return m.goBackInstallStep()

	case "enter", " ":
		return m.handleSelection()
	}

	return m, nil
}

// goBackInstallStep handles going back during installation wizard
func (m Model) goBackInstallStep() (tea.Model, tea.Cmd) {
	switch m.Screen {
	case ScreenOSSelect:
		// Go back to main menu
		m.Screen = ScreenMainMenu
		m.Cursor = 0
		// Reset choices
		m.Choices = UserChoices{}

	case ScreenTerminalSelect:
		m.Screen = ScreenOSSelect
		m.Cursor = 0
		// Reset terminal choice
		m.Choices.Terminal = ""

	case ScreenFontSelect:
		m.Screen = ScreenTerminalSelect
		m.Cursor = 0
		// Reset font choice
		m.Choices.InstallFont = false

	case ScreenShellSelect:
		// Termux: go back to OS selection (skipped terminal and font)
		if m.SystemInfo.IsTermux {
			m.Screen = ScreenOSSelect
		} else if m.Choices.Terminal == "none" {
			// If we skipped font selection (terminal = none), go back to terminal
			m.Screen = ScreenTerminalSelect
		} else {
			m.Screen = ScreenFontSelect
		}
		m.Cursor = 0
		m.Choices.Shell = ""

	case ScreenWMSelect:
		m.Screen = ScreenShellSelect
		m.Cursor = 0
		m.Choices.WindowMgr = ""

	case ScreenNvimSelect:
		m.Screen = ScreenWMSelect
		m.Cursor = 0
		m.Choices.InstallNvim = false
	}

	return m, nil
}

func (m Model) handleSelection() (tea.Model, tea.Cmd) {
	options := m.GetCurrentOptions()
	if m.Cursor >= len(options) {
		return m, nil
	}

	selected := strings.ToLower(options[m.Cursor])

	// Check for "learn" options
	if strings.Contains(selected, "learn about terminals") {
		m.PrevScreen = m.Screen
		m.Screen = ScreenLearnTerminals
		m.Cursor = 0
		return m, nil
	}
	if strings.Contains(selected, "learn about shells") {
		m.PrevScreen = m.Screen
		m.Screen = ScreenLearnShells
		m.Cursor = 0
		return m, nil
	}
	if strings.Contains(selected, "learn about multiplexers") {
		m.PrevScreen = m.Screen
		m.Screen = ScreenLearnWM
		m.Cursor = 0
		return m, nil
	}
	if strings.Contains(selected, "learn about neovim") {
		m.PrevScreen = m.Screen
		m.Screen = ScreenLearnNvim
		m.Cursor = 0
		return m, nil
	}
	if strings.Contains(selected, "view keymaps") {
		m.PrevScreen = m.Screen
		m.Screen = ScreenKeymaps
		m.Cursor = 0
		return m, nil
	}
	if strings.Contains(selected, "lazyvim guide") {
		m.PrevScreen = m.Screen
		m.Screen = ScreenLearnLazyVim
		m.Cursor = 0
		return m, nil
	}

	// Skip separators
	if strings.HasPrefix(selected, "───") {
		return m, nil
	}

	switch m.Screen {
	case ScreenOSSelect:
		selectedLower := strings.ToLower(selected)
		if strings.Contains(selectedLower, "mac") {
			m.Choices.OS = "mac"
		} else if strings.Contains(selectedLower, "termux") {
			m.Choices.OS = "termux"
		} else {
			m.Choices.OS = "linux"
		}
		// Termux: skip Terminal selection (you're already in a terminal!)
		// But allow font installation (Termux supports custom fonts)
		if m.Choices.OS == "termux" {
			m.Choices.Terminal = "none"
			m.Choices.InstallFont = true // Install Nerd Font for Termux
			m.Screen = ScreenShellSelect
		} else {
			m.Screen = ScreenTerminalSelect
		}
		m.Cursor = 0

	case ScreenTerminalSelect:
		term := strings.ToLower(strings.Split(options[m.Cursor], " ")[0])
		m.Choices.Terminal = term

		// Check if Ghostty on Debian/Ubuntu - show warning
		if term == "ghostty" && m.Choices.OS == "linux" && m.SystemInfo.OS == system.OSDebian && !system.CommandExists("ghostty") {
			m.Screen = ScreenGhosttyWarning
			m.Cursor = 0
			return m, nil
		}

		if term != "none" {
			m.Screen = ScreenFontSelect
		} else {
			m.Screen = ScreenShellSelect
		}
		m.Cursor = 0

	case ScreenFontSelect:
		m.Choices.InstallFont = m.Cursor == 0
		m.Screen = ScreenShellSelect
		m.Cursor = 0

	case ScreenShellSelect:
		m.Choices.Shell = strings.ToLower(options[m.Cursor])
		m.Screen = ScreenWMSelect
		m.Cursor = 0

	case ScreenGhosttyWarning:
		switch m.Cursor {
		case 0: // Continue with Ghostty anyway
			m.Screen = ScreenFontSelect
			m.Cursor = 0
		case 1: // Choose different terminal
			m.Screen = ScreenTerminalSelect
			m.Cursor = 0
		case 2: // Cancel
			m.Screen = ScreenMainMenu
			m.Cursor = 0
		}

	case ScreenWMSelect:
		m.Choices.WindowMgr = strings.ToLower(options[m.Cursor])
		m.Screen = ScreenNvimSelect
		m.Cursor = 0

	case ScreenNvimSelect:
		m.Choices.InstallNvim = m.Cursor == 0
		// Detect existing configs before proceeding
		m.ExistingConfigs = system.DetectExistingConfigs()
		if len(m.ExistingConfigs) > 0 {
			// Show backup confirmation screen
			m.Screen = ScreenBackupConfirm
			m.Cursor = 0
		} else {
			// No existing configs, proceed directly
			m.SetupInstallSteps()
			m.Screen = ScreenInstalling
			m.CurrentStep = 0
			return m, func() tea.Msg { return installStartMsg{} }
		}
	}

	return m, nil
}

func (m Model) handleLearnMenuKeys(key string) (tea.Model, tea.Cmd) {
	options := m.GetCurrentOptions()

	switch key {
	case "up", "k":
		if m.Cursor > 0 {
			m.Cursor--
			if strings.HasPrefix(options[m.Cursor], "───") && m.Cursor > 0 {
				m.Cursor--
			}
		}
	case "down", "j":
		if m.Cursor < len(options)-1 {
			m.Cursor++
			if strings.HasPrefix(options[m.Cursor], "───") && m.Cursor < len(options)-1 {
				m.Cursor++
			}
		}
	case "enter", " ":
		selected := options[m.Cursor]
		if strings.Contains(selected, "Back") {
			m.Screen = m.PrevScreen
			m.Cursor = 0
			m.ViewingTool = ""
			return m, nil
		}
		if strings.HasPrefix(selected, "───") {
			return m, nil
		}

		// Handle Learn Nvim special options
		if m.Screen == ScreenLearnNvim {
			switch m.Cursor {
			case 0: // View Features
				m.ViewingTool = "features"
			case 1: // View Keymaps
				m.Screen = ScreenKeymaps
				m.PrevScreen = ScreenLearnNvim
				m.Cursor = 0
				return m, nil
			case 2: // LazyVim Guide
				m.Screen = ScreenLearnLazyVim
				m.PrevScreen = ScreenLearnNvim
				m.Cursor = 0
				return m, nil
			}
			return m, nil
		}

		// Set viewing tool for other learn screens
		m.ViewingTool = strings.ToLower(selected)
	}

	return m, nil
}

func (m Model) handleKeymapsMenuKeys(key string) (tea.Model, tea.Cmd) {
	options := m.GetCurrentOptions()

	switch key {
	case "up", "k":
		if m.Cursor > 0 {
			m.Cursor--
			if strings.HasPrefix(options[m.Cursor], "───") && m.Cursor > 0 {
				m.Cursor--
			}
		}
	case "down", "j":
		if m.Cursor < len(options)-1 {
			m.Cursor++
			if strings.HasPrefix(options[m.Cursor], "───") && m.Cursor < len(options)-1 {
				m.Cursor++
			}
		}
	case "enter", " ":
		selected := options[m.Cursor]
		if strings.Contains(selected, "Back") {
			m.Screen = m.PrevScreen
			m.Cursor = 0
			return m, nil
		}
		if strings.HasPrefix(selected, "───") {
			return m, nil
		}

		// Select category and show keymaps
		m.SelectedCategory = m.Cursor
		m.Screen = ScreenKeymapCategory
		m.KeymapScroll = 0
	}

	return m, nil
}

func (m Model) handleKeymapCategoryKeys(key string) (tea.Model, tea.Cmd) {
	category := m.KeymapCategories[m.SelectedCategory]

	// Calculate visible items based on terminal height (same as view)
	visibleItems := m.Height - 9
	if visibleItems < 5 {
		visibleItems = 5
	}

	maxScroll := len(category.Keymaps) - visibleItems
	if maxScroll < 0 {
		maxScroll = 0
	}

	switch key {
	case "up", "k":
		if m.KeymapScroll > 0 {
			m.KeymapScroll--
		}
	case "down", "j":
		if m.KeymapScroll < maxScroll {
			m.KeymapScroll++
		}
	case "enter", " ", "q", "esc":
		m.Screen = ScreenKeymaps
		m.KeymapScroll = 0
	}

	return m, nil
}

// handleToolKeymapsMenuKeys handles the tool selection menu (Neovim, Tmux, Zellij, Ghostty)
func (m Model) handleToolKeymapsMenuKeys(key string) (tea.Model, tea.Cmd) {
	options := m.GetCurrentOptions()

	switch key {
	case "up", "k":
		if m.Cursor > 0 {
			m.Cursor--
			if strings.HasPrefix(options[m.Cursor], "───") && m.Cursor > 0 {
				m.Cursor--
			}
		}
	case "down", "j":
		if m.Cursor < len(options)-1 {
			m.Cursor++
			if strings.HasPrefix(options[m.Cursor], "───") && m.Cursor < len(options)-1 {
				m.Cursor++
			}
		}
	case "enter", " ":
		selected := options[m.Cursor]
		if strings.Contains(selected, "Back") {
			m.Screen = m.PrevScreen
			m.Cursor = 0
			return m, nil
		}
		if strings.HasPrefix(selected, "───") {
			return m, nil
		}

		// Navigate to specific tool's keymaps
		switch m.Cursor {
		case 0: // Neovim
			m.Screen = ScreenKeymaps
			m.Cursor = 0
		case 1: // Tmux
			m.Screen = ScreenKeymapsTmux
			m.Cursor = 0
		case 2: // Zellij
			m.Screen = ScreenKeymapsZellij
			m.Cursor = 0
		case 3: // Ghostty
			m.Screen = ScreenKeymapsGhostty
			m.Cursor = 0
		}
	}

	return m, nil
}

// handleTmuxKeymapsMenuKeys handles Tmux keymap category selection
func (m Model) handleTmuxKeymapsMenuKeys(key string) (tea.Model, tea.Cmd) {
	options := m.GetCurrentOptions()

	switch key {
	case "up", "k":
		if m.Cursor > 0 {
			m.Cursor--
			if strings.HasPrefix(options[m.Cursor], "───") && m.Cursor > 0 {
				m.Cursor--
			}
		}
	case "down", "j":
		if m.Cursor < len(options)-1 {
			m.Cursor++
			if strings.HasPrefix(options[m.Cursor], "───") && m.Cursor < len(options)-1 {
				m.Cursor++
			}
		}
	case "enter", " ":
		selected := options[m.Cursor]
		if strings.Contains(selected, "Back") {
			m.Screen = ScreenKeymapsMenu
			m.Cursor = 0
			return m, nil
		}
		if strings.HasPrefix(selected, "───") {
			return m, nil
		}

		// Select category and show keymaps
		m.TmuxSelectedCategory = m.Cursor
		m.Screen = ScreenKeymapsTmuxCat
		m.TmuxKeymapScroll = 0
	}

	return m, nil
}

// handleTmuxKeymapCategoryKeys handles scrolling in Tmux keymap category view
func (m Model) handleTmuxKeymapCategoryKeys(key string) (tea.Model, tea.Cmd) {
	category := m.TmuxKeymapCategories[m.TmuxSelectedCategory]

	visibleItems := m.Height - 9
	if visibleItems < 5 {
		visibleItems = 5
	}

	maxScroll := len(category.Keymaps) - visibleItems
	if maxScroll < 0 {
		maxScroll = 0
	}

	switch key {
	case "up", "k":
		if m.TmuxKeymapScroll > 0 {
			m.TmuxKeymapScroll--
		}
	case "down", "j":
		if m.TmuxKeymapScroll < maxScroll {
			m.TmuxKeymapScroll++
		}
	case "enter", " ", "q", "esc":
		m.Screen = ScreenKeymapsTmux
		m.TmuxKeymapScroll = 0
	}

	return m, nil
}

// handleZellijKeymapsMenuKeys handles Zellij keymap category selection
func (m Model) handleZellijKeymapsMenuKeys(key string) (tea.Model, tea.Cmd) {
	options := m.GetCurrentOptions()

	switch key {
	case "up", "k":
		if m.Cursor > 0 {
			m.Cursor--
			if strings.HasPrefix(options[m.Cursor], "───") && m.Cursor > 0 {
				m.Cursor--
			}
		}
	case "down", "j":
		if m.Cursor < len(options)-1 {
			m.Cursor++
			if strings.HasPrefix(options[m.Cursor], "───") && m.Cursor < len(options)-1 {
				m.Cursor++
			}
		}
	case "enter", " ":
		selected := options[m.Cursor]
		if strings.Contains(selected, "Back") {
			m.Screen = ScreenKeymapsMenu
			m.Cursor = 0
			return m, nil
		}
		if strings.HasPrefix(selected, "───") {
			return m, nil
		}

		// Select category and show keymaps
		m.ZellijSelectedCategory = m.Cursor
		m.Screen = ScreenKeymapsZellijCat
		m.ZellijKeymapScroll = 0
	}

	return m, nil
}

// handleZellijKeymapCategoryKeys handles scrolling in Zellij keymap category view
func (m Model) handleZellijKeymapCategoryKeys(key string) (tea.Model, tea.Cmd) {
	category := m.ZellijKeymapCategories[m.ZellijSelectedCategory]

	visibleItems := m.Height - 9
	if visibleItems < 5 {
		visibleItems = 5
	}

	maxScroll := len(category.Keymaps) - visibleItems
	if maxScroll < 0 {
		maxScroll = 0
	}

	switch key {
	case "up", "k":
		if m.ZellijKeymapScroll > 0 {
			m.ZellijKeymapScroll--
		}
	case "down", "j":
		if m.ZellijKeymapScroll < maxScroll {
			m.ZellijKeymapScroll++
		}
	case "enter", " ", "q", "esc":
		m.Screen = ScreenKeymapsZellij
		m.ZellijKeymapScroll = 0
	}

	return m, nil
}

// handleGhosttyKeymapsMenuKeys handles Ghostty keymap category selection
func (m Model) handleGhosttyKeymapsMenuKeys(key string) (tea.Model, tea.Cmd) {
	options := m.GetCurrentOptions()

	switch key {
	case "up", "k":
		if m.Cursor > 0 {
			m.Cursor--
			if strings.HasPrefix(options[m.Cursor], "───") && m.Cursor > 0 {
				m.Cursor--
			}
		}
	case "down", "j":
		if m.Cursor < len(options)-1 {
			m.Cursor++
			if strings.HasPrefix(options[m.Cursor], "───") && m.Cursor < len(options)-1 {
				m.Cursor++
			}
		}
	case "enter", " ":
		selected := options[m.Cursor]
		if strings.Contains(selected, "Back") {
			m.Screen = ScreenKeymapsMenu
			m.Cursor = 0
			return m, nil
		}
		if strings.HasPrefix(selected, "───") {
			return m, nil
		}

		// Select category and show keymaps
		m.GhosttySelectedCategory = m.Cursor
		m.Screen = ScreenKeymapsGhosttyCat
		m.GhosttyKeymapScroll = 0
	}

	return m, nil
}

// handleGhosttyKeymapCategoryKeys handles scrolling in Ghostty keymap category view
func (m Model) handleGhosttyKeymapCategoryKeys(key string) (tea.Model, tea.Cmd) {
	category := m.GhosttyKeymapCategories[m.GhosttySelectedCategory]

	visibleItems := m.Height - 9
	if visibleItems < 5 {
		visibleItems = 5
	}

	maxScroll := len(category.Keymaps) - visibleItems
	if maxScroll < 0 {
		maxScroll = 0
	}

	switch key {
	case "up", "k":
		if m.GhosttyKeymapScroll > 0 {
			m.GhosttyKeymapScroll--
		}
	case "down", "j":
		if m.GhosttyKeymapScroll < maxScroll {
			m.GhosttyKeymapScroll++
		}
	case "enter", " ", "q", "esc":
		m.Screen = ScreenKeymapsGhostty
		m.GhosttyKeymapScroll = 0
	}

	return m, nil
}

func (m Model) handleLazyVimMenuKeys(key string) (tea.Model, tea.Cmd) {
	options := m.GetCurrentOptions()

	switch key {
	case "up", "k":
		if m.Cursor > 0 {
			m.Cursor--
			if strings.HasPrefix(options[m.Cursor], "───") && m.Cursor > 0 {
				m.Cursor--
			}
		}
	case "down", "j":
		if m.Cursor < len(options)-1 {
			m.Cursor++
			if strings.HasPrefix(options[m.Cursor], "───") && m.Cursor < len(options)-1 {
				m.Cursor++
			}
		}
	case "enter", " ":
		selected := options[m.Cursor]
		if strings.Contains(selected, "Back") {
			m.Screen = m.PrevScreen
			m.Cursor = 0
			return m, nil
		}
		if strings.HasPrefix(selected, "───") {
			return m, nil
		}

		// Select topic and show content
		m.SelectedLazyVimTopic = m.Cursor
		m.Screen = ScreenLazyVimTopic
		m.LazyVimScroll = 0
	}

	return m, nil
}

func (m Model) handleLazyVimTopicKeys(key string) (tea.Model, tea.Cmd) {
	topic := m.LazyVimTopics[m.SelectedLazyVimTopic]

	// Calculate view height based on terminal size (same as view)
	// Reserve space for: title(1) + description(1) + blank(2) + scroll info(2) + help(2) = 8 lines
	viewHeight := m.Height - 8
	if viewHeight < 10 {
		viewHeight = 10 // Minimum
	}

	// Calculate content height: content lines + code example lines + tips
	contentLines := len(topic.Content) + strings.Count(topic.CodeExample, "\n") + len(topic.Tips) + 10
	maxScroll := contentLines - viewHeight
	if maxScroll < 0 {
		maxScroll = 0
	}

	switch key {
	case "up", "k":
		if m.LazyVimScroll > 0 {
			m.LazyVimScroll--
		}
	case "down", "j":
		if m.LazyVimScroll < maxScroll {
			m.LazyVimScroll++
		}
	case "pgup":
		m.LazyVimScroll -= 10
		if m.LazyVimScroll < 0 {
			m.LazyVimScroll = 0
		}
	case "pgdown":
		m.LazyVimScroll += 10
		if m.LazyVimScroll > maxScroll {
			m.LazyVimScroll = maxScroll
		}
	case "enter", " ", "q", "esc":
		m.Screen = ScreenLearnLazyVim
		m.LazyVimScroll = 0
	}

	return m, nil
}

func (m Model) handleBackupConfirmKeys(key string) (tea.Model, tea.Cmd) {
	options := m.GetCurrentOptions()

	switch key {
	case "up", "k":
		if m.Cursor > 0 {
			m.Cursor--
		}
	case "down", "j":
		if m.Cursor < len(options)-1 {
			m.Cursor++
		}
	case "enter", " ":
		switch m.Cursor {
		case 0: // Install with Backup
			m.Choices.CreateBackup = true
			m.SetupInstallSteps()
			m.Screen = ScreenInstalling
			m.CurrentStep = 0
			return m, func() tea.Msg { return installStartMsg{} }
		case 1: // Install without Backup
			m.Choices.CreateBackup = false
			m.SetupInstallSteps()
			m.Screen = ScreenInstalling
			m.CurrentStep = 0
			return m, func() tea.Msg { return installStartMsg{} }
		case 2: // Cancel - abort the entire wizard
			m.Screen = ScreenMainMenu
			m.Cursor = 0
			// Reset choices when canceling
			m.Choices = UserChoices{}
		}
	case "esc", "backspace":
		// Go back to Nvim selection
		m.Screen = ScreenNvimSelect
		m.Cursor = 0
	}

	return m, nil
}

func (m Model) handleRestoreBackupKeys(key string) (tea.Model, tea.Cmd) {
	options := m.GetCurrentOptions()

	switch key {
	case "up", "k":
		if m.Cursor > 0 {
			m.Cursor--
			// Skip separator
			if strings.HasPrefix(options[m.Cursor], "───") && m.Cursor > 0 {
				m.Cursor--
			}
		}
	case "down", "j":
		if m.Cursor < len(options)-1 {
			m.Cursor++
			// Skip separator
			if strings.HasPrefix(options[m.Cursor], "───") && m.Cursor < len(options)-1 {
				m.Cursor++
			}
		}
	case "enter", " ":
		// Check if Back option
		if strings.Contains(options[m.Cursor], "Back") {
			m.Screen = ScreenMainMenu
			m.Cursor = 0
			return m, nil
		}
		// Skip separator
		if strings.HasPrefix(options[m.Cursor], "───") {
			return m, nil
		}
		// Select a backup
		if m.Cursor < len(m.AvailableBackups) {
			m.SelectedBackup = m.Cursor
			m.Screen = ScreenRestoreConfirm
			m.Cursor = 0
		}
	case "esc":
		m.Screen = ScreenMainMenu
		m.Cursor = 0
	}

	return m, nil
}

func (m Model) handleRestoreConfirmKeys(key string) (tea.Model, tea.Cmd) {
	options := m.GetCurrentOptions()

	switch key {
	case "up", "k":
		if m.Cursor > 0 {
			m.Cursor--
		}
	case "down", "j":
		if m.Cursor < len(options)-1 {
			m.Cursor++
		}
	case "enter", " ":
		backup := m.AvailableBackups[m.SelectedBackup]
		switch m.Cursor {
		case 0: // Restore
			err := system.RestoreBackup(backup.Path)
			if err != nil {
				m.Screen = ScreenError
				m.ErrorMsg = "Failed to restore backup: " + err.Error()
				return m, nil
			}
			// Refresh backups list
			m.AvailableBackups = system.ListBackups()
			m.Screen = ScreenComplete
			m.Choices = UserChoices{} // Clear choices to indicate restore
		case 1: // Delete
			_ = system.DeleteBackup(backup.Path)
			// Refresh backups list
			m.AvailableBackups = system.ListBackups()
			m.Screen = ScreenRestoreBackup
			m.Cursor = 0
			m.SelectedBackup = 0
		case 2: // Cancel
			m.Screen = ScreenRestoreBackup
			m.Cursor = m.SelectedBackup
		}
	case "esc":
		m.Screen = ScreenRestoreBackup
		m.Cursor = m.SelectedBackup
	}

	return m, nil
}

// runNextStep starts the next installation step
func (m Model) runNextStep() tea.Cmd {
	if m.CurrentStep >= len(m.Steps) {
		return func() tea.Msg {
			return installCompleteMsg{totalTime: 0}
		}
	}

	step := &m.Steps[m.CurrentStep]
	step.Status = StatusRunning

	// Check if this step needs interactive input (sudo, chsh, etc)
	if step.Interactive {
		return runInteractiveStep(step.ID, &m)
	}

	return func() tea.Msg {
		// Execute the step
		err := executeStep(step.ID, &m)
		return stepCompleteMsg{stepID: step.ID, err: err}
	}
}

// ============================================================================
// Trainer Handlers
// ============================================================================

// handleTrainerMenuKeys handles module selection in the trainer
func (m Model) handleTrainerMenuKeys(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "up", "k":
		if m.TrainerCursor > 0 {
			m.TrainerCursor--
		}
	case "down", "j":
		if m.TrainerCursor < len(m.TrainerModules)-1 {
			m.TrainerCursor++
		}
	case "enter", " ":
		// Select module and start lesson
		module := m.TrainerModules[m.TrainerCursor]

		if !m.TrainerStats.IsModuleUnlocked(module.ID) {
			m.TrainerMessage = "🔒 Module locked! Complete previous boss first."
			return m, nil
		}

		// Start lessons for the module
		lessons := trainer.GetLessons(module.ID)
		if len(lessons) == 0 {
			m.TrainerMessage = "No lessons available for this module yet."
			return m, nil
		}

		// Initialize game state with lesson count
		m.TrainerGameState = trainer.NewGameStateWithStats(m.TrainerStats)
		progress := m.TrainerStats.GetModuleProgress(module.ID)
		progress.LessonsTotal = len(lessons)
		m.TrainerGameState.StartLesson(module.ID)
		m.TrainerInput = ""
		m.TrainerMessage = ""
		m.Screen = ScreenTrainerLesson
	case "l":
		// L key for Lesson mode (if unlocked)
		if m.TrainerCursor < len(m.TrainerModules) {
			module := m.TrainerModules[m.TrainerCursor]
			if m.TrainerStats.IsModuleUnlocked(module.ID) {
				lessons := trainer.GetLessons(module.ID)
				if len(lessons) > 0 {
					m.TrainerGameState = trainer.NewGameStateWithStats(m.TrainerStats)
					progress := m.TrainerStats.GetModuleProgress(module.ID)
					progress.LessonsTotal = len(lessons)
					m.TrainerGameState.StartLesson(module.ID)
					m.TrainerInput = ""
					m.TrainerMessage = ""
					m.Screen = ScreenTrainerLesson
				}
			}
		}
	case "p":
		// P key for Practice mode (if ready)
		if m.TrainerCursor < len(m.TrainerModules) {
			module := m.TrainerModules[m.TrainerCursor]
			if m.TrainerStats.IsPracticeReady(module.ID) {
				// Check if practice is complete
				progress := m.TrainerStats.GetModuleProgress(module.ID)
				if progress.IsPracticeComplete(module.ID) {
					m.TrainerMessage = "🎉 Practice complete! All exercises mastered! Press [r] to reset."
					return m, nil
				}

				m.TrainerGameState = trainer.NewGameStateWithStats(m.TrainerStats)
				m.TrainerGameState.StartPractice(module.ID)

				// Check if we got an exercise (shouldn't fail if not complete, but safety check)
				if m.TrainerGameState.CurrentExercise == nil {
					m.TrainerMessage = "🎉 Practice complete! All exercises mastered! Press [r] to reset."
					return m, nil
				}

				m.TrainerInput = ""
				m.TrainerMessage = ""
				m.Screen = ScreenTrainerPractice
			} else {
				m.TrainerMessage = "Complete all lessons first to unlock practice!"
			}
		}
	case "r":
		// R key to reset practice progress for selected module
		if m.TrainerCursor < len(m.TrainerModules) {
			module := m.TrainerModules[m.TrainerCursor]
			if m.TrainerStats.IsModuleUnlocked(module.ID) {
				progress := m.TrainerStats.GetModuleProgress(module.ID)
				progress.ResetModulePractice()
				trainer.SaveStats(m.TrainerStats)
				m.TrainerMessage = "🔄 Practice progress reset for " + module.Name + ". Try again!"
			} else {
				m.TrainerMessage = "🔒 Module locked. Complete previous boss first."
			}
		}
	case "b":
		// B key for Boss fight (if ready)
		if m.TrainerCursor < len(m.TrainerModules) {
			module := m.TrainerModules[m.TrainerCursor]
			if m.TrainerStats.IsBossReady(module.ID) {
				boss := trainer.GetBoss(module.ID)
				if boss != nil {
					m.TrainerGameState = trainer.NewGameStateWithStats(m.TrainerStats)
					m.TrainerGameState.StartBoss(module.ID)
					m.TrainerInput = ""
					m.TrainerMessage = ""
					m.Screen = ScreenTrainerBoss
				} else {
					m.TrainerMessage = "Boss not implemented yet!"
				}
			} else {
				m.TrainerMessage = "Complete lessons + 80% practice accuracy to fight boss!"
			}
		}
	case "esc", "q":
		// Save stats and go back to main menu
		if m.TrainerStats != nil {
			trainer.SaveStats(m.TrainerStats)
		}
		m.Screen = ScreenMainMenu
		m.Cursor = 0
	}

	return m, nil
}

// handleTrainerExerciseKeys handles input during lesson/practice exercises
func (m Model) handleTrainerExerciseKeys(key string) (tea.Model, tea.Cmd) {
	if m.TrainerGameState == nil {
		m.Screen = ScreenTrainerMenu
		return m, nil
	}

	exercise := m.TrainerGameState.CurrentExercise
	if exercise == nil {
		m.Screen = ScreenTrainerMenu
		return m, nil
	}

	switch key {
	case "esc":
		// Exit to menu, save progress
		if m.TrainerStats != nil {
			trainer.SaveStats(m.TrainerStats)
		}
		m.Screen = ScreenTrainerMenu
		m.TrainerMessage = ""
		return m, nil

	case "backspace":
		// Remove last character from input
		if len(m.TrainerInput) > 0 {
			m.TrainerInput = m.TrainerInput[:len(m.TrainerInput)-1]
		}
		return m, nil

	case "enter":
		// Submit answer
		if m.TrainerInput == "" {
			return m, nil
		}

		// Validate answer using detailed validation
		validation := trainer.ValidateAnswerDetailed(exercise, m.TrainerInput)

		if validation.IsCorrect {
			// Record correct answer - time and optimal flag
			// Using a fixed time of 10 seconds for now (can add actual timing later)
			m.TrainerGameState.RecordCorrectAnswer(10.0, validation.IsOptimal)
			m.TrainerLastCorrect = true

			if validation.IsOptimal {
				m.TrainerMessage = "✨ Perfect! Optimal solution!"
			} else if validation.IsInSolutions {
				// Valid predefined solution but not optimal
				m.TrainerMessage = "✓ Correct! But " + exercise.Optimal + " is more efficient."
			} else {
				// Creative solution that works but not in predefined list
				m.TrainerMessage = "✓ Correct! Creative solution! Optimal: " + exercise.Optimal
			}
		} else {
			m.TrainerGameState.RecordIncorrectAnswer()
			m.TrainerLastCorrect = false
			// Show all valid solutions, not just optimal
			m.TrainerMessage = "✗ Incorrect. Solutions: " + trainer.FormatSolutionsHint(exercise)
		}

		// Record practice result for intelligent practice system
		if m.TrainerGameState.IsPracticeMode && exercise.ID != "" {
			progress := m.TrainerStats.GetModuleProgress(m.TrainerGameState.CurrentModule)
			progress.RecordPracticeResult(exercise.ID, validation.IsCorrect)
			trainer.SaveStats(m.TrainerStats)
		}

		m.Screen = ScreenTrainerResult
		return m, nil

	case "tab":
		// Show hint
		m.TrainerMessage = "💡 Hint: " + exercise.Hint
		return m, nil

	default:
		// Add character to input (filter control keys)
		// Accept single chars and specific ctrl combinations used in Vim
		validCtrlKeys := map[string]bool{
			"ctrl+a": true, "ctrl+e": true, "ctrl+w": true,
			"ctrl+d": true, "ctrl+u": true, "ctrl+f": true, "ctrl+b": true,
		}
		if len(key) == 1 || validCtrlKeys[key] {
			// Handle ctrl combinations - convert to control character
			if strings.HasPrefix(key, "ctrl+") {
				// Convert ctrl+X to actual control character for simulator
				switch key {
				case "ctrl+d":
					m.TrainerInput += "\x04"
				case "ctrl+u":
					m.TrainerInput += "\x15"
				case "ctrl+f":
					m.TrainerInput += "\x06"
				case "ctrl+b":
					m.TrainerInput += "\x02"
				default:
					m.TrainerInput += key
				}
			} else if len(key) == 1 {
				m.TrainerInput += key
			}
		} else if key == "space" {
			m.TrainerInput += " "
		}
	}

	return m, nil
}

// handleTrainerBossKeys handles input during boss fights
func (m Model) handleTrainerBossKeys(key string) (tea.Model, tea.Cmd) {
	if m.TrainerGameState == nil || m.TrainerGameState.CurrentBoss == nil {
		m.Screen = ScreenTrainerMenu
		return m, nil
	}

	switch key {
	case "esc":
		// Forfeit boss fight
		if m.TrainerStats != nil {
			trainer.SaveStats(m.TrainerStats)
		}
		m.Screen = ScreenTrainerMenu
		m.TrainerMessage = "Boss fight abandoned!"
		return m, nil

	case "backspace":
		if len(m.TrainerInput) > 0 {
			m.TrainerInput = m.TrainerInput[:len(m.TrainerInput)-1]
		}
		return m, nil

	case "enter":
		if m.TrainerInput == "" {
			return m, nil
		}

		// Get current boss step
		boss := m.TrainerGameState.CurrentBoss
		if m.TrainerGameState.BossStep >= len(boss.Steps) {
			// Boss complete!
			m.TrainerGameState.RecordBossVictory()
			m.TrainerLastCorrect = true
			m.TrainerMessage = "🏆 VICTORY! You defeated " + boss.Name + "!"
			m.Screen = ScreenTrainerBossResult
			return m, nil
		}

		step := boss.Steps[m.TrainerGameState.BossStep]
		isCorrect := trainer.ValidateAnswer(&step.Exercise, m.TrainerInput)
		isOptimal := trainer.IsOptimalAnswer(&step.Exercise, m.TrainerInput)

		if isCorrect {
			// Move to next step
			m.TrainerGameState.BossStep++
			m.TrainerInput = ""

			if m.TrainerGameState.BossStep >= len(boss.Steps) {
				// Boss defeated!
				m.TrainerGameState.RecordBossVictory()
				m.TrainerLastCorrect = true
				m.TrainerMessage = "🏆 VICTORY! You defeated " + boss.Name + "!"
				m.Screen = ScreenTrainerBossResult
			} else {
				if isOptimal {
					m.TrainerMessage = "✨ Perfect! Next challenge..."
				} else {
					m.TrainerMessage = "✓ Good! (Optimal: " + step.Exercise.Optimal + ") Next..."
				}
			}
		} else {
			// Lose a life - SHOW THE CORRECT SOLUTION
			m.TrainerGameState.BossLives--
			m.TrainerInput = ""

			// Format the solution hint
			solutionHint := trainer.FormatSolutionsHint(&step.Exercise)

			if m.TrainerGameState.BossLives <= 0 {
				// Game over - show final solution
				m.TrainerLastCorrect = false
				m.TrainerMessage = "💀 DEFEATED! Solution was: " + solutionHint
				m.Screen = ScreenTrainerBossResult
			} else {
				// Still has lives - show solution and remaining lives
				livesStr := strings.Repeat("❤️", m.TrainerGameState.BossLives)
				m.TrainerMessage = "✗ Wrong! Was: " + solutionHint + " | Lives: " + livesStr
			}
		}

		return m, nil

	default:
		// Add character to input
		// Accept single chars and specific ctrl combinations used in Vim
		validCtrlKeys := map[string]bool{
			"ctrl+d": true, "ctrl+u": true, "ctrl+f": true, "ctrl+b": true,
		}
		if len(key) == 1 {
			m.TrainerInput += key
		} else if key == "space" {
			m.TrainerInput += " "
		} else if validCtrlKeys[key] {
			// Convert ctrl+X to actual control character for simulator
			switch key {
			case "ctrl+d":
				m.TrainerInput += "\x04"
			case "ctrl+u":
				m.TrainerInput += "\x15"
			case "ctrl+f":
				m.TrainerInput += "\x06"
			case "ctrl+b":
				m.TrainerInput += "\x02"
			}
		}
	}

	return m, nil
}

// handleTrainerResultKeys handles the result screen after an exercise
func (m Model) handleTrainerResultKeys(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "enter", " ":
		// Continue to next exercise
		if m.TrainerGameState == nil {
			m.Screen = ScreenTrainerMenu
			return m, nil
		}

		var hasNext bool
		if m.TrainerGameState.IsPracticeMode {
			// Use intelligent practice selection
			hasNext = m.TrainerGameState.NextPracticeExercise()
		} else {
			// Lesson mode uses sequential
			hasNext = m.TrainerGameState.NextExercise()
		}

		if hasNext {
			m.TrainerInput = ""
			m.TrainerMessage = ""
			if m.TrainerGameState.IsLessonMode {
				m.Screen = ScreenTrainerLesson
			} else {
				m.Screen = ScreenTrainerPractice
			}
		} else {
			// Session complete
			if m.TrainerStats != nil {
				trainer.SaveStats(m.TrainerStats)
			}

			if m.TrainerGameState.IsPracticeMode {
				m.TrainerMessage = "🎉 All exercises mastered! You're a Vim master! 🏆"
			} else {
				m.TrainerMessage = "🎉 Lesson complete! Practice mode unlocked!"
			}
			m.Screen = ScreenTrainerMenu
		}

	case "esc", "q":
		// Return to menu
		if m.TrainerStats != nil {
			trainer.SaveStats(m.TrainerStats)
		}
		m.Screen = ScreenTrainerMenu
	}

	return m, nil
}

// handleTrainerBossResultKeys handles the result screen after a boss fight
func (m Model) handleTrainerBossResultKeys(key string) (tea.Model, tea.Cmd) {
	switch key {
	case "enter", " ", "esc", "q":
		// Return to menu
		if m.TrainerStats != nil {
			trainer.SaveStats(m.TrainerStats)
		}
		m.Screen = ScreenTrainerMenu
		m.TrainerMessage = ""
	}

	return m, nil
}
