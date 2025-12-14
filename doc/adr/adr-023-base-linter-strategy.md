# 📄 ADR-023: Linter Configuration Rendering Strategy

**Tags:** `baselint`, `linting`, `templates`, `rendering`, `golangci`

---

## Status

✅ Accepted

---

## Context

Linting configuration files (like `.golangci.yml`, `.markdownlint.yml`, etc.) need project-specific metadata injected
(e.g., module name, linter sets). Instead of copying static files, these templates must be rendered with context.

The `baseLint` service is responsible for generating these configurations dynamically,
ensuring they are tailored to each project.

---

## Decision

The `baseLint` service uses **Go’s built-in `text/template`**
to render linter configuration files from source templates, injecting:

- Project name or module path
- Enabled linters (from config)
- Optional service-related metadata

Rendering is strictly **file-by-file**, avoiding deep templating logic or runtime dependencies.

---

## Advantages

- Predictable and audit-friendly template system.
- Enables consistent config output across projects.
- No external dependencies or language extensions.
- Easier to test and reason about.

---

## Disadvantages

- While text/template supports range, if, and method calls, it lacks built-in helpers for different logics
  (e.g., trimming, joins, case conversion) — which must be handled in Go.
- Template debugging can be less ergonomic than richer engines.

---

## Alternatives Considered

- **Mustache:** Lacks logic, requiring pre-computed template input structs. Not a net win.
- **Sprig with Helm-style templates:** Powerful but introduces YAML logic bleed, extra cognitive load,
  and non-standard Go behavior.
