# 📄 ADR-018: Static Service Registration and Orchestration

**Tags:** `services`, `registration`, `explicit-architecture`

---

## Status

✅ Accepted

---

## Context

Dynamic service loading via reflection, plugins, or dependency injection frameworks introduces complexity,
weakens compile-time guarantees, and complicates extensibility. `goboot` favors static,
explicit definitions to align with idiomatic Go and clarity in OSS.

---

## Decision

- Each service (e.g., `base_project`, `base_lint`) is registered statically
  in a centralized registry function (e.g., `RegisterServices()`).
- Services implement a shared `Service` interface and are matched by declared ID (e.g., `base_lint`) in `types`.
- No dynamic discovery or auto-wiring is used.

---

## Advantages

- Compile-time safety, no runtime surprises
- Easy grepping and traceability
- Controlled and predictable service execution
- Encourages clear responsibility boundaries per service

---

## Disadvantages

- Slightly more manual for new service integration
- Less "magical" extensibility (but this is intentional)

---

## Alternatives Considered

- DI/Reflection: rejected due to runtime risk and reduced clarity
