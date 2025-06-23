# 📄 ADR-003: Intentional Use of `pkg/utils` as Pure Functional Set

**Tags:** `utils`, `hygiene`, `modularity`, `stateless`

---

## Status

✅ Accepted

---

## Context

Utility packages are frequently abused as dumping grounds for unrelated helpers, side-effect-laden code,
or improperly scoped functions.
This undermines testability, introduces circular imports, and confuses ownership.

---

## Decision

Restrict `pkg/utils` to **pure, stateless helper functions** with:

- No logging
- No config access
- No global variables or shared state

Only generic helpers (e.g., file-safe mkdir, path cleaning) are permitted.

---

## Advantages

- Prevents misuse and “god-package” growth
- Enables safe, dependency-free reuse
- Avoid import loops and testing side effects
- Clarifies that `utils` is **not** a place for business logic

---

## Disadvantages

- Requires active enforcement or review discipline
- May initially confuse new contributors (“where should I put this?”)
