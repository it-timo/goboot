# 📄 ADR-022: Dedicated Linting via `baseLint` Service

**Tags:** `service`, `linting`, `quality`, `separation-of-concerns`

---

## Status

✅ Accepted

---

## Context

Linting is a critical part of long-term maintainability and code quality. However, many project generators either:

- omit it entirely,
- bake in linter files directly,
- or mix linting setup with unrelated logic like testing or CI bootstrapping.

In `goboot`, linting is treated as a **first-class modular service**, provided through `baseLint`.

---

## Decision

Introduce a dedicated `baseLint` service that:

- Selectively enables predefined linters (Go, YAML, Markdown, etc.) from configuration.
- Generates config files with project-specific rendering via `text/template`.
- Delegates optional script registrations (e.g. `golangci-lint run`) to the `baseLocal` system, if enabled.
- Performs all logic in a standalone, clearly scoped unit (`pkg/baselint`).

---

## Advantages

- Ensures every generated project starts with consistent linting standards.
- Makes it trivial to evolve linter strategies over time without polluting unrelated services.
- Avoids premature assumptions: no CI coupling, no automatic script wiring unless explicitly opted in.
- Enables future expansion: custom linters, template variants, style presets.

---

## Disadvantages

- Slight increase in complexity.
- Requires a second level of awareness from users to enable/disable linters.

---

## Alternatives Considered

- **Bundling linting into `baseProject`:** Would lead to unclear service scope and potential config bloat.
- **Baking static files without templating:** Would limit reusability and require duplication across templates.
- **Relying on external plugins/hooks:** Would violate the goal of having reliable, predictable OSS bootstrapping.
