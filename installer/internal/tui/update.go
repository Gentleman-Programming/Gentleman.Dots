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

	case ScreenOSSelect, ScreenTerminalSelect, ScreenFontSelect, ScreenShellSelect, ScreenWMSelect, ScreenNvimSelect, ScreenAIFrameworkConfirm, ScreenAIFrameworkPreset, ScreenGhosttyWarning:
		return m.handleSelectionKeys(key)

	case ScreenAIToolsSelect:
		return m.handleAIToolsKeys(key)

	case ScreenAIFrameworkCategories:
		return m.handleAICategoriesKeys(key)

	case ScreenAIFrameworkCategoryItems:
		return m.handleAICategoryItemsKeys(key)

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
	case ScreenOSSelect, ScreenTerminalSelect, ScreenFontSelect, ScreenShellSelect, ScreenWMSelect, ScreenNvimSelect, ScreenAIToolsSelect, ScreenAIFrameworkConfirm, ScreenAIFrameworkPreset, ScreenAIFrameworkCategories, ScreenAIFrameworkCategoryItems:
		return m.goBackInstallStep()
	case ScreenGhosttyWarning:
		// Go back to terminal selection
		m.Screen = ScreenTerminalSelect
		m.Cursor = 0
	case ScreenBackupConfirm:
		// Go back to last AI screen in the wizard flow
		if len(m.Choices.AITools) > 0 && m.Choices.InstallAIFramework && m.AICategorySelected != nil {
			m.Screen = ScreenAIFrameworkCategories
		} else if len(m.Choices.AITools) > 0 && m.Choices.InstallAIFramework {
			m.Screen = ScreenAIFrameworkPreset
		} else if len(m.Choices.AITools) > 0 {
			m.Screen = ScreenAIFrameworkConfirm
		} else {
			m.Screen = ScreenAIToolsSelect
		}
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

	case ScreenAIToolsSelect:
		m.Screen = ScreenNvimSelect
		m.Cursor = 0
		m.Choices.AITools = nil
		m.AIToolSelected = nil

	case ScreenAIFrameworkConfirm:
		m.Screen = ScreenAIToolsSelect
		m.Cursor = 0
		m.Choices.InstallAIFramework = false

	case ScreenAIFrameworkPreset:
		m.Screen = ScreenAIFrameworkConfirm
		m.Cursor = 0
		m.Choices.AIFrameworkPreset = ""

	case ScreenAIFrameworkCategories:
		m.Screen = ScreenAIFrameworkPreset
		m.Cursor = 0
		m.Choices.AIFrameworkModules = nil
		m.AICategorySelected = nil

	case ScreenAIFrameworkCategoryItems:
		// Back to categories — restore cursor to this category
		m.Screen = ScreenAIFrameworkCategories
		m.Cursor = m.SelectedModuleCategory
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
		// Proceed to AI tools selection (skip on Termux)
		if m.SystemInfo.IsTermux {
			// Termux doesn't support AI tools, skip to backup/install
			return m.proceedToBackupOrInstall()
		}
		m.Screen = ScreenAIToolsSelect
		m.Cursor = 0
		m.AIToolSelected = make([]bool, len(aiToolIDMap))

	case ScreenAIFrameworkConfirm:
		m.Choices.InstallAIFramework = m.Cursor == 0
		if m.Choices.InstallAIFramework {
			m.Screen = ScreenAIFrameworkPreset
			m.Cursor = 0
		} else {
			return m.proceedToBackupOrInstall()
		}

	case ScreenAIFrameworkPreset:
		if m.Cursor == 0 { // Custom — first option
			m.Choices.AIFrameworkPreset = ""
			// Initialize category selection map
			m.AICategorySelected = make(map[string][]bool)
			for _, cat := range moduleCategories {
				m.AICategorySelected[cat.ID] = make([]bool, len(cat.Items))
			}
			m.Screen = ScreenAIFrameworkCategories
			m.Cursor = 0
		} else if m.Cursor >= 2 && m.Cursor <= 7 {
			// Presets at indices 2-7 (after separator at 1)
			presets := []string{"minimal", "frontend", "backend", "fullstack", "data", "complete"}
			presetIdx := m.Cursor - 2
			if presetIdx < len(presets) {
				m.Choices.AIFrameworkPreset = presets[presetIdx]
				m.Choices.AIFrameworkModules = nil
				return m.proceedToBackupOrInstall()
			}
		}
	}

	return m, nil
}

