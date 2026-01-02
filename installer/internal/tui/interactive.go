package tui

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/Gentleman-Programming/Gentleman.Dots/installer/internal/system"
	tea "github.com/charmbracelet/bubbletea"
)

// needsExecProcessMsg signals that we need to run tea.ExecProcess
type needsExecProcessMsg struct {
	stepID string
	cmd    *exec.Cmd
}

// runInteractiveStep creates a tea.Cmd that runs an interactive step
// This suspends the TUI and gives full terminal control to the process
func runInteractiveStep(stepID string, m *Model) tea.Cmd {
	return func() tea.Msg {
		script, err := getInteractiveScript(stepID, m)
		if err != nil {
			return execFinishedMsg{stepID: stepID, err: fmt.Errorf("failed to get script for %s: %w", stepID, err)}
		}

		// If no script needed (e.g., already installed), just succeed
		if script == "" {
			return execFinishedMsg{stepID: stepID, err: nil}
		}

		cmd, err := createTempScriptCommand(script)
		if err != nil {
			return execFinishedMsg{stepID: stepID, err: fmt.Errorf("failed to create script for %s: %w", stepID, err)}
		}

		// Return message that tells Update to use tea.ExecProcess
		return needsExecProcessMsg{stepID: stepID, cmd: cmd}
	}
}

// getInteractiveScript returns the bash script for interactive steps only
// Interactive steps are those that NEED user input (sudo password, chsh, etc)
func getInteractiveScript(stepID string, m *Model) (string, error) {
	switch stepID {
	case "homebrew":
		return getHomebrewScript(m)
	case "deps":
		return getDepsScript(m)
	case "terminal":
		return getTerminalScript(m)
	case "setshell":
		return getSetShellScript(m)
	default:
		return "", fmt.Errorf("unknown interactive step: %s", stepID)
	}
}

// getHomebrewScript returns script to install Homebrew (needs password on first install)
func getHomebrewScript(m *Model) (string, error) {
	if system.CommandExists("brew") {
		return "", nil // Already installed
	}

	brewPrefix := system.GetBrewPrefix()
	script := fmt.Sprintf(`#!/bin/bash
set -e
echo ""
echo "🍺 Installing Homebrew package manager..."
echo "   (You may be prompted for your password)"
echo ""
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

echo ""
echo "📝 Configuring shell to use Homebrew..."

# Add to shell configs
BREW_CONFIG='eval "$(%s/bin/brew shellenv)"'

for RC_FILE in "$HOME/.bashrc" "$HOME/.zshrc"; do
    if [ -f "$RC_FILE" ]; then
        if ! grep -q "brew shellenv" "$RC_FILE" 2>/dev/null; then
            echo "" >> "$RC_FILE"
            echo "$BREW_CONFIG" >> "$RC_FILE"
        fi
    fi
done

# Source it now
eval "$(%s/bin/brew shellenv)"

echo ""
echo "✅ Homebrew installed successfully!"
echo ""
echo "Press Enter to continue..."
read
`, brewPrefix, brewPrefix)

	return script, nil
}

// getDepsScript returns script to install dependencies on Linux (needs sudo)
func getDepsScript(m *Model) (string, error) {
	var script string

	if m.SystemInfo.OS == system.OSArch {
		script = `#!/bin/bash
set -e
echo ""
echo "🔄 Updating Arch Linux packages..."
echo "   (You may be prompted for your password)"
echo ""
sudo pacman -Syu --noconfirm
echo ""
echo "📦 Installing base dependencies..."
sudo pacman -S --needed --noconfirm base-devel curl file git wget unzip fontconfig
echo ""
echo "✅ Dependencies installed successfully!"
echo ""
echo "Press Enter to continue..."
read
`
	} else {
		// Debian/Ubuntu
		script = `#!/bin/bash
set -e
echo ""
echo "🔄 Updating apt package list..."
echo "   (You may be prompted for your password)"
echo ""
sudo apt-get update
echo ""
echo "📦 Installing base dependencies..."
sudo apt-get install -y build-essential curl file git unzip fontconfig
echo ""
echo "✅ Dependencies installed successfully!"
echo ""
echo "Press Enter to continue..."
read
`
	}

	return script, nil
}

