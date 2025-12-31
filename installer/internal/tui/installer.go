package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/Gentleman-Programming/Gentleman.Dots/installer/internal/system"
)

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
	if len(m.ExistingConfigs) == 0 {
		return nil
	}

	// Extract just the config keys from the ExistingConfigs slice
	configKeys := make([]string, len(m.ExistingConfigs))
	for i, config := range m.ExistingConfigs {
		configKeys[i] = config
	}

	backupDir, err := system.CreateBackup(configKeys)
	if err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	m.BackupDir = backupDir
	return nil
}

func stepCloneRepo(m *Model) error {
	// Check if already exists
	if _, err := os.Stat("Gentleman.Dots"); err == nil {
		// Remove existing
		system.Run("rm -rf Gentleman.Dots", nil)
	}

	result := system.Run("git clone https://github.com/Gentleman-Programming/Gentleman.Dots.git Gentleman.Dots", nil)
	if result.Error != nil {
		return fmt.Errorf("failed to clone repository: %w", result.Error)
	}

	return nil
}

func stepInstallHomebrew(m *Model) error {
	if system.CommandExists("brew") {
		return nil
	}

	result := system.Run(`/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"`, nil)
	if result.Error != nil {
		return fmt.Errorf("failed to install Homebrew: %w", result.Error)
	}

	// Add to PATH
	homeDir := os.Getenv("HOME")
	brewPrefix := system.GetBrewPrefix()

	shellConfig := fmt.Sprintf(`eval "$(%s/bin/brew shellenv)"`, brewPrefix)

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

	return nil
}

func stepInstallDeps(m *Model) error {
	if m.SystemInfo.OS == system.OSArch {
		result := system.RunSudo("pacman -Syu --noconfirm", nil)
		if result.Error != nil {
			return result.Error
		}
		result = system.RunSudo("pacman -S --needed --noconfirm base-devel curl file git wget unzip fontconfig", nil)
		return result.Error
	}

	// Debian/Ubuntu
	result := system.RunSudo("apt-get update", nil)
	if result.Error != nil {
		return result.Error
	}
	result = system.RunSudo("apt-get install -y build-essential curl file git unzip fontconfig", nil)
	return result.Error
}

func stepInstallXcode(m *Model) error {
	result := system.Run("xcode-select --install", nil)
	return result.Error
}

func stepInstallTerminal(m *Model) error {
	terminal := m.Choices.Terminal
	homeDir := os.Getenv("HOME")
	repoDir := "Gentleman.Dots"

	switch terminal {
	case "alacritty":
		if !system.CommandExists("alacritty") {
			var result *system.ExecResult
			if m.SystemInfo.OS == system.OSArch {
				result = system.RunSudo("pacman -S --noconfirm alacritty", nil)
			} else if m.SystemInfo.OS == system.OSMac {
				result = system.RunBrew("install --cask alacritty", nil)
			} else {
				system.RunSudo("add-apt-repository -y ppa:aslatter/ppa", nil)
				system.RunSudo("apt update", nil)
				result = system.RunSudo("apt install -y alacritty", nil)
			}
			if result.Error != nil {
				return result.Error
			}
		}
		system.EnsureDir(filepath.Join(homeDir, ".config/alacritty"))
		system.CopyFile(filepath.Join(repoDir, "alacritty.toml"), filepath.Join(homeDir, ".config/alacritty/alacritty.toml"))

	case "wezterm":
		if !system.CommandExists("wezterm") {
			var result *system.ExecResult
			if m.SystemInfo.OS == system.OSArch {
				result = system.RunSudo("pacman -S --noconfirm wezterm", nil)
			} else if m.SystemInfo.OS == system.OSMac {
				result = system.RunBrew("install --cask wezterm", nil)
			} else {
				system.Run("brew tap wez/wezterm-linuxbrew", nil)
				result = system.RunBrew("install wezterm", nil)
			}
			if result.Error != nil {
				return result.Error
			}
		}
		system.EnsureDir(filepath.Join(homeDir, ".config/wezterm"))
		system.CopyFile(filepath.Join(repoDir, ".wezterm.lua"), filepath.Join(homeDir, ".config/wezterm/wezterm.lua"))

	case "kitty":
		if !system.CommandExists("kitty") && m.SystemInfo.OS == system.OSMac {
			result := system.RunBrew("install --cask kitty", nil)
			if result.Error != nil {
				return result.Error
			}
		}
		system.EnsureDir(filepath.Join(homeDir, ".config/kitty"))
		system.CopyDir(filepath.Join(repoDir, "GentlemanKitty/*"), filepath.Join(homeDir, ".config/kitty/"))

	case "ghostty":
		if !system.CommandExists("ghostty") {
			var result *system.ExecResult
			if m.SystemInfo.OS == system.OSArch {
				result = system.RunSudo("pacman -S --noconfirm ghostty", nil)
			} else if m.SystemInfo.OS == system.OSMac {
				result = system.RunBrew("install --cask ghostty", nil)
			} else {
				result = system.Run(`/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/mkasberg/ghostty-ubuntu/HEAD/install.sh)"`, nil)
			}
			if result.Error != nil {
				return result.Error
			}
		}
		system.EnsureDir(filepath.Join(homeDir, ".config/ghostty"))
		system.CopyDir(filepath.Join(repoDir, "GentlemanGhostty/*"), filepath.Join(homeDir, ".config/ghostty/"))
	}

	return nil
}

