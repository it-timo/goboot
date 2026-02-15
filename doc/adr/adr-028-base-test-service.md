# ADR-028: Dedicated Test Scaffolding via `baseTest` Service

**Tags:** `service`, `testing`, `templates`, `scripts`, `scaffolding`

---

## Status

Accepted

---

## Context

Generated projects should include runnable test scaffolding without folding test concerns into unrelated services.

---

## Decision

Introduce `baseTest` as a dedicated service for test scaffolding.

- Use typed `BaseTestConfig` as service input.
- Render test templates with a two-pass path/content flow.
- Keep writes root-scoped.
- Optionally register test command lines through shared registrar integration.

---

## Advantages

- Test scaffold generation is modular and configurable.
- Template approach supports style variants without heavy runtime branching.
- Integrates with local scripts while preserving service boundaries.

---

## Disadvantages

- Adds service and config maintenance overhead.
- Default style choices can pull in optional external test dependencies.

---

## Alternatives Considered

- **Embed test files in `baseProject`:** rejected due to boundary expansion in base scaffolding service.
- **No generated tests:** rejected because users would repeatedly recreate baseline test setup manually.
