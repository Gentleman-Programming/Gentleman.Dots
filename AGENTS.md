# Agents Guide for Gentleman.Dots

## Build/Test Commands

- `nix run github:nix-community/home-manager -- switch --flake .#gentleman -b backup` - Apply configuration
- No traditional test commands - this is a dotfiles configuration repository
- Validate Nix syntax: `nix flake check`

## Repository Structure

This is a **dotfiles repository** using Nix flakes for declarative system configuration management. The main entry point is `flake.nix` which orchestrates all configurations.

## Code Style Guidelines

### Nix Files (.nix)

- Use kebab-case for file names (`nushell.nix`, `fish.nix`)
- 2-space indentation
- Use `lib.mkIf` for conditional logic
- Platform detection with `pkgs.stdenv.isDarwin`
- File paths should use `./relative/path` syntax
- Always include proper error handling for cross-platform compatibility

### Lua Files (Neovim config)

- Use snake_case for variables and functions
- 2-space indentation
- Prefer `require()` over `vim.cmd()`
- Comment complex configurations
- Use lazy loading for plugins

### Configuration Files

- Follow original tool conventions (Fish shell, Nushell, etc.)
- Use appropriate comment syntax for each language
- Maintain existing color schemes and themes
- Keep platform-specific configurations clearly separated

### Error Handling

- Always check for tool availability before configuration
- Use conditional blocks for OS-specific features
- Provide fallbacks for missing dependencies

## AI Agent Skills

This repo ships skills for Claude Code and OpenCode. Skills provide on-demand context and patterns.

### Generic Skills (installed to ~/.claude/skills/ and ~/.opencode/skills/)

| Skill | Description |
|-------|-------------|
| `react-19` | React 19 patterns, hooks, components |
| `nextjs-15` | Next.js 15, App Router, Server Components |
| `typescript` | TypeScript patterns, types, generics |
| `tailwind-4` | Tailwind CSS v4 patterns |
| `zod-4` | Zod validation schemas |
| `zustand-5` | Zustand state management |
| `ai-sdk-5` | Vercel AI SDK 5 |
| `django-drf` | Django REST Framework |
| `playwright` | Playwright E2E testing |
| `pytest` | Python pytest patterns |
| `vitest` | Vitest + React Testing Library |
| `tdd` | Test-Driven Development workflow |
| `pr-review` | GitHub PR review |
| `chained-pr` | Chained/stacked PRs workflow |
| `jira-epic` | Jira epic management |
| `jira-task` | Jira story/task management |
| `skill-creator` | Create new AI agent skills |
| `notion-adr` | Architecture Decision Records in Notion |
| `notion-prd` | Product Requirement Documents in Notion |
| `notion-product-brain` | Product ideation in Notion |
| `notion-rfc` | RFCs in Notion |
| `notion-to-jira` | Bridge Notion RFCs to Jira |
| `transcript-processor` | Meeting transcript processing |

### SDD (Spec-Driven Development) Skills

The SDD system uses an orchestrator + sub-agent architecture for structured feature development.

| Skill | Description |
|-------|-------------|
| `sdd-init` | Bootstrap openspec/ directory in a project |
| `sdd-explore` | Investigate ideas before committing to a change |
| `sdd-propose` | Create change proposals with intent, scope, approach |
| `sdd-spec` | Write specifications with requirements and scenarios |
| `sdd-design` | Create technical design with architecture decisions |
| `sdd-tasks` | Break down changes into implementation task checklists |
| `sdd-apply` | Implement tasks, writing actual code |
| `sdd-verify` | Validate implementation matches specs and design |
| `sdd-archive` | Sync delta specs to main specs and archive |

### SDD Commands

| Command | Action |
|---------|--------|
| `/sdd:init` | Bootstrap openspec/ in current project |
| `/sdd:explore <topic>` | Think through an idea |
| `/sdd:new <name>` | Start a new change |
| `/sdd:continue [name]` | Create next artifact |
| `/sdd:ff [name]` | Fast-forward all planning artifacts |
| `/sdd:apply [name]` | Implement tasks |
| `/sdd:verify [name]` | Validate implementation |
| `/sdd:archive [name]` | Sync specs + archive |
