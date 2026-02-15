# ADR-021: Extensibility Strategy for New Services and Features

**Tags:** `extensibility`, `oss`, `architecture`

---

## Status

Accepted

---

## Context

The project needs a repeatable way to add services and config fields without destabilizing existing behavior.

---

## Decision

Define extension points through existing architecture contracts.

- New services implement the shared service interface and register explicitly.
- New config fields are typed and validated in config modules.
- Templates remain data-driven through `text/template` in dedicated directories.

---

## Advantages

- Extension path is explicit and consistent.
- Risk of unreviewed runtime behavior is reduced.
- New features align with existing test and config workflows.

---

## Disadvantages

- Adding capabilities requires updates in multiple explicit places.
- Fast experimentation is slower than dynamic plugin-style approaches.

---

## Alternatives Considered

- **Dynamic registration extension points:** rejected due to traceability and validation concerns.
- **Feature-specific one-off patterns:** rejected because architectural drift would increase over time.
