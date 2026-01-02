package tui

import (
	"fmt"
	"strings"

	"github.com/Gentleman-Programming/Gentleman.Dots/installer/internal/tui/trainer"
	"github.com/charmbracelet/lipgloss"
)

const logo = `
                    ░░░░░░      ░░░░░░                        
                  ░░░░░░░░░░  ░░░░░░░░░░                      
                ░░░░░░░░░░░░░░░░░░░░░░░░░░                    
              ░░░░░░░░░░▒▒▒▒░░▒▒▒▒░░░░░░░░░░                  
  ░░░░      ░░░░░░▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒░░░░░░        ░░░░    
▒▒░░      ░░░░░░▒▒▒▒▒▒▒▒▒▒██▒▒██▒▒▒▒▒▒▒▒▒▒░░░░░░        ▒▒░░  
▒▒░░    ░░░░░░░░▒▒▒▒▒▒▒▒▒▒████▒▒████▒▒▒▒▒▒▒▒▒▒░░░░░░░░  ▒▒░░  
▒▒▒▒░░░░░░▒▒▒▒▒▒▒▒▒▒▒▒▒▒██████▒▒██████▒▒▒▒▒▒▒▒▒▒▒▒▒▒░░░░░░▒▒▒▒
██▒▒▒▒▒▒▒▒▒▒▒▒▒▒██▒▒▒▒██████▓▓██▒▒██████▒▒▓▓██▒▒▒▒▒▒▒▒▒▒▒▒▒▒██
████▒▒▒▒▒▒████▒▒▒▒██████████  ██████████▒▒▒▒████▒▒▒▒▒▒▒▒██    
  ████████████████████████      ████████████████████████      
    ██████████████████              ██████████████████        
        ██████████                      ██████████            
`

const gentlemanText = `
 ██████╗ ███████╗███╗   ██╗████████╗██╗     ███████╗███╗   ███╗ █████╗ ███╗   ██╗
██╔════╝ ██╔════╝████╗  ██║╚══██╔══╝██║     ██╔════╝████╗ ████║██╔══██╗████╗  ██║
██║  ███╗█████╗  ██╔██╗ ██║   ██║   ██║     █████╗  ██╔████╔██║███████║██╔██╗ ██║
██║   ██║██╔══╝  ██║╚██╗██║   ██║   ██║     ██╔══╝  ██║╚██╔╝██║██╔══██║██║╚██╗██║
╚██████╔╝███████╗██║ ╚████║   ██║   ███████╗███████╗██║ ╚═╝ ██║██║  ██║██║ ╚████║
 ╚═════╝ ╚══════╝╚═╝  ╚═══╝   ╚═╝   ╚══════╝╚══════╝╚═╝     ╚═╝╚═╝  ╚═╝╚═╝  ╚═══╝
                        ██████╗  ██████╗ ████████╗███████╗
                        ██╔══██╗██╔═══██╗╚══██╔══╝██╔════╝
                        ██║  ██║██║   ██║   ██║   ███████╗
                        ██║  ██║██║   ██║   ██║   ╚════██║
                        ██████╔╝╚██████╔╝   ██║   ███████║
                        ╚═════╝  ╚═════╝    ╚═╝   ╚══════╝
`

// View implements tea.Model
func (m Model) View() string {
	if m.Quitting {
		return ""
	}

	var s strings.Builder

	switch m.Screen {
	case ScreenWelcome:
		s.WriteString(m.renderWelcome())
	case ScreenMainMenu:
		s.WriteString(m.renderMainMenu())
	case ScreenOSSelect, ScreenTerminalSelect, ScreenFontSelect, ScreenShellSelect, ScreenWMSelect, ScreenNvimSelect, ScreenGhosttyWarning:
		s.WriteString(m.renderSelection())
	case ScreenLearnTerminals:
		s.WriteString(m.renderLearnTerminals())
	case ScreenLearnShells:
		s.WriteString(m.renderLearnShells())
	case ScreenLearnWM:
		s.WriteString(m.renderLearnWM())
	case ScreenLearnNvim:
		s.WriteString(m.renderLearnNvim())
	case ScreenKeymaps:
		s.WriteString(m.renderKeymapsMenu())
	case ScreenKeymapCategory:
		s.WriteString(m.renderKeymapCategory())
	case ScreenKeymapsMenu:
		s.WriteString(m.renderToolKeymapsMenu())
	case ScreenKeymapsTmux:
		s.WriteString(m.renderTmuxKeymapsMenu())
	case ScreenKeymapsTmuxCat:
		s.WriteString(m.renderTmuxKeymapCategory())
	case ScreenKeymapsZellij:
		s.WriteString(m.renderZellijKeymapsMenu())
	case ScreenKeymapsZellijCat:
		s.WriteString(m.renderZellijKeymapCategory())
	case ScreenKeymapsGhostty:
		s.WriteString(m.renderGhosttyKeymapsMenu())
	case ScreenKeymapsGhosttyCat:
		s.WriteString(m.renderGhosttyKeymapCategory())
	case ScreenLearnLazyVim:
		s.WriteString(m.renderLazyVimMenu())
	case ScreenLazyVimTopic:
		s.WriteString(m.renderLazyVimTopic())
	case ScreenBackupConfirm:
		s.WriteString(m.renderBackupConfirm())
	case ScreenRestoreBackup:
		s.WriteString(m.renderRestoreBackup())
	case ScreenRestoreConfirm:
		s.WriteString(m.renderRestoreConfirm())
	case ScreenInstalling:
		s.WriteString(m.renderInstalling())
	case ScreenComplete:
		s.WriteString(m.renderComplete())
	case ScreenError:
		s.WriteString(m.renderError())
	// Trainer screens
	case ScreenTrainerMenu:
		s.WriteString(m.renderTrainerMenu())
	case ScreenTrainerLesson:
		s.WriteString(m.renderTrainerExercise("Lesson"))
	case ScreenTrainerPractice:
		s.WriteString(m.renderTrainerExercise("Practice"))
	case ScreenTrainerBoss:
		s.WriteString(m.renderTrainerBoss())
	case ScreenTrainerResult:
		s.WriteString(m.renderTrainerResult())
	case ScreenTrainerBossResult:
		s.WriteString(m.renderTrainerBossResult())
	}

	// Leader mode indicator
	if m.LeaderMode {
		s.WriteString("\n")
		s.WriteString(WarningStyle.Render("▶ LEADER MODE - Press: q=quit, d=details"))
	}

	// Apply global padding (top: 1, right: 2, bottom: 0, left: 2)
	paddedStyle := lipgloss.NewStyle().Padding(1, 2, 0, 2)
	return paddedStyle.Render(s.String())
}

func (m Model) renderWelcome() string {
	var s strings.Builder

	// Logo
	s.WriteString(LogoStyle.Render(logo))
	s.WriteString("\n")
	s.WriteString(TitleStyle.Render(gentlemanText))
	s.WriteString("\n\n")

	// System info
	info := fmt.Sprintf("Detected: %s", m.SystemInfo.OSName)
	if m.SystemInfo.IsWSL {
		info += " (WSL)"
	}
	if m.SystemInfo.HasBrew {
		info += " | Homebrew ✓"
	}
	s.WriteString(InfoStyle.Render(info))
	s.WriteString("\n\n")

	// Instructions
	s.WriteString(SubtitleStyle.Render("Your terminal environment, configured in minutes."))
	s.WriteString("\n\n")
	s.WriteString(HelpStyle.Render("Press [Enter] to start • [Space q] to quit"))

	// Center both horizontally and vertically
	return CenterBoth(s.String(), m.Width, m.Height)
}

func (m Model) renderMainMenu() string {
	var s strings.Builder

	// Title
	s.WriteString(TitleStyle.Render("🎩 Gentleman.Dots"))
	s.WriteString("\n")
	s.WriteString(MutedStyle.Render("What would you like to do?"))
	s.WriteString("\n\n")

	// Options
	options := m.GetCurrentOptions()
	for i, opt := range options {
		cursor := "  "
		style := UnselectedStyle
		if i == m.Cursor {
			cursor = "▸ "
			style = SelectedStyle
		}
		s.WriteString(style.Render(cursor + opt))
		s.WriteString("\n")
	}

	s.WriteString("\n")
	s.WriteString(HelpStyle.Render("↑/k up • ↓/j down • [Enter] select • [Space q] quit"))

	return s.String()
}