// proceedToBackupOrInstall handles the transition from the last wizard screen to installation
func (m Model) proceedToBackupOrInstall() (tea.Model, tea.Cmd) {
	m.ExistingConfigs = system.DetectExistingConfigs()
	if len(m.ExistingConfigs) > 0 {
		m.Screen = ScreenBackupConfirm
		m.Cursor = 0
	} else {
		m.SetupInstallSteps()
		m.Screen = ScreenInstalling
		m.CurrentStep = 0
		return m, func() tea.Msg { return installStartMsg{} }
	}
	return m, nil
}

// aiToolIDMap maps AI tool option index to tool ID
var aiToolIDMap = []string{"claude", "opencode", "gemini", "copilot"}

// ModuleCategory groups related module items for the category drill-down UI
type ModuleCategory struct {
	ID       string       // Category identifier (e.g. "scripts")
	Label    string       // Display name
	Icon     string       // Emoji icon
	Items    []ModuleItem // Individual selectable items
	IsAtomic bool         // If true, selecting ANY sub-item sends the parent ID to the framework script
}

// ModuleItem represents a single selectable module within a category
type ModuleItem struct {
	ID    string // Module identifier sent to --modules flag
	Label string // Display label in the TUI
}

// moduleCategories is the data-driven registry of all AI framework module categories.
// Items mirror the real project-starter-framework repository structure.
// setup-global.sh installs features at the category level (--features=hooks,skills,...).
var moduleCategories = []ModuleCategory{
	{
		ID: "hooks", Label: "Hooks", Icon: "🪝",
		Items: []ModuleItem{
			{ID: "block-dangerous-commands", Label: "Block Dangerous Commands"},
			{ID: "commit-guard", Label: "Commit Guard"},
			{ID: "context-loader", Label: "Context Loader"},
			{ID: "improve-prompt", Label: "Improve Prompt"},
			{ID: "learning-log", Label: "Learning Log"},
			{ID: "model-router", Label: "Model Router"},
			{ID: "secret-scanner", Label: "Secret Scanner"},
			{ID: "skill-validator", Label: "Skill Validator"},
			{ID: "task-artifact", Label: "Task Artifact"},
			{ID: "validate-workflow", Label: "Validate Workflow"},
		},
	},
	{
		ID: "commands", Label: "Commands", Icon: "⚡",
		Items: []ModuleItem{
			// Git
			{ID: "git:changelog", Label: "Git: Changelog"},
			{ID: "git:ci-local", Label: "Git: CI Local"},
			{ID: "git:commit", Label: "Git: Commit"},
			{ID: "git:fix-issue", Label: "Git: Fix Issue"},
			{ID: "git:pr-create", Label: "Git: PR Create"},
			{ID: "git:pr-review", Label: "Git: PR Review"},
			{ID: "git:worktree", Label: "Git: Worktree"},
			// Refactoring
			{ID: "refactoring:cleanup", Label: "Refactoring: Cleanup"},
			{ID: "refactoring:dead-code", Label: "Refactoring: Dead Code"},
			{ID: "refactoring:extract", Label: "Refactoring: Extract"},
			// Testing
			{ID: "testing:e2e", Label: "Testing: E2E"},
			{ID: "testing:tdd", Label: "Testing: TDD"},
			{ID: "testing:test-coverage", Label: "Testing: Coverage"},
			{ID: "testing:test-fix", Label: "Testing: Fix Tests"},
			// Workflow
			{ID: "workflow:generate-agents-md", Label: "Workflow: Generate Agents"},
			{ID: "workflow:planning", Label: "Workflow: Planning"},
			{ID: "workflows:compound", Label: "Workflows: Compound"},
			{ID: "workflows:plan", Label: "Workflows: Plan"},
			{ID: "workflows:review", Label: "Workflows: Review"},
			{ID: "workflows:work", Label: "Workflows: Work"},
		},
	},
	{
		ID: "agents", Label: "Agents", Icon: "🤖",
		Items: []ModuleItem{
			// General
			{ID: "orchestrator", Label: "General: Orchestrator"},
			// Business
			{ID: "business-api-designer", Label: "Business: API Designer"},
			{ID: "business-business-analyst", Label: "Business: Business Analyst"},
			{ID: "business-product-strategist", Label: "Business: Product Strategist"},
			{ID: "business-project-manager", Label: "Business: Project Manager"},
			{ID: "business-requirements-analyst", Label: "Business: Requirements Analyst"},
			{ID: "business-technical-writer", Label: "Business: Technical Writer"},
			// Creative
			{ID: "creative-ux-designer", Label: "Creative: UX Designer"},
			// Data & AI
			{ID: "data-ai-ai-engineer", Label: "Data & AI: AI Engineer"},
			{ID: "data-ai-analytics-engineer", Label: "Data & AI: Analytics Engineer"},
			{ID: "data-ai-data-engineer", Label: "Data & AI: Data Engineer"},
			{ID: "data-ai-data-scientist", Label: "Data & AI: Data Scientist"},
			{ID: "data-ai-mlops-engineer", Label: "Data & AI: MLOps Engineer"},
			{ID: "data-ai-prompt-engineer", Label: "Data & AI: Prompt Engineer"},
			// Development
			{ID: "development-angular-expert", Label: "Development: Angular Expert"},
			{ID: "development-backend-architect", Label: "Development: Backend Architect"},
			{ID: "development-database-specialist", Label: "Development: Database Specialist"},
			{ID: "development-frontend-specialist", Label: "Development: Frontend Specialist"},
			{ID: "development-fullstack-engineer", Label: "Development: Fullstack Engineer"},
			{ID: "development-golang-pro", Label: "Development: Go Pro"},
			{ID: "development-java-enterprise", Label: "Development: Java Enterprise"},
			{ID: "development-javascript-pro", Label: "Development: JavaScript Pro"},
			{ID: "development-nextjs-pro", Label: "Development: Next.js Pro"},
			{ID: "development-python-pro", Label: "Development: Python Pro"},
			{ID: "development-react-pro", Label: "Development: React Pro"},
			{ID: "development-rust-pro", Label: "Development: Rust Pro"},
			{ID: "development-spring-boot-4-expert", Label: "Development: Spring Boot 4"},
			{ID: "development-typescript-pro", Label: "Development: TypeScript Pro"},
			{ID: "development-vue-specialist", Label: "Development: Vue Specialist"},
			// Infrastructure
			{ID: "infrastructure-cloud-architect", Label: "Infrastructure: Cloud Architect"},
			{ID: "infrastructure-deployment-manager", Label: "Infrastructure: Deployment Manager"},
			{ID: "infrastructure-devops-engineer", Label: "Infrastructure: DevOps Engineer"},
			{ID: "infrastructure-incident-responder", Label: "Infrastructure: Incident Responder"},
			{ID: "infrastructure-kubernetes-expert", Label: "Infrastructure: Kubernetes Expert"},
			{ID: "infrastructure-monitoring-specialist", Label: "Infrastructure: Monitoring Specialist"},
			{ID: "infrastructure-performance-engineer", Label: "Infrastructure: Performance Engineer"},
			// Quality
			{ID: "quality-accessibility-auditor", Label: "Quality: Accessibility Auditor"},
			{ID: "quality-code-reviewer-compact", Label: "Quality: Code Reviewer (Compact)"},
			{ID: "quality-code-reviewer", Label: "Quality: Code Reviewer"},
			{ID: "quality-dependency-manager", Label: "Quality: Dependency Manager"},
			{ID: "quality-e2e-test-specialist", Label: "Quality: E2E Test Specialist"},
			{ID: "quality-performance-tester", Label: "Quality: Performance Tester"},
			{ID: "quality-security-auditor", Label: "Quality: Security Auditor"},
			{ID: "quality-test-engineer", Label: "Quality: Test Engineer"},
			// Specialists
			{ID: "specialists-api-designer", Label: "Specialists: API Designer"},
			{ID: "specialists-backend-architect", Label: "Specialists: Backend Architect"},
			{ID: "specialists-code-reviewer", Label: "Specialists: Code Reviewer"},
			{ID: "specialists-db-optimizer", Label: "Specialists: DB Optimizer"},
			{ID: "specialists-devops-engineer", Label: "Specialists: DevOps Engineer"},
			{ID: "specialists-documentation-writer", Label: "Specialists: Documentation Writer"},
			{ID: "specialists-frontend-developer", Label: "Specialists: Frontend Developer"},
			{ID: "specialists-performance-analyst", Label: "Specialists: Performance Analyst"},
			{ID: "specialists-refactor-specialist", Label: "Specialists: Refactor Specialist"},
			{ID: "specialists-security-auditor", Label: "Specialists: Security Auditor"},
			{ID: "specialists-test-engineer", Label: "Specialists: Test Engineer"},
			{ID: "specialists-ux-consultant", Label: "Specialists: UX Consultant"},
			// Specialized
			{ID: "specialized-agent-generator", Label: "Specialized: Agent Generator"},
			{ID: "specialized-blockchain-developer", Label: "Specialized: Blockchain Developer"},
			{ID: "specialized-code-migrator", Label: "Specialized: Code Migrator"},
			{ID: "specialized-context-manager", Label: "Specialized: Context Manager"},
			{ID: "specialized-documentation-writer", Label: "Specialized: Documentation Writer"},
			{ID: "specialized-ecommerce-expert", Label: "Specialized: E-Commerce Expert"},
			{ID: "specialized-embedded-engineer", Label: "Specialized: Embedded Engineer"},
			{ID: "specialized-error-detective", Label: "Specialized: Error Detective"},
			{ID: "specialized-fintech-specialist", Label: "Specialized: Fintech Specialist"},
			{ID: "specialized-freelance-planner", Label: "Specialized: Freelance Planner"},
			{ID: "specialized-freelance-planner-v2", Label: "Specialized: Freelance Planner v2"},
			{ID: "specialized-freelance-planner-v3", Label: "Specialized: Freelance Planner v3"},
			{ID: "specialized-freelance-planner-v4", Label: "Specialized: Freelance Planner v4"},
			{ID: "specialized-game-developer", Label: "Specialized: Game Developer"},
			{ID: "specialized-healthcare-dev", Label: "Specialized: Healthcare Dev"},
			{ID: "specialized-mobile-developer", Label: "Specialized: Mobile Developer"},
			{ID: "specialized-parallel-plan-executor", Label: "Specialized: Parallel Plan Executor"},
			{ID: "specialized-plan-executor", Label: "Specialized: Plan Executor"},
			{ID: "specialized-solo-dev-planner", Label: "Specialized: Solo Dev Planner"},
			{ID: "specialized-template-writer", Label: "Specialized: Template Writer"},
			{ID: "specialized-test-runner", Label: "Specialized: Test Runner"},
			{ID: "specialized-vibekanban-worker", Label: "Specialized: VibeKanban Worker"},
			{ID: "specialized-wave-executor", Label: "Specialized: Wave Executor"},
			{ID: "specialized-workflow-optimizer", Label: "Specialized: Workflow Optimizer"},
		},
	},
	{
		ID: "skills", Label: "Skills", Icon: "🎯",
		Items: []ModuleItem{
			// Backend (21)
			{ID: "backend-api-gateway", Label: "Backend: API Gateway"},
			{ID: "backend-bff-concepts", Label: "Backend: BFF Concepts"},
			{ID: "backend-bff-spring", Label: "Backend: BFF Spring"},
			{ID: "backend-chi-router", Label: "Backend: Chi Router"},
			{ID: "backend-error-handling", Label: "Backend: Error Handling"},
			{ID: "backend-exceptions-spring", Label: "Backend: Exceptions Spring"},
			{ID: "backend-fastapi", Label: "Backend: FastAPI"},
			{ID: "backend-gateway-spring", Label: "Backend: Gateway Spring"},
			{ID: "backend-go-backend", Label: "Backend: Go Backend"},
			{ID: "backend-gradle-multimodule", Label: "Backend: Gradle Multi-Module"},
			{ID: "backend-graphql-concepts", Label: "Backend: GraphQL Concepts"},
			{ID: "backend-graphql-spring", Label: "Backend: GraphQL Spring"},
			{ID: "backend-grpc-concepts", Label: "Backend: gRPC Concepts"},
			{ID: "backend-grpc-spring", Label: "Backend: gRPC Spring"},
			{ID: "backend-jwt-auth", Label: "Backend: JWT Auth"},
			{ID: "backend-notifications-concepts", Label: "Backend: Notifications"},
			{ID: "backend-recommendations-concepts", Label: "Backend: Recommendations"},
			{ID: "backend-search-concepts", Label: "Backend: Search Concepts"},
			{ID: "backend-search-spring", Label: "Backend: Search Spring"},
			{ID: "backend-spring-boot-4", Label: "Backend: Spring Boot 4"},
			{ID: "backend-websockets", Label: "Backend: WebSockets"},
			// Data & AI (11)
			{ID: "data-ai-ai-ml", Label: "Data & AI: AI/ML"},
			{ID: "data-ai-analytics-concepts", Label: "Data & AI: Analytics Concepts"},
			{ID: "data-ai-analytics-spring", Label: "Data & AI: Analytics Spring"},
			{ID: "data-ai-duckdb-analytics", Label: "Data & AI: DuckDB Analytics"},
			{ID: "data-ai-langchain", Label: "Data & AI: LangChain"},
			{ID: "data-ai-mlflow", Label: "Data & AI: MLflow"},
			{ID: "data-ai-onnx-inference", Label: "Data & AI: ONNX Inference"},
			{ID: "data-ai-powerbi", Label: "Data & AI: Power BI"},
			{ID: "data-ai-pytorch", Label: "Data & AI: PyTorch"},
			{ID: "data-ai-scikit-learn", Label: "Data & AI: scikit-learn"},
			{ID: "data-ai-vector-db", Label: "Data & AI: Vector DB"},
			// Database (6)
			{ID: "database-graph-databases", Label: "Database: Graph Databases"},
			{ID: "database-graph-spring", Label: "Database: Graph Spring"},
			{ID: "database-pgx-postgres", Label: "Database: PGX Postgres"},
			{ID: "database-redis-cache", Label: "Database: Redis Cache"},
			{ID: "database-sqlite-embedded", Label: "Database: SQLite Embedded"},
			{ID: "database-timescaledb", Label: "Database: TimescaleDB"},
			// Docs (4)
			{ID: "docs-api-documentation", Label: "Docs: API Documentation"},
			{ID: "docs-docs-spring", Label: "Docs: Spring Docs"},
			{ID: "docs-mustache-templates", Label: "Docs: Mustache Templates"},
			{ID: "docs-technical-docs", Label: "Docs: Technical Docs"},
			// Frontend (7)
			{ID: "frontend-astro-ssr", Label: "Frontend: Astro SSR"},
			{ID: "frontend-frontend-design", Label: "Frontend: Design Patterns"},
			{ID: "frontend-frontend-web", Label: "Frontend: Web Development"},
			{ID: "frontend-mantine-ui", Label: "Frontend: Mantine UI"},
			{ID: "frontend-tanstack-query", Label: "Frontend: TanStack Query"},
			{ID: "frontend-zod-validation", Label: "Frontend: Zod Validation"},
			{ID: "frontend-zustand-state", Label: "Frontend: Zustand State"},
			// Infrastructure (8)
			{ID: "infra-chaos-engineering", Label: "Infrastructure: Chaos Engineering"},
			{ID: "infra-chaos-spring", Label: "Infrastructure: Chaos Spring"},
			{ID: "infra-devops-infra", Label: "Infrastructure: DevOps"},
			{ID: "infra-docker-containers", Label: "Infrastructure: Docker"},
			{ID: "infra-kubernetes", Label: "Infrastructure: Kubernetes"},
			{ID: "infra-opentelemetry", Label: "Infrastructure: OpenTelemetry"},
			{ID: "infra-traefik-proxy", Label: "Infrastructure: Traefik Proxy"},
			{ID: "infra-woodpecker-ci", Label: "Infrastructure: Woodpecker CI"},
			// Mobile (2)
			{ID: "mobile-ionic-capacitor", Label: "Mobile: Ionic Capacitor"},
			{ID: "mobile-mobile-ionic", Label: "Mobile: Mobile Ionic"},
			// Prompt & Quality (2)
			{ID: "prompt-improver", Label: "Prompt: Prompt Improver"},
			{ID: "quality-ghagga-review", Label: "Quality: Ghagga Review"},
			// References (5)
			{ID: "references-hooks-patterns", Label: "References: Hooks Patterns"},
			{ID: "references-mcp-servers", Label: "References: MCP Servers"},
			{ID: "references-plugins-reference", Label: "References: Plugins Reference"},
			{ID: "references-skills-reference", Label: "References: Skills Reference"},
			{ID: "references-subagent-templates", Label: "References: Subagent Templates"},
			// Systems & IoT (4)
			{ID: "systems-modbus-protocol", Label: "Systems: Modbus Protocol"},
			{ID: "systems-mqtt-rumqttc", Label: "Systems: MQTT rumqttc"},
			{ID: "systems-rust-systems", Label: "Systems: Rust Systems"},
			{ID: "systems-tokio-async", Label: "Systems: Tokio Async"},
			// Testing (3)
			{ID: "testing-playwright-e2e", Label: "Testing: Playwright E2E"},
			{ID: "testing-testcontainers", Label: "Testing: Testcontainers"},
			{ID: "testing-vitest-testing", Label: "Testing: Vitest Testing"},
			// Workflow (12)
			{ID: "workflow-ci-local-guide", Label: "Workflow: CI Local Guide"},
			{ID: "workflow-claude-automation", Label: "Workflow: Claude Automation"},
			{ID: "workflow-claude-md-improver", Label: "Workflow: CLAUDE.md Improver"},
			{ID: "workflow-finish-dev-branch", Label: "Workflow: Finish Dev Branch"},
			{ID: "workflow-git-github", Label: "Workflow: Git & GitHub"},
			{ID: "workflow-git-workflow", Label: "Workflow: Git Workflow"},
			{ID: "workflow-ide-plugins", Label: "Workflow: IDE Plugins"},
			{ID: "workflow-ide-plugins-intellij", Label: "Workflow: IDE Plugins IntelliJ"},
			{ID: "workflow-obsidian-brain", Label: "Workflow: Obsidian Brain"},
			{ID: "workflow-git-worktrees", Label: "Workflow: Git Worktrees"},
			{ID: "workflow-verification", Label: "Workflow: Verification"},
			{ID: "workflow-wave-workflow", Label: "Workflow: Wave Workflow"},
		},
	},
	{
		ID: "sdd", Label: "SDD (Spec-Driven Development)", Icon: "📐", IsAtomic: true,
		Items: []ModuleItem{
			{ID: "sdd-init", Label: "Init"},
			{ID: "sdd-explore", Label: "Explore"},
			{ID: "sdd-propose", Label: "Propose"},
			{ID: "sdd-spec", Label: "Spec"},
			{ID: "sdd-design", Label: "Design"},
			{ID: "sdd-tasks", Label: "Tasks"},
			{ID: "sdd-apply", Label: "Apply"},
			{ID: "sdd-verify", Label: "Verify"},
			{ID: "sdd-archive", Label: "Archive"},
		},
	},
	{
		ID: "mcp", Label: "MCP Servers", Icon: "🔌", IsAtomic: true,
		Items: []ModuleItem{
			{ID: "mcp-context7", Label: "Context7"},
			{ID: "mcp-engram", Label: "Engram"},
			{ID: "mcp-jira", Label: "Jira"},
			{ID: "mcp-atlassian", Label: "Atlassian"},
			{ID: "mcp-figma", Label: "Figma"},
			{ID: "mcp-notion", Label: "Notion"},
		},
	},
}

