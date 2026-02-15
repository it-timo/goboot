# ADR-031: Validate Generated Projects with Lint & Test Runs

**Tags:** `templates`, `quality`, `ci`, `generated-project`, `linting`, `testing`

---

## Status

Accepted

---

## Context

Generator unit tests do not fully prove generated project usability.
Template or wiring changes can break generated output even when generator tests pass.

---

## Decision

For changes affecting templates/service wiring/defaults:

1. generate a project from current config or fixture,
2. run generated project lint/test commands, and
3. treat failures as merge-blocking until fixed.

---

## Advantages

- Detects generated-output regressions earlier.
- Validates real usage paths beyond internal generator tests.
- Keeps templates and defaults runnable over time.

---

## Disadvantages

- Adds execution time to contributor workflows.
- Requires local environment parity for generated toolchain commands.

---

## Alternatives Considered

- **Generator tests only:** rejected because rendered-output integration issues can be missed.
- **CI-only generated-project checks:** rejected as sole control due to slower feedback loops.