func (m Model) renderSelection() string {
	var s strings.Builder

	// Progress indicator
	s.WriteString(m.renderStepProgress())
	s.WriteString("\n\n")

	// Title
	s.WriteString(TitleStyle.Render(m.GetScreenTitle()))
	s.WriteString("\n")
	s.WriteString(MutedStyle.Render(m.GetScreenDescription()))
	s.WriteString("\n\n")

	// Options
	options := m.GetCurrentOptions()
	for i, opt := range options {
		// Separator line
		if strings.HasPrefix(opt, "───") {
			s.WriteString(MutedStyle.Render(opt))
			s.WriteString("\n")
			continue
		}

		cursor := "  "
		style := UnselectedStyle
		if i == m.Cursor {
			cursor = "▸ "
			style = SelectedStyle
		}
		s.WriteString(style.Render(cursor + opt))
		s.WriteString("\n")
	}

	s.WriteString("\n")
	s.WriteString(HelpStyle.Render("↑/k up • ↓/j down • [Enter] select • [Esc] back"))

	return s.String()
}

func (m Model) renderStepProgress() string {
	steps := []string{"OS", "Terminal", "Font", "Shell", "WM", "Nvim"}
	currentIdx := 0

	switch m.Screen {
	case ScreenOSSelect:
		currentIdx = 0
	case ScreenTerminalSelect:
		currentIdx = 1
	case ScreenFontSelect:
		currentIdx = 2
	case ScreenShellSelect:
		currentIdx = 3
	case ScreenWMSelect:
		currentIdx = 4
	case ScreenNvimSelect:
		currentIdx = 5
	}

	var parts []string
	for i, step := range steps {
		var style lipgloss.Style
		if i < currentIdx {
			style = StepDoneStyle
			parts = append(parts, style.Render("✓ "+step))
		} else if i == currentIdx {
			style = StepActiveStyle
			parts = append(parts, style.Render("● "+step))
		} else {
			style = StepPendingStyle
			parts = append(parts, style.Render("○ "+step))
		}
	}

	return strings.Join(parts, MutedStyle.Render(" → "))
}

func (m Model) renderLearnTerminals() string {
	var s strings.Builder

	s.WriteString(TitleStyle.Render(m.GetScreenTitle()))
	s.WriteString("\n")
	s.WriteString(MutedStyle.Render("Select a terminal to learn more about it"))
	s.WriteString("\n\n")

	// If viewing a specific tool, show its info
	if m.ViewingTool != "" {
		return m.renderToolInfo(GetTerminalInfo(), m.ViewingTool, "terminal")
	}

	// Menu
	options := m.GetCurrentOptions()
	for i, opt := range options {
		if strings.HasPrefix(opt, "───") {
			s.WriteString(MutedStyle.Render(opt))
			s.WriteString("\n")
			continue
		}

		cursor := "  "
		style := UnselectedStyle
		if i == m.Cursor {
			cursor = "▸ "
			style = SelectedStyle
		}
		s.WriteString(style.Render(cursor + opt))
		s.WriteString("\n")
	}

	s.WriteString("\n")
	s.WriteString(HelpStyle.Render("↑/k up • ↓/j down • [Enter] select • [Esc] back"))

	return s.String()
}

func (m Model) renderLearnShells() string {
	var s strings.Builder

	s.WriteString(TitleStyle.Render(m.GetScreenTitle()))
	s.WriteString("\n")
	s.WriteString(MutedStyle.Render("Select a shell to learn more about it"))
	s.WriteString("\n\n")

	// If viewing a specific tool, show its info
	if m.ViewingTool != "" {
		return m.renderToolInfo(GetShellInfo(), m.ViewingTool, "shell")
	}

	// Menu
	options := m.GetCurrentOptions()
	for i, opt := range options {
		if strings.HasPrefix(opt, "───") {
			s.WriteString(MutedStyle.Render(opt))
			s.WriteString("\n")
			continue
		}

		cursor := "  "
		style := UnselectedStyle
		if i == m.Cursor {
			cursor = "▸ "
			style = SelectedStyle
		}
		s.WriteString(style.Render(cursor + opt))
		s.WriteString("\n")
	}

	s.WriteString("\n")
	s.WriteString(HelpStyle.Render("↑/k up • ↓/j down • [Enter] select • [Esc] back"))

	return s.String()
}

func (m Model) renderLearnWM() string {
	var s strings.Builder

	s.WriteString(TitleStyle.Render(m.GetScreenTitle()))
	s.WriteString("\n")
	s.WriteString(MutedStyle.Render("Select a window manager to learn more about it"))
	s.WriteString("\n\n")

	// If viewing a specific tool, show its info
	if m.ViewingTool != "" {
		return m.renderToolInfo(GetWMInfo(), m.ViewingTool, "wm")
	}

	// Menu
	options := m.GetCurrentOptions()
	for i, opt := range options {
		if strings.HasPrefix(opt, "───") {
			s.WriteString(MutedStyle.Render(opt))
			s.WriteString("\n")
			continue
		}

		cursor := "  "
		style := UnselectedStyle
		if i == m.Cursor {
			cursor = "▸ "
			style = SelectedStyle
		}
		s.WriteString(style.Render(cursor + opt))
		s.WriteString("\n")
	}

	s.WriteString("\n")
	s.WriteString(HelpStyle.Render("↑/k up • ↓/j down • [Enter] select • [Esc] back"))

	return s.String()
}

func (m Model) renderLearnNvim() string {
	var s strings.Builder

	s.WriteString(TitleStyle.Render(m.GetScreenTitle()))
	s.WriteString("\n")
	s.WriteString(MutedStyle.Render("Explore Neovim features and keybindings"))
	s.WriteString("\n\n")

	// If viewing features, show Nvim info
	if m.ViewingTool == "features" {
		info := GetNvimInfo()
		return m.renderSingleToolInfo(info)
	}

	// Menu
	options := m.GetCurrentOptions()
	for i, opt := range options {
		if strings.HasPrefix(opt, "───") {
			s.WriteString(MutedStyle.Render(opt))
			s.WriteString("\n")
			continue
		}

		cursor := "  "
		style := UnselectedStyle
		if i == m.Cursor {
			cursor = "▸ "
			style = SelectedStyle
		}
		s.WriteString(style.Render(cursor + opt))
		s.WriteString("\n")
	}

	s.WriteString("\n")
	s.WriteString(HelpStyle.Render("↑/k up • ↓/j down • [Enter] select • [Esc] back"))

	return s.String()
}

func (m Model) renderToolInfo(tools map[string]ToolInfo, toolKey string, category string) string {
	var s strings.Builder

	info, exists := tools[toolKey]
	if !exists {
		s.WriteString(ErrorStyle.Render("Tool not found"))
		return s.String()
	}

	return m.renderSingleToolInfo(info)
}

func (m Model) renderSingleToolInfo(info ToolInfo) string {
	var s strings.Builder

	// Tool name and description
	s.WriteString(TitleStyle.Render(info.Name))
	s.WriteString("\n")
	s.WriteString(MutedStyle.Render(info.Description))
	s.WriteString("\n")
	s.WriteString(MutedStyle.Render(info.Website))
	s.WriteString("\n\n")

	// Pros
	s.WriteString(SuccessStyle.Render("✓ Pros"))
	s.WriteString("\n")
	for _, pro := range info.Pros {
		s.WriteString(InfoStyle.Render("  • " + pro))
		s.WriteString("\n")
	}

	s.WriteString("\n")

	// Cons
	s.WriteString(WarningStyle.Render("✗ Cons"))
	s.WriteString("\n")
	for _, con := range info.Cons {
		s.WriteString(MutedStyle.Render("  • " + con))
		s.WriteString("\n")
	}

	s.WriteString("\n")
	s.WriteString(HelpStyle.Render("↑/k up • ↓/j down • [Enter] select • [Esc] back • [Space q] quit"))

	return s.String()
}

