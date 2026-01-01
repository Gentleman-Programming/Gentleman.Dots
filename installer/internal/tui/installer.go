package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Gentleman-Programming/Gentleman.Dots/installer/internal/system"
)

// StepError provides context about which step failed and why
type StepError struct {
	StepID      string
	StepName    string
	Description string
	Cause       error
}

func (e *StepError) Error() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Step '%s' failed\n", e.StepName))
	sb.WriteString(fmt.Sprintf("Description: %s\n", e.Description))
	if e.Cause != nil {
		sb.WriteString(fmt.Sprintf("\nDetails:\n%v", e.Cause))
	}
	return sb.String()
}

func (e *StepError) Unwrap() error {
	return e.Cause
}

// wrapStepError creates a detailed error for a step failure
func wrapStepError(stepID, stepName, description string, cause error) error {
	return &StepError{
		StepID:      stepID,
		StepName:    stepName,
		Description: description,
		Cause:       cause,
	}
}

// executeStep runs the actual installation for a step
func executeStep(stepID string, m *Model) error {
	switch stepID {
	case "backup":
		return stepBackupConfigs(m)
	case "clone":
		return stepCloneRepo(m)
	case "homebrew":
		return stepInstallHomebrew(m)
	case "deps":
		return stepInstallDeps(m)
	case "xcode":
		return stepInstallXcode(m)
	case "terminal":
		return stepInstallTerminal(m)
	case "font":
		return stepInstallFont(m)
	case "shell":
		return stepInstallShell(m)
	case "wm":
		return stepInstallWM(m)
	case "nvim":
		return stepInstallNvim(m)
	case "cleanup":
		return stepCleanup(m)
	default:
		return fmt.Errorf("unknown step: %s", stepID)
	}
}

func stepBackupConfigs(m *Model) error {
	stepID := "backup"
	if len(m.ExistingConfigs) == 0 {
		SendLog(stepID, "No existing configs to backup")
		return nil
	}

	SendLog(stepID, fmt.Sprintf("Backing up %d existing configs...", len(m.ExistingConfigs)))

	// Extract just the config keys from the ExistingConfigs slice
	configKeys := make([]string, len(m.ExistingConfigs))
	for i, config := range m.ExistingConfigs {
		configKeys[i] = config
		SendLog(stepID, fmt.Sprintf("  → %s", config))
	}

	backupDir, err := system.CreateBackup(configKeys)
	if err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	m.BackupDir = backupDir
	SendLog(stepID, fmt.Sprintf("✓ Backup created at: %s", backupDir))
	return nil
}

func stepCloneRepo(m *Model) error {
	stepID := "clone"

	// Check if already exists
	if _, err := os.Stat("Gentleman.Dots"); err == nil {
		SendLog(stepID, "Removing existing Gentleman.Dots directory...")
		result := system.RunWithLogs("rm -rf Gentleman.Dots", nil, func(line string) {
			SendLog(stepID, line)
		})
		if result.Error != nil {
			return wrapStepError("clone", "Clone Repository",
				"Failed to remove existing Gentleman.Dots directory",
				result.Error)
		}
	}

	SendLog(stepID, "Cloning repository from GitHub...")
	result := system.RunWithLogs("git clone --progress https://github.com/Gentleman-Programming/Gentleman.Dots.git Gentleman.Dots", nil, func(line string) {
		SendLog(stepID, line)
	})
	if result.Error != nil {
		return wrapStepError("clone", "Clone Repository",
			"Failed to clone the repository. Check your internet connection and git installation.",
			result.Error)
	}

	// Verify clone was successful
	if _, err := os.Stat("Gentleman.Dots"); os.IsNotExist(err) {
		return wrapStepError("clone", "Clone Repository",
			"Repository was cloned but directory not found",
			fmt.Errorf("Gentleman.Dots directory does not exist after clone"))
	}

	SendLog(stepID, "✓ Repository cloned successfully")
	return nil
}

