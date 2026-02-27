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

Hooks are automated actions triggered by events in your AI coding workflow. They run before or after specific actions (commits, command execution, prompt submission) to enforce quality, security, and workflow standards.

| ID | Label | Description |
|----|-------|-------------|
| `block-dangerous-commands` | Block Dangerous Commands | Intercepts and blocks destructive shell commands (`rm -rf /`, `DROP TABLE`, force pushes) before execution |
| `commit-guard` | Commit Guard | Validates commit messages follow Conventional Commits format and staged changes pass linting |
| `context-loader` | Context Loader | Auto-loads relevant project context (CLAUDE.md, package.json, tsconfig) at session start |
| `improve-prompt` | Improve Prompt | Enhances vague or incomplete prompts with context, constraints, and structure before AI processing |
| `learning-log` | Learning Log | Records insights, patterns, and corrections from AI interactions into a persistent learnings file |
| `model-router` | Model Router | Routes AI requests to the optimal model based on task complexity (Haiku for simple, Opus for complex) |
| `secret-scanner` | Secret Scanner | Scans staged changes for leaked secrets, API keys, tokens, and credentials before commit |
| `skill-validator` | Skill Validator | Validates SKILL.md files conform to the agent skills specification format |
| `task-artifact` | Task Artifact | Manages SDD task artifacts, tracking status and outputs across workflow phases |
| `validate-workflow` | Validate Workflow | Validates CI/CD workflow files (GitHub Actions, etc.) for syntax and best practices |

---

## ⚡ Commands (20 items)

Slash commands available in your AI coding tool. Type the command name (e.g., `/git:commit`) to trigger structured workflows for common development tasks.

### Git Commands (7)

Version control workflows for commits, PRs, changelogs, CI, and branch management.

| ID | Label | Description |
|----|-------|-------------|
| `git:changelog` | Git: Changelog | Generate a CHANGELOG from recent commits following Keep a Changelog format |
| `git:ci-local` | Git: CI Local | Run GitHub Actions CI pipeline locally before pushing using `act` or `wrkflw` |
| `git:commit` | Git: Commit | Analyze staged changes and generate a Conventional Commits message with scope and body |
| `git:fix-issue` | Git: Fix Issue | Implement the fix for a specific GitHub issue following the full branch→fix→test→PR flow |
| `git:pr-create` | Git: PR Create | Create a Pull Request for the current branch with structured description and checklist |
| `git:pr-review` | Git: PR Review | Review the current PR focusing on correctness, security, performance, and maintainability |
| `git:worktree` | Git: Worktree | Manage Git worktrees for working on multiple branches in parallel |

### Refactoring Commands (3)

Code improvement workflows that preserve behavior while improving structure and cleanliness.

| ID | Label | Description |
|----|-------|-------------|
| `refactoring:cleanup` | Refactoring: Cleanup | Clean up code in the current file/module without changing behavior (formatting, naming, simplification) |
| `refactoring:dead-code` | Refactoring: Dead Code | Identify and remove unreachable or unused code, imports, variables, and functions |
| `refactoring:extract` | Refactoring: Extract | Extract duplicated or complex logic into separate functions, modules, or components |

### Testing Commands (4)

Testing workflows for TDD, E2E, coverage analysis, and test fixes.

| ID | Label | Description |
|----|-------|-------------|
| `testing:e2e` | Testing: E2E | Write or run end-to-end tests to verify complete system flows (Playwright, Cypress) |
| `testing:tdd` | Testing: TDD | Implement a feature using Test-Driven Development: write failing test → implement → refactor |
| `testing:test-coverage` | Testing: Coverage | Analyze test coverage, identify critical gaps, and suggest tests for uncovered paths |
| `testing:test-fix` | Testing: Fix Tests | Fix failing tests without breaking business logic — diagnose root cause first |

### Workflow Commands (6)

Meta-workflows for planning, reviewing, executing tasks, and generating documentation.