func (m Model) renderKeymapsMenu() string {
	var s strings.Builder

	s.WriteString(TitleStyle.Render(m.GetScreenTitle()))
	s.WriteString("\n")
	s.WriteString(MutedStyle.Render("Select a category to view keybindings"))
	s.WriteString("\n\n")

	// Menu
	options := m.GetCurrentOptions()
	for i, opt := range options {
		if strings.HasPrefix(opt, "───") {
			s.WriteString(MutedStyle.Render(opt))
			s.WriteString("\n")
			continue
		}

		cursor := "  "
		style := UnselectedStyle
		if i == m.Cursor {
			cursor = "▸ "
			style = SelectedStyle
		}
		s.WriteString(style.Render(cursor + opt))
		s.WriteString("\n")
	}

	s.WriteString("\n")
	s.WriteString(HelpStyle.Render("↑/k up • ↓/j down • [Enter] select • [Esc/q] back"))

	return s.String()
}

func (m Model) renderKeymapCategory() string {
	var s strings.Builder

	if m.SelectedCategory >= len(m.KeymapCategories) {
		return ErrorStyle.Render("Category not found")
	}

	category := m.KeymapCategories[m.SelectedCategory]

	s.WriteString(TitleStyle.Render(category.Name))
	s.WriteString("\n")
	s.WriteString(MutedStyle.Render(category.Description))
	s.WriteString("\n\n")

	// Table header
	header := fmt.Sprintf("%-15s %-6s %s", "Keys", "Mode", "Description")
	s.WriteString(SubtitleStyle.Render(header))
	s.WriteString("\n")
	s.WriteString(MutedStyle.Render(strings.Repeat("─", 60)))
	s.WriteString("\n")

	// Calculate visible items based on terminal height
	// Reserve space for: title(1) + description(1) + blank(1) + header(1) + separator(1) + scroll info(2) + help(2) = 9 lines
	visibleItems := m.Height - 9
	if visibleItems < 5 {
		visibleItems = 5 // Minimum 5 items
	}
	if visibleItems > len(category.Keymaps) {
		visibleItems = len(category.Keymaps)
	}

	// Keymaps with scrolling
	start := m.KeymapScroll
	end := start + visibleItems
	if end > len(category.Keymaps) {
		end = len(category.Keymaps)
		start = end - visibleItems
		if start < 0 {
			start = 0
		}
	}

	for i := start; i < end; i++ {
		km := category.Keymaps[i]
		s.WriteString(KeyStyle.Render(km.Keys))
		s.WriteString(MutedStyle.Render(fmt.Sprintf(" %-6s ", km.Mode)))
		s.WriteString(InfoStyle.Render(km.Description))
		s.WriteString("\n")
	}

	// Scroll indicator
	if len(category.Keymaps) > visibleItems {
		s.WriteString("\n")
		scrollInfo := fmt.Sprintf("Showing %d-%d of %d", start+1, end, len(category.Keymaps))
		s.WriteString(MutedStyle.Render(scrollInfo))
	}

	s.WriteString("\n\n")
	s.WriteString(HelpStyle.Render("↑/k up • ↓/j down • [Enter/Esc/q] back"))

	return s.String()
}

// renderToolKeymapsMenu renders the tool selection menu (Neovim, Tmux, Zellij, Ghostty)
func (m Model) renderToolKeymapsMenu() string {
	var s strings.Builder

	s.WriteString(TitleStyle.Render(m.GetScreenTitle()))
	s.WriteString("\n")
	s.WriteString(MutedStyle.Render("Select a tool to view its keybindings"))
	s.WriteString("\n\n")

	// Menu
	options := m.GetCurrentOptions()
	for i, opt := range options {
		if strings.HasPrefix(opt, "───") {
			s.WriteString(MutedStyle.Render(opt))
			s.WriteString("\n")
			continue
		}

		cursor := "  "
		style := UnselectedStyle
		if i == m.Cursor {
			cursor = "▸ "
			style = SelectedStyle
		}
		s.WriteString(style.Render(cursor + opt))
		s.WriteString("\n")
	}

	s.WriteString("\n")
	s.WriteString(HelpStyle.Render("↑/k up • ↓/j down • [Enter] select • [Esc/q] back"))

	return s.String()
}

// renderTmuxKeymapsMenu renders the Tmux keymap categories menu
func (m Model) renderTmuxKeymapsMenu() string {
	var s strings.Builder

	s.WriteString(TitleStyle.Render(m.GetScreenTitle()))
	s.WriteString("\n")
	s.WriteString(MutedStyle.Render("Select a category to view Tmux keybindings"))
	s.WriteString("\n\n")

	// Menu
	options := m.GetCurrentOptions()
	for i, opt := range options {
		if strings.HasPrefix(opt, "───") {
			s.WriteString(MutedStyle.Render(opt))
			s.WriteString("\n")
			continue
		}

		cursor := "  "
		style := UnselectedStyle
		if i == m.Cursor {
			cursor = "▸ "
			style = SelectedStyle
		}
		s.WriteString(style.Render(cursor + opt))
		s.WriteString("\n")
	}

	s.WriteString("\n")
	s.WriteString(HelpStyle.Render("↑/k up • ↓/j down • [Enter] select • [Esc/q] back"))

	return s.String()
}

// renderTmuxKeymapCategory renders a specific Tmux keymap category
func (m Model) renderTmuxKeymapCategory() string {
	var s strings.Builder

	if m.TmuxSelectedCategory >= len(m.TmuxKeymapCategories) {
		return ErrorStyle.Render("Category not found")
	}

	category := m.TmuxKeymapCategories[m.TmuxSelectedCategory]

	s.WriteString(TitleStyle.Render(category.Name))
	s.WriteString("\n")
	s.WriteString(MutedStyle.Render(category.Description))
	s.WriteString("\n\n")

	// Table header
	header := fmt.Sprintf("%-20s %-6s %s", "Keys", "Mode", "Description")
	s.WriteString(SubtitleStyle.Render(header))
	s.WriteString("\n")
	s.WriteString(MutedStyle.Render(strings.Repeat("─", 60)))
	s.WriteString("\n")

	// Calculate visible items
	visibleItems := m.Height - 9
	if visibleItems < 5 {
		visibleItems = 5
	}
	if visibleItems > len(category.Keymaps) {
		visibleItems = len(category.Keymaps)
	}

	// Keymaps with scrolling
	start := m.TmuxKeymapScroll
	end := start + visibleItems
	if end > len(category.Keymaps) {
		end = len(category.Keymaps)
		start = end - visibleItems
		if start < 0 {
			start = 0
		}
	}

	for i := start; i < end; i++ {
		km := category.Keymaps[i]
		s.WriteString(KeyStyle.Render(km.Keys))
		s.WriteString(MutedStyle.Render(fmt.Sprintf(" %-6s ", km.Mode)))
		s.WriteString(InfoStyle.Render(km.Description))
		s.WriteString("\n")
	}

	// Scroll indicator
	if len(category.Keymaps) > visibleItems {
		s.WriteString("\n")
		scrollInfo := fmt.Sprintf("Showing %d-%d of %d", start+1, end, len(category.Keymaps))
		s.WriteString(MutedStyle.Render(scrollInfo))
	}

	s.WriteString("\n\n")
	s.WriteString(HelpStyle.Render("↑/k up • ↓/j down • [Enter/Esc/q] back"))

	return s.String()
}

// renderZellijKeymapsMenu renders the Zellij keymap categories menu
func (m Model) renderZellijKeymapsMenu() string {
	var s strings.Builder

	s.WriteString(TitleStyle.Render(m.GetScreenTitle()))
	s.WriteString("\n")
	s.WriteString(MutedStyle.Render("Select a category to view Zellij keybindings"))
	s.WriteString("\n\n")

	// Menu
	options := m.GetCurrentOptions()
	for i, opt := range options {
		if strings.HasPrefix(opt, "───") {
			s.WriteString(MutedStyle.Render(opt))
			s.WriteString("\n")
			continue
		}

		cursor := "  "
		style := UnselectedStyle
		if i == m.Cursor {
			cursor = "▸ "
			style = SelectedStyle
		}
		s.WriteString(style.Render(cursor + opt))
		s.WriteString("\n")
	}

	s.WriteString("\n")
	s.WriteString(HelpStyle.Render("↑/k up • ↓/j down • [Enter] select • [Esc/q] back"))

	return s.String()
}

