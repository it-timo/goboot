# ADR-024: Decoupled Script Coordination via `Registrar` Interface

**Tags:** `baseLocal`, `scripts`, `interface`, `coordination`, `registrar`, `extensibility`, `separation-of-concerns`

---

## Status

Accepted

---

## Context

Multiple services contribute commands that end up in shared local developer assets
(such as `Makefile`, `Taskfile.yml`, and shell scripts).
Direct file ownership by each service would create write conflicts and formatting drift.

---

## Decision

Use a `Registrar` contract implemented by the local-output service.
Contributing services receive the hook during orchestration and register commands/files, while `baseLocal` owns rendering.

---

## Advantages

- Separates command declaration from file rendering.
- Keeps shared output formatting in one place.
- Supports adding contributors without exposing `baseLocal` internals.

---

## Disadvantages

- Registration order and naming conflicts must be handled carefully.
- Contributors need to understand orchestration lifecycle hooks.

---

## Alternatives Considered

- **Each service writes target files directly:** rejected due to conflict risk and duplicated formatting logic.
- **Single global mutation helper with no typed contract:** rejected because validation and testability would be weaker.
