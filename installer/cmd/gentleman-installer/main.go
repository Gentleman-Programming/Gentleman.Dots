package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Gentleman-Programming/Gentleman.Dots/installer/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

var Version = "1.0.0"

// CLI flags for non-interactive mode
type cliFlags struct {
	version        bool
	help           bool
	test           bool
	dryRun         bool
	nonInteractive bool
	terminal       string
	shell          string
	windowMgr      string
	nvim           bool
	font           bool
	backup         bool
	aiTools        string
	aiFramework    bool
	aiPreset       string
	aiModules      string
}

func parseFlags() *cliFlags {
	flags := &cliFlags{}

	flag.BoolVar(&flags.version, "version", false, "Show version information")
	flag.BoolVar(&flags.version, "v", false, "Show version information (shorthand)")
	flag.BoolVar(&flags.help, "help", false, "Show help message")
	flag.BoolVar(&flags.help, "h", false, "Show help message (shorthand)")
	flag.BoolVar(&flags.test, "test", false, "Run in test mode (uses temporary directory)")
	flag.BoolVar(&flags.test, "t", false, "Run in test mode (shorthand)")
	flag.BoolVar(&flags.dryRun, "dry-run", false, "Show what would be installed without doing it")
	flag.BoolVar(&flags.nonInteractive, "non-interactive", false, "Run without TUI, use CLI flags")
	flag.StringVar(&flags.terminal, "terminal", "", "Terminal: alacritty, wezterm, kitty, ghostty, none")
	flag.StringVar(&flags.shell, "shell", "", "Shell: fish, zsh, nushell")
	flag.StringVar(&flags.windowMgr, "wm", "", "Window manager: tmux, zellij, none")
	flag.BoolVar(&flags.nvim, "nvim", false, "Install Neovim configuration")
	flag.BoolVar(&flags.font, "font", false, "Install Nerd Font")
	flag.BoolVar(&flags.backup, "backup", true, "Backup existing configs (default: true)")
	flag.StringVar(&flags.aiTools, "ai-tools", "", "AI tools: claude,opencode,gemini,copilot (comma-separated)")
	flag.BoolVar(&flags.aiFramework, "ai-framework", false, "Install AI coding framework")
	flag.StringVar(&flags.aiPreset, "ai-preset", "", "Framework preset: minimal, frontend, backend, fullstack, data, complete")
	flag.StringVar(&flags.aiModules, "ai-modules", "", "Framework modules (comma-separated, use --list-modules for names)")

	flag.Parse()
	return flags
}

func main() {
	flags := parseFlags()

	if flags.version {
		fmt.Printf("gentleman.dots v%s\n", Version)
		os.Exit(0)
	}

	if flags.help {
		printHelp()
		os.Exit(0)
	}

	if flags.test {
		setupTestMode()
	}

	if flags.dryRun {
		os.Setenv("GENTLEMAN_DRY_RUN", "1")
		fmt.Println("🧪 Dry-run mode: No actual installations will be performed")
	}

	// Non-interactive mode: run installation directly with provided flags
	if flags.nonInteractive {
		if err := runNonInteractive(flags); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	// Interactive TUI mode
	model := tui.NewModel()
	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	tui.SetGlobalProgram(p)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running installer: %v\n", err)
		os.Exit(1)
	}
}

