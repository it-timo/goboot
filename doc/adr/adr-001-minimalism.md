# ADR-001: Minimalism over Generalization

**Tags:** `design-philosophy`, `minimalism`, `abstraction`

---

## Status

Accepted

---

## Context

Early abstractions (interfaces, reflection, plugin maps) can make a small codebase harder to follow.
At this stage, `goboot` prioritizes predictable control flow and low onboarding overhead.

---

## Decision

Prefer concrete types and linear execution paths.

- Use explicit structs and direct function calls.
- Avoid plugin maps and injected handler chains.
- Keep helper usage limited to repeated, generic operations.
- Introduce interfaces only when at least two concrete implementations are required by current behavior.

---

## Advantages

- Lower indirection in core flows.
- Easier debugging and tracing during changes.
- Fewer abstraction layers to keep consistent.

---

## Disadvantages

- Refactoring cost increases if runtime-extensible behavior is introduced later.
- Some duplicated patterns may remain until a real abstraction need appears.

---

## Alternatives Considered

- **Abstract-first architecture:** rejected because current scope does not justify the added indirection.
- **Plugin-oriented extension model from day one:** rejected because extension is not a current product goal.
