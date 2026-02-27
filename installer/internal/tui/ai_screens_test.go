package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// --- Screen Options Tests ---

func TestAIToolsSelectOptions(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIToolsSelect
	opts := m.GetCurrentOptions()

	// 4 tools + separator + confirm = 6
	if len(opts) != 6 {
		t.Fatalf("Expected 6 AI tools options (4 tools + separator + confirm), got %d: %v", len(opts), opts)
	}
	if opts[0] != "Claude Code" {
		t.Errorf("Expected first option 'Claude Code', got %s", opts[0])
	}
	if opts[1] != "OpenCode" {
		t.Errorf("Expected second option 'OpenCode', got %s", opts[1])
	}
	if opts[2] != "Gemini CLI" {
		t.Errorf("Expected third option 'Gemini CLI', got %s", opts[2])
	}
	if opts[3] != "GitHub Copilot" {
		t.Errorf("Expected fourth option 'GitHub Copilot', got %s", opts[3])
	}
	if !strings.HasPrefix(opts[4], "───") {
		t.Errorf("Expected separator at index 4, got %s", opts[4])
	}
	if !strings.Contains(opts[5], "Confirm") {
		t.Errorf("Expected last option to be Confirm, got %s", opts[5])
	}
}

func TestAIFrameworkConfirmOptions(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIFrameworkConfirm
	opts := m.GetCurrentOptions()

	if len(opts) != 2 {
		t.Fatalf("Expected 2 framework confirm options, got %d", len(opts))
	}
}

func TestAIFrameworkPresetOptions(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIFrameworkPreset
	opts := m.GetCurrentOptions()

	// 6 presets + separator + custom = 8
	if len(opts) != 8 {
		t.Fatalf("Expected 8 preset options (6 presets + separator + custom), got %d", len(opts))
	}
	if !strings.Contains(opts[0], "Minimal") {
		t.Errorf("Expected first preset to be Minimal, got %s", opts[0])
	}
	if !strings.Contains(opts[7], "Custom") {
		t.Errorf("Expected last option to be Custom, got %s", opts[7])
	}
}

func TestAIFrameworkModulesOptions(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIFrameworkModules
	opts := m.GetCurrentOptions()

	// 27 modules + separator + confirm = 29
	if len(opts) != 29 {
		t.Fatalf("Expected 29 module options (27 + separator + confirm), got %d", len(opts))
	}
	if !strings.Contains(opts[len(opts)-1], "Confirm") {
		t.Errorf("Expected last option to be Confirm, got %s", opts[len(opts)-1])
	}
}

// --- Screen Titles Tests ---

func TestAIScreenTitles(t *testing.T) {
	m := NewModel()

	tests := []struct {
		screen Screen
		expect string
	}{
		{ScreenAIToolsSelect, "Step 7"},
		{ScreenAIFrameworkConfirm, "Step 8"},
		{ScreenAIFrameworkPreset, "Step 8"},
		{ScreenAIFrameworkModules, "Step 8"},
	}

	for _, tt := range tests {
		m.Screen = tt.screen
		title := m.GetScreenTitle()
		if !strings.Contains(title, tt.expect) {
			t.Errorf("Screen %v: expected title containing %q, got %q", tt.screen, tt.expect, title)
		}
	}
}

// --- Screen Flow Tests ---

func TestNvimSelectToAIToolsTransition(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenNvimSelect
	m.Cursor = 0 // Yes, install Neovim

	result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	newModel := result.(Model)

	if !newModel.Choices.InstallNvim {
		t.Error("Expected InstallNvim to be true")
	}
	if newModel.Screen != ScreenAIToolsSelect {
		t.Errorf("Expected ScreenAIToolsSelect, got %v", newModel.Screen)
	}
	if newModel.AIToolSelected == nil {
		t.Error("Expected AIToolSelected to be initialized on transition")
	}
	if len(newModel.AIToolSelected) != len(aiToolIDMap) {
		t.Errorf("Expected AIToolSelected length %d, got %d", len(aiToolIDMap), len(newModel.AIToolSelected))
	}
}

