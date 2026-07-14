# 📁 Project Structure — `goboot` (v0.3.0)

This document reflects the current repository structure for `goboot` at `v0.3.0`.
It documents what exists now; planned additions live in [`ROADMAP.md`](../ROADMAP.md).

## ✅ Implemented Directories and Files

### `/cmd/`

- `cmd/goboot/main.go` — CLI entry point

### `/pkg/`

- `pkg/baseproject/` — Base project scaffolding service
- `pkg/baselint/` — Lint configuration service (dockerized linters)
- `pkg/baselocal/` — Local development scripts service
- `pkg/basetest/` — Testing scaffold service (Ginkgo/Gomega suites and helpers)
- `pkg/baseci/` — CI scaffolding service (GitLab/GitHub generation with policy modes)
- `pkg/basedocker/` — CLI containerization service (Dockerfile, compose, dockerignore)
- `pkg/baserelease/` — GoReleaser and provider release automation service
- `pkg/config/` — Config types and loading logic
- `pkg/goboot/` — Core execution engine
- `pkg/goboottypes/` — Shared constants and interfaces (service IDs, linter definitions, etc.)
- `pkg/gobootutils/` — Path/FS safety, template helpers, secure root handling

### `/configs/`

- `goboot.yml` — Main config entry point
- `base_project.yml` — Base project service config
- `base_lint.yml` — Lint service config (dockerized linters incl. shellcheck/shfmt/editorconfig-checker)
- `base_local.yml` — Local scripts config
- `base_test.yml` — Test scaffold config
- `base_logger.yml` — Logger scaffold settings config
- `base_ci.yml` — CI scaffold config
- `base_docker.yml` — Docker scaffold config
- `base_release.yml` — Release automation config

### `/templates/`

- `project_base/` — Project scaffolding templates
- `lint_base/` — Lint configuration templates (golangci-lint, yamllint, checkmake, markdownlint, shellcheck, shfmt, editorconfig-checker)
- `local_base/` — Local helper scripts/templates
- `test_base/` — Testing templates (suite bootstrap, utils, sample specs)
- `ci_base/` — CI templates (GitLab and GitHub providers)
- `docker_base/` — Dockerfile, compose, and dockerignore templates
- `release_base/` — GoReleaser and release guide templates

### `/doc/adr/`

- ADRs (architecture decision records) for key technical choices  
  Example: config structure, service registry, no reflection, etc.

### `/doc/img/` and `/doc/diagram/`

- Visual documentation (Draw.io `.drawio` files and `.png` exports)

### `/.github/`

- GitHub workflow files and optional sponsor metadata
- Container workflow for building and smoke-testing the `goboot` image
- Release workflow for tag-driven binary distribution

### `/scripts/`

- Developer and verification scripts (`lint.sh`, `test.sh`, `verify_introproject.sh`,
  `verify_ci_canary.sh`)

### Tests

- BDD test suites (Ginkgo/Gomega) co-located with packages, covering services, utilities, and secure FS handling
- Testing guide at [`doc/TESTING.md`](./TESTING.md)

### Top-Level Files

- `README.md` — Project description and purpose
- `ROADMAP.md` — Versioned goals and features
- `doc/VERSIONING.md` — Semantic version strategy
- `doc/WORKFLOW.md` — Project lifecycle & contributor expectations
- `doc/TESTING.md` — Testing philosophy, commands, and coverage notes
- `LICENSE`, `NOTICE` — Legal OSS declarations
- `.editorconfig`, `.gitignore`, `.gitattributes` — Development consistency
- `Dockerfile`, `.dockerignore` — Container image definition for the `goboot` CLI
- `.goreleaser.yml` — Cross-platform binary release definition
- `.nvmrc` — Tooling hints
- `go.mod`, `go.sum` — Go module metadata
- **`Makefile` — Common developer tasks**
- **`Taskfile.yml` — Task runner configuration**
- **`.golangci.yml` — Go linting configuration**
- **`.markdownlint.yml` — Markdown linting configuration**
- **`.yamllint.yml` — YAML linting configuration**
- **`.shellcheckrc` — Shell lint configuration**
- **`.pre-commit-config.yaml` — Optional pre-commit hooks metadata**
- **`.version` — Current project version**

---

## 🔜 Not Yet Present (Planned in Later Versions)

These directories are **not yet introduced** but are part of the intended long-term structure.  
See [`ROADMAP.md`](../ROADMAP.md) for targeted milestones.

- `test/` — Additional integration/e2e harnesses
- `benchmarks/` — Performance regression tracking
- Contribution templates

---

## 🔄 Philosophy

This structure is:

- ✅ **Minimal by default**
- ✅ **Expanded only when needed**
- ✅ **Documented at every versioned step**

`goboot` aims to remain **predictable**, **clear**, and **scalable**,
without overwhelming new contributors or hiding logic behind automation.

---

_Last updated: v0.3.0 — matches real files in the repository._