func stepInstallFont(m *Model) error {
	homeDir := os.Getenv("HOME")

	if m.SystemInfo.OS == system.OSMac {
		result := system.RunBrew("install --cask font-iosevka-term-nerd-font", nil)
		return result.Error
	}

	// Linux
	fontDir := filepath.Join(homeDir, ".local/share/fonts")
	system.EnsureDir(fontDir)

	result := system.Run(fmt.Sprintf("wget -O %s/IosevkaTerm.zip https://github.com/ryanoasis/nerd-fonts/releases/download/v3.3.0/IosevkaTerm.zip", fontDir), nil)
	if result.Error != nil {
		return result.Error
	}

	result = system.Run(fmt.Sprintf("unzip -o %s/IosevkaTerm.zip -d %s/", fontDir, fontDir), nil)
	if result.Error != nil {
		return result.Error
	}

	system.Run("fc-cache -fv", nil)
	return nil
}

func stepInstallShell(m *Model) error {
	homeDir := os.Getenv("HOME")
	repoDir := "Gentleman.Dots"
	shell := m.Choices.Shell

	// Common dependencies
	system.EnsureDir(filepath.Join(homeDir, ".cache/starship"))
	system.EnsureDir(filepath.Join(homeDir, ".cache/carapace"))
	system.EnsureDir(filepath.Join(homeDir, ".local/share/atuin"))

	switch shell {
	case "fish":
		result := system.RunBrew("install fish carapace zoxide atuin starship", nil)
		if result.Error != nil {
			return result.Error
		}
		system.CopyFile(filepath.Join(repoDir, "starship.toml"), filepath.Join(homeDir, ".config/starship.toml"))
		system.CopyDir(filepath.Join(repoDir, "GentlemanFish/fish"), filepath.Join(homeDir, ".config/"))

	case "zsh":
		result := system.RunBrew("install zsh carapace zoxide atuin zsh-autosuggestions zsh-syntax-highlighting zsh-autocomplete powerlevel10k", nil)
		if result.Error != nil {
			return result.Error
		}
		system.CopyFile(filepath.Join(repoDir, "GentlemanZsh/.zshrc"), filepath.Join(homeDir, ".zshrc"))
		system.CopyFile(filepath.Join(repoDir, "GentlemanZsh/.p10k.zsh"), filepath.Join(homeDir, ".p10k.zsh"))
		system.CopyDir(filepath.Join(repoDir, "GentlemanZsh/.oh-my-zsh"), filepath.Join(homeDir, "/"))

	case "nushell":
		result := system.RunBrew("install nushell carapace zoxide atuin jq bash starship", nil)
		if result.Error != nil {
			return result.Error
		}
		system.CopyFile(filepath.Join(repoDir, "starship.toml"), filepath.Join(homeDir, ".config/starship.toml"))
		system.CopyFile(filepath.Join(repoDir, "bash-env-json"), filepath.Join(homeDir, ".config/bash-env-json"))
		system.CopyFile(filepath.Join(repoDir, "bash-env.nu"), filepath.Join(homeDir, ".config/bash-env.nu"))

		if runtime.GOOS == "darwin" {
			nuDir := filepath.Join(homeDir, "Library/Application Support/nushell")
			system.EnsureDir(nuDir)
			system.CopyDir(filepath.Join(repoDir, "GentlemanNushell/*"), nuDir+"/")
		} else {
			nuDir := filepath.Join(homeDir, ".config/nushell")
			system.EnsureDir(nuDir)
			system.CopyDir(filepath.Join(repoDir, "GentlemanNushell/*"), nuDir+"/")
		}
	}

	return nil
}