func TestNvimSelectSkipToAITools(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenNvimSelect
	m.Cursor = 1 // No, skip Neovim

	result, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	newModel := result.(Model)

	if newModel.Choices.InstallNvim {
		t.Error("Expected InstallNvim to be false")
	}
	if newModel.Screen != ScreenAIToolsSelect {
		t.Errorf("Expected ScreenAIToolsSelect even when skipping nvim, got %v", newModel.Screen)
	}
}

// --- AI Tools Multi-Select Tests ---

func TestAIToolsToggle(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIToolsSelect
	m.AIToolSelected = make([]bool, len(aiToolIDMap))
	m.Cursor = 0 // Claude Code

	// Toggle on
	result, _ := m.handleAIToolsKeys("enter")
	newModel := result.(Model)
	if !newModel.AIToolSelected[0] {
		t.Error("Expected tool 0 (Claude Code) to be toggled ON")
	}

	// Toggle off
	result, _ = newModel.handleAIToolsKeys("enter")
	newModel = result.(Model)
	if newModel.AIToolSelected[0] {
		t.Error("Expected tool 0 (Claude Code) to be toggled OFF")
	}
}

func TestAIToolsSelectAllTools(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIToolsSelect
	m.AIToolSelected = make([]bool, len(aiToolIDMap))
	// Select all tools
	for i := range m.AIToolSelected {
		m.AIToolSelected[i] = true
	}

	opts := m.GetCurrentOptions()
	m.Cursor = len(opts) - 1 // "Confirm selection"

	result, _ := m.handleAIToolsKeys("enter")
	newModel := result.(Model)

	if len(newModel.Choices.AITools) != 4 {
		t.Fatalf("Expected 4 AI tools, got %d: %v", len(newModel.Choices.AITools), newModel.Choices.AITools)
	}
	if newModel.Screen != ScreenAIFrameworkConfirm {
		t.Errorf("Expected ScreenAIFrameworkConfirm, got %v", newModel.Screen)
	}
}

func TestAIToolsSelectSingleTool(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIToolsSelect
	m.AIToolSelected = make([]bool, len(aiToolIDMap))
	m.AIToolSelected[0] = true // Claude Code only

	opts := m.GetCurrentOptions()
	m.Cursor = len(opts) - 1 // "Confirm selection"

	result, _ := m.handleAIToolsKeys("enter")
	newModel := result.(Model)

	if len(newModel.Choices.AITools) != 1 || newModel.Choices.AITools[0] != "claude" {
		t.Errorf("Expected [claude], got %v", newModel.Choices.AITools)
	}
	if newModel.Screen != ScreenAIFrameworkConfirm {
		t.Errorf("Expected ScreenAIFrameworkConfirm, got %v", newModel.Screen)
	}
}

func TestAIToolsSelectGeminiAndCopilot(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIToolsSelect
	m.AIToolSelected = make([]bool, len(aiToolIDMap))
	m.AIToolSelected[2] = true // Gemini CLI
	m.AIToolSelected[3] = true // GitHub Copilot

	opts := m.GetCurrentOptions()
	m.Cursor = len(opts) - 1 // "Confirm selection"

	result, _ := m.handleAIToolsKeys("enter")
	newModel := result.(Model)

	if len(newModel.Choices.AITools) != 2 {
		t.Fatalf("Expected 2 AI tools, got %d: %v", len(newModel.Choices.AITools), newModel.Choices.AITools)
	}
	if newModel.Choices.AITools[0] != "gemini" {
		t.Errorf("Expected first tool 'gemini', got %s", newModel.Choices.AITools[0])
	}
	if newModel.Choices.AITools[1] != "copilot" {
		t.Errorf("Expected second tool 'copilot', got %s", newModel.Choices.AITools[1])
	}
	if newModel.Screen != ScreenAIFrameworkConfirm {
		t.Errorf("Expected ScreenAIFrameworkConfirm, got %v", newModel.Screen)
	}
}

func TestAIToolsSelectNone(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIToolsSelect
	m.AIToolSelected = make([]bool, len(aiToolIDMap))
	// No tools toggled

	opts := m.GetCurrentOptions()
	m.Cursor = len(opts) - 1 // "Confirm selection"

	result, _ := m.handleAIToolsKeys("enter")
	newModel := result.(Model)

	if len(newModel.Choices.AITools) != 0 {
		t.Errorf("Expected no AI tools, got %v", newModel.Choices.AITools)
	}
	// Should skip framework and go to backup/install
	if newModel.Screen != ScreenBackupConfirm && newModel.Screen != ScreenInstalling {
		t.Errorf("Expected ScreenBackupConfirm or ScreenInstalling, got %v", newModel.Screen)
	}
}

