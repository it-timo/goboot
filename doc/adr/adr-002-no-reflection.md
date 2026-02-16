# ADR-002: Explicit Separation of Concerns - No Runtime Reflection or DI

**Tags:** `philosophy`, `idioms`, `anti-patterns`

---

## Status

Accepted

---

## Context

Runtime reflection and DI frameworks can reduce explicitness in control flow.
For a scaffolding tool, deterministic behavior and easy inspection are prioritized over dynamic wiring.

---

## Decision

Do not use reflection-based discovery, runtime plugin loading, or container-style dependency injection.

Use:

- explicit struct construction,
- explicit registration, and
- explicit orchestration in the top-level service manager.

---

## Advantages

- Execution path is visible in code.
- Fewer runtime failure modes from misconfigured containers/registries.
- Easier static review and traceability.

---

## Disadvantages

- Adding new modules requires explicit wiring updates.
- Third-party runtime extension points are intentionally limited.

---

## Alternatives Considered

- **Reflection-based service discovery:** rejected due to weaker compile-time guarantees and harder debugging.
- **DI container frameworks:** rejected due to additional runtime complexity without clear payoff for current scope.
