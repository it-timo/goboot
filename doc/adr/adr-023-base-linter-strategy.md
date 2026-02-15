# ADR-023: Linter Configuration Rendering Strategy

**Tags:** `baselint`, `linting`, `templates`, `rendering`, `golangci`

---

## Status

Accepted

---

## Context

Lint configs require project-specific values and selected linter sets.
Static copy alone is insufficient for these generated artifacts.

---

## Decision

Render lint config files with Go `text/template` on a file-by-file basis.
Template data includes project metadata and enabled lint selections.

---

## Advantages

- Deterministic rendering with stdlib tooling.
- Template behavior aligns with other generator services.
- Easier integration testing of rendered output.

---

## Disadvantages

- Complex formatting helpers are limited compared with richer template engines.
- Template errors surface at render time and require good tests.

---

## Alternatives Considered

- **Mustache-like logic-less templates:** rejected because required computed values would shift complexity into pre-processing.
- **Sprig/Helm-style function sets:** rejected to avoid additional abstraction and non-stdlib dependency overhead.