func stepInstallHomebrew(m *Model) error {
	stepID := "homebrew"

	if system.CommandExists("brew") {
		SendLog(stepID, "Homebrew already installed, skipping...")
		return nil
	}

	SendLog(stepID, "Installing Homebrew package manager...")
	result := system.RunWithLogs(`/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"`, nil, func(line string) {
		SendLog(stepID, line)
	})
	if result.Error != nil {
		return wrapStepError("homebrew", "Install Homebrew",
			"Failed to install Homebrew package manager. Check your internet connection.",
			result.Error)
	}

	// Add to PATH
	homeDir := os.Getenv("HOME")
	brewPrefix := system.GetBrewPrefix()

	shellConfig := fmt.Sprintf(`eval "$(%s/bin/brew shellenv)"`, brewPrefix)

	SendLog(stepID, "Configuring shell to use Homebrew...")
	// Add to common shell configs
	for _, rcFile := range []string{".bashrc", ".zshrc"} {
		rcPath := filepath.Join(homeDir, rcFile)
		if f, err := os.OpenFile(rcPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
			f.WriteString("\n" + shellConfig + "\n")
			f.Close()
		}
	}

	// Source it now
	system.Run(shellConfig, nil)

	SendLog(stepID, "✓ Homebrew installed successfully")
	return nil
}

func stepInstallDeps(m *Model) error {
	if m.SystemInfo.OS == system.OSArch {
		result := system.RunSudo("pacman -Syu --noconfirm", nil)
		if result.Error != nil {
			return wrapStepError("deps", "Install Dependencies",
				"Failed to update Arch Linux packages",
				result.Error)
		}
		result = system.RunSudo("pacman -S --needed --noconfirm base-devel curl file git wget unzip fontconfig", nil)
		if result.Error != nil {
			return wrapStepError("deps", "Install Dependencies",
				"Failed to install base dependencies on Arch Linux",
				result.Error)
		}
		return nil
	}

	// Debian/Ubuntu
	result := system.RunSudo("apt-get update", nil)
	if result.Error != nil {
		return wrapStepError("deps", "Install Dependencies",
			"Failed to update apt package list",
			result.Error)
	}
	result = system.RunSudo("apt-get install -y build-essential curl file git unzip fontconfig", nil)
	if result.Error != nil {
		return wrapStepError("deps", "Install Dependencies",
			"Failed to install base dependencies on Debian/Ubuntu",
			result.Error)
	}
	return nil
}

func stepInstallXcode(m *Model) error {
	result := system.Run("xcode-select --install", nil)
	if result.Error != nil {
		// xcode-select returns error if already installed, which is fine
		if result.ExitCode == 1 && strings.Contains(result.Stderr, "already installed") {
			return nil
		}
		return wrapStepError("xcode", "Install Xcode CLI",
			"Failed to install Xcode Command Line Tools. You may need to install them manually from the App Store.",
			result.Error)
	}
	return nil
}