// renderZellijKeymapCategory renders a specific Zellij keymap category
func (m Model) renderZellijKeymapCategory() string {
	var s strings.Builder

	if m.ZellijSelectedCategory >= len(m.ZellijKeymapCategories) {
		return ErrorStyle.Render("Category not found")
	}

	category := m.ZellijKeymapCategories[m.ZellijSelectedCategory]

	s.WriteString(TitleStyle.Render(category.Name))
	s.WriteString("\n")
	s.WriteString(MutedStyle.Render(category.Description))
	s.WriteString("\n\n")

	// Table header
	header := fmt.Sprintf("%-15s %-8s %s", "Keys", "Mode", "Description")
	s.WriteString(SubtitleStyle.Render(header))
	s.WriteString("\n")
	s.WriteString(MutedStyle.Render(strings.Repeat("─", 60)))
	s.WriteString("\n")

	// Calculate visible items
	visibleItems := m.Height - 9
	if visibleItems < 5 {
		visibleItems = 5
	}
	if visibleItems > len(category.Keymaps) {
		visibleItems = len(category.Keymaps)
	}

	// Keymaps with scrolling
	start := m.ZellijKeymapScroll
	end := start + visibleItems
	if end > len(category.Keymaps) {
		end = len(category.Keymaps)
		start = end - visibleItems
		if start < 0 {
			start = 0
		}
	}

	for i := start; i < end; i++ {
		km := category.Keymaps[i]
		s.WriteString(KeyStyle.Render(km.Keys))
		s.WriteString(MutedStyle.Render(fmt.Sprintf(" %-8s ", km.Mode)))
		s.WriteString(InfoStyle.Render(km.Description))
		s.WriteString("\n")
	}

	// Scroll indicator
	if len(category.Keymaps) > visibleItems {
		s.WriteString("\n")
		scrollInfo := fmt.Sprintf("Showing %d-%d of %d", start+1, end, len(category.Keymaps))
		s.WriteString(MutedStyle.Render(scrollInfo))
	}

	s.WriteString("\n\n")
	s.WriteString(HelpStyle.Render("↑/k up • ↓/j down • [Enter/Esc/q] back"))

	return s.String()
}

// renderGhosttyKeymapsMenu renders the Ghostty keymap categories menu
func (m Model) renderGhosttyKeymapsMenu() string {
	var s strings.Builder

	s.WriteString(TitleStyle.Render(m.GetScreenTitle()))
	s.WriteString("\n")
	s.WriteString(MutedStyle.Render("Select a category to view Ghostty keybindings"))
	s.WriteString("\n\n")

	// Menu
	options := m.GetCurrentOptions()
	for i, opt := range options {
		if strings.HasPrefix(opt, "───") {
			s.WriteString(MutedStyle.Render(opt))
			s.WriteString("\n")
			continue
		}

		cursor := "  "
		style := UnselectedStyle
		if i == m.Cursor {
			cursor = "▸ "
			style = SelectedStyle
		}
		s.WriteString(style.Render(cursor + opt))
		s.WriteString("\n")
	}

	s.WriteString("\n")
	s.WriteString(HelpStyle.Render("↑/k up • ↓/j down • [Enter] select • [Esc/q] back"))

	return s.String()
}

// renderGhosttyKeymapCategory renders a specific Ghostty keymap category
func (m Model) renderGhosttyKeymapCategory() string {
	var s strings.Builder

	if m.GhosttySelectedCategory >= len(m.GhosttyKeymapCategories) {
		return ErrorStyle.Render("Category not found")
	}

	category := m.GhosttyKeymapCategories[m.GhosttySelectedCategory]

	s.WriteString(TitleStyle.Render(category.Name))
	s.WriteString("\n")
	s.WriteString(MutedStyle.Render(category.Description))
	s.WriteString("\n\n")

	// Table header
	header := fmt.Sprintf("%-18s %-6s %s", "Keys", "Mode", "Description")
	s.WriteString(SubtitleStyle.Render(header))
	s.WriteString("\n")
	s.WriteString(MutedStyle.Render(strings.Repeat("─", 60)))
	s.WriteString("\n")

	// Calculate visible items
	visibleItems := m.Height - 9
	if visibleItems < 5 {
		visibleItems = 5
	}
	if visibleItems > len(category.Keymaps) {
		visibleItems = len(category.Keymaps)
	}

	// Keymaps with scrolling
	start := m.GhosttyKeymapScroll
	end := start + visibleItems
	if end > len(category.Keymaps) {
		end = len(category.Keymaps)
		start = end - visibleItems
		if start < 0 {
			start = 0
		}
	}

	for i := start; i < end; i++ {
		km := category.Keymaps[i]
		s.WriteString(KeyStyle.Render(km.Keys))
		s.WriteString(MutedStyle.Render(fmt.Sprintf(" %-6s ", km.Mode)))
		s.WriteString(InfoStyle.Render(km.Description))
		s.WriteString("\n")
	}

	// Scroll indicator
	if len(category.Keymaps) > visibleItems {
		s.WriteString("\n")
		scrollInfo := fmt.Sprintf("Showing %d-%d of %d", start+1, end, len(category.Keymaps))
		s.WriteString(MutedStyle.Render(scrollInfo))
	}

	s.WriteString("\n\n")
	s.WriteString(HelpStyle.Render("↑/k up • ↓/j down • [Enter/Esc/q] back"))

	return s.String()
}

func (m Model) renderLazyVimMenu() string {
	var s strings.Builder

	s.WriteString(TitleStyle.Render(m.GetScreenTitle()))
	s.WriteString("\n")
	s.WriteString(MutedStyle.Render("Learn how to use and customize LazyVim"))
	s.WriteString("\n\n")

	// Menu
	options := m.GetCurrentOptions()
	for i, opt := range options {
		if strings.HasPrefix(opt, "───") {
			s.WriteString(MutedStyle.Render(opt))
			s.WriteString("\n")
			continue
		}

		cursor := "  "
		style := UnselectedStyle
		if i == m.Cursor {
			cursor = "▸ "
			style = SelectedStyle
		}
		s.WriteString(style.Render(cursor + opt))
		s.WriteString("\n")
	}

	s.WriteString("\n")
	s.WriteString(HelpStyle.Render("↑/k up • ↓/j down • [Enter] select • [Esc/q] back"))

	return s.String()
}

func (m Model) renderLazyVimTopic() string {
	var s strings.Builder

	if m.SelectedLazyVimTopic >= len(m.LazyVimTopics) {
		return ErrorStyle.Render("Topic not found")
	}

	topic := m.LazyVimTopics[m.SelectedLazyVimTopic]

	s.WriteString(TitleStyle.Render(topic.Title))
	s.WriteString("\n")
	s.WriteString(SubtitleStyle.Render(topic.Description))
	s.WriteString("\n\n")

	// Build all content
	var allLines []string

	// Content
	allLines = append(allLines, topic.Content...)
	allLines = append(allLines, "") // Empty line

	// Code example
	if topic.CodeExample != "" {
		allLines = append(allLines, "📝 Example:")
		allLines = append(allLines, "")
		codeLines := strings.Split(topic.CodeExample, "\n")
		allLines = append(allLines, codeLines...)
		allLines = append(allLines, "") // Empty line
	}

	// Tips
	if len(topic.Tips) > 0 {
		allLines = append(allLines, "💡 Tips:")
		for _, tip := range topic.Tips {
			allLines = append(allLines, "  • "+tip)
		}
	}

	// Calculate view height based on terminal size
	// Reserve space for: title(1) + description(1) + blank(2) + scroll info(2) + help(2) = 8 lines
	viewHeight := m.Height - 8
	if viewHeight < 10 {
		viewHeight = 10 // Minimum
	}

	// Apply scrolling
	start := m.LazyVimScroll
	end := start + viewHeight
	if end > len(allLines) {
		end = len(allLines)
	}
	if start > len(allLines) {
		start = 0
	}

	for i := start; i < end; i++ {
		line := allLines[i]
		// Style code lines differently
		if strings.HasPrefix(line, "--") || strings.HasPrefix(line, "local") ||
			strings.HasPrefix(line, "return") || strings.HasPrefix(line, "{") ||
			strings.HasPrefix(line, "}") || strings.HasPrefix(line, "  ") ||
			strings.HasPrefix(line, "map(") || strings.HasPrefix(line, "vim.") ||
			strings.HasPrefix(line, "require") {
			s.WriteString(CodeStyle.Render(line))
		} else if strings.HasPrefix(line, "📝") || strings.HasPrefix(line, "💡") {
			s.WriteString(SubtitleStyle.Render(line))
		} else if strings.HasPrefix(line, "  •") {
			s.WriteString(InfoStyle.Render(line))
		} else if strings.HasPrefix(line, "•") {
			s.WriteString(MutedStyle.Render(line))
		} else {
			s.WriteString(InfoStyle.Render(line))
		}
		s.WriteString("\n")
	}

	// Scroll indicator
	if len(allLines) > viewHeight {
		s.WriteString("\n")
		scrollInfo := fmt.Sprintf("Lines %d-%d of %d (↑↓ to scroll, PgUp/PgDn for fast scroll)", start+1, end, len(allLines))
		s.WriteString(MutedStyle.Render(scrollInfo))
	}

	s.WriteString("\n\n")
	s.WriteString(HelpStyle.Render("↑/k up • ↓/j down • PgUp/PgDn • [Enter/Esc/q] back"))

	return s.String()
}

