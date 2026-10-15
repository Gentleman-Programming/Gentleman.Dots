package trainer

import "testing"

func TestUnicodeNavigation(t *testing.T) {
	// Basic sanity: ASCII word motion still works
	line := []string{"hello world 🌍 emoji"}
	pos := SimulatedPosition{Line: 0, Col: 0}
	// Move to next word (should skip "hello" and land on "world")
	pos = moveWordForward(pos, line, false)
	if pos.Col == 0 {
		t.Errorf("word forward failed on ASCII")
	}
}