func stepInstallTerminal(m *Model) error {
	terminal := m.Choices.Terminal
	homeDir := os.Getenv("HOME")
	repoDir := "Gentleman.Dots"
	stepID := "terminal"

	switch terminal {
	case "alacritty":
		if !system.CommandExists("alacritty") {
			SendLog(stepID, "Installing Alacritty...")
			var result *system.ExecResult
			if m.SystemInfo.OS == system.OSArch {
				result = system.RunSudoWithLogs("pacman -S --noconfirm alacritty", nil, func(line string) {
					SendLog(stepID, line)
				})
			} else if m.SystemInfo.OS == system.OSMac {
				result = system.RunBrewWithLogs("install --cask alacritty", nil, func(line string) {
					SendLog(stepID, line)
				})
			} else {
				// Ubuntu/Debian - install software-properties-common first (provides add-apt-repository)
				SendLog(stepID, "Installing required packages for PPA support...")
				system.RunSudo("apt-get install -y software-properties-common", nil)
				system.RunSudo("add-apt-repository -y ppa:aslatter/ppa", nil)
				system.RunSudo("apt-get update", nil)
				result = system.RunSudoWithLogs("apt-get install -y alacritty", nil, func(line string) {
					SendLog(stepID, line)
				})
			}
			if result.Error != nil {
				return wrapStepError("terminal", "Install Alacritty",
					"Failed to install Alacritty terminal emulator",
					result.Error)
			}
		} else {
			SendLog(stepID, "Alacritty already installed")
		}
		SendLog(stepID, "Copying Alacritty configuration...")
		if err := system.EnsureDir(filepath.Join(homeDir, ".config/alacritty")); err != nil {
			return wrapStepError("terminal", "Install Alacritty",
				"Failed to create Alacritty config directory",
				err)
		}
		if err := system.CopyFile(filepath.Join(repoDir, "alacritty.toml"), filepath.Join(homeDir, ".config/alacritty/alacritty.toml")); err != nil {
			return wrapStepError("terminal", "Install Alacritty",
				"Failed to copy Alacritty configuration",
				err)
		}
		SendLog(stepID, "✓ Alacritty configured")

	case "wezterm":
		if !system.CommandExists("wezterm") {
			SendLog(stepID, "Installing WezTerm...")
			var result *system.ExecResult
			if m.SystemInfo.OS == system.OSArch {
				result = system.RunSudoWithLogs("pacman -S --noconfirm wezterm", nil, func(line string) {
					SendLog(stepID, line)
				})
			} else if m.SystemInfo.OS == system.OSMac {
				result = system.RunBrewWithLogs("install --cask wezterm", nil, func(line string) {
					SendLog(stepID, line)
				})
			} else {
				system.Run("brew tap wez/wezterm-linuxbrew", nil)
				result = system.RunBrewWithLogs("install wezterm", nil, func(line string) {
					SendLog(stepID, line)
				})
			}
			if result.Error != nil {
				return wrapStepError("terminal", "Install WezTerm",
					"Failed to install WezTerm terminal emulator",
					result.Error)
			}
		} else {
			SendLog(stepID, "WezTerm already installed")
		}
		SendLog(stepID, "Copying WezTerm configuration...")
		if err := system.EnsureDir(filepath.Join(homeDir, ".config/wezterm")); err != nil {
			return wrapStepError("terminal", "Install WezTerm",
				"Failed to create WezTerm config directory",
				err)
		}
		if err := system.CopyFile(filepath.Join(repoDir, ".wezterm.lua"), filepath.Join(homeDir, ".config/wezterm/wezterm.lua")); err != nil {
			return wrapStepError("terminal", "Install WezTerm",
				"Failed to copy WezTerm configuration",
				err)
		}
		SendLog(stepID, "✓ WezTerm configured")

	case "kitty":
		if !system.CommandExists("kitty") && m.SystemInfo.OS == system.OSMac {
			SendLog(stepID, "Installing Kitty...")
			result := system.RunBrewWithLogs("install --cask kitty", nil, func(line string) {
				SendLog(stepID, line)
			})
			if result.Error != nil {
				return wrapStepError("terminal", "Install Kitty",
					"Failed to install Kitty terminal emulator",
					result.Error)
			}
		} else {
			SendLog(stepID, "Kitty already installed")
		}
		SendLog(stepID, "Copying Kitty configuration...")
		if err := system.EnsureDir(filepath.Join(homeDir, ".config/kitty")); err != nil {
			return wrapStepError("terminal", "Install Kitty",
				"Failed to create Kitty config directory",
				err)
		}
		if err := system.CopyDir(filepath.Join(repoDir, "GentlemanKitty/*"), filepath.Join(homeDir, ".config/kitty/")); err != nil {
			return wrapStepError("terminal", "Install Kitty",
				"Failed to copy Kitty configuration",
				err)
		}
		SendLog(stepID, "✓ Kitty configured")

	case "ghostty":
		if !system.CommandExists("ghostty") {
			SendLog(stepID, "Installing Ghostty...")
			var result *system.ExecResult
			if m.SystemInfo.OS == system.OSArch {
				result = system.RunSudoWithLogs("pacman -S --noconfirm ghostty", nil, func(line string) {
					SendLog(stepID, line)
				})
			} else if m.SystemInfo.OS == system.OSMac {
				result = system.RunBrewWithLogs("install --cask ghostty", nil, func(line string) {
					SendLog(stepID, line)
				})
			} else {
				result = system.RunWithLogs(`/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/mkasberg/ghostty-ubuntu/HEAD/install.sh)"`, nil, func(line string) {
					SendLog(stepID, line)
				})
			}
			if result.Error != nil {
				return wrapStepError("terminal", "Install Ghostty",
					"Failed to install Ghostty terminal emulator",
					result.Error)
			}
		} else {
			SendLog(stepID, "Ghostty already installed")
		}
		SendLog(stepID, "Copying Ghostty configuration...")
		if err := system.EnsureDir(filepath.Join(homeDir, ".config/ghostty")); err != nil {
			return wrapStepError("terminal", "Install Ghostty",
				"Failed to create Ghostty config directory",
				err)
		}
		if err := system.CopyDir(filepath.Join(repoDir, "GentlemanGhostty/*"), filepath.Join(homeDir, ".config/ghostty/")); err != nil {
			return wrapStepError("terminal", "Install Ghostty",
				"Failed to copy Ghostty configuration",
				err)
		}
		SendLog(stepID, "✓ Ghostty configured")
	}

	return nil
}

