# 📄 ADR-002: Explicit Separation of Concerns — No Runtime Reflection or DI

**Tags:** `philosophy`, `idioms`, `anti-patterns`

---

## Status

✅ Accepted

---

## Context

Modern Go OSS often overuses runtime abstractions such as:

- Global registries
- Interface-based plugin discovery
- Reflection or generic container types

These increase complexity without adding tangible value, especially in tools like `goboot` which aim for deterministic,
high-trust generation logic.

---

## Decision

Reject all the following:

- No global plugin systems
- No reflection-based discovery
- No runtime service loading
- No interface-driven dependency injection

Instead:

- Explicit structs
- Explicit registration
- Top-level service mapping

---

## Advantages

- Enforces clarity and predictability
- Encourages meaningful, localized logic
- Easier for contributors to follow and extend

---

## Disadvantages

- Not extensible via external plugin systems (intended limitation)
- Slightly more maintenance effort to onboard new modules
