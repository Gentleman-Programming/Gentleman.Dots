package system

import (
	"os/exec"
	"testing"
)

func TestGetShell(t *testing.T) {
	t.Run("should return a valid shell path", func(t *testing.T) {
		shell := GetShell()

		if shell == "" {
			t.Error("GetShell should return a non-empty path")
		}

		// Verify the shell exists
		if _, err := exec.LookPath(shell); err != nil {
			t.Errorf("GetShell returned '%s' which is not executable: %v", shell, err)
		}
	})

	t.Run("should return same value on repeated calls (cached)", func(t *testing.T) {
		shell1 := GetShell()
		shell2 := GetShell()

		if shell1 != shell2 {
			t.Errorf("GetShell should return cached value: got '%s' and '%s'", shell1, shell2)
		}
	})

	t.Run("returned shell should support -c flag", func(t *testing.T) {
		shell := GetShell()

		// Test that the shell can execute a simple command
		cmd := exec.Command(shell, "-c", "echo hello")
		output, err := cmd.Output()

		if err != nil {
			t.Errorf("Shell '%s' failed to execute 'echo hello': %v", shell, err)
		}

		if len(output) == 0 {
			t.Error("Shell command produced no output")
		}
	})
}

func TestRunWithDynamicShell(t *testing.T) {
	t.Run("Run should work with detected shell", func(t *testing.T) {
		result := Run("echo 'test'", nil)

		if result.Error != nil {
			t.Errorf("Run failed: %v", result.Error)
		}

		if result.ExitCode != 0 {
			t.Errorf("Expected exit code 0, got %d", result.ExitCode)
		}
	})

	t.Run("Run should handle complex commands", func(t *testing.T) {
		result := Run("test -d / && echo 'exists'", nil)

		if result.Error != nil {
			t.Errorf("Run failed with complex command: %v", result.Error)
		}
	})

	t.Run("Run should handle environment variables", func(t *testing.T) {
		result := Run("echo $HOME", nil)

		if result.Error != nil {
			t.Errorf("Run failed with env var: %v", result.Error)
		}

		if result.Output == "" {
			t.Error("Expected HOME to be set")
		}
	})
}
