# ADR-016: Service Responsibility & Scope - `baseProject`

**Tags:** `service`, `responsibility`, `structure`

---

## Status

Accepted

---

## Context

`baseProject` creates the baseline repository structure and template output.
Without a clear boundary, base scaffolding concerns can spread into other services.

---

## Decision

Keep `base_project` focused on foundational project generation.

- Input: `BaseProjectConfig`.
- Work: render template paths and file contents.
- Filesystem boundary: all writes occur in scoped root handling.

---

## Advantages

- Clear ownership of initial scaffold generation.
- Lower coupling with lint/test/local/CI services.
- Predictable generation lifecycle in orchestrator ordering.

---

## Disadvantages

- Feature requests touching base structure often require this service to change.
- Additional output modes would require explicit extension work.

---

## Alternatives Considered

- **Fold base structure generation into orchestrator:** rejected because orchestration and generation concerns would mix.
- **Split into multiple micro-services immediately:** rejected because current scope is still cohesive within one service.
