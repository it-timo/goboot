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

Current release.

**Focus:** reproducible automation

- CI pipeline for goboot itself
  - lint, test, coverage
- CI scaffolding for generated projects
- Status badges (tests, coverage, lint)

---

## In Progress

### v0.1.1 — Structured Logging

**Focus:** observability without noise

- Replace `fmt` usage with structured logging
- Configurable logging in generated projects
- ADR update for log handling

## Planned Milestones

### v0.2.0 — Containerization

**Focus:** deployment-ready outputs

- Dockerfile templates (multi-stage)
- docker-compose for multi-service setups
- Container-based integration testing
- Local dev parity with CI

---

### v0.2.1 — Release Automation

- GoReleaser integration
- Automated versioning and changelogs
- Binary distribution

---

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
