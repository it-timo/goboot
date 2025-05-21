# 📄 ADR-003: Template Engine, Structure, and Naming Rules

**Tags:** `templates`, `text/template`, `scaffolding`, `structure`

---

## Status

✅ Accepted

---

## Context

The `goboot` tool must generate a complete, idiomatic Go project by:

- Creating nested folders (e.g. `cmd/{{.LowerProjectName}}/main.go`)
- Injecting structured content into files (e.g. `README.md`, `LICENSE`)
- Adapting naming conventions consistently across files

This requires a predictable, traceable, and testable templating strategy for both **file paths** and **file contents**.

---

## Decision

Use Go’s standard `text/template` engine to:

- Render all **file and directory paths** from template names
- Render all **file contents** with the same templating engine
- Inject values from the parsed `BaseProjectConfig`

All template files live under `templates/project_base/` and are rendered into a secure `*os.Root` under the target directory.

Template placeholders include:

| Placeholder             | Description                               |
|-------------------------|-------------------------------------------|
| `{{.ProjectName}}`      | Original project name                     |
| `{{.LowerProjectName}}` | Lowercase-safe project name               |
| `{{.CapsProjectName}}`  | Uppercase variant (e.g., for headers)     |
| `{{.UsedGoVersion}}`    | Go version to inject into `go.mod`, etc.  |
| `{{.Author}}`           | For LICENSE, NOTICE, README               |
| ...                     | All other fields from `BaseProjectConfig` |

### Structural Flow

1. Walk all source files under `templates/project_base/` using `fs.WalkDir`
2. Render paths into target dir via `renderPath(...)`
3. Copy raw files into `*os.Root`
4. Second pass: apply `renderContent(...)` using `text/template`

---

## Advantages

- ✅ No external dependency — stdlib-based
- ✅ Debuggable and inspectable by any Go developer
- ✅ Consistent rendering engine for structure and logic
- ✅ Easier to add new fields by expanding config

---

## Disadvantages

- ❗ Template logic must be straightforward — no `.tmpl` extensions or third-party functions
- ❗ Errors in template variables cause runtime errors (e.g., missing `LowerProjectName`)

---

## Alternatives Considered

- Using `.tmpl` extensions and parsing logic based on extension: rejected for simplicity.
- Using `sprig` (Helm-like templating functions): rejected for now; might revisit if advanced use-cases emerge.

---
