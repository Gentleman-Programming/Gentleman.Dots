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

func TestIsTermux(t *testing.T) {
	t.Run("should not panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("isTermux panicked: %v", r)
			}
		}()
		_ = isTermux()
	})

	t.Run("should detect TERMUX_VERSION env", func(t *testing.T) {
		// Save original value
		original := os.Getenv("TERMUX_VERSION")
		defer os.Setenv("TERMUX_VERSION", original)

		// Set Termux env
		os.Setenv("TERMUX_VERSION", "0.118.0")
		if !isTermux() {
			t.Error("Should detect Termux when TERMUX_VERSION is set")
		}

		// Unset
		os.Unsetenv("TERMUX_VERSION")
		// Note: might still be true if PREFIX contains termux
	})

	t.Run("should detect PREFIX with termux path", func(t *testing.T) {
		// Save original values
		originalVersion := os.Getenv("TERMUX_VERSION")
		originalPrefix := os.Getenv("PREFIX")
		defer func() {
			os.Setenv("TERMUX_VERSION", originalVersion)
			os.Setenv("PREFIX", originalPrefix)
		}()

		// Clear TERMUX_VERSION, set PREFIX
		os.Unsetenv("TERMUX_VERSION")
		os.Setenv("PREFIX", "/data/data/com.termux/files/usr")

		if !isTermux() {
			t.Error("Should detect Termux when PREFIX contains 'com.termux'")
		}
	})
}

func TestCheckPkg(t *testing.T) {
	t.Run("should not panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("checkPkg panicked: %v", r)
			}
		}()
		_ = checkPkg()
	})
}

func TestDetectTermuxFields(t *testing.T) {
	t.Run("SystemInfo should have Termux fields", func(t *testing.T) {
		info := Detect()
		// Just verify the fields exist and are initialized
		_ = info.IsTermux
		_ = info.HasPkg
		_ = info.Prefix
	})

	t.Run("Non-Termux system should have IsTermux=false", func(t *testing.T) {
		// Save original values
		originalVersion := os.Getenv("TERMUX_VERSION")
		originalPrefix := os.Getenv("PREFIX")
		defer func() {
			if originalVersion != "" {
				os.Setenv("TERMUX_VERSION", originalVersion)
			}
			if originalPrefix != "" {
				os.Setenv("PREFIX", originalPrefix)
			}
		}()

		// Clear Termux env vars
		os.Unsetenv("TERMUX_VERSION")
		os.Setenv("PREFIX", "/usr/local") // Non-termux prefix

		// On non-Termux systems, IsTermux should be false
		// (unless /data/data/com.termux exists, which it won't on normal systems)
		info := Detect()
		if info.IsTermux && info.Prefix != "" && !containsTermux(info.Prefix) {
			t.Error("IsTermux should be false on non-Termux systems")
		}
	})
}

// Helper to check if string contains termux
func containsTermux(s string) bool {
	return len(s) > 0 && (s == "/data/data/com.termux/files/usr" ||
		(len(s) > 10 && s[:10] == "/data/data"))
}
