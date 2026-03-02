package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Gentleman-Programming/Gentleman.Dots/installer/internal/system"
	tea "github.com/charmbracelet/bubbletea"
)

// TestAllUserSelectionPaths tests every possible user selection path
func TestAllUserSelectionPaths(t *testing.T) {
	// All possible OS choices
	osChoices := []string{"mac", "linux"}

	// All possible terminal choices (varies by OS)
	terminalChoicesMac := []string{"alacritty", "wezterm", "kitty", "ghostty", ""}
	terminalChoicesLinux := []string{"alacritty", "wezterm", "ghostty", ""}

	// All possible shell choices
	shellChoices := []string{"fish", "zsh", "nushell"}

	// All possible WM choices
	wmChoices := []string{"tmux", "zellij", ""}

	// All possible nvim choices
	nvimChoices := []bool{true, false}

	// All possible font choices
	fontChoices := []bool{true, false}

	for _, os := range osChoices {
		terminals := terminalChoicesMac
		if os == "linux" {
			terminals = terminalChoicesLinux
		}

		for _, terminal := range terminals {
			for _, shell := range shellChoices {
				for _, wm := range wmChoices {
					for _, nvim := range nvimChoices {
						for _, font := range fontChoices {
							// Skip font test if terminal is none
							// Skip font selection if terminal is skipped (empty string)
							if terminal == "" && font {
								continue
							}

							testName := generateTestName(os, terminal, shell, wm, nvim, font)
							t.Run(testName, func(t *testing.T) {
								choices := UserChoices{
									OS:          os,
									Terminal:    terminal,
									Shell:       shell,
									WindowMgr:   wm,
									InstallNvim: nvim,
									InstallFont: font,
								}

								validateChoices(t, choices)
							})
						}
					}
				}
			}
		}
	}
}

func generateTestName(os, terminal, shell, wm string, nvim, font bool) string {
	nvimStr := "no-nvim"
	if nvim {
		nvimStr = "nvim"
	}
	fontStr := "no-font"
	if font {
		fontStr = "font"
	}
	return os + "/" + terminal + "/" + shell + "/" + wm + "/" + nvimStr + "/" + fontStr
}

func validateChoices(t *testing.T, choices UserChoices) {
	m := NewModel()
	m.Choices = choices
	m.SystemInfo = &system.SystemInfo{
		OS:       system.OSMac,
		HasBrew:  true,
		HasXcode: true,
	}
	if choices.OS == "linux" {
		m.SystemInfo.OS = system.OSLinux
		m.SystemInfo.HasBrew = false
	}

	// Setup install steps
	m.SetupInstallSteps()

	// Validate steps are created
	if len(m.Steps) == 0 {
		t.Error("No steps were created")
		return
	}

	// Validate required steps exist
	requiredSteps := []string{"clone", "shell", "setshell", "cleanup"}
	for _, req := range requiredSteps {
		found := false
		for _, step := range m.Steps {
			if step.ID == req {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Required step '%s' not found in steps", req)
		}
	}

	// Validate conditional steps
	if choices.Terminal != "none" && choices.Terminal != "" {
		assertStepExists(t, m.Steps, "terminal", "Terminal step should exist when terminal is selected")
	}

	if choices.WindowMgr != "none" && choices.WindowMgr != "" {
		assertStepExists(t, m.Steps, "wm", "WM step should exist when WM is selected")
	}

	if choices.InstallNvim {
		assertStepExists(t, m.Steps, "nvim", "Nvim step should exist when nvim is selected")
	}

	if choices.InstallFont {
		assertStepExists(t, m.Steps, "font", "Font step should exist when font is selected")
	}

	if choices.OS == "linux" && !m.SystemInfo.HasBrew {
		assertStepExists(t, m.Steps, "homebrew", "Homebrew step should exist on Linux without brew")
		assertStepExists(t, m.Steps, "deps", "Deps step should exist on Linux")
	}
}

func assertStepExists(t *testing.T, steps []InstallStep, id, msg string) {
	for _, step := range steps {
		if step.ID == id {
			return
		}
	}
	t.Error(msg)
}

// TestStepCloneRepository tests the clone step
func TestStepCloneRepository(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping clone test in short mode")
	}

	t.Run("clone creates directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		originalWd, _ := os.Getwd()
		defer os.Chdir(originalWd)
		os.Chdir(tmpDir)

		m := NewModel()
		m.Choices = UserChoices{OS: "mac", Shell: "fish"}

		err := stepCloneRepo(&m)

		// Check if Gentleman.Dots directory exists
		if _, statErr := os.Stat(filepath.Join(tmpDir, "Gentleman.Dots")); os.IsNotExist(statErr) {
			if err == nil {
				t.Error("Clone reported success but directory doesn't exist")
			}
		}
	})
}

