package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Gentleman-Programming/Gentleman.Dots/installer/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

var Version = "1.0.0"

func main() {
	// Handle flags
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--version", "-v", "version":
			fmt.Printf("gentleman.dots v%s\n", Version)
			os.Exit(0)
		case "--help", "-h", "help":
			printHelp()
			os.Exit(0)
		case "--test", "-t", "test":
			setupTestMode()
		case "--dry-run":
			os.Setenv("GENTLEMAN_DRY_RUN", "1")
			fmt.Println("🧪 Dry-run mode: No actual installations will be performed")
		}
	}

	p := tea.NewProgram(
		tui.NewModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running installer: %v\n", err)
		os.Exit(1)
	}
}

func setupTestMode() {
	// Create a temporary test directory
	testDir := filepath.Join(os.TempDir(), "gentleman-dots-test")
	testHome := filepath.Join(testDir, "home")
	testConfig := filepath.Join(testHome, ".config")

	// Create directories
	os.MkdirAll(testConfig, 0755)

	// Override HOME to use test directory
	os.Setenv("HOME", testHome)
	os.Setenv("GENTLEMAN_TEST_MODE", "1")

	fmt.Printf("🧪 Test mode enabled!\n")
	fmt.Printf("   Test HOME: %s\n", testHome)
	fmt.Printf("   Your real configs are SAFE.\n")
	fmt.Printf("   Press Enter to continue...\n")

	// Wait for user to acknowledge
	fmt.Scanln()
}

func printHelp() {
	fmt.Println(`gentleman.dots - TUI installer for Gentleman.Dots terminal environment

Usage:
  gentleman.dots [flags]

Flags:
  -h, --help      Show this help message
  -v, --version   Show version information
  -t, --test      Run in test mode (uses temporary directory, safe for testing)
  --dry-run       Show what would be installed without actually doing it

Navigation:
  ↑/k, ↓/j        Navigate up/down
  Enter/Space     Select option
  Esc             Go back
  q               Quit
  d               Toggle details (during installation)

For more info: https://github.com/Gentleman-Programming/Gentleman.Dots`)
}
