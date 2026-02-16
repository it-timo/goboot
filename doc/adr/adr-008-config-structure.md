# ADR-008: Config System Structure and Philosophy

**Tags:** `config`, `modular-design`, `idiomatic-go`

---

## Status

Accepted

---

## Context

Configuration must remain typed and predictable as services increase.
Dynamic map-based config handling makes validation and refactoring harder.

---

## Decision

Use a centralized `config` package with typed service config modules and a manager for load/validate/access flows.

- Define service-specific config structs.
- Validate during load before service execution.
- Keep orchestration explicit in `GoBoot`.

---

## Advantages

- Typed config boundaries for each service.
- Centralized validation behavior.
- Predictable startup and config lifecycle.

---

## Disadvantages

- New config types require explicit registration updates.
- Boilerplate cost is higher than dynamic decoding.

---

## Alternatives Considered

- **Map-based/untyped config model:** rejected due to weaker safety and tooling support.
- **Runtime plugin config registration:** rejected due to reduced traceability and higher runtime complexity.
