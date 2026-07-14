# 🛠 Development Workflow — `goboot`

This file describes the current development workflow for `goboot`.
It prioritizes explicit process, clear layering, and repeatable quality checks.

As `goboot` evolves, this workflow will expand (see [`ROADMAP.md`](../ROADMAP.md)).

---

## 📦 Project Standards

| Area          | Practice                                                                                                           |
|---------------|--------------------------------------------------------------------------------------------------------------------|
| Structure     | Use `cmd/`, `pkg/`, `configs/`, `templates/`, and `doc/` directories                                               |
| Versioning    | Semantic versioning (`v0.x.x`) with clearly defined milestones                                                     |
| Licensing     | MIT license in `LICENSE` and `NOTICE`                                                                              |
| Documentation | Markdown-based (`README.md`, `ROADMAP.md`, ADRs)                                                                   |
| Tooling       | `Makefile`, `Taskfile.yml`, and `scripts/` for common tasks                                                        |
| Linting       | Dockerized `golangci-lint`, `yamllint`, `checkmake`, `markdownlint`, `shellcheck`, `shfmt`, `editorconfig-checker` |
| Testing       | BDD-style suites with Ginkgo/Gomega via `make test` (race + coverage)                                              |
| Profiles      | Explicit `minimal`, `standard`, `enterprise`, or `oss` scaffold baseline                                          |
| Philosophy    | No runtime magic, minimal indirection, deterministic scaffolding                                                   |

---

## 🚧 Development Flow

### 🔀 Branching

| Type       | Prefix    | Description                                    |
|------------|-----------|------------------------------------------------|
| Features   | `feat/*`  | New features or enhancements                   |
| Fixes      | `fix/*`   | Bug fixes                                      |
| Chores     | `chore/*` | Maintenance tasks (e.g., updates, refactoring) |
| Planning   | `docs/*`  | Documentation updates                          |

> All work happens on branches. Merge into `main` only when stable and reviewed.

---

## ✅ Commit Conventions

```text
feat: Add support for base project target dir
fix: Correct YAML parsing edge case
docs: Add ADR for service registration
chore: Prepare v0.0.1 release tag
```

---

## 📦 Releases

1. Complete and test the milestone
2. Run `make release_check`
3. Update `.version`, `ROADMAP.md`, and release-facing docs
4. Create and push an annotated semantic-version tag; provider CI runs
   GoReleaser and publishes the binary release:

  ```bash
  git tag v0.1.1 -m "Release v0.1.1 — Adds structured logging and release hardening"
  git push origin v0.1.1
  ```

---

## 📚 Planning Files

This project uses:

- [`ROADMAP.md`](../ROADMAP.md) — Development layers and milestones
- [`PROJECT_STRUCTURE.md`](./PROJECT_STRUCTURE.md) — Directory design rationale
- [`ci.md`](./ci.md) — CI policy modes and generated CI expectations
- [`profiles.md`](./profiles.md) — Template profile behavior and precedence
- [`releasing.md`](./releasing.md) — Tag-driven release automation
- [`adr/`](./adr) — Architecture decisions (ADRs)

---

## ⚠️ Notes

- This is a **developer tool**, not a runtime framework
- All logic must be explicit, safe, and maintainable
- Avoid runtime "magic," reflection, or abstract factories
- Linting is designed to be **containerized**:
  - `make lint` runs all linters via Docker (no local installs needed; includes shellcheck/shfmt/editorconfig-checker)
- Run `make test` (Ginkgo/Gomega suites) before submitting PRs
- **Use all `make`, `task` and `scripts` files before submitting PRs**