func stepInstallFont(m *Model) error {
	homeDir := os.Getenv("HOME")
	stepID := "font"

	if m.SystemInfo.OS == system.OSMac {
		SendLog(stepID, "Installing Iosevka Term Nerd Font...")
		result := system.RunBrewWithLogs("install --cask font-iosevka-term-nerd-font", nil, func(line string) {
			SendLog(stepID, line)
		})
		if result.Error != nil {
			return wrapStepError("font", "Install Iosevka Nerd Font",
				"Failed to install font via Homebrew. Try installing manually from https://www.nerdfonts.com/",
				result.Error)
		}
		SendLog(stepID, "✓ Font installed")
		return nil
	}

	// Linux
	fontDir := filepath.Join(homeDir, ".local/share/fonts")
	SendLog(stepID, "Creating fonts directory...")
	if err := system.EnsureDir(fontDir); err != nil {
		return wrapStepError("font", "Install Iosevka Nerd Font",
			"Failed to create fonts directory",
			err)
	}

	SendLog(stepID, "Downloading Iosevka Term Nerd Font...")
	result := system.RunWithLogs(fmt.Sprintf("wget -O %s/IosevkaTerm.zip https://github.com/ryanoasis/nerd-fonts/releases/download/v3.3.0/IosevkaTerm.zip", fontDir), nil, func(line string) {
		SendLog(stepID, line)
	})
	if result.Error != nil {
		return wrapStepError("font", "Install Iosevka Nerd Font",
			"Failed to download font. Check your internet connection.",
			result.Error)
	}

	SendLog(stepID, "Extracting font archive...")
	result = system.RunWithLogs(fmt.Sprintf("unzip -o %s/IosevkaTerm.zip -d %s/", fontDir, fontDir), nil, func(line string) {
		SendLog(stepID, line)
	})
	if result.Error != nil {
		return wrapStepError("font", "Install Iosevka Nerd Font",
			"Failed to extract font archive",
			result.Error)
	}

	SendLog(stepID, "Updating font cache...")
	system.RunWithLogs("fc-cache -fv", nil, func(line string) {
		SendLog(stepID, line)
	})
	SendLog(stepID, "✓ Font installed")
	return nil
}

