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
	case "setshell":
		return stepSetDefaultShell(m)
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

	// Termux doesn't use Homebrew - it uses pkg
	if m.SystemInfo.IsTermux {
		SendLog(stepID, "Skipping Homebrew (Termux uses pkg package manager)")
		return nil
	}

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
	stepID := "deps"

	// Termux: use pkg (no sudo needed)
	// Check both SystemInfo and Choices.OS for redundancy
	isTermux := m.SystemInfo.IsTermux || m.Choices.OS == "termux"
	if isTermux {
		SendLog(stepID, "Updating Termux packages...")
		result := system.RunPkgWithLogs("update", nil, func(line string) {
			SendLog(stepID, line)
		})
		if result.Error != nil {
			return wrapStepError("deps", "Install Dependencies",
				"Failed to update Termux packages",
				result.Error)
		}
		result = system.RunPkgWithLogs("upgrade -y", nil, func(line string) {
			SendLog(stepID, line)
		})
		if result.Error != nil {
			// Upgrade failures are not critical
			SendLog(stepID, "Warning: package upgrade had issues, continuing...")
		}
		SendLog(stepID, "Installing base dependencies...")
		result = system.RunPkgInstall("git curl", nil, func(line string) {
			SendLog(stepID, line)
		})
		if result.Error != nil {
			return wrapStepError("deps", "Install Dependencies",
				"Failed to install base dependencies on Termux",
				result.Error)
		}
		return nil
	}

	// Arch Linux
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

	// Fedora/RHEL
	if m.SystemInfo.OS == system.OSFedora {
		result := system.RunSudo("dnf check-update || true", nil) // dnf check-update returns 100 if updates available
		result = system.RunSudo("dnf install -y @development-tools curl file git wget unzip fontconfig", nil)
		if result.Error != nil {
			return wrapStepError("deps", "Install Dependencies",
				"Failed to install base dependencies on Fedora/RHEL",
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
			} else if m.SystemInfo.OS == system.OSFedora {
				// Fedora: install from dnf
				result = system.RunSudoWithLogs("dnf install -y alacritty", nil, func(line string) {
					SendLog(stepID, line)
				})
			} else if m.SystemInfo.OS == system.OSDebian || m.SystemInfo.OS == system.OSLinux {
				// Debian/Ubuntu: compile from source (PPAs are unreliable)
				SendLog(stepID, "Building Alacritty from source...")
				SendLog(stepID, "Installing build dependencies...")
				result = system.RunSudoWithLogs("apt-get install -y cmake pkg-config libfreetype6-dev libfontconfig1-dev libxcb-xfixes0-dev libxkbcommon-dev python3 gzip scdoc git curl", nil, func(line string) {
					SendLog(stepID, line)
				})
				if result.Error != nil {
					return wrapStepError("terminal", "Install Alacritty",
						"Failed to install build dependencies",
						result.Error)
				}
				// Install Rust/Cargo only for this build
				cargoPath := filepath.Join(homeDir, ".cargo/bin/cargo")
				if !system.CommandExists("cargo") && !system.CommandExists(cargoPath) {
					SendLog(stepID, "Installing Rust/Cargo toolchain...")
					result = system.RunWithLogs("curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y", nil, func(line string) {
						SendLog(stepID, line)
					})
					if result.Error != nil {
						return wrapStepError("terminal", "Install Alacritty",
							"Failed to install Rust",
							result.Error)
					}
					cargoPath = filepath.Join(homeDir, ".cargo/bin/cargo")
				}
				// Clone and build Alacritty
				SendLog(stepID, "Cloning Alacritty repository...")
				alacrittyDir := filepath.Join(os.TempDir(), "alacritty-build")
				os.RemoveAll(alacrittyDir)
				result = system.RunWithLogs(fmt.Sprintf("git clone https://github.com/alacritty/alacritty.git %s", alacrittyDir), nil, func(line string) {
					SendLog(stepID, line)
				})
				if result.Error != nil {
					return wrapStepError("terminal", "Install Alacritty",
						"Failed to clone Alacritty repository",
						result.Error)
				}
				SendLog(stepID, "Building Alacritty (this may take 5-10 minutes)...")
				if !system.CommandExists("cargo") {
					cargoPath = filepath.Join(homeDir, ".cargo/bin/cargo")
				} else {
					cargoPath = "cargo"
				}
				result = system.RunWithLogs(fmt.Sprintf("%s build --release --manifest-path %s/Cargo.toml", cargoPath, alacrittyDir), nil, func(line string) {
					SendLog(stepID, line)
				})
				if result.Error != nil {
					return wrapStepError("terminal", "Install Alacritty",
						"Failed to build Alacritty",
						result.Error)
				}
				SendLog(stepID, "Installing Alacritty binary...")
				result = system.RunSudoWithLogs(fmt.Sprintf("cp %s/target/release/alacritty /usr/local/bin/alacritty", alacrittyDir), nil, func(line string) {
					SendLog(stepID, line)
				})
				if result.Error != nil {
					return wrapStepError("terminal", "Install Alacritty",
						"Failed to install Alacritty binary",
						result.Error)
				}
				system.RunSudoWithLogs(fmt.Sprintf("cp %s/extra/linux/Alacritty.desktop /usr/share/applications/", alacrittyDir), nil, func(line string) {
					SendLog(stepID, line)
				})
				os.RemoveAll(alacrittyDir)
				SendLog(stepID, "✓ Alacritty built and installed from source")
			} else {
				return wrapStepError("terminal", "Install Alacritty",
					"Unsupported operating system for Alacritty installation",
					fmt.Errorf("OS type: %v", m.SystemInfo.OS))
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
			} else if m.SystemInfo.OS == system.OSFedora {
				// Fedora: enable COPR and install
				system.RunSudo("dnf copr enable -y wezfurlong/wezterm-nightly", nil)
				result = system.RunSudoWithLogs("dnf install -y wezterm", nil, func(line string) {
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
		if err := system.CopyDir(filepath.Join(repoDir, "GentlemanKitty"), filepath.Join(homeDir, ".config", "kitty")); err != nil {
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
			} else if m.SystemInfo.OS == system.OSFedora {
				// Fedora: enable COPR and install
				system.RunSudo("dnf copr enable -y pgdev/ghostty", nil)
				result = system.RunSudoWithLogs("dnf install -y ghostty", nil, func(line string) {
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
		if err := system.CopyDir(filepath.Join(repoDir, "GentlemanGhostty"), filepath.Join(homeDir, ".config", "ghostty")); err != nil {
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

	// Termux: fonts work differently - copy to ~/.termux/font.ttf
	isTermux := m.SystemInfo.IsTermux || m.Choices.OS == "termux"
	if isTermux {
		SendLog(stepID, "Downloading JetBrainsMono Nerd Font for Termux...")
		termuxDir := filepath.Join(homeDir, ".termux")
		if err := system.EnsureDir(termuxDir); err != nil {
			return wrapStepError("font", "Install Nerd Font",
				"Failed to create .termux directory",
				err)
		}

		// Download a single TTF file for Termux
		result := system.RunWithLogs(fmt.Sprintf("curl -fsSL -o %s/font.ttf https://github.com/ryanoasis/nerd-fonts/raw/HEAD/patched-fonts/JetBrainsMono/Ligatures/Regular/JetBrainsMonoNerdFont-Regular.ttf", termuxDir), nil, func(line string) {
			SendLog(stepID, line)
		})
		if result.Error != nil {
			return wrapStepError("font", "Install Nerd Font",
				"Failed to download font. Check your internet connection.",
				result.Error)
		}

		SendLog(stepID, "Reloading Termux settings...")
		system.Run("termux-reload-settings", nil)
		SendLog(stepID, "✓ Font installed - restart Termux to apply")
		return nil
	}

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
	result := system.RunWithLogs(fmt.Sprintf("curl -fsSL -o %s/IosevkaTerm.zip https://github.com/ryanoasis/nerd-fonts/releases/download/v3.3.0/IosevkaTerm.zip", fontDir), nil, func(line string) {
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
	SendLog(stepID, "Creating required directories...")
	system.EnsureDir(filepath.Join(homeDir, ".config"))
	system.EnsureDir(filepath.Join(homeDir, ".cache/starship"))
	system.EnsureDir(filepath.Join(homeDir, ".cache/carapace"))
	system.EnsureDir(filepath.Join(homeDir, ".local/share/atuin"))

	switch shell {
	case "fish":
		SendLog(stepID, "Installing Fish shell and plugins...")
		var result *system.ExecResult
		if m.SystemInfo.IsTermux {
			result = system.RunPkgInstall("fish starship zoxide", nil, func(line string) {
				SendLog(stepID, line)
			})
		} else {
			result = system.RunBrewWithLogs("install fish carapace zoxide atuin starship", nil, func(line string) {
				SendLog(stepID, line)
			})
		}
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
		if err := system.CopyDir(filepath.Join(repoDir, "GentlemanFish", "fish"), filepath.Join(homeDir, ".config", "fish")); err != nil {
			return wrapStepError("shell", "Install Fish",
				"Failed to copy Fish configuration",
				err)
		}
		// Patch config.fish based on WM choice
		SendLog(stepID, "Configuring shell for window manager...")
		if err := system.PatchFishForWM(filepath.Join(homeDir, ".config/fish/config.fish"), m.Choices.WindowMgr, m.Choices.InstallNvim); err != nil {
			return wrapStepError("shell", "Install Fish",
				"Failed to configure config.fish for window manager",
				err)
		}
		// Remove tmux.fish function if not using tmux
		if m.Choices.WindowMgr != "tmux" {
			os.Remove(filepath.Join(homeDir, ".config/fish/functions/tmux.fish"))
		}
		// Termux: Add fish to $PREFIX/etc/shells so tmux doesn't complain
		if m.SystemInfo.IsTermux {
			SendLog(stepID, "Adding fish to Termux shells...")
			prefix := os.Getenv("PREFIX")
			if prefix == "" {
				prefix = "/data/data/com.termux/files/usr"
			}
			shellsFile := filepath.Join(prefix, "etc", "shells")
			system.EnsureDir(filepath.Join(prefix, "etc"))
			f, err := os.OpenFile(shellsFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err == nil {
				f.WriteString(filepath.Join(prefix, "bin", "fish") + "\n")
				f.Close()
			}
		}
		SendLog(stepID, "✓ Fish shell configured")

	case "zsh":
		SendLog(stepID, "Installing Zsh and plugins...")
		var result *system.ExecResult
		if m.SystemInfo.IsTermux {
			// Termux has zsh in pkg, but plugins need to be installed differently
			result = system.RunPkgInstall("zsh starship zoxide", nil, func(line string) {
				SendLog(stepID, line)
			})
		} else {
			result = system.RunBrewWithLogs("install zsh carapace zoxide atuin zsh-autosuggestions zsh-syntax-highlighting zsh-autocomplete powerlevel10k", nil, func(line string) {
				SendLog(stepID, line)
			})
		}
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
		// Patch .zshrc based on WM choice
		SendLog(stepID, "Configuring shell for window manager...")
		if err := system.PatchZshForWM(filepath.Join(homeDir, ".zshrc"), m.Choices.WindowMgr, m.Choices.InstallNvim); err != nil {
			return wrapStepError("shell", "Install Zsh",
				"Failed to configure .zshrc for window manager",
				err)
		}
		if err := system.CopyFile(filepath.Join(repoDir, "GentlemanZsh/.p10k.zsh"), filepath.Join(homeDir, ".p10k.zsh")); err != nil {
			return wrapStepError("shell", "Install Zsh",
				"Failed to copy Powerlevel10k configuration",
				err)
		}
		if err := system.CopyDir(filepath.Join(repoDir, "GentlemanZsh", ".oh-my-zsh"), filepath.Join(homeDir, ".oh-my-zsh")); err != nil {
			return wrapStepError("shell", "Install Zsh",
				"Failed to copy Oh-My-Zsh directory",
				err)
		}
		// Termux: Add zsh to $PREFIX/etc/shells so tmux doesn't complain
		if m.SystemInfo.IsTermux {
			SendLog(stepID, "Adding zsh to Termux shells...")
			prefix := os.Getenv("PREFIX")
			if prefix == "" {
				prefix = "/data/data/com.termux/files/usr"
			}
			shellsFile := filepath.Join(prefix, "etc", "shells")
			system.EnsureDir(filepath.Join(prefix, "etc"))
			f, err := os.OpenFile(shellsFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err == nil {
				f.WriteString(filepath.Join(prefix, "bin", "zsh") + "\n")
				f.Close()
			}
		}
		SendLog(stepID, "✓ Zsh configured with Powerlevel10k")

	case "nushell":
		SendLog(stepID, "Installing Nushell and dependencies...")
		var result *system.ExecResult
		if m.SystemInfo.IsTermux {
			result = system.RunPkgInstall("nushell starship zoxide jq", nil, func(line string) {
				SendLog(stepID, line)
			})
		} else {
			result = system.RunBrewWithLogs("install nushell carapace zoxide atuin jq bash starship", nil, func(line string) {
				SendLog(stepID, line)
			})
		}
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
		if err := system.CopyDir(filepath.Join(repoDir, "GentlemanNushell"), nuDir); err != nil {
			return wrapStepError("shell", "Install Nushell",
				"Failed to copy Nushell configuration",
				err)
		}
		// Patch config.nu based on WM choice
		SendLog(stepID, "Configuring shell for window manager...")
		if err := system.PatchNushellForWM(filepath.Join(nuDir, "config.nu"), m.Choices.WindowMgr); err != nil {
			return wrapStepError("shell", "Install Nushell",
				"Failed to configure config.nu for window manager",
				err)
		}
		// Termux: Add nu to $PREFIX/etc/shells so tmux doesn't complain
		if m.SystemInfo.IsTermux {
			SendLog(stepID, "Adding nushell to Termux shells...")
			prefix := os.Getenv("PREFIX")
			if prefix == "" {
				prefix = "/data/data/com.termux/files/usr"
			}
			shellsFile := filepath.Join(prefix, "etc", "shells")
			system.EnsureDir(filepath.Join(prefix, "etc"))
			f, err := os.OpenFile(shellsFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err == nil {
				f.WriteString(filepath.Join(prefix, "bin", "nu") + "\n")
				f.Close()
			}
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
			var result *system.ExecResult
			if m.SystemInfo.IsTermux {
				result = system.RunPkgInstall("tmux", nil, func(line string) {
					SendLog(stepID, line)
				})
			} else {
				result = system.RunBrewWithLogs("install tmux", nil, func(line string) {
					SendLog(stepID, line)
				})
			}
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
		if err := system.CopyDir(filepath.Join(repoDir, "GentlemanTmux", "plugins"), filepath.Join(homeDir, ".tmux", "plugins")); err != nil {
			return wrapStepError("wm", "Install Tmux",
				"Failed to copy Tmux plugins",
				err)
		}
		if err := system.CopyFile(filepath.Join(repoDir, "GentlemanTmux/tmux.conf"), filepath.Join(homeDir, ".tmux.conf")); err != nil {
			return wrapStepError("wm", "Install Tmux",
				"Failed to copy tmux.conf",
				err)
		}

		// Configure tmux to use the user's chosen shell
		SendLog(stepID, "Configuring tmux default shell...")
		tmuxConfPath := filepath.Join(homeDir, ".tmux.conf")
		shellName := ""
		switch m.Choices.Shell {
		case "fish":
			shellName = "fish"
		case "zsh":
			shellName = "zsh"
		case "nushell":
			shellName = "nu"
		}
		if shellName != "" {
			// Find the full path to the shell
			shellFullPath := ""
			if m.SystemInfo.IsTermux {
				// In Termux, construct the path directly (which command has issues)
				prefix := os.Getenv("PREFIX")
				if prefix == "" {
					prefix = "/data/data/com.termux/files/usr"
				}
				shellFullPath = filepath.Join(prefix, "bin", shellName)
			} else {
				result := system.Run(fmt.Sprintf("which %s", shellName), nil)
				if result.Error == nil && result.Output != "" {
					shellFullPath = strings.TrimSpace(result.Output)
				}
			}
			if shellFullPath == "" {
				shellFullPath = shellName // Fallback
			}

			// Replace placeholder in tmux.conf with actual shell config
			content, err := os.ReadFile(tmuxConfPath)
			if err == nil {
				shellConfig := fmt.Sprintf("set -g default-command \"%s\"\nset -g default-shell \"%s\"", shellFullPath, shellFullPath)
				newContent := strings.Replace(string(content), "# GENTLEMAN_DEFAULT_SHELL", shellConfig, 1)
				os.WriteFile(tmuxConfPath, []byte(newContent), 0644)
			}
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
			var result *system.ExecResult
			if m.SystemInfo.IsTermux {
				result = system.RunPkgInstall("zellij", nil, func(line string) {
					SendLog(stepID, line)
				})
			} else {
				result = system.RunBrewWithLogs("install zellij", nil, func(line string) {
					SendLog(stepID, line)
				})
			}
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
		if err := system.CopyDir(filepath.Join(repoDir, "GentlemanZellij", "zellij"), zellijDir); err != nil {
			return wrapStepError("wm", "Install Zellij",
				"Failed to copy Zellij configuration",
				err)
		}

		// Configure zellij to use the user's chosen shell
		SendLog(stepID, "Configuring zellij default shell...")
		zellijConfPath := filepath.Join(zellijDir, "config.kdl")
		shellPath := ""
		switch m.Choices.Shell {
		case "fish":
			shellPath = "fish"
		case "zsh":
			shellPath = "zsh"
		case "nushell":
			shellPath = "nu"
		}
		if shellPath != "" {
			// Append default_shell config to zellij config.kdl
			f, err := os.OpenFile(zellijConfPath, os.O_APPEND|os.O_WRONLY, 0644)
			if err == nil {
				f.WriteString(fmt.Sprintf("\n// Default shell (configured by Gentleman.Dots)\ndefault_shell \"%s\"\n", shellPath))
				f.Close()
			}
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
		var result *system.ExecResult
		if m.SystemInfo.IsTermux {
			result = system.RunPkgInstall("nodejs", nil, func(line string) {
				SendLog(stepID, line)
			})
		} else {
			result = system.RunBrewWithLogs("install node", nil, func(line string) {
				SendLog(stepID, line)
			})
		}
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
	var result *system.ExecResult
	if m.SystemInfo.IsTermux {
		// Termux package names (neovim instead of nvim, clang instead of gcc)
		result = system.RunPkgInstall("neovim git clang fzf fd ripgrep bat curl lazygit", nil, func(line string) {
			SendLog(stepID, line)
		})
	} else {
		result = system.RunBrewWithLogs("install nvim git gcc fzf fd ripgrep coreutils bat curl lazygit tree-sitter", nil, func(line string) {
			SendLog(stepID, line)
		})
	}
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
	// Copy nvim config directory
	srcNvim := filepath.Join(repoDir, "GentlemanNvim", "nvim")
	if err := system.CopyDir(srcNvim, nvimDir); err != nil {
		return wrapStepError("nvim", "Install Neovim",
			"Failed to copy Neovim configuration",
			err)
	}

	// Install Claude Code (optional, don't fail on error)
	// Skip on Termux - Claude Code doesn't support Android
	if !m.SystemInfo.IsTermux {
		SendLog(stepID, "Installing Claude Code (optional)...")
		system.RunWithLogs(`curl -fsSL https://claude.ai/install.sh | bash`, nil, func(line string) {
			SendLog(stepID, line)
		})
	} else {
		SendLog(stepID, "Skipping Claude Code (not supported on Termux)")
	}

	// Configure Claude Code
	SendLog(stepID, "Configuring Claude Code...")
	claudeDir := filepath.Join(homeDir, ".claude")
	system.EnsureDir(claudeDir)
	system.EnsureDir(filepath.Join(claudeDir, "output-styles"))
	system.EnsureDir(filepath.Join(claudeDir, "skills"))
	system.CopyFile(filepath.Join(repoDir, "GentlemanClaude/CLAUDE.md"), filepath.Join(claudeDir, "CLAUDE.md"))
	system.CopyFile(filepath.Join(repoDir, "GentlemanClaude/settings.json"), filepath.Join(claudeDir, "settings.json"))
	system.CopyFile(filepath.Join(repoDir, "GentlemanClaude/statusline.sh"), filepath.Join(claudeDir, "statusline.sh"))
	system.Run(fmt.Sprintf("chmod +x %s", filepath.Join(claudeDir, "statusline.sh")), nil)
	system.CopyFile(filepath.Join(repoDir, "GentlemanClaude/output-styles/gentleman.md"), filepath.Join(claudeDir, "output-styles/gentleman.md"))
	system.CopyFile(filepath.Join(repoDir, "GentlemanClaude/mcp-servers.template.json"), filepath.Join(claudeDir, "mcp-servers.template.json"))
	system.CopyFile(filepath.Join(repoDir, "GentlemanClaude/tweakcc-theme.json"), filepath.Join(claudeDir, "tweakcc-theme.json"))
	// Copy skills (excluding prowler-* which are work-specific)
	skillsToCopy := []string{"ai-sdk-5", "django-drf", "nextjs-15", "playwright", "pytest", "react-19", "tailwind-4", "typescript", "zod-4", "zustand-5"}
	for _, skill := range skillsToCopy {
		skillSrc := filepath.Join(repoDir, "GentlemanClaude/skills", skill)
		skillDst := filepath.Join(claudeDir, "skills", skill)
		system.EnsureDir(skillDst)
		system.CopyFile(filepath.Join(skillSrc, "SKILL.md"), filepath.Join(skillDst, "SKILL.md"))
	}
	SendLog(stepID, "⚙️ Copied CLAUDE.md")
	SendLog(stepID, "📊 Copied statusline.sh")
	SendLog(stepID, "🎨 Copied output styles")
	SendLog(stepID, "🧠 Copied Claude skills")

	// Apply tweakcc theme (only if Claude Code was installed)
	if !m.SystemInfo.IsTermux {
		SendLog(stepID, "Applying tweakcc theme...")
		result := system.Run("npx tweakcc --apply", nil)
		if result.Error == nil {
			SendLog(stepID, "🎨 Applied tweakcc theme")
		} else {
			SendLog(stepID, "⚠️ Could not apply tweakcc theme (run 'npx tweakcc --apply' manually)")
		}
	}

	// Install OpenCode (optional, don't fail on error)
	// Skip on Termux - OpenCode doesn't support Android
	if !m.SystemInfo.IsTermux {
		SendLog(stepID, "Installing OpenCode (optional)...")
		system.RunWithLogs(`curl -fsSL https://opencode.ai/install | bash`, nil, func(line string) {
			SendLog(stepID, line)
		})
	} else {
		SendLog(stepID, "Skipping OpenCode (not supported on Termux)")
	}

	// Configure OpenCode
	SendLog(stepID, "Configuring OpenCode...")
	openCodeDir := filepath.Join(homeDir, ".config/opencode")
	system.EnsureDir(openCodeDir)
	system.EnsureDir(filepath.Join(openCodeDir, "themes"))
	system.EnsureDir(filepath.Join(openCodeDir, "skill"))
	system.CopyFile(filepath.Join(repoDir, "GentlemanOpenCode/opencode.json"), filepath.Join(openCodeDir, "opencode.json"))
	system.CopyFile(filepath.Join(repoDir, "GentlemanOpenCode/themes/gentleman.json"), filepath.Join(openCodeDir, "themes/gentleman.json"))
	system.CopyDir(filepath.Join(repoDir, "GentlemanOpenCode", "skill"), filepath.Join(openCodeDir, "skill"))
	SendLog(stepID, "🧠 Copied OpenCode skills")

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

// stepSetDefaultShell sets the selected shell as the user's default shell
// In non-interactive mode, this handles Termux specially (via .bashrc)
// and attempts to set the shell on other systems if possible
func stepSetDefaultShell(m *Model) error {
	stepID := "setshell"
	shell := m.Choices.Shell
	homeDir := os.Getenv("HOME")

	var shellCmd string
	switch shell {
	case "fish":
		shellCmd = "fish"
	case "zsh":
		shellCmd = "zsh"
	case "nushell":
		shellCmd = "nu"
	default:
		SendLog(stepID, fmt.Sprintf("Unknown shell: %s, skipping", shell))
		return nil
	}

	// Termux: no chsh available, modify .bashrc to auto-start shell
	if m.SystemInfo.IsTermux {
		SendLog(stepID, "Configuring shell auto-start for Termux...")

		// Find the shell path
		shellPath := system.Run(fmt.Sprintf("which %s", shellCmd), nil)
		if shellPath.Error != nil || strings.TrimSpace(shellPath.Output) == "" {
			SendLog(stepID, fmt.Sprintf("Shell '%s' not found in PATH, skipping", shellCmd))
			return nil
		}
		shellPathStr := strings.TrimSpace(shellPath.Output)

		// Read existing .bashrc
		bashrcPath := filepath.Join(homeDir, ".bashrc")
		existingContent := ""
		if data, err := os.ReadFile(bashrcPath); err == nil {
			existingContent = string(data)
		}

		// Check if already configured
		if strings.Contains(existingContent, "# Gentleman.Dots shell auto-start") {
			SendLog(stepID, "Shell auto-start already configured in ~/.bashrc")
			return nil
		}

		// Append auto-start configuration
		autoStartConfig := fmt.Sprintf(`
# Gentleman.Dots shell auto-start
if [ -x "%s" ] && [ -z "$GENTLEMANDOTS_SHELL_STARTED" ]; then
    export GENTLEMANDOTS_SHELL_STARTED=1
    exec %s
fi
`, shellPathStr, shellPathStr)

		f, err := os.OpenFile(bashrcPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return wrapStepError("setshell", "Set Default Shell",
				"Failed to open ~/.bashrc for writing",
				err)
		}
		defer f.Close()

		if _, err := f.WriteString(autoStartConfig); err != nil {
			return wrapStepError("setshell", "Set Default Shell",
				"Failed to write shell auto-start to ~/.bashrc",
				err)
		}

		SendLog(stepID, fmt.Sprintf("✓ Configured %s to auto-start in ~/.bashrc", shell))
		SendLog(stepID, "Close and reopen Termux for changes to take effect")
		return nil
	}

	// Non-Termux: Try to set shell using sudo usermod (works if NOPASSWD configured)
	// Find the shell path first
	shellPath := system.Run(fmt.Sprintf("which %s", shellCmd), nil)
	if shellPath.Error != nil || strings.TrimSpace(shellPath.Output) == "" {
		SendLog(stepID, fmt.Sprintf("Shell '%s' not found in PATH, skipping", shellCmd))
		return nil
	}
	shellPathStr := strings.TrimSpace(shellPath.Output)

	// Get current username
	currentUser := os.Getenv("USER")
	if currentUser == "" {
		currentUser = os.Getenv("LOGNAME")
	}
	if currentUser == "" {
		// Fallback to whoami command (useful in Docker containers)
		whoamiResult := system.Run("whoami", nil)
		if whoamiResult.Error == nil {
			currentUser = strings.TrimSpace(whoamiResult.Output)
		}
	}
	if currentUser == "" {
		SendLog(stepID, "Could not determine current user, skipping shell change")
		return nil
	}

	// First, ensure shell is in /etc/shells
	SendLog(stepID, fmt.Sprintf("Adding %s to /etc/shells if needed...", shellPathStr))
	checkShells := system.Run(fmt.Sprintf("grep -q '^%s$' /etc/shells", shellPathStr), nil)
	if checkShells.Error != nil {
		// Shell not in /etc/shells, try to add it
		addResult := system.RunSudo(fmt.Sprintf("sh -c 'echo \"%s\" >> /etc/shells'", shellPathStr), nil)
		if addResult.Error != nil {
			SendLog(stepID, fmt.Sprintf("Could not add %s to /etc/shells (may need manual setup)", shellPathStr))
		}
	}

	// Try sudo usermod first (more reliable than chsh in scripts)
	SendLog(stepID, fmt.Sprintf("Setting %s as default shell for %s...", shell, currentUser))
	result := system.RunSudo(fmt.Sprintf("usermod -s %s %s", shellPathStr, currentUser), nil)
	if result.Error != nil {
		// usermod failed, try chsh as fallback
		SendLog(stepID, "usermod failed, trying chsh...")
		result = system.RunSudo(fmt.Sprintf("chsh -s %s %s", shellPathStr, currentUser), nil)
		if result.Error != nil {
			// Both failed - not critical, just inform user
			SendLog(stepID, fmt.Sprintf("Could not set default shell automatically"))
			SendLog(stepID, fmt.Sprintf("Run manually: chsh -s %s", shellPathStr))
			return nil
		}
	}

	SendLog(stepID, fmt.Sprintf("✓ Default shell set to %s", shell))
	SendLog(stepID, "Log out and log back in for changes to take effect")
	return nil
}
