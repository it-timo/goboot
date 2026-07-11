# goboot — Project Roadmap

Deterministic scaffolding for long-lived Go repositories.

This roadmap describes goboot as a service-oriented CLI generator with explicit,
ADR-backed contracts per capability.

Priority is reproducibility and long-term maintainability over one-click convenience.

---

## Design Principles (Invariant)

- Deterministic generation (same input → same output)
- Explicit service composition via YAML
- No reflection, no hidden globals, no magic defaults
- Dockerized tooling for reproducible quality gates
- Secure filesystem model (root confinement + path validation)
- Documentation as architecture (ADR-driven design)

Generated projects are expected to be **lint-clean and test-passing on the first commit**.

---

## Released Milestones

### v0.0.0 — Bootstrap (Released)

- CLI entrypoint and flag handling
- YAML config loader
- Initial project scaffolding
- ADR framework introduced

---

### v0.0.1 — Tooling Baseline (Released)

- Dockerized linting stack (golangci-lint, yamllint, markdownlint, shellcheck, shfmt, checkmake)
- Makefile and Taskfile generation
- Pre-commit hooks (container-based)
- Base services: `baseproject`, `baselint`, `baselocal`
- Template system with strict `.tmpl` policy

---

### v0.0.2 — Testing & Validation (Released)

- `basetest` service with automatic test scaffolding
- stdlib or Ginkgo/Gomega BDD support
- ~90% coverage with race detection
- Conditional linter configuration (BDD-aware)
- Hardened filesystem model
- Internal utils/types refactor
- End-to-end validation: generated projects pass all quality gates
- 33 ADRs documenting architecture

---

### v0.1.0 — CI/CD Foundation (Released)

**Focus:** reproducible automation

- CI pipeline for goboot itself
  - lint, test, coverage
- CI scaffolding for generated projects
- Status badges (tests, coverage, lint)

---

### v0.1.1 — Structured Logging & Release Hardening (Released)

Current release.

**Focus:** observability without noise

- Replace `fmt` usage with structured logging
- Configurable logging in generated projects
- ADR update for log handling
- Refreshed Go, tool, action, and Docker image versions
- Stronger config validation and generated-output guardrails
- Release checks aligned across Make, Task, pre-commit, and CI canaries

---

### v0.2.0 — Containerization

**Focus:** deployment-ready outputs

- Dockerfile templates (multi-stage)
- docker-compose for CLI-container execution
- Container build and compose validation
- Local dev parity with CI

---

## In Progress

### v0.2.1 — Release Automation

- GoReleaser v2 integration for goboot and generated projects
- Explicit semantic versioning through immutable `v*` Git tags
- Automated changelogs, archives, and checksum manifests
- GitHub and GitLab release jobs generated through `base_ci`
- Linux, macOS, and Windows binary distribution for AMD64 and ARM64

---

## Planned Milestones

### v0.3.0 — Template Profiles

**Focus:** controlled flexibility

- Profiles: minimal / standard / enterprise / OSS
- Profile-specific lint/test baselines
- Profile-aware documentation

---

### v0.4.0 — Governance & Contribution

- Issue / PR templates
- CODEOWNERS
- SECURITY.md
- Contribution workflows

---

### v0.5.0 — Supply Chain Security

- CodeQL
- License compliance scanning
- Dependency vulnerability checks
- SBOM generation

---

### v0.6.0 — Performance & Scale

- Template rendering benchmarks
- Config parsing performance
- Large-project generation profiling
- Parallel execution optimizations

---

## 1.0 Vision

### v1.0.0 — Stable Public Release

- Stable template registry
- Hardened documentation
- Public announcement

**Exit criteria:**

- ≥90% test coverage
- No breaking changes for 6 months
- Real-world usage across multiple projects
- Performance baselines established

---

## Beyond 1.0

- External service plugins
- Community templates
- `goboot doctor` (project health checks)
- Optional TUI
- Declarative init pipelines
