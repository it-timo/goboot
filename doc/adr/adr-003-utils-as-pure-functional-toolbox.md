# ADR-003: Intentional Use of `pkg/gobootutils` as Pure Functional Set

**Tags:** `utils`, `hygiene`, `modularity`, `stateless`

---

## Status

Accepted

---

## Context

Utility packages commonly accumulate mixed responsibilities.
That pattern increases ownership ambiguity and can introduce hidden dependencies.

---

## Decision

Keep `pkg/gobootutils` restricted to stateless helpers.

- No logging side effects.
- No config/state ownership.
- No global mutable variables.
- No domain orchestration logic.

---

## Advantages

- Clear package boundary for reusable helpers.
- Reduced risk of import cycles.
- Easier unit testing of helper behavior.

---

## Disadvantages

- Enforcing boundaries requires review discipline.
- Some helpers may need relocation as domain boundaries evolve.

---

## Alternatives Considered

- **Broad catch-all utils package:** rejected due to long-term coupling and discoverability issues.
- **No shared utils at all:** rejected because repeated low-level code would grow across services.
