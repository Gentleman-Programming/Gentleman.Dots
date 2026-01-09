# AI Configuration for Neovim

> ⚠️ **Important Notice (January 2026)**: Anthropic has blocked third-party tools from using Claude Max subscriptions. OAuth tokens are now restricted to Claude Code only. This config now uses **Claude Code as the primary AI assistant**.

This configuration includes several AI assistants integrated with Neovim. By default, **Claude Code is enabled** as the primary AI assistant with the custom Gentleman personality.

## Available AI Assistants

| Plugin | Description | Status |
|--------|-------------|--------|
| **Claude Code.nvim** | Claude AI integration (official) | ✅ Enabled by default |
| **OpenCode.nvim** | OpenCode AI integration | Disabled (requires own API keys) |
| **Avante.nvim** | AI-powered coding assistant | Disabled |
| **CopilotChat.nvim** | GitHub Copilot chat interface | Disabled |
| **CodeCompanion.nvim** | Multi-AI provider support | Disabled |
| **Gemini.nvim** | Google Gemini integration | Disabled |

## Switching AI Plugins

All plugin states are managed in a single file:

```bash
nvim ~/.config/nvim/lua/plugins/disabled.lua
```

To switch plugins:

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
  enabled = true,   -- Enable OpenCode (requires API keys)
},
```

> **Important:** Only enable ONE AI plugin at a time to avoid conflicts and keybinding issues.

## Required CLI Tools

These are automatically installed by the script:

- **OpenCode CLI**: `curl -fsSL https://opencode.ai/install | bash`
- **Claude Code CLI**: `curl -fsSL https://claude.ai/install.sh | bash`
- **Gemini CLI**: `brew install gemini-cli`

API keys may be required for some services - check each plugin's documentation.

## Recommended AI Assistants

| Use Case | Recommended Plugin |
|----------|-------------------|
| Full Gentleman experience | **Claude Code.nvim** (default) |
| OpenAI/other API keys | **OpenCode.nvim** |
| GitHub Copilot users | **CopilotChat.nvim** |
| Multi-provider flexibility | **CodeCompanion.nvim** |
| Google Gemini users | **Gemini.nvim** |

---

## Claude Code Configuration (Primary)

Claude Code is installed automatically with the custom **Gentleman** output style, skills, and configuration.

### What's Included

The installer configures:

- **CLAUDE.md**: Global instructions with Gentleman personality
- **settings.json**: Permissions, output style, status line config
- **statusline.sh**: Custom status bar script
- **output-styles/gentleman.md**: The Gentleman persona definition
- **skills/**: 10 framework-specific coding standards (React 19, Next.js 15, TypeScript, Tailwind 4, Zod 4, Zustand 5, AI SDK 5, Django DRF, Playwright, Pytest)
- **mcp-servers.template.json**: MCP server templates (Context7, Jira, Figma)

### Gentleman Output Style

The Gentleman output style is a Senior Architect persona with 15+ years of experience:

- **Never a Yes-Man**: Won't say "you're right" without verifying first
- **Collaborative Partner**: Like Jarvis to Tony Stark - provides data, alternatives, and pushes back
- **Proposes Alternatives**: Always presents options with tradeoffs
- **Bilingual**: Responds in Rioplatense Spanish or confrontational English based on your language

### Using Claude Code with Gentleman

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

3. The Gentleman personality will now be applied to all responses.

### Configuration Location

All Claude Code config is stored in:
```bash
~/.claude/
├── CLAUDE.md              # Global instructions
├── settings.json          # Settings and permissions
├── statusline.sh          # Status bar script
├── output-styles/         # Output style definitions
│   └── gentleman.md
└── skills/                # Framework coding standards
    ├── react-19/
    ├── nextjs-15/
    ├── typescript/
    └── ...
```

### Gentleman Theme (Visual Colors)

Claude Code supports custom color themes via [tweakcc](https://github.com/Piebald-AI/tweakcc). The Gentleman theme provides Kanagawa-inspired colors matching the OpenCode theme.

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

**After Claude Code updates:**

Tweakcc patches Claude Code's CLI directly. After updates, re-apply:
```bash
npx tweakcc --apply
```

**Theme colors include:**
- Primary blue: `rgb(127,180,202)` - main UI elements
- Accent gold: `rgb(224,193,90)` - permissions, highlights, spinner shimmer
- Success green: `rgb(183,204,133)` - diffs, confirmations
- Error pink: `rgb(203,124,148)` - errors, removals
- Purple: `rgb(201,154,214)` - plan mode
- Dark background: `rgb(6,8,15)` - terminal background

**Custom thinking verbs (Spanish/Rioplatense):**

The config also includes custom "thinking" messages like:
- *Remontando el Barrilete Cósmico…*
- *Preguntándole al Patito…*
- *Haciendo Magia Negra…*
- *Bancame un Toque…*
- *En Modo Jarvis…*
- And 40+ more!

---

## OpenCode Configuration

> ⚠️ **Note**: As of January 2026, OpenCode can no longer use Claude Max subscriptions. You'll need your own API keys (OpenAI, Anthropic API, etc.) or wait for OpenCode's upcoming subscription service.

OpenCode is installed automatically with a custom **Gentleman** agent and theme.

### Gentleman Agent Philosophy

The Gentleman agent is a Senior Architect persona with 15+ years of experience:

- **Never a Yes-Man**: Won't say "you're right" without verifying first
- **Collaborative Partner**: Like Jarvis to Tony Stark - provides data, alternatives, and pushes back
- **Proposes Alternatives**: Always presents options with tradeoffs
- **Verifies Before Agreeing**: Investigates before accepting challenges to suggestions
- **Bilingual**: Responds in Rioplatense Spanish or confrontational English based on your language

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
      ...
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

You can also set a specific model per agent:

```json
{
  "agent": {
    "gentleman": {
      "model": "anthropic/claude-sonnet-4-20250514",
      ...
    }
  }
}
```

### OpenCode Theme

The configuration includes a custom **Gentleman** theme with a dark background and Kanagawa-inspired colors. The theme is automatically applied when you run OpenCode.

### MCP Integrations

The Gentleman OpenCode config includes MCP (Model Context Protocol) integration:

- **Context7**: Remote MCP for fetching up-to-date documentation

This is enabled by default and enhances the agent's ability to verify information with current docs.
