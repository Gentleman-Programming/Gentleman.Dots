---
name: technical-review
description: >
  Review technical exercises and candidate submissions with structured evaluation.
  Trigger: When reviewing technical exercises, code assessments, candidate submissions, or take-home tests.
license: MIT
metadata:
  author: gentleman-programming
  version: "1.0"
---

## When to Use

- Reviewing take-home technical exercises from candidates
- Evaluating code submissions for hiring decisions
- Assessing technical tests before interviews

## Review Process

1. **Explore structure first** - Use Task tool with explore agent to understand project layout
2. **Read key files in parallel** - models, views, serializers, tasks, tests, README, docker-compose
3. **Check for tests** - Presence/absence of tests is a major signal for senior roles
4. **Look for red flags** - Security issues, leaked corporate data, no error handling
5. **Score each factor 0-10** with specific evidence from code
6. **Output as Markdown table** per candidate

## Evaluation Factors (Always These 6)

| Factor | What to Look For | Red Flags |
|--------|------------------|-----------|
| **Styling** | Consistent formatting, naming conventions, file organization, language idioms | Mixed styles, inconsistent naming, messy structure |
| **Technical expertise** | Correct use of primitives, sensible architecture, tradeoff awareness | Cargo-cult patterns, overengineering, security gaps |
| **Code Quality** | Maintainability, testability, small functions, error handling, tests | Giant functions, tight coupling, no tests |
| **Go beyond what was asked** | Useful additions: tests, validation, docs, UX improvements | Scope creep, unrelated features |
| **Detailed explanations** | README quality, design rationale, setup instructions | No context, missing setup instructions |
| **Other comments or notes** | Security, operational concerns, collaboration signals | Leaked corporate data, hardcoded secrets |

## Red Flags Checklist

- [ ] Secrets/API keys in code or config
- [ ] Employer data exposed (AWS accounts, internal URLs, company names)
- [ ] No tests at all (critical for senior roles)
- [ ] Copy-pasted code without understanding
- [ ] Missing README or setup instructions
- [ ] Security gaps (SQL injection, no input validation, unsafe defaults)
- [ ] Over-engineering without justification
- [ ] Giant functions (>50 lines)

## Output Format (Markdown Table)

For each candidate, output a table with these columns:

```markdown
# Review Candidato: {Name}

| Factor | Guidance | Scoring (0-10) | Notes |
|--------|----------|----------------|-------|
| **Styling** | {guidance text} | {score} | {specific observations} |
| **Technical expertise** | {guidance text} | {score} | {specific observations} |
| **Code Quality** | {guidance text} | {score} | {specific observations} |
| **Go beyond what was asked** | {guidance text} | {score} | {specific observations} |
| **Detailed explanations** | {guidance text} | {score} | {specific observations} |
| **Other comments or notes** | {guidance text} | | {Strengths / Concerns / Questions for follow-up} |
```

End with comparative summary if multiple candidates.

## Guidance Text for Each Factor

### Styling
What to look for: Consistent formatting, naming, file organization, and adherence to the project's conventions. Good signals: following language naming schema, predictable folder structure, lint/formatter passes, readable layout, no "random" naming. Red flags: Mixed styles in the same codebase, inconsistent naming, messy diffs, unclear structure, ignores existing conventions. How to score: Judge consistency and readability, not personal preference.

### Technical expertise
What to look for: Correct use of language/framework primitives, sensible architecture, and awareness of tradeoffs. Good signals: Chooses appropriate data structures, avoids unnecessary complexity, uses correct async/concurrency patterns, handles edge cases, understands performance/security implications. Red flags: Cargo-cult patterns, misuse of framework/libraries, overengineering without justification, obvious security gaps, no awareness of complexity. How to score: Reward correctness + good judgment. Ask: "Would I trust this person to work without constant supervision?"

### Code Quality
What to look for: Maintainability, testability, clarity, and correctness. Good signals: Small focused functions, clear interfaces, minimal duplication, good error handling, meaningful tests, good separation of concerns. Red flags: Giant functions, tight coupling, magic numbers, fragile logic, unclear responsibilities, no tests where they matter. How to score: Ask: "Can someone else extend this in 3 months without rewriting it?"

### Go beyond what was asked
What to look for: Useful additions that improve the solution without derailing scope. Good signals: Adds tests, better error messages, sensible validation, small UX improvements, optional enhancements behind flags, documentation for setup/run. Red flags: Adds unrelated features, expands scope massively, introduces new dependencies "just because", changes requirements without agreement. How to score: Reward high leverage extras. Penalize scope creep.

### Detailed explanations
What to look for: Can they communicate intent, rationale, and tradeoffs. Good signals: Clear README or PR description, explains design choices, calls out assumptions/limitations, describes how to run/test, notes future improvements. Red flags: "It works" with no context, no explanation of non-obvious decisions, missing setup instructions, can't justify tradeoffs. How to score: Prefer concise clarity over essay-length. If it's complex, explanations should match complexity.

### Other comments or notes
What to include: Anything that affects hiring confidence but doesn't fit above. Examples: Risk areas (security concerns, input validation gaps, missing auth, unsafe defaults), Operational concerns (logging, observability, configuration handling), Collaboration signals (commit hygiene, responsiveness), Constraints (time spent, what they intentionally didn't do, known limitations). Useful format: "Strengths / Concerns / Questions to ask in follow-up interview".

## Commands

```bash
# Explore project structure
eza -la --tree --level=3

# Find key files
fd -e py -e md -e yml -e yaml

# Check for tests
fd test

# Look for secrets (red flag check)
rg -i "secret|password|api_key|token" --type py

# Check for hardcoded AWS accounts or corporate data
rg -i "arn:aws|account.*[0-9]{12}" --type py
```

## Comparative Summary Template

```markdown
## Resumen Comparativo

| Factor | {Candidate 1} | {Candidate 2} |
|--------|---------------|---------------|
| Styling | {score} | {score} |
| Technical expertise | {score} | {score} |
| Code Quality | {score} | {score} |
| Go beyond what was asked | {score} | {score} |
| Detailed explanations | {score} | {score} |
| **TOTAL** | **{total}** | **{total}** |

**Recomendación:** {Clear recommendation with reasoning}
```