// Spinner frames for running steps
var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

func (m Model) renderInstalling() string {
	var s strings.Builder

	s.WriteString(TitleStyle.Render("🚀 Installing Gentleman.Dots"))
	s.WriteString("\n\n")

	// Progress steps
	for i, step := range m.Steps {
		var icon string
		var style lipgloss.Style

		switch step.Status {
		case StatusPending:
			icon = "○"
			style = MutedStyle
		case StatusRunning:
			// Animated spinner
			icon = spinnerFrames[m.SpinnerFrame%len(spinnerFrames)]
			style = WarningStyle
		case StatusDone:
			icon = "✓"
			style = SuccessStyle
		case StatusFailed:
			icon = "✗"
			style = ErrorStyle
		case StatusSkipped:
			icon = "⊘"
			style = MutedStyle
		}

		line := fmt.Sprintf("%s %s", icon, step.Name)
		s.WriteString(style.Render(line))
		s.WriteString("\n")

		// Show current step description
		if i == m.CurrentStep && step.Status == StatusRunning {
			s.WriteString(MutedStyle.Render("   " + step.Description))
			s.WriteString("\n")
		}
	}

	// Log output if details enabled
	if m.ShowDetails && len(m.LogLines) > 0 {
		s.WriteString("\n")
		s.WriteString(BoxStyle.Render(strings.Join(m.LogLines[max(0, len(m.LogLines)-10):], "\n")))
	}

	s.WriteString("\n")
	s.WriteString(HelpStyle.Render("[space+d] toggle details"))

	return s.String()
}

func (m Model) renderComplete() string {
	var s strings.Builder

	s.WriteString(SuccessStyle.Render("✨ Installation Complete! ✨"))
	s.WriteString("\n\n")

	// Summary
	s.WriteString(TitleStyle.Render("Summary"))
	s.WriteString("\n")

	items := []string{
		fmt.Sprintf("OS: %s", m.Choices.OS),
		fmt.Sprintf("Terminal: %s", m.Choices.Terminal),
		fmt.Sprintf("Shell: %s", m.Choices.Shell),
		fmt.Sprintf("Window Manager: %s", m.Choices.WindowMgr),
	}

	if m.Choices.InstallFont {
		items = append(items, "Font: Iosevka Term Nerd Font")
	}
	if m.Choices.InstallNvim {
		items = append(items, "Editor: Neovim with Gentleman config")
	}

	for _, item := range items {
		s.WriteString(InfoStyle.Render("  • " + item))
		s.WriteString("\n")
	}

	// Shell change instructions
	shell := m.Choices.Shell
	shellCmd := shell
	if shell == "nushell" {
		shellCmd = "nu"
	}

	s.WriteString("\n")
	s.WriteString(TitleStyle.Render("Next Steps"))
	s.WriteString("\n\n")

	s.WriteString(InfoStyle.Render("1. To use your new shell now, run:"))
	s.WriteString("\n")
	s.WriteString(HighlightStyle.Render(fmt.Sprintf("   exec %s", shellCmd)))
	s.WriteString("\n\n")

	s.WriteString(InfoStyle.Render("2. To make it your default shell, run:"))
	s.WriteString("\n")

	// Show the commands needed to set default shell
	s.WriteString(MutedStyle.Render(fmt.Sprintf("   # Find shell path\n   which %s", shellCmd)))
	s.WriteString("\n")
	s.WriteString(MutedStyle.Render("   # Add to allowed shells (needs sudo)"))
	s.WriteString("\n")
	s.WriteString(HighlightStyle.Render(fmt.Sprintf("   sudo sh -c 'echo $(which %s) >> /etc/shells'", shellCmd)))
	s.WriteString("\n")
	s.WriteString(MutedStyle.Render("   # Set as default"))
	s.WriteString("\n")
	s.WriteString(HighlightStyle.Render(fmt.Sprintf("   chsh -s $(which %s)", shellCmd)))
	s.WriteString("\n\n")

	s.WriteString(HelpStyle.Render("Press [Enter] or [q] to exit"))

	return s.String()
}

