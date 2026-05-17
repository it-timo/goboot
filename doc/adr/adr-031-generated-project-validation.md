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

The generated test files are scaffold smoke tests. Their job is to prove the
rendered repository compiles and that generated Make, Task, and script entry
points remain executable. They are not expected to assert product/domain
behavior that only the downstream application can define.

Running multiple wrappers over the same underlying lint/test commands is
intentional parity evidence, not duplicate release evidence. A generated
repository should be usable through whichever local entry point a team chooses.

---

## Advantages

- Detects generated-output regressions earlier.
- Validates real usage paths beyond internal generator tests.
- Validates wrapper parity for Make, Task, and scripts.
- Keeps templates and defaults runnable over time.

---

## Disadvantages

- Adds execution time to contributor workflows.
- Requires local environment parity for generated toolchain commands.

---

## Alternatives Considered

- **Generator tests only:** rejected because rendered-output integration issues can be missed.
- **CI-only generated-project checks:** rejected as sole control due to slower feedback loops.
