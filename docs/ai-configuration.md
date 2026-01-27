# AI Configuration for Neovim

> ℹ️ **Update (January 2026)**: OpenCode now supports Claude Max/Pro subscriptions via the `opencode-anthropic-auth` plugin (included in this config). Both Claude Code and OpenCode work with your Claude subscription. *Note: This workaround is stable for now, but Anthropic could block it in the future.*

This configuration includes several AI assistants integrated with Neovim. By default, **Claude Code is enabled** as the primary AI assistant with the custom Gentleman personality.

## Table of Contents

- [Available AI Assistants](#available-ai-assistants)
- [Switching AI Plugins](#switching-ai-plugins)
- [Required CLI Tools](#required-cli-tools)
- [Recommended by Use Case](#recommended-by-use-case)
- [Claude Code Configuration](#claude-code-configuration)
  - [What's Included](#whats-included)
  - [Gentleman Persona](#gentleman-persona)
  - [Using Claude Code](#using-claude-code)
  - [Configuration Location](#configuration-location)
  - [Gentleman Theme](#gentleman-theme-visual-colors)
- [OpenCode Configuration](#opencode-configuration)
  - [Using the Gentleman Agent](#using-the-gentleman-agent)
  - [Configuring the Default Model](#configuring-the-default-model)
  - [Available Models](#available-models)
  - [OpenCode Theme](#opencode-theme)
  - [MCP Integrations](#mcp-integrations)

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

> ⚠️ **Important:** Only enable ONE AI plugin at a time to avoid conflicts and keybinding issues.

## Required CLI Tools

These are automatically installed by the installer script:

| Tool | Installation Command |
|------|---------------------|
| OpenCode CLI | `curl -fsSL https://opencode.ai/install \| bash` |
| Claude Code CLI | `curl -fsSL https://claude.ai/install.sh \| bash` |
| Gemini CLI | `brew install gemini-cli` |

> Some services require API keys. Check each plugin's documentation for details.

## Recommended by Use Case

| Use Case | Recommended Plugin |
|----------|-------------------|
| Full Gentleman experience | **Claude Code.nvim** (default) |
| OpenCode CLI in terminal | **OpenCode.nvim** |
| GitHub Copilot users | **CopilotChat.nvim** |
| Multi-provider flexibility | **CodeCompanion.nvim** |
| Google Gemini users | **Gemini.nvim** |

---

## Claude Code Configuration

Claude Code is installed automatically with the custom **Gentleman** output style, skills, and configuration.

### What's Included

The installer configures:

| Component | Description |
|-----------|-------------|
| `CLAUDE.md` | Global instructions with Gentleman personality |
| `settings.json` | Permissions, output style, status line config |
| `statusline.sh` | Custom status bar script |
| `output-styles/gentleman.md` | The Gentleman persona definition |
| `skills/` | 10 framework-specific coding standards |
| `mcp-servers.template.json` | MCP server templates (Context7, Jira, Figma) |

**Included Skills:**

React 19, Next.js 15, TypeScript, Tailwind 4, Zod 4, Zustand 5, AI SDK 5, Django DRF, Playwright, Pytest

### Gentleman Persona

The Gentleman persona is a Senior Architect with 15+ years of experience. Both Claude Code and OpenCode share this personality:

| Trait | Description |
|-------|-------------|
| **Never a Yes-Man** | Won't agree without verifying first |
| **Collaborative Partner** | Like Jarvis to Tony Stark - provides data, alternatives, and pushes back |
| **Proposes Alternatives** | Always presents options with tradeoffs |
| **Verifies Claims** | Investigates before accepting challenges to suggestions |
| **Bilingual** | Rioplatense Spanish or direct English based on your input language |

### Using Claude Code

1. Open Claude Code in your terminal:

   ```bash
   claude
   # or use the alias
   cc
   ```

2. Select the Gentleman output style:

   ```bash
   /config
   # Navigate to "Output style" and select "Gentleman"
   ```

   Or set it directly:

   ```bash
   claude config set outputStyle Gentleman
   ```

### Configuration Location

```
~/.claude/
├── CLAUDE.md              # Global instructions
├── settings.json          # Settings and permissions
├── statusline.sh          # Status bar script
├── output-styles/
│   └── gentleman.md       # Gentleman persona definition
└── skills/                # Framework coding standards
    ├── react-19/
    ├── nextjs-15/
    ├── typescript/
    └── ...
```

### Gentleman Theme (Visual Colors)

Claude Code supports custom color themes via [tweakcc](https://github.com/Piebald-AI/tweakcc). The Gentleman theme provides Kanagawa-inspired colors.

**Installation:**

```bash
# 1. Install tweakcc and apply the Gentleman theme
npx tweakcc

# 2. Go to "Themes" > Create new theme or import
# 3. Import from: GentlemanClaude/tweakcc-theme.json
# 4. Select "Apply customizations"
```

**Or manually merge the theme:**

```bash
# Add Gentleman theme to tweakcc config
jq '.settings.themes += [input]' ~/.tweakcc/config.json GentlemanClaude/tweakcc-theme.json > tmp.json && mv tmp.json ~/.tweakcc/config.json
jq '.settings.selectedTheme = "gentleman"' ~/.tweakcc/config.json > tmp.json && mv tmp.json ~/.tweakcc/config.json

# Apply the patch
npx tweakcc --apply
```

> ⚠️ **After Claude Code updates:** Tweakcc patches Claude Code's CLI directly. Re-apply after updates with `npx tweakcc --apply`

**Theme Colors:**

| Color | RGB | Usage |
|-------|-----|-------|
| Primary blue | `rgb(127,180,202)` | Main UI elements |
| Accent gold | `rgb(224,193,90)` | Permissions, highlights, spinner |
| Success green | `rgb(183,204,133)` | Diffs, confirmations |
| Error pink | `rgb(203,124,148)` | Errors, removals |
| Purple | `rgb(201,154,214)` | Plan mode |
| Dark background | `rgb(6,8,15)` | Terminal background |

**Custom Thinking Verbs (Spanish/Rioplatense):**

The config includes 40+ custom "thinking" messages:

- *Remontando el Barrilete Cósmico…*
- *Preguntándole al Patito…*
- *Haciendo Magia Negra…*
- *Bancame un Toque…*
- *En Modo Jarvis…*

---

## OpenCode Configuration

OpenCode is installed automatically with a custom **Gentleman** agent and theme.

> ✅ **Claude Max/Pro Support**: OpenCode supports Claude subscriptions via the `opencode-anthropic-auth` plugin, already configured in `GentlemanOpenCode/opencode.json`. Just run `opencode` and authenticate with your Claude account.

### Using the Gentleman Agent

1. Open OpenCode in your terminal:

   ```bash
   opencode
   ```

2. Type `/agent` and press Enter

3. Select **gentleman** from the list

### Configuring the Default Model

Edit your OpenCode configuration:

```bash
nvim ~/.config/opencode/opencode.json
```

Example configuration:

```json
{
  "$schema": "https://opencode.ai/config.json",
  "theme": "gentleman",
  "model": "anthropic/claude-sonnet-4-20250514",
  "agent": {
    "gentleman": {
      "model": "anthropic/claude-sonnet-4-20250514"
    }
  }
}
```

### Available Models

| Provider | Model ID |
|----------|----------|
| Anthropic | `anthropic/claude-sonnet-4-20250514` |
| Anthropic | `anthropic/claude-haiku-4-20250514` |
| OpenAI | `openai/gpt-4o` |
| OpenAI | `openai/gpt-4o-mini` |
| Google | `google/gemini-2.0-flash` |
| Google | `google/gemini-2.5-pro-preview-06-05` |

> You can also set a specific model per agent in the `agent` section.

### OpenCode Theme

The configuration includes a custom **Gentleman** theme with a dark background and Kanagawa-inspired colors. The theme is automatically applied when you run OpenCode.

### MCP Integrations

The Gentleman OpenCode config includes MCP (Model Context Protocol) integration:

| Server | Description |
|--------|-------------|
| **Context7** | Remote MCP for fetching up-to-date documentation |

This is enabled by default and enhances the agent's ability to verify information with current docs.