func TestAIFrameworkConfirmYes(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIFrameworkConfirm
	m.Choices.AITools = []string{"claude"}
	m.Cursor = 0 // Yes

	result, _ := m.handleSelection()
	newModel := result.(Model)

	if !newModel.Choices.InstallAIFramework {
		t.Error("Expected InstallAIFramework to be true")
	}
	if newModel.Screen != ScreenAIFrameworkPreset {
		t.Errorf("Expected ScreenAIFrameworkPreset, got %v", newModel.Screen)
	}
}

func TestAIFrameworkConfirmNo(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIFrameworkConfirm
	m.Choices.AITools = []string{"claude"}
	m.Cursor = 1 // No

	result, _ := m.handleSelection()
	newModel := result.(Model)

	if newModel.Choices.InstallAIFramework {
		t.Error("Expected InstallAIFramework to be false")
	}
	// Should go to backup/install
	if newModel.Screen != ScreenBackupConfirm && newModel.Screen != ScreenInstalling {
		t.Errorf("Expected ScreenBackupConfirm or ScreenInstalling, got %v", newModel.Screen)
	}
}

func TestAIFrameworkPresetSelection(t *testing.T) {
	presets := []string{"minimal", "frontend", "backend", "fullstack", "data", "complete"}
	for i, preset := range presets {
		m := NewModel()
		m.Screen = ScreenAIFrameworkPreset
		m.Choices.AITools = []string{"claude"}
		m.Choices.InstallAIFramework = true
		m.Cursor = i

		result, _ := m.handleSelection()
		newModel := result.(Model)

		if newModel.Choices.AIFrameworkPreset != preset {
			t.Errorf("Cursor %d: expected preset %q, got %q", i, preset, newModel.Choices.AIFrameworkPreset)
		}
		// Should proceed to backup/install
		if newModel.Screen != ScreenBackupConfirm && newModel.Screen != ScreenInstalling {
			t.Errorf("Preset %s: expected ScreenBackupConfirm or ScreenInstalling, got %v", preset, newModel.Screen)
		}
	}
}

func TestAIFrameworkPresetCustom(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIFrameworkPreset
	m.Choices.AITools = []string{"claude"}
	m.Choices.InstallAIFramework = true
	m.Cursor = 7 // Custom

	result, _ := m.handleSelection()
	newModel := result.(Model)

	if newModel.Screen != ScreenAIFrameworkModules {
		t.Errorf("Expected ScreenAIFrameworkModules, got %v", newModel.Screen)
	}
	if newModel.AIModuleSelected == nil {
		t.Error("Expected AIModuleSelected to be initialized")
	}
	if len(newModel.AIModuleSelected) != 27 {
		t.Errorf("Expected 27 module toggles, got %d", len(newModel.AIModuleSelected))
	}
}

// --- Module Multi-Select Tests ---

func TestAIModulesToggle(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIFrameworkModules
	m.AIModuleSelected = make([]bool, 27)
	m.Cursor = 0 // First module

	// Toggle on
	result, _ := m.handleAIModulesKeys("enter")
	newModel := result.(Model)
	if !newModel.AIModuleSelected[0] {
		t.Error("Expected module 0 to be toggled ON")
	}

	// Toggle off
	result, _ = newModel.handleAIModulesKeys("enter")
	newModel = result.(Model)
	if newModel.AIModuleSelected[0] {
		t.Error("Expected module 0 to be toggled OFF")
	}
}