func stepInstallWM(m *Model) error {
	homeDir := os.Getenv("HOME")
	repoDir := "Gentleman.Dots"
	wm := m.Choices.WindowMgr

	switch wm {
	case "tmux":
		if !system.CommandExists("tmux") {
			result := system.RunBrew("install tmux", nil)
			if result.Error != nil {
				return result.Error
			}
		}

		// TPM
		tpmDir := filepath.Join(homeDir, ".tmux/plugins/tpm")
		if _, err := os.Stat(tpmDir); os.IsNotExist(err) {
			system.Run(fmt.Sprintf("git clone https://github.com/tmux-plugins/tpm %s", tpmDir), nil)
		}

		system.EnsureDir(filepath.Join(homeDir, ".tmux"))
		system.CopyDir(filepath.Join(repoDir, "GentlemanTmux/plugins"), filepath.Join(homeDir, ".tmux/"))
		system.CopyFile(filepath.Join(repoDir, "GentlemanTmux/tmux.conf"), filepath.Join(homeDir, ".tmux.conf"))

		// Install plugins
		system.Run(filepath.Join(homeDir, ".tmux/plugins/tpm/bin/install_plugins"), nil)

	case "zellij":
		if !system.CommandExists("zellij") {
			result := system.RunBrew("install zellij", nil)
			if result.Error != nil {
				return result.Error
			}
		}

		zellijDir := filepath.Join(homeDir, ".config/zellij")
		system.EnsureDir(zellijDir)
		system.CopyDir(filepath.Join(repoDir, "GentlemanZellij/zellij/*"), zellijDir+"/")
	}

	return nil
}

func stepInstallNvim(m *Model) error {
	homeDir := os.Getenv("HOME")
	repoDir := "Gentleman.Dots"

	// Obsidian path
	obsidianDir := filepath.Join(homeDir, ".config/obsidian")
	system.EnsureDir(obsidianDir)
	system.EnsureDir(filepath.Join(obsidianDir, "templates"))

	// Check Node.js
	if !system.CommandExists("node") {
		result := system.RunBrew("install node", nil)
		if result.Error != nil {
			return result.Error
		}
	}

	// Install dependencies
	result := system.RunBrew("install nvim git gcc fzf fd ripgrep coreutils bat curl lazygit tree-sitter", nil)
	if result.Error != nil {
		return result.Error
	}

	// Copy config
	nvimDir := filepath.Join(homeDir, ".config/nvim")
	system.EnsureDir(nvimDir)
	system.CopyDir(filepath.Join(repoDir, "GentlemanNvim/nvim/*"), nvimDir+"/")

	// Install Claude Code
	system.Run(`curl -fsSL https://claude.ai/install.sh | bash`, nil)

	// Install OpenCode
	system.Run(`curl -fsSL https://opencode.ai/install | bash`, nil)

	// Configure OpenCode
	openCodeDir := filepath.Join(homeDir, ".config/opencode")
	system.EnsureDir(openCodeDir)
	system.EnsureDir(filepath.Join(openCodeDir, "themes"))
	system.CopyFile(filepath.Join(repoDir, "GentlemanOpenCode/opencode.json"), filepath.Join(openCodeDir, "opencode.json"))
	system.CopyFile(filepath.Join(repoDir, "GentlemanOpenCode/themes/gentleman.json"), filepath.Join(openCodeDir, "themes/gentleman.json"))

	return nil
}

func stepCleanup(m *Model) error {
	// Fix permissions
	brewPrefix := system.GetBrewPrefix()
	system.RunSudo(fmt.Sprintf("chown -R $(whoami) %s/*", brewPrefix), nil)

	// Remove cloned repo
	system.Run("rm -rf Gentleman.Dots", nil)
	return nil
}

func stepSetDefaultShell(m *Model) error {
	shell := m.Choices.Shell
	if shell == "nushell" {
		shell = "nu"
	}

	// Get shell path
	result := system.Run(fmt.Sprintf("which %s", shell), nil)
	if result.Error != nil {
		return fmt.Errorf("could not find shell: %s", shell)
	}

	shellPath := result.Output
	shellPath = shellPath[:len(shellPath)-1] // Remove newline

	// Add to /etc/shells if not present
	system.RunSudo(fmt.Sprintf(`grep -Fxq "%s" /etc/shells || echo "%s" >> /etc/shells`, shellPath, shellPath), nil)

	// Change default shell
	result = system.RunSudo(fmt.Sprintf("chsh -s %s $USER", shellPath), nil)
	if result.Error != nil {
		return fmt.Errorf("failed to set default shell: %w", result.Error)
	}

	return nil
}
