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

func TestAIFrameworkPresetCustomIsFirst(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIFrameworkPreset
	opts := m.GetCurrentOptions()

	// Custom first, separator, then 6 presets = 8
	if len(opts) != 8 {
		t.Fatalf("Expected 8 preset options (custom + separator + 6 presets), got %d", len(opts))
	}
	if !strings.Contains(opts[0], "Custom") {
		t.Errorf("Expected first option to be Custom, got %s", opts[0])
	}
	if !strings.HasPrefix(opts[1], "───") {
		t.Errorf("Expected separator at index 1, got %s", opts[1])
	}
	if !strings.Contains(opts[2], "Minimal") {
		t.Errorf("Expected index 2 to be Minimal, got %s", opts[2])
	}
}

func TestAIFrameworkPresetSelection(t *testing.T) {
	// Presets are at indices 2-7 (after Custom at 0 and separator at 1)
	presets := []string{"minimal", "frontend", "backend", "fullstack", "data", "complete"}
	for i, preset := range presets {
		m := NewModel()
		m.Screen = ScreenAIFrameworkPreset
		m.Choices.AITools = []string{"claude"}
		m.Choices.InstallAIFramework = true
		m.Cursor = i + 2 // presets start at index 2

		result, _ := m.handleSelection()
		newModel := result.(Model)

		if newModel.Choices.AIFrameworkPreset != preset {
			t.Errorf("Cursor %d: expected preset %q, got %q", i+2, preset, newModel.Choices.AIFrameworkPreset)
		}
		// Should proceed to backup/install
		if newModel.Screen != ScreenBackupConfirm && newModel.Screen != ScreenInstalling {
			t.Errorf("Preset %s: expected ScreenBackupConfirm or ScreenInstalling, got %v", preset, newModel.Screen)
		}
	}
}

// --- Module Categories Tests ---

func TestModuleCategoriesCount(t *testing.T) {
	if len(moduleCategories) != 6 {
		t.Errorf("Expected 6 module categories, got %d", len(moduleCategories))
	}
}

func TestModuleCategoriesDataIntegrity(t *testing.T) {
	seen := make(map[string]bool)
	for _, cat := range moduleCategories {
		if cat.ID == "" {
			t.Error("Category has empty ID")
		}
		if cat.Label == "" {
			t.Errorf("Category %s has empty Label", cat.ID)
		}
		if cat.Icon == "" {
			t.Errorf("Category %s has empty Icon", cat.ID)
		}
		if len(cat.Items) == 0 {
			t.Errorf("Category %s has no items", cat.ID)
		}
		for _, item := range cat.Items {
			if item.ID == "" {
				t.Errorf("Category %s has item with empty ID", cat.ID)
			}
			if item.Label == "" {
				t.Errorf("Category %s item %s has empty Label", cat.ID, item.ID)
			}
			if seen[item.ID] {
				t.Errorf("Duplicate item ID across categories: %s", item.ID)
			}
			seen[item.ID] = true
		}
	}
}

func TestModuleCategoriesItemCounts(t *testing.T) {
	expected := map[string]int{
		"hooks":    10,
		"commands": 20,
		"agents":   80,
		"skills":   85,
		"sdd":      2,
		"mcp":      6,
	}
	for _, cat := range moduleCategories {
		exp, ok := expected[cat.ID]
		if !ok {
			t.Errorf("Unexpected category %s", cat.ID)
			continue
		}
		if len(cat.Items) != exp {
			t.Errorf("Category %s: expected %d items, got %d", cat.ID, exp, len(cat.Items))
		}
	}
}

func TestModuleCategoriesAtomicFlag(t *testing.T) {
	for _, cat := range moduleCategories {
		switch cat.ID {
		case "mcp":
			if !cat.IsAtomic {
				t.Errorf("Category %s should be atomic", cat.ID)
			}
		default:
			if cat.IsAtomic {
				t.Errorf("Category %s should NOT be atomic", cat.ID)
			}
		}
	}
}

// --- Category Menu Tests ---

func TestAICategoryMenuOptions(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIFrameworkCategories
	m.AICategorySelected = make(map[string][]bool)
	for _, cat := range moduleCategories {
		m.AICategorySelected[cat.ID] = make([]bool, len(cat.Items))
	}

	opts := m.GetCurrentOptions()
	// 6 categories + separator + confirm = 8
	if len(opts) != 8 {
		t.Fatalf("Expected 8 category options (6 + separator + confirm), got %d: %v", len(opts), opts)
	}
	if !strings.Contains(opts[len(opts)-1], "Confirm") {
		t.Errorf("Expected last option to be Confirm, got %s", opts[len(opts)-1])
	}
}

func TestAICategoryMenuShowsCounts(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIFrameworkCategories
	m.AICategorySelected = make(map[string][]bool)
	for _, cat := range moduleCategories {
		m.AICategorySelected[cat.ID] = make([]bool, len(cat.Items))
	}
	// Select 2 out of 10 hooks
	m.AICategorySelected["hooks"][0] = true
	m.AICategorySelected["hooks"][2] = true

	opts := m.GetCurrentOptions()
	// First option should be Hooks with (2/10 selected)
	if !strings.Contains(opts[0], "(2/10 selected)") {
		t.Errorf("Expected Hooks to show (2/10 selected), got %s", opts[0])
	}
	// Commands should show (0/20 selected)
	if !strings.Contains(opts[1], "(0/20 selected)") {
		t.Errorf("Expected Commands to show (0/20 selected), got %s", opts[1])
	}
}

// --- Category Drill-Down Tests ---