| ID | Label | Description |
|----|-------|-------------|
| `workflow:generate-agents-md` | Workflow: Generate Agents | Inspect the project and auto-generate a cross-agent compatible AGENTS.md file |
| `workflow:planning` | Workflow: Planning | Create a structured plan before implementing a feature or complex change |
| `workflow:compound` | Workflow: Compound | Post-completion workflow: summarize learnings, update learnings.md, suggest CLAUDE.md improvements |
| `workflow:review` | Workflow: Review | Multi-perspective code review using code-reviewer and security-auditor agents on staged changes |
| `workflow:work` | Workflow: Work | Execute the current plan by creating a git worktree, checking off tasks, and tracking files edited |
| `workflow:plan` | Workflow: Plan | Decompose a feature into structured tasks with goals, acceptance criteria, estimates, and dependencies |

---

## 🤖 Agents (80 items)

Specialized AI personas with deep domain expertise. Each agent has a system prompt that shapes its behavior, knowledge focus, and communication style. Switch agents based on the task at hand.

### Business Agents (9)

Agents focused on project management, business analysis, stakeholder communication, and technical writing. Useful for non-code tasks like planning sprints, writing specs, and estimating effort.

| ID | Label | Description |
|----|-------|-------------|
| `business-agile-pm` | Business: Agile PM | Agile project management with sprint planning, velocity tracking, and retrospectives |
| `business-business-analyst` | Business: Business Analyst | Requirements gathering, user stories, acceptance criteria, and stakeholder interviews |
| `business-customer-success` | Business: Customer Success | Customer onboarding flows, support documentation, and satisfaction metrics |
| `business-data-analyst` | Business: Data Analyst | Data analysis, reporting dashboards, SQL queries, and business intelligence |
| `business-product-manager` | Business: Product Manager | Product roadmaps, feature prioritization, market analysis, and PRDs |
| `business-project-estimator` | Business: Project Estimator | Task estimation using story points, T-shirt sizing, and historical velocity |
| `business-scrum-master` | Business: Scrum Master | Scrum ceremonies, impediment removal, team velocity, and process improvement |
| `business-stakeholder-communicator` | Business: Stakeholder Communicator | Status reports, executive summaries, and non-technical explanations of technical work |
| `business-tech-writer` | Business: Tech Writer | Technical documentation, API docs, user guides, and README files |

### Data & AI Agents (10)

Agents specialized in data engineering, machine learning, analytics, and AI research. From ETL pipelines to model training and deployment.

| ID | Label | Description |
|----|-------|-------------|
| `data-ai-ai-researcher` | Data & AI: AI Researcher | AI/ML research, paper summaries, model comparison, and emerging techniques |
| `data-ai-analytics-engineer` | Data & AI: Analytics Engineer | Analytics pipelines, dbt models, data warehousing, and metrics definitions |
| `data-ai-data-engineer` | Data & AI: Data Engineer | ETL/ELT pipelines, data lakes, Spark, Airflow, and data quality |
| `data-ai-data-pipeline-architect` | Data & AI: Data Pipeline Architect | End-to-end data pipeline design, stream processing, and batch architectures |
| `data-ai-data-scientist` | Data & AI: Data Scientist | Statistical analysis, hypothesis testing, feature engineering, and model selection |
| `data-ai-data-viz-specialist` | Data & AI: Data Viz Specialist | Data visualization with D3, Plotly, Matplotlib, and dashboard design |
| `data-ai-feature-engineer` | Data & AI: Feature Engineer | Feature extraction, transformation pipelines, and feature stores |
| `data-ai-ml-engineer` | Data & AI: ML Engineer | Model training, hyperparameter tuning, distributed training, and serving |
| `data-ai-mlops-engineer` | Data & AI: MLOps Engineer | ML model deployment, monitoring, A/B testing, and CI/CD for ML |
| `data-ai-prompt-engineer` | Data & AI: Prompt Engineer | Prompt optimization, few-shot examples, chain-of-thought, and evaluation |

### Development Agents (15)

Core software development agents covering frontend, backend, and full-stack across multiple languages and frameworks.