func TestAIModulesConfirmWithSelection(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIFrameworkModules
	m.AIModuleSelected = make([]bool, 27)
	m.AIModuleSelected[0] = true // scripts-project
	m.AIModuleSelected[4] = true // hooks-security
	m.Choices.AITools = []string{"claude"}
	m.Choices.InstallAIFramework = true

	opts := m.GetCurrentOptions()
	m.Cursor = len(opts) - 1 // "Confirm selection"

	result, _ := m.handleAIModulesKeys("enter")
	newModel := result.(Model)

	if len(newModel.Choices.AIFrameworkModules) != 2 {
		t.Fatalf("Expected 2 selected modules, got %d: %v", len(newModel.Choices.AIFrameworkModules), newModel.Choices.AIFrameworkModules)
	}
	if newModel.Choices.AIFrameworkModules[0] != "scripts-project" {
		t.Errorf("Expected first module 'scripts-project', got %s", newModel.Choices.AIFrameworkModules[0])
	}
	if newModel.Choices.AIFrameworkModules[1] != "hooks-security" {
		t.Errorf("Expected second module 'hooks-security', got %s", newModel.Choices.AIFrameworkModules[1])
	}
}

func TestAIModulesConfirmEmpty(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIFrameworkModules
	m.AIModuleSelected = make([]bool, 27)
	m.Choices.AITools = []string{"claude"}
	m.Choices.InstallAIFramework = true

	opts := m.GetCurrentOptions()
	m.Cursor = len(opts) - 1 // "Confirm selection"

	result, _ := m.handleAIModulesKeys("enter")
	newModel := result.(Model)

	// No modules selected = skip framework
	if newModel.Choices.InstallAIFramework {
		t.Error("Expected InstallAIFramework to be false when no modules selected")
	}
}

// --- Back Navigation Tests ---

func TestAIToolsSelectGoBack(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIToolsSelect
	m.AIToolSelected = make([]bool, len(aiToolIDMap))

	result, _ := m.goBackInstallStep()
	newModel := result.(Model)

	if newModel.Screen != ScreenNvimSelect {
		t.Errorf("Expected ScreenNvimSelect, got %v", newModel.Screen)
	}
	if newModel.AIToolSelected != nil {
		t.Error("Expected AIToolSelected to be cleared on back")
	}
}

func TestAIFrameworkConfirmGoBack(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIFrameworkConfirm

	result, _ := m.goBackInstallStep()
	newModel := result.(Model)

	if newModel.Screen != ScreenAIToolsSelect {
		t.Errorf("Expected ScreenAIToolsSelect, got %v", newModel.Screen)
	}
}

func TestAIFrameworkPresetGoBack(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIFrameworkPreset

	result, _ := m.goBackInstallStep()
	newModel := result.(Model)

	if newModel.Screen != ScreenAIFrameworkConfirm {
		t.Errorf("Expected ScreenAIFrameworkConfirm, got %v", newModel.Screen)
	}
}

func TestAIFrameworkModulesGoBack(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIFrameworkModules
	m.AIModuleSelected = make([]bool, 27)

	result, _ := m.goBackInstallStep()
	newModel := result.(Model)

	if newModel.Screen != ScreenAIFrameworkPreset {
		t.Errorf("Expected ScreenAIFrameworkPreset, got %v", newModel.Screen)
	}
}

// --- SetupInstallSteps Tests ---

func TestSetupInstallStepsWithAITools(t *testing.T) {
	m := NewModel()
	m.Choices.Shell = "fish"
	m.Choices.AITools = []string{"claude", "opencode"}

	m.SetupInstallSteps()

	hasAIToolsStep := false
	for _, step := range m.Steps {
		if step.ID == "aitools" {
			hasAIToolsStep = true
			if !strings.Contains(step.Description, "claude") {
				t.Errorf("Expected description to mention claude, got %s", step.Description)
			}
		}
	}
	if !hasAIToolsStep {
		t.Error("Expected aitools step when AITools are selected")
	}
}

func TestSetupInstallStepsWithoutAITools(t *testing.T) {
	m := NewModel()
	m.Choices.Shell = "fish"
	m.Choices.AITools = nil

	m.SetupInstallSteps()

	for _, step := range m.Steps {
		if step.ID == "aitools" {
			t.Error("Expected NO aitools step when AITools is nil")
		}
	}
}

