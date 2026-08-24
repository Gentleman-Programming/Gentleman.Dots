package tui

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetHomebrewScriptUsesTTYSafePrompt(t *testing.T) {
	originalPath := os.Getenv("PATH")
	t.Cleanup(func() {
		_ = os.Setenv("PATH", originalPath)
	})

	if err := os.Setenv("PATH", ""); err != nil {
		t.Fatalf("failed to clear PATH: %v", err)
	}

	script, err := getHomebrewScript(&Model{})
	if err != nil {
		t.Fatalf("getHomebrewScript returned error: %v", err)
	}

	if !strings.Contains(script, `if [ -t 0 ] && [ -r /dev/tty ]; then`) {
		t.Fatalf("expected Homebrew script to guard prompt with tty check")
	}

	if !strings.Contains(script, `IFS= read -r dummy < /dev/tty`) {
		t.Fatalf("expected Homebrew script to read from /dev/tty")
	}

	if strings.Contains(script, "echo \"Press Enter to continue...\"\nread dummy") {
		t.Fatalf("expected Homebrew script to avoid unconditional read dummy")
	}
}

func TestDescribeInteractiveStepError(t *testing.T) {
	t.Run("homebrew install failure stays specific when brew absent", func(t *testing.T) {
		originalPath := os.Getenv("PATH")
		t.Cleanup(func() {
			_ = os.Setenv("PATH", originalPath)
		})

		if err := os.Setenv("PATH", ""); err != nil {
			t.Fatalf("failed to clear PATH: %v", err)
		}

		err := describeInteractiveStepError("homebrew", errors.New("exit status 1"))
		if !strings.Contains(err.Error(), "Homebrew installation command failed") {
			t.Fatalf("expected install failure message, got %q", err.Error())
		}
	})

	t.Run("post-install failure does not blame homebrew when brew exists", func(t *testing.T) {
		tmpDir := t.TempDir()
		brewPath := filepath.Join(tmpDir, "brew")
		if err := os.WriteFile(brewPath, []byte("#!/bin/sh\nexit 0\n"), 0755); err != nil {
			t.Fatalf("failed to create fake brew: %v", err)
		}

		originalPath := os.Getenv("PATH")
		t.Cleanup(func() {
			_ = os.Setenv("PATH", originalPath)
		})

		if err := os.Setenv("PATH", tmpDir); err != nil {
			t.Fatalf("failed to set PATH: %v", err)
		}

		err := describeInteractiveStepError("homebrew", errors.New("exit status 1"))
		if !strings.Contains(err.Error(), "Homebrew installed, but post-install shell setup did not complete cleanly") {
			t.Fatalf("expected post-install failure message, got %q", err.Error())
		}
	})

	t.Run("non-homebrew steps pass through unchanged", func(t *testing.T) {
		original := errors.New("boom")
		if got := describeInteractiveStepError("deps", original); !errors.Is(got, original) {
			t.Fatalf("expected original error to pass through")
		}
	})
}