// TestStepInstallShellZsh tests zsh installation step
func TestStepInstallShellZsh(t *testing.T) {
	t.Run("zsh step patches config based on WM choice - none", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create mock .zshrc
		zshrc := `WM_VAR="/$TMUX"
WM_CMD="tmux"
function start_if_needed() {
    exec $WM_CMD
}
eval "$(fzf --zsh)"
start_if_needed`

		zshrcPath := filepath.Join(tmpDir, ".zshrc")
		os.WriteFile(zshrcPath, []byte(zshrc), 0644)

		// Patch for WM=none, nvim=false
		err := system.PatchZshForWM(zshrcPath, "none", false)
		if err != nil {
			t.Fatalf("Patch failed: %v", err)
		}

		// Verify WM lines removed
		content, _ := os.ReadFile(zshrcPath)
		contentStr := string(content)

		if contains(contentStr, "WM_VAR") {
			t.Error("WM_VAR should be removed when WM=none")
		}
		if contains(contentStr, "start_if_needed") {
			t.Error("start_if_needed should be removed when WM=none")
		}
		if !contains(contentStr, "if command -v fzf") {
			t.Error("fzf should be wrapped with command check when nvim=false")
		}
	})

	t.Run("zsh step patches config based on WM choice - zellij", func(t *testing.T) {
		tmpDir := t.TempDir()

		zshrc := `WM_VAR="/$TMUX"
WM_CMD="tmux"
function start_if_needed() {
    exec $WM_CMD
}
start_if_needed`

		zshrcPath := filepath.Join(tmpDir, ".zshrc")
		os.WriteFile(zshrcPath, []byte(zshrc), 0644)

		err := system.PatchZshForWM(zshrcPath, "zellij", true)
		if err != nil {
			t.Fatalf("Patch failed: %v", err)
		}

		content, _ := os.ReadFile(zshrcPath)
		contentStr := string(content)

		if !contains(contentStr, "ZELLIJ") {
			t.Error("Should contain ZELLIJ when WM=zellij")
		}
		if !contains(contentStr, `WM_CMD="zellij"`) {
			t.Error("Should have WM_CMD=zellij")
		}
	})

	t.Run("zsh step patches config based on WM choice - tmux", func(t *testing.T) {
		tmpDir := t.TempDir()

		zshrc := `WM_VAR="/$TMUX"
WM_CMD="tmux"
start_if_needed`

		zshrcPath := filepath.Join(tmpDir, ".zshrc")
		os.WriteFile(zshrcPath, []byte(zshrc), 0644)

		err := system.PatchZshForWM(zshrcPath, "tmux", true)
		if err != nil {
			t.Fatalf("Patch failed: %v", err)
		}

		content, _ := os.ReadFile(zshrcPath)
		contentStr := string(content)

		if !contains(contentStr, "TMUX") {
			t.Error("Should keep TMUX when WM=tmux")
		}
		if !contains(contentStr, `WM_CMD="tmux"`) {
			t.Error("Should keep WM_CMD=tmux")
		}
	})
}

