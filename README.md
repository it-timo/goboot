# goboot

A modular, versioned scaffold for reproducible Go project generation.

[![License](https://img.shields.io/github/license/it-timo/goboot)](LICENSE)
[![Version](https://img.shields.io/github/v/release/it-timo/goboot?include_prereleases)](https://github.com/it-timo/goboot/releases)
[![Test](https://github.com/it-timo/goboot/actions/workflows/test.yml/badge.svg)](https://github.com/it-timo/goboot/actions/workflows/test.yml)
[![Lint](https://github.com/it-timo/goboot/actions/workflows/lint.yml/badge.svg)](https://github.com/it-timo/goboot/actions/workflows/lint.yml)
[![Coverage](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/it-timo/goboot/main/badges/coverage.json)](https://github.com/it-timo/goboot/actions/workflows/test.yml)
[![Security](https://github.com/it-timo/goboot/actions/workflows/security.yml/badge.svg)](https://github.com/it-timo/goboot/actions/workflows/security.yml)

---

## 📦 What is `goboot`?

`goboot` is a deterministic scaffolding tool for Go repositories.
It focuses on explicit, layered project setup that stays maintainable as projects grow.

It is not a framework or IDE. It targets teams that want reproducible output,
clear service boundaries, and auditable generation behavior.

---

## 📁 Current State

`v0.6.0` is the active performance and scale milestone. It adds measured
generation baselines, repeatable profiling, and bounded parallel service execution.

### Core Capabilities

- **Modular Service Architecture**: Logic is split into isolated services
(`base_project`, `base_lint`, `base_test`, `base_logger`, `base_docker`,
`base_release`, `base_governance`, `base_supplychain`, `base_local`, `base_ci`)
with strict contracts.
- **Containerized Lint Tooling**: Lint jobs run via Docker by default, while CI simulation uses host tools (`act`, `gitlab-ci-local`).
- **Secure Scaffolding**: Built-in protection against path traversal and strict root confinement.
- **BDD Testing**: Full Ginkgo/Gomega suite covering core packages and E2E flows.
- **CI Generation**: GitLab and GitHub CI templates generated from explicit contracts and image policies.
- **Logger-Aware Scaffolding**: Generated projects can use `slog` or `zerolog`
through explicit settings owned by the project templates.
- **CLI Container Packaging**: Generated projects can include a multi-stage
Dockerfile, compose file, `.dockerignore`, local Docker commands, and CI image
build validation through `base_docker`.
- **Containerized Generator**: The `goboot` CLI itself can be built as a Docker
image for mounted-workspace generation runs.
- **Tag-Driven Releases**: GoReleaser produces Linux, macOS, and Windows binary
  archives and checksums from explicit semantic-version tags.
- **Explicit Profiles**: Named project baselines adjust owned lint/test defaults
  without silently enabling services or overriding explicit service values.
- **Repository Governance**: Generated GitHub and GitLab projects receive
  profile-aware issue/change templates, CODEOWNERS, contribution guidance, and
  a private-first security policy.
- **Supply-Chain Security**: Generated CI enforces pinned vulnerability and
  license checks, emits CycloneDX SBOM artifacts, and adds CodeQL on GitHub.
- **Measured Scale**: Template, config, and service benchmarks track generator
  costs, while opt-in bounded parallelism accelerates independent services.
For file layout details, see [`doc/PROJECT_STRUCTURE.md`](./doc/PROJECT_STRUCTURE.md).

---

## 📐 Intended Design Principles

Even in early stages, `goboot` is being built with:

- Layered versioning and changelog visibility
- Clear module boundaries (`cmd/`, `pkg/`, `configs/`, etc.)
- Incremental support for Docker, CI/CD, and template-driven code generation

You can follow the structural milestones in [`ROADMAP.md`](./ROADMAP.md).

---

## 🛠️ Getting Started (For Contributors Only)

### Requirements

- [Go 1.26.5+](https://go.dev/doc/install)
- [Make](https://www.gnu.org/software/make/) for `make` targets
- [Task](https://taskfile.dev) for `task` targets and full `verify_intro`
- [Docker](https://www.docker.com/) for containerized lint tooling
  and container image builds

Lint tools such as GolangCI-Lint, Yamllint, Checkmake, Markdownlint, ShellCheck,
shfmt, and EditorConfig Checker run from pinned container images by default.

For `make verify_ci_canary`, [act](https://github.com/nektos/act),
[gitlab-ci-local](https://github.com/firecow/gitlab-ci-local), Docker daemon access
(including `/var/run/docker.sock`), and a non-restricted host runtime are required.
The Make/Task canary targets run both providers. Direct script runs can override
provider selection with `./scripts/verify_ci_canary.sh --provider=<github|gitlab|both|config>`.

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
- [CHANGELOG.md](./CHANGELOG.md) for release-facing change summaries
- [doc/VERSIONING.md](./doc/VERSIONING.md) for semantic version handling
- [doc/WORKFLOW.md](./doc/WORKFLOW.md) to define long-term contribution and CI logic
- [doc/ci.md](./doc/ci.md) for CI policy modes and provider-specific generated CI behavior
- [doc/containerization.md](./doc/containerization.md) for generated Docker behavior
- [doc/releasing.md](./doc/releasing.md) for tag-driven release behavior
- [doc/profiles.md](./doc/profiles.md) for template profile behavior
- [doc/governance.md](./doc/governance.md) for generated contribution and security policy
- [doc/supply-chain-security.md](./doc/supply-chain-security.md) for generated
  source, dependency, license, and SBOM controls
- [doc/performance.md](./doc/performance.md) for parallel execution, benchmarks,
  and profiling
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