// collectSelectedFeatures converts the category selection map into feature flags for setup-global.sh.
// If ANY item within a category is selected, the category's feature flag is included.
// setup-global.sh operates at the feature level: --features=hooks,skills,agents,sdd,mcp
func collectSelectedFeatures(sel map[string][]bool) []string {
	var features []string
	for _, cat := range moduleCategories {
		bools, ok := sel[cat.ID]
		if !ok {
			continue
		}
		for _, b := range bools {
			if b {
				features = append(features, cat.ID)
				break
			}
		}
	}
	return features
}

func (m Model) handleAIToolsKeys(key string) (tea.Model, tea.Cmd) {
	options := m.GetCurrentOptions()
	lastToolIdx := len(aiToolIDMap) - 1 // Last toggleable tool index
	confirmIdx := len(options) - 1      // "Confirm selection" is last option

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
			if strings.HasPrefix(options[m.Cursor], "───") && m.Cursor < len(options)-1 {
				m.Cursor++
			}
		}
	case "enter", " ":
		if m.Cursor <= lastToolIdx {
			// Toggle tool selection
			if m.AIToolSelected != nil && m.Cursor < len(m.AIToolSelected) {
				m.AIToolSelected[m.Cursor] = !m.AIToolSelected[m.Cursor]
			}
		} else if m.Cursor == confirmIdx {
			// Confirm — collect selected tools
			var selected []string
			for i, sel := range m.AIToolSelected {
				if sel && i < len(aiToolIDMap) {
					selected = append(selected, aiToolIDMap[i])
				}
			}
			m.Choices.AITools = selected
			// If any AI tools selected, ask about framework
			if len(m.Choices.AITools) > 0 {
				m.Screen = ScreenAIFrameworkConfirm
				m.Cursor = 0
			} else {
				// No AI tools, skip framework too
				m.Choices.InstallAIFramework = false
				return m.proceedToBackupOrInstall()
			}
		}
	case "esc", "backspace":
		return m.goBackInstallStep()
	}

	return m, nil
}