| ID | Label | Description |
|----|-------|-------------|
| `development-angular-expert` | Development: Angular Expert | Angular 17+ with signals, standalone components, control flow, and RxJS |
| `development-backend-architect` | Development: Backend Architect | Backend architecture: microservices, event-driven, CQRS, domain-driven design |
| `development-database-specialist` | Development: Database Specialist | Database design, query optimization, indexing strategies, and migrations |
| `development-frontend-specialist` | Development: Frontend Specialist | Frontend architecture: component patterns, state management, and performance |
| `development-fullstack-engineer` | Development: Fullstack Engineer | End-to-end development: API design, frontend integration, and deployment |
| `development-golang-pro` | Development: Go Pro | Go development with concurrency patterns, error handling, and standard library |
| `development-java-enterprise` | Development: Java Enterprise | Java enterprise development with Spring, JPA, and enterprise patterns |
| `development-javascript-pro` | Development: JavaScript Pro | Modern JavaScript (ES2024+), async patterns, and runtime optimization |
| `development-nextjs-pro` | Development: Next.js Pro | Next.js 15 with App Router, Server Components, Server Actions, and ISR |
| `development-python-pro` | Development: Python Pro | Python development with typing, async, and modern tooling (uv, ruff) |
| `development-react-pro` | Development: React Pro | React 19 with React Compiler, Server Components, advanced hooks, and performance |
| `development-rust-pro` | Development: Rust Pro | Rust development with ownership, lifetimes, traits, and async runtime |
| `development-spring-boot-4-expert` | Development: Spring Boot 4 | Spring Boot 4 with virtual threads, GraalVM native, and Spring Security 7 |
| `development-typescript-pro` | Development: TypeScript Pro | Advanced TypeScript: generics, conditional types, mapped types, and type safety |
| `development-vue-specialist` | Development: Vue Specialist | Vue 3 with Composition API, Pinia, and TypeScript integration |

### Infrastructure Agents (7)

Agents for cloud architecture, DevOps, Kubernetes, monitoring, and incident response. From infrastructure-as-code to production troubleshooting.

| ID | Label | Description |
|----|-------|-------------|
| `infrastructure-cloud-architect` | Infrastructure: Cloud Architect | Cloud architecture on AWS/GCP/Azure: networking, IAM, and cost optimization |
| `infrastructure-deployment-manager` | Infrastructure: Deployment Manager | Deployment strategies: blue-green, canary, rolling, and feature flags |
| `infrastructure-devops-engineer` | Infrastructure: DevOps Engineer | CI/CD pipelines, infrastructure as code (Terraform, Pulumi), and automation |
| `infrastructure-incident-responder` | Infrastructure: Incident Responder | Incident triage, root cause analysis, postmortems, and runbook creation |
| `infrastructure-kubernetes-expert` | Infrastructure: Kubernetes Expert | Kubernetes: manifests, Helm charts, operators, networking, and troubleshooting |
| `infrastructure-monitoring-specialist` | Infrastructure: Monitoring Specialist | Observability: Prometheus, Grafana, alerting rules, SLOs, and dashboards |
| `infrastructure-performance-engineer` | Infrastructure: Performance Engineer | Performance profiling, load testing, bottleneck identification, and optimization |

### Quality Agents (8)

Agents focused on code quality, testing, security, accessibility, and dependency management.

| ID | Label | Description |
|----|-------|-------------|
| `quality-accessibility-auditor` | Quality: Accessibility Auditor | WCAG 2.1 AA/AAA compliance, ARIA patterns, screen reader testing, and color contrast |
| `quality-code-reviewer-compact` | Quality: Code Reviewer (Compact) | Lightweight code review focused on critical issues only (security, bugs, performance) |
| `quality-code-reviewer` | Quality: Code Reviewer | Comprehensive code review: architecture, patterns, naming, error handling, and tests |
| `quality-dependency-manager` | Quality: Dependency Manager | Dependency auditing, version updates, vulnerability scanning, and lock file management |
| `quality-e2e-test-specialist` | Quality: E2E Test Specialist | End-to-end test strategy, Page Object patterns, test data management, and CI integration |
| `quality-performance-tester` | Quality: Performance Tester | Load testing with k6/Artillery, performance benchmarks, and regression detection |
| `quality-security-auditor` | Quality: Security Auditor | OWASP Top 10 auditing, dependency vulnerabilities, auth/authz review, and hardening |
| `quality-test-engineer` | Quality: Test Engineer | Test strategy, unit/integration/E2E testing, mocking, and coverage analysis |