func TestSetupInstallStepsWithAIFramework(t *testing.T) {
	m := NewModel()
	m.Choices.Shell = "fish"
	m.Choices.AITools = []string{"claude"}
	m.Choices.InstallAIFramework = true
	m.Choices.AIFrameworkPreset = "frontend"

	m.SetupInstallSteps()

	hasFrameworkStep := false
	for _, step := range m.Steps {
		if step.ID == "aiframework" {
			hasFrameworkStep = true
			if !strings.Contains(step.Description, "frontend") {
				t.Errorf("Expected description to mention preset, got %s", step.Description)
			}
		}
	}
	if !hasFrameworkStep {
		t.Error("Expected aiframework step when InstallAIFramework is true")
	}
}

func TestSetupInstallStepsWithoutAIFramework(t *testing.T) {
	m := NewModel()
	m.Choices.Shell = "fish"
	m.Choices.AITools = []string{"claude"}
	m.Choices.InstallAIFramework = false

	m.SetupInstallSteps()

	for _, step := range m.Steps {
		if step.ID == "aiframework" {
			t.Error("Expected NO aiframework step when InstallAIFramework is false")
		}
	}
}

// --- ID Map Tests ---

func TestAIToolIDMapLength(t *testing.T) {
	if len(aiToolIDMap) != 4 {
		t.Errorf("Expected 4 tool IDs in aiToolIDMap, got %d", len(aiToolIDMap))
	}
}

func TestModuleIDMapLength(t *testing.T) {
	if len(moduleIDMap) != 27 {
		t.Errorf("Expected 27 module IDs in moduleIDMap, got %d", len(moduleIDMap))
	}
}

func TestModuleIDMapContainsExpected(t *testing.T) {
	expected := []string{"scripts-project", "hooks-security", "agents-development", "skills-frontend", "commands-git", "sdd", "mcp"}
	for _, e := range expected {
		found := false
		for _, id := range moduleIDMap {
			if id == e {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected moduleIDMap to contain %q", e)
		}
	}
}

// --- hasAITool Helper Tests ---

func TestHasAITool(t *testing.T) {
	tools := []string{"claude", "opencode", "gemini", "copilot"}
	if !hasAITool(tools, "claude") {
		t.Error("Expected hasAITool to find 'claude'")
	}
	if !hasAITool(tools, "opencode") {
		t.Error("Expected hasAITool to find 'opencode'")
	}
	if !hasAITool(tools, "gemini") {
		t.Error("Expected hasAITool to find 'gemini'")
	}
	if !hasAITool(tools, "copilot") {
		t.Error("Expected hasAITool to find 'copilot'")
	}
	if hasAITool(tools, "cursor") {
		t.Error("Expected hasAITool NOT to find 'cursor'")
	}
	if hasAITool(nil, "claude") {
		t.Error("Expected hasAITool to return false for nil slice")
	}
}

// --- Progress Bar Tests ---

func TestProgressBarIncludesAISteps(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIToolsSelect
	progress := m.renderStepProgress()

	if !strings.Contains(progress, "AI Tools") {
		t.Error("Expected progress bar to contain 'AI Tools'")
	}
	if !strings.Contains(progress, "Framework") {
		t.Error("Expected progress bar to contain 'Framework'")
	}
}

// --- View Render Tests ---

func TestRenderAIToolSelectionShowsCheckboxes(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIToolsSelect
	m.Width = 100
	m.Height = 40
	m.AIToolSelected = make([]bool, len(aiToolIDMap))
	m.AIToolSelected[0] = true // Claude Code checked
	m.Cursor = 0

	rendered := m.renderAIToolSelection()

	if !strings.Contains(rendered, "[✓]") {
		t.Error("Expected rendered view to contain checked checkbox [✓]")
	}
	if !strings.Contains(rendered, "[ ]") {
		t.Error("Expected rendered view to contain unchecked checkbox [ ]")
	}
}

func TestRenderAIModuleSelectionShowsCheckboxes(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIFrameworkModules
	m.Width = 100
	m.Height = 40
	m.AIModuleSelected = make([]bool, 27)
	m.AIModuleSelected[0] = true
	m.Cursor = 0

	rendered := m.renderAIModuleSelection()

	if !strings.Contains(rendered, "[✓]") {
		t.Error("Expected rendered view to contain checked checkbox [✓]")
	}
	if !strings.Contains(rendered, "[ ]") {
		t.Error("Expected rendered view to contain unchecked checkbox [ ]")
	}
}
