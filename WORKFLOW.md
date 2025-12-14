# 🛠 Development Workflow — `goboot`

This file describes the early development approach for `goboot`.  
It prioritizes simplicity, clarity, and clean layering during the initial stages.

> As `goboot` evolves, this workflow will expand — see [`ROADMAP.md`](./ROADMAP.md).

---

## 📦 Project Standards

| Area          | Practice                                                                                   |
|---------------|--------------------------------------------------------------------------------------------|
| Structure     | Use `cmd/`, `pkg/`, `configs/`, `templates/`, and `doc/` directories                       |
| Versioning    | Semantic versioning (`v0.x.x`) with clearly defined milestones                             |
| Licensing     | MIT license in `LICENSE` and `NOTICE`                                                      |
| Documentation | Markdown-based (`README.md`, `ROADMAP.md`, ADRs)                                           |
| Tooling       | `Makefile`, `Taskfile.yml`, and `scripts/` for common tasks                                |
| Linting       | Dockerized `golangci-lint`, `yamllint`, `checkmake`, `markdownlint`, `shellcheck`, `shfmt` |
| Testing       | BDD-style suites with Ginkgo/Gomega via `make test` (race + coverage)                      |
| Philosophy    | No runtime magic, minimal indirection, deterministic scaffolding                           |

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

## 📦 Releases (Manual for Now)

1. Complete and test the milestone
2. Update `ROADMAP.md` and optionally `CHANGELOG.md`
3. Create annotated tag:

   ```bash
   git tag v0.0.2 -m "Release v0.0.2 — Adds testing scaffold and coverage"
   git push origin v0.0.2
   ```

---

## 📚 Planning Files

This project uses:

- [`ROADMAP.md`](./ROADMAP.md) — Development layers and milestones
- [`PROJECT_STRUCTURE.md`](./PROJECT_STRUCTURE.md) — Directory design rationale
- [`doc/adr/`](./doc/adr) — Architecture decisions (ADRs)

---

## ⚠️ Notes

- This is a **developer tool**, not a runtime framework
- All logic must be explicit, safe, and maintainable
- Avoid runtime "magic," reflection, or abstract factories
- Linting is designed to be **containerized**:
  - `make lint` runs all linters via Docker (no local installs needed; includes shellcheck/shfmt)
- Run `make test` (Ginkgo/Gomega suites) before submitting PRs
- **Use all `make`, `task` and `scripts` files before submitting PRs**
