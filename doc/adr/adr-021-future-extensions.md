# 📄 ADR-021: Extensibility Strategy for New Services and Features

**Tags:** `extensibility`, `oss`, `architecture`

---

## Status

✅ Accepted

---

## Context

As an OSS project, `goboot` must enable contributors to add new services, templates,
or config fields with minimal risk of breaking existing behavior.

---

## Decision

- All services must:
  - Implement a shared `Service` interface
  - Register via `RegisterServices()` based on their declared ID
  - Define a distinct config struct (`BaseXConfig`)

- All config options must:
  - Be added via the `config.Manager`
  - Pass validation before execution

- All templates must:
  - Be rendered via `text/template`
  - Live under dedicated, discoverable directories

---

## Advantages

- Encourages community contribution
- Low risk of regressions
- Predictable points of integration

---

## Disadvantages

- Slight manual effort to extend registry
- Requires discipline across contributors

---

## Alternatives Considered

- Dynamic service registration — rejected due to validation and traceability concerns
- Global service loader maps — rejected due to tight coupling and testability loss
