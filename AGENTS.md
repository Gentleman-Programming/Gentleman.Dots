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
| `vitest` | Vitest + React Testing Library | [GentlemanClaude/skills/vitest](GentlemanClaude/skills/vitest/SKILL.md) |
| `tdd` | Test-Driven Development workflow | [GentlemanClaude/skills/tdd](GentlemanClaude/skills/tdd/SKILL.md) |
| `pr-review` | GitHub PR review | [GentlemanClaude/skills/pr-review](GentlemanClaude/skills/pr-review/SKILL.md) |
| `chained-pr` | Chained/stacked PRs workflow | [GentlemanClaude/skills/chained-pr](GentlemanClaude/skills/chained-pr/SKILL.md) |
| `jira-epic` | Jira epic management | [GentlemanClaude/skills/jira-epic](GentlemanClaude/skills/jira-epic/SKILL.md) |
| `jira-task` | Jira story/task management | [GentlemanClaude/skills/jira-task](GentlemanClaude/skills/jira-task/SKILL.md) |
| `notion-adr` | Architecture Decision Records in Notion | [GentlemanClaude/skills/notion-adr](GentlemanClaude/skills/notion-adr/SKILL.md) |
| `notion-prd` | Product Requirement Documents in Notion | [GentlemanClaude/skills/notion-prd](GentlemanClaude/skills/notion-prd/SKILL.md) |
| `notion-product-brain` | Product ideation in Notion | [GentlemanClaude/skills/notion-product-brain](GentlemanClaude/skills/notion-product-brain/SKILL.md) |
| `notion-rfc` | RFCs in Notion | [GentlemanClaude/skills/notion-rfc](GentlemanClaude/skills/notion-rfc/SKILL.md) |
| `notion-to-jira` | Bridge Notion RFCs to Jira | [GentlemanClaude/skills/notion-to-jira](GentlemanClaude/skills/notion-to-jira/SKILL.md) |
| `transcript-processor` | Meeting transcript processing | [GentlemanClaude/skills/transcript-processor](GentlemanClaude/skills/transcript-processor/SKILL.md) |
| `skill-creator` | Create new AI agent skills | [GentlemanClaude/skills/skill-creator](GentlemanClaude/skills/skill-creator/SKILL.md) |

### SDD (Spec-Driven Development)

Sub-agent skills used by the SDD orchestrator for structured planning and implementation.

| Skill | Description | Source |
|-------|-------------|--------|
| `sdd-init` | Bootstrap openspec/ directory | [GentlemanClaude/skills/sdd-init](GentlemanClaude/skills/sdd-init/SKILL.md) |
| `sdd-explore` | Investigate ideas before committing | [GentlemanClaude/skills/sdd-explore](GentlemanClaude/skills/sdd-explore/SKILL.md) |
| `sdd-propose` | Create change proposals | [GentlemanClaude/skills/sdd-propose](GentlemanClaude/skills/sdd-propose/SKILL.md) |
| `sdd-spec` | Write specifications | [GentlemanClaude/skills/sdd-spec](GentlemanClaude/skills/sdd-spec/SKILL.md) |
| `sdd-design` | Technical design documents | [GentlemanClaude/skills/sdd-design](GentlemanClaude/skills/sdd-design/SKILL.md) |
| `sdd-tasks` | Implementation task checklists | [GentlemanClaude/skills/sdd-tasks](GentlemanClaude/skills/sdd-tasks/SKILL.md) |
| `sdd-apply` | Implement tasks, write code | [GentlemanClaude/skills/sdd-apply](GentlemanClaude/skills/sdd-apply/SKILL.md) |
| `sdd-verify` | Validate implementation | [GentlemanClaude/skills/sdd-verify](GentlemanClaude/skills/sdd-verify/SKILL.md) |
| `sdd-archive` | Sync specs and archive | [GentlemanClaude/skills/sdd-archive](GentlemanClaude/skills/sdd-archive/SKILL.md) |

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
| Starting new feature/change | `sdd-propose` | Structured planning before coding |
| Writing specs/requirements | `sdd-spec` | Spec format and scenarios |
| Implementing from specs | `sdd-apply` | Task-driven implementation |

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
├── sdd-init/SKILL.md                # SDD sub-agent skills
├── sdd-explore/SKILL.md
├── sdd-propose/SKILL.md
├── sdd-spec/SKILL.md
├── sdd-design/SKILL.md
├── sdd-tasks/SKILL.md
├── sdd-apply/SKILL.md
├── sdd-verify/SKILL.md
├── sdd-archive/SKILL.md
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
