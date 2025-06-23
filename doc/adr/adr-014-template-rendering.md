# ADR-014: Template Rendering Strategy using `text/template`

**Tags:** `templates`, `text/template`, `scaffolding`

---

## Status

✅ Accepted

---

## Context

Many templating engines exist (`mustache`, `sprig`, `goja`), but for Go projects focused on filesystem scaffolding
and config rendering, Go's `text/template` is sufficient, robust, and idiomatic.

---

## Decision

- All templates (files, paths, scripts) are rendered using Go's standard `text/template` engine.
- No custom functions or external engines are introduced.
- The rendering context is always a strongly typed struct (e.g., `BaseLintConfig`, `scriptRegistry`).

---

## Advantages

- Minimal dependencies
- Fully Go-native and safe
- Predictable output, no hidden magic

---

## Disadvantages

- Less expressive than engines like `sprig`
- Requires more verbose templates in complex cases

---

## Alternatives Considered

- `sprig` for extra helpers — rejected to avoid dependencies
- `mustache` or `handlebars` — non-native and overkill for needs
