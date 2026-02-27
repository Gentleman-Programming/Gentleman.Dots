# AI Framework Module Registry

Complete reference of all 203 modules across 6 categories available in the [project-starter-framework](https://github.com/JNZader/project-starter-framework). These modules are configured via the TUI installer's Custom mode or the `--ai-modules` CLI flag.

## Table of Contents

- [Overview](#overview)
- [How Features Work](#how-features-work)
- [Categories Summary](#categories-summary)
- [Hooks (10 items)](#-hooks-10-items)
- [Commands (20 items)](#-commands-20-items)
- [Agents (80 items)](#-agents-80-items)
- [Skills (85 items)](#-skills-85-items)
- [SDD — Spec-Driven Development (2 items)](#-sdd--spec-driven-development-2-items)
- [MCP Servers (6 items)](#-mcp-servers-6-items)

---

## Overview

The installer presents 203 individual modules organized into 6 categories. In the TUI, you can browse and toggle individual items within each category. However, `setup-global.sh` operates at the **feature level** — selecting ANY item within a category enables that entire feature.

## How Features Work

```
TUI Category Drill-Down        →    setup-global.sh
━━━━━━━━━━━━━━━━━━━━━━━        →    ━━━━━━━━━━━━━━━━
🪝 Hooks (3/10 selected)       →    --features=hooks
⚡ Commands (0/20 selected)    →    (not included)
🤖 Agents (1/80 selected)     →    --features=agents
🎯 Skills (5/85 selected)     →    --features=skills
📐 SDD: OpenSpec ✓             →    --features=sdd
🔌 MCP (2/6 selected)         →    --features=mcp
```

**Result:** `setup-global.sh --features=hooks,agents,skills,sdd,mcp`

The individual item selection in the TUI is **informational** — it helps you see what's in each category. The actual installation installs ALL items within the enabled features.

**Exception — SDD category:** SDD has special behavior. See [SDD section](#-sdd--spec-driven-development-2-items).

## Categories Summary

| Category | Icon | Items | Feature Flag | Atomic |
|----------|------|------:|:-------------|:------:|
| Hooks | 🪝 | 10 | `hooks` | No |
| Commands | ⚡ | 20 | `commands` | No |
| Agents | 🤖 | 80 | `agents` | No |
| Skills | 🎯 | 85 | `skills` | No |
| SDD | 📐 | 2 | `sdd` (OpenSpec only) | No |
| MCP Servers | 🔌 | 6 | `mcp` | Yes |

**Total: 203 modules**

---

## 🪝 Hooks (10 items)

Hooks are automated actions triggered by events in your AI coding workflow (commit guards, secret scanning, prompt improvement, etc.).

| ID | Label | Description |
|----|-------|-------------|
| `block-dangerous-commands` | Block Dangerous Commands | Prevents execution of destructive shell commands |
| `commit-guard` | Commit Guard | Validates commit messages and staged changes |
| `context-loader` | Context Loader | Auto-loads relevant context for AI sessions |
| `improve-prompt` | Improve Prompt | Enhances prompts before sending to AI |
| `learning-log` | Learning Log | Records learnings from AI interactions |
| `model-router` | Model Router | Routes requests to optimal AI model |
| `secret-scanner` | Secret Scanner | Scans for leaked secrets in code |
| `skill-validator` | Skill Validator | Validates skill file format and structure |
| `task-artifact` | Task Artifact | Manages task artifacts during SDD workflow |
| `validate-workflow` | Validate Workflow | Validates workflow configurations |

---

## ⚡ Commands (20 items)

Slash commands available in your AI coding tool. Organized into 4 groups: Git, Refactoring, Testing, and Workflows.

### Git Commands (7)

| ID | Label |
|----|-------|
| `git:changelog` | Git: Changelog |
| `git:ci-local` | Git: CI Local |
| `git:commit` | Git: Commit |
| `git:fix-issue` | Git: Fix Issue |
| `git:pr-create` | Git: PR Create |
| `git:pr-review` | Git: PR Review |
| `git:worktree` | Git: Worktree |

### Refactoring Commands (3)

| ID | Label |
|----|-------|
| `refactoring:cleanup` | Refactoring: Cleanup |
| `refactoring:dead-code` | Refactoring: Dead Code |
| `refactoring:extract` | Refactoring: Extract |

### Testing Commands (4)

| ID | Label |
|----|-------|
| `testing:e2e` | Testing: E2E |
| `testing:tdd` | Testing: TDD |
| `testing:test-coverage` | Testing: Coverage |
| `testing:test-fix` | Testing: Fix Tests |

### Workflow Commands (6)

| ID | Label |
|----|-------|
| `workflow:generate-agents-md` | Workflow: Generate Agents |
| `workflow:planning` | Workflow: Planning |
| `workflow:compound` | Workflow: Compound |
| `workflow:review` | Workflow: Review |
| `workflow:work` | Workflow: Work |
| `workflow:plan` | Workflow: Plan |

---

## 🤖 Agents (80 items)

Specialized AI agents for different domains. Organized into 8 subcategories.

### Business Agents (9)

| ID | Label |
|----|-------|
| `business-agile-pm` | Business: Agile PM |
| `business-business-analyst` | Business: Business Analyst |
| `business-customer-success` | Business: Customer Success |
| `business-data-analyst` | Business: Data Analyst |
| `business-product-manager` | Business: Product Manager |
| `business-project-estimator` | Business: Project Estimator |
| `business-scrum-master` | Business: Scrum Master |
| `business-stakeholder-communicator` | Business: Stakeholder Communicator |
| `business-tech-writer` | Business: Tech Writer |

### Data & AI Agents (10)

| ID | Label |
|----|-------|
| `data-ai-ai-researcher` | Data & AI: AI Researcher |
| `data-ai-analytics-engineer` | Data & AI: Analytics Engineer |
| `data-ai-data-engineer` | Data & AI: Data Engineer |
| `data-ai-data-pipeline-architect` | Data & AI: Data Pipeline Architect |
| `data-ai-data-scientist` | Data & AI: Data Scientist |
| `data-ai-data-viz-specialist` | Data & AI: Data Viz Specialist |
| `data-ai-feature-engineer` | Data & AI: Feature Engineer |
| `data-ai-ml-engineer` | Data & AI: ML Engineer |
| `data-ai-mlops-engineer` | Data & AI: MLOps Engineer |
| `data-ai-prompt-engineer` | Data & AI: Prompt Engineer |

### Development Agents (15)

| ID | Label |
|----|-------|
| `development-angular-expert` | Development: Angular Expert |
| `development-backend-architect` | Development: Backend Architect |
| `development-database-specialist` | Development: Database Specialist |
| `development-frontend-specialist` | Development: Frontend Specialist |
| `development-fullstack-engineer` | Development: Fullstack Engineer |
| `development-golang-pro` | Development: Go Pro |
| `development-java-enterprise` | Development: Java Enterprise |
| `development-javascript-pro` | Development: JavaScript Pro |
| `development-nextjs-pro` | Development: Next.js Pro |
| `development-python-pro` | Development: Python Pro |
| `development-react-pro` | Development: React Pro |
| `development-rust-pro` | Development: Rust Pro |
| `development-spring-boot-4-expert` | Development: Spring Boot 4 |
| `development-typescript-pro` | Development: TypeScript Pro |
| `development-vue-specialist` | Development: Vue Specialist |

### Infrastructure Agents (7)

| ID | Label |
|----|-------|
| `infrastructure-cloud-architect` | Infrastructure: Cloud Architect |
| `infrastructure-deployment-manager` | Infrastructure: Deployment Manager |
| `infrastructure-devops-engineer` | Infrastructure: DevOps Engineer |
| `infrastructure-incident-responder` | Infrastructure: Incident Responder |
| `infrastructure-kubernetes-expert` | Infrastructure: Kubernetes Expert |
| `infrastructure-monitoring-specialist` | Infrastructure: Monitoring Specialist |
| `infrastructure-performance-engineer` | Infrastructure: Performance Engineer |

### Quality Agents (8)

| ID | Label |
|----|-------|
| `quality-accessibility-auditor` | Quality: Accessibility Auditor |
| `quality-code-reviewer-compact` | Quality: Code Reviewer (Compact) |
| `quality-code-reviewer` | Quality: Code Reviewer |
| `quality-dependency-manager` | Quality: Dependency Manager |
| `quality-e2e-test-specialist` | Quality: E2E Test Specialist |
| `quality-performance-tester` | Quality: Performance Tester |
| `quality-security-auditor` | Quality: Security Auditor |
| `quality-test-engineer` | Quality: Test Engineer |

### Specialists Agents (12)

| ID | Label |
|----|-------|
| `specialists-api-designer` | Specialists: API Designer |
| `specialists-backend-architect` | Specialists: Backend Architect |
| `specialists-code-reviewer` | Specialists: Code Reviewer |
| `specialists-db-optimizer` | Specialists: DB Optimizer |
| `specialists-devops-engineer` | Specialists: DevOps Engineer |
| `specialists-documentation-writer` | Specialists: Documentation Writer |
| `specialists-frontend-developer` | Specialists: Frontend Developer |
| `specialists-performance-analyst` | Specialists: Performance Analyst |
| `specialists-refactor-specialist` | Specialists: Refactor Specialist |
| `specialists-security-auditor` | Specialists: Security Auditor |
| `specialists-test-engineer` | Specialists: Test Engineer |
| `specialists-ux-consultant` | Specialists: UX Consultant |

### Specialized Agents (19)

| ID | Label |
|----|-------|
| `specialized-agent-generator` | Specialized: Agent Generator |
| `specialized-blockchain-developer` | Specialized: Blockchain Developer |
| `specialized-code-migrator` | Specialized: Code Migrator |
| `specialized-context-manager` | Specialized: Context Manager |
| `specialized-documentation-writer` | Specialized: Documentation Writer |
| `specialized-ecommerce-expert` | Specialized: E-Commerce Expert |
| `specialized-embedded-engineer` | Specialized: Embedded Engineer |
| `specialized-error-detective` | Specialized: Error Detective |
| `specialized-fintech-specialist` | Specialized: Fintech Specialist |
| `specialized-freelance-planner` | Specialized: Freelance Planner |
| `specialized-freelance-planner-v2` | Specialized: Freelance Planner v2 |
| `specialized-freelance-planner-v3` | Specialized: Freelance Planner v3 |
| `specialized-freelance-planner-v4` | Specialized: Freelance Planner v4 |
| `specialized-game-developer` | Specialized: Game Developer |
| `specialized-healthcare-dev` | Specialized: Healthcare Dev |
| `specialized-mobile-developer` | Specialized: Mobile Developer |
| `specialized-parallel-plan-executor` | Specialized: Parallel Plan Executor |
| `specialized-plan-executor` | Specialized: Plan Executor |
| `specialized-solo-dev-planner` | Specialized: Solo Dev Planner |
| `specialized-template-writer` | Specialized: Template Writer |
| `specialized-test-runner` | Specialized: Test Runner |
| `specialized-vibekanban-worker` | Specialized: VibeKanban Worker |
| `specialized-wave-executor` | Specialized: Wave Executor |
| `specialized-workflow-optimizer` | Specialized: Workflow Optimizer |

---

## 🎯 Skills (85 items)

Framework standards and coding patterns. Organized into 11 subcategories.

### Backend Skills (21)

| ID | Label |
|----|-------|
| `backend-api-gateway` | Backend: API Gateway |
| `backend-bff-concepts` | Backend: BFF Concepts |
| `backend-bff-spring` | Backend: BFF Spring |
| `backend-chi-router` | Backend: Chi Router |
| `backend-error-handling` | Backend: Error Handling |
| `backend-exceptions-spring` | Backend: Exceptions Spring |
| `backend-fastapi` | Backend: FastAPI |
| `backend-gateway-spring` | Backend: Gateway Spring |
| `backend-go-backend` | Backend: Go Backend |
| `backend-gradle-multimodule` | Backend: Gradle Multi-Module |
| `backend-graphql-concepts` | Backend: GraphQL Concepts |
| `backend-graphql-spring` | Backend: GraphQL Spring |
| `backend-grpc-concepts` | Backend: gRPC Concepts |
| `backend-grpc-spring` | Backend: gRPC Spring |
| `backend-jwt-auth` | Backend: JWT Auth |
| `backend-notifications-concepts` | Backend: Notifications |
| `backend-recommendations-concepts` | Backend: Recommendations |
| `backend-search-concepts` | Backend: Search Concepts |
| `backend-search-spring` | Backend: Search Spring |
| `backend-spring-boot-4` | Backend: Spring Boot 4 |
| `backend-websockets` | Backend: WebSockets |

### Data & AI Skills (11)

| ID | Label |
|----|-------|
| `data-ai-ai-ml` | Data & AI: AI/ML |
| `data-ai-analytics-concepts` | Data & AI: Analytics Concepts |
| `data-ai-analytics-spring` | Data & AI: Analytics Spring |
| `data-ai-duckdb-analytics` | Data & AI: DuckDB Analytics |
| `data-ai-langchain` | Data & AI: LangChain |
| `data-ai-mlflow` | Data & AI: MLflow |
| `data-ai-onnx-inference` | Data & AI: ONNX Inference |
| `data-ai-powerbi` | Data & AI: Power BI |
| `data-ai-pytorch` | Data & AI: PyTorch |
| `data-ai-scikit-learn` | Data & AI: scikit-learn |
| `data-ai-vector-db` | Data & AI: Vector DB |

### Database Skills (6)

| ID | Label |
|----|-------|
| `database-graph-databases` | Database: Graph Databases |
| `database-graph-spring` | Database: Graph Spring |
| `database-pgx-postgres` | Database: PGX Postgres |
| `database-redis-cache` | Database: Redis Cache |
| `database-sqlite-embedded` | Database: SQLite Embedded |
| `database-timescaledb` | Database: TimescaleDB |

### Documentation Skills (4)

| ID | Label |
|----|-------|
| `docs-api-documentation` | Docs: API Documentation |
| `docs-docs-spring` | Docs: Spring Docs |
| `docs-mustache-templates` | Docs: Mustache Templates |
| `docs-technical-docs` | Docs: Technical Docs |

### Frontend Skills (7)

| ID | Label |
|----|-------|
| `frontend-astro-ssr` | Frontend: Astro SSR |
| `frontend-frontend-design` | Frontend: Design Patterns |
| `frontend-frontend-web` | Frontend: Web Development |
| `frontend-mantine-ui` | Frontend: Mantine UI |
| `frontend-tanstack-query` | Frontend: TanStack Query |
| `frontend-zod-validation` | Frontend: Zod Validation |
| `frontend-zustand-state` | Frontend: Zustand State |

### Infrastructure Skills (8)

| ID | Label |
|----|-------|
| `infra-chaos-engineering` | Infrastructure: Chaos Engineering |
| `infra-chaos-spring` | Infrastructure: Chaos Spring |
| `infra-devops-infra` | Infrastructure: DevOps |
| `infra-docker-containers` | Infrastructure: Docker |
| `infra-kubernetes` | Infrastructure: Kubernetes |
| `infra-opentelemetry` | Infrastructure: OpenTelemetry |
| `infra-traefik-proxy` | Infrastructure: Traefik Proxy |
| `infra-woodpecker-ci` | Infrastructure: Woodpecker CI |

### Mobile Skills (2)

| ID | Label |
|----|-------|
| `mobile-ionic-capacitor` | Mobile: Ionic Capacitor |
| `mobile-mobile-ionic` | Mobile: Mobile Ionic |

### Prompt & Quality Skills (2)

| ID | Label |
|----|-------|
| `prompt-improver` | Prompt: Prompt Improver |
| `quality-ghagga-review` | Quality: Ghagga Review |

### References Skills (5)

| ID | Label |
|----|-------|
| `references-hooks-patterns` | References: Hooks Patterns |
| `references-mcp-servers` | References: MCP Servers |
| `references-plugins-reference` | References: Plugins Reference |
| `references-skills-reference` | References: Skills Reference |
| `references-subagent-templates` | References: Subagent Templates |

### Systems & IoT Skills (4)

| ID | Label |
|----|-------|
| `systems-modbus-protocol` | Systems: Modbus Protocol |
| `systems-mqtt-rumqttc` | Systems: MQTT rumqttc |
| `systems-rust-systems` | Systems: Rust Systems |
| `systems-tokio-async` | Systems: Tokio Async |

### Testing Skills (3)

| ID | Label |
|----|-------|
| `testing-playwright-e2e` | Testing: Playwright E2E |
| `testing-testcontainers` | Testing: Testcontainers |
| `testing-vitest-testing` | Testing: Vitest Testing |

### Workflow Skills (12)

| ID | Label |
|----|-------|
| `workflow-ci-local-guide` | Workflow: CI Local Guide |
| `workflow-claude-automation` | Workflow: Claude Automation |
| `workflow-claude-md-improver` | Workflow: CLAUDE.md Improver |
| `workflow-finish-dev-branch` | Workflow: Finish Dev Branch |
| `workflow-git-github` | Workflow: Git & GitHub |
| `workflow-git-workflow` | Workflow: Git Workflow |
| `workflow-ide-plugins` | Workflow: IDE Plugins |
| `workflow-ide-plugins-intellij` | Workflow: IDE Plugins IntelliJ |
| `workflow-obsidian-brain` | Workflow: Obsidian Brain |
| `workflow-git-worktrees` | Workflow: Git Worktrees |
| `workflow-verification` | Workflow: Verification |
| `workflow-wave-workflow` | Workflow: Wave Workflow |

---

## 📐 SDD — Spec-Driven Development (2 items)

The SDD category is **special**. Instead of showing individual SDD phases, it presents two **implementation choices**:

| ID | Label | Description |
|----|-------|-------------|
| `sdd-openspec` | OpenSpec (project-starter-framework) | File-based SDD with YAML schema, installed via `--features=sdd` |
| `sdd-agent-teams` | Agent Teams Lite | Lightweight SDD with 9 Markdown skill files, zero dependencies |

**Behavior:**
- Selecting **OpenSpec** → includes `sdd` in the `--features=` flag for `setup-global.sh`
- Selecting **Agent Teams Lite** → triggers a separate install via `agent-teams-lite/install.sh`
- Selecting **both** → runs both installations
- Selecting **neither** → SDD is not installed

This is different from other categories where any selection enables the whole feature. SDD items have distinct install paths.

See [Agent Teams Lite documentation](agent-teams-lite.md) for details on the Agent Teams Lite option.

---

## 🔌 MCP Servers (6 items)

Model Context Protocol servers for enhanced AI capabilities. MCP is **atomic** — selecting any item enables the `mcp` feature, which installs ALL MCP server configurations.

| ID | Label | Description |
|----|-------|-------------|
| `mcp-context7` | Context7 | Remote MCP for fetching up-to-date library documentation |
| `mcp-engram` | Engram | Local MCP backend for persistent SDD artifacts and memory |
| `mcp-jira` | Jira | Jira integration via Atlassian MCP |
| `mcp-atlassian` | Atlassian | Confluence and Atlassian suite integration |
| `mcp-figma` | Figma | Figma design file integration |
| `mcp-notion` | Notion | Notion workspace integration |
