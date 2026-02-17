---
name: sdd-init
description: >
  Bootstrap the openspec/ directory structure for Spec-Driven Development in any project.
  Trigger: When user wants to initialize SDD in a project, or says "sdd init", "iniciar sdd", "openspec init".
license: MIT
metadata:
  author: gentleman-programming
  version: "1.0"
---

## Purpose

You are a sub-agent responsible for bootstrapping the Spec-Driven Development (SDD) structure in a project. You initialize the `openspec/` directory and optionally create the project config.

## What to Do

### Step 1: Detect Project Context

Read the project to understand:
- Tech stack (check package.json, go.mod, pyproject.toml, etc.)
- Existing conventions (linters, test frameworks, CI)
- Architecture patterns in use

### Step 2: Create Directory Structure

```
openspec/
├── config.yaml              ← Project-specific SDD config
├── specs/                   ← Source of truth (empty initially)
└── changes/                 ← Active changes
    └── archive/             ← Completed changes
```

### Step 3: Generate config.yaml

Based on what you detected, create the config:

```yaml
# openspec/config.yaml
schema: spec-driven

context: |
  Tech stack: {detected stack}
  Architecture: {detected patterns}
  Testing: {detected test framework}
  Style: {detected linting/formatting}

rules:
  proposal:
    - Include rollback plan for risky changes
    - Identify affected modules/packages
  specs:
    - Use Given/When/Then format for scenarios
    - Use RFC 2119 keywords (MUST, SHALL, SHOULD, MAY)
  design:
    - Include sequence diagrams for complex flows
    - Document architecture decisions with rationale
  tasks:
    - Group tasks by phase (infrastructure, implementation, testing)
    - Use hierarchical numbering (1.1, 1.2, etc.)
    - Keep tasks small enough to complete in one session
```

### Step 4: Return Summary

Return a structured summary:

```
## SDD Initialized

**Project**: {project name}
**Stack**: {detected stack}
**Location**: openspec/

### Structure Created
- openspec/config.yaml ← Project config with detected context
- openspec/specs/      ← Ready for specifications
- openspec/changes/    ← Ready for change proposals

### Next Steps
Tell the orchestrator: "Ready for /sdd:new <change-name>"
```

## Rules

- NEVER create placeholder spec files - specs are created via sdd-spec during a change
- ALWAYS detect the real tech stack, don't guess
- If the project already has an `openspec/` directory, report what exists and ask the orchestrator if it should be updated
- Keep config.yaml context CONCISE - no more than 10 lines
