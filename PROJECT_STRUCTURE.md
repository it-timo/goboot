# 📁 Project Structure — `goboot` (v0.0.1)

This document reflects the **current structure** of the `goboot` project as of version `v0.0.1`.

It is intentionally scoped to **what exists**, not what’s imagined.  
As new versions introduce layers (e.g., CI, tests, Docker), this file will be updated accordingly.

For planned features, see [`ROADMAP.md`](./ROADMAP.md).

---

## ✅ Implemented Directories and Files

### `/cmd/`

- `cmd/goboot/main.go` — CLI entry point

### `/pkg/`

- `pkg/config/` — Config types and loading logic
- `pkg/goboot/` — Core execution engine
- `pkg/baseProject/` — The first built-in service
- `pkg/types/` — Shared constants (e.g., service identifiers)
- `pkg/utils/` — General-purpose helpers

### `/configs/`

- `goboot.yml` — Main config entry point
- `base_project.yml` — Service-specific config

### `/templates/project_base/`

- Template tree for new project scaffolding
- Contains real templates like `README.md`, `LICENSE`, `cmd/{{.LowerProjectName}}`, etc.

### `/doc/adr/`

- ADRs (architecture decision records) for key technical choices  
  Example: config structure, service registry, no reflection, etc.

### `/doc/img/` and `/doc/diagram/`

- Visual documentation (Draw.io `.drawio` files and `.png` exports)

### `/.github/`

- GitHub `FUNDING.yml` file for sponsor links

### `/scripts/`

- Developer scripts (e.g., `lint`, `format`, `bootstrap`)

### Top-Level Files

- `README.md` — Project description and purpose
- `ROADMAP.md` — Versioned goals and features
- `VERSIONING.md` — Semantic version strategy
- `WORKFLOW.md` — Project lifecycle & contributor expectations
- `LICENSE`, `NOTICE` — Legal OSS declarations
- `.editorconfig`, `.gitignore`, `.gitattributes` — Development consistency
- `.env.example`, `.env.ci` — Placeholder environments
- `.nvmrc`, `.tool-versions` — Tooling hints
- `go.mod`, `go.sum` — Go module metadata
- **`Makefile` — Common developer tasks**
- **`Taskfile.yml` — Task runner configuration**
- **`.golangci.yml` — Go linting configuration**
- **`.markdownlint.yaml` — Markdown linting configuration**
- **`.yamllint.yaml` — YAML linting configuration**

---

## 🔜 Not Yet Present (Planned in Later Versions)

These directories are **not yet introduced** but are part of the intended long-term structure.  
See [`ROADMAP.md`](./ROADMAP.md) for targeted milestones.

- `test/` — Test structure (unit/integration harnesses)
- `benchmarks/` — Performance regression tracking
- `.github/` / `.gitlab/` — CI workflows, issue templates, etc.

---

## 🔄 Philosophy

This structure is:

- ✅ **Minimal by default**
- ✅ **Expanded only when needed**
- ✅ **Documented at every versioned step**

`goboot` aims to remain **predictable**, **clear**, and **scalable**,
without overwhelming new contributors or hiding logic behind automation.

---

_Last updated: v0.0.1 — matches real files in the repository._
