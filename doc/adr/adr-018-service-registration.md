# ADR-018: Static Service Registration and Orchestration

**Tags:** `services`, `registration`, `explicit-architecture`

---

## Status

Accepted

---

## Context

Service registration determines what can run in a generation flow.
Dynamic registration mechanisms can obscure behavior and produce environment-dependent startup paths.

---

## Decision

Register services statically in centralized orchestration code.

- Services implement the shared `Service` interface.
- Registry wiring is explicit.
- No reflection-based auto-registration.

---

## Advantages

- Startup behavior is reviewable in one place.
- Compile-time references support easier refactoring.
- Lower variance between environments.

---

## Disadvantages

- Registration list must be maintained manually.
- Third-party extension requires source-level integration.

---

## Alternatives Considered

- **Reflection/scan-based registration:** rejected due to weaker traceability.
- **Init-time self-registration maps:** rejected because implicit side effects make startup harder to reason about.
