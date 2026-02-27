# Agent Teams Lite

[Agent Teams Lite](https://github.com/Gentleman-Programming/agent-teams-lite) is a zero-dependency, Markdown-based SDD (Spec-Driven Development) orchestration framework for AI coding assistants. It replaces the "one giant conversation" paradigm with a delegate-only orchestrator that launches 9 specialized sub-agents.

## Table of Contents

- [What is Agent Teams Lite?](#what-is-agent-teams-lite)
- [How It Differs from OpenSpec](#how-it-differs-from-openspec)
- [The 9 SDD Sub-Agents](#the-9-sdd-sub-agents)
- [Supported AI Tools](#supported-ai-tools)
- [How the Installer Handles It](#how-the-installer-handles-it)
- [SDD Workflow](#sdd-workflow)
- [Persistence Backends](#persistence-backends)
- [Manual Installation](#manual-installation)

---

## What is Agent Teams Lite?

The core philosophy: long conversations with AI lead to lost details and hallucinations due to context compression. Instead, a lightweight orchestrator coordinates focused sub-agents that each do ONE job well, return structured results, and hand off to the next phase.

**Key characteristics:**
- **Zero dependencies** — no npm, no binaries, just Markdown files
- **9 specialized sub-agents** — each with a SKILL.md file
- **Delegate-only orchestrator** — never executes work itself, only coordinates
- **Fresh context per agent** — each sub-agent starts with clean context to avoid hallucination
- **Multi-tool support** — works with Claude Code, OpenCode, Gemini CLI, and more

## How It Differs from OpenSpec

Both are SDD implementations available in the installer, but they serve different roles:

| Aspect | OpenSpec | Agent Teams Lite |
|--------|----------|-----------------|
| **Source** | [project-starter-framework](https://github.com/JNZader/project-starter-framework) | [agent-teams-lite](https://github.com/Gentleman-Programming/agent-teams-lite) |
| **Install method** | `setup-global.sh --features=sdd` | `install.sh --agent <tool>` |
| **Format** | Part of the larger framework with 203 modules | Standalone, focused on SDD only |
| **Dependencies** | Installed alongside other framework features | Zero dependencies, just Markdown |
| **Persistence** | File-based YAML schema in `openspec/` | Engram (recommended), OpenSpec, or none |
| **Weight** | Heavier (part of full framework) | Lightweight (9 Markdown files) |

**You can install both.** They complement each other — OpenSpec provides the file-based artifact structure, while Agent Teams Lite provides the orchestration pattern.

## The 9 SDD Sub-Agents

Each sub-agent is a SKILL.md file installed in your tool's skills directory:

| Phase | Sub-Agent | Purpose |
|-------|-----------|---------|
| 1 | `sdd-init` | Bootstrap project context, detect stack and existing conventions |
| 2 | `sdd-explore` | Investigate codebase, compare approaches, clarify requirements |
| 3 | `sdd-propose` | Create change proposal with intent, scope, and approach |
| 4 | `sdd-spec` | Write delta specifications with GIVEN/WHEN/THEN scenarios |
| 5 | `sdd-design` | Technical architecture decisions and approach documentation |
| 6 | `sdd-tasks` | Break work into phased implementation task checklists |
| 7 | `sdd-apply` | Implement code following specs and design (supports TDD) |
| 8 | `sdd-verify` | Validate completeness with actual test execution |
| 9 | `sdd-archive` | Merge delta specs into main specs, archive completed change |

### Dependency Graph

```
Explorer → Proposer → (Specs ⟷ Design) → Tasks → Apply → Verify → Archive
```

## Supported AI Tools

Agent Teams Lite supports multiple AI coding assistants:

| AI Tool | Agent Name | Install Path |
|---------|-----------|-------------|
| Claude Code | `claude-code` | `~/.claude/skills/sdd-*/SKILL.md` |
| OpenCode | `opencode` | `~/.config/opencode/skill/sdd-*/SKILL.md` |
| Gemini CLI | `gemini-cli` | `~/.gemini/skills/sdd-*/SKILL.md` |
| Codex | `codex` | `~/.codex/skills/sdd-*/SKILL.md` |
| VS Code | `vscode` | `./.vscode/skills/sdd-*/SKILL.md` |
| Antigravity | `antigravity` | `~/.gemini/antigravity/skills/sdd-*/SKILL.md` |
| Cursor | `cursor` | `~/.cursor/skills/sdd-*/SKILL.md` |

## How the Installer Handles It

### TUI Mode (Interactive)

1. In the **Custom** category drill-down, navigate to the **SDD** category
2. You'll see two checkboxes:
   - `[ ] OpenSpec (project-starter-framework)`
   - `[ ] Agent Teams Lite`
3. Toggle one or both with `Enter`
4. Confirm your selection

### CLI Mode (Non-Interactive)

```bash
# Agent Teams Lite only
gentleman.dots --non-interactive --shell=zsh \
  --ai-tools=claude --ai-framework --agent-teams-lite

# OpenSpec only
gentleman.dots --non-interactive --shell=zsh \
  --ai-tools=claude --ai-modules=sdd

# Both
gentleman.dots --non-interactive --shell=zsh \
  --ai-tools=claude --ai-modules=hooks,skills,sdd --agent-teams-lite
```

### What the Installer Does

When Agent Teams Lite is selected, the installer:

1. **Clones** the repo: `git clone --depth 1 https://github.com/Gentleman-Programming/agent-teams-lite.git /tmp/agent-teams-lite-install`
2. **Runs** `install.sh` for each supported AI tool:
   - `claude` → `install.sh --agent claude-code`
   - `opencode` → `install.sh --agent opencode`
   - `gemini` → `install.sh --agent gemini-cli`
   - GitHub Copilot is not supported
3. **Cleans up** the cloned repo

### Tool ID Mapping

The installer maps its internal tool IDs to Agent Teams Lite agent names:

| Installer ID | Agent Teams Lite Name |
|-------------|----------------------|
| `claude` | `claude-code` |
| `opencode` | `opencode` |
| `gemini` | `gemini-cli` |
| `copilot` | _(not supported)_ |

## SDD Workflow

Once installed, you can use SDD commands in your AI tool:

| Command | Action |
|---------|--------|
| `/sdd-init` | Bootstrap project context |
| `/sdd-explore <topic>` | Research and investigate |
| `/sdd-new <name>` | Start a new change (explore + propose) |
| `/sdd-continue [name]` | Next artifact in chain |
| `/sdd-ff [name]` | Fast-forward all planning phases |
| `/sdd-apply [name]` | Implement tasks |
| `/sdd-verify [name]` | Validate implementation |
| `/sdd-archive [name]` | Close and archive change |

## Persistence Backends

Agent Teams Lite supports three persistence modes:

| Mode | Description | When to Use |
|------|-------------|-------------|
| **Engram** (recommended) | External persistent memory via MCP | Default — no files in project |
| **OpenSpec** | File-based artifacts in `openspec/changes/` | When you need file-based specs |
| **None** | Ephemeral, inline results only | Quick explorations |

> **Important:** OpenSpec mode is NEVER auto-selected. It requires explicit user request. Engram is the recommended default.

## Manual Installation

If you prefer to install Agent Teams Lite manually (outside the TUI installer):

```bash
# Clone the repository
git clone https://github.com/Gentleman-Programming/agent-teams-lite.git
cd agent-teams-lite

# Interactive menu
./scripts/install.sh

# Non-interactive for a specific tool
./scripts/install.sh --agent claude-code

# Install for all tools globally
./scripts/install.sh --agent all-global
```

After installation, add the orchestrator block to your tool's configuration file:
- **Claude Code:** Copy from `examples/claude-code/CLAUDE.md` into `~/.claude/CLAUDE.md`
- **Gemini CLI:** Copy from `examples/gemini-cli/GEMINI.md` into `~/.gemini/GEMINI.md`
- **OpenCode:** Merge from `examples/opencode/opencode.json` into `~/.config/opencode/opencode.json`
