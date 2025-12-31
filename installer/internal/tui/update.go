package tui

import (
	"strings"
	"time"

	"github.com/Gentleman-Programming/Gentleman.Dots/installer/internal/system"
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
					m.ErrorMsg = msg.err.Error()
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
	}

	return m, nil
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	// Global keys
	switch key {
	case "ctrl+c":
		m.Quitting = true
		return m, tea.Quit

	case "q":
		// q to quit only when not installing and not in content view
		if m.Screen != ScreenInstalling && m.Screen != ScreenKeymapCategory && m.Screen != ScreenLazyVimTopic {
			m.Quitting = true
			return m, tea.Quit
		}
		// In content views, q goes back
		if m.Screen == ScreenKeymapCategory {
			m.Screen = ScreenKeymaps
			m.KeymapScroll = 0
			return m, nil
		}
		if m.Screen == ScreenLazyVimTopic {
			m.Screen = ScreenLearnLazyVim
			m.LazyVimScroll = 0
			return m, nil
		}
		return m, nil

	case "d":
		// Toggle details only during installation
		if m.Screen == ScreenInstalling {
			m.ShowDetails = !m.ShowDetails
		}
		return m, nil

	case "esc":
		// ESC goes back from content/learn screens
		return m.handleEscape()
	}

	// Screen-specific keys
	switch m.Screen {
	case ScreenWelcome:
		switch key {
		case "enter", " ":
			m.Screen = ScreenMainMenu
			m.Cursor = 0
		}

	case ScreenMainMenu:
		return m.handleMainMenuKeys(key)

	case ScreenOSSelect, ScreenTerminalSelect, ScreenFontSelect, ScreenShellSelect, ScreenWMSelect, ScreenNvimSelect:
		return m.handleSelectionKeys(key)

	case ScreenLearnTerminals, ScreenLearnShells, ScreenLearnWM, ScreenLearnNvim:
		return m.handleLearnMenuKeys(key)

	case ScreenKeymaps:
		return m.handleKeymapsMenuKeys(key)

	case ScreenKeymapCategory:
		return m.handleKeymapCategoryKeys(key)

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
	case ScreenKeymapCategory:
		m.Screen = ScreenKeymaps
		m.KeymapScroll = 0
	case ScreenLazyVimTopic:
		m.Screen = ScreenLearnLazyVim
		m.LazyVimScroll = 0
	case ScreenLearnTerminals, ScreenLearnShells, ScreenLearnWM, ScreenLearnNvim:
		m.Screen = m.PrevScreen
		m.Cursor = 0
		m.ViewingTool = ""
	case ScreenKeymaps, ScreenLearnLazyVim:
		m.Screen = m.PrevScreen
		m.Cursor = 0
	case ScreenBackupConfirm:
		// Go back to Nvim selection (not abort)
		m.Screen = ScreenNvimSelect
		m.Cursor = 0
	case ScreenRestoreBackup, ScreenRestoreConfirm:
		m.Screen = ScreenMainMenu
		m.Cursor = 0
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
			m.Cursor = 0
		case strings.Contains(selected, "Learn About Tools"):
			m.Screen = ScreenLearnTerminals
			m.PrevScreen = ScreenMainMenu
			m.Cursor = 0
		case strings.Contains(selected, "Neovim Keymaps"):
			m.Screen = ScreenKeymaps
			m.PrevScreen = ScreenMainMenu
			m.Cursor = 0
		case strings.Contains(selected, "LazyVim Guide"):
			m.Screen = ScreenLearnLazyVim
			m.PrevScreen = ScreenMainMenu
			m.Cursor = 0
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
		// If we skipped font selection (terminal = none), go back to terminal
		if m.Choices.Terminal == "none" {
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
		if strings.Contains(selected, "mac") {
			m.Choices.OS = "mac"
		} else {
			m.Choices.OS = "linux"
		}
		m.Screen = ScreenTerminalSelect
		m.Cursor = 0

	case ScreenTerminalSelect:
		term := strings.ToLower(strings.Split(options[m.Cursor], " ")[0])
		m.Choices.Terminal = term
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

	return func() tea.Msg {
		// Execute the step
		err := executeStep(step.ID, &m)
		return stepCompleteMsg{stepID: step.ID, err: err}
	}
}
