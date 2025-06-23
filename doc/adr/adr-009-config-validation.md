# 📄 ADR-009: Typed Config Structure and Validation Strategy

**Tags:** `config`, `validation`, `typed-structure`

---

## Status

✅ Accepted

---

## Context

Handling config via `map[string]interface{}` or dynamic YAML decoding leads to fragile code and poor IDE support.
`goboot` adopts strict, typed configuration via Go structs with validation logic.

---

## Decision

- Every service defines a dedicated config struct (e.g., `BaseLintConfig`)
- A central config manager validates configs at load time
- Only validated `ServiceConfig` types are passed to services

---

## Advantages

- Prevents runtime panics from missing fields
- IDE auto-completion and refactor support
- Easier testing, documentation, and migration

---

## Disadvantages

- Slightly more boilerplate per config
- Changes require struct updates and revalidation logic

---

## Alternatives Considered

- Unstructured map config — rejected due to brittleness and low maintainability
