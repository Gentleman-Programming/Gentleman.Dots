# AI Tools & Framework Integration

The Gentleman.Dots installer includes a complete AI tools and framework integration system. This guide covers the TUI wizard flow, preset/custom selection, category drill-down, viewport scrolling, and non-interactive CLI usage.

## Table of Contents

- [Overview](#overview)
- [Installation Flow](#installation-flow)
  - [Step 7: AI Tools Selection](#step-7-ai-tools-selection)
  - [Step 8a: Framework Confirmation](#step-8a-framework-confirmation)
  - [Step 8b: Preset Selection](#step-8b-preset-selection)
  - [Step 8c: Custom Category Drill-Down](#step-8c-custom-category-drill-down)
- [Viewport Scrolling](#viewport-scrolling)
- [SDD Choice: OpenSpec vs Agent Teams Lite](#sdd-choice-openspec-vs-agent-teams-lite)
- [Non-Interactive CLI](#non-interactive-cli)
- [Presets Reference](#presets-reference)
- [UX Flow Diagram](#ux-flow-diagram)
- [Progress Bar](#progress-bar)
- [Back Navigation](#back-navigation)
- [Installation Steps](#installation-steps)

---

## Overview

The AI integration adds two new installation steps to the wizard:

| Step | Screen | Description |
|------|--------|-------------|
| **Step 7** | AI Tools Selection | Multi-select: Claude Code, OpenCode, Gemini CLI, GitHub Copilot |
| **Step 8** | AI Framework | Preset or custom module selection from the [project-starter-framework](https://github.com/JNZader/project-starter-framework) |

These steps appear after Neovim configuration (Step 6) and before the backup confirmation screen.

---

## Installation Flow

### Step 7: AI Tools Selection

**Screen:** `ScreenAIToolsSelect`

Multi-select checkboxes for 4 AI coding tools:

| Tool | ID | Install Method |
|------|-----|----------------|
| Claude Code | `claude` | `curl -fsSL https://claude.ai/install.sh \| bash` + configs, skills, theme |
| OpenCode | `opencode` | `curl -fsSL https://opencode.ai/install \| bash` + configs, theme |
| Gemini CLI | `gemini` | `npm install -g @google/gemini-cli` |
| GitHub Copilot | `copilot` | `gh extension install github/gh-copilot` |

**Behavior:**
- Toggle individual tools with `Enter` or `Space`
- `[✓]` / `[ ]` checkboxes show selection state
- "Confirm selection" collects all toggled tools
- If **no tools** are selected, the framework step is skipped entirely
- This step is **skipped on Termux** (AI tools not supported on Android)

**What gets installed for Claude Code:**
- Binary via official installer
- `CLAUDE.md` global instructions with Gentleman personality
- `settings.json` with permissions and output style
- `statusline.sh` custom status bar
- `output-styles/gentleman.md` persona definition
- 10+ skill directories (React 19, Next.js 15, TypeScript, etc.)
- `mcp-servers.template.json` for Context7, Jira, Figma
- `tweakcc-theme.json` Kanagawa-inspired color theme

**What gets installed for OpenCode:**
- Binary via official installer
- `opencode.json` with Gentleman agent and SDD orchestrator
- Gentleman theme
- Skills directory

### Step 8a: Framework Confirmation

**Screen:** `ScreenAIFrameworkConfirm`

Simple Yes/No to install the [project-starter-framework](https://github.com/JNZader/project-starter-framework):

- **Yes** → proceeds to preset selection
- **No** → skips framework → proceeds to backup/install

### Step 8b: Preset Selection

**Screen:** `ScreenAIFrameworkPreset`

Custom is the **first option** (index 0), followed by a separator, then 6 presets:

```
🔧 Custom — Pick individual modules        ← Custom (index 0)
─────────────                               ← Separator
🎯 Minimal — Core + git commands only       ← Preset (index 2)
🖥️  Frontend — React, Vue, Angular...       ← Preset (index 3)
⚙️  Backend — APIs, databases...             ← Preset (index 4)
🔄 Fullstack — Frontend + Backend...        ← Preset (index 5)
📊 Data — Data engineering, ML/AI...        ← Preset (index 6)
📦 Complete — Everything included           ← Preset (index 7)
```

- Selecting a **preset** sets `AIFrameworkPreset` and proceeds to backup/install
- Selecting **Custom** initializes the category selection map and enters the drill-down

### Step 8c: Custom Category Drill-Down

The custom selection uses a **two-level drill-down** instead of a flat checkbox list, making it possible to navigate 203 individual modules across 6 categories.

#### Category Menu (`ScreenAIFrameworkCategories`)

Shows 6 categories with live selection counts. No checkboxes — cursor navigation only. Press `Enter` to drill into a category.

```
▸ 🪝 Hooks (2/10 selected)
  ⚡ Commands (0/20 selected)
  🤖 Agents (3/80 selected)
  🎯 Skills (5/85 selected)
  📐 SDD (1/2 selected)
  🔌 MCP (2/6 selected)
  ─────────────
  ✅ Confirm selection
```

#### Category Items (`ScreenAIFrameworkCategoryItems`)

Shows individual items within a category with checkboxes. The title dynamically shows the category name and icon.

```
Step 8: 🪝 Hooks

▸ [✓] Block Dangerous Commands
  [ ] Commit Guard
  [✓] Context Loader
  [ ] Improve Prompt
  ...
  ─────────────
  ← Back
```

**Navigation:**
- `Enter` / `Space` toggles items on/off
- `Esc` or "← Back" returns to category menu with **cursor preserved** on the originating category
- "Confirm selection" on the category menu collects all selections

---

## Viewport Scrolling

Categories with many items (Agents: 80, Skills: 85) use **viewport scrolling** to fit within the terminal height.

**How it works:**
- The visible area is calculated as `terminal_height - 8` lines (reserving space for chrome)
- Minimum viewport: 5 items visible at all times
- Scroll indicators show items above/below the viewport:

```
  ▲ 12 more above
▸ [✓] Development: Frontend Specialist
  [ ] Development: Fullstack Engineer
  [ ] Development: Go Pro
  ...
  ▼ 45 more below
```

**Scroll sync:** The viewport automatically follows the cursor:
- Moving cursor above viewport → scrolls up
- Moving cursor below viewport → scrolls down
- Scroll resets to 0 when entering a category or going back

---

## SDD Choice: OpenSpec vs Agent Teams Lite

The SDD category offers a choice between two SDD implementations:

| Option | Description | Install Method |
|--------|-------------|----------------|
| **OpenSpec** | SDD from project-starter-framework | Included via `--features=sdd` in `setup-global.sh` |
| **Agent Teams Lite** | Lightweight SDD with 9 sub-agent skills | Separate clone + `install.sh --agent <tool>` |

You can select **one or both**.

**When OpenSpec is selected:**
- The `sdd` feature flag is included in the `setup-global.sh --features=...` command
- This installs the SDD phases from the project-starter-framework

**When Agent Teams Lite is selected:**
- The installer clones [agent-teams-lite](https://github.com/Gentleman-Programming/agent-teams-lite)
- Runs `install.sh --agent <tool>` for each supported AI tool
- Tool mapping: `claude` → `claude-code`, `opencode` → `opencode`, `gemini` → `gemini-cli`
- GitHub Copilot is not supported by Agent Teams Lite

**When both are selected:**
- Both installations run sequentially

See [Agent Teams Lite documentation](agent-teams-lite.md) for more details.

---

## Non-Interactive CLI

For CI/CD or scripted installations, all AI options are available as CLI flags:

```bash
gentleman.dots --non-interactive --shell=<shell> [AI options]
```

### AI Flags

| Flag | Values | Description |
|------|--------|-------------|
| `--ai-tools=<tools>` | `claude,opencode,gemini,copilot` | AI tools (comma-separated) |
| `--ai-framework` | | Install AI coding framework |
| `--ai-preset=<name>` | `minimal,frontend,backend,fullstack,data,complete` | Framework preset |
| `--ai-modules=<feats>` | `hooks,commands,skills,agents,sdd,mcp` | Feature flags (comma-separated) |
| `--agent-teams-lite` | | Install Agent Teams Lite SDD framework |

### Examples

```bash
# Full setup with preset
gentleman.dots --non-interactive --shell=fish --nvim \
  --ai-tools=claude,opencode,gemini,copilot --ai-preset=fullstack

# Custom feature selection with Agent Teams Lite
gentleman.dots --non-interactive --shell=zsh --ai-tools=claude --ai-framework \
  --ai-modules=hooks,skills --agent-teams-lite

# Both SDD options: OpenSpec + Agent Teams Lite
gentleman.dots --non-interactive --shell=zsh --ai-tools=claude \
  --ai-modules=hooks,skills,sdd --agent-teams-lite

# Dry run to preview
gentleman.dots --dry-run --non-interactive --shell=zsh \
  --ai-tools=claude --ai-preset=complete
```

### Validation

- AI tools validated against: `{claude, opencode, gemini, copilot}`
- Presets validated against: `{minimal, frontend, backend, fullstack, data, complete}`
- Feature flags validated against: `{hooks, commands, skills, agents, sdd, mcp}`
- Framework is auto-enabled if `--ai-preset`, `--ai-modules`, or `--agent-teams-lite` is provided

---

## Presets Reference

Each preset maps to a set of feature flags passed to `setup-global.sh --features=`:

| Preset | hooks | commands | skills | agents | sdd | mcp |
|--------|:-----:|:--------:|:------:|:------:|:---:|:---:|
| **Minimal** | ✅ | ✅ | | | ✅ | |
| **Frontend** | ✅ | ✅ | ✅ | ✅ | ✅ | |
| **Backend** | ✅ | ✅ | ✅ | ✅ | ✅ | |
| **Fullstack** | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| **Data** | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| **Complete** | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |

> **Note:** Presets install OpenSpec SDD by default (via the `sdd` feature). For Agent Teams Lite, use Custom mode or the `--agent-teams-lite` CLI flag.

---

## UX Flow Diagram

```
ScreenNvimSelect (Step 6)
  └→ ScreenAIToolsSelect (Step 7)
       ├→ [Confirm with tools] → ScreenAIFrameworkConfirm (Step 8a)
       │     ├→ Yes → ScreenAIFrameworkPreset (Step 8b)
       │     │     ├→ Custom (idx 0) → ScreenAIFrameworkCategories (Step 8c)
       │     │     │     ├→ [Enter category] → ScreenAIFrameworkCategoryItems
       │     │     │     │     ├→ [Toggle items] → stays in items
       │     │     │     │     └→ [Esc/Back] → back to Categories (cursor preserved)
       │     │     │     └→ [Confirm] → proceedToBackupOrInstall
       │     │     └→ Preset (idx 2-7) → proceedToBackupOrInstall
       │     └→ No → proceedToBackupOrInstall
       └→ [Confirm with no tools] → proceedToBackupOrInstall (skip framework)
```

---

## Progress Bar

The step progress indicator shows 8 steps:

```
✓ OS → ✓ Terminal → ✓ Font → ✓ Shell → ✓ WM → ✓ Nvim → ● AI Tools → ○ Framework
```

All AI screens (tools, framework confirm, preset, categories, category items) show the progress bar.

---

## Back Navigation

Complete bidirectional navigation through the entire wizard:

```
Nvim ← AIToolsSelect ← AIFrameworkConfirm ← AIFrameworkPreset ← AIFrameworkCategories ← AIFrameworkCategoryItems
```

**Special cases:**

| From | Back goes to | Notes |
|------|-------------|-------|
| `AIToolsSelect` | `NvimSelect` | Clears `AIToolSelected` and `AITools` |
| `AIFrameworkConfirm` | `AIToolsSelect` | |
| `AIFrameworkPreset` | `AIFrameworkConfirm` | |
| `AIFrameworkCategories` | `AIFrameworkPreset` | Clears `AICategorySelected` map |
| `AIFrameworkCategoryItems` | `AIFrameworkCategories` | **Cursor preserved** on originating category |
| `BackupConfirm` (Esc) | Smart routing based on wizard state | See below |

**Smart back-routing from BackupConfirm:**
- Has AI tools + framework + custom mode → goes to Categories
- Has AI tools + framework → goes to Preset
- Has AI tools only → goes to Framework Confirm
- No AI tools → goes to AI Tools Select

---

## Installation Steps

### `stepInstallAITools` (step ID: `aitools`)

Installs each selected tool independently:

| Tool | Process |
|------|---------|
| **Claude Code** | Install binary, copy CLAUDE.md, settings.json, statusline.sh, output-styles, 10+ skills, mcp-servers template, tweakcc theme, apply theme |
| **OpenCode** | Install binary, copy opencode.json, gentleman theme, skills |
| **Gemini CLI** | `npm install -g @google/gemini-cli` |
| **GitHub Copilot** | `gh extension install github/gh-copilot` |

Individual tool install failures are logged as **warnings**, not errors. The installation continues.

### `stepInstallAIFramework` (step ID: `aiframework`)

Two potential sub-steps:

**1. Project Starter Framework** (if features selected):
- Clean up leftover clone directory
- Shallow clone `project-starter-framework` to `/tmp/project-starter-framework-install`
- Run `setup-global.sh --auto --skip-install --clis=<tools> --features=<features>`
- Clean up clone

**2. Agent Teams Lite** (if selected):
- Clean up leftover clone directory
- Shallow clone `agent-teams-lite` to `/tmp/agent-teams-lite-install`
- Run `install.sh --agent <tool>` for each supported AI tool
- Clean up clone

If only Agent Teams Lite is selected (no other features), the project-starter-framework clone is skipped entirely.