func TestAICategoryDrillDown(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIFrameworkCategories
	m.AICategorySelected = make(map[string][]bool)
	for _, cat := range moduleCategories {
		m.AICategorySelected[cat.ID] = make([]bool, len(cat.Items))
	}
	m.Cursor = 0 // Hooks category (first)

	result, _ := m.handleAICategoriesKeys("enter")
	newModel := result.(Model)

	if newModel.Screen != ScreenAIFrameworkCategoryItems {
		t.Errorf("Expected ScreenAIFrameworkCategoryItems, got %v", newModel.Screen)
	}
	if newModel.SelectedModuleCategory != 0 {
		t.Errorf("Expected SelectedModuleCategory 0, got %d", newModel.SelectedModuleCategory)
	}
	if newModel.Cursor != 0 {
		t.Errorf("Expected cursor reset to 0, got %d", newModel.Cursor)
	}
}

func TestAICategoryItemsToggle(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIFrameworkCategoryItems
	m.SelectedModuleCategory = 0 // Hooks (first category)
	m.AICategorySelected = make(map[string][]bool)
	for _, cat := range moduleCategories {
		m.AICategorySelected[cat.ID] = make([]bool, len(cat.Items))
	}
	m.Cursor = 0 // First item

	// Toggle on
	result, _ := m.handleAICategoryItemsKeys("enter")
	newModel := result.(Model)
	if !newModel.AICategorySelected["hooks"][0] {
		t.Error("Expected first hooks item to be toggled ON")
	}

	// Toggle off
	result, _ = newModel.handleAICategoryItemsKeys("enter")
	newModel = result.(Model)
	if newModel.AICategorySelected["hooks"][0] {
		t.Error("Expected first hooks item to be toggled OFF")
	}
}

func TestAICategoryItemsBack(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIFrameworkCategoryItems
	m.SelectedModuleCategory = 2 // Agents
	m.AICategorySelected = make(map[string][]bool)
	for _, cat := range moduleCategories {
		m.AICategorySelected[cat.ID] = make([]bool, len(cat.Items))
	}

	result, _ := m.handleAICategoryItemsKeys("esc")
	newModel := result.(Model)

	if newModel.Screen != ScreenAIFrameworkCategories {
		t.Errorf("Expected ScreenAIFrameworkCategories, got %v", newModel.Screen)
	}
	// Cursor should be preserved to the category we came from
	if newModel.Cursor != 2 {
		t.Errorf("Expected cursor preserved at 2 (Agents), got %d", newModel.Cursor)
	}
}

func TestAICategoryItemsBackButton(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIFrameworkCategoryItems
	m.SelectedModuleCategory = 1 // Hooks
	m.AICategorySelected = make(map[string][]bool)
	for _, cat := range moduleCategories {
		m.AICategorySelected[cat.ID] = make([]bool, len(cat.Items))
	}

	opts := m.GetCurrentOptions()
	m.Cursor = len(opts) - 1 // "← Back" button

	result, _ := m.handleAICategoryItemsKeys("enter")
	newModel := result.(Model)

	if newModel.Screen != ScreenAIFrameworkCategories {
		t.Errorf("Expected ScreenAIFrameworkCategories, got %v", newModel.Screen)
	}
	if newModel.Cursor != 1 {
		t.Errorf("Expected cursor preserved at 1 (Hooks), got %d", newModel.Cursor)
	}
}

// --- collectSelectedFeatures Tests ---

func TestCollectSelectedFeaturesNormal(t *testing.T) {
	sel := make(map[string][]bool)
	for _, cat := range moduleCategories {
		sel[cat.ID] = make([]bool, len(cat.Items))
	}
	sel["hooks"][0] = true    // any hook selected → "hooks" feature
	sel["commands"][0] = true // any command selected → "commands" feature
	sel["skills"][5] = true   // any skill selected → "skills" feature

	result := collectSelectedFeatures(sel)

	// Should produce feature-level IDs, one per category with selections
	expected := []string{"hooks", "commands", "skills"}
	if len(result) != len(expected) {
		t.Fatalf("Expected %d features, got %d: %v", len(expected), len(result), result)
	}
	for i, exp := range expected {
		if result[i] != exp {
			t.Errorf("Feature %d: expected %q, got %q", i, exp, result[i])
		}
	}
}

func TestCollectSelectedFeaturesAtomic(t *testing.T) {
	sel := make(map[string][]bool)
	for _, cat := range moduleCategories {
		sel[cat.ID] = make([]bool, len(cat.Items))
	}
	// Select OpenSpec (index 0) in SDD — should produce "sdd"
	sel["sdd"][0] = true
	// Select some MCP sub-items — should produce "mcp"
	sel["mcp"][1] = true

	result := collectSelectedFeatures(sel)

	if len(result) != 2 {
		t.Fatalf("Expected 2 features (sdd, mcp), got %d: %v", len(result), result)
	}
	if result[0] != "sdd" {
		t.Errorf("Expected first feature 'sdd', got %s", result[0])
	}
	if result[1] != "mcp" {
		t.Errorf("Expected second feature 'mcp', got %s", result[1])
	}
}

func TestCollectSelectedFeaturesSDDAgentTeamsOnly(t *testing.T) {
	sel := make(map[string][]bool)
	for _, cat := range moduleCategories {
		sel[cat.ID] = make([]bool, len(cat.Items))
	}
	// Select ONLY Agent Teams Lite (index 1) — should NOT produce "sdd" feature
	sel["sdd"][1] = true

	result := collectSelectedFeatures(sel)

	for _, f := range result {
		if f == "sdd" {
			t.Error("Selecting only Agent Teams Lite should NOT produce 'sdd' feature")
		}
	}
}

func TestCollectSelectedFeaturesSDDBoth(t *testing.T) {
	sel := make(map[string][]bool)
	for _, cat := range moduleCategories {
		sel[cat.ID] = make([]bool, len(cat.Items))
	}
	// Select BOTH OpenSpec and Agent Teams Lite
	sel["sdd"][0] = true // OpenSpec → "sdd"
	sel["sdd"][1] = true // Agent Teams Lite → no extra feature

	result := collectSelectedFeatures(sel)

	sddCount := 0
	for _, f := range result {
		if f == "sdd" {
			sddCount++
		}
	}
	if sddCount != 1 {
		t.Errorf("Expected exactly 1 'sdd' feature, got %d in %v", sddCount, result)
	}
}