// getTerminalScript returns script to install terminal on Linux (needs sudo)
func getTerminalScript(m *Model) (string, error) {
	terminal := m.Choices.Terminal
	homeDir := os.Getenv("HOME")

	var installCmd string
	var configCmd string

	switch terminal {
	case "alacritty":
		if system.CommandExists("alacritty") {
			installCmd = `echo "✓ Alacritty already installed"`
		} else if m.SystemInfo.OS == system.OSArch {
			installCmd = `sudo pacman -S --noconfirm alacritty`
		} else {
			// Debian/Ubuntu
			installCmd = `sudo apt-get install -y software-properties-common
sudo add-apt-repository -y ppa:aslatter/ppa
sudo apt-get update
sudo apt-get install -y alacritty`
		}
		configCmd = fmt.Sprintf(`mkdir -p "%s/.config/alacritty"
cp "Gentleman.Dots/alacritty.toml" "%s/.config/alacritty/alacritty.toml"`, homeDir, homeDir)

	case "wezterm":
		if system.CommandExists("wezterm") {
			installCmd = `echo "✓ WezTerm already installed"`
		} else if m.SystemInfo.OS == system.OSArch {
			installCmd = `sudo pacman -S --noconfirm wezterm`
		} else {
			// Debian uses brew, not interactive
			return "", nil
		}
		configCmd = fmt.Sprintf(`mkdir -p "%s/.config/wezterm"
cp "Gentleman.Dots/.wezterm.lua" "%s/.config/wezterm/wezterm.lua"`, homeDir, homeDir)

	case "ghostty":
		if system.CommandExists("ghostty") {
			installCmd = `echo "✓ Ghostty already installed"`
		} else if m.SystemInfo.OS == system.OSArch {
			installCmd = `sudo pacman -S --noconfirm ghostty`
		} else {
			// Debian uses install script
			installCmd = `curl -fsSL https://raw.githubusercontent.com/mkasberg/ghostty-ubuntu/HEAD/install.sh | bash`
		}
		configCmd = fmt.Sprintf(`mkdir -p "%s/.config/ghostty"
cp -r Gentleman.Dots/GentlemanGhostty/* "%s/.config/ghostty/"`, homeDir, homeDir)

	default:
		return "", nil
	}

	script := fmt.Sprintf(`#!/bin/bash
set -e
echo ""
echo "🖥️  Installing %s..."
echo "   (You may be prompted for your password)"
echo ""
%s
echo ""
echo "📝 Copying %s configuration..."
%s
echo ""
echo "✅ %s configured!"
echo ""
echo "Press Enter to continue..."
read
`, terminal, installCmd, terminal, configCmd, terminal)

	return script, nil
}

// getSetShellScript returns script to set the default shell (needs chsh password)
func getSetShellScript(m *Model) (string, error) {
	shell := m.Choices.Shell
	var shellCmd string

	switch shell {
	case "fish":
		shellCmd = "fish"
	case "zsh":
		shellCmd = "zsh"
	case "nushell":
		shellCmd = "nu"
	default:
		return "", fmt.Errorf("unknown shell: %s", shell)
	}

	brewPrefix := system.GetBrewPrefix()

	script := fmt.Sprintf(`#!/bin/bash
set -e

# Add brew to PATH for this script
export PATH="%s/bin:$PATH"

SHELL_PATH=$(which %s 2>/dev/null)

if [ -z "$SHELL_PATH" ]; then
    echo "❌ Shell '%s' not found in PATH"
    echo ""
    echo "Press Enter to continue..."
    read
    exit 1
fi

echo ""
echo "🐚 Setting $SHELL_PATH as your default shell..."
echo ""

# Check if shell is already in /etc/shells
if ! grep -q "^$SHELL_PATH$" /etc/shells 2>/dev/null; then
    echo "📝 Adding $SHELL_PATH to /etc/shells (requires sudo)..."
    echo "$SHELL_PATH" | sudo tee -a /etc/shells > /dev/null
fi

# Change shell
echo ""
echo "🔐 Changing default shell..."
echo "   (You may need to enter your password)"
echo ""
chsh -s "$SHELL_PATH"

echo ""
echo "✅ Default shell changed to $SHELL_PATH"
echo "   Please log out and log back in for changes to take effect."
echo ""
echo "Press Enter to continue..."
read
`, brewPrefix, shellCmd, shellCmd)

	return script, nil
}

// createTempScriptCommand creates a temporary bash script and returns a command to execute it
func createTempScriptCommand(script string) (*exec.Cmd, error) {
	// Create temp file
	tmpFile, err := os.CreateTemp("", "gentleman-install-*.sh")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp script: %w", err)
	}

	// Write script
	if _, err := tmpFile.WriteString(script); err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return nil, fmt.Errorf("failed to write script: %w", err)
	}
	tmpFile.Close()

	// Make executable
	if err := os.Chmod(tmpFile.Name(), 0755); err != nil {
		os.Remove(tmpFile.Name())
		return nil, fmt.Errorf("failed to make script executable: %w", err)
	}

	// Return command - use bash explicitly
	cmd := exec.Command("bash", tmpFile.Name())
	cmd.Env = os.Environ()

	return cmd, nil
}
