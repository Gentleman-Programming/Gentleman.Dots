package tui

import (
	"fmt"
	"runtime"

	"github.com/Gentleman-Programming/Gentleman.Dots/installer/internal/system"
)

// RunNonInteractive executes the installation without TUI
func RunNonInteractive(choices UserChoices) error {
	// Enable non-interactive mode for logging
	SetNonInteractiveMode(true)

	// Detect system info
	sysInfo := system.Detect()

	// Determine OS choice based on system
	osChoice := "linux"
	if runtime.GOOS == "darwin" {
		osChoice = "mac"
	}
	choices.OS = osChoice

	// Create a minimal model for the installation functions
	model := &Model{
		SystemInfo: sysInfo,
		Choices:    choices,
		LogLines:   []string{},
	}

	// Detect existing configs for backup functionality
	if choices.CreateBackup {
		model.ExistingConfigs = system.DetectExistingConfigs()
	}

	// Define steps to run based on choices
	steps := buildStepsForChoices(model)

	fmt.Printf("📋 Running %d installation steps...\n\n", len(steps))

	// Execute each step
	for i, step := range steps {
		fmt.Printf("[%d/%d] %s...\n", i+1, len(steps), step.Name)

		err := executeStep(step.ID, model)
		if err != nil {
			fmt.Printf("    ❌ FAILED: %v\n", err)
			return fmt.Errorf("step '%s' failed: %w", step.Name, err)
		}
		fmt.Printf("    ✓ Done\n")
	}

	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("✅ Installation complete!")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	return nil
}

// buildStepsForChoices creates the list of steps based on user choices
func buildStepsForChoices(m *Model) []InstallStep {
	var steps []InstallStep

	// Always backup first if enabled
	if m.Choices.CreateBackup {
		steps = append(steps, InstallStep{ID: "backup", Name: "Backup existing configs"})
	}

	// Clone repo
	steps = append(steps, InstallStep{ID: "clone", Name: "Clone Gentleman.Dots repository"})

	// Homebrew (for Mac and Debian/Ubuntu Linux - NOT Fedora/Arch which use native package managers)
	if m.SystemInfo.OS == system.OSMac || m.SystemInfo.OS == system.OSDebian || m.SystemInfo.OS == system.OSLinux {
		steps = append(steps, InstallStep{ID: "homebrew", Name: "Install/Update Homebrew"})
	}

	// Dependencies
	steps = append(steps, InstallStep{ID: "deps", Name: "Install dependencies"})

	// Xcode (Mac only)
	if m.SystemInfo.OS == system.OSMac {
		steps = append(steps, InstallStep{ID: "xcode", Name: "Install Xcode CLI tools"})
	}

	// Terminal
	if m.Choices.Terminal != "none" {
		steps = append(steps, InstallStep{ID: "terminal", Name: fmt.Sprintf("Install %s terminal", m.Choices.Terminal)})
	}

	// Font
	if m.Choices.InstallFont {
		steps = append(steps, InstallStep{ID: "font", Name: "Install Nerd Font"})
	}

	// Shell
	steps = append(steps, InstallStep{ID: "shell", Name: fmt.Sprintf("Install %s shell", m.Choices.Shell)})

	// Window Manager
	if m.Choices.WindowMgr != "none" {
		steps = append(steps, InstallStep{ID: "wm", Name: fmt.Sprintf("Install %s", m.Choices.WindowMgr)})
	}

	// Neovim
	if m.Choices.InstallNvim {
		steps = append(steps, InstallStep{ID: "nvim", Name: "Install Neovim configuration"})
	}

	// AI Tools
	if len(m.Choices.AITools) > 0 {
		steps = append(steps, InstallStep{ID: "aitools", Name: "Install AI tools"})
	}

	// AI Framework
	if m.Choices.InstallAIFramework {
		steps = append(steps, InstallStep{ID: "aiframework", Name: "Install AI framework"})
	}

	// Set shell as default
	steps = append(steps, InstallStep{ID: "setshell", Name: "Set shell as default"})

	// Cleanup
	steps = append(steps, InstallStep{ID: "cleanup", Name: "Cleanup"})

	return steps
}