### Specialists Agents (12)

Focused domain specialists for API design, database optimization, security, documentation, and UX consulting.

| ID | Label | Description |
|----|-------|-------------|
| `specialists-api-designer` | Specialists: API Designer | RESTful API design: resource modeling, versioning, pagination, and OpenAPI specs |
| `specialists-backend-architect` | Specialists: Backend Architect | Backend system design: service boundaries, data flow, and integration patterns |
| `specialists-code-reviewer` | Specialists: Code Reviewer | In-depth code review with architecture and design pattern focus |
| `specialists-db-optimizer` | Specialists: DB Optimizer | Database query optimization, index tuning, execution plans, and schema design |
| `specialists-devops-engineer` | Specialists: DevOps Engineer | DevOps practices: GitOps, infrastructure automation, and deployment pipelines |
| `specialists-documentation-writer` | Specialists: Documentation Writer | Technical writing: architecture docs, runbooks, ADRs, and onboarding guides |
| `specialists-frontend-developer` | Specialists: Frontend Developer | Frontend development: component architecture, CSS-in-JS, and build optimization |
| `specialists-performance-analyst` | Specialists: Performance Analyst | Application performance analysis: profiling, flame graphs, and optimization strategies |
| `specialists-refactor-specialist` | Specialists: Refactor Specialist | Code refactoring: extract method, replace conditional, introduce pattern, and migration |
| `specialists-security-auditor` | Specialists: Security Auditor | Security review: authentication flows, input validation, encryption, and compliance |
| `specialists-test-engineer` | Specialists: Test Engineer | Testing strategy and implementation across all levels (unit, integration, E2E, contract) |
| `specialists-ux-consultant` | Specialists: UX Consultant | UX review: usability heuristics, user flows, accessibility, and design system alignment |

### Specialized Agents (19)

Niche agents for specific domains like blockchain, game dev, mobile, healthcare, fintech, and advanced workflow patterns.

| ID | Label | Description |
|----|-------|-------------|
| `specialized-agent-generator` | Specialized: Agent Generator | Creates new AI agent definitions with system prompts, tools, and skill files |
| `specialized-blockchain-developer` | Specialized: Blockchain Developer | Smart contract development (Solidity, Rust), DeFi protocols, and Web3 integration |
| `specialized-code-migrator` | Specialized: Code Migrator | Code migration between frameworks, languages, or major versions (e.g., Angular→React, v2→v3) |
| `specialized-context-manager` | Specialized: Context Manager | Manages AI context windows: summarization, compression, and relevant context selection |
| `specialized-documentation-writer` | Specialized: Documentation Writer | Project documentation generation: README, API docs, architecture diagrams, and changelogs |
| `specialized-ecommerce-expert` | Specialized: E-Commerce Expert | E-commerce systems: payment integration, cart flows, inventory, and order management |
| `specialized-embedded-engineer` | Specialized: Embedded Engineer | Embedded systems: firmware, RTOS, hardware interfaces, and resource-constrained development |
| `specialized-error-detective` | Specialized: Error Detective | Bug investigation: stack trace analysis, reproduction steps, and root cause identification |
| `specialized-fintech-specialist` | Specialized: Fintech Specialist | Financial systems: payment processing, compliance (PCI-DSS, SOX), and trading platforms |
| `specialized-freelance-planner` | Specialized: Freelance Planner | Freelance project planning: scope, timeline, milestones, and client communication |
| `specialized-freelance-planner-v2` | Specialized: Freelance Planner v2 | Enhanced freelance planning with risk assessment and resource allocation |
| `specialized-freelance-planner-v3` | Specialized: Freelance Planner v3 | Advanced freelance planning with portfolio management and pricing strategies |
| `specialized-freelance-planner-v4` | Specialized: Freelance Planner v4 | Full freelance business planning with contracts, invoicing, and client lifecycle |
| `specialized-game-developer` | Specialized: Game Developer | Game development: Unity/Godot, game loops, physics, ECS, and asset pipelines |
| `specialized-healthcare-dev` | Specialized: Healthcare Dev | Healthcare systems: HIPAA compliance, HL7/FHIR standards, and EHR integration |
| `specialized-mobile-developer` | Specialized: Mobile Developer | Mobile development: React Native, Flutter, native APIs, and app store deployment |
| `specialized-parallel-plan-executor` | Specialized: Parallel Plan Executor | Executes multiple independent plan tasks concurrently using sub-agents |
| `specialized-plan-executor` | Specialized: Plan Executor | Sequential plan execution: reads task list, implements each step, and tracks progress |
| `specialized-solo-dev-planner` | Specialized: Solo Dev Planner | Solo developer project planning with realistic timelines and MVP prioritization |
| `specialized-template-writer` | Specialized: Template Writer | Generates project templates, boilerplate, and scaffold configurations |
| `specialized-test-runner` | Specialized: Test Runner | Runs test suites, analyzes failures, and reports results with fix suggestions |
| `specialized-vibekanban-worker` | Specialized: VibeKanban Worker | Kanban-style task execution with WIP limits and flow optimization |
| `specialized-wave-executor` | Specialized: Wave Executor | Executes tasks in dependency-ordered waves (parallelizes within each wave) |
| `specialized-workflow-optimizer` | Specialized: Workflow Optimizer | Analyzes and optimizes development workflows, CI pipelines, and automation |

