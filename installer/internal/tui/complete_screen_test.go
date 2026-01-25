package tui

import (
	"strings"
	"testing"
)

// TestRenderComplete_FullInstallation tests complete screen with everything selected
func TestRenderComplete_FullInstallation(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenComplete
	m.Choices.OS = "darwin"
	m.Choices.Terminal = "ghostty"
	m.Choices.InstallFont = true
	m.Choices.Shell = "zsh"
	m.Choices.WindowMgr = "tmux"
	m.Choices.InstallNvim = true
	m.Choices.AIAssistants = []string{"opencode"}
	m.SkippedSteps = make(map[Screen]bool)
	m.AIAssistantsList = GetAvailableAIAssistants()

	output := m.renderComplete()

	// Must contain installation complete message
	if !strings.Contains(output, "Installation Complete") {
		t.Error("Should contain 'Installation Complete'")
	}

	// Must show OS
	if !strings.Contains(output, "OS: darwin") {
		t.Error("Should show OS selection")
	}

	// Must show Terminal with checkmark
	if !strings.Contains(output, "✓ Terminal: Ghostty") {
		t.Error("Should show Terminal: Ghostty with checkmark")
	}

	// Must show Font (sub-item of Terminal)
	if !strings.Contains(output, "Iosevka Nerd Font") {
		t.Error("Should show font when InstallFont=true")
	}

	// Must show Shell with checkmark
	if !strings.Contains(output, "✓ Shell: Zsh") {
		t.Error("Should show Shell: Zsh with checkmark")
	}

	// Must show Multiplexer with checkmark
	if !strings.Contains(output, "✓ Multiplexer: Tmux") {
		t.Error("Should show Multiplexer: Tmux with checkmark")
	}

	// Must show Neovim with checkmark
	if !strings.Contains(output, "✓ Neovim: LazyVim configuration") {
		t.Error("Should show Neovim with checkmark")
	}

	// Must show AI Assistant with checkmark
	if !strings.Contains(output, "✓ AI Assistant: OpenCode") {
		t.Error("Should show AI Assistant: OpenCode with checkmark")
	}

	// Must show shell exec command
	if !strings.Contains(output, "exec zsh") {
		t.Error("Should show 'exec zsh' command")
	}

	// Must show exit instructions
	if !strings.Contains(output, "Press [Enter] or [q] to exit") {
		t.Error("Should show exit instructions")
	}
}

// TestRenderComplete_OnlyAIAssistant tests when only AI Assistant is selected
func TestRenderComplete_OnlyAIAssistant(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenComplete
	m.Choices.OS = "darwin"
	m.Choices.AIAssistants = []string{"opencode"}
	m.SkippedSteps = make(map[Screen]bool)
	m.SkippedSteps[ScreenTerminalSelect] = true
	m.SkippedSteps[ScreenShellSelect] = true
	m.SkippedSteps[ScreenWMSelect] = true
	m.SkippedSteps[ScreenNvimSelect] = true
	m.SkippedSteps[ScreenAIAssistants] = false
	m.AIAssistantsList = GetAvailableAIAssistants()

	output := m.renderComplete()

	// Must show OS
	if !strings.Contains(output, "OS: darwin") {
		t.Error("Should show OS")
	}

	// Must show Terminal as skipped
	if !strings.Contains(output, "✗ Terminal (skipped)") {
		t.Error("Should show Terminal as skipped")
	}

	// Must show Shell as skipped
	if !strings.Contains(output, "✗ Shell (skipped)") {
		t.Error("Should show Shell as skipped")
	}

	// Must show Multiplexer as skipped
	if !strings.Contains(output, "✗ Multiplexer (skipped)") {
		t.Error("Should show Multiplexer as skipped")
	}

	// Must show Neovim as skipped
	if !strings.Contains(output, "✗ Neovim (skipped)") {
		t.Error("Should show Neovim as skipped")
	}

	// Must show AI Assistant with checkmark
	if !strings.Contains(output, "✓ AI Assistant: OpenCode") {
		t.Error("Should show AI Assistant: OpenCode with checkmark")
	}

	// Should NOT show exec command (shell was skipped)
	if strings.Contains(output, "exec") {
		t.Error("Should NOT show 'exec' command when shell is skipped")
	}

	// Should show generic message instead
	if !strings.Contains(output, "Your dotfiles have been configured") {
		t.Error("Should show generic message when shell is skipped")
	}
}