func TestCollectSelectedFeaturesEmpty(t *testing.T) {
	sel := make(map[string][]bool)
	for _, cat := range moduleCategories {
		sel[cat.ID] = make([]bool, len(cat.Items))
	}
	// Nothing selected

	result := collectSelectedFeatures(sel)

	if len(result) != 0 {
		t.Errorf("Expected empty slice, got %v", result)
	}
}

func TestCollectSelectedFeaturesMixed(t *testing.T) {
	sel := make(map[string][]bool)
	for _, cat := range moduleCategories {
		sel[cat.ID] = make([]bool, len(cat.Items))
	}
	sel["hooks"][0] = true   // hooks feature
	sel["agents"][5] = true  // agents feature
	sel["sdd"][0] = true     // sdd feature
	sel["mcp"][2] = true     // mcp feature

	result := collectSelectedFeatures(sel)

	if len(result) != 4 {
		t.Fatalf("Expected 4 features, got %d: %v", len(result), result)
	}
	// Order follows moduleCategories order: hooks, agents, sdd, mcp
	expected := []string{"hooks", "agents", "sdd", "mcp"}
	for i, exp := range expected {
		if result[i] != exp {
			t.Errorf("Feature %d: expected %q, got %q", i, exp, result[i])
		}
	}
}

// --- Category Confirm Tests ---

func TestAICategoryConfirmWithSelection(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIFrameworkCategories
	m.Choices.AITools = []string{"claude"}
	m.Choices.InstallAIFramework = true
	m.AICategorySelected = make(map[string][]bool)
	for _, cat := range moduleCategories {
		m.AICategorySelected[cat.ID] = make([]bool, len(cat.Items))
	}
	m.AICategorySelected["hooks"][0] = true
	m.AICategorySelected["skills"][0] = true

	opts := m.GetCurrentOptions()
	m.Cursor = len(opts) - 1 // "Confirm selection"

	result, _ := m.handleAICategoriesKeys("enter")
	newModel := result.(Model)

	// collectSelectedFeatures produces feature-level IDs
	if len(newModel.Choices.AIFrameworkModules) != 2 {
		t.Fatalf("Expected 2 features, got %d: %v", len(newModel.Choices.AIFrameworkModules), newModel.Choices.AIFrameworkModules)
	}
	if newModel.Choices.AIFrameworkModules[0] != "hooks" {
		t.Errorf("Expected 'hooks', got %s", newModel.Choices.AIFrameworkModules[0])
	}
	if newModel.Choices.AIFrameworkModules[1] != "skills" {
		t.Errorf("Expected 'skills', got %s", newModel.Choices.AIFrameworkModules[1])
	}
}

func TestAICategoryConfirmAgentTeamsLite(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIFrameworkCategories
	m.Choices.AITools = []string{"claude"}
	m.Choices.InstallAIFramework = true
	m.AICategorySelected = make(map[string][]bool)
	for _, cat := range moduleCategories {
		m.AICategorySelected[cat.ID] = make([]bool, len(cat.Items))
	}
	// Select only Agent Teams Lite in SDD
	m.AICategorySelected["sdd"][1] = true

	opts := m.GetCurrentOptions()
	m.Cursor = len(opts) - 1 // "Confirm selection"

	result, _ := m.handleAICategoriesKeys("enter")
	newModel := result.(Model)

	// No features for setup-global.sh
	if len(newModel.Choices.AIFrameworkModules) != 0 {
		t.Errorf("Expected 0 features, got %d: %v", len(newModel.Choices.AIFrameworkModules), newModel.Choices.AIFrameworkModules)
	}
	// But Agent Teams Lite flag should be set
	if !newModel.Choices.InstallAgentTeamsLite {
		t.Error("Expected InstallAgentTeamsLite to be true")
	}
	// And install should still proceed
	if !newModel.Choices.InstallAIFramework {
		t.Error("Expected InstallAIFramework to remain true when Agent Teams Lite selected")
	}
}

func TestAICategoryConfirmBothSDDOptions(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIFrameworkCategories
	m.Choices.AITools = []string{"claude"}
	m.Choices.InstallAIFramework = true
	m.AICategorySelected = make(map[string][]bool)
	for _, cat := range moduleCategories {
		m.AICategorySelected[cat.ID] = make([]bool, len(cat.Items))
	}
	// Select both OpenSpec and Agent Teams Lite
	m.AICategorySelected["sdd"][0] = true // OpenSpec
	m.AICategorySelected["sdd"][1] = true // Agent Teams Lite

	opts := m.GetCurrentOptions()
	m.Cursor = len(opts) - 1

	result, _ := m.handleAICategoriesKeys("enter")
	newModel := result.(Model)

	// Should have "sdd" in features
	hasSDD := false
	for _, f := range newModel.Choices.AIFrameworkModules {
		if f == "sdd" {
			hasSDD = true
		}
	}
	if !hasSDD {
		t.Error("Expected 'sdd' in features when OpenSpec selected")
	}
	// And Agent Teams Lite flag
	if !newModel.Choices.InstallAgentTeamsLite {
		t.Error("Expected InstallAgentTeamsLite to be true")
	}
}

func TestIsAgentTeamsLiteSelected(t *testing.T) {
	sel := make(map[string][]bool)
	for _, cat := range moduleCategories {
		sel[cat.ID] = make([]bool, len(cat.Items))
	}

	// Nothing selected
	if isAgentTeamsLiteSelected(sel) {
		t.Error("Expected false when nothing selected")
	}

	// Only OpenSpec
	sel["sdd"][0] = true
	if isAgentTeamsLiteSelected(sel) {
		t.Error("Expected false when only OpenSpec selected")
	}

	// Agent Teams Lite selected
	sel["sdd"][1] = true
	if !isAgentTeamsLiteSelected(sel) {
		t.Error("Expected true when Agent Teams Lite selected")
	}
}