---

## 🎯 Skills (85 items)

Framework standards, coding patterns, and technology-specific best practices. Skills are reference documents that AI agents consult when working with specific technologies. They define patterns, anti-patterns, and conventions to follow.

### Backend Skills (21)

Server-side development patterns covering API design, microservices, authentication, and multiple frameworks (Spring Boot, FastAPI, Go, GraphQL, gRPC).

| ID | Label | Description |
|----|-------|-------------|
| `backend-api-gateway` | Backend: API Gateway | API gateway patterns: routing, rate limiting, authentication, and service mesh |
| `backend-bff-concepts` | Backend: BFF Concepts | Backend-for-Frontend pattern: client-specific APIs, aggregation, and data shaping |
| `backend-bff-spring` | Backend: BFF Spring | BFF implementation with Spring Boot: WebClient, reactive endpoints, and circuit breakers |
| `backend-chi-router` | Backend: Chi Router | Go HTTP routing with chi: middleware chains, URL parameters, and group routing |
| `backend-error-handling` | Backend: Error Handling | Backend error handling patterns: error types, propagation, HTTP mapping, and logging |
| `backend-exceptions-spring` | Backend: Exceptions Spring | Spring exception handling: @ControllerAdvice, ProblemDetail, and RFC 7807 responses |
| `backend-fastapi` | Backend: FastAPI | FastAPI development with Pydantic v2, async services, dependency injection, and OpenAPI |
| `backend-gateway-spring` | Backend: Gateway Spring | Spring Cloud Gateway: route predicates, filters, load balancing, and circuit breaking |
| `backend-go-backend` | Backend: Go Backend | Go backend patterns: project layout, interfaces, error handling, and testing |
| `backend-gradle-multimodule` | Backend: Gradle Multi-Module | Gradle multi-module project structure with version catalogs and convention plugins |
| `backend-graphql-concepts` | Backend: GraphQL Concepts | GraphQL fundamentals: schema design, resolvers, N+1 prevention, and subscriptions |
| `backend-graphql-spring` | Backend: GraphQL Spring | Spring for GraphQL: schema-first, DataLoader, security, and subscription support |
| `backend-grpc-concepts` | Backend: gRPC Concepts | gRPC fundamentals: protobuf, service definitions, streaming, and error handling |
| `backend-grpc-spring` | Backend: gRPC Spring | gRPC with Spring Boot: server/client config, interceptors, and health checks |
| `backend-jwt-auth` | Backend: JWT Auth | JWT authentication: token generation, validation, refresh tokens, and RBAC |
| `backend-notifications-concepts` | Backend: Notifications | Notification systems: channels (email, push, SMS), templates, and delivery tracking |
| `backend-recommendations-concepts` | Backend: Recommendations | Recommendation engine patterns: collaborative filtering, content-based, and hybrid |
| `backend-search-concepts` | Backend: Search Concepts | Search fundamentals: full-text search, facets, fuzzy matching, and relevance tuning |
| `backend-search-spring` | Backend: Search Spring | Search with Spring: Elasticsearch/OpenSearch integration, indexing, and query DSL |
| `backend-spring-boot-4` | Backend: Spring Boot 4 | Spring Boot 4 with virtual threads, GraalVM native image, and Spring Security 7 |
| `backend-websockets` | Backend: WebSockets | WebSocket implementation: connection management, rooms, heartbeats, and scaling |

