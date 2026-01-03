package system

import (
	"os"
	"runtime"
	"testing"
)

func TestDetect(t *testing.T) {
	info := Detect()

	t.Run("should detect OS", func(t *testing.T) {
		if info.OS == OSUnknown && runtime.GOOS != "windows" {
			t.Error("OS should not be unknown on unix systems")
		}
	})

	t.Run("should have home directory", func(t *testing.T) {
		if info.HomeDir == "" {
			t.Error("HomeDir should not be empty")
		}
	})

	t.Run("should detect current shell", func(t *testing.T) {
		// Only test if SHELL env is set
		if os.Getenv("SHELL") != "" && info.UserShell == "unknown" {
			t.Error("UserShell should be detected when SHELL env is set")
		}
	})

	t.Run("OSName should match OS type", func(t *testing.T) {
		switch runtime.GOOS {
		case "darwin":
			if info.OSName != "macOS" {
				t.Errorf("Expected OSName to be 'macOS', got '%s'", info.OSName)
			}
		case "linux":
			validNames := []string{"Linux", "Arch Linux", "Debian/Ubuntu"}
			found := false
			for _, name := range validNames {
				if info.OSName == name {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Unexpected OSName for Linux: '%s'", info.OSName)
			}
		}
	})
}

func TestCommandExists(t *testing.T) {
	t.Run("should find common commands", func(t *testing.T) {
		// These should exist on any unix system
		commonCmds := []string{"ls", "echo", "cat"}
		for _, cmd := range commonCmds {
			if !CommandExists(cmd) {
				t.Errorf("Command '%s' should exist", cmd)
			}
		}
	})

	t.Run("should not find non-existent commands", func(t *testing.T) {
		if CommandExists("this-command-definitely-does-not-exist-xyz123") {
			t.Error("Should not find non-existent command")
		}
	})
}

func TestGetBrewPrefix(t *testing.T) {
	prefix := GetBrewPrefix()

	t.Run("should return valid prefix based on OS and arch", func(t *testing.T) {
		switch runtime.GOOS {
		case "darwin":
			if runtime.GOARCH == "arm64" {
				if prefix != "/opt/homebrew" {
					t.Errorf("Expected '/opt/homebrew' on macOS ARM64, got '%s'", prefix)
				}
			} else {
				if prefix != "/usr/local" {
					t.Errorf("Expected '/usr/local' on macOS Intel, got '%s'", prefix)
				}
			}
		case "linux":
			if prefix != "/home/linuxbrew/.linuxbrew" {
				t.Errorf("Expected '/home/linuxbrew/.linuxbrew' on Linux, got '%s'", prefix)
			}
		}
	})
}

func TestCheckWSL(t *testing.T) {
	// This test just ensures the function doesn't panic
	t.Run("should not panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("checkWSL panicked: %v", r)
			}
		}()
		_ = checkWSL()
	})
}

func TestIsArchLinux(t *testing.T) {
	t.Run("should not panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("isArchLinux panicked: %v", r)
			}
		}()
		_ = isArchLinux()
	})
}

func TestIsDebian(t *testing.T) {
	t.Run("should not panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("isDebian panicked: %v", r)
			}
		}()
		_ = isDebian()
	})
}
