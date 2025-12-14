# 📄 ADR-007: Centralized Service Name Registry in `types/names.go`

**Tags:** `constants`, `service-names`, `structure`, `decoupling`

---

## Status

✅ Accepted

---

### Context

Service identifiers such as `"base_project"` are used across config declarations,
runtime registration, and execution logic.
Allowing these identifiers to be spread across the codebase as raw strings introduces coupling, typos,
and inconsistency.

---

## Decision

Define **all service names as constants** in a centralized file:
`pkg/types/names.go`

This registry provides a **single source of truth** for cross-cutting service identifiers.

---

## Advantages

- Prevents config ↔ logic mismatch due to typos
- Enables future tooling, e.g., service listing or generation
- Promotes discoverability across the codebase
- Safe for refactoring and reuse

---

## Disadvantages

- Introduces indirect coupling between services and `types`
- Slight overhead from needing to import a common constants file