### Data & AI Skills (11)

Machine learning, analytics, and data processing patterns. Covers model training, inference, data pipelines, and visualization tools.

| ID | Label | Description |
|----|-------|-------------|
| `data-ai-ai-ml` | Data & AI: AI/ML | Machine learning fundamentals: model selection, training loops, evaluation metrics |
| `data-ai-analytics-concepts` | Data & AI: Analytics Concepts | Analytics patterns: event tracking, funnel analysis, cohort analysis, and A/B testing |
| `data-ai-analytics-spring` | Data & AI: Analytics Spring | Analytics with Spring: event collection, aggregation pipelines, and reporting endpoints |
| `data-ai-duckdb-analytics` | Data & AI: DuckDB Analytics | DuckDB for analytics: in-process OLAP, Parquet files, and SQL analytics |
| `data-ai-langchain` | Data & AI: LangChain | LangChain patterns: chains, agents, RAG, vector stores, and memory management |
| `data-ai-mlflow` | Data & AI: MLflow | MLflow experiment tracking, model registry, and deployment workflows |
| `data-ai-onnx-inference` | Data & AI: ONNX Inference | ONNX Runtime inference: model conversion, optimization, and cross-platform deployment |
| `data-ai-powerbi` | Data & AI: Power BI | Power BI: DAX formulas, data modeling, report design, and embedded analytics |
| `data-ai-pytorch` | Data & AI: PyTorch | PyTorch: custom datasets, training loops, distributed training, and model export |
| `data-ai-scikit-learn` | Data & AI: scikit-learn | scikit-learn: pipelines, cross-validation, feature selection, and model persistence |
| `data-ai-vector-db` | Data & AI: Vector DB | Vector databases: embeddings, similarity search, indexing (HNSW, IVF), and RAG |

### Database Skills (6)

Database-specific patterns for relational, graph, time-series, and caching systems.

| ID | Label | Description |
|----|-------|-------------|
| `database-graph-databases` | Database: Graph Databases | Graph database patterns: modeling, Cypher/Gremlin queries, and traversal optimization |
| `database-graph-spring` | Database: Graph Spring | Spring Data Neo4j: entity mapping, repository patterns, and graph projections |
| `database-pgx-postgres` | Database: PGX Postgres | PostgreSQL with pgx (Go): connection pools, prepared statements, and COPY protocol |
| `database-redis-cache` | Database: Redis Cache | Redis caching strategies: TTL, invalidation, pub/sub, and data structures |
| `database-sqlite-embedded` | Database: SQLite Embedded | SQLite for embedded use: WAL mode, migrations, FTS5, and concurrent access |
| `database-timescaledb` | Database: TimescaleDB | TimescaleDB: hypertables, continuous aggregates, compression, and time-series queries |

### Documentation Skills (4)

Documentation generation and templating patterns for APIs, architecture, and project docs.

| ID | Label | Description |
|----|-------|-------------|
| `docs-api-documentation` | Docs: API Documentation | API documentation: OpenAPI/Swagger, endpoint descriptions, examples, and versioning |
| `docs-docs-spring` | Docs: Spring Docs | Spring REST Docs: test-driven API documentation with Asciidoctor |
| `docs-mustache-templates` | Docs: Mustache Templates | Mustache/Handlebars templates for code generation and document rendering |
| `docs-technical-docs` | Docs: Technical Docs | Technical documentation: ADRs, architecture guides, runbooks, and onboarding docs |

