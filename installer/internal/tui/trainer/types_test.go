package trainer

import (
	"testing"
)

// =============================================================================
// TIPOS BASE: ModuleID, Position, ExerciseType
// =============================================================================

func TestModuleID_Constants(t *testing.T) {
	// Verificar que todos los módulos están definidos correctamente
	modules := []ModuleID{
		ModuleHorizontal,
		ModuleVertical,
		ModuleTextObjects,
		ModuleChangeRepeat,
		ModuleSubstitution,
		ModuleRegex,
		ModuleMacros,
	}

	if len(modules) != 7 {
		t.Errorf("Expected 7 modules, got %d", len(modules))
	}

	// Verificar valores únicos
	seen := make(map[ModuleID]bool)
	for _, m := range modules {
		if seen[m] {
			t.Errorf("Duplicate module ID: %s", m)
		}
		seen[m] = true
	}
}

func TestModuleID_StringValues(t *testing.T) {
	tests := []struct {
		module   ModuleID
		expected string
	}{
		{ModuleHorizontal, "horizontal"},
		{ModuleVertical, "vertical"},
		{ModuleTextObjects, "textobjects"},
		{ModuleChangeRepeat, "cgn"},
		{ModuleSubstitution, "substitution"},
		{ModuleRegex, "regex"},
		{ModuleMacros, "macros"},
	}

	for _, tt := range tests {
		if string(tt.module) != tt.expected {
			t.Errorf("ModuleID %v: expected %q, got %q", tt.module, tt.expected, string(tt.module))
		}
	}
}

func TestExerciseType_Constants(t *testing.T) {
	// Verificar tipos de ejercicio
	types := []ExerciseType{
		ExerciseLesson,
		ExercisePractice,
		ExerciseBoss,
	}

	if len(types) != 3 {
		t.Errorf("Expected 3 exercise types, got %d", len(types))
	}

	// Verificar valores string
	if string(ExerciseLesson) != "lesson" {
		t.Errorf("ExerciseLesson should be 'lesson', got %q", ExerciseLesson)
	}
	if string(ExercisePractice) != "practice" {
		t.Errorf("ExercisePractice should be 'practice', got %q", ExercisePractice)
	}
	if string(ExerciseBoss) != "boss" {
		t.Errorf("ExerciseBoss should be 'boss', got %q", ExerciseBoss)
	}
}

func TestPosition_Creation(t *testing.T) {
	pos := Position{Line: 5, Col: 10}

	if pos.Line != 5 {
		t.Errorf("Position.Line: expected 5, got %d", pos.Line)
	}
	if pos.Col != 10 {
		t.Errorf("Position.Col: expected 10, got %d", pos.Col)
	}
}

func TestPosition_ZeroValue(t *testing.T) {
	var pos Position

	if pos.Line != 0 || pos.Col != 0 {
		t.Errorf("Zero Position should be {0, 0}, got {%d, %d}", pos.Line, pos.Col)
	}
}

// =============================================================================
// EXERCISE STRUCT
// =============================================================================

func TestExercise_Creation(t *testing.T) {
	exercise := Exercise{
		ID:          "horizontal_001",
		Module:      ModuleHorizontal,
		Level:       1,
		Type:        ExerciseLesson,
		Code:        []string{"const foo = 'bar';"},
		CursorPos:   Position{Line: 0, Col: 0},
		Mission:     "Move to 'foo'",
		Solutions:   []string{"w", "W"},
		Optimal:     "w",
		Hint:        "Use word motion",
		Explanation: "w moves to next word",
		TimeoutSecs: 30,
		Points:      10,
	}

	if exercise.ID != "horizontal_001" {
		t.Errorf("Exercise.ID: expected 'horizontal_001', got %q", exercise.ID)
	}
	if exercise.Module != ModuleHorizontal {
		t.Errorf("Exercise.Module: expected ModuleHorizontal, got %v", exercise.Module)
	}
	if len(exercise.Solutions) != 2 {
		t.Errorf("Exercise.Solutions: expected 2 solutions, got %d", len(exercise.Solutions))
	}
}

func TestExercise_WithCursorTarget(t *testing.T) {
	target := &Position{Line: 0, Col: 6}
	exercise := Exercise{
		ID:           "test",
		CursorPos:    Position{Line: 0, Col: 0},
		CursorTarget: target,
	}

	if exercise.CursorTarget == nil {
		t.Error("CursorTarget should not be nil")
	}
	if exercise.CursorTarget.Col != 6 {
		t.Errorf("CursorTarget.Col: expected 6, got %d", exercise.CursorTarget.Col)
	}
}

func TestExercise_WithoutCursorTarget(t *testing.T) {
	// Para ejercicios de text objects, no hay target de posición
	exercise := Exercise{
		ID:           "textobj_001",
		CursorTarget: nil,
	}

	if exercise.CursorTarget != nil {
		t.Error("CursorTarget should be nil for text object exercises")
	}
}

// =============================================================================
// MODULE INFO
// =============================================================================

func TestGetAllModules_ReturnsCorrectCount(t *testing.T) {
	modules := GetAllModules()

	if len(modules) != 7 {
		t.Errorf("Expected 7 modules, got %d", len(modules))
	}
}

func TestGetAllModules_CorrectOrder(t *testing.T) {
	modules := GetAllModules()

	expectedOrder := []ModuleID{
		ModuleHorizontal,
		ModuleVertical,
		ModuleTextObjects,
		ModuleChangeRepeat,
		ModuleSubstitution,
		ModuleRegex,
		ModuleMacros,
	}

	for i, expected := range expectedOrder {
		if modules[i].ID != expected {
			t.Errorf("Module %d: expected %s, got %s", i, expected, modules[i].ID)
		}
	}
}

func TestGetAllModules_HasRequiredFields(t *testing.T) {
	modules := GetAllModules()

	for _, mod := range modules {
		if mod.ID == "" {
			t.Error("Module ID should not be empty")
		}
		if mod.Name == "" {
			t.Errorf("Module %s: Name should not be empty", mod.ID)
		}
		if mod.Icon == "" {
			t.Errorf("Module %s: Icon should not be empty", mod.ID)
		}
		if mod.Description == "" {
			t.Errorf("Module %s: Description should not be empty", mod.ID)
		}
		if mod.BossName == "" {
			t.Errorf("Module %s: BossName should not be empty", mod.ID)
		}
	}
}

func TestGetAllModules_HorizontalIsFirst(t *testing.T) {
	modules := GetAllModules()

	if modules[0].ID != ModuleHorizontal {
		t.Errorf("First module should be Horizontal, got %s", modules[0].ID)
	}
	if modules[0].Icon != "🏃" {
		t.Errorf("Horizontal icon should be 🏃, got %s", modules[0].Icon)
	}
}

func TestGetAllModules_BossNames(t *testing.T) {
	modules := GetAllModules()

	expectedBosses := map[ModuleID]string{
		ModuleHorizontal:   "The Line Walker",
		ModuleVertical:     "The Code Tower",
		ModuleTextObjects:  "The Bracket Demon",
		ModuleChangeRepeat: "The Clone Army",
		ModuleSubstitution: "The Transformer",
		ModuleRegex:        "The Pattern Master",
		ModuleMacros:       "The Automaton",
	}

	for _, mod := range modules {
		expected, ok := expectedBosses[mod.ID]
		if !ok {
			t.Errorf("Unexpected module ID: %s", mod.ID)
			continue
		}
		if mod.BossName != expected {
			t.Errorf("Module %s: expected boss %q, got %q", mod.ID, expected, mod.BossName)
		}
	}
}
