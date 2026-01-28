package tui

import (
	"strings"
	"testing"

	"github.com/Gentleman-Programming/Gentleman.Dots/installer/internal/system"
)

// TestInstallationSummaryWithNeovimAndOpenCode tests that Claude Code only appears once
func TestInstallationSummaryWithNeovimAndOpenCode(t *testing.T) {
	m := NewModel()
	m.SystemInfo = &system.SystemInfo{
		OS:       system.OSMac,
		IsTermux: false,
	}

	// Simulate user selections
	m.Choices.InstallNvim = true
	m.Choices.AIAssistants = []string{"opencode"} // Only OpenCode selected

	summary := m.GetInstallationSummary()

	// Count how many times "Claude Code" appears
	claudeCodeCount := 0
	for _, line := range summary {
		if strings.Contains(line, "Claude Code") {
			claudeCodeCount++
		}
	}

	// Should appear exactly ONCE (with Neovim)
	if claudeCodeCount != 1 {
		t.Errorf("Claude Code should appear exactly once, but appeared %d times", claudeCodeCount)
		t.Logf("Summary:")
		for _, line := range summary {
			t.Logf("  %s", line)
		}
	}

	// Verify Claude Code line mentions Neovim
	foundCorrectLine := false
	for _, line := range summary {
		if strings.Contains(line, "Claude Code") && strings.Contains(line, "Neovim") {
			foundCorrectLine = true
			break
		}
	}
	if !foundCorrectLine {
		t.Error("Claude Code should be shown as '(with Neovim)'")
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
