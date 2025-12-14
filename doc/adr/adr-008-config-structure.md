# 📄 ADR-008: Config System Structure and Philosophy

**Tags:** `config`, `modular-design`, `idiomatic-go`

---

## Status

✅ Accepted

---

## Context

Modern Go tools benefit from modular and declarative configuration, but many systems introduce runtime in-direction,
reflection, or plugin hooks that obscure behavior and reduce maintainability.

The `goboot` project needs a configuration system that is:

- Predictable
- Fully typed
- Declarative but explicit
- Idiomatic and traceable

This applies both to how configs are defined and how they are loaded, validated, and used.

---

## Decision

Implement a centralized `config` package that provides:

- A `ServiceConfig` interface for static config modules (e.g., base project, linting, Docker)
- A `Manager` that registers, validates, and exposes these modules
- A `GoBoot` type that acts as entry point and orchestrator

---

## Advantages

- Fully typed: All configs are concrete Go structs, no dynamic maps or generics
- Centralized: Validation, loading, and usage happen in well-defined places
- Predictable: No runtime injection or service hooks
- Scalable: New config modules can be added without changing the core behavior

---

## Disadvantages

- All config modules must be hardcoded in `createServiceConfig()`
- Cannot dynamically register config types at runtime