### Frontend Skills (7)

Client-side development patterns for modern web applications, including component design, state management, data fetching, and validation.

| ID | Label | Description |
|----|-------|-------------|
| `frontend-astro-ssr` | Frontend: Astro SSR | Astro SSR: islands architecture, content collections, and hybrid rendering |
| `frontend-frontend-design` | Frontend: Design Patterns | Frontend design patterns: atomic design, compound components, and render props |
| `frontend-frontend-web` | Frontend: Web Development | Core web development: HTML semantics, CSS layout, accessibility, and performance |
| `frontend-mantine-ui` | Frontend: Mantine UI | Mantine UI components: theming, form handling, and custom component composition |
| `frontend-tanstack-query` | Frontend: TanStack Query | TanStack Query: cache invalidation, optimistic updates, infinite queries, and prefetching |
| `frontend-zod-validation` | Frontend: Zod Validation | Zod schema validation: type inference, transforms, refinements, and form integration |
| `frontend-zustand-state` | Frontend: Zustand State | Zustand state management: slices, middleware, persistence, and devtools integration |

### Infrastructure Skills (8)

Infrastructure, DevOps, and platform engineering patterns. Covers containers, orchestration, CI/CD, observability, and chaos engineering.

| ID | Label | Description |
|----|-------|-------------|
| `infra-chaos-engineering` | Infrastructure: Chaos Engineering | Chaos engineering principles: steady state, hypothesis, blast radius, and game days |
| `infra-chaos-spring` | Infrastructure: Chaos Spring | Chaos engineering with Spring: Chaos Monkey, fault injection, and resilience testing |
| `infra-devops-infra` | Infrastructure: DevOps | DevOps practices: Terraform, Ansible, GitOps workflows, and infrastructure as code |
| `infra-docker-containers` | Infrastructure: Docker | Docker: multi-stage builds, compose, security scanning, and image optimization |
| `infra-kubernetes` | Infrastructure: Kubernetes | Kubernetes: deployments, services, ingress, RBAC, Helm charts, and troubleshooting |
| `infra-opentelemetry` | Infrastructure: OpenTelemetry | OpenTelemetry: traces, metrics, logs, SDK configuration, and collector pipelines |
| `infra-traefik-proxy` | Infrastructure: Traefik Proxy | Traefik reverse proxy: dynamic routing, Let's Encrypt, middleware, and Docker provider |
| `infra-woodpecker-ci` | Infrastructure: Woodpecker CI | Woodpecker CI: pipeline syntax, plugins, secrets management, and matrix builds |

### Mobile Skills (2)

Mobile and hybrid app development patterns with Ionic and Capacitor.

| ID | Label | Description |
|----|-------|-------------|
| `mobile-ionic-capacitor` | Mobile: Ionic Capacitor | Capacitor native plugins: camera, filesystem, push notifications, and deep links |
| `mobile-mobile-ionic` | Mobile: Mobile Ionic | Ionic framework: components, navigation, theming, and platform-specific styling |

### Prompt & Quality Skills (2)

Prompt engineering and code quality review patterns.

| ID | Label | Description |
|----|-------|-------------|
| `prompt-improver` | Prompt: Prompt Improver | Prompt improvement techniques: structure, context, constraints, and few-shot examples |
| `quality-ghagga-review` | Quality: Ghagga Review | Code review following Ghagga methodology: incremental, checklist-based quality gates |

### References Skills (5)

Meta-references and templates for creating new hooks, skills, plugins, MCP servers, and sub-agents.

| ID | Label | Description |
|----|-------|-------------|
| `references-hooks-patterns` | References: Hooks Patterns | Reference patterns for creating new hooks: event types, validation, and side effects |
| `references-mcp-servers` | References: MCP Servers | Reference for configuring MCP servers: transport, tools, resources, and prompts |
| `references-plugins-reference` | References: Plugins Reference | Plugin development reference: lifecycle, configuration, and integration points |
| `references-skills-reference` | References: Skills Reference | Skill file specification: required sections, format, and best practices |
| `references-subagent-templates` | References: Subagent Templates | Templates for creating sub-agent skill files with delegation patterns |

