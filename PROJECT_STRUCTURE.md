# 📁 Project Structure — `goboot` (v0.0.2)

This document reflects the **current structure** of the `goboot` project as of version `v0.0.2`.

It is intentionally scoped to **what exists**, not what’s imagined.  
As new versions introduce layers (e.g., CI, Dockerization), this file will be updated accordingly.

For planned features, see [`ROADMAP.md`](./ROADMAP.md).

---

## ✅ Implemented Directories and Files

### `/cmd/`

- `cmd/goboot/main.go` — CLI entry point

### `/pkg/`

- `pkg/baseproject/` — Base project scaffolding service
- `pkg/baselint/` — Lint configuration service (dockerized linters)
- `pkg/baselocal/` — Local development scripts service
- `pkg/basetest/` — Testing scaffold service (Ginkgo/Gomega suites and helpers)
- `pkg/config/` — Config types and loading logic
- `pkg/goboot/` — Core execution engine
- `pkg/goboottypes/` — Shared constants and interfaces (service IDs, linter definitions, etc.)
- `pkg/gobootutils/` — Path/FS safety, template helpers, secure root handling

### `/configs/`

- `goboot.yml` — Main config entry point
- `base_project.yml` — Base project service config
- `base_lint.yml` — Lint service config (dockerized linters incl. shellcheck/shfmt)
- `base_local.yml` — Local scripts config
- `base_test.yml` — Test scaffold config

### `/templates/`

- `project_base/` — Project scaffolding templates
- `lint_base/` — Lint configuration templates (golangci-lint, yamllint, checkmake, markdownlint, shellcheck, shfmt)
- `local_base/` — Local helper scripts/templates
- `test_base/` — Testing templates (suite bootstrap, utils, sample specs)

### `/doc/adr/`

- ADRs (architecture decision records) for key technical choices  
  Example: config structure, service registry, no reflection, etc.

### `/doc/img/` and `/doc/diagram/`

- Visual documentation (Draw.io `.drawio` files and `.png` exports)

### `/.github/`

- GitHub `FUNDING.yml` file for sponsor links

### `/scripts/`

- Developer scripts (e.g., `lint`, `format`, `bootstrap`)

### Tests

- BDD test suites (Ginkgo/Gomega) co-located with packages, covering services, utilities, and secure FS handling
- Testing guide at [`TESTING.md`](./TESTING.md)

### Top-Level Files

- `README.md` — Project description and purpose
- `ROADMAP.md` — Versioned goals and features
- `VERSIONING.md` — Semantic version strategy
- `WORKFLOW.md` — Project lifecycle & contributor expectations
- `TESTING.md` — Testing philosophy, commands, and coverage notes
- `LICENSE`, `NOTICE` — Legal OSS declarations
- `.editorconfig`, `.gitignore`, `.gitattributes` — Development consistency
- `.env.example`, `.env.ci` — Placeholder environments
- `.nvmrc` — Tooling hints
- `go.mod`, `go.sum` — Go module metadata
- **`Makefile` — Common developer tasks**
- **`Taskfile.yml` — Task runner configuration**
- **`.golangci.yml` — Go linting configuration**
- **`.markdownlint.yaml` — Markdown linting configuration**
- **`.yamllint.yaml` — YAML linting configuration**
- **`.shellcheckrc` — Shell lint configuration**
- **`.pre-commit-config.yaml` — Optional pre-commit hooks metadata**
- **`.version` — Current project version**

---

## 🔜 Not Yet Present (Planned in Later Versions)

These directories are **not yet introduced** but are part of the intended long-term structure.  
See [`ROADMAP.md`](./ROADMAP.md) for targeted milestones.

- `test/` — Additional integration/e2e harnesses
- `benchmarks/` — Performance regression tracking
- CI workflows and contribution templates

---

## 🔄 Philosophy

This structure is:

- ✅ **Minimal by default**
- ✅ **Expanded only when needed**
- ✅ **Documented at every versioned step**

`goboot` aims to remain **predictable**, **clear**, and **scalable**,
without overwhelming new contributors or hiding logic behind automation.

---

_Last updated: v0.0.2 — matches real files in the repository._