func runNonInteractive(flags *cliFlags) error {
	// Validate required flags
	if flags.shell == "" {
		return fmt.Errorf("--shell is required (fish, zsh, nushell)")
	}

	// Normalize inputs
	terminal := strings.ToLower(flags.terminal)
	shell := strings.ToLower(flags.shell)
	wm := strings.ToLower(flags.windowMgr)

	// Validate values
	validTerminals := map[string]bool{"alacritty": true, "wezterm": true, "kitty": true, "ghostty": true, "none": true, "": true}
	validShells := map[string]bool{"fish": true, "zsh": true, "nushell": true}
	validWMs := map[string]bool{"tmux": true, "zellij": true, "none": true, "": true}

	if !validTerminals[terminal] {
		return fmt.Errorf("invalid terminal: %s (valid: alacritty, wezterm, kitty, ghostty, none)", terminal)
	}
	if !validShells[shell] {
		return fmt.Errorf("invalid shell: %s (valid: fish, zsh, nushell)", shell)
	}
	if !validWMs[wm] {
		return fmt.Errorf("invalid window manager: %s (valid: tmux, zellij, none)", wm)
	}

	// Default empty values to "none"
	if terminal == "" {
		terminal = "none"
	}
	if wm == "" {
		wm = "none"
	}

	// Parse AI tools
	var aiTools []string
	if flags.aiTools != "" {
		validAITools := map[string]bool{"claude": true, "opencode": true, "gemini": true, "copilot": true}
		for _, tool := range strings.Split(flags.aiTools, ",") {
			tool = strings.TrimSpace(strings.ToLower(tool))
			if !validAITools[tool] {
				return fmt.Errorf("invalid AI tool: %s (valid: claude, opencode, gemini, copilot)", tool)
			}
			aiTools = append(aiTools, tool)
		}
	}

	// Validate AI preset
	aiPreset := strings.ToLower(flags.aiPreset)
	validPresets := map[string]bool{"minimal": true, "frontend": true, "backend": true, "fullstack": true, "data": true, "complete": true, "": true}
	if !validPresets[aiPreset] {
		return fmt.Errorf("invalid AI preset: %s (valid: minimal, frontend, backend, fullstack, data, complete)", aiPreset)
	}

	// Parse AI modules
	var aiModules []string
	if flags.aiModules != "" {
		for _, mod := range strings.Split(flags.aiModules, ",") {
			mod = strings.TrimSpace(mod)
			if mod != "" {
				aiModules = append(aiModules, mod)
			}
		}
	}

	// Determine if framework should be installed
	installFramework := flags.aiFramework || aiPreset != "" || len(aiModules) > 0

	// Create choices
	choices := tui.UserChoices{
		Terminal:           terminal,
		Shell:              shell,
		WindowMgr:          wm,
		InstallNvim:        flags.nvim,
		InstallFont:        flags.font,
		CreateBackup:       flags.backup,
		AITools:            aiTools,
		InstallAIFramework: installFramework,
		AIFrameworkPreset:  aiPreset,
		AIFrameworkModules: aiModules,
	}

	fmt.Println("🚀 Gentleman.Dots Non-Interactive Installer")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("  Terminal:    %s\n", choices.Terminal)
	fmt.Printf("  Shell:       %s\n", choices.Shell)
	fmt.Printf("  Window Mgr:  %s\n", choices.WindowMgr)
	fmt.Printf("  Neovim:      %v\n", choices.InstallNvim)
	fmt.Printf("  Font:        %v\n", choices.InstallFont)
	fmt.Printf("  Backup:      %v\n", choices.CreateBackup)
	if len(choices.AITools) > 0 {
		fmt.Printf("  AI Tools:    %s\n", strings.Join(choices.AITools, ", "))
	}
	if choices.InstallAIFramework {
		if choices.AIFrameworkPreset != "" {
			fmt.Printf("  AI Framework: preset=%s\n", choices.AIFrameworkPreset)
		} else if len(choices.AIFrameworkModules) > 0 {
			fmt.Printf("  AI Framework: modules=%s\n", strings.Join(choices.AIFrameworkModules, ","))
		} else {
			fmt.Printf("  AI Framework: yes\n")
		}
	}
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	// Run the installation
	return tui.RunNonInteractive(choices)
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

Interactive Mode (default):
  Just run 'gentleman.dots' to start the TUI installer.

Non-Interactive Mode:
  gentleman.dots --non-interactive --shell=<shell> [options]

Flags:
  -h, --help           Show this help message
  -v, --version        Show version information
  -t, --test           Run in test mode (uses temporary directory)
  --dry-run            Show what would be installed without doing it
  --non-interactive    Run without TUI, use CLI flags instead

Non-Interactive Options:
  --shell=<shell>      Shell to install (required): fish, zsh, nushell
  --terminal=<term>    Terminal: alacritty, wezterm, kitty, ghostty, none
  --wm=<wm>            Window manager: tmux, zellij, none
  --nvim               Install Neovim configuration
  --font               Install Nerd Font
  --backup=false       Disable config backup (default: true)

AI Options:
  --ai-tools=<tools>   AI tools (comma-separated): claude, opencode, gemini, copilot
  --ai-framework       Install AI coding framework
  --ai-preset=<name>   Framework preset: minimal, frontend, backend, fullstack, data, complete
  --ai-modules=<mods>  Framework modules (comma-separated, granular IDs)
                       Categories: scripts-*, hooks-*, agents-*, skill-*, commands-*
                       Atomic: sdd, mcp (any sub-item selects the whole category)

Examples:
  # Interactive TUI
  gentleman.dots

  # Non-interactive with Fish + Zellij + Neovim
  gentleman.dots --non-interactive --shell=fish --wm=zellij --nvim

  # Full setup with AI tools and framework
  gentleman.dots --non-interactive --shell=fish --nvim --ai-tools=claude,opencode,gemini,copilot --ai-preset=fullstack

  # Custom modules with granular skill IDs
  gentleman.dots --non-interactive --shell=zsh --ai-tools=claude --ai-framework \
    --ai-modules=scripts-project,hooks-security,skill-react-19,skill-typescript,sdd

  # Test mode with Zsh + Tmux (no terminal, no nvim)
  gentleman.dots --test --non-interactive --shell=zsh --wm=tmux

  # Verbose output (shows all command logs)
  GENTLEMAN_VERBOSE=1 gentleman.dots --non-interactive --shell=fish --nvim

Navigation (TUI mode):
  ↑/k, ↓/j        Navigate up/down
  Enter/Space     Select option
  Esc             Go back
  q               Quit
  d               Toggle details (during installation)

For more info: https://github.com/Gentleman-Programming/Gentleman.Dots`)
}
