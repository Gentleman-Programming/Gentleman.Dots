package system

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestPOSIXCompatibility tests that our commands work with POSIX sh
func TestPOSIXCompatibility(t *testing.T) {
	shell := GetShell()

	tests := []struct {
		name    string
		command string
	}{
		{
			name:    "echo command",
			command: "echo 'hello world'",
		},
		{
			name:    "variable assignment",
			command: "VAR=test; echo $VAR",
		},
		{
			name:    "command substitution",
			command: "echo $(echo nested)",
		},
		{
			name:    "test command",
			command: "test -d / && echo 'root exists'",
		},
		{
			name:    "if statement",
			command: "if [ -d / ]; then echo 'yes'; fi",
		},
		{
			name:    "for loop",
			command: "for i in a b c; do echo $i; done",
		},
		{
			name:    "which command",
			command: "which echo >/dev/null 2>&1 && echo 'found'",
		},
		{
			name:    "read with variable (POSIX)",
			command: "echo 'test' | read dummy; echo 'done'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(shell, "-c", tt.command)
			output, err := cmd.CombinedOutput()

			if err != nil {
				t.Errorf("Command failed with shell '%s': %v\nOutput: %s", shell, err, output)
			}
		})
	}
}

// TestScriptExecution tests that scripts can be executed with detected shell
func TestScriptExecution(t *testing.T) {
	shell := GetShell()
	tmpDir := t.TempDir()

	t.Run("execute simple script", func(t *testing.T) {
		scriptPath := filepath.Join(tmpDir, "test.sh")
		script := `#!/bin/sh
set -e
echo "Script running"
VAR="test value"
echo "VAR is: $VAR"
exit 0
`
		if err := os.WriteFile(scriptPath, []byte(script), 0755); err != nil {
			t.Fatalf("Failed to write script: %v", err)
		}

		cmd := exec.Command(shell, scriptPath)
		output, err := cmd.CombinedOutput()

		if err != nil {
			t.Errorf("Script execution failed: %v\nOutput: %s", err, output)
		}
	})

	t.Run("execute script with functions", func(t *testing.T) {
		scriptPath := filepath.Join(tmpDir, "func.sh")
		script := `#!/bin/sh
set -e

my_function() {
    echo "Function called with: $1"
}

my_function "test argument"
exit 0
`
		if err := os.WriteFile(scriptPath, []byte(script), 0755); err != nil {
			t.Fatalf("Failed to write script: %v", err)
		}

		cmd := exec.Command(shell, scriptPath)
		output, err := cmd.CombinedOutput()

		if err != nil {
			t.Errorf("Script with functions failed: %v\nOutput: %s", err, output)
		}
	})
}

// TestRunFunction tests the Run function with various commands
func TestRunFunction(t *testing.T) {
	t.Run("simple command", func(t *testing.T) {
		result := Run("echo 'hello'", nil)
		if result.Error != nil {
			t.Errorf("Run failed: %v", result.Error)
		}
		if result.ExitCode != 0 {
			t.Errorf("Expected exit code 0, got %d", result.ExitCode)
		}
	})

	t.Run("command with pipes", func(t *testing.T) {
		result := Run("echo 'hello world' | wc -w", nil)
		if result.Error != nil {
			t.Errorf("Pipe command failed: %v", result.Error)
		}
	})

	t.Run("command that fails", func(t *testing.T) {
		result := Run("exit 1", nil)
		if result.Error == nil {
			t.Error("Expected error for failing command")
		}
		if result.ExitCode != 1 {
			t.Errorf("Expected exit code 1, got %d", result.ExitCode)
		}
	})

	t.Run("command with working directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		result := Run("pwd", &ExecOptions{WorkDir: tmpDir})
		if result.Error != nil {
			t.Errorf("WorkDir command failed: %v", result.Error)
		}
	})

	t.Run("command with environment", func(t *testing.T) {
		result := Run("echo $MY_TEST_VAR", &ExecOptions{
			Env: []string{"MY_TEST_VAR=test_value"},
		})
		if result.Error != nil {
			t.Errorf("Env command failed: %v", result.Error)
		}
	})
}

// TestGitCloneCompatibility tests git clone works with our shell setup
func TestGitCloneCompatibility(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping git clone test in short mode")
	}

	// Check if git is available
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available, skipping")
	}

	tmpDir := t.TempDir()

	t.Run("git clone with progress", func(t *testing.T) {
		// Use a small public repo for testing
		result := Run("git clone --depth 1 https://github.com/octocat/Hello-World.git test-repo", &ExecOptions{
			WorkDir: tmpDir,
		})

		if result.Error != nil {
			t.Errorf("Git clone failed: %v\nStderr: %s", result.Error, result.Stderr)
		}

		// Verify repo was cloned
		if _, err := os.Stat(filepath.Join(tmpDir, "test-repo")); os.IsNotExist(err) {
			t.Error("Repository was not cloned")
		}
	})
}

// TestBrewCommandCompatibility tests brew commands work (if brew is installed)
func TestBrewCommandCompatibility(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping brew test in short mode")
	}

	// Check if brew is available
	brewPath := GetBrewPrefix() + "/bin/brew"
	if _, err := os.Stat(brewPath); os.IsNotExist(err) {
		t.Skip("brew not available, skipping")
	}

	t.Run("brew --version", func(t *testing.T) {
		result := RunBrew("--version", nil)
		if result.Error != nil {
			t.Errorf("brew --version failed: %v", result.Error)
		}
	})

	t.Run("brew list (should not fail)", func(t *testing.T) {
		result := RunBrew("list --versions | head -1", nil)
		// This might fail if brew has no packages, but shouldn't error
		if result.Error != nil && result.ExitCode != 0 && result.ExitCode != 1 {
			t.Errorf("brew list failed unexpectedly: %v", result.Error)
		}
	})
}