func stepInstallShell(m *Model) error {
	homeDir := os.Getenv("HOME")
	repoDir := "Gentleman.Dots"
	shell := m.Choices.Shell
	stepID := "shell"

	// Common dependencies
	SendLog(stepID, "Creating cache directories...")
	system.EnsureDir(filepath.Join(homeDir, ".cache/starship"))
	system.EnsureDir(filepath.Join(homeDir, ".cache/carapace"))
	system.EnsureDir(filepath.Join(homeDir, ".local/share/atuin"))

	switch shell {
	case "fish":
		SendLog(stepID, "Installing Fish shell and plugins...")
		result := system.RunBrewWithLogs("install fish carapace zoxide atuin starship", nil, func(line string) {
			SendLog(stepID, line)
		})
		if result.Error != nil {
			return wrapStepError("shell", "Install Fish",
				"Failed to install Fish shell and dependencies",
				result.Error)
		}
		SendLog(stepID, "Copying Fish configuration...")
		if err := system.CopyFile(filepath.Join(repoDir, "starship.toml"), filepath.Join(homeDir, ".config/starship.toml")); err != nil {
			return wrapStepError("shell", "Install Fish",
				"Failed to copy starship configuration",
				err)
		}
		if err := system.CopyDir(filepath.Join(repoDir, "GentlemanFish/fish"), filepath.Join(homeDir, ".config/")); err != nil {
			return wrapStepError("shell", "Install Fish",
				"Failed to copy Fish configuration",
				err)
		}
		SendLog(stepID, "✓ Fish shell configured")

	case "zsh":
		SendLog(stepID, "Installing Zsh and plugins...")
		result := system.RunBrewWithLogs("install zsh carapace zoxide atuin zsh-autosuggestions zsh-syntax-highlighting zsh-autocomplete powerlevel10k", nil, func(line string) {
			SendLog(stepID, line)
		})
		if result.Error != nil {
			return wrapStepError("shell", "Install Zsh",
				"Failed to install Zsh and plugins",
				result.Error)
		}
		SendLog(stepID, "Copying Zsh configuration...")
		if err := system.CopyFile(filepath.Join(repoDir, "GentlemanZsh/.zshrc"), filepath.Join(homeDir, ".zshrc")); err != nil {
			return wrapStepError("shell", "Install Zsh",
				"Failed to copy .zshrc configuration",
				err)
		}
		if err := system.CopyFile(filepath.Join(repoDir, "GentlemanZsh/.p10k.zsh"), filepath.Join(homeDir, ".p10k.zsh")); err != nil {
			return wrapStepError("shell", "Install Zsh",
				"Failed to copy Powerlevel10k configuration",
				err)
		}
		if err := system.CopyDir(filepath.Join(repoDir, "GentlemanZsh/.oh-my-zsh"), filepath.Join(homeDir, "/")); err != nil {
			return wrapStepError("shell", "Install Zsh",
				"Failed to copy Oh-My-Zsh directory",
				err)
		}
		SendLog(stepID, "✓ Zsh configured with Powerlevel10k")

	case "nushell":
		SendLog(stepID, "Installing Nushell and dependencies...")
		result := system.RunBrewWithLogs("install nushell carapace zoxide atuin jq bash starship", nil, func(line string) {
			SendLog(stepID, line)
		})
		if result.Error != nil {
			return wrapStepError("shell", "Install Nushell",
				"Failed to install Nushell and dependencies",
				result.Error)
		}
		SendLog(stepID, "Copying Nushell configuration...")
		if err := system.CopyFile(filepath.Join(repoDir, "starship.toml"), filepath.Join(homeDir, ".config/starship.toml")); err != nil {
			return wrapStepError("shell", "Install Nushell",
				"Failed to copy starship configuration",
				err)
		}
		if err := system.CopyFile(filepath.Join(repoDir, "bash-env-json"), filepath.Join(homeDir, ".config/bash-env-json")); err != nil {
			return wrapStepError("shell", "Install Nushell",
				"Failed to copy bash-env-json",
				err)
		}
		if err := system.CopyFile(filepath.Join(repoDir, "bash-env.nu"), filepath.Join(homeDir, ".config/bash-env.nu")); err != nil {
			return wrapStepError("shell", "Install Nushell",
				"Failed to copy bash-env.nu",
				err)
		}

		var nuDir string
		if runtime.GOOS == "darwin" {
			nuDir = filepath.Join(homeDir, "Library/Application Support/nushell")
		} else {
			nuDir = filepath.Join(homeDir, ".config/nushell")
		}
		if err := system.EnsureDir(nuDir); err != nil {
			return wrapStepError("shell", "Install Nushell",
				"Failed to create Nushell config directory",
				err)
		}
		if err := system.CopyDir(filepath.Join(repoDir, "GentlemanNushell/*"), nuDir+"/"); err != nil {
			return wrapStepError("shell", "Install Nushell",
				"Failed to copy Nushell configuration",
				err)
		}
		SendLog(stepID, "✓ Nushell configured")
	}

	return nil
}