// TestStepInstallShellFish tests fish installation step
func TestStepInstallShellFish(t *testing.T) {
	t.Run("fish step patches config based on WM choice - none", func(t *testing.T) {
		tmpDir := t.TempDir()

		fish := `if not set -q TMUX
    tmux
end
fzf --fish | source`

		fishPath := filepath.Join(tmpDir, "config.fish")
		os.WriteFile(fishPath, []byte(fish), 0644)

		err := system.PatchFishForWM(fishPath, "none", false)
		if err != nil {
			t.Fatalf("Patch failed: %v", err)
		}

		content, _ := os.ReadFile(fishPath)
		contentStr := string(content)

		if contains(contentStr, "TMUX") {
			t.Error("TMUX block should be removed when WM=none")
		}
		if !contains(contentStr, "if command -v fzf") {
			t.Error("fzf should be wrapped when nvim=false")
		}
	})

	t.Run("fish step patches config based on WM choice - zellij", func(t *testing.T) {
		tmpDir := t.TempDir()

		fish := `if not set -q TMUX
    tmux
end

#if not set -q ZELLIJ 
#  zellij
#end`

		fishPath := filepath.Join(tmpDir, "config.fish")
		os.WriteFile(fishPath, []byte(fish), 0644)

		err := system.PatchFishForWM(fishPath, "zellij", true)
		if err != nil {
			t.Fatalf("Patch failed: %v", err)
		}

		content, _ := os.ReadFile(fishPath)
		contentStr := string(content)

		if contains(contentStr, "if not set -q TMUX") {
			t.Error("TMUX block should be removed")
		}
		if !contains(contentStr, "if not set -q ZELLIJ") {
			t.Error("ZELLIJ block should be uncommented")
		}
	})
}

// TestStepInstallShellNushell tests nushell installation step
func TestStepInstallShellNushell(t *testing.T) {
	t.Run("nushell step patches config based on WM choice - none", func(t *testing.T) {
		tmpDir := t.TempDir()

		nu := `let MULTIPLEXER = "tmux" 
let MULTIPLEXER_ENV_PREFIX = "TMUX"

def start_multiplexer [] {
  if $MULTIPLEXER_ENV_PREFIX not-in ($env | columns) {
    run-external $MULTIPLEXER
  }
}

start_multiplexer`

		nuPath := filepath.Join(tmpDir, "config.nu")
		os.WriteFile(nuPath, []byte(nu), 0644)

		err := system.PatchNushellForWM(nuPath, "none")
		if err != nil {
			t.Fatalf("Patch failed: %v", err)
		}

		content, _ := os.ReadFile(nuPath)
		contentStr := string(content)

		if contains(contentStr, "MULTIPLEXER") {
			t.Error("MULTIPLEXER should be removed when WM=none")
		}
		if contains(contentStr, "start_multiplexer") {
			t.Error("start_multiplexer should be removed when WM=none")
		}
	})

	t.Run("nushell step patches config based on WM choice - zellij", func(t *testing.T) {
		tmpDir := t.TempDir()

		nu := `let MULTIPLEXER = "tmux" 
let MULTIPLEXER_ENV_PREFIX = "TMUX"
start_multiplexer`

		nuPath := filepath.Join(tmpDir, "config.nu")
		os.WriteFile(nuPath, []byte(nu), 0644)

		err := system.PatchNushellForWM(nuPath, "zellij")
		if err != nil {
			t.Fatalf("Patch failed: %v", err)
		}

		content, _ := os.ReadFile(nuPath)
		contentStr := string(content)

		if !contains(contentStr, `"zellij"`) {
			t.Error("Should contain zellij")
		}
		if !contains(contentStr, `"ZELLIJ"`) {
			t.Error("Should contain ZELLIJ")
		}
	})
}

