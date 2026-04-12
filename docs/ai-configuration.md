# AI Configuration for Neovim

This configuration includes several AI assistants integrated with Neovim. By default, **Claude Code is enabled** as the primary AI assistant.

> **Note:** AI tool configurations (persona, skills, themes, MCP servers) are now managed by [gentle-ai](https://github.com/Gentleman-Programming/gentle-ai). This document covers the **Neovim AI plugins** and CLI tool installation only.

## Table of Contents

- [Available AI Assistants](#available-ai-assistants)
- [Switching AI Plugins](#switching-ai-plugins)
- [Required CLI Tools](#required-cli-tools)
- [Recommended by Use Case](#recommended-by-use-case)

---

## Available AI Assistants

| Plugin | Description | Status |
|--------|-------------|--------|
| **Claude Code.nvim** | Claude AI integration (official) | ✅ Enabled by default |
| **OpenCode.nvim** | OpenCode AI integration | Disabled |
| **Avante.nvim** | AI-powered coding assistant | Disabled |
| **CopilotChat.nvim** | GitHub Copilot chat interface | Disabled |
| **CodeCompanion.nvim** | Multi-AI provider support | Disabled |
| **Gemini.nvim** | Google Gemini integration | Disabled |

## Switching AI Plugins

All plugin states are managed in a single file:

```bash
nvim ~/.config/nvim/lua/plugins/disabled.lua
```

**Steps to switch plugins:**

1. Find the plugin you want to disable and set `enabled = false`
2. Find the plugin you want to enable and set `enabled = true`
3. Save and restart Neovim

**Example - switching from Claude Code to OpenCode:**

```lua
{
  "coder/claudecode.nvim",
  enabled = false,  -- Disable Claude Code
},
{
  "NickvanDyke/opencode.nvim",
  enabled = true,   -- Enable OpenCode
},
```

> **Important:** Only enable ONE AI plugin at a time to avoid conflicts and keybinding issues.

## Required CLI Tools

These CLI tools are automatically installed by the Gentleman.Dots installer:

| Tool | Installation Command |
|------|---------------------|
| Claude Code CLI | `curl -fsSL https://claude.ai/install.sh \| bash` |
| OpenCode CLI | `curl -fsSL https://opencode.ai/install \| bash` |
| Gemini CLI | `brew install gemini-cli` |

> Some services require API keys. Check each plugin's documentation for details.

For AI tool **configuration** (persona, skills, themes, MCP servers), use [gentle-ai](https://github.com/Gentleman-Programming/gentle-ai):

```bash
brew install Gentleman-Programming/tap/gentle-ai
gentle-ai
```

## Recommended by Use Case

| Use Case | Recommended Plugin |
|----------|-------------------|
| Full Gentleman experience | **Claude Code.nvim** (default) |
| OpenCode CLI in terminal | **OpenCode.nvim** |
| GitHub Copilot users | **CopilotChat.nvim** |
| Multi-provider flexibility | **CodeCompanion.nvim** |
| Google Gemini users | **Gemini.nvim** |