func stepInstallWM(m *Model) error {
	homeDir := os.Getenv("HOME")
	repoDir := "Gentleman.Dots"
	wm := m.Choices.WindowMgr
	stepID := "wm"

	switch wm {
	case "tmux":
		if !system.CommandExists("tmux") {
			SendLog(stepID, "Installing Tmux...")
			result := system.RunBrewWithLogs("install tmux", nil, func(line string) {
				SendLog(stepID, line)
			})
			if result.Error != nil {
				return wrapStepError("wm", "Install Tmux",
					"Failed to install Tmux",
					result.Error)
			}
		} else {
			SendLog(stepID, "Tmux already installed")
		}

		// TPM
		tpmDir := filepath.Join(homeDir, ".tmux/plugins/tpm")
		if _, err := os.Stat(tpmDir); os.IsNotExist(err) {
			SendLog(stepID, "Cloning TPM (Tmux Plugin Manager)...")
			result := system.RunWithLogs(fmt.Sprintf("git clone https://github.com/tmux-plugins/tpm %s", tpmDir), nil, func(line string) {
				SendLog(stepID, line)
			})
			if result.Error != nil {
				return wrapStepError("wm", "Install Tmux",
					"Failed to clone TPM (Tmux Plugin Manager)",
					result.Error)
			}
		}

		SendLog(stepID, "Copying Tmux configuration...")
		if err := system.EnsureDir(filepath.Join(homeDir, ".tmux")); err != nil {
			return wrapStepError("wm", "Install Tmux",
				"Failed to create .tmux directory",
				err)
		}
		if err := system.CopyDir(filepath.Join(repoDir, "GentlemanTmux/plugins"), filepath.Join(homeDir, ".tmux/")); err != nil {
			return wrapStepError("wm", "Install Tmux",
				"Failed to copy Tmux plugins",
				err)
		}
		if err := system.CopyFile(filepath.Join(repoDir, "GentlemanTmux/tmux.conf"), filepath.Join(homeDir, ".tmux.conf")); err != nil {
			return wrapStepError("wm", "Install Tmux",
				"Failed to copy tmux.conf",
				err)
		}

		// Install plugins
		SendLog(stepID, "Installing Tmux plugins...")
		system.RunWithLogs(filepath.Join(homeDir, ".tmux/plugins/tpm/bin/install_plugins"), nil, func(line string) {
			SendLog(stepID, line)
		})
		SendLog(stepID, "✓ Tmux configured")

	case "zellij":
		if !system.CommandExists("zellij") {
			SendLog(stepID, "Installing Zellij...")
			result := system.RunBrewWithLogs("install zellij", nil, func(line string) {
				SendLog(stepID, line)
			})
			if result.Error != nil {
				return wrapStepError("wm", "Install Zellij",
					"Failed to install Zellij",
					result.Error)
			}
		} else {
			SendLog(stepID, "Zellij already installed")
		}

		SendLog(stepID, "Copying Zellij configuration...")
		zellijDir := filepath.Join(homeDir, ".config/zellij")
		if err := system.EnsureDir(zellijDir); err != nil {
			return wrapStepError("wm", "Install Zellij",
				"Failed to create Zellij config directory",
				err)
		}
		if err := system.CopyDir(filepath.Join(repoDir, "GentlemanZellij/zellij/*"), zellijDir+"/"); err != nil {
			return wrapStepError("wm", "Install Zellij",
				"Failed to copy Zellij configuration",
				err)
		}
		SendLog(stepID, "✓ Zellij configured")
	}

	return nil
}