### Systems & IoT Skills (4)

Low-level systems programming, embedded protocols, and async runtime patterns for Rust.

| ID | Label | Description |
|----|-------|-------------|
| `systems-modbus-protocol` | Systems: Modbus Protocol | Modbus TCP/RTU: register mapping, function codes, and industrial device communication |
| `systems-mqtt-rumqttc` | Systems: MQTT rumqttc | MQTT with rumqttc (Rust): pub/sub, QoS levels, retained messages, and TLS |
| `systems-rust-systems` | Systems: Rust Systems | Rust systems programming: FFI, unsafe blocks, memory layout, and no_std |
| `systems-tokio-async` | Systems: Tokio Async | Tokio async runtime: tasks, channels, select!, I/O operations, and graceful shutdown |

### Testing Skills (3)

Testing frameworks and strategies for frontend and backend applications.

| ID | Label | Description |
|----|-------|-------------|
| `testing-playwright-e2e` | Testing: Playwright E2E | Playwright E2E: Page Objects, auto-waiting, fixtures, and parallel execution |
| `testing-testcontainers` | Testing: Testcontainers | Testcontainers: Docker-based integration tests for databases, queues, and services |
| `testing-vitest-testing` | Testing: Vitest Testing | Vitest: fast unit tests, mocking, snapshot testing, and coverage with v8 |

### Workflow Skills (12)

Development workflow patterns for Git, CI, automation, and AI-assisted development.

| ID | Label | Description |
|----|-------|-------------|
| `workflow-ci-local-guide` | Workflow: CI Local Guide | Guide for running CI pipelines locally using `act` or `wrkflw` before pushing |
| `workflow-claude-automation` | Workflow: Claude Automation | Claude Code automation patterns: headless mode, batch processing, and scripting |
| `workflow-claude-md-improver` | Workflow: CLAUDE.md Improver | Best practices for writing and maintaining CLAUDE.md project instructions |
| `workflow-finish-dev-branch` | Workflow: Finish Dev Branch | Branch completion workflow: squash, rebase, cleanup, and merge strategies |
| `workflow-git-github` | Workflow: Git & GitHub | Git and GitHub workflows: branching strategies, PR conventions, and automation |
| `workflow-git-workflow` | Workflow: Git Workflow | Git workflow patterns: trunk-based, GitFlow, GitHub Flow, and release management |
| `workflow-ide-plugins` | Workflow: IDE Plugins | IDE plugin recommendations and configuration for VS Code and Neovim |
| `workflow-ide-plugins-intellij` | Workflow: IDE Plugins IntelliJ | IntelliJ IDEA plugins and configuration for optimal Java/Kotlin development |
| `workflow-obsidian-brain` | Workflow: Obsidian Brain | Obsidian knowledge management: daily notes, MOCs, templates, and AI integration |
| `workflow-git-worktrees` | Workflow: Git Worktrees | Git worktrees: parallel branch development, setup, and cleanup workflows |
| `workflow-verification` | Workflow: Verification | Post-implementation verification checklist: tests, linting, types, and manual review |
| `workflow-wave-workflow` | Workflow: Wave Workflow | Wave-based execution: dependency analysis, parallel waves, and progress tracking |

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

Model Context Protocol servers for enhanced AI capabilities. MCP is **atomic** — selecting any item enables the `mcp` feature, which installs ALL MCP server configurations. These servers extend your AI tool with external data sources and services.

| ID | Label | Description |
|----|-------|-------------|
| `mcp-context7` | Context7 | Remote MCP for fetching up-to-date library documentation |
| `mcp-engram` | Engram | Local MCP backend for persistent SDD artifacts and memory |
| `mcp-jira` | Jira | Jira integration via Atlassian MCP |
| `mcp-atlassian` | Atlassian | Confluence and Atlassian suite integration |
| `mcp-figma` | Figma | Figma design file integration |
| `mcp-notion` | Notion | Notion workspace integration |