func (m Model) handleAICategoriesKeys(key string) (tea.Model, tea.Cmd) {
	options := m.GetCurrentOptions()
	lastCategoryIdx := len(moduleCategories) - 1
	confirmIdx := len(options) - 1

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
		if m.Cursor <= lastCategoryIdx {
			// Drill into category
			m.SelectedModuleCategory = m.Cursor
			m.Screen = ScreenAIFrameworkCategoryItems
			m.Cursor = 0
		} else if m.Cursor == confirmIdx {
			// Confirm — collect selected features for setup-global.sh
			m.Choices.AIFrameworkModules = collectSelectedFeatures(m.AICategorySelected)
			if len(m.Choices.AIFrameworkModules) == 0 {
				m.Choices.InstallAIFramework = false
			}
			return m.proceedToBackupOrInstall()
		}
	case "esc", "backspace":
		return m.goBackInstallStep()
	}

	return m, nil
}

func (m Model) handleAICategoryItemsKeys(key string) (tea.Model, tea.Cmd) {
	if m.SelectedModuleCategory < 0 || m.SelectedModuleCategory >= len(moduleCategories) {
		return m, nil
	}
	cat := moduleCategories[m.SelectedModuleCategory]
	options := m.GetCurrentOptions()
	lastItemIdx := len(cat.Items) - 1
	backIdx := len(options) - 1

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
		if m.Cursor <= lastItemIdx {
			// Toggle item within category
			bools := m.AICategorySelected[cat.ID]
			if bools != nil && m.Cursor < len(bools) {
				bools[m.Cursor] = !bools[m.Cursor]
				m.AICategorySelected[cat.ID] = bools
			}
		} else if m.Cursor == backIdx {
			// Back to categories — restore cursor to this category
			m.Screen = ScreenAIFrameworkCategories
			m.Cursor = m.SelectedModuleCategory
		}
	case "esc", "backspace":
		// Back to categories — restore cursor to this category
		m.Screen = ScreenAIFrameworkCategories
		m.Cursor = m.SelectedModuleCategory
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
		// Go back to the last AI screen in the wizard flow
		if len(m.Choices.AITools) > 0 && m.Choices.InstallAIFramework && m.AICategorySelected != nil {
			// Was in custom mode — go back to categories
			m.Screen = ScreenAIFrameworkCategories
		} else if len(m.Choices.AITools) > 0 && m.Choices.InstallAIFramework {
			m.Screen = ScreenAIFrameworkPreset
		} else if len(m.Choices.AITools) > 0 {
			m.Screen = ScreenAIFrameworkConfirm
		} else {
			m.Screen = ScreenAIToolsSelect
		}
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
