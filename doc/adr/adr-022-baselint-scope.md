# ADR-022: Dedicated Linting via `baseLint` Service

**Tags:** `service`, `linting`, `quality`, `separation-of-concerns`

---

## Status

Accepted

---

## Context

Linting setup can be mixed into unrelated generation logic, making ownership unclear.
`goboot` requires lint scaffolding while keeping service boundaries explicit.

---

## Decision

Use a dedicated `baseLint` service.

- Select enabled linters from config.
- Render lint config templates with project context.
- Optionally register local script commands through shared registrar interfaces.

---

## Advantages

- Lint concerns are isolated from base project and CI concerns.
- Lint defaults can evolve without changing unrelated services.
- Users can configure lint behavior through one service boundary.

---

## Disadvantages

- Adds another service and config surface area.
- Cross-service script integration requires orchestration coordination.

---

## Alternatives Considered

- **Put lint generation in `baseProject`:** rejected because it merges distinct responsibilities.
- **Ship static lint files only:** rejected because project-specific rendering requirements exist.
