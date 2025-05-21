## 📄 ADR-012: Template Rendering Strategy (Paths and Content)

**Tags:** `templates`, `rendering`, `text/template`

## Status
✅ Accepted

### Context

Both file paths (e.g., `cmd/{{.LowerProjectName}}`) and file contents (e.g., `README.md`) contain placeholders 
that must be rendered from configuration.

### Decision

Use Go's standard `text/template` engine to process:

* File paths during structure creation
* File contents in a second pass

### Advantages

* Familiar, built-in templating engine — no dependency overhead
* Supports logic (`if`, `eq`, etc.) directly in templates
* Rendering logic is debuggable and transparent

### Disadvantages

* More complex templates (e.g., nested loops or data access) may be harder to express
* No out-of-the-box support for `.tmpl` conditionals like other engines
