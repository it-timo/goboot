# ADR-029: Test Template Styles - Ginkgo by Default, Stdlib as Opt-Out

**Tags:** `testing`, `templates`, `ginkgo`, `stdlib`, `flexibility`

---

## Status

Accepted

---

## Context

Consumers have different test framework constraints.
The generator needs an explicit style choice without runtime ambiguity.

---

## Decision

Support `useStyle` in `BaseTestConfig` with:

- `ginkgo` (default), and
- `go` (stdlib style).

Templates include both variants; generation selects files based on `useStyle`.

---

## Advantages

- Preserves one default while supporting stdlib-only consumers.
- Keeps style selection explicit in config.
- Avoids style-specific branching spread across code paths.

---

## Disadvantages

- Maintaining multiple style templates increases upkeep.
- Shared defaults (such as test command strings) may not fit every team without override.

---

## Alternatives Considered

- **Ginkgo-only output:** rejected because some environments disallow third-party test frameworks.
- **Single hybrid template with many conditional branches:** rejected due to readability and maintenance costs.