// TestAllShellAndWMCombinations tests all shell+WM combinations
func TestAllShellAndWMCombinations(t *testing.T) {
	shells := []string{"zsh", "fish", "nushell"}
	wms := []string{"tmux", "zellij", "none"}
	nvimOptions := []bool{true, false}

	for _, shell := range shells {
		for _, wm := range wms {
			for _, nvim := range nvimOptions {
				testName := shell + "/" + wm
				if nvim {
					testName += "/with-nvim"
				} else {
					testName += "/no-nvim"
				}

				t.Run(testName, func(t *testing.T) {
					tmpDir := t.TempDir()
					var err error

					switch shell {
					case "zsh":
						content := `WM_VAR="/$TMUX"
WM_CMD="tmux"
function start_if_needed() { exec $WM_CMD; }
eval "$(fzf --zsh)"
start_if_needed`
						path := filepath.Join(tmpDir, ".zshrc")
						os.WriteFile(path, []byte(content), 0644)
						err = system.PatchZshForWM(path, wm, nvim)

					case "fish":
						content := `if not set -q TMUX
    tmux
end
#if not set -q ZELLIJ 
#  zellij
#end
fzf --fish | source`
						path := filepath.Join(tmpDir, "config.fish")
						os.WriteFile(path, []byte(content), 0644)
						err = system.PatchFishForWM(path, wm, nvim)

					case "nushell":
						content := `let MULTIPLEXER = "tmux" 
let MULTIPLEXER_ENV_PREFIX = "TMUX"
def start_multiplexer [] { run-external $MULTIPLEXER }
start_multiplexer`
						path := filepath.Join(tmpDir, "config.nu")
						os.WriteFile(path, []byte(content), 0644)
						err = system.PatchNushellForWM(path, wm)
					}

					if err != nil {
						t.Errorf("Patch failed for %s/%s: %v", shell, wm, err)
					}
				})
			}
		}
	}
}

// TestNavigationToEachScreen tests navigation to every screen
func TestNavigationToEachScreen(t *testing.T) {
	t.Run("welcome -> main menu", func(t *testing.T) {
		m := NewModel()
		if m.Screen != ScreenWelcome {
			t.Fatal("Should start at welcome")
		}

		m, _ = simulateKeyPress(m, "enter")
		if m.Screen != ScreenMainMenu {
			t.Errorf("Expected MainMenu, got %v", m.Screen)
		}
	})

	t.Run("main menu -> start installation -> OS select", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenMainMenu
		m.Cursor = 0 // Start Installation

		m, _ = simulateKeyPress(m, "enter")
		if m.Screen != ScreenOSSelect {
			t.Errorf("Expected OSSelect, got %v", m.Screen)
		}
	})

	t.Run("OS select mac -> terminal select", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenOSSelect
		m.Cursor = 0 // macOS

		m, _ = simulateKeyPress(m, "enter")
		if m.Screen != ScreenTerminalSelect {
			t.Errorf("Expected TerminalSelect, got %v", m.Screen)
		}
		if m.Choices.OS != "mac" {
			t.Errorf("Expected OS=mac, got %s", m.Choices.OS)
		}
	})

	t.Run("OS select linux -> terminal select", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenOSSelect
		m.Cursor = 1 // Linux

		m, _ = simulateKeyPress(m, "enter")
		if m.Screen != ScreenTerminalSelect {
			t.Errorf("Expected TerminalSelect, got %v", m.Screen)
		}
		if m.Choices.OS != "linux" {
			t.Errorf("Expected OS=linux, got %s", m.Choices.OS)
		}
	})

	t.Run("terminal select -> font select (when terminal selected)", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenTerminalSelect
		m.Choices.OS = "mac"
		m.Cursor = 0 // Alacritty

		m, _ = simulateKeyPress(m, "enter")
		if m.Screen != ScreenFontSelect {
			t.Errorf("Expected FontSelect, got %v", m.Screen)
		}
	})

	t.Run("terminal select none -> shell select (skip font)", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenTerminalSelect
		m.Choices.OS = "mac"

		// Find Skip option
		options := m.GetCurrentOptions()
		for i, opt := range options {
			if strings.Contains(opt, "Skip this step") {
				m.Cursor = i
				break
			}
		}

		m, _ = simulateKeyPress(m, "enter")
		if m.Screen != ScreenShellSelect {
			t.Errorf("Expected ShellSelect when terminal skipped, got %v", m.Screen)
		}
	})

	t.Run("font select -> shell select", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenFontSelect
		m.Cursor = 0 // Yes

		m, _ = simulateKeyPress(m, "enter")
		if m.Screen != ScreenShellSelect {
			t.Errorf("Expected ShellSelect, got %v", m.Screen)
		}
	})

	t.Run("shell select -> WM select", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenShellSelect
		m.Cursor = 0 // Fish

		m, _ = simulateKeyPress(m, "enter")
		if m.Screen != ScreenWMSelect {
			t.Errorf("Expected WMSelect, got %v", m.Screen)
		}
	})

	t.Run("WM select -> nvim select", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenWMSelect
		m.Cursor = 0 // Tmux

		m, _ = simulateKeyPress(m, "enter")
		if m.Screen != ScreenNvimSelect {
			t.Errorf("Expected NvimSelect, got %v", m.Screen)
		}
	})
}

