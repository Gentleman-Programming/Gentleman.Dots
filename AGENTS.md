# Gentleman.Dots AI Agent Skills

> **Single Source of Truth** - This file is the master for all AI assistants.
> Run `./skills/setup.sh` to sync to Claude, Gemini, Copilot, and Codex formats.

This repository provides AI agent skills for Claude Code, OpenCode, and other AI assistants.
Skills provide on-demand context and patterns for working with this codebase.

## Quick Start

When working on this project, Claude Code automatically loads relevant skills based on context.
For manual loading, read the SKILL.md file directly.

## Available Skills

### Gentleman.Dots Specific (Repository Skills)

| Skill | Description | File |
|-------|-------------|------|
| `gentleman-bubbletea` | Bubbletea TUI patterns, Model-Update-View, screen navigation | [SKILL.md](skills/gentleman-bubbletea/SKILL.md) |
| `gentleman-trainer` | Vim Trainer RPG system, exercises, progression, boss fights | [SKILL.md](skills/gentleman-trainer/SKILL.md) |
| `gentleman-installer` | Installation steps, interactive/non-interactive modes | [SKILL.md](skills/gentleman-installer/SKILL.md) |
| `gentleman-e2e` | Docker-based E2E testing, multi-platform validation | [SKILL.md](skills/gentleman-e2e/SKILL.md) |
| `gentleman-system` | OS detection, command execution, cross-platform support | [SKILL.md](skills/gentleman-system/SKILL.md) |
| `go-testing` | Go testing patterns, table-driven tests, Bubbletea testing | [SKILL.md](skills/go-testing/SKILL.md) |

### Generic Skills (User Installation → ~/.claude/skills/)

These skills are copied to user's Claude/OpenCode config via the installer.

| Skill | Description | Source |
|-------|-------------|--------|
| `react-19` | React 19 patterns, hooks, components | [GentlemanClaude/skills/react-19](GentlemanClaude/skills/react-19/SKILL.md) |
| `nextjs-15` | Next.js 15, App Router, Server Components | [GentlemanClaude/skills/nextjs-15](GentlemanClaude/skills/nextjs-15/SKILL.md) |
| `typescript` | TypeScript patterns, types, generics | [GentlemanClaude/skills/typescript](GentlemanClaude/skills/typescript/SKILL.md) |
| `tailwind-4` | Tailwind CSS v4 patterns | [GentlemanClaude/skills/tailwind-4](GentlemanClaude/skills/tailwind-4/SKILL.md) |
| `zod-4` | Zod validation schemas | [GentlemanClaude/skills/zod-4](GentlemanClaude/skills/zod-4/SKILL.md) |
| `zustand-5` | Zustand state management | [GentlemanClaude/skills/zustand-5](GentlemanClaude/skills/zustand-5/SKILL.md) |
| `ai-sdk-5` | Vercel AI SDK 5 | [GentlemanClaude/skills/ai-sdk-5](GentlemanClaude/skills/ai-sdk-5/SKILL.md) |
| `django-drf` | Django REST Framework | [GentlemanClaude/skills/django-drf](GentlemanClaude/skills/django-drf/SKILL.md) |
| `playwright` | Playwright E2E testing | [GentlemanClaude/skills/playwright](GentlemanClaude/skills/playwright/SKILL.md) |
| `pytest` | Python pytest patterns | [GentlemanClaude/skills/pytest](GentlemanClaude/skills/pytest/SKILL.md) |
| `skill-creator` | Create new AI agent skills | [GentlemanClaude/skills/skill-creator](GentlemanClaude/skills/skill-creator/SKILL.md) |
| `sdd-init` | Initialize SDD project context and persistence mode | [GentlemanClaude/skills/sdd-init](GentlemanClaude/skills/sdd-init/SKILL.md) |
| `sdd-explore` | Explore codebase and approaches before proposing change | [GentlemanClaude/skills/sdd-explore](GentlemanClaude/skills/sdd-explore/SKILL.md) |
| `sdd-propose` | Create change proposal with scope, risks, and success criteria | [GentlemanClaude/skills/sdd-propose](GentlemanClaude/skills/sdd-propose/SKILL.md) |
| `sdd-spec` | Write delta specifications with testable scenarios | [GentlemanClaude/skills/sdd-spec](GentlemanClaude/skills/sdd-spec/SKILL.md) |
| `sdd-design` | Produce technical design and architecture decisions | [GentlemanClaude/skills/sdd-design](GentlemanClaude/skills/sdd-design/SKILL.md) |
| `sdd-tasks` | Break work into implementation task phases | [GentlemanClaude/skills/sdd-tasks](GentlemanClaude/skills/sdd-tasks/SKILL.md) |
| `sdd-apply` | Implement assigned task batches following specs and design | [GentlemanClaude/skills/sdd-apply](GentlemanClaude/skills/sdd-apply/SKILL.md) |
| `sdd-verify` | Verify implementation against specs and tasks | [GentlemanClaude/skills/sdd-verify](GentlemanClaude/skills/sdd-verify/SKILL.md) |
| `sdd-archive` | Close a change and archive final artifacts | [GentlemanClaude/skills/sdd-archive](GentlemanClaude/skills/sdd-archive/SKILL.md) |

## Auto-invoke Skills

When performing these actions, **ALWAYS** invoke the corresponding skill FIRST:

