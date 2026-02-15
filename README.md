# goboot

A modular, versioned starting point for production-grade Go projects.

[![License](https://img.shields.io/github/license/it-timo/goboot)](LICENSE)
[![Version](https://img.shields.io/github/v/release/it-timo/goboot?include_prereleases)](https://github.com/it-timo/goboot/releases)

---

## 📦 What is `goboot`?

`goboot` is a deterministic scaffolding tool for Go repositories.
It focuses on explicit, layered project setup that stays maintainable as projects grow.

It is not a framework or IDE. It targets teams that want reproducible output,
clear service boundaries, and auditable generation behavior.

---

## 📁 Current State (v0.1.0)

`v0.1.0` is the CI foundation milestone with provider-aware CI generation
and policy-based image pinning.

### Core Capabilities

- **Modular Service Architecture**: Logic is split into isolated services
(`baseproject`, `baselint`, `basetest`, `baseci`) with strict contracts.
- **Containerized Lint Tooling**: Lint jobs run via Docker by default, while CI simulation uses host tools (`act`, `gitlab-ci-local`).
- **Secure Scaffolding**: Built-in protection against path traversal and strict root confinement.
- **BDD Testing**: Full Ginkgo/Gomega suite covering core packages and E2E flows.
- **CI Generation**: GitLab and GitHub CI templates generated from explicit contracts and image policies.

For file layout details, see [`doc/PROJECT_STRUCTURE.md`](./doc/PROJECT_STRUCTURE.md).

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

- [Go 1.25+](https://go.dev/doc/install)
- [Make](https://www.gnu.org/software/make/) (optional, for `Makefile` tasks)
- [Task](https://taskfile.dev) (optional, for `Taskfile.yml` tasks)
- [GolangCI-Lint](https://golangci-lint.run/) (for Go linting, see `.golangci.yml`)
- [Yamllint](https://yamllint.readthedocs.io/) (for YAML linting, see `.yamllint.yaml`)
- [Checkmake](https://github.com/mrtazz/checkmake) (for Makefile linting)
- [Docker](https://www.docker.com/) (for running Markdown linting via container)
- [act](https://github.com/nektos/act) (for local GitHub Actions simulation)
- [gitlab-ci-local](https://github.com/firecow/gitlab-ci-local) (for local GitLab CI simulation)
- [Markdownlint](https://github.com/DavidAnson/markdownlint) (for Markdown linting, see `.markdownlint.yaml`)
- [EditorConfig Checker](https://github.com/editorconfig-checker/editorconfig-checker) (for `.editorconfig` validation)

For `make verify_ci_canary`, Docker daemon access
(including `/var/run/docker.sock`) and a non-restricted host runtime are required.
Provider selection for canary checks defaults to `gitProvider` in `configs/goboot.yml`
and can be overridden with `./scripts/verify_ci_canary.sh --provider=<github|gitlab|both>`.

### Clone and Explore

```bash
git clone https://github.com/it-timo/goboot.git
cd goboot
make lint
make test
# or, using Task
task lint
task test
```

`make test` runs BDD suites (Ginkgo/Gomega) with race detection and coverage,
excluding `/test/noauto` and `/templates` by default. See [`doc/TESTING.md`](./doc/TESTING.md).

There’s still no “one-click project generator” here — the goal is deterministic scaffolding with visible layers.

### First Run

If you are new to the repository, start with:

- [`doc/quickstart.md`](./doc/quickstart.md)
- [`doc/examples.md`](./doc/examples.md)

These documents focus on input config -> expected output behavior.

---

## 📚 Project Planning

This repository uses:

- [ROADMAP.md](./ROADMAP.md) for planned milestones
- [doc/VERSIONING.md](./doc/VERSIONING.md) for semantic version handling
- [doc/WORKFLOW.md](./doc/WORKFLOW.md) to define long-term contribution and CI logic
- [doc/ci.md](./doc/ci.md) for CI policy modes and provider-specific generated CI behavior
- [doc/quickstart.md](./doc/quickstart.md) for first-run usage without reading internals
- [doc/examples.md](./doc/examples.md) for concrete config-to-output scenarios
- [doc/PROJECT_STRUCTURE.md](./doc/PROJECT_STRUCTURE.md) to track how the folder layout evolves over time
- [doc/README.md](./doc/README.md) for the full docs index
- [doc/adr/](./doc/adr) for architecture decisions in ADR format
- [doc/img/](./doc/img) for flow visualizations

These documents evolve alongside the project.

---

## 🧱 Template Limits

To protect generation runs from unexpectedly large template trees, goboot enforces runtime guardrails:

- max template files per service source: `5000`
- max total template bytes per service source: `64 MiB`

Limits are validated before rendering and fail fast when exceeded.

---

## ⚖️ License

Licensed under the MIT License. See [LICENSE](./LICENSE).
Includes attribution in [NOTICE](./NOTICE) (if applicable).

---

## 🚧 Status

`goboot` is **pre-alpha** and intended for structural exploration and reproducible setup.
It is not yet suitable as a production generator baseline.

---

## 💖 Support This Project

If `goboot` helps you or saves you time, consider supporting its development:

- [💖 GitHub Sponsors](https://github.com/sponsors/it-timo)
- [🎁 Ko-Fi](https://ko-fi.com/ittimo)
- [☕ BuyMeACoffee](https://buymeacoffee.com/ittimo)

> No pressure — just a small way to say "thanks" if it brought you value.