// TestEachTerminalSelection tests each terminal option
func TestEachTerminalSelection(t *testing.T) {
	terminalsMac := []struct {
		name   string
		cursor int
		expect string
	}{
		{"alacritty", 0, "alacritty"},
		{"wezterm", 1, "wezterm"},
		{"kitty", 2, "kitty"},
		{"ghostty", 3, "ghostty"},
	}

	for _, tc := range terminalsMac {
		t.Run("mac/"+tc.name, func(t *testing.T) {
			m := NewModel()
			m.Screen = ScreenTerminalSelect
			m.Choices.OS = "mac"
			m.Cursor = tc.cursor

			m, _ = simulateKeyPress(m, "enter")
			if m.Choices.Terminal != tc.expect {
				t.Errorf("Expected terminal=%s, got %s", tc.expect, m.Choices.Terminal)
			}
		})
	}

	terminalsLinux := []struct {
		name   string
		cursor int
		expect string
	}{
		{"alacritty", 0, "alacritty"},
		{"wezterm", 1, "wezterm"},
		{"ghostty", 2, "ghostty"},
	}

	for _, tc := range terminalsLinux {
		t.Run("linux/"+tc.name, func(t *testing.T) {
			m := NewModel()
			m.Screen = ScreenTerminalSelect
			m.Choices.OS = "linux"
			m.Cursor = tc.cursor

			m, _ = simulateKeyPress(m, "enter")
			if m.Choices.Terminal != tc.expect {
				t.Errorf("Expected terminal=%s, got %s", tc.expect, m.Choices.Terminal)
			}
		})
	}
}

// TestEachShellSelection tests each shell option
func TestEachShellSelection(t *testing.T) {
	shells := []struct {
		name   string
		cursor int
		expect string
	}{
		{"fish", 0, "fish"},
		{"zsh", 1, "zsh"},
		{"nushell", 2, "nushell"},
	}

	for _, tc := range shells {
		t.Run(tc.name, func(t *testing.T) {
			m := NewModel()
			m.Screen = ScreenShellSelect
			m.Cursor = tc.cursor

			m, _ = simulateKeyPress(m, "enter")
			if m.Choices.Shell != tc.expect {
				t.Errorf("Expected shell=%s, got %s", tc.expect, m.Choices.Shell)
			}
		})
	}
}

// TestEachWMSelection tests each WM option
func TestEachWMSelection(t *testing.T) {
	wms := []struct {
		name   string
		cursor int
		expect string
	}{
		{"tmux", 0, "tmux"},
		{"zellij", 1, "zellij"},
	}

	for _, tc := range wms {
		t.Run(tc.name, func(t *testing.T) {
			m := NewModel()
			m.Screen = ScreenWMSelect
			m.Cursor = tc.cursor

			m, _ = simulateKeyPress(m, "enter")
			if m.Choices.WindowMgr != tc.expect {
				t.Errorf("Expected wm=%s, got %s", tc.expect, m.Choices.WindowMgr)
			}
		})
	}
}