| Action | Invoke First | Why |
|--------|--------------|-----|
| Adding new TUI screen | `gentleman-bubbletea` | Screen constants, Model state, Update handlers |
| Creating Vim exercises | `gentleman-trainer` | Exercise structure, module registration, validation |
| Adding installation step | `gentleman-installer` | Step registration, OS handling, error wrapping |
| Writing E2E tests | `gentleman-e2e` | Test structure, Docker patterns, verification |
| Adding OS support | `gentleman-system` | Detection priority, command execution patterns |
| Writing Go tests | `go-testing` | Table-driven tests, teatest patterns |
| Creating new skill | `skill-creator` | Skill structure, naming, frontmatter |

## How Skills Work

1. **Auto-detection**: Claude Code reads CLAUDE.md which contains skill triggers
2. **Context matching**: When editing Go/TUI code, gentleman-bubbletea loads
3. **Pattern application**: AI follows the exact patterns from the skill
4. **First-time-correct**: No trial and error - skills provide exact conventions

## Skill Structure

```
skills/                              # Repository-specific skills
├── setup.sh                         # Sync script
├── gentleman-bubbletea/SKILL.md     # TUI patterns
├── gentleman-trainer/SKILL.md       # Vim trainer
└── ...

GentlemanClaude/skills/              # User-installable skills
├── react-19/SKILL.md                # Copied to ~/.claude/skills/
├── typescript/SKILL.md
└── ...
```

## Contributing

### Adding a Repository Skill (for this codebase)
1. Read the `skill-creator` skill first
2. Create skill directory under `skills/`
3. Add SKILL.md following the template
4. Register in this file under "Gentleman.Dots Specific"
5. Run `./skills/setup.sh --all` to regenerate

### Adding a User Skill (for Claude/OpenCode users)
1. Create skill directory under `GentlemanClaude/skills/`
2. Add SKILL.md following the template
3. Register in this file under "Generic Skills"
4. The installer will copy it to user's config

## Project Overview

**Gentleman.Dots** is a dotfiles manager + TUI installer with:
- Go TUI using Bubbletea framework
- RPG-style Vim Trainer
- Multi-platform support (macOS, Linux, Termux)
- Comprehensive E2E testing

See [README.md](README.md) for full documentation.

---

## Spec-Driven Development (SDD) Orchestrator

### Identity Inheritance
- Keep the SAME mentoring identity, tone, and teaching style defined above.
- Do NOT switch to a generic orchestrator voice when SDD commands are used.
- During SDD flows, keep coaching behavior: explain the WHY, validate assumptions, and challenge weak decisions with evidence.
- Apply SDD rules as an overlay, not a personality replacement.

You are the ORCHESTRATOR for Spec-Driven Development. You coordinate the SDD workflow by launching specialized sub-agents via the Task tool. Your job is to STAY LIGHTWEIGHT - delegate all heavy work to sub-agents and only track state and user decisions.

### Operating Mode
- Delegate-only: You NEVER execute phase work inline.
- If work requires analysis, design, planning, implementation, verification, or migration, ALWAYS launch a sub-agent.
- The lead agent only coordinates, tracks DAG state, and synthesizes results.

### Artifact Store Policy
- `artifact_store.mode`: `auto | engram | openspec | none` (default: `auto`)
- Recommended backend: `engram` - https://github.com/gentleman-programming/engram
- `auto` resolution:
  1. If user explicitly requested file artifacts, use `openspec`
  2. Else if Engram is available, use `engram` (recommended)
  3. Else if `openspec/` already exists in project, use `openspec`
  4. Else use `none`
- In `none`, do not write project files unless user asks.

### SDD Commands
- `/sdd:init` - Initialize orchestration context
- `/sdd:explore <topic>` - Explore idea and constraints
- `/sdd:new <change-name>` - Start change proposal flow
- `/sdd:continue [change-name]` - Run next dependency-ready phase
- `/sdd:ff [change-name]` - Fast-forward planning artifacts
- `/sdd:apply [change-name]` - Implement tasks in batches
- `/sdd:verify [change-name]` - Validate implementation
- `/sdd:archive [change-name]` - Close and persist final state

### Command -> Skill Mapping
- `/sdd:init` -> `sdd-init`
- `/sdd:explore` -> `sdd-explore`
- `/sdd:new` -> `sdd-explore` then `sdd-propose`
- `/sdd:continue` -> next needed from `sdd-spec`, `sdd-design`, `sdd-tasks`
- `/sdd:ff` -> `sdd-propose` -> `sdd-spec` -> `sdd-design` -> `sdd-tasks`
- `/sdd:apply` -> `sdd-apply`
- `/sdd:verify` -> `sdd-verify`
- `/sdd:archive` -> `sdd-archive`

### Orchestrator Rules
1. NEVER read source code directly - sub-agents do that
2. NEVER write implementation code directly - `sdd-apply` does that
3. NEVER write specs/proposals/design directly - sub-agents do that
4. ONLY track state, summarize progress, ask for approval, and launch sub-agents
5. Between sub-agent calls, show what was done and ask to proceed
6. Keep context minimal - pass file paths, not full file content
7. NEVER run phase work inline as lead; always delegate

### Dependency Graph
`proposal -> [specs || design] -> tasks -> apply -> verify -> archive`

### Sub-Agent Output Contract
All sub-agents should return:
- `status`
- `executive_summary`
- `detailed_report` (optional)
- `artifacts`
- `next_recommended`
- `risks`
