package tui

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Gentleman-Programming/Gentleman.Dots/installer/internal/system"
)

// AIAssistant represents an AI coding assistant that can be installed
type AIAssistant struct {
	ID           string   // Unique identifier (e.g., "opencode")
	Name         string   // Display name (e.g., "OpenCode")
	Description  string   // Short description
	LongDesc     string   // Detailed description for the selection screen
	Available    bool     // Whether this assistant is currently available
	SkillsPath   string   // Path to skills in repo (e.g., "GentlemanOpenCode/skill")
	ConfigPath   string   // Installation path relative to $HOME (e.g., ".config/opencode")
	InstallCmd   string   // Command to install the assistant
	Skills       []string // List of skill names
	ConfigFiles  []string // Additional config files to copy (relative to SkillsPath parent)
	RequiresNvim bool     // Whether this assistant requires Neovim
}

// GetAvailableAIAssistants returns the list of all AI assistants
func GetAvailableAIAssistants() []AIAssistant {
	return []AIAssistant{
		{
			ID:          "claudecode",
			Name:        "Claude Code",
			Description: "Official Claude AI integration (installed with Neovim)",
			LongDesc: `Claude Code is the official Anthropic CLI with:
  • Native Claude integration for terminal
  • Skills system for custom patterns
  • MCP server support
  • Custom themes and statusline

Note: Automatically installed when you select Neovim.`,
			Available:    true,
			SkillsPath:   "GentlemanClaude/skills",
			ConfigPath:   ".claude",
			InstallCmd:   "curl -fsSL https://claude.ai/install.sh | bash",
			RequiresNvim: true,
			Skills: []string{
				"react-19", "nextjs-15", "typescript", "tailwind-4",
				"ai-sdk-5", "django-drf", "playwright", "pytest",
				"zod-4", "zustand-5",
			},
			ConfigFiles: []string{
				"CLAUDE.md",
				"settings.json",
				"statusline.sh",
				"output-styles/gentleman.md",
				"mcp-servers.template.json",
				"tweakcc-theme.json",
			},
		},
		{
			ID:          "opencode",
			Name:        "OpenCode",
			Description: "Terminal-native AI with Claude Max/Pro support",
			LongDesc: `OpenCode is a terminal-native AI coding assistant with support for:
  • Claude Max/Pro subscription integration
  • Skills system for React, Next.js, TypeScript, and more
  • Terminal-first workflow
  • Custom themes and configurations`,
			Available:    true,
			SkillsPath:   "GentlemanOpenCode/skill",
			ConfigPath:   ".config/opencode",
			InstallCmd:   "curl -fsSL https://opencode.ai/install | bash",
			RequiresNvim: false,
			Skills: []string{
				"react-19", "nextjs-15", "typescript", "tailwind-4",
				"ai-sdk-5", "django-drf", "playwright", "pytest",
				"zod-4", "zustand-5", "skill-creator", "jira-task", "jira-epic",
			},
			ConfigFiles: []string{
				"opencode.json",
				"themes/gentleman.json",
			},
		},
		{
			ID:          "gemini-cli",
			Name:        "Gemini CLI",
			Description: "Google Gemini AI CLI with terminal and Neovim support",
			LongDesc: `Gemini CLI provides Google's Gemini AI in your terminal:
  • Official Google CLI for Gemini models
  • Terminal-native chat interface
  • Integration with Neovim via gemini.lua plugin
  • Multiple Gemini models support (Pro, Flash, etc.)
  • Free tier available

Requires: Node.js/npm (for CLI installation)
Note: Automatically installed when you select Neovim.`,
			Available:    true,
			SkillsPath:   "",
			ConfigPath:   "",
			InstallCmd:   "npm install -g @google/gemini-cli",
			RequiresNvim: true,
			Skills:       []string{},
			ConfigFiles:  []string{},
		},
		{
			ID:          "copilot-cli",
			Name:        "GitHub Copilot CLI",
			Description: "GitHub Copilot CLI with terminal and Neovim support",
			LongDesc: `GitHub Copilot CLI brings AI pair programming to your terminal:
  • Official GitHub CLI for Copilot
  • Terminal-native AI assistance
  • Integration with Neovim via copilot.lua and copilot-chat.lua
  • Code suggestions and completions
  • Chat interface for questions
  • Multi-language support

Requires: GitHub Copilot subscription + Node.js/npm
Note: Automatically installed when you select Neovim.`,
			Available:    true,
			SkillsPath:   "",
			ConfigPath:   "",
			InstallCmd:   "npm install -g @github/copilot",
			RequiresNvim: true,
			Skills:       []string{},
			ConfigFiles:  []string{},
		},
	}
}

// GetAIAssistantByID returns an AI assistant by its ID
func GetAIAssistantByID(id string) *AIAssistant {
	for _, ai := range GetAvailableAIAssistants() {
		if ai.ID == id {
			return &ai
		}
	}
	return nil
}