// TestEachNvimSelection tests each nvim option
func TestEachNvimSelection(t *testing.T) {
	t.Run("yes", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenNvimSelect
		m.Choices = UserChoices{OS: "mac", Shell: "fish"}
		m.SystemInfo = &system.SystemInfo{OS: system.OSMac, HasBrew: true, HasXcode: true}
		m.Cursor = 0 // Yes

		m, _ = simulateKeyPress(m, "enter")
		if !m.Choices.InstallNvim {
			t.Error("Expected InstallNvim=true")
		}
	})

	t.Run("no", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenNvimSelect
		m.Choices = UserChoices{OS: "mac", Shell: "fish"}
		m.SystemInfo = &system.SystemInfo{OS: system.OSMac, HasBrew: true, HasXcode: true}
		m.Cursor = 1 // No

		m, _ = simulateKeyPress(m, "enter")
		if m.Choices.InstallNvim {
			t.Error("Expected InstallNvim=false")
		}
	})
}

// TestEachFontSelection tests each font option
func TestEachFontSelection(t *testing.T) {
	t.Run("yes", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenFontSelect
		m.Cursor = 0 // Yes

		m, _ = simulateKeyPress(m, "enter")
		if !m.Choices.InstallFont {
			t.Error("Expected InstallFont=true")
		}
	})

	t.Run("no", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenFontSelect
		m.Cursor = 1 // No

		m, _ = simulateKeyPress(m, "enter")
		if m.Choices.InstallFont {
			t.Error("Expected InstallFont=false")
		}
	})
}

// TestBackupOptions tests backup confirmation options
func TestBackupOptions(t *testing.T) {
	t.Run("install with backup", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenBackupConfirm
		m.ExistingConfigs = []string{"fish: /home/user/.config/fish"}
		m.SystemInfo = &system.SystemInfo{OS: system.OSMac, HasBrew: true, HasXcode: true}
		m.Choices = UserChoices{OS: "mac", Shell: "fish"}
		m.Cursor = 0 // Install with Backup

		m, _ = simulateKeyPress(m, "enter")
		if !m.Choices.CreateBackup {
			t.Error("Expected CreateBackup=true")
		}
	})

	t.Run("install without backup", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenBackupConfirm
		m.ExistingConfigs = []string{"fish: /home/user/.config/fish"}
		m.SystemInfo = &system.SystemInfo{OS: system.OSMac, HasBrew: true, HasXcode: true}
		m.Choices = UserChoices{OS: "mac", Shell: "fish"}
		m.Cursor = 1 // Install without Backup

		m, _ = simulateKeyPress(m, "enter")
		if m.Choices.CreateBackup {
			t.Error("Expected CreateBackup=false")
		}
	})

	t.Run("cancel", func(t *testing.T) {
		m := NewModel()
		m.Screen = ScreenBackupConfirm
		m.ExistingConfigs = []string{"fish: /home/user/.config/fish"}
		m.Choices = UserChoices{Shell: "fish"} // User chose fish, so 3 options available
		m.Cursor = 2                           // Cancel

		m, _ = simulateKeyPress(m, "enter")
		if m.Screen != ScreenMainMenu {
			t.Errorf("Expected MainMenu after cancel, got %v", m.Screen)
		}
	})
}

// Helper functions
func simulateKeyPress(m Model, key string) (Model, interface{}) {
	var msg tea.Msg
	switch key {
	case "enter":
		msg = tea.KeyMsg{Type: tea.KeyEnter}
	case "up":
		msg = tea.KeyMsg{Type: tea.KeyUp}
	case "down":
		msg = tea.KeyMsg{Type: tea.KeyDown}
	case "esc":
		msg = tea.KeyMsg{Type: tea.KeyEscape}
	}

	result, _ := m.Update(msg)
	return result.(Model), nil
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