// TestRenderComplete_MultipleAIAssistants tests multiple AI assistants
func TestRenderComplete_MultipleAIAssistants(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenComplete
	m.Choices.OS = "darwin"
	m.Choices.Shell = "fish"
	m.Choices.AIAssistants = []string{"opencode", "continue"}
	m.SkippedSteps = make(map[Screen]bool)
	m.SkippedSteps[ScreenTerminalSelect] = true
	m.SkippedSteps[ScreenShellSelect] = false
	m.SkippedSteps[ScreenWMSelect] = true
	m.SkippedSteps[ScreenNvimSelect] = true
	m.SkippedSteps[ScreenAIAssistants] = false
	m.AIAssistantsList = GetAvailableAIAssistants()

	output := m.renderComplete()

	// Must show Shell
	if !strings.Contains(output, "✓ Shell: Fish") {
		t.Error("Should show Shell: Fish")
	}

	// Must show BOTH AI Assistants
	if !strings.Contains(output, "✓ AI Assistant: OpenCode") {
		t.Error("Should show AI Assistant: OpenCode")
	}
	if !strings.Contains(output, "✓ AI Assistant: Continue.dev") {
		t.Error("Should show AI Assistant: Continue.dev")
	}

	// Must show shell exec command
	if !strings.Contains(output, "exec fish") {
		t.Error("Should show 'exec fish' command")
	}
}

// TestRenderComplete_NushellCommand tests nushell "nu" command conversion
func TestRenderComplete_NushellCommand(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenComplete
	m.Choices.OS = "darwin"
	m.Choices.Shell = "nushell"
	m.SkippedSteps = make(map[Screen]bool)
	m.SkippedSteps[ScreenShellSelect] = false
	m.SkippedSteps[ScreenTerminalSelect] = true
	m.SkippedSteps[ScreenWMSelect] = true
	m.SkippedSteps[ScreenNvimSelect] = true
	m.SkippedSteps[ScreenAIAssistants] = true

	output := m.renderComplete()

	// Must show Nushell
	if !strings.Contains(output, "✓ Shell: Nushell") {
		t.Error("Should show Shell: Nushell")
	}

	// Must show "exec nu" (not "exec nushell")
	if !strings.Contains(output, "exec nu") {
		t.Error("Should show 'exec nu' command (not 'exec nushell')")
	}

	if strings.Contains(output, "exec nushell") {
		t.Error("Should NOT show 'exec nushell', should be 'exec nu'")
	}
}

// TestRenderComplete_NoAIAssistants tests when AI is skipped
func TestRenderComplete_NoAIAssistants(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenComplete
	m.Choices.OS = "darwin"
	m.Choices.Shell = "zsh"
	m.Choices.AIAssistants = []string{} // No AI selected
	m.SkippedSteps = make(map[Screen]bool)
	m.SkippedSteps[ScreenShellSelect] = false
	m.SkippedSteps[ScreenAIAssistants] = true // Skipped
	m.AIAssistantsList = GetAvailableAIAssistants()

	output := m.renderComplete()

	// Must show Shell
	if !strings.Contains(output, "✓ Shell: Zsh") {
		t.Error("Should show Shell: Zsh")
	}

	// Must show AI Assistants as skipped
	if !strings.Contains(output, "✗ AI Assistants (skipped)") {
		t.Error("Should show AI Assistants as skipped")
	}

	// Should NOT show any specific AI assistant
	if strings.Contains(output, "OpenCode") {
		t.Error("Should NOT show OpenCode when skipped")
	}
}