func (m Model) renderError() string {
	var s strings.Builder

	s.WriteString(ErrorStyle.Render("❌ Installation Failed"))
	s.WriteString("\n\n")

	s.WriteString(MutedStyle.Render("Error:"))
	s.WriteString("\n")
	s.WriteString(ErrorStyle.Render(m.ErrorMsg))
	s.WriteString("\n\n")

	// Show last few log lines for context
	if len(m.LogLines) > 0 {
		s.WriteString(MutedStyle.Render("Recent logs:"))
		s.WriteString("\n")
		// Show last 5 log lines
		startIdx := len(m.LogLines) - 5
		if startIdx < 0 {
			startIdx = 0
		}
		for _, line := range m.LogLines[startIdx:] {
			s.WriteString(InfoStyle.Render("  " + line))
			s.WriteString("\n")
		}
		s.WriteString("\n")
	}

	s.WriteString(HelpStyle.Render("[r] retry • [space+q] quit"))

	return s.String()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (m Model) renderBackupConfirm() string {
	var s strings.Builder

	s.WriteString(TitleStyle.Render(m.GetScreenTitle()))
	s.WriteString("\n")
	s.WriteString(MutedStyle.Render("The following configs will be overwritten:"))
	s.WriteString("\n\n")

	// List existing configs
	for _, config := range m.ExistingConfigs {
		s.WriteString(WarningStyle.Render("  ⚠️  " + config))
		s.WriteString("\n")
	}

	s.WriteString("\n")
	s.WriteString(InfoStyle.Render("Creating a backup allows you to restore later if needed."))
	s.WriteString("\n\n")

	// Options
	options := m.GetCurrentOptions()
	for i, opt := range options {
		cursor := "  "
		style := UnselectedStyle
		if i == m.Cursor {
			cursor = "▸ "
			style = SelectedStyle
		}
		s.WriteString(style.Render(cursor + opt))
		s.WriteString("\n")
	}

	s.WriteString("\n")
	s.WriteString(HelpStyle.Render("↑/k up • ↓/j down • [Enter] select • [Esc] back"))

	return s.String()
}

func (m Model) renderRestoreBackup() string {
	var s strings.Builder

	s.WriteString(TitleStyle.Render(m.GetScreenTitle()))
	s.WriteString("\n")
	s.WriteString(MutedStyle.Render("Select a backup to restore or delete"))
	s.WriteString("\n\n")

	if len(m.AvailableBackups) == 0 {
		s.WriteString(MutedStyle.Render("No backups found."))
		s.WriteString("\n")
	} else {
		// List backups
		for i, backup := range m.AvailableBackups {
			cursor := "  "
			style := UnselectedStyle
			if i == m.Cursor {
				cursor = "▸ "
				style = SelectedStyle
			}

			// Format: timestamp + item count
			label := fmt.Sprintf("📁 %s (%d items)", backup.Timestamp.Format("2006-01-02 15:04:05"), len(backup.Files))
			s.WriteString(style.Render(cursor + label))
			s.WriteString("\n")
		}
	}

	// Separator and Back
	s.WriteString(MutedStyle.Render("─────────────"))
	s.WriteString("\n")

	backIdx := len(m.AvailableBackups) + 1
	cursor := "  "
	style := UnselectedStyle
	if m.Cursor == backIdx {
		cursor = "▸ "
		style = SelectedStyle
	}
	s.WriteString(style.Render(cursor + "← Back"))
	s.WriteString("\n")

	s.WriteString("\n")
	s.WriteString(HelpStyle.Render("↑/k up • ↓/j down • [Enter] select • [Esc] back"))

	return s.String()
}

func (m Model) renderRestoreConfirm() string {
	var s strings.Builder

	if m.SelectedBackup >= len(m.AvailableBackups) {
		return ErrorStyle.Render("No backup selected")
	}

	backup := m.AvailableBackups[m.SelectedBackup]

	s.WriteString(TitleStyle.Render(m.GetScreenTitle()))
	s.WriteString("\n")
	s.WriteString(MutedStyle.Render("Backup from: " + backup.Timestamp.Format("2006-01-02 15:04:05")))
	s.WriteString("\n\n")

	// List files in backup
	s.WriteString(SubtitleStyle.Render("Contents:"))
	s.WriteString("\n")
	for _, file := range backup.Files {
		s.WriteString(InfoStyle.Render("  • " + file))
		s.WriteString("\n")
	}

	s.WriteString("\n")
	s.WriteString(WarningStyle.Render("⚠️  Restoring will overwrite your current configs!"))
	s.WriteString("\n\n")

	// Options
	options := m.GetCurrentOptions()
	for i, opt := range options {
		cursor := "  "
		style := UnselectedStyle
		if i == m.Cursor {
			cursor = "▸ "
			style = SelectedStyle
		}
		s.WriteString(style.Render(cursor + opt))
		s.WriteString("\n")
	}

	s.WriteString("\n")
	s.WriteString(HelpStyle.Render("↑/k up • ↓/j down • [Enter] select • [Esc] cancel"))

	return s.String()
}

// ============================================================================
// Trainer Views
// ============================================================================

func (m Model) renderTrainerMenu() string {
	var s strings.Builder

	// Header
	s.WriteString(TitleStyle.Render("🎮 Vim Mastery Trainer"))
	s.WriteString("\n")
	s.WriteString(MutedStyle.Render("Master Vim motions through progressive challenges"))
	s.WriteString("\n\n")

	// Stats bar
	if m.TrainerStats != nil {
		score := fmt.Sprintf("Score: %d", m.TrainerStats.TotalScore)
		streak := fmt.Sprintf("Streak: %d", m.TrainerStats.CurrentStreak)
		bosses := fmt.Sprintf("Bosses: %d/7", len(m.TrainerStats.BossesDefeated))
		s.WriteString(InfoStyle.Render(fmt.Sprintf("📊 %s  |  🔥 %s  |  👑 %s", score, streak, bosses)))
		s.WriteString("\n\n")
	}

	// Module list
	s.WriteString(SubtitleStyle.Render("Select a Module:"))
	s.WriteString("\n\n")

	for i, module := range m.TrainerModules {
		isUnlocked := m.TrainerStats != nil && m.TrainerStats.IsModuleUnlocked(module.ID)
		isBossDefeated := m.TrainerStats != nil && m.TrainerStats.IsBossDefeated(module.ID)
		isLessonsComplete := m.TrainerStats != nil && m.TrainerStats.IsLessonsComplete(module.ID)
		isPracticeReady := m.TrainerStats != nil && m.TrainerStats.IsPracticeReady(module.ID)
		isBossReady := m.TrainerStats != nil && m.TrainerStats.IsBossReady(module.ID)

		cursor := "  "
		style := UnselectedStyle
		if i == m.TrainerCursor {
			cursor = "▸ "
			style = SelectedStyle
		}

		// Module name with status indicators
		status := ""
		if !isUnlocked {
			status = "🔒"
			style = MutedStyle
		} else if isBossDefeated {
			status = "👑"
		} else if isBossReady {
			status = "⚔️"
		} else if isPracticeReady {
			status = "🎯"
		} else if isLessonsComplete {
			status = "📚"
		} else {
			status = "📖"
		}

		line := fmt.Sprintf("%s %s %s - %s", status, module.Icon, module.Name, module.Description)
		s.WriteString(style.Render(cursor + line))
		s.WriteString("\n")

		// Show progress for selected module
		if i == m.TrainerCursor && isUnlocked && m.TrainerStats != nil {
			progress := m.TrainerStats.GetModuleProgress(module.ID)
			var progressLine string
			if progress.LessonsTotal > 0 {
				lessonsPercent := float64(progress.LessonsCompleted) / float64(progress.LessonsTotal) * 100
				progressLine = fmt.Sprintf("     Lessons: %d/%d (%.0f%%)", progress.LessonsCompleted, progress.LessonsTotal, lessonsPercent)
			} else {
				progressLine = "     Lessons: 0/0"
			}
			if progress.PracticeAttempts > 0 {
				progressLine += fmt.Sprintf("  |  Practice: %.0f%%", progress.PracticeAccuracy*100)
			}

			// Show mastery progress for practice mode
			if isPracticeReady {
				practiceStats := trainer.GetPracticeStatsForModule(module.ID, progress)
				if practiceStats.TotalExercises > 0 {
					progressLine += fmt.Sprintf("  |  Mastered: %d/%d", practiceStats.MasteredCount, practiceStats.TotalExercises)
					if practiceStats.PracticeComplete {
						progressLine += " ✅"
					}
				}
			}

			s.WriteString(MutedStyle.Render(progressLine))
			s.WriteString("\n")
		}
	}

	// Show message if any
	if m.TrainerMessage != "" {
		s.WriteString("\n")
		s.WriteString(WarningStyle.Render(m.TrainerMessage))
		s.WriteString("\n")
	}

	// Help
	s.WriteString("\n")
	s.WriteString(HelpStyle.Render("↑/k up • ↓/j down • [Enter/l] lesson • [p] practice • [b] boss • [r] reset • [q/Esc] back"))

	return s.String()
}

func (m Model) renderTrainerExercise(mode string) string {
	var s strings.Builder

	if m.TrainerGameState == nil || m.TrainerGameState.CurrentExercise == nil {
		s.WriteString(ErrorStyle.Render("No exercise loaded"))
		return s.String()
	}

	exercise := m.TrainerGameState.CurrentExercise

	// Header with mode
	title := fmt.Sprintf("🎮 %s Mode: %s", mode, string(m.TrainerGameState.CurrentModule))
	s.WriteString(TitleStyle.Render(title))
	s.WriteString("\n")

	// Progress bar
	var progressText string
	if m.TrainerGameState.IsLessonMode {
		current := m.TrainerGameState.ExerciseIndex + 1
		total := len(m.TrainerGameState.Exercises)
		progressText = fmt.Sprintf("Exercise %d of %d", current, total)
	} else {
		progressText = fmt.Sprintf("Score: %d | Streak: %d", m.TrainerGameState.SessionScore, m.TrainerGameState.CurrentStreak)
	}
	s.WriteString(MutedStyle.Render(progressText))
	s.WriteString("\n\n")

	// Mission
	s.WriteString(SubtitleStyle.Render("📋 Mission:"))
	s.WriteString("\n")
	s.WriteString(InfoStyle.Render("   " + exercise.Mission))
	s.WriteString("\n\n")

	// Calculate simulated cursor position and selection based on current input
	startPos := exercise.CursorPos
	simResult := trainer.SimulateMotionsWithSelection(startPos, exercise.Code, m.TrainerInput)
	simPos := simResult.Position
	selection := simResult.Selection

	// Code display with cursors and selection
	s.WriteString(SubtitleStyle.Render("📝 Code:"))
	s.WriteString("\n")
	s.WriteString(MutedStyle.Render(strings.Repeat("─", 60)))
	s.WriteString("\n")

	for lineNum, line := range exercise.Code {
		lineNumStr := fmt.Sprintf("%2d │ ", lineNum+1)
		s.WriteString(MutedStyle.Render(lineNumStr))

		// Check if there's an active selection on this line
		if selection.Active && lineNum == selection.StartLine {
			s.WriteString(renderLineWithSelection(line, startPos, selection))
			s.WriteString("\n")
			continue
		}

		// Determine which cursors are on this line
		startOnLine := lineNum == startPos.Line
		currentOnLine := lineNum == simPos.Line
		samePosition := startOnLine && currentOnLine && startPos.Col == simPos.Col

		// Helper to get cursor character (use space block for empty lines)
		getCursorChar := func(line string, col int) string {
			if col < len(line) {
				return string(line[col])
			}
			return " " // Space with background color for empty line cursor
		}

		if samePosition {
			// Both cursors at same position - show start cursor (they haven't moved yet)
			if startPos.Col < len(line) {
				before := line[:startPos.Col]
				cursor := getCursorChar(line, startPos.Col)
				after := ""
				if startPos.Col+1 < len(line) {
					after = line[startPos.Col+1:]
				}
				s.WriteString(CodeStyle.Render(before))
				s.WriteString(StartCursorStyle.Render(cursor))
				s.WriteString(CodeStyle.Render(after))
			} else {
				// Cursor at end of line or empty line
				s.WriteString(CodeStyle.Render(line))
				s.WriteString(StartCursorStyle.Render(" "))
			}
		} else if startOnLine && currentOnLine {
			// Both cursors on same line but different positions
			s.WriteString(renderLineWithTwoCursors(line, startPos.Col, simPos.Col))
		} else if startOnLine {
			// Only start cursor on this line
			if startPos.Col < len(line) {
				before := line[:startPos.Col]
				cursor := getCursorChar(line, startPos.Col)
				after := ""
				if startPos.Col+1 < len(line) {
					after = line[startPos.Col+1:]
				}
				s.WriteString(CodeStyle.Render(before))
				s.WriteString(StartCursorStyle.Render(cursor))
				s.WriteString(CodeStyle.Render(after))
			} else {
				// Cursor at end of line or empty line
				s.WriteString(CodeStyle.Render(line))
				s.WriteString(StartCursorStyle.Render(" "))
			}
		} else if currentOnLine {
			// Only current cursor on this line
			if simPos.Col < len(line) {
				before := line[:simPos.Col]
				cursor := getCursorChar(line, simPos.Col)
				after := ""
				if simPos.Col+1 < len(line) {
					after = line[simPos.Col+1:]
				}
				s.WriteString(CodeStyle.Render(before))
				s.WriteString(CurrentCursorStyle.Render(cursor))
				s.WriteString(CodeStyle.Render(after))
			} else {
				// Cursor at end of line or empty line
				s.WriteString(CodeStyle.Render(line))
				s.WriteString(CurrentCursorStyle.Render(" "))
			}
		} else {
			// No cursors on this line
			s.WriteString(CodeStyle.Render(line))
		}
		s.WriteString("\n")
	}

	s.WriteString(MutedStyle.Render(strings.Repeat("─", 60)))
	s.WriteString("\n\n")

	// Input field
	s.WriteString(SubtitleStyle.Render("⌨️  Your answer:"))
	s.WriteString("\n")
	inputDisplay := m.TrainerInput
	if inputDisplay == "" {
		inputDisplay = "..."
	}
	s.WriteString(BoxStyle.Render(KeyStyle.Render(inputDisplay)))
	s.WriteString("\n")

	// Show message/hint if any
	if m.TrainerMessage != "" {
		s.WriteString("\n")
		s.WriteString(InfoStyle.Render(m.TrainerMessage))
		s.WriteString("\n")
	}

	// Help
	s.WriteString("\n")
	s.WriteString(HelpStyle.Render("Type command • [Enter] submit • [Tab] hint • [Backspace] clear • [Esc] quit"))

	return s.String()
}

// renderLineWithTwoCursors renders a line with both start and current cursor
func renderLineWithTwoCursors(line string, startCol, currentCol int) string {
	var result strings.Builder

	// Helper to get cursor character (use space for empty/end of line)
	getCursorChar := func(col int) string {
		if col < len(line) {
			return string(line[col])
		}
		return " "
	}

	// Determine order of cursors
	firstCol, secondCol := startCol, currentCol
	firstStyle, secondStyle := StartCursorStyle, CurrentCursorStyle
	if currentCol < startCol {
		firstCol, secondCol = currentCol, startCol
		firstStyle, secondStyle = CurrentCursorStyle, StartCursorStyle
	}

	// Handle empty line case
	if len(line) == 0 {
		// Both cursors on empty line at col 0
		if firstCol == secondCol {
			result.WriteString(firstStyle.Render(" "))
		} else {
			// This shouldn't happen on empty line, but handle it
			result.WriteString(firstStyle.Render(" "))
			result.WriteString(secondStyle.Render(" "))
		}
		return result.String()
	}

	// Build the line piece by piece
	// Part before first cursor
	if firstCol > 0 && firstCol <= len(line) {
		result.WriteString(CodeStyle.Render(line[:firstCol]))
	}

	// First cursor
	result.WriteString(firstStyle.Render(getCursorChar(firstCol)))

	// Part between cursors (if any)
	if firstCol+1 < secondCol && firstCol+1 < len(line) {
		endIdx := secondCol
		if endIdx > len(line) {
			endIdx = len(line)
		}
		if firstCol+1 < endIdx {
			result.WriteString(CodeStyle.Render(line[firstCol+1 : endIdx]))
		}
	}

	// Second cursor (only if different position from first)
	if secondCol > firstCol {
		result.WriteString(secondStyle.Render(getCursorChar(secondCol)))
	}

	// Part after second cursor
	if secondCol+1 < len(line) {
		result.WriteString(CodeStyle.Render(line[secondCol+1:]))
	}

	return result.String()
}

// renderLineWithSelection renders a line with visual selection highlighted
func renderLineWithSelection(line string, startPos trainer.Position, sel trainer.Selection) string {
	var result strings.Builder

	if len(line) == 0 {
		// Empty line with selection
		result.WriteString(SelectionStyle.Render(" "))
		return result.String()
	}

	// Clamp selection bounds to line length
	selStart := sel.StartCol
	selEnd := sel.EndCol
	if selStart < 0 {
		selStart = 0
	}
	if selEnd >= len(line) {
		selEnd = len(line) - 1
	}
	if selEnd < selStart {
		// Invalid selection, just render the line normally with start cursor
		if startPos.Col < len(line) {
			result.WriteString(CodeStyle.Render(line[:startPos.Col]))
			result.WriteString(StartCursorStyle.Render(string(line[startPos.Col])))
			if startPos.Col+1 < len(line) {
				result.WriteString(CodeStyle.Render(line[startPos.Col+1:]))
			}
		} else {
			result.WriteString(CodeStyle.Render(line))
		}
		return result.String()
	}

	// Render: [before selection] [SELECTION] [after selection]
	if selStart > 0 {
		result.WriteString(CodeStyle.Render(line[:selStart]))
	}

	// The selection itself
	selectedText := line[selStart : selEnd+1]
	result.WriteString(SelectionStyle.Render(selectedText))

	// After selection
	if selEnd+1 < len(line) {
		result.WriteString(CodeStyle.Render(line[selEnd+1:]))
	}

	return result.String()
}

func (m Model) renderTrainerBoss() string {
	var s strings.Builder

	if m.TrainerGameState == nil || m.TrainerGameState.CurrentBoss == nil {
		s.WriteString(ErrorStyle.Render("No boss loaded"))
		return s.String()
	}

	boss := m.TrainerGameState.CurrentBoss
	currentStep := m.TrainerGameState.BossStep

	// Boss header
	s.WriteString(DangerStyle.Render("⚔️  BOSS FIGHT: " + boss.Name))
	s.WriteString("\n")

	// Lives and progress
	lives := strings.Repeat("❤️ ", m.TrainerGameState.BossLives)
	lostLives := strings.Repeat("🖤 ", boss.Lives-m.TrainerGameState.BossLives)
	s.WriteString(fmt.Sprintf("Lives: %s%s  |  Step: %d/%d", lives, lostLives, currentStep+1, len(boss.Steps)))
	s.WriteString("\n\n")

	if currentStep < len(boss.Steps) {
		step := boss.Steps[currentStep]
		exercise := &step.Exercise

		// Mission
		s.WriteString(SubtitleStyle.Render("📋 Challenge:"))
		s.WriteString("\n")
		s.WriteString(InfoStyle.Render("   " + exercise.Mission))
		s.WriteString("\n\n")

		// Calculate simulated cursor position and selection based on current input
		startPos := exercise.CursorPos
		simResult := trainer.SimulateMotionsWithSelection(startPos, exercise.Code, m.TrainerInput)
		simPos := simResult.Position
		selection := simResult.Selection

		// Code display with cursors and selection
		s.WriteString(SubtitleStyle.Render("📝 Code:"))
		s.WriteString("\n")
		s.WriteString(MutedStyle.Render(strings.Repeat("─", 60)))
		s.WriteString("\n")

		for lineNum, line := range exercise.Code {
			lineNumStr := fmt.Sprintf("%2d │ ", lineNum+1)
			s.WriteString(MutedStyle.Render(lineNumStr))

			// Check if there's an active selection on this line
			if selection.Active && lineNum == selection.StartLine {
				s.WriteString(renderLineWithSelection(line, startPos, selection))
				s.WriteString("\n")
				continue
			}

			// Determine which cursors are on this line
			startOnLine := lineNum == startPos.Line
			currentOnLine := lineNum == simPos.Line
			samePosition := startOnLine && currentOnLine && startPos.Col == simPos.Col

			// Helper to get cursor character (use space block for empty lines)
			getCursorChar := func(line string, col int) string {
				if col < len(line) {
					return string(line[col])
				}
				return " " // Space with background color for empty line cursor
			}

			if samePosition {
				// Both cursors at same position
				if startPos.Col < len(line) {
					before := line[:startPos.Col]
					cursor := getCursorChar(line, startPos.Col)
					after := ""
					if startPos.Col+1 < len(line) {
						after = line[startPos.Col+1:]
					}
					s.WriteString(CodeStyle.Render(before))
					s.WriteString(StartCursorStyle.Render(cursor))
					s.WriteString(CodeStyle.Render(after))
				} else {
					// Cursor at end of line or empty line
					s.WriteString(CodeStyle.Render(line))
					s.WriteString(StartCursorStyle.Render(" "))
				}
			} else if startOnLine && currentOnLine {
				// Both cursors on same line but different positions
				s.WriteString(renderLineWithTwoCursors(line, startPos.Col, simPos.Col))
			} else if startOnLine {
				// Only start cursor on this line
				if startPos.Col < len(line) {
					before := line[:startPos.Col]
					cursor := getCursorChar(line, startPos.Col)
					after := ""
					if startPos.Col+1 < len(line) {
						after = line[startPos.Col+1:]
					}
					s.WriteString(CodeStyle.Render(before))
					s.WriteString(StartCursorStyle.Render(cursor))
					s.WriteString(CodeStyle.Render(after))
				} else {
					// Cursor at end of line or empty line
					s.WriteString(CodeStyle.Render(line))
					s.WriteString(StartCursorStyle.Render(" "))
				}
			} else if currentOnLine {
				// Only current cursor on this line
				if simPos.Col < len(line) {
					before := line[:simPos.Col]
					cursor := getCursorChar(line, simPos.Col)
					after := ""
					if simPos.Col+1 < len(line) {
						after = line[simPos.Col+1:]
					}
					s.WriteString(CodeStyle.Render(before))
					s.WriteString(CurrentCursorStyle.Render(cursor))
					s.WriteString(CodeStyle.Render(after))
				} else {
					// Cursor at end of line or empty line
					s.WriteString(CodeStyle.Render(line))
					s.WriteString(CurrentCursorStyle.Render(" "))
				}
			} else {
				// No cursors on this line
				s.WriteString(CodeStyle.Render(line))
			}
			s.WriteString("\n")
		}

		s.WriteString(MutedStyle.Render(strings.Repeat("─", 60)))
		s.WriteString("\n\n")

		// Input field
		s.WriteString(SubtitleStyle.Render("⌨️  Your answer:"))
		s.WriteString("\n")
		inputDisplay := m.TrainerInput
		if inputDisplay == "" {
			inputDisplay = "..."
		}
		s.WriteString(BoxStyle.Render(KeyStyle.Render(inputDisplay)))
		s.WriteString("\n")
	}

	// Show message if any
	if m.TrainerMessage != "" {
		s.WriteString("\n")
		s.WriteString(WarningStyle.Render(m.TrainerMessage))
		s.WriteString("\n")
	}

	// Help
	s.WriteString("\n")
	s.WriteString(HelpStyle.Render("Type command • [Enter] submit • [Esc] forfeit"))

	return s.String()
}

func (m Model) renderTrainerResult() string {
	var s strings.Builder

	// Result header
	if m.TrainerLastCorrect {
		s.WriteString(SuccessStyle.Render("✨ CORRECT! ✨"))
	} else {
		s.WriteString(ErrorStyle.Render("❌ INCORRECT"))
	}
	s.WriteString("\n\n")

	// Show message/explanation
	s.WriteString(InfoStyle.Render(m.TrainerMessage))
	s.WriteString("\n")

	if m.TrainerGameState != nil && m.TrainerGameState.CurrentExercise != nil {
		exercise := m.TrainerGameState.CurrentExercise
		if exercise.Explanation != "" {
			s.WriteString("\n")
			s.WriteString(SubtitleStyle.Render("📖 Explanation:"))
			s.WriteString("\n")
			s.WriteString(MutedStyle.Render("   " + exercise.Explanation))
			s.WriteString("\n")
		}
	}

	// Score info
	if m.TrainerGameState != nil {
		s.WriteString("\n")
		s.WriteString(MutedStyle.Render(fmt.Sprintf("Session Score: %d  |  Streak: %d", m.TrainerGameState.SessionScore, m.TrainerGameState.CurrentStreak)))
		s.WriteString("\n")
	}

	// Help
	s.WriteString("\n")
	s.WriteString(HelpStyle.Render("[Enter] continue • [Esc] back"))

	return s.String()
}

func (m Model) renderTrainerBossResult() string {
	var s strings.Builder

	// Victory or defeat
	if m.TrainerLastCorrect {
		s.WriteString(SuccessStyle.Render("🏆 VICTORY! 🏆"))
		s.WriteString("\n\n")
		if m.TrainerGameState != nil && m.TrainerGameState.CurrentBoss != nil {
			s.WriteString(TitleStyle.Render("You defeated " + m.TrainerGameState.CurrentBoss.Name + "!"))
			s.WriteString("\n\n")
			s.WriteString(InfoStyle.Render(fmt.Sprintf("Lives remaining: %s", strings.Repeat("❤️ ", m.TrainerGameState.BossLives))))
			s.WriteString("\n")
		}
		s.WriteString("\n")
		s.WriteString(SuccessStyle.Render("🎉 +500 bonus points!"))
		s.WriteString("\n")
		s.WriteString(SuccessStyle.Render("🔓 Next module unlocked!"))
	} else {
		s.WriteString(DangerStyle.Render("💀 DEFEATED 💀"))
		s.WriteString("\n\n")
		if m.TrainerGameState != nil && m.TrainerGameState.CurrentBoss != nil {
			s.WriteString(MutedStyle.Render(m.TrainerGameState.CurrentBoss.Name + " wins this time..."))
			s.WriteString("\n\n")
		}
		s.WriteString(InfoStyle.Render("Keep practicing and try again!"))
	}

	// Show message
	if m.TrainerMessage != "" {
		s.WriteString("\n\n")
		s.WriteString(MutedStyle.Render(m.TrainerMessage))
	}

	// Stats
	if m.TrainerStats != nil {
		s.WriteString("\n\n")
		s.WriteString(MutedStyle.Render(fmt.Sprintf("Total Score: %d  |  Bosses Defeated: %d/7", m.TrainerStats.TotalScore, len(m.TrainerStats.BossesDefeated))))
	}

	// Help
	s.WriteString("\n\n")
	s.WriteString(HelpStyle.Render("[Enter/Space/Esc/q] return to menu"))

	return s.String()
}