func TestAICategoryConfirmEmpty(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIFrameworkCategories
	m.Choices.AITools = []string{"claude"}
	m.Choices.InstallAIFramework = true
	m.AICategorySelected = make(map[string][]bool)
	for _, cat := range moduleCategories {
		m.AICategorySelected[cat.ID] = make([]bool, len(cat.Items))
	}

	opts := m.GetCurrentOptions()
	m.Cursor = len(opts) - 1

	result, _ := m.handleAICategoriesKeys("enter")
	newModel := result.(Model)

	if newModel.Choices.InstallAIFramework {
		t.Error("Expected InstallAIFramework to be false when no modules selected")
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
		{ScreenAIFrameworkCategories, "Step 8"},
		{ScreenAIFrameworkCategoryItems, "Step 8"},
	}

	for _, tt := range tests {
		m.Screen = tt.screen
		title := m.GetScreenTitle()
		if !strings.Contains(title, tt.expect) {
			t.Errorf("Screen %v: expected title containing %q, got %q", tt.screen, tt.expect, title)
		}
	}
}

func TestAICategoryItemsScreenTitle(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIFrameworkCategoryItems
	m.SelectedModuleCategory = 0 // Hooks (first category)

	title := m.GetScreenTitle()
	if !strings.Contains(title, "Hooks") {
		t.Errorf("Expected title to contain category name 'Hooks', got %q", title)
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

func TestAIFrameworkPresetCustomGoesToCategories(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIFrameworkPreset
	m.Choices.AITools = []string{"claude"}
	m.Choices.InstallAIFramework = true
	m.Cursor = 0 // Custom (first option)

	result, _ := m.handleSelection()
	newModel := result.(Model)

	if newModel.Screen != ScreenAIFrameworkCategories {
		t.Errorf("Expected ScreenAIFrameworkCategories, got %v", newModel.Screen)
	}
	if newModel.AICategorySelected == nil {
		t.Error("Expected AICategorySelected to be initialized")
	}
	// Should have entries for all 6 categories
	if len(newModel.AICategorySelected) != 6 {
		t.Errorf("Expected 6 category entries, got %d", len(newModel.AICategorySelected))
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

func TestAICategoriesGoBack(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIFrameworkCategories
	m.AICategorySelected = make(map[string][]bool)

	result, _ := m.goBackInstallStep()
	newModel := result.(Model)

	if newModel.Screen != ScreenAIFrameworkPreset {
		t.Errorf("Expected ScreenAIFrameworkPreset, got %v", newModel.Screen)
	}
	if newModel.AICategorySelected != nil {
		t.Error("Expected AICategorySelected to be cleared on back")
	}
}

func TestAICategoryItemsGoBack(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIFrameworkCategoryItems
	m.SelectedModuleCategory = 3 // Skills
	m.AICategorySelected = make(map[string][]bool)

	result, _ := m.goBackInstallStep()
	newModel := result.(Model)

	if newModel.Screen != ScreenAIFrameworkCategories {
		t.Errorf("Expected ScreenAIFrameworkCategories, got %v", newModel.Screen)
	}
	if newModel.Cursor != 3 {
		t.Errorf("Expected cursor preserved at 3 (Skills), got %d", newModel.Cursor)
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

func TestRenderAICategoryItemsCheckboxes(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIFrameworkCategoryItems
	m.Width = 100
	m.Height = 40
	m.SelectedModuleCategory = 0 // Hooks (first category)
	m.AICategorySelected = make(map[string][]bool)
	for _, cat := range moduleCategories {
		m.AICategorySelected[cat.ID] = make([]bool, len(cat.Items))
	}
	m.AICategorySelected["hooks"][0] = true
	m.Cursor = 0

	rendered := m.renderAICategoryItems()

	if !strings.Contains(rendered, "[✓]") {
		t.Error("Expected rendered view to contain checked checkbox [✓]")
	}
	if !strings.Contains(rendered, "[ ]") {
		t.Error("Expected rendered view to contain unchecked checkbox [ ]")
	}
}

func TestRenderAICategoryMenuNoCounts(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIFrameworkCategories
	m.Width = 100
	m.Height = 40
	m.AICategorySelected = make(map[string][]bool)
	for _, cat := range moduleCategories {
		m.AICategorySelected[cat.ID] = make([]bool, len(cat.Items))
	}
	m.Cursor = 0

	rendered := m.renderAICategoryMenu()

	// Should contain category names
	if !strings.Contains(rendered, "Hooks") {
		t.Error("Expected rendered view to contain 'Hooks'")
	}
	// Should NOT contain checkboxes (category menu uses cursor only)
	if strings.Contains(rendered, "[✓]") || strings.Contains(rendered, "[ ]") {
		t.Error("Category menu should NOT contain checkboxes")
	}
}

// ==========================================================================
// Separator Skipping Tests
// ==========================================================================

func TestAIToolsSeparatorSkipDown(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIToolsSelect
	m.AIToolSelected = make([]bool, len(aiToolIDMap))
	m.Cursor = len(aiToolIDMap) - 1 // Last tool (GitHub Copilot, index 3)

	// Navigate down — should skip separator (index 4) and land on Confirm (index 5)
	result, _ := m.handleAIToolsKeys("down")
	newModel := result.(Model)

	opts := newModel.GetCurrentOptions()
	confirmIdx := len(opts) - 1
	if newModel.Cursor != confirmIdx {
		t.Errorf("Expected cursor to skip separator and land on Confirm (index %d), got %d", confirmIdx, newModel.Cursor)
	}
}

func TestAIToolsSeparatorSkipUp(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIToolsSelect
	m.AIToolSelected = make([]bool, len(aiToolIDMap))

	opts := m.GetCurrentOptions()
	m.Cursor = len(opts) - 1 // "Confirm selection" (last index)

	// Navigate up — should skip separator and land on last tool
	result, _ := m.handleAIToolsKeys("up")
	newModel := result.(Model)

	if newModel.Cursor != len(aiToolIDMap)-1 {
		t.Errorf("Expected cursor to skip separator and land on last tool (index %d), got %d", len(aiToolIDMap)-1, newModel.Cursor)
	}
}

func TestAICategoriesSeparatorSkipDown(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIFrameworkCategories
	m.AICategorySelected = make(map[string][]bool)
	for _, cat := range moduleCategories {
		m.AICategorySelected[cat.ID] = make([]bool, len(cat.Items))
	}
	m.Cursor = len(moduleCategories) - 1 // Last category (MCP, index 5)

	// Navigate down — should skip separator and land on Confirm
	result, _ := m.handleAICategoriesKeys("down")
	newModel := result.(Model)

	opts := m.GetCurrentOptions()
	confirmIdx := len(opts) - 1
	if newModel.Cursor != confirmIdx {
		t.Errorf("Expected cursor at Confirm (index %d), got %d", confirmIdx, newModel.Cursor)
	}
}

func TestAICategoriesSeparatorSkipUp(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIFrameworkCategories
	m.AICategorySelected = make(map[string][]bool)
	for _, cat := range moduleCategories {
		m.AICategorySelected[cat.ID] = make([]bool, len(cat.Items))
	}
	opts := m.GetCurrentOptions()
	m.Cursor = len(opts) - 1 // Confirm selection

	// Navigate up — should skip separator
	result, _ := m.handleAICategoriesKeys("up")
	newModel := result.(Model)

	if newModel.Cursor != len(moduleCategories)-1 {
		t.Errorf("Expected cursor at last category (index %d), got %d", len(moduleCategories)-1, newModel.Cursor)
	}
}

func TestAICategoryItemsSeparatorSkipDown(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIFrameworkCategoryItems
	m.SelectedModuleCategory = 5 // MCP (6 items)
	m.Height = 50
	m.AICategorySelected = make(map[string][]bool)
	for _, cat := range moduleCategories {
		m.AICategorySelected[cat.ID] = make([]bool, len(cat.Items))
	}
	cat := moduleCategories[5]
	m.Cursor = len(cat.Items) - 1 // Last MCP item

	// Navigate down — should skip separator and land on "← Back"
	result, _ := m.handleAICategoryItemsKeys("down")
	newModel := result.(Model)

	opts := m.GetCurrentOptions()
	backIdx := len(opts) - 1
	if newModel.Cursor != backIdx {
		t.Errorf("Expected cursor at Back (index %d), got %d", backIdx, newModel.Cursor)
	}
}

func TestAICategoryItemsSeparatorSkipUp(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIFrameworkCategoryItems
	m.SelectedModuleCategory = 5 // MCP (6 items)
	m.Height = 50
	m.AICategorySelected = make(map[string][]bool)
	for _, cat := range moduleCategories {
		m.AICategorySelected[cat.ID] = make([]bool, len(cat.Items))
	}
	opts := m.GetCurrentOptions()
	m.Cursor = len(opts) - 1 // "← Back"

	// Navigate up — should skip separator
	result, _ := m.handleAICategoryItemsKeys("up")
	newModel := result.(Model)

	cat := moduleCategories[5]
	if newModel.Cursor != len(cat.Items)-1 {
		t.Errorf("Expected cursor at last item (index %d), got %d", len(cat.Items)-1, newModel.Cursor)
	}
}

// ==========================================================================
// Vim-style Navigation (j/k) Tests
// ==========================================================================

func TestAIToolsVimNavigation(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIToolsSelect
	m.AIToolSelected = make([]bool, len(aiToolIDMap))
	m.Cursor = 0

	// k at top should stay at 0
	result, _ := m.handleAIToolsKeys("k")
	newModel := result.(Model)
	if newModel.Cursor != 0 {
		t.Errorf("k at top: expected cursor 0, got %d", newModel.Cursor)
	}

	// j should move down
	result, _ = newModel.handleAIToolsKeys("j")
	newModel = result.(Model)
	if newModel.Cursor != 1 {
		t.Errorf("j from 0: expected cursor 1, got %d", newModel.Cursor)
	}

	// k should move up
	result, _ = newModel.handleAIToolsKeys("k")
	newModel = result.(Model)
	if newModel.Cursor != 0 {
		t.Errorf("k from 1: expected cursor 0, got %d", newModel.Cursor)
	}
}

// ==========================================================================
// Scroll Sync Tests (viewport follows cursor)
// ==========================================================================

func TestCategoryItemsScrollSyncDown(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIFrameworkCategoryItems
	m.SelectedModuleCategory = 3 // Skills (85 items)
	m.Height = 20                // Small height to trigger scrolling
	m.AICategorySelected = make(map[string][]bool)
	for _, cat := range moduleCategories {
		m.AICategorySelected[cat.ID] = make([]bool, len(cat.Items))
	}
	m.Cursor = 0
	m.CategoryItemsScroll = 0

	// Navigate down past the visible area
	visibleItems := m.Height - 8 // 12
	for i := 0; i < visibleItems+2; i++ {
		result, _ := m.handleAICategoryItemsKeys("down")
		m = result.(Model)
	}

	// Cursor should be past visible area, scroll should have adjusted
	if m.CategoryItemsScroll == 0 {
		t.Error("Expected CategoryItemsScroll to advance when cursor moves past visible area")
	}
	if m.Cursor < m.CategoryItemsScroll || m.Cursor >= m.CategoryItemsScroll+visibleItems {
		t.Errorf("Cursor %d should be within visible range [%d, %d)", m.Cursor, m.CategoryItemsScroll, m.CategoryItemsScroll+visibleItems)
	}
}

func TestCategoryItemsScrollSyncUp(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIFrameworkCategoryItems
	m.SelectedModuleCategory = 3 // Skills (85 items)
	m.Height = 20
	m.AICategorySelected = make(map[string][]bool)
	for _, cat := range moduleCategories {
		m.AICategorySelected[cat.ID] = make([]bool, len(cat.Items))
	}
	// Start with cursor and scroll offset in the middle
	m.Cursor = 30
	m.CategoryItemsScroll = 30

	// Navigate up past the scroll offset
	for i := 0; i < 5; i++ {
		result, _ := m.handleAICategoryItemsKeys("up")
		m = result.(Model)
	}

	// Scroll should follow cursor up
	if m.CategoryItemsScroll > m.Cursor {
		t.Errorf("Scroll (%d) should not be above cursor (%d)", m.CategoryItemsScroll, m.Cursor)
	}
}

func TestCategoryItemsScrollMinimumVisible(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIFrameworkCategoryItems
	m.SelectedModuleCategory = 0 // Hooks (10 items)
	m.Height = 5                 // Very small — should cap at minimum 5 visible
	m.AICategorySelected = make(map[string][]bool)
	for _, cat := range moduleCategories {
		m.AICategorySelected[cat.ID] = make([]bool, len(cat.Items))
	}
	m.Cursor = 0
	m.CategoryItemsScroll = 0

	// Navigate down — scroll sync uses minimum 5 visible items
	for i := 0; i < 7; i++ {
		result, _ := m.handleAICategoryItemsKeys("down")
		m = result.(Model)
	}

	// Should not crash or go negative
	if m.CategoryItemsScroll < 0 {
		t.Errorf("Scroll should not be negative, got %d", m.CategoryItemsScroll)
	}
}

func TestCategoryItemsScrollResetOnScreenChange(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIFrameworkCategoryItems
	m.SelectedModuleCategory = 3 // Skills
	m.Height = 20
	m.CategoryItemsScroll = 25 // Non-zero scroll
	m.AICategorySelected = make(map[string][]bool)
	for _, cat := range moduleCategories {
		m.AICategorySelected[cat.ID] = make([]bool, len(cat.Items))
	}

	// Press esc to go back to categories
	result, _ := m.handleAICategoryItemsKeys("esc")
	newModel := result.(Model)

	if newModel.CategoryItemsScroll != 0 {
		t.Errorf("Expected scroll reset to 0 on screen change, got %d", newModel.CategoryItemsScroll)
	}
}

// ==========================================================================
// handleEscape from BackupConfirm Tests
// ==========================================================================

func TestHandleEscapeFromBackupConfirmToCategories(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenBackupConfirm
	m.Choices.AITools = []string{"claude"}
	m.Choices.InstallAIFramework = true
	m.AICategorySelected = make(map[string][]bool) // Custom mode

	result, _ := m.handleEscape()
	newModel := result.(Model)

	if newModel.Screen != ScreenAIFrameworkCategories {
		t.Errorf("Expected ScreenAIFrameworkCategories, got %v", newModel.Screen)
	}
}

func TestHandleEscapeFromBackupConfirmToPreset(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenBackupConfirm
	m.Choices.AITools = []string{"claude"}
	m.Choices.InstallAIFramework = true
	m.AICategorySelected = nil // Preset mode (not custom)

	result, _ := m.handleEscape()
	newModel := result.(Model)

	if newModel.Screen != ScreenAIFrameworkPreset {
		t.Errorf("Expected ScreenAIFrameworkPreset, got %v", newModel.Screen)
	}
}

func TestHandleEscapeFromBackupConfirmToFrameworkConfirm(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenBackupConfirm
	m.Choices.AITools = []string{"claude"}
	m.Choices.InstallAIFramework = false

	result, _ := m.handleEscape()
	newModel := result.(Model)

	if newModel.Screen != ScreenAIFrameworkConfirm {
		t.Errorf("Expected ScreenAIFrameworkConfirm, got %v", newModel.Screen)
	}
}

func TestHandleEscapeFromBackupConfirmToAITools(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenBackupConfirm
	m.Choices.AITools = nil // No AI tools

	result, _ := m.handleEscape()
	newModel := result.(Model)

	if newModel.Screen != ScreenAIToolsSelect {
		t.Errorf("Expected ScreenAIToolsSelect, got %v", newModel.Screen)
	}
}

// ==========================================================================
// handleBackupConfirmKeys esc Tests (same logic, different entry point)
// ==========================================================================

func TestBackupConfirmEscToCategories(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenBackupConfirm
	m.Choices.AITools = []string{"claude"}
	m.Choices.InstallAIFramework = true
	m.AICategorySelected = make(map[string][]bool)

	result, _ := m.handleBackupConfirmKeys("esc")
	newModel := result.(Model)

	if newModel.Screen != ScreenAIFrameworkCategories {
		t.Errorf("Expected ScreenAIFrameworkCategories, got %v", newModel.Screen)
	}
}

func TestBackupConfirmEscToPreset(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenBackupConfirm
	m.Choices.AITools = []string{"claude"}
	m.Choices.InstallAIFramework = true
	m.AICategorySelected = nil

	result, _ := m.handleBackupConfirmKeys("esc")
	newModel := result.(Model)

	if newModel.Screen != ScreenAIFrameworkPreset {
		t.Errorf("Expected ScreenAIFrameworkPreset, got %v", newModel.Screen)
	}
}

func TestBackupConfirmEscToFrameworkConfirm(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenBackupConfirm
	m.Choices.AITools = []string{"opencode"}
	m.Choices.InstallAIFramework = false

	result, _ := m.handleBackupConfirmKeys("esc")
	newModel := result.(Model)

	if newModel.Screen != ScreenAIFrameworkConfirm {
		t.Errorf("Expected ScreenAIFrameworkConfirm, got %v", newModel.Screen)
	}
}

func TestBackupConfirmEscToAITools(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenBackupConfirm
	m.Choices.AITools = nil

	result, _ := m.handleBackupConfirmKeys("esc")
	newModel := result.(Model)

	if newModel.Screen != ScreenAIToolsSelect {
		t.Errorf("Expected ScreenAIToolsSelect, got %v", newModel.Screen)
	}
}

// ==========================================================================
// Out-of-Bounds Guard Tests
// ==========================================================================

func TestCategoryItemsOutOfBoundsGuard(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIFrameworkCategoryItems
	m.SelectedModuleCategory = 999 // Invalid
	m.Height = 40
	m.AICategorySelected = make(map[string][]bool)

	// Should not panic
	result, _ := m.handleAICategoryItemsKeys("enter")
	newModel := result.(Model)

	// Should be a no-op
	if newModel.Screen != ScreenAIFrameworkCategoryItems {
		t.Errorf("Expected screen unchanged, got %v", newModel.Screen)
	}
}

func TestCategoryItemsNegativeIndexGuard(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIFrameworkCategoryItems
	m.SelectedModuleCategory = -1 // Negative
	m.Height = 40
	m.AICategorySelected = make(map[string][]bool)

	result, _ := m.handleAICategoryItemsKeys("down")
	newModel := result.(Model)
	_ = newModel // Should not panic
}

func TestGetCurrentOptionsOutOfBoundsCategory(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIFrameworkCategoryItems
	m.SelectedModuleCategory = 999

	opts := m.GetCurrentOptions()
	if len(opts) != 0 {
		t.Errorf("Expected empty options for invalid category, got %d items", len(opts))
	}
}

func TestGetScreenTitleOutOfBoundsCategory(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIFrameworkCategoryItems
	m.SelectedModuleCategory = -1

	title := m.GetScreenTitle()
	if !strings.Contains(title, "Step 8") {
		t.Errorf("Expected fallback title with 'Step 8', got %q", title)
	}
	if strings.Contains(title, "Hooks") || strings.Contains(title, "Skills") {
		t.Error("Fallback title should not contain a category name")
	}
}

// ==========================================================================
// Nil/Empty Guard Tests
// ==========================================================================

func TestAIToolsToggleWithNilSelected(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIToolsSelect
	m.AIToolSelected = nil // nil slice
	m.Cursor = 0

	// Should not panic
	result, _ := m.handleAIToolsKeys("enter")
	newModel := result.(Model)
	_ = newModel
}

func TestCategoryItemsToggleWithNilBools(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIFrameworkCategoryItems
	m.SelectedModuleCategory = 0
	m.Height = 40
	m.AICategorySelected = make(map[string][]bool)
	// Don't initialize "hooks" key
	m.AICategorySelected["hooks"] = nil
	m.Cursor = 0

	// Should not panic
	result, _ := m.handleAICategoryItemsKeys("enter")
	newModel := result.(Model)
	_ = newModel
}

func TestCollectSelectedFeaturesMissingKey(t *testing.T) {
	sel := make(map[string][]bool)
	// Only provide some categories, not all
	sel["hooks"] = []bool{true, false, false, false, false, false, false, false, false, false}

	// Should not panic — missing keys are skipped
	result := collectSelectedFeatures(sel)

	if len(result) != 1 || result[0] != "hooks" {
		t.Errorf("Expected [hooks], got %v", result)
	}
}

// ==========================================================================
// GetScreenDescription Tests
// ==========================================================================

func TestAIScreenDescriptions(t *testing.T) {
	tests := []struct {
		screen Screen
		expect string
	}{
		{ScreenAIToolsSelect, "Toggle tools"},
		{ScreenAIFrameworkConfirm, "Agents, skills"},
		{ScreenAIFrameworkPreset, "Presets bundle"},
		{ScreenAIFrameworkCategories, "Select a category"},
		{ScreenAIFrameworkCategoryItems, "Toggle modules"},
	}

	m := NewModel()
	for _, tt := range tests {
		m.Screen = tt.screen
		desc := m.GetScreenDescription()
		if !strings.Contains(desc, tt.expect) {
			t.Errorf("Screen %v: expected description containing %q, got %q", tt.screen, tt.expect, desc)
		}
	}
}

// ==========================================================================
// SetupInstallSteps Custom Preset Label Test
// ==========================================================================

func TestSetupInstallStepsCustomPresetLabel(t *testing.T) {
	m := NewModel()
	m.Choices.Shell = "fish"
	m.Choices.AITools = []string{"claude"}
	m.Choices.InstallAIFramework = true
	m.Choices.AIFrameworkPreset = "" // Custom mode

	m.SetupInstallSteps()

	for _, step := range m.Steps {
		if step.ID == "aiframework" {
			if !strings.Contains(step.Description, "custom") {
				t.Errorf("Expected custom preset label in description, got %q", step.Description)
			}
			return
		}
	}
	t.Error("Expected aiframework step to exist")
}

// ==========================================================================
// Termux AI Skip Test
// ==========================================================================

func TestTermuxSkipsAIScreens(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenNvimSelect
	m.SystemInfo.IsTermux = true
	m.Cursor = 0 // Yes, install Neovim

	result, _ := m.handleSelection()
	newModel := result.(Model)

	// Should skip AI tools and go straight to backup/install
	if newModel.Screen == ScreenAIToolsSelect {
		t.Error("Termux should skip AI tools screen")
	}
	if newModel.Screen != ScreenBackupConfirm && newModel.Screen != ScreenInstalling {
		t.Errorf("Expected ScreenBackupConfirm or ScreenInstalling for Termux, got %v", newModel.Screen)
	}
}

// ==========================================================================
// Preset Separator No-Op Test
// ==========================================================================

func TestPresetSeparatorIsNoOp(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIFrameworkPreset
	m.Choices.AITools = []string{"claude"}
	m.Choices.InstallAIFramework = true
	m.Cursor = 1 // Separator index

	result, _ := m.handleSelection()
	newModel := result.(Model)

	// Should remain on same screen (no-op for separator)
	if newModel.Screen != ScreenAIFrameworkPreset {
		t.Errorf("Selecting separator should be no-op, but screen changed to %v", newModel.Screen)
	}
}

// ==========================================================================
// Render Viewport Scroll Indicators Tests
// ==========================================================================

func TestRenderCategoryItemsScrollIndicators(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIFrameworkCategoryItems
	m.Width = 100
	m.Height = 15 // Small enough to trigger scrolling for Skills (85 items)
	m.SelectedModuleCategory = 3 // Skills
	m.AICategorySelected = make(map[string][]bool)
	for _, cat := range moduleCategories {
		m.AICategorySelected[cat.ID] = make([]bool, len(cat.Items))
	}
	m.Cursor = 0
	m.CategoryItemsScroll = 0

	rendered := m.renderAICategoryItems()

	// At top: should NOT have "above" indicator but SHOULD have "below"
	if strings.Contains(rendered, "more above") {
		t.Error("At top of list, should not show 'more above' indicator")
	}
	if !strings.Contains(rendered, "more below") {
		t.Error("With 85 items in small viewport, should show 'more below' indicator")
	}
}

func TestRenderCategoryItemsScrollUpIndicator(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIFrameworkCategoryItems
	m.Width = 100
	m.Height = 15
	m.SelectedModuleCategory = 3 // Skills (85 items)
	m.AICategorySelected = make(map[string][]bool)
	for _, cat := range moduleCategories {
		m.AICategorySelected[cat.ID] = make([]bool, len(cat.Items))
	}
	m.Cursor = 40
	m.CategoryItemsScroll = 35

	rendered := m.renderAICategoryItems()

	// In middle: should have both indicators
	if !strings.Contains(rendered, "more above") {
		t.Error("In middle of scrolled list, should show 'more above' indicator")
	}
	if !strings.Contains(rendered, "more below") {
		t.Error("In middle of scrolled list, should show 'more below' indicator")
	}
}

// ==========================================================================
// Space Key vs Leader Mode Tests
// ==========================================================================

func TestSpaceActivatesLeaderModeOnAIScreens(t *testing.T) {
	// Space should activate leader mode on AI screens, NOT toggle items
	// This verifies that handleAIToolsKeys " " branch is unreachable via Update()
	m := NewModel()
	m.Screen = ScreenAIToolsSelect
	m.AIToolSelected = make([]bool, len(aiToolIDMap))
	m.Cursor = 0

	result, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
	newModel := result.(Model)

	// Space should have activated leader mode, NOT toggled the tool
	if !newModel.LeaderMode {
		t.Error("Space should activate leader mode on AI tools screen")
	}
	if newModel.AIToolSelected[0] {
		t.Error("Space should NOT toggle AI tool selection (leader mode intercepts)")
	}
}

func TestSpaceActivatesLeaderModeOnCategoryItems(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIFrameworkCategoryItems
	m.SelectedModuleCategory = 0
	m.Height = 40
	m.AICategorySelected = make(map[string][]bool)
	for _, cat := range moduleCategories {
		m.AICategorySelected[cat.ID] = make([]bool, len(cat.Items))
	}
	m.Cursor = 0

	result, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
	newModel := result.(Model)

	if !newModel.LeaderMode {
		t.Error("Space should activate leader mode on category items screen")
	}
	if newModel.AICategorySelected["hooks"][0] {
		t.Error("Space should NOT toggle item selection (leader mode intercepts)")
	}
}

// ==========================================================================
// handleAIToolsKeys("esc") delegation test
// ==========================================================================

func TestAIToolsEscDelegatesToGoBack(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIToolsSelect
	m.AIToolSelected = make([]bool, len(aiToolIDMap))

	result, _ := m.handleAIToolsKeys("esc")
	newModel := result.(Model)

	if newModel.Screen != ScreenNvimSelect {
		t.Errorf("Expected esc to go back to ScreenNvimSelect, got %v", newModel.Screen)
	}
}

func TestAICategoriesEscDelegatesToGoBack(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIFrameworkCategories
	m.AICategorySelected = make(map[string][]bool)

	result, _ := m.handleAICategoriesKeys("esc")
	newModel := result.(Model)

	if newModel.Screen != ScreenAIFrameworkPreset {
		t.Errorf("Expected esc to go back to ScreenAIFrameworkPreset, got %v", newModel.Screen)
	}
}

// ==========================================================================
// View() Dispatch for AI Screens
// ==========================================================================

func TestViewDispatchAIScreens(t *testing.T) {
	screens := []struct {
		screen Screen
		expect string
	}{
		{ScreenAIToolsSelect, "AI Coding Tools"},
		{ScreenAIFrameworkCategories, "Module Categories"},
	}

	for _, tt := range screens {
		m := NewModel()
		m.Screen = tt.screen
		m.Width = 100
		m.Height = 40
		m.AIToolSelected = make([]bool, len(aiToolIDMap))
		m.AICategorySelected = make(map[string][]bool)
		for _, cat := range moduleCategories {
			m.AICategorySelected[cat.ID] = make([]bool, len(cat.Items))
		}

		rendered := m.View()
		if !strings.Contains(rendered, tt.expect) {
			t.Errorf("View() for screen %v should contain %q", tt.screen, tt.expect)
		}
	}
}

func TestViewDispatchCategoryItems(t *testing.T) {
	m := NewModel()
	m.Screen = ScreenAIFrameworkCategoryItems
	m.Width = 100
	m.Height = 40
	m.SelectedModuleCategory = 0 // Hooks
	m.AICategorySelected = make(map[string][]bool)
	for _, cat := range moduleCategories {
		m.AICategorySelected[cat.ID] = make([]bool, len(cat.Items))
	}

	rendered := m.View()
	if !strings.Contains(rendered, "Hooks") {
		t.Error("View() for category items should contain category name 'Hooks'")
	}
}