// TestRenderComplete_AllSkipped tests when everything is skipped
func TestRenderComplete_AllSkipped(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenComplete
	m.Choices.OS = "darwin"
	m.SkippedSteps = make(map[Screen]bool)
	m.SkippedSteps[ScreenTerminalSelect] = true
	m.SkippedSteps[ScreenShellSelect] = true
	m.SkippedSteps[ScreenWMSelect] = true
	m.SkippedSteps[ScreenNvimSelect] = true
	m.SkippedSteps[ScreenAIAssistants] = true

	output := m.renderComplete()

	// Must show all as skipped
	if !strings.Contains(output, "✗ Terminal (skipped)") {
		t.Error("Should show Terminal as skipped")
	}
	if !strings.Contains(output, "✗ Shell (skipped)") {
		t.Error("Should show Shell as skipped")
	}
	if !strings.Contains(output, "✗ Multiplexer (skipped)") {
		t.Error("Should show Multiplexer as skipped")
	}
	if !strings.Contains(output, "✗ Neovim (skipped)") {
		t.Error("Should show Neovim as skipped")
	}
	if !strings.Contains(output, "✗ AI Assistants (skipped)") {
		t.Error("Should show AI Assistants as skipped")
	}

	// Should NOT show exec command
	if strings.Contains(output, "exec") {
		t.Error("Should NOT show 'exec' command when shell is skipped")
	}

	// Should show generic message
	if !strings.Contains(output, "Your dotfiles have been configured") {
		t.Error("Should show generic message when everything is skipped")
	}
}

// TestRenderComplete_TerminalWithFont tests terminal with font
func TestRenderComplete_TerminalWithFont(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenComplete
	m.Choices.OS = "darwin"
	m.Choices.Terminal = "alacritty"
	m.Choices.InstallFont = true
	m.SkippedSteps = make(map[Screen]bool)
	m.SkippedSteps[ScreenShellSelect] = true
	m.SkippedSteps[ScreenWMSelect] = true
	m.SkippedSteps[ScreenNvimSelect] = true
	m.SkippedSteps[ScreenAIAssistants] = true

	output := m.renderComplete()

	// Must show Terminal
	if !strings.Contains(output, "✓ Terminal: Alacritty") {
		t.Error("Should show Terminal: Alacritty")
	}

	// Must show Font as sub-item (with indentation)
	if !strings.Contains(output, "└─ Iosevka Nerd Font") {
		t.Error("Should show font as sub-item with └─")
	}
}

// TestRenderComplete_TerminalWithoutFont tests terminal without font
func TestRenderComplete_TerminalWithoutFont(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenComplete
	m.Choices.OS = "darwin"
	m.Choices.Terminal = "wezterm"
	m.Choices.InstallFont = false
	m.SkippedSteps = make(map[Screen]bool)
	m.SkippedSteps[ScreenShellSelect] = true
	m.SkippedSteps[ScreenWMSelect] = true
	m.SkippedSteps[ScreenNvimSelect] = true
	m.SkippedSteps[ScreenAIAssistants] = true

	output := m.renderComplete()

	// Must show Terminal
	if !strings.Contains(output, "✓ Terminal: Wezterm") {
		t.Error("Should show Terminal: Wezterm")
	}

	// Should NOT show Font
	if strings.Contains(output, "Iosevka Nerd Font") {
		t.Error("Should NOT show font when InstallFont=false")
	}
}

// TestRenderComplete_MixedSelections tests a realistic mix of selections
func TestRenderComplete_MixedSelections(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenComplete
	m.Choices.OS = "linux"
	m.Choices.Terminal = "ghostty"
	m.Choices.InstallFont = true
	m.Choices.Shell = "fish"
	m.Choices.WindowMgr = "zellij"
	m.Choices.InstallNvim = false
	m.Choices.AIAssistants = []string{"opencode"}
	m.SkippedSteps = make(map[Screen]bool)
	m.SkippedSteps[ScreenNvimSelect] = false // User explicitly said "No"
	m.AIAssistantsList = GetAvailableAIAssistants()

	output := m.renderComplete()

	// Must show selected items
	if !strings.Contains(output, "✓ Terminal: Ghostty") {
		t.Error("Should show Terminal: Ghostty")
	}
	if !strings.Contains(output, "✓ Shell: Fish") {
		t.Error("Should show Shell: Fish")
	}
	if !strings.Contains(output, "✓ Multiplexer: Zellij") {
		t.Error("Should show Multiplexer: Zellij")
	}
	if !strings.Contains(output, "✓ AI Assistant: OpenCode") {
		t.Error("Should show AI Assistant: OpenCode")
	}

	// Should NOT show Neovim (user said No, not skipped)
	// GetInstallationSummary() doesn't add it if InstallNvim=false and not skipped
	if strings.Contains(output, "Neovim") {
		t.Error("Should NOT show Neovim when user explicitly said No")
	}

	// Must show exec fish command
	if !strings.Contains(output, "exec fish") {
		t.Error("Should show 'exec fish' command")
	}
}
