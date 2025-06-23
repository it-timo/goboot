# goboot

> A modular, versioned starting point for building production-grade Go projects.

[![License](https://img.shields.io/github/license/it-timo/goboot)](LICENSE)
[![Version](https://img.shields.io/github/v/release/it-timo/goboot?include_prereleases)](https://github.com/it-timo/goboot/releases)

---

## 📦 What is `goboot`?

> `goboot` is a deterministic scaffolding tool that provides a clean, modular foundation for real-world Go projects
> — including applications, tools, and infrastructure code.
>
> It’s **not** a framework. It’s **not** an IDE.
> Instead, `goboot` is built for OSS maintainers, contributors,
> and enterprise teams who care about long-term project hygiene, structure, and reproducibility.
>
> This repository focuses on **layered, progressive structure scaffolding**, not one-click demos or opinionated code generation.
>
> The goal isn’t just to get you started — it’s to help you grow Go projects that remain clean, consistent,
> and scalable over time.

---

## 📁 Current State (v0.0.1)

This version (`v0.0.1`) includes:

- CLI entry point: `cmd/goboot/main.go`
- Configuration logic: `pkg/config/`
- Core engine: `pkg/goboot/`
- First service example: `pkg/baseProject/`
- Shared types: `pkg/types/`
- Utility helpers: `pkg/utils/`
- Configuration files: `configs/goboot.yml`, `configs/base_project.yml`
- Project templates: `templates/project_base/`
- Architecture decisions: `doc/adr/`
- Visuals: `doc/img/`, `doc/diagram/`
- GitHub funding: `.github/FUNDING.yml`
- Central docs: `README.md`, `ROADMAP.md`, `VERSIONING.md`, `WORKFLOW.md`, `LICENSE`, `NOTICE`
- Dev files: `.editorconfig`, `.gitignore`, `.gitattributes`, `.env.example`, `.env.ci`, `.nvmrc`, `.version`
- **New in v0.0.1:**  
  - `Makefile` for common tasks  
  - `Taskfile.yml` for task runner support  
  - `.golangci.yml` for Go linting configuration  
  - `.markdownlint.yaml` for Markdown linting  
  - `.yamllint.yaml` for YAML linting  
  - `scripts/` directory for developer scripts (e.g., lint)

> The project is in an **early layout and planning phase**. The structure is minimal but already modular and documented.

---

## 📐 Intended Design Principles

Even in early stages, `goboot` is being built with:

- Layered versioning and changelog visibility
- Clear module boundaries (`cmd/`, `pkg/`, `configs/`, etc.)
- Future support for Docker, CI/CD, and template-driven code generation

You can follow the structural milestones in [`ROADMAP.md`](./ROADMAP.md).

---

## 🛠️ Getting Started (For Contributors Only)

### Requirements

- [Go 1.24+](https://go.dev/doc/install)
- [Make](https://www.gnu.org/software/make/) (optional, for `Makefile` tasks)
- [Task](https://taskfile.dev) (optional, for `Taskfile.yml` tasks)
- [GolangCI-Lint](https://golangci-lint.run/) (for Go linting, see `.golangci.yml`)
- [Yamllint](https://yamllint.readthedocs.io/) (for YAML linting, see `.yamllint.yaml`)
- [Checkmake](https://github.com/mrtazz/checkmake) (for Makefile linting)
- [Docker](https://www.docker.com/) (for running Markdown linting via container)
- [Markdownlint](https://github.com/DavidAnson/markdownlint) (for Markdown linting, see `.markdownlint.yaml`)

### Clone and Explore

```bash
git clone https://github.com/it-timo/goboot.git
cd goboot
make lint
# or, using Task
task lint
```

There’s not much output yet — but that’s intentional.

---

## 📚 Project Planning

This repository uses:

- [ROADMAP.md](./ROADMAP.md) for planned milestones
- [VERSIONING.md](./VERSIONING.md) for semantic version handling
- [WORKFLOW.md](./WORKFLOW.md) to define long-term contribution and CI logic
- [PROJECT_STRUCTURE.md](./PROJECT_STRUCTURE.md) to track how the folder layout evolves over time
- [doc/adr/](./doc/adr) for architecture decisions in ADR format
- [doc/img/](./doc/img) for flow visualizations

These documents evolve alongside the project.

---

## ⚖️ License

Licensed under the MIT License. See [LICENSE](./LICENSE).
Includes attribution in [NOTICE](./NOTICE) (if applicable).

---

## 🚧 Status

> `goboot` is in **pre-alpha**.
> Intended for structural exploration and reproducible project setup. Not yet suitable for generating production-ready projects.

---

## 💖 Support This Project

If `goboot` helps you or saves you time, consider supporting its development:

- [💖 GitHub Sponsors](https://github.com/sponsors/it-timo)
- [🎁 Ko-Fi](https://ko-fi.com/ittimo)
- [☕ BuyMeACoffee](https://buymeacoffee.com/ittimo)

> No pressure — just a small way to say "thanks" if it brought you value.
