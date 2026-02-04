package tui

import (
	"strings"
	"testing"

	"github.com/Gentleman-Programming/Gentleman.Dots/installer/internal/system"
)

// TestInstallationSummaryWithNeovimAndOpenCode tests that AI assistants with Neovim appear correctly
func TestInstallationSummaryWithNeovimAndOpenCode(t *testing.T) {
	m := NewModel()
	m.SystemInfo = &system.SystemInfo{
		OS:       system.OSMac,
		IsTermux: false,
	}

	// Simulate user selections - Neovim + OpenCode
	// Note: Gemini and Copilot should NOT be in AIAssistants list because they're auto-installed with Neovim
	m.Choices.InstallNvim = true
	m.Choices.AIAssistants = []string{"opencode"}

	summary := m.GetInstallationSummary()

	// Count how many times each AI assistant appears
	claudeCodeCount := 0
	openCodeCount := 0
	geminiCount := 0
	copilotCount := 0
	for _, line := range summary {
		if strings.Contains(line, "Claude Code") {
			claudeCodeCount++
		}
		if strings.Contains(line, "OpenCode") {
			openCodeCount++
		}
		if strings.Contains(line, "Gemini CLI") {
			geminiCount++
		}
		if strings.Contains(line, "GitHub Copilot CLI") || strings.Contains(line, "Copilot CLI") {
			copilotCount++
		}
	}

	// Claude Code, Gemini, and Copilot should each appear exactly once (auto-installed with Neovim)
	if claudeCodeCount != 1 {
		t.Errorf("Claude Code should appear exactly once, but appeared %d times", claudeCodeCount)
		t.Logf("Summary:")
		for _, line := range summary {
			t.Logf("  %s", line)
		}
	}
	if geminiCount != 1 {
		t.Errorf("Gemini CLI should appear exactly once, but appeared %d times", geminiCount)
	}
	if copilotCount != 1 {
		t.Errorf("GitHub Copilot CLI should appear exactly once, but appeared %d times", copilotCount)
	}

	// OpenCode should appear once (explicitly selected)
	if openCodeCount != 1 {
		t.Errorf("OpenCode should appear exactly once, but appeared %d times", openCodeCount)
	}

	// Verify Claude Code, Gemini, and Copilot mention Neovim (auto-installed)
	for _, line := range summary {
		if strings.Contains(line, "Claude Code") || strings.Contains(line, "Gemini CLI") || strings.Contains(line, "Copilot CLI") {
			if !strings.Contains(line, "Neovim") {
				t.Errorf("AI assistant line should mention Neovim: %s", line)
			}
		}
	}
}

// TestInstallationSummaryWithClaudeCodeOnly tests Claude Code without Neovim
func TestInstallationSummaryWithClaudeCodeOnly(t *testing.T) {
	m := NewModel()
	m.SystemInfo = &system.SystemInfo{
		OS:       system.OSMac,
		IsTermux: false,
	}

	// Simulate user selections - NO Neovim, but Claude Code selected
	m.Choices.InstallNvim = false
	m.Choices.AIAssistants = []string{"claudecode"}

	summary := m.GetInstallationSummary()

	// Count how many times "Claude Code" appears
	claudeCodeCount := 0
	for _, line := range summary {
		if strings.Contains(line, "Claude Code") {
			claudeCodeCount++
		}
	}

	// Should appear exactly ONCE (as AI Assistant)
	if claudeCodeCount != 1 {
		t.Errorf("Claude Code should appear exactly once, but appeared %d times", claudeCodeCount)
		t.Logf("Summary:")
		for _, line := range summary {
			t.Logf("  %s", line)
		}
	}

	// Verify Claude Code does NOT mention Neovim
	for _, line := range summary {
		if strings.Contains(line, "Claude Code") && strings.Contains(line, "Neovim") {
			t.Error("Claude Code should NOT mention Neovim when Neovim is not installed")
		}
	}
}

// TestInstallationSummaryNoClaudeCode tests when neither Neovim nor Claude Code selected
func TestInstallationSummaryNoClaudeCode(t *testing.T) {
	m := NewModel()
	m.SystemInfo = &system.SystemInfo{
		OS:       system.OSMac,
		IsTermux: false,
	}

	// Simulate user selections - NO Neovim, NO Claude Code
	m.Choices.InstallNvim = false
	m.Choices.AIAssistants = []string{"opencode"}

	summary := m.GetInstallationSummary()

	// Claude Code should NOT appear at all
	for _, line := range summary {
		if strings.Contains(line, "Claude Code") {
			t.Errorf("Claude Code should not appear in summary, but found: %s", line)
		}
	}
}