func stepInstallNvim(m *Model) error {
	homeDir := os.Getenv("HOME")
	repoDir := "Gentleman.Dots"
	stepID := "nvim"

	// Obsidian path
	SendLog(stepID, "Creating Obsidian directories...")
	obsidianDir := filepath.Join(homeDir, ".config/obsidian")
	system.EnsureDir(obsidianDir)
	system.EnsureDir(filepath.Join(obsidianDir, "templates"))

	// Check Node.js
	if !system.CommandExists("node") {
		SendLog(stepID, "Installing Node.js...")
		result := system.RunBrewWithLogs("install node", nil, func(line string) {
			SendLog(stepID, line)
		})
		if result.Error != nil {
			return wrapStepError("nvim", "Install Neovim",
				"Failed to install Node.js (required for LSP servers)",
				result.Error)
		}
	} else {
		SendLog(stepID, "Node.js already installed")
	}

	// Install dependencies
	SendLog(stepID, "Installing Neovim and dependencies...")
	result := system.RunBrewWithLogs("install nvim git gcc fzf fd ripgrep coreutils bat curl lazygit tree-sitter", nil, func(line string) {
		SendLog(stepID, line)
	})
	if result.Error != nil {
		return wrapStepError("nvim", "Install Neovim",
			"Failed to install Neovim and dependencies",
			result.Error)
	}

	// Copy config
	SendLog(stepID, "Copying Neovim configuration...")
	nvimDir := filepath.Join(homeDir, ".config/nvim")
	if err := system.EnsureDir(nvimDir); err != nil {
		return wrapStepError("nvim", "Install Neovim",
			"Failed to create Neovim config directory",
			err)
	}
	if err := system.CopyDir(filepath.Join(repoDir, "GentlemanNvim/nvim/*"), nvimDir+"/"); err != nil {
		return wrapStepError("nvim", "Install Neovim",
			"Failed to copy Neovim configuration",
			err)
	}

	// Install Claude Code (optional, don't fail on error)
	SendLog(stepID, "Installing Claude Code (optional)...")
	system.RunWithLogs(`curl -fsSL https://claude.ai/install.sh | bash`, nil, func(line string) {
		SendLog(stepID, line)
	})

	// Install OpenCode (optional, don't fail on error)
	SendLog(stepID, "Installing OpenCode (optional)...")
	system.RunWithLogs(`curl -fsSL https://opencode.ai/install | bash`, nil, func(line string) {
		SendLog(stepID, line)
	})

	// Configure OpenCode
	SendLog(stepID, "Configuring OpenCode...")
	openCodeDir := filepath.Join(homeDir, ".config/opencode")
	system.EnsureDir(openCodeDir)
	system.EnsureDir(filepath.Join(openCodeDir, "themes"))
	system.CopyFile(filepath.Join(repoDir, "GentlemanOpenCode/opencode.json"), filepath.Join(openCodeDir, "opencode.json"))
	system.CopyFile(filepath.Join(repoDir, "GentlemanOpenCode/themes/gentleman.json"), filepath.Join(openCodeDir, "themes/gentleman.json"))

	SendLog(stepID, "✓ Neovim configured with Gentleman setup")
	return nil
}

func stepCleanup(m *Model) error {
	stepID := "cleanup"
	SendLog(stepID, "Removing temporary files...")
	// Only remove the cloned repo - no sudo needed
	result := system.Run("rm -rf Gentleman.Dots", nil)
	if result.Error != nil {
		// Non-critical error, just log it
		SendLog(stepID, "Warning: Could not remove temporary directory")
		return nil
	}
	SendLog(stepID, "✓ Cleanup complete")
	return nil
}
