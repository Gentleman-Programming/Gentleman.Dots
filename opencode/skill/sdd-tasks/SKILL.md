---
name: sdd-tasks
description: >
  Break down a change into an implementation task checklist.
  Trigger: When the orchestrator launches you to create or update the task breakdown for a change.
license: MIT
metadata:
  author: gentleman-programming
  version: "1.0"
---

## Purpose

You are a sub-agent responsible for creating the TASK BREAKDOWN. You take the proposal, specs, and design, then produce a `tasks.md` with concrete, actionable implementation steps organized by phase.

## What You Receive

From the orchestrator:
- Change name
- The `proposal.md` content
- The delta specs from `specs/`
- The `design.md` content
- Project config from `openspec/config.yaml`

## What to Do

### Step 1: Analyze the Design

From the design document, identify:
- All files that need to be created/modified/deleted
- The dependency order (what must come first)
- Testing requirements per component

### Step 2: Write tasks.md

Create the task file:

```
openspec/changes/{change-name}/
├── proposal.md
├── specs/
├── design.md
└── tasks.md               ← You create this
```

#### Task File Format

```markdown
# Tasks: {Change Title}

## Phase 1: {Phase Name} (e.g., Infrastructure / Foundation)

- [ ] 1.1 {Concrete action — what file, what change}
- [ ] 1.2 {Concrete action}
- [ ] 1.3 {Concrete action}

## Phase 2: {Phase Name} (e.g., Core Implementation)

- [ ] 2.1 {Concrete action}
- [ ] 2.2 {Concrete action}
- [ ] 2.3 {Concrete action}
- [ ] 2.4 {Concrete action}

## Phase 3: {Phase Name} (e.g., Testing / Verification)

- [ ] 3.1 {Write tests for ...}
- [ ] 3.2 {Write tests for ...}
- [ ] 3.3 {Verify integration between ...}

## Phase 4: {Phase Name} (e.g., Cleanup / Documentation)

- [ ] 4.1 {Update docs/comments}
- [ ] 4.2 {Remove temporary code}
```

### Task Writing Rules

Each task MUST be:

| Criteria | Example ✅ | Anti-example ❌ |
|----------|-----------|----------------|
| **Specific** | "Create `internal/auth/middleware.go` with JWT validation" | "Add auth" |
| **Actionable** | "Add `ValidateToken()` method to `AuthService`" | "Handle tokens" |
| **Verifiable** | "Test: `POST /login` returns 401 without token" | "Make sure it works" |
| **Small** | One file or one logical unit of work | "Implement the feature" |

### Phase Organization Guidelines

```
Phase 1: Foundation / Infrastructure
  └─ New types, interfaces, database changes, config
  └─ Things other tasks depend on

Phase 2: Core Implementation
  └─ Main logic, business rules, core behavior
  └─ The meat of the change

Phase 3: Integration / Wiring
  └─ Connect components, routes, UI wiring
  └─ Make everything work together

Phase 4: Testing
  └─ Unit tests, integration tests, e2e tests
  └─ Verify against spec scenarios

Phase 5: Cleanup (if needed)
  └─ Documentation, remove dead code, polish
```

### Step 3: Return Summary

Return to the orchestrator:

```markdown
## Tasks Created

**Change**: {change-name}
**Location**: openspec/changes/{change-name}/tasks.md

### Breakdown
| Phase | Tasks | Focus |
|-------|-------|-------|
| Phase 1 | {N} | {Phase name} |
| Phase 2 | {N} | {Phase name} |
| Phase 3 | {N} | {Phase name} |
| Total | {N} | |

### Implementation Order
{Brief description of the recommended order and why}

### Next Step
Ready for implementation (sdd-apply).
```

## Rules

- ALWAYS reference concrete file paths in tasks
- Tasks MUST be ordered by dependency — Phase 1 tasks shouldn't depend on Phase 2
- Testing tasks should reference specific scenarios from the specs
- Each task should be completable in ONE session (if a task feels too big, split it)
- Use hierarchical numbering: 1.1, 1.2, 2.1, 2.2, etc.
- NEVER include vague tasks like "implement feature" or "add tests"
- Apply any `rules.tasks` from `openspec/config.yaml`
- If the TDD skill applies, integrate test-first tasks: RED task → GREEN task → REFACTOR task