// GetAvailableAIAssistantsOnly returns only the assistants that are currently available
func GetAvailableAIAssistantsOnly() []AIAssistant {
	assistants := GetAvailableAIAssistants()
	available := make([]AIAssistant, 0)
	for _, ai := range assistants {
		if ai.Available {
			available = append(available, ai)
		}
	}
	return available
}

// InstallAIAssistant installs a single AI assistant
func InstallAIAssistant(assistant AIAssistant, m *Model, repoDir string, stepID string) error {
	homeDir := os.Getenv("HOME")

	// Install the AI assistant binary (if not on Termux)
	if !m.SystemInfo.IsTermux && assistant.InstallCmd != "" {
		SendLog(stepID, fmt.Sprintf("Installing %s...", assistant.Name))
		result := system.RunWithLogs(assistant.InstallCmd, nil, func(line string) {
			SendLog(stepID, line)
		})
		if result.Error != nil {
			// Non-critical error - continue with configuration
			SendLog(stepID, fmt.Sprintf("⚠️  Could not install %s binary (you may need to install manually)", assistant.Name))
		} else {
			SendLog(stepID, fmt.Sprintf("✓ %s installed", assistant.Name))
		}
	} else if m.SystemInfo.IsTermux {
		SendLog(stepID, fmt.Sprintf("Skipping %s installation (not supported on Termux)", assistant.Name))
	}

	// Configure - create directories
	SendLog(stepID, fmt.Sprintf("Configuring %s...", assistant.Name))
	configDir := filepath.Join(homeDir, assistant.ConfigPath)
	if err := system.EnsureDir(configDir); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Create skill directory
	skillDir := filepath.Join(configDir, "skill")
	if err := system.EnsureDir(skillDir); err != nil {
		return fmt.Errorf("failed to create skill directory: %w", err)
	}

	// Create themes directory if needed
	themesDir := filepath.Join(configDir, "themes")
	if err := system.EnsureDir(themesDir); err != nil {
		return fmt.Errorf("failed to create themes directory: %w", err)
	}

	// Copy skills
	if assistant.SkillsPath != "" {
		srcSkills := filepath.Join(repoDir, assistant.SkillsPath)
		if _, err := os.Stat(srcSkills); err == nil {
			if err := system.CopyDir(srcSkills, skillDir); err != nil {
				return fmt.Errorf("failed to copy skills: %w", err)
			}
			SendLog(stepID, fmt.Sprintf("✓ Copied %d skills", len(assistant.Skills)))
		}
	}

	// Copy config files
	for _, configFile := range assistant.ConfigFiles {
		// Config files are relative to the SkillsPath parent directory
		srcDir := filepath.Dir(filepath.Join(repoDir, assistant.SkillsPath))
		srcFile := filepath.Join(srcDir, configFile)
		dstFile := filepath.Join(configDir, configFile)

		if _, err := os.Stat(srcFile); err == nil {
			// Ensure destination directory exists
			if err := system.EnsureDir(filepath.Dir(dstFile)); err != nil {
				return fmt.Errorf("failed to create directory for %s: %w", configFile, err)
			}
			if err := system.CopyFile(srcFile, dstFile); err != nil {
				return fmt.Errorf("failed to copy %s: %w", configFile, err)
			}
			SendLog(stepID, fmt.Sprintf("✓ Copied %s", configFile))
		}
	}

	SendLog(stepID, fmt.Sprintf("✓ %s configured successfully", assistant.Name))
	return nil
}

// stepInstallAIAssistants installs all selected AI assistants
func stepInstallAIAssistants(m *Model) error {
	stepID := "ai"
	repoDir := "Gentleman.Dots"

	if len(m.Choices.AIAssistants) == 0 {
		SendLog(stepID, "No AI assistants selected, skipping...")
		return nil
	}

	SendLog(stepID, fmt.Sprintf("Installing %d AI assistant(s)...", len(m.Choices.AIAssistants)))

	for _, aiID := range m.Choices.AIAssistants {
		assistant := GetAIAssistantByID(aiID)
		if assistant == nil {
			SendLog(stepID, fmt.Sprintf("⚠️  Unknown AI assistant: %s", aiID))
			continue
		}

		if !assistant.Available {
			SendLog(stepID, fmt.Sprintf("⚠️  %s is not available yet", assistant.Name))
			continue
		}

		// Check if Neovim is required but not installed
		if assistant.RequiresNvim && !m.Choices.InstallNvim {
			SendLog(stepID, fmt.Sprintf("⚠️  %s requires Neovim but it's not selected for installation", assistant.Name))
			SendLog(stepID, fmt.Sprintf("   Skipping %s...", assistant.Name))
			continue
		}

		// Install the assistant
		if err := InstallAIAssistant(*assistant, m, repoDir, stepID); err != nil {
			// Non-critical error - log and continue with other assistants
			SendLog(stepID, fmt.Sprintf("⚠️  Error installing %s: %v", assistant.Name, err))
			continue
		}
	}

	SendLog(stepID, "✓ AI assistants configuration complete")
	return nil
}
